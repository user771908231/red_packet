package main

import (
	"casino_tools/white_list_utils/conf"
	"casino_tools/white_list_utils/conf/config"
	"casino_tools/white_list_utils/modules"
	"casino_tools/white_list_utils/routers"
	"casino_common/common/sys"
	"casino_common/proto/ddproto"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"os"
	"casino_common/common/service/rpcService"
	"html/template"
	"fmt"
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
		"admin",
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

	//初始化转转rpc
	rpcService.ZzHzPool.Init(conf.Server.MjZzHzRpcAddr, 1)
	rpcService.PdkPool.Init(conf.Server.PdkRpcAddr, 1)
	rpcService.HallPool.Init(conf.Server.HallRpcAddr, 1)
}

func main() {
	m := macaron.Classic()
	//注册模板
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "templates", IndentJSON: true,
		Funcs:[]template.FuncMap{
			template.FuncMap{
				"fmt": func(str string,v ...interface{}) string {
					return fmt.Sprintf(str, v);
				},
			},
		},
	}))
	//注册微信机器人二维码目录
	m.Use(macaron.Static("../web/public/qrcode"))
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
