package logDao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"time"
)

type LogData struct {
	DeskId    string `json:"deskid" binding:"Required"` //桌子id
	UserId    string `json:"userid" binding:"Required"` //用户id
	Level     string `json:"level" binding:"Required"` //日志级别
	Data      string `json:"data" binding:"Required"` //日志内容
	createdAt int64 //创建时间
}

type ReqLog struct {
	Code      string `json:"code" binding:"Required;Include(casino)"` //验证码casino
	Logs      []LogData `json:"logs" binding:"Required"` //数组
}

var TABLE string = "t_log"

//func (l ReqLog) GetReadTime() string {
//
//}

func FindLogsByKV(key string, v interface{}) []LogData {
	logData := []LogData{}
	db.Query(func(d *mgo.Database) {
		d.C(TABLE).Find(bson.M{key: v}).All(&logData)
	})
	if len(logData) > 0 {
		return logData
	} else {
		return nil
	}
}

func FindLogsByMap(m bson.M) []LogData {
	logData := []LogData{}
	db.Query(func(d *mgo.Database) {
		d.C(TABLE).Find(m).All(&logData)
	})
	if len(logData) > 0 {
		return logData
	} else {
		return nil
	}
}

//通过level 查找一个log
func FindLogsByLevel(level string) []LogData {
	return FindLogsByKV("level", level)
}

//通过desk id查找一个log
func FindLogsByDeskId(deskId string) []LogData {
	return FindLogsByKV("deskid", deskId)
}

//通过user id查找一个log
func FindLogsByUserId(userId string) []LogData {
	return FindLogsByKV("userid", userId)
}

//
func SaveLog2Mgo(logData LogData) error {
	logData.createdAt = time.Now().UnixNano()
	return db.InsertMgoData(TABLE, logData)
}

//
func DeleteAllLogs2Mgo() error {
	db.Query(func(d *mgo.Database) {
		d.C(TABLE).RemoveAll(bson.M{})
	})
	return nil
}
