package majiang

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

// 麻将牌的结构
type MJPai struct {
	Index            *int32  `protobuf:"varint,1,opt,name=index" json:"index,omitempty"`
	Flower           *int32  `protobuf:"varint,2,opt,name=flower" json:"flower,omitempty"`
	Value            *int32  `protobuf:"varint,3,opt,name=value" json:"value,omitempty"`
	Des              *string `protobuf:"bytes,4,opt,name=des" json:"des,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MJPai) Reset()                    { *m = MJPai{} }
func (m *MJPai) String() string            { return proto.CompactTextString(m) }
func (*MJPai) ProtoMessage()               {}
func (*MJPai) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MJPai) GetIndex() int32 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *MJPai) GetFlower() int32 {
	if m != nil && m.Flower != nil {
		return *m.Flower
	}
	return 0
}

func (m *MJPai) GetValue() int32 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

func (m *MJPai) GetDes() string {
	if m != nil && m.Des != nil {
		return *m.Des
	}
	return ""
}

// 手里的一副牌
type MJHandPai struct {
	Pais      []*MJPai `protobuf:"bytes,1,rep,name=Pais" json:"Pais,omitempty"`
	PengPais  []*MJPai `protobuf:"bytes,2,rep,name=PengPais" json:"PengPais,omitempty"`
	GangPais  []*MJPai `protobuf:"bytes,3,rep,name=GangPais" json:"GangPais,omitempty"`
	HuPais    []*MJPai `protobuf:"bytes,4,rep,name=HuPais" json:"HuPais,omitempty"`
	InPai     *MJPai   `protobuf:"bytes,5,opt,name=inPai" json:"inPai,omitempty"`
	QueFlower *int32   `protobuf:"varint,6,opt,name=queFlower" json:"queFlower,omitempty"`
	OutPais   []*MJPai `protobuf:"bytes,7,rep,name=OutPais" json:"OutPais,omitempty"`
	ChiPais   []*MJPai `protobuf:"bytes,8,rep,name=ChiPais" json:"ChiPais,omitempty"`
	InPai2    *MJPai   `protobuf:"bytes,9,opt,name=inPai2" json:"inPai2,omitempty"`
	gangInfos []*GangPaiInfo //杠牌的信息
}

func (m *MJHandPai) Reset()                    { *m = MJHandPai{} }
func (m *MJHandPai) String() string            { return proto.CompactTextString(m) }
func (*MJHandPai) ProtoMessage()               {}
func (*MJHandPai) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MJHandPai) GetPais() []*MJPai {
	if m != nil {
		return m.Pais
	}
	return nil
}

func (m *MJHandPai) GetPengPais() []*MJPai {
	if m != nil {
		return m.PengPais
	}
	return nil
}

func (m *MJHandPai) GetGangPais() []*MJPai {
	if m != nil {
		return m.GangPais
	}
	return nil
}

func (m *MJHandPai) GetHuPais() []*MJPai {
	if m != nil {
		return m.HuPais
	}
	return nil
}

func (m *MJHandPai) GetInPai() *MJPai {
	if m != nil {
		return m.InPai
	}
	return nil
}

func (m *MJHandPai) GetQueFlower() int32 {
	if m != nil && m.QueFlower != nil {
		return *m.QueFlower
	}
	return 0
}

func (m *MJHandPai) GetOutPais() []*MJPai {
	if m != nil {
		return m.OutPais
	}
	return nil
}

func (m *MJHandPai) GetChiPais() []*MJPai {
	if m != nil {
		return m.ChiPais
	}
	return nil
}

func (m *MJHandPai) GetInPai2() *MJPai {
	if m != nil {
		return m.InPai2
	}
	return nil
}

type HuPaiInfo struct {
	SendUserId       *uint32 `protobuf:"varint,1,opt,name=sendUserId" json:"sendUserId,omitempty"`
	GetUserId        *uint32 `protobuf:"varint,2,opt,name=getUserId" json:"getUserId,omitempty"`
	ByWho            *int32  `protobuf:"varint,3,opt,name=byWho" json:"byWho,omitempty"`
	HuType           *int32  `protobuf:"varint,4,opt,name=huType" json:"huType,omitempty"`
	PaiType          []int32 `protobuf:"varint,5,rep,name=paiType" json:"paiType,omitempty"`
	HuDesc           *string `protobuf:"bytes,6,opt,name=huDesc" json:"huDesc,omitempty"`
	Pai              *MJPai  `protobuf:"bytes,7,opt,name=pai" json:"pai,omitempty"`
	Fan              *int32  `protobuf:"varint,8,opt,name=fan" json:"fan,omitempty"`
	Score            *int64  `protobuf:"varint,9,opt,name=score" json:"score,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *HuPaiInfo) Reset()                    { *m = HuPaiInfo{} }
func (m *HuPaiInfo) String() string            { return proto.CompactTextString(m) }
func (*HuPaiInfo) ProtoMessage()               {}
func (*HuPaiInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *HuPaiInfo) GetSendUserId() uint32 {
	if m != nil && m.SendUserId != nil {
		return *m.SendUserId
	}
	return 0
}

func (m *HuPaiInfo) GetGetUserId() uint32 {
	if m != nil && m.GetUserId != nil {
		return *m.GetUserId
	}
	return 0
}

func (m *HuPaiInfo) GetByWho() int32 {
	if m != nil && m.ByWho != nil {
		return *m.ByWho
	}
	return 0
}

func (m *HuPaiInfo) GetHuType() int32 {
	if m != nil && m.HuType != nil {
		return *m.HuType
	}
	return 0
}

func (m *HuPaiInfo) GetPaiType() []int32 {
	if m != nil {
		return m.PaiType
	}
	return nil
}

func (m *HuPaiInfo) GetHuDesc() string {
	if m != nil && m.HuDesc != nil {
		return *m.HuDesc
	}
	return ""
}

func (m *HuPaiInfo) GetPai() *MJPai {
	if m != nil {
		return m.Pai
	}
	return nil
}

func (m *HuPaiInfo) GetFan() int32 {
	if m != nil && m.Fan != nil {
		return *m.Fan
	}
	return 0
}

func (m *HuPaiInfo) GetScore() int64 {
	if m != nil && m.Score != nil {
		return *m.Score
	}
	return 0
}

type GangPaiInfo struct {
	SendUserId       *uint32 `protobuf:"varint,1,opt,name=sendUserId" json:"sendUserId,omitempty"`
	GetUserId        *uint32 `protobuf:"varint,2,opt,name=getUserId" json:"getUserId,omitempty"`
	ByWho            *int32  `protobuf:"varint,3,opt,name=byWho" json:"byWho,omitempty"`
	GangType         *int32  `protobuf:"varint,4,opt,name=gangType" json:"gangType,omitempty"`
	Pai              *MJPai  `protobuf:"bytes,5,opt,name=pai" json:"pai,omitempty"`
	Transfer         *int32  `protobuf:"varint,6,opt,name=transfer" json:"transfer,omitempty"`
	Bu               *bool   `protobuf:"varint,7,opt,name=bu" json:"bu,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GangPaiInfo) Reset()                    { *m = GangPaiInfo{} }
func (m *GangPaiInfo) String() string            { return proto.CompactTextString(m) }
func (*GangPaiInfo) ProtoMessage()               {}
func (*GangPaiInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GangPaiInfo) GetSendUserId() uint32 {
	if m != nil && m.SendUserId != nil {
		return *m.SendUserId
	}
	return 0
}

func (m *GangPaiInfo) GetGetUserId() uint32 {
	if m != nil && m.GetUserId != nil {
		return *m.GetUserId
	}
	return 0
}

func (m *GangPaiInfo) GetByWho() int32 {
	if m != nil && m.ByWho != nil {
		return *m.ByWho
	}
	return 0
}

func (m *GangPaiInfo) GetGangType() int32 {
	if m != nil && m.GangType != nil {
		return *m.GangType
	}
	return 0
}

func (m *GangPaiInfo) GetPai() *MJPai {
	if m != nil {
		return m.Pai
	}
	return nil
}

func (m *GangPaiInfo) GetTransfer() int32 {
	if m != nil && m.Transfer != nil {
		return *m.Transfer
	}
	return 0
}

func (m *GangPaiInfo) GetBu() bool {
	if m != nil && m.Bu != nil {
		return *m.Bu
	}
	return false
}

type GuoHuInfo struct {
	SendUserId       *uint32 `protobuf:"varint,1,opt,name=sendUserId" json:"sendUserId,omitempty"`
	GetUserId        *uint32 `protobuf:"varint,2,opt,name=getUserId" json:"getUserId,omitempty"`
	Pai              *MJPai  `protobuf:"bytes,3,opt,name=pai" json:"pai,omitempty"`
	FanShu           *int32  `protobuf:"varint,4,opt,name=fanShu" json:"fanShu,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GuoHuInfo) Reset()                    { *m = GuoHuInfo{} }
func (m *GuoHuInfo) String() string            { return proto.CompactTextString(m) }
func (*GuoHuInfo) ProtoMessage()               {}
func (*GuoHuInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GuoHuInfo) GetSendUserId() uint32 {
	if m != nil && m.SendUserId != nil {
		return *m.SendUserId
	}
	return 0
}

func (m *GuoHuInfo) GetGetUserId() uint32 {
	if m != nil && m.GetUserId != nil {
		return *m.GetUserId
	}
	return 0
}

func (m *GuoHuInfo) GetPai() *MJPai {
	if m != nil {
		return m.Pai
	}
	return nil
}

func (m *GuoHuInfo) GetFanShu() int32 {
	if m != nil && m.FanShu != nil {
		return *m.FanShu
	}
	return 0
}

// 一个玩家的游戏信息
type PlayerGameData struct {
	HandPai          *MJHandPai     `protobuf:"bytes,1,opt,name=handPai" json:"handPai,omitempty"`
	HuInfo           []*HuPaiInfo   `protobuf:"bytes,2,rep,name=huInfo" json:"huInfo,omitempty"`
	DianHuInfo       []*HuPaiInfo   `protobuf:"bytes,3,rep,name=dianHuInfo" json:"dianHuInfo,omitempty"`
	GangInfo         []*GangPaiInfo `protobuf:"bytes,4,rep,name=gangInfo" json:"gangInfo,omitempty"`
	DianGangInfo     []*GangPaiInfo `protobuf:"bytes,5,rep,name=dianGangInfo" json:"dianGangInfo,omitempty"`
	GangTransfer     []*GangPaiInfo `protobuf:"bytes,6,rep,name=gangTransfer" json:"gangTransfer,omitempty"`
	TotalScore       *int32         `protobuf:"varint,7,opt,name=totalScore" json:"totalScore,omitempty"`
	GuoHuInfo        []*GuoHuInfo   `protobuf:"bytes,8,rep,name=guoHuInfo" json:"guoHuInfo,omitempty"`
	ExchangeCardsOut []*MJPai       `protobuf:"bytes,9,rep,name=exchangeCardsOut" json:"exchangeCardsOut,omitempty"`
	ExchangeCardsIn  []*MJPai       `protobuf:"bytes,10,rep,name=exchangeCardsIn" json:"exchangeCardsIn,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *PlayerGameData) Reset()                    { *m = PlayerGameData{} }
func (m *PlayerGameData) String() string            { return proto.CompactTextString(m) }
func (*PlayerGameData) ProtoMessage()               {}
func (*PlayerGameData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *PlayerGameData) GetHandPai() *MJHandPai {
	if m != nil {
		return m.HandPai
	}
	return nil
}

func (m *PlayerGameData) GetHuInfo() []*HuPaiInfo {
	if m != nil {
		return m.HuInfo
	}
	return nil
}

func (m *PlayerGameData) GetDianHuInfo() []*HuPaiInfo {
	if m != nil {
		return m.DianHuInfo
	}
	return nil
}

func (m *PlayerGameData) GetGangInfo() []*GangPaiInfo {
	if m != nil {
		return m.GangInfo
	}
	return nil
}

func (m *PlayerGameData) GetDianGangInfo() []*GangPaiInfo {
	if m != nil {
		return m.DianGangInfo
	}
	return nil
}

func (m *PlayerGameData) GetGangTransfer() []*GangPaiInfo {
	if m != nil {
		return m.GangTransfer
	}
	return nil
}

func (m *PlayerGameData) GetTotalScore() int32 {
	if m != nil && m.TotalScore != nil {
		return *m.TotalScore
	}
	return 0
}

func (m *PlayerGameData) GetGuoHuInfo() []*GuoHuInfo {
	if m != nil {
		return m.GuoHuInfo
	}
	return nil
}

func (m *PlayerGameData) GetExchangeCardsOut() []*MJPai {
	if m != nil {
		return m.ExchangeCardsOut
	}
	return nil
}

func (m *PlayerGameData) GetExchangeCardsIn() []*MJPai {
	if m != nil {
		return m.ExchangeCardsIn
	}
	return nil
}

// 需要确认的事件
type CheckBean struct {
	UserId           *uint32  `protobuf:"varint,1,opt,name=userId" json:"userId,omitempty"`
	CanPeng          *bool    `protobuf:"varint,2,opt,name=CanPeng" json:"CanPeng,omitempty"`
	CanGang          *bool    `protobuf:"varint,3,opt,name=CanGang" json:"CanGang,omitempty"`
	CanHu            *bool    `protobuf:"varint,4,opt,name=CanHu" json:"CanHu,omitempty"`
	CheckStatus      *int32   `protobuf:"varint,5,opt,name=CheckStatus" json:"CheckStatus,omitempty"`
	CanChi           *bool    `protobuf:"varint,6,opt,name=canChi" json:"canChi,omitempty"`
	CanGuo           *bool    `protobuf:"varint,7,opt,name=canGuo" json:"canGuo,omitempty"`
	ChiCards         []*MJPai `protobuf:"bytes,8,rep,name=chiCards" json:"chiCards,omitempty"`
	CheckPai         *MJPai   `protobuf:"bytes,9,opt,name=checkPai" json:"checkPai,omitempty"`
	CanBu            *bool    `protobuf:"varint,10,opt,name=CanBu" json:"CanBu,omitempty"`
	GangCards        []*MJPai `protobuf:"bytes,11,rep,name=gangCards" json:"gangCards,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *CheckBean) Reset()                    { *m = CheckBean{} }
func (m *CheckBean) String() string            { return proto.CompactTextString(m) }
func (*CheckBean) ProtoMessage()               {}
func (*CheckBean) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CheckBean) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *CheckBean) GetCanPeng() bool {
	if m != nil && m.CanPeng != nil {
		return *m.CanPeng
	}
	return false
}

func (m *CheckBean) GetCanGang() bool {
	if m != nil && m.CanGang != nil {
		return *m.CanGang
	}
	return false
}

func (m *CheckBean) GetCanHu() bool {
	if m != nil && m.CanHu != nil {
		return *m.CanHu
	}
	return false
}

func (m *CheckBean) GetCheckStatus() int32 {
	if m != nil && m.CheckStatus != nil {
		return *m.CheckStatus
	}
	return 0
}

func (m *CheckBean) GetCanChi() bool {
	if m != nil && m.CanChi != nil {
		return *m.CanChi
	}
	return false
}

func (m *CheckBean) GetCanGuo() bool {
	if m != nil && m.CanGuo != nil {
		return *m.CanGuo
	}
	return false
}

func (m *CheckBean) GetChiCards() []*MJPai {
	if m != nil {
		return m.ChiCards
	}
	return nil
}

func (m *CheckBean) GetCheckPai() *MJPai {
	if m != nil {
		return m.CheckPai
	}
	return nil
}

func (m *CheckBean) GetCanBu() bool {
	if m != nil && m.CanBu != nil {
		return *m.CanBu
	}
	return false
}

func (m *CheckBean) GetGangCards() []*MJPai {
	if m != nil {
		return m.GangCards
	}
	return nil
}

type CheckCase struct {
	CheckB           []*CheckBean `protobuf:"bytes,1,rep,name=CheckB" json:"CheckB,omitempty"`
	CheckStatus      *int32       `protobuf:"varint,2,opt,name=CheckStatus" json:"CheckStatus,omitempty"`
	CheckMJPai       *MJPai       `protobuf:"bytes,3,opt,name=CheckMJPai" json:"CheckMJPai,omitempty"`
	UserIdOut        *uint32      `protobuf:"varint,4,opt,name=UserIdOut" json:"UserIdOut,omitempty"`
	PreOutGangInfo   *GangPaiInfo `protobuf:"bytes,5,opt,name=PreOutGangInfo" json:"PreOutGangInfo,omitempty"`
	CheckMJPai2      *MJPai       `protobuf:"bytes,6,opt,name=CheckMJPai2" json:"CheckMJPai2,omitempty"`
	DianPaoCount     *int32       `protobuf:"varint,7,opt,name=dianPaoCount" json:"dianPaoCount,omitempty"`
	QiangGang        *bool        `protobuf:"varint,8,opt,name=qiangGang" json:"qiangGang,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *CheckCase) Reset()                    { *m = CheckCase{} }
func (m *CheckCase) String() string            { return proto.CompactTextString(m) }
func (*CheckCase) ProtoMessage()               {}
func (*CheckCase) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *CheckCase) GetCheckB() []*CheckBean {
	if m != nil {
		return m.CheckB
	}
	return nil
}

func (m *CheckCase) GetCheckStatus() int32 {
	if m != nil && m.CheckStatus != nil {
		return *m.CheckStatus
	}
	return 0
}

func (m *CheckCase) GetCheckMJPai() *MJPai {
	if m != nil {
		return m.CheckMJPai
	}
	return nil
}

func (m *CheckCase) GetUserIdOut() uint32 {
	if m != nil && m.UserIdOut != nil {
		return *m.UserIdOut
	}
	return 0
}

func (m *CheckCase) GetPreOutGangInfo() *GangPaiInfo {
	if m != nil {
		return m.PreOutGangInfo
	}
	return nil
}

func (m *CheckCase) GetCheckMJPai2() *MJPai {
	if m != nil {
		return m.CheckMJPai2
	}
	return nil
}

func (m *CheckCase) GetDianPaoCount() int32 {
	if m != nil && m.DianPaoCount != nil {
		return *m.DianPaoCount
	}
	return 0
}

func (m *CheckCase) GetQiangGang() bool {
	if m != nil && m.QiangGang != nil {
		return *m.QiangGang
	}
	return false
}

// 麻将牌的结构
type PMjDesk struct {
	MjRoomType       *int32     `protobuf:"varint,1,opt,name=mjRoomType" json:"mjRoomType,omitempty"`
	CapMax           *int64     `protobuf:"varint,2,opt,name=capMax" json:"capMax,omitempty"`
	CardsNum         *int32     `protobuf:"varint,3,opt,name=cardsNum" json:"cardsNum,omitempty"`
	Settlement       *int32     `protobuf:"varint,4,opt,name=settlement" json:"settlement,omitempty"`
	ZiMoRadio        *int32     `protobuf:"varint,5,opt,name=ziMoRadio" json:"ziMoRadio,omitempty"`
	DianGangHuaRadio *int32     `protobuf:"varint,6,opt,name=dianGangHuaRadio" json:"dianGangHuaRadio,omitempty"`
	OthersCheckBox   []int32    `protobuf:"varint,7,rep,name=othersCheckBox" json:"othersCheckBox,omitempty"`
	HuRadio          *int32     `protobuf:"varint,8,opt,name=huRadio" json:"huRadio,omitempty"`
	AllMJPai         []*MJPai   `protobuf:"bytes,9,rep,name=AllMJPai" json:"AllMJPai,omitempty"`
	DeskMJPai        []*MJPai   `protobuf:"bytes,10,rep,name=DeskMJPai" json:"DeskMJPai,omitempty"`
	MJPaiCursor      *int32     `protobuf:"varint,11,opt,name=MJPaiCursor" json:"MJPaiCursor,omitempty"`
	TotalPlayCount   *int32     `protobuf:"varint,12,opt,name=TotalPlayCount" json:"TotalPlayCount,omitempty"`
	CurrPlayCount    *int32     `protobuf:"varint,13,opt,name=CurrPlayCount" json:"CurrPlayCount,omitempty"`
	CheckCase        *CheckCase `protobuf:"bytes,14,opt,name=CheckCase" json:"CheckCase,omitempty"`
	ActiveUser       *uint32    `protobuf:"varint,15,opt,name=ActiveUser" json:"ActiveUser,omitempty"`
	ActUser          *uint32    `protobuf:"varint,16,opt,name=ActUser" json:"ActUser,omitempty"`
	ActType          *int32     `protobuf:"varint,17,opt,name=ActType" json:"ActType,omitempty"`
	NextBanker       *uint32    `protobuf:"varint,18,opt,name=nextBanker" json:"nextBanker,omitempty"`
	FangCountLimit   *int32     `protobuf:"varint,19,opt,name=fangCountLimit" json:"fangCountLimit,omitempty"`
	RoomType         *int32     `protobuf:"varint,20,opt,name=roomType" json:"roomType,omitempty"`
	CoinLimit        *int64     `protobuf:"varint,21,opt,name=coinLimit" json:"coinLimit,omitempty"`
	RoomLevel        *int32     `protobuf:"varint,22,opt,name=roomLevel" json:"roomLevel,omitempty"`
	ApplyDis         *bool      `protobuf:"varint,23,opt,name=applyDis" json:"applyDis,omitempty"`
	CoinLimitUL      *int64     `protobuf:"varint,24,opt,name=coinLimitUL" json:"coinLimitUL,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *PMjDesk) Reset()                    { *m = PMjDesk{} }
func (m *PMjDesk) String() string            { return proto.CompactTextString(m) }
func (*PMjDesk) ProtoMessage()               {}
func (*PMjDesk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *PMjDesk) GetMjRoomType() int32 {
	if m != nil && m.MjRoomType != nil {
		return *m.MjRoomType
	}
	return 0
}

func (m *PMjDesk) GetCapMax() int64 {
	if m != nil && m.CapMax != nil {
		return *m.CapMax
	}
	return 0
}

func (m *PMjDesk) GetCardsNum() int32 {
	if m != nil && m.CardsNum != nil {
		return *m.CardsNum
	}
	return 0
}

func (m *PMjDesk) GetSettlement() int32 {
	if m != nil && m.Settlement != nil {
		return *m.Settlement
	}
	return 0
}

func (m *PMjDesk) GetZiMoRadio() int32 {
	if m != nil && m.ZiMoRadio != nil {
		return *m.ZiMoRadio
	}
	return 0
}

func (m *PMjDesk) GetDianGangHuaRadio() int32 {
	if m != nil && m.DianGangHuaRadio != nil {
		return *m.DianGangHuaRadio
	}
	return 0
}

func (m *PMjDesk) GetOthersCheckBox() []int32 {
	if m != nil {
		return m.OthersCheckBox
	}
	return nil
}

func (m *PMjDesk) GetHuRadio() int32 {
	if m != nil && m.HuRadio != nil {
		return *m.HuRadio
	}
	return 0
}

func (m *PMjDesk) GetAllMJPai() []*MJPai {
	if m != nil {
		return m.AllMJPai
	}
	return nil
}

func (m *PMjDesk) GetDeskMJPai() []*MJPai {
	if m != nil {
		return m.DeskMJPai
	}
	return nil
}

func (m *PMjDesk) GetMJPaiCursor() int32 {
	if m != nil && m.MJPaiCursor != nil {
		return *m.MJPaiCursor
	}
	return 0
}

func (m *PMjDesk) GetTotalPlayCount() int32 {
	if m != nil && m.TotalPlayCount != nil {
		return *m.TotalPlayCount
	}
	return 0
}

func (m *PMjDesk) GetCurrPlayCount() int32 {
	if m != nil && m.CurrPlayCount != nil {
		return *m.CurrPlayCount
	}
	return 0
}

func (m *PMjDesk) GetCheckCase() *CheckCase {
	if m != nil {
		return m.CheckCase
	}
	return nil
}

func (m *PMjDesk) GetActiveUser() uint32 {
	if m != nil && m.ActiveUser != nil {
		return *m.ActiveUser
	}
	return 0
}

func (m *PMjDesk) GetActUser() uint32 {
	if m != nil && m.ActUser != nil {
		return *m.ActUser
	}
	return 0
}

func (m *PMjDesk) GetActType() int32 {
	if m != nil && m.ActType != nil {
		return *m.ActType
	}
	return 0
}

func (m *PMjDesk) GetNextBanker() uint32 {
	if m != nil && m.NextBanker != nil {
		return *m.NextBanker
	}
	return 0
}

func (m *PMjDesk) GetFangCountLimit() int32 {
	if m != nil && m.FangCountLimit != nil {
		return *m.FangCountLimit
	}
	return 0
}

func (m *PMjDesk) GetRoomType() int32 {
	if m != nil && m.RoomType != nil {
		return *m.RoomType
	}
	return 0
}

func (m *PMjDesk) GetCoinLimit() int64 {
	if m != nil && m.CoinLimit != nil {
		return *m.CoinLimit
	}
	return 0
}

func (m *PMjDesk) GetRoomLevel() int32 {
	if m != nil && m.RoomLevel != nil {
		return *m.RoomLevel
	}
	return 0
}

func (m *PMjDesk) GetApplyDis() bool {
	if m != nil && m.ApplyDis != nil {
		return *m.ApplyDis
	}
	return false
}

func (m *PMjDesk) GetCoinLimitUL() int64 {
	if m != nil && m.CoinLimitUL != nil {
		return *m.CoinLimitUL
	}
	return 0
}

type BillBean struct {
	UserId           *uint32 `protobuf:"varint,1,opt,name=UserId" json:"UserId,omitempty"`
	OutUserId        *uint32 `protobuf:"varint,2,opt,name=OutUserId" json:"OutUserId,omitempty"`
	Type             *int32  `protobuf:"varint,3,opt,name=type" json:"type,omitempty"`
	Des              *string `protobuf:"bytes,4,opt,name=des" json:"des,omitempty"`
	Amount           *int64  `protobuf:"varint,5,opt,name=amount" json:"amount,omitempty"`
	Pai              *MJPai  `protobuf:"bytes,6,opt,name=pai" json:"pai,omitempty"`
	IsBird           *bool   `protobuf:"varint,7,opt,name=isBird" json:"isBird,omitempty"`
	IsQiShouHu       *bool   `protobuf:"varint,8,opt,name=isQiShouHu" json:"isQiShouHu,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *BillBean) Reset()                    { *m = BillBean{} }
func (m *BillBean) String() string            { return proto.CompactTextString(m) }
func (*BillBean) ProtoMessage()               {}
func (*BillBean) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *BillBean) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *BillBean) GetOutUserId() uint32 {
	if m != nil && m.OutUserId != nil {
		return *m.OutUserId
	}
	return 0
}

func (m *BillBean) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *BillBean) GetDes() string {
	if m != nil && m.Des != nil {
		return *m.Des
	}
	return ""
}

func (m *BillBean) GetAmount() int64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *BillBean) GetPai() *MJPai {
	if m != nil {
		return m.Pai
	}
	return nil
}

func (m *BillBean) GetIsBird() bool {
	if m != nil && m.IsBird != nil {
		return *m.IsBird
	}
	return false
}

func (m *BillBean) GetIsQiShouHu() bool {
	if m != nil && m.IsQiShouHu != nil {
		return *m.IsQiShouHu
	}
	return false
}

type Bill struct {
	WinAmount        *int64      `protobuf:"varint,1,opt,name=WinAmount" json:"WinAmount,omitempty"`
	Bills            []*BillBean `protobuf:"bytes,2,rep,name=bills" json:"bills,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *Bill) Reset()                    { *m = Bill{} }
func (m *Bill) String() string            { return proto.CompactTextString(m) }
func (*Bill) ProtoMessage()               {}
func (*Bill) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Bill) GetWinAmount() int64 {
	if m != nil && m.WinAmount != nil {
		return *m.WinAmount
	}
	return 0
}

func (m *Bill) GetBills() []*BillBean {
	if m != nil {
		return m.Bills
	}
	return nil
}

type StatiscRound struct {
	Round             *int32  `protobuf:"varint,1,opt,name=Round" json:"Round,omitempty"`
	GameNumber        *int32  `protobuf:"varint,2,opt,name=GameNumber" json:"GameNumber,omitempty"`
	Result            *string `protobuf:"bytes,3,opt,name=Result" json:"Result,omitempty"`
	WinAmount         *int64  `protobuf:"varint,4,opt,name=WinAmount" json:"WinAmount,omitempty"`
	CountHu           *int32  `protobuf:"varint,5,opt,name=CountHu" json:"CountHu,omitempty"`
	CountDianPao      *int32  `protobuf:"varint,6,opt,name=CountDianPao" json:"CountDianPao,omitempty"`
	CountMingGang     *int32  `protobuf:"varint,7,opt,name=CountMingGang" json:"CountMingGang,omitempty"`
	CountDianGang     *int32  `protobuf:"varint,8,opt,name=CountDianGang" json:"CountDianGang,omitempty"`
	CountBaGnag       *int32  `protobuf:"varint,9,opt,name=CountBaGnag" json:"CountBaGnag,omitempty"`
	CountBeiBaGang    *int32  `protobuf:"varint,10,opt,name=CountBeiBaGang" json:"CountBeiBaGang,omitempty"`
	CountAnGang       *int32  `protobuf:"varint,11,opt,name=CountAnGang" json:"CountAnGang,omitempty"`
	CountBeiAnGang    *int32  `protobuf:"varint,12,opt,name=CountBeiAnGang" json:"CountBeiAnGang,omitempty"`
	CountZiMo         *int32  `protobuf:"varint,13,opt,name=CountZiMo" json:"CountZiMo,omitempty"`
	CountBeiZiMo      *int32  `protobuf:"varint,14,opt,name=CountBeiZiMo" json:"CountBeiZiMo,omitempty"`
	CountChaDaJiao    *int32  `protobuf:"varint,15,opt,name=CountChaDaJiao" json:"CountChaDaJiao,omitempty"`
	CountBeiChaJiao   *int32  `protobuf:"varint,16,opt,name=CountBeiChaJiao" json:"CountBeiChaJiao,omitempty"`
	CountChaHuaZhu    *int32  `protobuf:"varint,17,opt,name=CountChaHuaZhu" json:"CountChaHuaZhu,omitempty"`
	CountBeiChaHuaZhu *int32  `protobuf:"varint,18,opt,name=CountBeiChaHuaZhu" json:"CountBeiChaHuaZhu,omitempty"`
	CountCatchBird    *int32  `protobuf:"varint,19,opt,name=CountCatchBird" json:"CountCatchBird,omitempty"`
	CountCaughtBird   *int32  `protobuf:"varint,20,opt,name=CountCaughtBird" json:"CountCaughtBird,omitempty"`
	XXX_unrecognized  []byte  `json:"-"`
}

func (m *StatiscRound) Reset()                    { *m = StatiscRound{} }
func (m *StatiscRound) String() string            { return proto.CompactTextString(m) }
func (*StatiscRound) ProtoMessage()               {}
func (*StatiscRound) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *StatiscRound) GetRound() int32 {
	if m != nil && m.Round != nil {
		return *m.Round
	}
	return 0
}

func (m *StatiscRound) GetGameNumber() int32 {
	if m != nil && m.GameNumber != nil {
		return *m.GameNumber
	}
	return 0
}

func (m *StatiscRound) GetResult() string {
	if m != nil && m.Result != nil {
		return *m.Result
	}
	return ""
}

func (m *StatiscRound) GetWinAmount() int64 {
	if m != nil && m.WinAmount != nil {
		return *m.WinAmount
	}
	return 0
}

func (m *StatiscRound) GetCountHu() int32 {
	if m != nil && m.CountHu != nil {
		return *m.CountHu
	}
	return 0
}

func (m *StatiscRound) GetCountDianPao() int32 {
	if m != nil && m.CountDianPao != nil {
		return *m.CountDianPao
	}
	return 0
}

func (m *StatiscRound) GetCountMingGang() int32 {
	if m != nil && m.CountMingGang != nil {
		return *m.CountMingGang
	}
	return 0
}

func (m *StatiscRound) GetCountDianGang() int32 {
	if m != nil && m.CountDianGang != nil {
		return *m.CountDianGang
	}
	return 0
}

func (m *StatiscRound) GetCountBaGnag() int32 {
	if m != nil && m.CountBaGnag != nil {
		return *m.CountBaGnag
	}
	return 0
}

func (m *StatiscRound) GetCountBeiBaGang() int32 {
	if m != nil && m.CountBeiBaGang != nil {
		return *m.CountBeiBaGang
	}
	return 0
}

func (m *StatiscRound) GetCountAnGang() int32 {
	if m != nil && m.CountAnGang != nil {
		return *m.CountAnGang
	}
	return 0
}

func (m *StatiscRound) GetCountBeiAnGang() int32 {
	if m != nil && m.CountBeiAnGang != nil {
		return *m.CountBeiAnGang
	}
	return 0
}

func (m *StatiscRound) GetCountZiMo() int32 {
	if m != nil && m.CountZiMo != nil {
		return *m.CountZiMo
	}
	return 0
}

func (m *StatiscRound) GetCountBeiZiMo() int32 {
	if m != nil && m.CountBeiZiMo != nil {
		return *m.CountBeiZiMo
	}
	return 0
}

func (m *StatiscRound) GetCountChaDaJiao() int32 {
	if m != nil && m.CountChaDaJiao != nil {
		return *m.CountChaDaJiao
	}
	return 0
}

func (m *StatiscRound) GetCountBeiChaJiao() int32 {
	if m != nil && m.CountBeiChaJiao != nil {
		return *m.CountBeiChaJiao
	}
	return 0
}

func (m *StatiscRound) GetCountChaHuaZhu() int32 {
	if m != nil && m.CountChaHuaZhu != nil {
		return *m.CountChaHuaZhu
	}
	return 0
}

func (m *StatiscRound) GetCountBeiChaHuaZhu() int32 {
	if m != nil && m.CountBeiChaHuaZhu != nil {
		return *m.CountBeiChaHuaZhu
	}
	return 0
}

func (m *StatiscRound) GetCountCatchBird() int32 {
	if m != nil && m.CountCatchBird != nil {
		return *m.CountCatchBird
	}
	return 0
}

func (m *StatiscRound) GetCountCaughtBird() int32 {
	if m != nil && m.CountCaughtBird != nil {
		return *m.CountCaughtBird
	}
	return 0
}

type MjUserStatisc struct {
	RoundBean         []*StatiscRound `protobuf:"bytes,1,rep,name=roundBean" json:"roundBean,omitempty"`
	WinCoin           *int64          `protobuf:"varint,2,opt,name=WinCoin" json:"WinCoin,omitempty"`
	CountHu           *int32          `protobuf:"varint,3,opt,name=CountHu" json:"CountHu,omitempty"`
	CountZiMo         *int32          `protobuf:"varint,4,opt,name=CountZiMo" json:"CountZiMo,omitempty"`
	CountBeiZiMo      *int32          `protobuf:"varint,5,opt,name=CountBeiZiMo" json:"CountBeiZiMo,omitempty"`
	CountDianPao      *int32          `protobuf:"varint,6,opt,name=CountDianPao" json:"CountDianPao,omitempty"`
	CountAnGang       *int32          `protobuf:"varint,7,opt,name=CountAnGang" json:"CountAnGang,omitempty"`
	CountBeiAnGang    *int32          `protobuf:"varint,8,opt,name=CountBeiAnGang" json:"CountBeiAnGang,omitempty"`
	CountMingGang     *int32          `protobuf:"varint,9,opt,name=CountMingGang" json:"CountMingGang,omitempty"`
	CountDianGang     *int32          `protobuf:"varint,10,opt,name=CountDianGang" json:"CountDianGang,omitempty"`
	CountBaGang       *int32          `protobuf:"varint,11,opt,name=CountBaGang" json:"CountBaGang,omitempty"`
	CountBeiBaGang    *int32          `protobuf:"varint,12,opt,name=CountBeiBaGang" json:"CountBeiBaGang,omitempty"`
	CountChaDaJiao    *int32          `protobuf:"varint,13,opt,name=CountChaDaJiao" json:"CountChaDaJiao,omitempty"`
	CountBeiChaJiao   *int32          `protobuf:"varint,14,opt,name=CountBeiChaJiao" json:"CountBeiChaJiao,omitempty"`
	CountChaHuaZhu    *int32          `protobuf:"varint,15,opt,name=CountChaHuaZhu" json:"CountChaHuaZhu,omitempty"`
	CountBeiChaHuaZhu *int32          `protobuf:"varint,16,opt,name=CountBeiChaHuaZhu" json:"CountBeiChaHuaZhu,omitempty"`
	CountCatchBird    *int32          `protobuf:"varint,17,opt,name=CountCatchBird" json:"CountCatchBird,omitempty"`
	CountCaughtBird   *int32          `protobuf:"varint,18,opt,name=CountCaughtBird" json:"CountCaughtBird,omitempty"`
	XXX_unrecognized  []byte          `json:"-"`
}

func (m *MjUserStatisc) Reset()                    { *m = MjUserStatisc{} }
func (m *MjUserStatisc) String() string            { return proto.CompactTextString(m) }
func (*MjUserStatisc) ProtoMessage()               {}
func (*MjUserStatisc) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *MjUserStatisc) GetRoundBean() []*StatiscRound {
	if m != nil {
		return m.RoundBean
	}
	return nil
}

