package weixin

import "casino_super/modules"

func MainHandler(ctx *modules.Context) {

	ctx.HTML(200, "weixin/agent/index")
}

//个人信息
func InfoHandler(ctx *modules.Context) {
	ctx.HTML(200, "weixin/agent/info")
}
