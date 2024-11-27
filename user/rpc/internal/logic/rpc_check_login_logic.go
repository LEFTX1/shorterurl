package logic

import (
	"context"

	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcCheckLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRpcCheckLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcCheckLoginLogic {
	return &RpcCheckLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查用户是否登录
func (l *RpcCheckLoginLogic) RpcCheckLogin(in *user.CheckUsernameRequest) (*user.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &user.CommonResponse{}, nil
}
