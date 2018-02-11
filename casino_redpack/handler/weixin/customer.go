package weixin

import (
	"casino_redpack/modules"
	"casino_common/proto/ddproto"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/model/agentModel"
	"math"
	"fmt"
)

//我的客户列表
func CustomersListHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	user_list := []*ddproto.User{}

	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	//return
	//ctx.Dump(wx_info)
	//ctx.Dump("agentid:", agent_id)

	count := 0
	if agent_id >0 {
		_,count = db.C(tableName.DBT_T_USER).Page(bson.M{
			"agentid": agent_id,
		}, &user_list, "", page, 10)
	}

	//ctx.Ajax(1, "", data)
	list := []bson.M{}
	for _,user := range user_list {
		week_sales_count := agentModel.GetAgentThisWeekOneUserSalesCount(agent_id, user.GetId())
		all_sales_count := agentModel.GetAgentUserSalesCount(agent_id, user.GetId())
		list = append(list, bson.M{
			"Id": user.GetId(),
			"NickName": user.GetNickName(),
			"AllSalesCount": all_sales_count,
			"WeekSalesCount": week_sales_count,
		})
	}
	ctx.Data["list"] = list
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(user_list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	ctx.HTML(200, "weixin/agent/customers")
}

//我的客户列表
func CustomersCoinListHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	user_list := []*ddproto.User{}

	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	//return
	//ctx.Dump(wx_info)
	//ctx.Dump("agentid:", agent_id)

	count := 0
	if agent_id >0 {
		_,count = db.C(tableName.DBT_T_USER).Page(bson.M{
			"invcode": fmt.Sprintf("%v", agent_id),
		}, &user_list, "", page, 10)
	}

	//ctx.Ajax(1, "", data)
	list := []bson.M{}
	for _,user := range user_list {
		week_sales_count := agentModel.GetAgentCustomCoinFeeRebateWeekSum(agent_id, user.GetId())
		all_sales_count := agentModel.GetAgentCustomCoinFeeRebateAllSum(agent_id, user.GetId())
		list = append(list, bson.M{
			"Id": user.GetId(),
			"NickName": user.GetNickName(),
			"AllSalesCount": all_sales_count,
			"WeekSalesCount": week_sales_count,
		})
	}
	ctx.Data["list"] = list
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(user_list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	ctx.HTML(200, "weixin/agent/customers_coin")
}

