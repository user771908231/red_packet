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
	EProtoId_ZJHROOM                EProtoId = 10
	EProtoId_ZJHBET                 EProtoId = 11
	EProtoId_ZJHMSG                 EProtoId = 12
	EProtoId_ZJHQUERYNOSEATUSER     EProtoId = 13
	EProtoId_ZJHREQSEAT             EProtoId = 14
	EProtoId_ZJHLOTTERY             EProtoId = 15
	EProtoId_ZJHBROADCASTBEGINBET   EProtoId = 16
	EProtoId_LOGINSIGNINBONUS       EProtoId = 17
	EProtoId_LOGINTURNTABLEBONUS    EProtoId = 18
	EProtoId_OLINEBONUS             EProtoId = 19
	EProtoId_TIMINGBONUS            EProtoId = 20
	EProtoId_THROOM                 EProtoId = 21
	EProtoId_THBET                  EProtoId = 22
	EProtoId_THBETBEGIN             EProtoId = 23
	EProtoId_THBETBROADCAST         EProtoId = 24
	EProtoId_THROOMADDUSERBROADCAST EProtoId = 25
	EProtoId_THLOTTERY              EProtoId = 26
	EProtoId_GAME_FOLLOWBET         EProtoId = 41
	EProtoId_GAME_ACKFOLLOWBET      EProtoId = 42
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
	26: "THLOTTERY",
	41: "GAME_FOLLOWBET",
	42: "GAME_ACKFOLLOWBET",
}
var EProtoId_value = map[string]int32{
	"TESTP1":                 0,
	"REG":                    1,
	"REQAUTHUSER":            2,
	"HEATBEAT":               3,
	"GETINTOROOM":            4,
	"ROOMMSG":                5,
	"GETREWARDS":             6,
	"SHUIGUOJI":              7,
	"SHUIGUOJIHILOMP":        8,
	"SHUIGUOJIRES":           9,
	"ZJHROOM":                10,
	"ZJHBET":                 11,
	"ZJHMSG":                 12,
	"ZJHQUERYNOSEATUSER":     13,
	"ZJHREQSEAT":             14,
	"ZJHLOTTERY":             15,
	"ZJHBROADCASTBEGINBET":   16,
	"LOGINSIGNINBONUS":       17,
	"LOGINTURNTABLEBONUS":    18,
	"OLINEBONUS":             19,
	"TIMINGBONUS":            20,
	"THROOM":                 21,
	"THBET":                  22,
	"THBETBEGIN":             23,
	"THBETBROADCAST":         24,
	"THROOMADDUSERBROADCAST": 25,
	"THLOTTERY":              26,
	"GAME_FOLLOWBET":         41,
	"GAME_ACKFOLLOWBET":      42,
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
	// 505 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x44, 0x51, 0x4b, 0x73, 0xd3, 0x3c,
	0x14, 0xfd, 0xf2, 0x74, 0x7c, 0xf3, 0x52, 0x94, 0x34, 0xf5, 0xd7, 0x15, 0xd3, 0x15, 0x74, 0xc1,
	0x0c, 0x3f, 0xc1, 0x69, 0x84, 0xad, 0x62, 0x5b, 0x89, 0x2d, 0xd3, 0x81, 0x0d, 0xe3, 0x24, 0x02,
	0x32, 0xb4, 0x71, 0xc6, 0x4e, 0x18, 0xb6, 0xfc, 0x6d, 0x56, 0x5c, 0xc9, 0x26, 0x5d, 0x78, 0xe6,
	0x9e, 0x73, 0x1f, 0xe7, 0x1c, 0x19, 0x60, 0x93, 0x95, 0xea, 0xed, 0xb1, 0xc8, 0x4f, 0x39, 0xb5,
	0x36, 0x1b, 0x53, 0xdc, 0xfe, 0x6e, 0x40, 0x7f, 0xa5, 0x2b, 0x5f, 0x65, 0x3b, 0x55, 0xd0, 0x11,
	0x74, 0xb3, 0xe3, 0xfe, 0xa3, 0x2a, 0x9c, 0xc6, 0xab, 0xc6, 0x6b, 0x9b, 0x4e, 0xc0, 0x2e, 0x55,
	0x59, 0xee, 0xf3, 0x03, 0xdf, 0x39, 0x4d, 0x43, 0xe1, 0xc8, 0xb9, 0x54, 0x05, 0xe2, 0x16, 0xe2,
	0x21, 0x25, 0xd0, 0x3b, 0x66, 0xdb, 0x1f, 0xea, 0x84, 0x4c, 0x1b, 0x99, 0x0e, 0x1d, 0x40, 0x7b,
	0x9b, 0xef, 0x94, 0xd3, 0x31, 0x68, 0x08, 0x1d, 0x55, 0x14, 0x79, 0xe1, 0x74, 0xcd, 0x3a, 0x8e,
	0xab, 0x5f, 0xa7, 0x22, 0x93, 0xd9, 0x37, 0xc7, 0xd2, 0x03, 0xb7, 0x21, 0x0c, 0xa4, 0x2a, 0x9e,
	0xf7, 0x87, 0xec, 0x89, 0x1f, 0xbe, 0xe6, 0x74, 0x0c, 0xd6, 0xf6, 0x7b, 0x76, 0x38, 0xa8, 0xa7,
	0xda, 0x04, 0x05, 0xd8, 0xa9, 0x9f, 0xfb, 0xad, 0x8a, 0xb2, 0x67, 0x55, 0xbb, 0x40, 0x8d, 0xf3,
	0x79, 0x5f, 0x79, 0xb0, 0x71, 0xa0, 0x99, 0x97, 0x46, 0xdd, 0xbe, 0xfb, 0xd3, 0x82, 0x1e, 0x33,
	0x99, 0xf8, 0x0e, 0x1b, 0x5d, 0xc9, 0x12, 0xb9, 0x7a, 0x47, 0xfe, 0xa3, 0x16, 0xb4, 0x62, 0xe6,
	0x91, 0x06, 0x0a, 0xf4, 0x63, 0xb6, 0x76, 0x53, 0xe9, 0xa7, 0x09, 0x8b, 0x49, 0x13, 0x8f, 0xf5,
	0x7c, 0xe6, 0xca, 0x05, 0x7e, 0xa4, 0xa5, 0xdb, 0x1e, 0x93, 0x3c, 0x92, 0x22, 0x16, 0x22, 0x24,
	0x6d, 0xda, 0x07, 0x4b, 0x57, 0x61, 0xe2, 0x91, 0x0e, 0xc6, 0x07, 0xec, 0xc6, 0xec, 0xd1, 0x8d,
	0x97, 0x09, 0xe9, 0x62, 0x3c, 0x3b, 0xf1, 0x53, 0xee, 0xa5, 0xe2, 0x81, 0x13, 0x8b, 0x4e, 0x61,
	0x7c, 0x81, 0x3e, 0x0f, 0x44, 0xb8, 0x22, 0x3d, 0xcc, 0x3c, 0xb8, 0x90, 0x31, 0x4b, 0x88, 0xad,
	0x4f, 0x7e, 0x7e, 0xf0, 0xcd, 0x7d, 0xd0, 0x26, 0x11, 0x2c, 0x98, 0x24, 0xfd, 0xba, 0xd6, 0x52,
	0x03, 0x3a, 0x07, 0x8a, 0xf5, 0x3a, 0x65, 0xf1, 0xa7, 0x48, 0x24, 0x68, 0xce, 0xd8, 0x1d, 0x6a,
	0x0b, 0x7a, 0x99, 0xad, 0x35, 0x47, 0x46, 0x35, 0x0e, 0x84, 0x94, 0x38, 0x49, 0xc6, 0xd4, 0x81,
	0x99, 0xbe, 0x17, 0x0b, 0x77, 0x79, 0xef, 0x26, 0x18, 0xcb, 0xe3, 0x91, 0xbe, 0x4e, 0xe8, 0x0c,
	0x48, 0x20, 0x10, 0x25, 0xdc, 0x8b, 0x90, 0x13, 0x51, 0x9a, 0x90, 0x09, 0xbd, 0x86, 0xa9, 0x61,
	0x65, 0x1a, 0x47, 0xd2, 0x5d, 0x04, 0xac, 0x6a, 0x50, 0x7d, 0x58, 0x04, 0x3c, 0xaa, 0xf1, 0x54,
	0xbf, 0x8c, 0xe4, 0x21, 0x8f, 0xbc, 0x8a, 0x98, 0x99, 0xe7, 0xad, 0x52, 0x5c, 0x51, 0x1b, 0x3a,
	0xd2, 0x84, 0x98, 0xeb, 0x3d, 0x53, 0x1a, 0x65, 0x72, 0x8d, 0x3f, 0x70, 0x54, 0xe1, 0x7f, 0x96,
	0x88, 0x43, 0x6f, 0x60, 0x5e, 0xad, 0xba, 0xcb, 0xa5, 0xce, 0xf5, 0xd2, 0xfb, 0x5f, 0xbf, 0xa9,
	0xbc, 0xe4, 0xb9, 0xd1, 0xeb, 0x9e, 0x1b, 0xb2, 0x2f, 0xef, 0x45, 0x10, 0x88, 0x47, 0x2d, 0xf1,
	0x86, 0x5e, 0xc1, 0xc4, 0x70, 0xee, 0xfd, 0x87, 0x17, 0xfa, 0xee, 0x6f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xcf, 0x84, 0xa6, 0x16, 0xe5, 0x02, 0x00, 0x00,
}
