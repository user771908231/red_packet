package logDao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"time"
	"casino_common/utils/timeUtils"
	"errors"
)

type LogData struct {
	DeskId    string `json:"deskid" binding:"Required"` //桌子id
	UserId    string `json:"userid" binding:"Required"` //用户id
	Level     string `json:"level" binding:"Required"`  //日志级别
	Data      string `json:"data" binding:"Required"`   //日志内容
	CreatedAt string                                    //创建时间
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
		d.C(TABLE).Find(bson.M{key: v}).Sort("createdAt").All(&logData)
	})
	if len(logData) > 0 {
		return logData
	} else {
		return nil
	}
}

func FindLogsByMap(m bson.M, skip, limit int) []LogData {
	logData := []LogData{}
	db.Query(func(d *mgo.Database) {
		//d.C(TABLE).Find(m).Sort("createdAt").All(&logData)
		d.C(TABLE).Find(m).Sort("createdAt").Skip(skip).Limit(limit).All(&logData)
	})
	if len(logData) > 0 {
		return logData
	} else {
		return nil
	}
}

func FindLogsByMapCount(m bson.M) (int, error) {
	err := errors.New("")
	c := 0
	db.Query(func(d *mgo.Database) {
		c, err = d.C(TABLE).Find(m).Count()
	})
	if err != nil {
		return -1, err
	}
	return c, nil
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
	logData.CreatedAt = timeUtils.Format(time.Now())
	return db.InsertMgoData(TABLE, logData)
}

//
func DeleteAllLogs2Mgo() error {
	db.Query(func(d *mgo.Database) {
		d.C(TABLE).RemoveAll(bson.M{})
	})
	return nil
}
