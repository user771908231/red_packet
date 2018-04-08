package weixin

import (
	"casino_redpack/modules"
	"casino_common/common/model/agentModel"
	"casino_common/common/userService"
)

//代理个人信息
func AgentInfoHandler(ctx *modules.Context) {
	agent_id := uint32(ctx.QueryInt("agent_id"))
	if agent_id == 0 {
		wx_info := ctx.IsWxLogin()
		agent_id = agentModel.GetUserIdByUnionId(wx_info.UnionId)
	}
	agent_info := agentModel.GetAgentInfoById(agent_id)
	if agent_info == nil {
		ctx.Error("查看代理信息失败！", "", 0)
		return
	}
	my_roomcard := userService.GetUserRoomCard(agent_id)
	ctx.Data["my_roomcard"] = my_roomcard
	ctx.Data["all_sales_count"] = agentModel.GetAgentSalesCount(agent_id)
	//ctx.Data["all_rebate_count"] = agentModel.GetAgentAllRebateRoomCardNum(agent_id)
	ctx.Data["info"] = agent_info
	ctx.HTML(200, "weixin/agent/info")
}
