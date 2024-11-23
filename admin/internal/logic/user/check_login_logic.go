package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-shorterurl/admin/internal/svc"
)

type CheckLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckLoginLogic {
	return &CheckLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckLoginLogic) CheckLogin() (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
