package gate

import (
	"casino_common/proto/ddproto"
	"casino_templet/msg"
	"casino_templet/login"
)

func init() {
	msg.Processor.SetRouter(&ddproto.QuickConn{}, login.ChanRPC)

}
