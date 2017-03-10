package weixin

import (
	"casino_super/modules"
	"casino_super/model/agentModel"
)

//需要微信登录验证
func NeedWxLogin(ctx *modules.Context) {
	wx_info := ctx.IsWxLogin()
	if wx_info == nil {
		ctx.Redirect("/weixin/oauth/login", 302)
		return
	}
	ctx.Data["wx_user"] = wx_info
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
