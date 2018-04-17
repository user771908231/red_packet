package main

import (
	"gopkg.in/macaron.v1"
	"html/template"
	"github.com/go-macaron/session"
	"sendlinks/modules"
	"sendlinks/routers"
	"sendlinks/conf"
)

func main() {
	m := macaron.Classic()
	//注册模板
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "templates",
		IndentJSON: true,
		Funcs: []template.FuncMap{template.FuncMap{
			"add": func(a,b int) int {return a+b},
		}},
	}))
	//Session
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

	m.Run(conf.Server.HttpIp,conf.Server.HttpPort)
}