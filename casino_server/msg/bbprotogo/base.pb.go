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
	EProtoId_TESTP1          EProtoId = 0
	EProtoId_REG             EProtoId = 1
	EProtoId_REQAUTHUSER     EProtoId = 2
	EProtoId_HEATBEAT        EProtoId = 3
	EProtoId_GETINTOROOM     EProtoId = 4
	EProtoId_ROOMMSG         EProtoId = 5
	EProtoId_GETREWARDS      EProtoId = 6
	EProtoId_SHUIGUOJI       EProtoId = 7
	EProtoId_SHUIGUOJIHILOMP EProtoId = 8
	EProtoId_SHUIGUOJIRES    EProtoId = 9
	// 扎金花相关的逻辑
	EProtoId_ZJHROOM                          EProtoId = 10
	EProtoId_ZJHBET                           EProtoId = 11
	EProtoId_ZJHMSG                           EProtoId = 12
	EProtoId_ZJHQUERYNOSEATUSER               EProtoId = 13
	EProtoId_ZJHREQSEAT                       EProtoId = 14
	EProtoId_ZJHLOTTERY                       EProtoId = 15
	EProtoId_ZJHBROADCASTBEGINBET             EProtoId = 16
	EProtoId_LOGINSIGNINBONUS                 EProtoId = 17
	EProtoId_LOGINTURNTABLEBONUS              EProtoId = 18
	EProtoId_OLINEBONUS                       EProtoId = 19
	EProtoId_TIMINGBONUS                      EProtoId = 20
	EProtoId_THROOM                           EProtoId = 21
	EProtoId_THBET                            EProtoId = 22
	EProtoId_THBETBEGIN                       EProtoId = 23
	EProtoId_THBETBROADCAST                   EProtoId = 24
	EProtoId_THROOMADDUSERBROADCAST           EProtoId = 25
	EProtoId_REQQUICKCONN                     EProtoId = 26
	EProtoId_PID_ACKQUICKCONN                 EProtoId = 27
	EProtoId_PID_NULLMSG                      EProtoId = 28
	EProtoId_PID_MATCHLIST_REQMOBILEMATCHLIST EProtoId = 29
	EProtoId_PID_GAME_LOGINGAME               EProtoId = 30
	EProtoId_PID_GAME_ENTERMATCH              EProtoId = 31
	EProtoId_PID_GAME_ACKENTERMATCH           EProtoId = 32
	EProtoId_PID_GAME_SENDGAMEINFO            EProtoId = 33
	EProtoId_PID_GAME_BLINDCOIN               EProtoId = 34
	EProtoId_PID_GAME_INITCARD                EProtoId = 35
	EProtoId_PID_GAME_SENDFLOPCARD            EProtoId = 36
	EProtoId_PID_GAME_SENDTURNCARD            EProtoId = 37
	EProtoId_PID_GAME_SENDRIVERCARD           EProtoId = 38
	EProtoId_PID_GAME_RAISEBET                EProtoId = 39
	EProtoId_PID_GAME_ACKRAISEBET             EProtoId = 40
	EProtoId_PID_GAME_FOLLOWBET               EProtoId = 41
	EProtoId_PID_GAME_ACKFOLLOWBET            EProtoId = 42
	EProtoId_PID_GAME_FOLDBET                 EProtoId = 43
	EProtoId_PID_GAME_ACKFOLDBET              EProtoId = 44
	EProtoId_PID_GAME_CHECKBET                EProtoId = 45
	EProtoId_PID_GAME_ACKCHECKBET             EProtoId = 46
	EProtoId_PID_GAME_SENDOVERTURN            EProtoId = 47
	EProtoId_PID_GAME_SENDADDUSER             EProtoId = 48
	EProtoId_PID_GAME_GAME_SHOWCARD           EProtoId = 49
	EProtoId_PID_GAME_GAME_ACKSHOWCARD        EProtoId = 50
	EProtoId_PID_GAME_GAME_TESTRESULT         EProtoId = 51
)

