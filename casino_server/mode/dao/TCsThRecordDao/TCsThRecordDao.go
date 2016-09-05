package TCsThRecordDao

import (
	"casino_server/mode"
	"casino_server/utils/db"
	"gopkg.in/mgo.v2"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
)

//数据查询全部放在dao层

func GetByMatchId(matchId int32) *mode.T_cs_th_record {
	ret := &mode.T_cs_th_record{}
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_CS_TH_RECORD).Find(bson.M{"id":matchId}).One(ret)
	})
	if ret.Id <= 0 {
		return nil
	} else {
		return ret
	}
}
