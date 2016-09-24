package internal

import (
	"reflect"
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	majiangProto "casino_majiang/msg/bbprotogo"
	"casino_majiang/service/MJService"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&bbproto.NullMsg{}, handlerNull)
	handler(&majiangProto.Game_CreateRoom{}, handlerGame_CreateRoom)

}

func handlerNull(args []interface{}) {
	log.T("进入到 game.handlerNull()")
	m := args[0].(*bbproto.NullMsg)
	a := args[1].(gate.Agent)
	a.WriteMsg(m)
}

//处理创建房间
func handlerGame_CreateRoom(args []interface{}) {
	m := args[0].(*majiangProto.Game_CreateRoom)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_CreateRoom(m, a)
}
