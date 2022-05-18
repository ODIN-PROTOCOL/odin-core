// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: oracle/v1/params.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

// Params is the data structure that keeps the parameters of the oracle module.
type Params struct {
	// MaxRawRequestCount is the maximum number of data source raw requests a
	// request can make.
	MaxRawRequestCount uint64 `protobuf:"varint,1,opt,name=max_raw_request_count,json=maxRawRequestCount,proto3" json:"max_raw_request_count,omitempty"`
	// MaxAskCount is the maximum number of validators a request can target.
	MaxAskCount uint64 `protobuf:"varint,2,opt,name=max_ask_count,json=maxAskCount,proto3" json:"max_ask_count,omitempty"`
	// ExpirationBlockCount is the number of blocks a request stays valid before
	// it gets expired due to insufficient reports.
	ExpirationBlockCount uint64 `protobuf:"varint,3,opt,name=expiration_block_count,json=expirationBlockCount,proto3" json:"expiration_block_count,omitempty"`
	// BaseOwasmGas is the base amount of Cosmos-SDK gas charged for owasm
	// execution.
	BaseOwasmGas uint64 `protobuf:"varint,4,opt,name=base_owasm_gas,json=baseOwasmGas,proto3" json:"base_owasm_gas,omitempty"`
	// PerValidatorRequestGas is the amount of Cosmos-SDK gas charged per
	// requested validator.
	PerValidatorRequestGas uint64 `protobuf:"varint,5,opt,name=per_validator_request_gas,json=perValidatorRequestGas,proto3" json:"per_validator_request_gas,omitempty"`
	// SamplingTryCount the number of validator sampling tries to pick the highest
	// voting power subset of validators to perform an oracle task.
	SamplingTryCount uint64 `protobuf:"varint,6,opt,name=sampling_try_count,json=samplingTryCount,proto3" json:"sampling_try_count,omitempty"`
	// OracleRewardPercentage is the percentage of block rewards allocated to
	// active oracle validators.
	OracleRewardPercentage uint64 `protobuf:"varint,7,opt,name=oracle_reward_percentage,json=oracleRewardPercentage,proto3" json:"oracle_reward_percentage,omitempty"`
	// InactivePenaltyDuration is the duration period where a validator cannot
	// activate back after missing an oracle report.
	InactivePenaltyDuration uint64 `protobuf:"varint,8,opt,name=inactive_penalty_duration,json=inactivePenaltyDuration,proto3" json:"inactive_penalty_duration,omitempty"`
	// MaxDataSize is the maximum number of bytes that can be present in the
	// report as the result
	MaxDataSize uint64 `protobuf:"varint,9,opt,name=max_data_size,json=maxDataSize,proto3" json:"max_data_size,omitempty"`
	// MaxCalldataSize is the maximum number of bytes that can be present in the
	// calldata
	MaxCalldataSize uint64 `protobuf:"varint,10,opt,name=max_calldata_size,json=maxCalldataSize,proto3" json:"max_calldata_size,omitempty"`
	// DataProviderRewardPerByte is the amount of tokens, user gets for the byte of data provided
	DataProviderRewardPerByte github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,11,rep,name=data_provider_reward_per_byte,json=dataProviderRewardPerByte,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"data_provider_reward_per_byte"`
	// DataProviderRewardThreshold is the maximum amount of tokens that can be paid for data per time
	DataProviderRewardThreshold RewardThreshold `protobuf:"bytes,12,opt,name=data_provider_reward_threshold,json=dataProviderRewardThreshold,proto3" json:"data_provider_reward_threshold"`
	// RewardDecreasingFraction is the percentage by which the cost of data per byte is reduced when the limit is reached
	RewardDecreasingFraction github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,13,opt,name=reward_decreasing_fraction,json=rewardDecreasingFraction,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reward_decreasing_fraction"`
	// Denominations that can be used for withdrawing fee from data requesters
	DataRequesterFeeDenoms []string `protobuf:"bytes,14,rep,name=data_requester_fee_denoms,json=dataRequesterFeeDenoms,proto3" json:"data_requester_fee_denoms,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_d7000dc69c8e604b, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetMaxRawRequestCount() uint64 {
	if m != nil {
		return m.MaxRawRequestCount
	}
	return 0
}

func (m *Params) GetMaxAskCount() uint64 {
	if m != nil {
		return m.MaxAskCount
	}
	return 0
}

func (m *Params) GetExpirationBlockCount() uint64 {
	if m != nil {
		return m.ExpirationBlockCount
	}
	return 0
}

func (m *Params) GetBaseOwasmGas() uint64 {
	if m != nil {
		return m.BaseOwasmGas
	}
	return 0
}

func (m *Params) GetPerValidatorRequestGas() uint64 {
	if m != nil {
		return m.PerValidatorRequestGas
	}
	return 0
}

func (m *Params) GetSamplingTryCount() uint64 {
	if m != nil {
		return m.SamplingTryCount
	}
	return 0
}

func (m *Params) GetOracleRewardPercentage() uint64 {
	if m != nil {
		return m.OracleRewardPercentage
	}
	return 0
}

func (m *Params) GetInactivePenaltyDuration() uint64 {
	if m != nil {
		return m.InactivePenaltyDuration
	}
	return 0
}

func (m *Params) GetMaxDataSize() uint64 {
	if m != nil {
		return m.MaxDataSize
	}
	return 0
}

func (m *Params) GetMaxCalldataSize() uint64 {
	if m != nil {
		return m.MaxCalldataSize
	}
	return 0
}

func (m *Params) GetDataProviderRewardPerByte() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.DataProviderRewardPerByte
	}
	return nil
}

func (m *Params) GetDataProviderRewardThreshold() RewardThreshold {
	if m != nil {
		return m.DataProviderRewardThreshold
	}
	return RewardThreshold{}
}

func (m *Params) GetDataRequesterFeeDenoms() []string {
	if m != nil {
		return m.DataRequesterFeeDenoms
	}
	return nil
}

// RewardThreshold
type RewardThreshold struct {
	// Amount is the maximum amount of tokens that can be paid for data
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	// Blocks is the number of blocks during which the sum of the reward should not exceed total_reward_amount
	Blocks uint64 `protobuf:"varint,2,opt,name=blocks,proto3" json:"blocks,omitempty"`
}

func (m *RewardThreshold) Reset()         { *m = RewardThreshold{} }
func (m *RewardThreshold) String() string { return proto.CompactTextString(m) }
func (*RewardThreshold) ProtoMessage()    {}
func (*RewardThreshold) Descriptor() ([]byte, []int) {
	return fileDescriptor_d7000dc69c8e604b, []int{1}
}
func (m *RewardThreshold) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RewardThreshold) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RewardThreshold.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RewardThreshold) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardThreshold.Merge(m, src)
}
func (m *RewardThreshold) XXX_Size() int {
	return m.Size()
}
func (m *RewardThreshold) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardThreshold.DiscardUnknown(m)
}

var xxx_messageInfo_RewardThreshold proto.InternalMessageInfo

func (m *RewardThreshold) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *RewardThreshold) GetBlocks() uint64 {
	if m != nil {
		return m.Blocks
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "oracle.v1.Params")
	proto.RegisterType((*RewardThreshold)(nil), "oracle.v1.RewardThreshold")
}

func init() { proto.RegisterFile("oracle/v1/params.proto", fileDescriptor_d7000dc69c8e604b) }

var fileDescriptor_d7000dc69c8e604b = []byte{
	// 698 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x4d, 0x6f, 0x13, 0x49,
	0x10, 0xf5, 0x6c, 0xbc, 0x4e, 0xd2, 0xce, 0xc7, 0xee, 0x28, 0xeb, 0x1d, 0x7b, 0xb5, 0x63, 0x2b,
	0x5a, 0xad, 0xac, 0xd5, 0x66, 0x06, 0x07, 0x0e, 0x90, 0x1b, 0xb6, 0x15, 0x84, 0x40, 0xd8, 0x1a,
	0x22, 0x0e, 0x5c, 0x5a, 0xed, 0x99, 0x8a, 0x33, 0xca, 0xcc, 0xf4, 0xd0, 0xdd, 0xfe, 0xca, 0x7f,
	0x40, 0xe2, 0x90, 0x03, 0xc7, 0x9c, 0xf9, 0x25, 0x39, 0xe6, 0x88, 0x38, 0x04, 0x94, 0x5c, 0xf8,
	0x19, 0xa8, 0x3f, 0xec, 0x44, 0xc0, 0x81, 0x03, 0x27, 0x7b, 0xea, 0xbd, 0x57, 0xf5, 0xaa, 0xba,
	0xba, 0x51, 0x85, 0x32, 0x12, 0x26, 0xe0, 0x8f, 0x5b, 0x7e, 0x4e, 0x18, 0x49, 0xb9, 0x97, 0x33,
	0x2a, 0xa8, 0xbd, 0xaa, 0xe3, 0xde, 0xb8, 0x55, 0xdb, 0x1a, 0xd2, 0x21, 0x55, 0x51, 0x5f, 0xfe,
	0xd3, 0x84, 0x9a, 0x1b, 0x52, 0x9e, 0x52, 0xee, 0x0f, 0x08, 0x97, 0xea, 0x01, 0x08, 0xd2, 0xf2,
	0x43, 0x1a, 0x67, 0x1a, 0xdf, 0x3e, 0x5d, 0x46, 0xa5, 0xbe, 0xca, 0x68, 0xb7, 0xd0, 0x1f, 0x29,
	0x99, 0x62, 0x46, 0x26, 0x98, 0xc1, 0xab, 0x11, 0x70, 0x81, 0x43, 0x3a, 0xca, 0x84, 0x63, 0x35,
	0xac, 0x66, 0x31, 0xb0, 0x53, 0x32, 0x0d, 0xc8, 0x24, 0xd0, 0x50, 0x47, 0x22, 0xf6, 0x36, 0x5a,
	0x97, 0x12, 0xc2, 0x8f, 0x0d, 0xf5, 0x17, 0x45, 0x2d, 0xa7, 0x64, 0xfa, 0x90, 0x1f, 0x6b, 0xce,
	0x3d, 0x54, 0x81, 0x69, 0x1e, 0x33, 0x22, 0x62, 0x9a, 0xe1, 0x41, 0x42, 0xc3, 0x39, 0x79, 0x49,
	0x91, 0xb7, 0x6e, 0xd0, 0xb6, 0x04, 0xb5, 0xea, 0x1f, 0xb4, 0x21, 0x2d, 0x63, 0x3a, 0x21, 0x3c,
	0xc5, 0x43, 0xc2, 0x9d, 0xa2, 0x62, 0xaf, 0xc9, 0x68, 0x4f, 0x06, 0x1f, 0x11, 0x6e, 0x3f, 0x40,
	0xd5, 0x1c, 0x18, 0x1e, 0x93, 0x24, 0x8e, 0x88, 0xa0, 0x6c, 0x61, 0x5c, 0x0a, 0x7e, 0x55, 0x82,
	0x4a, 0x0e, 0xec, 0xc5, 0x1c, 0x37, 0xe6, 0xa5, 0xf4, 0x7f, 0x64, 0x73, 0x92, 0xe6, 0x49, 0x9c,
	0x0d, 0xb1, 0x60, 0x33, 0x63, 0xa9, 0xa4, 0x34, 0xbf, 0xcd, 0x91, 0x03, 0x36, 0xd3, 0x76, 0xee,
	0x23, 0x47, 0x4f, 0x1a, 0x33, 0x98, 0x10, 0x16, 0xe1, 0x1c, 0x58, 0x08, 0x99, 0x20, 0x43, 0x70,
	0x96, 0x75, 0x1d, 0x8d, 0x07, 0x0a, 0xee, 0x2f, 0x50, 0x7b, 0x0f, 0x55, 0xe3, 0x8c, 0x84, 0x22,
	0x1e, 0x03, 0xce, 0x21, 0x23, 0x89, 0x98, 0xe1, 0x68, 0xa4, 0xfb, 0x75, 0x56, 0x94, 0xf4, 0xcf,
	0x39, 0xa1, 0xaf, 0xf1, 0xae, 0x81, 0xe7, 0xe3, 0x8d, 0x88, 0x20, 0x98, 0xc7, 0x27, 0xe0, 0xac,
	0x2e, 0xc6, 0xdb, 0x25, 0x82, 0x3c, 0x8f, 0x4f, 0xc0, 0xfe, 0x0f, 0xfd, 0x2e, 0x39, 0x21, 0x49,
	0x92, 0x1b, 0x1e, 0x52, 0xbc, 0xcd, 0x94, 0x4c, 0x3b, 0x26, 0xae, 0xb8, 0xaf, 0x2d, 0xf4, 0xb7,
	0x22, 0xe5, 0x8c, 0x8e, 0xe3, 0x08, 0xd8, 0xad, 0x6e, 0xf0, 0x60, 0x26, 0xc0, 0x29, 0x37, 0x96,
	0x9a, 0xe5, 0xdd, 0xaa, 0xa7, 0xb7, 0xc6, 0x93, 0xc3, 0xf6, 0xcc, 0xd6, 0x78, 0x1d, 0x1a, 0x67,
	0xed, 0x3b, 0xe7, 0x97, 0xf5, 0xc2, 0xbb, 0x8f, 0xf5, 0xe6, 0x30, 0x16, 0x47, 0xa3, 0x81, 0x17,
	0xd2, 0xd4, 0x37, 0x2b, 0xa6, 0x7f, 0x76, 0x78, 0x74, 0xec, 0x8b, 0x59, 0x0e, 0x5c, 0x09, 0x78,
	0x50, 0x95, 0x15, 0xfb, 0xa6, 0xe0, 0x62, 0x3c, 0xed, 0x99, 0x00, 0x1b, 0x90, 0xfb, 0x5d, 0x3b,
	0xe2, 0x88, 0x01, 0x3f, 0xa2, 0x49, 0xe4, 0xac, 0x35, 0xac, 0x66, 0x79, 0xb7, 0xe6, 0x2d, 0xd6,
	0xdc, 0xd3, 0x19, 0x0e, 0xe6, 0x8c, 0x76, 0x51, 0x1a, 0x0a, 0xfe, 0xfa, 0xb6, 0xc8, 0x82, 0x62,
	0x27, 0xa8, 0x66, 0x12, 0x47, 0x10, 0x32, 0x20, 0x5c, 0x9e, 0xf9, 0x21, 0x93, 0x33, 0xa7, 0x99,
	0xb3, 0xde, 0xb0, 0x9a, 0x6b, 0x6d, 0x4f, 0xa6, 0xf9, 0x70, 0x59, 0xff, 0xf7, 0x07, 0xfa, 0xea,
	0x42, 0x18, 0x38, 0x3a, 0x63, 0x77, 0x91, 0x70, 0xdf, 0xe4, 0x93, 0x3b, 0xa9, 0x9a, 0x32, 0xab,
	0x08, 0x0c, 0x1f, 0x02, 0xe0, 0x08, 0x32, 0x9a, 0x72, 0x67, 0xa3, 0xb1, 0xd4, 0x5c, 0x0d, 0x2a,
	0x92, 0x10, 0xcc, 0xf1, 0x7d, 0x80, 0xae, 0x42, 0xf7, 0x56, 0xde, 0x9e, 0xd5, 0x0b, 0x9f, 0xcf,
	0xea, 0xd6, 0xf6, 0xa9, 0x85, 0x36, 0xbf, 0x6e, 0x23, 0x44, 0x25, 0x92, 0x9a, 0x0b, 0xf9, 0xd3,
	0x4f, 0xc9, 0xa4, 0xb6, 0x2b, 0xa8, 0xa4, 0xae, 0x28, 0x37, 0x57, 0xd9, 0x7c, 0xed, 0x15, 0xa5,
	0xad, 0xf6, 0x93, 0xf3, 0x2b, 0xd7, 0xba, 0xb8, 0x72, 0xad, 0x4f, 0x57, 0xae, 0xf5, 0xe6, 0xda,
	0x2d, 0x5c, 0x5c, 0xbb, 0x85, 0xf7, 0xd7, 0x6e, 0xe1, 0x65, 0xeb, 0x56, 0xa5, 0x5e, 0xf7, 0xf1,
	0xb3, 0x9d, 0x7e, 0xd0, 0x3b, 0xe8, 0x75, 0x7a, 0x4f, 0x7d, 0x1a, 0xc5, 0xd9, 0x4e, 0x48, 0x19,
	0xf8, 0x53, 0xdf, 0xbc, 0x62, 0xaa, 0xf0, 0xa0, 0xa4, 0x5e, 0xa0, 0xbb, 0x5f, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x7b, 0xf1, 0xda, 0x16, 0xdc, 0x04, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
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
	if this.MaxRawRequestCount != that1.MaxRawRequestCount {
		return false
	}
	if this.MaxAskCount != that1.MaxAskCount {
		return false
	}
	if this.ExpirationBlockCount != that1.ExpirationBlockCount {
		return false
	}
	if this.BaseOwasmGas != that1.BaseOwasmGas {
		return false
	}
	if this.PerValidatorRequestGas != that1.PerValidatorRequestGas {
		return false
	}
	if this.SamplingTryCount != that1.SamplingTryCount {
		return false
	}
	if this.OracleRewardPercentage != that1.OracleRewardPercentage {
		return false
	}
	if this.InactivePenaltyDuration != that1.InactivePenaltyDuration {
		return false
	}
	if this.MaxDataSize != that1.MaxDataSize {
		return false
	}
	if this.MaxCalldataSize != that1.MaxCalldataSize {
		return false
	}
	if len(this.DataProviderRewardPerByte) != len(that1.DataProviderRewardPerByte) {
		return false
	}
	for i := range this.DataProviderRewardPerByte {
		if !this.DataProviderRewardPerByte[i].Equal(&that1.DataProviderRewardPerByte[i]) {
			return false
		}
	}
	if !this.DataProviderRewardThreshold.Equal(&that1.DataProviderRewardThreshold) {
		return false
	}
	if !this.RewardDecreasingFraction.Equal(that1.RewardDecreasingFraction) {
		return false
	}
	if len(this.DataRequesterFeeDenoms) != len(that1.DataRequesterFeeDenoms) {
		return false
	}
	for i := range this.DataRequesterFeeDenoms {
		if this.DataRequesterFeeDenoms[i] != that1.DataRequesterFeeDenoms[i] {
			return false
		}
	}
	return true
}
func (this *RewardThreshold) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RewardThreshold)
	if !ok {
		that2, ok := that.(RewardThreshold)
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
	if this.Blocks != that1.Blocks {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DataRequesterFeeDenoms) > 0 {
		for iNdEx := len(m.DataRequesterFeeDenoms) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.DataRequesterFeeDenoms[iNdEx])
			copy(dAtA[i:], m.DataRequesterFeeDenoms[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.DataRequesterFeeDenoms[iNdEx])))
			i--
			dAtA[i] = 0x72
		}
	}
	{
		size := m.RewardDecreasingFraction.Size()
		i -= size
		if _, err := m.RewardDecreasingFraction.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x6a
	{
		size, err := m.DataProviderRewardThreshold.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	if len(m.DataProviderRewardPerByte) > 0 {
		for iNdEx := len(m.DataProviderRewardPerByte) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DataProviderRewardPerByte[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x5a
		}
	}
	if m.MaxCalldataSize != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxCalldataSize))
		i--
		dAtA[i] = 0x50
	}
	if m.MaxDataSize != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxDataSize))
		i--
		dAtA[i] = 0x48
	}
	if m.InactivePenaltyDuration != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.InactivePenaltyDuration))
		i--
		dAtA[i] = 0x40
	}
	if m.OracleRewardPercentage != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.OracleRewardPercentage))
		i--
		dAtA[i] = 0x38
	}
	if m.SamplingTryCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SamplingTryCount))
		i--
		dAtA[i] = 0x30
	}
	if m.PerValidatorRequestGas != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.PerValidatorRequestGas))
		i--
		dAtA[i] = 0x28
	}
	if m.BaseOwasmGas != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.BaseOwasmGas))
		i--
		dAtA[i] = 0x20
	}
	if m.ExpirationBlockCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.ExpirationBlockCount))
		i--
		dAtA[i] = 0x18
	}
	if m.MaxAskCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxAskCount))
		i--
		dAtA[i] = 0x10
	}
	if m.MaxRawRequestCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxRawRequestCount))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *RewardThreshold) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RewardThreshold) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RewardThreshold) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Blocks != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Blocks))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxRawRequestCount != 0 {
		n += 1 + sovParams(uint64(m.MaxRawRequestCount))
	}
	if m.MaxAskCount != 0 {
		n += 1 + sovParams(uint64(m.MaxAskCount))
	}
	if m.ExpirationBlockCount != 0 {
		n += 1 + sovParams(uint64(m.ExpirationBlockCount))
	}
	if m.BaseOwasmGas != 0 {
		n += 1 + sovParams(uint64(m.BaseOwasmGas))
	}
	if m.PerValidatorRequestGas != 0 {
		n += 1 + sovParams(uint64(m.PerValidatorRequestGas))
	}
	if m.SamplingTryCount != 0 {
		n += 1 + sovParams(uint64(m.SamplingTryCount))
	}
	if m.OracleRewardPercentage != 0 {
		n += 1 + sovParams(uint64(m.OracleRewardPercentage))
	}
	if m.InactivePenaltyDuration != 0 {
		n += 1 + sovParams(uint64(m.InactivePenaltyDuration))
	}
	if m.MaxDataSize != 0 {
		n += 1 + sovParams(uint64(m.MaxDataSize))
	}
	if m.MaxCalldataSize != 0 {
		n += 1 + sovParams(uint64(m.MaxCalldataSize))
	}
	if len(m.DataProviderRewardPerByte) > 0 {
		for _, e := range m.DataProviderRewardPerByte {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	l = m.DataProviderRewardThreshold.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.RewardDecreasingFraction.Size()
	n += 1 + l + sovParams(uint64(l))
	if len(m.DataRequesterFeeDenoms) > 0 {
		for _, s := range m.DataRequesterFeeDenoms {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func (m *RewardThreshold) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if m.Blocks != 0 {
		n += 1 + sovParams(uint64(m.Blocks))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxRawRequestCount", wireType)
			}
			m.MaxRawRequestCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxRawRequestCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxAskCount", wireType)
			}
			m.MaxAskCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxAskCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpirationBlockCount", wireType)
			}
			m.ExpirationBlockCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpirationBlockCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseOwasmGas", wireType)
			}
			m.BaseOwasmGas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BaseOwasmGas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PerValidatorRequestGas", wireType)
			}
			m.PerValidatorRequestGas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PerValidatorRequestGas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SamplingTryCount", wireType)
			}
			m.SamplingTryCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SamplingTryCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OracleRewardPercentage", wireType)
			}
			m.OracleRewardPercentage = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OracleRewardPercentage |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InactivePenaltyDuration", wireType)
			}
			m.InactivePenaltyDuration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.InactivePenaltyDuration |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxDataSize", wireType)
			}
			m.MaxDataSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxDataSize |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxCalldataSize", wireType)
			}
			m.MaxCalldataSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxCalldataSize |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataProviderRewardPerByte", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataProviderRewardPerByte = append(m.DataProviderRewardPerByte, types.Coin{})
			if err := m.DataProviderRewardPerByte[len(m.DataProviderRewardPerByte)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataProviderRewardThreshold", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DataProviderRewardThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardDecreasingFraction", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardDecreasingFraction.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataRequesterFeeDenoms", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataRequesterFeeDenoms = append(m.DataRequesterFeeDenoms, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *RewardThreshold) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: RewardThreshold: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RewardThreshold: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Blocks", wireType)
			}
			m.Blocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Blocks |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
