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
	handler(&ddzproto.Game_Pull{}, handlePull)    //
	handler(&ddzproto.Game_OutCards{}, handlerChuPai)    //
	handler(&ddzproto.Game_ActGuo{}, handlerActGuo)    //
	handler(&ddzproto.Game_DissolveDesk{}, handlerDissolveDesk)    //
	handler(&ddzproto.Game_LeaveDesk{}, handlerLeaveDesk)    //
	handler(&ddzproto.Game_Message{}, handlerMessage)    //
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
	a := args[1].(gate.Agent)

	var enterType int32 = 0
	DdzService.HandlerEnterDesk(m.GetUserId(), m.GetPassWord(), enterType, a)

}

//准备
func handlerReady(args []interface{}) {
	m := args[0].(*ddzproto.Game_Ready)
	DdzService.HandlerFDdzReady(m.GetUserId())
}

func handlerJiaoDiZhu(args []interface{}) {
	m := args[0].(*ddzproto.Game_JiaoDiZhu)
	DdzService.HandlerJiaoDiZhu(m.GetUserId())
}

//明牌
func handlerShowHandPokers(args []interface{}) {
	m := args[0].(*ddzproto.Game_ShowHandPokers)
	DdzService.HandlerShowHandPokers(m.GetUserId())
}

//闷牌
func handlerMenuZhua(args []interface{}) {
	m := args[0].(*ddzproto.Game_MenuZhua)
	DdzService.HandlerMenuZhua(m.GetUserId())

}
func handlerSeeCards(args []interface{}) {
	m := args[0].(*ddzproto.Game_SeeCards)
	DdzService.HandlerSeeCards(m.GetUserId())
}

func handlePull(args []interface{}) {
	m := args[0].(*ddzproto.Game_Pull)
	DdzService.HandlePull(m.GetUserId())
}

//抢地主
func handlerQiangDiZhu(args []interface{}) {
	m := args[0].(*ddzproto.Game_RobDiZhu)
	DdzService.HandlerQiangDiZhu(m.GetUserId())

}

func handlerActGuo(args []interface{}) {
	m := args[0].(*ddzproto.Game_ActGuo)
	DdzService.HandlerActPass(m.GetUserId())
}

func handlerDissolveDesk(args []interface{}) {
	m := args[0].(*ddzproto.Game_DissolveDesk)
	DdzService.HandlerDissolveDesk(m.GetUserId())
}

func handlerLeaveDesk(args []interface{}) {
	m := args[0].(*ddzproto.Game_LeaveDesk)
	DdzService.HandlerLeaveDesk(m.GetHeader().GetUserId())
}

func handlerMessage(args []interface{}) {
	//m := args[0].(*ddzproto.Game_Message)
	DdzService.HandlerMessage()

}
func handlerGameRecord(args []interface{}) {
	//m := args[0].(*ddzproto.Game_GameRecord)
	DdzService.HandlerGameRecord()

}

//加倍
func handlerJiaBei(args []interface{}) {
	//m := args[0].(*ddzproto.Game_Double)
	DdzService.HandlerJiaBei()
}

//出牌
func handlerChuPai(args []interface{}) {
	m := args[0].(*ddzproto.Game_OutCards)
	DdzService.HandlerActOut(m.GetHeader().GetUserId())
}







