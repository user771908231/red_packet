package mode

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//一局德州游戏是一条数据
type T_cs_th_record struct {
	Mid      bson.ObjectId		`json:"mid" bson:"_id"`
	Id       int32   	//id
	BeginTime	time.Time
	EndTime 	time.Time
}