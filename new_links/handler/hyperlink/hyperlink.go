package hyperlink

import (
	"new_links/modules"
	"fmt"
	"new_links/model/hyperlinkModel"
)

func Indexhandler(ctx *modules.Context) {
	host := ctx.Query("host")

	fmt.Println(host)
	hyperlinkModel.GetGroup(host)

	ctx.Redirect("",320)

}
