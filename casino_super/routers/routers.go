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
		}, admin.NeedLogin)
	}, admin.ShowPanel)

	//微信
	m.Group("/weixin", func() {
		m.Get("/oath", weixin.Oath)
		m.Get("/user/info", weixin.UserInfo)
	})

	//首页
	m.Get("/", func(ctx *modules.Context) {
		//ctx.Success("即将跳转至后台！", "/admin", 3)
		ctx.Redirect("/admin", 302)
	})
}
