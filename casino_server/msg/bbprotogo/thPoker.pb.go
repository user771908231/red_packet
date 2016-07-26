// Code generated by protoc-gen-go.
// source: thPoker.proto
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

// Ignoring public import of User from user.proto

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

// 德州扑克的
type ETHType int32

const (
	ETHType_TH_TYPE_GAOPAI              ETHType = 1
	ETHType_TH_TYPE_YIDUI               ETHType = 2
	ETHType_TH_TYPE_LIANGDUI            ETHType = 3
	ETHType_TH_TYPE_THREE               ETHType = 4
	ETHType_TH_TYPE_SHUNZI              ETHType = 5
	ETHType_TH_TYPE_TONGHUA             ETHType = 6
	ETHType_TH_TYPE_HULU                ETHType = 7
	ETHType_TH_TYPE_SITIAO              ETHType = 8
	ETHType_TH_TYPE_TONGHUASHUN         ETHType = 9
	ETHType_TH_TYPE_HUANGSHITONGHUASHUN ETHType = 10
)

var ETHType_name = map[int32]string{
	1:  "TH_TYPE_GAOPAI",
	2:  "TH_TYPE_YIDUI",
	3:  "TH_TYPE_LIANGDUI",
	4:  "TH_TYPE_THREE",
	5:  "TH_TYPE_SHUNZI",
	6:  "TH_TYPE_TONGHUA",
	7:  "TH_TYPE_HULU",
	8:  "TH_TYPE_SITIAO",
	9:  "TH_TYPE_TONGHUASHUN",
	10: "TH_TYPE_HUANGSHITONGHUASHUN",
}
var ETHType_value = map[string]int32{
	"TH_TYPE_GAOPAI":              1,
	"TH_TYPE_YIDUI":               2,
	"TH_TYPE_LIANGDUI":            3,
	"TH_TYPE_THREE":               4,
	"TH_TYPE_SHUNZI":              5,
	"TH_TYPE_TONGHUA":             6,
	"TH_TYPE_HULU":                7,
	"TH_TYPE_SITIAO":              8,
	"TH_TYPE_TONGHUASHUN":         9,
	"TH_TYPE_HUANGSHITONGHUASHUN": 10,
}

func (x ETHType) Enum() *ETHType {
	p := new(ETHType)
	*p = x
	return p
}
func (x ETHType) String() string {
	return proto.EnumName(ETHType_name, int32(x))
}
func (x *ETHType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ETHType_value, data, "ETHType")
	if err != nil {
		return err
	}
	*x = ETHType(value)
	return nil
}
func (ETHType) EnumDescriptor() ([]byte, []int) { return fileDescriptor13, []int{0} }

type ETHAction int32

const (
	ETHAction_TH_DESK_BET_TYPE_BET      ETHAction = 1
	ETHAction_TH_DESK_BET_TYPE_CALL     ETHAction = 2
	ETHAction_TH_DESK_BET_TYPE_FOLD     ETHAction = 3
	ETHAction_TH_DESK_BET_TYPE_CHECK    ETHAction = 4
	ETHAction_TH_DESK_BET_TYPE_RAISE    ETHAction = 5
	ETHAction_TH_DESK_BET_TYPE_RERRAISE ETHAction = 6
	ETHAction_TH_DESK_BET_TYPE_ALLIN    ETHAction = 7
)

var ETHAction_name = map[int32]string{
	1: "TH_DESK_BET_TYPE_BET",
	2: "TH_DESK_BET_TYPE_CALL",
	3: "TH_DESK_BET_TYPE_FOLD",
	4: "TH_DESK_BET_TYPE_CHECK",
	5: "TH_DESK_BET_TYPE_RAISE",
	6: "TH_DESK_BET_TYPE_RERRAISE",
	7: "TH_DESK_BET_TYPE_ALLIN",
}
var ETHAction_value = map[string]int32{
	"TH_DESK_BET_TYPE_BET":      1,
	"TH_DESK_BET_TYPE_CALL":     2,
	"TH_DESK_BET_TYPE_FOLD":     3,
	"TH_DESK_BET_TYPE_CHECK":    4,
	"TH_DESK_BET_TYPE_RAISE":    5,
	"TH_DESK_BET_TYPE_RERRAISE": 6,
	"TH_DESK_BET_TYPE_ALLIN":    7,
}

