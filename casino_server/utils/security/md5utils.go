package security

import (
	"crypto/md5"
	"fmt"
)

//byte[] SECRET_KEY = new byte[] { 0x93, 0x46, 0x78, 0x20 };
var  SECRET_KEY  = []byte{0x93, 0x46, 0x78, 0x20 }

func Md5(data []byte) []byte{
	fmt.Println("data:",data)
	fmt.Println("SECRET_KEY",SECRET_KEY)

	md5data := append(data,SECRET_KEY[0],SECRET_KEY[1],SECRET_KEY[2],SECRET_KEY[3])
	fmt.Println("md5data",md5data)

	h := md5.New()
	h.Write(md5data)
	//result := hex.EncodeToString(h.Sum(nil))
	resultByte := h.Sum(nil)
	fmt.Println("md5:",resultByte)

	var resultByte4 []byte
	resultByte4 = append(resultByte4,resultByte[4],resultByte[6],resultByte[8],resultByte[10])
	fmt.Println("resultByte4:",resultByte4)
	return resultByte4
}
