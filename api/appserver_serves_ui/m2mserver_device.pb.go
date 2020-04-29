// Code generated by protoc-gen-go. DO NOT EDIT.
// source: m2mserver_device.proto

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

type DeviceMode int32

const (
	DeviceMode_DV_INACTIVE              DeviceMode = 0
	DeviceMode_DV_FREE_GATEWAYS_LIMITED DeviceMode = 1
	DeviceMode_DV_WHOLE_NETWORK         DeviceMode = 2
	DeviceMode_DV_DELETED               DeviceMode = 3
)

var DeviceMode_name = map[int32]string{
	0: "DV_INACTIVE",
	1: "DV_FREE_GATEWAYS_LIMITED",
	2: "DV_WHOLE_NETWORK",
	3: "DV_DELETED",
}

var DeviceMode_value = map[string]int32{
	"DV_INACTIVE":              0,
	"DV_FREE_GATEWAYS_LIMITED": 1,
	"DV_WHOLE_NETWORK":         2,
	"DV_DELETED":               3,
}

func (x DeviceMode) String() string {
	return proto.EnumName(DeviceMode_name, int32(x))
}

func (DeviceMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{0}
}

type GetDeviceListRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int64    `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDeviceListRequest) Reset()         { *m = GetDeviceListRequest{} }
func (m *GetDeviceListRequest) String() string { return proto.CompactTextString(m) }
func (*GetDeviceListRequest) ProtoMessage()    {}
func (*GetDeviceListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{0}
}

func (m *GetDeviceListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDeviceListRequest.Unmarshal(m, b)
}
func (m *GetDeviceListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDeviceListRequest.Marshal(b, m, deterministic)
}
func (m *GetDeviceListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDeviceListRequest.Merge(m, src)
}
func (m *GetDeviceListRequest) XXX_Size() int {
	return xxx_messageInfo_GetDeviceListRequest.Size(m)
}
func (m *GetDeviceListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDeviceListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetDeviceListRequest proto.InternalMessageInfo

func (m *GetDeviceListRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

func (m *GetDeviceListRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *GetDeviceListRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type GetDeviceListResponse struct {
	DevProfile           []*DSDeviceProfile `protobuf:"bytes,1,rep,name=dev_profile,json=devProfile,proto3" json:"dev_profile,omitempty"`
	Count                int64              `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	UserProfile          *ProfileResponse   `protobuf:"bytes,3,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetDeviceListResponse) Reset()         { *m = GetDeviceListResponse{} }
func (m *GetDeviceListResponse) String() string { return proto.CompactTextString(m) }
func (*GetDeviceListResponse) ProtoMessage()    {}
func (*GetDeviceListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{1}
}

func (m *GetDeviceListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDeviceListResponse.Unmarshal(m, b)
}
func (m *GetDeviceListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDeviceListResponse.Marshal(b, m, deterministic)
}
func (m *GetDeviceListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDeviceListResponse.Merge(m, src)
}
func (m *GetDeviceListResponse) XXX_Size() int {
	return xxx_messageInfo_GetDeviceListResponse.Size(m)
}
func (m *GetDeviceListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDeviceListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetDeviceListResponse proto.InternalMessageInfo

func (m *GetDeviceListResponse) GetDevProfile() []*DSDeviceProfile {
	if m != nil {
		return m.DevProfile
	}
	return nil
}

func (m *GetDeviceListResponse) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *GetDeviceListResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

type GetDSDeviceProfileRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	DevId                int64    `protobuf:"varint,2,opt,name=dev_id,json=devId,proto3" json:"dev_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDSDeviceProfileRequest) Reset()         { *m = GetDSDeviceProfileRequest{} }
func (m *GetDSDeviceProfileRequest) String() string { return proto.CompactTextString(m) }
func (*GetDSDeviceProfileRequest) ProtoMessage()    {}
func (*GetDSDeviceProfileRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{2}
}

