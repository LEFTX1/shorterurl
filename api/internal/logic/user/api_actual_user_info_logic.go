package user

import (
	"context"

	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiActualUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiActualUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiActualUserInfoLogic {
	return &ApiActualUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiActualUserInfoLogic) ApiActualUserInfo(req *types.UserUsernameReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
