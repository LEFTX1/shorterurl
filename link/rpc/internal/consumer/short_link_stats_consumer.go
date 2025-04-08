package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"shorterurl/link/rpc/internal/model"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// Redis Stream 相关常量
const (
	// Redis Stream Key
	ShortLinkStatsStreamKey = "short-link:stats:stream"
	// 消费者组名称
	ShortLinkStatsGroupName = "short-link-stats-group"
	// 消费者名称前缀
	ShortLinkStatsConsumerPrefix = "consumer-"
	// 批量获取消息数量
	BatchCount = 10
	// 轮询间隔(毫秒)
	PollInterval = 100
)

// StatsRecord 统计记录数据结构
type StatsRecord struct {
	FullShortUrl string    `json:"full_short_url"`
	Gid          string    `json:"gid"`
	User         string    `json:"user"`
	UvFirstFlag  bool      `json:"uv_first_flag"`
	UipFirstFlag bool      `json:"uip_first_flag"`
	Ip           string    `json:"ip"`
	Browser      string    `json:"browser"`
	Os           string    `json:"os"`
	Device       string    `json:"device"`
	Network      string    `json:"network"`
	Locale       string    `json:"locale"`
	CurrentDate  time.Time `json:"current_date"`
}

// ShortLinkStatsConsumer 短链接统计消费者
type ShortLinkStatsConsumer struct {
	serviceCtx  ServiceContext
	redisStream *RedisStream
	consumerID  string
	running     bool
	stopChan    chan struct{}
	wg          sync.WaitGroup
	logger      logx.Logger
}

// NewShortLinkStatsConsumer 创建新的短链接统计消费者
func NewShortLinkStatsConsumer(serviceCtx ServiceContext) *ShortLinkStatsConsumer {
	// 生成唯一的消费者ID
	consumerID := fmt.Sprintf("%s%d", ShortLinkStatsConsumerPrefix, time.Now().UnixNano())
	logx.Infof("[统计消费者] 创建消费者: ID=%s", consumerID)

	// 创建 Redis Stream 包装器
	redisStream := NewRedisStream(serviceCtx.GetRedis())

	return &ShortLinkStatsConsumer{
		serviceCtx:  serviceCtx,
		redisStream: redisStream,
		consumerID:  consumerID,
		stopChan:    make(chan struct{}),
		wg:          sync.WaitGroup{},
		logger:      logx.WithContext(context.Background()),
	}
}

// Start 启动消费者
func (c *ShortLinkStatsConsumer) Start() {
	logx.Infof("[统计消费者] 启动消费者: ID=%s", c.consumerID)

	// 确保消费者组存在
	if err := c.ensureConsumerGroup(); err != nil {
		c.logger.Errorf("创建消费者组失败: %v", err)
		return
	}
	logx.Infof("[统计消费者] 消费者组已就绪: %s", ShortLinkStatsGroupName)

	// 启动消费协程
	c.wg.Add(1)
	c.running = true
	go c.consume()

	logx.Infof("[统计消费者] 消费者已启动并开始监听消息")
}

// Stop 停止消费者处理循环
func (c *ShortLinkStatsConsumer) Stop() {
	if !c.running {
		return
	}
	c.running = false
	close(c.stopChan)
	c.wg.Wait()
	logx.Infof("短链接统计消费者已停止: %s", c.consumerID)
}

