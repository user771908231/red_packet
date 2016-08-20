package mongodb

import (
	"testing"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"time"
	"casino_server/mode"
)

func TestRedis(t *testing.T)  {
	f1(t)
}

func f1(t *testing.T){
	redis, err := cache.NewCache("redis", `{"conn":"192.168.199.120:6379", "key":"casinoRedis"}`)
	if err != nil {
		t.Error(err)
	}

	user:=mode.T_user{
	}

	v := "a,b,c,d"
	err2 := redis.Put("key",v,10*time.Second)
	if err2 !=nil{
		t.Error(err2)
	}
	redis.Put("key",user,10*time.Second)
	result := redis.Get("key")
	t.Log("结果",string(result.([]byte)))
}


