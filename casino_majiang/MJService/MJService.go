package MJService

import (
	majiangProto "casino_majiang/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_majiang/msg/funcsInit"
)


//service的作用就是handler的具体实现


/*
	创建room

 */
func HandlerGame_CreateRoom(m *majiangProto.Game_CreateRoom, a gate.Agent) {
	log.T("收到请求，HandlerGame_CreateRoom(m[%v],a[%v])", m, a)
	result := newProto.NewGame_AckCreateRoom()
	a.WriteMsg(result)
}