// Submit 提交统计记录到Redis Stream
func (c *ShortLinkStatsConsumer) Submit(record *StatsRecord) {
	// 如果当前日期为空，设置为当前时间
	if record.CurrentDate.IsZero() {
		record.CurrentDate = time.Now()
	}

	logx.Infof("[统计] 接收统计记录: 短链接=%s, GID=%s, UV首次=%v, UIP首次=%v, 设备=%s, 浏览器=%s, 系统=%s",
		record.FullShortUrl, record.Gid, record.UvFirstFlag, record.UipFirstFlag,
		record.Device, record.Browser, record.Os)

	// 将记录直接放到字段映射中
	values := map[string]string{
		"full_short_url": record.FullShortUrl,
		"gid":            record.Gid,
		"user":           record.User,
		"uv_first_flag":  fmt.Sprintf("%t", record.UvFirstFlag),
		"uip_first_flag": fmt.Sprintf("%t", record.UipFirstFlag),
		"ip":             record.Ip,
		"browser":        record.Browser,
		"os":             record.Os,
		"device":         record.Device,
		"network":        record.Network,
		"locale":         record.Locale,
		"current_date":   record.CurrentDate.Format(time.RFC3339),
	}

	// 发送到Redis Stream
	msgID, err := c.redisStream.Xadd(ShortLinkStatsStreamKey, "*", values)
	if err != nil {
		logx.Errorf("发送统计记录到Redis Stream失败: %v", err)
		return
	}

	logx.Infof("[统计] 成功发送到Redis Stream, 消息ID: %s", msgID)
}

// ensureConsumerGroup 确保消费者组存在
func (c *ShortLinkStatsConsumer) ensureConsumerGroup() error {
	// 直接创建消费者组，如果 Stream 不存在会自动创建
	_, err := c.redisStream.Xgroup(ShortLinkStatsStreamKey, ShortLinkStatsGroupName, "0", true)
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("创建消费者组失败: %v", err)
	}

	return nil
}

// consume 消费循环
func (c *ShortLinkStatsConsumer) consume() {
	defer c.wg.Done()
	logx.Infof("[统计消费者] 开始消费循环: ID=%s", c.consumerID)

	// 用于保存待确认的消息ID
	pending := make([]string, 0, BatchCount)
	retryCount := 0
	maxRetries := 3

	for {
		select {
		case <-c.stopChan:
			// 处理剩余的消息
			logx.Infof("[统计消费者] 收到停止信号，处理剩余消息")
			if len(pending) > 0 {
				c.processPendingMessages(pending)
			}
			return

		default:
			// 读取消息
			logx.Infof("[统计消费者] 尝试从Redis Stream读取消息")
			messages, err := c.redisStream.Xreadgroup(ShortLinkStatsStreamKey, ShortLinkStatsGroupName, c.consumerID, ">", BatchCount, 1000)
			if err != nil {
				if strings.Contains(err.Error(), "circuit breaker is open") {
					// 熔断器打开，等待一段时间后重试
					if retryCount < maxRetries {
						retryCount++
						logx.Infof("[统计消费者] 熔断器打开，等待后重试 (%d/%d)", retryCount, maxRetries)
						time.Sleep(time.Second * time.Duration(retryCount))
						continue
					}
				}
				c.logger.Errorf("读取待处理消息失败: %v", err)
				time.Sleep(time.Second)
				continue
			}

			// 重置重试计数
			retryCount = 0

			// 处理消息
			if len(messages) > 0 {
				logx.Infof("[统计消费者] 读取到 %d 条待处理消息", len(messages))
			}

			for _, msg := range messages {
				// 解析消息
				statsRecord, err := c.parseStreamMessage(msg)
				if err != nil {
					c.logger.Errorf("解析消息失败: %v", err)
					// 即使解析失败，也添加到待确认列表，避免死信
					pending = append(pending, msg.ID)
					continue
				}

				// 处理统计记录
				if err := c.processMessage(statsRecord); err != nil {
					c.logger.Errorf("处理统计记录失败: %v", err)
					// 处理失败的记录也确认，避免死信
					pending = append(pending, msg.ID)
					continue
				}

				// 处理成功，添加到待确认列表
				pending = append(pending, msg.ID)

				// 如果待确认消息达到批量大小，进行确认
				if len(pending) >= BatchCount {
					logx.Infof("[统计消费者] 确认 %d 条消息", len(pending))
					if err := c.ackMessages(pending); err != nil {
						c.logger.Errorf("确认消息失败: %v", err)
					} else {
						logx.Infof("[统计消费者] 成功确认 %d 条消息", len(pending))
					}
					pending = pending[:0]
				}
			}

			// 如果有待确认消息，也进行确认，避免消息堆积
			if len(pending) > 0 {
				logx.Infof("[统计消费者] 确认剩余 %d 条消息", len(pending))
				if err := c.ackMessages(pending); err != nil {
					c.logger.Errorf("确认剩余消息失败: %v", err)
				} else {
					logx.Infof("[统计消费者] 成功确认剩余 %d 条消息", len(pending))
				}
				pending = pending[:0]
			}

			// 如果没有消息，等待一段时间
			if len(messages) == 0 {
				time.Sleep(time.Second)
			}
		}
	}
}

