// Code generated by protoc-gen-go.
// source: mahjong_play.proto
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

// Ignoring public import of HuPaiType from base.proto

// 开局（接收服务端消息）
type Game_Opening struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	MatchId          *int32       `protobuf:"varint,2,opt,name=matchId" json:"matchId,omitempty"`
	TableId          *int32       `protobuf:"varint,3,opt,name=tableId" json:"tableId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_Opening) Reset()                    { *m = Game_Opening{} }
func (m *Game_Opening) String() string            { return proto.CompactTextString(m) }
func (*Game_Opening) ProtoMessage()               {}
func (*Game_Opening) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *Game_Opening) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_Opening) GetMatchId() int32 {
	if m != nil && m.MatchId != nil {
		return *m.MatchId
	}
	return 0
}

func (m *Game_Opening) GetTableId() int32 {
	if m != nil && m.TableId != nil {
		return *m.TableId
	}
	return 0
}

// 发牌
type Game_DealCards struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	PlayerCard       *PlayerCard  `protobuf:"bytes,2,opt,name=playerCard" json:"playerCard,omitempty"`
	OtherPlayerCards []int32      `protobuf:"varint,3,rep,name=otherPlayerCards" json:"otherPlayerCards,omitempty"`
	DealerUserId     *uint32      `protobuf:"varint,4,opt,name=dealerUserId" json:"dealerUserId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_DealCards) Reset()                    { *m = Game_DealCards{} }
func (m *Game_DealCards) String() string            { return proto.CompactTextString(m) }
func (*Game_DealCards) ProtoMessage()               {}
func (*Game_DealCards) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *Game_DealCards) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_DealCards) GetPlayerCard() *PlayerCard {
	if m != nil {
		return m.PlayerCard
	}
	return nil
}

func (m *Game_DealCards) GetOtherPlayerCards() []int32 {
	if m != nil {
		return m.OtherPlayerCards
	}
	return nil
}

func (m *Game_DealCards) GetDealerUserId() uint32 {
	if m != nil && m.DealerUserId != nil {
		return *m.DealerUserId
	}
	return 0
}

// 换牌（3张）
type Game_ExchangeCards struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	MatchId          *int32       `protobuf:"varint,2,opt,name=matchId" json:"matchId,omitempty"`
	TableId          *int32       `protobuf:"varint,3,opt,name=tableId" json:"tableId,omitempty"`
	Seat             *int32       `protobuf:"varint,4,opt,name=seat" json:"seat,omitempty"`
	ExchangeOutCards []*CardInfo  `protobuf:"bytes,5,rep,name=exchangeOutCards" json:"exchangeOutCards,omitempty"`
	UserId           *uint32      `protobuf:"varint,6,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_ExchangeCards) Reset()                    { *m = Game_ExchangeCards{} }
func (m *Game_ExchangeCards) String() string            { return proto.CompactTextString(m) }
func (*Game_ExchangeCards) ProtoMessage()               {}
func (*Game_ExchangeCards) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

func (m *Game_ExchangeCards) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_ExchangeCards) GetMatchId() int32 {
	if m != nil && m.MatchId != nil {
		return *m.MatchId
	}
	return 0
}

func (m *Game_ExchangeCards) GetTableId() int32 {
	if m != nil && m.TableId != nil {
		return *m.TableId
	}
	return 0
}

func (m *Game_ExchangeCards) GetSeat() int32 {
	if m != nil && m.Seat != nil {
		return *m.Seat
	}
	return 0
}

func (m *Game_ExchangeCards) GetExchangeOutCards() []*CardInfo {
	if m != nil {
		return m.ExchangeOutCards
	}
	return nil
}

func (m *Game_ExchangeCards) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

type Game_AckExchangeCards struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	MatchId          *int32       `protobuf:"varint,2,opt,name=matchId" json:"matchId,omitempty"`
	TableId          *int32       `protobuf:"varint,3,opt,name=tableId" json:"tableId,omitempty"`
	Seat             *int32       `protobuf:"varint,4,opt,name=seat" json:"seat,omitempty"`
	ExchangeOutCards []*CardInfo  `protobuf:"bytes,5,rep,name=exchangeOutCards" json:"exchangeOutCards,omitempty"`
	ExchangeOutseat  *int32       `protobuf:"varint,6,opt,name=exchangeOutseat" json:"exchangeOutseat,omitempty"`
	ExchangeInCards  []*CardInfo  `protobuf:"bytes,7,rep,name=exchangeInCards" json:"exchangeInCards,omitempty"`
	ExchangeInseat   *int32       `protobuf:"varint,8,opt,name=exchangeInseat" json:"exchangeInseat,omitempty"`
	Dice             *int32       `protobuf:"varint,9,opt,name=dice" json:"dice,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckExchangeCards) Reset()                    { *m = Game_AckExchangeCards{} }
