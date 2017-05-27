package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_common/proto/ddproto"
)

var PDKProcessor = protobuf.NewProcessor()

func init() {
	PDKProcessor.Register(&ddproto.Heartbeat{})        //0
	PDKProcessor.Register(&ddproto.PdkReqCreateDesk{}) //1 跑得快创建房间
}
