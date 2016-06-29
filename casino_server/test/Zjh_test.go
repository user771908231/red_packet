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
	//zjhMsg()
	//zjhQueryNoSeatUser()
	//zjhReqSeat()
	//zjhZjhLottery()
	//zjhBet()
}


func zjhRoom() {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide := bbproto.EUnitProtoId_value[bbproto.EUnitProtoId_ZJHROOM.String()]
	fmt.Println("proto 得到的id ",ide)
	var userid uint32 = 10001
	var reqType int32 = 0
	data := &bbproto.ZjhRoom{}
	h := &bbproto.ProtoHeader{}

	h.UserId = &userid
	data.Header = h
	data.ReqType =&reqType

	m := test.AssembleData(uint16(ide), data)
	conn.Write(m)
	result := test.Read(conn).(*bbproto.ZjhRoom)
	fmt.Println("读取的结果:", result.GetBanker().GetName())
	fmt.Println("读取的结果:", result.GetBanker().GetBalance())
	fmt.Println("读取的结果:", result.GetMe().GetName())
	fmt.Println("读取的结果:", result.GetMe().GetBalance())
	fmt.Println("读取的结果:", result.GetJackpot())
}


func zjhMsg() {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide := bbproto.EUnitProtoId_value[bbproto.EUnitProtoId_ZJHMSG.String()]
	fmt.Println("proto 得到的id ",ide)
	var userid uint32 = 10001
	data := &bbproto.ZjhMsg{}
	h := &bbproto.ProtoHeader{}

	h.UserId = &userid
	data.Header = h

	m := test.AssembleData(uint16(ide), data)
	conn.Write(m)
	//result := utils.Read(conn).(*bbproto.ShuiguojiRes)
	//fmt.Println("读取的结果:", result.Result)
}


func zjhQueryNoSeatUser() {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide := bbproto.EUnitProtoId_value[bbproto.EUnitProtoId_ZJHQUERYNOSEATUSER.String()]
	fmt.Println("proto 得到的id ",ide)
	var userid uint32 = 10001
	data := &bbproto.ZjhQueryNoSeatUser{}
	h := &bbproto.ProtoHeader{}

	h.UserId = &userid
	data.Header = h

	m := test.AssembleData(uint16(ide), data)
	conn.Write(m)
	//result := utils.Read(conn).(*bbproto.ShuiguojiRes)
	//fmt.Println("读取的结果:", result.Result)
}


//
func zjhReqSeat() {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide := bbproto.EUnitProtoId_value[bbproto.EUnitProtoId_ZJHREQSEAT.String()]
	fmt.Println("proto 得到的id ",ide)
	var userid uint32 = 10001
	data := &bbproto.ZjhReqSeat{}
	h := &bbproto.ProtoHeader{}

	h.UserId = &userid
	data.Header = h

	m := test.AssembleData(uint16(ide), data)
	conn.Write(m)
	//result := utils.Read(conn).(*bbproto.ShuiguojiRes)
	//fmt.Println("读取的结果:", result.Result)
}

func zjhZjhLottery(){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide := bbproto.EUnitProtoId_value[bbproto.EUnitProtoId_ZJHLOTTERY.String()]
	fmt.Println("proto 得到的id ",ide)
	var userid uint32 = 10001
	data := &bbproto.ZjhLottery{}
	h := &bbproto.ProtoHeader{}

	h.UserId = &userid
	data.Header = h

	m := test.AssembleData(uint16(ide), data)
	conn.Write(m)
	//result := utils.Read(conn).(*bbproto.ShuiguojiRes)
	//fmt.Println("读取的结果:", result.Result)
}


func zjhBet(){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ide := bbproto.EUnitProtoId_value[bbproto.EUnitProtoId_ZJHBET.String()]
	fmt.Println("proto 得到的id ",ide)
	var userid uint32 = 10001
	data := &bbproto.ZjhBet{}
	h := &bbproto.ProtoHeader{}

	h.UserId = &userid
	data.Header = h

	m := test.AssembleData(uint16(ide), data)
	conn.Write(m)

}