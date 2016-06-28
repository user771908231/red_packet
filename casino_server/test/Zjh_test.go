package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"fmt"
	"casino_server/utils/test"
)

func TestZjhMain(t *testing.T) {
	zjhRoom()
}

/**
	PortoProcessor.Register(&bbproto.ZjhBet{})	//11	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhMsg{})	//12	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhQueryNoSeatUser{})	//13	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhReqSeat{})	//14	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhLottery{})	//15	进入扎金花的房间
 */

func zjhRoom() {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide := bbproto.EUnitProtoId_value[bbproto.EUnitProtoId_ZJHROOM.String()]
	fmt.Println("proto 得到的id ",ide)
	var userid uint32 = 10001
	data := &bbproto.ZjhRoom{}
	h := &bbproto.ProtoHeader{}

	h.UserId = &userid
	data.Header = h

	m := utils.AssembleData(uint16(ide), data)
	conn.Write(m)
	//result := utils.Read(conn).(*bbproto.ShuiguojiRes)
	//fmt.Println("读取的结果:", result.Result)
}
