package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	casinoProto "casino_server/msg/bbprotogo"
	//majiangProto "casino_majiang/msg/bbprotogo"


)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&casinoProto.NullMsg{})                          //0
}



