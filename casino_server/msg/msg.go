package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_server/msg/bbproto"
)

// protobuf 消息处理器
var PortoProcessor = protobuf.NewProcessor()

func init() {
	// 这里我们注册了一个 JSON 消息 Hello

	//次处注册proto 的消息
	PortoProcessor.Register(&bbproto.TestP1{})	//0	测试用
	PortoProcessor.Register(&bbproto.Reg{})		//1	注册协议(已经废弃)
	PortoProcessor.Register(&bbproto.ReqAuthUser{})	//2	登陆、注册的协议
	PortoProcessor.Register(&bbproto.HeatBeat{})	//3	心跳协议,检测网络是否联通
	PortoProcessor.Register(&bbproto.GetIntoRoom{})	//4	进入房间时候的请求



}
