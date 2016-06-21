package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbproto"
	"casino_server/utils/test"
)

func TestLogin(t *testing.T){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var id uint16 = 2
	var uId uint32 = 1000
	//var uu string ="2384028409320424"
	var h *bbproto.ProtoHeader = &bbproto.ProtoHeader{
		UserId:&uId,
	}
	var data bbproto.ReqAuthUser
	data.Header = h
	//data.Uuid = &uu
	m := utils.AssembleData(id,&data)
	conn.Write(m)
}

