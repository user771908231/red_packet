package admin

import (
	"casino_admin/modules"
	"casino_admin/model/dataModel"
	"fmt"
)

//在房间玩家
func AtHome(ctx *modules.Context) {
	result := dataModel.AtHome()
	if result !=nil{
		fmt.Println("success",result)
		ctx.Data["user"] = result
		ctx.HTML(200,"admin/data/atHome")
	}

}


//条件查询
func AtHomeList(ctx *modules.Context) {
	GameID := ctx.Query("gameID")
	result :=dataModel.AtHomeList(GameID)
	ctx.Data["user"] = result
	fmt.Println("success2222",result)
	ctx.HTML(200,"admin/data/atHome")
}

//在线统计
func OnlineStatic(ctx *modules.Context) {
	ctx.HTML(200,"admin/data/onlineStatic")
}