func (x ETHAction) Enum() *ETHAction {
	p := new(ETHAction)
	*p = x
	return p
}
func (x ETHAction) String() string {
	return proto.EnumName(ETHAction_name, int32(x))
}
func (x *ETHAction) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ETHAction_value, data, "ETHAction")
	if err != nil {
		return err
	}
	*x = ETHAction(value)
	return nil
}
func (ETHAction) EnumDescriptor() ([]byte, []int) { return fileDescriptor13, []int{1} }

type ThRoom struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	ReqType          *int32       `protobuf:"varint,2,opt,name=reqType" json:"reqType,omitempty"`
	Jackpot          *int64       `protobuf:"varint,3,opt,name=jackpot" json:"jackpot,omitempty"`
	Users            []*THUser    `protobuf:"bytes,4,rep,name=users" json:"users,omitempty"`
	BetTime          *int32       `protobuf:"varint,5,opt,name=betTime" json:"betTime,omitempty"`
	DeskStatus       *int32       `protobuf:"varint,6,opt,name=deskStatus" json:"deskStatus,omitempty"`
	DeskNumber       *int32       `protobuf:"varint,7,opt,name=deskNumber" json:"deskNumber,omitempty"`
	PublicPais       []*Pai       `protobuf:"bytes,8,rep,name=publicPais" json:"publicPais,omitempty"`
	BetUserId        *uint32      `protobuf:"varint,9,opt,name=betUserId" json:"betUserId,omitempty"`
	BigBlind         *uint32      `protobuf:"varint,10,opt,name=bigBlind" json:"bigBlind,omitempty"`
	SmallBlind       *uint32      `protobuf:"varint,11,opt,name=smallBlind" json:"smallBlind,omitempty"`
	Butten           *uint32      `protobuf:"varint,12,opt,name=butten" json:"butten,omitempty"`
	BetRemainTime    *int32       `protobuf:"varint,13,opt,name=betRemainTime" json:"betRemainTime,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ThRoom) Reset()                    { *m = ThRoom{} }
func (m *ThRoom) String() string            { return proto.CompactTextString(m) }
func (*ThRoom) ProtoMessage()               {}
func (*ThRoom) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{0} }

func (m *ThRoom) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ThRoom) GetReqType() int32 {
	if m != nil && m.ReqType != nil {
		return *m.ReqType
	}
	return 0
}

func (m *ThRoom) GetJackpot() int64 {
	if m != nil && m.Jackpot != nil {
		return *m.Jackpot
	}
	return 0
}

func (m *ThRoom) GetUsers() []*THUser {
	if m != nil {
		return m.Users
	}
	return nil
}

func (m *ThRoom) GetBetTime() int32 {
	if m != nil && m.BetTime != nil {
		return *m.BetTime
	}
	return 0
}

func (m *ThRoom) GetDeskStatus() int32 {
	if m != nil && m.DeskStatus != nil {
		return *m.DeskStatus
	}
	return 0
}

func (m *ThRoom) GetDeskNumber() int32 {
	if m != nil && m.DeskNumber != nil {
		return *m.DeskNumber
	}
	return 0
}

func (m *ThRoom) GetPublicPais() []*Pai {
	if m != nil {
		return m.PublicPais
	}
	return nil
}

func (m *ThRoom) GetBetUserId() uint32 {
	if m != nil && m.BetUserId != nil {
		return *m.BetUserId
	}
	return 0
}

func (m *ThRoom) GetBigBlind() uint32 {
	if m != nil && m.BigBlind != nil {
		return *m.BigBlind
	}
	return 0
}

func (m *ThRoom) GetSmallBlind() uint32 {
	if m != nil && m.SmallBlind != nil {
		return *m.SmallBlind
	}
	return 0
}

func (m *ThRoom) GetButten() uint32 {
	if m != nil && m.Butten != nil {
		return *m.Butten
	}
	return 0
}

func (m *ThRoom) GetBetRemainTime() int32 {
	if m != nil && m.BetRemainTime != nil {
		return *m.BetRemainTime
	}
	return 0
}

// 刚进来的人
type THRoomAddUserBroadcast struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	User             *THUser      `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *THRoomAddUserBroadcast) Reset()                    { *m = THRoomAddUserBroadcast{} }
func (m *THRoomAddUserBroadcast) String() string            { return proto.CompactTextString(m) }
func (*THRoomAddUserBroadcast) ProtoMessage()               {}
func (*THRoomAddUserBroadcast) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{1} }

