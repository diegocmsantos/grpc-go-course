// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.17.3
// source: findmaximum/findmaximumpb/findmaximum.proto

package findmaximumpb

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

// FindMaximumServiceClient is the client API for FindMaximumService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FindMaximumServiceClient interface {
	FindMaximum(ctx context.Context, opts ...grpc.CallOption) (FindMaximumService_FindMaximumClient, error)
}

type findMaximumServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFindMaximumServiceClient(cc grpc.ClientConnInterface) FindMaximumServiceClient {
	return &findMaximumServiceClient{cc}
}

func (c *findMaximumServiceClient) FindMaximum(ctx context.Context, opts ...grpc.CallOption) (FindMaximumService_FindMaximumClient, error) {
	stream, err := c.cc.NewStream(ctx, &FindMaximumService_ServiceDesc.Streams[0], "/findmaximum.FindMaximumService/FindMaximum", opts...)
	if err != nil {
		return nil, err
	}
	x := &findMaximumServiceFindMaximumClient{stream}
	return x, nil
}

type FindMaximumService_FindMaximumClient interface {
	Send(*FindMaximumRequest) error
	Recv() (*FindMaximumResponse, error)
	grpc.ClientStream
}

type findMaximumServiceFindMaximumClient struct {
	grpc.ClientStream
}

func (x *findMaximumServiceFindMaximumClient) Send(m *FindMaximumRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *findMaximumServiceFindMaximumClient) Recv() (*FindMaximumResponse, error) {
	m := new(FindMaximumResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FindMaximumServiceServer is the server API for FindMaximumService service.
// All implementations must embed UnimplementedFindMaximumServiceServer
// for forward compatibility
type FindMaximumServiceServer interface {
	FindMaximum(FindMaximumService_FindMaximumServer) error
	mustEmbedUnimplementedFindMaximumServiceServer()
}

// UnimplementedFindMaximumServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFindMaximumServiceServer struct {
}

func (UnimplementedFindMaximumServiceServer) FindMaximum(FindMaximumService_FindMaximumServer) error {
	return status.Errorf(codes.Unimplemented, "method FindMaximum not implemented")
}
func (UnimplementedFindMaximumServiceServer) mustEmbedUnimplementedFindMaximumServiceServer() {}

// UnsafeFindMaximumServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FindMaximumServiceServer will
// result in compilation errors.
type UnsafeFindMaximumServiceServer interface {
	mustEmbedUnimplementedFindMaximumServiceServer()
}

func RegisterFindMaximumServiceServer(s grpc.ServiceRegistrar, srv FindMaximumServiceServer) {
	s.RegisterService(&FindMaximumService_ServiceDesc, srv)
}

func _FindMaximumService_FindMaximum_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FindMaximumServiceServer).FindMaximum(&findMaximumServiceFindMaximumServer{stream})
}

type FindMaximumService_FindMaximumServer interface {
	Send(*FindMaximumResponse) error
	Recv() (*FindMaximumRequest, error)
	grpc.ServerStream
}

type findMaximumServiceFindMaximumServer struct {
	grpc.ServerStream
}

func (x *findMaximumServiceFindMaximumServer) Send(m *FindMaximumResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *findMaximumServiceFindMaximumServer) Recv() (*FindMaximumRequest, error) {
	m := new(FindMaximumRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FindMaximumService_ServiceDesc is the grpc.ServiceDesc for FindMaximumService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FindMaximumService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "findmaximum.FindMaximumService",
	HandlerType: (*FindMaximumServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FindMaximum",
			Handler:       _FindMaximumService_FindMaximum_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "findmaximum/findmaximumpb/findmaximum.proto",
}
