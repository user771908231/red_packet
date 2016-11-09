// Code generated by protoc-gen-go.
// source: ddz_server.proto
// DO NOT EDIT!

/*
Package doudizhu is a generated protocol buffer package.

It is generated from these files:
	ddz_server.proto

It has these top-level messages:
	PPokerPai
	PDdzDesk
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

// desk
type PDdzDesk struct {
	DeskId           *int32  `protobuf:"varint,1,opt,name=deskId" json:"deskId,omitempty"`
	Key              *string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	UserCountLimit   *int32  `protobuf:"varint,3,opt,name=userCountLimit" json:"userCountLimit,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PDdzDesk) Reset()                    { *m = PDdzDesk{} }
func (m *PDdzDesk) String() string            { return proto.CompactTextString(m) }
func (*PDdzDesk) ProtoMessage()               {}
func (*PDdzDesk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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

// user
type PDdzUser struct {
	UserId           *uint32 `protobuf:"varint,1,opt,name=userId" json:"userId,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PDdzUser) Reset()                    { *m = PDdzUser{} }
func (m *PDdzUser) String() string            { return proto.CompactTextString(m) }
func (*PDdzUser) ProtoMessage()               {}
func (*PDdzUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PDdzUser) GetUserId() uint32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// room
type PDdzRoom struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *PDdzRoom) Reset()                    { *m = PDdzRoom{} }
func (m *PDdzRoom) String() string            { return proto.CompactTextString(m) }
func (*PDdzRoom) ProtoMessage()               {}
func (*PDdzRoom) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// 备份专用...
type PDdzbak struct {
	Desk             *PDdzDesk   `protobuf:"bytes,1,opt,name=desk" json:"desk,omitempty"`
	Users            []*PDdzUser `protobuf:"bytes,2,rep,name=users" json:"users,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *PDdzbak) Reset()                    { *m = PDdzbak{} }
func (m *PDdzbak) String() string            { return proto.CompactTextString(m) }
func (*PDdzbak) ProtoMessage()               {}
func (*PDdzbak) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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
func (*PDdzSession) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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
	proto.RegisterType((*PDdzDesk)(nil), "doudizhu.PDdzDesk")
	proto.RegisterType((*PDdzUser)(nil), "doudizhu.PDdzUser")
	proto.RegisterType((*PDdzRoom)(nil), "doudizhu.PDdzRoom")
	proto.RegisterType((*PDdzbak)(nil), "doudizhu.PDdzbak")
	proto.RegisterType((*PDdzSession)(nil), "doudizhu.PDdzSession")
}

var fileDescriptor0 = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x90, 0x4d, 0x4b, 0x33, 0x41,
	0x10, 0x84, 0xc9, 0x7e, 0xbc, 0x6f, 0xd2, 0x6b, 0x82, 0xcc, 0x41, 0x16, 0x4f, 0x71, 0x4e, 0x9e,
	0xf6, 0xe0, 0x1f, 0xf0, 0xe0, 0x82, 0x08, 0x22, 0xab, 0x21, 0x67, 0x19, 0x99, 0x16, 0x87, 0xc9,
	0xee, 0xc8, 0x7c, 0x44, 0xdc, 0x5f, 0xef, 0xf4, 0x4e, 0x22, 0x01, 0x6f, 0x5d, 0x74, 0xf1, 0x54,
	0x77, 0xc1, 0xb9, 0x94, 0xe3, 0xab, 0x43, 0xbb, 0x47, 0xdb, 0x7c, 0x5a, 0xe3, 0x0d, 0x9b, 0x4b,
	0x13, 0xa4, 0x1a, 0x3f, 0x02, 0x7f, 0x86, 0x45, 0xd7, 0x19, 0x8d, 0xb6, 0x13, 0x8a, 0x01, 0x64,
	0x4a, 0xd6, 0xb3, 0xf5, 0xec, 0xba, 0x64, 0x15, 0xe4, 0x12, 0x5d, 0x9d, 0x45, 0xb1, 0x60, 0x4b,
	0x28, 0xf7, 0x62, 0x17, 0xb0, 0xce, 0xa7, 0xdd, 0x0a, 0xfe, 0xbd, 0xef, 0xcc, 0x17, 0xda, 0xba,
	0x98, 0xf4, 0x19, 0x14, 0x83, 0xe8, 0xb1, 0x2e, 0xc9, 0xcc, 0x6f, 0x61, 0xde, 0xb5, 0x72, 0x6c,
	0xd1, 0x69, 0x72, 0x46, 0x8a, 0x7e, 0x38, 0xa1, 0x6a, 0xfc, 0x3e, 0x50, 0x2f, 0x60, 0x15, 0xe2,
	0x59, 0x77, 0x26, 0x0c, 0xfe, 0x51, 0xf5, 0xca, 0x27, 0x3c, 0xbf, 0x4c, 0x80, 0x6d, 0xdc, 0x11,
	0x80, 0x3c, 0x07, 0xc0, 0x92, 0x43, 0xda, 0xbd, 0x18, 0xd3, 0xf3, 0x27, 0xf8, 0x4f, 0xf3, 0x9b,
	0xd0, 0x6c, 0x0d, 0x05, 0xe5, 0x4c, 0xa6, 0xea, 0x86, 0x35, 0xc7, 0xff, 0x9a, 0xdf, 0x4b, 0xae,
	0xa0, 0x24, 0x10, 0x7d, 0x94, 0xff, 0xb5, 0x50, 0x56, 0xec, 0xa2, 0xa2, 0x79, 0x83, 0xce, 0x29,
	0x33, 0x50, 0xf4, 0xf6, 0x24, 0x9a, 0x34, 0xc5, 0x46, 0x9d, 0x1d, 0x5b, 0x68, 0xd3, 0x6f, 0xa9,
	0x95, 0x58, 0xdf, 0x7d, 0x6c, 0x61, 0xe3, 0x85, 0x0f, 0x2e, 0x35, 0xf3, 0x13, 0x00, 0x00, 0xff,
	0xff, 0xbc, 0xa8, 0x44, 0x9a, 0x7b, 0x01, 0x00, 0x00,
}