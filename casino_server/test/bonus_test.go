package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"casino_server/msg/bbprotoFuncs"
	"fmt"
	"casino_server/utils/test"
)

func TestBonus(t *testing.T){
	testLoginTurntable()
}

//测试转盘奖励
func testLoginTurntable(){

	//得到连接
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//请求登陆奖励
	ide := bbproto.EProtoId_value[bbproto.EProtoId_LOGINTURNTABLEBONUS.String()]
	var userId uint32 = 10005
	res := &bbproto.LoginTurntableBonus{}
	res.Header = protoUtils.GetSuccHeaderwithUserid(&userId)

	fmt.Println("开始请求")
	conn.Write(test.AssembleData(uint16(ide), res))
	fmt.Println("请求结束")

	result := test.Read(conn).(*bbproto.LoginTurntableBonus)
	fmt.Println("读取到的结果-code:",result.GetHeader().GetCode())
	fmt.Println("读取到的结果-msg:",result.GetHeader().GetError())

}
