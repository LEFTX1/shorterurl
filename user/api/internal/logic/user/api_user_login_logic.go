package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/rpc/userservice"

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
	// 调用RPC服务登录
	rpcResp, err := l.svcCtx.UserRpc.UserLogin(l.ctx, &userservice.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		logx.Errorf("用户登录失败 username: %s, error: %v", req.Username, err)
		return nil, err
	}

	return &types.UserLoginResp{
		Token:      rpcResp.Token,
		Username:   rpcResp.Username,
		RealName:   rpcResp.RealName,
		CreateTime: rpcResp.CreateTime,
	}, nil
}
