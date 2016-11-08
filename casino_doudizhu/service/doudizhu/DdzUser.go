package doudizhu

import "sync"

type TestModel struct {
	TId   *int32 `protobuf:"varint,2,opt,name=TId" json:"TId,omitempty"`
	TUser *string `protobuf:"bytes,5,opt,name=TUser" json:"TUser,omitempty"`
}

type DdzUser struct {
	sync.Mutex
	Id   *int32 `protobuf:"varint,2,opt,name=Id" json:"Id,omitempty"`
	User *string `protobuf:"bytes,5,opt,name=User" json:"User,omitempty"`
	Test *TestModel `protobuf:"bytes,18,opt,name=Test" json:"Test,omitempty"`
}

//实现proto的接口，方便redis存储
func (d *DdzUser)  Reset() {

}

func (d *DdzUser)  String() string {
	return ""

}

func (d *DdzUser)  ProtoMessage() {

}

//清楚session
func (u *DdzUser)ClearAgentGameData() {

}
