package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_doudizhu/msg/protogo"
	"casino_server/msg/bbprotogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(bbproto.NullMsg{})        //0连接服务器
	Processor.Register(&ddzproto.Game_QuickConn{})        //1连接服务器
	Processor.Register(&ddzproto.Game_AckQuickConn{})//2登录游戏

}
