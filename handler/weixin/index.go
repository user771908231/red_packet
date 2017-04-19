package weixin

import (
	"casino_admin/modules"
)

func MainHandler(ctx *modules.Context) {

	ctx.HTML(200, "weixin/agent/index")
}
