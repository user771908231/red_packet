package routers

import (
	"gopkg.in/macaron.v1"
	"casino_game/handler/game"
)

//注册路由
func Regist(m *macaron.Macaron) {
	//日志
	m.Group("/game", func() {
		m.Get("",game.GameTest)
		m.Post("/edit",game.GameEdit)
	})


}
