package keys

import (
	"new_links/modules"
	"gopkg.in/mgo.v2/bson"
	"new_links/model/keysModel"
	"encoding/json"
	"math"

)

func IndexHandler(ctx *modules.Context) {
	status := ctx.QueryInt("status")
	page := ctx.QueryInt("page")
	switch status {
	case 1:
		query := bson.M{}
		count,list := keysModel.GetKeysAll(query,page,10)
		data := []bson.M{}
		for _,itme := range list{
			row := bson.M{
				"id":itme.ObjId.Hex(),
				"name":itme.Keys,
				"remarks":itme.Remarks,
				"status":itme.Status,
				"time":itme.Time.Unix(),
			}
			data = append(data,row)
		}
		ctx.Data["keys"] = data
		ctx.Data["keys_page"] = bson.M{
			"count":      count,
			"list_count": len(list),
			"limit":      10,
			"page":       page,
			"page_count": math.Ceil(float64(count) / float64(10)),
		}

	}
	ctx.HTML(200,"admin/keys/index")
}

func AddHandler(ctx *modules.Context) {
	ctx.HTML(200,"admin/keys/add")
}

func EditHandler(ctx *modules.Context)  {
	id := ctx.Query("id")
	row := keysModel.IdKeyRow(id)
	data := bson.M{
		"id":row.ObjId.Hex(),
		"name":row.Keys,
		"remarks":row.Remarks,
	}

	ctx.Data["row"] = data
	ctx.HTML(200,"admin/keys/add")
}

func PostAddHandler(ctx *modules.Context)  {
	Name := ctx.Query("keys")
	Remarks := ctx.Query("remarks")
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"关键词不能为空！",
	}
	if Name != "" {
		K := keysModel.Keys{
			Keys:Name,
			Remarks:Remarks,
		}
		err := K.Insert()
		if err != nil {
			res["msg"] = "新增关键词失败！"
		}
		res["message"] = "success"
		res["msg"] = "新增关键词成功！"
	}
	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}

func StatusHandler(ctx *modules.Context){

}

func UpdateHandler(ctx *modules.Context)  {

}

func DelHandler(ctx *modules.Context)  {

}

