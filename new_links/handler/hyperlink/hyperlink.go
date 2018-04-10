package hyperlink

import (
	"new_links/modules"
	"new_links/model/hyperlinkModel"
)

func Indexhandler(ctx *modules.Context) {
	host := ctx.Query("host")
	url := hyperlinkModel.GetGroup(host)
	ctx.Redirect(url,302)

}


