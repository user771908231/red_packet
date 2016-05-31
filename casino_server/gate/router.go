package gate

import (
)
import (
	"casino_server/msg"
	"casino_server/msg/bbproto"
	"casino_server/game"
)

func init() {
	// 这里指定消息 Hello 路由到 game 模块
	// 模块间使用 ChanRPC 通讯，消息路由也不例外
	msg.Processor.SetRouter(&msg.Hello{}, game.ChanRPC)
	msg.PortoProcessor.SetRouter(&bbproto.N{},game.ChanRPC)
}