func (m *MjUserStatisc) GetWinCoin() int64 {
	if m != nil && m.WinCoin != nil {
		return *m.WinCoin
	}
	return 0
}

func (m *MjUserStatisc) GetCountHu() int32 {
	if m != nil && m.CountHu != nil {
		return *m.CountHu
	}
	return 0
}

func (m *MjUserStatisc) GetCountZiMo() int32 {
	if m != nil && m.CountZiMo != nil {
		return *m.CountZiMo
	}
	return 0
}

func (m *MjUserStatisc) GetCountBeiZiMo() int32 {
	if m != nil && m.CountBeiZiMo != nil {
		return *m.CountBeiZiMo
	}
	return 0
}

func (m *MjUserStatisc) GetCountDianPao() int32 {
	if m != nil && m.CountDianPao != nil {
		return *m.CountDianPao
	}
	return 0
}

func (m *MjUserStatisc) GetCountAnGang() int32 {
	if m != nil && m.CountAnGang != nil {
		return *m.CountAnGang
	}
	return 0
}

func (m *MjUserStatisc) GetCountBeiAnGang() int32 {
	if m != nil && m.CountBeiAnGang != nil {
		return *m.CountBeiAnGang
	}
	return 0
}

func (m *MjUserStatisc) GetCountMingGang() int32 {
	if m != nil && m.CountMingGang != nil {
		return *m.CountMingGang
	}
	return 0
}

