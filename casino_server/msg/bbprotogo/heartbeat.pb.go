// Code generated by protoc-gen-go.
// source: heartbeat.proto
// DO NOT EDIT!

package bbproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type HeatBeat struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *HeatBeat) Reset()                    { *m = HeatBeat{} }
func (m *HeatBeat) String() string            { return proto.CompactTextString(m) }
func (*HeatBeat) ProtoMessage()               {}
func (*HeatBeat) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func init() {
	proto.RegisterType((*HeatBeat)(nil), "bbproto.HeatBeat")
}

var fileDescriptor3 = []byte{
	// 56 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0x48, 0x4d, 0x2c,
	0x2a, 0x49, 0x4a, 0x4d, 0x2c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f, 0x4a, 0x02,
	0x33, 0x94, 0xb8, 0xb8, 0x38, 0x3c, 0x80, 0xc2, 0x4e, 0x40, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff,
	0x49, 0x63, 0x69, 0x10, 0x26, 0x00, 0x00, 0x00,
}