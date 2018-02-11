package config

import (
	"casino_redpack/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/service/taskService/taskType"
	"casino_redpack/model/configModel"
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
func TaskEditPost(ctx *modules.Context) {
	id := ctx.Query("id")
	obj_id :=bson.ObjectIdHex(id)
	TaskId := ctx.QueryFloat64("TaskId")
	CateId := ctx.QueryFloat64("CateId")
	TaskType := ctx.Query("TaskType")
	GameId := ctx.QueryFloat64("GameId")
	Title := ctx.Query("Title")
	Description := ctx.Query("Description")
	Sort := ctx.QueryFloat64("Sort")
	TaskSum := ctx.QueryFloat64("TaskSum")
	RepeatSum := ctx.QueryFloat64("RepeatSum")
	configModel.GameConfigTask(obj_id,TaskId,CateId,TaskType,GameId,Title,Description,Sort,TaskSum,RepeatSum)
}

//编辑奖励
func TaskEditReward(ctx *modules.Context) {

}
