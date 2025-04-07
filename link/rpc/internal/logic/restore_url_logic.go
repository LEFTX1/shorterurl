package logic

import (
	"context"
	"fmt"
	"io"
	"time"

	"shorterurl/link/rpc/internal/consumer"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"crypto/md5"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/codes"
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

		// 从上下文中获取请求信息（在实际项目中可能需要通过中间件传递）
		// 这里模拟获取，实际项目中应从 context 或请求中获取
		ip := l.getValueFromContext(l.ctx, "ip", "")
		browser := l.getValueFromContext(l.ctx, "browser", "未知浏览器")
		os := l.getValueFromContext(l.ctx, "os", "未知系统")
		device := l.getValueFromContext(l.ctx, "device", "未知设备")
		network := l.getValueFromContext(l.ctx, "network", "未知网络")
		userAgent := l.getValueFromContext(l.ctx, "user-agent", "")

		// 获取或生成用户标识（可以是 cookie 中的值或根据 IP+UserAgent 生成的哈希）
		user := l.getUserIdentifier(ip, userAgent)

		// 检查是否是新的 UV 和 UIP
		uvFirstFlag := l.checkFirstUv(fullShortUrl, user)
		uipFirstFlag := l.checkFirstUip(fullShortUrl, ip)

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
			logx.Infof("开始获取访问IP的地理位置信息: %s", ip)

			// 创建IP位置查询逻辑
			ipLocationLogic := NewGetIpLocationLogic(ctx, l.svcCtx)
			formattedLocation, err := ipLocationLogic.GetFormattedLocation(ip)

			if err == nil && formattedLocation != "" {
				logx.Infof("IP地理位置解析成功: %s -> %s", ip, formattedLocation)
				statsRecord.Locale = formattedLocation
			} else {
				logx.Errorf("获取IP地理位置信息失败: %v", err)
			}
		}

		// 提交统计记录到消费者队列
		l.svcCtx.StatsConsumer.Submit(statsRecord)
	})
}

// 从上下文中获取值，如果不存在则返回默认值
func (l *RestoreUrlLogic) getValueFromContext(ctx context.Context, key, defaultValue string) string {
	if value, ok := ctx.Value(key).(string); ok && value != "" {
		return value
	}
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