func (m *MjUserStatisc) GetCountDianGang() int32 {
	if m != nil && m.CountDianGang != nil {
		return *m.CountDianGang
	}
	return 0
}

func (m *MjUserStatisc) GetCountBaGang() int32 {
	if m != nil && m.CountBaGang != nil {
		return *m.CountBaGang
	}
	return 0
}

func (m *MjUserStatisc) GetCountBeiBaGang() int32 {
	if m != nil && m.CountBeiBaGang != nil {
		return *m.CountBeiBaGang
	}
	return 0
}

func (m *MjUserStatisc) GetCountChaDaJiao() int32 {
	if m != nil && m.CountChaDaJiao != nil {
		return *m.CountChaDaJiao
	}
	return 0
}

func (m *MjUserStatisc) GetCountBeiChaJiao() int32 {
	if m != nil && m.CountBeiChaJiao != nil {
		return *m.CountBeiChaJiao
	}
	return 0
}

func (m *MjUserStatisc) GetCountChaHuaZhu() int32 {
	if m != nil && m.CountChaHuaZhu != nil {
		return *m.CountChaHuaZhu
	}
	return 0
}

func (m *MjUserStatisc) GetCountBeiChaHuaZhu() int32 {
	if m != nil && m.CountBeiChaHuaZhu != nil {
		return *m.CountBeiChaHuaZhu
	}
	return 0
}

