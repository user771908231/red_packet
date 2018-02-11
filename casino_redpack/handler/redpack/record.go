package redpack

import "casino_redpack/modules"

//游戏记录
func RecordHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/record")
}
