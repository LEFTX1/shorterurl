package logic

import (
	"context"

	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcCheckUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRpcCheckUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcCheckUsernameLogic {
	return &RpcCheckUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查用户名是否存在
func (l *RpcCheckUsernameLogic) RpcCheckUsername(in *user.CheckUsernameRequest) (*user.CheckUsernameResponse, error) {
	// todo: add your logic here and delete this line

	return &user.CheckUsernameResponse{}, nil
}
