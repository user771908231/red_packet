package gate

import (
	"casino_majiang/msg"
	"casino_server/msg/bbprotogo"
	"casino_majiang/game"
	"casino_majiang/msg/protogo"
)

func init() {
	msg.Processor.SetRouter(&bbproto.NullMsg{}, game.ChanRPC)
	//msg.Processor.SetRouter(&bbproto.REQQuickConn{}, login.ChanRPC)
	msg.Processor.SetRouter(&mjproto.Game_CreateRoom{},game.ChanRPC)


}
