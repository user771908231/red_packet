package manage

import (
	"casino_super/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/service/exchangeService"
	"time"
	"math"
	"casino_super/model/agentModel"
)

//红包与实物兑换

//列表
func ApplyListHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	status := ctx.QueryInt("status")
	if page == 0 {
		page = 1
	}
	if status == 0 {
		status = 1
	}

	query := bson.M{
		"$and": []bson.M{
			bson.M{"status": status,},
		},
	}

	start_time := ctx.Query("start")
	if start_time != "" {
		start,_ := time.Parse("2006-01-02", start_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"requesttime": bson.M{"$gte": start},
		})
	}
	end_time := ctx.Query("end")
	if end_time != "" {
		end,_ := time.Parse("2006-01-02", end_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"requesttime": bson.M{"$lt": end.AddDate(0,0,1)},
		})
	}

	list := []*exchangeService.ExchangeRecord{}
	_,count := db.C(tableName.DBT_AGENT_APPLY_LOG).Page(query, &list, "-requesttime", page, 10)
	ctx.Data["status"] = status
	ctx.Data["list"] = list
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	ctx.HTML(200, "admin/manage/apply/index")
}

//切换状态
func ApplySwitchState(ctx *modules.Context) {
	id := ctx.Query("id")
	status := ctx.QueryInt("status")
	if len(id) != 24 || status <= 0 || status > 3 {
		ctx.Ajax(-1, "参数错误", nil)
		return
	}
	row := agentModel.GetApplyRecordById(id)
	if row == nil {
		ctx.Ajax(-2, "切换状态失败！", nil)
		return
	}
	row.Status = exchangeService.ExchangeState(status)
	row.ProcessTime = time.Now()
	err := row.Save()
	if err != nil {
		ctx.Ajax(-3, "切换状态失败！", nil)
		return
	}
	ctx.Ajax(1,"切换状态成功！", nil)
}
