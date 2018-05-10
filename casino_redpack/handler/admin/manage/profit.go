package manage

import (
	"casino_redpack/modules"
	"casino_redpack/model/profitModel"
	"gopkg.in/mgo.v2/bson"
)

func ProfitHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	if page < 1 {
		page = 1
	}
	//status := ctx.QueryInt("status")
	var data bson.M
	data = profitModel.GetSelectProfit(bson.M{},page)

	ctx.Data["data"] = data
	ctx.HTML(200,"admin/manage/profit/index")
}

func ProfitDetailedHandler(ctx *modules.Context) {
	id := uint32(ctx.QueryInt("uid"))
	page := ctx.QueryInt("page")
	if page < 1 {
		page = 1
	}
	data := profitModel.GetUserIdSelectProfitDetailedAll(bson.M{"userid":id},page)
	ctx.Data["data"] = data
	ctx.HTML(200,"admin/manage/profit/index")

}
