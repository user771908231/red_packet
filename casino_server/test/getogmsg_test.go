package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/test"
)

//次协议用来返回房间的信息
func TestOgg(t *testing.T){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var s string = "ogg"
	ide2 := int32(0)
	data2 := &bbproto.TestP1{}
	data2.Name2 = &s
	m2 := test.AssembleDataNomd5(uint16(ide2), data2)
	conn.Write(m2)

	for ; ;  {

	}

}
