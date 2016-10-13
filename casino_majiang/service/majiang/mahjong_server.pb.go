// Code generated by protoc-gen-go.
// source: mahjong_server.proto
// DO NOT EDIT!

/*
Package majiang is a generated protocol buffer package.

It is generated from these files:
	mahjong_server.proto

It has these top-level messages:
	MJPai
	MJHandPai
	HuPaiInfo
	GangPaiInfo
	PlayerGameData
	CheckBean
	CheckCase
	MjDesk
	MjUser
	MjSession
	MjRoom
*/
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
	Pais             []*MJPai `protobuf:"bytes,1,rep,name=Pais" json:"Pais,omitempty"`
	PengPais         []*MJPai `protobuf:"bytes,2,rep,name=PengPais" json:"PengPais,omitempty"`
	GangPais         []*MJPai `protobuf:"bytes,3,rep,name=GangPais" json:"GangPais,omitempty"`
	HuPais           []*MJPai `protobuf:"bytes,4,rep,name=HuPais" json:"HuPais,omitempty"`
	InPai            *MJPai   `protobuf:"bytes,5,opt,name=inPai" json:"inPai,omitempty"`
	QueFlower        *int32   `protobuf:"varint,6,opt,name=queFlower" json:"queFlower,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
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

type HuPaiInfo struct {
	SendUserId       *uint32 `protobuf:"varint,1,opt,name=sendUserId" json:"sendUserId,omitempty"`
	ByWho            *int32  `protobuf:"varint,2,opt,name=byWho" json:"byWho,omitempty"`
	HuType           *int32  `protobuf:"varint,3,opt,name=huType" json:"huType,omitempty"`
	CardType         *int32  `protobuf:"varint,4,opt,name=cardType" json:"cardType,omitempty"`
	Fan              *int32  `protobuf:"varint,5,opt,name=fan" json:"fan,omitempty"`
	Pai              *MJPai  `protobuf:"bytes,6,opt,name=pai" json:"pai,omitempty"`
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

func (m *HuPaiInfo) GetCardType() int32 {
	if m != nil && m.CardType != nil {
		return *m.CardType
	}
	return 0
}

func (m *HuPaiInfo) GetFan() int32 {
	if m != nil && m.Fan != nil {
		return *m.Fan
	}
	return 0
}

func (m *HuPaiInfo) GetPai() *MJPai {
	if m != nil {
		return m.Pai
	}
	return nil
}

type GangPaiInfo struct {
	SendUserId       *uint32 `protobuf:"varint,1,opt,name=sendUserId" json:"sendUserId,omitempty"`
	ByWho            *int32  `protobuf:"varint,2,opt,name=byWho" json:"byWho,omitempty"`
	GangType         *int32  `protobuf:"varint,3,opt,name=gangType" json:"gangType,omitempty"`
	Pai              *MJPai  `protobuf:"bytes,4,opt,name=pai" json:"pai,omitempty"`
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

// 一个玩家的游戏信息
type PlayerGameData struct {
	HandPai          *MJHandPai     `protobuf:"bytes,1,opt,name=handPai" json:"handPai,omitempty"`
	HuInfo           []*HuPaiInfo   `protobuf:"bytes,2,rep,name=huInfo" json:"huInfo,omitempty"`
	DianHuInfo       []*HuPaiInfo   `protobuf:"bytes,3,rep,name=dianHuInfo" json:"dianHuInfo,omitempty"`
	GangInfo         []*GangPaiInfo `protobuf:"bytes,4,rep,name=gangInfo" json:"gangInfo,omitempty"`
	DianGangInfo     []*GangPaiInfo `protobuf:"bytes,5,rep,name=dianGangInfo" json:"dianGangInfo,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *PlayerGameData) Reset()                    { *m = PlayerGameData{} }
