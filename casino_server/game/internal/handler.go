package internal

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/gate"
	"reflect"
	"server/msg"
	"server/msg/bbproto"
	"time"
)

func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&msg.Hello{}, handleHello)
	handler(&bbproto.N{},handleProtHello)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleHello(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*msg.Hello)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("hello %v", m.Name)

	// 给发送者回应一个 Hello 消息

	a.WriteMsg(&msg.Hello{
		Name: "bbgogogogogogo",
	})

}

func handleProtHello(args []interface{}){

	log.Debug("进入处理函数")
	// 收到的 Hello 消息
	m := args[0].(*bbproto.N)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("接收到的name %v", *m.Name)

	 //给发送者回应一个 Hello 消息
	var data bbproto.N
	var n string = "go------!"+ time.Now().String()
	data.Name = &n
	a.WriteMsg(&data)
}