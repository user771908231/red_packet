package mode

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

/**
	转盘奖励的表
 */
type T_bonus_turntable  struct{
	Mid		bson.ObjectId		`json:"mid" bson:"_id"`
	Id		uint32		//id
	Amount 		int32		//金额
	Time		time.Time	//时间
	UserId		uint32		//用户id
}