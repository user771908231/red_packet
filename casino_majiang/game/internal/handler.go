package internal

import (
	"reflect"
	"casino_server/common/log"
	"github.com/name5566/leaf/gate"
	mjProto "casino_majiang/msg/protogo"
	"casino_majiang/service/MJService"
	"casino_server/service/noticeServer"
	"casino_majiang/msg/funcsInit"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&mjProto.Heartbeat{}, handlerNull)
	handler(&mjProto.Game_CreateRoom{}, handlerGame_CreateDesk)
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


	//战绩相关
	handler(&mjProto.Game_GameRecord{}, handlerGame_GameRecord)

	//聊天
	handler(&mjProto.Game_Message{}, handlerGame_Message)

	//通知信息
	handler(&mjProto.Game_Notice{}, HandlerNotice)

}

func handlerNull(args []interface{}) {
	log.T("进入到 game.Heartbeat()")
	m := args[0].(*mjProto.Heartbeat)
	a := args[1].(gate.Agent)
	a.WriteMsg(m)
}

//处理创建房间
func handlerGame_CreateDesk(args []interface{}) {
	m := args[0].(*mjProto.Game_CreateRoom)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_CreateDesk(m, a)
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
	MJService.HandlerGame_ExchangeCards(m)
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
	MJService.HandlerGame_ActHu(m.GetUserId())
}



//查询战绩
func handlerGame_GameRecord(args []interface{}) {
	m := args[0].(*mjProto.Game_GameRecord)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_GameRecord(m.GetUserId(), a)
}

//聊天的协议
func handlerGame_Message(args []interface{}) {
	m := args[0].(*mjProto.Game_Message)
	MJService.HandlerGame_Message(m)
}

//通知信息
func HandlerNotice(args []interface{}) {
	m := args[0].(*mjProto.Game_Notice)
	a := args[1].(gate.Agent)
	bback := noticeServer.GetNoticeByType(m.GetNoticeType())
	ack := newProto.NewGame_AckNotice()
	*ack.NoticeType = bback.GetNoticeType()
	*ack.NoticeTitle = bback.GetNoticeTitle()
	*ack.NoticeMemo = bback.GetNoticeMemo()
	*ack.NoticeContent = bback.GetNoticeContent()
	ack.Fileds = bback.Fileds

	a.WriteMsg(ack)
}
