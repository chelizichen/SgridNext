// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: NodeServer.proto

package protocol

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NodeServantClient is the client API for NodeServant service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeServantClient interface {
	// 节点心跳
	KeepAlive(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 获取节点服务状态
	GetNodeStat(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetNodeStatRsp, error)
	// 同步所有节点状态
	SyncAllNodeStat(ctx context.Context, in *SyncStatReq, opts ...grpc.CallOption) (*BasicRes, error)
	// 唤起
	ActivateServant(ctx context.Context, in *ActivateReq, opts ...grpc.CallOption) (*BasicRes, error)
	// 关闭
	DeactivateServant(ctx context.Context, in *ActivateReq, opts ...grpc.CallOption) (*BasicRes, error)
	// 告知同步配置文件(异步调用，结束后同步主节点)
	SyncConfigFile(ctx context.Context, in *SyncReq, opts ...grpc.CallOption) (*BasicRes, error)
	// 告知同步服务包(异步调用，结束后同步主节点)
	SyncServicePackage(ctx context.Context, in *SyncReq, opts ...grpc.CallOption) (*BasicRes, error)
	// Cgroup限制
	CgroupLimit(ctx context.Context, in *CgroupLimitReq, opts ...grpc.CallOption) (*BasicRes, error)
	// 检查状态
	CheckStat(ctx context.Context, in *CheckStatReq, opts ...grpc.CallOption) (*BasicRes, error)
	// 下载文件
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (NodeServant_DownloadFileClient, error)
	// 获取文件列表
	GetFileList(ctx context.Context, in *GetFileListReq, opts ...grpc.CallOption) (*GetFileListResponse, error)
	// 获取日志
	GetLog(ctx context.Context, in *GetLogReq, opts ...grpc.CallOption) (*GetLogRes, error)
}

type nodeServantClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeServantClient(cc grpc.ClientConnInterface) NodeServantClient {
	return &nodeServantClient{cc}
}

func (c *nodeServantClient) KeepAlive(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/KeepAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) GetNodeStat(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetNodeStatRsp, error) {
	out := new(GetNodeStatRsp)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/GetNodeStat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) SyncAllNodeStat(ctx context.Context, in *SyncStatReq, opts ...grpc.CallOption) (*BasicRes, error) {
	out := new(BasicRes)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/SyncAllNodeStat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) ActivateServant(ctx context.Context, in *ActivateReq, opts ...grpc.CallOption) (*BasicRes, error) {
	out := new(BasicRes)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/ActivateServant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) DeactivateServant(ctx context.Context, in *ActivateReq, opts ...grpc.CallOption) (*BasicRes, error) {
	out := new(BasicRes)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/DeactivateServant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) SyncConfigFile(ctx context.Context, in *SyncReq, opts ...grpc.CallOption) (*BasicRes, error) {
	out := new(BasicRes)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/SyncConfigFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) SyncServicePackage(ctx context.Context, in *SyncReq, opts ...grpc.CallOption) (*BasicRes, error) {
	out := new(BasicRes)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/SyncServicePackage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) CgroupLimit(ctx context.Context, in *CgroupLimitReq, opts ...grpc.CallOption) (*BasicRes, error) {
	out := new(BasicRes)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/CgroupLimit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) CheckStat(ctx context.Context, in *CheckStatReq, opts ...grpc.CallOption) (*BasicRes, error) {
	out := new(BasicRes)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/CheckStat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (NodeServant_DownloadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeServant_ServiceDesc.Streams[0], "/SgridProtocol.NodeServant/DownloadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeServantDownloadFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NodeServant_DownloadFileClient interface {
	Recv() (*DownloadFileResponse, error)
	grpc.ClientStream
}

type nodeServantDownloadFileClient struct {
	grpc.ClientStream
}

func (x *nodeServantDownloadFileClient) Recv() (*DownloadFileResponse, error) {
	m := new(DownloadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nodeServantClient) GetFileList(ctx context.Context, in *GetFileListReq, opts ...grpc.CallOption) (*GetFileListResponse, error) {
	out := new(GetFileListResponse)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/GetFileList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServantClient) GetLog(ctx context.Context, in *GetLogReq, opts ...grpc.CallOption) (*GetLogRes, error) {
	out := new(GetLogRes)
	err := c.cc.Invoke(ctx, "/SgridProtocol.NodeServant/GetLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServantServer is the server API for NodeServant service.
// All implementations must embed UnimplementedNodeServantServer
// for forward compatibility
type NodeServantServer interface {
	// 节点心跳
	KeepAlive(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// 获取节点服务状态
	GetNodeStat(context.Context, *emptypb.Empty) (*GetNodeStatRsp, error)
	// 同步所有节点状态
	SyncAllNodeStat(context.Context, *SyncStatReq) (*BasicRes, error)
	// 唤起
	ActivateServant(context.Context, *ActivateReq) (*BasicRes, error)
	// 关闭
	DeactivateServant(context.Context, *ActivateReq) (*BasicRes, error)
	// 告知同步配置文件(异步调用，结束后同步主节点)
	SyncConfigFile(context.Context, *SyncReq) (*BasicRes, error)
	// 告知同步服务包(异步调用，结束后同步主节点)
	SyncServicePackage(context.Context, *SyncReq) (*BasicRes, error)
	// Cgroup限制
	CgroupLimit(context.Context, *CgroupLimitReq) (*BasicRes, error)
	// 检查状态
	CheckStat(context.Context, *CheckStatReq) (*BasicRes, error)
	// 下载文件
	DownloadFile(*DownloadFileRequest, NodeServant_DownloadFileServer) error
	// 获取文件列表
	GetFileList(context.Context, *GetFileListReq) (*GetFileListResponse, error)
	// 获取日志
	GetLog(context.Context, *GetLogReq) (*GetLogRes, error)
	mustEmbedUnimplementedNodeServantServer()
}

// UnimplementedNodeServantServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServantServer struct {
}

func (UnimplementedNodeServantServer) KeepAlive(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KeepAlive not implemented")
}
func (UnimplementedNodeServantServer) GetNodeStat(context.Context, *emptypb.Empty) (*GetNodeStatRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeStat not implemented")
}
func (UnimplementedNodeServantServer) SyncAllNodeStat(context.Context, *SyncStatReq) (*BasicRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncAllNodeStat not implemented")
}
func (UnimplementedNodeServantServer) ActivateServant(context.Context, *ActivateReq) (*BasicRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateServant not implemented")
}
func (UnimplementedNodeServantServer) DeactivateServant(context.Context, *ActivateReq) (*BasicRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateServant not implemented")
}
func (UnimplementedNodeServantServer) SyncConfigFile(context.Context, *SyncReq) (*BasicRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncConfigFile not implemented")
}
func (UnimplementedNodeServantServer) SyncServicePackage(context.Context, *SyncReq) (*BasicRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncServicePackage not implemented")
}
func (UnimplementedNodeServantServer) CgroupLimit(context.Context, *CgroupLimitReq) (*BasicRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CgroupLimit not implemented")
}
func (UnimplementedNodeServantServer) CheckStat(context.Context, *CheckStatReq) (*BasicRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckStat not implemented")
}
func (UnimplementedNodeServantServer) DownloadFile(*DownloadFileRequest, NodeServant_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedNodeServantServer) GetFileList(context.Context, *GetFileListReq) (*GetFileListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileList not implemented")
}
func (UnimplementedNodeServantServer) GetLog(context.Context, *GetLogReq) (*GetLogRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLog not implemented")
}
func (UnimplementedNodeServantServer) mustEmbedUnimplementedNodeServantServer() {}

// UnsafeNodeServantServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServantServer will
// result in compilation errors.
type UnsafeNodeServantServer interface {
	mustEmbedUnimplementedNodeServantServer()
}

func RegisterNodeServantServer(s grpc.ServiceRegistrar, srv NodeServantServer) {
	s.RegisterService(&NodeServant_ServiceDesc, srv)
}

func _NodeServant_KeepAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).KeepAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/KeepAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).KeepAlive(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_GetNodeStat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).GetNodeStat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/GetNodeStat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).GetNodeStat(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_SyncAllNodeStat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncStatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).SyncAllNodeStat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/SyncAllNodeStat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).SyncAllNodeStat(ctx, req.(*SyncStatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_ActivateServant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).ActivateServant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/ActivateServant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).ActivateServant(ctx, req.(*ActivateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_DeactivateServant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).DeactivateServant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/DeactivateServant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).DeactivateServant(ctx, req.(*ActivateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_SyncConfigFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).SyncConfigFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/SyncConfigFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).SyncConfigFile(ctx, req.(*SyncReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_SyncServicePackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).SyncServicePackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/SyncServicePackage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).SyncServicePackage(ctx, req.(*SyncReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_CgroupLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CgroupLimitReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).CgroupLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/CgroupLimit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).CgroupLimit(ctx, req.(*CgroupLimitReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_CheckStat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckStatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).CheckStat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/CheckStat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).CheckStat(ctx, req.(*CheckStatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NodeServantServer).DownloadFile(m, &nodeServantDownloadFileServer{stream})
}

type NodeServant_DownloadFileServer interface {
	Send(*DownloadFileResponse) error
	grpc.ServerStream
}

type nodeServantDownloadFileServer struct {
	grpc.ServerStream
}

func (x *nodeServantDownloadFileServer) Send(m *DownloadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NodeServant_GetFileList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).GetFileList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/GetFileList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).GetFileList(ctx, req.(*GetFileListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeServant_GetLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServantServer).GetLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.NodeServant/GetLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServantServer).GetLog(ctx, req.(*GetLogReq))
	}
	return interceptor(ctx, in, info, handler)
}

