package mongodb

import (
	"net"
	"casino_server/msg/bbproto"
	"testing"
	"casino_server/utils/test"
	"fmt"
	"casino_server/msg"
)

const url = "192.168.199.120:3563"
const TCP = "tcp"

func TestGetIntoRoom(t *testing.T) {
	fun1(t)
}

func fun1(t *testing.T) {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}

	var id uint16 = 4
	var data bbproto.GetIntoRoom
	var header bbproto.ProtoHeader
	var userId uint32 = 989
	header.UserId = &userId

	m := utils.AssembleData(id,&data)
	conn.Write(m)

	var res [250]byte
	count, err := conn.Read(res[0:])
	if err != nil {
		fmt.Println("err != nil")
	}

	t.Log("读取到的 res %v", res)
	msg2, err := msg.PortoProcessor.Unmarshal(res[2:count])
	if err != nil {
	}
	m5 := msg2.(*bbproto.GetIntoRoom)
	fmt.Println("m5.UserId:", m5)

	for ; ; {

	}
}
