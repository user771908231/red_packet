package autoLoginModel

import (
	"encoding/base64"
	"casino_common/common/log"
	"casino_common/utils/encodingUtils"
	"strings"
	"time"
	"strconv"
)

type Data struct {
	Time time.Time
	Id uint32
}

// base64编码
func EnBase64(enc *base64.Encoding,string []byte) string{

	encStr := enc.EncodeToString(string)
	return encStr
}
//  base64解码
func DeBase64(enc *base64.Encoding,string string) ([]byte,error) {
	data,err := enc.DecodeString(string)
	if err != nil {
		log.T("base64解密错误:%s",err)

	}
	return data,err
}

//
func DataDecoding(str string,key []byte) ([]byte,error) {
	//base64解码
	de,err := DeBase64(base64.StdEncoding,str)
	if err != nil {
		log.T("base64解密错误:%s",err)
		return de ,err
	}
	//AES解码
	decode_str, err := encodingUtils.AesDecrypt(de, key)
	if err != nil {
		log.T("AES解密错误:%s",err)

	}
	return decode_str, err
}

// 字符分割处理

func CharacterSplit(str string) *Data {
	arr := strings.Split(string(str), "-")
	if len(arr) == 2 {
		str0,_ := strconv.Atoi(arr[0])
		str1,_ := strconv.Atoi(arr[1])
		id := uint32(str1)
		time := time.Unix(int64(str0),0)
		res := Data{
			Time:time,
			Id:id,
		}
		return &res
	}
	return nil
}

//一个小时之内

func OneTimeWithin() time.Time {
	now := time.Now()
	h, _ := time.ParseDuration("1h")
	h1 := now.Add(1 * h)
	return h1
}

