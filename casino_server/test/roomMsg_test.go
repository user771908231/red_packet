package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/test"
)

func TestRoomMsg(t *testing.T){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var id uint16 = 5
	var roomId int32 = 1000
	var data bbproto.RoomMsg
	sendMsg := "测试发送的信息"
	data.RoomId = &roomId
	data.Msg    = &sendMsg
	m := utils.AssembleData(id,&data)
	conn.Write(m)
}
