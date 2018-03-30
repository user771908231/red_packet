package goods

import "casino_redpack/modules"

func IndexHandler(ctx *modules.Context) {

	ctx.HTML(200, "weixin/paywap/paymethod")
}
