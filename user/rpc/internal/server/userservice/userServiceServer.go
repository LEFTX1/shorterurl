// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"shorterurl/user/rpc/internal/logic/userservice"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/pb"
)

type UserServiceServer struct {
	svcCtx *svc.ServiceContext
	__.UnimplementedUserServiceServer
}

func NewUserServiceServer(svcCtx *svc.ServiceContext) *UserServiceServer {
	return &UserServiceServer{
		svcCtx: svcCtx,
	}
}

// 用户注册
func (s *UserServiceServer) UserRegister(ctx context.Context, in *__.RegisterRequest) (*__.RegisterResponse, error) {
	l := userservicelogic.NewUserRegisterLogic(ctx, s.svcCtx)
	return l.UserRegister(in)
}

// 用户登录
func (s *UserServiceServer) UserLogin(ctx context.Context, in *__.LoginRequest) (*__.LoginResponse, error) {
	l := userservicelogic.NewUserLoginLogic(ctx, s.svcCtx)
	return l.UserLogin(in)
}

// 获取用户信息
func (s *UserServiceServer) UserGetInfo(ctx context.Context, in *__.CheckUsernameRequest) (*__.UserInfoResponse, error) {
	l := userservicelogic.NewUserGetInfoLogic(ctx, s.svcCtx)
	return l.UserGetInfo(in)
}

// 获取无脱敏用户信息
func (s *UserServiceServer) UserGetActualInfo(ctx context.Context, in *__.CheckUsernameRequest) (*__.UserInfoResponse, error) {
	l := userservicelogic.NewUserGetActualInfoLogic(ctx, s.svcCtx)
	return l.UserGetActualInfo(in)
}

// 检查用户名是否存在
func (s *UserServiceServer) UserCheckUsername(ctx context.Context, in *__.CheckUsernameRequest) (*__.CheckUsernameResponse, error) {
	l := userservicelogic.NewUserCheckUsernameLogic(ctx, s.svcCtx)
	return l.UserCheckUsername(in)
}

// 更新用户信息
func (s *UserServiceServer) UserUpdate(ctx context.Context, in *__.UpdateRequest) (*__.CommonResponse, error) {
	l := userservicelogic.NewUserUpdateLogic(ctx, s.svcCtx)
	return l.UserUpdate(in)
}

// 检查用户是否登录
func (s *UserServiceServer) UserCheckLogin(ctx context.Context, in *__.LogoutRequest) (*__.CommonResponse, error) {
	l := userservicelogic.NewUserCheckLoginLogic(ctx, s.svcCtx)
	return l.UserCheckLogin(in)
}

// 用户退出登录
func (s *UserServiceServer) UserLogout(ctx context.Context, in *__.LogoutRequest) (*__.CommonResponse, error) {
	l := userservicelogic.NewUserLogoutLogic(ctx, s.svcCtx)
	return l.UserLogout(in)
}