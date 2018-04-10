package keys

import (
	"new_links/modules"
	"gopkg.in/mgo.v2/bson"
	"new_links/model/keysModel"
	"encoding/json"
	"math"

	"fmt"
	"new_links/utils"
)
//列表页
func IndexHandler(ctx *modules.Context) {
	status := ctx.QueryInt("status")
	page := ctx.QueryInt("page")
	query := bson.M{}
	switch status {
	case 1:
		query = bson.M{}
	case 2:
		query = bson.M{
			"status":1,
		}
	case 3:
		query = bson.M{"status":0,}
	default:
		query = bson.M{}
	}
	count,list := keysModel.GetKeysAll(query,page,10)
	data := []bson.M{}
	for _,itme := range list{
		row := bson.M{
			"id":itme.ObjId.Hex(),
			"name":itme.Keys,
			"remarks":itme.Remarks,
			"status":itme.Status,
			"time":itme.Time,
		}
		data = append(data,row)
	}
	ctx.Data["Keys"] = data
	ctx.Data["Keys_page"] = bson.M{
		"count":      count,
		"list_count": len(list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	ctx.HTML(200,"admin/keys/index")
}
//新增
func AddHandler(ctx *modules.Context) {
	ctx.HTML(200,"admin/keys/add")
}
//编辑
func EditHandler(ctx *modules.Context)  {
	id := ctx.Query("id")
	row := keysModel.IdKeyRow(id)
	data := bson.M{
		"id":row.ObjId.Hex(),
		"name":row.Keys,
		"remarks":row.Remarks,
	}

	ctx.Data["Row"] = data
	ctx.HTML(200,"admin/keys/add")
}
//POST提交
func PostAddHandler(ctx *modules.Context)  {
	Id := ctx.Query("id")
	Name := ctx.Query("keys")
	Remarks := ctx.Query("remarks")
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"关键词不能为空！",
	}
	if Id == "" {

		if Name != "" {
			K := keysModel.Keys{
				Keys:Name,
				Remarks:Remarks,
			}
			err := K.Insert()
			if err != nil {
				res["msg"] = "新增关键词失败！"
			}
			res["code"] = 1
			res["message"] = "success"
			res["msg"] = "新增关键词成功！"
		}

	}else {
		K := keysModel.IdKeyRow(Id)
		K.Keys = Name
		K.Remarks = Remarks
		err := K.Update()
		if err != nil {
			res["msg"] = "修改关键词失败！"
		}
		res["code"] = 1
		res["message"] = "success"
		res["msg"] = "编辑关键词成功！"
	}
	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}
//修改状态
func StatusHandler(ctx *modules.Context){
	id := ctx.Query("id")
	Type := ctx.Query("types")
	fmt.Println(id,Type)
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"修改失败",
	}
	if Type == "ok" {
		err := keysModel.KeysStatus(id,1)
		if err == nil {
			res["code"] = 1
			res["message"] = "success"
			res["msg"] = "修改成功！"
		}
	}
	if Type == "no" {
		err := keysModel.KeysStatus(id,0)
		if err == nil {
			res["code"] = 1
			res["message"] = "success"
			res["msg"] = "修改成功！"
		}
	}
	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))
}

func UpdateHandler(ctx *modules.Context)  {

}
//删除
func DelHandler(ctx *modules.Context)  {
	id := ctx.Query("id")
	K := keysModel.GetkeysId(id)
	err := K.Del()
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
//文件上传
func  Uploadhandler(ctx *modules.Context) {
	res := bson.M{
		"code":0,
		"message": "faid",
		"msg":"上传文件失败！",
	}
	_,File,err := ctx.GetFile("file")
	if err == nil {
		err := utils.SaveFileTo(File,"upload/file/",File.Filename)
		if err == nil {
			res["code"] = 1
			res["message"] = "success"
			res["msg"] = "上传文件成功！"

		}
	}

	data,_ := json.Marshal(res)
	ctx.Write([]byte(data))

	str := fmt.Sprintf("upload/file/%s",File.Filename)
	keysModel.OpenFiles(str)
}

