package TThDeskRecordDao

import (
	"casino_server/mode"
	"gopkg.in/mgo.v2"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/utils/db"
	"casino_server/utils/numUtils"
	"casino_server/common/log"
)

//找有某个人的战绩
func Find(userId uint32, limit int) []mode.T_th_desk_record {
	var deskRecords []mode.T_th_desk_record
	queryKey, _ := numUtils.Uint2String(userId)
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_TH_DESK_RECORD).Find(bson.M{"userids": bson.RegEx{queryKey, "."}}).Sort("-id").Limit(20).All(&deskRecords)
	})
	log.T("找到的战绩[%v]",deskRecords)
	return deskRecords
}
