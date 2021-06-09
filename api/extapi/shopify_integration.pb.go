// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.13.0
// source: shopify_integration.proto

package extapi

import (
	context "context"
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

type GetOrdersByUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// user's email address for supernode account
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *GetOrdersByUserRequest) Reset() {
	*x = GetOrdersByUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shopify_integration_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetOrdersByUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrdersByUserRequest) ProtoMessage() {}

func (x *GetOrdersByUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shopify_integration_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrdersByUserRequest.ProtoReflect.Descriptor instead.
func (*GetOrdersByUserRequest) Descriptor() ([]byte, []int) {
	return file_shopify_integration_proto_rawDescGZIP(), []int{0}
}

func (x *GetOrdersByUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// email address user used for shopify account,
	ShopifyAccount string `protobuf:"bytes,1,opt,name=shopify_account,json=shopifyAccount,proto3" json:"shopify_account,omitempty"`
	// order id is generated and maintained on shopify service side, appserver saves this as a reference
	OrderId   string `protobuf:"bytes,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	CreatedAt string `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// used to identify a specific product created in shopify
	ProductId string `protobuf:"bytes,4,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	// amount of product with given product_id from an order with given order_id
	AmountProduct int64 `protobuf:"varint,5,opt,name=amount_product,json=amountProduct,proto3" json:"amount_product,omitempty"`
	// when bonus_status is 'done'
	//  users who request refund will get ( number of returned good * bonus_per_piece_usd ) less
	// when bonus_status is 'pending'
	//  users will get refund with full amount
	BonusStatus string `protobuf:"bytes,6,opt,name=bonus_status,json=bonusStatus,proto3" json:"bonus_status,omitempty"`
	// amount of USD rewarded to user for purchasing one product with given product id
	BonusPerPieceUsd string `protobuf:"bytes,7,opt,name=bonus_per_piece_usd,json=bonusPerPieceUsd,proto3" json:"bonus_per_piece_usd,omitempty"`
}

func (x *Order) Reset() {
	*x = Order{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shopify_integration_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_shopify_integration_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_shopify_integration_proto_rawDescGZIP(), []int{1}
}

func (x *Order) GetShopifyAccount() string {
	if x != nil {
		return x.ShopifyAccount
	}
	return ""
}

func (x *Order) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *Order) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Order) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *Order) GetAmountProduct() int64 {
	if x != nil {
		return x.AmountProduct
	}
	return 0
}

func (x *Order) GetBonusStatus() string {
	if x != nil {
		return x.BonusStatus
	}
	return ""
}

func (x *Order) GetBonusPerPieceUsd() string {
	if x != nil {
		return x.BonusPerPieceUsd
	}
	return ""
}

type GetOrdersByUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Orders []*Order `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
}

func (x *GetOrdersByUserResponse) Reset() {
	*x = GetOrdersByUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shopify_integration_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetOrdersByUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrdersByUserResponse) ProtoMessage() {}

func (x *GetOrdersByUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shopify_integration_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrdersByUserResponse.ProtoReflect.Descriptor instead.
func (*GetOrdersByUserResponse) Descriptor() ([]byte, []int) {
	return file_shopify_integration_proto_rawDescGZIP(), []int{2}
}

func (x *GetOrdersByUserResponse) GetOrders() []*Order {
	if x != nil {
		return x.Orders
	}
	return nil
}

var File_shopify_integration_proto protoreflect.FileDescriptor

var file_shopify_integration_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x68, 0x6f, 0x70, 0x69, 0x66, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x65, 0x78, 0x74,
	0x61, 0x70, 0x69, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x2e, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x22, 0x82, 0x02, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x73,
	0x68, 0x6f, 0x70, 0x69, 0x66, 0x79, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x68, 0x6f, 0x70, 0x69, 0x66, 0x79, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x25, 0x0a,
	0x0e, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x6f, 0x6e, 0x75,
	0x73, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2d, 0x0a, 0x13, 0x62, 0x6f, 0x6e, 0x75, 0x73,
	0x5f, 0x70, 0x65, 0x72, 0x5f, 0x70, 0x69, 0x65, 0x63, 0x65, 0x5f, 0x75, 0x73, 0x64, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x50, 0x65, 0x72, 0x50, 0x69,
	0x65, 0x63, 0x65, 0x55, 0x73, 0x64, 0x22, 0x40, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x73, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x25, 0x0a, 0x06, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x65, 0x78, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x52, 0x06, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x32, 0x91, 0x01, 0x0a, 0x12, 0x53, 0x68, 0x6f,
	0x70, 0x69, 0x66, 0x79, 0x49, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x7b, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x1e, 0x2e, 0x65, 0x78, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x65, 0x78, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x73, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21, 0x12, 0x1f, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x73, 0x68, 0x6f, 0x70, 0x69, 0x66, 0x79, 0x2d, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x42, 0x3e, 0x5a, 0x3c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x78, 0x63, 0x2d, 0x66,
	0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6c, 0x70, 0x77, 0x61, 0x6e, 0x2d,
	0x61, 0x70, 0x70, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65,
	0x78, 0x74, 0x61, 0x70, 0x69, 0x3b, 0x65, 0x78, 0x74, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shopify_integration_proto_rawDescOnce sync.Once
	file_shopify_integration_proto_rawDescData = file_shopify_integration_proto_rawDesc
)

func file_shopify_integration_proto_rawDescGZIP() []byte {
	file_shopify_integration_proto_rawDescOnce.Do(func() {
		file_shopify_integration_proto_rawDescData = protoimpl.X.CompressGZIP(file_shopify_integration_proto_rawDescData)
	})
	return file_shopify_integration_proto_rawDescData
}

var file_shopify_integration_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_shopify_integration_proto_goTypes = []interface{}{
	(*GetOrdersByUserRequest)(nil),  // 0: extapi.GetOrdersByUserRequest
	(*Order)(nil),                   // 1: extapi.Order
	(*GetOrdersByUserResponse)(nil), // 2: extapi.GetOrdersByUserResponse
}
var file_shopify_integration_proto_depIdxs = []int32{
	1, // 0: extapi.GetOrdersByUserResponse.orders:type_name -> extapi.Order
	0, // 1: extapi.ShopifyIntegration.GetOrdersByUser:input_type -> extapi.GetOrdersByUserRequest
	2, // 2: extapi.ShopifyIntegration.GetOrdersByUser:output_type -> extapi.GetOrdersByUserResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_shopify_integration_proto_init() }
func file_shopify_integration_proto_init() {
	if File_shopify_integration_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shopify_integration_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetOrdersByUserRequest); i {
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
		file_shopify_integration_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Order); i {
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
		file_shopify_integration_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetOrdersByUserResponse); i {
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
			RawDescriptor: file_shopify_integration_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shopify_integration_proto_goTypes,
		DependencyIndexes: file_shopify_integration_proto_depIdxs,
		MessageInfos:      file_shopify_integration_proto_msgTypes,
	}.Build()
	File_shopify_integration_proto = out.File
	file_shopify_integration_proto_rawDesc = nil
	file_shopify_integration_proto_goTypes = nil
	file_shopify_integration_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ShopifyIntegrationClient is the client API for ShopifyIntegration service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ShopifyIntegrationClient interface {
	// GetOrdersByUser returns a list of shopify orders filtered by given email, this API is only open for global admin user
	GetOrdersByUser(ctx context.Context, in *GetOrdersByUserRequest, opts ...grpc.CallOption) (*GetOrdersByUserResponse, error)
}

type shopifyIntegrationClient struct {
	cc grpc.ClientConnInterface
}

func NewShopifyIntegrationClient(cc grpc.ClientConnInterface) ShopifyIntegrationClient {
	return &shopifyIntegrationClient{cc}
}

func (c *shopifyIntegrationClient) GetOrdersByUser(ctx context.Context, in *GetOrdersByUserRequest, opts ...grpc.CallOption) (*GetOrdersByUserResponse, error) {
	out := new(GetOrdersByUserResponse)
	err := c.cc.Invoke(ctx, "/extapi.ShopifyIntegration/GetOrdersByUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShopifyIntegrationServer is the server API for ShopifyIntegration service.
type ShopifyIntegrationServer interface {
	// GetOrdersByUser returns a list of shopify orders filtered by given email, this API is only open for global admin user
	GetOrdersByUser(context.Context, *GetOrdersByUserRequest) (*GetOrdersByUserResponse, error)
}

// UnimplementedShopifyIntegrationServer can be embedded to have forward compatible implementations.
type UnimplementedShopifyIntegrationServer struct {
}

func (*UnimplementedShopifyIntegrationServer) GetOrdersByUser(context.Context, *GetOrdersByUserRequest) (*GetOrdersByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrdersByUser not implemented")
}

func RegisterShopifyIntegrationServer(s *grpc.Server, srv ShopifyIntegrationServer) {
	s.RegisterService(&_ShopifyIntegration_serviceDesc, srv)
}

func _ShopifyIntegration_GetOrdersByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrdersByUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShopifyIntegrationServer).GetOrdersByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/extapi.ShopifyIntegration/GetOrdersByUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShopifyIntegrationServer).GetOrdersByUser(ctx, req.(*GetOrdersByUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ShopifyIntegration_serviceDesc = grpc.ServiceDesc{
	ServiceName: "extapi.ShopifyIntegration",
	HandlerType: (*ShopifyIntegrationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOrdersByUser",
			Handler:    _ShopifyIntegration_GetOrdersByUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shopify_integration.proto",
}
