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
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ShortLinkCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkCreateLogic {
	return &ShortLinkCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --------------------- 短链接管理接口 ---------------------
func (l *ShortLinkCreateLogic) ShortLinkCreate(in *pb.CreateShortLinkRequest) (*pb.CreateShortLinkResponse, error) {
	// 参数校验
	if in.OriginUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "原始链接不能为空")
	}

	// 验证白名单
	if err := l.verificationWhitelist(in.OriginUrl); err != nil {
		return nil, err
	}

	// 获取域名，如果没有提供，使用配置中的默认域名
	domain := in.Domain
	if domain == "" {
		domain = l.svcCtx.Config.DefaultDomain
	}

	// 生成短链接后缀
	shortUri, err := l.generateShortUri(in.OriginUrl)
	if err != nil {
		l.Logger.Errorf("生成短链接失败: %v", err)
		return nil, status.Error(codes.Internal, "生成短链接失败")
	}

	// 构建完整的短链接
	fullShortUrl := util.Create(domain).Append("/").Append(shortUri).String()

	// 解析有效期
	var validDate time.Time
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

	// 创建短链接对象
	link := &model.Link{
		Domain:        domain,
		ShortUri:      shortUri,
		FullShortUrl:  fullShortUrl,
		OriginUrl:     in.OriginUrl,
		Gid:           in.Gid,
		Favicon:       util.GetFavicon(in.OriginUrl),
		EnableStatus:  0, // 默认启用
		CreatedType:   int(in.CreatedType),
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

	// 开始事务
	tx := l.svcCtx.RepoManager.GetLinkDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建短链接记录
	if err := l.svcCtx.RepoManager.Link.Create(l.ctx, link); err != nil {
		tx.Rollback()
		l.Logger.Errorf("创建短链接记录失败: %v", err)
		return nil, status.Error(codes.Internal, "创建短链接记录失败")
	}

	// 创建短链接跳转记录
	if err := l.svcCtx.RepoManager.LinkGoto.Create(l.ctx, linkGoto); err != nil {
		tx.Rollback()
		l.Logger.Errorf("创建短链接跳转记录失败: %v", err)
		return nil, status.Error(codes.Internal, "创建短链接跳转记录失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, status.Error(codes.Internal, "提交事务失败")
	}

	// 添加到布隆过滤器
	if err := l.svcCtx.BloomFilterMgr.Add(l.ctx, fullShortUrl); err != nil {
		l.Logger.Errorf("添加到布隆过滤器失败: %v", err)
		// 继续执行，不影响主流程
	}

	// 设置Redis缓存
	cacheKey := fmt.Sprintf("link:goto:%s", fullShortUrl)
	cacheExpire := util.GetLinkCacheValidTime(validDate)
	if err := l.svcCtx.BizRedis.SetexCtx(l.ctx, cacheKey, in.OriginUrl, int(cacheExpire/1000)); err != nil {
		l.Logger.Errorf("设置Redis缓存失败: %v", err)
		// 继续执行，不影响主流程
	}

	// 返回结果
	return &pb.CreateShortLinkResponse{
		FullShortUrl: "http://" + fullShortUrl,
		OriginUrl:    in.OriginUrl,
		Gid:          in.Gid,
	}, nil
}

// 验证白名单
func (l *ShortLinkCreateLogic) verificationWhitelist(originUrl string) error {
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

	// 添加日志输出
	l.Logger.Infof("正在验证域名: %s", domain)
	l.Logger.Infof("白名单域名列表: %v", details)

	for _, whiteDomain := range details {
		// 支持子域名匹配
		if domain == whiteDomain || strings.HasSuffix(domain, "."+whiteDomain) {
			l.Logger.Infof("域名 %s 匹配白名单域名 %s", domain, whiteDomain)
			return nil
		}
	}

	// 如果不在白名单中，返回错误
	errMsg := fmt.Sprintf("演示环境为避免恶意攻击，请生成以下网站跳转链接：%s", l.svcCtx.Config.GotoDomainWhiteList.Names)
	l.Logger.Errorf("域名 %s 不在白名单中", domain)
	return status.Error(codes.PermissionDenied, errMsg)
}

// generateShortUri 生成短链接后缀
func (l *ShortLinkCreateLogic) generateShortUri(originUrl string) (string, error) {
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