func (m *Game_AckExchangeCards) String() string            { return proto.CompactTextString(m) }
func (*Game_AckExchangeCards) ProtoMessage()               {}
func (*Game_AckExchangeCards) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *Game_AckExchangeCards) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckExchangeCards) GetMatchId() int32 {
	if m != nil && m.MatchId != nil {
		return *m.MatchId
	}
	return 0
}

func (m *Game_AckExchangeCards) GetTableId() int32 {
	if m != nil && m.TableId != nil {
		return *m.TableId
	}
	return 0
}

func (m *Game_AckExchangeCards) GetSeat() int32 {
	if m != nil && m.Seat != nil {
		return *m.Seat
	}
	return 0
}

func (m *Game_AckExchangeCards) GetExchangeOutCards() []*CardInfo {
	if m != nil {
		return m.ExchangeOutCards
	}
	return nil
}

func (m *Game_AckExchangeCards) GetExchangeOutseat() int32 {
	if m != nil && m.ExchangeOutseat != nil {
		return *m.ExchangeOutseat
	}
	return 0
}

func (m *Game_AckExchangeCards) GetExchangeInCards() []*CardInfo {
	if m != nil {
		return m.ExchangeInCards
	}
	return nil
}

func (m *Game_AckExchangeCards) GetExchangeInseat() int32 {
	if m != nil && m.ExchangeInseat != nil {
		return *m.ExchangeInseat
	}
	return 0
}

func (m *Game_AckExchangeCards) GetDice() int32 {
	if m != nil && m.Dice != nil {
		return *m.Dice
	}
	return 0
}

// 定缺（和个人玩家ACK）
type Game_DingQue struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	MatchId          *int32       `protobuf:"varint,2,opt,name=matchId" json:"matchId,omitempty"`
	TableId          *int32       `protobuf:"varint,3,opt,name=tableId" json:"tableId,omitempty"`
	Color            *int32       `protobuf:"varint,4,opt,name=color" json:"color,omitempty"`
	UserId           *uint32      `protobuf:"varint,5,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_DingQue) Reset()                    { *m = Game_DingQue{} }
func (m *Game_DingQue) String() string            { return proto.CompactTextString(m) }
func (*Game_DingQue) ProtoMessage()               {}
func (*Game_DingQue) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *Game_DingQue) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_DingQue) GetMatchId() int32 {
	if m != nil && m.MatchId != nil {
		return *m.MatchId
	}
	return 0
}

func (m *Game_DingQue) GetTableId() int32 {
	if m != nil && m.TableId != nil {
		return *m.TableId
	}
	return 0
}

func (m *Game_DingQue) GetColor() int32 {
	if m != nil && m.Color != nil {
		return *m.Color
	}
	return 0
}

func (m *Game_DingQue) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

type DingQueEndBean struct {
	UserId           *uint32 `protobuf:"varint,1,opt,name=userId" json:"userId,omitempty"`
	Flower           *int32  `protobuf:"varint,2,opt,name=flower" json:"flower,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DingQueEndBean) Reset()                    { *m = DingQueEndBean{} }
func (m *DingQueEndBean) String() string            { return proto.CompactTextString(m) }
func (*DingQueEndBean) ProtoMessage()               {}
func (*DingQueEndBean) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *DingQueEndBean) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *DingQueEndBean) GetFlower() int32 {
	if m != nil && m.Flower != nil {
		return *m.Flower
	}
	return 0
}

type Game_DingQueEnd struct {
	Header           *ProtoHeader      `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Ques             []*DingQueEndBean `protobuf:"bytes,2,rep,name=Ques" json:"Ques,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Game_DingQueEnd) Reset()                    { *m = Game_DingQueEnd{} }
func (m *Game_DingQueEnd) String() string            { return proto.CompactTextString(m) }
func (*Game_DingQueEnd) ProtoMessage()               {}
func (*Game_DingQueEnd) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

func (m *Game_DingQueEnd) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_DingQueEnd) GetQues() []*DingQueEndBean {
	if m != nil {
		return m.Ques
	}
	return nil
}

