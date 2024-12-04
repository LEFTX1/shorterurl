package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiCheckLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCheckLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCheckLoginLogic {
	return &ApiCheckLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCheckLoginLogic) ApiCheckLogin(req *types.UserCheckLoginReq) (resp *types.SuccessResp, err error) {
	// 调用rpc服务
	Response, err := l.svcCtx.UserRpc.UserCheckLogin(l.ctx,
		&userservice.CheckUsernameRequest{
			Username: req.Username,
		})
	if err != nil {
		logx.Errorf("rpc调用失败")
		return nil, err
	}

	return &types.SuccessResp{
		Code:    Response.GetMessage(),
		Success: Response.GetSuccess(),
	}, nil
}
