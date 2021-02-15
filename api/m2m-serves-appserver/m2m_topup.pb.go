// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: m2m_topup.proto

package m2m_serves_appserver

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type GetTopUpHistoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrgId int64 `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	// int64 offset = 2;
	// int64 limit = 3;
	Currency string               `protobuf:"bytes,4,opt,name=currency,proto3" json:"currency,omitempty"`
	From     *timestamp.Timestamp `protobuf:"bytes,5,opt,name=from,proto3" json:"from,omitempty"`
	Till     *timestamp.Timestamp `protobuf:"bytes,6,opt,name=till,proto3" json:"till,omitempty"`
}

func (x *GetTopUpHistoryRequest) Reset() {
	*x = GetTopUpHistoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_m2m_topup_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTopUpHistoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTopUpHistoryRequest) ProtoMessage() {}

func (x *GetTopUpHistoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_m2m_topup_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTopUpHistoryRequest.ProtoReflect.Descriptor instead.
func (*GetTopUpHistoryRequest) Descriptor() ([]byte, []int) {
	return file_m2m_topup_proto_rawDescGZIP(), []int{0}
}

func (x *GetTopUpHistoryRequest) GetOrgId() int64 {
	if x != nil {
		return x.OrgId
	}
	return 0
}

func (x *GetTopUpHistoryRequest) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *GetTopUpHistoryRequest) GetFrom() *timestamp.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *GetTopUpHistoryRequest) GetTill() *timestamp.Timestamp {
	if x != nil {
		return x.Till
	}
	return nil
}

type TopUpHistory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// double amount = 1;
	// string created_at = 2;
	TxHash    string               `protobuf:"bytes,3,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	Timestamp *timestamp.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Amount    string               `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *TopUpHistory) Reset() {
	*x = TopUpHistory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_m2m_topup_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopUpHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopUpHistory) ProtoMessage() {}

func (x *TopUpHistory) ProtoReflect() protoreflect.Message {
	mi := &file_m2m_topup_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopUpHistory.ProtoReflect.Descriptor instead.
func (*TopUpHistory) Descriptor() ([]byte, []int) {
	return file_m2m_topup_proto_rawDescGZIP(), []int{1}
}

func (x *TopUpHistory) GetTxHash() string {
	if x != nil {
		return x.TxHash
	}
	return ""
}

func (x *TopUpHistory) GetTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *TopUpHistory) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

type GetTopUpHistoryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TopupHistory []*TopUpHistory `protobuf:"bytes,2,rep,name=topup_history,json=topupHistory,proto3" json:"topup_history,omitempty"`
}

func (x *GetTopUpHistoryResponse) Reset() {
	*x = GetTopUpHistoryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_m2m_topup_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTopUpHistoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTopUpHistoryResponse) ProtoMessage() {}

