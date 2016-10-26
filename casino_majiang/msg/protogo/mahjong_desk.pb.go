// Code generated by protoc-gen-go.
// source: mahjong_desk.proto
// DO NOT EDIT!

package mjproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Ignoring public import of ProtoHeader from base.proto

// Ignoring public import of WeixinInfo from base.proto

// Ignoring public import of CardInfo from base.proto

// Ignoring public import of PlayOptions from base.proto

// Ignoring public import of RoomTypeInfo from base.proto

// Ignoring public import of ComposeCard from base.proto

// Ignoring public import of PlayerCard from base.proto

// Ignoring public import of PlayerInfo from base.proto

// Ignoring public import of DeskGameInfo from base.proto

// Ignoring public import of EProtoId from base.proto

// Ignoring public import of ErrorCode from base.proto

// Ignoring public import of MJOption from base.proto

// Ignoring public import of MJRoomType from base.proto

// Ignoring public import of MahjongColor from base.proto

// Ignoring public import of GangType from base.proto

// Ignoring public import of ComposeCardType from base.proto

// Ignoring public import of HuPaiType from base.proto

// Ignoring public import of MJUserGameStatus from base.proto

// Ignoring public import of DeskGameStatus from base.proto

// 房主解散房间(未开局)
type Game_DissolveDesk struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_DissolveDesk) Reset()                    { *m = Game_DissolveDesk{} }
func (m *Game_DissolveDesk) String() string            { return proto.CompactTextString(m) }
func (*Game_DissolveDesk) ProtoMessage()               {}
func (*Game_DissolveDesk) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Game_DissolveDesk) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_DissolveDesk) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 解散房间回复
type Game_AckDissolveDesk struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	DeskId           *int32       `protobuf:"varint,3,opt,name=deskId" json:"deskId,omitempty"`
	PassWord         *string      `protobuf:"bytes,4,opt,name=passWord" json:"passWord,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckDissolveDesk) Reset()                    { *m = Game_AckDissolveDesk{} }
func (m *Game_AckDissolveDesk) String() string            { return proto.CompactTextString(m) }
func (*Game_AckDissolveDesk) ProtoMessage()               {}
func (*Game_AckDissolveDesk) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *Game_AckDissolveDesk) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckDissolveDesk) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_AckDissolveDesk) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *Game_AckDissolveDesk) GetPassWord() string {
	if m != nil && m.PassWord != nil {
		return *m.PassWord
	}
	return ""
}

// 申请解散房间(游戏中)
type Game_ReqDissolveDesk struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_ReqDissolveDesk) Reset()                    { *m = Game_ReqDissolveDesk{} }
func (m *Game_ReqDissolveDesk) String() string            { return proto.CompactTextString(m) }
func (*Game_ReqDissolveDesk) ProtoMessage()               {}
func (*Game_ReqDissolveDesk) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *Game_ReqDissolveDesk) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_ReqDissolveDesk) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

type Game_AckReqDissolveDesk struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserIdAgree      []uint32     `protobuf:"varint,2,rep,name=userIdAgree" json:"userIdAgree,omitempty"`
	UserIdWait       []uint32     `protobuf:"varint,3,rep,name=userIdWait" json:"userIdWait,omitempty"`
	UserIdDisagree   []uint32     `protobuf:"varint,4,rep,name=userIdDisagree" json:"userIdDisagree,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckReqDissolveDesk) Reset()                    { *m = Game_AckReqDissolveDesk{} }
func (m *Game_AckReqDissolveDesk) String() string            { return proto.CompactTextString(m) }
func (*Game_AckReqDissolveDesk) ProtoMessage()               {}
func (*Game_AckReqDissolveDesk) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *Game_AckReqDissolveDesk) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckReqDissolveDesk) GetUserIdAgree() []uint32 {
	if m != nil {
		return m.UserIdAgree
	}
	return nil
}

func (m *Game_AckReqDissolveDesk) GetUserIdWait() []uint32 {
	if m != nil {
		return m.UserIdWait
	}
	return nil
}

