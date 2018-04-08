package links

import (
	"new_links/modules"
	"fmt"
)

func IndexHandler(ctx *modules.Context) {
	status := ctx.QueryInt("status")
	page := ctx.QueryInt("page")
	fmt.Println(status,page)
}