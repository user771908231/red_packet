// Code generated by protoc-gen-go.
// source: zhajinhua.proto
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

// Ignoring public import of EUnitType from base.proto

// Ignoring public import of EUnitRace from base.proto

// Ignoring public import of EUnitProtoId from base.proto

// Ignoring public import of User from user.proto

// 进入扎金花的房间:押注中(剩余time）、开奖中（剩余time）、jackpot奖池金额、balance、庄家信息、在座玩家
type ZjhRoom struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	ReqType          *int32       `protobuf:"varint,2,opt,name=reqType" json:"reqType,omitempty"`
	Jackpot          *int64       `protobuf:"varint,3,opt,name=jackpot" json:"jackpot,omitempty"`
	Banker           *User        `protobuf:"bytes,4,opt,name=banker" json:"banker,omitempty"`
	Me               *User        `protobuf:"bytes,5,opt,name=me" json:"me,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ZjhRoom) Reset()                    { *m = ZjhRoom{} }
func (m *ZjhRoom) String() string            { return proto.CompactTextString(m) }
func (*ZjhRoom) ProtoMessage()               {}
func (*ZjhRoom) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{0} }

func (m *ZjhRoom) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ZjhRoom) GetReqType() int32 {
	if m != nil && m.ReqType != nil {
		return *m.ReqType
	}
	return 0
}

func (m *ZjhRoom) GetJackpot() int64 {
	if m != nil && m.Jackpot != nil {
		return *m.Jackpot
	}
	return 0
}

func (m *ZjhRoom) GetBanker() *User {
	if m != nil {
		return m.Banker
	}
	return nil
}

func (m *ZjhRoom) GetMe() *User {
	if m != nil {
		return m.Me
	}
	return nil
}

// 押注：押注区betzone=1,2,3,4   押注金额betAmount （若当前时刻已停止押注，返回失败）
type ZjhBet struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Betzone          []int32      `protobuf:"varint,2,rep,name=betzone" json:"betzone,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ZjhBet) Reset()                    { *m = ZjhBet{} }
func (m *ZjhBet) String() string            { return proto.CompactTextString(m) }
func (*ZjhBet) ProtoMessage()               {}
func (*ZjhBet) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{1} }

func (m *ZjhBet) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ZjhBet) GetBetzone() []int32 {
	if m != nil {
		return m.Betzone
	}
	return nil
}

// 开奖,广播
type ZjhLottery struct {
	// 牌面
	Header *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Zjhpai []*ZjhPai    `protobuf:"bytes,2,rep,name=zjhpai" json:"zjhpai,omitempty"`
	// 自己的输赢
	Balance          *int32 `protobuf:"varint,7,opt,name=balance" json:"balance,omitempty"`
	WinAmount        *int32 `protobuf:"varint,8,opt,name=winAmount" json:"winAmount,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ZjhLottery) Reset()                    { *m = ZjhLottery{} }
func (m *ZjhLottery) String() string            { return proto.CompactTextString(m) }
func (*ZjhLottery) ProtoMessage()               {}
func (*ZjhLottery) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{2} }

func (m *ZjhLottery) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ZjhLottery) GetZjhpai() []*ZjhPai {
	if m != nil {
		return m.Zjhpai
	}
	return nil
}

func (m *ZjhLottery) GetBalance() int32 {
	if m != nil && m.Balance != nil {
		return *m.Balance
	}
	return 0
}

func (m *ZjhLottery) GetWinAmount() int32 {
	if m != nil && m.WinAmount != nil {
		return *m.WinAmount
	}
	return 0
}

// 扎金花牌的结构
type ZjhPai struct {
	PaiType          *int32 `protobuf:"varint,1,opt,name=PaiType" json:"PaiType,omitempty"`
	LotteryType      *int32 `protobuf:"varint,2,opt,name=lotteryType" json:"lotteryType,omitempty"`
	Result           *bool  `protobuf:"varint,3,opt,name=result" json:"result,omitempty"`
	Pai              []*Pai `protobuf:"bytes,4,rep,name=pai" json:"pai,omitempty"`
	WinAmount        *int32 `protobuf:"varint,7,opt,name=winAmount" json:"winAmount,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ZjhPai) Reset()                    { *m = ZjhPai{} }
func (m *ZjhPai) String() string            { return proto.CompactTextString(m) }
func (*ZjhPai) ProtoMessage()               {}
func (*ZjhPai) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{3} }

