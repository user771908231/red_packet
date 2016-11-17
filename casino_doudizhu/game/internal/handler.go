package internal

import (
	"reflect"
	"casino_doudizhu/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_doudizhu/service/DdzService"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)

}

func init() {

	handler(&ddzproto.Game_CreateRoom{}, handlerCreateDesk)    //创建房间
	handler(&ddzproto.Game_EnterRoom{}, HandlerEnterRoom)    //进入房间
	handler(&ddzproto.Game_Ready{}, handlerReady)    //准备
	handler(&ddzproto.Game_JiaoDiZhu{}, handlerJiaoDiZhu)    //
	handler(&ddzproto.Game_RobDiZhu{}, handlerQiangDiZhu)    //
	handler(&ddzproto.Game_Double{}, handlerJiaBei)    //
	handler(&ddzproto.Game_ShowHandPokers{}, handlerShowHandPokers)    //
	handler(&ddzproto.Game_MenuZhua{}, handlerMenuZhua)    //
	handler(&ddzproto.Game_SeeCards{}, handlerSeeCards)    //
	handler(&ddzproto.Game_Pull{}, handlerull)    //
	handler(&ddzproto.Game_OutCards{}, handlerChuPai)    //
	handler(&ddzproto.Game_ActGuo{}, handlerActGuo)    //
	handler(&ddzproto.Game_SendCurrentResult{}, handlerSendCurrentResult)    //
	handler(&ddzproto.Game_DissolveDesk{}, handlerDissolveDesk)    //
	handler(&ddzproto.Game_LeaveDesk{}, handlerLeaveDesk)    //
	handler(&ddzproto.Game_Message{}, hanlerMessage)    //
	handler(&ddzproto.Game_GameRecord{}, handlerGameRecord)    //
}


//创建房间
func handlerCreateDesk(args []interface{}) {
	m := args[0].(*ddzproto.Game_CreateRoom)
	a := args[1].(gate.Agent)
	DdzService.HandlerCreateDesk(m.GetHeader().GetUserId(), a)
}

//进入房间
func HandlerEnterRoom(args []interface{}) {
	m := args[0].(*ddzproto.Game_EnterRoom)

}

//准备
func handlerReady(args []interface{}) {
	m := args[0].(*ddzproto.Game_Ready)

}

func handlerJiaoDiZhu(args []interface{}) {
	m := args[0].(*ddzproto.Game_JiaoDiZhu)
}

//明牌
func handlerShowHandPokers(args []interface{}) {
	m := args[0].(*ddzproto.Game_ShowHandPokers)
}

//闷牌
func handlerMenuZhua(args []interface{}) {
	m := args[0].(*ddzproto.Game_MenuZhua)
}
func handlerSeeCards(args []interface{}) {
	m := args[0].(*ddzproto.Game_SeeCards)
}

func handlerull(args []interface{}) { // rull or pull ?
	m := args[0].(*ddzproto.Game_Pull)
}

//抢地主
func handlerQiangDiZhu(args []interface{}) {
	m := args[0].(*ddzproto.Game_)
}

func handlerActGuo(args []interface{}) {
	m := args[0].(*ddzproto.Game_ActGuo)
}

func handlerSendCurrentResult(args []interface{}) {
	m := args[0].(*ddzproto.Game_SendCurrentResult)
}

func handlerDissolveDesk(args []interface{}) {
	m := args[0].(*ddzproto.Game_DissolveDesk)
}

func handlerLeaveDesk(args []interface{}) {
	m := args[0].(*ddzproto.Game_LeaveDesk)
}

func hanlerMessage(args []interface{}) {
	m := args[0].(*ddzproto.Game_Message)
}
func handlerGameRecord(args []interface{}) {
	m := args[0].(*ddzproto.Game_GameRecord)
}

//加倍
func handlerJiaBei(args []interface{}) {
	m := args[0].(*ddzproto.Game_)
}

//出牌
func handlerChuPai(args []interface{}) {
	m := args[0].(*ddzproto.Game_)
}







