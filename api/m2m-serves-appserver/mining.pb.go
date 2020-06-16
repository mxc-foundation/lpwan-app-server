// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mining.proto

package m2m_serves_appserver

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// MiningRequest sends gateway list to m2m
type MiningRequest struct {
	GatewayMac           []string `protobuf:"bytes,1,rep,name=gateway_mac,json=gatewayMac,proto3" json:"gateway_mac,omitempty"`
	MiningRevenue        float64  `protobuf:"fixed64,2,opt,name=mining_revenue,json=miningRevenue,proto3" json:"mining_revenue,omitempty"`
	MxcPrice             float64  `protobuf:"fixed64,3,opt,name=mxc_price,json=mxcPrice,proto3" json:"mxc_price,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MiningRequest) Reset()         { *m = MiningRequest{} }
func (m *MiningRequest) String() string { return proto.CompactTextString(m) }
func (*MiningRequest) ProtoMessage()    {}
func (*MiningRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee9e6e83dd861d31, []int{0}
}

func (m *MiningRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MiningRequest.Unmarshal(m, b)
}
func (m *MiningRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MiningRequest.Marshal(b, m, deterministic)
}
func (m *MiningRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MiningRequest.Merge(m, src)
}
func (m *MiningRequest) XXX_Size() int {
	return xxx_messageInfo_MiningRequest.Size(m)
}
func (m *MiningRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MiningRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MiningRequest proto.InternalMessageInfo

func (m *MiningRequest) GetGatewayMac() []string {
	if m != nil {
		return m.GatewayMac
	}
	return nil
}

func (m *MiningRequest) GetMiningRevenue() float64 {
	if m != nil {
		return m.MiningRevenue
	}
	return 0
}

func (m *MiningRequest) GetMxcPrice() float64 {
	if m != nil {
		return m.MxcPrice
	}
	return 0
}

type MiningResponse struct {
	Status               bool     `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MiningResponse) Reset()         { *m = MiningResponse{} }
func (m *MiningResponse) String() string { return proto.CompactTextString(m) }
func (*MiningResponse) ProtoMessage()    {}
func (*MiningResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee9e6e83dd861d31, []int{1}
}

func (m *MiningResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MiningResponse.Unmarshal(m, b)
}
func (m *MiningResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MiningResponse.Marshal(b, m, deterministic)
}
func (m *MiningResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MiningResponse.Merge(m, src)
}
func (m *MiningResponse) XXX_Size() int {
	return xxx_messageInfo_MiningResponse.Size(m)
}
func (m *MiningResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MiningResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MiningResponse proto.InternalMessageInfo

func (m *MiningResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func init() {
	proto.RegisterType((*MiningRequest)(nil), "m2m_serves_appserver.MiningRequest")
	proto.RegisterType((*MiningResponse)(nil), "m2m_serves_appserver.MiningResponse")
}

func init() { proto.RegisterFile("mining.proto", fileDescriptor_ee9e6e83dd861d31) }

var fileDescriptor_ee9e6e83dd861d31 = []byte{
	// 213 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xc1, 0x4a, 0xc4, 0x30,
	0x18, 0x84, 0x89, 0x0b, 0x65, 0xf7, 0xd7, 0xed, 0x21, 0x88, 0x04, 0x3d, 0x58, 0xaa, 0x42, 0x4e,
	0x3d, 0xd4, 0xe7, 0x28, 0x48, 0xfa, 0x00, 0x21, 0xc6, 0x9f, 0x92, 0x43, 0xd2, 0x98, 0xa4, 0xb5,
	0xbe, 0xbd, 0x98, 0xb6, 0x07, 0x41, 0xf6, 0x96, 0x19, 0x3e, 0x66, 0x32, 0x3f, 0xdc, 0x58, 0xe3,
	0x8c, 0x1b, 0x1a, 0x1f, 0xc6, 0x34, 0xd2, 0x5b, 0xdb, 0x5a, 0x19, 0x31, 0xcc, 0x18, 0xa5, 0xf2,
	0x3e, 0xbf, 0x42, 0x9d, 0xe0, 0xdc, 0x65, 0x4a, 0xe0, 0xe7, 0x84, 0x31, 0xd1, 0x47, 0xb8, 0x1e,
	0x54, 0xc2, 0x2f, 0xf5, 0x2d, 0xad, 0xd2, 0x8c, 0x54, 0x07, 0x7e, 0x12, 0xb0, 0x59, 0x9d, 0xd2,
	0xf4, 0x05, 0xca, 0x35, 0x57, 0x06, 0x9c, 0xd1, 0x4d, 0xc8, 0xae, 0x2a, 0xc2, 0x89, 0x38, 0xdb,
	0x2d, 0x27, 0x9b, 0xf4, 0x01, 0x4e, 0x76, 0xd1, 0xd2, 0x07, 0xa3, 0x91, 0x1d, 0x32, 0x71, 0xb4,
	0x8b, 0x7e, 0xfb, 0xd5, 0x35, 0x87, 0x72, 0x6f, 0x8d, 0x7e, 0x74, 0x11, 0xe9, 0x1d, 0x14, 0x31,
	0xa9, 0x34, 0x45, 0x46, 0x2a, 0xc2, 0x8f, 0x62, 0x53, 0xed, 0xc7, 0xfe, 0xbf, 0x1e, 0xc3, 0x6c,
	0x34, 0xd2, 0x1e, 0x8a, 0xd5, 0xa0, 0x4f, 0xcd, 0x7f, 0x8b, 0x9a, 0x3f, 0x73, 0xee, 0x9f, 0x2f,
	0x43, 0x6b, 0xfb, 0x7b, 0x91, 0x4f, 0xf4, 0xfa, 0x13, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xfb, 0x4b,
	0x91, 0x32, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MiningServiceClient is the client API for MiningService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MiningServiceClient interface {
	Mining(ctx context.Context, in *MiningRequest, opts ...grpc.CallOption) (*MiningResponse, error)
}

type miningServiceClient struct {
	cc *grpc.ClientConn
}

func NewMiningServiceClient(cc *grpc.ClientConn) MiningServiceClient {
	return &miningServiceClient{cc}
}

func (c *miningServiceClient) Mining(ctx context.Context, in *MiningRequest, opts ...grpc.CallOption) (*MiningResponse, error) {
	out := new(MiningResponse)
	err := c.cc.Invoke(ctx, "/m2m_serves_appserver.MiningService/Mining", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MiningServiceServer is the server API for MiningService service.
type MiningServiceServer interface {
	Mining(context.Context, *MiningRequest) (*MiningResponse, error)
}

// UnimplementedMiningServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMiningServiceServer struct {
}

func (*UnimplementedMiningServiceServer) Mining(ctx context.Context, req *MiningRequest) (*MiningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mining not implemented")
}

func RegisterMiningServiceServer(s *grpc.Server, srv MiningServiceServer) {
	s.RegisterService(&_MiningService_serviceDesc, srv)
}

func _MiningService_Mining_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MiningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiningServiceServer).Mining(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_serves_appserver.MiningService/Mining",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiningServiceServer).Mining(ctx, req.(*MiningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MiningService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m2m_serves_appserver.MiningService",
	HandlerType: (*MiningServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Mining",
			Handler:    _MiningService_Mining_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mining.proto",
}
