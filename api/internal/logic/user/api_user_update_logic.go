package user

import (
	"context"

	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiUserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUserUpdateLogic {
	return &ApiUserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiUserUpdateLogic) ApiUserUpdate(req *types.UserUpdateReq) (resp *types.UserUpdateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
