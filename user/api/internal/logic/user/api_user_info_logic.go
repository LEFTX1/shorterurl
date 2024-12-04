package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/rpc/userservice"

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
	// 调用RPC服务获取用户信息
	rpcResp, err := l.svcCtx.UserRpc.UserGetInfo(l.ctx, &userservice.CheckUsernameRequest{
		Username: req.Username,
	})
	if err != nil {
		logx.Errorf("获取用户信息失败 username: %s, error: %v", req.Username, err)
		return nil, err
	}

	return &types.UserInfoResp{
		Username:   rpcResp.Username,
		RealName:   rpcResp.RealName,
		Phone:      rpcResp.Phone,
		Mail:       rpcResp.Mail,
		CreateTime: rpcResp.CreateTime,
		UpdateTime: rpcResp.UpdateTime,
	}, nil
}
