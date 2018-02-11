package redpack

import "casino_redpack/modules"

//代理中心
func AgentHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/agent/index")
}
