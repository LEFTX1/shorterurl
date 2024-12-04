package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/rpc/userservice"

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

func (l *ApiUserUpdateLogic) ApiUserUpdate(req *types.UserUpdateReq) (resp *types.SuccessResp, err error) {
	// 调用RPC服务更新用户信息
	_, err = l.svcCtx.UserRpc.UserUpdate(l.ctx, &userservice.UpdateRequest{
		Username: req.Username,
		Password: req.Password,
		RealName: req.RealName,
		Phone:    req.Phone,
		Mail:     req.Mail,
	})
	if err != nil {
		logx.Errorf("更新用户信息失败 username: %s, error: %v", req.Username, err)
		return nil, err
	}

	return &types.SuccessResp{
		Code:    "0",
		Success: true,
	}, nil
}
