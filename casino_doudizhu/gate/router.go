package gate

import (
	"casino_majiang/msg"
	"casino_majiang/game"
	"casino_server/msg/bbprotogo"
)

func init() {
	msg.Processor.SetRouter(&bbproto.NullMsg{}, game.ChanRPC)
}
