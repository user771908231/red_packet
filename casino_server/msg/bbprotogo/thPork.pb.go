// Code generated by protoc-gen-go.
// source: thPork.proto
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
func (ETHType) EnumDescriptor() ([]byte, []int) { return fileDescriptor10, []int{0} }

type ThRoom struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	ReqType          *int32       `protobuf:"varint,2,opt,name=reqType" json:"reqType,omitempty"`
	Jackpot          *int64       `protobuf:"varint,3,opt,name=jackpot" json:"jackpot,omitempty"`
	Banker           *User        `protobuf:"bytes,4,opt,name=banker" json:"banker,omitempty"`
	Me               *User        `protobuf:"bytes,5,opt,name=me" json:"me,omitempty"`
	Zjhpai           []*ZjhPai    `protobuf:"bytes,6,rep,name=zjhpai" json:"zjhpai,omitempty"`
	BetTime          *int32       `protobuf:"varint,7,opt,name=betTime" json:"betTime,omitempty"`
	RoomStatus       *int32       `protobuf:"varint,8,opt,name=roomStatus" json:"roomStatus,omitempty"`
	PublicPais       []*Pai       `protobuf:"bytes,9,rep,name=publicPais" json:"publicPais,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ThRoom) Reset()                    { *m = ThRoom{} }
func (m *ThRoom) String() string            { return proto.CompactTextString(m) }
func (*ThRoom) ProtoMessage()               {}
func (*ThRoom) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{0} }

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

func (m *ThRoom) GetBanker() *User {
	if m != nil {
		return m.Banker
	}
	return nil
}

func (m *ThRoom) GetMe() *User {
	if m != nil {
		return m.Me
	}
	return nil
}

func (m *ThRoom) GetZjhpai() []*ZjhPai {
	if m != nil {
		return m.Zjhpai
	}
	return nil
}

func (m *ThRoom) GetBetTime() int32 {
	if m != nil && m.BetTime != nil {
		return *m.BetTime
	}
	return 0
}

func (m *ThRoom) GetRoomStatus() int32 {
	if m != nil && m.RoomStatus != nil {
		return *m.RoomStatus
	}
	return 0
}

func (m *ThRoom) GetPublicPais() []*Pai {
	if m != nil {
		return m.PublicPais
	}
	return nil
}

type THPork struct {
	ThType           *int32 `protobuf:"varint,1,opt,name=thType" json:"thType,omitempty"`
	Pais             []*Pai `protobuf:"bytes,2,rep,name=pais" json:"pais,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *THPork) Reset()                    { *m = THPork{} }
func (m *THPork) String() string            { return proto.CompactTextString(m) }
func (*THPork) ProtoMessage()               {}
func (*THPork) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{1} }

func (m *THPork) GetThType() int32 {
	if m != nil && m.ThType != nil {
		return *m.ThType
	}
	return 0
}

func (m *THPork) GetPais() []*Pai {
	if m != nil {
		return m.Pais
	}
	return nil
}

type THBet struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	BetType          *int32       `protobuf:"varint,2,opt,name=betType" json:"betType,omitempty"`
	BetAmount        *int32       `protobuf:"varint,3,opt,name=betAmount" json:"betAmount,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *THBet) Reset()                    { *m = THBet{} }
func (m *THBet) String() string            { return proto.CompactTextString(m) }
func (*THBet) ProtoMessage()               {}
func (*THBet) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{2} }

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

func (m *THBet) GetBetAmount() int32 {
	if m != nil && m.BetAmount != nil {
		return *m.BetAmount
	}
	return 0
}

type THBegin struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Users            []*THUser    `protobuf:"bytes,2,rep,name=users" json:"users,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *THBegin) Reset()                    { *m = THBegin{} }
func (m *THBegin) String() string            { return proto.CompactTextString(m) }
func (*THBegin) ProtoMessage()               {}
func (*THBegin) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{3} }

func (m *THBegin) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *THBegin) GetUsers() []*THUser {
	if m != nil {
		return m.Users
	}
	return nil
}

