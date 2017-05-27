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
	"github.com/go-macaron/session"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/cache"
	"casino_common/common/service/pushService"
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
		config.SUPER_DBNAM,
		[]string{
			config.DBT_SUPER_LOGS,
			config.DBT_T_TH_GAMENUMBER_SEQ,
			config.DB_USER_SEQ,
		})
	//判断初始化是否成功
	if err != nil {
		os.Exit(-1)        //推出系统
	}

	//初始化pushService
	pushService.PoolInit(conf.Server.HallTcpAddr)

	time.Sleep(time.Second * 3)        //初始化3秒之后启动程序
}

func main() {
	m := macaron.Classic()
	//注册模板
	m.Use(macaron.Renderer(macaron.RenderOptions{Directory:"templates",IndentJSON:true}))
	//注册Session
	m.Use(session.Sessioner())
	//验证码依赖缓存组件
	m.Use(cache.Cacher())
	//验证码
	m.Use(captcha.Captchaer(captcha.Options{
		FieldCaptchaName: "captcha",
		ChallengeNums:4,
	}))
	//注册Context
	m.Use(func(ctx *macaron.Context, session session.Store) {
		ctx.Map(&modules.Context{Context:ctx,Session:session})
	})

	//注册路由
	routers.Regist(m)

	m.NotFound(func(ctx *modules.Context) {
		ctx.Error("对不起未找到该页面！", "", 0)
	})

	m.Run(conf.Server.HttpIp, conf.Server.HttpPort)

}
