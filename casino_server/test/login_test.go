package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/test"
)

func TestLogin(t *testing.T) {
	login2()
	//login1(10007)
	//login1(10008)
	//login1(10009)
	//login1(10010)
	//login1(10011)
}


//指定id登陆--
func login1(userId uint32) {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var id uint16 = uint16(bbproto.EProtoId_REQQUICKCONN)

	data := &bbproto.REQQuickConn{}
	data.UserId = &userId
	m := test.AssembleDataNomd5(id, data)
	conn.Write(m)
	_ = test.Read(conn).(*bbproto.ACKQuickConn)
}

//游客登陆

func login2() {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var id uint16 = uint16(bbproto.EProtoId_REQQUICKCONN)
	data := &bbproto.REQQuickConn{}
	m := test.AssembleDataNomd5(id, data)
	conn.Write(m)
	_ = test.Read(conn).(*bbproto.ACKQuickConn)
}
