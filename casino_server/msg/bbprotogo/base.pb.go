// Code generated by protoc-gen-go.
// source: base.proto
// DO NOT EDIT!

package bbproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EProtoId int32

const (
	EProtoId_GETINTOROOM     EProtoId = 1
	EProtoId_SHUIGUOJI       EProtoId = 2
	EProtoId_SHUIGUOJIHILOMP EProtoId = 3
	EProtoId_SHUIGUOJIRES    EProtoId = 4
	// 扎金花相关的逻辑
	EProtoId_ZJHROOM                          EProtoId = 5
	EProtoId_ZJHBET                           EProtoId = 6
	EProtoId_ZJHMSG                           EProtoId = 7
	EProtoId_ZJHQUERYNOSEATUSER               EProtoId = 8
	EProtoId_ZJHREQSEAT                       EProtoId = 9
	EProtoId_ZJHLOTTERY                       EProtoId = 10
	EProtoId_ZJHBROADCASTBEGINBET             EProtoId = 11
	EProtoId_LOGINSIGNINBONUS                 EProtoId = 12
	EProtoId_LOGINTURNTABLEBONUS              EProtoId = 13
	EProtoId_OLINEBONUS                       EProtoId = 14
	EProtoId_TIMINGBONUS                      EProtoId = 15
	EProtoId_THROOM                           EProtoId = 16
	EProtoId_THBET                            EProtoId = 17
	EProtoId_THBETBEGIN                       EProtoId = 18
	EProtoId_THBETBROADCAST                   EProtoId = 19
	EProtoId_THROOMADDUSERBROADCAST           EProtoId = 20
	EProtoId_PID_REQQUICKCONN                 EProtoId = 21
	EProtoId_PID_ACKQUICKCONN                 EProtoId = 22
	EProtoId_PID_NULLMSG                      EProtoId = 23
	EProtoId_PID_MATCHLIST_REQMOBILEMATCHLIST EProtoId = 24
	EProtoId_PID_GAME_LOGINGAME               EProtoId = 25
	EProtoId_PID_GAME_ENTERMATCH              EProtoId = 26
	EProtoId_PID_GAME_ACKENTERMATCH           EProtoId = 27
	EProtoId_PID_GAME_SENDGAMEINFO            EProtoId = 28
	EProtoId_PID_GAME_BLINDCOIN               EProtoId = 29
	EProtoId_PID_GAME_INITCARD                EProtoId = 30
	EProtoId_PID_GAME_SENDFLOPCARD            EProtoId = 31
	EProtoId_PID_GAME_SENDTURNCARD            EProtoId = 32
	EProtoId_PID_GAME_SENDRIVERCARD           EProtoId = 33
	EProtoId_PID_GAME_RAISEBET                EProtoId = 34
	EProtoId_PID_GAME_ACKRAISEBET             EProtoId = 35
	EProtoId_PID_GAME_FOLLOWBET               EProtoId = 36
	EProtoId_PID_GAME_ACKFOLLOWBET            EProtoId = 37
	EProtoId_PID_GAME_FOLDBET                 EProtoId = 38
	EProtoId_PID_GAME_ACKFOLDBET              EProtoId = 39
	EProtoId_PID_GAME_CHECKBET                EProtoId = 40
	EProtoId_PID_GAME_ACKCHECKBET             EProtoId = 41
	EProtoId_PID_GAME_SENDOVERTURN            EProtoId = 42
	EProtoId_PID_GAME_SENDADDUSER             EProtoId = 43
	EProtoId_PID_GAME_GAME_SHOWCARD           EProtoId = 44
	EProtoId_PID_GAME_GAME_ACKSHOWCARD        EProtoId = 45
	EProtoId_PID_GAME_GAME_TESTRESULT         EProtoId = 46
	EProtoId_PID_GAME_PRECOIN                 EProtoId = 47
	EProtoId_PID_GAME_GAMENOTICE              EProtoId = 48
	EProtoId_PID_GAME_GAME_ACKNOTICE          EProtoId = 49
	EProtoId_PID_GAME_GAME_CREATEDESK         EProtoId = 50
	EProtoId_PID_GAME_GAME_ACKCREATEDESK      EProtoId = 51
	EProtoId_PID_GAME_GAME_READY              EProtoId = 52
	EProtoId_PID_GAME_GAME_ACKREADY           EProtoId = 53
	EProtoId_PID_GAME_GAME_BEGIN              EProtoId = 54
	EProtoId_PID_GAME_GAME_GAMERECORD         EProtoId = 55
	EProtoId_PID_GAME_GAME_ACKGAMERECORD      EProtoId = 56
	EProtoId_PID_GAME_GAME_BEANGAMERECORD     EProtoId = 57
	EProtoId_PID_GAME_GAME_DISSOLVEDESK       EProtoId = 58
	EProtoId_PID_GAME_GAME_ACKDISSOLVEDESK    EProtoId = 59
	EProtoId_PID_GAME_LEAVEDESK               EProtoId = 60
	EProtoId_PID_GAME_ACKLEAVEDESK            EProtoId = 61
	EProtoId_PID_GAME_SENDDESKENDLOTTERY      EProtoId = 62
	EProtoId_PID_GAME_MESSAGE                 EProtoId = 63
	EProtoId_PID_GAME_SENDMESSAGE             EProtoId = 64
	EProtoId_PID_GAME_TOUNAMENTBLIND          EProtoId = 65
	EProtoId_PID_GAME_TOUNAMENTREWARDS        EProtoId = 66
	EProtoId_PID_GAME_TOUNAMENTRANK           EProtoId = 67
	EProtoId_PID_GAME_TOUNAMENTSUMMARY        EProtoId = 68
	EProtoId_PID_GAME_MATCHLIST               EProtoId = 69
	EProtoId_PID_GAME_TOUNAMENTPLAYERRANK     EProtoId = 70
	EProtoId_PID_GAME_REBUY                   EProtoId = 71
	EProtoId_PID_GAME_ACKREBUY                EProtoId = 72
	EProtoId_PID_GAME_LOGIN                   EProtoId = 73
	EProtoId_PID_GAME_ACKLOGIN                EProtoId = 74
	EProtoId_PID_GAME_FEEDBACK                EProtoId = 75
	EProtoId_PID_GAME_NOTREBUY                EProtoId = 76
	EProtoId_PID_GAME_ACKNOTREBUY             EProtoId = 77
	EProtoId_PID_GAME_SENDCHANGEDESKOWNER     EProtoId = 78
)

