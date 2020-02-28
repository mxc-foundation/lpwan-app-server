// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wallet.proto

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

type GetWalletBalanceRequest struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	OrgId                int64    `protobuf:"varint,2,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWalletBalanceRequest) Reset()         { *m = GetWalletBalanceRequest{} }
func (m *GetWalletBalanceRequest) String() string { return proto.CompactTextString(m) }
func (*GetWalletBalanceRequest) ProtoMessage()    {}
func (*GetWalletBalanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{0}
}

func (m *GetWalletBalanceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWalletBalanceRequest.Unmarshal(m, b)
}
func (m *GetWalletBalanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWalletBalanceRequest.Marshal(b, m, deterministic)
}
func (m *GetWalletBalanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWalletBalanceRequest.Merge(m, src)
}
func (m *GetWalletBalanceRequest) XXX_Size() int {
	return xxx_messageInfo_GetWalletBalanceRequest.Size(m)
}
func (m *GetWalletBalanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWalletBalanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetWalletBalanceRequest proto.InternalMessageInfo

func (m *GetWalletBalanceRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *GetWalletBalanceRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

type GetWalletBalanceResponse struct {
	Balance              float64          `protobuf:"fixed64,1,opt,name=balance,proto3" json:"balance,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,2,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetWalletBalanceResponse) Reset()         { *m = GetWalletBalanceResponse{} }
func (m *GetWalletBalanceResponse) String() string { return proto.CompactTextString(m) }
func (*GetWalletBalanceResponse) ProtoMessage()    {}
func (*GetWalletBalanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{1}
}

