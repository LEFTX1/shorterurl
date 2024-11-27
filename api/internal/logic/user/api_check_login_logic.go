package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"shorterurl/api/internal/svc"
)

type ApiCheckLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCheckLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCheckLoginLogic {
	return &ApiCheckLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCheckLoginLogic) ApiCheckLogin() (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
