package autoLogin

import (
	"testing"
	"casino_common/utils/encodingUtils"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"casino_redpack/model/autoLoginModel"
	"casino_common/common/userService"
)

func TestAcceptData(t *testing.T){
	str := "SBZCxV42UKSn6e7DflvBlM7 mgtiwRjUyTC/etY1zuc="
	str = strings.Replace(str, " ", "+", -1)
	t.Log(str)
}


func TestAesEncrypt(t *testing.T) {
	now := time.Now()
	one := fmt.Sprintf("%d-12345",now.Unix())
	plat_text := []byte(one)
	key := []byte("123asdssssssssss")

	t.Log("plat_text:[%v]", plat_text)

	encode_str, err := encodingUtils.AesEncrypt(plat_text, key)
	if err != nil {
		t.Log("AES加密错误！")
	}
	t.Log("AES加密",encode_str)
	en := EnBase64(base64.StdEncoding,encode_str)
	t.Log("base64加密",en)
	de := DeBase64(base64.StdEncoding,en)
	t.Log("base64解密",de)
	decode_str, err := encodingUtils.AesDecrypt(de, key)
	t.Log("decode_str:[%v] err:[%v]", string(decode_str), err)
	arr := strings.Split(string(decode_str), "-")
	t.Log("分给",arr)
	t.Log("长度",len(arr))
	data := autoLoginModel.CharacterSplit(string(decode_str))
	t.Log("时间：",data.Time)

	subM := now.Sub(data.Time)
	t.Log("时间差：",subM.Hours())
	if subM.Minutes() <  float64(10) {
		t.Log("时间差小于10分钟")
		user_info := userService.GetUserById(data.Id)
		if &user_info != nil {
			t.Log("获取到的用户",user_info.Id)
		}
		t.Log("获取用户失败")
	}else {
		t.Log("时间差大于10分钟")
	}


}
// base64编码
func EnBase64(enc *base64.Encoding,string []byte) string{

	encStr := enc.EncodeToString(string)
	return encStr
}
//  base64解码
func DeBase64(enc *base64.Encoding,string string) []byte {
	data,err := enc.DecodeString(string)
	if err != nil {
		fmt.Println("base64解密错误:%s",err)
		return []byte("")
	}
	return data
}

