package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_server/msg/bbprotogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&bbproto.NullMsg{})                          //0
	Processor.Register(&bbproto.REQQuickConn{})                        //1
}



