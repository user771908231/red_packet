package main

import (
	"testing"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_laowangye/service/niuniu"
)

func GetDesk() *laowangye.Desk {
	option := &ddproto.LwyDeskOption{
		MinUser: proto.Int32(2),
		MaxUser: proto.Int32(6),
		MaxCircle: proto.Int32(20),
		HasFlower: proto.Bool(true),
		BankRule: ddproto.LwyEnumBankerRule_DING_ZHUANG.Enum(),
		IsFlowerPlay:proto.Bool(true),
		IsJiaoFenJiaBei: proto.Bool(true),
	}
	room,_ := laowangye.Rooms.GetRoomById(0)
	new_desk,_ := room.CreateFriendDesk(option, 1)
	return new_desk
}

//测试：创建desk
func TestRoom_ChreateDesk(t *testing.T) {
	option := &ddproto.LwyDeskOption{
		MinUser: proto.Int32(2),
		MaxUser: proto.Int32(6),
		MaxCircle: proto.Int32(20),
		HasFlower: proto.Bool(true),
		BankRule: ddproto.LwyEnumBankerRule_DING_ZHUANG.Enum(),
		IsFlowerPlay:proto.Bool(true),
		IsJiaoFenJiaBei: proto.Bool(true),
	}
	t.Log(laowangye.Rooms)
	room,_ := laowangye.Rooms.GetRoomById(0)
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
