package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"casino_majiang/conf"
	"casino_majiang/game"
	"casino_majiang/gate"
	"casino_majiang/login"
	"os"
	"casino_majiang/conf/config"
	"casino_majiang/conf/log"
)

func init() {
	//初始化系统
	e := config.InitConfig(false)
	if e.IsError() {
		log.Error("config init failed.", e)
		os.Exit(-1)
	}
}

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
