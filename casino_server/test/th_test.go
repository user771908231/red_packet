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
	intoRoom()
	//bet()
}


/** 182.92.179.230
	测试进入房间
 */
func intoRoom(){
	//获得连接
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide2 := int32(bbproto.EProtoId_THROOM)
	fmt.Println("proto 得到的id ",ide2)
	var userid uint32 = 10010
	var reqType int32 = 1
	data2 := &bbproto.ThRoom{}
	h2 := &bbproto.ProtoHeader{}

	h2.UserId = &userid
	data2.Header = h2
	data2.ReqType =&reqType
	m2 := test.AssembleData(uint16(ide2), data2)
	conn.Write(m2)

	result := test.Read(conn).(*bbproto.ThRoom)
	fmt.Println("读取的结果:", result.GetHeader())		//测试服务器同意返回98989

}

func  bet(){
	//获得连接
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide2 := int32(bbproto.EProtoId_THBET)
	fmt.Println("proto 得到的id ",ide2)
	var userid uint32 = 10008
	var amount int32 = 999
	data2 := &bbproto.THBet{}
	h2 := &bbproto.ProtoHeader{}

	h2.UserId = &userid
	data2.Header = h2
	data2.BetAmount = &amount
	data2.BetType = &(room.TH_DESK_BET_TYPE_CALL)
	m2 := test.AssembleData(uint16(ide2), data2)
	conn.Write(m2)

	result := test.Read(conn).(*bbproto.THBet)
	fmt.Println("读取的结果:", result.GetHeader())		//测试服务器同意返回98989
}