var EProtoId_name = map[int32]string{
	1:  "GETINTOROOM",
	2:  "SHUIGUOJI",
	3:  "SHUIGUOJIHILOMP",
	4:  "SHUIGUOJIRES",
	5:  "ZJHROOM",
	6:  "ZJHBET",
	7:  "ZJHMSG",
	8:  "ZJHQUERYNOSEATUSER",
	9:  "ZJHREQSEAT",
	10: "ZJHLOTTERY",
	11: "ZJHBROADCASTBEGINBET",
	12: "LOGINSIGNINBONUS",
	13: "LOGINTURNTABLEBONUS",
	14: "OLINEBONUS",
	15: "TIMINGBONUS",
	16: "THROOM",
	17: "THBET",
	18: "THBETBEGIN",
	19: "THBETBROADCAST",
	20: "THROOMADDUSERBROADCAST",
	21: "PID_REQQUICKCONN",
	22: "PID_ACKQUICKCONN",
	23: "PID_NULLMSG",
	24: "PID_MATCHLIST_REQMOBILEMATCHLIST",
	25: "PID_GAME_LOGINGAME",
	26: "PID_GAME_ENTERMATCH",
	27: "PID_GAME_ACKENTERMATCH",
	28: "PID_GAME_SENDGAMEINFO",
	29: "PID_GAME_BLINDCOIN",
	30: "PID_GAME_INITCARD",
	31: "PID_GAME_SENDFLOPCARD",
	32: "PID_GAME_SENDTURNCARD",
	33: "PID_GAME_SENDRIVERCARD",
	34: "PID_GAME_RAISEBET",
	35: "PID_GAME_ACKRAISEBET",
	36: "PID_GAME_FOLLOWBET",
	37: "PID_GAME_ACKFOLLOWBET",
	38: "PID_GAME_FOLDBET",
	39: "PID_GAME_ACKFOLDBET",
	40: "PID_GAME_CHECKBET",
	41: "PID_GAME_ACKCHECKBET",
	42: "PID_GAME_SENDOVERTURN",
	43: "PID_GAME_SENDADDUSER",
	44: "PID_GAME_GAME_SHOWCARD",
	45: "PID_GAME_GAME_ACKSHOWCARD",
	46: "PID_GAME_GAME_TESTRESULT",
	47: "PID_GAME_PRECOIN",
	48: "PID_GAME_GAMENOTICE",
	49: "PID_GAME_GAME_ACKNOTICE",
	50: "PID_GAME_GAME_CREATEDESK",
	51: "PID_GAME_GAME_ACKCREATEDESK",
	52: "PID_GAME_GAME_READY",
	53: "PID_GAME_GAME_ACKREADY",
	54: "PID_GAME_GAME_BEGIN",
	55: "PID_GAME_GAME_GAMERECORD",
	56: "PID_GAME_GAME_ACKGAMERECORD",
	57: "PID_GAME_GAME_BEANGAMERECORD",
	58: "PID_GAME_GAME_DISSOLVEDESK",
	59: "PID_GAME_GAME_ACKDISSOLVEDESK",
	60: "PID_GAME_LEAVEDESK",
	61: "PID_GAME_ACKLEAVEDESK",
	62: "PID_GAME_SENDDESKENDLOTTERY",
	63: "PID_GAME_MESSAGE",
	64: "PID_GAME_SENDMESSAGE",
	65: "PID_GAME_TOUNAMENTBLIND",
	66: "PID_GAME_TOUNAMENTREWARDS",
	67: "PID_GAME_TOUNAMENTRANK",
	68: "PID_GAME_TOUNAMENTSUMMARY",
	69: "PID_GAME_MATCHLIST",
	70: "PID_GAME_TOUNAMENTPLAYERRANK",
	71: "PID_GAME_REBUY",
	72: "PID_GAME_ACKREBUY",
	73: "PID_GAME_LOGIN",
	74: "PID_GAME_ACKLOGIN",
	75: "PID_GAME_FEEDBACK",
	76: "PID_GAME_NOTREBUY",
	77: "PID_GAME_ACKNOTREBUY",
	78: "PID_GAME_SENDCHANGEDESKOWNER",
}
var EProtoId_value = map[string]int32{
	"GETINTOROOM":                      1,
	"SHUIGUOJI":                        2,
	"SHUIGUOJIHILOMP":                  3,
	"SHUIGUOJIRES":                     4,
	"ZJHROOM":                          5,
	"ZJHBET":                           6,
	"ZJHMSG":                           7,
	"ZJHQUERYNOSEATUSER":               8,
	"ZJHREQSEAT":                       9,
	"ZJHLOTTERY":                       10,
	"ZJHBROADCASTBEGINBET":             11,
	"LOGINSIGNINBONUS":                 12,
	"LOGINTURNTABLEBONUS":              13,
	"OLINEBONUS":                       14,
	"TIMINGBONUS":                      15,
	"THROOM":                           16,
	"THBET":                            17,
	"THBETBEGIN":                       18,
	"THBETBROADCAST":                   19,
	"THROOMADDUSERBROADCAST":           20,
	"PID_REQQUICKCONN":                 21,
	"PID_ACKQUICKCONN":                 22,
	"PID_NULLMSG":                      23,
	"PID_MATCHLIST_REQMOBILEMATCHLIST": 24,
	"PID_GAME_LOGINGAME":               25,
	"PID_GAME_ENTERMATCH":              26,
	"PID_GAME_ACKENTERMATCH":           27,
	"PID_GAME_SENDGAMEINFO":            28,
	"PID_GAME_BLINDCOIN":               29,
	"PID_GAME_INITCARD":                30,
	"PID_GAME_SENDFLOPCARD":            31,
	"PID_GAME_SENDTURNCARD":            32,
	"PID_GAME_SENDRIVERCARD":           33,
	"PID_GAME_RAISEBET":                34,
	"PID_GAME_ACKRAISEBET":             35,
	"PID_GAME_FOLLOWBET":               36,
	"PID_GAME_ACKFOLLOWBET":            37,
	"PID_GAME_FOLDBET":                 38,
	"PID_GAME_ACKFOLDBET":              39,
	"PID_GAME_CHECKBET":                40,
	"PID_GAME_ACKCHECKBET":             41,
	"PID_GAME_SENDOVERTURN":            42,
	"PID_GAME_SENDADDUSER":             43,
	"PID_GAME_GAME_SHOWCARD":           44,
	"PID_GAME_GAME_ACKSHOWCARD":        45,
	"PID_GAME_GAME_TESTRESULT":         46,
	"PID_GAME_PRECOIN":                 47,
	"PID_GAME_GAMENOTICE":              48,
	"PID_GAME_GAME_ACKNOTICE":          49,
	"PID_GAME_GAME_CREATEDESK":         50,
	"PID_GAME_GAME_ACKCREATEDESK":      51,
	"PID_GAME_GAME_READY":              52,
	"PID_GAME_GAME_ACKREADY":           53,
	"PID_GAME_GAME_BEGIN":              54,
	"PID_GAME_GAME_GAMERECORD":         55,
	"PID_GAME_GAME_ACKGAMERECORD":      56,
	"PID_GAME_GAME_BEANGAMERECORD":     57,
	"PID_GAME_GAME_DISSOLVEDESK":       58,
	"PID_GAME_GAME_ACKDISSOLVEDESK":    59,
	"PID_GAME_LEAVEDESK":               60,
	"PID_GAME_ACKLEAVEDESK":            61,
	"PID_GAME_SENDDESKENDLOTTERY":      62,
	"PID_GAME_MESSAGE":                 63,
	"PID_GAME_SENDMESSAGE":             64,
	"PID_GAME_TOUNAMENTBLIND":          65,
	"PID_GAME_TOUNAMENTREWARDS":        66,
	"PID_GAME_TOUNAMENTRANK":           67,
	"PID_GAME_TOUNAMENTSUMMARY":        68,
	"PID_GAME_MATCHLIST":               69,
	"PID_GAME_TOUNAMENTPLAYERRANK":     70,
	"PID_GAME_REBUY":                   71,
	"PID_GAME_ACKREBUY":                72,
	"PID_GAME_LOGIN":                   73,
	"PID_GAME_ACKLOGIN":                74,
	"PID_GAME_FEEDBACK":                75,
	"PID_GAME_NOTREBUY":                76,
	"PID_GAME_ACKNOTREBUY":             77,
	"PID_GAME_SENDCHANGEDESKOWNER":     78,
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
func (EProtoId) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

type DDErrorCode int32

const (
	DDErrorCode_ERRORCODE_SUCC DDErrorCode = 0
	// -101   -200	游戏异常
	DDErrorCode_ERRORCODE_CREATE_DESK_DIAMOND_NOTENOUGH DDErrorCode = -101
	DDErrorCode_ERRORCODE_CREATE_DESK_USER_NOTFOUND     DDErrorCode = -102
	DDErrorCode_ERRORCODE_INTO_DESK_NOTFOUND            DDErrorCode = -103
	DDErrorCode_ERRORCODE_TOURNAMENT_CANNOT_JOIN        DDErrorCode = -104
	DDErrorCode_ERRORCODE_TOURNAMENT_CANNOT_REBUY       DDErrorCode = -105
	DDErrorCode_ERRORCODE_GAME_READY_REPEAT             DDErrorCode = -110
	DDErrorCode_ERRORCODE_GAME_READY_CHIP_NOT_ENOUGH    DDErrorCode = -111
)

var DDErrorCode_name = map[int32]string{
	0:    "ERRORCODE_SUCC",
	-101: "ERRORCODE_CREATE_DESK_DIAMOND_NOTENOUGH",
	-102: "ERRORCODE_CREATE_DESK_USER_NOTFOUND",
	-103: "ERRORCODE_INTO_DESK_NOTFOUND",
	-104: "ERRORCODE_TOURNAMENT_CANNOT_JOIN",
	-105: "ERRORCODE_TOURNAMENT_CANNOT_REBUY",
	-110: "ERRORCODE_GAME_READY_REPEAT",
	-111: "ERRORCODE_GAME_READY_CHIP_NOT_ENOUGH",
}
var DDErrorCode_value = map[string]int32{
	"ERRORCODE_SUCC":                          0,
	"ERRORCODE_CREATE_DESK_DIAMOND_NOTENOUGH": -101,
	"ERRORCODE_CREATE_DESK_USER_NOTFOUND":     -102,
	"ERRORCODE_INTO_DESK_NOTFOUND":            -103,
	"ERRORCODE_TOURNAMENT_CANNOT_JOIN":        -104,
	"ERRORCODE_TOURNAMENT_CANNOT_REBUY":       -105,
	"ERRORCODE_GAME_READY_REPEAT":             -110,
	"ERRORCODE_GAME_READY_CHIP_NOT_ENOUGH":    -111,
}

func (x DDErrorCode) Enum() *DDErrorCode {
	p := new(DDErrorCode)
	*p = x
	return p
}
func (x DDErrorCode) String() string {
	return proto.EnumName(DDErrorCode_name, int32(x))
}
func (x *DDErrorCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(DDErrorCode_value, data, "DDErrorCode")
	if err != nil {
		return err
	}
	*x = DDErrorCode(value)
	return nil
}
func (DDErrorCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

// general response protocol
type ProtoHeader struct {
	ApiVer           *string `protobuf:"bytes,1,opt,name=apiVer" json:"apiVer,omitempty"`
	SessionId        *string `protobuf:"bytes,2,opt,name=sessionId" json:"sessionId,omitempty"`
	UserId           *uint32 `protobuf:"varint,3,opt,name=userId" json:"userId,omitempty"`
	PacketId         *int32  `protobuf:"varint,4,opt,name=packetId" json:"packetId,omitempty"`
	Code             *int32  `protobuf:"varint,5,opt,name=code" json:"code,omitempty"`
	Error            *string `protobuf:"bytes,6,opt,name=error" json:"error,omitempty"`
	ExtraTag         *int32  `protobuf:"varint,7,opt,name=extraTag" json:"extraTag,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ProtoHeader) Reset()                    { *m = ProtoHeader{} }
func (m *ProtoHeader) String() string            { return proto.CompactTextString(m) }
func (*ProtoHeader) ProtoMessage()               {}
func (*ProtoHeader) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *ProtoHeader) GetApiVer() string {
	if m != nil && m.ApiVer != nil {
		return *m.ApiVer
	}
	return ""
}

func (m *ProtoHeader) GetSessionId() string {
	if m != nil && m.SessionId != nil {
		return *m.SessionId
	}
	return ""
}

func (m *ProtoHeader) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *ProtoHeader) GetPacketId() int32 {
	if m != nil && m.PacketId != nil {
		return *m.PacketId
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

func (m *ProtoHeader) GetExtraTag() int32 {
	if m != nil && m.ExtraTag != nil {
		return *m.ExtraTag
	}
	return 0
}

type TerminalInfo struct {
	Channel          *string `protobuf:"bytes,1,opt,name=channel" json:"channel,omitempty"`
	DeviceName       *string `protobuf:"bytes,2,opt,name=deviceName" json:"deviceName,omitempty"`
	Uuid             *string `protobuf:"bytes,3,opt,name=uuid" json:"uuid,omitempty"`
	Os               *string `protobuf:"bytes,4,opt,name=os" json:"os,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *TerminalInfo) Reset()                    { *m = TerminalInfo{} }
func (m *TerminalInfo) String() string            { return proto.CompactTextString(m) }
func (*TerminalInfo) ProtoMessage()               {}
func (*TerminalInfo) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *TerminalInfo) GetChannel() string {
	if m != nil && m.Channel != nil {
		return *m.Channel
	}
	return ""
}

func (m *TerminalInfo) GetDeviceName() string {
	if m != nil && m.DeviceName != nil {
		return *m.DeviceName
	}
	return ""
}

func (m *TerminalInfo) GetUuid() string {
	if m != nil && m.Uuid != nil {
		return *m.Uuid
	}
	return ""
}

func (m *TerminalInfo) GetOs() string {
	if m != nil && m.Os != nil {
		return *m.Os
	}
	return ""
}

type User struct {
	Mid                *string `protobuf:"bytes,1,opt,name=mid" json:"mid,omitempty"`
	Pwd                *string `protobuf:"bytes,2,opt,name=pwd" json:"pwd,omitempty"`
	Coin               *int64  `protobuf:"varint,3,opt,name=coin" json:"coin,omitempty"`
	Id                 *uint32 `protobuf:"varint,4,opt,name=id" json:"id,omitempty"`
	NickName           *string `protobuf:"bytes,5,opt,name=nickName" json:"nickName,omitempty"`
	LoginTurntable     *bool   `protobuf:"varint,6,opt,name=loginTurntable" json:"loginTurntable,omitempty"`
	LoginTurntableTime *string `protobuf:"bytes,7,opt,name=loginTurntableTime" json:"loginTurntableTime,omitempty"`
	Scores             *int32  `protobuf:"varint,8,opt,name=scores" json:"scores,omitempty"`
	LastSignTime       *string `protobuf:"bytes,9,opt,name=lastSignTime" json:"lastSignTime,omitempty"`
	SignCount          *int32  `protobuf:"varint,10,opt,name=signCount" json:"signCount,omitempty"`
	Diamond            *int64  `protobuf:"varint,11,opt,name=Diamond" json:"Diamond,omitempty"`
	OpenId             *string `protobuf:"bytes,12,opt,name=openId" json:"openId,omitempty"`
	HeadUrl            *string `protobuf:"bytes,13,opt,name=headUrl" json:"headUrl,omitempty"`
	XXX_unrecognized   []byte  `json:"-"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *User) GetMid() string {
	if m != nil && m.Mid != nil {
		return *m.Mid
	}
	return ""
}

func (m *User) GetPwd() string {
	if m != nil && m.Pwd != nil {
		return *m.Pwd
	}
	return ""
}

func (m *User) GetCoin() int64 {
	if m != nil && m.Coin != nil {
		return *m.Coin
	}
	return 0
}

func (m *User) GetId() uint32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *User) GetNickName() string {
	if m != nil && m.NickName != nil {
		return *m.NickName
	}
	return ""
}

func (m *User) GetLoginTurntable() bool {
	if m != nil && m.LoginTurntable != nil {
		return *m.LoginTurntable
	}
	return false
}

func (m *User) GetLoginTurntableTime() string {
	if m != nil && m.LoginTurntableTime != nil {
		return *m.LoginTurntableTime
	}
	return ""
}

func (m *User) GetScores() int32 {
	if m != nil && m.Scores != nil {
		return *m.Scores
	}
	return 0
}

func (m *User) GetLastSignTime() string {
	if m != nil && m.LastSignTime != nil {
		return *m.LastSignTime
	}
	return ""
}

func (m *User) GetSignCount() int32 {
	if m != nil && m.SignCount != nil {
		return *m.SignCount
	}
	return 0
}

func (m *User) GetDiamond() int64 {
	if m != nil && m.Diamond != nil {
		return *m.Diamond
	}
	return 0
}

func (m *User) GetOpenId() string {
	if m != nil && m.OpenId != nil {
		return *m.OpenId
	}
	return ""
}

func (m *User) GetHeadUrl() string {
	if m != nil && m.HeadUrl != nil {
		return *m.HeadUrl
	}
	return ""
}

func init() {
	proto.RegisterType((*ProtoHeader)(nil), "bbproto.ProtoHeader")
	proto.RegisterType((*TerminalInfo)(nil), "bbproto.TerminalInfo")
	proto.RegisterType((*User)(nil), "bbproto.User")
	proto.RegisterEnum("bbproto.EProtoId", EProtoId_name, EProtoId_value)
	proto.RegisterEnum("bbproto.DDErrorCode", DDErrorCode_name, DDErrorCode_value)
}

var fileDescriptor2 = []byte{
	// 1249 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x56, 0x59, 0x77, 0xdb, 0x36,
	0x13, 0xfd, 0xbc, 0x5b, 0xf0, 0x12, 0x04, 0xd9, 0x94, 0xd8, 0x49, 0x1c, 0x27, 0x5f, 0xb3, 0xb4,
	0x49, 0x93, 0x36, 0xdd, 0x57, 0x8a, 0x84, 0x25, 0xd8, 0x24, 0xa0, 0x90, 0x60, 0x7c, 0xdc, 0x17,
	0x1f, 0x59, 0x62, 0x13, 0x9d, 0xd8, 0xa2, 0x0f, 0x25, 0xb7, 0x7d, 0xed, 0x4f, 0x68, 0x5f, 0xba,
	0xfd, 0xc4, 0xbe, 0xf4, 0x1f, 0xa4, 0x33, 0xa0, 0x44, 0x1a, 0xa6, 0x4f, 0xfd, 0x20, 0x03, 0xf7,
	0xde, 0x01, 0xee, 0xcc, 0x00, 0x38, 0x24, 0xe4, 0xa0, 0x33, 0x4c, 0x9e, 0x1c, 0x67, 0xe9, 0x28,
	0x65, 0x0b, 0x07, 0x07, 0x66, 0xb0, 0xf9, 0xf3, 0x14, 0x59, 0x6a, 0xe3, 0xa8, 0x95, 0x74, 0x7a,
	0x49, 0xc6, 0x56, 0xc9, 0x7c, 0xe7, 0xb8, 0xff, 0x32, 0xc9, 0xea, 0x53, 0x1b, 0x53, 0x0f, 0x6a,
	0xec, 0x22, 0xa9, 0x0d, 0x93, 0xe1, 0xb0, 0x9f, 0x0e, 0x44, 0xaf, 0x3e, 0x6d, 0x20, 0x90, 0x9c,
	0x0c, 0x93, 0x0c, 0xe6, 0x33, 0x30, 0x5f, 0x61, 0x94, 0x2c, 0x1e, 0x77, 0xba, 0x6f, 0x92, 0x11,
	0x20, 0xb3, 0x80, 0xcc, 0xb1, 0x65, 0x32, 0xdb, 0x4d, 0x7b, 0x49, 0x7d, 0xce, 0xcc, 0x56, 0xc8,
	0x5c, 0x92, 0x65, 0x69, 0x56, 0x9f, 0x37, 0xe1, 0x20, 0x4f, 0x7e, 0x1a, 0x65, 0x1d, 0xdd, 0x79,
	0x55, 0x5f, 0x40, 0xc1, 0x66, 0x40, 0x96, 0x75, 0x92, 0x1d, 0xf5, 0x07, 0x9d, 0x43, 0x31, 0xf8,
	0x3e, 0x65, 0x17, 0xc8, 0x42, 0xf7, 0x75, 0x67, 0x30, 0x48, 0x0e, 0xc7, 0x26, 0x18, 0x21, 0xbd,
	0xe4, 0x87, 0x7e, 0x37, 0x91, 0x9d, 0xa3, 0x64, 0xec, 0x02, 0xf6, 0x38, 0x39, 0xe9, 0xe7, 0x1e,
	0x6a, 0x20, 0x98, 0x4e, 0x87, 0x66, 0xf7, 0xda, 0xe6, 0xdf, 0x53, 0x64, 0x36, 0x06, 0x83, 0x6c,
	0x89, 0xcc, 0x1c, 0x81, 0x22, 0x5f, 0x03, 0x26, 0xc7, 0x3f, 0xf6, 0xca, 0xe0, 0x6e, 0xda, 0x1f,
	0x98, 0xe0, 0x19, 0x0c, 0xee, 0xe7, 0xd6, 0x4d, 0x32, 0x83, 0x7e, 0xf7, 0x8d, 0xd9, 0x68, 0xce,
	0x68, 0xaf, 0x92, 0xd5, 0xc3, 0xf4, 0x55, 0x7f, 0xa0, 0x4f, 0xb2, 0xc1, 0xa8, 0x73, 0x70, 0x98,
	0x98, 0x3c, 0x16, 0xd9, 0x0d, 0xc2, 0x6c, 0x5c, 0xf7, 0x21, 0x66, 0x61, 0x52, 0xa2, 0x61, 0x37,
	0xcd, 0x92, 0x61, 0x7d, 0xd1, 0x94, 0xe0, 0x32, 0x59, 0x3e, 0xec, 0x0c, 0x47, 0x51, 0xff, 0xd5,
	0xc0, 0xa8, 0x6a, 0x45, 0x6d, 0x01, 0x71, 0xd3, 0x93, 0xc1, 0xa8, 0x4e, 0x8c, 0x10, 0x52, 0xf7,
	0xfa, 0x9d, 0xa3, 0x74, 0xd0, 0xab, 0x2f, 0x19, 0x6f, 0xb0, 0x52, 0x7a, 0x9c, 0x60, 0xf1, 0x97,
	0x4d, 0x0c, 0x08, 0x5e, 0x43, 0xa7, 0xe2, 0xec, 0xb0, 0xbe, 0x82, 0xc0, 0xa3, 0xb7, 0xab, 0x64,
	0x91, 0x9b, 0x0e, 0x8a, 0x1e, 0xb0, 0x4b, 0x4d, 0xae, 0x85, 0xd4, 0x2a, 0x54, 0x2a, 0xa0, 0x53,
	0x50, 0xfb, 0x5a, 0xd4, 0x8a, 0x45, 0x33, 0x56, 0xdb, 0x82, 0x4e, 0xb3, 0x4b, 0xe4, 0x42, 0x31,
	0x6d, 0x09, 0x5f, 0x05, 0x6d, 0x3a, 0x03, 0x29, 0x2f, 0x17, 0x60, 0xc8, 0x23, 0x3a, 0x0b, 0xb5,
	0x5a, 0xf8, 0x6e, 0xbb, 0x65, 0x96, 0x98, 0x83, 0xea, 0xcc, 0xc3, 0xa4, 0xc1, 0x35, 0x9d, 0x1f,
	0x8f, 0x83, 0xa8, 0x49, 0x17, 0xa0, 0x2e, 0x0c, 0xc6, 0x2f, 0x62, 0x1e, 0xee, 0x49, 0x15, 0x71,
	0x47, 0xc7, 0x11, 0x0f, 0xe9, 0x22, 0x38, 0x26, 0x18, 0xcc, 0x5f, 0x20, 0x46, 0x6b, 0xe3, 0xb9,
	0xaf, 0xb4, 0x06, 0x25, 0x25, 0xac, 0x4e, 0x2e, 0xe3, 0x7a, 0xa1, 0x72, 0x3c, 0xd7, 0x89, 0x74,
	0x83, 0x37, 0x85, 0xc4, 0xd5, 0x97, 0xa0, 0x4a, 0xd4, 0x57, 0x30, 0x8b, 0x44, 0x53, 0x02, 0xa6,
	0x64, 0x1c, 0xd1, 0x65, 0x76, 0x8d, 0x5c, 0x32, 0xa8, 0x8e, 0x43, 0xa9, 0x9d, 0x86, 0xcf, 0x73,
	0x62, 0x05, 0x17, 0x56, 0xbe, 0x90, 0xe3, 0xf9, 0x2a, 0x26, 0xaf, 0x45, 0x20, 0x64, 0x33, 0x07,
	0x2e, 0xa0, 0x5b, 0x9d, 0x67, 0x41, 0x59, 0x8d, 0xcc, 0x69, 0x93, 0xc4, 0x45, 0x8c, 0x33, 0x43,
	0xb3, 0x33, 0x65, 0x70, 0xba, 0x56, 0xf3, 0xf9, 0xc4, 0x12, 0xbd, 0x04, 0xcd, 0xbd, 0x9a, 0x87,
	0x3a, 0x9e, 0x87, 0x79, 0x95, 0xdc, 0x65, 0xb4, 0xd9, 0x16, 0xde, 0x3e, 0x64, 0xf8, 0x22, 0x16,
	0xee, 0x8e, 0xab, 0xa4, 0xa4, 0x57, 0x26, 0xa8, 0xe3, 0xee, 0x94, 0xe8, 0x55, 0xf4, 0x84, 0xa8,
	0x8c, 0x7d, 0x1f, 0xab, 0x76, 0x8d, 0xdd, 0x23, 0x1b, 0x08, 0x04, 0x8e, 0x76, 0x5b, 0xbe, 0x88,
	0x34, 0x2e, 0x13, 0xa8, 0x86, 0xf0, 0x79, 0x01, 0xd1, 0x3a, 0xd6, 0x16, 0x55, 0x4d, 0x27, 0xe0,
	0xfb, 0x26, 0x79, 0x1c, 0xd1, 0xeb, 0x58, 0x8b, 0x02, 0xe7, 0x12, 0x0a, 0x6a, 0x82, 0xe8, 0x0d,
	0xf4, 0x5b, 0x10, 0x60, 0xe1, 0x14, 0xb7, 0xc6, 0xae, 0x93, 0x2b, 0x05, 0x17, 0x71, 0xe9, 0xe1,
	0x40, 0xc8, 0x2d, 0x45, 0xd7, 0xad, 0x7d, 0x1a, 0x50, 0x4b, 0xcf, 0x55, 0x50, 0x92, 0x9b, 0xec,
	0x0a, 0xb9, 0x58, 0xe0, 0x42, 0x0a, 0xed, 0x3a, 0xa1, 0x47, 0x6f, 0x55, 0x56, 0xda, 0xf2, 0x55,
	0xdb, 0x50, 0xb7, 0x2b, 0x14, 0x76, 0xcb, 0x50, 0x1b, 0x96, 0x37, 0xa4, 0x42, 0xf1, 0x92, 0x87,
	0x86, 0xbb, 0x63, 0x6d, 0x14, 0x3a, 0x22, 0xe2, 0xd8, 0xa2, 0x4d, 0x3c, 0x23, 0xa7, 0xd3, 0x29,
	0x98, 0xbb, 0x96, 0xe3, 0x2d, 0xe5, 0xfb, 0x6a, 0x17, 0xf1, 0x7b, 0xd6, 0xfe, 0x10, 0x51, 0x52,
	0xff, 0x9f, 0x74, 0x66, 0x12, 0xe2, 0x21, 0xfa, 0x8e, 0x55, 0xca, 0x3c, 0xc0, 0x10, 0xf7, 0x2d,
	0x4b, 0x6e, 0x8b, 0xbb, 0x3b, 0x08, 0x3f, 0x38, 0x6b, 0xa9, 0x60, 0x1e, 0x56, 0x52, 0x57, 0x90,
	0x1e, 0xa6, 0x4f, 0x1f, 0x59, 0x41, 0x48, 0x8d, 0x4f, 0x13, 0x7d, 0xd7, 0x2a, 0x4a, 0x4e, 0xb7,
	0xd4, 0xae, 0x29, 0xca, 0x7b, 0xec, 0x26, 0xb9, 0x6e, 0x73, 0xb0, 0x5f, 0x41, 0x3f, 0x66, 0xeb,
	0xa4, 0x6e, 0xd3, 0x9a, 0x47, 0x1a, 0x2e, 0x6e, 0xec, 0x6b, 0xfa, 0xc4, 0xca, 0xb6, 0x1d, 0x72,
	0xd3, 0xd0, 0xf7, 0xad, 0x6c, 0xf1, 0x47, 0x2a, 0x2d, 0x5c, 0x4e, 0x9f, 0xb2, 0x35, 0x72, 0xad,
	0xb2, 0xd7, 0x98, 0x7c, 0x56, 0xdd, 0xc9, 0x0d, 0xe1, 0x4e, 0x73, 0x8f, 0x47, 0x3b, 0xf4, 0x03,
	0x76, 0x9b, 0xac, 0x55, 0x42, 0x4f, 0x09, 0x3e, 0xac, 0x6c, 0x0a, 0x67, 0xdd, 0xf1, 0xf6, 0xe8,
	0xf3, 0x6a, 0xf2, 0xd8, 0x63, 0xc3, 0x7d, 0x54, 0x0d, 0xca, 0xaf, 0xe9, 0xc7, 0x55, 0x33, 0xf8,
	0x83, 0x09, 0x42, 0x51, 0x3e, 0x39, 0xd7, 0xcc, 0x29, 0xc1, 0xa7, 0x6c, 0x83, 0xac, 0x9f, 0x5d,
	0xd7, 0x91, 0xa7, 0x14, 0x9f, 0xb1, 0x5b, 0xe4, 0x86, 0xad, 0xf0, 0x44, 0x14, 0x29, 0xff, 0x65,
	0x9e, 0xce, 0xe7, 0xec, 0x0e, 0xb9, 0x59, 0xd9, 0xc2, 0x92, 0x7c, 0x61, 0xdf, 0x5b, 0xee, 0x8c,
	0xf1, 0x2f, 0xcf, 0x9e, 0xce, 0x92, 0xfa, 0xca, 0x32, 0x8e, 0x47, 0x04, 0x61, 0xf8, 0x37, 0x79,
	0x2f, 0xbf, 0xb6, 0x1a, 0x1a, 0xf0, 0x28, 0x72, 0x9a, 0x9c, 0x7e, 0x53, 0x39, 0x59, 0x13, 0xe6,
	0x5b, 0xab, 0xa3, 0x5a, 0xc5, 0x12, 0xbb, 0xad, 0xcd, 0xe5, 0xa6, 0x8e, 0x75, 0xb4, 0x0a, 0x32,
	0xe4, 0xbb, 0x70, 0xb2, 0x22, 0xda, 0xb0, 0x1a, 0x53, 0xd2, 0x8e, 0xdc, 0xa1, 0xee, 0xf9, 0xa1,
	0x51, 0x1c, 0x04, 0x0e, 0xd8, 0xf4, 0xac, 0xd4, 0xcb, 0xa7, 0x8c, 0x5b, 0x75, 0x2f, 0xc2, 0xda,
	0xbe, 0xb3, 0xc7, 0x43, 0xb3, 0xf0, 0x16, 0xbe, 0xbf, 0xe5, 0x1b, 0xc0, 0x1b, 0xf1, 0x1e, 0x6d,
	0x5a, 0x97, 0xd0, 0x1c, 0x0e, 0x84, 0x5b, 0x96, 0xd4, 0xbc, 0x8b, 0x54, 0x9c, 0x95, 0xe6, 0xf0,
	0xb6, 0x05, 0x6f, 0x71, 0xee, 0x35, 0x80, 0xa3, 0x3b, 0x16, 0x0c, 0xe7, 0x3c, 0x5f, 0xd8, 0x3f,
	0x7b, 0xbb, 0x0b, 0x26, 0xb0, 0xfc, 0x63, 0xa1, 0xdd, 0x16, 0x1c, 0x1b, 0xd3, 0x3c, 0xb5, 0x2b,
	0xe1, 0x2a, 0xcb, 0x47, 0xff, 0x4c, 0x93, 0x25, 0x68, 0x1b, 0x7e, 0xe2, 0xb8, 0xf0, 0xd5, 0x83,
	0x26, 0x21, 0x37, 0x15, 0xba, 0xca, 0x83, 0x90, 0xd8, 0x75, 0xe9, 0xff, 0xd8, 0x73, 0x72, 0xbf,
	0xc4, 0xf2, 0x4b, 0xb2, 0x8f, 0x6b, 0xc0, 0x19, 0x73, 0x02, 0x25, 0x3d, 0xf4, 0x02, 0x97, 0x32,
	0x6e, 0xb6, 0xe8, 0x5f, 0x6f, 0xc7, 0x7f, 0x53, 0xec, 0x29, 0xb9, 0x7b, 0x7e, 0x14, 0x3e, 0x22,
	0x18, 0xb2, 0x05, 0x05, 0xf5, 0xe8, 0x9f, 0x65, 0xc4, 0x43, 0xb2, 0x5e, 0x46, 0xe0, 0x77, 0x40,
	0xae, 0x2f, 0xa4, 0x7f, 0x94, 0xd2, 0xc7, 0x64, 0xa3, 0x94, 0x42, 0x67, 0xc2, 0xbc, 0x35, 0xfb,
	0xae, 0x23, 0x41, 0xbe, 0xbf, 0x8d, 0x0f, 0xc7, 0xef, 0xa5, 0xfc, 0x09, 0xb9, 0xf3, 0x5f, 0xf2,
	0xbc, 0x5c, 0xbf, 0x95, 0xfa, 0x07, 0x64, 0xad, 0xd4, 0x97, 0xb7, 0x1f, 0x7e, 0xdb, 0xf8, 0x5d,
	0xf0, 0x6b, 0xa9, 0x7c, 0x46, 0xee, 0x9d, 0xab, 0x74, 0x5b, 0xa2, 0x8d, 0xce, 0xf7, 0xc7, 0x85,
	0xf9, 0xa5, 0x08, 0xf9, 0x37, 0x00, 0x00, 0xff, 0xff, 0x8b, 0x04, 0xf8, 0x83, 0xcb, 0x0a, 0x00,
	0x00,
}
