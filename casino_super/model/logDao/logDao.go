package logDao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/utils/timeUtils"
	"time"
	"errors"
	"casino_super/conf/config"
	"strconv"
	"casino_common/common/log"
)

//校验上传日志数据的结构
type LogValidater struct {
	Time   string  `json:"time" binding:"Required"`  //日志时间 unix 毫秒
	DeskId string `json:"deskid" binding:"Required"` //桌子id
	UserId string `json:"userid" binding:"Required"` //用户id
	Level  string `json:"level" binding:"Required"`  //日志级别
	Data   string `json:"data" binding:"Required"`   //日志内容
}

//实际存储日志数据的结构
type LogData struct {
	id        int32                                     //自增键
	Time      int64  `json:"time" binding:"Required"`   //日志时间 unix 毫秒
	DeskId    string `json:"deskid" binding:"Required"` //桌子id
	UserId    string `json:"userid" binding:"Required"` //用户id
	Level     string `json:"level" binding:"Required"`  //日志级别
	Data      string `json:"data" binding:"Required"`   //日志内容
	CreatedAt string                                    //创建时间
}

type ReqLog struct {
	Code string `json:"code" binding:"Required;Include(casino)"` //验证码casino
	Logs []LogValidater `json:"logs" binding:"Required"`         //数组
}

//func (l ReqLog) GetReadTime() string {
//
//}

//按记录查询
func FindLogsByKV(key string, v interface{}) []LogValidater {
	logData := []LogValidater{}
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
		d.C(config.DBT_SUPER_LOGS).Find(m).Sort("time").Skip(skip).Limit(limit).All(&logData)
	})
	if len(logData) > 0 {
		return logData
	} else {
		return nil
	}
}



//批量插入的方法
func SaveLogs2Mgo(logValidaters []LogValidater) int {
	new := make([]interface{}, len(logValidaters))
	for i, logValidater := range logValidaters {
		seqId, _ := db.GetNextSeq(config.DBT_SUPER_LOGS) //自增键
		t, err := strconv.ParseInt(logValidater.Time, 10, 64)
		if err != nil {
			t = int64(0)
		}
		logData := LogData{
			id        :seqId,
			Time      :t,
			DeskId    :logValidater.DeskId,
			UserId    :logValidater.UserId,
			Level     :logValidater.Level,
			Data      :logValidater.Data,
			CreatedAt :timeUtils.Format(time.Now()),
		}
		log.T("insert logData %v", logData)
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
