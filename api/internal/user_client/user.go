// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package user_client

import (
	"context"

	"shorterurl/user/rpc/user"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CheckUsernameRequest  = user.CheckUsernameRequest
	CheckUsernameResponse = user.CheckUsernameResponse
	CommonResponse        = user.CommonResponse
	LoginRequest          = user.LoginRequest
	LoginResponse         = user.LoginResponse
	RegisterRequest       = user.RegisterRequest
	RegisterResponse      = user.RegisterResponse
	UpdatePasswordRequest = user.UpdatePasswordRequest
	UpdateRequest         = user.UpdateRequest
	UserInfoResponse      = user.UserInfoResponse

	User interface {
		// 用户注册
		RpcRegister(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
		// 用户登录
		RpcLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
		// 获取用户信息
		RpcGetUserInfo(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*UserInfoResponse, error)
		// 获取无脱敏用户信息
		RpcGetActualUserInfo(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*UserInfoResponse, error)
		// 检查用户名是否存在
		RpcCheckUsername(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CheckUsernameResponse, error)
		// 更新用户信息
		RpcUpdateUser(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*CommonResponse, error)
		// 修改密码
		RpcUpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*CommonResponse, error)
		// 检查用户是否登录
		RpcCheckLogin(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CommonResponse, error)
		// 用户退出登录
		RpcLogout(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CommonResponse, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

// 用户注册
func (m *defaultUser) RpcRegister(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcRegister(ctx, in, opts...)
}

// 用户登录
func (m *defaultUser) RpcLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcLogin(ctx, in, opts...)
}

// 获取用户信息
func (m *defaultUser) RpcGetUserInfo(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*UserInfoResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcGetUserInfo(ctx, in, opts...)
}

// 获取无脱敏用户信息
func (m *defaultUser) RpcGetActualUserInfo(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*UserInfoResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcGetActualUserInfo(ctx, in, opts...)
}

// 检查用户名是否存在
func (m *defaultUser) RpcCheckUsername(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CheckUsernameResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcCheckUsername(ctx, in, opts...)
}

// 更新用户信息
func (m *defaultUser) RpcUpdateUser(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcUpdateUser(ctx, in, opts...)
}

// 修改密码
func (m *defaultUser) RpcUpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcUpdatePassword(ctx, in, opts...)
}

// 检查用户是否登录
func (m *defaultUser) RpcCheckLogin(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcCheckLogin(ctx, in, opts...)
}

// 用户退出登录
func (m *defaultUser) RpcLogout(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.RpcLogout(ctx, in, opts...)
}