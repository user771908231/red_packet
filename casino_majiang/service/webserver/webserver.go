package webserver

import (
	"gopkg.in/macaron.v1"
	"casino_majiang/service/webserver/handler/mjDeskHandler"
)

type Validate struct {
	code string
}

func Run() {
	macaron.Env = macaron.PROD        //设置环境变量
	m := macaron.Classic()
	m.Use(macaron.Renderer())        //使用模板

	m.Use(macaron.Renderer(macaron.RenderOptions{
		// 模板文件目录，默认为 "templates"
		Directory: "templates/majiang",
	}))
	//routers
	m.Get("mjdesks", mjDeskHandler.Get) //get
	m.Get("mjdesk-users/:id", mjDeskHandler.GetUsers)

	m.NotFound(func() string {
		return "not found 233..."
	})

	//launch server
	m.Run()
}