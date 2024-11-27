package user

import (
	"context"
	"shorterurl/user/rpc/user"

	"shorterurl/api/internal/svc"
	"shorterurl/api/internal/types"

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
	// 调用rpc服务
	// 用户注册
	//	RpcRegister(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Response, err := l.svcCtx.UserRpc.RpcRegister(l.ctx,
		&user.RegisterRequest{
			Username: req.Username,
			Password: req.Password,
			RealName: req.RealName,
			Phone:    req.Phone,
			Mail:     req.Mail,
		})
	if err != nil {
		logx.Errorf("rpc调用失败")
		return nil, err
	}

	return &types.UserRegisterResp{
		Username:   Response.GetUsername(),
		CreateTime: Response.GetCreateTime(),
		Message:    Response.GetMessage(),
	}, nil
}
