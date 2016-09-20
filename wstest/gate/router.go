package gate

import (
	"wstest/msg"
	"casino_server/msg/bbprotogo"
	"casino_server/game"
)

func init() {
	msg.Processor.Route(&bbproto.NullMsg{},game.ChanRPC)
}
