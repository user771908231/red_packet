package routers

import (
	"gopkg.in/macaron.v1"
	"casino_admin/handler/admin"
	"github.com/go-macaron/binding"
	"casino_admin/model/logDao"
	"casino_admin/handler/logHandler"
	"casino_admin/modules"
	"casino_admin/handler/admin/manage"
	"casino_admin/handler/weixin"
	"casino_admin/model/weixinModel"
	"casino_admin/handler/admin/config"
	"casino_common/common/model"
	"casino_admin/handler/qrLoginHandler"
	"casino_common/common/service/taskService/taskType"
	"casino_common/common/service/configService"
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
				m.Get("/", manage.UserIndexHnadler)
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
			//代理商申请审核
			m.Group("/apply", func() {
				m.Get("/", manage.ApplyListHandler)
				m.Get("/switch", manage.ApplySwitchState)
			})
			//代理商管理
			m.Group("/agent", func() {
				m.Get("/", manage.AgentListHandler)
				m.Get("/switch", manage.AgentSwitchHandler)
				m.Get("/del", manage.AgentDelHandler)
			})
		}, admin.NeedLogin)

		//配置管理
		m.Group("/config", func() {
			//商城商品配置
			m.Group("/goods", func() {
				m.Get("/list", config.GoodsListHandler)
				m.Post("/edit", binding.Bind(model.T_Goods_Row{}), config.GoodsEditPost)
				m.Post("/insert", binding.Bind(model.T_Goods_Row{}), config.GoodsInsertPost)
				m.Get("/remove", config.GoodsRemoveHnadler)
			})
			//任务信息配置
			m.Group("/task", func() {
				m.Get("/list", config.TaskListHandler)
				m.Post("/edit", binding.Bind(taskType.TaskInfo{}), config.TaskEditPost)
			})
			//游戏配置
			m.Group("/game", func() {
				m.Get("/list", config.GameConfigListHandler)
				m.Post("/edit", config.GameConfigEdit)
				m.Post("/editUpdate", config.GameConfigUpdate)
				m.Get("/add", config.GameConfigAddHandler)
				m.Post("/addServerInfo",binding.Bind(configService.LoginServerInfo{}), config.GameServerInfoAddPost)
			})
		})

		//数据分析
		m.Group("/data", func() {
			m.Get("/atHome",admin.AtHome)
			m.Post("/atHomeList",admin.AtHomeList)
			m.Get("/onlineStatic",admin.OnlineStatic)
			m.Post("/onlineStaticList",admin.OnlineStaticList)
		})
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

		//二维码登陆
		m.Group("/game", func() {
			m.Get("/qrlogin", qrLoginHandler.QrLoginHandler)
		}, weixin.NeedWxLogin)

	})
	//代理申请
	m.Get("/weixin/agent/apply", weixin.NeedWxLogin, weixin.ApplyHandler)
	m.Post("/weixin/agent/apply", weixin.NeedWxLogin, weixin.NeedIsGamer, admin.NeedCaptcha, binding.BindIgnErr(weixin.ApplyForm{}), weixin.ApplyPostHandler)
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
