package logic

import (
	"context"
	"fmt"
	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"shorterurl/link/rpc/pkg/util"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShortLinkUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkUpdateLogic {
	return &ShortLinkUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShortLinkUpdateLogic) ShortLinkUpdate(in *pb.UpdateShortLinkRequest) (*pb.UpdateShortLinkResponse, error) {
	// 参数校验
	if in.FullShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "短链接不能为空")
	}

	if in.OriginUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "原始链接不能为空")
	}

	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}

	// 验证白名单
	if err := l.verificationWhitelist(in.OriginUrl); err != nil {
		return nil, err
	}

	l.Logger.Infof("处理短链接更新请求，原始短链接: %s", in.FullShortUrl)

	// 处理完整短链接
	fullShortUrl := in.FullShortUrl
	// 如果URL以http://或https://开头，去掉前缀
	if strings.HasPrefix(fullShortUrl, "http://") {
		fullShortUrl = fullShortUrl[7:]
	} else if strings.HasPrefix(fullShortUrl, "https://") {
		fullShortUrl = fullShortUrl[8:]
	}

	l.Logger.Infof("去掉前缀后的短链接: %s", fullShortUrl)

	// 获取短链接的URI部分
	var shortUri string

	// 检查是否有路径分隔符
	if strings.Contains(fullShortUrl, "/") {
		parts := strings.Split(fullShortUrl, "/")
		shortUri = parts[len(parts)-1]
		l.Logger.Infof("从路径中提取的短链接URI: %s", shortUri)
	} else {
		// 没有路径分隔符，可能是纯短链接或者只有域名
		l.Logger.Infof("无法从路径中提取短链接，尝试其他方法")
		if strings.Contains(fullShortUrl, ".") {
			// 可能是域名，尝试提取最后一部分作为短链接
			parts := strings.Split(fullShortUrl, ".")
			lastPart := parts[len(parts)-1]
			if len(lastPart) > 2 {
				shortUri = lastPart
				l.Logger.Infof("从域名中提取的短链接URI: %s", shortUri)
			} else {
				// 最后一部分太短可能是域名后缀(如cn, com)，尝试直接使用全部
				shortUri = fullShortUrl
				l.Logger.Infof("直接使用完整字符串作为短链接URI: %s", shortUri)
			}
		} else {
			// 没有分隔符，直接使用
			shortUri = fullShortUrl
			l.Logger.Infof("直接使用完整字符串作为短链接URI: %s", shortUri)
		}
	}

	// 尝试查找短链接
	l.Logger.Infof("准备通过URI '%s' 和分组ID '%s' 查找短链接记录", shortUri, in.Gid)
	link, err := l.svcCtx.RepoManager.Link.FindByShortUriAndGid(l.ctx, shortUri, in.Gid)
	if err != nil {
		l.Logger.Errorf("通过URI和分组ID查找短链接失败: %v", err)

		// 如果找不到，尝试通过完整URL和分组ID查找
		l.Logger.Infof("尝试通过完整URL '%s' 和分组ID '%s' 查找短链接记录", fullShortUrl, in.Gid)
		link, err = l.svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(l.ctx, fullShortUrl, in.Gid)
		if err != nil {
			l.Logger.Errorf("通过完整URL和分组ID查找短链接失败: %v", err)
			return nil, status.Error(codes.NotFound, "短链接不存在")
		}
		l.Logger.Infof("通过完整URL和分组ID找到短链接记录")
	} else {
		l.Logger.Infof("通过URI和分组ID找到短链接记录")
	}

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

	// 记录原始分组ID，用于判断是否需要更新t_link_goto表
	oldGid := link.Gid

	// 更新链接信息
	link.OriginUrl = in.OriginUrl
	link.Gid = in.Gid
	link.ValidDateType = int(in.ValidDateType)
	link.ValidDate = validDate
	link.Describe = in.Describe
	link.UpdateTime = time.Now()

	// 开始事务，使用正确的分片数据库对象
	tx := l.svcCtx.DBs.LinkDB.WithContext(l.ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新短链接记录
	if err := l.svcCtx.RepoManager.Link.Update(l.ctx, link); err != nil {
		tx.Rollback()
		l.Logger.Errorf("更新短链接记录失败: %v", err)
		return nil, status.Error(codes.Internal, "更新短链接记录失败")
	}

	// 如果分组发生变化，需要更新短链接跳转表
	if oldGid != in.Gid {
		// 查找原有的跳转记录
		linkGoto, err := l.svcCtx.RepoManager.LinkGoto.FindByFullShortUrl(l.ctx, fullShortUrl)
		if err != nil {
			tx.Rollback()
			l.Logger.Errorf("查找短链接跳转记录失败: %v", err)
			return nil, status.Error(codes.Internal, "查找短链接跳转记录失败")
		}

		// 删除原有记录并创建新记录
		// 注意：由于LinkGotoRepo接口没有Update方法，我们使用删除旧记录并创建新记录的方式
		if err := l.svcCtx.RepoManager.LinkGoto.Delete(l.ctx, linkGoto.ID); err != nil {
			tx.Rollback()
			l.Logger.Errorf("删除旧短链接跳转记录失败: %v", err)
			return nil, status.Error(codes.Internal, "删除旧短链接跳转记录失败")
		}

		// 创建新的短链接跳转记录
		newLinkGoto := &model.LinkGoto{
			FullShortUrl: fullShortUrl, // 分片键
			Gid:          in.Gid,
		}
		if err := l.svcCtx.RepoManager.LinkGoto.Create(l.ctx, newLinkGoto); err != nil {
			tx.Rollback()
			l.Logger.Errorf("创建新短链接跳转记录失败: %v", err)
			return nil, status.Error(codes.Internal, "创建新短链接跳转记录失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		l.Logger.Errorf("提交事务失败: %v", err)
		return nil, status.Error(codes.Internal, "提交事务失败")
	}

	// 更新Redis缓存
	cacheKey := fmt.Sprintf("link:goto:%s", fullShortUrl)
	cacheExpire := util.GetLinkCacheValidTime(validDate)
	if err := l.svcCtx.BizRedis.SetexCtx(l.ctx, cacheKey, in.OriginUrl, int(cacheExpire/1000)); err != nil {
		l.Logger.Errorf("更新Redis缓存失败: %v", err)
		// 继续执行，不影响主流程
	}

	return &pb.UpdateShortLinkResponse{}, nil
}

// 验证白名单
func (l *ShortLinkUpdateLogic) verificationWhitelist(originUrl string) error {
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
