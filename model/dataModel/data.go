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
	time        time.Time      `bson:"time"`		//进入时间
}
func AtHome() []*User{
	info := []*User{}
	db.C(tableName.ADMIN_USER_ATHOME).FindAll(bson.M{},&info)
	return info
}


func AtHomeList(GameID string) []*User{
	info := []*User{}
	db.C(tableName.ADMIN_USER_ATHOME).FindAll(bson.M{
		"GameID" : GameID,
	},&info)
	return info
}



