package logDao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"time"
)

type ReqLog struct {
	Code      string `form:"code" binding:"Required;Include(casino)"` //验证码casino
	DeskId    string `form:"deskid" binding:"Required"` //桌子id
	UserId    string `form:"userid" binding:"Required"` //用户id
	Level     string `form:"level" binding:"Required"` //日志级别
	Data      string `form:"data" binding:"Required"` //日志内容
	createdAt int64 //创建时间
}

var TABLE string = "t_log"

func (l ReqLog) GetReadTime() string {
	
}

func FindLogsByKV(key string, v interface{}) []ReqLog {
	reqLogs := []ReqLog{}
	db.Query(func(d *mgo.Database) {
		d.C(TABLE).Find(bson.M{key: v}).All(&reqLogs)
	})
	if len(reqLogs) > 0 {
		return reqLogs
	} else {
		return nil
	}
}

func FindLogsByMap(m bson.M) []ReqLog {
	reqLogs := []ReqLog{}
	db.Query(func(d *mgo.Database) {
		d.C(TABLE).Find(m).All(&reqLogs)
	})
	if len(reqLogs) > 0 {
		return reqLogs
	} else {
		return nil
	}
}

//通过level 查找一个log
func FindLogsByLevel(level string) []ReqLog {
	return FindLogsByKV("level", level)
}

//通过desk id查找一个log
func FindLogsByDeskId(deskId string) []ReqLog {
	return FindLogsByKV("deskid", deskId)
}

//通过user id查找一个log
func FindLogsByUserId(userId string) []ReqLog {
	return FindLogsByKV("userid", userId)
}

//
func SaveLog2Mgo(reqLog ReqLog) error {
	reqLog.createdAt = time.Now().UnixNano()
	return db.InsertMgoData(TABLE, reqLog)
}
