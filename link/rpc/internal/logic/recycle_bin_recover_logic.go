package logic

import (
	"context"
	"fmt"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	link, err := l.svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(l.ctx, in.FullShortUrl, in.Gid)
	if err != nil {
		l.Logger.Errorf("查询短链接失败: %v", err)
		return nil, status.Error(codes.NotFound, "短链接不存在")
	}

	// 检查是否在回收站中 (EnableStatus == 1表示在回收站中)
	if link.EnableStatus != 1 {
		l.Logger.Infof("短链接不在回收站中: %s", in.FullShortUrl)
		return &pb.RecoverFromRecycleBinResponse{
			Success: true, // 已经是正常状态，视为成功
		}, nil
	}

	// 检查是否已被永久删除 (DelFlag == 1表示已永久删除)
	if link.DelFlag == 1 {
		l.Logger.Errorf("短链接已被永久删除，无法恢复: %s", in.FullShortUrl)
		return nil, status.Error(codes.FailedPrecondition, "短链接已被永久删除，无法恢复")
	}

	// 将短链接恢复为正常状态 (设置EnableStatus = 0表示启用状态，非回收站)
	link.EnableStatus = 0
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
