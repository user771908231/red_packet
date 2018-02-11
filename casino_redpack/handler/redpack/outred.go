package redpack

import "casino_redpack/modules"

//提现
func OutredHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/outred")
}
