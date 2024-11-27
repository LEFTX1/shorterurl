package logic

import (
	"context"

	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcUpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRpcUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcUpdateUserLogic {
	return &RpcUpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户信息
func (l *RpcUpdateUserLogic) RpcUpdateUser(in *user.UpdateRequest) (*user.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &user.CommonResponse{}, nil
}
