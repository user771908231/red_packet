// Code generated by protoc-gen-go.
// source: n.proto
// DO NOT EDIT!

package bbproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type N struct {
	Name             *string `protobuf:"bytes,2,req,name=name" json:"name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *N) Reset()                    { *m = N{} }
func (m *N) String() string            { return proto.CompactTextString(m) }
func (*N) ProtoMessage()               {}
func (*N) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *N) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*N)(nil), "bbproto.N")
}

var fileDescriptor3 = []byte{
	// 58 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x62, 0xcf, 0xd3, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f, 0x4a, 0x02, 0x33, 0x94, 0x04, 0xb9, 0x18, 0xfd, 0x84, 0x78,
	0xb8, 0x58, 0xf2, 0x12, 0x73, 0x53, 0x25, 0x98, 0x14, 0x98, 0x34, 0x38, 0x01, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xf4, 0x20, 0xac, 0x6a, 0x25, 0x00, 0x00, 0x00,
}
