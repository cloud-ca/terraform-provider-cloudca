// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc_broker.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ConnInfo struct {
	ServiceId            uint32   `protobuf:"varint,1,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	Network              string   `protobuf:"bytes,2,opt,name=network,proto3" json:"network,omitempty"`
	Address              string   `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnInfo) Reset()         { *m = ConnInfo{} }
func (m *ConnInfo) String() string { return proto.CompactTextString(m) }
func (*ConnInfo) ProtoMessage()    {}
func (*ConnInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_802e9beed3ec3b28, []int{0}
}

func (m *ConnInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnInfo.Unmarshal(m, b)
}
func (m *ConnInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnInfo.Marshal(b, m, deterministic)
}
func (m *ConnInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnInfo.Merge(m, src)
}
func (m *ConnInfo) XXX_Size() int {
	return xxx_messageInfo_ConnInfo.Size(m)
}
func (m *ConnInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ConnInfo proto.InternalMessageInfo

func (m *ConnInfo) GetServiceId() uint32 {
	if m != nil {
		return m.ServiceId
	}
	return 0
}

func (m *ConnInfo) GetNetwork() string {
	if m != nil {
		return m.Network
	}
	return ""
}

func (m *ConnInfo) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func init() {
	proto.RegisterType((*ConnInfo)(nil), "proto.ConnInfo")
}

func init() { proto.RegisterFile("grpc_broker.proto", fileDescriptor_802e9beed3ec3b28) }

var fileDescriptor_802e9beed3ec3b28 = []byte{
	// 164 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0x2f, 0x2a, 0x48,
	0x8e, 0x4f, 0x2a, 0xca, 0xcf, 0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05,
	0x53, 0x4a, 0xb1, 0x5c, 0x1c, 0xce, 0xf9, 0x79, 0x79, 0x9e, 0x79, 0x69, 0xf9, 0x42, 0xb2, 0x5c,
	0x5c, 0xc5, 0xa9, 0x45, 0x65, 0x99, 0xc9, 0xa9, 0xf1, 0x99, 0x29, 0x12, 0x8c, 0x0a, 0x8c, 0x1a,
	0xbc, 0x41, 0x9c, 0x50, 0x11, 0xcf, 0x14, 0x21, 0x09, 0x2e, 0xf6, 0xbc, 0xd4, 0x92, 0xf2, 0xfc,
	0xa2, 0x6c, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x18, 0x17, 0x24, 0x93, 0x98, 0x92, 0x52,
	0x94, 0x5a, 0x5c, 0x2c, 0xc1, 0x0c, 0x91, 0x81, 0x72, 0x8d, 0x1c, 0xb9, 0xb8, 0xdc, 0x83, 0x02,
	0x9c, 0x9d, 0xc0, 0x36, 0x0b, 0x19, 0x73, 0x71, 0x07, 0x97, 0x24, 0x16, 0x95, 0x04, 0x97, 0x14,
	0xa5, 0x26, 0xe6, 0x0a, 0xf1, 0x43, 0x9c, 0xa2, 0x07, 0x73, 0x80, 0x14, 0xba, 0x80, 0x06, 0xa3,
	0x01, 0x63, 0x12, 0x1b, 0x58, 0xcc, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x7a, 0xda, 0xd5, 0x84,
	0xc4, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GRPCBrokerClient is the client API for GRPCBroker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GRPCBrokerClient interface {
	StartStream(ctx context.Context, opts ...grpc.CallOption) (GRPCBroker_StartStreamClient, error)
}

type gRPCBrokerClient struct {
	cc *grpc.ClientConn
}

func NewGRPCBrokerClient(cc *grpc.ClientConn) GRPCBrokerClient {
	return &gRPCBrokerClient{cc}
}

func (c *gRPCBrokerClient) StartStream(ctx context.Context, opts ...grpc.CallOption) (GRPCBroker_StartStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GRPCBroker_serviceDesc.Streams[0], "/proto.GRPCBroker/StartStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCBrokerStartStreamClient{stream}
	return x, nil
}

type GRPCBroker_StartStreamClient interface {
	Send(*ConnInfo) error
	Recv() (*ConnInfo, error)
	grpc.ClientStream
}

type gRPCBrokerStartStreamClient struct {
	grpc.ClientStream
}

func (x *gRPCBrokerStartStreamClient) Send(m *ConnInfo) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gRPCBrokerStartStreamClient) Recv() (*ConnInfo, error) {
	m := new(ConnInfo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GRPCBrokerServer is the server API for GRPCBroker service.
type GRPCBrokerServer interface {
	StartStream(GRPCBroker_StartStreamServer) error
}

func RegisterGRPCBrokerServer(s *grpc.Server, srv GRPCBrokerServer) {
	s.RegisterService(&_GRPCBroker_serviceDesc, srv)
}

func _GRPCBroker_StartStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GRPCBrokerServer).StartStream(&gRPCBrokerStartStreamServer{stream})
}

type GRPCBroker_StartStreamServer interface {
	Send(*ConnInfo) error
	Recv() (*ConnInfo, error)
	grpc.ServerStream
}

type gRPCBrokerStartStreamServer struct {
	grpc.ServerStream
}

func (x *gRPCBrokerStartStreamServer) Send(m *ConnInfo) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gRPCBrokerStartStreamServer) Recv() (*ConnInfo, error) {
	m := new(ConnInfo)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _GRPCBroker_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GRPCBroker",
	HandlerType: (*GRPCBrokerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StartStream",
			Handler:       _GRPCBroker_StartStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpc_broker.proto",
}
