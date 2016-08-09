
package mode

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


/**

	warn 需要注意的事情:
	model中比较ObjectId
 */

type T_user struct {
	Mid      bson.ObjectId		`json:"mid" bson:"_id"`
	Id       uint32   //id
	Name     string
	NickName string   //昵称
	Sex      string
	Icon     string
	Desc     string   //个人描述
	Location string   //地区
	Password string
	Mobile   string

	Coin     int64    //金币
	Diamond	 int64		//钻石
	Scores   int32    //积分

	Viplevel string   //vip的等级
	Level    string   //等级
	Rank     string   //称号

	Gift     []string //礼物的记录
	TimeLastSign	time.Time	//最后一次签到时间
	ContinueSignCount int32		//连续签到的次数
}



