package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"os"
	"casino_doudizhu/conf/config"
	"casino_doudizhu/conf"
	"casino_doudizhu/service/webserver"
	"casino_doudizhu/game"
	"casino_doudizhu/gate"
	"casino_doudizhu/login"
	"time"
)

func init() {
	//初始化系统
	e := config.InitConfig()
	if e.IsError() {
		//log.Error("config init failed.", e)
		os.Exit(-1)
	}

	time.Sleep(time.Second * 3)        //初始化3秒之后启动程序
}

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath

	//后台管理
	go func() {
		webserver.InitCms()
	}()

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
