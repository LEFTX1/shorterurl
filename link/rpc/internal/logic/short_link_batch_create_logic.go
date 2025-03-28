package logic

import (
	"context"
	"fmt"
	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/internal/types/errorx"
	"shorterurl/link/rpc/pb"
	"shorterurl/link/rpc/pkg/hash"
	"shorterurl/link/rpc/pkg/util"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ShortLinkBatchCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkBatchCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkBatchCreateLogic {
	return &ShortLinkBatchCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShortLinkBatchCreateLogic) ShortLinkBatchCreate(in *pb.BatchCreateShortLinkRequest) (*pb.BatchCreateShortLinkResponse, error) {
	// 参数校验
	if len(in.OriginUrls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "原始链接列表不能为空")
	}

	// 获取域名，如果没有提供，使用配置中的默认域名
	domain := in.Domain
	if domain == "" {
		domain = l.svcCtx.Config.DefaultDomain
	}

	// 解析有效期
	var validDate time.Time
	var err error
	if in.ValidDateType == util.ValidDateTypeCustom && in.ValidDate != "" {
		validDate, err = time.Parse(time.RFC3339, in.ValidDate)
		if err != nil {
			l.Logger.Errorf("解析有效期失败: %v", err)
			return nil, status.Error(codes.InvalidArgument, "有效期格式错误，请使用ISO-8601格式")
		}
	} else {
		// 当不是自定义有效期时，设置一个有效的默认日期，比如10年后
		validDate = time.Now().AddDate(10, 0, 0)
	}

	// 开始事务
	tx := l.svcCtx.RepoManager.GetLinkDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建返回结果
	resp := &pb.BatchCreateShortLinkResponse{
		Results: make([]*pb.BatchCreateResult, 0, len(in.OriginUrls)),
	}

	// 创建批量短链接
	var links []*model.Link
	var linkGotos []*model.LinkGoto

	// 使用for循环批量处理
	for _, originUrl := range in.OriginUrls {
		if originUrl == "" {
			continue // 跳过空链接
		}

		// 验证白名单
		if err := l.verificationWhitelist(originUrl); err != nil {
			l.Logger.Errorf("白名单验证失败: %v", err)
			continue // 跳过不符合白名单的链接
		}

		// 生成短链接后缀
		shortUri, err := l.generateShortUri(originUrl)
		if err != nil {
			l.Logger.Errorf("生成短链接失败: %v", err)
			tx.Rollback()
			return nil, status.Error(codes.Internal, "生成短链接失败: "+err.Error())
		}

		// 构建完整的短链接
		fullShortUrl := util.Create(domain).Append("/").Append(shortUri).String()

		// 创建短链接对象
		link := &model.Link{
			Domain:        domain,
			ShortUri:      shortUri,
			FullShortUrl:  fullShortUrl,
			OriginUrl:     originUrl,
			Gid:           in.Gid,
			Favicon:       util.GetFavicon(originUrl),
			EnableStatus:  0, // 默认启用
			CreatedType:   0, // 默认接口创建
			ValidDateType: int(in.ValidDateType),
			ValidDate:     validDate,
			Describe:      in.Describe,
			ClickNum:      0,
			TotalPv:       0,
			TotalUv:       0,
			TotalUip:      0,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			DelFlag:       0,
			DelTime:       0,
		}

		// 创建短链接跳转对象
		linkGoto := &model.LinkGoto{
			FullShortUrl: fullShortUrl,
			Gid:          in.Gid,
		}

		links = append(links, link)
		linkGotos = append(linkGotos, linkGoto)

		// 添加到响应结果
		resp.Results = append(resp.Results, &pb.BatchCreateResult{
			FullShortUrl: "http://" + fullShortUrl,
			OriginUrl:    originUrl,
			Gid:          in.Gid,
		})
	}

	// 如果没有有效链接，则返回空结果
	if len(links) == 0 {
		return &pb.BatchCreateShortLinkResponse{
			Results: []*pb.BatchCreateResult{},
		}, nil
	}

	// 批量创建短链接记录
	if err := l.svcCtx.RepoManager.Link.BatchCreate(l.ctx, links); err != nil {
		tx.Rollback()
		l.Logger.Errorf("批量创建短链接记录失败: %v", err)
		return nil, status.Error(codes.Internal, "批量创建短链接记录失败")
	}

	// 批量创建短链接跳转记录 - 由于表分片问题，逐个创建
	for _, linkGoto := range linkGotos {
		if err := l.svcCtx.RepoManager.LinkGoto.Create(l.ctx, linkGoto); err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建短链接跳转记录失败: %v", err)
			return nil, status.Error(codes.Internal, "创建短链接跳转记录失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, status.Error(codes.Internal, "提交事务失败")
	}

	// 异步添加到布隆过滤器和Redis缓存
	threading.GoSafe(func() {
		for i, link := range links {
			// 添加到布隆过滤器
			if err := l.svcCtx.BloomFilterMgr.Add(context.Background(), link.FullShortUrl); err != nil {
				l.Logger.Errorf("添加到布隆过滤器失败: %v", err)
			}

			// 设置Redis缓存
			cacheKey := fmt.Sprintf("link:goto:%s", link.FullShortUrl)
			cacheExpire := util.GetLinkCacheValidTime(validDate)
			if err := l.svcCtx.BizRedis.Setex(cacheKey, in.OriginUrls[i], int(cacheExpire/1000)); err != nil {
				l.Logger.Errorf("设置Redis缓存失败: %v", err)
			}
		}
	})

	return resp, nil
}

// 验证白名单
func (l *ShortLinkBatchCreateLogic) verificationWhitelist(originUrl string) error {
	// 检查白名单是否启用
	if !l.svcCtx.Config.GotoDomainWhiteList.Enable {
		return nil
	}

	// 提取域名
	domain := util.ExtractDomain(originUrl)
	if domain == "" {
		return status.Error(codes.InvalidArgument, "跳转链接填写错误")
	}

	// 检查域名是否在白名单中
	details := l.svcCtx.Config.GotoDomainWhiteList.Details
	if len(details) == 0 {
		return nil
	}

	for _, whiteDomain := range details {
		if whiteDomain == domain {
			return nil
		}
	}

	// 如果不在白名单中，返回错误
	errMsg := fmt.Sprintf("演示环境为避免恶意攻击，请生成以下网站跳转链接：%s", l.svcCtx.Config.GotoDomainWhiteList.Names)
	return status.Error(codes.PermissionDenied, errMsg)
}

// generateShortUri 生成短链接后缀
func (l *ShortLinkBatchCreateLogic) generateShortUri(originUrl string) (string, error) {
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		// 每次尝试时，为了避免冲突，添加一些随机性
		suffix := hash.HashToBase62(originUrl + fmt.Sprintf("%d", time.Now().UnixNano()))

		// 检查短链接是否已存在（先检查布隆过滤器，再查数据库）
		domain := l.svcCtx.Config.DefaultDomain
		fullShortUrl := util.Create(domain).Append("/").Append(suffix).String()

		// 检查布隆过滤器
		if exists, _ := l.svcCtx.BloomFilterMgr.Exists(l.ctx, fullShortUrl); exists {
			// 如果布隆过滤器中存在，则进一步检查数据库
			_, err := l.svcCtx.RepoManager.Link.FindByShortUri(l.ctx, suffix)
			if err == nil {
				// 如果数据库中确实存在，则继续下一次尝试
				continue
			} else if err != gorm.ErrRecordNotFound {
				// 如果是其他错误，则返回错误
				return "", err
			}
			// 如果数据库中不存在（布隆过滤器误判），则可以使用该短链接
		}

		// 如果布隆过滤器中不存在，或者数据库中不存在（布隆过滤器误判），则可以使用该短链接
		return suffix, nil
	}

	// 如果多次尝试后仍然无法生成唯一的短链接
	return "", errorx.NewCodeError(errorx.ErrShortLinkExists, errorx.ErrShortLinkExists, "短链接生成重复，请稍后再试")
}
