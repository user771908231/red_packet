package redpack

import "casino_redpack/modules"

//牛牛首页
func NiuniuHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/niuniu")
}

//牛牛发红包页面
func NiuniuAddHandler(ctx *modules.Context) {
	ctx.HTML(200, "redpack/home/niuniu_add")
}
