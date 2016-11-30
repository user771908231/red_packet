package test

import (
	"casino_common/common/sys"
	"os"
	"casino_majiang/conf"
	"casino_majiang/conf/config"
)

func init() {
	e := sys.SysInit(
		conf.Server.RedisAddr,
		"test",
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