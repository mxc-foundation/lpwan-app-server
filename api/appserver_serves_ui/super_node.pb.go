// Code generated by protoc-gen-go. DO NOT EDIT.
// source: super_node.proto

package appserver_serves_ui

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type GetSuperNodeActiveMoneyAccountRequest struct {
	MoneyAbbr            Money    `protobuf:"varint,1,opt,name=money_abbr,json=moneyAbbr,proto3,enum=appserver_serves_ui.Money" json:"money_abbr,omitempty"`
	OrgId                int64    `protobuf:"varint,2,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSuperNodeActiveMoneyAccountRequest) Reset()         { *m = GetSuperNodeActiveMoneyAccountRequest{} }
func (m *GetSuperNodeActiveMoneyAccountRequest) String() string { return proto.CompactTextString(m) }
func (*GetSuperNodeActiveMoneyAccountRequest) ProtoMessage()    {}
func (*GetSuperNodeActiveMoneyAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_02e142dc5bc4ebd3, []int{0}
}

func (m *GetSuperNodeActiveMoneyAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSuperNodeActiveMoneyAccountRequest.Unmarshal(m, b)
}
func (m *GetSuperNodeActiveMoneyAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSuperNodeActiveMoneyAccountRequest.Marshal(b, m, deterministic)
}
func (m *GetSuperNodeActiveMoneyAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSuperNodeActiveMoneyAccountRequest.Merge(m, src)
}
func (m *GetSuperNodeActiveMoneyAccountRequest) XXX_Size() int {
	return xxx_messageInfo_GetSuperNodeActiveMoneyAccountRequest.Size(m)
}
func (m *GetSuperNodeActiveMoneyAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSuperNodeActiveMoneyAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSuperNodeActiveMoneyAccountRequest proto.InternalMessageInfo

func (m *GetSuperNodeActiveMoneyAccountRequest) GetMoneyAbbr() Money {
	if m != nil {
		return m.MoneyAbbr
	}
	return Money_ETH
}

func (m *GetSuperNodeActiveMoneyAccountRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

type GetSuperNodeActiveMoneyAccountResponse struct {
	SupernodeActiveAccount string           `protobuf:"bytes,1,opt,name=supernode_active_account,json=supernodeActiveAccount,proto3" json:"supernode_active_account,omitempty"`
	UserProfile            *ProfileResponse `protobuf:"bytes,2,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral   struct{}         `json:"-"`
	XXX_unrecognized       []byte           `json:"-"`
	XXX_sizecache          int32            `json:"-"`
}

func (m *GetSuperNodeActiveMoneyAccountResponse) Reset() {
	*m = GetSuperNodeActiveMoneyAccountResponse{}
}
func (m *GetSuperNodeActiveMoneyAccountResponse) String() string { return proto.CompactTextString(m) }
func (*GetSuperNodeActiveMoneyAccountResponse) ProtoMessage()    {}
func (*GetSuperNodeActiveMoneyAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_02e142dc5bc4ebd3, []int{1}
}