func (m *GetDSDeviceProfileRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDSDeviceProfileRequest.Unmarshal(m, b)
}
func (m *GetDSDeviceProfileRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDSDeviceProfileRequest.Marshal(b, m, deterministic)
}
func (m *GetDSDeviceProfileRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDSDeviceProfileRequest.Merge(m, src)
}
func (m *GetDSDeviceProfileRequest) XXX_Size() int {
	return xxx_messageInfo_GetDSDeviceProfileRequest.Size(m)
}
func (m *GetDSDeviceProfileRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDSDeviceProfileRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetDSDeviceProfileRequest proto.InternalMessageInfo

func (m *GetDSDeviceProfileRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

func (m *GetDSDeviceProfileRequest) GetDevId() int64 {
	if m != nil {
		return m.DevId
	}
	return 0
}

type DSDeviceProfile struct {
	Id                   int64      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DevEui               string     `protobuf:"bytes,2,opt,name=dev_eui,json=devEui,proto3" json:"dev_eui,omitempty"`
	FkWallet             int64      `protobuf:"varint,3,opt,name=fk_wallet,json=fkWallet,proto3" json:"fk_wallet,omitempty"`
	Mode                 DeviceMode `protobuf:"varint,4,opt,name=mode,proto3,enum=appserver_serves_ui.DeviceMode" json:"mode,omitempty"`
	CreatedAt            string     `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	LastSeenAt           string     `protobuf:"bytes,6,opt,name=last_seen_at,json=lastSeenAt,proto3" json:"last_seen_at,omitempty"`
	ApplicationId        int64      `protobuf:"varint,7,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	Name                 string     `protobuf:"bytes,8,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *DSDeviceProfile) Reset()         { *m = DSDeviceProfile{} }
func (m *DSDeviceProfile) String() string { return proto.CompactTextString(m) }
func (*DSDeviceProfile) ProtoMessage()    {}
func (*DSDeviceProfile) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{3}
}

func (m *DSDeviceProfile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DSDeviceProfile.Unmarshal(m, b)
}
func (m *DSDeviceProfile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DSDeviceProfile.Marshal(b, m, deterministic)
}
func (m *DSDeviceProfile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DSDeviceProfile.Merge(m, src)
}
func (m *DSDeviceProfile) XXX_Size() int {
	return xxx_messageInfo_DSDeviceProfile.Size(m)
}
func (m *DSDeviceProfile) XXX_DiscardUnknown() {
	xxx_messageInfo_DSDeviceProfile.DiscardUnknown(m)
}

var xxx_messageInfo_DSDeviceProfile proto.InternalMessageInfo

func (m *DSDeviceProfile) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *DSDeviceProfile) GetDevEui() string {
	if m != nil {
		return m.DevEui
	}
	return ""
}

func (m *DSDeviceProfile) GetFkWallet() int64 {
	if m != nil {
		return m.FkWallet
	}
	return 0
}

func (m *DSDeviceProfile) GetMode() DeviceMode {
	if m != nil {
		return m.Mode
	}
	return DeviceMode_DV_INACTIVE
}

func (m *DSDeviceProfile) GetCreatedAt() string {
	if m != nil {
		return m.CreatedAt
	}
	return ""
}

func (m *DSDeviceProfile) GetLastSeenAt() string {
	if m != nil {
		return m.LastSeenAt
	}
	return ""
}

func (m *DSDeviceProfile) GetApplicationId() int64 {
	if m != nil {
		return m.ApplicationId
	}
	return 0
}

