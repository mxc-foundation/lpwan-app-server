// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ext_account.proto

package m2m_serves_appserver

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Money int32

const (
	Money_ETH  Money = 0
	Money_MXC  Money = 1
	Money_TETH Money = 2
)

var Money_name = map[int32]string{
	0: "ETH",
	1: "MXC",
	2: "TETH",
}

var Money_value = map[string]int32{
	"ETH":  0,
	"MXC":  1,
	"TETH": 2,
}

func (x Money) String() string {
	return proto.EnumName(Money_name, int32(x))
}

func (Money) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d6ddbf9312a3f483, []int{0}
}

type GetActiveMoneyAccountRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetActiveMoneyAccountRequest) Reset()         { *m = GetActiveMoneyAccountRequest{} }
func (m *GetActiveMoneyAccountRequest) String() string { return proto.CompactTextString(m) }
func (*GetActiveMoneyAccountRequest) ProtoMessage()    {}
func (*GetActiveMoneyAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d6ddbf9312a3f483, []int{0}
}

func (m *GetActiveMoneyAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetActiveMoneyAccountRequest.Unmarshal(m, b)
}
func (m *GetActiveMoneyAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetActiveMoneyAccountRequest.Marshal(b, m, deterministic)
}
func (m *GetActiveMoneyAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetActiveMoneyAccountRequest.Merge(m, src)
}
func (m *GetActiveMoneyAccountRequest) XXX_Size() int {
	return xxx_messageInfo_GetActiveMoneyAccountRequest.Size(m)
}
func (m *GetActiveMoneyAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetActiveMoneyAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetActiveMoneyAccountRequest proto.InternalMessageInfo

func (m *GetActiveMoneyAccountRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

type GetActiveMoneyAccountResponse struct {
	ActiveAccount        string   `protobuf:"bytes,1,opt,name=active_account,json=activeAccount,proto3" json:"active_account,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetActiveMoneyAccountResponse) Reset()         { *m = GetActiveMoneyAccountResponse{} }
func (m *GetActiveMoneyAccountResponse) String() string { return proto.CompactTextString(m) }
func (*GetActiveMoneyAccountResponse) ProtoMessage()    {}
func (*GetActiveMoneyAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d6ddbf9312a3f483, []int{1}
}

func (m *GetActiveMoneyAccountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetActiveMoneyAccountResponse.Unmarshal(m, b)
}
func (m *GetActiveMoneyAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetActiveMoneyAccountResponse.Marshal(b, m, deterministic)
}
func (m *GetActiveMoneyAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetActiveMoneyAccountResponse.Merge(m, src)
}
func (m *GetActiveMoneyAccountResponse) XXX_Size() int {
	return xxx_messageInfo_GetActiveMoneyAccountResponse.Size(m)
}
func (m *GetActiveMoneyAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetActiveMoneyAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetActiveMoneyAccountResponse proto.InternalMessageInfo

func (m *GetActiveMoneyAccountResponse) GetActiveAccount() string {
	if m != nil {
		return m.ActiveAccount
	}
	return ""
}

func init() {
	proto.RegisterEnum("m2m_serves_appserver.Money", Money_name, Money_value)
	proto.RegisterType((*GetActiveMoneyAccountRequest)(nil), "m2m_serves_appserver.GetActiveMoneyAccountRequest")
	proto.RegisterType((*GetActiveMoneyAccountResponse)(nil), "m2m_serves_appserver.GetActiveMoneyAccountResponse")
}

func init() { proto.RegisterFile("ext_account.proto", fileDescriptor_d6ddbf9312a3f483) }

var fileDescriptor_d6ddbf9312a3f483 = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xad, 0x28, 0x89,
	0x4f, 0x4c, 0x4e, 0xce, 0x2f, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0xc9,
	0x35, 0xca, 0x8d, 0x2f, 0x4e, 0x2d, 0x2a, 0x4b, 0x2d, 0x8e, 0x4f, 0x2c, 0x28, 0x00, 0xb3, 0x8a,
	0x94, 0x4c, 0xb9, 0x64, 0xdc, 0x53, 0x4b, 0x1c, 0x93, 0x4b, 0x32, 0xcb, 0x52, 0x7d, 0xf3, 0xf3,
	0x52, 0x2b, 0x1d, 0x21, 0x9a, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x44, 0xb9, 0xd8,
	0xf2, 0x8b, 0xd2, 0xe3, 0x33, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x98, 0x83, 0x58, 0xf3, 0x8b,
	0xd2, 0x3d, 0x53, 0x94, 0xdc, 0xb8, 0x64, 0x71, 0x68, 0x2b, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x15,
	0x52, 0xe5, 0xe2, 0x4b, 0x04, 0xcb, 0xc2, 0x5c, 0x01, 0xd6, 0xcf, 0x19, 0xc4, 0x0b, 0x11, 0x85,
	0x2a, 0xd7, 0x52, 0xe6, 0x62, 0x05, 0x6b, 0x17, 0x62, 0xe7, 0x62, 0x76, 0x0d, 0xf1, 0x10, 0x60,
	0x00, 0x31, 0x7c, 0x23, 0x9c, 0x05, 0x18, 0x85, 0x38, 0xb8, 0x58, 0x42, 0x40, 0x42, 0x4c, 0x46,
	0x13, 0x19, 0xb9, 0x78, 0xc0, 0xaa, 0x82, 0x53, 0x8b, 0xca, 0x32, 0x93, 0x53, 0x85, 0x1a, 0x18,
	0xb9, 0x44, 0xb1, 0x5a, 0x2f, 0x64, 0xa4, 0x87, 0xcd, 0x97, 0x7a, 0xf8, 0xbc, 0x28, 0x65, 0x4c,
	0x92, 0x1e, 0x88, 0xff, 0x92, 0xd8, 0xc0, 0x81, 0x6a, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xcd,
	0xe8, 0xfd, 0x58, 0x69, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MoneyServiceClient is the client API for MoneyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MoneyServiceClient interface {
	GetActiveMoneyAccount(ctx context.Context, in *GetActiveMoneyAccountRequest, opts ...grpc.CallOption) (*GetActiveMoneyAccountResponse, error)
}

type moneyServiceClient struct {
	cc *grpc.ClientConn
}

func NewMoneyServiceClient(cc *grpc.ClientConn) MoneyServiceClient {
	return &moneyServiceClient{cc}
}

func (c *moneyServiceClient) GetActiveMoneyAccount(ctx context.Context, in *GetActiveMoneyAccountRequest, opts ...grpc.CallOption) (*GetActiveMoneyAccountResponse, error) {
	out := new(GetActiveMoneyAccountResponse)
	err := c.cc.Invoke(ctx, "/m2m_serves_appserver.MoneyService/GetActiveMoneyAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MoneyServiceServer is the server API for MoneyService service.
type MoneyServiceServer interface {
	GetActiveMoneyAccount(context.Context, *GetActiveMoneyAccountRequest) (*GetActiveMoneyAccountResponse, error)
}

// UnimplementedMoneyServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMoneyServiceServer struct {
}

func (*UnimplementedMoneyServiceServer) GetActiveMoneyAccount(ctx context.Context, req *GetActiveMoneyAccountRequest) (*GetActiveMoneyAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActiveMoneyAccount not implemented")
}

func RegisterMoneyServiceServer(s *grpc.Server, srv MoneyServiceServer) {
	s.RegisterService(&_MoneyService_serviceDesc, srv)
}

func _MoneyService_GetActiveMoneyAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetActiveMoneyAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoneyServiceServer).GetActiveMoneyAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_serves_appserver.MoneyService/GetActiveMoneyAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoneyServiceServer).GetActiveMoneyAccount(ctx, req.(*GetActiveMoneyAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MoneyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m2m_serves_appserver.MoneyService",
	HandlerType: (*MoneyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetActiveMoneyAccount",
			Handler:    _MoneyService_GetActiveMoneyAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ext_account.proto",
}
