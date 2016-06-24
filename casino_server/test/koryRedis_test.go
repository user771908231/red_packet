package mongodb

import (
	"testing"
	"casino_server/utils/redis"
	"fmt"
)

func TestReids(t *testing.T){
	data.InitRedis()
	conn := data.Data{}
	conn.Open("test")
	conn.Set("t1", []byte("test_insert_t1"))
	str,_ :=conn.Get("t1")
	fmt.Println(str)

	//user := mode.User{
	//	Name:"哈哈哈哈",
	//}
	//conn.Set()


}
