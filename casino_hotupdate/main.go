package main

import (
	"casino_hotupdate/conf"
	"casino_hotupdate/modules"
	"casino_hotupdate/routers"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"casino_common/common/db"
	"casino_common/utils/redis"
)

func init() {
	//初始化系统
	db.InitMongoDb(conf.Server.MongoIp, conf.Server.MongoLogIp, "test", "id", []string{})
	data.InitRedis(conf.Server.RedisAddr, "test")
}

func main() {
	m := macaron.Classic()
	//注册模板
	m.Use(macaron.Renderer(macaron.RenderOptions{Directory: "templates", IndentJSON: true}))
	//注册Session
	m.Use(session.Sessioner())
	//验证码依赖缓存组件
	m.Use(cache.Cacher())
	//验证码
	m.Use(captcha.Captchaer(captcha.Options{
		FieldCaptchaName: "captcha",
		ChallengeNums:    4,
	}))
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