func (m *Game_AckReqDissolveDesk) GetUserIdDisagree() []uint32 {
	if m != nil {
		return m.UserIdDisagree
	}
	return nil
}

// 离开房间
type Game_LeaveDesk struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_LeaveDesk) Reset()                    { *m = Game_LeaveDesk{} }
func (m *Game_LeaveDesk) String() string            { return proto.CompactTextString(m) }
func (*Game_LeaveDesk) ProtoMessage()               {}
func (*Game_LeaveDesk) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *Game_LeaveDesk) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

// 离开房间的回复
type Game_AckLeaveDesk struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckLeaveDesk) Reset()                    { *m = Game_AckLeaveDesk{} }
func (m *Game_AckLeaveDesk) String() string            { return proto.CompactTextString(m) }
func (*Game_AckLeaveDesk) ProtoMessage()               {}
func (*Game_AckLeaveDesk) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *Game_AckLeaveDesk) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

// 准备游戏
type Game_Ready struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_Ready) Reset()                    { *m = Game_Ready{} }
func (m *Game_Ready) String() string            { return proto.CompactTextString(m) }
func (*Game_Ready) ProtoMessage()               {}
func (*Game_Ready) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *Game_Ready) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_Ready) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 准备游戏的结果
type Game_AckReady struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Msg              *string      `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	UserId           *uint32      `protobuf:"varint,3,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckReady) Reset()                    { *m = Game_AckReady{} }
func (m *Game_AckReady) String() string            { return proto.CompactTextString(m) }
func (*Game_AckReady) ProtoMessage()               {}
func (*Game_AckReady) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *Game_AckReady) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckReady) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

func (m *Game_AckReady) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 聊天的内容
type Game_Message struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	MsgType          *int32       `protobuf:"varint,2,opt,name=msgType" json:"msgType,omitempty"`
	Id               *int32       `protobuf:"varint,3,opt,name=id" json:"id,omitempty"`
	Msg              *string      `protobuf:"bytes,4,opt,name=msg" json:"msg,omitempty"`
	UserId           *uint32      `protobuf:"varint,5,opt,name=userId" json:"userId,omitempty"`
	DeskId           *int32       `protobuf:"varint,6,opt,name=deskId" json:"deskId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_Message) Reset()                    { *m = Game_Message{} }
func (m *Game_Message) String() string            { return proto.CompactTextString(m) }
func (*Game_Message) ProtoMessage()               {}
func (*Game_Message) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

func (m *Game_Message) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_Message) GetMsgType() int32 {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return 0
}

func (m *Game_Message) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Game_Message) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

func (m *Game_Message) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_Message) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

// 消息广播
type Game_SendMessage struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	MsgType          *int32       `protobuf:"varint,2,opt,name=msgType" json:"msgType,omitempty"`
	Id               *int32       `protobuf:"varint,3,opt,name=id" json:"id,omitempty"`
	Msg              *string      `protobuf:"bytes,4,opt,name=msg" json:"msg,omitempty"`
	UserId           *uint32      `protobuf:"varint,5,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_SendMessage) Reset()                    { *m = Game_SendMessage{} }
func (m *Game_SendMessage) String() string            { return proto.CompactTextString(m) }
func (*Game_SendMessage) ProtoMessage()               {}
func (*Game_SendMessage) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{9} }

func (m *Game_SendMessage) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_SendMessage) GetMsgType() int32 {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return 0
}

func (m *Game_SendMessage) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Game_SendMessage) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

