// Code generated by protoc-gen-go.
// source: ddz_server.proto
// DO NOT EDIT!

/*
Package doudizhu is a generated protocol buffer package.

It is generated from these files:
	ddz_server.proto

It has these top-level messages:
	PPokerPai
	POutPokerPais
	PDdzDeskTongJi
	PDdzDesk
	PGameData
	PDdzBillBean
	PDdzBill
	PDdzUser
	PDdzRoom
	PDdzbak
	PDdzSession
*/
package doudizhu

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

type EDDZTYPE int32

const (
	EDDZTYPE_DDZ_TYPE_HUANLEDIZHU EDDZTYPE = 1
	EDDZTYPE_DDZ_TYPE_SICHUAN     EDDZTYPE = 2
)

var EDDZTYPE_name = map[int32]string{
	1: "DDZ_TYPE_HUANLEDIZHU",
	2: "DDZ_TYPE_SICHUAN",
}
var EDDZTYPE_value = map[string]int32{
	"DDZ_TYPE_HUANLEDIZHU": 1,
	"DDZ_TYPE_SICHUAN":     2,
}

func (x EDDZTYPE) Enum() *EDDZTYPE {
	p := new(EDDZTYPE)
	*p = x
	return p
}
func (x EDDZTYPE) String() string {
	return proto.EnumName(EDDZTYPE_name, int32(x))
}
func (x *EDDZTYPE) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(EDDZTYPE_value, data, "EDDZTYPE")
	if err != nil {
		return err
	}
	*x = EDDZTYPE(value)
	return nil
}
func (EDDZTYPE) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// 牌 ,单张扑克牌
type PPokerPai struct {
	Id               *int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Des              *string `protobuf:"bytes,2,opt,name=des" json:"des,omitempty"`
	Value            *int32  `protobuf:"varint,3,opt,name=value" json:"value,omitempty"`
	Flower           *int32  `protobuf:"varint,4,opt,name=flower" json:"flower,omitempty"`
	Name             *string `protobuf:"bytes,5,opt,name=name" json:"name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PPokerPai) Reset()                    { *m = PPokerPai{} }
func (m *PPokerPai) String() string            { return proto.CompactTextString(m) }
func (*PPokerPai) ProtoMessage()               {}
func (*PPokerPai) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PPokerPai) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *PPokerPai) GetDes() string {
	if m != nil && m.Des != nil {
		return *m.Des
	}
	return ""
}

func (m *PPokerPai) GetValue() int32 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

func (m *PPokerPai) GetFlower() int32 {
	if m != nil && m.Flower != nil {
		return *m.Flower
	}
	return 0
}

func (m *PPokerPai) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

// 打出去的牌
type POutPokerPais struct {
	KeyValue         *int32       `protobuf:"varint,1,opt,name=keyValue" json:"keyValue,omitempty"`
	PokerPais        []*PPokerPai `protobuf:"bytes,2,rep,name=pokerPais" json:"pokerPais,omitempty"`
	Type             *int32       `protobuf:"varint,3,opt,name=type" json:"type,omitempty"`
	IsBomb           *bool        `protobuf:"varint,4,opt,name=isBomb" json:"isBomb,omitempty"`
	CountDuizi       *int32       `protobuf:"varint,5,opt,name=countDuizi" json:"countDuizi,omitempty"`
	CountSanzhang    *int32       `protobuf:"varint,6,opt,name=countSanzhang" json:"countSanzhang,omitempty"`
	CountSizhang     *int32       `protobuf:"varint,7,opt,name=countSizhang" json:"countSizhang,omitempty"`
	CountYizhang     *int32       `protobuf:"varint,8,opt,name=countYizhang" json:"countYizhang,omitempty"`
	UserId           *uint32      `protobuf:"varint,9,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *POutPokerPais) Reset()                    { *m = POutPokerPais{} }
func (m *POutPokerPais) String() string            { return proto.CompactTextString(m) }
func (*POutPokerPais) ProtoMessage()               {}
func (*POutPokerPais) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *POutPokerPais) GetKeyValue() int32 {
	if m != nil && m.KeyValue != nil {
		return *m.KeyValue
	}
	return 0
}

func (m *POutPokerPais) GetPokerPais() []*PPokerPai {
	if m != nil {
		return m.PokerPais
	}
	return nil
}

func (m *POutPokerPais) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *POutPokerPais) GetIsBomb() bool {
	if m != nil && m.IsBomb != nil {
		return *m.IsBomb
	}
	return false
}