// 定缺开始广播（和ACK）
type Game_BroadcastBeginDingQue struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Color            []int32      `protobuf:"varint,2,rep,name=color" json:"color,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_BroadcastBeginDingQue) Reset()                    { *m = Game_BroadcastBeginDingQue{} }
func (m *Game_BroadcastBeginDingQue) String() string            { return proto.CompactTextString(m) }
func (*Game_BroadcastBeginDingQue) ProtoMessage()               {}
func (*Game_BroadcastBeginDingQue) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{7} }

func (m *Game_BroadcastBeginDingQue) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_BroadcastBeginDingQue) GetColor() []int32 {
	if m != nil {
		return m.Color
	}
	return nil
}

// 换牌开始(广播)
type Game_BroadcastBeginExchange struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_BroadcastBeginExchange) Reset()                    { *m = Game_BroadcastBeginExchange{} }
func (m *Game_BroadcastBeginExchange) String() string            { return proto.CompactTextString(m) }
func (*Game_BroadcastBeginExchange) ProtoMessage()               {}
func (*Game_BroadcastBeginExchange) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{8} }

func (m *Game_BroadcastBeginExchange) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

// 摸牌
type Game_GetInCard struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Card             *CardInfo    `protobuf:"bytes,2,opt,name=card" json:"card,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_GetInCard) Reset()                    { *m = Game_GetInCard{} }
func (m *Game_GetInCard) String() string            { return proto.CompactTextString(m) }
func (*Game_GetInCard) ProtoMessage()               {}
func (*Game_GetInCard) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{9} }

func (m *Game_GetInCard) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_GetInCard) GetCard() *CardInfo {
	if m != nil {
		return m.Card
	}
	return nil
}

// 出牌
type Game_SendOutCard struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	CardId           *int32       `protobuf:"varint,2,opt,name=cardId" json:"cardId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_SendOutCard) Reset()                    { *m = Game_SendOutCard{} }
func (m *Game_SendOutCard) String() string            { return proto.CompactTextString(m) }
func (*Game_SendOutCard) ProtoMessage()               {}
func (*Game_SendOutCard) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{10} }

func (m *Game_SendOutCard) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_SendOutCard) GetCardId() int32 {
	if m != nil && m.CardId != nil {
		return *m.CardId
	}
	return 0
}

type Game_AckSendOutCard struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Result           *int32       `protobuf:"varint,2,opt,name=result" json:"result,omitempty"`
	UserId           *uint32      `protobuf:"varint,3,opt,name=userId" json:"userId,omitempty"`
	Card             *CardInfo    `protobuf:"bytes,4,opt,name=card" json:"card,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckSendOutCard) Reset()                    { *m = Game_AckSendOutCard{} }
func (m *Game_AckSendOutCard) String() string            { return proto.CompactTextString(m) }
func (*Game_AckSendOutCard) ProtoMessage()               {}
func (*Game_AckSendOutCard) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{11} }

func (m *Game_AckSendOutCard) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckSendOutCard) GetResult() int32 {
	if m != nil && m.Result != nil {
		return *m.Result
	}
	return 0
}

func (m *Game_AckSendOutCard) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_AckSendOutCard) GetCard() *CardInfo {
	if m != nil {
		return m.Card
	}
	return nil
}

// 碰牌
type Game_ActPeng struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	PengCard         *CardInfo    `protobuf:"bytes,3,opt,name=pengCard" json:"pengCard,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_ActPeng) Reset()                    { *m = Game_ActPeng{} }
func (m *Game_ActPeng) String() string            { return proto.CompactTextString(m) }
func (*Game_ActPeng) ProtoMessage()               {}
func (*Game_ActPeng) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{12} }

func (m *Game_ActPeng) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_ActPeng) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_ActPeng) GetPengCard() *CardInfo {
	if m != nil {
		return m.PengCard
	}
	return nil
}

type Game_AckActPeng struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	PengCard         []*CardInfo  `protobuf:"bytes,2,rep,name=pengCard" json:"pengCard,omitempty"`
	UserIdOut        *uint32      `protobuf:"varint,3,opt,name=userIdOut" json:"userIdOut,omitempty"`
	UserIdIn         *uint32      `protobuf:"varint,4,opt,name=userIdIn" json:"userIdIn,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckActPeng) Reset()                    { *m = Game_AckActPeng{} }