func (m *PlayerGameData) String() string            { return proto.CompactTextString(m) }
func (*PlayerGameData) ProtoMessage()               {}
func (*PlayerGameData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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

// 需要确认的时间
type CheckBean struct {
	UserId           *uint32 `protobuf:"varint,1,opt,name=userId" json:"userId,omitempty"`
	CanPeng          *bool   `protobuf:"varint,2,opt,name=CanPeng" json:"CanPeng,omitempty"`
	CanGang          *bool   `protobuf:"varint,3,opt,name=CanGang" json:"CanGang,omitempty"`
	CanHu            *bool   `protobuf:"varint,4,opt,name=CanHu" json:"CanHu,omitempty"`
	CheckStatus      *int32  `protobuf:"varint,5,opt,name=CheckStatus" json:"CheckStatus,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CheckBean) Reset()                    { *m = CheckBean{} }
func (m *CheckBean) String() string            { return proto.CompactTextString(m) }
func (*CheckBean) ProtoMessage()               {}
func (*CheckBean) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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

type CheckCase struct {
	CheckB           []*CheckBean `protobuf:"bytes,1,rep,name=CheckB" json:"CheckB,omitempty"`
	CheckStatus      *int32       `protobuf:"varint,2,opt,name=CheckStatus" json:"CheckStatus,omitempty"`
	CheckMJPai       *MJPai       `protobuf:"bytes,3,opt,name=CheckMJPai" json:"CheckMJPai,omitempty"`
	UserIdOut        *uint32      `protobuf:"varint,4,opt,name=UserIdOut" json:"UserIdOut,omitempty"`
	PreOutGangInfo   *GangPaiInfo `protobuf:"bytes,5,opt,name=PreOutGangInfo" json:"PreOutGangInfo,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *CheckCase) Reset()                    { *m = CheckCase{} }
func (m *CheckCase) String() string            { return proto.CompactTextString(m) }
func (*CheckCase) ProtoMessage()               {}
func (*CheckCase) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

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

// 麻将牌的结构
type MjDesk struct {
	DeskId           *int32     `protobuf:"varint,1,opt,name=DeskId" json:"DeskId,omitempty"`
	RoomId           *int32     `protobuf:"varint,2,opt,name=RoomId" json:"RoomId,omitempty"`
	Status           *int32     `protobuf:"varint,3,opt,name=Status" json:"Status,omitempty"`
	Users            []*MjUser  `protobuf:"bytes,4,rep,name=users" json:"users,omitempty"`
	Password         *string    `protobuf:"bytes,5,opt,name=Password" json:"Password,omitempty"`
	Owner            *uint32    `protobuf:"varint,6,opt,name=Owner" json:"Owner,omitempty"`
	CreateFee        *int64     `protobuf:"varint,7,opt,name=CreateFee" json:"CreateFee,omitempty"`
	MjRoomType       *int32     `protobuf:"varint,8,opt,name=mjRoomType" json:"mjRoomType,omitempty"`
	BoardsCout       *int32     `protobuf:"varint,9,opt,name=boardsCout" json:"boardsCout,omitempty"`
	CapMax           *int64     `protobuf:"varint,10,opt,name=capMax" json:"capMax,omitempty"`
	CardsNum         *int32     `protobuf:"varint,11,opt,name=cardsNum" json:"cardsNum,omitempty"`
	Settlement       *int32     `protobuf:"varint,12,opt,name=settlement" json:"settlement,omitempty"`
	BaseValue        *int64     `protobuf:"varint,13,opt,name=baseValue" json:"baseValue,omitempty"`
	ZiMoRadio        *int32     `protobuf:"varint,14,opt,name=ziMoRadio" json:"ziMoRadio,omitempty"`
	DianGangHuaRadio *int32     `protobuf:"varint,15,opt,name=dianGangHuaRadio" json:"dianGangHuaRadio,omitempty"`
	OthersCheckBox   []int32    `protobuf:"varint,16,rep,name=othersCheckBox" json:"othersCheckBox,omitempty"`
	HuRadio          *int32     `protobuf:"varint,17,opt,name=huRadio" json:"huRadio,omitempty"`
	AllMJPai         []*MJPai   `protobuf:"bytes,18,rep,name=AllMJPai" json:"AllMJPai,omitempty"`
	DeskMJPai        []*MJPai   `protobuf:"bytes,19,rep,name=DeskMJPai" json:"DeskMJPai,omitempty"`
	MJPaiCursor      *int32     `protobuf:"varint,20,opt,name=MJPaiCursor" json:"MJPaiCursor,omitempty"`
	TotalPlayCount   *int32     `protobuf:"varint,21,opt,name=TotalPlayCount" json:"TotalPlayCount,omitempty"`
	CurrPlayCount    *int32     `protobuf:"varint,22,opt,name=CurrPlayCount" json:"CurrPlayCount,omitempty"`
	Banker           *uint32    `protobuf:"varint,23,opt,name=Banker" json:"Banker,omitempty"`
	CheckCase        *CheckCase `protobuf:"bytes,24,opt,name=CheckCase" json:"CheckCase,omitempty"`
	NextUserCursor   *uint32    `protobuf:"varint,25,opt,name=NextUserCursor" json:"NextUserCursor,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *MjDesk) Reset()                    { *m = MjDesk{} }
func (m *MjDesk) String() string            { return proto.CompactTextString(m) }
func (*MjDesk) ProtoMessage()               {}
func (*MjDesk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *MjDesk) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *MjDesk) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *MjDesk) GetStatus() int32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return 0
}

func (m *MjDesk) GetUsers() []*MjUser {
	if m != nil {
		return m.Users
	}
	return nil
}

func (m *MjDesk) GetPassword() string {
	if m != nil && m.Password != nil {
		return *m.Password
	}
	return ""
}

func (m *MjDesk) GetOwner() uint32 {
	if m != nil && m.Owner != nil {
		return *m.Owner
	}
	return 0
}

func (m *MjDesk) GetCreateFee() int64 {
	if m != nil && m.CreateFee != nil {
		return *m.CreateFee
	}
	return 0
}

func (m *MjDesk) GetMjRoomType() int32 {
	if m != nil && m.MjRoomType != nil {
		return *m.MjRoomType
	}
	return 0
}

func (m *MjDesk) GetBoardsCout() int32 {
	if m != nil && m.BoardsCout != nil {
		return *m.BoardsCout
	}
	return 0
}

func (m *MjDesk) GetCapMax() int64 {
	if m != nil && m.CapMax != nil {
		return *m.CapMax
	}
	return 0
}

func (m *MjDesk) GetCardsNum() int32 {
	if m != nil && m.CardsNum != nil {
		return *m.CardsNum
	}
	return 0
}

func (m *MjDesk) GetSettlement() int32 {
	if m != nil && m.Settlement != nil {
		return *m.Settlement
	}
	return 0
}

func (m *MjDesk) GetBaseValue() int64 {
	if m != nil && m.BaseValue != nil {
		return *m.BaseValue
	}
	return 0
}

func (m *MjDesk) GetZiMoRadio() int32 {
	if m != nil && m.ZiMoRadio != nil {
		return *m.ZiMoRadio
	}
	return 0
}

func (m *MjDesk) GetDianGangHuaRadio() int32 {
	if m != nil && m.DianGangHuaRadio != nil {
		return *m.DianGangHuaRadio
	}
	return 0
}

func (m *MjDesk) GetOthersCheckBox() []int32 {
	if m != nil {
		return m.OthersCheckBox
	}
	return nil
}

func (m *MjDesk) GetHuRadio() int32 {
	if m != nil && m.HuRadio != nil {
		return *m.HuRadio
	}
	return 0
}

func (m *MjDesk) GetAllMJPai() []*MJPai {
	if m != nil {
		return m.AllMJPai
	}
	return nil
}

func (m *MjDesk) GetDeskMJPai() []*MJPai {
	if m != nil {
		return m.DeskMJPai
	}
	return nil
}

func (m *MjDesk) GetMJPaiCursor() int32 {
	if m != nil && m.MJPaiCursor != nil {
		return *m.MJPaiCursor
	}
	return 0
}

func (m *MjDesk) GetTotalPlayCount() int32 {
	if m != nil && m.TotalPlayCount != nil {
		return *m.TotalPlayCount
	}
	return 0
}

func (m *MjDesk) GetCurrPlayCount() int32 {
	if m != nil && m.CurrPlayCount != nil {
		return *m.CurrPlayCount
	}
	return 0
}

func (m *MjDesk) GetBanker() uint32 {
	if m != nil && m.Banker != nil {
		return *m.Banker
	}
	return 0
}

func (m *MjDesk) GetCheckCase() *CheckCase {
	if m != nil {
		return m.CheckCase
	}
	return nil
}

func (m *MjDesk) GetNextUserCursor() uint32 {
	if m != nil && m.NextUserCursor != nil {
		return *m.NextUserCursor
	}
	return 0
}

// 手里的一副牌
type MjUser struct {
	UserId           *uint32         `protobuf:"varint,1,opt,name=UserId" json:"UserId,omitempty"`
	Coin             *int64          `protobuf:"varint,2,opt,name=Coin" json:"Coin,omitempty"`
	GameData         *PlayerGameData `protobuf:"bytes,3,opt,name=GameData" json:"GameData,omitempty"`
	Status           *int32          `protobuf:"varint,4,opt,name=Status" json:"Status,omitempty"`
	IsBreak          *bool           `protobuf:"varint,5,opt,name=IsBreak" json:"IsBreak,omitempty"`
	IsLeave          *bool           `protobuf:"varint,6,opt,name=IsLeave" json:"IsLeave,omitempty"`
	DeskId           *int32          `protobuf:"varint,7,opt,name=DeskId" json:"DeskId,omitempty"`
	RoomId           *int32          `protobuf:"varint,8,opt,name=RoomId" json:"RoomId,omitempty"`
	DingQue          *bool           `protobuf:"varint,9,opt,name=DingQue" json:"DingQue,omitempty"`
	Exchanged        *bool           `protobuf:"varint,10,opt,name=Exchanged" json:"Exchanged,omitempty"`
	Ready            *bool           `protobuf:"varint,11,opt,name=Ready" json:"Ready,omitempty"`
	IsBanker         *bool           `protobuf:"varint,12,opt,name=IsBanker" json:"IsBanker,omitempty"`
	PreMoGangInfo    *GangPaiInfo    `protobuf:"bytes,13,opt,name=PreMoGangInfo" json:"PreMoGangInfo,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *MjUser) Reset()                    { *m = MjUser{} }
func (m *MjUser) String() string            { return proto.CompactTextString(m) }
func (*MjUser) ProtoMessage()               {}
func (*MjUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *MjUser) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *MjUser) GetCoin() int64 {
	if m != nil && m.Coin != nil {
		return *m.Coin
	}
	return 0
}

func (m *MjUser) GetGameData() *PlayerGameData {
	if m != nil {
		return m.GameData
	}
	return nil
}

func (m *MjUser) GetStatus() int32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return 0
}

func (m *MjUser) GetIsBreak() bool {
	if m != nil && m.IsBreak != nil {
		return *m.IsBreak
	}
	return false
}

func (m *MjUser) GetIsLeave() bool {
	if m != nil && m.IsLeave != nil {
		return *m.IsLeave
	}
	return false
}

func (m *MjUser) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *MjUser) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *MjUser) GetDingQue() bool {
	if m != nil && m.DingQue != nil {
		return *m.DingQue
	}
	return false
}

