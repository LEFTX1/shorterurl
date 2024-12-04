package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/rpc/userservice"

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
	// 调用rpc服务
	Response, err := l.svcCtx.UserRpc.UserGetActualInfo(l.ctx,
		&userservice.CheckUsernameRequest{
			Username: req.Username,
		})
	if err != nil {
		logx.Errorf("rpc调用失败")
		return nil, err
	}

	return &types.UserInfoResp{
		Username:   Response.GetUsername(),
		RealName:   Response.GetRealName(),
		Phone:      Response.GetPhone(),
		Mail:       Response.GetMail(),
		CreateTime: Response.GetCreateTime(),
		UpdateTime: Response.GetUpdateTime(),
	}, nil
}
