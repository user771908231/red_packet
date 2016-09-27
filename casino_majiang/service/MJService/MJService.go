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
	//result := newProto.NewGame_AckEnterRoom()
	//a.WriteMsg(result)
}

//
func HandlerGame_Ready(m *mjProto.Game_EnterRoom, a gate.Agent) {
	log.T("收到请求，HandlerGame_CreateRoom(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckEnterRoom()
	//a.WriteMsg(result)
}

//出牌
func HandlerGame_SendOutCard(m *mjProto.Game_EnterRoom, a gate.Agent) {
	log.T("收到请求，HandlerGame_CreateRoom(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckEnterRoom()
	//a.WriteMsg(result)
}

