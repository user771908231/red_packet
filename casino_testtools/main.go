package main

import (
	"casino_testtools/modules"
	"casino_testtools/routers"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"casino_common/common/sys"
	"casino_admin/conf"
	"casino_admin/conf/config"
	"casino_common/proto/ddproto"
	"os"
)

func init() {
	//初始化系统
	err := sys.SysInit(
		int32(ddproto.COMMON_ENUM_RELEASETAG_R_PRO),
		conf.Server.ProdMod,
		conf.Server.RedisAddr,
		"test",
		conf.Server.RedisPwd,
		conf.Server.LogPath,
		"super",
		conf.Server.LogFileSize,
		conf.Server.LogFileCount,
		conf.Server.MongoIp,
		conf.Server.MongoLogIp,
		config.SUPER_DBNAM,
		[]string{
			config.DBT_SUPER_LOGS,
			config.DB_USER_SEQ,
		})

	//判断初始化是否成功
	if err != nil {
		os.Exit(-1) //推出系统
	}
}

func main() {
	m := macaron.Classic()
	//注册模板
	m.Use(macaron.Renderer(macaron.RenderOptions{Directory: "templates", IndentJSON: true}))
	//注册Session
	m.Use(session.Sessioner())
	//注册Context
	m.Use(func(ctx *macaron.Context, session session.Store) {
		ctx.Map(&modules.Context{Context: ctx, Session: session})
	})

	//注册路由
	routers.Regist(m)

	m.NotFound(func(ctx *modules.Context) {
		ctx.Error("对不起未找到该页面！", "", 0)
	})

	m.Run(conf.Server.HttpIp, conf.Server.HttpPort)

}