// processPendingMessages 处理待确认的消息
func (c *ShortLinkStatsConsumer) processPendingMessages(pending []string) {
	if len(pending) == 0 {
		return
	}

	// 尝试确认消息
	if err := c.ackMessages(pending); err != nil {
		c.logger.Errorf("确认剩余消息失败: %v", err)
	}
}

// ackMessages 确认消息
func (c *ShortLinkStatsConsumer) ackMessages(ids []string) error {
	_, err := c.redisStream.Xack(ShortLinkStatsStreamKey, ShortLinkStatsGroupName, ids...)
	return err
}

// parseStreamMessage 解析 Stream 消息为统计记录
func (c *ShortLinkStatsConsumer) parseStreamMessage(msg StreamMessage) (*StatsRecord, error) {
	logx.Infof("[统计] 开始解析消息: ID=%s", msg.ID)

	record := &StatsRecord{}
	var err error

	// 解析字段
	if fullShortUrl, ok := msg.Fields["full_short_url"]; ok {
		record.FullShortUrl = fullShortUrl
	}
	if gid, ok := msg.Fields["gid"]; ok {
		record.Gid = gid
	}
	if user, ok := msg.Fields["user"]; ok {
		record.User = user
	}
	if uvFirstFlag, ok := msg.Fields["uv_first_flag"]; ok {
		record.UvFirstFlag = uvFirstFlag == "true"
	}
	if uipFirstFlag, ok := msg.Fields["uip_first_flag"]; ok {
		record.UipFirstFlag = uipFirstFlag == "true"
	}
	if ip, ok := msg.Fields["ip"]; ok {
		record.Ip = ip
	}
	if browser, ok := msg.Fields["browser"]; ok {
		record.Browser = browser
	}
	if os, ok := msg.Fields["os"]; ok {
		record.Os = os
	}
	if device, ok := msg.Fields["device"]; ok {
		record.Device = device
	}
	if network, ok := msg.Fields["network"]; ok {
		record.Network = network
	}
	if locale, ok := msg.Fields["locale"]; ok {
		record.Locale = locale
	}
	if currentDate, ok := msg.Fields["current_date"]; ok {
		record.CurrentDate, err = time.Parse(time.RFC3339, currentDate)
		if err != nil {
			return nil, fmt.Errorf("解析时间失败: %v", err)
		}
	}

	logx.Infof("[统计] 解析成功: 短链接=%s, GID=%s, UV首次=%v, UIP首次=%v",
		record.FullShortUrl, record.Gid, record.UvFirstFlag, record.UipFirstFlag)

	return record, nil
}