func (m *MjUser) GetExchanged() bool {
	if m != nil && m.Exchanged != nil {
		return *m.Exchanged
	}
	return false
}

func (m *MjUser) GetReady() bool {
	if m != nil && m.Ready != nil {
		return *m.Ready
	}
	return false
}

func (m *MjUser) GetIsBanker() bool {
	if m != nil && m.IsBanker != nil {
		return *m.IsBanker
	}
	return false
}

func (m *MjUser) GetPreMoGangInfo() *GangPaiInfo {
	if m != nil {
		return m.PreMoGangInfo
	}
	return nil
}

// 用户的session
type MjSession struct {
	UserId           *uint32 `protobuf:"varint,1,opt,name=UserId" json:"UserId,omitempty"`
	RoomId           *int32  `protobuf:"varint,2,opt,name=RoomId" json:"RoomId,omitempty"`
	DeskId           *int32  `protobuf:"varint,3,opt,name=DeskId" json:"DeskId,omitempty"`
	GameStatus       *int32  `protobuf:"varint,4,opt,name=GameStatus" json:"GameStatus,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MjSession) Reset()                    { *m = MjSession{} }
func (m *MjSession) String() string            { return proto.CompactTextString(m) }
func (*MjSession) ProtoMessage()               {}
func (*MjSession) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *MjSession) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *MjSession) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *MjSession) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *MjSession) GetGameStatus() int32 {
	if m != nil && m.GameStatus != nil {
		return *m.GameStatus
	}
	return 0
}

// 一个麻将room
type MjRoom struct {
	RoomId           *int32    `protobuf:"varint,1,opt,name=RoomId" json:"RoomId,omitempty"`
	RoomType         *int32    `protobuf:"varint,2,opt,name=RoomType" json:"RoomType,omitempty"`
	ReadyTime        *int64    `protobuf:"varint,3,opt,name=ReadyTime" json:"ReadyTime,omitempty"`
	BeginTime        *int64    `protobuf:"varint,4,opt,name=BeginTime" json:"BeginTime,omitempty"`
	EndTIme          *int64    `protobuf:"varint,5,opt,name=EndTIme" json:"EndTIme,omitempty"`
	Desks            []*MjDesk `protobuf:"bytes,6,rep,name=Desks" json:"Desks,omitempty"`
	PaiIndexNow      *int32    `protobuf:"varint,7,opt,name=PaiIndexNow" json:"PaiIndexNow,omitempty"`
	MjPais           []*MJPai  `protobuf:"bytes,8,rep,name=MjPais" json:"MjPais,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *MjRoom) Reset()                    { *m = MjRoom{} }
func (m *MjRoom) String() string            { return proto.CompactTextString(m) }
func (*MjRoom) ProtoMessage()               {}
func (*MjRoom) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *MjRoom) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *MjRoom) GetRoomType() int32 {
	if m != nil && m.RoomType != nil {
		return *m.RoomType
	}
	return 0
}

