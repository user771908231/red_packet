package main

import (
	"gopkg.in/macaron.v1"
	"casino_super/handler/userHandler"
)

//
func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())        //使用模板

	//m.Get("/recharge", userHandler.Recharge)
	//m.Post("/rechargePost", userHandler.RechargePost)

	// view handler
	//m.Get("/views/userList", viewHandler.getUserList)
	m.Get("/", viewHandler.getUserList)


	/**
	**	RESTful APIs
	**	expect domain name: api.dongdian.com/
	**/

	//user resource
	m.Combo("/users").
		Get(userHandler.getAll). //get all
		Post(userHandler.post).
		Put(userHandler.put).
		Delete(userHandler.delete) //soft delete
	m.Get("/users/:id", userHandler.getOneById) //add get one
	m.Get("/users/:id", userHandler.getOneById) //add get one

	//launch server
	m.Run()
}