// processMessage 处理统计记录
func (c *ShortLinkStatsConsumer) processMessage(record *StatsRecord) error {
	logx.Infof("[统计] 开始处理统计记录: 短链接=%s, GID=%s", record.FullShortUrl, record.Gid)

	ctx := context.Background()
	commonDB := c.serviceCtx.GetDBs().GetCommon()
	linkDB := c.serviceCtx.GetDBs().GetLinkDB()

	// 更新今日统计
	if err := c.updateTodayStats(ctx, commonDB, record, record.Gid); err != nil {
		return fmt.Errorf("更新今日统计失败: %v", err)
	}
	logx.Infof("[统计] 更新今日统计成功")

	// 更新基础统计
	if err := c.updateBaseStats(ctx, linkDB, record, record.Gid); err != nil {
		return fmt.Errorf("更新基础统计失败: %v", err)
	}
	logx.Infof("[统计] 更新基础统计成功: PV增加, UV增加=%v, UIP增加=%v",
		record.UvFirstFlag, record.UipFirstFlag)

	// 更新地区统计
	if err := c.updateLocaleStats(ctx, commonDB, record, record.Gid); err != nil {
		return fmt.Errorf("更新地区统计失败: %v", err)
	}
	if record.Locale != "" {
		logx.Infof("[统计] 更新地区统计成功: 地区=%s", record.Locale)
	}

	// 更新设备统计
	if err := c.updateDeviceStats(ctx, commonDB, record, record.Gid); err != nil {
		return fmt.Errorf("更新设备统计失败: %v", err)
	}
	logx.Infof("[统计] 更新设备统计成功: 设备=%s", record.Device)

	// 更新浏览器统计
	if record.Browser != "" {
		if err := c.updateBrowserStats(ctx, commonDB, record, record.Gid); err != nil {
			return fmt.Errorf("更新浏览器统计失败: %v", err)
		}
		logx.Infof("[统计] 更新浏览器统计成功: 浏览器=%s", record.Browser)
	}

	// 更新操作系统统计
	if record.Os != "" {
		if err := c.updateOsStats(ctx, commonDB, record, record.Gid); err != nil {
			return fmt.Errorf("更新操作系统统计失败: %v", err)
		}
		logx.Infof("[统计] 更新操作系统统计成功: 系统=%s", record.Os)
	}

	// 更新网络统计
	if record.Network != "" {
		if err := c.updateNetworkStats(ctx, commonDB, record, record.Gid); err != nil {
			return fmt.Errorf("更新网络统计失败: %v", err)
		}
		logx.Infof("[统计] 更新网络统计成功: 网络=%s", record.Network)
	}

	// 记录访问日志
	if err := c.insertAccessLog(ctx, commonDB, record); err != nil {
		return fmt.Errorf("记录访问日志失败: %v", err)
	}
	logx.Infof("[统计] 记录访问日志成功")

	logx.Infof("[统计] 统计记录处理完成: 短链接=%s", record.FullShortUrl)
	return nil
}

// updateBaseStats 更新基础统计
func (c *ShortLinkStatsConsumer) updateBaseStats(ctx context.Context, db *gorm.DB, record *StatsRecord, gid string) error {
	// 更新总访问量
	if err := db.WithContext(ctx).Model(&model.Link{}).
		Where("gid = ? AND full_short_url = ?", gid, record.FullShortUrl).
		UpdateColumn("total_pv", gorm.Expr("total_pv + ?", 1)).Error; err != nil {
		return err
	}

	// 更新总UV
	if record.UvFirstFlag {
		if err := db.WithContext(ctx).Model(&model.Link{}).
			Where("gid = ? AND full_short_url = ?", gid, record.FullShortUrl).
			UpdateColumn("total_uv", gorm.Expr("total_uv + ?", 1)).Error; err != nil {
			return err
		}
	}

	// 更新总UIP
	if record.UipFirstFlag {
		if err := db.WithContext(ctx).Model(&model.Link{}).
			Where("gid = ? AND full_short_url = ?", gid, record.FullShortUrl).
			UpdateColumn("total_uip", gorm.Expr("total_uip + ?", 1)).Error; err != nil {
			return err
		}
	}

	return nil
}