func (m *MjUserStatisc) GetCountCatchBird() int32 {
	if m != nil && m.CountCatchBird != nil {
		return *m.CountCatchBird
	}
	return 0
}

func (m *MjUserStatisc) GetCountCaughtBird() int32 {
	if m != nil && m.CountCaughtBird != nil {
		return *m.CountCaughtBird
	}
	return 0
}

type PMjUser struct {
	GameData         *PlayerGameData `protobuf:"bytes,1,opt,name=GameData" json:"GameData,omitempty"`
	DingQue          *bool           `protobuf:"varint,2,opt,name=DingQue" json:"DingQue,omitempty"`
	Exchanged        *bool           `protobuf:"varint,3,opt,name=Exchanged" json:"Exchanged,omitempty"`
	Ready            *bool           `protobuf:"varint,4,opt,name=Ready" json:"Ready,omitempty"`
	PreMoGangInfo    *GangPaiInfo    `protobuf:"bytes,5,opt,name=PreMoGangInfo" json:"PreMoGangInfo,omitempty"`
	Bill             *Bill           `protobuf:"bytes,6,opt,name=bill" json:"bill,omitempty"`
	Statisc          *MjUserStatisc  `protobuf:"bytes,7,opt,name=statisc" json:"statisc,omitempty"`
	IsBanker         *bool           `protobuf:"varint,8,opt,name=IsBanker" json:"IsBanker,omitempty"`
	ApplyDissolve    *int32          `protobuf:"varint,9,opt,name=applyDissolve" json:"applyDissolve,omitempty"`
	AgentMode        *bool           `protobuf:"varint,10,opt,name=agentMode" json:"agentMode,omitempty"`
	ActTimeoutCount  *int32          `protobuf:"varint,11,opt,name=actTimeoutCount" json:"actTimeoutCount,omitempty"`
	NeedHaidiStatus  *int32          `protobuf:"varint,12,opt,name=needHaidiStatus" json:"needHaidiStatus,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *PMjUser) Reset()                    { *m = PMjUser{} }
func (m *PMjUser) String() string            { return proto.CompactTextString(m) }
func (*PMjUser) ProtoMessage()               {}
func (*PMjUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *PMjUser) GetGameData() *PlayerGameData {
	if m != nil {
		return m.GameData
	}
	return nil
}

func (m *PMjUser) GetDingQue() bool {
	if m != nil && m.DingQue != nil {
		return *m.DingQue
	}
	return false
}

func (m *PMjUser) GetExchanged() bool {
	if m != nil && m.Exchanged != nil {
		return *m.Exchanged
	}
	return false
}

func (m *PMjUser) GetReady() bool {
	if m != nil && m.Ready != nil {
		return *m.Ready
	}
	return false
}

func (m *PMjUser) GetPreMoGangInfo() *GangPaiInfo {
	if m != nil {
		return m.PreMoGangInfo
	}
	return nil
}

func (m *PMjUser) GetBill() *Bill {
	if m != nil {
		return m.Bill
	}
	return nil
}

func (m *PMjUser) GetStatisc() *MjUserStatisc {
	if m != nil {
		return m.Statisc
	}
	return nil
}

func (m *PMjUser) GetIsBanker() bool {
	if m != nil && m.IsBanker != nil {
		return *m.IsBanker
	}
	return false
}

func (m *PMjUser) GetApplyDissolve() int32 {
	if m != nil && m.ApplyDissolve != nil {
		return *m.ApplyDissolve
	}
	return 0
}

func (m *PMjUser) GetAgentMode() bool {
	if m != nil && m.AgentMode != nil {
		return *m.AgentMode
	}
	return false
}

func (m *PMjUser) GetActTimeoutCount() int32 {
	if m != nil && m.ActTimeoutCount != nil {
		return *m.ActTimeoutCount
	}
	return 0
}

func (m *PMjUser) GetNeedHaidiStatus() int32 {
	if m != nil && m.NeedHaidiStatus != nil {
		return *m.NeedHaidiStatus
	}
	return 0
}

// 一个麻将room
type PMjRoom struct {
	RoomId           *int32   `protobuf:"varint,1,opt,name=RoomId" json:"RoomId,omitempty"`
	RoomType         *int32   `protobuf:"varint,2,opt,name=RoomType" json:"RoomType,omitempty"`
	RoomLevel        *int32   `protobuf:"varint,3,opt,name=RoomLevel" json:"RoomLevel,omitempty"`
	ReadyTime        *int64   `protobuf:"varint,4,opt,name=ReadyTime" json:"ReadyTime,omitempty"`
	BeginTime        *int64   `protobuf:"varint,5,opt,name=BeginTime" json:"BeginTime,omitempty"`
	EndTIme          *int64   `protobuf:"varint,6,opt,name=EndTIme" json:"EndTIme,omitempty"`
	PaiIndexNow      *int32   `protobuf:"varint,7,opt,name=PaiIndexNow" json:"PaiIndexNow,omitempty"`
	MjPais           []*MJPai `protobuf:"bytes,8,rep,name=MjPais" json:"MjPais,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *PMjRoom) Reset()                    { *m = PMjRoom{} }
func (m *PMjRoom) String() string            { return proto.CompactTextString(m) }
func (*PMjRoom) ProtoMessage()               {}
func (*PMjRoom) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *PMjRoom) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *PMjRoom) GetRoomType() int32 {
	if m != nil && m.RoomType != nil {
		return *m.RoomType
	}
	return 0
}

