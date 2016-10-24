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
	handler(&mjProto.Game_DissolveDesk{}, handlerDissolveDesk)
	handler(&mjProto.Game_Ready{}, handlerGame_Ready)

	handler(&mjProto.Game_DingQue{}, handlerGame_DingQue)        //定缺
	handler(&mjProto.Game_ExchangeCards{}, handlerGame_ExchangeCards)        //换3张

	handler(&mjProto.Game_SendOutCard{}, handlerGame_SendOutCard) //出牌

	//碰、杠、过、胡
	handler(&mjProto.Game_ActPeng{}, handlerGame_ActPeng)
	handler(&mjProto.Game_ActGang{}, handlerGame_ActGang)
	handler(&mjProto.Game_ActGuo{}, handlerGame_ActGuo)
	handler(&mjProto.Game_ActHu{}, handlerGame_ActHu)

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
	MJService.HandlerGame_EnterRoom(m.GetHeader().GetUserId(), m.GetPassWord(), a)
}


//解散房间
func handlerDissolveDesk(args []interface{}) {
	m := args[0].(*mjProto.Game_DissolveDesk)        //解散房间
	MJService.HandlerDissolveDesk(m.GetHeader().GetUserId())
}

//准备游戏
func handlerGame_Ready(args []interface{}) {
	m := args[0].(*mjProto.Game_Ready)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_Ready(m, a)
}

//定缺
func handlerGame_DingQue(args []interface{}) {
	m := args[0].(*mjProto.Game_DingQue)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_DingQue(m, a)
}

//换3张
func handlerGame_ExchangeCards(args []interface{}) {
	m := args[0].(*mjProto.Game_ExchangeCards)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_ExchangeCards(m, a)
}


//出牌
func handlerGame_SendOutCard(args []interface{}) {
	m := args[0].(*mjProto.Game_SendOutCard)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_SendOutCard(m, a)
}

//碰
func handlerGame_ActPeng(args []interface{}) {
	m := args[0].(*mjProto.Game_ActPeng)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_ActPeng(m, a)
}

//杠
func handlerGame_ActGang(args []interface{}) {
	m := args[0].(*mjProto.Game_ActGang)
	MJService.HandlerGame_ActGang(m)
}

//过
func handlerGame_ActGuo(args []interface{}) {
	m := args[0].(*mjProto.Game_ActGuo)
	MJService.HandlerGame_ActGuo(m)
}

//胡
func handlerGame_ActHu(args []interface{}) {
	m := args[0].(*mjProto.Game_ActHu)
	MJService.HandlerGame_ActHu(m)
}



