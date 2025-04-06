package recycle

import (
	"context"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 移除短链接
func NewRecycleBinRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinRemoveLogic {
	return &RecycleBinRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleBinRemoveLogic) RecycleBinRemove(req *types.RecycleBinOperateReq) (resp *types.SuccessResp, err error) {
	// 调用link RPC服务从回收站删除短链接
	result, err := l.svcCtx.LinkRpc.RecycleBinRemove(l.ctx, &shortlinkservice.RemoveFromRecycleBinRequest{
		Gid:          req.Gid,
		FullShortUrl: req.FullShortUrl,
	})
	if err != nil {
		logx.Errorf("从回收站删除短链接失败: %v", err)
		return nil, err
	}

	// 构建响应
	resp = &types.SuccessResp{
		Code:    "0",
		Success: result.Success,
	}

	return resp, nil
}
