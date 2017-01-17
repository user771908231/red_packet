package main

import (
	"gopkg.in/macaron.v1"
	"casino_super/handler/logHandler"
	"github.com/go-macaron/binding"
	"casino_super/model/logDao"
	"casino_common/common/sys"
	"os"
	"casino_super/conf"
	"casino_super/conf/config"
	"time"
	"casino_common/proto/ddproto"
)


func init() {
	//初始化系统
	err := sys.SysInit(
		int32(ddproto.COMMON_ENUM_RELEASETAG_R_PRO),
		conf.Server.ProdMod,
		conf.Server.RedisAddr,
		"test",
		conf.Server.LogPath,
		"super",
		conf.Server.MongoIp,
		conf.Server.MongoPort,
		config.SUPER_DBNAM,
		config.DB_ENSURECOUNTER_KEY,
		[]string{
			config.DBT_SUPER_LOGS,
			config.DBT_T_TH_GAMENUMBER_SEQ})
	//判断初始化是否成功
	if err != nil {
		os.Exit(-1)        //推出系统
	}

	time.Sleep(time.Second * 3)        //初始化3秒之后启动程序
}
//
func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())        //使用模板

	m.Use(macaron.Static("public/assets"))

	//log upload interface
	m.Post("log", binding.Json(logDao.ReqLog{}), logHandler.Post)

	m.Delete("logs", binding.Json(logHandler.CodeValidate{}), logHandler.Delete)

	m.Get("logs/:page", logHandler.Get)
	m.Get("logs", logHandler.Get)

	m.NotFound(func() string {
		return "not found 233..."
	})
	//launch server
	m.Run(conf.Server.HttpIp, conf.Server.HttpPort)


}