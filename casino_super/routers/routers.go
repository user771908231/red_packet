package routers

import (
	"gopkg.in/macaron.v1"
	"casino_super/handler/admin"
	"github.com/go-macaron/binding"
	"casino_super/model/logDao"
	"casino_super/handler/logHandler"
	"casino_super/modules"
	"casino_super/handler/admin/manage"
	"casino_super/handler/weixin"
	"casino_super/model/weixinModel"
)

//注册路由
func Regist(m *macaron.Macaron) {
	//日志
	m.Post("/log", binding.Json(logDao.ReqLog{}), logHandler.Post)
	m.Delete("/logs", binding.Json(logHandler.CodeValidate{}), logHandler.Delete)
	m.Get("/logs/:page", logHandler.Get)
	m.Get("/logs", logHandler.Get)

	//后台
	m.Group("/admin", func() {
		//需要登录
		m.Group("", func() {
			m.Get("/", admin.IndexHandler)
			m.Get("/main", admin.MainHandler)
			m.Get("/sign", admin.SignHandler)
		},admin.NeedLogin)
		m.Get("/login", admin.LoginHandler)
		m.Get("/logout", admin.LogoutHandler)
		m.Post("/login", admin.NeedCaptcha, binding.Bind(admin.LoginForm{}), admin.LoginPostHandler)

		//管理页
		m.Group("/manage", func() {
			//用户相关
			m.Group("/user", func() {
				//所有用户
				m.Get("/", func(ctx *modules.Context) {ctx.HTML(200, "admin/manage/user/index")})
				m.Get("/list", manage.UserListHandler)
				m.Get("/del/:id", manage.DelUserHandler)
				m.Post("/update", binding.BindIgnErr(manage.UserUpdateForm{}), manage.UpdateUserHandler)
				m.Post("/recharge", binding.BindIgnErr(manage.RechargeForm{}),manage.RechargeHandler)
			})
			//红包兑换相关
			m.Group("/exchange", func() {
				m.Get("/", manage.ExchangeListHandler)
				m.Get("/switch", manage.ExchangeSwitchState)
			})
			m.Group("/apply", func() {
				m.Get("/", manage.ApplyListHandler)
				m.Get("/switch", manage.ApplySwitchState)
			})
		}, admin.NeedLogin)
	}, admin.ShowPanel)

	//微信
	m.Group("/weixin", func() {
		m.Get("/oauth/login", weixinModel.OauthLogin)
		m.Get("/oauth/callback", weixinModel.OauthCallBack)
		//需要微信登录
		m.Group("/agent", func() {
			m.Get("/", weixin.MainHandler)
			//代理个人信息
			m.Get("/info", weixin.AgentInfoHandler)
			//充值
			m.Group("/recharge", func() {
				m.Get("/", weixin.RechargeListHandler)
				m.Get("/done", weixin.RechargeDoneHandler)
				m.Get("/wx_pay", weixin.RechargeAjaxWxTradeDataHandler)
				m.Get("/log", weixin.RechargeLogHandler)
			}, weixin.NeedIsRootAgent)
			//销售
			m.Group("/sales", func() {
				m.Get("/", weixin.SalesIndexHandler)
				m.Post("/",binding.BindIgnErr(weixin.SalesForm{}), weixin.SalesToUserHandler)
				m.Get("/log", weixin.SalesLogHandler)
			})
			//我的客户
			m.Get("/customers", weixin.CustomersListHandler)
			m.Group("/apply", func() {
				m.Get("/log", weixin.ApplyLogHandler)
				m.Get("/switch", weixin.ApplySwitchState)
			})

			//我的下线代理
			m.Get("/my_agents", weixin.MyAgentsHandler)
			//返利记录
			m.Group("/rebate", func() {
				m.Get("/log", weixin.RebateLogHandler)
				m.Get("/check", weixin.CheckRebateHandler)
			})
			//登录-登出
			m.Group("/user", func() {
				m.Get("/login", func(ctx *modules.Context) {
					ctx.Redirect("/weixin/oauth/login", 200)
				})
				m.Get("/logout", func(ctx *modules.Context) {
					ctx.Session.Delete("wx_user")
					ctx.Success("退出成功！", "/weixin/", 3)
				})
			})
		}, weixin.NeedWxLogin, weixin.NeedIsAgent)
	})
	//代理申请
	m.Get("/weixin/agent/apply", weixin.NeedWxLogin, weixin.ApplyHandler)
	m.Post("/weixin/agent/apply", weixin.NeedWxLogin, admin.NeedCaptcha, binding.BindIgnErr(weixin.ApplyForm{}), weixin.ApplyPostHandler)
	//微信充值回调
	m.Any("/mp/pay/callback", weixinModel.WxNotifyHandler)
	//微信领红包
	m.Get("/weixin/get_redpack", weixin.NeedWxLogin, weixin.NeedIsGamer, weixin.GetRedPackHandler)

	//首页
	m.Get("/", func(ctx *modules.Context) {
		//ctx.Success("即将跳转至后台！", "/admin", 3)
		//ctx.Redirect("/admin", 302)
		ctx.Redirect("/weixin/agent", 302)
	})
}
