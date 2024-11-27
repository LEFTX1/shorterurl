package logic

import (
	"context"

	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRpcLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcLoginLogic {
	return &RpcLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登录
func (l *RpcLoginLogic) RpcLogin(in *user.LoginRequest) (*user.LoginResponse, error) {
	// todo: add your logic here and delete this line

	return &user.LoginResponse{}, nil
}
