package main

import (
	"gopkg.in/macaron.v1"
	"casino_super/handler/userHandler"
)

//
func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())        //使用模板

	m.Get("/", userHandler.Users)
	m.Get("/recharge", userHandler.Recharge)
	m.Post("/rechargePost", userHandler.RechargePost)
	m.Run()
}