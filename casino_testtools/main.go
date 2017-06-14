package main

import (
	"casino_testtools/modules"
	"casino_testtools/routers"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func main() {
	m := macaron.Classic()
	//注册模板
	m.Use(macaron.Renderer(macaron.RenderOptions{Directory: "templates", IndentJSON: true}))
	//注册Session
	m.Use(session.Sessioner())
	//注册Context
	m.Use(func(ctx *macaron.Context, session session.Store) {
		ctx.Map(&modules.Context{Context: ctx, Session: session})
	})

	//注册路由
	routers.Regist(m)

	m.NotFound(func(ctx *modules.Context) {
		ctx.Error("对不起未找到该页面！", "", 0)
	})

	m.Run(9093)

}
