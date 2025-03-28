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

// 短链接跳转缓存键
const GotoShortLinkKey = "link:goto:%s"

type RecycleBinSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleBinSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinSaveLogic {
	return &RecycleBinSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecycleBinSave 保存到回收站
func (l *RecycleBinSaveLogic) RecycleBinSave(in *pb.SaveToRecycleBinRequest) (*pb.SaveToRecycleBinResponse, error) {
	// 参数校验
	if in.FullShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "短链接不能为空")
	}
	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}

	l.Logger.Infof("保存短链接到回收站, 短链接: %s, 分组: %s", in.FullShortUrl, in.Gid)

	// 更新短链接状态为未启用(1)
	// 查询短链接
	link, err := l.svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(l.ctx, in.FullShortUrl, in.Gid)
	if err != nil {
		l.Logger.Errorf("查询短链接失败: %v", err)
		return nil, status.Error(codes.NotFound, "短链接不存在")
	}

	// 检查是否已经在回收站中
	if link.EnableStatus != 0 {
		l.Logger.Infof("短链接已经在回收站中: %s", in.FullShortUrl)
		return &pb.SaveToRecycleBinResponse{
			Success: true,
		}, nil
	}

	// 更新为未启用状态
	link.EnableStatus = 1
	if err := l.svcCtx.RepoManager.Link.Update(l.ctx, link); err != nil {
		l.Logger.Errorf("更新短链接状态失败: %v", err)
		return nil, status.Error(codes.Internal, "保存到回收站失败")
	}

	// 删除跳转缓存
	cacheKey := fmt.Sprintf(GotoShortLinkKey, in.FullShortUrl)
	_, err = l.svcCtx.BizRedis.Del(cacheKey)
	if err != nil {
		l.Logger.Errorf("删除短链接跳转缓存失败: %v, 键: %s", err, cacheKey)
		// 继续执行，不影响主流程
	}

	return &pb.SaveToRecycleBinResponse{
		Success: true,
	}, nil
}
