package redpack

import "casino_redpack/modules"

//五人对战
func WurenDZHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/wurenduizhan")
}

//发红包
func WurenFahongbaoHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/wurenduizhan_add")
}

//开红包结果
func WurenKaiOkHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/wurenduizhan_kaiok")
}
