// Code generated by protoc-gen-go. DO NOT EDIT.
// source: topup.proto

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

type Money int32

const (
	Money_ETH     Money = 0
	Money_ETH_MXC Money = 1
	Money_TETH    Money = 2
)

var Money_name = map[int32]string{
	0: "ETH",
	1: "ETH_MXC",
	2: "TETH",
}

var Money_value = map[string]int32{
	"ETH":     0,
	"ETH_MXC": 1,
	"TETH":    2,
}

func (x Money) String() string {
	return proto.EnumName(Money_name, int32(x))
}

func (Money) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8eec749941d0cb6c, []int{0}
}

type GetTopUpHistoryRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int64    `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
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

func (m *GetTopUpHistoryRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
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
	Count                int64           `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	TopupHistory         []*TopUpHistory `protobuf:"bytes,2,rep,name=topup_history,json=topupHistory,proto3" json:"topup_history,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
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

type GetTopUpDestinationRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
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

func (m *GetTopUpDestinationRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

type GetTopUpDestinationResponse struct {
	ActiveAccount        string   `protobuf:"bytes,1,opt,name=active_account,json=activeAccount,proto3" json:"active_account,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
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

func init() {
	proto.RegisterEnum("appserver_serves_ui.Money", Money_name, Money_value)
	proto.RegisterType((*GetTopUpHistoryRequest)(nil), "appserver_serves_ui.GetTopUpHistoryRequest")
	proto.RegisterType((*TopUpHistory)(nil), "appserver_serves_ui.TopUpHistory")
	proto.RegisterType((*GetTopUpHistoryResponse)(nil), "appserver_serves_ui.GetTopUpHistoryResponse")
	proto.RegisterType((*GetTopUpDestinationRequest)(nil), "appserver_serves_ui.GetTopUpDestinationRequest")
	proto.RegisterType((*GetTopUpDestinationResponse)(nil), "appserver_serves_ui.GetTopUpDestinationResponse")
}

func init() { proto.RegisterFile("topup.proto", fileDescriptor_8eec749941d0cb6c) }

