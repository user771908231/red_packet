package links

import (
	"new_links/modules"
	"fmt"
	"new_links/model/linksModel"
	"gopkg.in/mgo.v2/bson"
	"math"
	"new_links/model/groupingModel"
)

func IndexHandler(ctx *modules.Context) {
	status := ctx.QueryInt("status")
	page := ctx.QueryInt("page")
	fmt.Println(status,page)
	switch status {
	case 1:
		query := bson.M{}
		count,list := linksModel.GetLinksAll(query,page,10)
		ctx.Data["list"] = list
		ctx.Data["page"] = bson.M{
			"count":      count,
			"list_count": len(list),
			"limit":      10,
			"page":       page,
			"page_count": math.Ceil(float64(count) / float64(10)),
		}
	}
	ctx.HTML(200,"admin/links/index")
}

func AddHandler(ctx *modules.Context) {
	query := bson.M{}
	list := groupingModel.GetGroup(query)
	data := []bson.M{}
	for _,item := range list {
		row := bson.M{
			"id":item.ObjId.Hex(),
			"name":item.GroupName,
		}
		data = append(data,row)
	}

	ctx.Data["Gourps"] =data
	ctx.HTML(200,"admin/links/add")
}
