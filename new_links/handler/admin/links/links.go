package links

import (
	"new_links/modules"
	"fmt"
	"new_links/model/linksModel"
	"gopkg.in/mgo.v2/bson"
	"math"
	"new_links/model/groupingModel"
	"encoding/json"
)

func IndexHandler(ctx *modules.Context) {
	status := ctx.QueryInt("status")
	page := ctx.QueryInt("page")
	group := ctx.Query("group")
	fmt.Println(group)
	query := bson.M{}
	switch status {
	case 1:
		if group != "" {
			query = bson.M{"gruopid":bson.ObjectIdHex(group)}
		}else{
			query = bson.M{}
		}
	}
	count,list := linksModel.GetLinksAll(query,page,10)
	data := []bson.M{}
	for _,item := range list {
		row := bson.M{
			"id":item.ObjId.Hex(),
			"group":groupingModel.GetGroupObjId(item.GruopId).GroupName,
			"push":item.Push*100,
			"remarks":item.Remarks,
			"time":item.Time,
			"status":item.Status,
			"visit":linksModel.GetDayVisit(item.ObjId),
			"weight":item.Weight,
			"id_number":item.Id,
			"quota":item.Quota,
			"excess_id":item.ExcessId,

		}
		data = append(data,row)
	}
	Q := bson.M{"status":1}
	grouping := groupingModel.GetGroup(Q)
	grouping_data := []bson.M{}
	for _,item := range grouping {
		row := bson.M{
			"id":item.ObjId.Hex(),
			"name":item.GroupName,
		}
		grouping_data = append(grouping_data,row)
	}
	ctx.Data["status"] = status
	ctx.Data["Gourps"] =grouping_data
	ctx.Data["list"] = data
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
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




func PostAddHandler(ctx *modules.Context,form linksModel.PostForm) {
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"新增链接错误！",
	}
	err := linksModel.Createlink(form)
	if err == nil {
		res["code"] = 1
		res["message"] = "success"
		res["msg"] = "新增链接成功！"
	}

	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}


func Delhandler(ctx *modules.Context) {
	id := ctx.Query("id")
	err := linksModel.LinksIdDel(bson.ObjectIdHex(id))
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"删除失败！",
	}
	if err == nil {
		res["code"] = 1
		res["message"] = "success"
		res["msg"] = "删除成功！"
	}

	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}

func Statushandler(ctx *modules.Context) {
	id := ctx.Query("id")
	Type := ctx.Query("types")
	fmt.Println(id,Type)
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"修改失败",
	}
	if Type == "ok" {
		err := linksModel.LinksStatus(id,1)
		if err == nil {
			res["code"] = 1
			res["message"] = "success"
			res["msg"] = "修改成功！"
		}
	}
	if Type == "no" {
		err := linksModel.LinksStatus(id,0)
		if err == nil {
			res["code"] = 1
			res["message"] = "success"
			res["msg"] = "修改成功！"
		}
	}
	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}

func Edithandler(ctx *modules.Context) {
	id := ctx.Query("id")
	query := bson.M{}
	if id != "" {
		L := linksModel.GetLinkId(id)
		list := groupingModel.GetGroup(query)
		data := []bson.M{}
		for _,item := range list {
			row := bson.M{
				"id":item.ObjId.Hex(),
				"name":item.GroupName,
			}
			data = append(data,row)
		}

		ctx.Data["Links"] =bson.M{
			"obj_id":L.ObjId.Hex(),
			"group_id":L.GruopId.Hex(),
			"url":L.Url,
			"id":L.Id,
			"weight":L.Weight,
			"Remarks":L.Remarks,
			"quota":L.Quota,
			"excess_id":L.ExcessId,
		}
		ctx.Data["Gourps"] =data
	}
	ctx.HTML(200,"admin/links/edit")
}

func Uploadhandler(ctx *modules.Context ,Upload linksModel.PostUpload) {
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"修改失败",
	}
	defer func() {
		data,_ := json.Marshal(res)
		ctx.Write([]byte(data))
	}()
	//link := linksModel.Links{
	//	ObjId:bson.ObjectIdHex(Upload.ObjId),
	//	GruopId:bson.ObjectIdHex(Upload.Group),
	//	Id:Upload.Id,
	//	Url:Upload.Url,
	//	KeysId:bson.ObjectIdHex(Upload.Keys),
	//	Weight:Upload.Push,
	//	Remarks:Upload.Remarks,
	//}
	link := linksModel.GetLinkId(Upload.ObjId)
	if link == nil {
		fmt.Println("kong")
		return
	}
	link.GruopId = bson.ObjectIdHex(Upload.Group)
	link.Id = Upload.Id
	link.Url = Upload.Url
	link.Weight = Upload.Push
	link.Remarks = Upload.Remarks
	link.Quota = Upload.Quota
	link.ExcessId = Upload.ExcessId
	err := link.Update()
	if err == nil {
		defer func() {
			linksModel.LInskPush(Upload.Group)
		}()
		res["code"] = 1
		res["message"] = "success"
		res["msg"] = "修改成功！"
	}


}