package recycle

import (
	"context"
	"shorterurl/link/rpc/shortlinkservice"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存到回收站
func NewRecycleBinSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinSaveLogic {
	return &RecycleBinSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleBinSaveLogic) RecycleBinSave(req *types.RecycleBinOperateReq) (resp *types.SuccessResp, err error) {
	// 调用link RPC服务将短链接保存到回收站
	result, err := l.svcCtx.LinkRpc.RecycleBinSave(l.ctx, &shortlinkservice.SaveToRecycleBinRequest{
		Gid:          req.Gid,
		FullShortUrl: req.FullShortUrl,
	})
	if err != nil {
		logx.Errorf("保存短链接到回收站失败: %v", err)
		return nil, err
	}

	// 构建响应
	resp = &types.SuccessResp{
		Code:    "0",
		Success: result.Success,
	}

	return resp, nil
}
