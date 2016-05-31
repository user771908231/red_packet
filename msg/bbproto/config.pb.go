// Code generated by protoc-gen-go.
// source: config.proto
// DO NOT EDIT!

package bbproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EPlayType int32

const (
	EPlayType_ONCE EPlayType = 1
	EPlayType_LOOP EPlayType = 2
)

var EPlayType_name = map[int32]string{
	1: "ONCE",
	2: "LOOP",
}
var EPlayType_value = map[string]int32{
	"ONCE": 1,
	"LOOP": 2,
}

func (x EPlayType) Enum() *EPlayType {
	p := new(EPlayType)
	*p = x
	return p
}
func (x EPlayType) String() string {
	return proto.EnumName(EPlayType_name, int32(x))
}
func (x *EPlayType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(EPlayType_value, data, "EPlayType")
	if err != nil {
		return err
	}
	*x = EPlayType(value)
	return nil
}
func (EPlayType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

// Audio config file
type AudioConfigItem struct {
	Version          *int32     `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Id               *int32     `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	Name             *string    `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	ResourcePath     *string    `protobuf:"bytes,4,opt,name=resourcePath" json:"resourcePath,omitempty"`
	Type             *EPlayType `protobuf:"varint,5,opt,name=type,enum=bbproto.EPlayType" json:"type,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *AudioConfigItem) Reset()                    { *m = AudioConfigItem{} }
func (m *AudioConfigItem) String() string            { return proto.CompactTextString(m) }
func (*AudioConfigItem) ProtoMessage()               {}
func (*AudioConfigItem) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *AudioConfigItem) GetVersion() int32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *AudioConfigItem) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *AudioConfigItem) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *AudioConfigItem) GetResourcePath() string {
	if m != nil && m.ResourcePath != nil {
		return *m.ResourcePath
	}
	return ""
}

func (m *AudioConfigItem) GetType() EPlayType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return EPlayType_ONCE
}

type AudioConfigFile struct {
	AudioConfig      []*AudioConfigItem `protobuf:"bytes,1,rep,name=audioConfig" json:"audioConfig,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *AudioConfigFile) Reset()                    { *m = AudioConfigFile{} }
func (m *AudioConfigFile) String() string            { return proto.CompactTextString(m) }
func (*AudioConfigFile) ProtoMessage()               {}
func (*AudioConfigFile) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *AudioConfigFile) GetAudioConfig() []*AudioConfigItem {
	if m != nil {
		return m.AudioConfig
	}
	return nil
}

func init() {
	proto.RegisterType((*AudioConfigItem)(nil), "bbproto.AudioConfigItem")
	proto.RegisterType((*AudioConfigFile)(nil), "bbproto.AudioConfigFile")
	proto.RegisterEnum("bbproto.EPlayType", EPlayType_name, EPlayType_value)
}

var fileDescriptor2 = []byte{
	// 194 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0x4b,
	0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f, 0x4a, 0x02, 0x33, 0x94, 0x4a,
	0xb8, 0xf8, 0x1d, 0x4b, 0x53, 0x32, 0xf3, 0x9d, 0xc1, 0xb2, 0x9e, 0x25, 0xa9, 0xb9, 0x42, 0xfc,
	0x5c, 0xec, 0x65, 0xa9, 0x45, 0xc5, 0x99, 0xf9, 0x79, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xac, 0x42,
	0x5c, 0x5c, 0x4c, 0x99, 0x29, 0x12, 0x4c, 0x60, 0x36, 0x0f, 0x17, 0x4b, 0x5e, 0x62, 0x6e, 0xaa,
	0x04, 0x33, 0x90, 0xc7, 0x29, 0x24, 0xc2, 0xc5, 0x53, 0x94, 0x5a, 0x9c, 0x5f, 0x5a, 0x94, 0x9c,
	0x1a, 0x90, 0x58, 0x92, 0x21, 0xc1, 0x02, 0x16, 0x55, 0xe0, 0x62, 0x29, 0xa9, 0x2c, 0x48, 0x95,
	0x60, 0x05, 0xf2, 0xf8, 0x8c, 0x84, 0xf4, 0xa0, 0x76, 0xe9, 0xb9, 0x06, 0xe4, 0x24, 0x56, 0x86,
	0x00, 0x65, 0x94, 0x1c, 0x50, 0x6c, 0x75, 0xcb, 0xcc, 0x49, 0x15, 0xd2, 0xe5, 0xe2, 0x4e, 0x44,
	0x08, 0x01, 0x6d, 0x66, 0xd6, 0xe0, 0x36, 0x92, 0x80, 0xeb, 0x45, 0x73, 0xa4, 0x96, 0x3c, 0x17,
	0x27, 0xdc, 0x38, 0x21, 0x0e, 0x2e, 0x16, 0x7f, 0x3f, 0x67, 0x57, 0x01, 0x46, 0x10, 0xcb, 0xc7,
	0xdf, 0x3f, 0x40, 0x80, 0x09, 0x10, 0x00, 0x00, 0xff, 0xff, 0xf7, 0x3b, 0x69, 0x1f, 0xf0, 0x00,
	0x00, 0x00,
}
