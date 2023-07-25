// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: db.proto

package glist

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
	GListService_LPush_FullMethodName  = "/glist.GListService/LPush"
	GListService_LPushs_FullMethodName = "/glist.GListService/LPushs"
	GListService_RPush_FullMethodName  = "/glist.GListService/RPush"
	GListService_RPushs_FullMethodName = "/glist.GListService/RPushs"
	GListService_LPop_FullMethodName   = "/glist.GListService/LPop"
	GListService_RPop_FullMethodName   = "/glist.GListService/RPop"
	GListService_LRange_FullMethodName = "/glist.GListService/LRange"
	GListService_LLen_FullMethodName   = "/glist.GListService/LLen"
	GListService_LRem_FullMethodName   = "/glist.GListService/LRem"
	GListService_LIndex_FullMethodName = "/glist.GListService/LIndex"
	GListService_LSet_FullMethodName   = "/glist.GListService/LSet"
	GListService_LTrim_FullMethodName  = "/glist.GListService/LTrim"
)

// GListServiceClient is the client API for GListService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GListServiceClient interface {
	// example
	LPush(ctx context.Context, in *GListLPushRequest, opts ...grpc.CallOption) (*GListLPushResponse, error)
	LPushs(ctx context.Context, in *GListLPushsRequest, opts ...grpc.CallOption) (*GListLPushsResponse, error)
	RPush(ctx context.Context, in *GListRPushRequest, opts ...grpc.CallOption) (*GListRPushResponse, error)
	RPushs(ctx context.Context, in *GListRPushsRequest, opts ...grpc.CallOption) (*GListRPushsResponse, error)
	LPop(ctx context.Context, in *GListLPopRequest, opts ...grpc.CallOption) (*GListLPopResponse, error)
	RPop(ctx context.Context, in *GListRPopRequest, opts ...grpc.CallOption) (*GListRPopResponse, error)
	LRange(ctx context.Context, in *GListLRangeRequest, opts ...grpc.CallOption) (*GListLRangeResponse, error)
	LLen(ctx context.Context, in *GListLLenRequest, opts ...grpc.CallOption) (*GListLLenResponse, error)
	LRem(ctx context.Context, in *GListLRemRequest, opts ...grpc.CallOption) (*GListLRemResponse, error)
	LIndex(ctx context.Context, in *GListLIndexRequest, opts ...grpc.CallOption) (*GListLIndexResponse, error)
	LSet(ctx context.Context, in *GListLSetRequest, opts ...grpc.CallOption) (*GListLSetResponse, error)
	LTrim(ctx context.Context, in *GListLTrimRequest, opts ...grpc.CallOption) (*GListLTrimResponse, error)
}

type gListServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGListServiceClient(cc grpc.ClientConnInterface) GListServiceClient {
	return &gListServiceClient{cc}
}

