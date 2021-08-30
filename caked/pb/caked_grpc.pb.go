// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// CakedClient is the client API for Caked service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CakedClient interface {
	RunSlice(ctx context.Context, in *Slice, opts ...grpc.CallOption) (*SliceStatus, error)
	StopSlice(ctx context.Context, in *Slice, opts ...grpc.CallOption) (*SliceStatus, error)
}

type cakedClient struct {
	cc grpc.ClientConnInterface
}

func NewCakedClient(cc grpc.ClientConnInterface) CakedClient {
	return &cakedClient{cc}
}

func (c *cakedClient) RunSlice(ctx context.Context, in *Slice, opts ...grpc.CallOption) (*SliceStatus, error) {
	out := new(SliceStatus)
	err := c.cc.Invoke(ctx, "/cake.Caked/RunSlice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cakedClient) StopSlice(ctx context.Context, in *Slice, opts ...grpc.CallOption) (*SliceStatus, error) {
	out := new(SliceStatus)
	err := c.cc.Invoke(ctx, "/cake.Caked/StopSlice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CakedServer is the server API for Caked service.
// All implementations must embed UnimplementedCakedServer
// for forward compatibility
type CakedServer interface {
	RunSlice(context.Context, *Slice) (*SliceStatus, error)
	StopSlice(context.Context, *Slice) (*SliceStatus, error)
	mustEmbedUnimplementedCakedServer()
}

// UnimplementedCakedServer must be embedded to have forward compatible implementations.
type UnimplementedCakedServer struct {
}

func (UnimplementedCakedServer) RunSlice(context.Context, *Slice) (*SliceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunSlice not implemented")
}
func (UnimplementedCakedServer) StopSlice(context.Context, *Slice) (*SliceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopSlice not implemented")
}
func (UnimplementedCakedServer) mustEmbedUnimplementedCakedServer() {}

// UnsafeCakedServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CakedServer will
// result in compilation errors.
type UnsafeCakedServer interface {
	mustEmbedUnimplementedCakedServer()
}

func RegisterCakedServer(s grpc.ServiceRegistrar, srv CakedServer) {
	s.RegisterService(&Caked_ServiceDesc, srv)
}

func _Caked_RunSlice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Slice)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakedServer).RunSlice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cake.Caked/RunSlice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakedServer).RunSlice(ctx, req.(*Slice))
	}
	return interceptor(ctx, in, info, handler)
}

func _Caked_StopSlice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Slice)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakedServer).StopSlice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cake.Caked/StopSlice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakedServer).StopSlice(ctx, req.(*Slice))
	}
	return interceptor(ctx, in, info, handler)
}

// Caked_ServiceDesc is the grpc.ServiceDesc for Caked service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Caked_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cake.Caked",
	HandlerType: (*CakedServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunSlice",
			Handler:    _Caked_RunSlice_Handler,
		},
		{
			MethodName: "StopSlice",
			Handler:    _Caked_StopSlice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "caked.proto",
}
