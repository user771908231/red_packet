// Code generated by protoc-gen-go.
// source: base.proto
// DO NOT EDIT!

/*
Package mjproto is a generated protocol buffer package.

It is generated from these files:
	base.proto
	mahjong_desk.proto
	mahjong_hall.proto
	mahjong_play.proto

It has these top-level messages:
	ProtoHeader
	WeixinInfo
	CardInfo
	PlayOptions
	RoomTypeInfo
	ComposeCard
	PlayerCard
	PlayerInfo
	DeskGameInfo
	Game_DissolveDesk
	Game_AckDissolveDesk
	Game_LeaveDesk
	Game_AckLeaveDesk
	Game_Ready
	Game_AckReady
	Game_Message
	Game_SendMessage
	WinCoinInfo
	EndLotteryInfo
	Game_SendCurrentResult
	Game_SendEndLottery
	ServerInfo
	Game_QuickConn
	Game_AckQuickConn
	Game_Login
	Game_AckLogin
	Game_Notice
	Game_AckNotice
	Game_GameRecord
	Game_BeanUserRecord
	Game_BeanGameRecord
	Game_AckGameRecord
	Game_Feedback
	Game_CreateRoom
	Game_AckCreateRoom
	Game_EnterRoom
	Game_AckEnterRoom
	Game_Opening
	Game_DealCards
	Game_ExchangeCards
	Game_AckExchangeCards
	Game_DingQue
	Game_BroadcastBeginDingQue
	Game_BroadcastBeginExchange
	Game_GetInCard
	Game_SendOutCard
	Game_AckSendOutCard
	Game_ActPeng
	Game_AckActPeng
	Game_ActGang
	Game_AckActGang
	Game_ActHu
	Game_AckActHu
	Game_ActGuo
	Game_AckActGuo
	Game_OverTurn
	Game_SendGameInfo
*/
package mjproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EProtoId int32

const (
	EProtoId_PID_QUICK_CONN               EProtoId = 1
	EProtoId_PID_QUICK_CONN_ACK           EProtoId = 2
	EProtoId_PID_GAME_LOGIN               EProtoId = 3
	EProtoId_PID_GAME_LOGIN_ACK           EProtoId = 4
	EProtoId_PID_CREATEROOM               EProtoId = 5
	EProtoId_PID_CREATEROOM_ACK           EProtoId = 6
	EProtoId_PID_ENTER_ROOM               EProtoId = 7
	EProtoId_PID_ENTER_ROOM_ACK           EProtoId = 8
	EProtoId_PID_SEND_GAMEINFO            EProtoId = 9
	EProtoId_PID_READY                    EProtoId = 10
	EProtoId_PID_READY_ACK                EProtoId = 11
	EProtoId_PID_EXCHANGECARDS            EProtoId = 12
	EProtoId_PID_EXCHANGECARDS_ACK        EProtoId = 13
	EProtoId_PID_DINGQUE                  EProtoId = 14
	EProtoId_PID_OPENING                  EProtoId = 15
	EProtoId_PID_DEAL_CARDS               EProtoId = 16
	EProtoId_PID_GET_IN_CARD              EProtoId = 17
	EProtoId_PID_SEND_OUT_CARD            EProtoId = 18
	EProtoId_PID_SEND_OUT_CARD_ACK        EProtoId = 19
	EProtoId_PID_PENG_CARD                EProtoId = 20
	EProtoId_PID_PENG_CARD_ACK            EProtoId = 21
	EProtoId_PID_GANG_CARD                EProtoId = 22
	EProtoId_PID_GANG_CARD_ACK            EProtoId = 23
	EProtoId_PID_GUO_CARD                 EProtoId = 24
	EProtoId_PID_GUO_CARD_ACK             EProtoId = 25
	EProtoId_PID_HU_CARD                  EProtoId = 26
	EProtoId_PID_HU_CARD_ACK              EProtoId = 27
	EProtoId_PID_BROADCAST_BEGIN_DINGQUE  EProtoId = 28
	EProtoId_PID_BROADCAST_BEGIN_EXCHANGE EProtoId = 29
	EProtoId_PID_OVERTURN                 EProtoId = 30
	EProtoId_PID_CURRENTRESULT            EProtoId = 31
	EProtoId_PID_SENDENDLOTTERY           EProtoId = 32
	EProtoId_PID_DISSOLVE_DESK            EProtoId = 33
	EProtoId_PID_DISSOLVE_DESK_ACK        EProtoId = 34
	EProtoId_PID_LEAVE_DESK               EProtoId = 35
	EProtoId_PID_LEAVE_DESK_ACK           EProtoId = 36
	EProtoId_PID_MESSAGE                  EProtoId = 37
	EProtoId_PID_SEND_MESSAGE             EProtoId = 38
)