// 更新地域统计
func (c *ShortLinkStatsConsumer) updateLocaleStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, gid string) error {
	if record.Locale == "" {
		return nil // 没有地域信息，跳过
	}

	// 从地域信息中提取省市信息
	// 格式通常是 "省份,城市" 或者 "直辖市"
	var province, city, adcode string
	localeInfo := parseLocaleInfo(record.Locale)
	if len(localeInfo) >= 1 {
		province = localeInfo["province"]
		city = localeInfo["city"]
		adcode = localeInfo["adcode"]
	}

	// 使用 INSERT ... ON DUPLICATE KEY UPDATE 语法
	sql := `INSERT INTO t_link_locale_stats 
            (full_short_url, cnt, province, city, adcode, country, create_time, update_time, del_flag) 
            VALUES (?, 1, ?, ?, ?, 'CN', ?, ?, 0)
            ON DUPLICATE KEY UPDATE 
            cnt = cnt + 1,
            update_time = ?`

	err := tx.Exec(sql,
		// INSERT 部分参数
		record.FullShortUrl,
		province,
		city,
		adcode,
		time.Now(),
		time.Now(),
		// UPDATE 部分参数
		time.Now()).Error

	return err
}

// 更新浏览器统计
func (c *ShortLinkStatsConsumer) updateBrowserStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, gid string) error {
	// 使用 INSERT ... ON DUPLICATE KEY UPDATE 语法
	sql := `INSERT INTO t_link_browser_stats 
            (full_short_url, cnt, browser, create_time, update_time, del_flag) 
            VALUES (?, 1, ?, ?, ?, 0)
            ON DUPLICATE KEY UPDATE 
            cnt = cnt + 1,
            update_time = ?`

	err := tx.Exec(sql,
		// INSERT 部分参数
		record.FullShortUrl,
		record.Browser,
		time.Now(),
		time.Now(),
		// UPDATE 部分参数
		time.Now()).Error

	return err
}

// 更新操作系统统计
func (c *ShortLinkStatsConsumer) updateOsStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, gid string) error {
	// 使用 INSERT ... ON DUPLICATE KEY UPDATE 语法
	sql := `INSERT INTO t_link_os_stats 
            (full_short_url, cnt, os, create_time, update_time, del_flag) 
            VALUES (?, 1, ?, ?, ?, 0)
            ON DUPLICATE KEY UPDATE 
            cnt = cnt + 1,
            update_time = ?`

	err := tx.Exec(sql,
		// INSERT 部分参数
		record.FullShortUrl,
		record.Os,
		time.Now(),
		time.Now(),
		// UPDATE 部分参数
		time.Now()).Error

	return err
}

// 更新设备统计
func (c *ShortLinkStatsConsumer) updateDeviceStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, gid string) error {
	// 使用 INSERT ... ON DUPLICATE KEY UPDATE 语法
	sql := `INSERT INTO t_link_device_stats 
            (full_short_url, cnt, device, create_time, update_time, del_flag) 
            VALUES (?, 1, ?, ?, ?, 0)
            ON DUPLICATE KEY UPDATE 
            cnt = cnt + 1,
            update_time = ?`

	err := tx.Exec(sql,
		// INSERT 部分参数
		record.FullShortUrl,
		record.Device,
		time.Now(),
		time.Now(),
		// UPDATE 部分参数
		time.Now()).Error

	return err
}

// 更新网络统计
func (c *ShortLinkStatsConsumer) updateNetworkStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, gid string) error {
	// 使用 INSERT ... ON DUPLICATE KEY UPDATE 语法
	sql := `INSERT INTO t_link_network_stats 
            (full_short_url, cnt, network, create_time, update_time, del_flag) 
            VALUES (?, 1, ?, ?, ?, 0)
            ON DUPLICATE KEY UPDATE 
            cnt = cnt + 1,
            update_time = ?`

	err := tx.Exec(sql,
		// INSERT 部分参数
		record.FullShortUrl,
		record.Network,
		time.Now(),
		time.Now(),
		// UPDATE 部分参数
		time.Now()).Error

	return err
}

