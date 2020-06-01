// Code generated by protoc-gen-go. DO NOT EDIT.
// source: settings.proto

package appserver_serves_ui

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

type GetSettingsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSettingsRequest) Reset()         { *m = GetSettingsRequest{} }
func (m *GetSettingsRequest) String() string { return proto.CompactTextString(m) }
func (*GetSettingsRequest) ProtoMessage()    {}
func (*GetSettingsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c7cab62fa432213, []int{0}
}

func (m *GetSettingsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSettingsRequest.Unmarshal(m, b)
}
func (m *GetSettingsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSettingsRequest.Marshal(b, m, deterministic)
}
func (m *GetSettingsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSettingsRequest.Merge(m, src)
}
func (m *GetSettingsRequest) XXX_Size() int {
	return xxx_messageInfo_GetSettingsRequest.Size(m)
}
func (m *GetSettingsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSettingsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSettingsRequest proto.InternalMessageInfo

type GetSettingsResponse struct {
	// when supernode income is lower than expected revenue, warn system owner to increase income
	LowBalanceWarning    string  `protobuf:"bytes,1,opt,name=low_balance_warning,json=lowBalanceWarning,proto3" json:"low_balance_warning,omitempty"`
	DownlinkPrice        float64 `protobuf:"fixed64,2,opt,name=downlink_price,json=downlinkPrice,proto3" json:"downlink_price,omitempty"`
	SupernodeIncomeRatio float64 `protobuf:"fixed64,3,opt,name=supernode_income_ratio,json=supernodeIncomeRatio,proto3" json:"supernode_income_ratio,omitempty"`
	// this is the percentage of reward from supernode income
	StakingPercentage float64 `protobuf:"fixed64,4,opt,name=staking_percentage,json=stakingPercentage,proto3" json:"staking_percentage,omitempty"`
	// this is the percentage of expected reward from staking amount
	StakingExpectedRevenuePercentage float64  `protobuf:"fixed64,5,opt,name=staking_expected_revenue_percentage,json=stakingExpectedRevenuePercentage,proto3" json:"staking_expected_revenue_percentage,omitempty"`
	XXX_NoUnkeyedLiteral             struct{} `json:"-"`
	XXX_unrecognized                 []byte   `json:"-"`
	XXX_sizecache                    int32    `json:"-"`
}

func (m *GetSettingsResponse) Reset()         { *m = GetSettingsResponse{} }
func (m *GetSettingsResponse) String() string { return proto.CompactTextString(m) }
func (*GetSettingsResponse) ProtoMessage()    {}
func (*GetSettingsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c7cab62fa432213, []int{1}
}

func (m *GetSettingsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSettingsResponse.Unmarshal(m, b)
}
func (m *GetSettingsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSettingsResponse.Marshal(b, m, deterministic)
}
func (m *GetSettingsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSettingsResponse.Merge(m, src)
}
func (m *GetSettingsResponse) XXX_Size() int {
	return xxx_messageInfo_GetSettingsResponse.Size(m)
}
func (m *GetSettingsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSettingsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSettingsResponse proto.InternalMessageInfo

func (m *GetSettingsResponse) GetLowBalanceWarning() string {
	if m != nil {
		return m.LowBalanceWarning
	}
	return ""
}

func (m *GetSettingsResponse) GetDownlinkPrice() float64 {
	if m != nil {
		return m.DownlinkPrice
	}
	return 0
}

func (m *GetSettingsResponse) GetSupernodeIncomeRatio() float64 {
	if m != nil {
		return m.SupernodeIncomeRatio
	}
	return 0
}

func (m *GetSettingsResponse) GetStakingPercentage() float64 {
	if m != nil {
		return m.StakingPercentage
	}
	return 0
}

func (m *GetSettingsResponse) GetStakingExpectedRevenuePercentage() float64 {
	if m != nil {
		return m.StakingExpectedRevenuePercentage
	}
	return 0
}

type ModifySettingsRequest struct {
	LowBalanceWarning          *wrappers.Int64Value `protobuf:"bytes,1,opt,name=lowBalanceWarning,proto3" json:"lowBalanceWarning,omitempty"`
	DownlinkFee                *wrappers.Int64Value `protobuf:"bytes,2,opt,name=downlinkFee,proto3" json:"downlinkFee,omitempty"`
	TransactionPercentageShare *wrappers.Int64Value `protobuf:"bytes,3,opt,name=transactionPercentageShare,proto3" json:"transactionPercentageShare,omitempty"`
	XXX_NoUnkeyedLiteral       struct{}             `json:"-"`
	XXX_unrecognized           []byte               `json:"-"`
	XXX_sizecache              int32                `json:"-"`
}

func (m *ModifySettingsRequest) Reset()         { *m = ModifySettingsRequest{} }
func (m *ModifySettingsRequest) String() string { return proto.CompactTextString(m) }
func (*ModifySettingsRequest) ProtoMessage()    {}
func (*ModifySettingsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c7cab62fa432213, []int{2}
}

func (m *ModifySettingsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ModifySettingsRequest.Unmarshal(m, b)
}
func (m *ModifySettingsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ModifySettingsRequest.Marshal(b, m, deterministic)
}
func (m *ModifySettingsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ModifySettingsRequest.Merge(m, src)
}
func (m *ModifySettingsRequest) XXX_Size() int {
	return xxx_messageInfo_ModifySettingsRequest.Size(m)
}
func (m *ModifySettingsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ModifySettingsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ModifySettingsRequest proto.InternalMessageInfo

func (m *ModifySettingsRequest) GetLowBalanceWarning() *wrappers.Int64Value {
	if m != nil {
		return m.LowBalanceWarning
	}
	return nil
}

func (m *ModifySettingsRequest) GetDownlinkFee() *wrappers.Int64Value {
	if m != nil {
		return m.DownlinkFee
	}
	return nil
}

func (m *ModifySettingsRequest) GetTransactionPercentageShare() *wrappers.Int64Value {
	if m != nil {
		return m.TransactionPercentageShare
	}
	return nil
}

type ModifySettingsResponse struct {
	Status               bool     `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ModifySettingsResponse) Reset()         { *m = ModifySettingsResponse{} }
func (m *ModifySettingsResponse) String() string { return proto.CompactTextString(m) }
func (*ModifySettingsResponse) ProtoMessage()    {}
func (*ModifySettingsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c7cab62fa432213, []int{3}
}

func (m *ModifySettingsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ModifySettingsResponse.Unmarshal(m, b)
}
func (m *ModifySettingsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ModifySettingsResponse.Marshal(b, m, deterministic)
}
func (m *ModifySettingsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ModifySettingsResponse.Merge(m, src)
}
func (m *ModifySettingsResponse) XXX_Size() int {
	return xxx_messageInfo_ModifySettingsResponse.Size(m)
}
func (m *ModifySettingsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ModifySettingsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ModifySettingsResponse proto.InternalMessageInfo

func (m *ModifySettingsResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func init() {
	proto.RegisterType((*GetSettingsRequest)(nil), "appserver_serves_ui.GetSettingsRequest")
	proto.RegisterType((*GetSettingsResponse)(nil), "appserver_serves_ui.GetSettingsResponse")
	proto.RegisterType((*ModifySettingsRequest)(nil), "appserver_serves_ui.ModifySettingsRequest")
	proto.RegisterType((*ModifySettingsResponse)(nil), "appserver_serves_ui.ModifySettingsResponse")
}

func init() { proto.RegisterFile("settings.proto", fileDescriptor_6c7cab62fa432213) }

var fileDescriptor_6c7cab62fa432213 = []byte{
	// 477 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xd1, 0x6a, 0x13, 0x41,
	0x14, 0x65, 0xa3, 0x16, 0x9d, 0xd0, 0x94, 0x4e, 0xda, 0xb0, 0xac, 0x22, 0x61, 0x45, 0x0c, 0x15,
	0x37, 0x52, 0x8b, 0x0f, 0x82, 0x2f, 0x82, 0x4a, 0x1e, 0x0a, 0x65, 0x03, 0xfa, 0xe0, 0xc3, 0x30,
	0xd9, 0xdc, 0xae, 0x43, 0xd7, 0x3b, 0xe3, 0xcc, 0x6c, 0x56, 0x5f, 0xf5, 0x07, 0x04, 0x7f, 0xc0,
	0x7f, 0xf2, 0x17, 0xfc, 0x0e, 0x91, 0x9d, 0x9d, 0x6d, 0x93, 0x36, 0x98, 0x3e, 0x2d, 0x33, 0xe7,
	0xdc, 0xbb, 0xe7, 0x9e, 0x7b, 0x86, 0xf4, 0x0c, 0x58, 0x2b, 0x30, 0x37, 0x89, 0xd2, 0xd2, 0x4a,
	0xda, 0xe7, 0x4a, 0x19, 0xd0, 0x0b, 0xd0, 0xcc, 0x7d, 0x0c, 0x2b, 0x45, 0x74, 0x2f, 0x97, 0x32,
	0x2f, 0x60, 0xcc, 0x95, 0x18, 0x73, 0x44, 0x69, 0xb9, 0x15, 0x12, 0x7d, 0x49, 0x74, 0xdf, 0xa3,
	0xee, 0x34, 0x2b, 0x4f, 0xc7, 0x95, 0xe6, 0x4a, 0x81, 0xf6, 0x78, 0xbc, 0x47, 0xe8, 0x5b, 0xb0,
	0x53, 0xff, 0x9f, 0x14, 0x3e, 0x97, 0x60, 0x6c, 0xfc, 0xab, 0x43, 0xfa, 0x2b, 0xd7, 0x46, 0x49,
	0x34, 0x40, 0x13, 0xd2, 0x2f, 0x64, 0xc5, 0x66, 0xbc, 0xe0, 0x98, 0x01, 0xab, 0xb8, 0x46, 0x81,
	0x79, 0x18, 0x0c, 0x83, 0xd1, 0x9d, 0x74, 0xb7, 0x90, 0xd5, 0xab, 0x06, 0x79, 0xdf, 0x00, 0xf4,
	0x21, 0xe9, 0xcd, 0x65, 0x85, 0x85, 0xc0, 0x33, 0xa6, 0xb4, 0xc8, 0x20, 0xec, 0x0c, 0x83, 0x51,
	0x90, 0x6e, 0xb7, 0xb7, 0x27, 0xf5, 0x25, 0x3d, 0x22, 0x03, 0x53, 0x2a, 0xd0, 0x28, 0xe7, 0xc0,
	0x04, 0x66, 0xf2, 0x13, 0x30, 0x5d, 0x8f, 0x11, 0xde, 0x70, 0xf4, 0xbd, 0x73, 0x74, 0xe2, 0xc0,
	0xb4, 0xc6, 0xe8, 0x13, 0x42, 0x8d, 0xe5, 0x67, 0x02, 0x73, 0xa6, 0x40, 0x67, 0x80, 0x96, 0xe7,
	0x10, 0xde, 0x74, 0x15, 0xbb, 0x1e, 0x39, 0x39, 0x07, 0xe8, 0x31, 0x79, 0xd0, 0xd2, 0xe1, 0x8b,
	0x82, 0xcc, 0xc2, 0x9c, 0x69, 0x58, 0x00, 0x96, 0xb0, 0x5c, 0x7f, 0xcb, 0xd5, 0x0f, 0x3d, 0xf5,
	0xb5, 0x67, 0xa6, 0x0d, 0xf1, 0xa2, 0x5d, 0xfc, 0x37, 0x20, 0xfb, 0xc7, 0x72, 0x2e, 0x4e, 0xbf,
	0x5e, 0x32, 0x8f, 0x4e, 0xc8, 0x55, 0x27, 0x9c, 0x45, 0xdd, 0xc3, 0xbb, 0x49, 0xb3, 0x8e, 0xa4,
	0x5d, 0x47, 0x32, 0x41, 0xfb, 0xfc, 0xe8, 0x1d, 0x2f, 0x4a, 0x58, 0xe7, 0xdf, 0x4b, 0xd2, 0x6d,
	0x9d, 0x7a, 0x03, 0x8d, 0x79, 0x1b, 0x9a, 0x2c, 0xf3, 0xe9, 0x07, 0x12, 0x59, 0xcd, 0xd1, 0xf0,
	0xac, 0x8e, 0xc4, 0x85, 0xf8, 0xe9, 0x47, 0xae, 0xc1, 0x79, 0xbb, 0xa1, 0xdb, 0x7f, 0xca, 0xe3,
	0xa7, 0x64, 0x70, 0x79, 0x7e, 0x9f, 0x92, 0x01, 0xd9, 0x32, 0x96, 0xdb, 0xd2, 0xb8, 0xa9, 0x6f,
	0xa7, 0xfe, 0x74, 0xf8, 0xa3, 0x43, 0x76, 0x5a, 0xf2, 0x14, 0xf4, 0xa2, 0x5e, 0x7d, 0x45, 0xba,
	0x4b, 0x41, 0xa3, 0x8f, 0x92, 0x35, 0x11, 0x4f, 0xae, 0x26, 0x34, 0x1a, 0x6d, 0x26, 0x36, 0x6a,
	0xe2, 0xfd, 0x6f, 0xbf, 0xff, 0xfc, 0xec, 0xec, 0xd0, 0x6d, 0xf7, 0x42, 0xda, 0x17, 0x45, 0xbf,
	0x07, 0xa4, 0xb7, 0xaa, 0x9f, 0x1e, 0xac, 0xed, 0xb9, 0x76, 0xc9, 0xd1, 0xe3, 0x6b, 0x71, 0xbd,
	0x84, 0xd0, 0x49, 0xa0, 0xd1, 0xaa, 0x84, 0x17, 0xc1, 0xc1, 0x6c, 0xcb, 0xb9, 0xfe, 0xec, 0x5f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x1a, 0xf5, 0x7f, 0xc0, 0xea, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SettingsServiceClient is the client API for SettingsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SettingsServiceClient interface {
	GetSettings(ctx context.Context, in *GetSettingsRequest, opts ...grpc.CallOption) (*GetSettingsResponse, error)
	ModifySettings(ctx context.Context, in *ModifySettingsRequest, opts ...grpc.CallOption) (*ModifySettingsResponse, error)
}

type settingsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSettingsServiceClient(cc grpc.ClientConnInterface) SettingsServiceClient {
	return &settingsServiceClient{cc}
}

func (c *settingsServiceClient) GetSettings(ctx context.Context, in *GetSettingsRequest, opts ...grpc.CallOption) (*GetSettingsResponse, error) {
	out := new(GetSettingsResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.SettingsService/GetSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingsServiceClient) ModifySettings(ctx context.Context, in *ModifySettingsRequest, opts ...grpc.CallOption) (*ModifySettingsResponse, error) {
	out := new(ModifySettingsResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.SettingsService/ModifySettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SettingsServiceServer is the server API for SettingsService service.
type SettingsServiceServer interface {
	GetSettings(context.Context, *GetSettingsRequest) (*GetSettingsResponse, error)
	ModifySettings(context.Context, *ModifySettingsRequest) (*ModifySettingsResponse, error)
}

// UnimplementedSettingsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSettingsServiceServer struct {
}

func (*UnimplementedSettingsServiceServer) GetSettings(ctx context.Context, req *GetSettingsRequest) (*GetSettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSettings not implemented")
}
func (*UnimplementedSettingsServiceServer) ModifySettings(ctx context.Context, req *ModifySettingsRequest) (*ModifySettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifySettings not implemented")
}

func RegisterSettingsServiceServer(s *grpc.Server, srv SettingsServiceServer) {
	s.RegisterService(&_SettingsService_serviceDesc, srv)
}

func _SettingsService_GetSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingsServiceServer).GetSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.SettingsService/GetSettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingsServiceServer).GetSettings(ctx, req.(*GetSettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SettingsService_ModifySettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifySettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingsServiceServer).ModifySettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.SettingsService/ModifySettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingsServiceServer).ModifySettings(ctx, req.(*ModifySettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SettingsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "appserver_serves_ui.SettingsService",
	HandlerType: (*SettingsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSettings",
			Handler:    _SettingsService_GetSettings_Handler,
		},
		{
			MethodName: "ModifySettings",
			Handler:    _SettingsService_ModifySettings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "settings.proto",
}
