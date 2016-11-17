package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_server/msg/bbprotogo"
	"casino_login/msg/protogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&bbproto.NullMsg{})        //0连接服务器
	Processor.Register(&loginproto.Game_QuickConn{})        //1连接服务器
	Processor.Register(&loginproto.Game_AckQuickConn{})//2登录游戏
}