func (m *Game_AckActPeng) String() string            { return proto.CompactTextString(m) }
func (*Game_AckActPeng) ProtoMessage()               {}
func (*Game_AckActPeng) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{13} }

func (m *Game_AckActPeng) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckActPeng) GetPengCard() []*CardInfo {
	if m != nil {
		return m.PengCard
	}
	return nil
}

func (m *Game_AckActPeng) GetUserIdOut() uint32 {
	if m != nil && m.UserIdOut != nil {
		return *m.UserIdOut
	}
	return 0
}

func (m *Game_AckActPeng) GetUserIdIn() uint32 {
	if m != nil && m.UserIdIn != nil {
		return *m.UserIdIn
	}
	return 0
}

// 杠牌
type Game_ActGang struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	GangCard         *CardInfo    `protobuf:"bytes,3,opt,name=gangCard" json:"gangCard,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_ActGang) Reset()                    { *m = Game_ActGang{} }
func (m *Game_ActGang) String() string            { return proto.CompactTextString(m) }
func (*Game_ActGang) ProtoMessage()               {}
func (*Game_ActGang) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{14} }

func (m *Game_ActGang) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_ActGang) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_ActGang) GetGangCard() *CardInfo {
	if m != nil {
		return m.GangCard
	}
	return nil
}

type Game_AckActGang struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	GangType         *int32       `protobuf:"varint,2,opt,name=gangType" json:"gangType,omitempty"`
	GangCard         []*CardInfo  `protobuf:"bytes,3,rep,name=GangCard" json:"GangCard,omitempty"`
	UserIdOut        *uint32      `protobuf:"varint,4,opt,name=userIdOut" json:"userIdOut,omitempty"`
	UserIdIn         *uint32      `protobuf:"varint,5,opt,name=userIdIn" json:"userIdIn,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckActGang) Reset()                    { *m = Game_AckActGang{} }
func (m *Game_AckActGang) String() string            { return proto.CompactTextString(m) }
func (*Game_AckActGang) ProtoMessage()               {}
func (*Game_AckActGang) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{15} }

func (m *Game_AckActGang) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckActGang) GetGangType() int32 {
	if m != nil && m.GangType != nil {
		return *m.GangType
	}
	return 0
}

func (m *Game_AckActGang) GetGangCard() []*CardInfo {
	if m != nil {
		return m.GangCard
	}
	return nil
}

func (m *Game_AckActGang) GetUserIdOut() uint32 {
	if m != nil && m.UserIdOut != nil {
		return *m.UserIdOut
	}
	return 0
}

func (m *Game_AckActGang) GetUserIdIn() uint32 {
	if m != nil && m.UserIdIn != nil {
		return *m.UserIdIn
	}
	return 0
}

// 胡牌
type Game_ActHu struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	HuCard           *CardInfo    `protobuf:"bytes,3,opt,name=huCard" json:"huCard,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_ActHu) Reset()                    { *m = Game_ActHu{} }
func (m *Game_ActHu) String() string            { return proto.CompactTextString(m) }
func (*Game_ActHu) ProtoMessage()               {}
func (*Game_ActHu) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{16} }

func (m *Game_ActHu) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_ActHu) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_ActHu) GetHuCard() *CardInfo {
	if m != nil {
		return m.HuCard
	}
	return nil
}

type Game_AckActHu struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	HuType           *int32       `protobuf:"varint,2,opt,name=huType" json:"huType,omitempty"`
	HuCard           *CardInfo    `protobuf:"bytes,3,opt,name=huCard" json:"huCard,omitempty"`
	UserIdOut        *uint32      `protobuf:"varint,4,opt,name=userIdOut" json:"userIdOut,omitempty"`
	UserIdIn         *uint32      `protobuf:"varint,5,opt,name=userIdIn" json:"userIdIn,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckActHu) Reset()                    { *m = Game_AckActHu{} }
func (m *Game_AckActHu) String() string            { return proto.CompactTextString(m) }
func (*Game_AckActHu) ProtoMessage()               {}
func (*Game_AckActHu) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{17} }

func (m *Game_AckActHu) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckActHu) GetHuType() int32 {
	if m != nil && m.HuType != nil {
		return *m.HuType
	}
	return 0
}

func (m *Game_AckActHu) GetHuCard() *CardInfo {
	if m != nil {
		return m.HuCard
	}
	return nil
}

