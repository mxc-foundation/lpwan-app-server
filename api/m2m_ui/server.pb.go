// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server.proto

package m2m_ui

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetVersionResponse struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetVersionResponse) Reset()         { *m = GetVersionResponse{} }
func (m *GetVersionResponse) String() string { return proto.CompactTextString(m) }
func (*GetVersionResponse) ProtoMessage()    {}
func (*GetVersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{0}
}

func (m *GetVersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetVersionResponse.Unmarshal(m, b)
}
func (m *GetVersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetVersionResponse.Marshal(b, m, deterministic)
}
func (m *GetVersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetVersionResponse.Merge(m, src)
}
func (m *GetVersionResponse) XXX_Size() int {
	return xxx_messageInfo_GetVersionResponse.Size(m)
}
func (m *GetVersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetVersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetVersionResponse proto.InternalMessageInfo

func (m *GetVersionResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*GetVersionResponse)(nil), "m2m_ui.GetVersionResponse")
}

func init() { proto.RegisterFile("server.proto", fileDescriptor_ad098daeda4239f7) }

var fileDescriptor_ad098daeda4239f7 = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcb, 0x35, 0xca, 0x8d, 0x2f, 0xcd,
	0x94, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xcc, 0xcb,
	0xcb, 0x2f, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0x86, 0xa8, 0x92, 0x92, 0x86, 0xca, 0x82, 0x79,
	0x49, 0xa5, 0x69, 0xfa, 0xa9, 0xb9, 0x05, 0x25, 0x95, 0x10, 0x49, 0x25, 0x3d, 0x2e, 0x21, 0xf7,
	0xd4, 0x92, 0xb0, 0xd4, 0xa2, 0xe2, 0xcc, 0xfc, 0xbc, 0xa0, 0xd4, 0xe2, 0x82, 0xfc, 0xbc, 0xe2,
	0x54, 0x21, 0x09, 0x2e, 0xf6, 0x32, 0x88, 0x90, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x8c,
	0x6b, 0x54, 0xce, 0x25, 0x18, 0x0c, 0x76, 0x82, 0x67, 0x5e, 0x5a, 0x3e, 0x88, 0x95, 0x99, 0x9c,
	0x2a, 0x94, 0xc4, 0xc5, 0x85, 0x30, 0x44, 0x48, 0x4c, 0x0f, 0x62, 0xa1, 0x1e, 0xcc, 0x42, 0x3d,
	0x57, 0x90, 0x85, 0x52, 0x52, 0x7a, 0x10, 0xe7, 0xea, 0x61, 0x5a, 0xa8, 0xa4, 0xd0, 0x74, 0xf9,
	0xc9, 0x64, 0x26, 0x29, 0x21, 0x09, 0xb0, 0x1f, 0x20, 0x9e, 0xd4, 0xcd, 0xcc, 0x4b, 0xcb, 0xd7,
	0x87, 0x5a, 0x9c, 0xc4, 0x06, 0x36, 0xcd, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x7a, 0x6e, 0x88,
	0x3a, 0x02, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServerInfoServiceClient is the client API for ServerInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServerInfoServiceClient interface {
	// get version
	GetVersion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetVersionResponse, error)
}

type serverInfoServiceClient struct {
	cc *grpc.ClientConn
}

func NewServerInfoServiceClient(cc *grpc.ClientConn) ServerInfoServiceClient {
	return &serverInfoServiceClient{cc}
}

func (c *serverInfoServiceClient) GetVersion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetVersionResponse, error) {
	out := new(GetVersionResponse)
	err := c.cc.Invoke(ctx, "/m2m_ui.ServerInfoService/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerInfoServiceServer is the server API for ServerInfoService service.
type ServerInfoServiceServer interface {
	// get version
	GetVersion(context.Context, *empty.Empty) (*GetVersionResponse, error)
}

// UnimplementedServerInfoServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServerInfoServiceServer struct {
}

func (*UnimplementedServerInfoServiceServer) GetVersion(ctx context.Context, req *empty.Empty) (*GetVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}

func RegisterServerInfoServiceServer(s *grpc.Server, srv ServerInfoServiceServer) {
	s.RegisterService(&_ServerInfoService_serviceDesc, srv)
}

func _ServerInfoService_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerInfoServiceServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_ui.ServerInfoService/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerInfoServiceServer).GetVersion(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServerInfoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m2m_ui.ServerInfoService",
	HandlerType: (*ServerInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _ServerInfoService_GetVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
