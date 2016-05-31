package main

import (
)
import (
	"net"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"fmt"
	"casino_server/msg/bbproto"
	"casino_server/msg"
	"github.com/name5566/leaf/log"
)

func main() {

	conn, err := net.Dial("tcp", "192.168.199.120:3563")
	if err != nil {
		panic(err)
	}

	var id  = []byte{0,0}
	var data bbproto.N
	var n string = "a"
	data.Name = &n
	data3 ,err :=  proto.Marshal(&data)
	m2 := make([]byte, 4+len(data3))

	//// 默认使用大端序
	binary.BigEndian.PutUint16(m2, uint16(2+len(data3)))
	copy(m2[2:4], id)
	copy(m2[4:], data3)
	fmt.Println("发送的m2:",m2)
	conn.Write(m2)

	var res [250]byte
	count,err := conn.Read(res[0:])
	if err != nil {
		fmt.Println("err != nil")
	}

	log.Debug("读取到的 res %v",res)
	msg2, err := msg.PortoProcessor.Unmarshal(res[2:count])
	if err != nil {
	}

	m5 :=  msg2.(*bbproto.N)
	fmt.Println("m5:",*m5.Name)

}