func (m *THRoomAddUserBroadcast) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *THRoomAddUserBroadcast) GetUser() *THUser {
	if m != nil {
		return m.User
	}
	return nil
}

type THPoker struct {
	ThType           *int32 `protobuf:"varint,1,opt,name=thType" json:"thType,omitempty"`
	Pais             []*Pai `protobuf:"bytes,2,rep,name=pais" json:"pais,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *THPoker) Reset()                    { *m = THPoker{} }
func (m *THPoker) String() string            { return proto.CompactTextString(m) }
func (*THPoker) ProtoMessage()               {}
func (*THPoker) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{2} }

func (m *THPoker) GetThType() int32 {
	if m != nil && m.ThType != nil {
		return *m.ThType
	}
	return 0
}

func (m *THPoker) GetPais() []*Pai {
	if m != nil {
		return m.Pais
	}
	return nil
}

type THBet struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	BetType          *int32       `protobuf:"varint,2,opt,name=betType" json:"betType,omitempty"`
	BetAmount        *int64       `protobuf:"varint,3,opt,name=betAmount" json:"betAmount,omitempty"`
	NextBetUser      *uint32      `protobuf:"varint,4,opt,name=nextBetUser" json:"nextBetUser,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *THBet) Reset()                    { *m = THBet{} }
func (m *THBet) String() string            { return proto.CompactTextString(m) }
func (*THBet) ProtoMessage()               {}
func (*THBet) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{3} }

func (m *THBet) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *THBet) GetBetType() int32 {
	if m != nil && m.BetType != nil {
		return *m.BetType
	}
	return 0
}

func (m *THBet) GetBetAmount() int64 {
	if m != nil && m.BetAmount != nil {
		return *m.BetAmount
	}
	return 0
}

func (m *THBet) GetNextBetUser() uint32 {
	if m != nil && m.NextBetUser != nil {
		return *m.NextBetUser
	}
	return 0
}

type THBetBegin struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Dealer           *uint32      `protobuf:"varint,2,opt,name=dealer" json:"dealer,omitempty"`
	BigBlind         *THUser      `protobuf:"bytes,3,opt,name=bigBlind" json:"bigBlind,omitempty"`
	SmallBlind       *THUser      `protobuf:"bytes,4,opt,name=smallBlind" json:"smallBlind,omitempty"`
	BetUserNow       *uint32      `protobuf:"varint,5,opt,name=betUserNow" json:"betUserNow,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *THBetBegin) Reset()                    { *m = THBetBegin{} }
func (m *THBetBegin) String() string            { return proto.CompactTextString(m) }
func (*THBetBegin) ProtoMessage()               {}
func (*THBetBegin) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{4} }

func (m *THBetBegin) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *THBetBegin) GetDealer() uint32 {
	if m != nil && m.Dealer != nil {
		return *m.Dealer
	}
	return 0
}

func (m *THBetBegin) GetBigBlind() *THUser {
	if m != nil {
		return m.BigBlind
	}
	return nil
}

func (m *THBetBegin) GetSmallBlind() *THUser {
	if m != nil {
		return m.SmallBlind
	}
	return nil
}

func (m *THBetBegin) GetBetUserNow() uint32 {
	if m != nil && m.BetUserNow != nil {
		return *m.BetUserNow
	}
	return 0
}

type THBetBroadcast struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	BetType          *int32       `protobuf:"varint,2,opt,name=betType" json:"betType,omitempty"`
	BetAmount        *int64       `protobuf:"varint,3,opt,name=betAmount" json:"betAmount,omitempty"`
	User             *THUser      `protobuf:"bytes,4,opt,name=user" json:"user,omitempty"`
	NextBetUserId    *uint32      `protobuf:"varint,5,opt,name=nextBetUserId" json:"nextBetUserId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *THBetBroadcast) Reset()                    { *m = THBetBroadcast{} }
func (m *THBetBroadcast) String() string            { return proto.CompactTextString(m) }
func (*THBetBroadcast) ProtoMessage()               {}
func (*THBetBroadcast) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{5} }

func (m *THBetBroadcast) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *THBetBroadcast) GetBetType() int32 {
	if m != nil && m.BetType != nil {
		return *m.BetType
	}
	return 0
}

func (m *THBetBroadcast) GetBetAmount() int64 {
	if m != nil && m.BetAmount != nil {
		return *m.BetAmount
	}
	return 0
}

