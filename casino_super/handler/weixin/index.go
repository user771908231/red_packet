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
	ctx.HTML(200, "weixin/agent/index")
}
