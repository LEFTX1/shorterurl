package user

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiUserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiUserRegisterLogic {
	return &ApiUserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiUserRegisterLogic) ApiUserRegister(req *types.UserRegisterReq) (resp *types.UserRegisterResp, err error) {
	// 调用RPC服务注册用户
	rpcResp, err := l.svcCtx.UserRpc.UserRegister(l.ctx, &userservice.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		RealName: req.RealName,
		Phone:    req.Phone,
		Mail:     req.Mail,
	})
	if err != nil {
		logx.Errorf("用户注册失败 username: %s, error: %v", req.Username, err)
		return nil, err
	}

	return &types.UserRegisterResp{
		Username:   rpcResp.Username,
		CreateTime: rpcResp.CreateTime,
	}, nil
}
