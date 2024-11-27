package user

import (
	"context"

	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUserInfoLogic {
	return &ApiUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiUserInfoLogic) ApiUserInfo(req *types.UserUsernameReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