// 德州扑克中的User
type THUser struct {
	U                *User  `protobuf:"bytes,1,opt,name=u" json:"u,omitempty"`
	HandPais         []*Pai `protobuf:"bytes,2,rep,name=handPais" json:"handPais,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *THUser) Reset()                    { *m = THUser{} }
func (m *THUser) String() string            { return proto.CompactTextString(m) }
func (*THUser) ProtoMessage()               {}
func (*THUser) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{4} }

func (m *THUser) GetU() *User {
	if m != nil {
		return m.U
	}
	return nil
}

func (m *THUser) GetHandPais() []*Pai {
	if m != nil {
		return m.HandPais
	}
	return nil
}

func init() {
	proto.RegisterType((*ThRoom)(nil), "bbproto.ThRoom")
	proto.RegisterType((*THPork)(nil), "bbproto.THPork")
	proto.RegisterType((*THBet)(nil), "bbproto.THBet")
	proto.RegisterType((*THBegin)(nil), "bbproto.THBegin")
	proto.RegisterType((*THUser)(nil), "bbproto.THUser")
	proto.RegisterEnum("bbproto.ETHType", ETHType_name, ETHType_value)
}

var fileDescriptor10 = []byte{
	// 460 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x52, 0xd1, 0x6e, 0xda, 0x30,
	0x14, 0x5d, 0x02, 0x49, 0xe0, 0x8e, 0x16, 0xcf, 0xad, 0xb4, 0xac, 0xd3, 0x3a, 0x84, 0xf6, 0x30,
	0xed, 0x81, 0x87, 0x69, 0x3f, 0x90, 0x6a, 0x08, 0x47, 0xaa, 0xc0, 0x2b, 0xce, 0x43, 0xfb, 0x52,
	0x39, 0xad, 0x35, 0x07, 0x46, 0x92, 0x05, 0xe7, 0x61, 0xfd, 0xdc, 0x7d, 0xc2, 0xbe, 0x60, 0xb6,
	0x0b, 0x84, 0x55, 0x9d, 0xc4, 0x0b, 0xba, 0x3a, 0xf7, 0x9e, 0xe3, 0x73, 0x0e, 0x81, 0x9e, 0x92,
	0xb4, 0xa8, 0x96, 0xa3, 0xb2, 0x2a, 0x54, 0x81, 0x83, 0x34, 0xb5, 0xc3, 0x19, 0xa4, 0x7c, 0x2d,
	0x46, 0x9b, 0xb9, 0x5e, 0x8b, 0x6a, 0x33, 0xf7, 0x1f, 0x24, 0x5f, 0x64, 0xb9, 0xac, 0xf9, 0x23,
	0x30, 0xfc, 0xe3, 0x80, 0xcf, 0xe4, 0x55, 0x51, 0xac, 0xf0, 0x07, 0xf0, 0xa5, 0xe0, 0xf7, 0xa2,
	0x0a, 0x9d, 0x81, 0xf3, 0xf1, 0xe5, 0xe7, 0xd3, 0xd1, 0x46, 0x6d, 0x44, 0xcd, 0x2f, 0xb1, 0x3b,
	0xdc, 0x87, 0xa0, 0x12, 0x3f, 0xd9, 0xaf, 0x52, 0x84, 0xae, 0x3e, 0xf3, 0x0c, 0xb0, 0xe0, 0x77,
	0xcb, 0xb2, 0x50, 0x61, 0x4b, 0x03, 0x2d, 0xfc, 0x0e, 0xfc, 0x94, 0xe7, 0x4b, 0xad, 0xd3, 0xb6,
	0x3a, 0x47, 0x3b, 0x9d, 0x44, 0x1b, 0xc1, 0x6f, 0xc0, 0x5d, 0x89, 0xd0, 0x7b, 0x6e, 0xf5, 0x1e,
	0xfc, 0x87, 0x85, 0x2c, 0x79, 0x16, 0xfa, 0x83, 0x96, 0x5e, 0xf7, 0x77, 0xeb, 0x9b, 0x85, 0xa4,
	0x3c, 0x33, 0x6f, 0xa5, 0x42, 0xb1, 0x4c, 0x0b, 0x04, 0xf6, 0x71, 0x0c, 0x50, 0x69, 0xef, 0x73,
	0xc5, 0x55, 0xbd, 0x0e, 0x3b, 0x16, 0x1b, 0x00, 0x94, 0x75, 0xfa, 0x23, 0xbb, 0xd3, 0x8c, 0x75,
	0xd8, 0xb5, 0x4a, 0xbd, 0x26, 0x0b, 0xcf, 0x86, 0x5f, 0x74, 0x66, 0x62, 0x6a, 0xc3, 0xc7, 0xe0,
	0x2b, 0x69, 0xc3, 0x38, 0x96, 0x7b, 0x06, 0xed, 0xd2, 0xb0, 0xdc, 0x67, 0x58, 0xdf, 0xc0, 0x63,
	0xe4, 0x42, 0xa8, 0xc3, 0x8b, 0x32, 0x5e, 0x9b, 0xa2, 0x5e, 0x41, 0x57, 0x03, 0xd1, 0xaa, 0xa8,
	0xf3, 0xc7, 0xaa, 0xbc, 0xe1, 0x0c, 0x02, 0x23, 0xf9, 0x3d, 0xcb, 0x0f, 0x14, 0x3d, 0x07, 0xcf,
	0xfc, 0x9b, 0x5b, 0x83, 0x4d, 0x41, 0x8c, 0x98, 0x06, 0x87, 0x17, 0x26, 0x99, 0xed, 0x32, 0x04,
	0xa7, 0xde, 0x48, 0x3d, 0x69, 0xf9, 0x1c, 0x3a, 0x92, 0xe7, 0xf7, 0xf4, 0x3f, 0x39, 0x3f, 0xfd,
	0x76, 0x20, 0x18, 0x33, 0x62, 0x9c, 0xeb, 0x7e, 0x8f, 0x19, 0xb9, 0x65, 0xd7, 0x74, 0x7c, 0x3b,
	0x89, 0x66, 0x34, 0x8a, 0x91, 0xa3, 0x73, 0x1c, 0x6d, 0xb1, 0xeb, 0xf8, 0x6b, 0x12, 0x23, 0x17,
	0x9f, 0x02, 0xda, 0x42, 0x97, 0x71, 0x34, 0x9d, 0x18, 0xb4, 0xb5, 0x7f, 0xc8, 0xc8, 0xd5, 0x78,
	0x8c, 0xda, 0xfb, 0x7a, 0x73, 0x92, 0x4c, 0x6f, 0x62, 0xe4, 0xe1, 0x13, 0xe8, 0xef, 0xce, 0x66,
	0xd3, 0x09, 0x49, 0x22, 0xe4, 0x63, 0x04, 0xbd, 0x2d, 0x48, 0x92, 0xcb, 0x04, 0x05, 0xff, 0x50,
	0x63, 0x16, 0x47, 0x33, 0xd4, 0xc1, 0xaf, 0xe1, 0xe4, 0x09, 0xd5, 0xa8, 0xa2, 0xae, 0xfe, 0x92,
	0xde, 0x36, 0x74, 0x6d, 0x68, 0x4e, 0xe2, 0xfd, 0x03, 0xa0, 0x2f, 0xa8, 0x43, 0xdd, 0xbf, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x20, 0x05, 0x81, 0x17, 0x3e, 0x03, 0x00, 0x00,
}