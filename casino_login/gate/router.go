package gate

import (
	"casino_login/msg"
	"casino_login/msg/protogo"
	"casino_login/login"
)

func init() {
	msg.Processor.SetRouter(&loginproto.Game_QuickConn{}, login.ChanRPC)

}
