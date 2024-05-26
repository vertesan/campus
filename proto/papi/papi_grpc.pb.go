// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: papi.proto

package papi

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	System_Check_FullMethodName = "/client.api.System/Check"
)

// SystemClient is the client API for System service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SystemClient interface {
	Check(ctx context.Context, in *SystemCheckRequest, opts ...grpc.CallOption) (*SystemCheckResponse, error)
}

type systemClient struct {
	cc grpc.ClientConnInterface
}

func NewSystemClient(cc grpc.ClientConnInterface) SystemClient {
	return &systemClient{cc}
}

func (c *systemClient) Check(ctx context.Context, in *SystemCheckRequest, opts ...grpc.CallOption) (*SystemCheckResponse, error) {
	out := new(SystemCheckResponse)
	err := c.cc.Invoke(ctx, System_Check_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SystemServer is the server API for System service.
// All implementations must embed UnimplementedSystemServer
// for forward compatibility
type SystemServer interface {
	Check(context.Context, *SystemCheckRequest) (*SystemCheckResponse, error)
	mustEmbedUnimplementedSystemServer()
}

// UnimplementedSystemServer must be embedded to have forward compatible implementations.
type UnimplementedSystemServer struct {
}

func (UnimplementedSystemServer) Check(context.Context, *SystemCheckRequest) (*SystemCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedSystemServer) mustEmbedUnimplementedSystemServer() {}

// UnsafeSystemServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SystemServer will
// result in compilation errors.
type UnsafeSystemServer interface {
	mustEmbedUnimplementedSystemServer()
}

func RegisterSystemServer(s grpc.ServiceRegistrar, srv SystemServer) {
	s.RegisterService(&System_ServiceDesc, srv)
}

func _System_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SystemCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: System_Check_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemServer).Check(ctx, req.(*SystemCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// System_ServiceDesc is the grpc.ServiceDesc for System service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var System_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.api.System",
	HandlerType: (*SystemServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _System_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "papi.proto",
}

const (
	Auth_Login_FullMethodName = "/client.api.Auth/Login"
)

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	Login(ctx context.Context, in *AuthLoginRequest, opts ...grpc.CallOption) (*AuthLoginResponse, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Login(ctx context.Context, in *AuthLoginRequest, opts ...grpc.CallOption) (*AuthLoginResponse, error) {
	out := new(AuthLoginResponse)
	err := c.cc.Invoke(ctx, Auth_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	Login(context.Context, *AuthLoginRequest) (*AuthLoginResponse, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) Login(context.Context, *AuthLoginRequest) (*AuthLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Login(ctx, req.(*AuthLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.api.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Auth_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "papi.proto",
}

const (
	Master_Get_FullMethodName = "/client.api.Master/Get"
)

// MasterClient is the client API for Master service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MasterClient interface {
	Get(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MasterGetResponse, error)
}

type masterClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterClient(cc grpc.ClientConnInterface) MasterClient {
	return &masterClient{cc}
}

func (c *masterClient) Get(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MasterGetResponse, error) {
	out := new(MasterGetResponse)
	err := c.cc.Invoke(ctx, Master_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterServer is the server API for Master service.
// All implementations must embed UnimplementedMasterServer
// for forward compatibility
type MasterServer interface {
	Get(context.Context, *Empty) (*MasterGetResponse, error)
	mustEmbedUnimplementedMasterServer()
}

// UnimplementedMasterServer must be embedded to have forward compatible implementations.
type UnimplementedMasterServer struct {
}

func (UnimplementedMasterServer) Get(context.Context, *Empty) (*MasterGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedMasterServer) mustEmbedUnimplementedMasterServer() {}

// UnsafeMasterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MasterServer will
// result in compilation errors.
type UnsafeMasterServer interface {
	mustEmbedUnimplementedMasterServer()
}

func RegisterMasterServer(s grpc.ServiceRegistrar, srv MasterServer) {
	s.RegisterService(&Master_ServiceDesc, srv)
}

func _Master_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Master_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).Get(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Master_ServiceDesc is the grpc.ServiceDesc for Master service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Master_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.api.Master",
	HandlerType: (*MasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Master_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "papi.proto",
}

const (
	User_Get_FullMethodName = "/client.api.User/Get"
)

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	Get(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*UserGetResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) Get(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*UserGetResponse, error) {
	out := new(UserGetResponse)
	err := c.cc.Invoke(ctx, User_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	Get(context.Context, *Empty) (*UserGetResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) Get(context.Context, *Empty) (*UserGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Get(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.api.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _User_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "papi.proto",
}

const (
	Home_Login_FullMethodName = "/client.api.Home/Login"
	Home_Enter_FullMethodName = "/client.api.Home/Enter"
)

// HomeClient is the client API for Home service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HomeClient interface {
	Login(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HomeLoginResponse, error)
	Enter(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HomeEnterResponse, error)
}

type homeClient struct {
	cc grpc.ClientConnInterface
}

func NewHomeClient(cc grpc.ClientConnInterface) HomeClient {
	return &homeClient{cc}
}

func (c *homeClient) Login(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HomeLoginResponse, error) {
	out := new(HomeLoginResponse)
	err := c.cc.Invoke(ctx, Home_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *homeClient) Enter(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HomeEnterResponse, error) {
	out := new(HomeEnterResponse)
	err := c.cc.Invoke(ctx, Home_Enter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HomeServer is the server API for Home service.
// All implementations must embed UnimplementedHomeServer
// for forward compatibility
type HomeServer interface {
	Login(context.Context, *Empty) (*HomeLoginResponse, error)
	Enter(context.Context, *Empty) (*HomeEnterResponse, error)
	mustEmbedUnimplementedHomeServer()
}

// UnimplementedHomeServer must be embedded to have forward compatible implementations.
type UnimplementedHomeServer struct {
}

func (UnimplementedHomeServer) Login(context.Context, *Empty) (*HomeLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedHomeServer) Enter(context.Context, *Empty) (*HomeEnterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enter not implemented")
}
func (UnimplementedHomeServer) mustEmbedUnimplementedHomeServer() {}

// UnsafeHomeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HomeServer will
// result in compilation errors.
type UnsafeHomeServer interface {
	mustEmbedUnimplementedHomeServer()
}

func RegisterHomeServer(s grpc.ServiceRegistrar, srv HomeServer) {
	s.RegisterService(&Home_ServiceDesc, srv)
}

func _Home_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HomeServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Home_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HomeServer).Login(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Home_Enter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HomeServer).Enter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Home_Enter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HomeServer).Enter(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Home_ServiceDesc is the grpc.ServiceDesc for Home service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Home_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.api.Home",
	HandlerType: (*HomeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Home_Login_Handler,
		},
		{
			MethodName: "Enter",
			Handler:    _Home_Enter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "papi.proto",
}

const (
	LoginBonus_Check_FullMethodName   = "/client.api.LoginBonus/Check"
	LoginBonus_Confirm_FullMethodName = "/client.api.LoginBonus/Confirm"
)

// LoginBonusClient is the client API for LoginBonus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoginBonusClient interface {
	Check(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LoginBonusCheckResponse, error)
	Confirm(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LoginBonusConfirmResponse, error)
}

type loginBonusClient struct {
	cc grpc.ClientConnInterface
}

func NewLoginBonusClient(cc grpc.ClientConnInterface) LoginBonusClient {
	return &loginBonusClient{cc}
}

func (c *loginBonusClient) Check(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LoginBonusCheckResponse, error) {
	out := new(LoginBonusCheckResponse)
	err := c.cc.Invoke(ctx, LoginBonus_Check_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loginBonusClient) Confirm(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LoginBonusConfirmResponse, error) {
	out := new(LoginBonusConfirmResponse)
	err := c.cc.Invoke(ctx, LoginBonus_Confirm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginBonusServer is the server API for LoginBonus service.
// All implementations must embed UnimplementedLoginBonusServer
// for forward compatibility
type LoginBonusServer interface {
	Check(context.Context, *Empty) (*LoginBonusCheckResponse, error)
	Confirm(context.Context, *Empty) (*LoginBonusConfirmResponse, error)
	mustEmbedUnimplementedLoginBonusServer()
}

// UnimplementedLoginBonusServer must be embedded to have forward compatible implementations.
type UnimplementedLoginBonusServer struct {
}

func (UnimplementedLoginBonusServer) Check(context.Context, *Empty) (*LoginBonusCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedLoginBonusServer) Confirm(context.Context, *Empty) (*LoginBonusConfirmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Confirm not implemented")
}
func (UnimplementedLoginBonusServer) mustEmbedUnimplementedLoginBonusServer() {}

// UnsafeLoginBonusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoginBonusServer will
// result in compilation errors.
type UnsafeLoginBonusServer interface {
	mustEmbedUnimplementedLoginBonusServer()
}

func RegisterLoginBonusServer(s grpc.ServiceRegistrar, srv LoginBonusServer) {
	s.RegisterService(&LoginBonus_ServiceDesc, srv)
}

func _LoginBonus_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginBonusServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoginBonus_Check_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginBonusServer).Check(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoginBonus_Confirm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginBonusServer).Confirm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoginBonus_Confirm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginBonusServer).Confirm(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// LoginBonus_ServiceDesc is the grpc.ServiceDesc for LoginBonus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoginBonus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.api.LoginBonus",
	HandlerType: (*LoginBonusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _LoginBonus_Check_Handler,
		},
		{
			MethodName: "Confirm",
			Handler:    _LoginBonus_Confirm_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "papi.proto",
}