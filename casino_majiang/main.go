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
	"casino_common/common/sys"
)

func init() {
	//redisAddr string, redisName string, logName string, mongoIp string, mongoPort int, mongoName string, mongoSeqKey string, mongoSeqTables []string
	e := sys.SysInit(
		conf.Server.RedisAddr,
		"test",
		conf.Server.LogPath,
		"majaing",
		conf.Server.MongoIp,
		conf.Server.MongoPort,
		config.MJ_DBNAM,
		config.DB_ENSURECOUNTER_KEY,
		[]string{
			config.DBT_MJ_DESK,
			config.DBT_T_TH_GAMENUMBER_SEQ})
	//初始化系统
	if e != nil {
		os.Exit(-1)
	}
}

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath

	//后台管理
	go func() {
		webserver.Run()
	}()

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
