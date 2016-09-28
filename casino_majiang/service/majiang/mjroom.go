package majiang

import (
	"casino_majiang/msg/funcsInit"
	"casino_majiang/conf/log"
)


//普通的麻将房间...

//room的结构定义在proto中
//type MjRoom struct {
//	roomType int32
//
//}


func init() {
	FMJRoomIns = newProto.NewMjRoom()
	FMJRoomIns.OnInit()
}

var FMJRoomIns *MjRoom

//初始化
func (r *MjRoom) OnInit() {
	log.T("初始化麻将的room")


}

func (r *MjRoom) CreateDesk() *MjDesk {
	//create 的时候，是否需要通过type 来判断,怎么样创建房间
	desk :=newProto.NewMjDesk()
	return desk
}

//room中增加一个desk
func (r *MjRoom) AddDesk() error {
	return nil

}

