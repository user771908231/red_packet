package weixin

import (
	"casino_redpack/modules"
	"casino_common/common/model/agentModel"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"time"
	"casino_common/common/consts/tableName"
	"math"
	"casino_common/utils/numUtils"
	"casino_common/common/log"
	"casino_common/common/service/exchangeService"
)

//房费返利记录
func CoinFeeRebateLogHandler(ctx *modules.Context)  {
	agent := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(agent.UnionId)
	agent_info := agentModel.GetAgentInfoById(agent_id)

	rebate_type := ctx.Query("rebate_type")
	if rebate_type == "" {
		rebate_type = "rebate_log"
	}

	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}

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

	switch rebate_type {
	case "rebate_log":
		list := []agentModel.CoinFeeRebateLog{}
		_, count := db.Log(tableName.DBT_AGENT_COIN_FEE_REBATE_LOG).Page(query, &list, "-time", page, 10)
		ctx.Data["list"] = list
		ctx.Data["page"] = bson.M{
			"count":      count,
			"list_count": len(list),
			"limit":      10,
			"page":       page,
			"page_count": math.Ceil(float64(count) / float64(10)),
		}
	case "withdraw_log":
		list := []agentModel.RebateCoinWithdrawLog{}
		_, count := db.C(tableName.DBT_AGENT_COIN_FEE_WITHDRAW_LOG).Page(query, &list, "-time", page, 10)
		ctx.Data["list"] = list
		ctx.Data["page"] = bson.M{
			"count":      count,
			"list_count": len(list),
			"limit":      10,
			"page":       page,
			"page_count": math.Ceil(float64(count) / float64(10)),
		}
	}

	ctx.Data["start_time"] = start_time
	ctx.Data["end_time"] = end_time
	ctx.Data["agent_info"] = agent_info
	ctx.Data["rebate_type"] = rebate_type
	ctx.Data["all_rebate"] = numUtils.Float64Format(agent_info.CoinFeeRebate)
	ctx.Data["yesterday_rebate"] = agentModel.GetAgentCoinFeeRebateYesterday(agent_id)

	ctx.HTML(200, "weixin/agent/coin_fee_rebate_log")
}

//提现POST
func CoinFeeRebateWithdrawPost(ctx *modules.Context) {
	agent := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(agent.UnionId)
	agent_info := agentModel.GetAgentInfoById(agent_id)

	withdraw_amount := ctx.QueryFloat64("withdraw")

	log.T("用户%d发起提现申请，当前余额：%v 提现金额%v", agent_id, agent_info.CoinFeeRebate, withdraw_amount)
	//验证数值
	if withdraw_amount <= 0 || withdraw_amount - float64(int64(withdraw_amount)) != 0 || int64(withdraw_amount) % 100 != 0 {
		log.E("提现金额异常，用户id %v 提现金额%v", agent_id, withdraw_amount)
		ctx.Ajax(-1, "提现金额必须大于等于100，且必须为100的整倍数。", nil)
		return
	}

	if agent_info.CoinFeeRebate < withdraw_amount {
		ctx.Ajax(-2, "余额不足，提现失败！", nil)
		return
	}

	//提现成功
	err := agent_info.Inc("coinfeerebate", -withdraw_amount)

	if err == nil{
		//插入记录
		agentModel.RebateCoinWithdrawLog{
			AgentId: agent_id,
			Amount: withdraw_amount,
			Status: exchangeService.PROCESS_ING,
		}.Insert()

		ctx.Ajax(1, "提现成功！", nil)
	}else {
		log.E("提现失败：Err:%v", err)
		ctx.Ajax(-3, "提现失败,请联系管理员。Err:"+err.Error(), nil)
	}
}
