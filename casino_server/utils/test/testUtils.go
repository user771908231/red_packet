package utils

import (
	"github.com/golang/protobuf/proto"
	"casino_server/utils/security"
	"encoding/binary"
	"fmt"
	"net"
	"casino_server/msg"
)


//
func AssembleData(idv uint16, data proto.Message) []byte {
	//1,把id转化成 []byte
	id := make([]byte, 2)
	binary.BigEndian.PutUint16(id, idv)
	//	//2,data 转化成[]byte
	data2, err := proto.Marshal(data)
	if err != nil {
		panic(err)
	}

	//3计算md5
	md5byte := security.Md5IdAndData(id, data2)
	len := len(data2)
	m2 := make([]byte, 4 + len + 4)
	binary.BigEndian.PutUint16(m2, uint16(2 + len + 4))                //// 默认使用大端序
	copy(m2[2:4], id)
	copy(m2[4:len + 4], data2)
	copy(m2[len + 4:], md5byte)

	//4,返回值
	fmt.Println("发送的m2:", m2)
	return m2
}

func Read(conn net.Conn, m proto.Message) proto.Message{
	fmt.Println("开始读取信息:")
	lenBuf := make([]byte, 2)
	conn.Read(lenBuf)

	msgLen := uint32(binary.BigEndian.Uint16(lenBuf))
	msgData := make([]byte, msgLen)
	conn.Read(msgData[0:])
	temp, _ := msg.PortoProcessor.Unmarshal(msgData)
	m = temp.(proto.Message)
	fmt.Println("读取信息完毕:",m)
	return m
}
