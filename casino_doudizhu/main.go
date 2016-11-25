package main

import (
	"github.com/name5566/leaf"
	"casino_doudizhu/service/webserver"
	"casino_doudizhu/game"
	"casino_doudizhu/gate"
	"casino_doudizhu/login"
	"time"
	"casino_common/common/sys"
	"casino_doudizhu/conf"
	"casino_doudizhu/conf/config"
)

func init() {
	//初始化系统
	sys.SysInit(conf.Server.RedisAddr,
		"dodizhu",
		conf.Server.MongoIp,
		conf.Server.MongoPort,
		config.MJ_DBNAM,
		config.DB_ENSURECOUNTER_KEY,
		[]string{config.DBT_DDZ_DESK, config.DBT_T_TH_GAMENUMBER_SEQ})
	time.Sleep(time.Second * 3)        //初始化3秒之后启动程序
}

func main() {
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
