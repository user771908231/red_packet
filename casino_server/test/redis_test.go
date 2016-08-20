package test

import (
	"testing"
	"casino_server/msg/bbprotogo"
	"fmt"
	"casino_server/utils/redisUtils"
)

func TestRedis(t *testing.T) {
	fmt.Println("begin")

	ret := redisUtils.GetObj("agent_session_10146", &bbproto.ThServerUserSession{})
	fmt.Println("ret", ret)
	p := ret.(*bbproto.ThServerUserSession)
	fmt.Println("p", p)

}
