package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"casino_doudizhu/service/webserver"
	"casino_doudizhu/game"
	"casino_doudizhu/gate"
	"casino_doudizhu/login"
	"time"
	"casino_common/common/sys"
	"casino_doudizhu/conf"
	"casino_doudizhu/conf/config"
	"os"
)

func init() {
	//初始化系统
	err := sys.SysInit(
		conf.Server.RedisAddr,
		"test",
		conf.Server.LogPath,
		"dodizhu",
		conf.Server.MongoIp,
		conf.Server.MongoPort,
		config.MJ_DBNAM,
		config.DB_ENSURECOUNTER_KEY,
		[]string{
			config.DBT_DDZ_DESK,
			config.DBT_T_TH_GAMENUMBER_SEQ})
	//判断初始化是否成功
	if err != nil {
		os.Exit(-1)        //推出系统
	}

	time.Sleep(time.Second * 3)        //初始化3秒之后启动程序
}

func main() {
	//leaf 的日志路径
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath

	//后台管理
	go func() {
		webserver.InitCms()
	}()

	//启动leaf
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