func (m *POutPokerPais) GetCountDuizi() int32 {
	if m != nil && m.CountDuizi != nil {
		return *m.CountDuizi
	}
	return 0
}

func (m *POutPokerPais) GetCountSanzhang() int32 {
	if m != nil && m.CountSanzhang != nil {
		return *m.CountSanzhang
	}
	return 0
}

func (m *POutPokerPais) GetCountSizhang() int32 {
	if m != nil && m.CountSizhang != nil {
		return *m.CountSizhang
	}
	return 0
}

func (m *POutPokerPais) GetCountYizhang() int32 {
	if m != nil && m.CountYizhang != nil {
		return *m.CountYizhang
	}
	return 0
}

func (m *POutPokerPais) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 对 desk的统计
type PDdzDeskTongJi struct {
	Bombs            []*POutPokerPais `protobuf:"bytes,1,rep,name=bombs" json:"bombs,omitempty"`
	CountQiangDiZhu  *int32           `protobuf:"varint,2,opt,name=countQiangDiZhu" json:"countQiangDiZhu,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *PDdzDeskTongJi) Reset()                    { *m = PDdzDeskTongJi{} }
func (m *PDdzDeskTongJi) String() string            { return proto.CompactTextString(m) }
func (*PDdzDeskTongJi) ProtoMessage()               {}
func (*PDdzDeskTongJi) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PDdzDeskTongJi) GetBombs() []*POutPokerPais {
	if m != nil {
		return m.Bombs
	}
	return nil
}

func (m *PDdzDeskTongJi) GetCountQiangDiZhu() int32 {
	if m != nil && m.CountQiangDiZhu != nil {
		return *m.CountQiangDiZhu
	}
	return 0
}

// desk
type PDdzDesk struct {
	DeskId           *int32          `protobuf:"varint,1,opt,name=deskId" json:"deskId,omitempty"`
	Key              *string         `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	UserCountLimit   *int32          `protobuf:"varint,3,opt,name=userCountLimit" json:"userCountLimit,omitempty"`
	AllPokerPai      []*PPokerPai    `protobuf:"bytes,4,rep,name=allPokerPai" json:"allPokerPai,omitempty"`
	DiPokerPai       []*PPokerPai    `protobuf:"bytes,5,rep,name=diPokerPai" json:"diPokerPai,omitempty"`
	OutPai           *POutPokerPais  `protobuf:"bytes,6,opt,name=outPai" json:"outPai,omitempty"`
	Owner            *uint32         `protobuf:"varint,7,opt,name=owner" json:"owner,omitempty"`
	DizhuPaiUser     *uint32         `protobuf:"varint,8,opt,name=dizhuPaiUser" json:"dizhuPaiUser,omitempty"`
	DiZhuUserId      *uint32         `protobuf:"varint,9,opt,name=diZhuUserId" json:"diZhuUserId,omitempty"`
	ActiveUserId     *uint32         `protobuf:"varint,10,opt,name=activeUserId" json:"activeUserId,omitempty"`
	Tongji           *PDdzDeskTongJi `protobuf:"bytes,11,opt,name=tongji" json:"tongji,omitempty"`
	BaseValue        *int64          `protobuf:"varint,12,opt,name=baseValue" json:"baseValue,omitempty"`
	QingDizhuValue   *int64          `protobuf:"varint,13,opt,name=qingDizhuValue" json:"qingDizhuValue,omitempty"`
	WinValue         *int64          `protobuf:"varint,14,opt,name=winValue" json:"winValue,omitempty"`
	DdzType          *int32          `protobuf:"varint,15,opt,name=ddzType" json:"ddzType,omitempty"`
	RoomType         *int32          `protobuf:"varint,16,opt,name=RoomType" json:"RoomType,omitempty"`
	BoardsCount      *int32          `protobuf:"varint,17,opt,name=BoardsCount" json:"BoardsCount,omitempty"`
	CapMax           *int64          `protobuf:"varint,18,opt,name=CapMax" json:"CapMax,omitempty"`
	IsJiaoFen        *bool           `protobuf:"varint,19,opt,name=IsJiaoFen" json:"IsJiaoFen,omitempty"`
	RoomId           *int32          `protobuf:"varint,20,opt,name=roomId" json:"roomId,omitempty"`
	FootRate         *int32          `protobuf:"varint,21,opt,name=footRate" json:"footRate,omitempty"`
	PlayRate         *int32          `protobuf:"varint,22,opt,name=playRate" json:"playRate,omitempty"`
	GameStatus       *int32          `protobuf:"varint,23,opt,name=GameStatus" json:"GameStatus,omitempty"`
	InitRoomCoin     *int64          `protobuf:"varint,24,opt,name=initRoomCoin" json:"initRoomCoin,omitempty"`
	CurrPlayCount    *int32          `protobuf:"varint,25,opt,name=currPlayCount" json:"currPlayCount,omitempty"`
	TotalPlayCount   *int32          `protobuf:"varint,26,opt,name=totalPlayCount" json:"totalPlayCount,omitempty"`
	PlayerNum        *int32          `protobuf:"varint,27,opt,name=playerNum" json:"playerNum,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *PDdzDesk) Reset()                    { *m = PDdzDesk{} }
func (m *PDdzDesk) String() string            { return proto.CompactTextString(m) }
func (*PDdzDesk) ProtoMessage()               {}
func (*PDdzDesk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PDdzDesk) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *PDdzDesk) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *PDdzDesk) GetUserCountLimit() int32 {
	if m != nil && m.UserCountLimit != nil {
		return *m.UserCountLimit
	}
	return 0
}

func (m *PDdzDesk) GetAllPokerPai() []*PPokerPai {
	if m != nil {
		return m.AllPokerPai
	}
	return nil
}

func (m *PDdzDesk) GetDiPokerPai() []*PPokerPai {
	if m != nil {
		return m.DiPokerPai
	}
	return nil
}

func (m *PDdzDesk) GetOutPai() *POutPokerPais {
	if m != nil {
		return m.OutPai
	}
	return nil
}

func (m *PDdzDesk) GetOwner() uint32 {
	if m != nil && m.Owner != nil {
		return *m.Owner
	}
	return 0
}

func (m *PDdzDesk) GetDizhuPaiUser() uint32 {
	if m != nil && m.DizhuPaiUser != nil {
		return *m.DizhuPaiUser
	}
	return 0
}

func (m *PDdzDesk) GetDiZhuUserId() uint32 {
	if m != nil && m.DiZhuUserId != nil {
		return *m.DiZhuUserId
	}
	return 0
}

func (m *PDdzDesk) GetActiveUserId() uint32 {
	if m != nil && m.ActiveUserId != nil {
		return *m.ActiveUserId
	}
	return 0
}

func (m *PDdzDesk) GetTongji() *PDdzDeskTongJi {
	if m != nil {
		return m.Tongji
	}
	return nil
}

func (m *PDdzDesk) GetBaseValue() int64 {
	if m != nil && m.BaseValue != nil {
		return *m.BaseValue
	}
	return 0
}

func (m *PDdzDesk) GetQingDizhuValue() int64 {
	if m != nil && m.QingDizhuValue != nil {
		return *m.QingDizhuValue
	}
	return 0
}

func (m *PDdzDesk) GetWinValue() int64 {
	if m != nil && m.WinValue != nil {
		return *m.WinValue
	}
	return 0
}

func (m *PDdzDesk) GetDdzType() int32 {
	if m != nil && m.DdzType != nil {
		return *m.DdzType
	}
	return 0
}

func (m *PDdzDesk) GetRoomType() int32 {
	if m != nil && m.RoomType != nil {
		return *m.RoomType
	}
	return 0
}

func (m *PDdzDesk) GetBoardsCount() int32 {
	if m != nil && m.BoardsCount != nil {
		return *m.BoardsCount
	}
	return 0
}

func (m *PDdzDesk) GetCapMax() int64 {
	if m != nil && m.CapMax != nil {
		return *m.CapMax
	}
	return 0
}

func (m *PDdzDesk) GetIsJiaoFen() bool {
	if m != nil && m.IsJiaoFen != nil {
		return *m.IsJiaoFen
	}
	return false
}

func (m *PDdzDesk) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *PDdzDesk) GetFootRate() int32 {
	if m != nil && m.FootRate != nil {
		return *m.FootRate
	}
	return 0
}

func (m *PDdzDesk) GetPlayRate() int32 {
	if m != nil && m.PlayRate != nil {
		return *m.PlayRate
	}
	return 0
}

func (m *PDdzDesk) GetGameStatus() int32 {
	if m != nil && m.GameStatus != nil {
		return *m.GameStatus
	}
	return 0
}

func (m *PDdzDesk) GetInitRoomCoin() int64 {
	if m != nil && m.InitRoomCoin != nil {
		return *m.InitRoomCoin
	}
	return 0
}

func (m *PDdzDesk) GetCurrPlayCount() int32 {
	if m != nil && m.CurrPlayCount != nil {
		return *m.CurrPlayCount
	}
	return 0
}

func (m *PDdzDesk) GetTotalPlayCount() int32 {
	if m != nil && m.TotalPlayCount != nil {
		return *m.TotalPlayCount
	}
	return 0
}

func (m *PDdzDesk) GetPlayerNum() int32 {
	if m != nil && m.PlayerNum != nil {
		return *m.PlayerNum
	}
	return 0
}

// 游戏数据
type PGameData struct {
	HandPokers       []*PPokerPai     `protobuf:"bytes,1,rep,name=handPokers" json:"handPokers,omitempty"`
	OutPaiList       []*POutPokerPais `protobuf:"bytes,2,rep,name=outPaiList" json:"outPaiList,omitempty"`
	OutPai           *POutPokerPais   `protobuf:"bytes,3,opt,name=outPai" json:"outPai,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *PGameData) Reset()                    { *m = PGameData{} }
func (m *PGameData) String() string            { return proto.CompactTextString(m) }
func (*PGameData) ProtoMessage()               {}
func (*PGameData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *PGameData) GetHandPokers() []*PPokerPai {
	if m != nil {
		return m.HandPokers
	}
	return nil
}

func (m *PGameData) GetOutPaiList() []*POutPokerPais {
	if m != nil {
		return m.OutPaiList
	}
	return nil
}

func (m *PGameData) GetOutPai() *POutPokerPais {
	if m != nil {
		return m.OutPai
	}
	return nil
}

type PDdzBillBean struct {
	// 斗地主的账单
	Coin             *int64  `protobuf:"varint,1,opt,name=coin" json:"coin,omitempty"`
	LoseUser         *uint32 `protobuf:"varint,2,opt,name=loseUser" json:"loseUser,omitempty"`
	WinUser          *uint32 `protobuf:"varint,3,opt,name=winUser" json:"winUser,omitempty"`
	Desc             *string `protobuf:"bytes,4,opt,name=desc" json:"desc,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PDdzBillBean) Reset()                    { *m = PDdzBillBean{} }
func (m *PDdzBillBean) String() string            { return proto.CompactTextString(m) }
func (*PDdzBillBean) ProtoMessage()               {}
func (*PDdzBillBean) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *PDdzBillBean) GetCoin() int64 {
	if m != nil && m.Coin != nil {
		return *m.Coin
	}
	return 0
}

func (m *PDdzBillBean) GetLoseUser() uint32 {
	if m != nil && m.LoseUser != nil {
		return *m.LoseUser
	}
	return 0
}

func (m *PDdzBillBean) GetWinUser() uint32 {
	if m != nil && m.WinUser != nil {
		return *m.WinUser
	}
	return 0
}

func (m *PDdzBillBean) GetDesc() string {
	if m != nil && m.Desc != nil {
		return *m.Desc
	}
	return ""
}

type PDdzBill struct {
	// 斗地主的账单
	WinCoin          *int64          `protobuf:"varint,1,opt,name=winCoin" json:"winCoin,omitempty"`
	BillBean         []*PDdzBillBean `protobuf:"bytes,2,rep,name=billBean" json:"billBean,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *PDdzBill) Reset()                    { *m = PDdzBill{} }
func (m *PDdzBill) String() string            { return proto.CompactTextString(m) }
func (*PDdzBill) ProtoMessage()               {}
func (*PDdzBill) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *PDdzBill) GetWinCoin() int64 {
	if m != nil && m.WinCoin != nil {
		return *m.WinCoin
	}
	return 0
}

func (m *PDdzBill) GetBillBean() []*PDdzBillBean {
	if m != nil {
		return m.BillBean
	}
	return nil
}

// user
type PDdzUser struct {
	UserId           *uint32    `protobuf:"varint,1,opt,name=userId" json:"userId,omitempty"`
	GameData         *PGameData `protobuf:"bytes,2,opt,name=gameData" json:"gameData,omitempty"`
	Status           *int32     `protobuf:"varint,3,opt,name=status" json:"status,omitempty"`
	IsBreak          *bool      `protobuf:"varint,4,opt,name=isBreak" json:"isBreak,omitempty"`
	IsLeave          *bool      `protobuf:"varint,5,opt,name=isLeave" json:"isLeave,omitempty"`
	QiangDiZhuStatus *int32     `protobuf:"varint,7,opt,name=qiangDiZhuStatus" json:"qiangDiZhuStatus,omitempty"`
	JiabeiStatus     *int32     `protobuf:"varint,8,opt,name=jiabeiStatus" json:"jiabeiStatus,omitempty"`
	Bill             *PDdzBill  `protobuf:"bytes,9,opt,name=bill" json:"bill,omitempty"`
	DeskId           *int32     `protobuf:"varint,10,opt,name=deskId" json:"deskId,omitempty"`
	RoomId           *int32     `protobuf:"varint,11,opt,name=roomId" json:"roomId,omitempty"`
	Coin             *int64     `protobuf:"varint,12,opt,name=coin" json:"coin,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *PDdzUser) Reset()                    { *m = PDdzUser{} }
func (m *PDdzUser) String() string            { return proto.CompactTextString(m) }
func (*PDdzUser) ProtoMessage()               {}
func (*PDdzUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *PDdzUser) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *PDdzUser) GetGameData() *PGameData {
	if m != nil {
		return m.GameData
	}
	return nil
}

func (m *PDdzUser) GetStatus() int32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return 0
}

func (m *PDdzUser) GetIsBreak() bool {
	if m != nil && m.IsBreak != nil {
		return *m.IsBreak
	}
	return false
}

func (m *PDdzUser) GetIsLeave() bool {
	if m != nil && m.IsLeave != nil {
		return *m.IsLeave
	}
	return false
}

func (m *PDdzUser) GetQiangDiZhuStatus() int32 {
	if m != nil && m.QiangDiZhuStatus != nil {
		return *m.QiangDiZhuStatus
	}
	return 0
}

func (m *PDdzUser) GetJiabeiStatus() int32 {
	if m != nil && m.JiabeiStatus != nil {
		return *m.JiabeiStatus
	}
	return 0
}

func (m *PDdzUser) GetBill() *PDdzBill {
	if m != nil {
		return m.Bill
	}
	return nil
}

func (m *PDdzUser) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *PDdzUser) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *PDdzUser) GetCoin() int64 {
	if m != nil && m.Coin != nil {
		return *m.Coin
	}
	return 0
}

// room
type PDdzRoom struct {
	RoomId           *int32 `protobuf:"varint,1,opt,name=roomId" json:"roomId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PDdzRoom) Reset()                    { *m = PDdzRoom{} }
func (m *PDdzRoom) String() string            { return proto.CompactTextString(m) }
func (*PDdzRoom) ProtoMessage()               {}
func (*PDdzRoom) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *PDdzRoom) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

// 备份专用...
type PDdzbak struct {
	Desk             *PDdzDesk   `protobuf:"bytes,1,opt,name=desk" json:"desk,omitempty"`
	Users            []*PDdzUser `protobuf:"bytes,2,rep,name=users" json:"users,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *PDdzbak) Reset()                    { *m = PDdzbak{} }
func (m *PDdzbak) String() string            { return proto.CompactTextString(m) }
func (*PDdzbak) ProtoMessage()               {}
func (*PDdzbak) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *PDdzbak) GetDesk() *PDdzDesk {
	if m != nil {
		return m.Desk
	}
	return nil
}

func (m *PDdzbak) GetUsers() []*PDdzUser {
	if m != nil {
		return m.Users
	}
	return nil
}

// session
type PDdzSession struct {
	UserId           *uint32 `protobuf:"varint,1,opt,name=UserId" json:"UserId,omitempty"`
	RoomId           *int32  `protobuf:"varint,2,opt,name=RoomId" json:"RoomId,omitempty"`
	DeskId           *int32  `protobuf:"varint,3,opt,name=DeskId" json:"DeskId,omitempty"`
	GameStatus       *int32  `protobuf:"varint,4,opt,name=GameStatus" json:"GameStatus,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PDdzSession) Reset()                    { *m = PDdzSession{} }
func (m *PDdzSession) String() string            { return proto.CompactTextString(m) }
func (*PDdzSession) ProtoMessage()               {}
func (*PDdzSession) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *PDdzSession) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *PDdzSession) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *PDdzSession) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *PDdzSession) GetGameStatus() int32 {
	if m != nil && m.GameStatus != nil {
		return *m.GameStatus
	}
	return 0
}

