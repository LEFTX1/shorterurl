package logic

import (
	"context"
	"time"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecycleBinRemoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleBinRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinRemoveLogic {
	return &RecycleBinRemoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecycleBinRemove 从回收站永久删除
func (l *RecycleBinRemoveLogic) RecycleBinRemove(in *pb.RemoveFromRecycleBinRequest) (*pb.RemoveFromRecycleBinResponse, error) {
	// 参数校验
	if in.FullShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "短链接不能为空")
	}
	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}

	l.Logger.Infof("从回收站永久删除短链接, 短链接: %s, 分组: %s", in.FullShortUrl, in.Gid)

	// 查询短链接
	link, err := l.svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(l.ctx, in.FullShortUrl, in.Gid)
	if err != nil {
		l.Logger.Errorf("查询短链接失败: %v", err)
		return nil, status.Error(codes.NotFound, "短链接不存在")
	}

	// 检查是否在回收站中（EnableStatus为1表示在回收站中）
	if link.EnableStatus != 1 {
		return nil, status.Error(codes.FailedPrecondition, "只能永久删除回收站中的短链接")
	}

	// 检查是否已经标记为永久删除
	if link.DelFlag == 1 {
		l.Logger.Infof("短链接已经标记为永久删除: %s", in.FullShortUrl)
		return &pb.RemoveFromRecycleBinResponse{
			Success: true,
		}, nil
	}

	// 执行永久删除操作（设置DelFlag = 1表示永久删除，不可恢复）
	link.DelFlag = 1
	link.DelTime = time.Now().Unix()

	if err := l.svcCtx.RepoManager.Link.Update(l.ctx, link); err != nil {
		l.Logger.Errorf("更新短链接状态失败: %v", err)
		return nil, status.Error(codes.Internal, "从回收站永久删除失败")
	}

	return &pb.RemoveFromRecycleBinResponse{
		Success: true,
	}, nil
}
