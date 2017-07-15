package logDao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"time"
	"errors"
	"casino_super/conf/config"
	"strconv"
	"casino_common/common/log"
	"fmt"
	"casino_common/utils/timeUtils"
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
	CreatedAt time.Time                                 //unix转time
}

type ReqLog struct {
	Code string `json:"code" binding:"Required;Include(casino)"` //验证码casino
	Logs []LogValidater `json:"logs" binding:"Required"`         //数组
}

func (l LogData) GetReadableCreatedAt() string {
	return timeUtils.Format(l.CreatedAt.UTC())
}

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
func FindLogsByMapCount(t string, m bson.M) (int, error) {
	err := errors.New("")
	c := 0
	log.T("开始从表[%v]中查询条数... m[%v]", t, m)
	db.Query(func(d *mgo.Database) {
		c, err = d.C(t).Find(m).Count()
	})
	return c, err
}

//分页查询
func FindLogsByMap(t string, m bson.M, skip, limit int) []LogData {
	logData := []LogData{}
	log.T("开始从表[%v]中查询数据... m[%v] skip[%v] limit[%v]", t, m, skip, limit)
	db.Query(func(d *mgo.Database) {
		d.C(t).Find(m).Sort("time").Skip(skip).Limit(limit).All(&logData)
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

	userId := "0"
	for i, logValidater := range logValidaters {
		seqId, _ := db.GetNextSeq(config.DBT_SUPER_LOGS) //自增键
		int64Time, err := strconv.ParseInt(logValidater.Time, 10, 64)
		if err != nil {
			int64Time = int64(0)
		}
		sec, _ := strconv.ParseInt(logValidater.Time[0:10], 10, 64)
		nsec, _ := strconv.ParseInt(logValidater.Time[10:], 10, 64)
		unixTime := timeUtils.String2YYYYMMDDHHMMSS(timeUtils.Format(time.Unix(sec, nsec)))
		println(fmt.Sprintf("unix %v unixTime %v", logValidater.Time, unixTime))
		if logValidater.UserId != "" {
			userId = logValidater.UserId
		}

		logData := LogData{
			id        :seqId,
			Time      :int64Time,
			DeskId    :logValidater.DeskId,
			UserId    :logValidater.UserId,
			Level     :logValidater.Level,
			Data      :logValidater.Data,
			CreatedAt :unixTime,
		}

		log.T("insert logData %v", logData)
		new[i] = logData
	}
	tb := GetTableName(config.DBT_SUPER_LOGS, timeUtils.String2YYYYMMDDHHMMSS(timeUtils.Format(time.Now())), userId)
	log.T("插入[%v]条数据到表[%v]中...", len(new), tb)
	err, count := db.InsertMgoDatas(tb, new)
	if err != nil {
		return -1
	}
	return count
}

func GetTableName(prefixName string, t time.Time, userId string) string {
	year, month, day := t.Date()
	userId64, _ := strconv.ParseInt(userId, 10, 64)
	userId64 = userId64 % 100
	//log.T("year %v month%v day%v", year, month, day)
	return fmt.Sprintf("%s_%d%d%d_%v", prefixName, year, month, day, userId64)
}

//删除所有
func DeleteAllLogs2Mgo() error {
	db.Query(func(d *mgo.Database) {
		d.C(config.DBT_SUPER_LOGS).RemoveAll(bson.M{})
	})
	return nil
}
