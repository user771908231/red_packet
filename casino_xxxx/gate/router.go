package gate

import (
	"casino_common/proto/ddproto"
	"casino_common/common/handlers"
	"casino_paodekuai/msg"
	"casino_paodekuai/game"
)

func init() {
	msg.PDKProcessor.SetHandler(&ddproto.Heartbeat{}, handlers.HandlerHeartBeat) //心跳
	msg.PDKProcessor.SetRouter(&ddproto.PdkReqCreateDesk{}, game.ChanRPC)        //牌得快
}
