package test

import (
	"testing"
	"casino_server/service/CSTHService"
	"casino_server/common/log"
	"casino_server/mode"
	"gopkg.in/mgo.v2"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/utils/db"
)


//测试竞标赛相关的api
func TestCSTHAPI(t *testing.T) {
	testGetMatchList()
}

func testGetMatchList() {

	//从数据库查询
	data1 := []mode.T_cs_th_record{}
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_CS_TH_RECORD).Find(bson.M{"id":bson.M{"$gt":0}}).Sort("-id").Limit(20).All(&data1)
	})
	log.T("从数据库中查询的结果[%v]", data1)


	//通过相关的方法来查询
	log.T("开始查询竞标赛列表")
	data2 := CSTHService.GetGameMatchList()
	log.T("开始查询竞标赛列表结果[%v]", data2)

}