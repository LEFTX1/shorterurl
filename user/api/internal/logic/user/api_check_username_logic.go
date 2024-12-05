package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiCheckUsernameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCheckUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCheckUsernameLogic {
	return &ApiCheckUsernameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCheckUsernameLogic) ApiCheckUsername(req *types.UserCheckUsernameReq) (resp *types.SuccessResp, err error) {
	// 调用rpc服务
	Response, err := l.svcCtx.UserRpc.UserCheckUsername(l.ctx,
		&userservice.CheckUsernameRequest{
			Username: req.Username,
		})
	if err != nil {
		logx.Errorf("rpc调用失败")
		return nil, err
	}

	return &types.SuccessResp{
		Code:    "200",
		Success: Response.GetExist(),
	}, nil
}
