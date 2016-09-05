package redisUtils

import (
	"casino_server/conf/casinoConf"
	"casino_server/utils/redis"
	"github.com/golang/protobuf/proto"
	"casino_server/common/log"
)

func init() {
	data.InitRedis()
}

//在需要的地方,需要自己关闭连接
func GetConn() data.Data {
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	return conn
}

/**
	获取redis中的一个对象,如果redis中没有值,那么返回nil
	如果有值,那么返回一个proto.message,需要取值的时候,需要强制转换
	eg.
		proto := GetObj(key,&bbproto.type{})
		if proto == nil{
			//没有找到值
		}else{
			result := proto.(*bbproto.type)
		}
 */

func GetObj(key string, p proto.Message) proto.Message {
	conn := GetConn()
	defer conn.Close()
	return conn.GetObjv2(key, p)
}

/**
	保存一个对象到redis
 */
func SetObj(key string, p proto.Message) {
	conn := GetConn()
	defer conn.Close()
	conn.SetObj(key, p)
}

func Del(key string){
	conn := GetConn()
	defer conn.Close()
	conn.Del(key)
}

func ZADD(key string, member string, score int64) {
	conn := GetConn()
	defer conn.Close()
	conn.ZAdd(key, member, score)
}

func ZREVRANK(key string, member string) int64 {
	conn := GetConn()
	defer conn.Close()
	values, err := conn.ZREVRANK(key, member)
	if values == nil || err != nil {
		return -1
	}
	return values.(int64)
}

//	LREM key count value 移除列表元素
func LREM(key string, value string) error {
	conn := GetConn()
	defer conn.Close()
	err := conn.LREM(key, value)
	if err != nil {
		log.E("redis lrem的时候出错 err[%v]", err)
	}
	return err
}
