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

const url  = "192.168.199.120:3563"
const TCP = "tcp"

func main() {
	//testN()

	//测试注册会员
	//testReqAuthUser()
	testReqAuthUserWithmd5()
}


func testN(){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}

	var id  = []byte{0,1}
	var data bbproto.Reg
	var n string = "aaappp"
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

	m5 :=  msg2.(*bbproto.Reg)
	fmt.Println("m5:",*m5.Name)
}


func testReqAuthUser(){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}

	var id  = []byte{0,2}
	var data bbproto.ReqAuthUser
	var header bbproto.ProtoHeader
	var v int32
	v = 119
	uid := "sufiuowiurw9er0wuo"

	var userId uint32
	userId = 989

	data.AppVersion = &v

	header.UserId = &userId
	data.Header = &header
	data.Uuid = &uid
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
	m5 :=  msg2.(*bbproto.ReqAuthUser)
	fmt.Println("m5.error:",*m5.Header.Error)
	fmt.Println("m5.UserId:",*m5.Header.UserId)
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

	var userId uint32
	userId = 989

	data.AppVersion = &v

	header.UserId = &userId
	data.Header = &header
	data.Uuid = &uid
	data3 ,err :=  proto.Marshal(&data)

	//根据data3 计算md5
	md5string := security.Md5(data3)

	len :=  len(data3)
	fmt.Println("msg.len:",len)
	m2 := make([]byte, 4+len+4)

	//// 默认使用大端序
	binary.BigEndian.PutUint16(m2, uint16(2+len+4))
	copy(m2[2:4], id)
	copy(m2[4:len], data3)
	copy(m2[len+4:],md5string)

	fmt.Println("发送的m2:",m2)
	conn.Write(m2)


	//var res [250]byte
	//count,err := conn.Read(res[0:])
	//if err != nil {
	//	fmt.Println("err != nil")
	//}
	//log.Debug("读取到的 res %v",res)
	//msg2, err := msg.PortoProcessor.Unmarshal(res[2:count])
	//if err != nil {
	//}
	//m5 :=  msg2.(*bbproto.ReqAuthUser)
	//fmt.Println("m5.error:",*m5.Header.Error)
	//fmt.Println("m5.UserId:",*m5.Header.UserId)
}
