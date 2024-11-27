package user

import (
	"context"

	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiUpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUpdatePasswordLogic {
	return &ApiUpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiUpdatePasswordLogic) ApiUpdatePassword(req *types.UserUpdatePasswordReq) error {
	// todo: add your logic here and delete this line

	return nil
}
