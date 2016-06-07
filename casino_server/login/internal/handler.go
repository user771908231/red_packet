package internal

import (
	"reflect"
	"casino_server/msg/bbproto"
	"github.com/name5566/leaf/log"

	"github.com/name5566/leaf/gate"
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


func handleReqAuthUser(args []interface{}){
	log.Debug("进入login.handler.handleReqAuthUser()")
	// 收到的 Hello 消息
	m := args[0].(*bbproto.ReqAuthUser)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("接收到的AppVersion %v", *m.AppVersion)
	//给发送者回应一个 Hello 消息

	var e string
	e = "收到了消息"
	var c int32
	c = 1
	var header bbproto.ProtoHeader
	header.Error = &e
	header.Code = &c;
	var data bbproto.ReqAuthUser
	data.Header = &header

	a.WriteMsg(&data)
}