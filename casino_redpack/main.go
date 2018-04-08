package main

import (
	"casino_redpack/conf"
	"casino_redpack/conf/config"
	"casino_redpack/modules"
	"casino_redpack/routers"
	"casino_common/common/sys"
	"casino_common/proto/ddproto"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"os"
	"casino_redpack/model/weixinModel"
	"html/template"
	"casino_common/common/service/rpcService"
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
		"weixin",
		conf.Server.LogFileSize,
		conf.Server.LogFileCount,
		conf.Server.MongoIp,
		conf.Server.MongoLogIp,
		config.SUPER_DBNAM,
		[]string{
			config.DBT_SUPER_LOGS,
			config.DB_USER_SEQ,
			"redpack_room_id",
			"redpack_redpack_id",
			config.ORDER_KEY_ID,
			config.WITHDRAWALS_KEY_ID,
		})

	//判断初始化是否成功
	if err != nil {
		os.Exit(-1) //推出系统
	}
}

func main() {
	m := macaron.Classic()
	//注册模板
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "templates",
		IndentJSON: true,
		Funcs: []template.FuncMap{template.FuncMap{
			"add": func(a,b int) int {return a+b},
		}},
		}))
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

	//初始化微信配置
	weixinModel.WX_APP_ID = conf.Server.WxAppId
	weixinModel.WX_MCH_ID = conf.Server.WxMchId
	weixinModel.WX_API_KEY = conf.Server.WxApiKey
	weixinModel.WX_APP_SECRET = conf.Server.WxAppSecret
	//初始化微信回调
	weixinModel.WxPayInit()

	//初始化大厅rpc
	rpcService.HallPool.Init(conf.Server.HallRpcAddr, 1)

	//开始监听http
	m.Run(conf.Server.HttpIp, 8089)

}
