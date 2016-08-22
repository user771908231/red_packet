package mode

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


//奖励的类型放在intCons里面
type T_rewards struct {
	Mid     bson.ObjectId                `json:"mid" bson:"_id"`
	rtype   int32     //奖励类型
	name    string    //
	time    time.Time //奖励时间
	num     int32     //奖励的额度
	userMid bson.ObjectId
}
