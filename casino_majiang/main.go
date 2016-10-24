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
	"casino_majiang/service/webserver"
)

func init() {
	//初始化系统
	e := config.InitConfig(false)
	if e.IsError() {
		//log.Error("config init failed.", e)
		os.Exit(-1)
	}
}

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath

	//后台管理
	go func() {
		//log.T("web start...")
		webserver.InitCms()
	}()

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
