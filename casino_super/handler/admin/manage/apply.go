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
	"casino_common/common/userService"
)

//红包与实物兑换

//列表
func ApplyListHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	status := ctx.QueryInt("status")
	invitedid := ctx.QueryInt("invited")

	if page == 0 {
		page = 1
	}

	//查询
	query := bson.M{
		"$and": []bson.M{},
	}

	//申请状态
	if status > 0 {
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"status": bson.M{"$eq": status},
		})
	}

	//代理类型
	switch invitedid {
	case 0:
		//总代理
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"invitedid": bson.M{"$eq": 0},
		})
	case 1:
		//子代理
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"invitedid": bson.M{"$ne": 0},
		})
	default:
		//子代理列表
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"invitedid": bson.M{"$eq": invitedid},
		})
	}

	start_time := ctx.Query("start")
	if start_time != "" {
		start,_ := time.Parse("2006-01-02", start_time)
		ctx.Data["start_time"] = start_time
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"requesttime": bson.M{"$gte": start},
		})
	}
	end_time := ctx.Query("end")
	if end_time != "" {
		end,_ := time.Parse("2006-01-02", end_time)
		ctx.Data["end_time"] = end_time
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"requesttime": bson.M{"$lt": end.AddDate(0,0,1)},
		})
	}

	list := []*agentModel.ApplyRecord{}
	_,count := db.C(tableName.DBT_AGENT_APPLY_LOG).Page(query, &list, "-requesttime", page, 10)
	ctx.Data["status"] = status
	ctx.Data["invited"] = invitedid
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
	types := ctx.QueryInt("types")
	var status exchangeService.ExchangeState
	switch types {
	case 1:
		status = exchangeService.PROCESS_TRUE
	case 2:
		status = exchangeService.PROCESS_TRUE
	case 3:
		status = exchangeService.PROCESS_TRUE
	case 4:
		status = exchangeService.PROCESS_FALSE
	}
	if len(id) != 24 || status <= 0 || status > 4 {
		ctx.Ajax(-1, "参数错误", nil)
		return
	}
	row := agentModel.GetApplyRecordById(id)
	if row == nil {
		ctx.Ajax(-2, "切换状态失败！", nil)
		return
	}
	row.Status = status
	row.ProcessTime = time.Now()
	err := row.Save()
	if err != nil {
		ctx.Ajax(-3, "切换状态失败！", nil)
		return
	}
	//切换为申请成功,则插入代理信息表
	if row.Status == exchangeService.PROCESS_TRUE {
		user_info := userService.GetUserById(row.UserId)
		new_agent_info := agentModel.AgentInfo{
			UserId: row.UserId,
			NickName: user_info.GetNickName(),
			RealName: row.Name,
			Phone: row.Phone,
			OpenId: user_info.GetOpenId(),
			UnionId: user_info.GetUnionId(),
			RootId: 0,
			Pid: row.InvitedId,
			Level: 1,
			Type: agentModel.AGENT_TYPE_2,
		}
		switch types {
		case 1:
			new_agent_info.Level = 1
			new_agent_info.RootId = 0
			new_agent_info.Type = agentModel.AGENT_TYPE_1
		case 2:
			new_agent_info.Level = 1
			new_agent_info.RootId = 0
			new_agent_info.Type = agentModel.AGENT_TYPE_2
		case 3:
			parent_info := agentModel.GetAgentInfoById(row.InvitedId)
			new_agent_info.RootId = parent_info.RootId
			new_agent_info.Level = parent_info.Level + 1
			new_agent_info.Type = agentModel.AGENT_TYPE_3
		}
		new_agent_info.Insert()
	}
	ctx.Ajax(1,"切换状态成功！", nil)
}