func (m *GetWalletBalanceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWalletBalanceResponse.Unmarshal(m, b)
}
func (m *GetWalletBalanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWalletBalanceResponse.Marshal(b, m, deterministic)
}
func (m *GetWalletBalanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWalletBalanceResponse.Merge(m, src)
}
func (m *GetWalletBalanceResponse) XXX_Size() int {
	return xxx_messageInfo_GetWalletBalanceResponse.Size(m)
}
func (m *GetWalletBalanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWalletBalanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetWalletBalanceResponse proto.InternalMessageInfo

func (m *GetWalletBalanceResponse) GetBalance() float64 {
	if m != nil {
		return m.Balance
	}
	return 0
}

func (m *GetWalletBalanceResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

type GetVmxcTxHistoryRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int64    `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetVmxcTxHistoryRequest) Reset()         { *m = GetVmxcTxHistoryRequest{} }
func (m *GetVmxcTxHistoryRequest) String() string { return proto.CompactTextString(m) }
func (*GetVmxcTxHistoryRequest) ProtoMessage()    {}
func (*GetVmxcTxHistoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{2}
}

func (m *GetVmxcTxHistoryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetVmxcTxHistoryRequest.Unmarshal(m, b)
}
func (m *GetVmxcTxHistoryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetVmxcTxHistoryRequest.Marshal(b, m, deterministic)
}
func (m *GetVmxcTxHistoryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetVmxcTxHistoryRequest.Merge(m, src)
}
func (m *GetVmxcTxHistoryRequest) XXX_Size() int {
	return xxx_messageInfo_GetVmxcTxHistoryRequest.Size(m)
}
func (m *GetVmxcTxHistoryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetVmxcTxHistoryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetVmxcTxHistoryRequest proto.InternalMessageInfo

func (m *GetVmxcTxHistoryRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

func (m *GetVmxcTxHistoryRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *GetVmxcTxHistoryRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type VmxcTxHistory struct {
	From                 string   `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To                   string   `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	TxType               string   `protobuf:"bytes,3,opt,name=tx_type,json=txType,proto3" json:"tx_type,omitempty"`
	Amount               float64  `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	CreatedAt            string   `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VmxcTxHistory) Reset()         { *m = VmxcTxHistory{} }
func (m *VmxcTxHistory) String() string { return proto.CompactTextString(m) }
func (*VmxcTxHistory) ProtoMessage()    {}
func (*VmxcTxHistory) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{3}
}

func (m *VmxcTxHistory) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VmxcTxHistory.Unmarshal(m, b)
}
func (m *VmxcTxHistory) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VmxcTxHistory.Marshal(b, m, deterministic)
}
func (m *VmxcTxHistory) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VmxcTxHistory.Merge(m, src)
}
func (m *VmxcTxHistory) XXX_Size() int {
	return xxx_messageInfo_VmxcTxHistory.Size(m)
}
func (m *VmxcTxHistory) XXX_DiscardUnknown() {
	xxx_messageInfo_VmxcTxHistory.DiscardUnknown(m)
}

var xxx_messageInfo_VmxcTxHistory proto.InternalMessageInfo

func (m *VmxcTxHistory) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *VmxcTxHistory) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *VmxcTxHistory) GetTxType() string {
	if m != nil {
		return m.TxType
	}
	return ""
}

func (m *VmxcTxHistory) GetAmount() float64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *VmxcTxHistory) GetCreatedAt() string {
	if m != nil {
		return m.CreatedAt
	}
	return ""
}

type GetVmxcTxHistoryResponse struct {
	Count                int64            `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	TxHistory            []*VmxcTxHistory `protobuf:"bytes,2,rep,name=tx_history,json=txHistory,proto3" json:"tx_history,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,3,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetVmxcTxHistoryResponse) Reset()         { *m = GetVmxcTxHistoryResponse{} }
func (m *GetVmxcTxHistoryResponse) String() string { return proto.CompactTextString(m) }
func (*GetVmxcTxHistoryResponse) ProtoMessage()    {}
func (*GetVmxcTxHistoryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{4}
}

func (m *GetVmxcTxHistoryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetVmxcTxHistoryResponse.Unmarshal(m, b)
}
func (m *GetVmxcTxHistoryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetVmxcTxHistoryResponse.Marshal(b, m, deterministic)
}
func (m *GetVmxcTxHistoryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetVmxcTxHistoryResponse.Merge(m, src)
}
func (m *GetVmxcTxHistoryResponse) XXX_Size() int {
	return xxx_messageInfo_GetVmxcTxHistoryResponse.Size(m)
}
func (m *GetVmxcTxHistoryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetVmxcTxHistoryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetVmxcTxHistoryResponse proto.InternalMessageInfo

func (m *GetVmxcTxHistoryResponse) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *GetVmxcTxHistoryResponse) GetTxHistory() []*VmxcTxHistory {
	if m != nil {
		return m.TxHistory
	}
	return nil
}

func (m *GetVmxcTxHistoryResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

type GetWalletUsageHistRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int64    `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWalletUsageHistRequest) Reset()         { *m = GetWalletUsageHistRequest{} }
func (m *GetWalletUsageHistRequest) String() string { return proto.CompactTextString(m) }
func (*GetWalletUsageHistRequest) ProtoMessage()    {}
func (*GetWalletUsageHistRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{5}
}

func (m *GetWalletUsageHistRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWalletUsageHistRequest.Unmarshal(m, b)
}
func (m *GetWalletUsageHistRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWalletUsageHistRequest.Marshal(b, m, deterministic)
}
func (m *GetWalletUsageHistRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWalletUsageHistRequest.Merge(m, src)
}
func (m *GetWalletUsageHistRequest) XXX_Size() int {
	return xxx_messageInfo_GetWalletUsageHistRequest.Size(m)
}
func (m *GetWalletUsageHistRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWalletUsageHistRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetWalletUsageHistRequest proto.InternalMessageInfo

func (m *GetWalletUsageHistRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

func (m *GetWalletUsageHistRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *GetWalletUsageHistRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type GetWalletUsageHist struct {
	StartAt              string   `protobuf:"bytes,1,opt,name=StartAt,proto3" json:"StartAt,omitempty"`
	DurationMinutes      int64    `protobuf:"varint,2,opt,name=DurationMinutes,proto3" json:"DurationMinutes,omitempty"`
	DlCntDv              int64    `protobuf:"varint,3,opt,name=DlCntDv,proto3" json:"DlCntDv,omitempty"`
	DlCntDvFree          int64    `protobuf:"varint,4,opt,name=DlCntDvFree,proto3" json:"DlCntDvFree,omitempty"`
	UlCntDv              int64    `protobuf:"varint,5,opt,name=UlCntDv,proto3" json:"UlCntDv,omitempty"`
	UlCntDvFree          int64    `protobuf:"varint,6,opt,name=UlCntDvFree,proto3" json:"UlCntDvFree,omitempty"`
	DlCntGw              int64    `protobuf:"varint,7,opt,name=DlCntGw,proto3" json:"DlCntGw,omitempty"`
	DlCntGwFree          int64    `protobuf:"varint,8,opt,name=DlCntGwFree,proto3" json:"DlCntGwFree,omitempty"`
	UlCntGw              int64    `protobuf:"varint,9,opt,name=UlCntGw,proto3" json:"UlCntGw,omitempty"`
	UlCntGwFree          int64    `protobuf:"varint,10,opt,name=UlCntGwFree,proto3" json:"UlCntGwFree,omitempty"`
	Spend                float64  `protobuf:"fixed64,11,opt,name=Spend,proto3" json:"Spend,omitempty"`
	Income               float64  `protobuf:"fixed64,12,opt,name=Income,proto3" json:"Income,omitempty"`
	UpdatedBalance       float64  `protobuf:"fixed64,13,opt,name=UpdatedBalance,proto3" json:"UpdatedBalance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWalletUsageHist) Reset()         { *m = GetWalletUsageHist{} }
func (m *GetWalletUsageHist) String() string { return proto.CompactTextString(m) }
func (*GetWalletUsageHist) ProtoMessage()    {}
func (*GetWalletUsageHist) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{6}
}

func (m *GetWalletUsageHist) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWalletUsageHist.Unmarshal(m, b)
}
func (m *GetWalletUsageHist) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWalletUsageHist.Marshal(b, m, deterministic)
}
func (m *GetWalletUsageHist) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWalletUsageHist.Merge(m, src)
}
func (m *GetWalletUsageHist) XXX_Size() int {
	return xxx_messageInfo_GetWalletUsageHist.Size(m)
}
func (m *GetWalletUsageHist) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWalletUsageHist.DiscardUnknown(m)
}

var xxx_messageInfo_GetWalletUsageHist proto.InternalMessageInfo

func (m *GetWalletUsageHist) GetStartAt() string {
	if m != nil {
		return m.StartAt
	}
	return ""
}

func (m *GetWalletUsageHist) GetDurationMinutes() int64 {
	if m != nil {
		return m.DurationMinutes
	}
	return 0
}

func (m *GetWalletUsageHist) GetDlCntDv() int64 {
	if m != nil {
		return m.DlCntDv
	}
	return 0
}

func (m *GetWalletUsageHist) GetDlCntDvFree() int64 {
	if m != nil {
		return m.DlCntDvFree
	}
	return 0
}

func (m *GetWalletUsageHist) GetUlCntDv() int64 {
	if m != nil {
		return m.UlCntDv
	}
	return 0
}

func (m *GetWalletUsageHist) GetUlCntDvFree() int64 {
	if m != nil {
		return m.UlCntDvFree
	}
	return 0
}

func (m *GetWalletUsageHist) GetDlCntGw() int64 {
	if m != nil {
		return m.DlCntGw
	}
	return 0
}

func (m *GetWalletUsageHist) GetDlCntGwFree() int64 {
	if m != nil {
		return m.DlCntGwFree
	}
	return 0
}

func (m *GetWalletUsageHist) GetUlCntGw() int64 {
	if m != nil {
		return m.UlCntGw
	}
	return 0
}

func (m *GetWalletUsageHist) GetUlCntGwFree() int64 {
	if m != nil {
		return m.UlCntGwFree
	}
	return 0
}

func (m *GetWalletUsageHist) GetSpend() float64 {
	if m != nil {
		return m.Spend
	}
	return 0
}

func (m *GetWalletUsageHist) GetIncome() float64 {
	if m != nil {
		return m.Income
	}
	return 0
}

func (m *GetWalletUsageHist) GetUpdatedBalance() float64 {
	if m != nil {
		return m.UpdatedBalance
	}
	return 0
}

type GetWalletUsageHistResponse struct {
	WalletUsageHis       []*GetWalletUsageHist `protobuf:"bytes,1,rep,name=wallet_usage_his,json=walletUsageHis,proto3" json:"wallet_usage_his,omitempty"`
	UserProfile          *ProfileResponse      `protobuf:"bytes,2,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	Count                int64                 `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *GetWalletUsageHistResponse) Reset()         { *m = GetWalletUsageHistResponse{} }
func (m *GetWalletUsageHistResponse) String() string { return proto.CompactTextString(m) }
func (*GetWalletUsageHistResponse) ProtoMessage()    {}
func (*GetWalletUsageHistResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{7}
}

func (m *GetWalletUsageHistResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWalletUsageHistResponse.Unmarshal(m, b)
}
func (m *GetWalletUsageHistResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWalletUsageHistResponse.Marshal(b, m, deterministic)
}
func (m *GetWalletUsageHistResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWalletUsageHistResponse.Merge(m, src)
}
func (m *GetWalletUsageHistResponse) XXX_Size() int {
	return xxx_messageInfo_GetWalletUsageHistResponse.Size(m)
}
func (m *GetWalletUsageHistResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWalletUsageHistResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetWalletUsageHistResponse proto.InternalMessageInfo

func (m *GetWalletUsageHistResponse) GetWalletUsageHis() []*GetWalletUsageHist {
	if m != nil {
		return m.WalletUsageHis
	}
	return nil
}

func (m *GetWalletUsageHistResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

func (m *GetWalletUsageHistResponse) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type GetDownLinkPriceRequest struct {
	OrgId                int64    `protobuf:"varint,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDownLinkPriceRequest) Reset()         { *m = GetDownLinkPriceRequest{} }
func (m *GetDownLinkPriceRequest) String() string { return proto.CompactTextString(m) }
func (*GetDownLinkPriceRequest) ProtoMessage()    {}
func (*GetDownLinkPriceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{8}
}

func (m *GetDownLinkPriceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDownLinkPriceRequest.Unmarshal(m, b)
}
func (m *GetDownLinkPriceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDownLinkPriceRequest.Marshal(b, m, deterministic)
}
func (m *GetDownLinkPriceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDownLinkPriceRequest.Merge(m, src)
}
func (m *GetDownLinkPriceRequest) XXX_Size() int {
	return xxx_messageInfo_GetDownLinkPriceRequest.Size(m)
}
func (m *GetDownLinkPriceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDownLinkPriceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetDownLinkPriceRequest proto.InternalMessageInfo

func (m *GetDownLinkPriceRequest) GetOrgId() int64 {
	if m != nil {
		return m.OrgId
	}
	return 0
}

type GetDownLinkPriceResponse struct {
	DownLinkPrice        float64          `protobuf:"fixed64,1,opt,name=down_link_price,json=downLinkPrice,proto3" json:"down_link_price,omitempty"`
	UserProfile          *ProfileResponse `protobuf:"bytes,2,opt,name=user_profile,json=userProfile,proto3" json:"user_profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetDownLinkPriceResponse) Reset()         { *m = GetDownLinkPriceResponse{} }
func (m *GetDownLinkPriceResponse) String() string { return proto.CompactTextString(m) }
func (*GetDownLinkPriceResponse) ProtoMessage()    {}
func (*GetDownLinkPriceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{9}
}

func (m *GetDownLinkPriceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDownLinkPriceResponse.Unmarshal(m, b)
}
func (m *GetDownLinkPriceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDownLinkPriceResponse.Marshal(b, m, deterministic)
}
func (m *GetDownLinkPriceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDownLinkPriceResponse.Merge(m, src)
}
func (m *GetDownLinkPriceResponse) XXX_Size() int {
	return xxx_messageInfo_GetDownLinkPriceResponse.Size(m)
}
func (m *GetDownLinkPriceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDownLinkPriceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetDownLinkPriceResponse proto.InternalMessageInfo

func (m *GetDownLinkPriceResponse) GetDownLinkPrice() float64 {
	if m != nil {
		return m.DownLinkPrice
	}
	return 0
}

func (m *GetDownLinkPriceResponse) GetUserProfile() *ProfileResponse {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

func init() {
	proto.RegisterType((*GetWalletBalanceRequest)(nil), "m2m_ui.GetWalletBalanceRequest")
	proto.RegisterType((*GetWalletBalanceResponse)(nil), "m2m_ui.GetWalletBalanceResponse")
	proto.RegisterType((*GetVmxcTxHistoryRequest)(nil), "m2m_ui.GetVmxcTxHistoryRequest")
	proto.RegisterType((*VmxcTxHistory)(nil), "m2m_ui.VmxcTxHistory")
	proto.RegisterType((*GetVmxcTxHistoryResponse)(nil), "m2m_ui.GetVmxcTxHistoryResponse")
	proto.RegisterType((*GetWalletUsageHistRequest)(nil), "m2m_ui.GetWalletUsageHistRequest")
	proto.RegisterType((*GetWalletUsageHist)(nil), "m2m_ui.GetWalletUsageHist")
	proto.RegisterType((*GetWalletUsageHistResponse)(nil), "m2m_ui.GetWalletUsageHistResponse")
	proto.RegisterType((*GetDownLinkPriceRequest)(nil), "m2m_ui.GetDownLinkPriceRequest")
	proto.RegisterType((*GetDownLinkPriceResponse)(nil), "m2m_ui.GetDownLinkPriceResponse")
}

func init() { proto.RegisterFile("wallet.proto", fileDescriptor_b88fd140af4deb6f) }

var fileDescriptor_b88fd140af4deb6f = []byte{
	// 788 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x41, 0x6f, 0xd3, 0x4a,
	0x10, 0x96, 0x93, 0x26, 0x7d, 0x99, 0x34, 0x6d, 0xb5, 0xaf, 0x7d, 0xf5, 0xcb, 0x7b, 0xa8, 0xa9,
	0x81, 0x52, 0x81, 0x68, 0x50, 0xe0, 0xc4, 0xad, 0x10, 0x11, 0x22, 0x81, 0x54, 0xb9, 0x0d, 0xdc,
	0x30, 0x6e, 0xbc, 0x09, 0x56, 0x1d, 0xaf, 0xb1, 0x37, 0x75, 0x2a, 0xa0, 0x07, 0x38, 0x70, 0xe0,
	0xc8, 0x9d, 0x03, 0xbf, 0x81, 0x7f, 0xc2, 0x5f, 0xe0, 0x87, 0xa0, 0x9d, 0x5d, 0x27, 0x76, 0x6a,
	0x2a, 0xa1, 0x72, 0xcb, 0x8c, 0xf7, 0xfb, 0xbe, 0xd9, 0x99, 0x6f, 0x27, 0xb0, 0x14, 0xdb, 0x9e,
	0x47, 0xf9, 0x6e, 0x10, 0x32, 0xce, 0x48, 0x79, 0xd4, 0x1a, 0x59, 0x63, 0xb7, 0xfe, 0xff, 0x90,
	0xb1, 0xa1, 0x47, 0x9b, 0x76, 0xe0, 0x36, 0x6d, 0xdf, 0x67, 0xdc, 0xe6, 0x2e, 0xf3, 0x23, 0x79,
	0xaa, 0x5e, 0x0b, 0x42, 0x36, 0x70, 0x3d, 0x2a, 0x43, 0xa3, 0x0b, 0x1b, 0x1d, 0xca, 0x9f, 0x23,
	0xcf, 0x03, 0xdb, 0xb3, 0xfd, 0x3e, 0x35, 0xe9, 0xeb, 0x31, 0x8d, 0x38, 0xd9, 0x80, 0xc5, 0x71,
	0x44, 0x43, 0xcb, 0x75, 0x74, 0xad, 0xa1, 0xed, 0x14, 0xcd, 0xb2, 0x08, 0xbb, 0x0e, 0x59, 0x87,
	0x32, 0x0b, 0x87, 0x22, 0x5f, 0xc0, 0x7c, 0x89, 0x85, 0xc3, 0xae, 0x63, 0x04, 0xa0, 0x9f, 0xa7,
	0x8a, 0x02, 0xe6, 0x47, 0x94, 0xe8, 0xb0, 0x78, 0x24, 0x53, 0xc8, 0xa5, 0x99, 0x49, 0x48, 0xee,
	0xc3, 0x12, 0xaa, 0xa8, 0xb2, 0x90, 0xb2, 0xda, 0xda, 0xd8, 0x95, 0x97, 0xd9, 0xdd, 0x97, 0xe9,
	0x84, 0xc8, 0xac, 0x8a, 0xc3, 0x2a, 0x69, 0xbc, 0xc0, 0xe2, 0x9f, 0x8d, 0x26, 0xfd, 0xc3, 0xc9,
	0x63, 0x37, 0xe2, 0x2c, 0x3c, 0x4d, 0x8a, 0x9f, 0xd5, 0xa8, 0xa5, 0x6a, 0x24, 0xff, 0x40, 0x99,
	0x0d, 0x06, 0x11, 0xe5, 0xaa, 0x74, 0x15, 0x91, 0x35, 0x28, 0x79, 0xee, 0xc8, 0xe5, 0x7a, 0x51,
	0x9e, 0xc6, 0xc0, 0xf8, 0xa0, 0x41, 0x2d, 0xc3, 0x4e, 0x08, 0x2c, 0x0c, 0x42, 0x36, 0x42, 0xd2,
	0x8a, 0x89, 0xbf, 0xc9, 0x32, 0x14, 0x38, 0x43, 0xbe, 0x8a, 0x59, 0xe0, 0x4c, 0xf4, 0x8d, 0x4f,
	0x2c, 0x7e, 0x1a, 0x50, 0x64, 0xab, 0x98, 0x65, 0x3e, 0x39, 0x3c, 0x0d, 0xa8, 0x10, 0xb7, 0x47,
	0x6c, 0xec, 0x73, 0x7d, 0x01, 0x7b, 0xa0, 0x22, 0x72, 0x05, 0xa0, 0x1f, 0x52, 0x9b, 0x53, 0xc7,
	0xb2, 0xb9, 0x5e, 0x42, 0x4c, 0x45, 0x65, 0xf6, 0xb8, 0xf1, 0x55, 0xc3, 0xc6, 0xce, 0x5d, 0x53,
	0x35, 0x76, 0x0d, 0x4a, 0x7d, 0xa4, 0x54, 0xd7, 0xc4, 0x80, 0xdc, 0x03, 0xe0, 0x13, 0xeb, 0x95,
	0x3c, 0xab, 0x17, 0x1a, 0xc5, 0x9d, 0x6a, 0x6b, 0x3d, 0x69, 0x69, 0x96, 0xa8, 0xc2, 0xa7, 0x97,
	0x9b, 0x1f, 0x45, 0xf1, 0x37, 0x46, 0xf1, 0x12, 0xfe, 0x9d, 0x0e, 0xbf, 0x17, 0xd9, 0x43, 0x2a,
	0x48, 0xff, 0xe8, 0x30, 0xbe, 0x14, 0x81, 0x9c, 0x97, 0x10, 0xce, 0x3a, 0xe0, 0x76, 0xc8, 0xf7,
	0xb8, 0x1a, 0x4a, 0x12, 0x92, 0x1d, 0x58, 0x69, 0x8f, 0x43, 0x34, 0xff, 0x53, 0xd7, 0x1f, 0x73,
	0x1a, 0x29, 0x9d, 0xf9, 0xb4, 0xe0, 0x68, 0x7b, 0x0f, 0x7d, 0xde, 0x3e, 0x51, 0x92, 0x49, 0x48,
	0x1a, 0x50, 0x55, 0x3f, 0x1f, 0x85, 0x94, 0xe2, 0xdc, 0x8a, 0x66, 0x3a, 0x25, 0xb0, 0x3d, 0x85,
	0x2d, 0x49, 0x6c, 0x6f, 0x86, 0xed, 0xa5, 0xb0, 0x65, 0x89, 0xed, 0x65, 0xb1, 0x48, 0xd5, 0x89,
	0xf5, 0xc5, 0x94, 0x6e, 0x27, 0x9e, 0xea, 0x76, 0x62, 0xc4, 0xfe, 0x95, 0xd2, 0x95, 0xa9, 0xa9,
	0x6e, 0x27, 0xd6, 0x2b, 0x29, 0x5d, 0x89, 0xed, 0xa5, 0xb0, 0x90, 0xd2, 0x55, 0xd8, 0x35, 0x28,
	0x1d, 0x04, 0xd4, 0x77, 0xf4, 0x2a, 0xfa, 0x50, 0x06, 0x62, 0x1c, 0x5d, 0xbf, 0xcf, 0x46, 0x54,
	0x5f, 0x92, 0xf6, 0x94, 0x11, 0xd9, 0x86, 0xe5, 0x5e, 0xe0, 0x08, 0x33, 0xaa, 0x57, 0xad, 0xd7,
	0xf0, 0xfb, 0x5c, 0xd6, 0xf8, 0xa6, 0x41, 0x3d, 0xcf, 0x03, 0xca, 0xa9, 0x6d, 0x58, 0x95, 0xeb,
	0xca, 0x1a, 0x8b, 0x6f, 0xc2, 0x9d, 0xba, 0x86, 0xce, 0xac, 0x27, 0x0e, 0xcb, 0x41, 0x2f, 0xc7,
	0x99, 0xc4, 0x65, 0xd6, 0xc5, 0xec, 0xad, 0x14, 0x53, 0x6f, 0xc5, 0xb8, 0x83, 0x4b, 0xa4, 0xcd,
	0x62, 0xff, 0x89, 0xeb, 0x1f, 0xef, 0x87, 0xee, 0x6c, 0x03, 0xe6, 0xfb, 0xd6, 0x38, 0xc3, 0xf7,
	0x38, 0x87, 0x50, 0xb7, 0xdc, 0x86, 0x15, 0x87, 0xc5, 0xbe, 0xe5, 0xb9, 0xfe, 0xb1, 0x15, 0x88,
	0x4f, 0x6a, 0xe1, 0xd5, 0x9c, 0xf4, 0xf9, 0xcb, 0xdc, 0xa3, 0xf5, 0x69, 0x01, 0x6a, 0xb2, 0x4f,
	0x07, 0x34, 0x3c, 0x11, 0x6c, 0x21, 0xac, 0xce, 0xaf, 0x5e, 0xb2, 0x79, 0xae, 0xab, 0xd9, 0xfd,
	0x5e, 0x6f, 0xfc, 0xfa, 0x80, 0x54, 0x35, 0xfe, 0x7b, 0xff, 0xfd, 0xc7, 0xe7, 0xc2, 0x3a, 0xf9,
	0x1b, 0xff, 0x4b, 0xe4, 0x24, 0x9a, 0xc9, 0xe2, 0x3e, 0x43, 0xcd, 0xec, 0x7a, 0x4c, 0x6b, 0xe6,
	0xad, 0xe5, 0x8c, 0x66, 0xee, 0x42, 0x33, 0x6e, 0xa0, 0xe6, 0x16, 0xd9, 0x4c, 0x6b, 0xbe, 0x91,
	0x63, 0x78, 0xd7, 0xe4, 0x93, 0xdb, 0x6a, 0xab, 0x91, 0x8f, 0x5a, 0xee, 0x3e, 0xd8, 0xba, 0xc0,
	0x4c, 0xaa, 0x08, 0xe3, 0xa2, 0x23, 0xaa, 0x8c, 0x9b, 0x58, 0xc6, 0x35, 0x62, 0xe4, 0x96, 0x81,
	0x0e, 0x9e, 0x56, 0xf2, 0x16, 0x40, 0xf8, 0xc1, 0x93, 0x93, 0x4d, 0xf7, 0x20, 0xcf, 0x55, 0x99,
	0x1e, 0xe4, 0x9a, 0xc8, 0xb8, 0x85, 0xe2, 0xd7, 0xc9, 0xd5, 0x5c, 0x71, 0x61, 0xa4, 0x99, 0xbd,
	0x8e, 0xca, 0xf8, 0x47, 0x7e, 0xf7, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0x46, 0xe9, 0x8a, 0x20,
	0x0d, 0x08, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WalletServiceClient is the client API for WalletService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WalletServiceClient interface {
	GetWalletBalance(ctx context.Context, in *GetWalletBalanceRequest, opts ...grpc.CallOption) (*GetWalletBalanceResponse, error)
	GetVmxcTxHistory(ctx context.Context, in *GetVmxcTxHistoryRequest, opts ...grpc.CallOption) (*GetVmxcTxHistoryResponse, error)
	GetWalletUsageHist(ctx context.Context, in *GetWalletUsageHistRequest, opts ...grpc.CallOption) (*GetWalletUsageHistResponse, error)
	GetDlPrice(ctx context.Context, in *GetDownLinkPriceRequest, opts ...grpc.CallOption) (*GetDownLinkPriceResponse, error)
}

type walletServiceClient struct {
	cc *grpc.ClientConn
}

func NewWalletServiceClient(cc *grpc.ClientConn) WalletServiceClient {
	return &walletServiceClient{cc}
}

func (c *walletServiceClient) GetWalletBalance(ctx context.Context, in *GetWalletBalanceRequest, opts ...grpc.CallOption) (*GetWalletBalanceResponse, error) {
	out := new(GetWalletBalanceResponse)
	err := c.cc.Invoke(ctx, "/m2m_ui.WalletService/GetWalletBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetVmxcTxHistory(ctx context.Context, in *GetVmxcTxHistoryRequest, opts ...grpc.CallOption) (*GetVmxcTxHistoryResponse, error) {
	out := new(GetVmxcTxHistoryResponse)
	err := c.cc.Invoke(ctx, "/m2m_ui.WalletService/GetVmxcTxHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetWalletUsageHist(ctx context.Context, in *GetWalletUsageHistRequest, opts ...grpc.CallOption) (*GetWalletUsageHistResponse, error) {
	out := new(GetWalletUsageHistResponse)
	err := c.cc.Invoke(ctx, "/m2m_ui.WalletService/GetWalletUsageHist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetDlPrice(ctx context.Context, in *GetDownLinkPriceRequest, opts ...grpc.CallOption) (*GetDownLinkPriceResponse, error) {
	out := new(GetDownLinkPriceResponse)
	err := c.cc.Invoke(ctx, "/m2m_ui.WalletService/GetDlPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WalletServiceServer is the server API for WalletService service.
type WalletServiceServer interface {
	GetWalletBalance(context.Context, *GetWalletBalanceRequest) (*GetWalletBalanceResponse, error)
	GetVmxcTxHistory(context.Context, *GetVmxcTxHistoryRequest) (*GetVmxcTxHistoryResponse, error)
	GetWalletUsageHist(context.Context, *GetWalletUsageHistRequest) (*GetWalletUsageHistResponse, error)
	GetDlPrice(context.Context, *GetDownLinkPriceRequest) (*GetDownLinkPriceResponse, error)
}

// UnimplementedWalletServiceServer can be embedded to have forward compatible implementations.
type UnimplementedWalletServiceServer struct {
}

func (*UnimplementedWalletServiceServer) GetWalletBalance(ctx context.Context, req *GetWalletBalanceRequest) (*GetWalletBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletBalance not implemented")
}
func (*UnimplementedWalletServiceServer) GetVmxcTxHistory(ctx context.Context, req *GetVmxcTxHistoryRequest) (*GetVmxcTxHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVmxcTxHistory not implemented")
}
func (*UnimplementedWalletServiceServer) GetWalletUsageHist(ctx context.Context, req *GetWalletUsageHistRequest) (*GetWalletUsageHistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletUsageHist not implemented")
}
func (*UnimplementedWalletServiceServer) GetDlPrice(ctx context.Context, req *GetDownLinkPriceRequest) (*GetDownLinkPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDlPrice not implemented")
}

func RegisterWalletServiceServer(s *grpc.Server, srv WalletServiceServer) {
	s.RegisterService(&_WalletService_serviceDesc, srv)
}

func _WalletService_GetWalletBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWalletBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetWalletBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_ui.WalletService/GetWalletBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetWalletBalance(ctx, req.(*GetWalletBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetVmxcTxHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVmxcTxHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetVmxcTxHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_ui.WalletService/GetVmxcTxHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetVmxcTxHistory(ctx, req.(*GetVmxcTxHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetWalletUsageHist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWalletUsageHistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetWalletUsageHist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_ui.WalletService/GetWalletUsageHist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetWalletUsageHist(ctx, req.(*GetWalletUsageHistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetDlPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDownLinkPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetDlPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m2m_ui.WalletService/GetDlPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetDlPrice(ctx, req.(*GetDownLinkPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WalletService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "m2m_ui.WalletService",
	HandlerType: (*WalletServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWalletBalance",
			Handler:    _WalletService_GetWalletBalance_Handler,
		},
		{
			MethodName: "GetVmxcTxHistory",
			Handler:    _WalletService_GetVmxcTxHistory_Handler,
		},
		{
			MethodName: "GetWalletUsageHist",
			Handler:    _WalletService_GetWalletUsageHist_Handler,
		},
		{
			MethodName: "GetDlPrice",
			Handler:    _WalletService_GetDlPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wallet.proto",
}
