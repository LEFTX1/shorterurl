package logic

import (
	"context"

	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcLogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRpcLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcLogoutLogic {
	return &RpcLogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户退出登录
func (l *RpcLogoutLogic) RpcLogout(in *user.CheckUsernameRequest) (*user.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &user.CommonResponse{}, nil
}
