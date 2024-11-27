package user

import (
	"context"

	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiCheckUsernameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCheckUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCheckUsernameLogic {
	return &ApiCheckUsernameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCheckUsernameLogic) ApiCheckUsername(req *types.UserCheckUsernameReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
