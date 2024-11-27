package server

import (
	"context"

	"shorterurl/rpc/user/internal/logic"
	"shorterurl/rpc/user/internal/svc"
	"shorterurl/rpc/user/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

// 实现 RPC 方法
func (s *UserServer) RpcRegister(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(req)
}

func (s *UserServer) RpcLogin(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(req)
}

func (s *UserServer) RpcGetUserInfo(ctx context.Context, req *user.CheckUsernameRequest) (*user.UserInfoResponse, error) {
	l := logic.NewGetUserInfoLogic(ctx, s.svcCtx)
	return l.GetUserInfo(req)
}

func (s *UserServer) RpcGetActualUserInfo(ctx context.Context, req *user.CheckUsernameRequest) (*user.UserInfoResponse, error) {
	l := logic.NewGetActualUserInfoLogic(ctx, s.svcCtx)
	return l.GetActualUserInfo(req)
}

func (s *UserServer) RpcCheckUsername(ctx context.Context, req *user.CheckUsernameRequest) (*user.CheckUsernameResponse, error) {
	l := logic.NewCheckUsernameLogic(ctx, s.svcCtx)
	return l.CheckUsername(req)
}

func (s *UserServer) RpcUpdateUser(ctx context.Context, req *user.UpdateRequest) (*user.CommonResponse, error) {
	l := logic.NewUpdateUserLogic(ctx, s.svcCtx)
	return l.UpdateUser(req)
}

func (s *UserServer) RpcUpdatePassword(ctx context.Context, req *user.UpdatePasswordRequest) (*user.CommonResponse, error) {
	l := logic.NewUpdatePasswordLogic(ctx, s.svcCtx)
	return l.UpdatePassword(req)
}

func (s *UserServer) RpcCheckLogin(ctx context.Context, req *user.CheckUsernameRequest) (*user.CommonResponse, error) {
	l := logic.NewCheckLoginLogic(ctx, s.svcCtx)
	return l.CheckLogin(req)
}

func (s *UserServer) RpcLogout(ctx context.Context, req *user.CheckUsernameRequest) (*user.CommonResponse, error) {
	l := logic.NewLogoutLogic(ctx, s.svcCtx)
	return l.Logout(req)
}
