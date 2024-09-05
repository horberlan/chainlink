// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: csa/v1/csa.proto

package v1

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
	CSAService_GetKeypair_FullMethodName   = "/api.csa.v1.CSAService/GetKeypair"
	CSAService_ListKeypairs_FullMethodName = "/api.csa.v1.CSAService/ListKeypairs"
)

// CSAServiceClient is the client API for CSAService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CSAServiceClient interface {
	GetKeypair(ctx context.Context, in *GetKeypairRequest, opts ...grpc.CallOption) (*GetKeypairResponse, error)
	ListKeypairs(ctx context.Context, in *ListKeypairsRequest, opts ...grpc.CallOption) (*ListKeypairsResponse, error)
}

type cSAServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCSAServiceClient(cc grpc.ClientConnInterface) CSAServiceClient {
	return &cSAServiceClient{cc}
}

func (c *cSAServiceClient) GetKeypair(ctx context.Context, in *GetKeypairRequest, opts ...grpc.CallOption) (*GetKeypairResponse, error) {
	out := new(GetKeypairResponse)
	err := c.cc.Invoke(ctx, CSAService_GetKeypair_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cSAServiceClient) ListKeypairs(ctx context.Context, in *ListKeypairsRequest, opts ...grpc.CallOption) (*ListKeypairsResponse, error) {
	out := new(ListKeypairsResponse)
	err := c.cc.Invoke(ctx, CSAService_ListKeypairs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CSAServiceServer is the server API for CSAService service.
// All implementations must embed UnimplementedCSAServiceServer
// for forward compatibility
type CSAServiceServer interface {
	GetKeypair(context.Context, *GetKeypairRequest) (*GetKeypairResponse, error)
	ListKeypairs(context.Context, *ListKeypairsRequest) (*ListKeypairsResponse, error)
	mustEmbedUnimplementedCSAServiceServer()
}

// UnimplementedCSAServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCSAServiceServer struct {
}

func (UnimplementedCSAServiceServer) GetKeypair(context.Context, *GetKeypairRequest) (*GetKeypairResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeypair not implemented")
}
func (UnimplementedCSAServiceServer) ListKeypairs(context.Context, *ListKeypairsRequest) (*ListKeypairsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListKeypairs not implemented")
}
func (UnimplementedCSAServiceServer) mustEmbedUnimplementedCSAServiceServer() {}

// UnsafeCSAServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CSAServiceServer will
// result in compilation errors.
type UnsafeCSAServiceServer interface {
	mustEmbedUnimplementedCSAServiceServer()
}

func RegisterCSAServiceServer(s grpc.ServiceRegistrar, srv CSAServiceServer) {
	s.RegisterService(&CSAService_ServiceDesc, srv)
}

func _CSAService_GetKeypair_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeypairRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CSAServiceServer).GetKeypair(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CSAService_GetKeypair_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CSAServiceServer).GetKeypair(ctx, req.(*GetKeypairRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CSAService_ListKeypairs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListKeypairsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CSAServiceServer).ListKeypairs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CSAService_ListKeypairs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CSAServiceServer).ListKeypairs(ctx, req.(*ListKeypairsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CSAService_ServiceDesc is the grpc.ServiceDesc for CSAService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CSAService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.csa.v1.CSAService",
	HandlerType: (*CSAServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetKeypair",
			Handler:    _CSAService_GetKeypair_Handler,
		},
		{
			MethodName: "ListKeypairs",
			Handler:    _CSAService_ListKeypairs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "csa/v1/csa.proto",
}