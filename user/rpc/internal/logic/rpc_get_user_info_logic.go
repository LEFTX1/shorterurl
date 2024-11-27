package logic

import (
	"context"

	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcGetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRpcGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcGetUserInfoLogic {
	return &RpcGetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息
func (l *RpcGetUserInfoLogic) RpcGetUserInfo(in *user.CheckUsernameRequest) (*user.UserInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &user.UserInfoResponse{}, nil
}
