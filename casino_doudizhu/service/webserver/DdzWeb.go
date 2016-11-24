package webserver

import (
	"gopkg.in/macaron.v1"
	"casino_doudizhu/service/webserver/handler/ddzDeskHandler"
)

func InitCms() {

	m := macaron.Classic()
	m.Use(macaron.Renderer())        //使用模板

	m.Use(macaron.Renderer(macaron.RenderOptions{
		// 模板文件目录，默认为 "templates"
		Directory: "service/webserver/templates",
	}))
	//routers
	m.Get("ddzdesk", ddzDeskHandler.Get) //get
	m.Get("ddzdesk-users/:id", ddzDeskHandler.GetUsers)

	m.NotFound(func() string {
		return "not found 233..."
	})

	//launch server
	m.Run("0.0.0.0", 40001)
}