func (m *Game_AckActHu) GetUserIdOut() uint32 {
	if m != nil && m.UserIdOut != nil {
		return *m.UserIdOut
	}
	return 0
}

func (m *Game_AckActHu) GetUserIdIn() uint32 {
	if m != nil && m.UserIdIn != nil {
		return *m.UserIdIn
	}
	return 0
}

// 过牌
type Game_ActGuo struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_ActGuo) Reset()                    { *m = Game_ActGuo{} }
func (m *Game_ActGuo) String() string            { return proto.CompactTextString(m) }
func (*Game_ActGuo) ProtoMessage()               {}
func (*Game_ActGuo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{18} }

func (m *Game_ActGuo) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_ActGuo) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 过牌收到服务器确认后，还是要协议 【PID_game_SENDOVERTURN ：game_SendOverTurn】 后结束此轮
type Game_AckActGuo struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckActGuo) Reset()                    { *m = Game_AckActGuo{} }
func (m *Game_AckActGuo) String() string            { return proto.CompactTextString(m) }
func (*Game_AckActGuo) ProtoMessage()               {}
func (*Game_AckActGuo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{19} }

func (m *Game_AckActGuo) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckActGuo) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 轮到谁操作(摸牌、碰/杠/过/胡)
type Game_OverTurn struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	NextUserId       *uint32      `protobuf:"varint,3,opt,name=nextUserId" json:"nextUserId,omitempty"`
	ActType          *int32       `protobuf:"varint,4,opt,name=actType" json:"actType,omitempty"`
	CanPeng          *bool        `protobuf:"varint,5,opt,name=canPeng" json:"canPeng,omitempty"`
	CanGang          *bool        `protobuf:"varint,6,opt,name=canGang" json:"canGang,omitempty"`
	CanHu            *bool        `protobuf:"varint,7,opt,name=canHu" json:"canHu,omitempty"`
	ActCard          *CardInfo    `protobuf:"bytes,8,opt,name=actCard" json:"actCard,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_OverTurn) Reset()                    { *m = Game_OverTurn{} }
func (m *Game_OverTurn) String() string            { return proto.CompactTextString(m) }
func (*Game_OverTurn) ProtoMessage()               {}
func (*Game_OverTurn) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{20} }

func (m *Game_OverTurn) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_OverTurn) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_OverTurn) GetNextUserId() uint32 {
	if m != nil && m.NextUserId != nil {
		return *m.NextUserId
	}
	return 0
}

func (m *Game_OverTurn) GetActType() int32 {
	if m != nil && m.ActType != nil {
		return *m.ActType
	}
	return 0
}

func (m *Game_OverTurn) GetCanPeng() bool {
	if m != nil && m.CanPeng != nil {
		return *m.CanPeng
	}
	return false
}

func (m *Game_OverTurn) GetCanGang() bool {
	if m != nil && m.CanGang != nil {
		return *m.CanGang
	}
	return false
}

func (m *Game_OverTurn) GetCanHu() bool {
	if m != nil && m.CanHu != nil {
		return *m.CanHu
	}
	return false
}

func (m *Game_OverTurn) GetActCard() *CardInfo {
	if m != nil {
		return m.ActCard
	}
	return nil
}

// 发送游戏信息(广播)
type Game_SendGameInfo struct {
	Header *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	// 1. 首先是牌桌的玩家数据（玩家数据包括其id昵称筹码头像等基本信息，其手牌数据，以及自己打出的牌的数据，还有状态是否已经胡牌了，玩家在整局的总输赢）
	PlayerInfo []*PlayerInfo `protobuf:"bytes,2,rep,name=playerInfo" json:"playerInfo,omitempty"`
	// 2. 桌面信息（包括：游戏是否结束，当前轮到哪个玩家，倒计时剩余时间）
	DeskGameInfo *DeskGameInfo `protobuf:"bytes,3,opt,name=deskGameInfo" json:"deskGameInfo,omitempty"`
	//
	SenderUserId     *uint32 `protobuf:"varint,4,opt,name=senderUserId" json:"senderUserId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Game_SendGameInfo) Reset()                    { *m = Game_SendGameInfo{} }
func (m *Game_SendGameInfo) String() string            { return proto.CompactTextString(m) }
func (*Game_SendGameInfo) ProtoMessage()               {}
func (*Game_SendGameInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{21} }

