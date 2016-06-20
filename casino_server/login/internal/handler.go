package internal

import (
	"reflect"
	"casino_server/msg/bbproto"
	"github.com/name5566/leaf/gate"
	"casino_server/conf/casinoConf"
	"casino_server/service"
	"casino_server/common/log"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&bbproto.Reg{},handleProtHello)
	handleMsg(&bbproto.ReqAuthUser{},handleReqAuthUser)

}



/**
	处理注册消息的方法
	此方法可能暂时没有使用,而使用handleReqAuthUser
 */
func handleProtHello(args []interface{}){
	log.Debug("进入login.handler.handleProtHello()")
	// 收到的 Hello 消息
	m := args[0].(*bbproto.Reg)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("接收到的name %v", *m.Name)
	//给发送者回应一个 Hello 消息
	var data bbproto.Reg
	//var n string = "a"+ time.Now().String()
	var n string = "hi leaf"
	data.Name = &n
	a.WriteMsg(&data)
}


/*

登陆注册的流程:
1,获取到请求信息
2,如果快速登录,则随机分配一个userId,并且自动注册
3,如果用户已经注册了,则走登录流程
 */
func handleReqAuthUser(args []interface{}){
	log.Debug("进入login.handler.handleReqAuthUser()")

	// 收到的 Hello 消息
	m := args[0].(*bbproto.ReqAuthUser)
	// 消息的发送者
	a := args[1].(gate.Agent)
	// 输出收到的消息的内容
	log.T("agent:",a)
	log.Debug("介绍到的reqAuthUser %v", *m)

	//判断是快速登陆还是
	loginWay := service.CheckUserId(m.GetHeader().GetUserId())
	switch loginWay {
	case casinoConf.LOGIN_WAY_QUICK:
		log.Debug("快速登录模式")
		service.QuickLogin(m,a)
	case casinoConf.LOGIN_WAY_LOGIN:
		log.Debug("普通登录模式")
		service.Login(m)
	default:
		log.Debug("没有找到合适的登录方式")
	}

}