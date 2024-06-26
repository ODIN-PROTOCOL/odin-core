// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mint/v1beta1/params.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// Params holds parameters for the mint module.
type Params struct {
	// type of coin to mint
	MintDenom string `protobuf:"bytes,1,opt,name=mint_denom,json=mintDenom,proto3" json:"mint_denom,omitempty"`
	// maximum annual change in inflation rate
	InflationRateChange cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=inflation_rate_change,json=inflationRateChange,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"inflation_rate_change" yaml:"inflation_rate_change"`
	// maximum inflation rate
	InflationMax cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=inflation_max,json=inflationMax,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"inflation_max" yaml:"inflation_max"`
	// minimum inflation rate
	InflationMin cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=inflation_min,json=inflationMin,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"inflation_min" yaml:"inflation_min"`
	// goal of percent bonded atoms
	GoalBonded cosmossdk_io_math.LegacyDec `protobuf:"bytes,5,opt,name=goal_bonded,json=goalBonded,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"goal_bonded" yaml:"goal_bonded"`
	// expected blocks per year
	BlocksPerYear uint64 `protobuf:"varint,6,opt,name=blocks_per_year,json=blocksPerYear,proto3" json:"blocks_per_year,omitempty" yaml:"blocks_per_year"`
	// max amount to withdraw per time
	MaxWithdrawalPerTime github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,7,rep,name=max_withdrawal_per_time,json=maxWithdrawalPerTime,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"max_withdrawal_per_time" yaml:"max_withdrawal_per_time"`
	// Deprecated: map with smart contracts addresses
	IntegrationAddresses map[string]string `protobuf:"bytes,8,rep,name=integration_addresses,json=integrationAddresses,proto3" json:"integration_addresses,omitempty" yaml:"integration_addresses" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// flag if minting from air
	MintAir bool `protobuf:"varint,9,opt,name=mint_air,json=mintAir,proto3" json:"mint_air,omitempty" yaml:"mint_air"`
	// eligible to withdraw accounts
	EligibleAccountsPool []string `protobuf:"bytes,10,rep,name=eligible_accounts_pool,json=eligibleAccountsPool,proto3" json:"eligible_accounts_pool,omitempty" yaml:"eligible_accounts_pool"`
	// max allowed mint volume
	MaxAllowedMintVolume github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,11,rep,name=max_allowed_mint_volume,json=maxAllowedMintVolume,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"max_allowed_mint_volume" yaml:"max_allowed_mint_volume"`
	// allowed mint denoms
	AllowedMintDenoms []*AllowedDenom `protobuf:"bytes,12,rep,name=allowed_mint_denoms,json=allowedMintDenoms,proto3" json:"allowed_mint_denoms,omitempty" yaml:"allowed_mint_denoms"`
	// allowed minter
	AllowedMinter []string `protobuf:"bytes,13,rep,name=allowed_minter,json=allowedMinter,proto3" json:"allowed_minter,omitempty" yaml:"allowed_minter"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce136b324133acfc, []int{0}
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

func (m *Params) GetMintDenom() string {
	if m != nil {
		return m.MintDenom
	}
	return ""
}

func (m *Params) GetBlocksPerYear() uint64 {
	if m != nil {
		return m.BlocksPerYear
	}
	return 0
}

func (m *Params) GetMaxWithdrawalPerTime() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.MaxWithdrawalPerTime
	}
	return nil
}

func (m *Params) GetIntegrationAddresses() map[string]string {
	if m != nil {
		return m.IntegrationAddresses
	}
	return nil
}

func (m *Params) GetMintAir() bool {
	if m != nil {
		return m.MintAir
	}
	return false
}

func (m *Params) GetEligibleAccountsPool() []string {
	if m != nil {
		return m.EligibleAccountsPool
	}
	return nil
}

func (m *Params) GetMaxAllowedMintVolume() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.MaxAllowedMintVolume
	}
	return nil
}

func (m *Params) GetAllowedMintDenoms() []*AllowedDenom {
	if m != nil {
		return m.AllowedMintDenoms
	}
	return nil
}

func (m *Params) GetAllowedMinter() []string {
	if m != nil {
		return m.AllowedMinter
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "mint.v1beta1.Params")
	proto.RegisterMapType((map[string]string)(nil), "mint.v1beta1.Params.IntegrationAddressesEntry")
}

func init() { proto.RegisterFile("mint/v1beta1/params.proto", fileDescriptor_ce136b324133acfc) }