var EProtoId_name = map[int32]string{
	1:  "PID_QUICK_CONN",
	2:  "PID_QUICK_CONN_ACK",
	3:  "PID_GAME_LOGIN",
	4:  "PID_GAME_LOGIN_ACK",
	5:  "PID_CREATEROOM",
	6:  "PID_CREATEROOM_ACK",
	7:  "PID_ENTER_ROOM",
	8:  "PID_ENTER_ROOM_ACK",
	9:  "PID_SEND_GAMEINFO",
	10: "PID_READY",
	11: "PID_READY_ACK",
	12: "PID_EXCHANGECARDS",
	13: "PID_EXCHANGECARDS_ACK",
	14: "PID_DINGQUE",
	15: "PID_OPENING",
	16: "PID_DEAL_CARDS",
	17: "PID_GET_IN_CARD",
	18: "PID_SEND_OUT_CARD",
	19: "PID_SEND_OUT_CARD_ACK",
	20: "PID_PENG_CARD",
	21: "PID_PENG_CARD_ACK",
	22: "PID_GANG_CARD",
	23: "PID_GANG_CARD_ACK",
	24: "PID_GUO_CARD",
	25: "PID_GUO_CARD_ACK",
	26: "PID_HU_CARD",
	27: "PID_HU_CARD_ACK",
	28: "PID_BROADCAST_BEGIN_DINGQUE",
	29: "PID_BROADCAST_BEGIN_EXCHANGE",
	30: "PID_OVERTURN",
	31: "PID_CURRENTRESULT",
	32: "PID_SENDENDLOTTERY",
	33: "PID_DISSOLVE_DESK",
	34: "PID_DISSOLVE_DESK_ACK",
	35: "PID_LEAVE_DESK",
	36: "PID_LEAVE_DESK_ACK",
	37: "PID_MESSAGE",
	38: "PID_SEND_MESSAGE",
}
var EProtoId_value = map[string]int32{
	"PID_QUICK_CONN":               1,
	"PID_QUICK_CONN_ACK":           2,
	"PID_GAME_LOGIN":               3,
	"PID_GAME_LOGIN_ACK":           4,
	"PID_CREATEROOM":               5,
	"PID_CREATEROOM_ACK":           6,
	"PID_ENTER_ROOM":               7,
	"PID_ENTER_ROOM_ACK":           8,
	"PID_SEND_GAMEINFO":            9,
	"PID_READY":                    10,
	"PID_READY_ACK":                11,
	"PID_EXCHANGECARDS":            12,
	"PID_EXCHANGECARDS_ACK":        13,
	"PID_DINGQUE":                  14,
	"PID_OPENING":                  15,
	"PID_DEAL_CARDS":               16,
	"PID_GET_IN_CARD":              17,
	"PID_SEND_OUT_CARD":            18,
	"PID_SEND_OUT_CARD_ACK":        19,
	"PID_PENG_CARD":                20,
	"PID_PENG_CARD_ACK":            21,
	"PID_GANG_CARD":                22,
	"PID_GANG_CARD_ACK":            23,
	"PID_GUO_CARD":                 24,
	"PID_GUO_CARD_ACK":             25,
	"PID_HU_CARD":                  26,
	"PID_HU_CARD_ACK":              27,
	"PID_BROADCAST_BEGIN_DINGQUE":  28,
	"PID_BROADCAST_BEGIN_EXCHANGE": 29,
	"PID_OVERTURN":                 30,
	"PID_CURRENTRESULT":            31,
	"PID_SENDENDLOTTERY":           32,
	"PID_DISSOLVE_DESK":            33,
	"PID_DISSOLVE_DESK_ACK":        34,
	"PID_LEAVE_DESK":               35,
	"PID_LEAVE_DESK_ACK":           36,
	"PID_MESSAGE":                  37,
	"PID_SEND_MESSAGE":             38,
}

func (x EProtoId) Enum() *EProtoId {
	p := new(EProtoId)
	*p = x
	return p
}
func (x EProtoId) String() string {
	return proto.EnumName(EProtoId_name, int32(x))
}
func (x *EProtoId) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(EProtoId_value, data, "EProtoId")
	if err != nil {
		return err
	}
	*x = EProtoId(value)
	return nil
}
func (EProtoId) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ErrorCode int32

const (
	ErrorCode_EC_SUCCESS ErrorCode = 0
	// -101   -200	游戏异常
	ErrorCode_EC_CREATE_DESK_DIAMOND_NOTENOUGH ErrorCode = -101
	ErrorCode_EC_CREATE_DESK_USER_NOTFOUND     ErrorCode = -102
	ErrorCode_EC_INTO_DESK_NOTFOUND            ErrorCode = -103
	ErrorCode_EC_GAME_READY_REPEAT             ErrorCode = -110
	ErrorCode_EC_GAME_READY_CHIP_NOT_ENOUGH    ErrorCode = -111
)

var ErrorCode_name = map[int32]string{
	0:    "EC_SUCCESS",
	-101: "EC_CREATE_DESK_DIAMOND_NOTENOUGH",
	-102: "EC_CREATE_DESK_USER_NOTFOUND",
	-103: "EC_INTO_DESK_NOTFOUND",
	-110: "EC_GAME_READY_REPEAT",
	-111: "EC_GAME_READY_CHIP_NOT_ENOUGH",
}
var ErrorCode_value = map[string]int32{
	"EC_SUCCESS":                       0,
	"EC_CREATE_DESK_DIAMOND_NOTENOUGH": -101,
	"EC_CREATE_DESK_USER_NOTFOUND":     -102,
	"EC_INTO_DESK_NOTFOUND":            -103,
	"EC_GAME_READY_REPEAT":             -110,
	"EC_GAME_READY_CHIP_NOT_ENOUGH":    -111,
}

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}
func (x ErrorCode) String() string {
	return proto.EnumName(ErrorCode_name, int32(x))
}
func (x *ErrorCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ErrorCode_value, data, "ErrorCode")
	if err != nil {
		return err
	}
	*x = ErrorCode(value)
	return nil
}
func (ErrorCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// 房间类型信息：包含房间类型和对应的局数、封顶、玩法等信息
// 房间类型枚举
type MJRoomType int32

const (
	MJRoomType_roomType_xueZhanDaoDi    MJRoomType = 0
	MJRoomType_roomType_sanRenLiangFang MJRoomType = 1
	MJRoomType_roomType_siRenLiangFang  MJRoomType = 2
	MJRoomType_roomType_deYangMaJiang   MJRoomType = 3
	MJRoomType_roomType_daoDaoHu        MJRoomType = 4
	MJRoomType_roomType_xueLiuChengHe   MJRoomType = 5
)

var MJRoomType_name = map[int32]string{
	0: "roomType_xueZhanDaoDi",
	1: "roomType_sanRenLiangFang",
	2: "roomType_siRenLiangFang",
	3: "roomType_deYangMaJiang",
	4: "roomType_daoDaoHu",
	5: "roomType_xueLiuChengHe",
}
var MJRoomType_value = map[string]int32{
	"roomType_xueZhanDaoDi":    0,
	"roomType_sanRenLiangFang": 1,
	"roomType_siRenLiangFang":  2,
	"roomType_deYangMaJiang":   3,
	"roomType_daoDaoHu":        4,
	"roomType_xueLiuChengHe":   5,
}

func (x MJRoomType) Enum() *MJRoomType {
	p := new(MJRoomType)
	*p = x
	return p
}
func (x MJRoomType) String() string {
	return proto.EnumName(MJRoomType_name, int32(x))
}
func (x *MJRoomType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MJRoomType_value, data, "MJRoomType")
	if err != nil {
		return err
	}
	*x = MJRoomType(value)
	return nil
}
func (MJRoomType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// 麻将花色
type MahjongColor int32

const (
	MahjongColor_WAN  MahjongColor = 1
	MahjongColor_TIAO MahjongColor = 2
	MahjongColor_TONG MahjongColor = 3
)

var MahjongColor_name = map[int32]string{
	1: "WAN",
	2: "TIAO",
	3: "TONG",
}
var MahjongColor_value = map[string]int32{
	"WAN":  1,
	"TIAO": 2,
	"TONG": 3,
}

func (x MahjongColor) Enum() *MahjongColor {
	p := new(MahjongColor)
	*p = x
	return p
}
func (x MahjongColor) String() string {
	return proto.EnumName(MahjongColor_name, int32(x))
}
func (x *MahjongColor) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MahjongColor_value, data, "MahjongColor")
	if err != nil {
		return err
	}
	*x = MahjongColor(value)
	return nil
}
func (MahjongColor) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// ProtoHeader 需要在每个 Message 中作为第一个字段
type ProtoHeader struct {
	Version          *string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
	UserId           *uint32 `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	Code             *int32  `protobuf:"varint,3,opt,name=code" json:"code,omitempty"`
	Error            *string `protobuf:"bytes,4,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ProtoHeader) Reset()                    { *m = ProtoHeader{} }
func (m *ProtoHeader) String() string            { return proto.CompactTextString(m) }
func (*ProtoHeader) ProtoMessage()               {}
func (*ProtoHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ProtoHeader) GetVersion() string {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return ""
}

func (m *ProtoHeader) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *ProtoHeader) GetCode() int32 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *ProtoHeader) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

