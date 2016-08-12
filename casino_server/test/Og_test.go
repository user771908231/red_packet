package mongodb

import (
	"net"
	"casino_server/msg/bbprotogo"
	"fmt"
	"casino_server/utils/test"
	"testing"
)

func TestOg(t *testing.T) {
	//game_EnterMatch(10007)
	//game_EnterMatch(10008)
	//game_EnterMatch(10009)
	//game_EnterMatch(10010)
	//game_EnterMatch(10011)
	//ogbet(0,20)
	//ogbet(1,20)
	//ogbet(2,20)
	//ogbet(3,20)

	//ogRaise(0,20)
	//ogRaise(1,20)
	//ogRaise(2,20)
	//ogRaise(3,20)

	//rEQQuickConn(10006)

	createDesk(10084)

	//getRecords(10084)

	for ; ; {
	}
}

func game_EnterMatch(userId uint32) {
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
func ogbet(seatId int32, coin int64) {
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
	//_ = test.Read(conn).(*bbproto.Game_AckFollowBet)
}



//用户登陆

func ogRaise(seatId int32, coin int64) {
	var tableId int32 = 0
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ide2 := int32(bbproto.EProtoId_PID_GAME_RAISEBET)

	fmt.Println("proto 得到的id ", ide2)
	raiseData := &bbproto.Game_RaiseBet{}
	raiseData.Tableid = &tableId
	raiseData.Seat = &seatId
	raiseData.Coin = &coin
	m2 := test.AssembleDataNomd5(uint16(ide2), raiseData)
	conn.Write(m2)
	//_ = test.Read(conn).(*bbproto.Game_AckRaiseBet)
}


//创建一个房间
func createDesk(userId uint32) {
	pid := int32(bbproto.EProtoId_PID_GAME_GAME_CREATEDESK)
	reqData := &bbproto.Game_CreateDesk{}
	reqData.BigBlind = new(int64)
	reqData.SmallBlind = new(int64)
	reqData.InitCount = new(int32)
	reqData.InitCoin = new(int64)
	reqData.UserId = new(uint32)
	reqData.Password = new(string)

	*reqData.BigBlind = 30
	*reqData.SmallBlind = 15
	*reqData.InitCount = 20
	*reqData.InitCoin = 1000
	*reqData.UserId = userId
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	m2 := test.AssembleDataNomd5(uint16(pid), reqData)
	conn.Write(m2)

	//读取放回的信息
	test.Read(conn)
}

func getRecords(userId uint32) {

	//userId := 10084
	//
	////1,获取数据库连接
	//c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	//if err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}
	//defer c.Close()
	//
	//s := c.Ref()
	//defer c.UnRef(s)
	//
	////开始查询
	//var rets []mode.T_th_record
	//s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TH_RECORD).Find(bson.M{"userid": userId}).Sort("-winamount").Skip(3).Limit(4).All(&rets)
	//fmt.Println("rets:【%v】", rets)
	//
	//for i := 0; i < len(rets); i++ {
	//	u := rets[i]
	//	fmt.Println("rets[%v]",u.Id, u.UserId,u.WinAmount)
	//}


	pid := int32(bbproto.EProtoId_PID_GAME_GAME_GAMERECORD)
	reqData := &bbproto.Game_GameRecord{}
	reqData.UserId = new(uint32)
	*reqData.UserId = userId

	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	m2 := test.AssembleDataNomd5(uint16(pid), reqData)
	conn.Write(m2)
	//读取放回的信息
	test.Read(conn)
}

