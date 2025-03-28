package logic

import (
	"context"
	"fmt"
	"time"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

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
		// 查询Gid
		linkGoto, err := l.svcCtx.RepoManager.LinkGoto.FindByFullShortUrl(context.Background(), fullShortUrl)
		if err != nil {
			logx.Errorf("查询短链接Gid失败: %v", err)
			return
		}

		// 调用统计方法
		statsLogic := NewShortLinkStatsLogic(context.Background(), l.svcCtx)
		_, err = statsLogic.ShortLinkStats(&pb.ShortLinkStatsRequest{
			FullShortUrl: fullShortUrl,
			Gid:          linkGoto.Gid,
			// 其他统计参数可以从上下文或请求中获取
			// 在实际项目中，可以通过中间件或context传递这些信息
		})

		if err != nil {
			logx.Errorf("记录短链接访问统计失败: %v", err)
		}
	})
}
