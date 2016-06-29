package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/test"
)

func TestReqUsr(t *testing.T){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var id uint16 = 5
	var data bbproto.ReqAuthUser
	sendMsg := "sdfajfoqiweu8198230hjksjlfj"
	data.Uuid = &sendMsg
	m := test.AssembleData(id,&data)
	conn.Write(m)

}