// 微信信息
type WeixinInfo struct {
	OpenId           *string `protobuf:"bytes,1,opt,name=openId" json:"openId,omitempty"`
	NickName         *string `protobuf:"bytes,2,opt,name=nickName" json:"nickName,omitempty"`
	HeadUrl          *string `protobuf:"bytes,3,opt,name=headUrl" json:"headUrl,omitempty"`
	Sex              *int32  `protobuf:"varint,4,opt,name=sex" json:"sex,omitempty"`
	City             *string `protobuf:"bytes,5,opt,name=city" json:"city,omitempty"`
	UnionId          *string `protobuf:"bytes,6,opt,name=unionId" json:"unionId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *WeixinInfo) Reset()                    { *m = WeixinInfo{} }
func (m *WeixinInfo) String() string            { return proto.CompactTextString(m) }
func (*WeixinInfo) ProtoMessage()               {}
func (*WeixinInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *WeixinInfo) GetOpenId() string {
	if m != nil && m.OpenId != nil {
		return *m.OpenId
	}
	return ""
}

func (m *WeixinInfo) GetNickName() string {
	if m != nil && m.NickName != nil {
		return *m.NickName
	}
	return ""
}

func (m *WeixinInfo) GetHeadUrl() string {
	if m != nil && m.HeadUrl != nil {
		return *m.HeadUrl
	}
	return ""
}

func (m *WeixinInfo) GetSex() int32 {
	if m != nil && m.Sex != nil {
		return *m.Sex
	}
	return 0
}

func (m *WeixinInfo) GetCity() string {
	if m != nil && m.City != nil {
		return *m.City
	}
	return ""
}

func (m *WeixinInfo) GetUnionId() string {
	if m != nil && m.UnionId != nil {
		return *m.UnionId
	}
	return ""
}

// 麻将牌
type CardInfo struct {
	Type             *int32 `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	Value            *int32 `protobuf:"varint,2,opt,name=value" json:"value,omitempty"`
	Id               *int32 `protobuf:"varint,3,opt,name=id" json:"id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *CardInfo) Reset()                    { *m = CardInfo{} }
func (m *CardInfo) String() string            { return proto.CompactTextString(m) }
func (*CardInfo) ProtoMessage()               {}
func (*CardInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CardInfo) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *CardInfo) GetValue() int32 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

func (m *CardInfo) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

// 玩法：包括自摸、点炮、以及可多选的玩法
type PlayOptions struct {
	ZiMoRadio        *int32  `protobuf:"varint,1,opt,name=ziMoRadio" json:"ziMoRadio,omitempty"`
	DianGangHuaRadio *int32  `protobuf:"varint,2,opt,name=dianGangHuaRadio" json:"dianGangHuaRadio,omitempty"`
	OthersCheckBox   []int32 `protobuf:"varint,3,rep,name=othersCheckBox" json:"othersCheckBox,omitempty"`
	HuRadio          *int32  `protobuf:"varint,4,opt,name=huRadio" json:"huRadio,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PlayOptions) Reset()                    { *m = PlayOptions{} }
func (m *PlayOptions) String() string            { return proto.CompactTextString(m) }
func (*PlayOptions) ProtoMessage()               {}
func (*PlayOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PlayOptions) GetZiMoRadio() int32 {
	if m != nil && m.ZiMoRadio != nil {
		return *m.ZiMoRadio
	}
	return 0
}

func (m *PlayOptions) GetDianGangHuaRadio() int32 {
	if m != nil && m.DianGangHuaRadio != nil {
		return *m.DianGangHuaRadio
	}
	return 0
}

func (m *PlayOptions) GetOthersCheckBox() []int32 {
	if m != nil {
		return m.OthersCheckBox
	}
	return nil
}

func (m *PlayOptions) GetHuRadio() int32 {
	if m != nil && m.HuRadio != nil {
		return *m.HuRadio
	}
	return 0
}

type RoomTypeInfo struct {
	MjRoomType       *MJRoomType  `protobuf:"varint,1,opt,name=mjRoomType,enum=mjproto.MJRoomType" json:"mjRoomType,omitempty"`
	BoardsCout       *int32       `protobuf:"varint,2,opt,name=boardsCout" json:"boardsCout,omitempty"`
	CapMax           *int64       `protobuf:"varint,3,opt,name=capMax" json:"capMax,omitempty"`
	PlayOptions      *PlayOptions `protobuf:"bytes,4,opt,name=playOptions" json:"playOptions,omitempty"`
	CardsNum         *int32       `protobuf:"varint,5,opt,name=cardsNum" json:"cardsNum,omitempty"`
	Settlement       *int32       `protobuf:"varint,6,opt,name=settlement" json:"settlement,omitempty"`
	BaseValue        *int64       `protobuf:"varint,7,opt,name=baseValue" json:"baseValue,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *RoomTypeInfo) Reset()                    { *m = RoomTypeInfo{} }
func (m *RoomTypeInfo) String() string            { return proto.CompactTextString(m) }
func (*RoomTypeInfo) ProtoMessage()               {}
func (*RoomTypeInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RoomTypeInfo) GetMjRoomType() MJRoomType {
	if m != nil && m.MjRoomType != nil {
		return *m.MjRoomType
	}
	return MJRoomType_roomType_xueZhanDaoDi
}

func (m *RoomTypeInfo) GetBoardsCout() int32 {
	if m != nil && m.BoardsCout != nil {
		return *m.BoardsCout
	}
	return 0
}

func (m *RoomTypeInfo) GetCapMax() int64 {
	if m != nil && m.CapMax != nil {
		return *m.CapMax
	}
	return 0
}

func (m *RoomTypeInfo) GetPlayOptions() *PlayOptions {
	if m != nil {
		return m.PlayOptions
	}
	return nil
}

func (m *RoomTypeInfo) GetCardsNum() int32 {
	if m != nil && m.CardsNum != nil {
		return *m.CardsNum
	}
	return 0
}

func (m *RoomTypeInfo) GetSettlement() int32 {
	if m != nil && m.Settlement != nil {
		return *m.Settlement
	}
	return 0
}

func (m *RoomTypeInfo) GetBaseValue() int64 {
	if m != nil && m.BaseValue != nil {
		return *m.BaseValue
	}
	return 0
}

type ComposeCard struct {
	Value            *int32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
	Type             *int32 `protobuf:"varint,2,opt,name=type" json:"type,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ComposeCard) Reset()                    { *m = ComposeCard{} }
func (m *ComposeCard) String() string            { return proto.CompactTextString(m) }
func (*ComposeCard) ProtoMessage()               {}
func (*ComposeCard) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ComposeCard) GetValue() int32 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

func (m *ComposeCard) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

type PlayerCard struct {
	HandCard         []*CardInfo    `protobuf:"bytes,1,rep,name=handCard" json:"handCard,omitempty"`
	ComposeCard      []*ComposeCard `protobuf:"bytes,2,rep,name=composeCard" json:"composeCard,omitempty"`
	OutCard          []int32        `protobuf:"varint,3,rep,name=outCard" json:"outCard,omitempty"`
	HuCard           *int32         `protobuf:"varint,4,opt,name=huCard" json:"huCard,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *PlayerCard) Reset()                    { *m = PlayerCard{} }
func (m *PlayerCard) String() string            { return proto.CompactTextString(m) }
func (*PlayerCard) ProtoMessage()               {}
func (*PlayerCard) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *PlayerCard) GetHandCard() []*CardInfo {
	if m != nil {
		return m.HandCard
	}
	return nil
}

func (m *PlayerCard) GetComposeCard() []*ComposeCard {
	if m != nil {
		return m.ComposeCard
	}
	return nil
}

func (m *PlayerCard) GetOutCard() []int32 {
	if m != nil {
		return m.OutCard
	}
	return nil
}

func (m *PlayerCard) GetHuCard() int32 {
	if m != nil && m.HuCard != nil {
		return *m.HuCard
	}
	return 0
}

type PlayerInfo struct {
	IsBanker         *bool       `protobuf:"varint,1,opt,name=isBanker" json:"isBanker,omitempty"`
	PlayerCard       *PlayerCard `protobuf:"bytes,2,opt,name=playerCard" json:"playerCard,omitempty"`
	Coin             *int64      `protobuf:"varint,3,opt,name=coin" json:"coin,omitempty"`
	NickName         *string     `protobuf:"bytes,4,opt,name=nickName" json:"nickName,omitempty"`
	Sex              *int32      `protobuf:"varint,5,opt,name=sex" json:"sex,omitempty"`
	UserId           *uint32     `protobuf:"varint,6,opt,name=userId" json:"userId,omitempty"`
	IsOwner          *bool       `protobuf:"varint,7,opt,name=isOwner" json:"isOwner,omitempty"`
	BReady           *int32      `protobuf:"varint,8,opt,name=bReady" json:"bReady,omitempty"`
	BDingQue         *int32      `protobuf:"varint,9,opt,name=bDingQue" json:"bDingQue,omitempty"`
	BExchanged       *int32      `protobuf:"varint,10,opt,name=bExchanged" json:"bExchanged,omitempty"`
	NHuPai           *int32      `protobuf:"varint,11,opt,name=nHuPai" json:"nHuPai,omitempty"`
	QuePai           *int32      `protobuf:"varint,12,opt,name=quePai" json:"quePai,omitempty"`
	WxInfo           *WeixinInfo `protobuf:"bytes,13,opt,name=wxInfo" json:"wxInfo,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *PlayerInfo) Reset()                    { *m = PlayerInfo{} }
func (m *PlayerInfo) String() string            { return proto.CompactTextString(m) }
func (*PlayerInfo) ProtoMessage()               {}
func (*PlayerInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *PlayerInfo) GetIsBanker() bool {
	if m != nil && m.IsBanker != nil {
		return *m.IsBanker
	}
	return false
}

func (m *PlayerInfo) GetPlayerCard() *PlayerCard {
	if m != nil {
		return m.PlayerCard
	}
	return nil
}

func (m *PlayerInfo) GetCoin() int64 {
	if m != nil && m.Coin != nil {
		return *m.Coin
	}
	return 0
}

func (m *PlayerInfo) GetNickName() string {
	if m != nil && m.NickName != nil {
		return *m.NickName
	}
	return ""
}

func (m *PlayerInfo) GetSex() int32 {
	if m != nil && m.Sex != nil {
		return *m.Sex
	}
	return 0
}

func (m *PlayerInfo) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *PlayerInfo) GetIsOwner() bool {
	if m != nil && m.IsOwner != nil {
		return *m.IsOwner
	}
	return false
}

func (m *PlayerInfo) GetBReady() int32 {
	if m != nil && m.BReady != nil {
		return *m.BReady
	}
	return 0
}

func (m *PlayerInfo) GetBDingQue() int32 {
	if m != nil && m.BDingQue != nil {
		return *m.BDingQue
	}
	return 0
}

func (m *PlayerInfo) GetBExchanged() int32 {
	if m != nil && m.BExchanged != nil {
		return *m.BExchanged
	}
	return 0
}

func (m *PlayerInfo) GetNHuPai() int32 {
	if m != nil && m.NHuPai != nil {
		return *m.NHuPai
	}
	return 0
}

func (m *PlayerInfo) GetQuePai() int32 {
	if m != nil && m.QuePai != nil {
		return *m.QuePai
	}
	return 0
}

func (m *PlayerInfo) GetWxInfo() *WeixinInfo {
	if m != nil {
		return m.WxInfo
	}
	return nil
}

type DeskGameInfo struct {
	GameStatus       *int32        `protobuf:"varint,1,opt,name=GameStatus" json:"GameStatus,omitempty"`
	RoomTypeInfo     *RoomTypeInfo `protobuf:"bytes,2,opt,name=roomTypeInfo" json:"roomTypeInfo,omitempty"`
	PlayerNum        *int32        `protobuf:"varint,3,opt,name=playerNum" json:"playerNum,omitempty"`
	ActiveUserId     *uint32       `protobuf:"varint,4,opt,name=activeUserId" json:"activeUserId,omitempty"`
	ActionTime       *int32        `protobuf:"varint,5,opt,name=actionTime" json:"actionTime,omitempty"`
	DelayTime        *int32        `protobuf:"varint,6,opt,name=delayTime" json:"delayTime,omitempty"`
	NInitActionTime  *int32        `protobuf:"varint,7,opt,name=nInitActionTime" json:"nInitActionTime,omitempty"`
	NInitDelayTime   *int32        `protobuf:"varint,8,opt,name=nInitDelayTime" json:"nInitDelayTime,omitempty"`
	InitRoomCoin     *int64        `protobuf:"varint,9,opt,name=initRoomCoin" json:"initRoomCoin,omitempty"`
	CurrPlayCount    *int32        `protobuf:"varint,10,opt,name=currPlayCount" json:"currPlayCount,omitempty"`
	TotalPlayCount   *int32        `protobuf:"varint,11,opt,name=totalPlayCount" json:"totalPlayCount,omitempty"`
	RoomNumber       *string       `protobuf:"bytes,12,opt,name=roomNumber" json:"roomNumber,omitempty"`
	RemainCards      *int32        `protobuf:"varint,13,opt,name=remainCards" json:"remainCards,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *DeskGameInfo) Reset()                    { *m = DeskGameInfo{} }
func (m *DeskGameInfo) String() string            { return proto.CompactTextString(m) }
func (*DeskGameInfo) ProtoMessage()               {}
func (*DeskGameInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *DeskGameInfo) GetGameStatus() int32 {
	if m != nil && m.GameStatus != nil {
		return *m.GameStatus
	}
	return 0
}

func (m *DeskGameInfo) GetRoomTypeInfo() *RoomTypeInfo {
	if m != nil {
		return m.RoomTypeInfo
	}
	return nil
}

func (m *DeskGameInfo) GetPlayerNum() int32 {
	if m != nil && m.PlayerNum != nil {
		return *m.PlayerNum
	}
	return 0
}

func (m *DeskGameInfo) GetActiveUserId() uint32 {
	if m != nil && m.ActiveUserId != nil {
		return *m.ActiveUserId
	}
	return 0
}

func (m *DeskGameInfo) GetActionTime() int32 {
	if m != nil && m.ActionTime != nil {
		return *m.ActionTime
	}
	return 0
}

func (m *DeskGameInfo) GetDelayTime() int32 {
	if m != nil && m.DelayTime != nil {
		return *m.DelayTime
	}
	return 0
}

func (m *DeskGameInfo) GetNInitActionTime() int32 {
	if m != nil && m.NInitActionTime != nil {
		return *m.NInitActionTime
	}
	return 0
}

func (m *DeskGameInfo) GetNInitDelayTime() int32 {
	if m != nil && m.NInitDelayTime != nil {
		return *m.NInitDelayTime
	}
	return 0
}

func (m *DeskGameInfo) GetInitRoomCoin() int64 {
	if m != nil && m.InitRoomCoin != nil {
		return *m.InitRoomCoin
	}
	return 0
}

func (m *DeskGameInfo) GetCurrPlayCount() int32 {
	if m != nil && m.CurrPlayCount != nil {
		return *m.CurrPlayCount
	}
	return 0
}

func (m *DeskGameInfo) GetTotalPlayCount() int32 {
	if m != nil && m.TotalPlayCount != nil {
		return *m.TotalPlayCount
	}
	return 0
}

func (m *DeskGameInfo) GetRoomNumber() string {
	if m != nil && m.RoomNumber != nil {
		return *m.RoomNumber
	}
	return ""
}

func (m *DeskGameInfo) GetRemainCards() int32 {
	if m != nil && m.RemainCards != nil {
		return *m.RemainCards
	}
	return 0
}

func init() {
	proto.RegisterType((*ProtoHeader)(nil), "mjproto.ProtoHeader")
	proto.RegisterType((*WeixinInfo)(nil), "mjproto.WeixinInfo")
	proto.RegisterType((*CardInfo)(nil), "mjproto.CardInfo")
	proto.RegisterType((*PlayOptions)(nil), "mjproto.PlayOptions")
	proto.RegisterType((*RoomTypeInfo)(nil), "mjproto.RoomTypeInfo")
	proto.RegisterType((*ComposeCard)(nil), "mjproto.ComposeCard")
	proto.RegisterType((*PlayerCard)(nil), "mjproto.PlayerCard")
	proto.RegisterType((*PlayerInfo)(nil), "mjproto.PlayerInfo")
	proto.RegisterType((*DeskGameInfo)(nil), "mjproto.DeskGameInfo")
	proto.RegisterEnum("mjproto.EProtoId", EProtoId_name, EProtoId_value)
	proto.RegisterEnum("mjproto.ErrorCode", ErrorCode_name, ErrorCode_value)
	proto.RegisterEnum("mjproto.MJRoomType", MJRoomType_name, MJRoomType_value)
	proto.RegisterEnum("mjproto.MahjongColor", MahjongColor_name, MahjongColor_value)
}

var fileDescriptor0 = []byte{
	// 1336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x56, 0xed, 0x76, 0xdb, 0x44,
	0x10, 0xc5, 0x71, 0x1c, 0xdb, 0x63, 0x27, 0xd9, 0xa8, 0x49, 0xeb, 0xb6, 0x29, 0x4d, 0x5d, 0x3e,
	0xc3, 0xa1, 0x3f, 0xca, 0x13, 0x28, 0x92, 0x6a, 0xab, 0xb1, 0xa5, 0x54, 0xb6, 0x5a, 0xca, 0x1f,
	0x1f, 0xc5, 0x5a, 0x12, 0x35, 0xb1, 0x64, 0xf4, 0xd1, 0x3a, 0x9c, 0xc3, 0x3b, 0x00, 0xff, 0x80,
	0xa7, 0xe0, 0x05, 0xe0, 0x09, 0x78, 0x0f, 0xde, 0x02, 0x66, 0x47, 0x5e, 0xc9, 0x2e, 0xf4, 0x34,
	0xe7, 0x48, 0x77, 0xef, 0xce, 0x9d, 0xb9, 0x33, 0xbb, 0x16, 0xc0, 0xb9, 0x97, 0xf0, 0x27, 0xf3,
	0x38, 0x4a, 0x23, 0xa5, 0x3e, 0x7b, 0x43, 0x0f, 0xdd, 0x53, 0x68, 0x9d, 0x89, 0x87, 0x3e, 0xf7,
	0x7c, 0x1e, 0x2b, 0xbb, 0x50, 0x7f, 0xcb, 0xe3, 0x24, 0x88, 0xc2, 0x4e, 0xe5, 0xa8, 0xf2, 0x59,
	0x53, 0xd9, 0x81, 0xad, 0x2c, 0xe1, 0xb1, 0xe9, 0x77, 0x36, 0xf0, 0x7d, 0x5b, 0x69, 0xc3, 0xe6,
	0x34, 0xf2, 0x79, 0xa7, 0x8a, 0x6f, 0x35, 0x65, 0x1b, 0x6a, 0x3c, 0x8e, 0xa3, 0xb8, 0xb3, 0x29,
	0xc8, 0xdd, 0x2b, 0x80, 0x57, 0x3c, 0x58, 0x04, 0xa1, 0x19, 0x7e, 0x1b, 0x89, 0xad, 0xd1, 0x9c,
	0x87, 0xb8, 0x35, 0x0f, 0xc5, 0xa0, 0x11, 0x06, 0xd3, 0x2b, 0xcb, 0x9b, 0x71, 0x0a, 0xd6, 0x14,
	0x6a, 0x97, 0xa8, 0xeb, 0xc6, 0xd7, 0x14, 0xaf, 0xa9, 0xb4, 0xa0, 0x9a, 0xf0, 0x05, 0x45, 0xab,
	0x91, 0x54, 0x90, 0xde, 0x74, 0x6a, 0x92, 0x9b, 0x85, 0x98, 0x17, 0x86, 0xdb, 0x22, 0xb1, 0xaf,
	0xa0, 0xa1, 0x79, 0xb1, 0x4f, 0x52, 0x48, 0x4d, 0x6f, 0xe6, 0x9c, 0x84, 0x28, 0xab, 0xb7, 0xde,
	0x75, 0x96, 0xab, 0xd4, 0x14, 0x80, 0x8d, 0xc0, 0xcf, 0x13, 0xee, 0x4e, 0xb1, 0xdc, 0x6b, 0xef,
	0xc6, 0x9e, 0xa7, 0x18, 0x2a, 0x51, 0xf6, 0xa0, 0xf9, 0x7d, 0x30, 0x8c, 0x1c, 0xcf, 0x0f, 0xa2,
	0xe5, 0xe6, 0x0e, 0x30, 0x3f, 0xf0, 0xc2, 0x9e, 0x17, 0x5e, 0xf4, 0x33, 0x2f, 0x5f, 0xc9, 0xe3,
	0xdc, 0x86, 0x9d, 0x28, 0xbd, 0x44, 0x77, 0xb4, 0x4b, 0x3e, 0xbd, 0x3a, 0x89, 0x16, 0x18, 0xb3,
	0x8a, 0xb8, 0xa8, 0x22, 0xcb, 0x89, 0x94, 0x78, 0xf7, 0xcf, 0x0a, 0xb4, 0x9d, 0x28, 0x9a, 0x8d,
	0x31, 0x25, 0x4a, 0xef, 0x53, 0x80, 0xd9, 0x1b, 0x89, 0x90, 0xce, 0xce, 0xd3, 0x5b, 0x4f, 0x96,
	0x2d, 0x78, 0x32, 0x7c, 0x2e, 0x97, 0x14, 0xcc, 0xf5, 0x3c, 0xc2, 0xa2, 0x12, 0x2d, 0xca, 0xd2,
	0xa5, 0x2c, 0xda, 0x38, 0xf5, 0xe6, 0x43, 0x6f, 0x41, 0x25, 0x54, 0x95, 0xcf, 0xa1, 0x35, 0x2f,
	0x4b, 0x20, 0xc9, 0xd6, 0xd3, 0xfd, 0x22, 0xda, 0x6a, 0x79, 0xe8, 0xf8, 0x54, 0x44, 0xb3, 0xb2,
	0x19, 0xb9, 0x58, 0x13, 0x02, 0x09, 0x4f, 0xd3, 0x6b, 0x3e, 0xe3, 0x61, 0x4a, 0x46, 0xd6, 0x84,
	0x09, 0x62, 0x32, 0x5e, 0x92, 0x65, 0x75, 0xa1, 0xd1, 0x3d, 0x86, 0x96, 0x16, 0xcd, 0xe6, 0x51,
	0xc2, 0x85, 0xc5, 0xa5, 0xa1, 0x15, 0xd9, 0x18, 0x72, 0x9b, 0xf2, 0xeb, 0xfe, 0x00, 0x20, 0x34,
	0x79, 0x4c, 0xd4, 0xc7, 0xd0, 0xb8, 0xf4, 0x42, 0x5f, 0x3c, 0x23, 0xbb, 0x8a, 0xa9, 0xed, 0x15,
	0xa9, 0x15, 0xed, 0xc2, 0x12, 0xa6, 0x65, 0x78, 0x8c, 0x53, 0x5d, 0x2b, 0x61, 0x55, 0x1a, 0xcd,
	0x45, 0x2b, 0x88, 0x96, 0xbb, 0x8d, 0x76, 0x5c, 0x66, 0xf4, 0x9e, 0x9b, 0xfd, 0xe3, 0x86, 0xd4,
	0xa7, 0xd0, 0x58, 0x72, 0x90, 0x9c, 0x78, 0xe1, 0x15, 0x8f, 0x29, 0xdb, 0x86, 0x30, 0x7f, 0x5e,
	0xe4, 0x47, 0x39, 0xb7, 0x56, 0xcc, 0x5f, 0x49, 0x9d, 0x46, 0x3b, 0x08, 0x97, 0x36, 0xaf, 0x4e,
	0xeb, 0xe6, 0xea, 0x70, 0xd6, 0x64, 0x57, 0x96, 0xe7, 0x62, 0x8b, 0xce, 0x05, 0xe6, 0x19, 0x24,
	0xf6, 0xbb, 0x10, 0x65, 0xeb, 0x24, 0x8b, 0x84, 0x73, 0x07, 0x87, 0xfb, 0xa6, 0xd3, 0xa0, 0x0d,
	0x18, 0xef, 0x5c, 0x0f, 0xc2, 0x8b, 0x17, 0x68, 0x63, 0x53, 0xf6, 0xe2, 0xdc, 0x58, 0x4c, 0xd1,
	0xae, 0x0b, 0xee, 0x77, 0x40, 0x86, 0x0d, 0xfb, 0xd9, 0x99, 0x17, 0x74, 0x5a, 0xf2, 0xfd, 0xbb,
	0x8c, 0x8b, 0xf7, 0x36, 0xbd, 0x3f, 0x86, 0xad, 0x77, 0x0b, 0x51, 0x68, 0x67, 0xfb, 0xbd, 0x42,
	0xca, 0x83, 0xd7, 0xfd, 0x63, 0x03, 0xda, 0x3a, 0x4f, 0xae, 0x7a, 0x98, 0x3b, 0x99, 0x82, 0x4a,
	0xe2, 0x79, 0x94, 0x7a, 0x69, 0x96, 0x2c, 0x9b, 0xf8, 0x05, 0xb4, 0xe3, 0x95, 0x19, 0x5d, 0x1a,
	0x73, 0x50, 0xc4, 0x5b, 0x1b, 0x60, 0x1c, 0x91, 0xdc, 0x43, 0x31, 0x49, 0xf9, 0xd1, 0xdf, 0x87,
	0xb6, 0x37, 0x4d, 0x83, 0xb7, 0xdc, 0xcd, 0x6d, 0xd8, 0x24, 0x1b, 0x50, 0x49, 0xa0, 0x51, 0x38,
	0x0e, 0xd0, 0xb7, 0x9a, 0x9c, 0x2f, 0x9f, 0xe3, 0x6e, 0x82, 0xf2, 0x91, 0xbb, 0x03, 0xbb, 0x98,
	0x6a, 0x90, 0xaa, 0x25, 0xb7, 0x2e, 0xcf, 0x18, 0x2d, 0xe8, 0xc5, 0x86, 0x86, 0x54, 0x0b, 0x10,
	0x16, 0x49, 0x69, 0xa2, 0x47, 0x4d, 0xea, 0xd1, 0x01, 0x6c, 0x4f, 0xb3, 0x38, 0x16, 0x3d, 0xc4,
	0x03, 0x83, 0x03, 0x0d, 0x32, 0x48, 0x1a, 0xa5, 0xde, 0x75, 0x89, 0xb7, 0xa4, 0xe1, 0xa2, 0x64,
	0xac, 0xe1, 0x1c, 0xdb, 0xd4, 0xa6, 0xa6, 0xde, 0x82, 0x56, 0xcc, 0x67, 0x5e, 0x10, 0x8a, 0x11,
	0x48, 0xc8, 0xd5, 0xda, 0xf1, 0x5f, 0x5b, 0xd0, 0x30, 0xe8, 0x5a, 0x34, 0x7d, 0xdc, 0xb5, 0x73,
	0x66, 0xea, 0x93, 0x17, 0xae, 0xa9, 0x9d, 0x4e, 0x34, 0xdb, 0xb2, 0x58, 0x05, 0x15, 0x94, 0x75,
	0x6c, 0xa2, 0x6a, 0xa7, 0x6c, 0x43, 0x72, 0x7b, 0xea, 0xd0, 0x98, 0x0c, 0xec, 0x9e, 0x69, 0xb1,
	0xaa, 0xe4, 0x96, 0x18, 0x71, 0x37, 0x25, 0x57, 0x73, 0x0c, 0x75, 0x6c, 0x38, 0xb6, 0x3d, 0x64,
	0x35, 0xc9, 0x2d, 0x31, 0xe2, 0x6e, 0x49, 0xae, 0x61, 0x21, 0x3a, 0x21, 0x6e, 0x5d, 0x72, 0x4b,
	0x8c, 0xb8, 0x0d, 0x34, 0x65, 0x4f, 0xe0, 0x23, 0xc3, 0xca, 0x45, 0x4d, 0xeb, 0x99, 0xcd, 0x9a,
	0x78, 0x86, 0x9b, 0x02, 0xc6, 0xc8, 0xfa, 0x6b, 0x06, 0xd8, 0x94, 0xed, 0xe2, 0x95, 0x36, 0xb6,
	0xe4, 0x46, 0xe3, 0x6b, 0xad, 0xaf, 0x5a, 0x3d, 0x43, 0x53, 0x1d, 0x7d, 0xc4, 0xda, 0xca, 0x5d,
	0x38, 0xf8, 0x0f, 0x4c, 0x3b, 0xc4, 0xd0, 0xb7, 0xc4, 0x92, 0x6e, 0x5a, 0xbd, 0x17, 0xae, 0xc1,
	0x76, 0x24, 0x60, 0x9f, 0x19, 0x16, 0x62, 0x6c, 0x57, 0x26, 0xae, 0x1b, 0xea, 0x60, 0x92, 0x07,
	0x64, 0x68, 0xf9, 0x2e, 0x19, 0x62, 0x8c, 0x27, 0x68, 0x86, 0x40, 0xd9, 0xde, 0x5a, 0xd6, 0xb6,
	0x3b, 0xce, 0x61, 0x45, 0x8a, 0xaf, 0xc1, 0x24, 0x7e, 0x4b, 0x56, 0x80, 0x52, 0xbd, 0x9c, 0xbd,
	0x2f, 0x83, 0x14, 0x10, 0x31, 0x0f, 0x24, 0xb3, 0xa7, 0x4a, 0xe6, 0x6d, 0xc9, 0x2c, 0x20, 0x62,
	0xde, 0xc1, 0x43, 0xda, 0x26, 0xd8, 0xb5, 0x73, 0x62, 0x07, 0x07, 0x8f, 0xad, 0x22, 0xc4, 0xbb,
	0x2b, 0xeb, 0xec, 0xbb, 0x39, 0xed, 0x9e, 0xac, 0x69, 0x09, 0x10, 0xeb, 0xbe, 0xf2, 0x10, 0xee,
	0x0b, 0xf0, 0xc4, 0xb1, 0x55, 0x5d, 0x53, 0x47, 0xe3, 0xc9, 0x89, 0x21, 0xda, 0x2f, 0xed, 0x3a,
	0x54, 0x8e, 0xe0, 0xf0, 0xff, 0x08, 0xd2, 0x6a, 0xf6, 0x40, 0x26, 0x64, 0xbf, 0x34, 0x9c, 0xb1,
	0xeb, 0x58, 0xec, 0x43, 0x99, 0xb9, 0xe6, 0x3a, 0x0e, 0xb6, 0xde, 0x31, 0x46, 0xee, 0x60, 0xcc,
	0x1e, 0xca, 0x69, 0x10, 0x46, 0xe1, 0xff, 0x81, 0x3d, 0xc6, 0xa9, 0x78, 0xcd, 0x8e, 0x24, 0x5d,
	0x37, 0x47, 0x23, 0x7b, 0xf0, 0xd2, 0xc0, 0x4e, 0x8c, 0x4e, 0xd9, 0x23, 0xe9, 0xeb, 0x1a, 0x4c,
	0x59, 0x77, 0x65, 0xcb, 0x06, 0x86, 0x2a, 0xe9, 0x8f, 0x65, 0xf4, 0x12, 0x23, 0xee, 0x47, 0xd2,
	0x87, 0xa1, 0x31, 0x1a, 0xa9, 0x98, 0xef, 0xc7, 0xd2, 0x2e, 0xea, 0x97, 0x44, 0x3f, 0x39, 0xfe,
	0xbb, 0x02, 0x4d, 0x43, 0x7c, 0x27, 0x68, 0xf8, 0xe9, 0x80, 0x77, 0x1a, 0x18, 0xda, 0x64, 0xe4,
	0x6a, 0x1a, 0x12, 0xd8, 0x07, 0xca, 0x97, 0x70, 0x84, 0xef, 0xf9, 0xcc, 0xe7, 0xc1, 0x75, 0x53,
	0x1d, 0xda, 0x18, 0xc1, 0xb2, 0xc7, 0x86, 0x65, 0xbb, 0xbd, 0x3e, 0xfb, 0xed, 0x9f, 0xe5, 0xbf,
	0x0a, 0xfe, 0x78, 0x1c, 0xbe, 0x47, 0x77, 0x47, 0x78, 0x02, 0x90, 0xfb, 0xcc, 0x76, 0x2d, 0x9d,
	0xfd, 0x5a, 0x52, 0xbb, 0x70, 0x80, 0x54, 0xd3, 0x1a, 0xdb, 0x39, 0xb1, 0xe0, 0xfc, 0x52, 0x72,
	0x1e, 0xc1, 0x3e, 0x72, 0xe8, 0x74, 0xe6, 0x87, 0xc1, 0x31, 0xce, 0x30, 0x34, 0xfb, 0xb9, 0xa4,
	0x1c, 0xc3, 0x83, 0x75, 0x8a, 0xd6, 0x37, 0xcf, 0x44, 0xb0, 0xc9, 0x32, 0xbb, 0x9f, 0x0a, 0xee,
	0xf1, 0xef, 0x15, 0x80, 0x95, 0x1f, 0x74, 0xf4, 0x59, 0xde, 0xb2, 0x93, 0x45, 0xc6, 0xbf, 0xc1,
	0xcb, 0x5e, 0xf7, 0x22, 0x3d, 0xc0, 0xb2, 0x0f, 0xa1, 0x53, 0x2c, 0x25, 0x5e, 0xe8, 0xf0, 0x70,
	0x80, 0xdf, 0x1d, 0x17, 0xcf, 0xf0, 0x0f, 0x6f, 0x98, 0xfb, 0x70, 0xa7, 0x5c, 0x0d, 0xd6, 0x16,
	0x37, 0x94, 0x7b, 0x70, 0xbb, 0x58, 0xf4, 0xf9, 0x6b, 0x04, 0x87, 0xde, 0x73, 0xb1, 0xce, 0xc4,
	0x9d, 0xb8, 0x57, 0xae, 0xa1, 0x94, 0x17, 0xf5, 0x33, 0xbc, 0x6d, 0x56, 0xb7, 0x60, 0x22, 0x83,
	0x20, 0xc3, 0xaf, 0x18, 0xfc, 0xc4, 0xe1, 0xac, 0x76, 0x8c, 0x3f, 0x05, 0x43, 0xef, 0xf2, 0x4d,
	0x14, 0x5e, 0x68, 0xd1, 0x75, 0x14, 0x2b, 0x75, 0xa8, 0xbe, 0x52, 0xc5, 0x35, 0xd7, 0x80, 0xcd,
	0xb1, 0xa9, 0xda, 0xa8, 0x28, 0x9e, 0x6c, 0x3c, 0xd1, 0xd5, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff,
	0xc6, 0xa8, 0x39, 0x28, 0x46, 0x0a, 0x00, 0x00,
}
