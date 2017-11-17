package weixin

import (
	"casino_super/modules"
	"casino_super/model/agentModel"
	"casino_common/common/service/exchangeService"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/model/weixinModel"
	"fmt"
)

//需要微信登录验证
func NeedWxLogin(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	if wx_info == nil {
		ctx.Session.Set("redirect", ctx.Req.RequestURI)
		ctx.Redirect("/weixin/oauth/login", 302)
		return
	}
	ctx.Data["wx_user"] = wx_info
}

//需要是游戏用户
func NeedIsGamer(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	//验证该微信是否已在游戏中注册
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	if agent_id == 0 {
		ctx.Error("您的微信账号还未与神经斗地主游戏关联，请先用手机登录游戏再访问本页面。","",0)
		return
	}
}

//需要是代理商
func NeedIsAgent(ctx *modules.Context)  {
	wx_info := ctx.IsWxLogin()
	//验证该微信是否已在游戏中注册
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	if agent_id == 0 {
		ctx.Error("您的微信账号还未与神经斗地主游戏关联，请先用手机登录游戏再刷新本页面。","",0)
		return
	}
	//验证是否为代理商
	if !agentModel.IsAgent(agent_id) {
		ctx.Error("您还不是代理商，请先填写申请表单。","/weixin/agent/apply",5)
		return
	}
}

//需要是一级代理
func NeedIsRootAgent(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	agent_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	agent_info := agentModel.GetAgentInfoById(agent_id)
	if agent_info == nil || agent_info.Level != 1 {
		ctx.Error("您没有权限访问本页面！", "", 0)
		return
	}
}

//领红包
func GetRedPackHandler(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	user_id := agentModel.GetUserIdByUnionId(wx_info.UnionId)
	rows_true := []*exchangeService.ExchangeRecord{}
	db.C(tableName.DBT_ADMIN_EXCHANGE_RECORD).FindAll(bson.M{
		"userid": user_id,
		"status": exchangeService.PROCESS_TRUE,
	}, &rows_true)
	if len(rows_true) == 0 {
		//查询审核中的红包数
		num_ing,_ := db.C(tableName.DBT_ADMIN_EXCHANGE_RECORD).Count(bson.M{
			"userid": user_id,
			"status": exchangeService.PROCESS_ING,
		})
		num_false,_ := db.C(tableName.DBT_ADMIN_EXCHANGE_RECORD).Count(bson.M{
			"userid": user_id,
			"status": exchangeService.PROCESS_FALSE,
		})
		if num_ing > 0 || num_false > 0 {
			ctx.Success(fmt.Sprintf("您有%d个兑换红包请求正在等待审核，%d个兑换红包请求被拒绝受理。", num_ing, num_false),"", 0)
			return
		}
		ctx.Error("您目前没有红包可供领取！", "", 0)
		return
	}else {
		//发放红包
		for _,item := range rows_true {
			err := weixinModel.SendRedPack(wx_info.OpenId, item.Money, item.Id.Hex())
			if err != nil {
				ctx.Error("抱歉，一个红包发送失败！请重试或联系管理员。","", 0)
				return
			}else {
				//设置状态为已发放
				item.Status = exchangeService.PROCESS_SENDED
				item.Save()
			}
		}
		ctx.Success(fmt.Sprintf("成功领取%d个红包！请返回公众号查看。", len(rows_true)), "", 0)
		return
	}
}