func (m *MjRoom) GetReadyTime() int64 {
	if m != nil && m.ReadyTime != nil {
		return *m.ReadyTime
	}
	return 0
}

func (m *MjRoom) GetBeginTime() int64 {
	if m != nil && m.BeginTime != nil {
		return *m.BeginTime
	}
	return 0
}

func (m *MjRoom) GetEndTIme() int64 {
	if m != nil && m.EndTIme != nil {
		return *m.EndTIme
	}
	return 0
}

func (m *MjRoom) GetDesks() []*MjDesk {
	if m != nil {
		return m.Desks
	}
	return nil
}

func (m *MjRoom) GetPaiIndexNow() int32 {
	if m != nil && m.PaiIndexNow != nil {
		return *m.PaiIndexNow
	}
	return 0
}

func (m *MjRoom) GetMjPais() []*MJPai {
	if m != nil {
		return m.MjPais
	}
	return nil
}

func init() {
	proto.RegisterType((*MJPai)(nil), "majiang.MJPai")
	proto.RegisterType((*MJHandPai)(nil), "majiang.MJHandPai")
	proto.RegisterType((*HuPaiInfo)(nil), "majiang.HuPaiInfo")
	proto.RegisterType((*GangPaiInfo)(nil), "majiang.GangPaiInfo")
	proto.RegisterType((*PlayerGameData)(nil), "majiang.PlayerGameData")
	proto.RegisterType((*CheckBean)(nil), "majiang.CheckBean")
	proto.RegisterType((*CheckCase)(nil), "majiang.CheckCase")
	proto.RegisterType((*MjDesk)(nil), "majiang.MjDesk")
	proto.RegisterType((*MjUser)(nil), "majiang.MjUser")
	proto.RegisterType((*MjSession)(nil), "majiang.MjSession")
	proto.RegisterType((*MjRoom)(nil), "majiang.MjRoom")
}

