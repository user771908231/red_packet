package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"casino_common/common/log"
	"casino_templet/conf"
	"casino_templet/game"
	"casino_templet/gate"
	"casino_templet/login"
)

func init() {
	log.InitLogger("log", "as")
}

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.ConsolePort = 3333
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