func (m *Game_SendGameInfo) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_SendGameInfo) GetPlayerInfo() []*PlayerInfo {
	if m != nil {
		return m.PlayerInfo
	}
	return nil
}

func (m *Game_SendGameInfo) GetDeskGameInfo() *DeskGameInfo {
	if m != nil {
		return m.DeskGameInfo
	}
	return nil
}

func (m *Game_SendGameInfo) GetSenderUserId() uint32 {
	if m != nil && m.SenderUserId != nil {
		return *m.SenderUserId
	}
	return 0
}

func init() {
	proto.RegisterType((*Game_Opening)(nil), "mjproto.game_Opening")
	proto.RegisterType((*Game_DealCards)(nil), "mjproto.game_DealCards")
	proto.RegisterType((*Game_ExchangeCards)(nil), "mjproto.game_ExchangeCards")
	proto.RegisterType((*Game_AckExchangeCards)(nil), "mjproto.game_AckExchangeCards")
	proto.RegisterType((*Game_DingQue)(nil), "mjproto.game_DingQue")
	proto.RegisterType((*DingQueEndBean)(nil), "mjproto.DingQueEndBean")
	proto.RegisterType((*Game_DingQueEnd)(nil), "mjproto.game_DingQueEnd")
	proto.RegisterType((*Game_BroadcastBeginDingQue)(nil), "mjproto.game_BroadcastBeginDingQue")
	proto.RegisterType((*Game_BroadcastBeginExchange)(nil), "mjproto.game_BroadcastBeginExchange")
	proto.RegisterType((*Game_GetInCard)(nil), "mjproto.game_GetInCard")
	proto.RegisterType((*Game_SendOutCard)(nil), "mjproto.game_SendOutCard")
	proto.RegisterType((*Game_AckSendOutCard)(nil), "mjproto.game_AckSendOutCard")
	proto.RegisterType((*Game_ActPeng)(nil), "mjproto.game_ActPeng")
	proto.RegisterType((*Game_AckActPeng)(nil), "mjproto.game_AckActPeng")
	proto.RegisterType((*Game_ActGang)(nil), "mjproto.game_ActGang")
	proto.RegisterType((*Game_AckActGang)(nil), "mjproto.game_AckActGang")
	proto.RegisterType((*Game_ActHu)(nil), "mjproto.game_ActHu")
	proto.RegisterType((*Game_AckActHu)(nil), "mjproto.game_AckActHu")
	proto.RegisterType((*Game_ActGuo)(nil), "mjproto.game_ActGuo")
	proto.RegisterType((*Game_AckActGuo)(nil), "mjproto.game_AckActGuo")
	proto.RegisterType((*Game_OverTurn)(nil), "mjproto.game_OverTurn")
	proto.RegisterType((*Game_SendGameInfo)(nil), "mjproto.game_SendGameInfo")
}

