package logic

import (
	"context"

	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcGetActualUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRpcGetActualUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcGetActualUserInfoLogic {
	return &RpcGetActualUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取无脱敏用户信息
func (l *RpcGetActualUserInfoLogic) RpcGetActualUserInfo(in *user.CheckUsernameRequest) (*user.UserInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &user.UserInfoResponse{}, nil
}
