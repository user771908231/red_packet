// Code generated by protoc-gen-go.
// source: mahjong_hall.proto
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

// Ignoring public import of EProtoId from base.proto

// Ignoring public import of DDErrorCode from base.proto

// Ignoring public import of MJRoomType from base.proto

// 接入服务器
type Game_QuickConn struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_QuickConn) Reset()                    { *m = Game_QuickConn{} }
func (m *Game_QuickConn) String() string            { return proto.CompactTextString(m) }
func (*Game_QuickConn) ProtoMessage()               {}
func (*Game_QuickConn) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *Game_QuickConn) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type Game_AckQuickConn struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckQuickConn) Reset()                    { *m = Game_AckQuickConn{} }
func (m *Game_AckQuickConn) String() string            { return proto.CompactTextString(m) }
func (*Game_AckQuickConn) ProtoMessage()               {}
func (*Game_AckQuickConn) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *Game_AckQuickConn) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

// 游戏登录请求
type Game_Login struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	ProtoVersion     *int32       `protobuf:"varint,3,opt,name=protoVersion" json:"protoVersion,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_Login) Reset()                    { *m = Game_Login{} }
func (m *Game_Login) String() string            { return proto.CompactTextString(m) }
func (*Game_Login) ProtoMessage()               {}
func (*Game_Login) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *Game_Login) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_Login) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_Login) GetProtoVersion() int32 {
	if m != nil && m.ProtoVersion != nil {
		return *m.ProtoVersion
	}
	return 0
}

// 游戏登录回复
type Game_AckLogin struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Notice           *string      `protobuf:"bytes,2,opt,name=notice" json:"notice,omitempty"`
	GameStatus       *int32       `protobuf:"varint,3,opt,name=gameStatus" json:"gameStatus,omitempty"`
	MatchId          *int32       `protobuf:"varint,4,opt,name=matchId" json:"matchId,omitempty"`
	TableId          *int32       `protobuf:"varint,5,opt,name=tableId" json:"tableId,omitempty"`
	RoomPassword     *string      `protobuf:"bytes,6,opt,name=roomPassword" json:"roomPassword,omitempty"`
	CostCreateRoom   *int64       `protobuf:"varint,7,opt,name=costCreateRoom" json:"costCreateRoom,omitempty"`
	CostRebuy        *int64       `protobuf:"varint,8,opt,name=costRebuy" json:"costRebuy,omitempty"`
	Championship     *bool        `protobuf:"varint,9,opt,name=championship" json:"championship,omitempty"`
	Chip             *int64       `protobuf:"varint,10,opt,name=chip" json:"chip,omitempty"`
	MailCount        *int32       `protobuf:"varint,11,opt,name=mailCount" json:"mailCount,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckLogin) Reset()                    { *m = Game_AckLogin{} }
func (m *Game_AckLogin) String() string            { return proto.CompactTextString(m) }
func (*Game_AckLogin) ProtoMessage()               {}
func (*Game_AckLogin) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *Game_AckLogin) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckLogin) GetNotice() string {
	if m != nil && m.Notice != nil {
		return *m.Notice
	}
	return ""
}

func (m *Game_AckLogin) GetGameStatus() int32 {
	if m != nil && m.GameStatus != nil {
		return *m.GameStatus
	}
	return 0
}

func (m *Game_AckLogin) GetMatchId() int32 {
	if m != nil && m.MatchId != nil {
		return *m.MatchId
	}
	return 0
}

func (m *Game_AckLogin) GetTableId() int32 {
	if m != nil && m.TableId != nil {
		return *m.TableId
	}
	return 0
}

func (m *Game_AckLogin) GetRoomPassword() string {
	if m != nil && m.RoomPassword != nil {
		return *m.RoomPassword
	}
	return ""
}

func (m *Game_AckLogin) GetCostCreateRoom() int64 {
	if m != nil && m.CostCreateRoom != nil {
		return *m.CostCreateRoom
	}
	return 0
}

func (m *Game_AckLogin) GetCostRebuy() int64 {
	if m != nil && m.CostRebuy != nil {
		return *m.CostRebuy
	}
	return 0
}

func (m *Game_AckLogin) GetChampionship() bool {
	if m != nil && m.Championship != nil {
		return *m.Championship
	}
	return false
}