func (m *Game_SendMessage) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 赢牌信息：谁赢了多少
type WinCoinInfo struct {
	NickName         *string     `protobuf:"bytes,1,opt,name=nickName" json:"nickName,omitempty"`
	UserId           *uint32     `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	WinCoin          *int64      `protobuf:"varint,3,opt,name=winCoin" json:"winCoin,omitempty"`
	Coin             *int64      `protobuf:"varint,4,opt,name=coin" json:"coin,omitempty"`
	CardTitle        *string     `protobuf:"bytes,5,opt,name=cardTitle" json:"cardTitle,omitempty"`
	Cards            *PlayerCard `protobuf:"bytes,6,opt,name=cards" json:"cards,omitempty"`
	IsDealer         *bool       `protobuf:"varint,7,opt,name=isDealer" json:"isDealer,omitempty"`
	HuCount          *int32      `protobuf:"varint,8,opt,name=huCount" json:"huCount,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *WinCoinInfo) Reset()                    { *m = WinCoinInfo{} }
func (m *WinCoinInfo) String() string            { return proto.CompactTextString(m) }
func (*WinCoinInfo) ProtoMessage()               {}
func (*WinCoinInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{10} }

func (m *WinCoinInfo) GetNickName() string {
	if m != nil && m.NickName != nil {
		return *m.NickName
	}
	return ""
}

func (m *WinCoinInfo) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *WinCoinInfo) GetWinCoin() int64 {
	if m != nil && m.WinCoin != nil {
		return *m.WinCoin
	}
	return 0
}

func (m *WinCoinInfo) GetCoin() int64 {
	if m != nil && m.Coin != nil {
		return *m.Coin
	}
	return 0
}

func (m *WinCoinInfo) GetCardTitle() string {
	if m != nil && m.CardTitle != nil {
		return *m.CardTitle
	}
	return ""
}

func (m *WinCoinInfo) GetCards() *PlayerCard {
	if m != nil {
		return m.Cards
	}
	return nil
}

func (m *WinCoinInfo) GetIsDealer() bool {
	if m != nil && m.IsDealer != nil {
		return *m.IsDealer
	}
	return false
}

func (m *WinCoinInfo) GetHuCount() int32 {
	if m != nil && m.HuCount != nil {
		return *m.HuCount
	}
	return 0
}

