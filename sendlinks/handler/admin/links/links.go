package links

import "sendlinks/modules"

func IndexHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/links/index")
}


func AddHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/links/add")
}

func DeleteHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/links/delete")
}


func SelectHandler(ctx *modules.Context) {
	//ctx.HTML(200,"x-admin/links/index")
}

func UplateHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/links/edit")
}

