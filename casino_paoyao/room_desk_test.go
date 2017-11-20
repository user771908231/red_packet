package main

import (
	"testing"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_paoyao/service/niuniu"
)

func GetDesk() *paoyao.Desk {
	option := &ddproto.PaoyaoniuDeskOption{
		MinUser: proto.Int32(2),
		MaxUser: proto.Int32(6),
		MaxCircle: proto.Int32(20),
		HasFlower: proto.Bool(true),
		BankRule: ddproto.PaoyaoniuEnumBankerRule_DING_ZHUANG.Enum(),
		IsFlowerPlay:proto.Bool(true),
		IsJiaoFenJiaBei: proto.Bool(true),
	}
	room,_ := paoyao.Rooms.GetRoomById(0)
	new_desk,_ := room.CreateFriendDesk(option, 1)
	return new_desk
}

//测试：创建desk
func TestRoom_ChreateDesk(t *testing.T) {
	option := &ddproto.PaoyaoniuDeskOption{
		MinUser: proto.Int32(2),
		MaxUser: proto.Int32(6),
		MaxCircle: proto.Int32(20),
		HasFlower: proto.Bool(true),
		BankRule: ddproto.PaoyaoniuEnumBankerRule_DING_ZHUANG.Enum(),
		IsFlowerPlay:proto.Bool(true),
		IsJiaoFenJiaBei: proto.Bool(true),
	}
	t.Log(paoyao.Rooms)
	room,_ := paoyao.Rooms.GetRoomById(0)
	t.Log(room)
	new_desk,err := room.CreateFriendDesk(option, 1)
	t.Log(new_desk, err)
}

//测试: 房间增加用户
func TestDeskAddUser(t *testing.T) {
	desk := GetDesk()
	t.Log(desk)
	user,err := desk.AddUser(1, nil)
	t.Log(user, err)
	user,err = desk.AddUser(1, nil)
	t.Log(user, err)
	user,err = desk.AddUser(2, nil)
	t.Log(user, err)
}