func (c *gListServiceClient) LPush(ctx context.Context, in *GListLPushRequest, opts ...grpc.CallOption) (*GListLPushResponse, error) {
	out := new(GListLPushResponse)
	err := c.cc.Invoke(ctx, GListService_LPush_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) LPushs(ctx context.Context, in *GListLPushsRequest, opts ...grpc.CallOption) (*GListLPushsResponse, error) {
	out := new(GListLPushsResponse)
	err := c.cc.Invoke(ctx, GListService_LPushs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) RPush(ctx context.Context, in *GListRPushRequest, opts ...grpc.CallOption) (*GListRPushResponse, error) {
	out := new(GListRPushResponse)
	err := c.cc.Invoke(ctx, GListService_RPush_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) RPushs(ctx context.Context, in *GListRPushsRequest, opts ...grpc.CallOption) (*GListRPushsResponse, error) {
	out := new(GListRPushsResponse)
	err := c.cc.Invoke(ctx, GListService_RPushs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) LPop(ctx context.Context, in *GListLPopRequest, opts ...grpc.CallOption) (*GListLPopResponse, error) {
	out := new(GListLPopResponse)
	err := c.cc.Invoke(ctx, GListService_LPop_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) RPop(ctx context.Context, in *GListRPopRequest, opts ...grpc.CallOption) (*GListRPopResponse, error) {
	out := new(GListRPopResponse)
	err := c.cc.Invoke(ctx, GListService_RPop_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) LRange(ctx context.Context, in *GListLRangeRequest, opts ...grpc.CallOption) (*GListLRangeResponse, error) {
	out := new(GListLRangeResponse)
	err := c.cc.Invoke(ctx, GListService_LRange_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) LLen(ctx context.Context, in *GListLLenRequest, opts ...grpc.CallOption) (*GListLLenResponse, error) {
	out := new(GListLLenResponse)
	err := c.cc.Invoke(ctx, GListService_LLen_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) LRem(ctx context.Context, in *GListLRemRequest, opts ...grpc.CallOption) (*GListLRemResponse, error) {
	out := new(GListLRemResponse)
	err := c.cc.Invoke(ctx, GListService_LRem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) LIndex(ctx context.Context, in *GListLIndexRequest, opts ...grpc.CallOption) (*GListLIndexResponse, error) {
	out := new(GListLIndexResponse)
	err := c.cc.Invoke(ctx, GListService_LIndex_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) LSet(ctx context.Context, in *GListLSetRequest, opts ...grpc.CallOption) (*GListLSetResponse, error) {
	out := new(GListLSetResponse)
	err := c.cc.Invoke(ctx, GListService_LSet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gListServiceClient) LTrim(ctx context.Context, in *GListLTrimRequest, opts ...grpc.CallOption) (*GListLTrimResponse, error) {
	out := new(GListLTrimResponse)
	err := c.cc.Invoke(ctx, GListService_LTrim_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GListServiceServer is the server API for GListService service.
// All implementations must embed UnimplementedGListServiceServer
// for forward compatibility
type GListServiceServer interface {
	// example
	LPush(context.Context, *GListLPushRequest) (*GListLPushResponse, error)
	LPushs(context.Context, *GListLPushsRequest) (*GListLPushsResponse, error)
	RPush(context.Context, *GListRPushRequest) (*GListRPushResponse, error)
	RPushs(context.Context, *GListRPushsRequest) (*GListRPushsResponse, error)
	LPop(context.Context, *GListLPopRequest) (*GListLPopResponse, error)
	RPop(context.Context, *GListRPopRequest) (*GListRPopResponse, error)
	LRange(context.Context, *GListLRangeRequest) (*GListLRangeResponse, error)
	LLen(context.Context, *GListLLenRequest) (*GListLLenResponse, error)
	LRem(context.Context, *GListLRemRequest) (*GListLRemResponse, error)
	LIndex(context.Context, *GListLIndexRequest) (*GListLIndexResponse, error)
	LSet(context.Context, *GListLSetRequest) (*GListLSetResponse, error)
	LTrim(context.Context, *GListLTrimRequest) (*GListLTrimResponse, error)
	mustEmbedUnimplementedGListServiceServer()
}

// UnimplementedGListServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGListServiceServer struct {
}

func (UnimplementedGListServiceServer) LPush(context.Context, *GListLPushRequest) (*GListLPushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LPush not implemented")
}
func (UnimplementedGListServiceServer) LPushs(context.Context, *GListLPushsRequest) (*GListLPushsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LPushs not implemented")
}
func (UnimplementedGListServiceServer) RPush(context.Context, *GListRPushRequest) (*GListRPushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPush not implemented")
}
func (UnimplementedGListServiceServer) RPushs(context.Context, *GListRPushsRequest) (*GListRPushsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPushs not implemented")
}
func (UnimplementedGListServiceServer) LPop(context.Context, *GListLPopRequest) (*GListLPopResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LPop not implemented")
}
func (UnimplementedGListServiceServer) RPop(context.Context, *GListRPopRequest) (*GListRPopResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPop not implemented")
}
func (UnimplementedGListServiceServer) LRange(context.Context, *GListLRangeRequest) (*GListLRangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LRange not implemented")
}
func (UnimplementedGListServiceServer) LLen(context.Context, *GListLLenRequest) (*GListLLenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LLen not implemented")
}
func (UnimplementedGListServiceServer) LRem(context.Context, *GListLRemRequest) (*GListLRemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LRem not implemented")
}
func (UnimplementedGListServiceServer) LIndex(context.Context, *GListLIndexRequest) (*GListLIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LIndex not implemented")
}
func (UnimplementedGListServiceServer) LSet(context.Context, *GListLSetRequest) (*GListLSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LSet not implemented")
}
func (UnimplementedGListServiceServer) LTrim(context.Context, *GListLTrimRequest) (*GListLTrimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LTrim not implemented")
}
func (UnimplementedGListServiceServer) mustEmbedUnimplementedGListServiceServer() {}

// UnsafeGListServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GListServiceServer will
// result in compilation errors.
type UnsafeGListServiceServer interface {
	mustEmbedUnimplementedGListServiceServer()
}

func RegisterGListServiceServer(s grpc.ServiceRegistrar, srv GListServiceServer) {
	s.RegisterService(&GListService_ServiceDesc, srv)
}

func _GListService_LPush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLPushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LPush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LPush_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LPush(ctx, req.(*GListLPushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_LPushs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLPushsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LPushs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LPushs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LPushs(ctx, req.(*GListLPushsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_RPush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListRPushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).RPush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_RPush_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).RPush(ctx, req.(*GListRPushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_RPushs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListRPushsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).RPushs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_RPushs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).RPushs(ctx, req.(*GListRPushsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_LPop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLPopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LPop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LPop_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LPop(ctx, req.(*GListLPopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_RPop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListRPopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).RPop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_RPop_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).RPop(ctx, req.(*GListRPopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_LRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLRangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LRange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LRange(ctx, req.(*GListLRangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_LLen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLLenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LLen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LLen_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LLen(ctx, req.(*GListLLenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_LRem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLRemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LRem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LRem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LRem(ctx, req.(*GListLRemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_LIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LIndex_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LIndex(ctx, req.(*GListLIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_LSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LSet(ctx, req.(*GListLSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GListService_LTrim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GListLTrimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GListServiceServer).LTrim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GListService_LTrim_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GListServiceServer).LTrim(ctx, req.(*GListLTrimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GListService_ServiceDesc is the grpc.ServiceDesc for GListService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GListService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "glist.GListService",
	HandlerType: (*GListServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LPush",
			Handler:    _GListService_LPush_Handler,
		},
		{
			MethodName: "LPushs",
			Handler:    _GListService_LPushs_Handler,
		},
		{
			MethodName: "RPush",
			Handler:    _GListService_RPush_Handler,
		},
		{
			MethodName: "RPushs",
			Handler:    _GListService_RPushs_Handler,
		},
		{
			MethodName: "LPop",
			Handler:    _GListService_LPop_Handler,
		},
		{
			MethodName: "RPop",
			Handler:    _GListService_RPop_Handler,
		},
		{
			MethodName: "LRange",
			Handler:    _GListService_LRange_Handler,
		},
		{
			MethodName: "LLen",
			Handler:    _GListService_LLen_Handler,
		},
		{
			MethodName: "LRem",
			Handler:    _GListService_LRem_Handler,
		},
		{
			MethodName: "LIndex",
			Handler:    _GListService_LIndex_Handler,
		},
		{
			MethodName: "LSet",
			Handler:    _GListService_LSet_Handler,
		},
		{
			MethodName: "LTrim",
			Handler:    _GListService_LTrim_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lib/proto/glist/db.proto",
}
