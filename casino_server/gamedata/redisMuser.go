package gamedata

import (
	"casino_server/utils/redis"
	"casino_server/conf/casinoConf"
)


/**
	存储在redis中的userData
*/


//存放User相关键值的信息
const(
	USERID  		=	"R_KEY_RUSER_USER_ID"
	USERNAME		=	"R_KEY_RUSER_USER_NAME"
)

type Ruser struct {
	userId 		int32
	userName 	string
}


/**
从redis中读取user信息
 */
func GetRuserById(id int32) (*Ruser,error){
	//1 连接redis
	conn := data.Data{}
	conn.Open(casinoConf.REDIS_DB_NAME)
	defer conn.Close()

	//2 得到user key
	return nil,nil
}


