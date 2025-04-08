package logic

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"shorterurl/link/rpc/internal/consumer"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"crypto/md5"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Redis键模板常量
const (
	// 短链接跳转前缀Key
	ShortLinkGotoKey = "short-link:goto:%s"
	// 短链接空值跳转前缀Key
	ShortLinkIsNullGotoKey = "short-link:is-null:goto_%s"
	// 短链接跳转锁前缀Key
	ShortLinkLockGotoKey = "short-link:lock:goto:%s"
)

type RestoreUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRestoreUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreUrlLogic {
	return &RestoreUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 短链接跳转
func (l *RestoreUrlLogic) RestoreUrl(in *pb.RestoreUrlRequest) (*pb.RestoreUrlResponse, error) {
	// 参数校验
	if in.ShortUri == "" {
		return nil, status.Error(codes.InvalidArgument, "短链接不能为空")
	}

	l.Logger.Infof("开始处理短链接跳转请求: %s", in.ShortUri)

	// 构建完整短链接 - 从配置中获取域名
	domain := l.svcCtx.Config.DefaultDomain
	fullShortUrl := fmt.Sprintf("%s/%s", domain, in.ShortUri)

	// 1. 首先尝试从Redis缓存中获取原始链接
	cacheKey := fmt.Sprintf(ShortLinkGotoKey, fullShortUrl)
	originUrl, err := l.svcCtx.BizRedis.Get(cacheKey)
	if err == nil && originUrl != "" {
		// 找到缓存的原始链接，进行访问统计并返回
		l.asyncRecordStats(fullShortUrl, in.ShortUri)
		return &pb.RestoreUrlResponse{
			OriginUrl: originUrl,
		}, nil
	}

	// 2. 检查空值缓存，避免无效短链接的重复查询
	nullCacheKey := fmt.Sprintf(ShortLinkIsNullGotoKey, fullShortUrl)
	isNull, _ := l.svcCtx.BizRedis.Get(nullCacheKey)
	if isNull != "" {
		return nil, status.Error(codes.NotFound, "短链接不存在")
	}

	// 3. 从数据库查询短链接映射
	// TODO: 这里可以添加布隆过滤器来减少数据库查询

	// 4. 先查询跳转表
	linkGoto, err := l.svcCtx.RepoManager.LinkGoto.FindByFullShortUrl(l.ctx, fullShortUrl)
	if err != nil {
		// 短链接不存在，设置空值缓存，避免缓存穿透
		l.svcCtx.BizRedis.Setex(nullCacheKey, "-", 30*60) // 缓存30分钟
		return nil, status.Error(codes.NotFound, "未找到对应的短链接")
	}

	// 5. 根据跳转表中的分组ID查询具体短链接信息
	link, err := l.svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(l.ctx, fullShortUrl, linkGoto.Gid)
	if err != nil {
		l.svcCtx.BizRedis.Setex(nullCacheKey, "-", 30*60) // 缓存30分钟
		return nil, status.Error(codes.NotFound, "未找到对应的短链接详情")
	}

	// 6. 检查链接是否有效
	if link.DelFlag > 0 {
		l.Logger.Errorf("链接已被删除")
		l.svcCtx.BizRedis.Setex(nullCacheKey, "-", 30*60) // 缓存30分钟
		return nil, status.Error(codes.NotFound, "短链接已被删除")
	}

	if link.EnableStatus > 0 {
		l.Logger.Errorf("链接已被禁用")
		l.svcCtx.BizRedis.Setex(nullCacheKey, "-", 30*60) // 缓存30分钟
		return nil, status.Error(codes.PermissionDenied, "短链接已被禁用")
	}

	// 7. 检查链接是否过期
	now := time.Now()
	if !link.ValidDate.IsZero() && link.ValidDate.Before(now) {
		l.Logger.Errorf("链接已过期")
		l.svcCtx.BizRedis.Setex(nullCacheKey, "-", 30*60) // 缓存30分钟
		return nil, status.Error(codes.PermissionDenied, "短链接已过期")
	}

	// 8. 将有效链接缓存到Redis，计算缓存过期时间
	var cacheExpireSeconds int
	if link.ValidDate.IsZero() {
		// 无过期时间，默认缓存一天
		cacheExpireSeconds = 24 * 60 * 60
	} else {
		// 计算链接的剩余有效时间（秒）
		cacheExpireSeconds = int(link.ValidDate.Sub(now).Seconds())
		if cacheExpireSeconds <= 0 {
			cacheExpireSeconds = 1 // 至少缓存1秒
		}
	}
	l.svcCtx.BizRedis.Setex(cacheKey, link.OriginUrl, cacheExpireSeconds)

	// 9. 记录访问统计
	l.asyncRecordStats(fullShortUrl, in.ShortUri)

	// 10. 返回原始链接
	return &pb.RestoreUrlResponse{
		OriginUrl: link.OriginUrl,
	}, nil
}

// 异步记录访问统计
func (l *RestoreUrlLogic) asyncRecordStats(fullShortUrl, shortUri string) {
	threading.GoSafe(func() {
		// 创建新的上下文
		ctx := context.Background()

		// 查询Gid
		linkGoto, err := l.svcCtx.RepoManager.LinkGoto.FindByFullShortUrl(ctx, fullShortUrl)
		if err != nil {
			logx.Errorf("查询短链接Gid失败: %v", err)
			return
		}

		// 从上下文中获取请求信息
		ip := l.getValueFromContext(l.ctx, "ip", "")
		userAgent := l.getValueFromContext(l.ctx, "user-agent", "")

		// 设备信息默认值
		browser := "未知浏览器"
		os := "未知系统"
		device := "未知设备"
		network := "未知网络"

		// 从 User-Agent 中提取设备信息
		if userAgent != "" {
			parsedBrowser, parsedOS, parsedDevice := parseUserAgent(userAgent)
			if parsedBrowser != "" {
				browser = parsedBrowser
			}
			if parsedOS != "" {
				os = parsedOS
			}
			if parsedDevice != "" {
				device = parsedDevice
			}
			logx.Infof("[访问统计] 从 User-Agent 解析设备信息: %s", userAgent)
		}

		// 获取或生成用户标识（可以是 cookie 中的值或根据 IP+UserAgent 生成的哈希）
		user := l.getUserIdentifier(ip, userAgent)

		// 检查是否是新的 UV 和 UIP
		uvFirstFlag := l.checkFirstUv(fullShortUrl, user)
		uipFirstFlag := l.checkFirstUip(fullShortUrl, ip)

		logx.Infof("[访问统计] 短链接=%s, GID=%s, IP=%s, 用户=%s",
			fullShortUrl, linkGoto.Gid, ip, user)
		logx.Infof("[访问统计] 设备信息: 浏览器=%s, 系统=%s, 设备=%s, 网络=%s",
			browser, os, device, network)
		logx.Infof("[访问统计] 统计标记: 首次UV=%v, 首次UIP=%v",
			uvFirstFlag, uipFirstFlag)

		// 构建统计记录
		statsRecord := &consumer.StatsRecord{
			FullShortUrl: fullShortUrl,
			Gid:          linkGoto.Gid,
			User:         user,
			UvFirstFlag:  uvFirstFlag,
			UipFirstFlag: uipFirstFlag,
			Ip:           ip,
			Browser:      browser,
			Os:           os,
			Device:       device,
			Network:      network,
			CurrentDate:  time.Now(),
		}

		// 如果请求中有IP信息，获取IP地理位置
		if ip != "" {
			logx.Infof("[访问统计] 开始获取访问IP的地理位置信息: %s", ip)

			// 创建IP位置查询逻辑
			ipLocationLogic := NewGetIpLocationLogic(ctx, l.svcCtx)
			formattedLocation, err := ipLocationLogic.GetFormattedLocation(ip)

			if err == nil && formattedLocation != "" {
				logx.Infof("[访问统计] IP地理位置解析成功: %s -> %s", ip, formattedLocation)
				statsRecord.Locale = formattedLocation
			} else {
				logx.Errorf("获取IP地理位置信息失败: %v", err)
			}
		}

		// 提交统计记录到消费者队列
		logx.Infof("[访问统计] 提交统计记录到消费队列")
		l.svcCtx.StatsConsumer.Submit(statsRecord)
	})
}

// parseUserAgent 从 User-Agent 字符串解析设备信息
func parseUserAgent(userAgent string) (browser, os, device string) {
	// 默认值
	browser = "未知浏览器"
	os = "未知系统"
	device = "未知设备"

	// 如果User-Agent为空，直接返回默认值
	if userAgent == "" {
		return
	}

	// 转换为小写以便于匹配
	ua := strings.ToLower(userAgent)

	// 检测浏览器
	switch {
	case strings.Contains(ua, "chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "safari"):
		browser = "Safari"
	case strings.Contains(ua, "edge"):
		browser = "Edge"
	case strings.Contains(ua, "opera"):
		browser = "Opera"
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		browser = "Internet Explorer"
	}

	// 检测操作系统
	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os x"):
		os = "macOS"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod"):
		os = "iOS"
	}

	// 检测设备类型
	switch {
	case strings.Contains(ua, "mobile"):
		device = "手机"
	case strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad"):
		device = "平板"
	case strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") || strings.Contains(ua, "spider"):
		device = "爬虫"
	default:
		device = "电脑"
	}

	return
}

// 从上下文中获取值
func (l *RestoreUrlLogic) getValueFromContext(ctx context.Context, key, defaultValue string) string {
	// 首先尝试从普通上下文获取
	if value, ok := ctx.Value(key).(string); ok && value != "" {
		return value
	}

	// 尝试从 gRPC 元数据中获取
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// 对于 User-Agent，优先使用原始User-Agent
		if key == "user-agent" {
			// 首先尝试获取原始User-Agent
			values := md.Get("original-user-agent")
			if len(values) > 0 && values[0] != "" {
				l.Logger.Infof("使用原始User-Agent: %s", values[0])
				return values[0]
			}

			// 尝试其他可能的 User-Agent 键名
			for _, k := range []string{"User-Agent", "user_agent", "User_Agent"} {
				values := md.Get(k)
				if len(values) > 0 && values[0] != "" {
					return values[0]
				}
			}
		}

		// 检查元数据中是否存在该键
		values := md.Get(key)
		if len(values) > 0 && values[0] != "" {
			return values[0]
		}
	}

	// 尝试从 HTTP 请求 gateway 中获取
	// 在很多 gRPC-Gateway 实现中，原始 HTTP 请求信息被保存在特定的键中
	if reqInfo, ok := ctx.Value("http-request").(map[string]interface{}); ok {
		if headerInfo, ok := reqInfo["headers"].(map[string]interface{}); ok {
			if val, ok := headerInfo[key].(string); ok && val != "" {
				return val
			}
		}
	}

	// 记录日志，帮助调试
	l.Logger.Infof("无法从上下文中获取 %s，使用默认值: %s", key, defaultValue)

	// 最后，返回默认值
	return defaultValue
}

