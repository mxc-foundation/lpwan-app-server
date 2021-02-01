// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.2
// source: m2m_distribute_bonus.proto

package m2m_serves_appserver

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type AddBonusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrgId       int64  `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Currency    string `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`
	Amount      string `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Discription string `protobuf:"bytes,4,opt,name=discription,proto3" json:"discription,omitempty"`
}

func (x *AddBonusRequest) Reset() {
	*x = AddBonusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_m2m_distribute_bonus_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddBonusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddBonusRequest) ProtoMessage() {}

func (x *AddBonusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_m2m_distribute_bonus_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddBonusRequest.ProtoReflect.Descriptor instead.
func (*AddBonusRequest) Descriptor() ([]byte, []int) {
	return file_m2m_distribute_bonus_proto_rawDescGZIP(), []int{0}
}

func (x *AddBonusRequest) GetOrgId() int64 {
	if x != nil {
		return x.OrgId
	}
	return 0
}

func (x *AddBonusRequest) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *AddBonusRequest) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *AddBonusRequest) GetDiscription() string {
	if x != nil {
		return x.Discription
	}
	return ""
}

type AddBonusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BonusId int64 `protobuf:"varint,1,opt,name=bonus_id,json=bonusId,proto3" json:"bonus_id,omitempty"`
}

func (x *AddBonusResponse) Reset() {
	*x = AddBonusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_m2m_distribute_bonus_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddBonusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddBonusResponse) ProtoMessage() {}

func (x *AddBonusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_m2m_distribute_bonus_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddBonusResponse.ProtoReflect.Descriptor instead.
func (*AddBonusResponse) Descriptor() ([]byte, []int) {
	return file_m2m_distribute_bonus_proto_rawDescGZIP(), []int{1}
}

func (x *AddBonusResponse) GetBonusId() int64 {
	if x != nil {
		return x.BonusId
	}
	return 0
}

var File_m2m_distribute_bonus_proto protoreflect.FileDescriptor

var file_m2m_distribute_bonus_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x6d, 0x32, 0x6d, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x5f, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x6d, 0x32,
	0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x7e, 0x0a, 0x0f, 0x41, 0x64, 0x64, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x69, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x2d, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x6f, 0x6e, 0x75, 0x73,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x62, 0x6f, 0x6e, 0x75, 0x73,
	0x49, 0x64, 0x32, 0x73, 0x0a, 0x16, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x42, 0x6f, 0x6e, 0x75, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x59, 0x0a, 0x08,
	0x41, 0x64, 0x64, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x12, 0x25, 0x2e, 0x6d, 0x32, 0x6d, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x41, 0x64, 0x64, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x26, 0x2e, 0x6d, 0x32, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x61, 0x70, 0x70,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x41, 0x64, 0x64, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x78, 0x63, 0x2d, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6d, 0x78, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2d,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x32, 0x6d, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_m2m_distribute_bonus_proto_rawDescOnce sync.Once
	file_m2m_distribute_bonus_proto_rawDescData = file_m2m_distribute_bonus_proto_rawDesc
)

func file_m2m_distribute_bonus_proto_rawDescGZIP() []byte {
	file_m2m_distribute_bonus_proto_rawDescOnce.Do(func() {
		file_m2m_distribute_bonus_proto_rawDescData = protoimpl.X.CompressGZIP(file_m2m_distribute_bonus_proto_rawDescData)
	})
	return file_m2m_distribute_bonus_proto_rawDescData
}

var file_m2m_distribute_bonus_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_m2m_distribute_bonus_proto_goTypes = []interface{}{
	(*AddBonusRequest)(nil),  // 0: m2m_serves_appserver.AddBonusRequest
	(*AddBonusResponse)(nil), // 1: m2m_serves_appserver.AddBonusResponse
}
var file_m2m_distribute_bonus_proto_depIdxs = []int32{
	0, // 0: m2m_serves_appserver.DistributeBonusService.AddBonus:input_type -> m2m_serves_appserver.AddBonusRequest
	1, // 1: m2m_serves_appserver.DistributeBonusService.AddBonus:output_type -> m2m_serves_appserver.AddBonusResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_m2m_distribute_bonus_proto_init() }
func file_m2m_distribute_bonus_proto_init() {
	if File_m2m_distribute_bonus_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_m2m_distribute_bonus_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddBonusRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_m2m_distribute_bonus_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddBonusResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_m2m_distribute_bonus_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_m2m_distribute_bonus_proto_goTypes,
		DependencyIndexes: file_m2m_distribute_bonus_proto_depIdxs,
		MessageInfos:      file_m2m_distribute_bonus_proto_msgTypes,
	}.Build()
	File_m2m_distribute_bonus_proto = out.File
	file_m2m_distribute_bonus_proto_rawDesc = nil
	file_m2m_distribute_bonus_proto_goTypes = nil
	file_m2m_distribute_bonus_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DistributeBonusServiceClient is the client API for DistributeBonusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DistributeBonusServiceClient interface {
	// AddBonus add new bonus record
	AddBonus(ctx context.Context, in *AddBonusRequest, opts ...grpc.CallOption) (*AddBonusResponse, error)
}

type distributeBonusServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDistributeBonusServiceClient(cc grpc.ClientConnInterface) DistributeBonusServiceClient {
	return &distributeBonusServiceClient{cc}
}

func (c *distributeBonusServiceClient) AddBonus(ctx context.Context, in *AddBonusRequest, opts ...grpc.CallOption) (*AddBonusResponse, error) {
	out := new(AddBonusResponse)
	err := c.cc.Invoke(ctx, "/m2m_serves_appserver.DistributeBonusService/AddBonus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DistributeBonusServiceServer is the server API for DistributeBonusService service.
type DistributeBonusServiceServer interface {
	// AddBonus add new bonus record
	AddBonus(context.Context, *AddBonusRequest) (*AddBonusResponse, error)
}

// UnimplementedDistributeBonusServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDistributeBonusServiceServer struct {
}

func (*UnimplementedDistributeBonusServiceServer) AddBonus(context.Context, *AddBonusRequest) (*AddBonusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBonus not implemented")
}

func RegisterDistributeBonusServiceServer(s *grpc.Server, srv DistributeBonusServiceServer) {
	s.RegisterService(&_DistributeBonusService_serviceDesc, srv)
}

func _DistributeBonusService_AddBonus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBonusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistributeBonusServiceServer).AddBonus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_serves_appserver.DistributeBonusService/AddBonus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistributeBonusServiceServer).AddBonus(ctx, req.(*AddBonusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DistributeBonusService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m2m_serves_appserver.DistributeBonusService",
	HandlerType: (*DistributeBonusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddBonus",
			Handler:    _DistributeBonusService_AddBonus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "m2m_distribute_bonus.proto",
}
