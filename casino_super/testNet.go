package main

import (
	"gopkg.in/fatih/pool.v2"
	"fmt"
	"net"
	"casino_server/utils/test"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
)

func main() {
	factory := func() (net.Conn, error) {
		return net.Dial("tcp", "192.168.199.240:2801")
	}

	p, err := pool.NewChannelPool(5, 30, factory)
	fmt.Printf("err:%v", err)

	conn, err := p.Get()
	if err != nil {
		panic(err)
	}

	data := &ddproto.CommonReqGameLogin{
		UserId: proto.Uint32(8988989),
	}
	m := test.AssembleData(uint16(ddproto.HallEnumProtoId_HALL_PID_GAME_LOGIN), data)
	conn.Write(m)

	conn.Close()

	p.Close()
}
