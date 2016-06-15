package security

import (
	"crypto/md5"
	"fmt"
	"github.com/name5566/leaf/log"
	"errors"
)

//byte[] SECRET_KEY = new byte[] { 0x93, 0x46, 0x78, 0x20 };
var  SECRET_KEY  = []byte{0x93, 0x46, 0x78, 0x20 }

func Md5IdAndData(id,data []byte) []byte{
	fmt.Println("id:",id)
	fmt.Println("data:",data)
	fmt.Println("SECRET_KEY",SECRET_KEY)
	m2 := make([]byte, 2+len(data))
	copy(m2[0:2],id)
	copy(m2[2:],data)
	return Md5(m2)
}

func Md5(data []byte) []byte{
	log.Debug("需要校验的data%v",data)
	log.Debug("SECRET_KEY%v",SECRET_KEY)

	md5data := append(data,SECRET_KEY[0],SECRET_KEY[1],SECRET_KEY[2],SECRET_KEY[3])
	h := md5.New()
	h.Write(md5data)
	resultByte := h.Sum(nil)
	log.Debug("校验计算出来的md5:",resultByte)

	var resultByte4 []byte
	resultByte4 = append(resultByte4,resultByte[4],resultByte[6],resultByte[8],resultByte[10])
	log.Debug("校验的结果sign:",resultByte4)
	return resultByte4
}


/**
验证tcp读取的数据是否有问题
 */
func CheckTcpData(data []byte)([]byte,error){
	md5data := make([]byte,len(data)-4)
	copy(md5data,data[:len(data)-4])
	sign := data[len(data)-4:]
	cmd5 := Md5(md5data)

	log.Debug("校验数据全%v",data)
	log.Debug("校验数据信息部分%v",md5data)
	log.Debug("校验数据全sign %v",sign)
	log.Debug("cmd5 %v",cmd5)

	result := true
	for i:= 0;i<len(sign) ;i++ {
		if sign[i] != cmd5[i] {
			log.Debug("sign[%v],cmd5[%v]",sign[i],cmd5[i])
			result = false
			break
		}
	}

	if !result {
		log.Error("校验失败的数据%v",data)
		return nil,errors.New("数据校验失败")
	}

	return  md5data,nil
}