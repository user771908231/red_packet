package routers

import (
	"gopkg.in/macaron.v1"
	"new_links/handler/admin"
	"new_links/handler/weixin"
	"new_links/handler/admin/links"
	"github.com/go-macaron/binding"
	"new_links/handler/admin/group"
	"new_links/handler/admin/keys"
)

//注册路由
func Regist(m *macaron.Macaron) {
	m.Group("", func() {
		//管理后台
		m.Group("/admin", func() {
			//主页
			m.Get("/", admin.IndexHandler)
			m.Get("/logout", admin.LoginOutHandler)
			//链接管理
			m.Group("/links", func() {
				m.Get("/",links.IndexHandler)
				m.Get("/add",links.AddHandler)
			})
			//关键词管理
			m.Group("/keys", func() {
				m.Get("/",keys.IndexHandler)
				m.Get("/add",keys.AddHandler)
				m.Get("/edit",keys.EditHandler)
				m.Post("/add",keys.PostAddHandler)
				m.Get("/status",keys.StatusHandler)
				m.Post("/update",keys.UpdateHandler)
				m.Get("/del",keys.DelHandler)
				m.Get("/upload",keys.Uploadhandler)
			})

			//分组管理
			m.Group("/group", func() {
				m.Get("/",group.IndexHandler)
				m.Get("/add",group.AddHandler)
				m.Post("/add",group.PostAddHandler)
				m.Get("/status",group.StatusHandler)
				m.Get("/del",group.DelHandler)
			})

		}, admin.NeedLogin(2))
		m.Post("/admin/login", admin.NeedCaptcha, binding.Bind(admin.LoginForm{}), admin.LoginPostHandler)
		//管理登录
		m.Get("/admin/login", admin.LoginHandler)



	}, weixin.RootMidware)
}