var fileDescriptor0 = []byte{
	// 948 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x55, 0x5d, 0x72, 0xdb, 0x46,
	0x0c, 0x1e, 0x59, 0xa2, 0x44, 0x42, 0x96, 0xec, 0x30, 0x4e, 0xc2, 0x4e, 0x5b, 0x4f, 0xca, 0x4e,
	0x3b, 0xe9, 0xcf, 0xf8, 0x21, 0x37, 0xa8, 0xe4, 0xc4, 0x76, 0xa6, 0xb2, 0x5d, 0xc7, 0x6d, 0x1f,
	0x33, 0x6b, 0x73, 0x2d, 0x51, 0x96, 0x76, 0x5d, 0x2e, 0x69, 0xcb, 0xb9, 0x47, 0x0f, 0xd1, 0xe9,
	0x7b, 0x5f, 0x7b, 0x8a, 0xde, 0xa7, 0x00, 0x76, 0x29, 0xb1, 0x0a, 0x9d, 0xe9, 0x93, 0xb4, 0x58,
	0x60, 0xf1, 0xe1, 0xc3, 0x07, 0x10, 0x76, 0xe6, 0x62, 0x32, 0xd5, 0x6a, 0xfc, 0xce, 0xc8, 0xec,
	0x56, 0x66, 0x7b, 0x37, 0x99, 0xce, 0x75, 0xd8, 0x99, 0x8b, 0x69, 0x2a, 0xd4, 0x38, 0x1e, 0x80,
	0x37, 0x7a, 0x73, 0x2a, 0xd2, 0xb0, 0x07, 0x5e, 0xaa, 0x12, 0xb9, 0x88, 0x1a, 0xcf, 0x1b, 0x2f,
	0xbc, 0xb0, 0x0f, 0xed, 0xab, 0x99, 0xbe, 0x93, 0x59, 0xb4, 0xc1, 0x67, 0xbc, 0xbe, 0x15, 0xb3,
	0x42, 0x46, 0x4d, 0x3e, 0x76, 0xa1, 0x99, 0x48, 0x13, 0xb5, 0xf0, 0x10, 0xc4, 0x7f, 0x37, 0x20,
	0x18, 0xbd, 0x39, 0x14, 0x2a, 0xa1, 0x87, 0x3e, 0x83, 0x16, 0xfe, 0x18, 0x7c, 0xa7, 0xf9, 0xa2,
	0xfb, 0xb2, 0xbf, 0xe7, 0x32, 0xed, 0xd9, 0x34, 0xcf, 0xc1, 0x3f, 0x95, 0x6a, 0xcc, 0x1e, 0x1b,
	0x0f, 0x79, 0x1c, 0x08, 0xe7, 0xd1, 0xac, 0xf5, 0xd8, 0x85, 0xf6, 0x61, 0xc1, 0xf7, 0xad, 0xda,
	0xfb, 0xcf, 0xa9, 0x14, 0xfc, 0x13, 0x79, 0x08, 0xef, 0xc3, 0xeb, 0x47, 0x10, 0xfc, 0x56, 0xc8,
	0xd7, 0xb6, 0xba, 0x36, 0x95, 0x13, 0xbf, 0x87, 0x80, 0x5f, 0x3c, 0x52, 0x57, 0x3a, 0x0c, 0x01,
	0x8c, 0x54, 0xc9, 0xcf, 0x48, 0xd8, 0x51, 0xc2, 0x74, 0xf4, 0xa8, 0xfc, 0x8b, 0xfb, 0x5f, 0x27,
	0xda, 0xb1, 0x81, 0xec, 0x4c, 0x8a, 0xf3, 0xfb, 0x9b, 0x92, 0x8e, 0x6d, 0xf0, 0x2f, 0x45, 0x96,
	0xb0, 0xa5, 0x55, 0x12, 0x74, 0x25, 0x14, 0x23, 0xf0, 0xc2, 0x4f, 0xa1, 0x79, 0x83, 0x70, 0xda,
	0x75, 0x70, 0xe2, 0x77, 0xd0, 0x75, 0xf5, 0xfe, 0xdf, 0xec, 0x98, 0x6d, 0x8c, 0x11, 0x95, 0xfc,
	0x2e, 0x41, 0xab, 0x36, 0xc1, 0x3f, 0x0d, 0xe8, 0x9f, 0xce, 0xc4, 0xbd, 0xcc, 0x0e, 0xc4, 0x5c,
	0xee, 0x8b, 0x5c, 0x84, 0x5f, 0x42, 0x67, 0x62, 0xdb, 0xc5, 0x19, 0xba, 0x2f, 0xc3, 0x4a, 0x4c,
	0xd9, 0xc8, 0x98, 0x8a, 0x24, 0x4c, 0xae, 0x51, 0x2b, 0x9f, 0x15, 0x57, 0x5f, 0x03, 0x24, 0x68,
	0x3a, 0xb4, 0x7e, 0xcd, 0x8f, 0xf8, 0x31, 0x64, 0xf6, 0xb2, 0x4d, 0xdb, 0x59, 0x7a, 0x55, 0xab,
	0xff, 0x16, 0x36, 0xe9, 0xbd, 0x83, 0xd2, 0xd7, 0x7b, 0xd8, 0x37, 0x16, 0x10, 0x0c, 0x27, 0xf2,
	0xf2, 0x7a, 0x20, 0x85, 0xa2, 0x8e, 0x14, 0x55, 0xca, 0xb6, 0xa0, 0x33, 0x14, 0x8a, 0xa4, 0xc6,
	0xa4, 0xf9, 0xce, 0x40, 0xf1, 0xcc, 0x99, 0x4f, 0xa4, 0x0e, 0x09, 0x39, 0xb3, 0xe6, 0x87, 0x8f,
	0xa1, 0xcb, 0xaf, 0xbd, 0xcd, 0x45, 0x5e, 0x18, 0xdb, 0xb8, 0xf8, 0x8f, 0x86, 0xcb, 0x31, 0x14,
	0x46, 0x12, 0x21, 0x36, 0xa1, 0xd3, 0xf6, 0xaa, 0xd0, 0x15, 0x8e, 0xb5, 0x67, 0x6c, 0xc3, 0x62,
	0x00, 0x36, 0x72, 0x3f, 0x38, 0x7d, 0xad, 0x2a, 0x6d, 0xcf, 0x4f, 0x8a, 0x9c, 0x21, 0xf5, 0xc2,
	0xef, 0xb1, 0x6f, 0x99, 0xc4, 0x73, 0x85, 0x8e, 0xc6, 0x83, 0x74, 0xfc, 0xd9, 0x82, 0xf6, 0x68,
	0xba, 0x2f, 0xcd, 0x35, 0x91, 0x41, 0xbf, 0x8e, 0x0c, 0x96, 0xeb, 0x99, 0xd6, 0x73, 0x3c, 0x2f,
	0xe5, 0xeb, 0xf0, 0x59, 0xf9, 0xec, 0x82, 0x47, 0xe4, 0x95, 0xf3, 0xb4, 0xb5, 0x82, 0x36, 0x25,
	0x4c, 0x24, 0xb8, 0x53, 0x61, 0xcc, 0x9d, 0xce, 0x12, 0x86, 0x10, 0x10, 0x79, 0x27, 0x77, 0xca,
	0xcd, 0x4f, 0x8f, 0xc0, 0x0f, 0x33, 0x29, 0x72, 0xf9, 0x5a, 0xca, 0xa8, 0x83, 0xa6, 0x26, 0xe9,
	0x78, 0x3e, 0xa5, 0xac, 0x2c, 0x53, 0x9f, 0xf3, 0xa0, 0xed, 0x42, 0xe3, 0x9c, 0x98, 0xa1, 0xc6,
	0x22, 0x83, 0x12, 0xcb, 0xa5, 0xb8, 0x19, 0x89, 0x45, 0x04, 0x1c, 0xe7, 0x46, 0xc9, 0x1c, 0x17,
	0xf3, 0xa8, 0x5b, 0x46, 0x19, 0x99, 0xe7, 0x33, 0x39, 0x97, 0x2a, 0x8f, 0x36, 0xd9, 0x86, 0x09,
	0x2f, 0xb0, 0x25, 0xbf, 0xf0, 0x4a, 0xea, 0x71, 0x20, 0x9a, 0xde, 0xa7, 0x23, 0x7d, 0x26, 0x92,
	0x54, 0x47, 0x7d, 0xf6, 0x8a, 0x60, 0xbb, 0x54, 0xd3, 0x61, 0x21, 0xec, 0xcd, 0x16, 0xdf, 0x3c,
	0x85, 0xbe, 0xce, 0x27, 0x58, 0xb2, 0xed, 0x9c, 0x5e, 0x44, 0xdb, 0x58, 0xba, 0x47, 0x2a, 0x99,
	0x14, 0xd6, 0xf1, 0x11, 0x3b, 0xe2, 0x36, 0xfa, 0x61, 0x36, 0xb3, 0x8d, 0x0b, 0x6b, 0xb7, 0xcd,
	0x17, 0x10, 0x10, 0xd9, 0xd6, 0xe5, 0x71, 0xad, 0x0b, 0x8a, 0x82, 0xff, 0x0c, 0x8b, 0xcc, 0xe8,
	0x2c, 0xda, 0x29, 0x21, 0x9c, 0xeb, 0x5c, 0xcc, 0x68, 0x34, 0x91, 0x0f, 0x2c, 0xed, 0x09, 0xdb,
	0x9f, 0x40, 0x0f, 0xfd, 0xb2, 0x95, 0xf9, 0x69, 0xc9, 0xd3, 0x40, 0xa8, 0x6b, 0xa4, 0xfc, 0x19,
	0x53, 0xfe, 0x55, 0x45, 0x99, 0x51, 0xb4, 0x36, 0xc4, 0x2b, 0xcd, 0x62, 0x96, 0x63, 0xb9, 0xc8,
	0xa9, 0x8d, 0x2e, 0xfb, 0x27, 0x14, 0x1e, 0xff, 0xbe, 0x41, 0x6a, 0xe1, 0xee, 0xe2, 0xcb, 0xff,
	0xd9, 0x36, 0x9b, 0xd0, 0x1a, 0xea, 0x54, 0xb1, 0x56, 0x9a, 0xe1, 0x37, 0xb4, 0x8e, 0xed, 0xda,
	0x70, 0xca, 0x7d, 0xb6, 0x4c, 0xb3, 0xb6, 0x55, 0x56, 0xb2, 0xb2, 0x3b, 0x10, 0xc9, 0x3c, 0x32,
	0x03, 0xd4, 0xc5, 0x35, 0xab, 0xc6, 0xb7, 0x86, 0x1f, 0xa5, 0xb8, 0x95, 0xac, 0x1b, 0xbf, 0x22,
	0xd4, 0xce, 0x9a, 0x50, 0xfd, 0xf2, 0x85, 0xfd, 0x54, 0x8d, 0x7f, 0xc2, 0x26, 0x07, 0x1c, 0x80,
	0x4d, 0x7e, 0xb5, 0xb8, 0xc4, 0xdd, 0x35, 0x96, 0x09, 0x0b, 0x86, 0xe7, 0xf8, 0x4c, 0x8a, 0xe4,
	0x9e, 0xd5, 0xe2, 0x93, 0x7e, 0x30, 0xa9, 0x65, 0x6a, 0x93, 0x2d, 0xdf, 0x41, 0x0f, 0xc7, 0x68,
	0xa4, 0x97, 0x53, 0xd4, 0xfb, 0xc8, 0x14, 0x9d, 0xe0, 0xa7, 0x6c, 0xfa, 0x56, 0x1a, 0x93, 0x6a,
	0xf5, 0x01, 0x33, 0x35, 0x73, 0xe4, 0xe0, 0x37, 0x4b, 0xa5, 0x12, 0x19, 0x55, 0x12, 0xe2, 0xbf,
	0x1a, 0x44, 0x34, 0x85, 0x55, 0xc2, 0x1b, 0xe5, 0x1e, 0x5f, 0x0e, 0xc8, 0x46, 0x29, 0x6b, 0xae,
	0xe5, 0x3c, 0x9d, 0xdb, 0xd5, 0xce, 0xb2, 0x1e, 0xc8, 0x71, 0xaa, 0xd8, 0xd4, 0x62, 0x13, 0xb2,
	0xf2, 0x4a, 0x25, 0xe7, 0x47, 0x68, 0xf0, 0xd8, 0x80, 0xf3, 0x4b, 0x38, 0x0c, 0xb2, 0xba, 0x3e,
	0xbf, 0xbc, 0x0f, 0x50, 0x7f, 0x5c, 0x1f, 0x7e, 0xde, 0x8f, 0xf5, 0x9d, 0xe3, 0x7a, 0x97, 0x70,
	0xf1, 0x57, 0xd4, 0xaf, 0x13, 0xed, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x12, 0x80, 0x6d, 0x46,
	0x39, 0x08, 0x00, 0x00,
}
