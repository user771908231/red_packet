package routers

import (
	"gopkg.in/macaron.v1"
	"casino_admin/handler/admin"
	"github.com/go-macaron/binding"
	"casino_admin/model/logDao"
	"casino_admin/handler/logHandler"
	"casino_admin/modules"
	"casino_admin/handler/admin/manage"
	"casino_admin/handler/admin/config"
	"casino_common/common/model"
	"casino_common/common/service/configService"
	"casino_admin/handler/game"
)

//注册路由
func Regist(m *macaron.Macaron) {
	//日志
	m.Post("/log", binding.Json(logDao.ReqLog{}), logHandler.Post)
	m.Delete("/logs", binding.Json(logHandler.CodeValidate{}), logHandler.Delete)
	m.Get("/logs/:page", logHandler.Get)
	m.Get("/logs", logHandler.Get)
	m.Group("/game", func() {
		m.Get("",game.GameTest)
		m.Post("/edit",game.GameEdit)
	})

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
				m.Post("/edit",config.GoodsEditPost)
				m.Post("/editUpdate", config.GoodsEditUpdate)
				m.Post("/insert", binding.Bind(model.T_Goods_Row{}), config.GoodsInsertPost)
				m.Get("/remove", config.GoodsRemoveHnadler)
			})
			//任务信息配置
			m.Group("/task", func() {
				m.Get("/list", config.TaskListHandler)
				m.Post("/edit", config.TaskEditPost)
			})
			//游戏配置
			m.Group("/game", func() {
				m.Get("/list", config.GameConfigListHandler)
				m.Get("/listAll", config.GameConfigList)
				m.Get("/listLogin", config.GameConfigLogin)
				m.Post("/edit", config.GameConfigEdit)
				m.Post("/editUpdate", config.GameConfigUpdate)
				m.Post("/editUpdateLogin", config.GameConfigUpdateLogin)
				m.Post("/editLogin", config.GameConfigEditLogin)
				m.Get("/add", config.GameConfigAddHandler)
				m.Post("/addServerInfo",binding.Bind(configService.LoginServerInfo{}), config.GameServerInfoAddPost)
				m.Get("/gameList", config.GameListHandler)
				m.Get("/gameListEdit", config.GameListEditHandler)

			})

			//Notice配置
			m.Group("/notice", func() {
				m.Get("/list", config.NoticeListHandler)
				m.Post("/edit", binding.Bind(config.NoticeForm{}), config.NoticeEditHandler)
			})
		})

		//数据分析
		m.Group("/data", func() {
			m.Get("/atHome",admin.AtHome)
			m.Post("/atHomeList",admin.AtHomeList)
			m.Get("/onlineStatic",admin.OnlineStatic)
			m.Post("/onlineStaticList",admin.OnlineStaticList)
			m.Get("/roomCard",admin.RoomCard)
			m.Post("/roomCardOne",admin.RoomCardOne)
			m.Get("/roomCardDay",admin.RoomCardDay)
			m.Post("/roomCardDayOne",admin.RoomCardDayOne)
			m.Get("/user_reg", manage.UserRegAllHandler)
			m.Post("/user_reg_one", manage.UserRegOneHandler)
			m.Get("/user_active", manage.UserActiveAllHandler)
			m.Post("/user_active_one", manage.UserActiveOneHandler)

		})
	}, admin.ShowPanel)

	//首页
	m.Get("/", func(ctx *modules.Context) {
		ctx.Redirect("/admin", 302)
	})
}
