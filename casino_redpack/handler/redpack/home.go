package redpack

import "casino_redpack/modules"

//首页-抢红包
func HomeHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/index")
}
