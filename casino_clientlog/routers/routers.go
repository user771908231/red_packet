package routers

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
	"casino_clientlog/model/logDao"
	"casino_clientlog/handler/logHandler"
	"casino_clientlog/modules"
	"casino_clientlog/handler/deskHandler"
)

//注册路由
func Regist(m *macaron.Macaron) {
	//日志
	m.Post("/log", binding.Json(logDao.ReqLog{}), logHandler.Post)
	m.Delete("/logs", binding.Json(logHandler.CodeValidate{}), logHandler.Delete)
	m.Get("/logs/:page", logHandler.Get)
	m.Get("/logs", logHandler.Get)

	//根据房号查战绩
	m.Get("/desk_users", deskHandler.FindUsersHandler)

	//首页
	m.Get("/", func(ctx *modules.Context) {
		ctx.Redirect("/logs", 302)
	})
}
