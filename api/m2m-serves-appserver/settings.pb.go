// Code generated by protoc-gen-go. DO NOT EDIT.
// source: settings.proto

package m2m_serves_appserver

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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
	Compensation                     float64  `protobuf:"fixed64,6,opt,name=compensation,proto3" json:"compensation,omitempty"`
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

func (m *GetSettingsResponse) GetCompensation() float64 {
	if m != nil {
		return m.Compensation
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
	proto.RegisterType((*GetSettingsRequest)(nil), "m2m_serves_appserver.GetSettingsRequest")
	proto.RegisterType((*GetSettingsResponse)(nil), "m2m_serves_appserver.GetSettingsResponse")
	proto.RegisterType((*ModifySettingsRequest)(nil), "m2m_serves_appserver.ModifySettingsRequest")
	proto.RegisterType((*ModifySettingsResponse)(nil), "m2m_serves_appserver.ModifySettingsResponse")
}

func init() { proto.RegisterFile("settings.proto", fileDescriptor_6c7cab62fa432213) }

var fileDescriptor_6c7cab62fa432213 = []byte{
	// 445 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x51, 0x8b, 0xd3, 0x40,
	0x10, 0x26, 0xa7, 0x16, 0x9d, 0x6a, 0xe5, 0xf6, 0x6a, 0x09, 0x15, 0xa4, 0x44, 0x84, 0x8a, 0x9a,
	0x93, 0x7a, 0xf8, 0xe6, 0x8b, 0xa0, 0xd2, 0x87, 0x83, 0x23, 0x05, 0x7d, 0xf0, 0x61, 0xd9, 0xa6,
	0x73, 0x31, 0x34, 0x9d, 0x5d, 0x77, 0x37, 0x8d, 0xfe, 0x3f, 0x7f, 0x89, 0x3f, 0x44, 0x24, 0x9b,
	0xcd, 0xd9, 0xde, 0x05, 0xef, 0xde, 0x92, 0xf9, 0xbe, 0x6f, 0x76, 0xe6, 0x9b, 0x0f, 0x06, 0x06,
	0xad, 0xcd, 0x29, 0x33, 0xb1, 0xd2, 0xd2, 0x4a, 0x36, 0xdc, 0xcc, 0x36, 0xdc, 0xa0, 0xde, 0xa2,
	0xe1, 0x42, 0x29, 0xf7, 0xa5, 0xc7, 0x4f, 0x32, 0x29, 0xb3, 0x02, 0x8f, 0x1d, 0x67, 0x59, 0x9e,
	0x1f, 0x57, 0x5a, 0x28, 0x85, 0xda, 0xab, 0xa2, 0x21, 0xb0, 0x4f, 0x68, 0x17, 0xbe, 0x55, 0x82,
	0xdf, 0x4b, 0x34, 0x36, 0xfa, 0x75, 0x00, 0x47, 0x7b, 0x65, 0xa3, 0x24, 0x19, 0x64, 0x31, 0x1c,
	0x15, 0xb2, 0xe2, 0x4b, 0x51, 0x08, 0x4a, 0x91, 0x57, 0x42, 0x53, 0x4e, 0x59, 0x18, 0x4c, 0x82,
	0xe9, 0xbd, 0xe4, 0xb0, 0x90, 0xd5, 0xfb, 0x06, 0xf9, 0xd2, 0x00, 0xec, 0x19, 0x0c, 0x56, 0xb2,
	0xa2, 0x22, 0xa7, 0x35, 0x57, 0x3a, 0x4f, 0x31, 0x3c, 0x98, 0x04, 0xd3, 0x20, 0x79, 0xd0, 0x56,
	0xcf, 0xea, 0x22, 0x3b, 0x81, 0x91, 0x29, 0x15, 0x6a, 0x92, 0x2b, 0xe4, 0x39, 0xa5, 0x72, 0x83,
	0x5c, 0x0b, 0x9b, 0xcb, 0xf0, 0x96, 0xa3, 0x0f, 0x2f, 0xd0, 0xb9, 0x03, 0x93, 0x1a, 0x63, 0xaf,
	0x80, 0x19, 0x2b, 0xd6, 0x39, 0x65, 0x5c, 0xa1, 0x4e, 0x91, 0xac, 0xc8, 0x30, 0xbc, 0xed, 0x14,
	0x87, 0x1e, 0x39, 0xbb, 0x00, 0xd8, 0x29, 0x3c, 0x6d, 0xe9, 0xf8, 0x43, 0x61, 0x6a, 0x71, 0xc5,
	0x35, 0x6e, 0x91, 0x4a, 0xdc, 0xd5, 0xdf, 0x71, 0xfa, 0x89, 0xa7, 0x7e, 0xf0, 0xcc, 0xa4, 0x21,
	0xee, 0xb4, 0x8b, 0xe0, 0x7e, 0x2a, 0x37, 0x0a, 0xc9, 0xd4, 0xc3, 0x50, 0xd8, 0x73, 0xba, 0xbd,
	0x5a, 0xf4, 0x27, 0x80, 0x47, 0xa7, 0x72, 0x95, 0x9f, 0xff, 0xbc, 0x64, 0x30, 0x9b, 0xc3, 0x55,
	0xb7, 0x9c, 0x8d, 0xfd, 0xd9, 0xe3, 0xb8, 0x39, 0x59, 0xdc, 0x9e, 0x2c, 0x9e, 0x93, 0x7d, 0x7b,
	0xf2, 0x59, 0x14, 0x25, 0x76, 0x79, 0xfc, 0x0e, 0xfa, 0xad, 0x9b, 0x1f, 0xb1, 0x31, 0xf8, 0x9a,
	0x26, 0xbb, 0x7c, 0xf6, 0x15, 0xc6, 0x56, 0x0b, 0x32, 0x22, 0xad, 0x47, 0xfe, 0xb7, 0xe0, 0xe2,
	0x9b, 0xd0, 0xe8, 0xfc, 0xbf, 0xa6, 0xdb, 0x7f, 0xe4, 0xd1, 0x6b, 0x18, 0x5d, 0xde, 0xdf, 0x27,
	0x69, 0x04, 0x3d, 0x63, 0x85, 0x2d, 0x8d, 0xdb, 0xfa, 0x6e, 0xe2, 0xff, 0x66, 0xbf, 0x03, 0x78,
	0xd8, 0x92, 0x17, 0xa8, 0xb7, 0x75, 0x3c, 0x96, 0xd0, 0xdf, 0x09, 0x23, 0x9b, 0xc6, 0x5d, 0x49,
	0x8f, 0xaf, 0xc6, 0x78, 0xfc, 0xfc, 0x06, 0x4c, 0x3f, 0xcf, 0x1a, 0x06, 0xfb, 0x93, 0xb2, 0x17,
	0xdd, 0xe2, 0xce, 0x7b, 0x8e, 0x5f, 0xde, 0x8c, 0xdc, 0x3c, 0xb6, 0xec, 0x39, 0x1f, 0xdf, 0xfc,
	0x0d, 0x00, 0x00, 0xff, 0xff, 0x78, 0xff, 0xda, 0x77, 0xc3, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SettingsServiceClient is the client API for SettingsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SettingsServiceClient interface {
	GetSettings(ctx context.Context, in *GetSettingsRequest, opts ...grpc.CallOption) (*GetSettingsResponse, error)
	ModifySettings(ctx context.Context, in *ModifySettingsRequest, opts ...grpc.CallOption) (*ModifySettingsResponse, error)
}

type settingsServiceClient struct {
	cc *grpc.ClientConn
}

func NewSettingsServiceClient(cc *grpc.ClientConn) SettingsServiceClient {
	return &settingsServiceClient{cc}
}

func (c *settingsServiceClient) GetSettings(ctx context.Context, in *GetSettingsRequest, opts ...grpc.CallOption) (*GetSettingsResponse, error) {
	out := new(GetSettingsResponse)
	err := c.cc.Invoke(ctx, "/m2m_serves_appserver.SettingsService/GetSettings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingsServiceClient) ModifySettings(ctx context.Context, in *ModifySettingsRequest, opts ...grpc.CallOption) (*ModifySettingsResponse, error) {
	out := new(ModifySettingsResponse)
	err := c.cc.Invoke(ctx, "/m2m_serves_appserver.SettingsService/ModifySettings", in, out, opts...)
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
		FullMethod: "/m2m_serves_appserver.SettingsService/GetSettings",
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
		FullMethod: "/m2m_serves_appserver.SettingsService/ModifySettings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingsServiceServer).ModifySettings(ctx, req.(*ModifySettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SettingsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m2m_serves_appserver.SettingsService",
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
