package test

import (
	"testing"
	"casino_server/msg/bbprotogo"
	"fmt"
	"casino_server/utils/redisUtils"
)

func TestRedis(t *testing.T) {
	testZset()
}

func testRediUtils() {
	fmt.Println("begin")

	ret := redisUtils.GetObj("agent_session_10146", &bbproto.ThServerUserSession{})
	fmt.Println("ret", ret)
	p := ret.(*bbproto.ThServerUserSession)
	fmt.Println("p", p)
}

func testZset() {
	key := "kkk"
	redisUtils.ZADD(key,"a",2)
	redisUtils.ZADD(key,"b",2)
	redisUtils.ZADD(key,"c",2)
	redisUtils.ZADD(key,"d",9)
	redisUtils.ZADD(key,"e",7)


	//返回a的排名
	fmt.Println("a排名:",redisUtils.ZREVRANK(key,"a"))
	fmt.Println("b排名:",redisUtils.ZREVRANK(key,"b"))
	fmt.Println("c排名:",redisUtils.ZREVRANK(key,"c"))
	fmt.Println("d排名:",redisUtils.ZREVRANK(key,"d"))
	fmt.Println("e排名:",redisUtils.ZREVRANK(key,"e"))
	fmt.Println("e排名:",redisUtils.ZREVRANK(key,"o"))



}