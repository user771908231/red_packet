package mongodb

import (
	"testing"
	"casino_server/utils/redis"
	"fmt"
	"casino_server/utils/numUtils"
	"casino_server/msg/bbprotogo"
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
	str2,_ := numUtils.Int2String(i2)
	fmt.Println(s2+str2)

	var hahaha string = "哈哈哈哈哈你"
	user := &bbproto.User{
		Name:&hahaha,
	}

	err := conn.SetObj("testObj",user)
	if err != nil {
		panic(err)
	}

	user2 := &bbproto.User{}

	err = conn.GetObj("testObj",user2)

	fmt.Println("user.name",user.GetName())
	fmt.Println("user2.name",user2.GetName())

}
