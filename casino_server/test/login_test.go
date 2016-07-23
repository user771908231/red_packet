package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/test"
	"fmt"
	"casino_server/msg/bbprotoFuncs"
)

func TestLogin(t *testing.T) {
	//login2()
	login1(10007)
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
	var id uint16 = 2

	data := &bbproto.ReqAuthUser{}
	data.Header = protoUtils.GetSuccHeaderwithUserid(&userId)
	m := test.AssembleData(id, data)
	conn.Write(m)

	result := test.Read(conn).(*bbproto.ReqAuthUser)
	fmt.Println("读取的结果:", result)

}

//游客登陆

func login2() {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var id uint16 = 2
	var h *bbproto.ProtoHeader = &bbproto.ProtoHeader{}
	var data bbproto.ReqAuthUser
	var uuidStr = "8029409jowejfiosjljfl"
	data.Header = h
	data.Uuid = &uuidStr
	m := test.AssembleData(id, &data)
	conn.Write(m)
}