func (m *DSDeviceProfile) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GetDSDeviceProfileResponse struct {
	DevProfile           *DSDeviceProfile `protobuf:"bytes,1,opt,name=dev_profile,json=devProfile,proto3" json:"dev_profile,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,3,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetDSDeviceProfileResponse) Reset()         { *m = GetDSDeviceProfileResponse{} }
func (m *GetDSDeviceProfileResponse) String() string { return proto.CompactTextString(m) }
func (*GetDSDeviceProfileResponse) ProtoMessage()    {}
func (*GetDSDeviceProfileResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{4}
}

func (m *GetDSDeviceProfileResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDSDeviceProfileResponse.Unmarshal(m, b)
}
func (m *GetDSDeviceProfileResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDSDeviceProfileResponse.Marshal(b, m, deterministic)
}
func (m *GetDSDeviceProfileResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDSDeviceProfileResponse.Merge(m, src)
}
func (m *GetDSDeviceProfileResponse) XXX_Size() int {
	return xxx_messageInfo_GetDSDeviceProfileResponse.Size(m)
}
func (m *GetDSDeviceProfileResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDSDeviceProfileResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetDSDeviceProfileResponse proto.InternalMessageInfo

func (m *GetDSDeviceProfileResponse) GetDevProfile() *DSDeviceProfile {
	if m != nil {
		return m.DevProfile
	}
	return nil
}

func (m *GetDSDeviceProfileResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

type GetDeviceHistoryRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	DevId                int64    `protobuf:"varint,2,opt,name=dev_id,json=devId,proto3" json:"dev_id,omitempty"`
	Offset               int64    `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int64    `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDeviceHistoryRequest) Reset()         { *m = GetDeviceHistoryRequest{} }
func (m *GetDeviceHistoryRequest) String() string { return proto.CompactTextString(m) }
func (*GetDeviceHistoryRequest) ProtoMessage()    {}
func (*GetDeviceHistoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{5}
}

func (m *GetDeviceHistoryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDeviceHistoryRequest.Unmarshal(m, b)
}
func (m *GetDeviceHistoryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDeviceHistoryRequest.Marshal(b, m, deterministic)
}
func (m *GetDeviceHistoryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDeviceHistoryRequest.Merge(m, src)
}
func (m *GetDeviceHistoryRequest) XXX_Size() int {
	return xxx_messageInfo_GetDeviceHistoryRequest.Size(m)
}
func (m *GetDeviceHistoryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDeviceHistoryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetDeviceHistoryRequest proto.InternalMessageInfo

func (m *GetDeviceHistoryRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

func (m *GetDeviceHistoryRequest) GetDevId() int64 {
	if m != nil {
		return m.DevId
	}
	return 0
}

func (m *GetDeviceHistoryRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *GetDeviceHistoryRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type GetDeviceHistoryResponse struct {
	DevHistory           string           `protobuf:"bytes,1,opt,name=dev_history,json=devHistory,proto3" json:"dev_history,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,2,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetDeviceHistoryResponse) Reset()         { *m = GetDeviceHistoryResponse{} }
func (m *GetDeviceHistoryResponse) String() string { return proto.CompactTextString(m) }
func (*GetDeviceHistoryResponse) ProtoMessage()    {}
func (*GetDeviceHistoryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{6}
}

func (m *GetDeviceHistoryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDeviceHistoryResponse.Unmarshal(m, b)
}
func (m *GetDeviceHistoryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDeviceHistoryResponse.Marshal(b, m, deterministic)
}
func (m *GetDeviceHistoryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDeviceHistoryResponse.Merge(m, src)
}
func (m *GetDeviceHistoryResponse) XXX_Size() int {
	return xxx_messageInfo_GetDeviceHistoryResponse.Size(m)
}
func (m *GetDeviceHistoryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDeviceHistoryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetDeviceHistoryResponse proto.InternalMessageInfo

func (m *GetDeviceHistoryResponse) GetDevHistory() string {
	if m != nil {
		return m.DevHistory
	}
	return ""
}

func (m *GetDeviceHistoryResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

type SetDeviceModeRequest struct {
	OrgId                int64      `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	DevId                int64      `protobuf:"varint,2,opt,name=dev_id,json=devId,proto3" json:"dev_id,omitempty"`
	DevMode              DeviceMode `protobuf:"varint,3,opt,name=dev_mode,json=devMode,proto3,enum=appserver_serves_ui.DeviceMode" json:"dev_mode,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *SetDeviceModeRequest) Reset()         { *m = SetDeviceModeRequest{} }
func (m *SetDeviceModeRequest) String() string { return proto.CompactTextString(m) }
func (*SetDeviceModeRequest) ProtoMessage()    {}
func (*SetDeviceModeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{7}
}

func (m *SetDeviceModeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetDeviceModeRequest.Unmarshal(m, b)
}
func (m *SetDeviceModeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetDeviceModeRequest.Marshal(b, m, deterministic)
}
func (m *SetDeviceModeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetDeviceModeRequest.Merge(m, src)
}
func (m *SetDeviceModeRequest) XXX_Size() int {
	return xxx_messageInfo_SetDeviceModeRequest.Size(m)
}
func (m *SetDeviceModeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetDeviceModeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetDeviceModeRequest proto.InternalMessageInfo

func (m *SetDeviceModeRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

func (m *SetDeviceModeRequest) GetDevId() int64 {
	if m != nil {
		return m.DevId
	}
	return 0
}

func (m *SetDeviceModeRequest) GetDevMode() DeviceMode {
	if m != nil {
		return m.DevMode
	}
	return DeviceMode_DV_INACTIVE
}

type SetDeviceModeResponse struct {
	Status               bool             `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,3,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SetDeviceModeResponse) Reset()         { *m = SetDeviceModeResponse{} }
func (m *SetDeviceModeResponse) String() string { return proto.CompactTextString(m) }
func (*SetDeviceModeResponse) ProtoMessage()    {}
func (*SetDeviceModeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8cd7858e525a121, []int{8}
}

func (m *SetDeviceModeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetDeviceModeResponse.Unmarshal(m, b)
}
func (m *SetDeviceModeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetDeviceModeResponse.Marshal(b, m, deterministic)
}
func (m *SetDeviceModeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetDeviceModeResponse.Merge(m, src)
}
func (m *SetDeviceModeResponse) XXX_Size() int {
	return xxx_messageInfo_SetDeviceModeResponse.Size(m)
}
func (m *SetDeviceModeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetDeviceModeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetDeviceModeResponse proto.InternalMessageInfo

func (m *SetDeviceModeResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *SetDeviceModeResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

func init() {
	proto.RegisterEnum("appserver_serves_ui.DeviceMode", DeviceMode_name, DeviceMode_value)
	proto.RegisterType((*GetDeviceListRequest)(nil), "appserver_serves_ui.GetDeviceListRequest")
	proto.RegisterType((*GetDeviceListResponse)(nil), "appserver_serves_ui.GetDeviceListResponse")
	proto.RegisterType((*GetDSDeviceProfileRequest)(nil), "appserver_serves_ui.GetDSDeviceProfileRequest")
	proto.RegisterType((*DSDeviceProfile)(nil), "appserver_serves_ui.DSDeviceProfile")
	proto.RegisterType((*GetDSDeviceProfileResponse)(nil), "appserver_serves_ui.GetDSDeviceProfileResponse")
	proto.RegisterType((*GetDeviceHistoryRequest)(nil), "appserver_serves_ui.GetDeviceHistoryRequest")
	proto.RegisterType((*GetDeviceHistoryResponse)(nil), "appserver_serves_ui.GetDeviceHistoryResponse")
	proto.RegisterType((*SetDeviceModeRequest)(nil), "appserver_serves_ui.SetDeviceModeRequest")
	proto.RegisterType((*SetDeviceModeResponse)(nil), "appserver_serves_ui.SetDeviceModeResponse")
}

func init() { proto.RegisterFile("m2mserver_device.proto", fileDescriptor_d8cd7858e525a121) }

var fileDescriptor_d8cd7858e525a121 = []byte{
	// 781 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xcd, 0x6e, 0xea, 0x46,
	0x14, 0xae, 0xf9, 0x0b, 0x1c, 0x12, 0x82, 0xa6, 0x90, 0xb8, 0x34, 0x55, 0x90, 0xd5, 0x4a, 0x24,
	0x4a, 0x40, 0x22, 0x5d, 0x65, 0x87, 0x6a, 0x37, 0xb1, 0x4a, 0x92, 0xca, 0x20, 0x50, 0xd5, 0x85,
	0xe5, 0xe2, 0x81, 0x8e, 0x62, 0x6c, 0xd7, 0x1e, 0xd3, 0x56, 0x51, 0xa4, 0xaa, 0xea, 0xa2, 0xcb,
	0x4a, 0x5d, 0x74, 0xd3, 0x6d, 0x17, 0x95, 0xfa, 0x0a, 0x7d, 0x8a, 0xfb, 0x0a, 0xf7, 0x41, 0xae,
	0x66, 0x6c, 0x13, 0x42, 0x0c, 0xe2, 0xa2, 0xac, 0x98, 0x73, 0x74, 0x7e, 0xbe, 0x73, 0xe6, 0xfb,
	0xc6, 0xc0, 0xc1, 0xb4, 0x3d, 0xf5, 0xb1, 0x37, 0xc3, 0x9e, 0x6e, 0xe2, 0x19, 0x19, 0xe1, 0xa6,
	0xeb, 0x39, 0xd4, 0x41, 0x1f, 0x1a, 0xae, 0x1b, 0xf9, 0xf9, 0x8f, 0xaf, 0x07, 0xa4, 0x76, 0x34,
	0x71, 0x9c, 0x89, 0x85, 0x5b, 0x86, 0x4b, 0x5a, 0x86, 0x6d, 0x3b, 0xd4, 0xa0, 0xc4, 0xb1, 0xfd,
	0x30, 0xa5, 0x56, 0x22, 0x36, 0xc5, 0x9e, 0x6d, 0x58, 0xa1, 0x2d, 0x7d, 0x0b, 0x95, 0x2b, 0x4c,
	0x65, 0x5e, 0xb5, 0x4b, 0x7c, 0xaa, 0xe1, 0x1f, 0x02, 0xec, 0x53, 0x54, 0x85, 0x9c, 0xe3, 0x4d,
	0x74, 0x62, 0x8a, 0x42, 0x5d, 0x68, 0xa4, 0xb5, 0xac, 0xe3, 0x4d, 0x54, 0x13, 0x1d, 0x40, 0xce,
	0x19, 0x8f, 0x7d, 0x4c, 0xc5, 0x14, 0x77, 0x47, 0x16, 0xaa, 0x40, 0xd6, 0x22, 0x53, 0x42, 0xc5,
	0x74, 0x18, 0xcd, 0x0d, 0xe9, 0x7f, 0x01, 0xaa, 0x4b, 0xd5, 0x7d, 0xd7, 0xb1, 0x7d, 0x8c, 0x14,
	0x28, 0x9a, 0x78, 0xa6, 0xbb, 0x9e, 0x33, 0x26, 0x16, 0x16, 0x85, 0x7a, 0xba, 0x51, 0x6c, 0x7f,
	0xda, 0x4c, 0x98, 0xa7, 0x29, 0xf7, 0xc2, 0xfc, 0xaf, 0xc3, 0x58, 0x0d, 0x4c, 0x3c, 0x8b, 0xce,
	0xac, 0xed, 0xc8, 0x09, 0xec, 0x18, 0x4d, 0x68, 0xa0, 0x2b, 0xd8, 0x0d, 0x7c, 0xec, 0xcd, 0xab,
	0x33, 0x4c, 0xab, 0xaa, 0xc7, 0x55, 0x23, 0x60, 0x5a, 0x91, 0x65, 0x46, 0x4e, 0x49, 0x85, 0x8f,
	0x18, 0xfc, 0x25, 0x00, 0xeb, 0x37, 0x54, 0x85, 0x1c, 0x9b, 0x8c, 0x98, 0x31, 0x26, 0x13, 0xcf,
	0x54, 0x53, 0xfa, 0x3d, 0x05, 0xfb, 0x4b, 0x85, 0x50, 0x09, 0x52, 0xf3, 0xec, 0x14, 0x31, 0xd1,
	0x21, 0xec, 0xb0, 0x54, 0x1c, 0x10, 0x9e, 0x5b, 0xd0, 0x58, 0x25, 0x25, 0x20, 0xe8, 0x63, 0x28,
	0x8c, 0xef, 0xf5, 0x1f, 0x0d, 0xcb, 0xc2, 0xf1, 0x86, 0xf3, 0xe3, 0xfb, 0x21, 0xb7, 0xd1, 0x05,
	0x64, 0xa6, 0x8e, 0x89, 0xc5, 0x4c, 0x5d, 0x68, 0x94, 0xda, 0xc7, 0xc9, 0x3b, 0xe4, 0x7d, 0x6f,
	0x1c, 0x13, 0x6b, 0x3c, 0x18, 0x7d, 0x02, 0x30, 0xf2, 0xb0, 0x41, 0xb1, 0xa9, 0x1b, 0x54, 0xcc,
	0xf2, 0x6e, 0x85, 0xc8, 0xd3, 0xa1, 0xa8, 0x0e, 0xbb, 0x96, 0xe1, 0x53, 0xdd, 0xc7, 0xd8, 0x66,
	0x01, 0x39, 0x1e, 0x00, 0xcc, 0xd7, 0xc3, 0xd8, 0xee, 0x50, 0xf4, 0x19, 0x94, 0x0c, 0xd7, 0xb5,
	0xc8, 0x88, 0xb3, 0x8b, 0x8d, 0xbb, 0xc3, 0x71, 0xed, 0x2d, 0x78, 0x55, 0x13, 0x21, 0xc8, 0xd8,
	0xc6, 0x14, 0x8b, 0x79, 0x5e, 0x80, 0x9f, 0xa5, 0xff, 0x04, 0xa8, 0x25, 0xad, 0x75, 0x15, 0x35,
	0x84, 0xad, 0xa8, 0xf1, 0x6a, 0x24, 0x08, 0xe0, 0x70, 0xce, 0xe1, 0x6b, 0xe2, 0x53, 0xc7, 0xfb,
	0x79, 0x2b, 0x0a, 0x2c, 0x68, 0x27, 0x9d, 0xac, 0x9d, 0xcc, 0xa2, 0x76, 0x7e, 0x13, 0x40, 0x7c,
	0xd9, 0x37, 0xda, 0xd1, 0x71, 0xb8, 0xa3, 0xef, 0x43, 0x37, 0xef, 0x5e, 0xe0, 0xd3, 0x47, 0x81,
	0x2f, 0xa6, 0x4f, 0x6d, 0x3b, 0xfd, 0x2f, 0x02, 0x54, 0x7a, 0x31, 0x0c, 0x4e, 0xa0, 0xad, 0x66,
	0xbf, 0x84, 0x3c, 0x73, 0x73, 0xa2, 0xa6, 0x37, 0x23, 0x2a, 0xd3, 0x02, 0x3b, 0x48, 0x3f, 0x41,
	0x75, 0x09, 0x41, 0xb4, 0x85, 0x03, 0xc8, 0xf9, 0xd4, 0xa0, 0x81, 0xcf, 0x21, 0xe4, 0xb5, 0xc8,
	0x7a, 0xb5, 0xab, 0x3f, 0x35, 0x00, 0x9e, 0xda, 0xa2, 0x7d, 0x28, 0xca, 0x03, 0x5d, 0xbd, 0xed,
	0x7c, 0xd1, 0x57, 0x07, 0x4a, 0xf9, 0x03, 0x74, 0x04, 0xa2, 0x3c, 0xd0, 0xbf, 0xd4, 0x14, 0x45,
	0xbf, 0xea, 0xf4, 0x95, 0x61, 0xe7, 0x9b, 0x9e, 0xde, 0x55, 0x6f, 0xd4, 0xbe, 0x22, 0x97, 0x05,
	0x54, 0x81, 0xb2, 0x3c, 0xd0, 0x87, 0xd7, 0x77, 0x5d, 0x45, 0xbf, 0x55, 0xfa, 0xc3, 0x3b, 0xed,
	0xab, 0x72, 0x0a, 0x95, 0x00, 0xe4, 0x81, 0x2e, 0x2b, 0x5d, 0x85, 0x45, 0xa5, 0xdb, 0x7f, 0x65,
	0x9f, 0xde, 0x85, 0x1e, 0xf6, 0xd8, 0x0f, 0xfa, 0x43, 0x80, 0xbd, 0x67, 0xcf, 0x26, 0x3a, 0x49,
	0xc4, 0x9e, 0xf4, 0x70, 0xd7, 0x4e, 0x37, 0x09, 0x0d, 0x87, 0x95, 0x1a, 0xbf, 0xbe, 0x79, 0xfb,
	0x67, 0x4a, 0x42, 0x75, 0xfe, 0xb1, 0x08, 0x3f, 0x2d, 0xad, 0x87, 0xf0, 0x5a, 0x1f, 0x23, 0xfb,
	0xdc, 0x62, 0x00, 0xfe, 0x15, 0xa0, 0x3c, 0xaf, 0x11, 0x4b, 0xac, 0xb9, 0xb2, 0x55, 0xe2, 0x8b,
	0x59, 0x6b, 0x6d, 0x1c, 0x1f, 0xe1, 0xfb, 0x9c, 0xe3, 0x6b, 0xa2, 0xb3, 0x75, 0xf8, 0xa2, 0x5b,
	0x6e, 0x3d, 0x84, 0xbc, 0x7b, 0x44, 0xff, 0x2c, 0x62, 0x8d, 0x05, 0x71, 0xb6, 0x7e, 0x2d, 0xcf,
	0x85, 0x5d, 0x3b, 0xdf, 0x30, 0xfa, 0x7d, 0x70, 0x46, 0x62, 0x7d, 0xc2, 0xf9, 0xb7, 0x00, 0x7b,
	0xcf, 0x88, 0xbd, 0xe2, 0x9a, 0x93, 0xe4, 0xb7, 0xe2, 0x9a, 0x13, 0x75, 0x12, 0xc3, 0x93, 0x4e,
	0xd6, 0xc1, 0x63, 0xd2, 0x9c, 0x63, 0xbb, 0x14, 0x4e, 0xbf, 0xcb, 0xf1, 0x3f, 0x08, 0x17, 0xef,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x31, 0xa1, 0x9a, 0x1e, 0x7d, 0x08, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DSDeviceServiceClient is the client API for DSDeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DSDeviceServiceClient interface {
	GetDeviceList(ctx context.Context, in *GetDeviceListRequest, opts ...grpc.CallOption) (*GetDeviceListResponse, error)
	GetDeviceProfile(ctx context.Context, in *GetDSDeviceProfileRequest, opts ...grpc.CallOption) (*GetDSDeviceProfileResponse, error)
	GetDeviceHistory(ctx context.Context, in *GetDeviceHistoryRequest, opts ...grpc.CallOption) (*GetDeviceHistoryResponse, error)
	SetDeviceMode(ctx context.Context, in *SetDeviceModeRequest, opts ...grpc.CallOption) (*SetDeviceModeResponse, error)
}

type dSDeviceServiceClient struct {
	cc *grpc.ClientConn
}

func NewDSDeviceServiceClient(cc *grpc.ClientConn) DSDeviceServiceClient {
	return &dSDeviceServiceClient{cc}
}

func (c *dSDeviceServiceClient) GetDeviceList(ctx context.Context, in *GetDeviceListRequest, opts ...grpc.CallOption) (*GetDeviceListResponse, error) {
	out := new(GetDeviceListResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.DSDeviceService/GetDeviceList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dSDeviceServiceClient) GetDeviceProfile(ctx context.Context, in *GetDSDeviceProfileRequest, opts ...grpc.CallOption) (*GetDSDeviceProfileResponse, error) {
	out := new(GetDSDeviceProfileResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.DSDeviceService/GetDeviceProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dSDeviceServiceClient) GetDeviceHistory(ctx context.Context, in *GetDeviceHistoryRequest, opts ...grpc.CallOption) (*GetDeviceHistoryResponse, error) {
	out := new(GetDeviceHistoryResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.DSDeviceService/GetDeviceHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dSDeviceServiceClient) SetDeviceMode(ctx context.Context, in *SetDeviceModeRequest, opts ...grpc.CallOption) (*SetDeviceModeResponse, error) {
	out := new(SetDeviceModeResponse)
	err := c.cc.Invoke(ctx, "/appserver_serves_ui.DSDeviceService/SetDeviceMode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DSDeviceServiceServer is the server API for DSDeviceService service.
type DSDeviceServiceServer interface {
	GetDeviceList(context.Context, *GetDeviceListRequest) (*GetDeviceListResponse, error)
	GetDeviceProfile(context.Context, *GetDSDeviceProfileRequest) (*GetDSDeviceProfileResponse, error)
	GetDeviceHistory(context.Context, *GetDeviceHistoryRequest) (*GetDeviceHistoryResponse, error)
	SetDeviceMode(context.Context, *SetDeviceModeRequest) (*SetDeviceModeResponse, error)
}

// UnimplementedDSDeviceServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDSDeviceServiceServer struct {
}

func (*UnimplementedDSDeviceServiceServer) GetDeviceList(ctx context.Context, req *GetDeviceListRequest) (*GetDeviceListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceList not implemented")
}
func (*UnimplementedDSDeviceServiceServer) GetDeviceProfile(ctx context.Context, req *GetDSDeviceProfileRequest) (*GetDSDeviceProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceProfile not implemented")
}
func (*UnimplementedDSDeviceServiceServer) GetDeviceHistory(ctx context.Context, req *GetDeviceHistoryRequest) (*GetDeviceHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceHistory not implemented")
}
func (*UnimplementedDSDeviceServiceServer) SetDeviceMode(ctx context.Context, req *SetDeviceModeRequest) (*SetDeviceModeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDeviceMode not implemented")
}

func RegisterDSDeviceServiceServer(s *grpc.Server, srv DSDeviceServiceServer) {
	s.RegisterService(&_DSDeviceService_serviceDesc, srv)
}

func _DSDeviceService_GetDeviceList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DSDeviceServiceServer).GetDeviceList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.DSDeviceService/GetDeviceList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DSDeviceServiceServer).GetDeviceList(ctx, req.(*GetDeviceListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DSDeviceService_GetDeviceProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDSDeviceProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DSDeviceServiceServer).GetDeviceProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.DSDeviceService/GetDeviceProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DSDeviceServiceServer).GetDeviceProfile(ctx, req.(*GetDSDeviceProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DSDeviceService_GetDeviceHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DSDeviceServiceServer).GetDeviceHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.DSDeviceService/GetDeviceHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DSDeviceServiceServer).GetDeviceHistory(ctx, req.(*GetDeviceHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DSDeviceService_SetDeviceMode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDeviceModeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DSDeviceServiceServer).SetDeviceMode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/appserver_serves_ui.DSDeviceService/SetDeviceMode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DSDeviceServiceServer).SetDeviceMode(ctx, req.(*SetDeviceModeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DSDeviceService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "appserver_serves_ui.DSDeviceService",
	HandlerType: (*DSDeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDeviceList",
			Handler:    _DSDeviceService_GetDeviceList_Handler,
		},
		{
			MethodName: "GetDeviceProfile",
			Handler:    _DSDeviceService_GetDeviceProfile_Handler,
		},
		{
			MethodName: "GetDeviceHistory",
			Handler:    _DSDeviceService_GetDeviceHistory_Handler,
		},
		{
			MethodName: "SetDeviceMode",
			Handler:    _DSDeviceService_SetDeviceMode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "m2mserver_device.proto",
}
