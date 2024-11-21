package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-shorterurl/admin/internal/svc"
)

type CheckUsernameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUsernameLogic {
	return &CheckUsernameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckUsernameLogic) CheckUsername() (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
