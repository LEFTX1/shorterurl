package user

import (
	"context"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActualUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActualUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActualUserInfoLogic {
	return &ActualUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActualUserInfoLogic) ActualUserInfo(req *types.UserUsernameReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
