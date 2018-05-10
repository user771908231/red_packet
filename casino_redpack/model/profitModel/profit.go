package profitModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"casino_redpack/model/agentModel"
	"casino_common/common/consts/tableName"
)
//总的收益
type Profit struct {
	Id bson.ObjectId `bson:"_id"`
	UserId uint32 //玩家ID
	Name string
	ProfitMoney float64
}
//每天的收益
type ProfitDetailed struct {
	Id bson.ObjectId `bson:"_id"`
	UserId uint32 //玩家ID
	Name string
	ProfitMoney float64
	Time time.Time
}

func (P *Profit) Insert() error {
	P.Id = bson.NewObjectId()
	err := db.C(tableName.TABLE_REDPACK_PROFIT_LOG).Upsert(bson.M{"_id":P.Id},P)
	return err
}

func (P *ProfitDetailed) Insert() error {
	P.Id = bson.NewObjectId()
	P.Time = agentModel.TimeObject()
	err := db.C(tableName.TABLE_REDPACK_PROFIT_ONE_DAY_LOG).Upsert(bson.M{"_id":P.Id},P)
	return err
}




