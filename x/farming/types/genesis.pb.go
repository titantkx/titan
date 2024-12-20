// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: titan/farming/genesis.proto

package types

import (
	fmt "fmt"
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

// GenesisState defines the farming module's genesis state.
type GenesisState struct {
	Params           Params            `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	FarmList         []Farm            `protobuf:"bytes,2,rep,name=farm_list,json=farmList,proto3" json:"farm_list"`
	StakingInfoList  []StakingInfo     `protobuf:"bytes,3,rep,name=staking_info_list,json=stakingInfoList,proto3" json:"staking_info_list"`
	DistributionInfo *DistributionInfo `protobuf:"bytes,4,opt,name=distribution_info,json=distributionInfo,proto3" json:"distribution_info,omitempty"`
	RewardList       []Reward          `protobuf:"bytes,5,rep,name=reward_list,json=rewardList,proto3" json:"reward_list"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_3376d6f5f93efbe0, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetFarmList() []Farm {
	if m != nil {
		return m.FarmList
	}
	return nil
}

func (m *GenesisState) GetStakingInfoList() []StakingInfo {
	if m != nil {
		return m.StakingInfoList
	}
	return nil
}

func (m *GenesisState) GetDistributionInfo() *DistributionInfo {
	if m != nil {
		return m.DistributionInfo
	}
	return nil
}

func (m *GenesisState) GetRewardList() []Reward {
	if m != nil {
		return m.RewardList
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "titan.farming.GenesisState")
}

func init() { proto.RegisterFile("titan/farming/genesis.proto", fileDescriptor_3376d6f5f93efbe0) }

var fileDescriptor_3376d6f5f93efbe0 = []byte{
	// 345 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xcf, 0x4a, 0xf3, 0x40,
	0x14, 0xc5, 0x93, 0xb6, 0x5f, 0xf9, 0x9c, 0x2a, 0x6a, 0x54, 0x08, 0x11, 0xd2, 0x22, 0x08, 0x75,
	0x93, 0x40, 0x0b, 0xae, 0x5c, 0x55, 0x51, 0x84, 0x2e, 0xa4, 0xdd, 0xb9, 0x29, 0x53, 0x9b, 0x8e,
	0x43, 0xcd, 0x4c, 0x99, 0xb9, 0xc5, 0xfa, 0x16, 0x3e, 0x56, 0x97, 0xdd, 0x08, 0xae, 0x44, 0xda,
	0x17, 0x91, 0xf9, 0x53, 0x4d, 0xa2, 0xab, 0x0c, 0xf3, 0x3b, 0xe7, 0xe4, 0xdc, 0x3b, 0xe8, 0x18,
	0x28, 0x60, 0x16, 0x8f, 0xb1, 0x48, 0x29, 0x23, 0x31, 0x49, 0x58, 0x22, 0xa9, 0x8c, 0xa6, 0x82,
	0x03, 0xf7, 0x76, 0x34, 0x8c, 0x2c, 0x0c, 0x0e, 0x09, 0x27, 0x5c, 0x93, 0x58, 0x9d, 0x8c, 0x28,
	0x38, 0xcd, 0x27, 0x8c, 0xa8, 0x04, 0x41, 0x87, 0x33, 0xa0, 0x9c, 0x0d, 0x28, 0x1b, 0x6f, 0x64,
	0x7e, 0x5e, 0xa6, 0xbe, 0x96, 0x04, 0x79, 0x32, 0xc5, 0x02, 0xa7, 0xf2, 0x6f, 0x26, 0x92, 0x67,
	0x2c, 0x46, 0x96, 0x35, 0xf2, 0x4c, 0x02, 0x9e, 0x50, 0x46, 0x32, 0xff, 0x3c, 0x79, 0x2b, 0xa1,
	0xed, 0x1b, 0x33, 0x51, 0x1f, 0x30, 0x24, 0x5e, 0x1b, 0x55, 0x4d, 0xbc, 0xef, 0x36, 0xdc, 0x66,
	0xad, 0x75, 0x14, 0xe5, 0x26, 0x8c, 0xee, 0x34, 0xec, 0x54, 0x16, 0x1f, 0x75, 0xa7, 0x67, 0xa5,
	0xde, 0x39, 0xda, 0x52, 0x7c, 0xf0, 0x44, 0x25, 0xf8, 0xa5, 0x46, 0xb9, 0x59, 0x6b, 0x1d, 0x14,
	0x7c, 0xd7, 0x58, 0xa4, 0xd6, 0xf5, 0x5f, 0xdd, 0x75, 0xa9, 0x04, 0xaf, 0x8b, 0xf6, 0xb3, 0x9d,
	0x8c, 0xbf, 0xac, 0xfd, 0x41, 0xc1, 0xdf, 0x37, 0xba, 0x5b, 0x36, 0xe6, 0x36, 0x66, 0x57, 0xfe,
	0x5c, 0x6d, 0xd2, 0x7e, 0xad, 0xd6, 0xaf, 0xe8, 0x29, 0xea, 0x85, 0xb4, 0xab, 0x8c, 0x4e, 0xf9,
	0x7b, 0x7b, 0xa3, 0xc2, 0x8d, 0x77, 0x81, 0x6a, 0x66, 0x97, 0xa6, 0xd5, 0x3f, 0xdd, 0xaa, 0xb8,
	0x8d, 0x9e, 0x56, 0xd8, 0x42, 0xc8, 0xe8, 0x55, 0x97, 0xce, 0xe5, 0x62, 0x15, 0xba, 0xcb, 0x55,
	0xe8, 0x7e, 0xae, 0x42, 0xf7, 0x75, 0x1d, 0x3a, 0xcb, 0x75, 0xe8, 0xbc, 0xaf, 0x43, 0xe7, 0xfe,
	0x8c, 0x50, 0x78, 0x9c, 0x0d, 0xa3, 0x07, 0x9e, 0xc6, 0x3a, 0x0c, 0x26, 0x73, 0x73, 0x88, 0xe7,
	0xdf, 0x2f, 0x05, 0x2f, 0xd3, 0x44, 0x0e, 0xab, 0xfa, 0x8d, 0xda, 0x5f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xcb, 0x71, 0xe6, 0x8c, 0x82, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RewardList) > 0 {
		for iNdEx := len(m.RewardList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RewardList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if m.DistributionInfo != nil {
		{
			size, err := m.DistributionInfo.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.StakingInfoList) > 0 {
		for iNdEx := len(m.StakingInfoList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.StakingInfoList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.FarmList) > 0 {
		for iNdEx := len(m.FarmList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.FarmList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.FarmList) > 0 {
		for _, e := range m.FarmList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.StakingInfoList) > 0 {
		for _, e := range m.StakingInfoList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.DistributionInfo != nil {
		l = m.DistributionInfo.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.RewardList) > 0 {
		for _, e := range m.RewardList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FarmList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FarmList = append(m.FarmList, Farm{})
			if err := m.FarmList[len(m.FarmList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakingInfoList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StakingInfoList = append(m.StakingInfoList, StakingInfo{})
			if err := m.StakingInfoList[len(m.StakingInfoList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DistributionInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DistributionInfo == nil {
				m.DistributionInfo = &DistributionInfo{}
			}
			if err := m.DistributionInfo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardList = append(m.RewardList, Reward{})
			if err := m.RewardList[len(m.RewardList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
