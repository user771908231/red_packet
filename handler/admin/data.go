package admin

import (
	"casino_admin/modules"
	"casino_admin/model/dataModel"
	"fmt"
	"time"
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
//在线统计列表
func OnlineStaticList(ctx *modules.Context) {
	date_start := ctx.Query("date_start")
	date_end := ctx.Query("date_end")
	date1,_ := time.Parse("2006-01-02",date_start)
	date2,_ := time.Parse("2006-01-02",date_end)

	ctx.Data["date_start"] =date_start
	ctx.Data["date_end"] =date_end

	err :=dataModel.OnlineStaticList(date1,date2)
	ctx.Data["static"] =err

	ctx.HTML(200,"admin/data/onlineStatic")
}
