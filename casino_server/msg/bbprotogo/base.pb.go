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
	EProtoId_ZJHROOM              EProtoId = 10
	EProtoId_ZJHBET               EProtoId = 11
	EProtoId_ZJHMSG               EProtoId = 12
	EProtoId_ZJHQUERYNOSEATUSER   EProtoId = 13
	EProtoId_ZJHREQSEAT           EProtoId = 14
	EProtoId_ZJHLOTTERY           EProtoId = 15
	EProtoId_ZJHBROADCASTBEGINBET EProtoId = 16
	EProtoId_LOGINSIGNINBONUS     EProtoId = 17
	EProtoId_LOGINTURNTABLEBONUS  EProtoId = 18
	EProtoId_OLINEBONUS           EProtoId = 19
	EProtoId_TIMINGBONUS          EProtoId = 20
	EProtoId_THROOM               EProtoId = 21
	EProtoId_THBET                EProtoId = 22
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
}
var EProtoId_value = map[string]int32{
	"TESTP1":               0,
	"REG":                  1,
	"REQAUTHUSER":          2,
	"HEATBEAT":             3,
	"GETINTOROOM":          4,
	"ROOMMSG":              5,
	"GETREWARDS":           6,
	"SHUIGUOJI":            7,
	"SHUIGUOJIHILOMP":      8,
	"SHUIGUOJIRES":         9,
	"ZJHROOM":              10,
	"ZJHBET":               11,
	"ZJHMSG":               12,
	"ZJHQUERYNOSEATUSER":   13,
	"ZJHREQSEAT":           14,
	"ZJHLOTTERY":           15,
	"ZJHBROADCASTBEGINBET": 16,
	"LOGINSIGNINBONUS":     17,
	"LOGINTURNTABLEBONUS":  18,
	"OLINEBONUS":           19,
	"TIMINGBONUS":          20,
	"THROOM":               21,
	"THBET":                22,
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
func (EProtoId) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

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
func (*ProtoHeader) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

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
func (*TerminalInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

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

var fileDescriptor1 = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x3c, 0x50, 0x4d, 0x8f, 0xd3, 0x30,
	0x14, 0xa4, 0x9f, 0x69, 0x5e, 0xdb, 0xad, 0xd7, 0x2d, 0x4b, 0x8e, 0x68, 0x4f, 0x88, 0x03, 0x12,
	0x3f, 0x21, 0x65, 0xad, 0xc4, 0xab, 0x24, 0x6e, 0x6d, 0x07, 0x04, 0xb7, 0xb4, 0x31, 0x10, 0xb1,
	0x9b, 0x54, 0x49, 0x8b, 0x38, 0xc2, 0x3f, 0xe7, 0xd9, 0x8d, 0xf6, 0x10, 0xe9, 0xcd, 0xbc, 0xc9,
	0xbc, 0x19, 0x03, 0x1c, 0x8a, 0xce, 0x7c, 0x38, 0xb5, 0xcd, 0xb9, 0xa1, 0xde, 0xe1, 0xe0, 0x86,
	0xfb, 0x7f, 0x03, 0x98, 0xef, 0xec, 0x14, 0x9b, 0xa2, 0x34, 0x2d, 0xbd, 0x81, 0x69, 0x71, 0xaa,
	0x3e, 0x9b, 0x36, 0x18, 0xbc, 0x1d, 0xbc, 0xf3, 0xe9, 0x2d, 0xf8, 0x9d, 0xe9, 0xba, 0xaa, 0xa9,
	0x79, 0x19, 0x0c, 0x1d, 0x85, 0x92, 0x4b, 0x67, 0x5a, 0xc4, 0x23, 0xc4, 0x4b, 0x4a, 0x60, 0x76,
	0x2a, 0x8e, 0xbf, 0xcc, 0x19, 0x99, 0x31, 0x32, 0x13, 0xba, 0x80, 0xf1, 0xb1, 0x29, 0x4d, 0x30,
	0x71, 0x68, 0x09, 0x13, 0xd3, 0xb6, 0x4d, 0x1b, 0x4c, 0xdd, 0xef, 0x28, 0x37, 0x7f, 0xce, 0x6d,
	0xa1, 0x8b, 0x1f, 0x81, 0x67, 0x05, 0xf7, 0x29, 0x2c, 0xb4, 0x69, 0x9f, 0xab, 0xba, 0x78, 0xe2,
	0xf5, 0xf7, 0x86, 0xae, 0xc0, 0x3b, 0xfe, 0x2c, 0xea, 0xda, 0x3c, 0xf5, 0x21, 0x28, 0x40, 0x69,
	0x7e, 0x57, 0x47, 0x93, 0x15, 0xcf, 0xa6, 0x4f, 0x81, 0x37, 0x2e, 0x97, 0xea, 0x9a, 0xc1, 0x47,
	0xc1, 0xb0, 0xe9, 0xdc, 0x75, 0xff, 0xfd, 0xdf, 0x11, 0xcc, 0x98, 0xeb, 0xc4, 0x4b, 0x5c, 0x4c,
	0x35, 0x53, 0x7a, 0xf7, 0x91, 0xbc, 0xa2, 0x1e, 0x8c, 0x24, 0x8b, 0xc8, 0x00, 0x0f, 0xcc, 0x25,
	0xdb, 0x87, 0xb9, 0x8e, 0x73, 0xc5, 0x24, 0x19, 0xa2, 0xd9, 0x2c, 0x66, 0xa1, 0xde, 0xe2, 0x47,
	0x46, 0x76, 0x1d, 0x31, 0xcd, 0x33, 0x2d, 0xa4, 0x10, 0x29, 0x19, 0xd3, 0x39, 0x78, 0x76, 0x4a,
	0x55, 0x44, 0x26, 0x58, 0x1f, 0x70, 0x2b, 0xd9, 0x97, 0x50, 0x3e, 0x28, 0x32, 0xc5, 0x7a, 0xbe,
	0x8a, 0x73, 0x1e, 0xe5, 0xe2, 0x91, 0x13, 0x8f, 0xae, 0x61, 0xf5, 0x02, 0x63, 0x9e, 0x88, 0x74,
	0x47, 0x66, 0xd8, 0x79, 0xf1, 0x42, 0x4a, 0xa6, 0x88, 0x6f, 0x2d, 0xbf, 0x3d, 0xc6, 0xce, 0x1f,
	0x6c, 0x48, 0x04, 0x5b, 0xa6, 0xc9, 0xbc, 0x9f, 0xed, 0xa9, 0x05, 0xbd, 0x03, 0x8a, 0xf3, 0x3e,
	0x67, 0xf2, 0x6b, 0x26, 0x14, 0x86, 0x73, 0x71, 0x97, 0x36, 0x82, 0xfd, 0x99, 0xed, 0x2d, 0x47,
	0x6e, 0x7a, 0x9c, 0x08, 0xad, 0x51, 0x49, 0x56, 0x34, 0x80, 0x8d, 0xf5, 0x93, 0x22, 0x7c, 0xf8,
	0x14, 0x2a, 0xac, 0x15, 0xf1, 0xcc, 0xba, 0x13, 0xba, 0x01, 0x92, 0x08, 0x44, 0x8a, 0x47, 0x19,
	0x72, 0x22, 0xcb, 0x15, 0xb9, 0xa5, 0x6f, 0x60, 0xed, 0x58, 0x9d, 0xcb, 0x4c, 0x87, 0xdb, 0x84,
	0x5d, 0x17, 0xd4, 0x1a, 0x8b, 0x84, 0x67, 0x3d, 0x5e, 0xdb, 0x97, 0xd1, 0x3c, 0xe5, 0x59, 0x74,
	0x25, 0x36, 0xee, 0x79, 0xaf, 0x2d, 0x5e, 0x53, 0x1f, 0x26, 0xda, 0x95, 0xb8, 0xfb, 0x1f, 0x00,
	0x00, 0xff, 0xff, 0xe5, 0x14, 0xfe, 0x32, 0x6b, 0x02, 0x00, 0x00,
}
