package weixin

import "casino_super/modules"

func MainHandler(ctx *modules.Context) {
	ctx.HTML(200, "weixin/user/index")
}

func InfoHandler(ctx *modules.Context) {
	ctx.HTML(200, "weixin/user/info")
}