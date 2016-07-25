package mongodb

import (
	"net"
	"casino_server/msg/bbprotogo"
	"fmt"
	"casino_server/utils/test"
	"testing"
)



func TestOg(t *testing.T) {
	//game_EnterMatch(10006)

	go game_EnterMatch(10007)
	go game_EnterMatch(10008)
	go game_EnterMatch(10009)
	//game_EnterMatch(10010)
	//game_EnterMatch(10011)
	//gamelogingame(1111)
	//ogbet(1,20)
	for ; ;  {
		
	}
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


func game_EnterMatch(userId int32){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ide2 := int32(31)

	fmt.Println("proto 得到的id ", ide2)
	data2 := &bbproto.Game_EnterMatch{}
	data2.Tableid = &userId
	m2 := test.AssembleDataNomd5(uint16(ide2), data2)
	conn.Write(m2)
	_ = test.Read(conn).(*bbproto.Game_SendGameInfo)
	_ = test.Read(conn).(*bbproto.Game_BlindCoin)
	_ = test.Read(conn).(*bbproto.Game_BlindCoin)
	_ = test.Read(conn).(*bbproto.Game_BlindCoin)
	_ = test.Read(conn).(*bbproto.Game_BlindCoin)
	_ = test.Read(conn).(*bbproto.Game_BlindCoin)
	_ = test.Read(conn).(*bbproto.Game_BlindCoin)

}


//押注
func ogbet(seatId int32,coin int64){
	var tableId int32 = 0
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ide2 := int32(bbproto.EProtoId_GAME_FOLLOWBET)

	fmt.Println("proto 得到的id ", ide2)
	followData := &bbproto.Game_FollowBet{}
	followData.Tableid = &tableId
	followData.Seat = &seatId
	m2 := test.AssembleDataNomd5(uint16(ide2), followData)
	conn.Write(m2)
	_ = test.Read(conn).(*bbproto.Game_AckFollowBet)
}


