package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"casino_majianagv2/conf"
	"casino_majianagv2/game"
	"casino_majianagv2/gate"
	"casino_majianagv2/login"
	"casino_common/common/sys"
	"casino_common/proto/ddproto"
	"casino_common/common/consts/tableName"
	"os"
)

func init() {
	e := sys.SysInit(
		int32(ddproto.COMMON_ENUM_RELEASETAG_R_PRO), //发布版本
		conf.Server.ProdMode,                        //开发模式
		conf.Server.RedisAddr,                       //redis地址
		"test",                                      //redis 数据哭名字
		conf.Server.LogPath,                         //日志路径
		"majaing",                                   //日志名字
		conf.Server.MongoIp,                         //mongo ip
		conf.Server.MongoPort,                       //mongo 端口
		"test",                                      //mongo 数据库名字
		tableName.DB_ENSURECOUNTER_KEY,              //统一的自增键名字
		//需要自增的表
		[]string{
			tableName.DBT_MJ_DESK,
			tableName.DBT_T_TH_GAMENUMBER_SEQ})
	//初始化系统
	if e != nil {
		os.Exit(-1)
	}
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
