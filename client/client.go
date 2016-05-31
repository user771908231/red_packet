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
)

func main() {



	conn, err := net.Dial("tcp", "192.168.199.120:3563")
	if err != nil {
		panic(err)
	}

	var id  = []byte{0,0}
	var data bbproto.N
	var n string = "hp"
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
	fmt.Println("res1:%v",res)

	count,err := conn.Read(res[0:])
	if err != nil {
		fmt.Println("err != nil")
	}

	msg2, err := msg.PortoProcessor.Unmarshal(res[2:count])
	if err != nil {
	}

	m5 :=  msg2.(*bbproto.N)
	fmt.Println("m5:",*m5.Name)


	//fmt.Println("发送的m2:",m2)
	//conn.Write(m2)


	//
	//var res2 [250]byte
	//fmt.Println("res1:%v",res2)
	//
	//count2,err2 := conn.Read(res[0:])
	//if err2 != nil {
	//	fmt.Println("err != nil")
	//}
	//
	//msg22, err := msg.PortoProcessor.Unmarshal(res[2:count2])
	//if err != nil {
	//}
	//
	//m52 :=  msg22.(*bbproto.N)
	//fmt.Println("m5:",*m52.Name)


//#############################json
	// Hello 消息（JSON 格式）
	// 对应游戏服务器 Hello 消息结构体
	//data := []byte(`{
	//	"Hello": {
	//	    "Name": "leaf"
	//	}
	//    }`)
	//
	//// len + data
	//m := make([]byte, 2+len(data))
	//
	//// 默认使用大端序
	//binary.BigEndian.PutUint16(m, uint16(len(data)))
	//
	//copy(m[2:], data)
	//
	//// 发送消息
	//conn.Write(m)
	//
	//
	//
	//var res [250]byte
	//fmt.Println("res1:%v",res)
	//count,err := conn.Read(res[0:])
	//if err != nil {
	//	fmt.Println("err != nil")
	//}
	//fmt.Println("count:%v",count)
	//fmt.Println("res2:%v",res)
	//
	//msg2, err := msg.Processor.Unmarshal(res[2:count])
	//if err != nil {
	//
	//}
	//
	//m5 :=  msg2.(*msg.Hello)
	//
	//fmt.Println("m5:",m5.Name)
	//fmt.Println("msg2:",msg2)

}