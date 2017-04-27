package logDao

import (
	"testing"
	"time"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/timeUtils"
	"casino_super/conf/config"
)

func TestSaveLogs2Mgo(t *testing.T) {
	t1 := time.Now()
	println(fmt.Sprintf("TestSaveLogs2Mgo begin %v", t1))
	for i := 0; i < 500; i++ {
		SaveLogs2Mgo(getTestLogs())
	}
	t2 := time.Now()
	println(fmt.Sprintf("TestSaveLogs2Mgo end %v", t2))
	println(fmt.Sprintf("TestSaveLogs2Mgo cost time %v", t2.Sub(t1)))
	time.Sleep(time.Duration(10) * time.Second)
}

func TestSaveLog2Mgo(t *testing.T) {
	t1 := time.Now()
	println(fmt.Sprintf("TestSaveLog2Mgo begin %v", t1))
	for i := 0; i < 500; i++ {
		//for _, l := range getTestLogs() {
		//	SaveLog2Mgo(l)
		//}
	}
	t2 := time.Now()
	println(fmt.Sprintf("TestSaveLog2Mgo end %v", t2))
	println(fmt.Sprintf("TestSaveLog2Mgo cost time %v", t2.Sub(t1)))
	time.Sleep(time.Duration(10) * time.Second)
}


func getTestLogs() []LogValidater {
	logs := []LogValidater{
		LogValidater{
			Time: "1490148861570",
			DeskId: "1",
			UserId: "2",
			Level: "1",
			Data: "过滤123吧",
		},
		LogValidater{
			Time: "1490148861572",
			DeskId: "2",
			UserId: "2",
			Level: "2",
			Data: "过滤我吧",
		},
	}

	for i := 0; i < 100; i++ {
		logs2 := []LogValidater{
			LogValidater{
				Time: "1490148861573",
				DeskId: "1",
				UserId: "2",
				Level: "2",
				Data: "过滤SsS吧",
			},
			LogValidater{
				Time: "1490148861574",
				DeskId: "2",
				UserId: "2",
				Level: "3",
				Data: "过滤SsS",
			},
		}
		logs = append(logs, logs2...)
	}
	return logs
}



func TestFindLogsByMap(t *testing.T) {
	m := bson.M{}
	//regString := "123"
	//m["data"] = bson.M{
	//	"$not" : bson.RegEx{
	//		regString,
	//		"i",
	//	},
	//}
	//m["data"] = bson.M{
	//	"$regex" : regString,
	//}
	timeBegin := timeUtils.StringYYYYMMDD2time("2017-01-23")
	timeEnd := timeBegin.AddDate(0, 0, 1)
	timeBeginS := timeUtils.Format(timeBegin)
	timeEndS := timeUtils.Format(timeEnd)
	println(fmt.Sprintf("begin %v end %v", timeBeginS, timeEndS))
	m["createdatt"] = bson.M{
		"$gte" : timeBegin,
		"$lt" : timeEnd,
	}

	logs := FindLogsByMap(GetTableName(config.DBT_SUPER_LOGS, time.Now(), "10001"), bson.M{}, 0, 100)
	for _, log := range logs {
		println(fmt.Sprintf("%v", log))
	}
	println(2)
	time.Sleep(time.Duration(2) * time.Second)
}

func TestGetTableName(t *testing.T) {
	println(GetTableName(config.DBT_SUPER_LOGS, time.Now(), "10751"))
}
