//protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=./  --go-grpc_opt=paths=source_relative ./internal/transport/grpc/server.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: internal/transport/grpc/server.proto

package gshort

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
	GRPC_GetBaseURL_FullMethodName     = "/grpc.gshort.GRPC/GetBaseURL"
	GRPC_Ping_FullMethodName           = "/grpc.gshort.GRPC/Ping"
	GRPC_SetURL_FullMethodName         = "/grpc.gshort.GRPC/SetURL"
	GRPC_SetUrls_FullMethodName        = "/grpc.gshort.GRPC/SetUrls"
	GRPC_GetUserURLS_FullMethodName    = "/grpc.gshort.GRPC/GetUserURLS"
	GRPC_DeleteUserURLS_FullMethodName = "/grpc.gshort.GRPC/DeleteUserURLS"
	GRPC_GetStats_FullMethodName       = "/grpc.gshort.GRPC/GetStats"
)

// GRPCClient is the client API for GRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GRPCClient interface {
	GetBaseURL(ctx context.Context, in *ShortUrlRequest, opts ...grpc.CallOption) (*UrlResponse, error)
	Ping(ctx context.Context, in *GetPingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	SetURL(ctx context.Context, in *UrlRequest, opts ...grpc.CallOption) (*UrlResponse, error)
	SetUrls(ctx context.Context, in *ShortenBatchRequest, opts ...grpc.CallOption) (*ShortenBatchRequest, error)
	GetUserURLS(ctx context.Context, in *GetUserURLSRequest, opts ...grpc.CallOption) (*UserUrlsResponse, error)
	DeleteUserURLS(ctx context.Context, in *DeleteUserURLSRequest, opts ...grpc.CallOption) (*DeleteUserURLSResponse, error)
	GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
}

type gRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewGRPCClient(cc grpc.ClientConnInterface) GRPCClient {
	return &gRPCClient{cc}
}

func (c *gRPCClient) GetBaseURL(ctx context.Context, in *ShortUrlRequest, opts ...grpc.CallOption) (*UrlResponse, error) {
	out := new(UrlResponse)
	err := c.cc.Invoke(ctx, GRPC_GetBaseURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCClient) Ping(ctx context.Context, in *GetPingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, GRPC_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCClient) SetURL(ctx context.Context, in *UrlRequest, opts ...grpc.CallOption) (*UrlResponse, error) {
	out := new(UrlResponse)
	err := c.cc.Invoke(ctx, GRPC_SetURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCClient) SetUrls(ctx context.Context, in *ShortenBatchRequest, opts ...grpc.CallOption) (*ShortenBatchRequest, error) {
	out := new(ShortenBatchRequest)
	err := c.cc.Invoke(ctx, GRPC_SetUrls_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCClient) GetUserURLS(ctx context.Context, in *GetUserURLSRequest, opts ...grpc.CallOption) (*UserUrlsResponse, error) {
	out := new(UserUrlsResponse)
	err := c.cc.Invoke(ctx, GRPC_GetUserURLS_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCClient) DeleteUserURLS(ctx context.Context, in *DeleteUserURLSRequest, opts ...grpc.CallOption) (*DeleteUserURLSResponse, error) {
	out := new(DeleteUserURLSResponse)
	err := c.cc.Invoke(ctx, GRPC_DeleteUserURLS_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCClient) GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, GRPC_GetStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GRPCServer is the server API for GRPC service.
// All implementations must embed UnimplementedGRPCServer
// for forward compatibility
type GRPCServer interface {
	GetBaseURL(context.Context, *ShortUrlRequest) (*UrlResponse, error)
	Ping(context.Context, *GetPingRequest) (*PingResponse, error)
	SetURL(context.Context, *UrlRequest) (*UrlResponse, error)
	SetUrls(context.Context, *ShortenBatchRequest) (*ShortenBatchRequest, error)
	GetUserURLS(context.Context, *GetUserURLSRequest) (*UserUrlsResponse, error)
	DeleteUserURLS(context.Context, *DeleteUserURLSRequest) (*DeleteUserURLSResponse, error)
	GetStats(context.Context, *GetStatsRequest) (*StatsResponse, error)
	mustEmbedUnimplementedGRPCServer()
}

// UnimplementedGRPCServer must be embedded to have forward compatible implementations.
type UnimplementedGRPCServer struct {
}

func (UnimplementedGRPCServer) GetBaseURL(context.Context, *ShortUrlRequest) (*UrlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBaseURL not implemented")
}
func (UnimplementedGRPCServer) Ping(context.Context, *GetPingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedGRPCServer) SetURL(context.Context, *UrlRequest) (*UrlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetURL not implemented")
}
func (UnimplementedGRPCServer) SetUrls(context.Context, *ShortenBatchRequest) (*ShortenBatchRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUrls not implemented")
}
func (UnimplementedGRPCServer) GetUserURLS(context.Context, *GetUserURLSRequest) (*UserUrlsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserURLS not implemented")
}
func (UnimplementedGRPCServer) DeleteUserURLS(context.Context, *DeleteUserURLSRequest) (*DeleteUserURLSResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserURLS not implemented")
}
func (UnimplementedGRPCServer) GetStats(context.Context, *GetStatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedGRPCServer) mustEmbedUnimplementedGRPCServer() {}

// UnsafeGRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GRPCServer will
// result in compilation errors.
type UnsafeGRPCServer interface {
	mustEmbedUnimplementedGRPCServer()
}

func RegisterGRPCServer(s grpc.ServiceRegistrar, srv GRPCServer) {
	s.RegisterService(&GRPC_ServiceDesc, srv)
}

func _GRPC_GetBaseURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortUrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCServer).GetBaseURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPC_GetBaseURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCServer).GetBaseURL(ctx, req.(*ShortUrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPC_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPC_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCServer).Ping(ctx, req.(*GetPingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPC_SetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCServer).SetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPC_SetURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCServer).SetURL(ctx, req.(*UrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPC_SetUrls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCServer).SetUrls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPC_SetUrls_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCServer).SetUrls(ctx, req.(*ShortenBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPC_GetUserURLS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserURLSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCServer).GetUserURLS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPC_GetUserURLS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCServer).GetUserURLS(ctx, req.(*GetUserURLSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPC_DeleteUserURLS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserURLSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCServer).DeleteUserURLS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPC_DeleteUserURLS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCServer).DeleteUserURLS(ctx, req.(*DeleteUserURLSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPC_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPC_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCServer).GetStats(ctx, req.(*GetStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GRPC_ServiceDesc is the grpc.ServiceDesc for GRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.gshort.GRPC",
	HandlerType: (*GRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBaseURL",
			Handler:    _GRPC_GetBaseURL_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _GRPC_Ping_Handler,
		},
		{
			MethodName: "SetURL",
			Handler:    _GRPC_SetURL_Handler,
		},
		{
			MethodName: "SetUrls",
			Handler:    _GRPC_SetUrls_Handler,
		},
		{
			MethodName: "GetUserURLS",
			Handler:    _GRPC_GetUserURLS_Handler,
		},
		{
			MethodName: "DeleteUserURLS",
			Handler:    _GRPC_DeleteUserURLS_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _GRPC_GetStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/transport/grpc/server.proto",
}
