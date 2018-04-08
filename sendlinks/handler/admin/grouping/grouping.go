package grouping

import "sendlinks/modules"

func IndexHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/grouping/index")
}


func AddHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/grouping/add")
}

func DeleteHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/grouping/delete")
}


func SelectHandler(ctx *modules.Context) {
	//ctx.HTML(200,"x-admin/links/index")
}

func UplateHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/grouping/edit")
}