// 记录访问日志
func (c *ShortLinkStatsConsumer) insertAccessLog(ctx context.Context, tx *gorm.DB, record *StatsRecord) error {
	// 插入访问日志记录
	sql := `INSERT INTO t_link_access_logs 
            (full_short_url, user, ip, browser, os, network, device, locale, create_time, update_time, del_flag) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)`

	err := tx.Exec(sql,
		record.FullShortUrl,
		record.User,
		record.Ip,
		record.Browser,
		record.Os,
		record.Network,
		record.Device,
		record.Locale,
		time.Now(),
		time.Now()).Error

	return err
}

// 更新短链接主表统计数据
func (c *ShortLinkStatsConsumer) updateLinkStats(ctx context.Context, tx *gorm.DB, record *StatsRecord) error {
	// 更新主表的点击量、PV、UV、UIP
	sql := `UPDATE t_link 
            SET click_num = click_num + 1,
                total_pv = total_pv + 1,
                total_uv = total_uv + ?,
                total_uip = total_uip + ?,
                update_time = ?
            WHERE full_short_url = ? AND gid = ? AND del_flag = 0`

	err := tx.Exec(sql,
		boolToInt(record.UvFirstFlag),
		boolToInt(record.UipFirstFlag),
		time.Now(),
		record.FullShortUrl,
		record.Gid).Error

	return err
}

// 更新今日统计
func (c *ShortLinkStatsConsumer) updateTodayStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, gid string) error {
	// 使用当前日期
	today := time.Now().Format("2006-01-02")
	date, _ := time.Parse("2006-01-02", today)

	// 使用 INSERT ... ON DUPLICATE KEY UPDATE 语法
	sql := `INSERT INTO t_link_stats_today 
            (full_short_url, date, today_pv, today_uv, today_uip, create_time, update_time, del_flag) 
            VALUES (?, ?, 1, ?, ?, ?, ?, 0)
            ON DUPLICATE KEY UPDATE 
            today_pv = today_pv + 1,
            today_uv = today_uv + ?,
            today_uip = today_uip + ?,
            update_time = ?`

	// 替换参数顺序，确保类型匹配
	err := tx.Exec(sql,
		// INSERT 部分的参数
		record.FullShortUrl,
		date,
		boolToInt(record.UvFirstFlag),
		boolToInt(record.UipFirstFlag),
		time.Now(),
		time.Now(),
		// UPDATE 部分的参数
		boolToInt(record.UvFirstFlag),
		boolToInt(record.UipFirstFlag),
		time.Now()).Error

	if err != nil {
		return err
	}

	return nil
}

// 工具函数：布尔值转为整数
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// 工具函数：解析地域信息
func parseLocaleInfo(locale string) map[string]string {
	result := make(map[string]string)

	// 尝试解析 JSON 格式的地域信息
	var localeData map[string]interface{}
	if err := json.Unmarshal([]byte(locale), &localeData); err == nil {
		// 是有效的 JSON 数据
		if province, ok := localeData["province"].(string); ok {
			result["province"] = province
		}
		if city, ok := localeData["city"].(string); ok {
			result["city"] = city
		}
		if adcode, ok := localeData["adcode"].(string); ok {
			result["adcode"] = adcode
		}
		return result
	}

	// 不是 JSON 格式，尝试从字符串中提取信息
	// 假设格式为: "省份,城市"
	parts := strings.Split(locale, ",")

	// 根据分隔的部分数量处理
	if len(parts) >= 2 {
		result["province"] = parts[0]
		result["city"] = parts[1]
		if len(parts) >= 3 {
			result["adcode"] = parts[2]
		}
	} else if len(parts) == 1 {
		// 可能只有省份或直辖市
		result["province"] = parts[0]
		result["city"] = parts[0]
	} else {
		// 未知格式，使用整个字符串作为省份和城市
		result["province"] = locale
		result["city"] = locale
	}

	return result
}
