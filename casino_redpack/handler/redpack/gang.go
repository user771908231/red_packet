package redpack

import "casino_redpack/modules"

//二八杠首页
func GangHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/gang")
}

//二八杠发红包页面
func GangAddHandler(ctx *modules.Context) {
	ctx.HTML(200, "redpack/home/gang_add")
}

