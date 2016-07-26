package mongodb

import (
	"testing"
	"casino_server/msg/bbprotogo"
	"fmt"
	"casino_server/utils/test"
	"net"
	"casino_server/service/room"
)

func TestTh(t *testing.T){
	//进入房间
	//intoRoom(10007)
	//intoRoom(10008)
	//intoRoom(10009)
	//intoRoom(10010)
	//intoRoom(10011)

	bet(10007,room.TH_DESK_BET_TYPE_CALL,0)	//跟注
	bet(10007,room.TH_DESK_BET_TYPE_RAISE,0)	//跟注
	bet(10007,room.TH_DESK_BET_TYPE_CHECK,0)	//跟注
	bet(10007,room.TH_DESK_BET_TYPE_FOLD,0)	//跟注
	bet(10007,room.TH_DESK_BET_TYPE_ALLIN,0)	//跟注

}



/** 182.92.179.230
	测试进入房间
 */
func intoRoom(userId uint32){
	//获得连接
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	protoId := int32(bbproto.EProtoId_THROOM)
	fmt.Println("proto 得到的id ", protoId)
	var userid uint32 = userId
	var reqType int32 = 1
	data := &bbproto.ThRoom{}
	h2 := &bbproto.ProtoHeader{}
	h2.UserId = &userid
	data.Header = h2
	data.ReqType =&reqType

	m2 := test.AssembleDataNomd5(uint16(protoId), data)
	conn.Write(m2)

	result := test.Read(conn).(*bbproto.ThRoom)
	fmt.Println("读取的结果:", result.GetHeader())		//测试服务器同意返回98989

}

func  bet(userId uint32,betType int32,amount int64){
	//获得连接
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide2 := int32(bbproto.EProtoId_THBET)
	fmt.Println("proto 得到的id ",ide2)
	reqData := &bbproto.THBet{}
	h2 := &bbproto.ProtoHeader{}
	h2.UserId = &userId

	reqData.Header = h2
	reqData.BetAmount = &amount
	reqData.BetType = &betType
	reqData.BetAmount = &amount

	m2 := test.AssembleData(uint16(ide2), reqData)
	conn.Write(m2)

	result := test.Read(conn).(*bbproto.THBet)
	fmt.Println("读取的结果:", result.GetHeader())		//测试服务器同意返回98989

}
