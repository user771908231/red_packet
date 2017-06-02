package routers

import (
	"casino_server/casino_testtools/handler/game"
	"gopkg.in/macaron.v1"
)

//注册路由
func Regist(m *macaron.Macaron) {
	//日志
	m.Group("/game", func() {
		m.Get("", game.GameTest)
		m.Post("/edit", game.GameEdit)
	})

}
