package user

import (
	"context"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserBasicInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserBasicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserBasicInfoLogic {
	return &GetUserBasicInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserBasicInfoLogic) GetUserBasicInfo(req *types.UserUsernameReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
