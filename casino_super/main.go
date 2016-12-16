package main

import (
	"gopkg.in/macaron.v1"
	"casino_super/handler/logHandler"
	"github.com/go-macaron/binding"
	"casino_super/model/logDao"
)

//
func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())        //使用模板

	//log upload interface
	m.Post("log", binding.Json(logDao.ReqLog{}), logHandler.Post)

	//m.Post("log", logHandler.Post)

	m.Get("logs", logHandler.Get)

	m.NotFound(func() string {
		return "not found 233..."
	})
	//launch server
	m.Run()


}