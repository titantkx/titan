// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: titan/farming/farm.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Farm defines the farming rewards for a token.
type Farm struct {
	Token   string        `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Rewards []*FarmReward `protobuf:"bytes,2,rep,name=rewards,proto3" json:"rewards,omitempty"`
}

func (m *Farm) Reset()         { *m = Farm{} }
func (m *Farm) String() string { return proto.CompactTextString(m) }
func (*Farm) ProtoMessage()    {}
func (*Farm) Descriptor() ([]byte, []int) {
	return fileDescriptor_a56e8f061d8bd3d8, []int{0}
}
func (m *Farm) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Farm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Farm.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Farm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Farm.Merge(m, src)
}
func (m *Farm) XXX_Size() int {
	return m.Size()
}
func (m *Farm) XXX_DiscardUnknown() {
	xxx_messageInfo_Farm.DiscardUnknown(m)
}

var xxx_messageInfo_Farm proto.InternalMessageInfo

func (m *Farm) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Farm) GetRewards() []*FarmReward {
	if m != nil {
		return m.Rewards
	}
	return nil
}

// FarmReward defines the farming rewards.
type FarmReward struct {
	Sender    string                                   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Amount    github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	EndTime   time.Time                                `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3,stdtime" json:"end_time"`
	StartTime time.Time                                `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3,stdtime" json:"start_time"`
}

func (m *FarmReward) Reset()         { *m = FarmReward{} }
func (m *FarmReward) String() string { return proto.CompactTextString(m) }
func (*FarmReward) ProtoMessage()    {}
func (*FarmReward) Descriptor() ([]byte, []int) {
	return fileDescriptor_a56e8f061d8bd3d8, []int{1}
}
func (m *FarmReward) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FarmReward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FarmReward.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FarmReward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FarmReward.Merge(m, src)
}
func (m *FarmReward) XXX_Size() int {
	return m.Size()
}
func (m *FarmReward) XXX_DiscardUnknown() {
	xxx_messageInfo_FarmReward.DiscardUnknown(m)
}

var xxx_messageInfo_FarmReward proto.InternalMessageInfo

func (m *FarmReward) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *FarmReward) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *FarmReward) GetEndTime() time.Time {
	if m != nil {
		return m.EndTime
	}
	return time.Time{}
}

func (m *FarmReward) GetStartTime() time.Time {
	if m != nil {
		return m.StartTime
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*Farm)(nil), "titan.farming.Farm")
	proto.RegisterType((*FarmReward)(nil), "titan.farming.FarmReward")
}

func init() { proto.RegisterFile("titan/farming/farm.proto", fileDescriptor_a56e8f061d8bd3d8) }

