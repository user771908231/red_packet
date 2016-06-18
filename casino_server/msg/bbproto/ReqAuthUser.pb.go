// Code generated by protoc-gen-go.
// source: ReqAuthUser.proto
// DO NOT EDIT!

/*
Package bbproto is a generated protocol buffer package.

It is generated from these files:
	ReqAuthUser.proto
	base.proto
	getIntoRoom.proto
	heartbeat.proto
	reg.proto
	roomMsg.proto
	testp1.proto
	user.proto

It has these top-level messages:
	ReqAuthUser
	ProtoHeader
	TerminalInfo
	GetIntoRoom
	HeatBeat
	Reg
	RoomMsg
	TestP1
	User
*/
package bbproto

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

// Ignoring public import of ProtoHeader from base.proto

// Ignoring public import of TerminalInfo from base.proto

// Ignoring public import of EUnitType from base.proto

// Ignoring public import of EUnitRace from base.proto

type ReqAuthUser struct {
	Header           *ProtoHeader  `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Terminal         *TerminalInfo `protobuf:"bytes,2,opt,name=terminal" json:"terminal,omitempty"`
	SelectRole       *uint32       `protobuf:"varint,3,opt,name=selectRole" json:"selectRole,omitempty"`
	AppVersion       *int32        `protobuf:"varint,4,opt,name=appVersion" json:"appVersion,omitempty"`
	Uuid             *string       `protobuf:"bytes,5,opt,name=uuid" json:"uuid,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *ReqAuthUser) Reset()                    { *m = ReqAuthUser{} }
func (m *ReqAuthUser) String() string            { return proto.CompactTextString(m) }
func (*ReqAuthUser) ProtoMessage()               {}
func (*ReqAuthUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ReqAuthUser) GetHeader() *ProtoHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ReqAuthUser) GetTerminal() *TerminalInfo {
	if m != nil {
		return m.Terminal
	}
	return nil
}

func (m *ReqAuthUser) GetSelectRole() uint32 {
	if m != nil && m.SelectRole != nil {
		return *m.SelectRole
	}
	return 0
}

func (m *ReqAuthUser) GetAppVersion() int32 {
	if m != nil && m.AppVersion != nil {
		return *m.AppVersion
	}
	return 0
}

func (m *ReqAuthUser) GetUuid() string {
	if m != nil && m.Uuid != nil {
		return *m.Uuid
	}
	return ""
}

func init() {
	proto.RegisterType((*ReqAuthUser)(nil), "bbproto.ReqAuthUser")
}

var fileDescriptor0 = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x0c, 0x4a, 0x2d, 0x74,
	0x2c, 0x2d, 0xc9, 0x08, 0x2d, 0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f,
	0x4a, 0x02, 0x33, 0xa4, 0xb8, 0x92, 0x12, 0x8b, 0x53, 0x21, 0x82, 0x4a, 0x93, 0x18, 0xb9, 0xb8,
	0x91, 0x94, 0x0a, 0xa9, 0x70, 0xb1, 0x65, 0xa4, 0x26, 0xa6, 0xa4, 0x16, 0x49, 0x30, 0x2a, 0x30,
	0x6a, 0x70, 0x1b, 0x89, 0xe8, 0x41, 0x75, 0xe9, 0x05, 0x80, 0x48, 0x0f, 0xb0, 0x9c, 0x90, 0x3a,
	0x17, 0x47, 0x49, 0x6a, 0x51, 0x6e, 0x66, 0x5e, 0x62, 0x8e, 0x04, 0x13, 0x58, 0x9d, 0x28, 0x5c,
	0x5d, 0x08, 0x54, 0xc2, 0x33, 0x2f, 0x2d, 0x5f, 0x48, 0x88, 0x8b, 0xab, 0x38, 0x35, 0x27, 0x35,
	0xb9, 0x24, 0x28, 0x3f, 0x27, 0x55, 0x82, 0x19, 0xa8, 0x94, 0x17, 0x24, 0x96, 0x58, 0x50, 0x10,
	0x96, 0x5a, 0x54, 0x9c, 0x99, 0x9f, 0x27, 0xc1, 0x02, 0x14, 0x63, 0x15, 0xe2, 0xe1, 0x62, 0x29,
	0x2d, 0xcd, 0x4c, 0x91, 0x60, 0x05, 0xf2, 0x38, 0x03, 0x18, 0x00, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xf1, 0xc5, 0x23, 0x27, 0xbf, 0x00, 0x00, 0x00,
}
