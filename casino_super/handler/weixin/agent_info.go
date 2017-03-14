package weixin

import (
	"casino_super/modules"
	"casino_super/model/agentModel"
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
	ctx.Data["info"] = agent_info
	ctx.HTML(200, "weixin/agent/info")
}
