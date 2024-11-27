package logic

import (
	"context"

	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcUpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRpcUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcUpdatePasswordLogic {
	return &RpcUpdatePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改密码
func (l *RpcUpdatePasswordLogic) RpcUpdatePassword(in *user.UpdatePasswordRequest) (*user.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &user.CommonResponse{}, nil
}
