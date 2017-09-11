package routers

import (
	"casino_testtools/handler/game"
	"gopkg.in/macaron.v1"
)

//注册路由
func Regist(m *macaron.Macaron) {
	//日志
	m.Group("/game", func() {
		m.Get("/pdk", game.GameTest)
		m.Get("/edit", game.GameEdit)
	})

}
