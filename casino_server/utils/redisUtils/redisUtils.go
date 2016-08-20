package redisUtils

import (
	"casino_server/conf/casinoConf"
	"casino_server/utils/redis"
	"github.com/golang/protobuf/proto"
)


func init(){
	data.InitRedis()
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

func GetObj(key string,p proto.Message) proto.Message{
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()
	return conn.GetObjv2(key, p)
}