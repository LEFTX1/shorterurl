package user

import (
	"context"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *CheckLoginLogic) CheckLogin(req *types.UserCheckLoginReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}