var fileDescriptor_a56e8f061d8bd3d8 = []byte{
	// 425 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xbf, 0x6e, 0xdb, 0x30,
	0x10, 0xc6, 0xc5, 0x24, 0x75, 0x12, 0xa6, 0x1d, 0x2a, 0x78, 0x50, 0x3c, 0x48, 0x46, 0x26, 0x35,
	0x40, 0xc8, 0x26, 0x79, 0x80, 0xa2, 0x32, 0x90, 0xbd, 0x6a, 0xa7, 0x2e, 0x06, 0x25, 0x31, 0x2c,
	0xe1, 0x90, 0x34, 0x48, 0xba, 0x8d, 0xdf, 0x22, 0x73, 0x9f, 0xa0, 0xe8, 0x94, 0xa1, 0x7b, 0xd7,
	0x8c, 0x41, 0xa7, 0x4e, 0x4d, 0x61, 0x0f, 0x7e, 0x8d, 0x82, 0x7f, 0xec, 0xb6, 0x63, 0x16, 0x91,
	0x77, 0xdf, 0x77, 0x77, 0xba, 0x9f, 0x04, 0x33, 0xcb, 0x2d, 0x91, 0xf8, 0x92, 0x68, 0xc1, 0x25,
	0xf3, 0x27, 0x9a, 0x6a, 0x65, 0x55, 0xfa, 0xcc, 0x2b, 0x28, 0x2a, 0x83, 0xe7, 0x44, 0x70, 0xa9,
	0xb0, 0x7f, 0x06, 0xc7, 0x20, 0x6f, 0x95, 0x11, 0xca, 0xe0, 0x86, 0x18, 0x8a, 0x3f, 0x9e, 0x36,
	0xd4, 0x92, 0x53, 0xdc, 0x2a, 0x2e, 0xa3, 0x7e, 0x18, 0xf4, 0xb1, 0x8f, 0x70, 0x08, 0xa2, 0xd4,
	0x67, 0x8a, 0xa9, 0x90, 0x77, 0xb7, 0x98, 0x2d, 0x98, 0x52, 0xec, 0x8a, 0x62, 0x1f, 0x35, 0xb3,
	0x4b, 0x6c, 0xb9, 0xa0, 0xc6, 0x12, 0x31, 0x0d, 0x86, 0xa3, 0x37, 0x70, 0xe7, 0x82, 0x68, 0x91,
	0xf6, 0xe1, 0x13, 0xab, 0x26, 0x54, 0x66, 0x60, 0x08, 0xca, 0xfd, 0x3a, 0x04, 0xe9, 0x39, 0xdc,
	0xd5, 0xf4, 0x13, 0xd1, 0x9d, 0xc9, 0xb6, 0x86, 0xdb, 0xe5, 0xc1, 0xd9, 0x21, 0xfa, 0x6f, 0x07,
	0xe4, 0x6a, 0x6b, 0xef, 0xa8, 0xd7, 0xce, 0xa3, 0xef, 0x5b, 0x10, 0xfe, 0xcd, 0xa7, 0x2f, 0x61,
	0xcf, 0x50, 0xd9, 0x51, 0x1d, 0x5a, 0x57, 0xd9, 0x8f, 0x6f, 0x27, 0xfd, 0xf8, 0xea, 0xaf, 0xbb,
	0x4e, 0x53, 0x63, 0xde, 0x5a, 0xcd, 0x25, 0xab, 0xa3, 0x2f, 0x9d, 0xc3, 0x1e, 0x11, 0x6a, 0x26,
	0xed, 0x66, 0x68, 0xb4, 0x3b, 0x2c, 0x28, 0x62, 0x41, 0x23, 0xc5, 0x65, 0x75, 0x71, 0xf7, 0xab,
	0x48, 0xbe, 0x3e, 0x14, 0x25, 0xe3, 0xf6, 0xc3, 0xac, 0x41, 0xad, 0x12, 0x11, 0x4b, 0x3c, 0x4e,
	0x4c, 0x37, 0xc1, 0x76, 0x3e, 0xa5, 0xc6, 0x17, 0x98, 0xcf, 0xab, 0xdb, 0xe3, 0xa7, 0x57, 0x94,
	0x91, 0x76, 0x3e, 0x76, 0x60, 0xcd, 0x97, 0xd5, 0xed, 0x31, 0xa8, 0xe3, 0xc0, 0xf4, 0x15, 0xdc,
	0xa3, 0xb2, 0x1b, 0x3b, 0x4a, 0xd9, 0xf6, 0x10, 0x94, 0x07, 0x67, 0x03, 0x14, 0x10, 0xa2, 0x35,
	0x42, 0xf4, 0x6e, 0x8d, 0xb0, 0xda, 0x73, 0xd3, 0x6f, 0x1e, 0x0a, 0x50, 0xef, 0x52, 0xd9, 0xb9,
	0x7c, 0x3a, 0x82, 0xd0, 0x58, 0xa2, 0x6d, 0x68, 0xb1, 0xf3, 0x88, 0x16, 0xfb, 0xbe, 0xce, 0x29,
	0xd5, 0xe8, 0x6e, 0x91, 0x83, 0xfb, 0x45, 0x0e, 0x7e, 0x2f, 0x72, 0x70, 0xb3, 0xcc, 0x93, 0xfb,
	0x65, 0x9e, 0xfc, 0x5c, 0xe6, 0xc9, 0xfb, 0x17, 0xff, 0xec, 0xe9, 0xbf, 0x84, 0x9d, 0x5c, 0x87,
	0x0b, 0xbe, 0xde, 0xfc, 0x72, 0x7e, 0xdd, 0xa6, 0xe7, 0xa7, 0x9d, 0xff, 0x09, 0x00, 0x00, 0xff,
	0xff, 0xac, 0x5a, 0x8d, 0x56, 0x90, 0x02, 0x00, 0x00,
}

func (m *Farm) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Farm) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Farm) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Rewards) > 0 {
		for iNdEx := len(m.Rewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Rewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintFarm(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintFarm(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *FarmReward) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FarmReward) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FarmReward) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.StartTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.StartTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintFarm(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.EndTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.EndTime):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintFarm(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x1a
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintFarm(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintFarm(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintFarm(dAtA []byte, offset int, v uint64) int {
	offset -= sovFarm(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Farm) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovFarm(uint64(l))
	}
	if len(m.Rewards) > 0 {
		for _, e := range m.Rewards {
			l = e.Size()
			n += 1 + l + sovFarm(uint64(l))
		}
	}
	return n
}

func (m *FarmReward) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovFarm(uint64(l))
	}
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovFarm(uint64(l))
		}
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.EndTime)
	n += 1 + l + sovFarm(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.StartTime)
	n += 1 + l + sovFarm(uint64(l))
	return n
}

func sovFarm(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFarm(x uint64) (n int) {
	return sovFarm(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Farm) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFarm
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
			return fmt.Errorf("proto: Farm: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Farm: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFarm
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
				return ErrInvalidLengthFarm
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFarm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFarm
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
				return ErrInvalidLengthFarm
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFarm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rewards = append(m.Rewards, &FarmReward{})
			if err := m.Rewards[len(m.Rewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFarm(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFarm
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
func (m *FarmReward) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFarm
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
			return fmt.Errorf("proto: FarmReward: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FarmReward: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFarm
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
				return ErrInvalidLengthFarm
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFarm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFarm
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
				return ErrInvalidLengthFarm
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFarm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFarm
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
				return ErrInvalidLengthFarm
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFarm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.EndTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFarm
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
				return ErrInvalidLengthFarm
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFarm
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.StartTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFarm(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFarm
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
func skipFarm(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFarm
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
					return 0, ErrIntOverflowFarm
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
					return 0, ErrIntOverflowFarm
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
				return 0, ErrInvalidLengthFarm
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFarm
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFarm
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFarm        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFarm          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFarm = fmt.Errorf("proto: unexpected end of group")
)