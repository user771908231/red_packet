package MJService

import (
	mjProto "casino_majiang/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_majiang/msg/funcsInit"
)


//service的作用就是handler的具体实现


/*
	创建room
	用户创建房间的逻辑
	1,如果用户之前已经创建了房间，怎么处理？
	2,余额不足怎么处理
	3,创建成功之后

 */
func HandlerGame_CreateRoom(m *mjProto.Game_CreateRoom, a gate.Agent) {
	log.T("收到请求，HandlerGame_CreateRoom(m[%v],a[%v])", m, a)
	result := newProto.NewGame_AckCreateRoom()
	result.Password = new(string);
	*result.Password = "MYPASS";
	a.WriteMsg(result)
}

func HandlerGame_EnterRoom(m *mjProto.Game_EnterRoom, a gate.Agent) {
	log.T("收到请求，HandlerGame_CreateRoom(m[%v],a[%v])", m, a)
	// TODO: 玩家进入房间

	result := &mjProto.Game_AckCreateRoom{}
	newProto.SuccessHeader(result.Header)

	a.WriteMsg(result)
}

//
func HandlerGame_Ready(m *mjProto.Game_Ready, a gate.Agent) {
	log.T("收到请求，game_Ready(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckCreateRoom{}
	newProto.SuccessHeader(result.Header)

	a.WriteMsg(result)
}



//定缺
func HandlerGame_DingQue(m *mjProto.Game_DingQue, a gate.Agent) {
	log.T("收到请求，HandlerGame_DingQue(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckEnterRoom()
	//a.WriteMsg(result)
}

//换3张
func HandlerGame_ExchangeCards(m *mjProto.Game_ExchangeCards, a gate.Agent) {
	log.T("收到请求，HandlerGame_ExchangeCards(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckExchangeCards()
	//a.WriteMsg(result)
}

//出牌
func HandlerGame_SendOutCard(m *mjProto.Game_SendOutCard, a gate.Agent) {
	log.T("收到请求，HandlerGame_SendOutCard(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckEnterRoom()
	//a.WriteMsg(result)
}

//碰
func HandlerGame_ActPeng(m *mjProto.Game_ActPeng, a gate.Agent) {
	log.T("收到请求，HandlerGame_ActPeng(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActPeng{}
	newProto.SuccessHeader(result.Header)

	a.WriteMsg(result)
}

//杠
func HandlerGame_ActGang(m *mjProto.Game_ActGang, a gate.Agent) {
	log.T("收到请求，game_ActGang(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActGang{}
	newProto.SuccessHeader(result.Header)

	a.WriteMsg(result)
}

//过
func HandlerGame_ActGuo(m *mjProto.Game_ActGuo, a gate.Agent) {
	log.T("收到请求，game_ActGuo(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActGuo{}
	newProto.SuccessHeader(result.Header)

	a.WriteMsg(result)
}

//胡
func HandlerGame_ActHu(m *mjProto.Game_ActHu, a gate.Agent) {
	log.T("收到请求，game_ActHu(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActHu{}
	newProto.SuccessHeader(result.Header)

	a.WriteMsg(result)
}