// 获取用户标识
func (l *RestoreUrlLogic) getUserIdentifier(ip, userAgent string) string {
	// 如果上下文中有用户标识，直接使用
	if user, ok := l.ctx.Value("user").(string); ok && user != "" {
		return user
	}

	// 否则根据 IP + UserAgent 生成用户标识
	if ip == "" {
		ip = "unknown"
	}
	if userAgent == "" {
		userAgent = "unknown"
	}

	// 使用 MD5 哈希生成用户标识
	h := md5.New()
	io.WriteString(h, ip)
	io.WriteString(h, userAgent)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 检查是否是新的 UV
func (l *RestoreUrlLogic) checkFirstUv(fullShortUrl, user string) bool {
	key := fmt.Sprintf("short-link:stats:uv:%s", fullShortUrl)
	added, err := l.svcCtx.BizRedis.Sadd(key, user)
	if err != nil {
		logx.Errorf("检查UV失败: %v", err)
		return false
	}
	// 设置过期时间 (90天)
	l.svcCtx.BizRedis.Expire(key, 90*24*60*60)
	return added > 0
}

// 检查是否是新的 UIP
func (l *RestoreUrlLogic) checkFirstUip(fullShortUrl, ip string) bool {
	if ip == "" {
		return false
	}
	key := fmt.Sprintf("short-link:stats:uip:%s", fullShortUrl)
	added, err := l.svcCtx.BizRedis.Sadd(key, ip)
	if err != nil {
		logx.Errorf("检查UIP失败: %v", err)
		return false
	}
	// 设置过期时间 (90天)
	l.svcCtx.BizRedis.Expire(key, 90*24*60*60)
	return added > 0
}