var fileDescriptor_8eec749941d0cb6c = []byte{
	// 430 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xcd, 0x6e, 0x13, 0x31,
	0x14, 0x85, 0x99, 0x19, 0x92, 0x90, 0x9b, 0x16, 0x2a, 0x87, 0xb6, 0x51, 0x0a, 0x52, 0x19, 0x09,
	0x51, 0xf1, 0x93, 0xa0, 0xf6, 0x09, 0x2a, 0x5a, 0x08, 0x8b, 0x6e, 0x86, 0x20, 0xb1, 0x01, 0xcb,
	0x4c, 0xdc, 0x19, 0x4b, 0xed, 0x5c, 0x63, 0xdf, 0x09, 0xed, 0x96, 0x1d, 0x6b, 0xc4, 0x8b, 0xf0,
	0x2a, 0xbc, 0x02, 0x0f, 0x82, 0x62, 0x9b, 0xdf, 0x8e, 0xa2, 0xac, 0x46, 0xf7, 0xc8, 0xe7, 0xdc,
	0xe3, 0x4f, 0x1e, 0xe8, 0x11, 0xea, 0x5a, 0x8f, 0xb4, 0x41, 0x42, 0xd6, 0x17, 0x5a, 0x5b, 0x69,
	0xe6, 0xd2, 0x70, 0xf7, 0xb1, 0xbc, 0x56, 0xc3, 0x3b, 0x05, 0x62, 0x71, 0x26, 0xc7, 0x42, 0xab,
	0xb1, 0xa8, 0x2a, 0x24, 0x41, 0x0a, 0x2b, 0xeb, 0x2d, 0xe9, 0x5b, 0xd8, 0x7a, 0x21, 0x69, 0x8a,
	0xfa, 0xb5, 0x9e, 0x28, 0x4b, 0x68, 0x2e, 0x33, 0xf9, 0xa1, 0x96, 0x96, 0xd8, 0x26, 0xb4, 0xd1,
	0x14, 0x5c, 0xcd, 0x06, 0xd1, 0x6e, 0xb4, 0x97, 0x64, 0x2d, 0x34, 0xc5, 0xcb, 0x19, 0xdb, 0x82,
	0x36, 0x9e, 0x9e, 0x5a, 0x49, 0x83, 0xd8, 0xc9, 0x61, 0x62, 0xb7, 0xa1, 0x75, 0xa6, 0xce, 0x15,
	0x0d, 0x12, 0x7f, 0xda, 0x0d, 0xe9, 0x3b, 0x58, 0xfb, 0x3b, 0x7b, 0xe1, 0x16, 0xe7, 0x58, 0x57,
	0xe4, 0x42, 0xa3, 0x2c, 0x4c, 0xec, 0x2e, 0x40, 0x6e, 0xa4, 0x20, 0x39, 0xe3, 0xc2, 0x27, 0x77,
	0xb3, 0x6e, 0x50, 0x0e, 0x89, 0x6d, 0x43, 0x87, 0x2e, 0x78, 0x29, 0x6c, 0xe9, 0xe2, 0xbb, 0x59,
	0x9b, 0x2e, 0x26, 0xc2, 0x96, 0xe9, 0x47, 0xd8, 0xbe, 0x52, 0xdf, 0x6a, 0xac, 0xac, 0x5c, 0x14,
	0xca, 0x7f, 0x6f, 0x4a, 0x32, 0x3f, 0xb0, 0xe7, 0xb0, 0xee, 0x88, 0xf1, 0xd2, 0x1f, 0x1f, 0xc4,
	0xbb, 0xc9, 0x5e, 0x6f, 0xff, 0xde, 0xa8, 0x01, 0xdd, 0xe8, 0x9f, 0xdc, 0x35, 0xe7, 0x0b, 0x53,
	0x7a, 0x00, 0xc3, 0x5f, 0x8b, 0x8f, 0xa4, 0x25, 0x55, 0x39, 0xaa, 0xcb, 0xd9, 0xa5, 0x47, 0xb0,
	0xd3, 0x68, 0x0a, 0x8d, 0xef, 0xc3, 0x4d, 0x91, 0x93, 0x9a, 0x4b, 0x2e, 0xf2, 0x3f, 0xd5, 0xbb,
	0xd9, 0xba, 0x57, 0x0f, 0xbd, 0xf8, 0xf0, 0x01, 0xb4, 0x4e, 0xb0, 0x92, 0x97, 0xac, 0x03, 0xc9,
	0xf1, 0x74, 0xb2, 0x71, 0x8d, 0xf5, 0xa0, 0x73, 0x3c, 0x9d, 0xf0, 0x93, 0x37, 0xcf, 0x36, 0x22,
	0x76, 0x03, 0xae, 0x4f, 0x17, 0x72, 0xbc, 0xff, 0x2d, 0x0e, 0xf4, 0x5f, 0x49, 0x33, 0x57, 0xb9,
	0x64, 0x9f, 0x23, 0xb8, 0xf5, 0x1f, 0x2e, 0xf6, 0xa8, 0xf1, 0xe6, 0xcd, 0x6f, 0x62, 0xf8, 0x78,
	0xb5, 0xc3, 0xfe, 0x3e, 0xe9, 0xce, 0xa7, 0xef, 0x3f, 0xbe, 0xc4, 0x9b, 0xac, 0xef, 0xde, 0x1e,
	0xa1, 0x7e, 0x52, 0xeb, 0x71, 0xe0, 0xce, 0xbe, 0x46, 0xd0, 0x6f, 0x80, 0xc1, 0xc6, 0x4b, 0x57,
	0x5c, 0x65, 0x3d, 0x7c, 0xba, 0xba, 0x61, 0x59, 0xaf, 0x80, 0xfc, 0x7d, 0xdb, 0xfd, 0x17, 0x07,
	0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x0f, 0x18, 0x47, 0xb5, 0x59, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TopUpServiceClient is the client API for TopUpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TopUpServiceClient interface {
	GetTopUpHistory(ctx context.Context, in *GetTopUpHistoryRequest, opts ...grpc.CallOption) (*GetTopUpHistoryResponse, error)
	GetTopUpDestination(ctx context.Context, in *GetTopUpDestinationRequest, opts ...grpc.CallOption) (*GetTopUpDestinationResponse, error)
}

type topUpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTopUpServiceClient(cc grpc.ClientConnInterface) TopUpServiceClient {
	return &topUpServiceClient{cc}
}

func (c *topUpServiceClient) GetTopUpHistory(ctx context.Context, in *GetTopUpHistoryRequest, opts ...grpc.CallOption) (*GetTopUpHistoryResponse, error) {
	out := new(GetTopUpHistoryResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.TopUpService/GetTopUpHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topUpServiceClient) GetTopUpDestination(ctx context.Context, in *GetTopUpDestinationRequest, opts ...grpc.CallOption) (*GetTopUpDestinationResponse, error) {
	out := new(GetTopUpDestinationResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.TopUpService/GetTopUpDestination", in, out, opts...)
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
		FullMethod: "/appserver_serves_ui.TopUpService/GetTopUpHistory",
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
		FullMethod: "/appserver_serves_ui.TopUpService/GetTopUpDestination",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopUpServiceServer).GetTopUpDestination(ctx, req.(*GetTopUpDestinationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TopUpService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "appserver_serves_ui.TopUpService",
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
