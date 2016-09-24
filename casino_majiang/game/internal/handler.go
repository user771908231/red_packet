package internal

import (
	"reflect"
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&bbproto.NullMsg{}, handlerNull)
}

func handlerNull(args []interface{}) {
	log.T("进入到 game.handlerNull()")
	m := args[0].(*bbproto.NullMsg)
	a := args[1].(gate.Agent)
	a.WriteMsg(m)
}
