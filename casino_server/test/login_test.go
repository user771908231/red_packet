package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/test"
)

func TestLogin(t *testing.T){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var id uint16 = 2
	var uId uint32 = 10001
	var h *bbproto.ProtoHeader = &bbproto.ProtoHeader{
		UserId:&uId,
	}
	var data bbproto.ReqAuthUser
	data.Header = h
	m := test.AssembleData(id,&data)
	conn.Write(m)
}

