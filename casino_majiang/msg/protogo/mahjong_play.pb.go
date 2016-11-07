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

// Ignoring public import of ComposeCardType from base.proto

// Ignoring public import of HuPaiType from base.proto

// Ignoring public import of MJUserGameStatus from base.proto

// Ignoring public import of DeskGameStatus from base.proto

// 开局（接收服务端消息）
type Game_Opening struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
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

// 发牌
type Game_DealCards struct {
	Header           *ProtoHeader  `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	PlayerCard       []*PlayerCard `protobuf:"bytes,2,rep,name=playerCard" json:"playerCard,omitempty"`
	DealerUserId     *uint32       `protobuf:"varint,3,opt,name=dealerUserId" json:"dealerUserId,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
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

func (m *Game_DealCards) GetPlayerCard() []*PlayerCard {
	if m != nil {
		return m.PlayerCard
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
	ExchangeOutCards []*CardInfo  `protobuf:"bytes,2,rep,name=exchangeOutCards" json:"exchangeOutCards,omitempty"`
	UserId           *uint32      `protobuf:"varint,3,opt,name=userId" json:"userId,omitempty"`
	ExchangeNum      *int32       `protobuf:"varint,4,opt,name=exchangeNum" json:"exchangeNum,omitempty"`
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

func (m *Game_ExchangeCards) GetExchangeNum() int32 {
	if m != nil && m.ExchangeNum != nil {
		return *m.ExchangeNum
	}
	return 0
}

//
type Game_AckExchangeCards struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
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

type Game_ExchangeCardsEnd struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	ExchangeInCards  []*CardInfo  `protobuf:"bytes,2,rep,name=exchangeInCards" json:"exchangeInCards,omitempty"`
	ExchangeNum      *int32       `protobuf:"varint,3,opt,name=exchangeNum" json:"exchangeNum,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_ExchangeCardsEnd) Reset()                    { *m = Game_ExchangeCardsEnd{} }
func (m *Game_ExchangeCardsEnd) String() string            { return proto.CompactTextString(m) }
func (*Game_ExchangeCardsEnd) ProtoMessage()               {}
func (*Game_ExchangeCardsEnd) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *Game_ExchangeCardsEnd) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_ExchangeCardsEnd) GetExchangeInCards() []*CardInfo {
	if m != nil {
		return m.ExchangeInCards
	}
	return nil
}

func (m *Game_ExchangeCardsEnd) GetExchangeNum() int32 {
	if m != nil && m.ExchangeNum != nil {
		return *m.ExchangeNum
	}
	return 0
}

// 定缺（和个人玩家ACK）
type Game_DingQue struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Color            *int32       `protobuf:"varint,2,opt,name=color" json:"color,omitempty"`
	UserId           *uint32      `protobuf:"varint,3,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_DingQue) Reset()                    { *m = Game_DingQue{} }
func (m *Game_DingQue) String() string            { return proto.CompactTextString(m) }
func (*Game_DingQue) ProtoMessage()               {}
func (*Game_DingQue) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *Game_DingQue) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
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
func (*DingQueEndBean) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

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
func (*Game_DingQueEnd) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{7} }

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
func (*Game_BroadcastBeginDingQue) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{8} }

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
func (*Game_BroadcastBeginExchange) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{9} }

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
func (*Game_GetInCard) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{10} }

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
func (*Game_SendOutCard) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{11} }

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
func (*Game_AckSendOutCard) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{12} }

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
func (*Game_ActPeng) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{13} }

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
func (*Game_AckActPeng) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{14} }

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
func (*Game_ActGang) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{15} }

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
func (*Game_AckActGang) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{16} }

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
func (*Game_ActHu) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{17} }

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
	IsZiMo           *bool        `protobuf:"varint,6,opt,name=isZiMo" json:"isZiMo,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckActHu) Reset()                    { *m = Game_AckActHu{} }
func (m *Game_AckActHu) String() string            { return proto.CompactTextString(m) }
func (*Game_AckActHu) ProtoMessage()               {}
func (*Game_AckActHu) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{18} }

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

func (m *Game_AckActHu) GetIsZiMo() bool {
	if m != nil && m.IsZiMo != nil {
		return *m.IsZiMo
	}
	return false
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
func (*Game_ActGuo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{19} }

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
func (*Game_AckActGuo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{20} }

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
	GangCards        []*CardInfo  `protobuf:"bytes,9,rep,name=gangCards" json:"gangCards,omitempty"`
	Time             *int32       `protobuf:"varint,10,opt,name=time" json:"time,omitempty"`
	PaiCount         *int32       `protobuf:"varint,11,opt,name=paiCount" json:"paiCount,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_OverTurn) Reset()                    { *m = Game_OverTurn{} }
