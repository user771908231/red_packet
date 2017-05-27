package dao

import (
	"casino_mj_changsha/gamedata/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/numUtils"
	"casino_common/utils/db"
	"casino_common/common/log"
	"casino_common/common/consts/tableName"
)

//更具userId查询战绩
func GetMjDeskRoundByUserId(userId uint32) []model.T_mj_desk_round {
	var deskRecords []model.T_mj_desk_round
	querKey, _ := numUtils.Uint2String(userId)
	db.Query(func(d *mgo.Database) {
		d.C(tableName.DBT_MJ_DESK_ROUND).Find(bson.M{"userids":bson.RegEx{querKey, "."}}).Sort("-deskid").Limit(20).All(&deskRecords)
	})

	if deskRecords == nil || len(deskRecords) <= 0 {
		log.T("没有找到玩家[%v]麻将相关的战绩...", userId)
		return nil
	} else {
		return deskRecords
	}
}

//更具userId查询战绩
func GetDdzDeskRoundByUserId(userId uint32) []model.T_ddz_desk_round {
	var deskRecords []model.T_ddz_desk_round
	querKey, _ := numUtils.Uint2String(userId)
	db.Query(func(d *mgo.Database) {
		d.C(tableName.DBT_DDZ_DESK_ROUND).Find(bson.M{"userids":bson.RegEx{querKey, "."}}).Sort("-deskid").Limit(20).All(&deskRecords)
	})

	if deskRecords == nil || len(deskRecords) <= 0 {
		log.T("没有找到玩家[%v]斗地主相关的战绩...", userId)
		return nil
	} else {
		return deskRecords
	}
}