func (m *Game_AckLogin) GetChip() int64 {
	if m != nil && m.Chip != nil {
		return *m.Chip
	}
	return 0
}

func (m *Game_AckLogin) GetMailCount() int32 {
	if m != nil && m.MailCount != nil {
		return *m.MailCount
	}
	return 0
}

type Game_Notice struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	NoticeType       *int32       `protobuf:"varint,2,opt,name=noticeType" json:"noticeType,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_Notice) Reset()                    { *m = Game_Notice{} }
func (m *Game_Notice) String() string            { return proto.CompactTextString(m) }
func (*Game_Notice) ProtoMessage()               {}
func (*Game_Notice) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *Game_Notice) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_Notice) GetNoticeType() int32 {
	if m != nil && m.NoticeType != nil {
		return *m.NoticeType
	}
	return 0
}

// 公告的内容
type Game_AckNotice struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	NoticeType       *int32       `protobuf:"varint,2,opt,name=noticeType" json:"noticeType,omitempty"`
	NoticeTitle      *string      `protobuf:"bytes,3,opt,name=noticeTitle" json:"noticeTitle,omitempty"`
	NoticeContent    *string      `protobuf:"bytes,4,opt,name=noticeContent" json:"noticeContent,omitempty"`
	NoticeMemo       *string      `protobuf:"bytes,5,opt,name=noticeMemo" json:"noticeMemo,omitempty"`
	Id               *int32       `protobuf:"varint,6,opt,name=id" json:"id,omitempty"`
	Fileds           []string     `protobuf:"bytes,7,rep,name=fileds" json:"fileds,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckNotice) Reset()                    { *m = Game_AckNotice{} }
func (m *Game_AckNotice) String() string            { return proto.CompactTextString(m) }
func (*Game_AckNotice) ProtoMessage()               {}
func (*Game_AckNotice) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *Game_AckNotice) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckNotice) GetNoticeType() int32 {
	if m != nil && m.NoticeType != nil {
		return *m.NoticeType
	}
	return 0
}

func (m *Game_AckNotice) GetNoticeTitle() string {
	if m != nil && m.NoticeTitle != nil {
		return *m.NoticeTitle
	}
	return ""
}

func (m *Game_AckNotice) GetNoticeContent() string {
	if m != nil && m.NoticeContent != nil {
		return *m.NoticeContent
	}
	return ""
}

func (m *Game_AckNotice) GetNoticeMemo() string {
	if m != nil && m.NoticeMemo != nil {
		return *m.NoticeMemo
	}
	return ""
}

func (m *Game_AckNotice) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Game_AckNotice) GetFileds() []string {
	if m != nil {
		return m.Fileds
	}
	return nil
}

// 游戏战绩
type Game_GameRecord struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Id               *int32       `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	UserId           *uint32      `protobuf:"varint,3,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_GameRecord) Reset()                    { *m = Game_GameRecord{} }
func (m *Game_GameRecord) String() string            { return proto.CompactTextString(m) }
func (*Game_GameRecord) ProtoMessage()               {}
func (*Game_GameRecord) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6} }

func (m *Game_GameRecord) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_GameRecord) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Game_GameRecord) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

//
type Game_BeanUserRecord struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32      `protobuf:"varint,2,opt,name=UserId" json:"UserId,omitempty"`
	NickName         *string      `protobuf:"bytes,3,opt,name=NickName" json:"NickName,omitempty"`
	WinAmount        *int64       `protobuf:"varint,4,opt,name=WinAmount" json:"WinAmount,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_BeanUserRecord) Reset()                    { *m = Game_BeanUserRecord{} }
func (m *Game_BeanUserRecord) String() string            { return proto.CompactTextString(m) }
func (*Game_BeanUserRecord) ProtoMessage()               {}
func (*Game_BeanUserRecord) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{7} }

func (m *Game_BeanUserRecord) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_BeanUserRecord) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_BeanUserRecord) GetNickName() string {
	if m != nil && m.NickName != nil {
		return *m.NickName
	}
	return ""
}

func (m *Game_BeanUserRecord) GetWinAmount() int64 {
	if m != nil && m.WinAmount != nil {
		return *m.WinAmount
	}
	return 0
}

//
type Game_BeanGameRecord struct {
	Header           *ProtoHeader           `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Id               *int32                 `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	DeskId           *int32                 `protobuf:"varint,3,opt,name=deskId" json:"deskId,omitempty"`
	BeginTime        *string                `protobuf:"bytes,4,opt,name=beginTime" json:"beginTime,omitempty"`
	Users            []*Game_BeanUserRecord `protobuf:"bytes,5,rep,name=users" json:"users,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *Game_BeanGameRecord) Reset()                    { *m = Game_BeanGameRecord{} }
