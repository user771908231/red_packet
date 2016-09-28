package internal

import (
	"reflect"
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	mjProto "casino_majiang/msg/protogo"
	"casino_majiang/service/MJService"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&bbproto.NullMsg{}, handlerNull)
	handler(&mjProto.Game_CreateRoom{}, handlerGame_CreateRoom)
	handler(&mjProto.Game_EnterRoom{}, handlerGame_EnterRoom)
	handler(&mjProto.Game_Ready{}, handlerGame_Ready)
	handler(&mjProto.Game_SendOutCard{}, handlerGame_SendOutCard)
}

func handlerNull(args []interface{}) {
	log.T("进入到 game.handlerNull()")
	m := args[0].(*bbproto.NullMsg)
	a := args[1].(gate.Agent)
	a.WriteMsg(m)
}

//处理创建房间
func handlerGame_CreateRoom(args []interface{}) {
	m := args[0].(*mjProto.Game_CreateRoom)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_CreateRoom(m, a)
}

//处理进入房间
func handlerGame_EnterRoom(args []interface{}) {
	m := args[0].(*mjProto.Game_EnterRoom)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_EnterRoom(m, a)
}

//准备游戏
func handlerGame_Ready(args []interface{}) {
	m := args[0].(*mjProto.Game_Ready)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_Ready(m, a)
}

//准备游戏
func handlerGame_SendOutCard(args []interface{}) {
	m := args[0].(*mjProto.Game_SendOutCard)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_SendOutCard(m, a)
}
