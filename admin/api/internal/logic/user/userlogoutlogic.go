package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-shorterurl/admin/api/internal/svc"
)

type UserLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogoutLogic {
	return &UserLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogoutLogic) UserLogout() error {
	// todo: add your logic here and delete this line

	return nil
}
