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
	// game_EnterMatch(10007)
	// game_EnterMatch(10008)
	// game_EnterMatch(10009)
	// game_EnterMatch(10010)
	 //game_EnterMatch(10011)
	//gamelogingame(1111)
	//ogbet(0,20)
	//ogbet(1,20)
	//ogbet(2,20)
	ogbet(3,20)

	//rEQQuickConn(10006)
	for ; ;  {
	}
}

func game_EnterMatch(userId uint32){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	ide2 := int32(31)

	fmt.Println("proto 得到的id ", ide2)
	data2 := &bbproto.Game_EnterMatch{}
	data2.UserId = &userId
	m2 := test.AssembleDataNomd5(uint16(ide2), data2)
	conn.Write(m2)
}

//押注
func ogbet(seatId int32,coin int64){
	var tableId int32 = 0
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ide2 := int32(bbproto.EProtoId_PID_GAME_FOLLOWBET)

	fmt.Println("proto 得到的id ", ide2)
	followData := &bbproto.Game_FollowBet{}
	followData.Tableid = &tableId
	followData.Seat = &seatId
	m2 := test.AssembleDataNomd5(uint16(ide2), followData)
	conn.Write(m2)
	_ = test.Read(conn).(*bbproto.Game_AckFollowBet)
}



//用户登陆
func rEQQuickConn(userId uint32){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ide2 := int32(bbproto.EProtoId_REQQUICKCONN)

	fmt.Println("proto 得到的id ", ide2)
	followData := &bbproto.REQQuickConn{}
	followData.UserId = &userId
	m2 := test.AssembleDataNomd5(uint16(ide2), followData)
	conn.Write(m2)
	_ = test.Read(conn).(*bbproto.ACKQuickConn)

}
