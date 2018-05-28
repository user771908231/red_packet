package agent

import (
	"casino_redpack/modules"
	"casino_redpack/model/agentModel"
	"gopkg.in/mgo.v2/bson"
	"math"
	"casino_redpack/model/userModel"
)

func AgentRebateLog(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	count,list := agentModel.GetAgentRebateLogPage(bson.M{},page,10)
	list_data := []bson.M{}
	for _,item := range list {
		row := bson.M{
			"id":item.Id.Hex(),
			"agent_id":item.AgentId,
			"rebate_money":agentModel.FloatValue(item.RebateMoeny,2),
			"rebate_id":item.RebateId,
			"nickname":userModel.GetUsernicknameById(item.AgentId),
			"rebate_nickname":userModel.GetUsernicknameById(item.RebateId),
		}
		list_data = append(list_data,row)
	}
	data := bson.M{
		"list": list_data,
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