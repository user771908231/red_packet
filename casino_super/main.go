package main

import (
	"gopkg.in/macaron.v1"
	"casino_common/common/sys"
	"os"
	"casino_super/conf"
	"casino_super/conf/config"
	"time"
	"casino_common/proto/ddproto"
	"casino_super/routers"
	"casino_super/modules"
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

func main() {
	m := macaron.New()
	//注册模板
	m.Use(macaron.Renderer(macaron.RenderOptions{Directory: "templates"}))
	//注册静态目录
	m.Use(macaron.Static("public"))

	//注册Alert模块
	m.Use(func(ctx *macaron.Context) {
		ctx.Map(modules.Alert{Contex:ctx})
	})

	//注册路由
	routers.Regist(m)

	m.NotFound(func(alert modules.Alert) {
		alert.Error("对不起未找到该页面！", "/admin", 3)
	})

	m.Run(conf.Server.HttpIp, conf.Server.HttpPort)

}
