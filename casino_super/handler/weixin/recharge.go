package weixin

import (
	"casino_super/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"casino_super/model/agentModel"
	"gopkg.in/mgo.v2/bson"
)

//充值列表
func RechargeListHandler(ctx *modules.Context) {
	list := []agentModel.Goods{}
	db.C(tableName.DBT_AGENT_GOODS).FindAll(bson.M{}, &list)
	ctx.Data["goods_list"] = list
	ctx.HTML(200, "weixin/agent/recharge")
}
