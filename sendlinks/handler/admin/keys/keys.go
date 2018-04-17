package keys

import "sendlinks/modules"

func IndexHandler(ctx *modules.Context) {

	ctx.HTML(200,"x-admin/keys/index")
}


func AddHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/keys/add")
}

func DeleteHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/keys/delete")
}


func SelectHandler(ctx *modules.Context) {
	//ctx.HTML(200,"x-admin/links/index")
}

func UplateHandler(ctx *modules.Context) {
	ctx.HTML(200,"x-admin/keys/edit")
}