func (m *ZjhPai) GetPaiType() int32 {
	if m != nil && m.PaiType != nil {
		return *m.PaiType
	}
	return 0
}

func (m *ZjhPai) GetLotteryType() int32 {
	if m != nil && m.LotteryType != nil {
		return *m.LotteryType
	}
	return 0
}

func (m *ZjhPai) GetResult() bool {
	if m != nil && m.Result != nil {
		return *m.Result
	}
	return false
}

func (m *ZjhPai) GetPai() []*Pai {
	if m != nil {
		return m.Pai
	}
	return nil
}

func (m *ZjhPai) GetWinAmount() int32 {
	if m != nil && m.WinAmount != nil {
		return *m.WinAmount
	}
	return 0
}

// 单张牌的结构
type Pai struct {
	MapKey           *int32  `protobuf:"varint,1,opt,name=mapKey" json:"mapKey,omitempty"`
	Mapdes           *string `protobuf:"bytes,2,opt,name=mapdes" json:"mapdes,omitempty"`
	Value            *int32  `protobuf:"varint,3,opt,name=value" json:"value,omitempty"`
	Flower           *string `protobuf:"bytes,4,opt,name=flower" json:"flower,omitempty"`
	Name             *string `protobuf:"bytes,5,opt,name=name" json:"name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Pai) Reset()                    { *m = Pai{} }
func (m *Pai) String() string            { return proto.CompactTextString(m) }
func (*Pai) ProtoMessage()               {}
func (*Pai) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{4} }

func (m *Pai) GetMapKey() int32 {
	if m != nil && m.MapKey != nil {
		return *m.MapKey
	}
	return 0
}

func (m *Pai) GetMapdes() string {
	if m != nil && m.Mapdes != nil {
		return *m.Mapdes
	}
	return ""
}

func (m *Pai) GetValue() int32 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

func (m *Pai) GetFlower() string {
	if m != nil && m.Flower != nil {
		return *m.Flower
	}
	return ""
}

func (m *Pai) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

type ZjhQueryNoSeatUser struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ZjhQueryNoSeatUser) Reset()                    { *m = ZjhQueryNoSeatUser{} }
func (m *ZjhQueryNoSeatUser) String() string            { return proto.CompactTextString(m) }
func (*ZjhQueryNoSeatUser) ProtoMessage()               {}
func (*ZjhQueryNoSeatUser) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{5} }

func (m *ZjhQueryNoSeatUser) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

// 座位
type ZjhReqSeat struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ZjhReqSeat) Reset()                    { *m = ZjhReqSeat{} }
func (m *ZjhReqSeat) String() string            { return proto.CompactTextString(m) }
func (*ZjhReqSeat) ProtoMessage()               {}
func (*ZjhReqSeat) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{6} }

func (m *ZjhReqSeat) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type ZjhMsg struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ZjhMsg) Reset()                    { *m = ZjhMsg{} }
func (m *ZjhMsg) String() string            { return proto.CompactTextString(m) }
func (*ZjhMsg) ProtoMessage()               {}
func (*ZjhMsg) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{7} }

func (m *ZjhMsg) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func init() {
	proto.RegisterType((*ZjhRoom)(nil), "bbproto.ZjhRoom")
	proto.RegisterType((*ZjhBet)(nil), "bbproto.ZjhBet")
	proto.RegisterType((*ZjhLottery)(nil), "bbproto.ZjhLottery")
	proto.RegisterType((*ZjhPai)(nil), "bbproto.ZjhPai")
	proto.RegisterType((*Pai)(nil), "bbproto.Pai")
	proto.RegisterType((*ZjhQueryNoSeatUser)(nil), "bbproto.ZjhQueryNoSeatUser")
	proto.RegisterType((*ZjhReqSeat)(nil), "bbproto.ZjhReqSeat")
	proto.RegisterType((*ZjhMsg)(nil), "bbproto.ZjhMsg")
}

var fileDescriptor10 = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x51, 0x4d, 0x4f, 0xab, 0x50,
	0x10, 0x7d, 0x94, 0x02, 0x65, 0xda, 0xbe, 0xe6, 0xdd, 0xe7, 0x02, 0x4d, 0x8c, 0x86, 0xb8, 0xe8,
	0x8a, 0x45, 0x97, 0x6e, 0x8c, 0xae, 0x4c, 0xfc, 0x08, 0x56, 0xdd, 0xb8, 0xbb, 0xb4, 0xa3, 0x40,
	0x81, 0x4b, 0xe1, 0x62, 0xd3, 0xfe, 0x05, 0xff, 0xb4, 0x97, 0x81, 0x54, 0x4d, 0x5c, 0xb0, 0x69,
	0x87, 0x99, 0x33, 0xe7, 0x9c, 0x7b, 0x06, 0x26, 0xbb, 0x90, 0xc7, 0x51, 0x16, 0x56, 0xdc, 0xcb,
	0x0b, 0x21, 0x05, 0xb3, 0x82, 0x80, 0x8a, 0x23, 0x08, 0x78, 0x89, 0x5e, 0x5b, 0x57, 0x25, 0x16,
	0x4d, 0xed, 0x7e, 0x68, 0x60, 0xbd, 0xc4, 0xe1, 0x5c, 0x88, 0x94, 0x9d, 0x81, 0x19, 0x22, 0x5f,
	0x62, 0xe1, 0x68, 0xa7, 0xda, 0x74, 0x38, 0x3b, 0xf0, 0xda, 0x6d, 0xcf, 0xaf, 0x7f, 0xaf, 0x69,
	0xc6, 0x26, 0x60, 0x15, 0xb8, 0x7e, 0xda, 0xe6, 0xe8, 0xf4, 0x14, 0xcc, 0xa8, 0x1b, 0x31, 0x5f,
	0xac, 0x72, 0x21, 0x1d, 0x5d, 0x35, 0x74, 0x76, 0x0c, 0x66, 0xc0, 0xb3, 0x95, 0xe2, 0xe9, 0x13,
	0xcf, 0x78, 0xcf, 0xf3, 0xac, 0x84, 0xd9, 0x21, 0xf4, 0x52, 0x74, 0x8c, 0x5f, 0x46, 0xee, 0x05,
	0x98, 0xca, 0xcc, 0x15, 0xca, 0xee, 0x5e, 0x02, 0x94, 0x3b, 0x91, 0xd5, 0x5e, 0xf4, 0xa9, 0xe1,
	0x6e, 0x00, 0x14, 0xc1, 0xad, 0x90, 0x12, 0x8b, 0x6d, 0x47, 0x92, 0x13, 0x30, 0x77, 0x71, 0x98,
	0xf3, 0x88, 0x38, 0x86, 0xb3, 0xc9, 0x1e, 0xa5, 0xa8, 0x7c, 0x1e, 0x91, 0x0a, 0x4f, 0x78, 0xb6,
	0x40, 0xc7, 0xa2, 0x17, 0xff, 0x03, 0x7b, 0x13, 0x65, 0x97, 0xa9, 0xa8, 0x32, 0xe9, 0x0c, 0xea,
	0x96, 0x9b, 0x90, 0xf3, 0x16, 0xad, 0xfe, 0x28, 0x1f, 0x8d, 0xd0, 0xff, 0x61, 0x98, 0x34, 0x86,
	0xbe, 0x85, 0xf6, 0x17, 0xcc, 0x02, 0xcb, 0x2a, 0x69, 0x32, 0x1b, 0xa8, 0x50, 0xf4, 0xda, 0x41,
	0x9f, 0x1c, 0x8c, 0xbe, 0x7c, 0x2a, 0xc2, 0x1f, 0x6a, 0x64, 0xc0, 0x9d, 0x83, 0x5e, 0x4f, 0x14,
	0x49, 0xca, 0xf3, 0x1b, 0xdc, 0xb6, 0x4a, 0xcd, 0xf7, 0x12, 0x4b, 0x12, 0xb1, 0xd9, 0x18, 0x8c,
	0x77, 0x9e, 0x54, 0x48, 0x1a, 0x34, 0x7e, 0x4d, 0xc4, 0xa6, 0xbd, 0x8b, 0xcd, 0x46, 0xd0, 0xcf,
	0x78, 0x7b, 0x0a, 0xdb, 0x3d, 0x07, 0xa6, 0x5e, 0xf0, 0x50, 0x29, 0x9f, 0xf7, 0xe2, 0x11, 0xb9,
	0xa4, 0x63, 0x75, 0x8a, 0xd0, 0x9d, 0x51, 0xec, 0x73, 0x5c, 0xd7, 0x7b, 0x1d, 0x77, 0x3c, 0x4a,
	0xec, 0xae, 0x7c, 0xeb, 0x86, 0xf7, 0xff, 0xf8, 0xda, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0x11,
	0x99, 0x41, 0x39, 0xe0, 0x02, 0x00, 0x00,
}
