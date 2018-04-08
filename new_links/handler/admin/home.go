package admin

import (
	"new_links/modules"
)

/**
	后台首页面板
 */
func IndexHandler(ctx *modules.Context) {

	ctx.HTML(200, "admin/index")
}

/**

 */
func MainHandler(ctx *modules.Context) {
	ctx.HTML(200, "admin/main")
}
