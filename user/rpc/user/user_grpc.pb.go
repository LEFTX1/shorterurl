// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.0--rc1
// source: user.proto

package user

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	User_RpcRegister_FullMethodName          = "/user.User/RpcRegister"
	User_RpcLogin_FullMethodName             = "/user.User/RpcLogin"
	User_RpcGetUserInfo_FullMethodName       = "/user.User/RpcGetUserInfo"
	User_RpcGetActualUserInfo_FullMethodName = "/user.User/RpcGetActualUserInfo"
	User_RpcCheckUsername_FullMethodName     = "/user.User/RpcCheckUsername"
	User_RpcUpdateUser_FullMethodName        = "/user.User/RpcUpdateUser"
	User_RpcUpdatePassword_FullMethodName    = "/user.User/RpcUpdatePassword"
	User_RpcCheckLogin_FullMethodName        = "/user.User/RpcCheckLogin"
	User_RpcLogout_FullMethodName            = "/user.User/RpcLogout"
)

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 用户服务
type UserClient interface {
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

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) RpcRegister(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, User_RpcRegister_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RpcLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, User_RpcLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RpcGetUserInfo(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*UserInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserInfoResponse)
	err := c.cc.Invoke(ctx, User_RpcGetUserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RpcGetActualUserInfo(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*UserInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserInfoResponse)
	err := c.cc.Invoke(ctx, User_RpcGetActualUserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RpcCheckUsername(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CheckUsernameResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckUsernameResponse)
	err := c.cc.Invoke(ctx, User_RpcCheckUsername_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RpcUpdateUser(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, User_RpcUpdateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RpcUpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, User_RpcUpdatePassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RpcCheckLogin(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, User_RpcCheckLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RpcLogout(ctx context.Context, in *CheckUsernameRequest, opts ...grpc.CallOption) (*CommonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommonResponse)
	err := c.cc.Invoke(ctx, User_RpcLogout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility.
//
// 用户服务
type UserServer interface {
	// 用户注册
	RpcRegister(context.Context, *RegisterRequest) (*RegisterResponse, error)
	// 用户登录
	RpcLogin(context.Context, *LoginRequest) (*LoginResponse, error)
	// 获取用户信息
	RpcGetUserInfo(context.Context, *CheckUsernameRequest) (*UserInfoResponse, error)
	// 获取无脱敏用户信息
	RpcGetActualUserInfo(context.Context, *CheckUsernameRequest) (*UserInfoResponse, error)
	// 检查用户名是否存在
	RpcCheckUsername(context.Context, *CheckUsernameRequest) (*CheckUsernameResponse, error)
	// 更新用户信息
	RpcUpdateUser(context.Context, *UpdateRequest) (*CommonResponse, error)
	// 修改密码
	RpcUpdatePassword(context.Context, *UpdatePasswordRequest) (*CommonResponse, error)
	// 检查用户是否登录
	RpcCheckLogin(context.Context, *CheckUsernameRequest) (*CommonResponse, error)
	// 用户退出登录
	RpcLogout(context.Context, *CheckUsernameRequest) (*CommonResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServer struct{}

func (UnimplementedUserServer) RpcRegister(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcRegister not implemented")
}
func (UnimplementedUserServer) RpcLogin(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcLogin not implemented")
}
func (UnimplementedUserServer) RpcGetUserInfo(context.Context, *CheckUsernameRequest) (*UserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcGetUserInfo not implemented")
}
func (UnimplementedUserServer) RpcGetActualUserInfo(context.Context, *CheckUsernameRequest) (*UserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcGetActualUserInfo not implemented")
}
func (UnimplementedUserServer) RpcCheckUsername(context.Context, *CheckUsernameRequest) (*CheckUsernameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcCheckUsername not implemented")
}
func (UnimplementedUserServer) RpcUpdateUser(context.Context, *UpdateRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcUpdateUser not implemented")
}
func (UnimplementedUserServer) RpcUpdatePassword(context.Context, *UpdatePasswordRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcUpdatePassword not implemented")
}
func (UnimplementedUserServer) RpcCheckLogin(context.Context, *CheckUsernameRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcCheckLogin not implemented")
}
func (UnimplementedUserServer) RpcLogout(context.Context, *CheckUsernameRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RpcLogout not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}
func (UnimplementedUserServer) testEmbeddedByValue()              {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	// If the following call pancis, it indicates UnimplementedUserServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_RpcRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcRegister_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcRegister(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RpcLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcLogin(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RpcGetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcGetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcGetUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcGetUserInfo(ctx, req.(*CheckUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RpcGetActualUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcGetActualUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcGetActualUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcGetActualUserInfo(ctx, req.(*CheckUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RpcCheckUsername_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcCheckUsername(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcCheckUsername_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcCheckUsername(ctx, req.(*CheckUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RpcUpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcUpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcUpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcUpdateUser(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RpcUpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcUpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcUpdatePassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcUpdatePassword(ctx, req.(*UpdatePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RpcCheckLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcCheckLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcCheckLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcCheckLogin(ctx, req.(*CheckUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RpcLogout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RpcLogout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RpcLogout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RpcLogout(ctx, req.(*CheckUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RpcRegister",
			Handler:    _User_RpcRegister_Handler,
		},
		{
			MethodName: "RpcLogin",
			Handler:    _User_RpcLogin_Handler,
		},
		{
			MethodName: "RpcGetUserInfo",
			Handler:    _User_RpcGetUserInfo_Handler,
		},
		{
			MethodName: "RpcGetActualUserInfo",
			Handler:    _User_RpcGetActualUserInfo_Handler,
		},
		{
			MethodName: "RpcCheckUsername",
			Handler:    _User_RpcCheckUsername_Handler,
		},
		{
			MethodName: "RpcUpdateUser",
			Handler:    _User_RpcUpdateUser_Handler,
		},
		{
			MethodName: "RpcUpdatePassword",
			Handler:    _User_RpcUpdatePassword_Handler,
		},
		{
			MethodName: "RpcCheckLogin",
			Handler:    _User_RpcCheckLogin_Handler,
		},
		{
			MethodName: "RpcLogout",
			Handler:    _User_RpcLogout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}