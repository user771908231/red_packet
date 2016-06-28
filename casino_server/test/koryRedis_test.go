package mongodb

import (
	"testing"
	"casino_server/utils/redis"
	"fmt"
	"casino_server/utils"
)

func TestReids(t *testing.T){
	fmt.Println("开始测试redis")
	data.InitRedis()
	conn := data.Data{}
	conn.Open("test")
	defer conn.Close()

	conn.Set("t1", []byte("test_insert_t1"))
	str,_ :=conn.Get("t1")

	fmt.Println(str)

	var s2  string = "我是s2"
	var i2  int32  = 89
	str2,_ := utils.Int2String(i2)
	fmt.Println(s2+str2)

	//user := mode.User{
	//	Name:"哈哈哈哈",
	//}
	//conn.Set()


}
