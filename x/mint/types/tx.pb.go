// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mint/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	_ "google.golang.org/protobuf/types/known/timestamppb"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MsgWithdrawCoinsToAccFromTreasury is a message for withdrawing from mint
// module.
type MsgWithdrawCoinsToAccFromTreasury struct {
	// Amount is the amoutn of coins to withdraw
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	// Receiver is for whom withdraw coins
	Receiver string `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	// Sender is the message signer who submits this report transaction
	Sender string `protobuf:"bytes,3,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (m *MsgWithdrawCoinsToAccFromTreasury) Reset()         { *m = MsgWithdrawCoinsToAccFromTreasury{} }
func (m *MsgWithdrawCoinsToAccFromTreasury) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawCoinsToAccFromTreasury) ProtoMessage()    {}
func (*MsgWithdrawCoinsToAccFromTreasury) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c467a85e368a1a7, []int{0}
}
func (m *MsgWithdrawCoinsToAccFromTreasury) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawCoinsToAccFromTreasury) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawCoinsToAccFromTreasury.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawCoinsToAccFromTreasury) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawCoinsToAccFromTreasury.Merge(m, src)
}
func (m *MsgWithdrawCoinsToAccFromTreasury) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawCoinsToAccFromTreasury) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawCoinsToAccFromTreasury.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawCoinsToAccFromTreasury proto.InternalMessageInfo

func (m *MsgWithdrawCoinsToAccFromTreasury) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *MsgWithdrawCoinsToAccFromTreasury) GetReceiver() string {
	if m != nil {
		return m.Receiver
	}
	return ""
}

func (m *MsgWithdrawCoinsToAccFromTreasury) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

// MsgWithdrawCoinsToAccFromTreasuryResponse
type MsgWithdrawCoinsToAccFromTreasuryResponse struct {
}

func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) Reset() {
	*m = MsgWithdrawCoinsToAccFromTreasuryResponse{}
}
func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) String() string {
	return proto.CompactTextString(m)
}
func (*MsgWithdrawCoinsToAccFromTreasuryResponse) ProtoMessage() {}
func (*MsgWithdrawCoinsToAccFromTreasuryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c467a85e368a1a7, []int{1}
}
func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawCoinsToAccFromTreasuryResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawCoinsToAccFromTreasuryResponse.Merge(m, src)
}
func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawCoinsToAccFromTreasuryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawCoinsToAccFromTreasuryResponse proto.InternalMessageInfo

