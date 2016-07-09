
package mode

import "gopkg.in/mgo.v2/bson"



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

	Coin     int32    //金币
	Scores   int32    //积分

	Viplevel string   //vip的等级
	Level    string   //等级
	Rank     string   //称号

	Gift     []string //礼物的记录
}



