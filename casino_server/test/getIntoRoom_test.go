package mongodb

import (
	"net"
	"casino_server/msg/bbprotogo"
	"testing"
	"casino_server/utils/test"
	"fmt"
	"casino_server/msg"
)

const url = "192.168.199.120:3797"
//const url = "182.92.179.230:3797"
const TCP = "tcp"

func TestGetIntoRoom(t *testing.T) {
	fun1(t)
}

func fun1(t *testing.T) {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var id uint16 = 4
	var data bbproto.GetIntoRoom
	var userId uint32 = 10013
	var inValue int32 = 1

	data.In = &inValue
	data.UserId = &userId
	m := test.AssembleData(id,&data)
	conn.Write(m)


	for ; ; {
		fmt.Println("开始读取广播消息")
		var res [250]byte
		count,err := conn.Read(res[0:])
		if err != nil {
			fmt.Println("err != nil")
		}
		fmt.Println("读取到广播消息",res)

		msg2, err := msg.ProtoProcessor.Unmarshal(res[2:count])
		if err != nil {
		}
		m5 :=  msg2.(*bbproto.RoomMsg)
		fmt.Println("读取广播消息结束roomid",m5.GetRoomId())
		fmt.Println("读取广播消息结束roomsg",m5.GetMsg())
		fmt.Println("读取广播消息结束")
	}
}
