// Code generated by protoc-gen-go. DO NOT EDIT.
// source: topup.proto

package m2m_ui

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

type GetTopUpHistoryRequest struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int64    `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	OrgId                int64    `protobuf:"varint,4,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTopUpHistoryRequest) Reset()         { *m = GetTopUpHistoryRequest{} }
func (m *GetTopUpHistoryRequest) String() string { return proto.CompactTextString(m) }
func (*GetTopUpHistoryRequest) ProtoMessage()    {}
func (*GetTopUpHistoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8eec749941d0cb6c, []int{0}
}

func (m *GetTopUpHistoryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTopUpHistoryRequest.Unmarshal(m, b)
}
func (m *GetTopUpHistoryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTopUpHistoryRequest.Marshal(b, m, deterministic)
}
func (m *GetTopUpHistoryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTopUpHistoryRequest.Merge(m, src)
}
func (m *GetTopUpHistoryRequest) XXX_Size() int {
	return xxx_messageInfo_GetTopUpHistoryRequest.Size(m)
}
func (m *GetTopUpHistoryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTopUpHistoryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTopUpHistoryRequest proto.InternalMessageInfo

func (m *GetTopUpHistoryRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *GetTopUpHistoryRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *GetTopUpHistoryRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *GetTopUpHistoryRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

type TopUpHistory struct {
	Amount               float64  `protobuf:"fixed64,1,opt,name=amount,proto3" json:"amount,omitempty"`
	CreatedAt            string   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	TxHash               string   `protobuf:"bytes,3,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TopUpHistory) Reset()         { *m = TopUpHistory{} }
func (m *TopUpHistory) String() string { return proto.CompactTextString(m) }
func (*TopUpHistory) ProtoMessage()    {}
func (*TopUpHistory) Descriptor() ([]byte, []int) {
	return fileDescriptor_8eec749941d0cb6c, []int{1}
}

func (m *TopUpHistory) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TopUpHistory.Unmarshal(m, b)
}
func (m *TopUpHistory) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TopUpHistory.Marshal(b, m, deterministic)
}
func (m *TopUpHistory) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TopUpHistory.Merge(m, src)
}
func (m *TopUpHistory) XXX_Size() int {
	return xxx_messageInfo_TopUpHistory.Size(m)
}
func (m *TopUpHistory) XXX_DiscardUnknown() {
	xxx_messageInfo_TopUpHistory.DiscardUnknown(m)
}

var xxx_messageInfo_TopUpHistory proto.InternalMessageInfo

func (m *TopUpHistory) GetAmount() float64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *TopUpHistory) GetCreatedAt() string {
	if m != nil {
		return m.CreatedAt
	}
	return ""
}

func (m *TopUpHistory) GetTxHash() string {
	if m != nil {
		return m.TxHash
	}
	return ""
}

type GetTopUpHistoryResponse struct {
	Count                int64            `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	TopupHistory         []*TopUpHistory  `protobuf:"bytes,2,rep,name=topup_history,json=topupHistory,proto3" json:"topup_history,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,3,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetTopUpHistoryResponse) Reset()         { *m = GetTopUpHistoryResponse{} }
func (m *GetTopUpHistoryResponse) String() string { return proto.CompactTextString(m) }
func (*GetTopUpHistoryResponse) ProtoMessage()    {}
func (*GetTopUpHistoryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8eec749941d0cb6c, []int{2}
}

func (m *GetTopUpHistoryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTopUpHistoryResponse.Unmarshal(m, b)
}
func (m *GetTopUpHistoryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTopUpHistoryResponse.Marshal(b, m, deterministic)
}
func (m *GetTopUpHistoryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTopUpHistoryResponse.Merge(m, src)
}
func (m *GetTopUpHistoryResponse) XXX_Size() int {
	return xxx_messageInfo_GetTopUpHistoryResponse.Size(m)
}
func (m *GetTopUpHistoryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTopUpHistoryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTopUpHistoryResponse proto.InternalMessageInfo

func (m *GetTopUpHistoryResponse) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *GetTopUpHistoryResponse) GetTopupHistory() []*TopUpHistory {
	if m != nil {
		return m.TopupHistory
	}
	return nil
}

func (m *GetTopUpHistoryResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

type GetTopUpDestinationRequest struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	MoneyAbbr            Money    `protobuf:"varint,2,opt,name=money_abbr,json=moneyAbbr,proto3,enum=m2m_ui.Money" json:"money_abbr,omitempty"`
	OrgId                int64    `protobuf:"varint,3,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTopUpDestinationRequest) Reset()         { *m = GetTopUpDestinationRequest{} }
func (m *GetTopUpDestinationRequest) String() string { return proto.CompactTextString(m) }
func (*GetTopUpDestinationRequest) ProtoMessage()    {}
func (*GetTopUpDestinationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8eec749941d0cb6c, []int{3}
}

func (m *GetTopUpDestinationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTopUpDestinationRequest.Unmarshal(m, b)
}
func (m *GetTopUpDestinationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTopUpDestinationRequest.Marshal(b, m, deterministic)
}
func (m *GetTopUpDestinationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTopUpDestinationRequest.Merge(m, src)
}
func (m *GetTopUpDestinationRequest) XXX_Size() int {
	return xxx_messageInfo_GetTopUpDestinationRequest.Size(m)
}
func (m *GetTopUpDestinationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTopUpDestinationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTopUpDestinationRequest proto.InternalMessageInfo

func (m *GetTopUpDestinationRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *GetTopUpDestinationRequest) GetMoneyAbbr() Money {
	if m != nil {
		return m.MoneyAbbr
	}
	return Money_ETH
}

func (m *GetTopUpDestinationRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

type GetTopUpDestinationResponse struct {
	ActiveAccount        string           `protobuf:"bytes,1,opt,name=active_account,json=activeAccount,proto3" json:"active_account,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,2,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetTopUpDestinationResponse) Reset()         { *m = GetTopUpDestinationResponse{} }
func (m *GetTopUpDestinationResponse) String() string { return proto.CompactTextString(m) }
func (*GetTopUpDestinationResponse) ProtoMessage()    {}
func (*GetTopUpDestinationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8eec749941d0cb6c, []int{4}
}

func (m *GetTopUpDestinationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTopUpDestinationResponse.Unmarshal(m, b)
}
func (m *GetTopUpDestinationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTopUpDestinationResponse.Marshal(b, m, deterministic)
}
func (m *GetTopUpDestinationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTopUpDestinationResponse.Merge(m, src)
}
func (m *GetTopUpDestinationResponse) XXX_Size() int {
	return xxx_messageInfo_GetTopUpDestinationResponse.Size(m)
}
func (m *GetTopUpDestinationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTopUpDestinationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTopUpDestinationResponse proto.InternalMessageInfo

func (m *GetTopUpDestinationResponse) GetActiveAccount() string {
	if m != nil {
		return m.ActiveAccount
	}
	return ""
}

func (m *GetTopUpDestinationResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

func init() {
	proto.RegisterType((*GetTopUpHistoryRequest)(nil), "m2m_ui.GetTopUpHistoryRequest")
	proto.RegisterType((*TopUpHistory)(nil), "m2m_ui.TopUpHistory")
	proto.RegisterType((*GetTopUpHistoryResponse)(nil), "m2m_ui.GetTopUpHistoryResponse")
	proto.RegisterType((*GetTopUpDestinationRequest)(nil), "m2m_ui.GetTopUpDestinationRequest")
	proto.RegisterType((*GetTopUpDestinationResponse)(nil), "m2m_ui.GetTopUpDestinationResponse")
}

func init() { proto.RegisterFile("topup.proto", fileDescriptor_8eec749941d0cb6c) }

var fileDescriptor_8eec749941d0cb6c = []byte{
	// 487 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xcd, 0x8e, 0xd3, 0x30,
	0x18, 0x54, 0x12, 0x08, 0xca, 0xd7, 0x76, 0x11, 0xde, 0x9f, 0x56, 0x59, 0x7e, 0xaa, 0x20, 0xa4,
	0x1e, 0xa0, 0x95, 0xc2, 0x09, 0x6e, 0x95, 0x90, 0xd8, 0x3d, 0x20, 0x21, 0x03, 0x57, 0x22, 0x27,
	0x75, 0x5b, 0x4b, 0x4d, 0x6c, 0x6c, 0xa7, 0xea, 0x82, 0x90, 0x10, 0x27, 0xc4, 0x95, 0x77, 0xe0,
	0x85, 0x78, 0x05, 0x1e, 0x04, 0xc5, 0x76, 0xa0, 0x6c, 0xcb, 0xcf, 0xf1, 0x9b, 0x4c, 0xbe, 0x19,
	0xcf, 0xd8, 0xd0, 0xd1, 0x5c, 0xd4, 0x62, 0x2c, 0x24, 0xd7, 0x1c, 0x85, 0x65, 0x5a, 0x66, 0x35,
	0x8b, 0x6f, 0x2e, 0x38, 0x5f, 0xac, 0xe8, 0x84, 0x08, 0x36, 0x21, 0x55, 0xc5, 0x35, 0xd1, 0x8c,
	0x57, 0xca, 0xb2, 0xe2, 0x9e, 0x90, 0x7c, 0xce, 0x56, 0xd4, 0x8d, 0x37, 0xe8, 0x46, 0x67, 0xa4,
	0x28, 0x78, 0x5d, 0x69, 0x0b, 0x25, 0x6b, 0x38, 0x79, 0x4a, 0xf5, 0x4b, 0x2e, 0x5e, 0x89, 0x33,
	0xa6, 0x34, 0x97, 0x17, 0x98, 0xbe, 0xa9, 0xa9, 0xd2, 0xa8, 0x0f, 0xd7, 0x6a, 0x45, 0x65, 0xc6,
	0x66, 0x03, 0x6f, 0xe8, 0x8d, 0x02, 0x1c, 0x36, 0xe3, 0xf9, 0x0c, 0x9d, 0x40, 0xc8, 0xe7, 0x73,
	0x45, 0xf5, 0xc0, 0xb7, 0xb8, 0x9d, 0xd0, 0x11, 0x5c, 0x5d, 0xb1, 0x92, 0xe9, 0x41, 0x60, 0x60,
	0x3b, 0xa0, 0x63, 0x08, 0xb9, 0x5c, 0x34, 0x5b, 0xae, 0x58, 0x98, 0xcb, 0xc5, 0xf9, 0x2c, 0x79,
	0x0d, 0xdd, 0x6d, 0xd1, 0x66, 0x29, 0x29, 0x1b, 0x5f, 0x46, 0xcc, 0xc3, 0x6e, 0x42, 0xb7, 0x00,
	0x0a, 0x49, 0x89, 0xa6, 0xb3, 0x8c, 0x58, 0xc1, 0x08, 0x47, 0x0e, 0x99, 0x1a, 0x93, 0x7a, 0x93,
	0x2d, 0x89, 0x5a, 0x1a, 0xd5, 0x08, 0x87, 0x7a, 0x73, 0x46, 0xd4, 0x32, 0xf9, 0xea, 0x41, 0x7f,
	0xe7, 0x60, 0x4a, 0xf0, 0x4a, 0xd1, 0xc6, 0x68, 0xf1, 0x53, 0x2a, 0xc0, 0x76, 0x40, 0x8f, 0xa0,
	0x67, 0x02, 0xce, 0x96, 0x96, 0x3e, 0xf0, 0x87, 0xc1, 0xa8, 0x93, 0x1e, 0x8d, 0x6d, 0xd2, 0xe3,
	0xdf, 0x56, 0x75, 0x0d, 0xb5, 0x35, 0xff, 0x18, 0xba, 0x26, 0x2a, 0x97, 0xb6, 0xb1, 0xd2, 0x49,
	0xfb, 0xed, 0x9f, 0xcf, 0x2d, 0xdc, 0xea, 0xe3, 0x4e, 0x43, 0x76, 0x60, 0xf2, 0x16, 0xe2, 0xd6,
	0xe7, 0x13, 0xaa, 0x34, 0xab, 0x4c, 0x81, 0xff, 0x2c, 0xe1, 0x3e, 0x40, 0xc9, 0x2b, 0x7a, 0x91,
	0x91, 0x3c, 0x97, 0x26, 0x97, 0x83, 0xb4, 0xd7, 0x0a, 0x3e, 0x6b, 0xbe, 0xe0, 0xc8, 0x10, 0xa6,
	0x79, 0x2e, 0xb7, 0x4a, 0x08, 0xb6, 0x4b, 0xf8, 0xe0, 0xc1, 0xe9, 0x5e, 0x71, 0x17, 0xd4, 0x3d,
	0x38, 0x20, 0x85, 0x66, 0x6b, 0xda, 0x5e, 0x1a, 0x63, 0x22, 0xc2, 0x3d, 0x8b, 0x4e, 0x2d, 0xb8,
	0x73, 0x7c, 0xff, 0xff, 0x8f, 0x9f, 0x7e, 0xf6, 0xdd, 0x45, 0x78, 0x41, 0xe5, 0x9a, 0x15, 0x14,
	0x71, 0xb8, 0x7e, 0xa9, 0x37, 0x74, 0xbb, 0xdd, 0xb4, 0xff, 0xa6, 0xc6, 0x77, 0xfe, 0xf8, 0xdd,
	0x2a, 0x26, 0xa7, 0x1f, 0xbf, 0x7d, 0xff, 0xe2, 0x1f, 0xa3, 0x43, 0xf3, 0x4c, 0x34, 0x17, 0x0f,
	0x6a, 0x31, 0x71, 0x35, 0xa3, 0x4f, 0x1e, 0x1c, 0xee, 0x09, 0x01, 0x25, 0x97, 0xb7, 0xee, 0xd6,
	0x13, 0xdf, 0xfd, 0x2b, 0xc7, 0xa9, 0x8f, 0x8c, 0x7a, 0x82, 0x86, 0xdb, 0xea, 0x2e, 0xd0, 0xc9,
	0xbb, 0x5f, 0x2d, 0xbe, 0xcf, 0x43, 0xf3, 0x26, 0x1f, 0xfe, 0x08, 0x00, 0x00, 0xff, 0xff, 0x2f,
	0x95, 0x53, 0x1a, 0xea, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TopUpServiceClient is the client API for TopUpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TopUpServiceClient interface {
	GetTopUpHistory(ctx context.Context, in *GetTopUpHistoryRequest, opts ...grpc.CallOption) (*GetTopUpHistoryResponse, error)
	GetTopUpDestination(ctx context.Context, in *GetTopUpDestinationRequest, opts ...grpc.CallOption) (*GetTopUpDestinationResponse, error)
}

type topUpServiceClient struct {
	cc *grpc.ClientConn
}

func NewTopUpServiceClient(cc *grpc.ClientConn) TopUpServiceClient {
	return &topUpServiceClient{cc}
}

func (c *topUpServiceClient) GetTopUpHistory(ctx context.Context, in *GetTopUpHistoryRequest, opts ...grpc.CallOption) (*GetTopUpHistoryResponse, error) {
	out := new(GetTopUpHistoryResponse)
	err := c.cc.Invoke(ctx, "/m2m_ui.TopUpService/GetTopUpHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topUpServiceClient) GetTopUpDestination(ctx context.Context, in *GetTopUpDestinationRequest, opts ...grpc.CallOption) (*GetTopUpDestinationResponse, error) {
	out := new(GetTopUpDestinationResponse)
	err := c.cc.Invoke(ctx, "/m2m_ui.TopUpService/GetTopUpDestination", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TopUpServiceServer is the server API for TopUpService service.
type TopUpServiceServer interface {
	GetTopUpHistory(context.Context, *GetTopUpHistoryRequest) (*GetTopUpHistoryResponse, error)
	GetTopUpDestination(context.Context, *GetTopUpDestinationRequest) (*GetTopUpDestinationResponse, error)
}

// UnimplementedTopUpServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTopUpServiceServer struct {
}

func (*UnimplementedTopUpServiceServer) GetTopUpHistory(ctx context.Context, req *GetTopUpHistoryRequest) (*GetTopUpHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopUpHistory not implemented")
}
func (*UnimplementedTopUpServiceServer) GetTopUpDestination(ctx context.Context, req *GetTopUpDestinationRequest) (*GetTopUpDestinationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopUpDestination not implemented")
}

func RegisterTopUpServiceServer(s *grpc.Server, srv TopUpServiceServer) {
	s.RegisterService(&_TopUpService_serviceDesc, srv)
}

func _TopUpService_GetTopUpHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTopUpHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopUpServiceServer).GetTopUpHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_ui.TopUpService/GetTopUpHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopUpServiceServer).GetTopUpHistory(ctx, req.(*GetTopUpHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TopUpService_GetTopUpDestination_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTopUpDestinationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopUpServiceServer).GetTopUpDestination(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_ui.TopUpService/GetTopUpDestination",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopUpServiceServer).GetTopUpDestination(ctx, req.(*GetTopUpDestinationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TopUpService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m2m_ui.TopUpService",
	HandlerType: (*TopUpServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTopUpHistory",
			Handler:    _TopUpService_GetTopUpHistory_Handler,
		},
		{
			MethodName: "GetTopUpDestination",
			Handler:    _TopUpService_GetTopUpDestination_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "topup.proto",
}
