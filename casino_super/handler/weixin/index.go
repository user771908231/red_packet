package weixin

import (
	"casino_super/modules"
	"casino_super/model/agentModel"
	"casino_common/common/userService"
)

func MainHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	my_roomcard := userService.GetUserRoomCard(agent_id)
	ctx.Data["my_roomcard"] = my_roomcard
	ctx.Data["all_sales_count"] = agentModel.GetAgentSalesCount(agent_id)
	ctx.Data["all_rebate_count"] = agentModel.GetAgentAllRebateRoomCardNum(agent_id)
	ctx.HTML(200, "weixin/agent/index")
}
