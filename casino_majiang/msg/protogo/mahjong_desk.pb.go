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

// Ignoring public import of ComposePoker from base.proto

// Ignoring public import of PlayerPoker from base.proto

// Ignoring public import of PlayerInfo from base.proto

// Ignoring public import of DeskGameInfo from base.proto

// Ignoring public import of EProtoId from base.proto

// Ignoring public import of ErrorCode from base.proto

// Ignoring public import of MJRoomType from base.proto

// 解散房间
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

// 离开房间
type Game_LeaveDesk struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_LeaveDesk) Reset()                    { *m = Game_LeaveDesk{} }
func (m *Game_LeaveDesk) String() string            { return proto.CompactTextString(m) }
func (*Game_LeaveDesk) ProtoMessage()               {}
func (*Game_LeaveDesk) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

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
func (*Game_AckLeaveDesk) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

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
func (*Game_Ready) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

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
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckReady) Reset()                    { *m = Game_AckReady{} }
func (m *Game_AckReady) String() string            { return proto.CompactTextString(m) }
func (*Game_AckReady) ProtoMessage()               {}
func (*Game_AckReady) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

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
func (*Game_Message) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

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
func (*Game_SendMessage) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

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

// Start: 本局结束
// 赢牌牌型和类型
type CardType struct {
	Cards            []*CardInfo `protobuf:"bytes,1,rep,name=cards" json:"cards,omitempty"`
	CardType         *int32      `protobuf:"varint,2,opt,name=cardType" json:"cardType,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *CardType) Reset()                    { *m = CardType{} }
func (m *CardType) String() string            { return proto.CompactTextString(m) }
func (*CardType) ProtoMessage()               {}
func (*CardType) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

func (m *CardType) GetCards() []*CardInfo {
	if m != nil {
		return m.Cards
	}
	return nil
}

func (m *CardType) GetCardType() int32 {
	if m != nil && m.CardType != nil {
		return *m.CardType
	}
	return 0
}

// 赢牌信息：谁赢了多少
type WinCoinInfo struct {
	NickName         *string   `protobuf:"bytes,1,opt,name=nickName" json:"nickName,omitempty"`
	Seat             *int32    `protobuf:"varint,2,opt,name=seat" json:"seat,omitempty"`
	WinCoin          *int64    `protobuf:"varint,3,opt,name=winCoin" json:"winCoin,omitempty"`
	Coin             *int64    `protobuf:"varint,4,opt,name=coin" json:"coin,omitempty"`
	CardType         *CardType `protobuf:"bytes,5,opt,name=cardType" json:"cardType,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *WinCoinInfo) Reset()                    { *m = WinCoinInfo{} }
func (m *WinCoinInfo) String() string            { return proto.CompactTextString(m) }
func (*WinCoinInfo) ProtoMessage()               {}
func (*WinCoinInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{9} }

func (m *WinCoinInfo) GetNickName() string {
	if m != nil && m.NickName != nil {
		return *m.NickName
	}
	return ""
}

func (m *WinCoinInfo) GetSeat() int32 {
	if m != nil && m.Seat != nil {
		return *m.Seat
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

func (m *WinCoinInfo) GetCardType() *CardType {
	if m != nil {
		return m.CardType
	}
	return nil
}

// 牌局的结果（收到后客户端可以先 play 赢牌、亮牌动画，再显示结果弹窗）
type Game_CurrentResult struct {
	Header           *ProtoHeader   `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	MatchId          *int32         `protobuf:"varint,2,opt,name=matchId" json:"matchId,omitempty"`
	TableId          *int32         `protobuf:"varint,3,opt,name=tableId" json:"tableId,omitempty"`
	WinCoinInfo      []*WinCoinInfo `protobuf:"bytes,4,rep,name=winCoinInfo" json:"winCoinInfo,omitempty"`
	Rank             *int32         `protobuf:"varint,5,opt,name=rank" json:"rank,omitempty"`
	RankUserCount    *int32         `protobuf:"varint,6,opt,name=rankUserCount" json:"rankUserCount,omitempty"`
	CanRebuy         *bool          `protobuf:"varint,7,opt,name=canRebuy" json:"canRebuy,omitempty"`
	RebuyCount       *int32         `protobuf:"varint,8,opt,name=rebuyCount" json:"rebuyCount,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *Game_CurrentResult) Reset()                    { *m = Game_CurrentResult{} }
func (m *Game_CurrentResult) String() string            { return proto.CompactTextString(m) }
func (*Game_CurrentResult) ProtoMessage()               {}
func (*Game_CurrentResult) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{10} }

func (m *Game_CurrentResult) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_CurrentResult) GetMatchId() int32 {
	if m != nil && m.MatchId != nil {
		return *m.MatchId
	}
	return 0
}

func (m *Game_CurrentResult) GetTableId() int32 {
	if m != nil && m.TableId != nil {
		return *m.TableId
	}
	return 0
}

func (m *Game_CurrentResult) GetWinCoinInfo() []*WinCoinInfo {
	if m != nil {
		return m.WinCoinInfo
	}
	return nil
}

func (m *Game_CurrentResult) GetRank() int32 {
	if m != nil && m.Rank != nil {
		return *m.Rank
	}
	return 0
}

func (m *Game_CurrentResult) GetRankUserCount() int32 {
	if m != nil && m.RankUserCount != nil {
		return *m.RankUserCount
	}
	return 0
}

func (m *Game_CurrentResult) GetCanRebuy() bool {
	if m != nil && m.CanRebuy != nil {
		return *m.CanRebuy
	}
	return false
}

func (m *Game_CurrentResult) GetRebuyCount() int32 {
	if m != nil && m.RebuyCount != nil {
		return *m.RebuyCount
	}
	return 0
}

// Start: 全场结束
type EndLotteryInfo struct {
	Seat             *int32  `protobuf:"varint,1,opt,name=seat" json:"seat,omitempty"`
	Coin             *int64  `protobuf:"varint,2,opt,name=coin" json:"coin,omitempty"`
	BigWin           *bool   `protobuf:"varint,3,opt,name=bigWin" json:"bigWin,omitempty"`
	Owner            *bool   `protobuf:"varint,4,opt,name=owner" json:"owner,omitempty"`
	Rolename         *string `protobuf:"bytes,5,opt,name=rolename" json:"rolename,omitempty"`
	UserId           *uint32 `protobuf:"varint,6,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *EndLotteryInfo) Reset()                    { *m = EndLotteryInfo{} }
func (m *EndLotteryInfo) String() string            { return proto.CompactTextString(m) }
func (*EndLotteryInfo) ProtoMessage()               {}
func (*EndLotteryInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{11} }

func (m *EndLotteryInfo) GetSeat() int32 {
	if m != nil && m.Seat != nil {
		return *m.Seat
	}
	return 0
}

func (m *EndLotteryInfo) GetCoin() int64 {
	if m != nil && m.Coin != nil {
		return *m.Coin
	}
	return 0
}

func (m *EndLotteryInfo) GetBigWin() bool {
	if m != nil && m.BigWin != nil {
		return *m.BigWin
	}
	return false
}

func (m *EndLotteryInfo) GetOwner() bool {
	if m != nil && m.Owner != nil {
		return *m.Owner
	}
	return false
}

func (m *EndLotteryInfo) GetRolename() string {
	if m != nil && m.Rolename != nil {
		return *m.Rolename
	}
	return ""
}

func (m *EndLotteryInfo) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
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
func (*Game_SendEndLottery) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{12} }

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
	proto.RegisterType((*Game_LeaveDesk)(nil), "mjproto.game_LeaveDesk")
	proto.RegisterType((*Game_AckLeaveDesk)(nil), "mjproto.game_AckLeaveDesk")
	proto.RegisterType((*Game_Ready)(nil), "mjproto.game_Ready")
	proto.RegisterType((*Game_AckReady)(nil), "mjproto.game_AckReady")
	proto.RegisterType((*Game_Message)(nil), "mjproto.game_Message")
	proto.RegisterType((*Game_SendMessage)(nil), "mjproto.game_SendMessage")
	proto.RegisterType((*CardType)(nil), "mjproto.CardType")
	proto.RegisterType((*WinCoinInfo)(nil), "mjproto.WinCoinInfo")
	proto.RegisterType((*Game_CurrentResult)(nil), "mjproto.game_CurrentResult")
	proto.RegisterType((*EndLotteryInfo)(nil), "mjproto.EndLotteryInfo")
	proto.RegisterType((*Game_SendEndLottery)(nil), "mjproto.game_SendEndLottery")
}

var fileDescriptor1 = []byte{
	// 522 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xbc, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x25, 0x71, 0x9c, 0x8f, 0x49, 0x13, 0xda, 0xa5, 0x88, 0x88, 0x53, 0x65, 0x38, 0xd0, 0x4b,
	0x0e, 0x1c, 0x90, 0xb8, 0x20, 0x91, 0x14, 0xa9, 0x91, 0x0a, 0xaa, 0x02, 0x28, 0xc7, 0x6a, 0x63,
	0x4f, 0x13, 0x37, 0xf1, 0x6e, 0xd8, 0xb5, 0x89, 0x72, 0xe1, 0x47, 0xf2, 0x8b, 0x98, 0x1d, 0xdb,
	0x69, 0xd2, 0x53, 0x03, 0x12, 0x97, 0xc8, 0xb3, 0xfb, 0xe6, 0xbd, 0x37, 0x6f, 0x36, 0x20, 0x12,
	0x39, 0xbf, 0xd3, 0x6a, 0x76, 0x13, 0xa1, 0x5d, 0xf4, 0x57, 0x46, 0xa7, 0x5a, 0x34, 0x92, 0x3b,
	0xfe, 0x78, 0x09, 0x53, 0x69, 0x31, 0x3f, 0x0c, 0x46, 0x70, 0x32, 0x93, 0x09, 0xde, 0x5c, 0xc4,
	0xd6, 0xea, 0xe5, 0x4f, 0xbc, 0x20, 0xbc, 0x78, 0x0d, 0xf5, 0x39, 0xca, 0x08, 0x4d, 0xaf, 0x72,
	0x56, 0x79, 0xd3, 0x7e, 0x7b, 0xda, 0x2f, 0x5a, 0xfb, 0xd7, 0xee, 0xf7, 0x92, 0xef, 0x44, 0x17,
	0xea, 0x99, 0x45, 0x33, 0x8a, 0x7a, 0x55, 0x42, 0x75, 0x02, 0x05, 0xa7, 0x4c, 0xf5, 0x31, 0x5c,
	0xfc, 0x3b, 0x9b, 0xab, 0x9d, 0x77, 0xaa, 0x3d, 0xaa, 0x7d, 0x71, 0x0c, 0xcd, 0x95, 0xb4, 0x76,
	0xa2, 0x4d, 0xd4, 0xab, 0xd1, 0x49, 0x2b, 0x78, 0x07, 0x5d, 0xd6, 0xbb, 0x42, 0x79, 0x88, 0x52,
	0xf0, 0xbe, 0x18, 0x99, 0x7c, 0x1e, 0xda, 0x3a, 0x00, 0xe0, 0xd6, 0x31, 0x95, 0x9b, 0xbf, 0x8c,
	0x69, 0x00, 0x9d, 0x52, 0xfe, 0x10, 0x9a, 0x36, 0x78, 0x89, 0x9d, 0x31, 0x47, 0x2b, 0xf8, 0x05,
	0x47, 0xcc, 0xf1, 0x19, 0xad, 0x95, 0x33, 0x7c, 0x24, 0xc5, 0x53, 0x68, 0x10, 0xc5, 0xb7, 0xcd,
	0x0a, 0x99, 0xc6, 0x17, 0x00, 0xd5, 0xb8, 0xcc, 0xb7, 0xe0, 0xe7, 0x68, 0x77, 0x3c, 0xfb, 0x0f,
	0x96, 0x51, 0x77, 0xe0, 0xc0, 0xc0, 0x31, 0xeb, 0x7f, 0x45, 0x15, 0xfd, 0x27, 0x0f, 0xc1, 0x07,
	0x68, 0x0e, 0xa5, 0x89, 0x5c, 0xab, 0x38, 0x03, 0x3f, 0xa4, 0x6f, 0x4b, 0x52, 0x1e, 0x49, 0x9d,
	0x6c, 0xa5, 0x1c, 0x62, 0xa4, 0x6e, 0xb5, 0x7b, 0x2e, 0x61, 0x81, 0xce, 0x85, 0x82, 0x1f, 0xd0,
	0x9e, 0xc4, 0x6a, 0xa8, 0x63, 0x55, 0x02, 0x54, 0x1c, 0x2e, 0xbe, 0xd0, 0x18, 0x6c, 0xb8, 0x25,
	0x8e, 0xa0, 0x66, 0x51, 0xa6, 0x85, 0x2f, 0x32, 0xba, 0xce, 0xe1, 0x6c, 0xce, 0x73, 0xd7, 0xa1,
	0xab, 0x6a, 0x5c, 0xbd, 0xda, 0xe1, 0xf7, 0x79, 0xde, 0x7d, 0x13, 0xee, 0x22, 0xf8, 0x5d, 0x01,
	0xc1, 0x39, 0x0d, 0x33, 0x63, 0x50, 0xa5, 0x63, 0xb4, 0xd9, 0x32, 0x3d, 0x20, 0x29, 0x99, 0x86,
	0xf3, 0xe2, 0xe1, 0xb0, 0xa3, 0x54, 0x4e, 0x97, 0xb8, 0xfd, 0x4b, 0x9c, 0x43, 0x7b, 0x7d, 0x3f,
	0x11, 0x19, 0xf3, 0xf6, 0xc8, 0x76, 0xa7, 0x25, 0xf3, 0x46, 0xaa, 0x05, 0x5b, 0xf5, 0xc5, 0x73,
	0xe8, 0xb8, 0xea, 0x3b, 0xc5, 0x3b, 0xd4, 0x99, 0x4a, 0xf3, 0xad, 0xe6, 0x99, 0xa9, 0x31, 0x4e,
	0xb3, 0x4d, 0xaf, 0x41, 0x27, 0x4d, 0x41, 0xdb, 0x31, 0xae, 0xcc, 0x51, 0x4d, 0xce, 0x31, 0x81,
	0xee, 0x27, 0x15, 0x5d, 0xe9, 0x34, 0x45, 0xb3, 0x29, 0xc9, 0x39, 0xb8, 0x0a, 0xb3, 0x94, 0x39,
	0x55, 0x39, 0x27, 0xda, 0xe2, 0x34, 0x9e, 0x4d, 0x8a, 0x14, 0x9b, 0xa2, 0x03, 0xbe, 0x5e, 0x2b,
	0x1a, 0xbd, 0xc6, 0x25, 0x49, 0x1a, 0xbd, 0x44, 0xe5, 0xb6, 0xe0, 0x3f, 0x58, 0x7b, 0x9d, 0xd7,
	0x7e, 0x0b, 0xcf, 0xb6, 0x4f, 0xed, 0x5e, 0xf7, 0x91, 0x19, 0x9e, 0xd3, 0x44, 0x65, 0x3c, 0x55,
	0x8e, 0xe7, 0xc5, 0x16, 0xb7, 0x3f, 0xc4, 0xa0, 0x7a, 0xe9, 0x5d, 0x3f, 0xf9, 0x13, 0x00, 0x00,
	0xff, 0xff, 0x6e, 0x25, 0x4c, 0x72, 0x38, 0x05, 0x00, 0x00,
}