// NodeServant_ServiceDesc is the grpc.ServiceDesc for NodeServant service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeServant_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SgridProtocol.NodeServant",
	HandlerType: (*NodeServantServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "KeepAlive",
			Handler:    _NodeServant_KeepAlive_Handler,
		},
		{
			MethodName: "GetNodeStat",
			Handler:    _NodeServant_GetNodeStat_Handler,
		},
		{
			MethodName: "SyncAllNodeStat",
			Handler:    _NodeServant_SyncAllNodeStat_Handler,
		},
		{
			MethodName: "ActivateServant",
			Handler:    _NodeServant_ActivateServant_Handler,
		},
		{
			MethodName: "DeactivateServant",
			Handler:    _NodeServant_DeactivateServant_Handler,
		},
		{
			MethodName: "SyncConfigFile",
			Handler:    _NodeServant_SyncConfigFile_Handler,
		},
		{
			MethodName: "SyncServicePackage",
			Handler:    _NodeServant_SyncServicePackage_Handler,
		},
		{
			MethodName: "CgroupLimit",
			Handler:    _NodeServant_CgroupLimit_Handler,
		},
		{
			MethodName: "CheckStat",
			Handler:    _NodeServant_CheckStat_Handler,
		},
		{
			MethodName: "GetFileList",
			Handler:    _NodeServant_GetFileList_Handler,
		},
		{
			MethodName: "GetLog",
			Handler:    _NodeServant_GetLog_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadFile",
			Handler:       _NodeServant_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "NodeServer.proto",
}
