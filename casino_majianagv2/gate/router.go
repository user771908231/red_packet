package gate

import (
	"casino_common/proto/ddproto"
	"casino_majianagv2/msg"
	"casino_majianagv2/login"
)

func init() {
	msg.Processor.SetRouter(&ddproto.QuickConn{}, login.ChanRPC)

}
