package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"shorterurl/api/internal/svc"
)

type ApiLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiLogoutLogic {
	return &ApiLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiLogoutLogic) ApiLogout() (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