func (m *THBetBroadcast) GetUser() *THUser {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *THBetBroadcast) GetNextBetUserId() uint32 {
	if m != nil && m.NextBetUserId != nil {
		return *m.NextBetUserId
	}
	return 0
}

// 德州扑克中的User
type THUser struct {
	User             *User    `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
	HandPais         []*Pai   `protobuf:"bytes,2,rep,name=handPais" json:"handPais,omitempty"`
	Status           *int32   `protobuf:"varint,3,opt,name=status" json:"status,omitempty"`
	SeatNumber       *int32   `protobuf:"varint,4,opt,name=seatNumber" json:"seatNumber,omitempty"`
	Thpoker          *THPoker `protobuf:"bytes,5,opt,name=thpoker" json:"thpoker,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *THUser) Reset()                    { *m = THUser{} }
func (m *THUser) String() string            { return proto.CompactTextString(m) }
func (*THUser) ProtoMessage()               {}
func (*THUser) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{6} }

func (m *THUser) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *THUser) GetHandPais() []*Pai {
	if m != nil {
		return m.HandPais
	}
	return nil
}

func (m *THUser) GetStatus() int32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return 0
}

func (m *THUser) GetSeatNumber() int32 {
	if m != nil && m.SeatNumber != nil {
		return *m.SeatNumber
	}
	return 0
}

func (m *THUser) GetThpoker() *THPoker {
	if m != nil {
		return m.Thpoker
	}
	return nil
}

