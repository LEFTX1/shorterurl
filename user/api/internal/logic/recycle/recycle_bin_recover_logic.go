package recycle

import (
	"context"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinRecoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 恢复短链接
func NewRecycleBinRecoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinRecoverLogic {
	return &RecycleBinRecoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleBinRecoverLogic) RecycleBinRecover(req *types.RecycleBinOperateReq) (resp *types.SuccessResp, err error) {
	// 调用link RPC服务从回收站恢复短链接
	result, err := l.svcCtx.LinkRpc.RecycleBinRecover(l.ctx, &shortlinkservice.RecoverFromRecycleBinRequest{
		Gid:          req.Gid,
		FullShortUrl: req.FullShortUrl,
	})
	if err != nil {
		logx.Errorf("从回收站恢复短链接失败: %v", err)
		return nil, err
	}

	// 构建响应
	resp = &types.SuccessResp{
		Code:    "0",
		Success: result.Success,
	}

	return resp, nil
}
