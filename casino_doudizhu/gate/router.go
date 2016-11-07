package gate

import (
	"casino_majiang/msg"
	"casino_majiang/game"
	"casino_server/msg/bbprotogo"
)

func init() {
	msg.Processor.SetRouter(&bbproto.NullMsg{}, game.ChanRPC)        //空协议

	//登陆相关

	//room 相关

	//desk 相关

}
