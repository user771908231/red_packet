package gate

import (
	"wstest/msg"
	"casino_server/msg/bbprotogo"
	"wstest/game"
	"wstest/login"
)

func init() {
	msg.Processor.SetRouter(&bbproto.NullMsg{}, game.ChanRPC)
	msg.Processor.SetRouter(&bbproto.REQQuickConn{}, login.ChanRPC)
}