type EndLotteryInfo struct {
	UserId           *uint32 `protobuf:"varint,1,opt,name=userId" json:"userId,omitempty"`
	NickName         *string `protobuf:"bytes,2,opt,name=nickName" json:"nickName,omitempty"`
	BigWin           *bool   `protobuf:"varint,3,opt,name=bigWin" json:"bigWin,omitempty"`
	IsOwner          *bool   `protobuf:"varint,4,opt,name=isOwner" json:"isOwner,omitempty"`
	WinCoin          *int64  `protobuf:"varint,5,opt,name=winCoin" json:"winCoin,omitempty"`
	CountHu          *int32  `protobuf:"varint,6,opt,name=countHu" json:"countHu,omitempty"`
	CountZiMo        *int32  `protobuf:"varint,7,opt,name=countZiMo" json:"countZiMo,omitempty"`
	CountDianPao     *int32  `protobuf:"varint,8,opt,name=countDianPao" json:"countDianPao,omitempty"`
	CountAnGang      *int32  `protobuf:"varint,9,opt,name=countAnGang" json:"countAnGang,omitempty"`
	CountMingGang    *int32  `protobuf:"varint,10,opt,name=countMingGang" json:"countMingGang,omitempty"`
	CountDianGang    *int32  `protobuf:"varint,11,opt,name=countDianGang" json:"countDianGang,omitempty"`
	CountChaJiao     *int32  `protobuf:"varint,12,opt,name=countChaJiao" json:"countChaJiao,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *EndLotteryInfo) Reset()                    { *m = EndLotteryInfo{} }
func (m *EndLotteryInfo) String() string            { return proto.CompactTextString(m) }
func (*EndLotteryInfo) ProtoMessage()               {}
func (*EndLotteryInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{11} }

func (m *EndLotteryInfo) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *EndLotteryInfo) GetNickName() string {
	if m != nil && m.NickName != nil {
		return *m.NickName
	}
	return ""
}

func (m *EndLotteryInfo) GetBigWin() bool {
	if m != nil && m.BigWin != nil {
		return *m.BigWin
	}
	return false
}

func (m *EndLotteryInfo) GetIsOwner() bool {
	if m != nil && m.IsOwner != nil {
		return *m.IsOwner
	}
	return false
}

func (m *EndLotteryInfo) GetWinCoin() int64 {
	if m != nil && m.WinCoin != nil {
		return *m.WinCoin
	}
	return 0
}

func (m *EndLotteryInfo) GetCountHu() int32 {
	if m != nil && m.CountHu != nil {
		return *m.CountHu
	}
	return 0
}

func (m *EndLotteryInfo) GetCountZiMo() int32 {
	if m != nil && m.CountZiMo != nil {
		return *m.CountZiMo
	}
	return 0
}

func (m *EndLotteryInfo) GetCountDianPao() int32 {
	if m != nil && m.CountDianPao != nil {
		return *m.CountDianPao
	}
	return 0
}

func (m *EndLotteryInfo) GetCountAnGang() int32 {
	if m != nil && m.CountAnGang != nil {
		return *m.CountAnGang
	}
	return 0
}

func (m *EndLotteryInfo) GetCountMingGang() int32 {
	if m != nil && m.CountMingGang != nil {
		return *m.CountMingGang
	}
	return 0
}

func (m *EndLotteryInfo) GetCountDianGang() int32 {
	if m != nil && m.CountDianGang != nil {
		return *m.CountDianGang
	}
	return 0
}

func (m *EndLotteryInfo) GetCountChaJiao() int32 {
	if m != nil && m.CountChaJiao != nil {
		return *m.CountChaJiao
	}
	return 0
}

// 本局结果(广播)
type Game_SendCurrentResult struct {
	Header           *ProtoHeader   `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	WinCoinInfo      []*WinCoinInfo `protobuf:"bytes,2,rep,name=winCoinInfo" json:"winCoinInfo,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *Game_SendCurrentResult) Reset()                    { *m = Game_SendCurrentResult{} }
func (m *Game_SendCurrentResult) String() string            { return proto.CompactTextString(m) }
func (*Game_SendCurrentResult) ProtoMessage()               {}
func (*Game_SendCurrentResult) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{12} }

func (m *Game_SendCurrentResult) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_SendCurrentResult) GetWinCoinInfo() []*WinCoinInfo {
	if m != nil {
		return m.WinCoinInfo
	}
	return nil
}

// 牌局结束(广播)
type Game_SendEndLottery struct {
	Header           *ProtoHeader      `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	CoinInfo         []*EndLotteryInfo `protobuf:"bytes,2,rep,name=coinInfo" json:"coinInfo,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Game_SendEndLottery) Reset()                    { *m = Game_SendEndLottery{} }
func (m *Game_SendEndLottery) String() string            { return proto.CompactTextString(m) }
func (*Game_SendEndLottery) ProtoMessage()               {}
func (*Game_SendEndLottery) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{13} }

func (m *Game_SendEndLottery) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_SendEndLottery) GetCoinInfo() []*EndLotteryInfo {
	if m != nil {
		return m.CoinInfo
	}
	return nil
}

func init() {
	proto.RegisterType((*Game_DissolveDesk)(nil), "mjproto.game_DissolveDesk")
	proto.RegisterType((*Game_AckDissolveDesk)(nil), "mjproto.game_AckDissolveDesk")
	proto.RegisterType((*Game_ReqDissolveDesk)(nil), "mjproto.game_ReqDissolveDesk")
	proto.RegisterType((*Game_AckReqDissolveDesk)(nil), "mjproto.game_AckReqDissolveDesk")
	proto.RegisterType((*Game_LeaveDesk)(nil), "mjproto.game_LeaveDesk")
	proto.RegisterType((*Game_AckLeaveDesk)(nil), "mjproto.game_AckLeaveDesk")
	proto.RegisterType((*Game_Ready)(nil), "mjproto.game_Ready")
	proto.RegisterType((*Game_AckReady)(nil), "mjproto.game_AckReady")
	proto.RegisterType((*Game_Message)(nil), "mjproto.game_Message")
	proto.RegisterType((*Game_SendMessage)(nil), "mjproto.game_SendMessage")
	proto.RegisterType((*WinCoinInfo)(nil), "mjproto.WinCoinInfo")
	proto.RegisterType((*EndLotteryInfo)(nil), "mjproto.EndLotteryInfo")
	proto.RegisterType((*Game_SendCurrentResult)(nil), "mjproto.game_SendCurrentResult")
	proto.RegisterType((*Game_SendEndLottery)(nil), "mjproto.game_SendEndLottery")
}

var fileDescriptor1 = []byte{
	// 582 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xbc, 0x54, 0x4f, 0x8f, 0xd2, 0x40,
	0x14, 0xb7, 0xb4, 0x05, 0xfa, 0x0a, 0xb8, 0x5b, 0x70, 0x97, 0x78, 0x32, 0x8d, 0x07, 0xbd, 0x70,
	0xd8, 0x83, 0x89, 0x47, 0x16, 0x8c, 0x60, 0x40, 0x09, 0x6e, 0x42, 0xe2, 0x65, 0x33, 0xdb, 0xce,
	0x96, 0x59, 0x60, 0x8a, 0x33, 0xad, 0x84, 0x8b, 0xf1, 0xe3, 0xf8, 0x0d, 0xbd, 0x3a, 0xf3, 0x5a,
	0xfe, 0xad, 0x17, 0x76, 0x35, 0x5e, 0x08, 0xf3, 0xeb, 0x9b, 0xdf, 0x9f, 0xf7, 0x5e, 0x0b, 0xde,
	0x82, 0x4c, 0xef, 0x62, 0x1e, 0x5d, 0x87, 0x54, 0xce, 0x5a, 0x4b, 0x11, 0x27, 0xb1, 0x57, 0x5a,
	0xdc, 0xe1, 0x9f, 0xe7, 0x70, 0x43, 0x24, 0xcd, 0x40, 0xbf, 0x0f, 0xa7, 0x11, 0x59, 0xd0, 0xeb,
	0x2e, 0x93, 0x32, 0x9e, 0x7f, 0xa3, 0x5d, 0x55, 0xef, 0xbd, 0x84, 0xe2, 0x94, 0x92, 0x90, 0x8a,
	0xa6, 0xf1, 0xc2, 0x78, 0xe5, 0x5e, 0x34, 0x5a, 0xf9, 0xd5, 0xd6, 0x48, 0xff, 0xf6, 0xf0, 0x99,
	0x57, 0x83, 0x62, 0x2a, 0xa9, 0xe8, 0x87, 0xcd, 0x82, 0xaa, 0xaa, 0xfa, 0x1c, 0x1a, 0x48, 0xd5,
	0x0e, 0x66, 0x7f, 0xcf, 0xa6, 0xcf, 0xda, 0xbb, 0x3a, 0x9b, 0xea, 0x6c, 0x7b, 0x27, 0x50, 0x5e,
	0x12, 0x29, 0x27, 0xb1, 0x08, 0x9b, 0x96, 0x42, 0x1c, 0x7f, 0x90, 0xeb, 0x8d, 0xe9, 0xd7, 0x7f,
	0xe0, 0xfe, 0x87, 0x01, 0xe7, 0x1b, 0xfb, 0x8f, 0x63, 0xac, 0x83, 0x9b, 0x31, 0xb6, 0x23, 0x41,
	0xa9, 0xa2, 0x35, 0x55, 0x0c, 0x0f, 0x20, 0x03, 0x27, 0x84, 0x25, 0x2a, 0x8a, 0xc6, 0xce, 0xa0,
	0x96, 0x61, 0x4a, 0x84, 0x60, 0xad, 0xa5, 0x71, 0xff, 0x0d, 0xd4, 0xd0, 0xc1, 0x80, 0x92, 0x87,
	0x08, 0xfb, 0x6f, 0xf3, 0x19, 0x2a, 0xe7, 0x0f, 0xbd, 0x7a, 0x09, 0x90, 0xf7, 0x90, 0x84, 0xeb,
	0x47, 0x76, 0x6e, 0x0c, 0xd5, 0x5d, 0xe3, 0x8e, 0xa7, 0x71, 0xc1, 0x5c, 0xc8, 0x08, 0x39, 0x9c,
	0x3d, 0x4e, 0x13, 0x39, 0xbf, 0x43, 0x05, 0x39, 0x87, 0x54, 0xaa, 0x16, 0xd1, 0x23, 0x29, 0x9f,
	0x42, 0x49, 0x51, 0x5e, 0xad, 0x97, 0x14, 0x69, 0x6d, 0xd5, 0xfc, 0x02, 0xdb, 0x2c, 0x50, 0xae,
	0x67, 0xdd, 0xd3, 0xb3, 0xef, 0x6d, 0x5b, 0x51, 0x17, 0xfb, 0x02, 0x4e, 0x50, 0xff, 0x33, 0xe5,
	0xe1, 0x7f, 0xf2, 0xe0, 0xff, 0x34, 0xc0, 0x9d, 0x30, 0xde, 0x89, 0x19, 0xef, 0xf3, 0xdb, 0x58,
	0x6f, 0x3c, 0x67, 0xc1, 0xec, 0xa3, 0xf2, 0x81, 0x8a, 0xce, 0x1f, 0xef, 0x88, 0xd2, 0x5a, 0x65,
	0x17, 0x90, 0xdf, 0xf4, 0x2a, 0x60, 0x05, 0xfa, 0x64, 0xe1, 0xe9, 0x14, 0x9c, 0x80, 0x88, 0xf0,
	0x8a, 0x25, 0x73, 0x8a, 0x1a, 0x8e, 0xe7, 0x83, 0xad, 0x21, 0x89, 0x31, 0xdd, 0x8b, 0xfa, 0x2e,
	0xc2, 0x9c, 0xac, 0xa9, 0xe8, 0xa8, 0x67, 0x5a, 0x97, 0xc9, 0x2e, 0x25, 0x73, 0x95, 0xb4, 0xa4,
	0xca, 0xca, 0x5a, 0x67, 0x9a, 0x76, 0xe2, 0x94, 0x27, 0xcd, 0x32, 0xb6, 0xe7, 0x97, 0x01, 0xb5,
	0x77, 0x3c, 0x1c, 0xc4, 0x49, 0x42, 0xc5, 0x1a, 0xdd, 0xee, 0xbc, 0x19, 0xe8, 0x6d, 0xdf, 0xfd,
	0x76, 0xc6, 0x37, 0x2c, 0x9a, 0xe4, 0x66, 0x91, 0x95, 0xc9, 0x4f, 0x2b, 0xae, 0x64, 0xac, 0x0d,
	0xb0, 0x89, 0x63, 0x63, 0x00, 0x05, 0x04, 0x5a, 0xb5, 0x97, 0x66, 0x63, 0xc1, 0x44, 0x1a, 0xf8,
	0xc2, 0x86, 0x31, 0x7a, 0xb3, 0xbd, 0x06, 0x54, 0x10, 0xea, 0x32, 0xc2, 0x47, 0x24, 0xce, 0x0c,
	0xea, 0x77, 0x11, 0xd1, 0x36, 0x7f, 0x4f, 0x78, 0xd4, 0x74, 0x10, 0x7c, 0x06, 0x55, 0x04, 0x87,
	0x8c, 0x47, 0x08, 0xc3, 0x01, 0xac, 0x19, 0x10, 0x76, 0x0f, 0x88, 0x3b, 0x53, 0xf2, 0x81, 0x29,
	0xe2, 0x0a, 0x26, 0x67, 0x70, 0xb6, 0x5d, 0x8c, 0x4e, 0x2a, 0x04, 0xe5, 0xc9, 0x98, 0xca, 0x74,
	0x9e, 0x1c, 0xb9, 0x1e, 0xaf, 0xc1, 0x5d, 0xed, 0x66, 0x8c, 0x1f, 0x89, 0xfd, 0xd2, 0xbd, 0xf9,
	0xfb, 0xb7, 0x50, 0xdf, 0x4a, 0xed, 0x9a, 0x7d, 0xb4, 0x4e, 0x39, 0x38, 0x14, 0x39, 0xdf, 0xd6,
	0x1d, 0x4e, 0xee, 0xb2, 0xd0, 0x33, 0x47, 0x4f, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0xc6, 0x58,
	0x5f, 0x20, 0x32, 0x06, 0x00, 0x00,
}
