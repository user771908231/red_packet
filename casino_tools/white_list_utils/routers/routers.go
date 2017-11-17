package routers

import (
	"gopkg.in/macaron.v1"
	"casino_tools/white_list_utils/modules"
	"casino_tools/white_list_utils/handler/admin"
	"github.com/go-macaron/binding"
	"casino_tools/white_list_utils/handler/admin/config"
)

//注册路由
func Regist(m *macaron.Macaron) {
	//后台
	m.Group("/admin", func() {
		m.Get("/login", admin.LoginHandler)
		m.Get("/logout", admin.LogoutHandler)
		m.Post("/login", admin.NeedCaptcha, binding.Bind(admin.LoginForm{}), admin.LoginPostHandler)
		//配置管理
		m.Group("/config", func() {
			m.Get("/white_list", config.WhiteListHandler)
			m.Get("/white_list_post", config.WhiteListPostHandler)
		}, admin.NeedLogin)
	}, admin.ShowPanel)

	//首页
	m.Get("/", func(ctx *modules.Context) {
		ctx.Redirect("/admin/config/white_list", 302)
	})
}
