// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: auction/params.proto

package types

import (
	fmt "fmt"
	types1 "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
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

// Params is the data structure that keeps the parameters of the auction module.
type Params struct {
	// AuctionStartThreshold is the threshold at which the auction starts
	AuctionStartThreshold github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=auction_start_threshold,json=auctionStartThreshold,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"auction_start_threshold"`
	// ExchangeRate is a rate for buying coins throw the auction
	ExchangeRates []types1.Exchange `protobuf:"bytes,2,rep,name=exchange_rates,json=exchangeRates,proto3" json:"exchange_rates"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96a95233ccbd0c2, []int{0}
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

func (m *Params) GetAuctionStartThreshold() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.AuctionStartThreshold
	}
	return nil
}

func (m *Params) GetExchangeRates() []types1.Exchange {
	if m != nil {
		return m.ExchangeRates
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "auction.Params")
}

func init() { proto.RegisterFile("auction/params.proto", fileDescriptor_e96a95233ccbd0c2) }

var fileDescriptor_e96a95233ccbd0c2 = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0xb1, 0x4e, 0xc2, 0x40,
	0x1c, 0xc6, 0x5b, 0x25, 0x98, 0xd4, 0xe8, 0xd0, 0x40, 0x04, 0x86, 0xc3, 0x38, 0xb1, 0x70, 0x7f,
	0xc1, 0xcd, 0xc5, 0x04, 0x74, 0x30, 0x21, 0x96, 0x54, 0x26, 0x17, 0x72, 0x2d, 0x67, 0xdb, 0x08,
	0xfd, 0x37, 0xbd, 0x43, 0x61, 0xf6, 0x05, 0x7c, 0x04, 0x67, 0x9f, 0x84, 0x91, 0xc5, 0xc4, 0x49,
	0x0d, 0x2c, 0x3e, 0x86, 0xe9, 0xf5, 0x9a, 0xe8, 0x74, 0x77, 0xdf, 0xf7, 0x5d, 0x7e, 0xdf, 0xdd,
	0xdf, 0xaa, 0xb0, 0xb9, 0x2f, 0x23, 0x8c, 0x21, 0x61, 0x29, 0x9b, 0x09, 0x9a, 0xa4, 0x28, 0xd1,
	0xde, 0xd3, 0x6a, 0xa3, 0x12, 0x60, 0x80, 0x4a, 0x83, 0x6c, 0x97, 0xdb, 0x0d, 0xe2, 0xa3, 0x98,
	0xa1, 0x00, 0x8f, 0x09, 0x0e, 0x8f, 0x1d, 0x8f, 0x4b, 0xd6, 0x01, 0x1f, 0xa3, 0x58, 0xfb, 0xf5,
	0x00, 0x31, 0x98, 0x72, 0x50, 0x27, 0x6f, 0x7e, 0x0f, 0x2c, 0x5e, 0x6a, 0xab, 0x9a, 0xc5, 0xc4,
	0x13, 0x4b, 0xfe, 0x01, 0x4f, 0xde, 0x4d, 0xab, 0x3c, 0x54, 0x82, 0xfd, 0x6c, 0x5a, 0x47, 0x1a,
	0x3f, 0x16, 0x92, 0xa5, 0x72, 0x2c, 0xc3, 0x94, 0x8b, 0x10, 0xa7, 0x93, 0x9a, 0x79, 0xbc, 0xdb,
	0xda, 0xef, 0xd6, 0x69, 0xce, 0xa7, 0x19, 0x9f, 0x6a, 0x3e, 0xed, 0x63, 0x14, 0xf7, 0x4e, 0x57,
	0x9f, 0x4d, 0xe3, 0xed, 0xab, 0xd9, 0x0a, 0x22, 0x19, 0xce, 0x3d, 0xea, 0xe3, 0x0c, 0x74, 0xd9,
	0x7c, 0x69, 0x8b, 0xc9, 0x03, 0xc8, 0x65, 0xc2, 0x85, 0xba, 0x20, 0xdc, 0xaa, 0x66, 0xdd, 0x66,
	0xa8, 0x51, 0x41, 0xb2, 0x2f, 0xac, 0x43, 0xbe, 0xf0, 0x43, 0x16, 0x07, 0x7c, 0x9c, 0x32, 0xc9,
	0x45, 0x6d, 0x47, 0xb1, 0x6d, 0x5a, 0x3c, 0x80, 0x5e, 0x69, 0xbf, 0x57, 0xca, 0xa0, 0xee, 0x41,
	0x91, 0x77, 0xb3, 0xf8, 0x79, 0xe9, 0xe7, 0xb5, 0x69, 0xf4, 0x06, 0xab, 0x0d, 0x31, 0xd7, 0x1b,
	0x62, 0x7e, 0x6f, 0x88, 0xf9, 0xb2, 0x25, 0xc6, 0x7a, 0x4b, 0x8c, 0x8f, 0x2d, 0x31, 0xee, 0xba,
	0x7f, 0x1a, 0x3a, 0x97, 0xd7, 0x37, 0xed, 0xa1, 0xeb, 0x8c, 0x9c, 0xbe, 0x33, 0x00, 0x9c, 0x44,
	0x71, 0xdb, 0xc7, 0x94, 0xc3, 0x02, 0x8a, 0xe9, 0xa8, 0xc6, 0x5e, 0x59, 0x7d, 0xd6, 0xd9, 0x6f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x31, 0x7a, 0x28, 0x07, 0xb5, 0x01, 0x00, 0x00,
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
	if len(m.ExchangeRates) > 0 {
		for iNdEx := len(m.ExchangeRates) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ExchangeRates[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.AuctionStartThreshold) > 0 {
		for iNdEx := len(m.AuctionStartThreshold) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AuctionStartThreshold[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.AuctionStartThreshold) > 0 {
		for _, e := range m.AuctionStartThreshold {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.ExchangeRates) > 0 {
		for _, e := range m.ExchangeRates {
			l = e.Size()
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
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionStartThreshold", wireType)
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
			m.AuctionStartThreshold = append(m.AuctionStartThreshold, types.Coin{})
			if err := m.AuctionStartThreshold[len(m.AuctionStartThreshold)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExchangeRates", wireType)
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
			m.ExchangeRates = append(m.ExchangeRates, types1.Exchange{})
			if err := m.ExchangeRates[len(m.ExchangeRates)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
