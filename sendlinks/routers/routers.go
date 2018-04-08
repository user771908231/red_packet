package routers

import (

	"gopkg.in/macaron.v1"
	"sendlinks/handler/admin"
	"sendlinks/handler/admin/links"
	"sendlinks/handler/admin/grouping"
	"sendlinks/handler/admin/keys"
)

func Regist(m *macaron.Macaron) {
	m.Group("", func() {
		m.Group("/", func(){

		})
		//admin
		m.Group("/admin", func() {
				//后台首页
				m.Get("/",admin.IndexHandler)

				//管理员列表
				m.Get("/lists",admin.ListsHandler)
				//新建管理员
				m.Get("/add",admin.AddHandler)
				//管理员权限
				m.Get("/power",admin.PowerHandler)
				//权限分类
				m.Get("/classification",admin.ClassiFicationHandler)

			//链接管理
			m.Group("/link_list", func() {
				//链接列表首页
				m.Get("/",links.IndexHandler)

				//添加新的链接
				m.Get("/add",links.AddHandler)

				//查询链接
				m.Get("/select",links.SelectHandler)

				//修改链接
				m.Get("/uplate",links.UplateHandler)

				//删除链接
				m.Get("/delete",links.DeleteHandler)
			})
			//分组管理
			m.Group("/grouping", func() {
				//列表首页
				m.Get("/",grouping.IndexHandler)

				//添加
				m.Get("/add",grouping.AddHandler)

				//查询
				m.Get("/select",grouping.SelectHandler)

				//修改
				m.Get("/uplate",grouping.UplateHandler)

				//删除
				m.Get("/delete",grouping.DeleteHandler)
			})
			//关键词管理
			m.Group("/keys", func() {
				//列表首页
				m.Get("/",keys.IndexHandler)

				//添加
				m.Get("/add",keys.AddHandler)

				//查询
				m.Get("/select",keys.SelectHandler)

				//修改
				m.Get("/uplate",keys.UplateHandler)

				//删除
				m.Get("/delete",keys.DeleteHandler)
			})

		})
		m.Get("admin/login",admin.LoginHandler)
		m.Post("admin/login",admin.IsLoginHandler)
	})
}

