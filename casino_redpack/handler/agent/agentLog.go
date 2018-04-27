package agent

import (
	"casino_redpack/modules"
	"casino_redpack/model/agentModel"
	"gopkg.in/mgo.v2/bson"
	"math"
)

func AgentRebateLog(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	count,list := agentModel.GetAgentRebateLogPage(bson.M{},page,10)
	data := bson.M{
		"list": list,
		"page": bson.M{
			"count":      count,
			"list_count": len(list),
			"limit":      10,
			"page":       page,
			"page_count": math.Ceil(float64(count) / float64(10)),
		},
	}
	ctx.Data["agent_rebate_list"] = data
	ctx.HTML(200,"admin/manage/rebate/index")
}