func (m *Game_BeanGameRecord) String() string            { return proto.CompactTextString(m) }
func (*Game_BeanGameRecord) ProtoMessage()               {}
func (*Game_BeanGameRecord) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{8} }

func (m *Game_BeanGameRecord) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_BeanGameRecord) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Game_BeanGameRecord) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *Game_BeanGameRecord) GetBeginTime() string {
	if m != nil && m.BeginTime != nil {
		return *m.BeginTime
	}
	return ""
}

func (m *Game_BeanGameRecord) GetUsers() []*Game_BeanUserRecord {
	if m != nil {
		return m.Users
	}
	return nil
}

//
type Game_AckGameRecord struct {
	Header           *ProtoHeader           `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	UserId           *uint32                `protobuf:"varint,2,opt,name=userId" json:"userId,omitempty"`
	Records          []*Game_BeanGameRecord `protobuf:"bytes,3,rep,name=records" json:"records,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *Game_AckGameRecord) Reset()                    { *m = Game_AckGameRecord{} }
func (m *Game_AckGameRecord) String() string            { return proto.CompactTextString(m) }
func (*Game_AckGameRecord) ProtoMessage()               {}
func (*Game_AckGameRecord) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{9} }

func (m *Game_AckGameRecord) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckGameRecord) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Game_AckGameRecord) GetRecords() []*Game_BeanGameRecord {
	if m != nil {
		return m.Records
	}
	return nil
}

// 反馈信息的协议
type Game_Feedback struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Message          *string      `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_Feedback) Reset()                    { *m = Game_Feedback{} }
func (m *Game_Feedback) String() string            { return proto.CompactTextString(m) }
func (*Game_Feedback) ProtoMessage()               {}
func (*Game_Feedback) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{10} }

func (m *Game_Feedback) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_Feedback) GetMessage() string {
	if m != nil && m.Message != nil {
		return *m.Message
	}
	return ""
}

// 创建房间
type Game_CreateRoom struct {
	Header           *ProtoHeader  `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	RoomTypeInfo     *RoomTypeInfo `protobuf:"bytes,2,opt,name=roomTypeInfo" json:"roomTypeInfo,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *Game_CreateRoom) Reset()                    { *m = Game_CreateRoom{} }
func (m *Game_CreateRoom) String() string            { return proto.CompactTextString(m) }
func (*Game_CreateRoom) ProtoMessage()               {}
func (*Game_CreateRoom) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{11} }

func (m *Game_CreateRoom) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_CreateRoom) GetRoomTypeInfo() *RoomTypeInfo {
	if m != nil {
		return m.RoomTypeInfo
	}
	return nil
}

// 创建房间回复的信息
type Game_AckCreateRoom struct {
	Header           *ProtoHeader  `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	DeskId           *int32        `protobuf:"varint,2,opt,name=deskId" json:"deskId,omitempty"`
	Password         *string       `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
	UserBalance      *int64        `protobuf:"varint,4,opt,name=userBalance" json:"userBalance,omitempty"`
	CreateFee        *int64        `protobuf:"varint,5,opt,name=createFee" json:"createFee,omitempty"`
	RoomTypeInfo     *RoomTypeInfo `protobuf:"bytes,6,opt,name=roomTypeInfo" json:"roomTypeInfo,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *Game_AckCreateRoom) Reset()                    { *m = Game_AckCreateRoom{} }
func (m *Game_AckCreateRoom) String() string            { return proto.CompactTextString(m) }
func (*Game_AckCreateRoom) ProtoMessage()               {}
func (*Game_AckCreateRoom) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{12} }

func (m *Game_AckCreateRoom) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_AckCreateRoom) GetDeskId() int32 {
	if m != nil && m.DeskId != nil {
		return *m.DeskId
	}
	return 0
}