func init() {
	proto.RegisterType((*PPokerPai)(nil), "doudizhu.PPokerPai")
	proto.RegisterType((*POutPokerPais)(nil), "doudizhu.POutPokerPais")
	proto.RegisterType((*PDdzDeskTongJi)(nil), "doudizhu.PDdzDeskTongJi")
	proto.RegisterType((*PDdzDesk)(nil), "doudizhu.PDdzDesk")
	proto.RegisterType((*PGameData)(nil), "doudizhu.PGameData")
	proto.RegisterType((*PDdzBillBean)(nil), "doudizhu.PDdzBillBean")
	proto.RegisterType((*PDdzBill)(nil), "doudizhu.PDdzBill")
	proto.RegisterType((*PDdzUser)(nil), "doudizhu.PDdzUser")
	proto.RegisterType((*PDdzRoom)(nil), "doudizhu.PDdzRoom")
	proto.RegisterType((*PDdzbak)(nil), "doudizhu.PDdzbak")
	proto.RegisterType((*PDdzSession)(nil), "doudizhu.PDdzSession")
	proto.RegisterEnum("doudizhu.EDDZTYPE", EDDZTYPE_name, EDDZTYPE_value)
}

var fileDescriptor0 = []byte{
	// 885 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x55, 0xc1, 0x52, 0x1b, 0x47,
	0x10, 0x2d, 0x21, 0x04, 0xab, 0x16, 0x02, 0x79, 0xc1, 0x30, 0x71, 0x2e, 0x8e, 0xaa, 0x92, 0xa8,
	0x92, 0x2a, 0x0e, 0x3e, 0xe6, 0x16, 0xb4, 0x24, 0x96, 0x0b, 0x13, 0x19, 0x50, 0xaa, 0xcc, 0x85,
	0x1a, 0xb1, 0x63, 0x18, 0x6b, 0xb5, 0x23, 0xef, 0xce, 0x42, 0xd0, 0x27, 0xe4, 0xe3, 0xf2, 0x27,
	0x39, 0xe5, 0x07, 0xd2, 0xdd, 0xb3, 0x23, 0xad, 0x6c, 0xc3, 0x4d, 0xf3, 0xa6, 0xb7, 0xfb, 0xf5,
	0x9b, 0xd7, 0x2d, 0xe8, 0xc4, 0xf1, 0xfc, 0x2a, 0x57, 0xd9, 0x9d, 0xca, 0x0e, 0x67, 0x99, 0xb1,
	0x26, 0x0c, 0x62, 0x53, 0xc4, 0x7a, 0x7e, 0x5b, 0x74, 0xdf, 0x41, 0x73, 0x38, 0x34, 0x13, 0x95,
	0x0d, 0xa5, 0x0e, 0x01, 0xd6, 0x74, 0x2c, 0x6a, 0x2f, 0x6b, 0xbd, 0x46, 0xd8, 0x82, 0x7a, 0xac,
	0x72, 0xb1, 0x86, 0x87, 0x66, 0xd8, 0x86, 0xc6, 0x9d, 0x4c, 0x0a, 0x25, 0xea, 0x7c, 0xb7, 0x0d,
	0x1b, 0x1f, 0x12, 0x73, 0xaf, 0x32, 0xb1, 0xce, 0xe7, 0x2d, 0x58, 0x4f, 0xe5, 0x54, 0x89, 0x06,
	0x05, 0x77, 0xff, 0xa9, 0x41, 0x7b, 0xf8, 0x47, 0x61, 0x7d, 0xda, 0x3c, 0xec, 0x40, 0x30, 0x51,
	0x0f, 0x7f, 0x72, 0x06, 0x97, 0xfd, 0x07, 0x68, 0xce, 0xfc, 0x35, 0xd6, 0xa8, 0xf7, 0x5a, 0xaf,
	0x76, 0x0f, 0x3d, 0xa9, 0xc3, 0x25, 0x23, 0xcc, 0x6c, 0x1f, 0x66, 0x95, 0xba, 0x3a, 0x3f, 0x32,
	0xd3, 0x31, 0xd7, 0x0d, 0x42, 0x24, 0x7c, 0x6d, 0x8a, 0xd4, 0x46, 0x85, 0x9e, 0x6b, 0xae, 0xde,
	0x08, 0x9f, 0x43, 0x9b, 0xb1, 0x73, 0x99, 0xce, 0x6f, 0x65, 0x7a, 0x23, 0x36, 0x18, 0xde, 0x83,
	0x2d, 0x07, 0x6b, 0x87, 0x6e, 0xae, 0xa0, 0xef, 0x4b, 0x34, 0xf0, 0x65, 0x0a, 0x94, 0x6b, 0x10,
	0x8b, 0x26, 0x9e, 0xdb, 0xa8, 0xd1, 0xf6, 0x30, 0x8a, 0xe7, 0x91, 0xca, 0x27, 0x17, 0x26, 0xbd,
	0x79, 0xa3, 0x91, 0x7e, 0x63, 0x8c, 0x34, 0x72, 0xec, 0x86, 0xa8, 0x1f, 0x54, 0xa8, 0xaf, 0x34,
	0x7e, 0x00, 0x3b, 0x9c, 0xff, 0x9d, 0xc6, 0xec, 0x91, 0xbe, 0xbc, 0x2d, 0x58, 0xd0, 0x46, 0xf7,
	0xdf, 0x75, 0x08, 0x7c, 0x4e, 0xaa, 0x87, 0x52, 0x4f, 0x06, 0x15, 0xe9, 0x51, 0xae, 0x52, 0xfa,
	0x7d, 0xd8, 0x26, 0x32, 0x7d, 0x4a, 0x73, 0xa2, 0xa7, 0xda, 0x96, 0x5a, 0xf4, 0xa0, 0x25, 0x93,
	0xc4, 0x97, 0x42, 0x41, 0x1e, 0xd5, 0xf0, 0x47, 0x80, 0x58, 0x2f, 0x02, 0x1b, 0x4f, 0x05, 0x6e,
	0x18, 0x64, 0x8f, 0x41, 0xa4, 0xd9, 0x13, 0x6d, 0xa1, 0x1d, 0xcc, 0x7d, 0x8a, 0xcf, 0x4f, 0x2a,
	0xb6, 0x49, 0x45, 0x8e, 0xc2, 0xbb, 0x11, 0x52, 0x65, 0x15, 0xdb, 0xe1, 0x2e, 0xb4, 0x62, 0xea,
	0x78, 0x54, 0x91, 0x92, 0x42, 0xe5, 0xb5, 0xd5, 0x77, 0xaa, 0x44, 0x81, 0xd1, 0x1e, 0x6c, 0x58,
	0x14, 0xf6, 0xa3, 0x16, 0x2d, 0x2e, 0x2c, 0x2a, 0x85, 0x57, 0x85, 0x7f, 0x06, 0xcd, 0xb1, 0xcc,
	0x95, 0xb3, 0xd2, 0x16, 0x06, 0xd7, 0x49, 0xa0, 0x4f, 0x9a, 0xd4, 0xc5, 0x70, 0x87, 0xb7, 0x19,
	0x47, 0xd3, 0xdd, 0xeb, 0xd4, 0x21, 0xdb, 0x8c, 0xec, 0xc0, 0x26, 0x4e, 0xc2, 0x05, 0xf9, 0x69,
	0x87, 0x35, 0xc4, 0x90, 0x33, 0x63, 0xa6, 0x8c, 0x74, 0x18, 0x41, 0xd2, 0x47, 0x46, 0x66, 0x71,
	0xce, 0x7a, 0x8b, 0x67, 0xde, 0x0f, 0x7d, 0x39, 0x7b, 0x2b, 0xff, 0x12, 0x21, 0xe7, 0x41, 0x12,
	0x83, 0xfc, 0x8d, 0x96, 0xe6, 0x37, 0x95, 0x8a, 0x5d, 0x76, 0x22, 0x86, 0x64, 0x98, 0x09, 0x3b,
	0xda, 0xf3, 0x99, 0x3f, 0x18, 0x63, 0xcf, 0xa4, 0x55, 0xe2, 0xb9, 0x47, 0x66, 0x89, 0x7c, 0x60,
	0x64, 0x9f, 0x11, 0x74, 0xef, 0xef, 0x38, 0x35, 0xe7, 0x56, 0xda, 0x22, 0x17, 0x07, 0xde, 0x90,
	0x3a, 0xd5, 0x96, 0x58, 0xf5, 0x8d, 0x4e, 0x85, 0xe0, 0x82, 0xe4, 0xe9, 0x22, 0xcb, 0x86, 0xf8,
	0xbd, 0xe3, 0xf5, 0x0d, 0x07, 0x63, 0xe7, 0xd6, 0x58, 0x99, 0x2c, 0xf1, 0x17, 0x8c, 0x23, 0x3f,
	0x2a, 0xa5, 0xb2, 0xd3, 0x62, 0x2a, 0xbe, 0x65, 0xbf, 0xfd, 0x5d, 0xc3, 0x39, 0xa7, 0x6a, 0x91,
	0xb4, 0x92, 0x1c, 0x81, 0x76, 0x8f, 0xf9, 0x41, 0xbd, 0x87, 0xbf, 0xea, 0x88, 0x9f, 0x01, 0x9c,
	0x23, 0x4e, 0x74, 0x6e, 0xcb, 0x39, 0x7d, 0xd4, 0x15, 0x4b, 0xfb, 0xd4, 0x9f, 0xb4, 0x4f, 0xf7,
	0x2d, 0x6c, 0xd1, 0xb3, 0x1e, 0xe9, 0x24, 0x39, 0x52, 0x32, 0xa5, 0x21, 0xbf, 0xa6, 0x66, 0x6b,
	0xfe, 0xdd, 0x12, 0x93, 0xb3, 0x41, 0x78, 0x04, 0xda, 0xf4, 0x6e, 0xf8, 0x92, 0x0c, 0xd4, 0x19,
	0xc0, 0x0f, 0x70, 0x60, 0xae, 0x79, 0x0b, 0x34, 0xbb, 0xc7, 0x6e, 0x94, 0x28, 0x5d, 0x19, 0xda,
	0x5f, 0x66, 0xeb, 0x41, 0x30, 0x2e, 0xeb, 0x94, 0xfc, 0xf7, 0x57, 0xcd, 0xe5, 0x59, 0x74, 0xff,
	0xab, 0xb9, 0x3c, 0x54, 0xa7, 0xb2, 0x02, 0x6a, 0x5c, 0xf1, 0x7b, 0x08, 0x6e, 0x4a, 0xf5, 0x98,
	0xd4, 0xaa, 0x5e, 0x0b, 0x61, 0xf1, 0xb3, 0xdc, 0x3d, 0xa7, 0x1b, 0x52, 0xa4, 0x83, 0x0b, 0x2b,
	0x53, 0x72, 0x52, 0x6e, 0x2c, 0x06, 0x4e, 0x94, 0xbc, 0x73, 0xcb, 0x32, 0x08, 0x05, 0x74, 0x3e,
	0x2d, 0x96, 0x43, 0x69, 0x85, 0xc5, 0x6e, 0xfa, 0xa8, 0xe5, 0x58, 0xe9, 0x12, 0x75, 0xbb, 0xe9,
	0x25, 0xac, 0x53, 0x3f, 0x3c, 0x4e, 0xad, 0x57, 0xe1, 0x97, 0xbd, 0x54, 0xb6, 0x09, 0x78, 0xf7,
	0x96, 0xd6, 0x6c, 0xf9, 0x65, 0xcd, 0x6a, 0xf3, 0xf4, 0x74, 0x5f, 0xb8, 0xa6, 0xc9, 0x70, 0x95,
	0x48, 0xde, 0x43, 0xdd, 0x53, 0xd8, 0xa4, 0xbb, 0xb1, 0x9c, 0x50, 0x59, 0x4a, 0xca, 0x17, 0x5f,
	0x94, 0xe5, 0x25, 0xf6, 0x1d, 0x34, 0x48, 0x31, 0xbf, 0xcd, 0x3f, 0x0b, 0x21, 0x51, 0x71, 0x8f,
	0xb6, 0xe8, 0xf7, 0xb9, 0xca, 0x73, 0x6d, 0x52, 0x2a, 0x37, 0xaa, 0x6a, 0x8c, 0xe7, 0x33, 0x57,
	0x7e, 0xcd, 0x13, 0x8f, 0x5c, 0x23, 0xf5, 0xaf, 0xcc, 0x0b, 0xff, 0xf3, 0xfc, 0xf4, 0x0b, 0x04,
	0xc7, 0x51, 0x74, 0x79, 0xf1, 0x7e, 0x78, 0x8c, 0x52, 0xee, 0xe1, 0xcf, 0x2b, 0xfa, 0x7d, 0xf5,
	0x7a, 0xf4, 0xeb, 0xe9, 0xc9, 0x71, 0x34, 0xb8, 0x7c, 0x3d, 0xea, 0xd4, 0x50, 0xca, 0xce, 0xe2,
	0xe6, 0x7c, 0xd0, 0xa7, 0xcb, 0xce, 0xda, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x41, 0xaa, 0xda,
	0x25, 0x17, 0x07, 0x00, 0x00,
}
