// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: serverInfo.proto

package appserver_serves_ui

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type ServerRegion int32

const (
	ServerRegion_NOT_DEFINED ServerRegion = 0
	ServerRegion_AVERAGE     ServerRegion = 1
	ServerRegion_RESTRICTED  ServerRegion = 2
)

// Enum value maps for ServerRegion.
var (
	ServerRegion_name = map[int32]string{
		0: "NOT_DEFINED",
		1: "AVERAGE",
		2: "RESTRICTED",
	}
	ServerRegion_value = map[string]int32{
		"NOT_DEFINED": 0,
		"AVERAGE":     1,
		"RESTRICTED":  2,
	}
)

func (x ServerRegion) Enum() *ServerRegion {
	p := new(ServerRegion)
	*p = x
	return p
}

func (x ServerRegion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServerRegion) Descriptor() protoreflect.EnumDescriptor {
	return file_serverInfo_proto_enumTypes[0].Descriptor()
}

func (ServerRegion) Type() protoreflect.EnumType {
	return &file_serverInfo_proto_enumTypes[0]
}

func (x ServerRegion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ServerRegion.Descriptor instead.
func (ServerRegion) EnumDescriptor() ([]byte, []int) {
	return file_serverInfo_proto_rawDescGZIP(), []int{0}
}

type GetAppserverVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetAppserverVersionResponse) Reset() {
	*x = GetAppserverVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_serverInfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAppserverVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAppserverVersionResponse) ProtoMessage() {}

func (x *GetAppserverVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_serverInfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAppserverVersionResponse.ProtoReflect.Descriptor instead.
func (*GetAppserverVersionResponse) Descriptor() ([]byte, []int) {
	return file_serverInfo_proto_rawDescGZIP(), []int{0}
}

func (x *GetAppserverVersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type GetServerRegionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerRegion string `protobuf:"bytes,1,opt,name=server_region,json=serverRegion,proto3" json:"server_region,omitempty"`
}

func (x *GetServerRegionResponse) Reset() {
	*x = GetServerRegionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_serverInfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetServerRegionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetServerRegionResponse) ProtoMessage() {}

func (x *GetServerRegionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_serverInfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetServerRegionResponse.ProtoReflect.Descriptor instead.
func (*GetServerRegionResponse) Descriptor() ([]byte, []int) {
	return file_serverInfo_proto_rawDescGZIP(), []int{1}
}

func (x *GetServerRegionResponse) GetServerRegion() string {
	if x != nil {
		return x.ServerRegion
	}
	return ""
}

type GetMxprotocolServerVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetMxprotocolServerVersionResponse) Reset() {
	*x = GetMxprotocolServerVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_serverInfo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMxprotocolServerVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMxprotocolServerVersionResponse) ProtoMessage() {}

func (x *GetMxprotocolServerVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_serverInfo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMxprotocolServerVersionResponse.ProtoReflect.Descriptor instead.
func (*GetMxprotocolServerVersionResponse) Descriptor() ([]byte, []int) {
	return file_serverInfo_proto_rawDescGZIP(), []int{2}
}

func (x *GetMxprotocolServerVersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_serverInfo_proto protoreflect.FileDescriptor

var file_serverInfo_proto_rawDesc = []byte{
	0x0a, 0x10, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x13, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x73, 0x5f, 0x75, 0x69, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x37, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3e, 0x0a, 0x17, 0x47,
	0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x5f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x3e, 0x0a, 0x22, 0x47,
	0x65, 0x74, 0x4d, 0x78, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2a, 0x3c, 0x0a, 0x0c, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x0f, 0x0a, 0x0b, 0x4e,
	0x4f, 0x54, 0x5f, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07,
	0x41, 0x56, 0x45, 0x52, 0x41, 0x47, 0x45, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x45, 0x53,
	0x54, 0x52, 0x49, 0x43, 0x54, 0x45, 0x44, 0x10, 0x02, 0x32, 0xc6, 0x03, 0x0a, 0x11, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x8b, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x30, 0x2e, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x73, 0x5f, 0x75, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x12, 0x22, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x61, 0x70, 0x70, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x7f, 0x0a,
	0x0f, 0x47, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x70, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x75, 0x69, 0x2e, 0x47,
	0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x26, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x20, 0x12, 0x1e,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x69, 0x6e, 0x66, 0x6f,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0xa1,
	0x01, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x4d, 0x78, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x37, 0x2e, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x75, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4d,
	0x78, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2c, 0x12, 0x2a, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2d, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x6d, 0x78, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x42, 0x58, 0x5a, 0x56, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6d, 0x78, 0x63, 0x2d, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x6c, 0x70, 0x77, 0x61, 0x6e, 0x2d, 0x61, 0x70, 0x70, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x73, 0x2d, 0x75, 0x69, 0x3b, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x75, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_serverInfo_proto_rawDescOnce sync.Once
	file_serverInfo_proto_rawDescData = file_serverInfo_proto_rawDesc
)

func file_serverInfo_proto_rawDescGZIP() []byte {
	file_serverInfo_proto_rawDescOnce.Do(func() {
		file_serverInfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_serverInfo_proto_rawDescData)
	})
	return file_serverInfo_proto_rawDescData
}

var file_serverInfo_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_serverInfo_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_serverInfo_proto_goTypes = []interface{}{
	(ServerRegion)(0),                          // 0: appserver_serves_ui.ServerRegion
	(*GetAppserverVersionResponse)(nil),        // 1: appserver_serves_ui.GetAppserverVersionResponse
	(*GetServerRegionResponse)(nil),            // 2: appserver_serves_ui.GetServerRegionResponse
	(*GetMxprotocolServerVersionResponse)(nil), // 3: appserver_serves_ui.GetMxprotocolServerVersionResponse
	(*empty.Empty)(nil),                        // 4: google.protobuf.Empty
}
var file_serverInfo_proto_depIdxs = []int32{
	4, // 0: appserver_serves_ui.ServerInfoService.GetAppserverVersion:input_type -> google.protobuf.Empty
	4, // 1: appserver_serves_ui.ServerInfoService.GetServerRegion:input_type -> google.protobuf.Empty
	4, // 2: appserver_serves_ui.ServerInfoService.GetMxprotocolServerVersion:input_type -> google.protobuf.Empty
	1, // 3: appserver_serves_ui.ServerInfoService.GetAppserverVersion:output_type -> appserver_serves_ui.GetAppserverVersionResponse
	2, // 4: appserver_serves_ui.ServerInfoService.GetServerRegion:output_type -> appserver_serves_ui.GetServerRegionResponse
	3, // 5: appserver_serves_ui.ServerInfoService.GetMxprotocolServerVersion:output_type -> appserver_serves_ui.GetMxprotocolServerVersionResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_serverInfo_proto_init() }
func file_serverInfo_proto_init() {
	if File_serverInfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_serverInfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAppserverVersionResponse); i {
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
		file_serverInfo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetServerRegionResponse); i {
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
		file_serverInfo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMxprotocolServerVersionResponse); i {
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
			RawDescriptor: file_serverInfo_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_serverInfo_proto_goTypes,
		DependencyIndexes: file_serverInfo_proto_depIdxs,
		EnumInfos:         file_serverInfo_proto_enumTypes,
		MessageInfos:      file_serverInfo_proto_msgTypes,
	}.Build()
	File_serverInfo_proto = out.File
	file_serverInfo_proto_rawDesc = nil
	file_serverInfo_proto_goTypes = nil
	file_serverInfo_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ServerInfoServiceClient is the client API for ServerInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServerInfoServiceClient interface {
	// get version
	GetAppserverVersion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetAppserverVersionResponse, error)
	GetServerRegion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetServerRegionResponse, error)
	GetMxprotocolServerVersion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetMxprotocolServerVersionResponse, error)
}

type serverInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewServerInfoServiceClient(cc grpc.ClientConnInterface) ServerInfoServiceClient {
	return &serverInfoServiceClient{cc}
}

func (c *serverInfoServiceClient) GetAppserverVersion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetAppserverVersionResponse, error) {
	out := new(GetAppserverVersionResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.ServerInfoService/GetAppserverVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverInfoServiceClient) GetServerRegion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetServerRegionResponse, error) {
	out := new(GetServerRegionResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.ServerInfoService/GetServerRegion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverInfoServiceClient) GetMxprotocolServerVersion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetMxprotocolServerVersionResponse, error) {
	out := new(GetMxprotocolServerVersionResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.ServerInfoService/GetMxprotocolServerVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerInfoServiceServer is the server API for ServerInfoService service.
type ServerInfoServiceServer interface {
	// get version
	GetAppserverVersion(context.Context, *empty.Empty) (*GetAppserverVersionResponse, error)
	GetServerRegion(context.Context, *empty.Empty) (*GetServerRegionResponse, error)
	GetMxprotocolServerVersion(context.Context, *empty.Empty) (*GetMxprotocolServerVersionResponse, error)
}

// UnimplementedServerInfoServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServerInfoServiceServer struct {
}

func (*UnimplementedServerInfoServiceServer) GetAppserverVersion(context.Context, *empty.Empty) (*GetAppserverVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAppserverVersion not implemented")
}
func (*UnimplementedServerInfoServiceServer) GetServerRegion(context.Context, *empty.Empty) (*GetServerRegionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServerRegion not implemented")
}
func (*UnimplementedServerInfoServiceServer) GetMxprotocolServerVersion(context.Context, *empty.Empty) (*GetMxprotocolServerVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMxprotocolServerVersion not implemented")
}

func RegisterServerInfoServiceServer(s *grpc.Server, srv ServerInfoServiceServer) {
	s.RegisterService(&_ServerInfoService_serviceDesc, srv)
}

func _ServerInfoService_GetAppserverVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerInfoServiceServer).GetAppserverVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.ServerInfoService/GetAppserverVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerInfoServiceServer).GetAppserverVersion(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerInfoService_GetServerRegion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerInfoServiceServer).GetServerRegion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.ServerInfoService/GetServerRegion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerInfoServiceServer).GetServerRegion(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerInfoService_GetMxprotocolServerVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerInfoServiceServer).GetMxprotocolServerVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.ServerInfoService/GetMxprotocolServerVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerInfoServiceServer).GetMxprotocolServerVersion(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServerInfoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "appserver_serves_ui.ServerInfoService",
	HandlerType: (*ServerInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAppserverVersion",
			Handler:    _ServerInfoService_GetAppserverVersion_Handler,
		},
		{
			MethodName: "GetServerRegion",
			Handler:    _ServerInfoService_GetServerRegion_Handler,
		},
		{
			MethodName: "GetMxprotocolServerVersion",
			Handler:    _ServerInfoService_GetMxprotocolServerVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "serverInfo.proto",
}
