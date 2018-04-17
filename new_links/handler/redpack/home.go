package redpack

import "casino_redpack/modules"

//首页-抢红包
func HomeHandler(ctx *modules.Context) {
	//强制跳转至扫雷
	//ctx.Invoke(SaoleiJLHandler)
	//return
	ctx.HTML(200, "redpack/home/index")
}
