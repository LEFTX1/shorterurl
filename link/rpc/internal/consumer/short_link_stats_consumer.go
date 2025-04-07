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
	// 确保消费者组存在
	if err := c.ensureConsumerGroup(); err != nil {
		c.logger.Errorf("创建消费者组失败: %v", err)
		return
	}

	// 启动消费协程
	c.wg.Add(1)
	go c.consume()
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
	// 将记录序列化为JSON
	recordJSON, err := json.Marshal(record)
	if err != nil {
		logx.Errorf("序列化统计记录失败: %v", err)
		return
	}

	// 发送到Redis Stream
	values := map[string]string{
		"data": string(recordJSON),
	}
	_, err = c.redisStream.Xadd(ShortLinkStatsStreamKey, "*", values)
	if err != nil {
		logx.Errorf("发送统计记录到Redis Stream失败: %v", err)
		return
	}
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

	// 用于保存待确认的消息ID
	pending := make([]string, 0, BatchCount)
	retryCount := 0
	maxRetries := 3

	for {
		select {
		case <-c.stopChan:
			// 处理剩余的消息
			if len(pending) > 0 {
				c.processPendingMessages(pending)
			}
			return

		default:
			// 读取消息
			messages, err := c.redisStream.Xreadgroup(ShortLinkStatsStreamKey, ShortLinkStatsGroupName, c.consumerID, ">", BatchCount, 1000)
			if err != nil {
				if strings.Contains(err.Error(), "circuit breaker is open") {
					// 熔断器打开，等待一段时间后重试
					if retryCount < maxRetries {
						retryCount++
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
			for _, msg := range messages {
				// 解析消息
				statsRecord, err := c.parseStreamMessage(msg)
				if err != nil {
					c.logger.Errorf("解析消息失败: %v", err)
					continue
				}

				// 处理统计记录
				if err := c.processMessage(statsRecord); err != nil {
					c.logger.Errorf("处理统计记录失败: %v", err)
					continue
				}

				// 添加到待确认列表
				pending = append(pending, msg.ID)

				// 如果待确认消息达到批量大小，进行确认
				if len(pending) >= BatchCount {
					if err := c.ackMessages(pending); err != nil {
						c.logger.Errorf("确认消息失败: %v", err)
					}
					pending = pending[:0]
				}
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

	return record, nil
}

// processMessage 处理统计记录
func (c *ShortLinkStatsConsumer) processMessage(record *StatsRecord) error {
	ctx := context.Background()
	commonDB := c.serviceCtx.GetDBs().GetCommon()
	linkDB := c.serviceCtx.GetDBs().GetLinkDB()

	// 更新今日统计
	if err := c.updateTodayStats(ctx, commonDB, record, record.Gid); err != nil {
		return fmt.Errorf("更新今日统计失败: %v", err)
	}

	// 更新基础统计
	if err := c.updateBaseStats(ctx, linkDB, record, record.Gid); err != nil {
		return fmt.Errorf("更新基础统计失败: %v", err)
	}

	// 更新地区统计
	if err := c.updateLocaleStats(ctx, commonDB, record, record.Gid); err != nil {
		return fmt.Errorf("更新地区统计失败: %v", err)
	}

	// 更新设备统计
	if err := c.updateDeviceStats(ctx, commonDB, record, record.Gid); err != nil {
		return fmt.Errorf("更新设备统计失败: %v", err)
	}

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

	// 查询是否已存在当天的地域统计记录
	var count int64
	err := tx.Table("t_link_locale_stats").
		Where("gid = ? AND full_short_url = ? AND province = ? AND city = ? AND del_flag = 0",
			gid, record.FullShortUrl, province, city).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		// 更新已存在的地域统计记录
		sql := `UPDATE t_link_locale_stats 
                SET cnt = cnt + 1,
                    update_time = ?
                WHERE gid = ? AND full_short_url = ? AND province = ? AND city = ? AND del_flag = 0`

		err = tx.Exec(sql,
			time.Now(),
			gid,
			record.FullShortUrl,
			province,
			city).Error
	} else {
		// 创建新的地域统计记录
		sql := `INSERT INTO t_link_locale_stats 
                (gid, full_short_url, cnt, province, city, adcode, country, create_time, update_time, del_flag) 
                VALUES (?, ?, 1, ?, ?, ?, 'CN', ?, ?, 0)`

		err = tx.Exec(sql,
			gid,
			record.FullShortUrl,
			province,
			city,
			adcode,
			time.Now(),
			time.Now()).Error
	}

	return err
}

// 更新浏览器统计
func (c *ShortLinkStatsConsumer) updateBrowserStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, gid string) error {
	// 查询是否已存在当天的浏览器统计记录
	var count int64
	err := tx.Table("t_link_browser_stats").
		Where("gid = ? AND full_short_url = ? AND browser = ? AND del_flag = 0",
			gid, record.FullShortUrl, record.Browser).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		// 更新已存在的浏览器统计记录
		sql := `UPDATE t_link_browser_stats 
                SET cnt = cnt + 1,
                    update_time = ?
                WHERE gid = ? AND full_short_url = ? AND browser = ? AND del_flag = 0`

		err = tx.Exec(sql,
			time.Now(),
			gid,
			record.FullShortUrl,
			record.Browser).Error
	} else {
		// 创建新的浏览器统计记录
		sql := `INSERT INTO t_link_browser_stats 
                (gid, full_short_url, cnt, browser, create_time, update_time, del_flag) 
                VALUES (?, ?, 1, ?, ?, ?, 0)`

		err = tx.Exec(sql,
			gid,
			record.FullShortUrl,
			record.Browser,
			time.Now(),
			time.Now()).Error
	}

	return err
}

// 更新操作系统统计
func (c *ShortLinkStatsConsumer) updateOsStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, dateStr string) error {
	date, _ := time.Parse("2006-01-02", dateStr)

	// 查询是否已存在当天的操作系统统计记录
	var count int64
	err := tx.Table("t_link_os_stats").
		Where("full_short_url = ? AND date = ? AND os = ? AND del_flag = 0",
			record.FullShortUrl, date, record.Os).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		// 更新已存在的操作系统统计记录
		sql := `UPDATE t_link_os_stats 
                SET cnt = cnt + 1,
                    update_time = ?
                WHERE full_short_url = ? AND date = ? AND os = ? AND del_flag = 0`

		err = tx.Exec(sql,
			time.Now(),
			record.FullShortUrl,
			date,
			record.Os).Error
	} else {
		// 创建新的操作系统统计记录
		sql := `INSERT INTO t_link_os_stats 
                (full_short_url, date, cnt, os, create_time, update_time, del_flag) 
                VALUES (?, ?, 1, ?, ?, ?, 0)`

		err = tx.Exec(sql,
			record.FullShortUrl,
			date,
			record.Os,
			time.Now(),
			time.Now()).Error
	}

	return err
}