var fileDescriptor3 = []byte{
	// 747 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xcc, 0x56, 0x4f, 0x53, 0xd3, 0x40,
	0x14, 0x37, 0x6d, 0x5a, 0xca, 0x2b, 0x2d, 0x25, 0x80, 0x74, 0xf0, 0xa0, 0x46, 0x1d, 0x1d, 0x99,
	0x61, 0x1c, 0xbe, 0x01, 0x05, 0x6c, 0x7b, 0x02, 0x46, 0xd0, 0x9b, 0xcc, 0x92, 0x3c, 0xd2, 0x42,
	0xbb, 0xe9, 0x24, 0x1b, 0x85, 0x19, 0x2f, 0x1e, 0x3c, 0x3a, 0x1e, 0xfc, 0x0c, 0x8e, 0x1f, 0xc3,
	0xaf, 0xe6, 0xee, 0xcb, 0xa6, 0x09, 0x62, 0x25, 0xad, 0x1e, 0xbc, 0x74, 0x9a, 0xfd, 0xf3, 0xfb,
	0xf3, 0xde, 0xdb, 0xb7, 0x0b, 0xd6, 0x90, 0xf5, 0xce, 0x7d, 0xee, 0x9d, 0x8c, 0x06, 0xec, 0x6a,
	0x73, 0x14, 0xf8, 0xc2, 0xb7, 0xe6, 0x86, 0xe7, 0xf4, 0x67, 0x1d, 0x4e, 0x59, 0x88, 0xf1, 0xa0,
	0xfd, 0x1a, 0x16, 0x3c, 0x36, 0xc4, 0x93, 0xfd, 0x11, 0xf2, 0x3e, 0xf7, 0xac, 0xc7, 0x50, 0xee,
	0x21, 0x73, 0x31, 0x68, 0x1a, 0x0f, 0x8c, 0x67, 0xd5, 0xad, 0x95, 0x4d, 0xbd, 0x6b, 0xf3, 0x40,
	0xfd, 0x76, 0x68, 0xce, 0x5a, 0x84, 0xb9, 0x21, 0x13, 0x4e, 0xaf, 0xeb, 0x36, 0x0b, 0x72, 0x59,
	0x49, 0x0d, 0x08, 0x76, 0x3a, 0x40, 0x39, 0x50, 0x54, 0x03, 0xf6, 0x17, 0x03, 0xea, 0x04, 0xbc,
	0x8b, 0x6c, 0xb0, 0xc3, 0x02, 0x37, 0xcc, 0x09, 0xfd, 0x14, 0x40, 0x69, 0xc6, 0x40, 0x6d, 0x22,
	0xf4, 0xea, 0xd6, 0x72, 0xba, 0x72, 0x3c, 0x65, 0x35, 0xa1, 0xe1, 0x8b, 0x1e, 0x06, 0xe9, 0x50,
	0x28, 0xb9, 0x8b, 0x52, 0xcc, 0x0a, 0x2c, 0xb8, 0x92, 0x15, 0x83, 0xe3, 0x10, 0x03, 0xa9, 0xc8,
	0x94, 0x20, 0x35, 0xfb, 0xbb, 0x01, 0x16, 0x29, 0xda, 0xbb, 0x74, 0x7a, 0x8c, 0x7b, 0x38, 0x8d,
	0xaa, 0x5b, 0x0d, 0x5b, 0x0b, 0x60, 0x86, 0xc8, 0x04, 0x91, 0x95, 0xac, 0x0d, 0x68, 0xa0, 0xa6,
	0xd9, 0x8f, 0x44, 0x2c, 0xae, 0x24, 0xc5, 0x55, 0xb7, 0x96, 0xc6, 0xf8, 0x6a, 0xb4, 0xcb, 0xcf,
	0x7c, 0xab, 0x0e, 0xe5, 0x28, 0x56, 0x5a, 0x26, 0xa5, 0x9f, 0x0a, 0xb0, 0x4a, 0x4a, 0xb7, 0x9d,
	0x8b, 0xff, 0x4e, 0xec, 0x1a, 0x2c, 0x66, 0x16, 0x13, 0x4a, 0x99, 0x50, 0x9e, 0xa7, 0x13, 0x5d,
	0x1e, 0x83, 0xcc, 0x4d, 0x02, 0xb9, 0x0b, 0xf5, 0x74, 0x2d, 0x61, 0x54, 0x12, 0x5d, 0x6e, 0xdf,
	0xc1, 0xe6, 0x3c, 0xd5, 0x50, 0xa4, 0x6b, 0x73, 0x57, 0x16, 0xe6, 0x61, 0x84, 0xff, 0xcc, 0x7d,
	0x0d, 0x4a, 0x8e, 0x3f, 0xf0, 0x03, 0x6d, 0x3f, 0x0d, 0x7f, 0x89, 0xc2, 0xff, 0x02, 0xea, 0x9a,
	0x71, 0x8f, 0xbb, 0x2d, 0x64, 0x3c, 0xb3, 0x42, 0x11, 0xd7, 0xd4, 0xf7, 0xd9, 0xc0, 0x7f, 0x2f,
	0x85, 0x10, 0x83, 0xfd, 0x16, 0x16, 0xb3, 0x42, 0xe5, 0xb6, 0x9c, 0x5a, 0x9f, 0x80, 0x29, 0xd7,
	0x87, 0x12, 0x46, 0x05, 0x6a, 0x6d, 0xbc, 0xe6, 0x3a, 0xbf, 0x7d, 0x08, 0xeb, 0x84, 0xdf, 0x0a,
	0x7c, 0xe6, 0x3a, 0x2c, 0x14, 0x2d, 0xf4, 0xfa, 0x7c, 0xba, 0xb0, 0x8c, 0x4d, 0x2b, 0xae, 0x92,
	0xbd, 0x03, 0xf7, 0x7e, 0x03, 0x99, 0x54, 0x5b, 0x3e, 0x4c, 0xfb, 0x8d, 0x3e, 0xe3, 0x6d, 0x14,
	0x71, 0xca, 0x73, 0x6a, 0xb9, 0x0f, 0xa6, 0x93, 0x9e, 0xee, 0x9b, 0xf5, 0x61, 0x77, 0xa0, 0x41,
	0xc0, 0xaf, 0x90, 0xbb, 0xba, 0x24, 0x73, 0x42, 0xcb, 0xd4, 0x28, 0xe8, 0x24, 0xf9, 0xf6, 0x07,
	0x58, 0x4e, 0x8e, 0xd2, 0x4c, 0x60, 0x01, 0x86, 0xd1, 0x40, 0xe8, 0x4a, 0x4a, 0xeb, 0xa0, 0x48,
	0x75, 0x90, 0xf8, 0x30, 0x27, 0xf9, 0xe8, 0xeb, 0x0a, 0xde, 0x76, 0xc4, 0x01, 0xe6, 0xee, 0xae,
	0x29, 0x4d, 0x81, 0x68, 0x1e, 0x41, 0x45, 0x76, 0x67, 0x8f, 0x1a, 0x62, 0x71, 0x12, 0xd5, 0x47,
	0x43, 0x17, 0xa1, 0x74, 0x3a, 0x1d, 0x5d, 0x16, 0xbe, 0x30, 0xe9, 0xc4, 0x2e, 0xc1, 0x7c, 0xac,
	0x49, 0x46, 0x50, 0xbb, 0x6f, 0x40, 0x25, 0x1e, 0xea, 0x72, 0xdd, 0x62, 0x33, 0x76, 0xdb, 0xec,
	0x6f, 0xec, 0x7a, 0xec, 0x36, 0xbb, 0x5f, 0xaf, 0xdb, 0x9d, 0x82, 0xae, 0x11, 0xc3, 0x1f, 0x5d,
	0x8d, 0x50, 0xa7, 0x55, 0x12, 0xb6, 0x53, 0xc2, 0x3c, 0x01, 0x30, 0x6f, 0x04, 0x20, 0x6e, 0x1d,
	0x08, 0x90, 0x04, 0xa0, 0x13, 0xcd, 0x68, 0xff, 0xa1, 0xdc, 0x15, 0xfd, 0xd9, 0xfc, 0x67, 0x03,
	0x6a, 0x19, 0xf3, 0xd3, 0x50, 0xf5, 0xa2, 0x8c, 0xf1, 0xdb, 0xa9, 0xf2, 0xd9, 0xde, 0x81, 0xea,
	0x38, 0xef, 0x91, 0x3f, 0x9b, 0x6f, 0xfb, 0xa5, 0x6e, 0x26, 0x3a, 0xa1, 0x33, 0xe3, 0xfc, 0x48,
	0x82, 0xb3, 0xff, 0x0e, 0x83, 0xa3, 0x28, 0xe0, 0x33, 0xe6, 0xc1, 0x02, 0xe0, 0x78, 0x29, 0x8e,
	0xb3, 0x07, 0x5e, 0x5e, 0x25, 0xcc, 0x11, 0x14, 0x41, 0x33, 0xb9, 0x5b, 0x1c, 0xc6, 0xd5, 0x61,
	0xa3, 0x50, 0x54, 0xf4, 0x80, 0x2a, 0x27, 0xba, 0x16, 0x2b, 0xd4, 0x77, 0x19, 0xef, 0x44, 0xf2,
	0x32, 0x54, 0x9f, 0x36, 0x21, 0x50, 0xcc, 0x2b, 0x93, 0xd2, 0xfb, 0xcd, 0x80, 0xa5, 0x71, 0xfb,
	0x6b, 0xcb, 0x3f, 0x94, 0x89, 0x29, 0x9f, 0x4f, 0x6a, 0x8f, 0x3e, 0xce, 0xbf, 0x3e, 0x9f, 0x08,
	0x6e, 0x43, 0x3d, 0x92, 0xc2, 0x8b, 0x04, 0x5e, 0x57, 0xc0, 0x6a, 0x7a, 0x05, 0x65, 0x26, 0xd5,
	0x8b, 0x2a, 0x94, 0x5a, 0xae, 0xbf, 0xa8, 0x5a, 0x85, 0x4e, 0xf1, 0xe0, 0xce, 0xcf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xcb, 0x56, 0xb0, 0x08, 0x6b, 0x0a, 0x00, 0x00,
}
