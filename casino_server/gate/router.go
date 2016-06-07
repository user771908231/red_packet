package gate

import (
)
import (
	"casino_server/msg"
	"casino_server/msg/bbproto"
	"casino_server/game"
	"casino_server/login"
)

func init() {
	// 这里指定消息 Hello 路由到 game 模块
	// 模块间使用 ChanRPC 通讯，消息路由也不例外,指定json格式的路由
	msg.Processor.SetRouter(&msg.Hello{}, game.ChanRPC)

	//指定protobuf格式的路由
	msg.PortoProcessor.SetRouter(&bbproto.TestP1{},game.ChanRPC)
	msg.PortoProcessor.SetRouter(&bbproto.Reg{},login.ChanRPC)
	msg.PortoProcessor.SetRouter(&bbproto.ReqAuthUser{},login.ChanRPC)
}