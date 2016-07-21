// Code generated by protoc-gen-go.
// source: serverModel.proto
// DO NOT EDIT!

package bbproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Ignoring public import of ProtoHeader from base.proto

// Ignoring public import of TerminalInfo from base.proto

// Ignoring public import of EProtoId from base.proto

// Ignoring public import of ZjhRoom from zhajinhua.proto

// Ignoring public import of ZjhBet from zhajinhua.proto

// Ignoring public import of ZjhLottery from zhajinhua.proto

// Ignoring public import of BroadcastBet from zhajinhua.proto

// Ignoring public import of ZjhPai from zhajinhua.proto

// Ignoring public import of Pai from zhajinhua.proto

// Ignoring public import of ZjhQueryNoSeatUser from zhajinhua.proto

// Ignoring public import of ZjhReqSeat from zhajinhua.proto

// Ignoring public import of ZjhMsg from zhajinhua.proto

// Ignoring public import of ZjhBroadcastBeginBet from zhajinhua.proto

// Ignoring public import of EPaiType from zhajinhua.proto

// 押注的记录,针对每轮扎金花,只有一次押注记录
type TBetRecord struct {
	UserId           *uint32 `protobuf:"varint,1,opt,name=userId" json:"userId,omitempty"`
	Betzone          []int32 `protobuf:"varint,2,rep,name=betzone" json:"betzone,omitempty"`
	ZjhRoundNumber   *string `protobuf:"bytes,3,opt,name=ZjhRoundNumber" json:"ZjhRoundNumber,omitempty"`
	WinAmount        *int32  `protobuf:"varint,4,opt,name=winAmount" json:"winAmount,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *TBetRecord) Reset()                    { *m = TBetRecord{} }
func (m *TBetRecord) String() string            { return proto.CompactTextString(m) }
func (*TBetRecord) ProtoMessage()               {}
func (*TBetRecord) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{0} }

func (m *TBetRecord) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *TBetRecord) GetBetzone() []int32 {
	if m != nil {
		return m.Betzone
	}
	return nil
}

func (m *TBetRecord) GetZjhRoundNumber() string {
	if m != nil && m.ZjhRoundNumber != nil {
		return *m.ZjhRoundNumber
	}
	return ""
}

func (m *TBetRecord) GetWinAmount() int32 {
	if m != nil && m.WinAmount != nil {
		return *m.WinAmount
	}
	return 0
}

// 每一轮扎金花的数据
type TZjhRound struct {
	Id               *uint32   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	BeginTime        *int64    `protobuf:"varint,2,opt,name=beginTime" json:"beginTime,omitempty"`
	BetEndTime       *int64    `protobuf:"varint,3,opt,name=betEndTime" json:"betEndTime,omitempty"`
	LotteryTime      *int64    `protobuf:"varint,4,opt,name=lotteryTime" json:"lotteryTime,omitempty"`
	EndTime          *int64    `protobuf:"varint,5,opt,name=endTime" json:"endTime,omitempty"`
	ZoneAmount       []int32   `protobuf:"varint,6,rep,name=zoneAmount" json:"zoneAmount,omitempty"`
	ZoneWinAmount    []int32   `protobuf:"varint,7,rep,name=zoneWinAmount" json:"zoneWinAmount,omitempty"`
	BankerUserId     *uint32   `protobuf:"varint,8,opt,name=BankerUserId" json:"BankerUserId,omitempty"`
	ZjhPaiList       []*ZjhPai `protobuf:"bytes,9,rep,name=ZjhPaiList" json:"ZjhPaiList,omitempty"`
	Number           *string   `protobuf:"bytes,10,opt,name=Number" json:"Number,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *TZjhRound) Reset()                    { *m = TZjhRound{} }
func (m *TZjhRound) String() string            { return proto.CompactTextString(m) }
func (*TZjhRound) ProtoMessage()               {}
func (*TZjhRound) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{1} }

func (m *TZjhRound) GetId() uint32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *TZjhRound) GetBeginTime() int64 {
	if m != nil && m.BeginTime != nil {
		return *m.BeginTime
	}
	return 0
}

func (m *TZjhRound) GetBetEndTime() int64 {
	if m != nil && m.BetEndTime != nil {
		return *m.BetEndTime
	}
	return 0
}

func (m *TZjhRound) GetLotteryTime() int64 {
	if m != nil && m.LotteryTime != nil {
		return *m.LotteryTime
	}
	return 0
}

func (m *TZjhRound) GetEndTime() int64 {
	if m != nil && m.EndTime != nil {
		return *m.EndTime
	}
	return 0
}

func (m *TZjhRound) GetZoneAmount() []int32 {
	if m != nil {
		return m.ZoneAmount
	}
	return nil
}

