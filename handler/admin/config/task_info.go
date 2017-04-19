package config

import (
	"casino_admin/modules"
	"github.com/go-macaron/binding"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/service/taskService/taskType"
)

//商品配置列表
func TaskListHandler(ctx *modules.Context) {
	list := []*taskType.TaskInfo{}
	cate_id := ctx.QueryInt("cate_id")
	query := bson.M{}
	if cate_id > 0 {
		query["cateid"] = cate_id
	}
	db.C(tableName.DBT_T_TASK_INFO).FindAll(query, &list)
	ctx.Data["cate_id"] = cate_id
	ctx.Data["list"] = list
	ctx.HTML(200, "admin/config/task/index")
}

//编辑
func TaskEditPost(ctx *modules.Context, form taskType.TaskInfo, errs binding.Errors) {
	if errs.Len() > 0 {
		ctx.Ajax(-1, "表单参数错误！", nil)
		return
	}
	//先还原obj_id
	form.Id = bson.ObjectIdHex(string(form.Id))

	//更新一维数据
	err := db.C(tableName.DBT_T_TASK_INFO).Update(bson.M{
		"_id": form.Id,
	}, bson.M{
		"$set": bson.M{
			"taskid": form.TaskId,
			"cateid": form.CateId,
			"tasktype": form.TaskType,
			//"gameid": form.GameId,
			"title": form.Title,
			"sort": form.Sort,
			"tasksum": form.TaskSum,
			"repeatsum": form.RepeatSum,
			"description": form.Description,
		},
	})

	if err != nil {
		ctx.Ajax(-1, "编辑配置失败！", err.Error())
		return
	}

	ctx.Ajax(1, "编辑配置成功！", nil)
}

//编辑奖励
func TaskEditReward(ctx *modules.Context) {

}