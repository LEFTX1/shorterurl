package user

import (
	"context"

	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiUserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUserLoginLogic {
	return &ApiUserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiUserLoginLogic) ApiUserLogin(req *types.UserLoginReq) (resp *types.UserLoginResp, err error) {
	// todo: add your logic here and delete this line

	return
}