// 更新设备统计
func (c *ShortLinkStatsConsumer) updateDeviceStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, dateStr string) error {
	date, _ := time.Parse("2006-01-02", dateStr)

	// 查询是否已存在当天的设备统计记录
	var count int64
	err := tx.Table("t_link_device_stats").
		Where("full_short_url = ? AND date = ? AND device = ? AND del_flag = 0",
			record.FullShortUrl, date, record.Device).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		// 更新已存在的设备统计记录
		sql := `UPDATE t_link_device_stats 
                SET cnt = cnt + 1,
                    update_time = ?
                WHERE full_short_url = ? AND date = ? AND device = ? AND del_flag = 0`

		err = tx.Exec(sql,
			time.Now(),
			record.FullShortUrl,
			date,
			record.Device).Error
	} else {
		// 创建新的设备统计记录
		sql := `INSERT INTO t_link_device_stats 
                (full_short_url, date, cnt, device, create_time, update_time, del_flag) 
                VALUES (?, ?, 1, ?, ?, ?, 0)`

		err = tx.Exec(sql,
			record.FullShortUrl,
			date,
			record.Device,
			time.Now(),
			time.Now()).Error
	}

	return err
}

// 更新网络统计
func (c *ShortLinkStatsConsumer) updateNetworkStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, dateStr string) error {
	date, _ := time.Parse("2006-01-02", dateStr)

	// 查询是否已存在当天的网络统计记录
	var count int64
	err := tx.Table("t_link_network_stats").
		Where("full_short_url = ? AND date = ? AND network = ? AND del_flag = 0",
			record.FullShortUrl, date, record.Network).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		// 更新已存在的网络统计记录
		sql := `UPDATE t_link_network_stats 
                SET cnt = cnt + 1,
                    update_time = ?
                WHERE full_short_url = ? AND date = ? AND network = ? AND del_flag = 0`

		err = tx.Exec(sql,
			time.Now(),
			record.FullShortUrl,
			date,
			record.Network).Error
	} else {
		// 创建新的网络统计记录
		sql := `INSERT INTO t_link_network_stats 
                (full_short_url, date, cnt, network, create_time, update_time, del_flag) 
                VALUES (?, ?, 1, ?, ?, ?, 0)`

		err = tx.Exec(sql,
			record.FullShortUrl,
			date,
			record.Network,
			time.Now(),
			time.Now()).Error
	}

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
func (c *ShortLinkStatsConsumer) updateTodayStats(ctx context.Context, tx *gorm.DB, record *StatsRecord, dateStr string) error {
	date, _ := time.Parse("2006-01-02", dateStr)

	// 检查当前日期是否是记录的日期
	today := time.Now().Format("2006-01-02")
	if dateStr != today {
		return nil // 不是今天的记录，跳过
	}

	// 查询是否已存在今日统计记录
	var count int64
	err := tx.Table("t_link_stats_today").
		Where("full_short_url = ? AND date = ? AND del_flag = 0",
			record.FullShortUrl, date).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		// 更新已存在的今日统计记录
		sql := `UPDATE t_link_stats_today 
                SET today_pv = today_pv + 1,
                    today_uv = today_uv + ?,
                    today_uip = today_uip + ?,
                    update_time = ?
                WHERE full_short_url = ? AND date = ? AND del_flag = 0`

		err = tx.Exec(sql,
			boolToInt(record.UvFirstFlag),
			boolToInt(record.UipFirstFlag),
			time.Now(),
			record.FullShortUrl,
			date).Error
	} else {
		// 创建新的今日统计记录
		sql := `INSERT INTO t_link_stats_today 
                (full_short_url, date, today_pv, today_uv, today_uip, create_time, update_time, del_flag) 
                VALUES (?, ?, 1, ?, ?, ?, ?, 0)`

		err = tx.Exec(sql,
			record.FullShortUrl,
			date,
			boolToInt(record.UvFirstFlag),
			boolToInt(record.UipFirstFlag),
			time.Now(),
			time.Now()).Error
	}

	return err
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
