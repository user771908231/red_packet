package configModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"


)
const ADMIN_GAME_TASK string = "t_config_task_info"
type TaskCongif struct {
	Id          bson.ObjectId  			 `bson:"_id"`		//用户ID
	TaskId      	float64 	 		`bson:"taskid"`		//游戏ID
	CateId		float64  		`bson:"cateid"`
	TaskType	string  		`bson:"tasktype"`
	GameId		float64  		`bson:"gameid"`
	Title		string  		`bson:"Title"`
	Description	string  		`bson:"Description"`
	Sort		float64  		`bson:"Sort"`
	TaskSum		float64  		`bson:"TaskSum"`
	RepeatSum	float64  		`bson:"RepeatSum"`
}

//登录服务器配置更新
func GameConfigTask(Id bson.ObjectId,TaskId float64,CateId float64,TaskType string,GameId float64,Title string,Description string,Sort float64,TaskSum float64,RepeatSum float64) (Task_Congif *TaskCongif) {
	err := db.C(ADMIN_GAME_TASK).Update(bson.M{"_id": Id},bson.M{
		"$set" : bson.M{
			"taskid" : TaskId,
			"cateid" : CateId,
			"tasktype" : TaskType,
			"gameid" : GameId,
			"title" : Title,
			"description" : Description,
			"sort" : Sort,
			"tasksum" : TaskSum,
			"repeatsum" : RepeatSum,
		},})
	if err != nil {
		return nil
	}
	return Task_Congif
}