func (x *GetTopUpHistoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_m2m_topup_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTopUpHistoryResponse.ProtoReflect.Descriptor instead.
func (*GetTopUpHistoryResponse) Descriptor() ([]byte, []int) {
	return file_m2m_topup_proto_rawDescGZIP(), []int{2}
}

func (x *GetTopUpHistoryResponse) GetTopupHistory() []*TopUpHistory {
	if x != nil {
		return x.TopupHistory
	}
	return nil
}

type GetTopUpDestinationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrgId    int64  `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Currency string `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`
}

func (x *GetTopUpDestinationRequest) Reset() {
	*x = GetTopUpDestinationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_m2m_topup_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTopUpDestinationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTopUpDestinationRequest) ProtoMessage() {}

func (x *GetTopUpDestinationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_m2m_topup_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTopUpDestinationRequest.ProtoReflect.Descriptor instead.
func (*GetTopUpDestinationRequest) Descriptor() ([]byte, []int) {
	return file_m2m_topup_proto_rawDescGZIP(), []int{3}
}

func (x *GetTopUpDestinationRequest) GetOrgId() int64 {
	if x != nil {
		return x.OrgId
	}
	return 0
}

func (x *GetTopUpDestinationRequest) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

type GetTopUpDestinationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActiveAccount string `protobuf:"bytes,1,opt,name=active_account,json=activeAccount,proto3" json:"active_account,omitempty"`
}

func (x *GetTopUpDestinationResponse) Reset() {
	*x = GetTopUpDestinationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_m2m_topup_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTopUpDestinationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTopUpDestinationResponse) ProtoMessage() {}

func (x *GetTopUpDestinationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_m2m_topup_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTopUpDestinationResponse.ProtoReflect.Descriptor instead.
func (*GetTopUpDestinationResponse) Descriptor() ([]byte, []int) {
	return file_m2m_topup_proto_rawDescGZIP(), []int{4}
}

func (x *GetTopUpDestinationResponse) GetActiveAccount() string {
	if x != nil {
		return x.ActiveAccount
	}
	return ""
}

var File_m2m_topup_proto protoreflect.FileDescriptor

var file_m2m_topup_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x32, 0x6d, 0x5f, 0x74, 0x6f, 0x70, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x14, 0x6d, 0x32, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x61, 0x70,
	0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xab, 0x01, 0x0a, 0x16, 0x47, 0x65, 0x74,
	0x54, 0x6f, 0x70, 0x55, 0x70, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x2e, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6c, 0x6c, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x04, 0x74, 0x69, 0x6c, 0x6c, 0x22, 0x79, 0x0a, 0x0c, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x48,
	0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x78, 0x5f, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x78, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x22, 0x62, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x48, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0d,
	0x74, 0x6f, 0x70, 0x75, 0x70, 0x5f, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6d, 0x32, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73,
	0x5f, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x6f, 0x70, 0x55, 0x70,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x0c, 0x74, 0x6f, 0x70, 0x75, 0x70, 0x48, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x22, 0x4f, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x55,
	0x70, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x22, 0x44, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70,
	0x55, 0x70, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0xfa, 0x01, 0x0a,
	0x0c, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6e, 0x0a,
	0x0f, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x12, 0x2c, 0x2e, 0x6d, 0x32, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x61, 0x70,
	0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x55, 0x70,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d,
	0x2e, 0x6d, 0x32, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x61, 0x70, 0x70, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x48, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7a, 0x0a,
	0x13, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x2e, 0x6d, 0x32, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x73, 0x5f, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x54,
	0x6f, 0x70, 0x55, 0x70, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x6d, 0x32, 0x6d, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x73, 0x5f, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x47, 0x65,
	0x74, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x78, 0x63, 0x2d, 0x66, 0x6f, 0x75, 0x6e,
	0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6d, 0x78, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x32, 0x6d,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x73, 0x5f, 0x61, 0x70, 0x70, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_m2m_topup_proto_rawDescOnce sync.Once
	file_m2m_topup_proto_rawDescData = file_m2m_topup_proto_rawDesc
)

func file_m2m_topup_proto_rawDescGZIP() []byte {
	file_m2m_topup_proto_rawDescOnce.Do(func() {
		file_m2m_topup_proto_rawDescData = protoimpl.X.CompressGZIP(file_m2m_topup_proto_rawDescData)
	})
	return file_m2m_topup_proto_rawDescData
}

var file_m2m_topup_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_m2m_topup_proto_goTypes = []interface{}{
	(*GetTopUpHistoryRequest)(nil),      // 0: m2m_serves_appserver.GetTopUpHistoryRequest
	(*TopUpHistory)(nil),                // 1: m2m_serves_appserver.TopUpHistory
	(*GetTopUpHistoryResponse)(nil),     // 2: m2m_serves_appserver.GetTopUpHistoryResponse
	(*GetTopUpDestinationRequest)(nil),  // 3: m2m_serves_appserver.GetTopUpDestinationRequest
	(*GetTopUpDestinationResponse)(nil), // 4: m2m_serves_appserver.GetTopUpDestinationResponse
	(*timestamp.Timestamp)(nil),         // 5: google.protobuf.Timestamp
}
var file_m2m_topup_proto_depIdxs = []int32{
	5, // 0: m2m_serves_appserver.GetTopUpHistoryRequest.from:type_name -> google.protobuf.Timestamp
	5, // 1: m2m_serves_appserver.GetTopUpHistoryRequest.till:type_name -> google.protobuf.Timestamp
	5, // 2: m2m_serves_appserver.TopUpHistory.timestamp:type_name -> google.protobuf.Timestamp
	1, // 3: m2m_serves_appserver.GetTopUpHistoryResponse.topup_history:type_name -> m2m_serves_appserver.TopUpHistory
	0, // 4: m2m_serves_appserver.TopUpService.GetTopUpHistory:input_type -> m2m_serves_appserver.GetTopUpHistoryRequest
	3, // 5: m2m_serves_appserver.TopUpService.GetTopUpDestination:input_type -> m2m_serves_appserver.GetTopUpDestinationRequest
	2, // 6: m2m_serves_appserver.TopUpService.GetTopUpHistory:output_type -> m2m_serves_appserver.GetTopUpHistoryResponse
	4, // 7: m2m_serves_appserver.TopUpService.GetTopUpDestination:output_type -> m2m_serves_appserver.GetTopUpDestinationResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_m2m_topup_proto_init() }
func file_m2m_topup_proto_init() {
	if File_m2m_topup_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_m2m_topup_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTopUpHistoryRequest); i {
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
		file_m2m_topup_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopUpHistory); i {
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
		file_m2m_topup_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTopUpHistoryResponse); i {
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
		file_m2m_topup_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTopUpDestinationRequest); i {
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
		file_m2m_topup_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTopUpDestinationResponse); i {
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
			RawDescriptor: file_m2m_topup_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_m2m_topup_proto_goTypes,
		DependencyIndexes: file_m2m_topup_proto_depIdxs,
		MessageInfos:      file_m2m_topup_proto_msgTypes,
	}.Build()
	File_m2m_topup_proto = out.File
	file_m2m_topup_proto_rawDesc = nil
	file_m2m_topup_proto_goTypes = nil
	file_m2m_topup_proto_depIdxs = nil
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
	err := c.cc.Invoke(ctx, "/m2m_serves_appserver.TopUpService/GetTopUpHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topUpServiceClient) GetTopUpDestination(ctx context.Context, in *GetTopUpDestinationRequest, opts ...grpc.CallOption) (*GetTopUpDestinationResponse, error) {
	out := new(GetTopUpDestinationResponse)
	err := c.cc.Invoke(ctx, "/m2m_serves_appserver.TopUpService/GetTopUpDestination", in, out, opts...)
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

func (*UnimplementedTopUpServiceServer) GetTopUpHistory(context.Context, *GetTopUpHistoryRequest) (*GetTopUpHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopUpHistory not implemented")
}
func (*UnimplementedTopUpServiceServer) GetTopUpDestination(context.Context, *GetTopUpDestinationRequest) (*GetTopUpDestinationResponse, error) {
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
		FullMethod: "/m2m_serves_appserver.TopUpService/GetTopUpHistory",
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
		FullMethod: "/m2m_serves_appserver.TopUpService/GetTopUpDestination",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopUpServiceServer).GetTopUpDestination(ctx, req.(*GetTopUpDestinationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TopUpService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m2m_serves_appserver.TopUpService",
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
	Metadata: "m2m_topup.proto",
}