// MsgMintCoins is a message for minting from mint module.
type MsgMintCoins struct {
	// Amount is the amount of coins to mint
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	// Sender is the message signer who submits this report transaction
	Sender string `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (m *MsgMintCoins) Reset()         { *m = MsgMintCoins{} }
func (m *MsgMintCoins) String() string { return proto.CompactTextString(m) }
func (*MsgMintCoins) ProtoMessage()    {}
func (*MsgMintCoins) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c467a85e368a1a7, []int{2}
}
func (m *MsgMintCoins) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgMintCoins) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgMintCoins.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgMintCoins) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgMintCoins.Merge(m, src)
}
func (m *MsgMintCoins) XXX_Size() int {
	return m.Size()
}
func (m *MsgMintCoins) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgMintCoins.DiscardUnknown(m)
}

var xxx_messageInfo_MsgMintCoins proto.InternalMessageInfo

func (m *MsgMintCoins) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *MsgMintCoins) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

// MsgMintCoinsResponse
type MsgMintCoinsResponse struct {
}

func (m *MsgMintCoinsResponse) Reset()         { *m = MsgMintCoinsResponse{} }
func (m *MsgMintCoinsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgMintCoinsResponse) ProtoMessage()    {}
func (*MsgMintCoinsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c467a85e368a1a7, []int{3}
}
func (m *MsgMintCoinsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgMintCoinsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgMintCoinsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgMintCoinsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgMintCoinsResponse.Merge(m, src)
}
func (m *MsgMintCoinsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgMintCoinsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgMintCoinsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgMintCoinsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgWithdrawCoinsToAccFromTreasury)(nil), "mint.MsgWithdrawCoinsToAccFromTreasury")
	proto.RegisterType((*MsgWithdrawCoinsToAccFromTreasuryResponse)(nil), "mint.MsgWithdrawCoinsToAccFromTreasuryResponse")
	proto.RegisterType((*MsgMintCoins)(nil), "mint.MsgMintCoins")
	proto.RegisterType((*MsgMintCoinsResponse)(nil), "mint.MsgMintCoinsResponse")
}

func init() { proto.RegisterFile("mint/tx.proto", fileDescriptor_6c467a85e368a1a7) }

var fileDescriptor_6c467a85e368a1a7 = []byte{
	// 429 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x53, 0x4f, 0x6b, 0xd4, 0x40,
	0x1c, 0xcd, 0x74, 0xcb, 0x62, 0x47, 0x45, 0x08, 0xa5, 0x6c, 0x73, 0x98, 0xd4, 0xbd, 0xb8, 0x22,
	0x9b, 0xb1, 0xf5, 0xa6, 0x27, 0x5b, 0x11, 0x0a, 0xc6, 0x48, 0x58, 0x10, 0xbc, 0x25, 0x93, 0x71,
	0x3a, 0x68, 0xe6, 0x17, 0x66, 0x26, 0xb5, 0xeb, 0xa7, 0xd0, 0x6f, 0xe0, 0xd9, 0x0f, 0xe1, 0xc1,
	0x53, 0x8f, 0x3d, 0x7a, 0x52, 0xd9, 0xbd, 0xf8, 0x31, 0x24, 0x93, 0xec, 0x76, 0x51, 0x61, 0x3d,
	0xf5, 0x92, 0xe4, 0xf7, 0xe7, 0xbd, 0xbc, 0xf7, 0x26, 0xc1, 0x37, 0x4b, 0xa9, 0x2c, 0xb5, 0x67,
	0x51, 0xa5, 0xc1, 0x82, 0xbf, 0xd9, 0x94, 0xc1, 0xb6, 0x00, 0x01, 0xae, 0x41, 0x9b, 0xa7, 0x76,
	0x16, 0x84, 0x02, 0x40, 0xbc, 0xe5, 0xd4, 0x55, 0x79, 0xfd, 0x9a, 0x5a, 0x59, 0x72, 0x63, 0xb3,
	0xb2, 0xea, 0x16, 0x76, 0xff, 0x5c, 0xc8, 0xd4, 0xb4, 0x1b, 0xdd, 0x72, 0xaf, 0x69, 0x2e, 0x5d,
	0x83, 0x30, 0x30, 0x25, 0x18, 0x9a, 0x67, 0x86, 0xd3, 0xd3, 0xfd, 0x9c, 0xdb, 0x6c, 0x9f, 0x32,
	0x90, 0xaa, 0x9d, 0x0f, 0xbf, 0x22, 0x7c, 0x3b, 0x36, 0xe2, 0xa5, 0xb4, 0x27, 0x85, 0xce, 0xde,
	0x1d, 0x81, 0x54, 0x66, 0x02, 0x8f, 0x19, 0x7b, 0xaa, 0xa1, 0x9c, 0x68, 0x9e, 0x99, 0x5a, 0x4f,
	0x7d, 0x86, 0xfb, 0x59, 0x09, 0xb5, 0xb2, 0x03, 0xb4, 0xd7, 0x1b, 0x5d, 0x3f, 0xd8, 0x8d, 0x5a,
	0xda, 0xa8, 0xa1, 0x8d, 0x3a, 0xda, 0xa8, 0x01, 0x1f, 0xde, 0x3f, 0xff, 0x1e, 0x7a, 0x9f, 0x7f,
	0x84, 0x23, 0x21, 0xed, 0x49, 0x9d, 0x47, 0x0c, 0x4a, 0xda, 0x69, 0x68, 0x6f, 0x63, 0x53, 0xbc,
	0xa1, 0x76, 0x5a, 0x71, 0xe3, 0x00, 0x26, 0xed, 0xa8, 0xfd, 0x00, 0x5f, 0xd3, 0x9c, 0x71, 0x79,
	0xca, 0xf5, 0x60, 0x63, 0x0f, 0x8d, 0xb6, 0xd2, 0x65, 0xed, 0xef, 0xe0, 0xbe, 0xe1, 0xaa, 0xe0,
	0x7a, 0xd0, 0x73, 0x93, 0xae, 0x7a, 0xb8, 0xf9, 0xeb, 0x53, 0x88, 0x86, 0xf7, 0xf0, 0xdd, 0xb5,
	0x1e, 0x52, 0x6e, 0x2a, 0x50, 0x86, 0x0f, 0x3f, 0x22, 0x7c, 0x23, 0x36, 0x22, 0x96, 0xca, 0xba,
	0xcd, 0xab, 0x31, 0x77, 0x69, 0x60, 0xe3, 0x1f, 0x06, 0x76, 0xf0, 0xf6, 0xaa, 0xa4, 0x85, 0xd6,
	0x83, 0x2f, 0x08, 0xf7, 0x62, 0x23, 0xfc, 0xf7, 0x98, 0xac, 0x39, 0xa1, 0x3b, 0x91, 0x3b, 0xf4,
	0xb5, 0x31, 0x04, 0xf4, 0x3f, 0x17, 0x17, 0x1a, 0xfc, 0x47, 0x78, 0xeb, 0x32, 0x2b, 0x7f, 0x89,
	0x5e, 0xf6, 0x82, 0xe0, 0xef, 0xde, 0x02, 0x7c, 0x78, 0x7c, 0x3e, 0x23, 0xe8, 0x62, 0x46, 0xd0,
	0xcf, 0x19, 0x41, 0x1f, 0xe6, 0xc4, 0xbb, 0x98, 0x13, 0xef, 0xdb, 0x9c, 0x78, 0xaf, 0xe8, 0x4a,
	0x84, 0xc9, 0x93, 0xe3, 0xe7, 0xe3, 0x17, 0x69, 0x32, 0x49, 0x8e, 0x92, 0x67, 0x14, 0x0a, 0xa9,
	0xc6, 0x0c, 0x34, 0xa7, 0x67, 0xb4, 0xfd, 0x6b, 0x9a, 0x3c, 0xf3, 0xbe, 0xfb, 0x60, 0x1f, 0xfc,
	0x0e, 0x00, 0x00, 0xff, 0xff, 0xe3, 0x4f, 0xf3, 0x81, 0x4a, 0x03, 0x00, 0x00,
}

func (this *MsgWithdrawCoinsToAccFromTreasury) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MsgWithdrawCoinsToAccFromTreasury)
	if !ok {
		that2, ok := that.(MsgWithdrawCoinsToAccFromTreasury)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Amount) != len(that1.Amount) {
		return false
	}
	for i := range this.Amount {
		if !this.Amount[i].Equal(&that1.Amount[i]) {
			return false
		}
	}
	if this.Receiver != that1.Receiver {
		return false
	}
	if this.Sender != that1.Sender {
		return false
	}
	return true
}
func (this *MsgMintCoins) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MsgMintCoins)
	if !ok {
		that2, ok := that.(MsgMintCoins)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Amount) != len(that1.Amount) {
		return false
	}
	for i := range this.Amount {
		if !this.Amount[i].Equal(&that1.Amount[i]) {
			return false
		}
	}
	if this.Sender != that1.Sender {
		return false
	}
	return true
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// WithdrawCoinsToAccFromTreasury defines a method for withdrawing from mint
	// module.
	WithdrawCoinsToAccFromTreasury(ctx context.Context, in *MsgWithdrawCoinsToAccFromTreasury, opts ...grpc.CallOption) (*MsgWithdrawCoinsToAccFromTreasuryResponse, error)
	// MintCoins defines a method for minting from mint module.
	MintCoins(ctx context.Context, in *MsgMintCoins, opts ...grpc.CallOption) (*MsgMintCoinsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) WithdrawCoinsToAccFromTreasury(ctx context.Context, in *MsgWithdrawCoinsToAccFromTreasury, opts ...grpc.CallOption) (*MsgWithdrawCoinsToAccFromTreasuryResponse, error) {
	out := new(MsgWithdrawCoinsToAccFromTreasuryResponse)
	err := c.cc.Invoke(ctx, "/mint.Msg/WithdrawCoinsToAccFromTreasury", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) MintCoins(ctx context.Context, in *MsgMintCoins, opts ...grpc.CallOption) (*MsgMintCoinsResponse, error) {
	out := new(MsgMintCoinsResponse)
	err := c.cc.Invoke(ctx, "/mint.Msg/MintCoins", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// WithdrawCoinsToAccFromTreasury defines a method for withdrawing from mint
	// module.
	WithdrawCoinsToAccFromTreasury(context.Context, *MsgWithdrawCoinsToAccFromTreasury) (*MsgWithdrawCoinsToAccFromTreasuryResponse, error)
	// MintCoins defines a method for minting from mint module.
	MintCoins(context.Context, *MsgMintCoins) (*MsgMintCoinsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) WithdrawCoinsToAccFromTreasury(ctx context.Context, req *MsgWithdrawCoinsToAccFromTreasury) (*MsgWithdrawCoinsToAccFromTreasuryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawCoinsToAccFromTreasury not implemented")
}
func (*UnimplementedMsgServer) MintCoins(ctx context.Context, req *MsgMintCoins) (*MsgMintCoinsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MintCoins not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_WithdrawCoinsToAccFromTreasury_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawCoinsToAccFromTreasury)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawCoinsToAccFromTreasury(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mint.Msg/WithdrawCoinsToAccFromTreasury",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawCoinsToAccFromTreasury(ctx, req.(*MsgWithdrawCoinsToAccFromTreasury))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_MintCoins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMintCoins)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).MintCoins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mint.Msg/MintCoins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).MintCoins(ctx, req.(*MsgMintCoins))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mint.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WithdrawCoinsToAccFromTreasury",
			Handler:    _Msg_WithdrawCoinsToAccFromTreasury_Handler,
		},
		{
			MethodName: "MintCoins",
			Handler:    _Msg_MintCoins_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mint/tx.proto",
}

func (m *MsgWithdrawCoinsToAccFromTreasury) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawCoinsToAccFromTreasury) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawCoinsToAccFromTreasury) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Receiver) > 0 {
		i -= len(m.Receiver)
		copy(dAtA[i:], m.Receiver)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Receiver)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgMintCoins) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgMintCoins) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgMintCoins) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MsgMintCoinsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgMintCoinsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgMintCoinsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgWithdrawCoinsToAccFromTreasury) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(m.Receiver)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgMintCoins) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgMintCoinsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgWithdrawCoinsToAccFromTreasury) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgWithdrawCoinsToAccFromTreasury: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawCoinsToAccFromTreasury: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receiver", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receiver = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgWithdrawCoinsToAccFromTreasuryResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgWithdrawCoinsToAccFromTreasuryResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawCoinsToAccFromTreasuryResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgMintCoins) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgMintCoins: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgMintCoins: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgMintCoinsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgMintCoinsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgMintCoinsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