var EProtoId_name = map[int32]string{
	0:  "TESTP1",
	1:  "REG",
	2:  "REQAUTHUSER",
	3:  "HEATBEAT",
	4:  "GETINTOROOM",
	5:  "ROOMMSG",
	6:  "GETREWARDS",
	7:  "SHUIGUOJI",
	8:  "SHUIGUOJIHILOMP",
	9:  "SHUIGUOJIRES",
	10: "ZJHROOM",
	11: "ZJHBET",
	12: "ZJHMSG",
	13: "ZJHQUERYNOSEATUSER",
	14: "ZJHREQSEAT",
	15: "ZJHLOTTERY",
	16: "ZJHBROADCASTBEGINBET",
	17: "LOGINSIGNINBONUS",
	18: "LOGINTURNTABLEBONUS",
	19: "OLINEBONUS",
	20: "TIMINGBONUS",
	21: "THROOM",
	22: "THBET",
	23: "THBETBEGIN",
	24: "THBETBROADCAST",
	25: "THROOMADDUSERBROADCAST",
	26: "REQQUICKCONN",
	27: "PID_ACKQUICKCONN",
	28: "PID_NULLMSG",
	29: "PID_MATCHLIST_REQMOBILEMATCHLIST",
	30: "PID_GAME_LOGINGAME",
	31: "PID_GAME_ENTERMATCH",
	32: "PID_GAME_ACKENTERMATCH",
	33: "PID_GAME_SENDGAMEINFO",
	34: "PID_GAME_BLINDCOIN",
	35: "PID_GAME_INITCARD",
	36: "PID_GAME_SENDFLOPCARD",
	37: "PID_GAME_SENDTURNCARD",
	38: "PID_GAME_SENDRIVERCARD",
	39: "PID_GAME_RAISEBET",
	40: "PID_GAME_ACKRAISEBET",
	41: "PID_GAME_FOLLOWBET",
	42: "PID_GAME_ACKFOLLOWBET",
	43: "PID_GAME_FOLDBET",
	44: "PID_GAME_ACKFOLDBET",
	45: "PID_GAME_CHECKBET",
	46: "PID_GAME_ACKCHECKBET",
	47: "PID_GAME_SENDOVERTURN",
	48: "PID_GAME_SENDADDUSER",
	49: "PID_GAME_GAME_SHOWCARD",
	50: "PID_GAME_GAME_ACKSHOWCARD",
	51: "PID_GAME_GAME_TESTRESULT",
}
var EProtoId_value = map[string]int32{
	"TESTP1":                           0,
	"REG":                              1,
	"REQAUTHUSER":                      2,
	"HEATBEAT":                         3,
	"GETINTOROOM":                      4,
	"ROOMMSG":                          5,
	"GETREWARDS":                       6,
	"SHUIGUOJI":                        7,
	"SHUIGUOJIHILOMP":                  8,
	"SHUIGUOJIRES":                     9,
	"ZJHROOM":                          10,
	"ZJHBET":                           11,
	"ZJHMSG":                           12,
	"ZJHQUERYNOSEATUSER":               13,
	"ZJHREQSEAT":                       14,
	"ZJHLOTTERY":                       15,
	"ZJHBROADCASTBEGINBET":             16,
	"LOGINSIGNINBONUS":                 17,
	"LOGINTURNTABLEBONUS":              18,
	"OLINEBONUS":                       19,
	"TIMINGBONUS":                      20,
	"THROOM":                           21,
	"THBET":                            22,
	"THBETBEGIN":                       23,
	"THBETBROADCAST":                   24,
	"THROOMADDUSERBROADCAST":           25,
	"REQQUICKCONN":                     26,
	"PID_ACKQUICKCONN":                 27,
	"PID_NULLMSG":                      28,
	"PID_MATCHLIST_REQMOBILEMATCHLIST": 29,
	"PID_GAME_LOGINGAME":               30,
	"PID_GAME_ENTERMATCH":              31,
	"PID_GAME_ACKENTERMATCH":           32,
	"PID_GAME_SENDGAMEINFO":            33,
	"PID_GAME_BLINDCOIN":               34,
	"PID_GAME_INITCARD":                35,
	"PID_GAME_SENDFLOPCARD":            36,
	"PID_GAME_SENDTURNCARD":            37,
	"PID_GAME_SENDRIVERCARD":           38,
	"PID_GAME_RAISEBET":                39,
	"PID_GAME_ACKRAISEBET":             40,
	"PID_GAME_FOLLOWBET":               41,
	"PID_GAME_ACKFOLLOWBET":            42,
	"PID_GAME_FOLDBET":                 43,
	"PID_GAME_ACKFOLDBET":              44,
	"PID_GAME_CHECKBET":                45,
	"PID_GAME_ACKCHECKBET":             46,
	"PID_GAME_SENDOVERTURN":            47,
	"PID_GAME_SENDADDUSER":             48,
	"PID_GAME_GAME_SHOWCARD":           49,
	"PID_GAME_GAME_ACKSHOWCARD":        50,
	"PID_GAME_GAME_TESTRESULT":         51,
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
func (EProtoId) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

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
func (*ProtoHeader) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

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
func (*TerminalInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

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

func init() {
	proto.RegisterType((*ProtoHeader)(nil), "bbproto.ProtoHeader")
	proto.RegisterType((*TerminalInfo)(nil), "bbproto.TerminalInfo")
	proto.RegisterEnum("bbproto.EProtoId", EProtoId_name, EProtoId_value)
}

var fileDescriptor3 = []byte{
	// 738 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x54, 0xdb, 0x52, 0xe3, 0x46,
	0x10, 0x8d, 0x01, 0xdf, 0xc6, 0x06, 0x86, 0xe1, 0x26, 0x08, 0x24, 0x84, 0x90, 0x84, 0x90, 0x84,
	0x84, 0xe4, 0x0b, 0x64, 0x69, 0x90, 0x06, 0xa4, 0x19, 0x5b, 0x1a, 0x41, 0x25, 0x2f, 0x94, 0xb0,
	0x95, 0xc4, 0xb5, 0x60, 0x53, 0x32, 0x6c, 0xed, 0xeb, 0xfe, 0xe6, 0x7e, 0xcd, 0x76, 0x8f, 0x8c,
	0x40, 0xcb, 0x03, 0x55, 0xdd, 0xe7, 0xf4, 0xe5, 0x9c, 0x1e, 0x61, 0x42, 0x6e, 0xd3, 0x59, 0x76,
	0xfa, 0x90, 0x4f, 0x1f, 0xa7, 0xac, 0x79, 0x7b, 0x6b, 0x82, 0xc3, 0x8f, 0x35, 0xd2, 0xe9, 0x63,
	0xe4, 0x67, 0xe9, 0x28, 0xcb, 0xd9, 0x0a, 0x69, 0xa4, 0x0f, 0xe3, 0xab, 0x2c, 0xb7, 0x6a, 0x07,
	0xb5, 0xe3, 0x36, 0x5b, 0x23, 0xed, 0x59, 0x36, 0x9b, 0x8d, 0xa7, 0x13, 0x31, 0xb2, 0x16, 0x0c,
	0x04, 0x25, 0x4f, 0xb3, 0x2c, 0x87, 0x7c, 0x11, 0xf2, 0x65, 0x46, 0x49, 0xeb, 0x21, 0x1d, 0xbe,
	0xcb, 0x1e, 0x01, 0x59, 0x02, 0xa4, 0xce, 0xba, 0x64, 0x69, 0x38, 0x1d, 0x65, 0x56, 0xdd, 0x64,
	0xcb, 0xa4, 0x9e, 0xe5, 0xf9, 0x34, 0xb7, 0x1a, 0xa6, 0x1d, 0xca, 0xb3, 0x0f, 0x8f, 0x79, 0xaa,
	0xd3, 0xff, 0xac, 0x26, 0x16, 0x1c, 0x86, 0xa4, 0xab, 0xb3, 0xfc, 0x7e, 0x3c, 0x49, 0xef, 0xc4,
	0xe4, 0xdf, 0x29, 0x5b, 0x25, 0xcd, 0xe1, 0xff, 0xe9, 0x64, 0x92, 0xdd, 0xcd, 0x45, 0x30, 0x42,
	0x46, 0xd9, 0xfb, 0xf1, 0x30, 0x93, 0xe9, 0x7d, 0x36, 0x57, 0x01, 0x3b, 0x9e, 0x9e, 0xc6, 0x85,
	0x86, 0x36, 0x14, 0x2c, 0x4c, 0x67, 0x66, 0x7b, 0xfb, 0xe4, 0x53, 0x8b, 0xb4, 0xb8, 0xf1, 0x24,
	0x46, 0x40, 0x34, 0x34, 0x8f, 0x75, 0xff, 0x8c, 0x7e, 0xc5, 0x9a, 0x64, 0x31, 0xe2, 0x1e, 0xad,
	0xc1, 0x82, 0x4e, 0xc4, 0x07, 0x76, 0xa2, 0xfd, 0x24, 0xe6, 0x11, 0x5d, 0x80, 0x61, 0x2d, 0x9f,
	0xdb, 0xba, 0x07, 0x7f, 0x74, 0x11, 0x69, 0x8f, 0x6b, 0x21, 0xb5, 0x8a, 0x94, 0x0a, 0xe9, 0x12,
	0xeb, 0x90, 0x26, 0x46, 0x61, 0xec, 0xd1, 0x3a, 0xd8, 0x27, 0xc0, 0x46, 0xfc, 0xda, 0x8e, 0xdc,
	0x98, 0x36, 0xc0, 0x5e, 0x3b, 0xf6, 0x13, 0xe1, 0x25, 0xea, 0x42, 0xd0, 0x26, 0x5b, 0x27, 0xab,
	0x65, 0xea, 0x8b, 0x40, 0x85, 0x7d, 0xda, 0x02, 0xcf, 0xdd, 0x12, 0x8c, 0x78, 0x4c, 0xdb, 0x38,
	0xf2, 0x9f, 0x0b, 0xdf, 0xcc, 0x27, 0x28, 0x12, 0x92, 0x1e, 0xd7, 0xb4, 0x33, 0x8f, 0x71, 0x55,
	0x97, 0x6d, 0x11, 0x06, 0xf1, 0x20, 0xe1, 0xd1, 0xdf, 0x52, 0xc5, 0x20, 0xce, 0xc8, 0x5d, 0x46,
	0x09, 0xd8, 0xcc, 0x07, 0x88, 0xd1, 0x95, 0x79, 0x1e, 0x28, 0xad, 0xa1, 0x92, 0xae, 0x32, 0x8b,
	0x6c, 0xe0, 0xbc, 0x48, 0xd9, 0xae, 0x63, 0xc7, 0x60, 0xcb, 0x13, 0x12, 0xa7, 0x53, 0xb6, 0x41,
	0x68, 0xa0, 0x20, 0x8b, 0x85, 0x27, 0x01, 0x53, 0x32, 0x89, 0xe9, 0x1a, 0xdb, 0x26, 0xeb, 0x06,
	0xd5, 0x49, 0x24, 0xb5, 0xdd, 0x0b, 0x78, 0x41, 0x30, 0x1c, 0xac, 0x02, 0x21, 0xe7, 0xf9, 0x3a,
	0x5e, 0x46, 0x8b, 0x50, 0x48, 0xaf, 0x00, 0x36, 0xcc, 0x79, 0x0b, 0x17, 0x9b, 0xac, 0x4d, 0xea,
	0xda, 0x98, 0xd8, 0xc2, 0x3e, 0x13, 0x9a, 0xcd, 0x74, 0x1b, 0x1e, 0x70, 0xa5, 0xc8, 0x9f, 0x25,
	0x51, 0x8b, 0xed, 0x92, 0xad, 0xa2, 0xd5, 0x76, 0x5d, 0xf4, 0xf5, 0xc2, 0xed, 0xe0, 0xbd, 0xc0,
	0xdd, 0x20, 0x11, 0xce, 0xa5, 0xa3, 0xa4, 0xa4, 0xbb, 0x28, 0xbc, 0x2f, 0xdc, 0x1b, 0xdb, 0xb9,
	0x7c, 0x41, 0xbf, 0x46, 0x3d, 0x88, 0xca, 0x24, 0x08, 0xf0, 0x62, 0x7b, 0xec, 0x88, 0x1c, 0x20,
	0x10, 0xda, 0xda, 0xf1, 0x03, 0x11, 0xeb, 0x1b, 0x18, 0x13, 0xaa, 0x9e, 0x08, 0x78, 0x09, 0xd1,
	0x7d, 0xbc, 0x2b, 0x56, 0x79, 0x76, 0xc8, 0x6f, 0x8c, 0x71, 0x8c, 0xe8, 0x37, 0x78, 0x87, 0x12,
	0xe7, 0x12, 0x8e, 0x69, 0x9a, 0xe8, 0xb7, 0xa8, 0xb5, 0x24, 0x40, 0xc2, 0x2b, 0xee, 0x80, 0xed,
	0x90, 0xcd, 0x92, 0x8b, 0xb9, 0x74, 0x31, 0x10, 0xf2, 0x5c, 0xd1, 0xef, 0x2a, 0x7b, 0x7a, 0x70,
	0x47, 0xd7, 0x51, 0x70, 0x8e, 0x43, 0xb6, 0x49, 0xd6, 0x4a, 0x5c, 0x48, 0xa1, 0x1d, 0xf8, 0x94,
	0xe8, 0xf7, 0x6f, 0x26, 0x9d, 0x07, 0xaa, 0x6f, 0xa8, 0xa3, 0x37, 0x14, 0xbe, 0x94, 0xa1, 0x7e,
	0xa8, 0x68, 0x43, 0x2a, 0x12, 0x57, 0x3c, 0x32, 0xdc, 0x8f, 0x95, 0x45, 0x91, 0x2d, 0x62, 0x8e,
	0xcf, 0xf3, 0x13, 0x7e, 0x1f, 0xaf, 0xed, 0x94, 0xcc, 0x71, 0x45, 0xf1, 0xb9, 0x0a, 0x02, 0x75,
	0x8d, 0xf8, 0xcf, 0x95, 0xfd, 0xd0, 0xf1, 0x42, 0x9d, 0x3c, 0xbf, 0xcc, 0x73, 0x8b, 0x8b, 0xe8,
	0x2f, 0x95, 0x53, 0x16, 0x0d, 0x86, 0xf8, 0xb5, 0x22, 0xc9, 0xf1, 0xb9, 0x73, 0x89, 0xf0, 0x6f,
	0x5f, 0x4a, 0x2a, 0x99, 0xd3, 0x37, 0xd6, 0x15, 0xd8, 0x43, 0xfb, 0xf4, 0xf7, 0x4a, 0x13, 0x52,
	0xf3, 0x2f, 0x89, 0xfe, 0x51, 0x39, 0x4a, 0x41, 0xfb, 0xea, 0xda, 0x1c, 0xe5, 0x8c, 0xed, 0x93,
	0x9d, 0x2a, 0x07, 0xfb, 0x4a, 0xfa, 0x4f, 0xb6, 0x47, 0xac, 0x2a, 0x8d, 0xbf, 0x1f, 0xf0, 0x4f,
	0x9b, 0x04, 0x9a, 0xfe, 0xf5, 0x39, 0x00, 0x00, 0xff, 0xff, 0x8e, 0x22, 0x80, 0x04, 0x45, 0x05,
	0x00, 0x00,
}
