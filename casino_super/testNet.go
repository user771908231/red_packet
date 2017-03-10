package main

import (
	"gopkg.in/fatih/pool.v2"
	"net"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"crypto/md5"
	"testing"
	"fmt"
	"casino_common/proto/funcsInit"
)

var  SECRET_KEY  = []byte{0x93, 0x46, 0x78, 0x20 }

func Md5(data []byte) []byte{
	//log.T("需要校验的data%v",data)
	md5data := append(data,SECRET_KEY[0],SECRET_KEY[1],SECRET_KEY[2],SECRET_KEY[3])
	h := md5.New()
	h.Write(md5data)
	resultByte := h.Sum(nil)
	//log.T("校验计算出来的完整md5:%v",resultByte)

	var resultByte4 []byte
	resultByte4 = append(resultByte4,resultByte[4],resultByte[6],resultByte[8],resultByte[10])
	//log.T("校验的结果sign:%v",resultByte4)
	return resultByte4
}

func Md5IdAndData(id,data []byte) []byte{
	//fmt.Println("id:",id)
	//fmt.Println("data:",data)
	//fmt.Println("SECRET_KEY",SECRET_KEY)
	m2 := make([]byte, 2+len(data))
	copy(m2[0:2],id)
	copy(m2[2:],data)
	return Md5(m2)
}

//
func AssembleData(idv ddproto.HallEnumProtoId, data proto.Message) []byte {
	//fmt.Println("需要转化的id",idv)
	//1,把id转化成 []byte
	id := make([]byte, 2)
	idv2 := uint16(idv)
	binary.BigEndian.PutUint16(id, idv2)
	//	//2,data 转化成[]byte
	data2, err := proto.Marshal(data)
	if err != nil {
		panic(err)
	}

	//3计算md5
	md5byte := Md5IdAndData(id, data2)
	leng := len(data2)
	m2 := make([]byte, 4 + leng + 4)
	binary.BigEndian.PutUint16(m2, uint16(2 + leng + 4))                //// 默认使用大端序
	copy(m2[2:4], id)
	copy(m2[4:leng + 4], data2)
	copy(m2[leng + 4:], md5byte)

	//4,返回值
	//fmt.Println("发送的m2:", m2)
	return m2
}

func TestRun(t *testing.T) {
	factory := func() (net.Conn, error) {
		return net.Dial("tcp", "192.168.199.155:2801")
	}

	p, err := pool.NewChannelPool(5, 30, factory)
	//fmt.Printf("err:%v", err)

	conn, err := p.Get()
	if err != nil {
		panic(err)
	}

	datas := AssembleData(ddproto.HallEnumProtoId_HALL_PID_USER_DATA_ACK, &ddproto.HallAckUserData{
		Header:commonNewPorot.NewHeader(),
		User:&ddproto.HallUserData{
			UserName:proto.String("james"),
		},
	})
	msg := &ddproto.Push{
		Id: proto.Uint32(13633),
		Data: datas,
	}
	m := AssembleData(ddproto.HallEnumProtoId_HALL_PID_PUSH_REQ, msg)
	_,err = conn.Write(m)
	fmt.Println(err)

	conn.Close()

	p.Close()
}

//推送
func push(user_id uint32, msg interface{}) {

}
