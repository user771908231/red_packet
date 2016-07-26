package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/test"
)

func TestBingfa(t *testing.T) {
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var s string = "ll"
	ide2 := int32(0)
	data2 := &bbproto.TestP1{}
	data2.Name2 = &s
	m2 := test.AssembleDataNomd5(uint16(ide2), data2)
	conn.Write(m2)

	for ; ;  {
		
	}
}