func (m *GetSuperNodeActiveMoneyAccountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSuperNodeActiveMoneyAccountResponse.Unmarshal(m, b)
}
func (m *GetSuperNodeActiveMoneyAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSuperNodeActiveMoneyAccountResponse.Marshal(b, m, deterministic)
}
func (m *GetSuperNodeActiveMoneyAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSuperNodeActiveMoneyAccountResponse.Merge(m, src)
}
func (m *GetSuperNodeActiveMoneyAccountResponse) XXX_Size() int {
	return xxx_messageInfo_GetSuperNodeActiveMoneyAccountResponse.Size(m)
}
func (m *GetSuperNodeActiveMoneyAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSuperNodeActiveMoneyAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSuperNodeActiveMoneyAccountResponse proto.InternalMessageInfo

func (m *GetSuperNodeActiveMoneyAccountResponse) GetSupernodeActiveAccount() string {
	if m != nil {
		return m.SupernodeActiveAccount
	}
	return ""
}

func (m *GetSuperNodeActiveMoneyAccountResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

type AddSuperNodeMoneyAccountRequest struct {
	MoneyAbbr            Money    `protobuf:"varint,1,opt,name=money_abbr,json=moneyAbbr,proto3,enum=appserver_serves_ui.Money" json:"money_abbr,omitempty"`
	AccountAddr          string   `protobuf:"bytes,2,opt,name=account_addr,json=accountAddr,proto3" json:"account_addr,omitempty"`
	OrgId                int64    `protobuf:"varint,3,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddSuperNodeMoneyAccountRequest) Reset()         { *m = AddSuperNodeMoneyAccountRequest{} }
func (m *AddSuperNodeMoneyAccountRequest) String() string { return proto.CompactTextString(m) }
func (*AddSuperNodeMoneyAccountRequest) ProtoMessage()    {}
func (*AddSuperNodeMoneyAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_02e142dc5bc4ebd3, []int{2}
}

func (m *AddSuperNodeMoneyAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddSuperNodeMoneyAccountRequest.Unmarshal(m, b)
}
func (m *AddSuperNodeMoneyAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddSuperNodeMoneyAccountRequest.Marshal(b, m, deterministic)
}
func (m *AddSuperNodeMoneyAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddSuperNodeMoneyAccountRequest.Merge(m, src)
}
func (m *AddSuperNodeMoneyAccountRequest) XXX_Size() int {
	return xxx_messageInfo_AddSuperNodeMoneyAccountRequest.Size(m)
}
func (m *AddSuperNodeMoneyAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddSuperNodeMoneyAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddSuperNodeMoneyAccountRequest proto.InternalMessageInfo

func (m *AddSuperNodeMoneyAccountRequest) GetMoneyAbbr() Money {
	if m != nil {
		return m.MoneyAbbr
	}
	return Money_ETH
}

func (m *AddSuperNodeMoneyAccountRequest) GetAccountAddr() string {
	if m != nil {
		return m.AccountAddr
	}
	return ""
}

func (m *AddSuperNodeMoneyAccountRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

type AddSuperNodeMoneyAccountResponse struct {
	Status               bool             `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,2,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *AddSuperNodeMoneyAccountResponse) Reset()         { *m = AddSuperNodeMoneyAccountResponse{} }
func (m *AddSuperNodeMoneyAccountResponse) String() string { return proto.CompactTextString(m) }
func (*AddSuperNodeMoneyAccountResponse) ProtoMessage()    {}
func (*AddSuperNodeMoneyAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_02e142dc5bc4ebd3, []int{3}
}

func (m *AddSuperNodeMoneyAccountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddSuperNodeMoneyAccountResponse.Unmarshal(m, b)
}
func (m *AddSuperNodeMoneyAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddSuperNodeMoneyAccountResponse.Marshal(b, m, deterministic)
}
func (m *AddSuperNodeMoneyAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddSuperNodeMoneyAccountResponse.Merge(m, src)
}
func (m *AddSuperNodeMoneyAccountResponse) XXX_Size() int {
	return xxx_messageInfo_AddSuperNodeMoneyAccountResponse.Size(m)
}
func (m *AddSuperNodeMoneyAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddSuperNodeMoneyAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddSuperNodeMoneyAccountResponse proto.InternalMessageInfo

func (m *AddSuperNodeMoneyAccountResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *AddSuperNodeMoneyAccountResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

func init() {
	proto.RegisterType((*GetSuperNodeActiveMoneyAccountRequest)(nil), "appserver_serves_ui.GetSuperNodeActiveMoneyAccountRequest")
	proto.RegisterType((*GetSuperNodeActiveMoneyAccountResponse)(nil), "appserver_serves_ui.GetSuperNodeActiveMoneyAccountResponse")
	proto.RegisterType((*AddSuperNodeMoneyAccountRequest)(nil), "appserver_serves_ui.AddSuperNodeMoneyAccountRequest")
	proto.RegisterType((*AddSuperNodeMoneyAccountResponse)(nil), "appserver_serves_ui.AddSuperNodeMoneyAccountResponse")
}

func init() { proto.RegisterFile("super_node.proto", fileDescriptor_02e142dc5bc4ebd3) }

var fileDescriptor_02e142dc5bc4ebd3 = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xd1, 0x6a, 0xd4, 0x40,
	0x14, 0x86, 0x99, 0x2e, 0x2e, 0xee, 0x6c, 0x29, 0x75, 0xc4, 0x12, 0x82, 0xe8, 0x1a, 0x54, 0x8a,
	0xe0, 0x46, 0xb6, 0x16, 0xb4, 0xbd, 0xca, 0x55, 0xf1, 0x42, 0x91, 0xf4, 0x01, 0x86, 0x49, 0xe6,
	0x18, 0x06, 0xd6, 0x99, 0x38, 0x33, 0x59, 0x2c, 0xa5, 0x37, 0xfa, 0x08, 0x82, 0x2f, 0xe1, 0x9d,
	0xcf, 0xe0, 0x1b, 0x08, 0x3e, 0x81, 0x0f, 0x22, 0x99, 0x99, 0xc6, 0x15, 0xb2, 0xae, 0x96, 0xbd,
	0x0a, 0x73, 0x72, 0xfe, 0x39, 0xdf, 0xf9, 0xcf, 0x49, 0xf0, 0xae, 0x69, 0x6a, 0xd0, 0x54, 0x2a,
	0x0e, 0xd3, 0x5a, 0x2b, 0xab, 0xc8, 0x4d, 0x56, 0xd7, 0x06, 0xf4, 0x02, 0x34, 0x75, 0x0f, 0x43,
	0x1b, 0x11, 0xdf, 0xae, 0x94, 0xaa, 0xe6, 0x90, 0xb2, 0x5a, 0xa4, 0x4c, 0x4a, 0x65, 0x99, 0x15,
	0x4a, 0x1a, 0x2f, 0x89, 0x6f, 0xc0, 0x7b, 0x4b, 0x59, 0x59, 0xaa, 0x46, 0xda, 0x10, 0xda, 0x11,
	0xd2, 0x82, 0x96, 0x6c, 0xee, 0xcf, 0xc9, 0x19, 0x7e, 0x70, 0x02, 0xf6, 0xb4, 0x2d, 0xf6, 0x4a,
	0x71, 0xc8, 0x4a, 0x2b, 0x16, 0xf0, 0x52, 0x49, 0x38, 0xcb, 0xbc, 0x2e, 0x87, 0x77, 0x0d, 0x18,
	0x4b, 0x9e, 0x63, 0xfc, 0xb6, 0x0d, 0x53, 0x56, 0x14, 0x3a, 0x42, 0x13, 0xb4, 0xbf, 0x33, 0x8b,
	0xa7, 0x3d, 0x4c, 0x53, 0xa7, 0xce, 0x47, 0x2e, 0x3b, 0x2b, 0x0a, 0x4d, 0x6e, 0xe1, 0xa1, 0xd2,
	0x15, 0x15, 0x3c, 0xda, 0x9a, 0xa0, 0xfd, 0x41, 0x7e, 0x4d, 0xe9, 0xea, 0x05, 0x4f, 0xbe, 0x20,
	0xfc, 0x70, 0x5d, 0x6d, 0x53, 0x2b, 0x69, 0x80, 0x3c, 0xc3, 0x91, 0xf3, 0xa3, 0xb5, 0x83, 0x32,
	0x97, 0x77, 0xd9, 0x97, 0x43, 0x19, 0xe5, 0x7b, 0xdd, 0x7b, 0x7f, 0x4d, 0xb8, 0x81, 0x9c, 0xe0,
	0xed, 0xc6, 0x80, 0xa6, 0xb5, 0x56, 0x6f, 0xc4, 0x1c, 0x1c, 0xc1, 0x78, 0x76, 0xbf, 0x17, 0xfc,
	0xb5, 0xcf, 0xb9, 0xac, 0x9a, 0x8f, 0x5b, 0x65, 0x08, 0x26, 0x9f, 0x11, 0xbe, 0x9b, 0x71, 0xde,
	0xd1, 0x6e, 0xd8, 0xa3, 0x7b, 0x78, 0x3b, 0x34, 0x44, 0x19, 0xe7, 0xda, 0x71, 0x8e, 0xf2, 0x71,
	0x88, 0x65, 0x9c, 0x2f, 0xdb, 0x38, 0x58, 0xb6, 0xf1, 0x23, 0xc2, 0x93, 0xd5, 0x60, 0xc1, 0xc0,
	0x3d, 0x3c, 0x34, 0x96, 0xd9, 0xc6, 0x38, 0xaa, 0xeb, 0x79, 0x38, 0x6d, 0xcc, 0x9e, 0xd9, 0xd7,
	0x01, 0xde, 0xed, 0x10, 0x4e, 0x41, 0x2f, 0x44, 0x09, 0xe4, 0x07, 0xc2, 0x77, 0xfe, 0x3e, 0x61,
	0x72, 0xd4, 0x5b, 0xea, 0x9f, 0x56, 0x32, 0x3e, 0xbe, 0x92, 0xd6, 0xd3, 0x27, 0xc7, 0x1f, 0xbe,
	0xff, 0xfc, 0xb4, 0x75, 0x48, 0x0e, 0xdc, 0xb7, 0xd3, 0x6d, 0x4f, 0x7a, 0xee, 0x3d, 0xbe, 0x48,
	0xfd, 0x9a, 0x3d, 0x0e, 0x13, 0x48, 0xcf, 0x7f, 0x4f, 0xf6, 0x82, 0x7c, 0x43, 0x38, 0x5a, 0xe5,
	0x39, 0x79, 0xda, 0x8b, 0xb5, 0x66, 0x77, 0xe2, 0xc3, 0xff, 0x54, 0xfd, 0xd9, 0x46, 0xf2, 0x64,
	0x65, 0x1b, 0x9c, 0xf7, 0xf6, 0x70, 0x84, 0x1e, 0x15, 0x43, 0xf7, 0x0f, 0x38, 0xf8, 0x15, 0x00,
	0x00, 0xff, 0xff, 0x7e, 0xfc, 0x24, 0xcb, 0x6d, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SuperNodeServiceClient is the client API for SuperNodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SuperNodeServiceClient interface {
	GetSuperNodeActiveMoneyAccount(ctx context.Context, in *GetSuperNodeActiveMoneyAccountRequest, opts ...grpc.CallOption) (*GetSuperNodeActiveMoneyAccountResponse, error)
	AddSuperNodeMoneyAccount(ctx context.Context, in *AddSuperNodeMoneyAccountRequest, opts ...grpc.CallOption) (*AddSuperNodeMoneyAccountResponse, error)
}

type superNodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSuperNodeServiceClient(cc grpc.ClientConnInterface) SuperNodeServiceClient {
	return &superNodeServiceClient{cc}
}

func (c *superNodeServiceClient) GetSuperNodeActiveMoneyAccount(ctx context.Context, in *GetSuperNodeActiveMoneyAccountRequest, opts ...grpc.CallOption) (*GetSuperNodeActiveMoneyAccountResponse, error) {
	out := new(GetSuperNodeActiveMoneyAccountResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.SuperNodeService/GetSuperNodeActiveMoneyAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *superNodeServiceClient) AddSuperNodeMoneyAccount(ctx context.Context, in *AddSuperNodeMoneyAccountRequest, opts ...grpc.CallOption) (*AddSuperNodeMoneyAccountResponse, error) {
	out := new(AddSuperNodeMoneyAccountResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.SuperNodeService/AddSuperNodeMoneyAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SuperNodeServiceServer is the server API for SuperNodeService service.
type SuperNodeServiceServer interface {
	GetSuperNodeActiveMoneyAccount(context.Context, *GetSuperNodeActiveMoneyAccountRequest) (*GetSuperNodeActiveMoneyAccountResponse, error)
	AddSuperNodeMoneyAccount(context.Context, *AddSuperNodeMoneyAccountRequest) (*AddSuperNodeMoneyAccountResponse, error)
}

// UnimplementedSuperNodeServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSuperNodeServiceServer struct {
}

func (*UnimplementedSuperNodeServiceServer) GetSuperNodeActiveMoneyAccount(ctx context.Context, req *GetSuperNodeActiveMoneyAccountRequest) (*GetSuperNodeActiveMoneyAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSuperNodeActiveMoneyAccount not implemented")
}
func (*UnimplementedSuperNodeServiceServer) AddSuperNodeMoneyAccount(ctx context.Context, req *AddSuperNodeMoneyAccountRequest) (*AddSuperNodeMoneyAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSuperNodeMoneyAccount not implemented")
}

func RegisterSuperNodeServiceServer(s *grpc.Server, srv SuperNodeServiceServer) {
	s.RegisterService(&_SuperNodeService_serviceDesc, srv)
}

func _SuperNodeService_GetSuperNodeActiveMoneyAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSuperNodeActiveMoneyAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SuperNodeServiceServer).GetSuperNodeActiveMoneyAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.SuperNodeService/GetSuperNodeActiveMoneyAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SuperNodeServiceServer).GetSuperNodeActiveMoneyAccount(ctx, req.(*GetSuperNodeActiveMoneyAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SuperNodeService_AddSuperNodeMoneyAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSuperNodeMoneyAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SuperNodeServiceServer).AddSuperNodeMoneyAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.SuperNodeService/AddSuperNodeMoneyAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SuperNodeServiceServer).AddSuperNodeMoneyAccount(ctx, req.(*AddSuperNodeMoneyAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SuperNodeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "appserver_serves_ui.SuperNodeService",
	HandlerType: (*SuperNodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSuperNodeActiveMoneyAccount",
			Handler:    _SuperNodeService_GetSuperNodeActiveMoneyAccount_Handler,
		},
		{
			MethodName: "AddSuperNodeMoneyAccount",
			Handler:    _SuperNodeService_AddSuperNodeMoneyAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "super_node.proto",
}
