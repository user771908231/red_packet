package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbproto"
	"casino_server/utils/test"
)

func TestShuiGuoJi(t *testing.T){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var id uint16 = 7
	var data bbproto.Shuiguoji

	var nApple int32 = 1						//苹果的数量
	data.ScoresApple = &nApple

	m := utils.AssembleData(id,&data)
	conn.Write(m)
}
