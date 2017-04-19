package weixin

import (
	"casino_admin/modules"
	"casino_admin/model/agentModel"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/consts/tableName"
	"time"
	"casino_common/utils/db"
	"math"
	"casino_common/common/userService"
)

//返利记录表
func RebateLogHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}
	agent := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(agent.UnionId)
	list := []agentModel.RebateRecord{}
	query := bson.M{
		"$and": []bson.M{
			bson.M{"agentid": bson.M{"$eq": agent_id}},
		},
	}
	start_time := ctx.Query("start")
	if start_time != "" {
		start,_ := time.Parse("2006-01-02", start_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"time": bson.M{"$gte": start},
		})
	}
	end_time := ctx.Query("end")
	if end_time != "" {
		end,_ := time.Parse("2006-01-02", end_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"time": bson.M{"$lt": end.AddDate(0,0,1)},
		})
	}

	_, count := db.C(tableName.DBT_AGENT_REBATE_LOG).Page(query, &list, "-time", page, 10)

	ctx.Data["list"] = list
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}

	ctx.HTML(200, "weixin/agent/rebate_log")
}

//领取返利
func CheckRebateHandler(ctx *modules.Context) {
	id := ctx.Query("id")
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	row := agentModel.RebateRecord{}
	err := db.C(tableName.DBT_AGENT_REBATE_LOG).Find(bson.M{
		"_id": bson.ObjectIdHex(id),
		"agentid": agent_id,
	}, &row)
	if err != nil {

	}
	if row.IsCheck == true {
		ctx.Ajax(-3, "无法重复领取！", nil)
		return
	}
	//更新状态
	row.IsCheck = true
	row.Save()
	//发放奖励
	_, err = userService.INCRUserRoomcard(agent_id, row.GiveNum)
	if err != nil {
		ctx.Ajax(-2, "领取返利失败！", nil)
		return
	}
	ctx.Ajax(1, "领取返利成功！", nil)
}