func (m *Game_OverTurn) String() string            { return proto.CompactTextString(m) }
func (*Game_OverTurn) ProtoMessage()               {}
func (*Game_OverTurn) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{21} }

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

func (m *Game_OverTurn) GetGangCards() []*CardInfo {
	if m != nil {
		return m.GangCards
	}
	return nil
}

func (m *Game_OverTurn) GetTime() int32 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func (m *Game_OverTurn) GetPaiCount() int32 {
	if m != nil && m.PaiCount != nil {
		return *m.PaiCount
	}
	return 0
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
func (*Game_SendGameInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{22} }

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
	proto.RegisterType((*Game_ExchangeCardsEnd)(nil), "mjproto.Game_ExchangeCardsEnd")
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
	// 730 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x56, 0xcf, 0x4f, 0xdb, 0x4c,
	0x10, 0xfd, 0x9c, 0x04, 0x08, 0x13, 0x02, 0xc1, 0x80, 0x3e, 0x8b, 0x1e, 0x4a, 0x5d, 0xaa, 0x56,
	0x45, 0x42, 0x55, 0xd5, 0x6b, 0x0f, 0x04, 0x68, 0x92, 0x43, 0x0b, 0x08, 0xaa, 0x4a, 0x3d, 0x14,
	0x6d, 0xed, 0xc1, 0x31, 0x24, 0xeb, 0xc8, 0xf6, 0xb6, 0x20, 0x55, 0x95, 0x7a, 0xe8, 0xad, 0xb7,
	0x9e, 0x7b, 0xec, 0xff, 0xd9, 0xdd, 0xf1, 0xda, 0x98, 0x1f, 0x29, 0x8e, 0x7b, 0x89, 0xe2, 0xf5,
	0xce, 0x7b, 0xf3, 0xe6, 0xcd, 0x8e, 0x17, 0xcc, 0x21, 0xeb, 0x9f, 0x06, 0xdc, 0x3b, 0x1e, 0x0d,
	0xd8, 0xc5, 0xe6, 0x28, 0x0c, 0xe2, 0xc0, 0x9c, 0x19, 0x9e, 0xd2, 0x9f, 0x55, 0xf8, 0xc8, 0x22,
	0x4c, 0x16, 0xed, 0x17, 0x30, 0xe7, 0xb1, 0x21, 0x1e, 0xef, 0x8d, 0x90, 0xfb, 0xdc, 0x33, 0xd7,
	0x61, 0xba, 0x8f, 0xcc, 0xc5, 0xd0, 0x32, 0xd6, 0x8c, 0x27, 0x8d, 0xe7, 0xcb, 0x9b, 0x3a, 0x6a,
	0x73, 0x5f, 0xfd, 0x76, 0xe9, 0x9d, 0x2d, 0x60, 0x9e, 0xa2, 0x76, 0x90, 0x0d, 0xb6, 0x59, 0xe8,
	0x46, 0xc5, 0xe2, 0xcc, 0xc7, 0x00, 0x2a, 0x21, 0x0c, 0x55, 0x90, 0x55, 0x59, 0xab, 0xca, 0x9d,
	0x4b, 0x97, 0x3b, 0xb3, 0x57, 0xe6, 0x32, 0xcc, 0xb9, 0x12, 0x1b, 0xc3, 0xb7, 0x11, 0x86, 0x3d,
	0xd7, 0xaa, 0x4a, 0xd0, 0xa6, 0xfd, 0xc3, 0x00, 0x93, 0x78, 0x77, 0xcf, 0x9d, 0x3e, 0xe3, 0x1e,
	0x4e, 0xc2, 0xbd, 0x01, 0x2d, 0xd4, 0x61, 0x7b, 0x22, 0xa6, 0x48, 0x9d, 0xc1, 0x62, 0xb6, 0x5f,
	0xad, 0xf6, 0xf8, 0x49, 0x60, 0xce, 0xc3, 0xb4, 0xc8, 0x31, 0x9b, 0x4b, 0xd0, 0x48, 0x83, 0xdf,
	0x88, 0xa1, 0x55, 0x93, 0x8b, 0x53, 0xf6, 0x4b, 0x58, 0xa1, 0x6c, 0xb6, 0x9c, 0xb3, 0x12, 0x09,
	0xd9, 0x5f, 0x61, 0xa5, 0x73, 0x43, 0xcc, 0x2e, 0x77, 0x0b, 0xea, 0x79, 0x0a, 0x0b, 0x69, 0x4a,
	0x3d, 0x7e, 0x87, 0x9c, 0x6b, 0xe9, 0x57, 0x29, 0xfd, 0x43, 0x6d, 0xfd, 0x8e, 0xf4, 0xfd, 0x40,
	0x60, 0x41, 0xda, 0x26, 0x4c, 0x39, 0xc1, 0x20, 0x08, 0x25, 0x99, 0x04, 0xb9, 0x5e, 0x28, 0xfb,
	0x19, 0xcc, 0x6b, 0x3c, 0xa9, 0xa4, 0x8d, 0x8c, 0xe7, 0x76, 0x18, 0x54, 0x4a, 0xf9, 0x7c, 0x32,
	0x08, 0x3e, 0xa3, 0x46, 0xb0, 0x3f, 0xc0, 0x42, 0x3e, 0x8d, 0xe2, 0x05, 0x78, 0x04, 0x35, 0xb9,
	0x3f, 0x55, 0xfd, 0x7f, 0xb6, 0xe7, 0x2a, 0xbf, 0x7d, 0x00, 0xab, 0x84, 0xdf, 0x0e, 0x03, 0xe6,
	0x3a, 0x2c, 0x8a, 0xdb, 0xe8, 0xf9, 0xbc, 0xb4, 0xe8, 0xaa, 0x4c, 0x79, 0x1b, 0xee, 0xdd, 0x02,
	0x99, 0xfa, 0x58, 0xd0, 0xfe, 0x77, 0xfa, 0x0c, 0x75, 0x30, 0x4e, 0xfc, 0x2b, 0x98, 0xcb, 0x7d,
	0xa8, 0x39, 0xc9, 0xe9, 0x31, 0x6e, 0x35, 0xdb, 0xee, 0x42, 0x8b, 0x80, 0x0f, 0x91, 0xbb, 0xba,
	0xd3, 0x0b, 0x42, 0x4b, 0x6b, 0x14, 0x74, 0xcf, 0xd5, 0xd6, 0x7c, 0x81, 0xa5, 0xb4, 0xc1, 0x4b,
	0x81, 0x85, 0x18, 0x89, 0x41, 0x7c, 0x7b, 0xa7, 0x64, 0x3a, 0x6a, 0xe3, 0x74, 0xf8, 0xba, 0x3f,
	0xb7, 0x9c, 0x78, 0x1f, 0x8b, 0x8e, 0xa6, 0x1c, 0x4d, 0x85, 0x68, 0x1e, 0x42, 0x5d, 0x8e, 0x36,
	0x8f, 0x06, 0x4e, 0x75, 0x1c, 0xd5, 0x37, 0x43, 0x37, 0xa1, 0x54, 0x3a, 0x19, 0x5d, 0x1e, 0x7e,
	0xec, 0xf1, 0x5b, 0x84, 0xd9, 0x24, 0x27, 0x59, 0x41, 0xad, 0xbe, 0x05, 0xf5, 0x64, 0xa9, 0xc7,
	0xa9, 0x02, 0xcd, 0xbc, 0xdc, 0x0e, 0xfb, 0x17, 0xb9, 0x1e, 0xbb, 0x4b, 0xee, 0xcf, 0xab, 0x72,
	0x27, 0xa0, 0x6b, 0x25, 0xf0, 0x47, 0x17, 0x23, 0xd4, 0xb6, 0x4a, 0xc2, 0xce, 0x25, 0x61, 0x91,
	0x02, 0xd4, 0x6e, 0x14, 0x60, 0x8a, 0x0a, 0x80, 0x00, 0x69, 0x01, 0xba, 0xa2, 0xa4, 0xfc, 0x07,
	0x32, 0x4a, 0xfc, 0x5d, 0xfc, 0x2f, 0x03, 0x9a, 0x39, 0xf1, 0x93, 0x50, 0xf5, 0x45, 0x4e, 0xf8,
	0xdd, 0x54, 0x85, 0x64, 0x2b, 0x5c, 0x3f, 0x7a, 0xef, 0xbf, 0x0e, 0xac, 0x69, 0xf9, 0x5c, 0x97,
	0xc3, 0xa5, 0x91, 0xf5, 0x81, 0x08, 0xca, 0xd5, 0xc1, 0x7e, 0xa5, 0x87, 0x8b, 0x36, 0xb8, 0x34,
	0xce, 0xf7, 0x8a, 0x2e, 0xd6, 0xde, 0x27, 0x0c, 0x8f, 0x44, 0xc8, 0x4b, 0xfa, 0x62, 0x02, 0x70,
	0x3c, 0x8f, 0xf3, 0x5f, 0x73, 0x73, 0x01, 0x66, 0x98, 0x13, 0x53, 0x45, 0xe9, 0x7b, 0xaa, 0x16,
	0x1c, 0xc6, 0xd5, 0xe1, 0xa3, 0xd2, 0xd4, 0xf5, 0x82, 0x6a, 0xaf, 0xa4, 0x36, 0x34, 0x87, 0x19,
	0xef, 0x0a, 0x6b, 0x86, 0x1e, 0x6d, 0x42, 0x20, 0x0f, 0xea, 0xe3, 0x3c, 0x58, 0x87, 0xd9, 0xf4,
	0x40, 0x44, 0xd6, 0xec, 0xb8, 0x06, 0x9d, 0x83, 0x5a, 0xec, 0x0f, 0xd1, 0x02, 0x4a, 0x44, 0x9a,
	0x34, 0x62, 0xfe, 0x76, 0x20, 0x78, 0x6c, 0x35, 0x68, 0x12, 0xfe, 0x36, 0x60, 0x31, 0x1b, 0xaa,
	0xea, 0xab, 0xad, 0xb1, 0x27, 0xba, 0xf4, 0xa8, 0x98, 0x31, 0x97, 0x1e, 0x82, 0xdb, 0x50, 0x97,
	0x9e, 0xe8, 0x2c, 0x85, 0xd7, 0x7d, 0xb5, 0x72, 0xf9, 0x61, 0xcb, 0xbd, 0x54, 0x37, 0xa4, 0x48,
	0xe6, 0x92, 0xdd, 0x90, 0xa8, 0xbd, 0xda, 0x95, 0x6e, 0x75, 0xff, 0xbf, 0x3f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xdf, 0xac, 0xc3, 0xce, 0xfe, 0x09, 0x00, 0x00,
}
