package logDao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/utils/timeUtils"
	"time"
	"errors"
	"casino_super/conf/config"
)

type LogData struct {
	id        int32                                     //自增键
	DeskId    string `json:"deskid" binding:"Required"` //桌子id
	UserId    string `json:"userid" binding:"Required"` //用户id
	Level     string `json:"level" binding:"Required"`  //日志级别
	Data      string `json:"data" binding:"Required"`   //日志内容
	CreatedAt string                                    //创建时间
}

type ReqLog struct {
	Code string `json:"code" binding:"Required;Include(casino)"` //验证码casino
	Logs []LogData `json:"logs" binding:"Required"`              //数组
}

//func (l ReqLog) GetReadTime() string {
//
//}

//按记录查询
func FindLogsByKV(key string, v interface{}) []LogData {
	logData := []LogData{}
	db.Query(func(d *mgo.Database) {
		d.C(config.DBT_SUPER_LOGS).Find(bson.M{key: v}).Sort("id").All(&logData)
	})
	if len(logData) > 0 {
		return logData
	} else {
		return nil
	}
}

//查询条数
func FindLogsByMapCount(m bson.M) (int, error) {
	err := errors.New("")
	c := 0
	db.Query(func(d *mgo.Database) {
		c, err = d.C(config.DBT_SUPER_LOGS).Find(m).Count()
	})
	return c, err
}

//分页查询
func FindLogsByMap(m bson.M, skip, limit int) []LogData {
	logData := []LogData{}
	db.Query(func(d *mgo.Database) {
		d.C(config.DBT_SUPER_LOGS).Find(m).Sort("id").Skip(skip).Limit(limit).All(&logData)
	})
	if len(logData) > 0 {
		return logData
	} else {
		return nil
	}
}

//单条插入的方法
func SaveLog2Mgo(logData LogData) error {
	logData.CreatedAt = timeUtils.Format(time.Now())
	logData.id, _ = db.GetNextSeq(config.DBT_SUPER_LOGS)
	return db.InsertMgoData(config.DBT_SUPER_LOGS, logData)
}



//批量插入的方法
func SaveLogs2Mgo(logDatas []LogData) int {
	new := make([]interface{}, len(logDatas))
	for i, logData := range logDatas {
		logData.CreatedAt = timeUtils.Format(time.Now())
		logData.id, _ = db.GetNextSeq(config.DBT_SUPER_LOGS) //自增键
		new[i] = logData
	}
	err, count := db.InsertMgoDatas(config.DBT_SUPER_LOGS, new)
	if err != nil {
		return -1
	}
	return count
}

//删除所有
func DeleteAllLogs2Mgo() error {
	db.Query(func(d *mgo.Database) {
		d.C(config.DBT_SUPER_LOGS).RemoveAll(bson.M{})
	})
	return nil
}
