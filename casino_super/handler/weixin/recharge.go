package weixin

import (
	"casino_super/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"casino_super/model/agentModel"
	"gopkg.in/mgo.v2/bson"
	"time"
	"math"
	"casino_super/model/weixinModel"
)

//充值列表
func RechargeListHandler(ctx *modules.Context) {
	list := []agentModel.Goods{}
	db.C(tableName.DBT_AGENT_GOODS).FindAll(bson.M{}, &list)
	ctx.Data["goods_list"] = list
	ctx.HTML(200, "weixin/agent/recharge")
}

//充值，订单确认页
func RechargeDoneHandler(ctx *modules.Context) {
	goods_id := ctx.QueryInt("goods_id")
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	//生成订单记录
	new_order := agentModel.AddNewRechargeLog(agent_id, int32(goods_id), 1)

	ctx.Data["order"] = new_order
	ctx.HTML(200, "weixin/agent/recharge_done")
}

//异步返回发起交易所需的参数
func RechargeAjaxWxTradeDataHandler(ctx *modules.Context) {
	order_id := ctx.Query("order_id")
	if len(order_id) != 24 {
		ctx.Ajax(-4, "订单号错误！", nil)
		return
	}
	order := agentModel.GetRechargeOrderId(order_id)
	wx_info := ctx.IsWxLogin()
	if order == nil {
		ctx.Ajax(-1, "订单号不存在！", nil)
		return
	}
	if wx_info == nil {
		ctx.Ajax(-2, "请先登录！", nil)
		return
	}
	//调用统一下单接口
	res,err := weixinModel.GetUnifiedOrderResponse(order.Id.Hex(), order.Money, order.Detail, ctx.RemoteAddr(), wx_info.OpenId)
	if err != nil {
		ctx.Ajax(-3, "统一下单失败！", nil)
		return
	}
	//返回发起交易需要的数据
	trade_data := weixinModel.GetTradeData(res.PrepayId)
	ctx.Ajax(1, "发起交易成功！", trade_data)
	return
}

//出售记录
func RechargeLogHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}
	agent := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(agent.UnionId)
	list := []agentModel.RechargeLog{}
	query := bson.M{
		"$and": []bson.M{
			bson.M{"agentid": bson.M{"$eq": agent_id}},
		},
	}
	start_time := ctx.Query("start")
	if start_time != "" {
		start,_ := time.Parse("2006-01-02", start_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"addtime": bson.M{"$gte": start},
		})
	}
	end_time := ctx.Query("end")
	if end_time != "" {
		end,_ := time.Parse("2006-01-02", end_time)
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"addtime": bson.M{"$lt": end.AddDate(0,0,1)},
		})
	}

	_, count := db.C(tableName.DBT_AGENT_RECHARGE_LOG).Page(query, &list, "-addtime", page, 10)

	data := bson.M{
		"list": list,
		"page": bson.M{
			"count":      count,
			"list_count": len(list),
			"limit":      10,
			"page":       page,
			"page_count": math.Ceil(float64(count) / float64(10)),
		},
	}

	ctx.Data["Logs"] = data
	ctx.Data["start_time"] = start_time
	ctx.Data["end_time"] = end_time
	ctx.HTML(200, "weixin/agent/recharge_log")
}
