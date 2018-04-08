package group

import (
	"new_links/modules"

	"new_links/model/groupingModel"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

func IndexHandler(ctx *modules.Context) {
	query := bson.M{}
	list := groupingModel.GetGroup(query)
	data := []bson.M{}
	for _,itme := range list{
		row := bson.M{
			"id":itme.ObjId.Hex(),
			"name":itme.GroupName,
			"time":itme.Time,
			"status":itme.Status,
		}
		data = append(data,row)
	}
	ctx.Data["Groups"] =data
	ctx.HTML(200,"admin/group/index")
}

func AddHandler(ctx *modules.Context) {
	ctx.HTML(200,"admin/group/add")
}

func PostAddHandler(ctx *modules.Context) {
	Name := ctx.Query("name")
	Remarks := ctx.Query("remarks")
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"分组名不能为空！",
	}
	if Name != "" {
		G := groupingModel.Grouping{
			GroupName:Name,
			Remarks:Remarks,
		}
		err := G.Insert()
		if err != nil {
			res["msg"] = "新增分组失败！"
		}
		res["code"] = 1
		res["message"] = "success"
		res["msg"] = "新增分组成功！"
	}
	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}

func StatusHandler(ctx *modules.Context) {
	id := ctx.Query("id")
	Type := ctx.Query("type")
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"删除失败",
	}
	if Type == "ok" {
		err := groupingModel.GroupStatus(id,1)
		if err == nil {
			res["code"] = 1
			res["message"] = "success"
			res["msg"] = "修改成功！"
		}
	}
	if Type == "no" {
		err := groupingModel.GroupStatus(id,0)
		if err == nil {
			res["code"] = 1
			res["message"] = "success"
			res["msg"] = "修改成功！"
		}
	}
	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}

func DelHandler(ctx *modules.Context) {
	str := ctx.Query("id")
	err := groupingModel.DelGroup(str)
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"删除失败",
	}
	if err == nil {
		res["code"] = 1
		res["message"] = "success"
		res["msg"] = "删除成功！"
	}
	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}