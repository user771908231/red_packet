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
	"casino_server/utils/security"
)

const url  = "192.168.199.111:3563"
const TCP = "tcp"

func main() {
	//testN()
	testReqAuthUserWithmd5()
}

func testReqAuthUserWithmd5(){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}

	var id  = []byte{0,2}
	var data bbproto.ReqAuthUser
	var header bbproto.ProtoHeader
	var v int32
	v = 119
	uid := "a"

	var userId uint32 = 989
	header.UserId = &userId
	data.AppVersion = &v
	data.Header = &header
	data.Uuid = &uid

	data3 ,err :=  proto.Marshal(&data)
	fmt.Println("data3:",data3)
	//根据data3 计算md5
	md5byte := security.Md5IdAndData(id,data3)

	len :=  len(data3)
	fmt.Println("msg.len:",len)
	m2 := make([]byte, 4+len+4)

	//// 默认使用大端序
	binary.BigEndian.PutUint16(m2, uint16(2+len+4))
	copy(m2[2:4], id)
	copy(m2[4:len+4], data3)
	copy(m2[len+4:], md5byte)

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
	m5 :=  msg2.(*bbproto.ReqAuthUser)
	fmt.Println("m5.UserId:",m5)

	for ; ;  {

	}
}