func (m *PMjRoom) GetRoomLevel() int32 {
	if m != nil && m.RoomLevel != nil {
		return *m.RoomLevel
	}
	return 0
}

func (m *PMjRoom) GetReadyTime() int64 {
	if m != nil && m.ReadyTime != nil {
		return *m.ReadyTime
	}
	return 0
}

func (m *PMjRoom) GetBeginTime() int64 {
	if m != nil && m.BeginTime != nil {
		return *m.BeginTime
	}
	return 0
}

func (m *PMjRoom) GetEndTIme() int64 {
	if m != nil && m.EndTIme != nil {
		return *m.EndTIme
	}
	return 0
}

func (m *PMjRoom) GetPaiIndexNow() int32 {
	if m != nil && m.PaiIndexNow != nil {
		return *m.PaiIndexNow
	}
	return 0
}

func (m *PMjRoom) GetMjPais() []*MJPai {
	if m != nil {
		return m.MjPais
	}
	return nil
}

//
type RunningDeskKeys struct {
	Keys             []int32 `protobuf:"varint,1,rep,name=keys" json:"keys,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RunningDeskKeys) Reset()                    { *m = RunningDeskKeys{} }
func (m *RunningDeskKeys) String() string            { return proto.CompactTextString(m) }
func (*RunningDeskKeys) ProtoMessage()               {}
func (*RunningDeskKeys) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *RunningDeskKeys) GetKeys() []int32 {
	if m != nil {
		return m.Keys
	}
	return nil
}

func init() {
	proto.RegisterType((*MJPai)(nil), "majiang.MJPai")
	proto.RegisterType((*MJHandPai)(nil), "majiang.MJHandPai")
	proto.RegisterType((*HuPaiInfo)(nil), "majiang.HuPaiInfo")
	proto.RegisterType((*GangPaiInfo)(nil), "majiang.GangPaiInfo")
	proto.RegisterType((*GuoHuInfo)(nil), "majiang.GuoHuInfo")
	proto.RegisterType((*PlayerGameData)(nil), "majiang.PlayerGameData")
	proto.RegisterType((*CheckBean)(nil), "majiang.CheckBean")
	proto.RegisterType((*CheckCase)(nil), "majiang.CheckCase")
	proto.RegisterType((*PMjDesk)(nil), "majiang.PMjDesk")
	proto.RegisterType((*BillBean)(nil), "majiang.BillBean")
	proto.RegisterType((*Bill)(nil), "majiang.Bill")
	proto.RegisterType((*StatiscRound)(nil), "majiang.StatiscRound")
	proto.RegisterType((*MjUserStatisc)(nil), "majiang.MjUserStatisc")
	proto.RegisterType((*PMjUser)(nil), "majiang.PMjUser")
	proto.RegisterType((*PMjRoom)(nil), "majiang.PMjRoom")
	proto.RegisterType((*RunningDeskKeys)(nil), "majiang.RunningDeskKeys")
}

var fileDescriptor0 = []byte{
	// 1519 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x57, 0x4d, 0x73, 0x1b, 0x45,
	0x13, 0x2e, 0x7d, 0x59, 0xd2, 0xc8, 0x92, 0xec, 0x8d, 0x93, 0xec, 0x5b, 0x79, 0x21, 0x61, 0x53,
	0x40, 0xf8, 0xa8, 0x1c, 0x7c, 0xe5, 0x64, 0x49, 0xc1, 0x76, 0x88, 0x82, 0x23, 0x27, 0x95, 0xaa,
	0x5c, 0x52, 0x13, 0xed, 0x58, 0x3b, 0xf6, 0x6a, 0xc6, 0xec, 0x87, 0x63, 0xf3, 0x27, 0xe0, 0x0a,
	0x37, 0xae, 0xf0, 0x03, 0x38, 0x73, 0xe3, 0x1f, 0x71, 0xa5, 0xbb, 0x67, 0x56, 0x2b, 0xc9, 0xbb,
	0x2e, 0x28, 0x4e, 0xd2, 0xf6, 0xf4, 0xf6, 0xf4, 0xf4, 0xf3, 0xf4, 0xd3, 0xb3, 0x6c, 0x67, 0xce,
	0x83, 0x53, 0xad, 0x66, 0x6f, 0x63, 0x11, 0x5d, 0x88, 0xe8, 0xf1, 0x79, 0xa4, 0x13, 0xed, 0x34,
	0xe7, 0xfc, 0x54, 0x72, 0x35, 0xf3, 0x06, 0xac, 0x31, 0x7e, 0x7a, 0xc4, 0xa5, 0xd3, 0x65, 0x0d,
	0xa9, 0x7c, 0x71, 0xe9, 0x56, 0x1e, 0x54, 0x1e, 0x35, 0x9c, 0x1e, 0xdb, 0x38, 0x09, 0xf5, 0x7b,
	0x11, 0xb9, 0x55, 0x7a, 0x86, 0xe5, 0x0b, 0x1e, 0xa6, 0xc2, 0xad, 0xd1, 0x63, 0x87, 0xd5, 0x7c,
	0x11, 0xbb, 0x75, 0x78, 0x68, 0x7b, 0xbf, 0x54, 0x59, 0x7b, 0xfc, 0xf4, 0x80, 0x2b, 0x1f, 0x03,
	0xfd, 0x9f, 0xd5, 0xe1, 0x27, 0x86, 0x38, 0xb5, 0x47, 0x9d, 0xdd, 0xde, 0x63, 0xbb, 0xd3, 0x63,
	0xb3, 0xcd, 0x03, 0xd6, 0x3a, 0x12, 0x6a, 0x46, 0x1e, 0xd5, 0x32, 0x8f, 0x7d, 0x6e, 0x3d, 0x6a,
	0x85, 0x1e, 0x1f, 0xb2, 0x8d, 0x83, 0x94, 0xd6, 0xeb, 0x85, 0xeb, 0x1f, 0xe0, 0x51, 0xe0, 0x8f,
	0xdb, 0x80, 0xf4, 0xae, 0x2f, 0x6f, 0xb3, 0xf6, 0x77, 0xa9, 0xf8, 0xda, 0x9c, 0x6e, 0x83, 0x8e,
	0x73, 0x9f, 0x35, 0xbf, 0x4d, 0x13, 0x0a, 0xd9, 0x2c, 0x0c, 0x09, 0x0e, 0xc3, 0x40, 0x92, 0x43,
	0xab, 0x2c, 0x27, 0xda, 0x73, 0xd7, 0x6d, 0x17, 0x6d, 0xea, 0xfd, 0x56, 0x61, 0x6d, 0x4a, 0xfa,
	0x50, 0x9d, 0x68, 0xc7, 0x61, 0x2c, 0x16, 0xca, 0x7f, 0x05, 0x98, 0x1c, 0xfa, 0x54, 0xf1, 0x2e,
	0xa6, 0x35, 0x13, 0x89, 0x35, 0x55, 0xc9, 0x04, 0x45, 0x7f, 0x77, 0xf5, 0x3a, 0xd0, 0xb6, 0xe8,
	0x80, 0x49, 0x90, 0xbe, 0xbc, 0x3a, 0x17, 0x54, 0xf7, 0x86, 0xd3, 0x67, 0xcd, 0x73, 0x2e, 0xc9,
	0xd0, 0x80, 0xa4, 0xac, 0xc3, 0x48, 0xc4, 0x53, 0x3a, 0x56, 0xdb, 0xb9, 0xc7, 0x6a, 0xe0, 0x00,
	0x47, 0x2a, 0x2a, 0x03, 0x40, 0x78, 0xc2, 0x15, 0x1c, 0xc7, 0xc2, 0x1b, 0x4f, 0x75, 0x24, 0x28,
	0xfb, 0x9a, 0xf7, 0x63, 0x85, 0x75, 0x2c, 0x08, 0xff, 0x21, 0xdf, 0x2d, 0xd6, 0x9a, 0x41, 0x90,
	0xa5, 0x8c, 0x6d, 0x42, 0xc5, 0xb8, 0x80, 0x7b, 0x12, 0x71, 0x15, 0x9f, 0x2c, 0x60, 0x61, 0xac,
	0xfa, 0x2e, 0xa5, 0xf4, 0x5b, 0xde, 0x5b, 0xd6, 0xde, 0x4f, 0xf5, 0x41, 0xfa, 0x6f, 0xf2, 0xb1,
	0xdb, 0xd5, 0x0a, 0xb7, 0x43, 0x86, 0x73, 0x75, 0x1c, 0xa4, 0x26, 0x37, 0xef, 0x87, 0x1a, 0xeb,
	0x1d, 0x85, 0xfc, 0x4a, 0x44, 0xfb, 0x7c, 0x2e, 0x46, 0x3c, 0xe1, 0xce, 0x43, 0xd6, 0x0c, 0x0c,
	0xab, 0x69, 0x8f, 0xce, 0xae, 0xb3, 0x14, 0x23, 0xe3, 0xbb, 0x87, 0x45, 0xc7, 0xac, 0x2c, 0x9f,
	0x73, 0x9f, 0x1c, 0xef, 0x4f, 0x18, 0xf3, 0xc1, 0x64, 0xb2, 0xb7, 0xac, 0x2e, 0xf6, 0xa3, 0x8a,
	0x91, 0x97, 0xe1, 0xf6, 0xce, 0xc2, 0x6b, 0x19, 0x8f, 0xcf, 0xd9, 0x26, 0xc6, 0xdb, 0xcf, 0x7c,
	0x1b, 0x37, 0xfb, 0x12, 0x0a, 0x79, 0x69, 0xcb, 0x7d, 0xa1, 0xae, 0x89, 0x4e, 0x78, 0x78, 0x4c,
	0x5c, 0x68, 0x12, 0x08, 0x1f, 0x43, 0x5d, 0xb3, 0xc2, 0x5b, 0xf2, 0xe7, 0xa9, 0xe7, 0x90, 0x3c,
	0x62, 0x5b, 0xe2, 0x72, 0x0a, 0xe5, 0x9a, 0x89, 0x21, 0x8f, 0xfc, 0x18, 0xfa, 0x09, 0xc8, 0x54,
	0xd4, 0x2a, 0x9f, 0xb2, 0xfe, 0x8a, 0xe7, 0xa1, 0x72, 0x59, 0x91, 0xa3, 0xf7, 0x17, 0xf4, 0xcc,
	0x30, 0x10, 0xd3, 0xb3, 0x81, 0xe0, 0x0a, 0xf1, 0x4a, 0x97, 0xf1, 0x06, 0xf6, 0x0f, 0xb9, 0x42,
	0x31, 0x21, 0xb4, 0x5b, 0xd6, 0x80, 0xc7, 0x21, 0xc4, 0x5b, 0x48, 0xc7, 0x21, 0x16, 0x9d, 0x00,
	0x6e, 0x39, 0xb7, 0x58, 0x87, 0xa2, 0x1d, 0x27, 0x3c, 0x49, 0x63, 0x22, 0x21, 0xb5, 0xcc, 0x94,
	0x2b, 0xe8, 0x6d, 0xa2, 0x5c, 0xcb, 0x3e, 0xc3, 0xb1, 0x0c, 0xed, 0x50, 0x8d, 0xa6, 0x81, 0xa4,
	0x3c, 0x4b, 0x3a, 0x9f, 0x3c, 0x20, 0x2c, 0xb2, 0xa4, 0xb0, 0xf7, 0x6d, 0x1e, 0x83, 0x14, 0x8e,
	0x89, 0x21, 0x3f, 0x82, 0x82, 0xc2, 0xa2, 0x89, 0xd9, 0xb9, 0xf9, 0xe4, 0x43, 0x1e, 0x0b, 0x64,
	0x98, 0x29, 0x83, 0xd5, 0xd4, 0xbc, 0xfc, 0x79, 0x75, 0xd6, 0x0e, 0x67, 0x44, 0xdb, 0x63, 0x8c,
	0x8c, 0x14, 0xb4, 0xa4, 0x0d, 0xa0, 0x6d, 0x4c, 0xcf, 0x20, 0x60, 0x75, 0xaa, 0xec, 0x97, 0xd0,
	0x08, 0x91, 0x80, 0xe7, 0x25, 0x7e, 0x55, 0x4a, 0x39, 0xf3, 0xd0, 0xee, 0x4c, 0xe1, 0x76, 0xa9,
	0x8c, 0xd7, 0x77, 0xd9, 0x31, 0x84, 0x3d, 0xe2, 0x7a, 0xa8, 0x53, 0x95, 0x58, 0x6a, 0xa1, 0x12,
	0xa3, 0x13, 0x61, 0xd6, 0xa2, 0x36, 0xff, 0xb9, 0xce, 0x9a, 0x47, 0xe3, 0x53, 0x10, 0xb1, 0x33,
	0x64, 0xe3, 0xfc, 0x74, 0xa2, 0xf5, 0x9c, 0x14, 0xa4, 0x92, 0xe3, 0x75, 0x3e, 0xe6, 0x97, 0x74,
	0xc4, 0x1a, 0x8a, 0xc6, 0x14, 0x0b, 0xf9, 0x3c, 0x9d, 0x5b, 0xd5, 0x21, 0x6d, 0x48, 0x92, 0x50,
	0xcc, 0x85, 0x4a, 0xac, 0xee, 0xc0, 0x46, 0xdf, 0xcb, 0xb1, 0x9e, 0x70, 0x5f, 0x6a, 0x0b, 0xbc,
	0xcb, 0xb6, 0xb2, 0x16, 0x3a, 0x48, 0xb9, 0x59, 0x31, 0xaa, 0x73, 0x87, 0xf5, 0x74, 0x12, 0x88,
	0x28, 0x36, 0xd5, 0xd5, 0x97, 0x34, 0x13, 0x48, 0x6e, 0x83, 0xd4, 0x38, 0x1a, 0xd1, 0x04, 0xe4,
	0xf7, 0xc2, 0xd0, 0x14, 0xb7, 0x98, 0xea, 0x00, 0x35, 0x9e, 0xc4, 0xb8, 0x14, 0x92, 0x1c, 0x81,
	0xa3, 0x3f, 0xc3, 0x34, 0x8a, 0x75, 0x04, 0x7c, 0xb0, 0x29, 0xbc, 0xc4, 0x3e, 0x44, 0x3d, 0x32,
	0x05, 0xdb, 0x24, 0xfb, 0x6d, 0xd6, 0x05, 0xbf, 0x28, 0x37, 0x77, 0xb3, 0x16, 0x5d, 0xb0, 0xc5,
	0xed, 0xad, 0x29, 0x55, 0xce, 0x23, 0xa8, 0xcc, 0xde, 0x34, 0x91, 0x17, 0x02, 0x01, 0x77, 0xfb,
	0x59, 0x17, 0x81, 0x8d, 0x0c, 0x5b, 0x4b, 0x06, 0xaa, 0xf8, 0x76, 0x56, 0x4f, 0x25, 0x2e, 0x93,
	0x01, 0x57, 0x67, 0xe0, 0xe4, 0x90, 0x13, 0xe4, 0x77, 0x82, 0x14, 0xc6, 0x1c, 0x9e, 0xc9, 0xb9,
	0x4c, 0xdc, 0x5b, 0x99, 0xe2, 0x47, 0x19, 0x5e, 0x3b, 0x59, 0xe5, 0xa7, 0x5a, 0x2a, 0xe3, 0x74,
	0x9b, 0x20, 0x03, 0x13, 0x3a, 0x3d, 0x13, 0x17, 0x22, 0x74, 0xef, 0x64, 0xef, 0xf1, 0xf3, 0xf3,
	0xf0, 0x6a, 0x04, 0xf3, 0xf5, 0x6e, 0xd6, 0xac, 0x8b, 0xf7, 0x5e, 0x3d, 0x73, 0x5d, 0x1a, 0x4b,
	0x3f, 0x55, 0x58, 0x6b, 0x20, 0xc3, 0x30, 0xd3, 0x83, 0x75, 0xfd, 0x07, 0xca, 0xae, 0xe8, 0xff,
	0x26, 0xab, 0x27, 0x98, 0xca, 0xf5, 0x3b, 0x0b, 0xbe, 0xcd, 0xe7, 0x54, 0xc2, 0x06, 0x25, 0x65,
	0x47, 0xc5, 0x46, 0xd9, 0xa8, 0x90, 0xf1, 0x40, 0x46, 0xbe, 0x15, 0x05, 0x28, 0x89, 0x8c, 0x5f,
	0xc8, 0xe3, 0x40, 0xa7, 0xa0, 0x2e, 0x86, 0xb8, 0x5f, 0xb1, 0x3a, 0xa6, 0x86, 0x69, 0xbc, 0x96,
	0x6a, 0xcf, 0xc4, 0xae, 0x50, 0xec, 0x07, 0x30, 0x16, 0x61, 0x29, 0xbb, 0xf0, 0x6c, 0x2f, 0xa2,
	0x67, 0x67, 0xf1, 0xfe, 0xa8, 0xb1, 0x4d, 0xec, 0x5c, 0x19, 0x4f, 0x27, 0xf0, 0xa2, 0x8f, 0x92,
	0x41, 0x7f, 0x2c, 0xeb, 0x61, 0x43, 0x1c, 0x4a, 0x40, 0xf2, 0x77, 0x8b, 0x1b, 0x19, 0x24, 0x35,
	0x11, 0x71, 0x1a, 0x26, 0x74, 0xbc, 0xf6, 0xea, 0xc6, 0x75, 0xda, 0x18, 0x15, 0x11, 0x1f, 0x21,
	0x49, 0x43, 0x7a, 0x68, 0x43, 0x32, 0x8c, 0x4c, 0x2f, 0x5a, 0xc2, 0x23, 0xab, 0xd0, 0x3a, 0x96,
	0xb6, 0x15, 0x9b, 0x2b, 0xe6, 0x91, 0x6d, 0x13, 0xcb, 0x7a, 0x54, 0x1a, 0x34, 0x0f, 0xf8, 0xbe,
	0xe2, 0x33, 0x92, 0x3c, 0x22, 0xac, 0x31, 0x0a, 0x09, 0x76, 0x74, 0x66, 0x2b, 0xce, 0x7b, 0x26,
	0x42, 0x67, 0xdd, 0xd9, 0xda, 0x37, 0x33, 0xae, 0x90, 0xfd, 0x0d, 0xb4, 0xaa, 0x65, 0x76, 0x96,
	0x30, 0xb8, 0x92, 0xb5, 0xb7, 0x12, 0x60, 0x18, 0xf0, 0x11, 0x7f, 0x2a, 0xe1, 0x20, 0x7d, 0xb2,
	0xdf, 0x65, 0xfd, 0xcc, 0x1b, 0x96, 0x68, 0x61, 0x6b, 0xfd, 0x05, 0x68, 0xf6, 0x37, 0x30, 0xf3,
	0x0d, 0xb7, 0xff, 0xc7, 0xb6, 0x97, 0x5e, 0xb0, 0x4b, 0xce, 0xea, 0x2b, 0x3c, 0x99, 0x06, 0x84,
	0xfd, 0xad, 0x95, 0x3d, 0x86, 0x3c, 0x9d, 0x05, 0x09, 0x2d, 0x10, 0xd3, 0xbd, 0xdf, 0x6b, 0xac,
	0x3b, 0x3e, 0x45, 0xfe, 0x59, 0x24, 0x61, 0x24, 0x02, 0xd1, 0x01, 0x44, 0x84, 0xd8, 0x4a, 0xf7,
	0xed, 0x05, 0xf6, 0x2b, 0x70, 0x03, 0x50, 0x80, 0xdd, 0x10, 0x08, 0x6f, 0x65, 0x6d, 0x09, 0xb9,
	0xda, 0xf5, 0xda, 0xd4, 0x0b, 0x6b, 0x73, 0x13, 0xc4, 0x6b, 0x38, 0x34, 0x4b, 0x70, 0x68, 0x15,
	0xf3, 0xa1, 0x5d, 0xcc, 0x07, 0xb6, 0xce, 0x87, 0x42, 0x88, 0xad, 0x7d, 0xb3, 0x04, 0xb9, 0x6e,
	0x19, 0x72, 0xbd, 0x12, 0xe4, 0xfa, 0xe5, 0xc8, 0x6d, 0x95, 0x20, 0xb7, 0x5d, 0x86, 0x1c, 0x41,
	0xed, 0xfd, 0x59, 0xa5, 0x99, 0x83, 0xd0, 0x39, 0x9f, 0xe1, 0xd7, 0x87, 0xb9, 0xfe, 0xd9, 0x3b,
	0xdf, 0xdd, 0x05, 0x64, 0x6b, 0xb7, 0x43, 0xc0, 0x68, 0x04, 0x15, 0x7a, 0x01, 0x1f, 0x45, 0xe6,
	0x02, 0x02, 0x18, 0x3d, 0xb1, 0x17, 0x1b, 0x3f, 0xbf, 0x82, 0x4c, 0x04, 0xf7, 0xaf, 0xec, 0x15,
	0xe4, 0x0b, 0xd6, 0x85, 0xc9, 0x3a, 0xd6, 0xff, 0x68, 0xb0, 0xde, 0x63, 0x75, 0x94, 0x0d, 0xab,
	0x49, 0xdd, 0x15, 0xd5, 0x80, 0x4b, 0x54, 0x33, 0x36, 0x0c, 0xb2, 0xd7, 0xfb, 0x3b, 0xb9, 0x66,
	0xad, 0x90, 0x10, 0xa4, 0xf5, 0x30, 0xb6, 0xe2, 0x4d, 0x4a, 0x85, 0x38, 0x66, 0x62, 0x1b, 0xeb,
	0xf0, 0x42, 0x58, 0x78, 0x21, 0x7b, 0x3e, 0x83, 0x91, 0x39, 0xd6, 0xbe, 0xb0, 0x37, 0x15, 0xa8,
	0x18, 0x87, 0x59, 0x20, 0xe7, 0x42, 0xa7, 0x89, 0x19, 0x38, 0x9d, 0xac, 0x94, 0x4a, 0x08, 0xff,
	0x80, 0x4b, 0x5f, 0xda, 0x1b, 0x07, 0xe1, 0xeb, 0xfd, 0x5a, 0xa1, 0x52, 0xe2, 0xd0, 0x26, 0x81,
	0x82, 0xdf, 0xc3, 0x4c, 0xc4, 0x20, 0x93, 0xc5, 0x30, 0xaf, 0x66, 0x5b, 0x4e, 0x16, 0x93, 0x60,
	0xc1, 0x73, 0x2a, 0x18, 0x6e, 0x6a, 0x55, 0x0c, 0x4c, 0x03, 0x31, 0x93, 0x8a, 0x4c, 0x8d, 0xac,
	0x3d, 0x9e, 0x28, 0xff, 0xe5, 0x21, 0x18, 0x36, 0xc8, 0x00, 0x24, 0xa4, 0xb2, 0xc1, 0x07, 0xec,
	0x73, 0xfd, 0xde, 0xf2, 0x1b, 0xbe, 0xc9, 0xc6, 0xa7, 0xe5, 0xdf, 0x6c, 0xde, 0x7d, 0xd6, 0x9f,
	0xa4, 0x4a, 0x01, 0x86, 0x38, 0xa4, 0xbf, 0x11, 0x57, 0x31, 0x4e, 0x8c, 0x33, 0xf8, 0xa5, 0x6e,
	0x6d, 0xfc, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x1e, 0x88, 0x77, 0x3c, 0x0f, 0x00, 0x00,
}
