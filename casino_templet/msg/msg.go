package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_common/proto/ddproto"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&ddproto.Heartbeat{})        //0连接服务器
	Processor.Register(&ddproto.QuickConn{})        //1连接服务器
	Processor.Register(&ddproto.AckQuickConn{})//2登录游戏
}
