package routers

import (
	"gopkg.in/macaron.v1"
	"casino_hotupdate/modules"
	"casino_hotupdate/handler/hotUpdate"
	"github.com/go-macaron/binding"
)

//注册路由
func Regist(m *macaron.Macaron) {

	//热更新
	m.Group("/hotupdate", func() {
		//版本列表
		m.Get("/list", hotUpdate.ListHandler)

		//上传
		m.Post("/upload", binding.MultipartForm(hotUpdate.UploadForm{}), hotUpdate.UploadHandler)
	})

	//首页
	m.Get("/", func(ctx *modules.Context) {
		ctx.Redirect("/hotupdate/list", 302)
	})
}
