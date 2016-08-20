package mode

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


//锦标赛游戏的状态
var T_CS_TH_RECORD_STATUS_0 int32  = 0
var T_CS_TH_RECORD_STATUS_1 int32  = 1
var T_CS_TH_RECORD_STATUS_2 int32  = 2



/**
	1,游戏开始的时候就应该保存数据,目前默认的状态是进行中
	2,竞标赛结束的时候,修改状态为已经完成...
 */

//一局德州游戏是一条数据
type T_cs_th_record struct {
	Mid     	 bson.ObjectId		`json:"mid" bson:"_id"`
	Id       	int32   	//id
	BeginTime	time.Time
	EndTime 	time.Time
	Status		int32			//游戏的状态,未开始,进行中,已经完成
}

func (t *T_cs_th_record) GetMid() bson.ObjectId{
	return t.Mid
}
