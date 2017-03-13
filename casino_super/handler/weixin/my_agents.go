package weixin

import (
	"casino_super/modules"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_super/model/agentModel"
	"math"
	"casino_common/utils/db"
)

func MyAgentsHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	my_agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	my_agent_info := agentModel.GetAgentInfoById(my_agent_id)
	agent_info := my_agent_info
	var agent_id uint32 = uint32(ctx.QueryInt("agent_id"))
	if agent_id == 0 {
		agent_id = my_agent_id
	}else {
		agent_info = agentModel.GetAgentInfoById(agent_id)
		my_agent_info = agentModel.GetAgentInfoById(my_agent_id)
		if agent_info == nil || my_agent_info.Level != 1 || my_agent_info.UserId != agent_info.RootId {
			ctx.Error("对不起，您没有权限访问此页面！", "", 0)
			return
		}
	}

	user_list := []*agentModel.AgentInfo{}

	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	count := 0
	if agent_id >0 {
		_,count = db.C(tableName.DBT_AGENT_INFO).Page(bson.M{
			"pid": agent_id,
		}, &user_list, "", page, 10)
	}

	list := []bson.M{}
	for _,user := range user_list {
		week_sales_count := agentModel.GetAgentThisWeekOneUserSalesCount(agent_id, user.UserId)
		all_sales_count := agentModel.GetAgentUserSalesCount(agent_id, user.UserId)
		list = append(list, bson.M{
			"Id": user.UserId,
			"NickName": user.NickName,
			"AllSalesCount": all_sales_count,
			"WeekSalesCount": week_sales_count,
			"ChildNum": agentModel.GetAgentChildNum(agent_id),
		})
	}
	ctx.Data["list"] = list
	ctx.Data["agent_info"] = agent_info
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(user_list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}

	ctx.HTML(200, "weixin/agent/my_agents")
}
