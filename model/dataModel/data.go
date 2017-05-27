package dataModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
)

type User struct {
	Id          bson.ObjectId   `bson:"_id"`		//用户ID
	UserID      string   	`bson:"UserID"`		//用户ID
	UserName    string    	 `bson:"UserName"`   	 	//用户名字
	UserNick    string    	`bson:"UserNick"`		//用户昵称
	GameID      string    	`bson:"GameID"`		//游戏ID
	GameIn      string	`bson:"GameIn"`		//所在游戏
	HomeIn      string    	`bson:"HomeIn"`		//所在房间
	IP          string    	`bson:"IP"`			//进入IP
	Time        time.Time      `bson:"time"`		//进入时间
}
type Online struct {
	Id          bson.ObjectId   	`bson:"_id"`		//用户ID
	OnlineCount      int32   		`bson:"OnlineCount"`		//在线数
	Time        int64     	 `bson:"Time"`		//记录时间
}
func AtHome() []*User{
	info := []*User{}
	db.C(tableName.ADMIN_USER_ATHOME).FindAll(bson.M{},&info)
	return info
}

//在房间玩家列表
func AtHomeList(GameID string) []*User{
	info := []*User{}
	db.C(tableName.ADMIN_USER_ATHOME).FindAll(bson.M{
		"GameID" : GameID,
	},&info)
	return info
}

//在线统计列表--(小时)
func OnlineStaticList(Time_start int64,Time_end int64) []*Online{
	info := []*Online{}
	db.C(tableName.ADMIN_USER_ONLINEHOUR).FindAll(bson.M{"Time": bson.M{"$gte": Time_start,"$lte": Time_end}},&info)
	return info
}
//在线统计列表--(天)
func OnlineStaticDay(Time_start int64,Time_end int64) []*Online{
	info := []*Online{}
	db.C(tableName.ADMIN_USER_ONLINEDAY).FindAll(bson.M{"time": bson.M{"$gte": Time_start,"$lte": Time_end}},&info)
	return info
}