func (m *TZjhRound) GetZoneWinAmount() []int32 {
	if m != nil {
		return m.ZoneWinAmount
	}
	return nil
}

func (m *TZjhRound) GetBankerUserId() uint32 {
	if m != nil && m.BankerUserId != nil {
		return *m.BankerUserId
	}
	return 0
}

func (m *TZjhRound) GetZjhPaiList() []*ZjhPai {
	if m != nil {
		return m.ZjhPaiList
	}
	return nil
}

func (m *TZjhRound) GetNumber() string {
	if m != nil && m.Number != nil {
		return *m.Number
	}
	return ""
}

func init() {
	proto.RegisterType((*TBetRecord)(nil), "bbproto.TBetRecord")
	proto.RegisterType((*TZjhRound)(nil), "bbproto.TZjhRound")
}

var fileDescriptor10 = []byte{
	// 279 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x3c, 0x90, 0xcf, 0x4e, 0x83, 0x40,
	0x10, 0x87, 0xa5, 0x94, 0x22, 0x53, 0x5b, 0xd2, 0xf5, 0x4f, 0x48, 0x4f, 0x4d, 0xbd, 0x70, 0xe2,
	0xe0, 0x1b, 0xd8, 0xc4, 0x83, 0x89, 0x1a, 0xd2, 0x60, 0x34, 0xde, 0x58, 0x99, 0xc8, 0xd6, 0xb2,
	0x6b, 0x96, 0x45, 0x63, 0x1f, 0xd5, 0xa7, 0x71, 0x58, 0xa0, 0xb7, 0xcd, 0x37, 0xbf, 0x9d, 0xf9,
	0x66, 0x60, 0x51, 0xa3, 0xfe, 0x46, 0xfd, 0xa8, 0x0a, 0xdc, 0x27, 0x5f, 0x5a, 0x19, 0xc5, 0x7c,
	0xce, 0xed, 0x63, 0x09, 0x3c, 0xaf, 0xb1, 0x83, 0xcb, 0xf0, 0x50, 0xe6, 0x3b, 0x21, 0xcb, 0x26,
	0xef, 0xc0, 0xfa, 0x15, 0x20, 0xdb, 0xa0, 0xd9, 0xe2, 0xbb, 0xd2, 0x05, 0x9b, 0xc3, 0xa4, 0xa1,
	0x4e, 0xf7, 0x45, 0xe4, 0xac, 0x9c, 0x78, 0xc6, 0x42, 0xf0, 0x39, 0x9a, 0x83, 0x92, 0x18, 0x8d,
	0x56, 0x6e, 0xec, 0xb1, 0x2b, 0x98, 0xbf, 0xed, 0xca, 0xad, 0x6a, 0x64, 0xf1, 0xd4, 0x54, 0x1c,
	0x75, 0xe4, 0x52, 0x30, 0x60, 0x0b, 0x08, 0x7e, 0x84, 0xbc, 0xad, 0xa8, 0x60, 0xa2, 0x31, 0x21,
	0x6f, 0xfd, 0xe7, 0x40, 0x90, 0x0d, 0x61, 0x06, 0x30, 0x12, 0x43, 0x57, 0x0a, 0x73, 0xfc, 0x10,
	0x32, 0x13, 0x55, 0xdb, 0xd7, 0x89, 0x5d, 0x46, 0x75, 0x1a, 0x74, 0x27, 0x0b, 0xcb, 0x5c, 0xcb,
	0xce, 0x61, 0xba, 0x57, 0xc6, 0xa0, 0xfe, 0xb5, 0x70, 0x6c, 0x21, 0x19, 0x61, 0x9f, 0xf2, 0x86,
	0x9f, 0xad, 0x5f, 0x3f, 0x7a, 0x62, 0x2d, 0x2f, 0x61, 0xd6, 0xb2, 0x97, 0xa3, 0x91, 0x6f, 0xf1,
	0x05, 0x9c, 0x6d, 0x72, 0xf9, 0x89, 0xfa, 0xb9, 0xdb, 0xf1, 0xd4, 0xda, 0x5c, 0x03, 0x90, 0x65,
	0x9a, 0x8b, 0x07, 0x51, 0x9b, 0x28, 0xa0, 0xe4, 0xf4, 0x26, 0x4c, 0xfa, 0xe3, 0x25, 0x5d, 0xa9,
	0x3d, 0x4c, 0xbf, 0x2f, 0xb4, 0xfb, 0xa6, 0x27, 0xa9, 0xf3, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x6e,
	0x2e, 0x02, 0x98, 0x74, 0x01, 0x00, 0x00,
}
