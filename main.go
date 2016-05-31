package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"server/conf"
	"server/game"
	"server/gate"
	"server/login"
	"github.com/name5566/leaf/log"
)

func main() {
	log.Debug("开始执行 main 函数");

	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath

	log.Debug("lconf.ConnAddrs %s",lconf.ConnAddrs)

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
