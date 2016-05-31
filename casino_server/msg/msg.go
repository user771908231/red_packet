package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"github.com/name5566/leaf/network/json"
	"casino_server/msg/bbproto"
)

// 使用默认的 JSON 消息处理器（默认还提供了 protobuf 消息处理器）
var Processor = json.NewProcessor()
var PortoProcessor = protobuf.NewProcessor()

func init() {
	// 这里我们注册了一个 JSON 消息 Hello
	Processor.Register(&Hello{})
	PortoProcessor.Register(&bbproto.N{})
}

// 一个结构体定义了一个 JSON 消息的格式
// 消息名为 Hello
type Hello struct {
	Name string
}