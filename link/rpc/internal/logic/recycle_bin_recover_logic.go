package logic

import (
	"context"
	"errors"
	"fmt"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// 回收站空值键
const GotoIsNullShortLinkKey = "link:is-null:goto_%s"

type RecycleBinRecoverLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleBinRecoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinRecoverLogic {
	return &RecycleBinRecoverLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecycleBinRecover 从回收站恢复
func (l *RecycleBinRecoverLogic) RecycleBinRecover(in *pb.RecoverFromRecycleBinRequest) (*pb.RecoverFromRecycleBinResponse, error) {
	// 参数校验
	if in.FullShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "短链接不能为空")
	}
	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}

	l.Logger.Infof("从回收站恢复短链接, 短链接: %s, 分组: %s", in.FullShortUrl, in.Gid)

	// 查询短链接
	link, err := l.svcCtx.RepoManager.Link.FindRecycleBinByFullShortUrlAndGid(l.ctx, in.FullShortUrl, in.Gid)
	if err != nil {
		l.Logger.Errorf("查询回收站短链接失败: %v", err)
		// 如果未找到，可能是因为链接不在回收站，或已被物理删除，或根本不存在
		// 需要根据错误类型判断，gorm.ErrRecordNotFound 表示不在回收站 (或已被物理删除)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 尝试查询非回收站的链接，看是否是已启用状态
			activeLink, findErr := l.svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(l.ctx, in.FullShortUrl, in.Gid)
			if findErr == nil && activeLink.EnableStatus == 0 {
				// 如果找到了且是启用状态，则恢复成功
				l.Logger.Infof("短链接已是启用状态(非回收站): %s", in.FullShortUrl)
				return &pb.RecoverFromRecycleBinResponse{
					Success: true,
				}, nil
			}
			// 否则，视为短链接不存在于回收站中
			return nil, status.Error(codes.NotFound, "短链接不在回收站中或不存在")
		}
		// 其他查询错误
		return nil, status.Error(codes.Internal, "查询回收站短链接失败")
	}

	// 检查是否 *确实* 在回收站中 (del_flag == 1)
	// FindRecycleBinByFullShortUrlAndGid 已经保证了这一点，所以无需额外检查 link.DelFlag

	// 更新为启用状态 (del_flag = 0, enable_status = 0)
	link.EnableStatus = 0
	link.DelFlag = 0
	link.DelTime = 0 // 重置删除时间
	if err := l.svcCtx.RepoManager.Link.Update(l.ctx, link); err != nil {
		l.Logger.Errorf("更新短链接状态失败: %v", err)
		return nil, status.Error(codes.Internal, "从回收站恢复失败")
	}

	// 删除空值缓存，以便能够重新使用
	cacheKey := fmt.Sprintf(GotoIsNullShortLinkKey, in.FullShortUrl)
	_, err = l.svcCtx.BizRedis.Del(cacheKey)
	if err != nil {
		l.Logger.Errorf("删除短链接空值缓存失败: %v, 键: %s", err, cacheKey)
		// 继续执行，不影响主流程
	}

	return &pb.RecoverFromRecycleBinResponse{
		Success: true,
	}, nil
}
