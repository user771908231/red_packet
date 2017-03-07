package weixin

import (
	"casino_super/modules"
	"casino_common/proto/ddproto"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_super/model/agentModel"
	"math"
)

//我的客户列表
func CustomersListHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	list := []*ddproto.User{}

	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	//ctx.Dump(wx_info)
	//ctx.Dump("agentid:", agent_id)

	count := 0
	if agent_id >0 {
		_,count = db.C(tableName.DBT_T_USER).Page(bson.M{
			"agentid": agent_id,
		}, &list, "", page, 10)
	}

	//ctx.Ajax(1, "", data)
	ctx.Data["list"] = list
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	ctx.HTML(200, "weixin/agent/customers")
}
