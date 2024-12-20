// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: titan/farming/staking_info.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// StakingInfo defines the staking info of a staker for a token.
type StakingInfo struct {
	Token  string                `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Staker string                `protobuf:"bytes,2,opt,name=staker,proto3" json:"staker,omitempty"`
	Amount cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=amount,proto3,customtype=cosmossdk.io/math.Int" json:"amount"`
}

func (m *StakingInfo) Reset()         { *m = StakingInfo{} }
func (m *StakingInfo) String() string { return proto.CompactTextString(m) }
func (*StakingInfo) ProtoMessage()    {}
func (*StakingInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_dc2cf0ffff1c61f5, []int{0}
}
func (m *StakingInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StakingInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StakingInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StakingInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StakingInfo.Merge(m, src)
}
func (m *StakingInfo) XXX_Size() int {
	return m.Size()
}
func (m *StakingInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_StakingInfo.DiscardUnknown(m)
}

var xxx_messageInfo_StakingInfo proto.InternalMessageInfo

func (m *StakingInfo) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *StakingInfo) GetStaker() string {
	if m != nil {
		return m.Staker
	}
	return ""
}

func init() {
	proto.RegisterType((*StakingInfo)(nil), "titan.farming.StakingInfo")
}

func init() { proto.RegisterFile("titan/farming/staking_info.proto", fileDescriptor_dc2cf0ffff1c61f5) }

var fileDescriptor_dc2cf0ffff1c61f5 = []byte{
	// 290 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x28, 0xc9, 0x2c, 0x49,
	0xcc, 0xd3, 0x4f, 0x4b, 0x2c, 0xca, 0xcd, 0xcc, 0x4b, 0xd7, 0x2f, 0x2e, 0x49, 0xcc, 0xce, 0xcc,
	0x4b, 0x8f, 0xcf, 0xcc, 0x4b, 0xcb, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x05, 0xab,
	0xd0, 0x83, 0xaa, 0x90, 0x12, 0x4c, 0xcc, 0xcd, 0xcc, 0xcb, 0xd7, 0x07, 0x93, 0x10, 0x15, 0x52,
	0x92, 0xc9, 0xf9, 0xc5, 0xb9, 0xf9, 0xc5, 0xf1, 0x60, 0x9e, 0x3e, 0x84, 0x03, 0x95, 0x12, 0x49,
	0xcf, 0x4f, 0xcf, 0x87, 0x88, 0x83, 0x58, 0x10, 0x51, 0xa5, 0xf9, 0x8c, 0x5c, 0xdc, 0xc1, 0x10,
	0x9b, 0x3c, 0xf3, 0xd2, 0xf2, 0x85, 0x44, 0xb8, 0x58, 0x4b, 0xf2, 0xb3, 0x53, 0xf3, 0x24, 0x18,
	0x15, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x21, 0x03, 0x2e, 0x36, 0x90, 0x73, 0x52, 0x8b, 0x24,
	0x98, 0x40, 0xc2, 0x4e, 0x12, 0x97, 0xb6, 0xe8, 0x8a, 0x40, 0x4d, 0x77, 0x4c, 0x49, 0x29, 0x4a,
	0x2d, 0x2e, 0x0e, 0x2e, 0x29, 0xca, 0xcc, 0x4b, 0x0f, 0x82, 0xaa, 0x13, 0xf2, 0xe0, 0x62, 0x4b,
	0xcc, 0xcd, 0x2f, 0xcd, 0x2b, 0x91, 0x60, 0x06, 0xeb, 0x30, 0x38, 0x71, 0x4f, 0x9e, 0xe1, 0xd6,
	0x3d, 0x79, 0x51, 0x88, 0xae, 0xe2, 0x94, 0x6c, 0xbd, 0xcc, 0x7c, 0xfd, 0xdc, 0xc4, 0x92, 0x0c,
	0x3d, 0xcf, 0xbc, 0x92, 0x4b, 0x5b, 0x74, 0xb9, 0xa0, 0xc6, 0x79, 0xe6, 0x95, 0xac, 0x78, 0xbe,
	0x41, 0x8b, 0x31, 0x08, 0xaa, 0xdf, 0xc9, 0xf9, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18,
	0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5,
	0x18, 0xa2, 0x34, 0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xc1, 0x21,
	0x53, 0x92, 0x5d, 0x01, 0x61, 0xe8, 0x57, 0xc0, 0x83, 0xb1, 0xa4, 0xb2, 0x20, 0xb5, 0x38, 0x89,
	0x0d, 0xec, 0x5b, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x64, 0x95, 0x66, 0x0f, 0x64, 0x01,
	0x00, 0x00,
}

func (m *StakingInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StakingInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StakingInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintStakingInfo(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Staker) > 0 {
		i -= len(m.Staker)
		copy(dAtA[i:], m.Staker)
		i = encodeVarintStakingInfo(dAtA, i, uint64(len(m.Staker)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintStakingInfo(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintStakingInfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovStakingInfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StakingInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovStakingInfo(uint64(l))
	}
	l = len(m.Staker)
	if l > 0 {
		n += 1 + l + sovStakingInfo(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovStakingInfo(uint64(l))
	return n
}

func sovStakingInfo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStakingInfo(x uint64) (n int) {
	return sovStakingInfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StakingInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStakingInfo
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
			return fmt.Errorf("proto: StakingInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StakingInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStakingInfo
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
				return ErrInvalidLengthStakingInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStakingInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Staker", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStakingInfo
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
				return ErrInvalidLengthStakingInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStakingInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Staker = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStakingInfo
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
				return ErrInvalidLengthStakingInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStakingInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStakingInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStakingInfo
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
func skipStakingInfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStakingInfo
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
					return 0, ErrIntOverflowStakingInfo
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
					return 0, ErrIntOverflowStakingInfo
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
				return 0, ErrInvalidLengthStakingInfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStakingInfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStakingInfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStakingInfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStakingInfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStakingInfo = fmt.Errorf("proto: unexpected end of group")
)