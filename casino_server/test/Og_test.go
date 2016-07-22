package mongodb

import (
	"net"
	"casino_server/msg/bbprotogo"
	"fmt"
	"casino_server/utils/test"
	"testing"
)

func TestOg(t *testing.T) {
	game_EnterMatch()
	//gamelogingame(1111)
}

func gamelogingame(userId uint32) {
	//获得连接
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//ide2 := int32(bbproto.EProtoId_THROOM)
	ide2 := int32(30)

	fmt.Println("proto 得到的id ", ide2)
	data2 := &bbproto.Game_LoginGame{}
	m2 := test.AssembleDataNomd5(uint16(ide2), data2)
	conn.Write(m2)
	result := test.Read(conn).(*bbproto.Game_LoginGame)
	fmt.Println("读取的结果:", result)                //测试服务器同意返回98989
}


func game_EnterMatch(){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ide2 := int32(31)

	fmt.Println("proto 得到的id ", ide2)
	data2 := &bbproto.Game_EnterMatch{}
	m2 := test.AssembleDataNomd5(uint16(ide2), data2)
	conn.Write(m2)
	test.Read(conn).(*bbproto.Game_SendGameInfo)
}
