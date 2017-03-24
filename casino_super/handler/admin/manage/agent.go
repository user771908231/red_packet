package manage

import (
	"casino_super/modules"
	"gopkg.in/mgo.v2/bson"
	"casino_super/model/agentModel"
	"casino_common/common/consts/tableName"
	"casino_common/utils/db"
	"math"
	"fmt"
)

//代理商管理
func AgentListHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	pid := ctx.QueryInt("pid")

	agent_info := agentModel.GetAgentInfoById(uint32(pid))

	if page == 0 {
		page = 1
	}

	//查询
	query := bson.M{
		"$and": []bson.M{
			bson.M{
				"pid": pid,
			},
		},
	}

	list := []*agentModel.AgentInfo{}
	_,count := db.C(tableName.DBT_AGENT_INFO).Page(query, &list, "", page, 10)

	ctx.Data["pid"] = pid
	ctx.Data["list"] = list
	ctx.Data["agent_info"] = agent_info
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	ctx.HTML(200, "admin/manage/agent/index")
}


//删除
func AgentDelHandler(ctx *modules.Context) {
	id := ctx.QueryInt("id")
	count,_ := db.C(tableName.DBT_AGENT_INFO).Count(bson.M{
		"pid": id,
	})
	if count > 0 {
		ctx.Ajax(-1, fmt.Sprintf("该代理商有%d个子代理，请先删除该代理下的所有子代理再执行此操作！", count), nil)
		return
	}
	err := db.C(tableName.DBT_AGENT_INFO).Remove(bson.M{
		"userid": id,
	})
	if err != nil {
		ctx.Ajax(-2, "删除失败！", nil)
	}
	ctx.Ajax(1, "删除成功！", nil)
}
