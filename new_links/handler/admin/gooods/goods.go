package goods

import (
	"casino_redpack/modules"
	"gopkg.in/mgo.v2/bson"
	"casino_redpack/model/googsModel"
	"encoding/json"

)


func IndexHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	status := ctx.QueryInt("status")
	query := bson.M{}
	switch status {
	case 1:
		query = bson.M{
			"status":1,
		}
	case 2:
		query = bson.M{
			"status" :0,
		}
	default:
		query = bson.M{

		}
	}
	list := googsModel.GetWhole(query,page)
	ctx.Data["GoogsCoin"] = list
	ctx.HTML(200, "admin/block/Googs/index")
}

func AddHandler(ctx *modules.Context) {
	ctx.HTML(200, "admin/block/Googs/add")
}

func UpdateHandler(ctx *modules.Context)  {
	ctx.HTML(200, "admin/block/Googs/uplate")
}

func AddPost(Form googsModel.GoogsForm,ctx *modules.Context){
	err := googsModel.AddOne(Form)
	if err == nil {
		ctx.Success("成功！", "/admin/manage/goods/add", 1)
		return
	}
	ctx.Success("失败！", "/admin/manage/goods/add", 1)
	return
}

func OperationHandler(ctx *modules.Context){
	Id := ctx.Query("id")
	types := ctx.Query("types")
	err := googsModel.Operation(Id,types)
	list := bson.M{
		"code": 0,
		"message": "fail",
		"msg": err,
	}
	if err == nil {
		list["code"] = 1
		list["message"] = "success"
		list["msg"] = "操作成功!"
	}

	data,_ := json.Marshal(list)
	ctx.Write([]byte(data))
}