var fileDescriptor_ce136b324133acfc = []byte{
	// 805 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0x41, 0x8f, 0xdb, 0x44,
	0x14, 0xc7, 0xe3, 0xa6, 0xdd, 0x26, 0x93, 0x5d, 0x4a, 0xbd, 0x69, 0xeb, 0x04, 0x6a, 0x07, 0x9f,
	0x22, 0xa4, 0xd8, 0x2a, 0x5c, 0xd0, 0x9e, 0x88, 0x37, 0x08, 0xad, 0xb4, 0x6d, 0x22, 0xab, 0x6a,
	0x05, 0x17, 0x6b, 0x62, 0x0f, 0xce, 0x10, 0xcf, 0x4c, 0x34, 0xe3, 0xec, 0x26, 0x1f, 0x80, 0x0b,
	0x27, 0x8e, 0x1c, 0x7b, 0xe0, 0x80, 0x38, 0x81, 0xc4, 0x87, 0xe8, 0x71, 0xc5, 0x09, 0x71, 0x30,
	0x68, 0xf7, 0x00, 0xe7, 0x7c, 0x02, 0xe4, 0x19, 0x6f, 0x36, 0xbb, 0x78, 0xc5, 0x42, 0x2f, 0xc9,
	0x78, 0xde, 0xff, 0xfd, 0xfe, 0xcf, 0x33, 0xf3, 0xc6, 0xa0, 0x45, 0x30, 0x4d, 0xdd, 0xa3, 0x27,
	0x63, 0x94, 0xc2, 0x27, 0xee, 0x0c, 0x72, 0x48, 0x84, 0x33, 0xe3, 0x2c, 0x65, 0xfa, 0x76, 0x1e,
	0x72, 0x8a, 0x50, 0xbb, 0x19, 0xb3, 0x98, 0xc9, 0x80, 0x9b, 0x8f, 0x94, 0xa6, 0x6d, 0x86, 0x4c,
	0x10, 0x26, 0xdc, 0x31, 0x14, 0x68, 0x4d, 0x09, 0x19, 0xa6, 0x45, 0xfc, 0xd1, 0x25, 0xbc, 0x04,
	0xaa, 0x40, 0x4b, 0x25, 0x06, 0x8a, 0xa8, 0x1e, 0x8a, 0xd0, 0x7d, 0x48, 0x30, 0x65, 0xae, 0xfc,
	0x55, 0x53, 0xf6, 0x4f, 0x0d, 0xb0, 0x35, 0x92, 0xb5, 0xe9, 0x8f, 0x01, 0xc8, 0x31, 0x41, 0x84,
	0x28, 0x23, 0x86, 0xd6, 0xd1, 0xba, 0x75, 0xbf, 0x9e, 0xcf, 0x0c, 0xf2, 0x09, 0xfd, 0x6b, 0x0d,
	0x3c, 0xc0, 0xf4, 0x8b, 0x04, 0xa6, 0x98, 0xd1, 0x80, 0xc3, 0x14, 0x05, 0xe1, 0x04, 0xd2, 0x18,
	0x19, 0xb7, 0x72, 0xa9, 0xf7, 0xe2, 0x75, 0x66, 0x55, 0x7e, 0xcb, 0xac, 0x77, 0x94, 0xa5, 0x88,
	0xa6, 0x0e, 0x66, 0x2e, 0x81, 0xe9, 0xc4, 0x39, 0x44, 0x31, 0x0c, 0x97, 0x03, 0x14, 0xae, 0x32,
	0xeb, 0xdd, 0x25, 0x24, 0xc9, 0x9e, 0x5d, 0x4a, 0xb2, 0x7f, 0xf9, 0xb9, 0x07, 0x8a, 0x8a, 0x07,
	0x28, 0xfc, 0xfe, 0xcf, 0x1f, 0xdf, 0xd7, 0xfc, 0xdd, 0xb5, 0xd4, 0x87, 0x29, 0xda, 0x97, 0x42,
	0x5d, 0x80, 0x9d, 0x0b, 0x02, 0x81, 0x0b, 0xa3, 0x2a, 0x6b, 0x78, 0x76, 0xb3, 0x1a, 0x9a, 0x57,
	0x6b, 0x20, 0x70, 0x51, 0xea, 0xbd, 0xbd, 0x96, 0x3c, 0x85, 0x8b, 0x2b, 0xa6, 0x98, 0x1a, 0xb7,
	0xdf, 0xcc, 0x14, 0xd3, 0x7f, 0x33, 0xc5, 0x54, 0x27, 0xa0, 0x11, 0x33, 0x98, 0x04, 0x63, 0x46,
	0x23, 0x14, 0x19, 0x77, 0xa4, 0xe5, 0xe1, 0xcd, 0x2c, 0x75, 0x65, 0xb9, 0x91, 0x5f, 0x6a, 0x08,
	0x72, 0x81, 0x27, 0xe3, 0xba, 0x07, 0xee, 0x8d, 0x13, 0x16, 0x4e, 0x45, 0x30, 0x43, 0x3c, 0x58,
	0x22, 0xc8, 0x8d, 0xad, 0x8e, 0xd6, 0xbd, 0xed, 0xb5, 0x57, 0x99, 0xf5, 0x50, 0xf1, 0xae, 0x08,
	0x6c, 0x7f, 0x47, 0xcd, 0x8c, 0x10, 0xff, 0x0c, 0x41, 0xae, 0x7f, 0xa7, 0x81, 0x47, 0x04, 0x2e,
	0x82, 0x63, 0x9c, 0x4e, 0x22, 0x0e, 0x8f, 0x61, 0x22, 0xb5, 0x29, 0x26, 0xc8, 0xb8, 0xdb, 0xa9,
	0x76, 0x1b, 0x1f, 0xb4, 0x9c, 0xa2, 0x86, 0xfc, 0x74, 0x9f, 0x37, 0x82, 0xb3, 0xcf, 0x30, 0xf5,
	0xfc, 0xfc, 0xd5, 0x56, 0x99, 0x65, 0x2a, 0xaf, 0x6b, 0x38, 0xf6, 0x0f, 0xbf, 0x5b, 0xdd, 0x18,
	0xa7, 0x93, 0xf9, 0xd8, 0x09, 0x19, 0x29, 0x8e, 0x79, 0xf1, 0xd7, 0x13, 0xd1, 0xd4, 0x4d, 0x97,
	0x33, 0x24, 0x24, 0x52, 0xf8, 0x4d, 0x02, 0x17, 0x2f, 0xd7, 0x90, 0x11, 0xe2, 0xcf, 0x31, 0x41,
	0xfa, 0x57, 0xf2, 0x40, 0xa7, 0x28, 0xe6, 0x6a, 0x3f, 0x60, 0x14, 0x71, 0x24, 0x04, 0x12, 0x46,
	0x4d, 0x16, 0xe9, 0x38, 0x9b, 0x6d, 0xea, 0xa8, 0x2e, 0x71, 0x0e, 0x2e, 0x32, 0xfa, 0xe7, 0x09,
	0x9f, 0xd0, 0x94, 0x2f, 0xbd, 0xce, 0xe6, 0xe9, 0x2e, 0xc1, 0xda, 0x7e, 0x13, 0x97, 0x24, 0xeb,
	0x0e, 0xa8, 0xc9, 0xbe, 0x83, 0x98, 0x1b, 0xf5, 0x8e, 0xd6, 0xad, 0x79, 0xbb, 0xab, 0xcc, 0xba,
	0x57, 0xbc, 0x7f, 0x11, 0xb1, 0xfd, 0xbb, 0xf9, 0xb0, 0x8f, 0xb9, 0xfe, 0x12, 0x3c, 0x44, 0x09,
	0x8e, 0xf1, 0x38, 0x41, 0x01, 0x0c, 0x43, 0x36, 0xa7, 0xa9, 0x08, 0x66, 0x8c, 0x25, 0x06, 0xe8,
	0x54, 0xbb, 0x75, 0xef, 0xbd, 0x55, 0x66, 0x3d, 0x56, 0xd9, 0xe5, 0x3a, 0xdb, 0x6f, 0x9e, 0x07,
	0xfa, 0xc5, 0xfc, 0x88, 0xb1, 0x64, 0xbd, 0x6f, 0x30, 0x49, 0xd8, 0x31, 0x8a, 0x02, 0xe9, 0x7d,
	0xc4, 0x92, 0x39, 0x41, 0x46, 0xe3, 0x7f, 0xec, 0x5b, 0x09, 0xe7, 0xbf, 0xef, 0x5b, 0x5f, 0x41,
	0x9e, 0x62, 0x9a, 0xbe, 0x90, 0x08, 0xfd, 0x4b, 0xb0, 0x7b, 0x89, 0x2c, 0xef, 0x2b, 0x61, 0x6c,
	0xcb, 0x0a, 0xdb, 0x97, 0x37, 0xad, 0xc8, 0x96, 0x37, 0x98, 0x67, 0xae, 0x32, 0xab, 0xad, 0xca,
	0x2b, 0x01, 0xd8, 0xfe, 0x7d, 0x78, 0xe1, 0x25, 0x33, 0x84, 0xfe, 0x31, 0x78, 0x6b, 0x53, 0x8a,
	0xb8, 0xb1, 0x23, 0xd7, 0xb8, 0xb5, 0xca, 0xac, 0x07, 0xff, 0x44, 0xa1, 0xbc, 0x19, 0x36, 0x28,
	0x88, 0xb7, 0x3f, 0x05, 0xad, 0x6b, 0x8f, 0x8c, 0xfe, 0x36, 0xa8, 0x4e, 0xd1, 0xb2, 0xb8, 0x6b,
	0xf3, 0xa1, 0xde, 0x04, 0x77, 0x8e, 0x60, 0x32, 0x2f, 0x2e, 0x55, 0x5f, 0x3d, 0xec, 0xdd, 0xfa,
	0x48, 0xdb, 0xab, 0x7d, 0xfb, 0xca, 0xaa, 0xfc, 0xf5, 0xca, 0xd2, 0xbc, 0x83, 0xd7, 0xa7, 0xa6,
	0x76, 0x72, 0x6a, 0x6a, 0x7f, 0x9c, 0x9a, 0xda, 0x37, 0x67, 0x66, 0xe5, 0xe4, 0xcc, 0xac, 0xfc,
	0x7a, 0x66, 0x56, 0x3e, 0x77, 0x37, 0x96, 0x76, 0x38, 0x38, 0x78, 0xd6, 0x1b, 0xf9, 0xc3, 0xe7,
	0xc3, 0xfd, 0xe1, 0xa1, 0xcb, 0x22, 0x4c, 0x7b, 0x21, 0xe3, 0xc8, 0x5d, 0xc8, 0x8f, 0x85, 0x5a,
	0xe7, 0xf1, 0x96, 0xfc, 0x0a, 0x7c, 0xf8, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7c, 0x98, 0xaf,
	0x12, 0xad, 0x06, 0x00, 0x00,
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
	if this.MintDenom != that1.MintDenom {
		return false
	}
	if !this.InflationRateChange.Equal(that1.InflationRateChange) {
		return false
	}
	if !this.InflationMax.Equal(that1.InflationMax) {
		return false
	}
	if !this.InflationMin.Equal(that1.InflationMin) {
		return false
	}
	if !this.GoalBonded.Equal(that1.GoalBonded) {
		return false
	}
	if this.BlocksPerYear != that1.BlocksPerYear {
		return false
	}
	if len(this.MaxWithdrawalPerTime) != len(that1.MaxWithdrawalPerTime) {
		return false
	}
	for i := range this.MaxWithdrawalPerTime {
		if !this.MaxWithdrawalPerTime[i].Equal(&that1.MaxWithdrawalPerTime[i]) {
			return false
		}
	}
	if len(this.IntegrationAddresses) != len(that1.IntegrationAddresses) {
		return false
	}
	for i := range this.IntegrationAddresses {
		if this.IntegrationAddresses[i] != that1.IntegrationAddresses[i] {
			return false
		}
	}
	if this.MintAir != that1.MintAir {
		return false
	}
	if len(this.EligibleAccountsPool) != len(that1.EligibleAccountsPool) {
		return false
	}
	for i := range this.EligibleAccountsPool {
		if this.EligibleAccountsPool[i] != that1.EligibleAccountsPool[i] {
			return false
		}
	}
	if len(this.MaxAllowedMintVolume) != len(that1.MaxAllowedMintVolume) {
		return false
	}
	for i := range this.MaxAllowedMintVolume {
		if !this.MaxAllowedMintVolume[i].Equal(&that1.MaxAllowedMintVolume[i]) {
			return false
		}
	}
	if len(this.AllowedMintDenoms) != len(that1.AllowedMintDenoms) {
		return false
	}
	for i := range this.AllowedMintDenoms {
		if !this.AllowedMintDenoms[i].Equal(that1.AllowedMintDenoms[i]) {
			return false
		}
	}
	if len(this.AllowedMinter) != len(that1.AllowedMinter) {
		return false
	}
	for i := range this.AllowedMinter {
		if this.AllowedMinter[i] != that1.AllowedMinter[i] {
			return false
		}
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
	if len(m.AllowedMinter) > 0 {
		for iNdEx := len(m.AllowedMinter) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AllowedMinter[iNdEx])
			copy(dAtA[i:], m.AllowedMinter[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.AllowedMinter[iNdEx])))
			i--
			dAtA[i] = 0x6a
		}
	}
	if len(m.AllowedMintDenoms) > 0 {
		for iNdEx := len(m.AllowedMintDenoms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AllowedMintDenoms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x62
		}
	}
	if len(m.MaxAllowedMintVolume) > 0 {
		for iNdEx := len(m.MaxAllowedMintVolume) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MaxAllowedMintVolume[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.EligibleAccountsPool) > 0 {
		for iNdEx := len(m.EligibleAccountsPool) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.EligibleAccountsPool[iNdEx])
			copy(dAtA[i:], m.EligibleAccountsPool[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.EligibleAccountsPool[iNdEx])))
			i--
			dAtA[i] = 0x52
		}
	}
	if m.MintAir {
		i--
		if m.MintAir {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x48
	}
	if len(m.IntegrationAddresses) > 0 {
		for k := range m.IntegrationAddresses {
			v := m.IntegrationAddresses[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintParams(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintParams(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintParams(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.MaxWithdrawalPerTime) > 0 {
		for iNdEx := len(m.MaxWithdrawalPerTime) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MaxWithdrawalPerTime[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.BlocksPerYear != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.BlocksPerYear))
		i--
		dAtA[i] = 0x30
	}
	{
		size := m.GoalBonded.Size()
		i -= size
		if _, err := m.GoalBonded.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.InflationMin.Size()
		i -= size
		if _, err := m.InflationMin.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.InflationMax.Size()
		i -= size
		if _, err := m.InflationMax.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.InflationRateChange.Size()
		i -= size
		if _, err := m.InflationRateChange.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.MintDenom) > 0 {
		i -= len(m.MintDenom)
		copy(dAtA[i:], m.MintDenom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.MintDenom)))
		i--
		dAtA[i] = 0xa
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
	l = len(m.MintDenom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.InflationRateChange.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InflationMax.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InflationMin.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.GoalBonded.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.BlocksPerYear != 0 {
		n += 1 + sovParams(uint64(m.BlocksPerYear))
	}
	if len(m.MaxWithdrawalPerTime) > 0 {
		for _, e := range m.MaxWithdrawalPerTime {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.IntegrationAddresses) > 0 {
		for k, v := range m.IntegrationAddresses {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovParams(uint64(len(k))) + 1 + len(v) + sovParams(uint64(len(v)))
			n += mapEntrySize + 1 + sovParams(uint64(mapEntrySize))
		}
	}
	if m.MintAir {
		n += 2
	}
	if len(m.EligibleAccountsPool) > 0 {
		for _, s := range m.EligibleAccountsPool {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.MaxAllowedMintVolume) > 0 {
		for _, e := range m.MaxAllowedMintVolume {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.AllowedMintDenoms) > 0 {
		for _, e := range m.AllowedMintDenoms {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.AllowedMinter) > 0 {
		for _, s := range m.AllowedMinter {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MintDenom", wireType)
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
			m.MintDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InflationRateChange", wireType)
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
			if err := m.InflationRateChange.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InflationMax", wireType)
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
			if err := m.InflationMax.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InflationMin", wireType)
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
			if err := m.InflationMin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GoalBonded", wireType)
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
			if err := m.GoalBonded.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlocksPerYear", wireType)
			}
			m.BlocksPerYear = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlocksPerYear |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxWithdrawalPerTime", wireType)
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
			m.MaxWithdrawalPerTime = append(m.MaxWithdrawalPerTime, types.Coin{})
			if err := m.MaxWithdrawalPerTime[len(m.MaxWithdrawalPerTime)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IntegrationAddresses", wireType)
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
			if m.IntegrationAddresses == nil {
				m.IntegrationAddresses = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowParams
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthParams
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthParams
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowParams
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthParams
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthParams
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipParams(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthParams
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.IntegrationAddresses[mapkey] = mapvalue
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MintAir", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.MintAir = bool(v != 0)
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EligibleAccountsPool", wireType)
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
			m.EligibleAccountsPool = append(m.EligibleAccountsPool, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxAllowedMintVolume", wireType)
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
			m.MaxAllowedMintVolume = append(m.MaxAllowedMintVolume, types.Coin{})
			if err := m.MaxAllowedMintVolume[len(m.MaxAllowedMintVolume)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllowedMintDenoms", wireType)
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
			m.AllowedMintDenoms = append(m.AllowedMintDenoms, &AllowedDenom{})
			if err := m.AllowedMintDenoms[len(m.AllowedMintDenoms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllowedMinter", wireType)
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
			m.AllowedMinter = append(m.AllowedMinter, string(dAtA[iNdEx:postIndex]))
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