// 开奖时需要用到的协议
type THLottery struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Users            []*THUser    `protobuf:"bytes,2,rep,name=users" json:"users,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *THLottery) Reset()                    { *m = THLottery{} }
func (m *THLottery) String() string            { return proto.CompactTextString(m) }
func (*THLottery) ProtoMessage()               {}
func (*THLottery) Descriptor() ([]byte, []int) { return fileDescriptor13, []int{7} }

func (m *THLottery) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *THLottery) GetUsers() []*THUser {
	if m != nil {
		return m.Users
	}
	return nil
}

func init() {
	proto.RegisterType((*ThRoom)(nil), "bbproto.ThRoom")
	proto.RegisterType((*THRoomAddUserBroadcast)(nil), "bbproto.THRoomAddUserBroadcast")
	proto.RegisterType((*THPoker)(nil), "bbproto.THPoker")
	proto.RegisterType((*THBet)(nil), "bbproto.THBet")
	proto.RegisterType((*THBetBegin)(nil), "bbproto.THBetBegin")
	proto.RegisterType((*THBetBroadcast)(nil), "bbproto.THBetBroadcast")
	proto.RegisterType((*THUser)(nil), "bbproto.THUser")
	proto.RegisterType((*THLottery)(nil), "bbproto.THLottery")
	proto.RegisterEnum("bbproto.ETHType", ETHType_name, ETHType_value)
	proto.RegisterEnum("bbproto.ETHAction", ETHAction_name, ETHAction_value)
}

var fileDescriptor13 = []byte{
	// 731 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x54, 0xdd, 0x52, 0xda, 0x50,
	0x10, 0x6e, 0xf8, 0x0b, 0x2c, 0x46, 0xd2, 0xa8, 0x6d, 0xc4, 0xd1, 0x2a, 0xed, 0x45, 0xc7, 0x0b,
	0x2e, 0x9c, 0xe9, 0x03, 0x04, 0x4d, 0x4d, 0x46, 0x06, 0x52, 0x08, 0x17, 0x76, 0xa6, 0xe3, 0x9c,
	0x90, 0x53, 0x13, 0x85, 0x84, 0x92, 0xc3, 0xb4, 0xf6, 0x19, 0x7a, 0xd1, 0x27, 0xe8, 0x0b, 0xf4,
	0x5d, 0xfa, 0x10, 0x7d, 0x92, 0xee, 0x39, 0x01, 0x44, 0xc1, 0x19, 0x7a, 0xe3, 0xc4, 0x6f, 0x77,
	0xbf, 0xf3, 0xed, 0xb7, 0xbb, 0x80, 0xc2, 0x02, 0x27, 0xbe, 0xa5, 0xe3, 0xfa, 0x68, 0x1c, 0xb3,
	0x58, 0x93, 0x3d, 0x4f, 0x7c, 0x54, 0xc1, 0x23, 0x09, 0xad, 0x4f, 0xbf, 0x27, 0xc9, 0x2c, 0xa1,
	0x5a, 0xf9, 0x1e, 0x90, 0x9b, 0x30, 0x0a, 0x26, 0x24, 0x05, 0x6a, 0xbf, 0x33, 0x50, 0x70, 0x83,
	0x4e, 0x1c, 0x0f, 0xb5, 0x37, 0x50, 0x08, 0x28, 0xf1, 0xe9, 0x58, 0x97, 0x0e, 0xa5, 0xb7, 0xe5,
	0x93, 0xed, 0xfa, 0x94, 0xad, 0xee, 0xf0, 0xbf, 0x96, 0x88, 0x69, 0x15, 0x90, 0xc7, 0xf4, 0x8b,
	0x7b, 0x37, 0xa2, 0x7a, 0x06, 0xd3, 0xf2, 0x1c, 0xb8, 0x21, 0xfd, 0xdb, 0x51, 0xcc, 0xf4, 0x2c,
	0x02, 0x59, 0xed, 0x00, 0xf2, 0xfc, 0xc5, 0x44, 0xcf, 0x1d, 0x66, 0x91, 0xa6, 0x32, 0xa7, 0x71,
	0xad, 0x5e, 0x92, 0x32, 0x78, 0x94, 0xb9, 0xe1, 0x90, 0xea, 0x79, 0xc1, 0xa0, 0x01, 0xf8, 0x34,
	0xb9, 0xed, 0x32, 0xc2, 0x26, 0x89, 0x5e, 0x58, 0xc4, 0x5a, 0x93, 0xa1, 0x87, 0x82, 0x64, 0x81,
	0x1d, 0x02, 0x8c, 0x26, 0xde, 0x20, 0xec, 0x3b, 0x24, 0x4c, 0xf4, 0xa2, 0x60, 0xdf, 0xb8, 0x17,
	0x49, 0x42, 0xed, 0x39, 0x94, 0x90, 0x9a, 0xbf, 0x62, 0xfb, 0x7a, 0x09, 0x8b, 0x14, 0x4d, 0x85,
	0xa2, 0x17, 0x5e, 0x37, 0x06, 0x61, 0xe4, 0xeb, 0x20, 0x10, 0xa4, 0x4e, 0x86, 0x64, 0x30, 0x48,
	0xb1, 0xb2, 0xc0, 0x36, 0xa1, 0xe0, 0x4d, 0x18, 0xa3, 0x91, 0xbe, 0x21, 0xfe, 0xdf, 0x01, 0x05,
	0x89, 0x3a, 0x74, 0x48, 0xc2, 0x48, 0x28, 0x55, 0xb8, 0x82, 0xda, 0x27, 0x78, 0xe1, 0x5a, 0xdc,
	0x2c, 0xc3, 0xf7, 0xf9, 0x2b, 0x8d, 0x71, 0x4c, 0xfc, 0x3e, 0x49, 0xd8, 0x9a, 0xe6, 0xed, 0x43,
	0x8e, 0x5b, 0x23, 0x9c, 0x5b, 0x76, 0xa6, 0xf6, 0x0e, 0x64, 0xd7, 0x12, 0xf3, 0xe4, 0x82, 0x58,
	0x20, 0x5c, 0x96, 0x44, 0xef, 0x55, 0xc8, 0x8d, 0x78, 0xd7, 0x99, 0xe5, 0xae, 0x6b, 0x9f, 0x21,
	0xef, 0x5a, 0x0d, 0xca, 0xd6, 0x9f, 0x20, 0xf7, 0xff, 0x7e, 0x82, 0xa9, 0x6b, 0xc6, 0x30, 0x9e,
	0x44, 0xb3, 0x19, 0x6e, 0x41, 0x39, 0xa2, 0xdf, 0x58, 0x23, 0x35, 0x13, 0x27, 0x89, 0xa6, 0xd4,
	0x7e, 0x49, 0x00, 0xe2, 0xa1, 0x06, 0xbd, 0x0e, 0xa3, 0x35, 0x5f, 0xc3, 0x46, 0x7c, 0x4a, 0x06,
	0xd3, 0xa6, 0x15, 0xed, 0x68, 0x61, 0x1e, 0xd9, 0x95, 0x36, 0x68, 0xaf, 0x1f, 0x0c, 0x28, 0xb7,
	0x3a, 0x09, 0xa7, 0x38, 0x1d, 0x75, 0x2b, 0xfe, 0x2a, 0x16, 0x49, 0xa9, 0xfd, 0x94, 0x60, 0x33,
	0x15, 0xf8, 0x9f, 0x73, 0x59, 0xc7, 0x92, 0xd9, 0xec, 0x9e, 0xd0, 0x83, 0x1b, 0xb3, 0xe0, 0x18,
	0xae, 0x5f, 0x2a, 0xe9, 0x87, 0x84, 0xf7, 0x95, 0x66, 0xec, 0x4d, 0x09, 0x52, 0x21, 0xca, 0x9c,
	0x40, 0x04, 0x0f, 0xa0, 0x18, 0x90, 0xc8, 0x77, 0x9e, 0x98, 0x31, 0xb7, 0x31, 0x49, 0xef, 0x23,
	0x3b, 0xbb, 0x8f, 0x84, 0x12, 0x36, 0xbd, 0x8f, 0x9c, 0xc0, 0x8e, 0x40, 0x66, 0xc1, 0x88, 0xaf,
	0x8f, 0x78, 0xbc, 0x7c, 0xa2, 0x2e, 0x88, 0x14, 0x6b, 0x55, 0xfb, 0x00, 0x25, 0xd7, 0x6a, 0xc6,
	0xb8, 0xe9, 0xe3, 0xbb, 0x35, 0xbd, 0x99, 0x9f, 0x73, 0x66, 0xe5, 0x39, 0x1f, 0xff, 0x95, 0x40,
	0x36, 0x5d, 0x8b, 0x9b, 0x87, 0xaa, 0xd0, 0xff, 0x2b, 0xf7, 0xd2, 0x31, 0xaf, 0xce, 0x8d, 0xb6,
	0x63, 0xd8, 0xaa, 0x84, 0x56, 0x2a, 0x33, 0xec, 0xd2, 0x3e, 0xeb, 0xd9, 0x6a, 0x46, 0xdb, 0x06,
	0x75, 0x06, 0x35, 0x6d, 0xa3, 0x75, 0xce, 0xd1, 0xec, 0x62, 0x22, 0x1e, 0x99, 0x69, 0xaa, 0xb9,
	0x45, 0xbe, 0xae, 0xd5, 0x6b, 0x7d, 0xb4, 0xd5, 0x3c, 0xae, 0x66, 0x65, 0x9e, 0xd6, 0x6e, 0x9d,
	0x5b, 0x3d, 0x43, 0x2d, 0xe0, 0x95, 0x6f, 0xcc, 0x40, 0xab, 0xd7, 0xec, 0xa9, 0xf2, 0x83, 0x52,
	0xdb, 0xb5, 0x8d, 0xb6, 0x5a, 0xd4, 0x5e, 0xc2, 0xd6, 0xa3, 0x52, 0xce, 0xaa, 0x96, 0xb4, 0x57,
	0xb0, 0x77, 0x5f, 0x8e, 0x82, 0xba, 0x96, 0xbd, 0x98, 0x00, 0xc7, 0x7f, 0x24, 0x28, 0x61, 0x93,
	0x46, 0x9f, 0x85, 0x71, 0xa4, 0xe9, 0xb0, 0x8d, 0xe9, 0x67, 0x66, 0xf7, 0xe2, 0xaa, 0x61, 0xba,
	0x69, 0x1d, 0x7e, 0x60, 0xb3, 0xbb, 0xb0, 0xb3, 0x14, 0x39, 0x35, 0x9a, 0x4d, 0x6c, 0x7a, 0x55,
	0xe8, 0x7d, 0xbb, 0x79, 0x86, 0x9d, 0x57, 0xf9, 0xcf, 0xca, 0xe3, 0x2a, 0xcb, 0x3c, 0xbd, 0x40,
	0x0b, 0x56, 0xc5, 0x3a, 0x86, 0xdd, 0x35, 0xd1, 0x8a, 0x7d, 0xd8, 0x5d, 0x8e, 0x99, 0x9d, 0x34,
	0x5c, 0x58, 0x59, 0x8a, 0x5a, 0xec, 0x96, 0x2a, 0x3b, 0xcf, 0x1c, 0xc9, 0xc9, 0xfc, 0x0b, 0x00,
	0x00, 0xff, 0xff, 0xee, 0x74, 0x38, 0x23, 0x3f, 0x06, 0x00, 0x00,
}