func (m *Game_AckCreateRoom) GetPassword() string {
	if m != nil && m.Password != nil {
		return *m.Password
	}
	return ""
}

func (m *Game_AckCreateRoom) GetUserBalance() int64 {
	if m != nil && m.UserBalance != nil {
		return *m.UserBalance
	}
	return 0
}

func (m *Game_AckCreateRoom) GetCreateFee() int64 {
	if m != nil && m.CreateFee != nil {
		return *m.CreateFee
	}
	return 0
}

func (m *Game_AckCreateRoom) GetRoomTypeInfo() *RoomTypeInfo {
	if m != nil {
		return m.RoomTypeInfo
	}
	return nil
}

// 客户端请求进入 room, 服务器返回game_SendGameInfo
type Game_EnterRoom struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	MatchId          *int32       `protobuf:"varint,2,opt,name=matchId" json:"matchId,omitempty"`
	TableId          *int32       `protobuf:"varint,3,opt,name=tableId" json:"tableId,omitempty"`
	PassWord         *string      `protobuf:"bytes,4,opt,name=PassWord" json:"PassWord,omitempty"`
	UserId           *uint32      `protobuf:"varint,5,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_EnterRoom) Reset()                    { *m = Game_EnterRoom{} }
func (m *Game_EnterRoom) String() string            { return proto.CompactTextString(m) }
func (*Game_EnterRoom) ProtoMessage()               {}
func (*Game_EnterRoom) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{13} }

func (m *Game_EnterRoom) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Game_EnterRoom) GetMatchId() int32 {
	if m != nil && m.MatchId != nil {
		return *m.MatchId
	}
	return 0
}

func (m *Game_EnterRoom) GetTableId() int32 {
	if m != nil && m.TableId != nil {
		return *m.TableId
	}
	return 0
}

func (m *Game_EnterRoom) GetPassWord() string {
	if m != nil && m.PassWord != nil {
		return *m.PassWord
	}
	return ""
}

func (m *Game_EnterRoom) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

type Game_AckEnterRoom struct {
	Header           *ProtoHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Game_AckEnterRoom) Reset()                    { *m = Game_AckEnterRoom{} }
func (m *Game_AckEnterRoom) String() string            { return proto.CompactTextString(m) }
func (*Game_AckEnterRoom) ProtoMessage()               {}
func (*Game_AckEnterRoom) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{14} }

func (m *Game_AckEnterRoom) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func init() {
	proto.RegisterType((*Game_QuickConn)(nil), "mjproto.game_QuickConn")
	proto.RegisterType((*Game_AckQuickConn)(nil), "mjproto.game_AckQuickConn")
	proto.RegisterType((*Game_Login)(nil), "mjproto.game_Login")
	proto.RegisterType((*Game_AckLogin)(nil), "mjproto.game_AckLogin")
	proto.RegisterType((*Game_Notice)(nil), "mjproto.game_Notice")
	proto.RegisterType((*Game_AckNotice)(nil), "mjproto.game_AckNotice")
	proto.RegisterType((*Game_GameRecord)(nil), "mjproto.game_GameRecord")
	proto.RegisterType((*Game_BeanUserRecord)(nil), "mjproto.game_BeanUserRecord")
	proto.RegisterType((*Game_BeanGameRecord)(nil), "mjproto.game_BeanGameRecord")
	proto.RegisterType((*Game_AckGameRecord)(nil), "mjproto.game_AckGameRecord")
	proto.RegisterType((*Game_Feedback)(nil), "mjproto.game_Feedback")
	proto.RegisterType((*Game_CreateRoom)(nil), "mjproto.game_CreateRoom")
	proto.RegisterType((*Game_AckCreateRoom)(nil), "mjproto.game_AckCreateRoom")
	proto.RegisterType((*Game_EnterRoom)(nil), "mjproto.game_EnterRoom")
	proto.RegisterType((*Game_AckEnterRoom)(nil), "mjproto.game_AckEnterRoom")
}

var fileDescriptor2 = []byte{
	// 644 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x55, 0xcf, 0x4e, 0xdb, 0x4e,
	0x10, 0xfe, 0x39, 0xc6, 0x49, 0x3c, 0x21, 0xe1, 0xc7, 0x02, 0x95, 0x55, 0xf5, 0x50, 0x59, 0x3d,
	0x20, 0xa1, 0xe6, 0xc0, 0xa1, 0x52, 0x8f, 0x24, 0x2a, 0x7f, 0xa4, 0x36, 0x4a, 0x03, 0x94, 0xde,
	0xd0, 0xc6, 0x5e, 0x62, 0x13, 0xdb, 0x6b, 0x79, 0x37, 0xaa, 0xd2, 0x17, 0xe9, 0x1b, 0xf4, 0xde,
	0x67, 0xeb, 0x0b, 0x74, 0x76, 0xd7, 0x26, 0x11, 0xfd, 0xa3, 0x58, 0x5c, 0x2c, 0xf8, 0x76, 0xe7,
	0x9b, 0xf9, 0x66, 0xe6, 0xdb, 0x00, 0x49, 0x69, 0x74, 0xcf, 0xb3, 0xd9, 0x6d, 0x44, 0x93, 0xa4,
	0x9f, 0x17, 0x5c, 0x72, 0xd2, 0x4a, 0xef, 0xf5, 0x1f, 0xcf, 0x61, 0x4a, 0x05, 0x33, 0xa0, 0xff,
	0x06, 0x7a, 0x33, 0x9a, 0xb2, 0xdb, 0x8f, 0x8b, 0x38, 0x98, 0x0f, 0x79, 0x96, 0x91, 0x57, 0xd0,
	0x8c, 0x18, 0x0d, 0x59, 0xe1, 0x59, 0x2f, 0xad, 0xc3, 0xce, 0xf1, 0x7e, 0xbf, 0x8c, 0xeb, 0x8f,
	0xd5, 0xf7, 0x5c, 0x9f, 0xf9, 0x6f, 0x61, 0x57, 0xc7, 0x9d, 0x04, 0xf3, 0xba, 0xa1, 0x9f, 0x01,
	0x74, 0xe8, 0x7b, 0x3e, 0x8b, 0x37, 0x8c, 0x21, 0x3d, 0x68, 0x2e, 0x04, 0x2b, 0x2e, 0x42, 0xaf,
	0x81, 0xb7, 0xba, 0x64, 0x1f, 0xb6, 0xf5, 0xa5, 0x4f, 0xac, 0x10, 0x31, 0xcf, 0x3c, 0x1b, 0x51,
	0xc7, 0xff, 0x69, 0x41, 0xb7, 0xaa, 0xaa, 0x26, 0x7b, 0xc6, 0x65, 0x1c, 0x30, 0xcd, 0xee, 0x12,
	0x62, 0x2a, 0xbc, 0x94, 0x54, 0x2e, 0x84, 0xe1, 0x26, 0x3b, 0xd0, 0x4a, 0xa9, 0x0c, 0x22, 0x2c,
	0x61, 0xab, 0x02, 0x24, 0x9d, 0x26, 0x0c, 0x01, 0x47, 0x03, 0x58, 0x53, 0xc1, 0x79, 0x3a, 0xa6,
	0x42, 0x7c, 0xe1, 0x45, 0xe8, 0x35, 0x35, 0xd7, 0x33, 0xe8, 0x05, 0x5c, 0xc8, 0x61, 0xc1, 0xa8,
	0x64, 0x13, 0x3c, 0xf7, 0x5a, 0x88, 0xdb, 0x64, 0x17, 0x5c, 0x85, 0x4f, 0xd8, 0x74, 0xb1, 0xf4,
	0xda, 0x1a, 0x42, 0x82, 0x20, 0xa2, 0x69, 0x8e, 0x82, 0x44, 0x14, 0xe7, 0x9e, 0x8b, 0x68, 0x9b,
	0x6c, 0xc3, 0x56, 0xa0, 0xfe, 0x83, 0x2a, 0x2c, 0xa5, 0x71, 0x32, 0xe4, 0x8b, 0x4c, 0x7a, 0x1d,
	0xad, 0xfa, 0x0c, 0x3a, 0x5a, 0xf4, 0x48, 0x4b, 0xd8, 0x50, 0x32, 0x4a, 0x34, 0x92, 0xaf, 0x96,
	0xb9, 0x91, 0xed, 0xf8, 0xdf, 0xad, 0x72, 0x19, 0xb0, 0x7d, 0x4f, 0x25, 0x23, 0x7b, 0xd0, 0x29,
	0xb1, 0x58, 0x26, 0x4c, 0x37, 0xd1, 0x25, 0x07, 0xd0, 0x35, 0x20, 0xae, 0x8b, 0x64, 0xa8, 0x60,
	0xab, 0xea, 0xb7, 0x81, 0x3f, 0xb0, 0x94, 0xeb, 0x6e, 0xba, 0x08, 0x35, 0x62, 0xd3, 0x43, 0x47,
	0xcd, 0xe7, 0x2e, 0x4e, 0x58, 0x28, 0xb0, 0x77, 0xf6, 0xa1, 0xeb, 0x5f, 0xc2, 0x8e, 0xae, 0xf3,
	0x0c, 0x3f, 0x13, 0x16, 0x60, 0xb3, 0x37, 0x2c, 0xd4, 0x90, 0x36, 0x2a, 0xd2, 0x72, 0xa5, 0x54,
	0x6d, 0x5d, 0x3f, 0x87, 0x3d, 0x4d, 0x3a, 0x60, 0x34, 0xbb, 0xc6, 0x83, 0x5a, 0xc4, 0x48, 0x76,
	0xbd, 0xbe, 0x9f, 0xff, 0x43, 0x7b, 0x84, 0xae, 0x18, 0x21, 0x61, 0x29, 0x1d, 0x07, 0x77, 0x13,
	0x67, 0x27, 0xa9, 0x1e, 0x9c, 0x92, 0x6d, 0xfb, 0xdf, 0xac, 0xb5, 0x94, 0x4f, 0xd5, 0x12, 0x32,
	0x31, 0x2f, 0xb5, 0x38, 0x2a, 0xd9, 0x94, 0xe1, 0xfe, 0x5f, 0xc5, 0x98, 0xdf, 0xf4, 0xf8, 0x08,
	0x1c, 0x25, 0x57, 0x60, 0x7b, 0x6d, 0xe4, 0x7c, 0xf1, 0xc0, 0xf9, 0x07, 0xd1, 0xfe, 0x12, 0x48,
	0xb5, 0x08, 0xb5, 0xeb, 0x7a, 0x6c, 0xd5, 0xd7, 0xd0, 0x2a, 0x74, 0xbc, 0x72, 0xd2, 0x5f, 0x52,
	0xaf, 0x92, 0xf8, 0xa7, 0xa5, 0x85, 0x4f, 0x19, 0x0b, 0xa7, 0x34, 0x98, 0x6f, 0x98, 0x55, 0xd9,
	0x93, 0x09, 0x41, 0x67, 0xa5, 0x87, 0xfd, 0xb0, 0xdc, 0x91, 0x95, 0xf1, 0x36, 0x64, 0x3a, 0x32,
	0x36, 0x56, 0xab, 0x7c, 0x91, 0xdd, 0x71, 0x4d, 0xd7, 0x39, 0x3e, 0x78, 0xb8, 0x3b, 0x59, 0x3b,
	0xf4, 0x7f, 0x58, 0xab, 0x4e, 0xd5, 0xce, 0xb4, 0x9a, 0x9a, 0x99, 0x22, 0x2e, 0x4d, 0x5e, 0x3d,
	0x1e, 0x66, 0x69, 0xd0, 0x44, 0xaa, 0x97, 0x03, 0x9a, 0xd0, 0x2c, 0x30, 0x93, 0x34, 0x2f, 0x87,
	0x4e, 0x85, 0x2d, 0xd2, 0x66, 0xb1, 0x7f, 0xab, 0xb9, 0xf9, 0xaf, 0x9a, 0xbf, 0x96, 0x2e, 0x7f,
	0x87, 0x16, 0x2c, 0x6a, 0x94, 0xbb, 0xf6, 0x02, 0x36, 0x1e, 0xbf, 0x80, 0x76, 0x25, 0x40, 0xbd,
	0x7e, 0x37, 0x4a, 0x80, 0xd9, 0xba, 0xd5, 0x32, 0x38, 0xda, 0x64, 0x6b, 0x3f, 0x1b, 0x35, 0xd3,
	0x0f, 0x1a, 0xe7, 0xf6, 0xf8, 0xbf, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe1, 0xcc, 0xa1, 0x3e,
	0xd9, 0x06, 0x00, 0x00,
}
