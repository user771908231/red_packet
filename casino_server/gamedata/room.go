package gamedata

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/msg/bbproto"
)

func init() {
	log.T("初始化CashOutRoom的AgentMap")
	CashOutRoom.AgentMap = make(map[uint32] gate.Agent)
}

var CashOutRoom room

/**
游戏房间
 */
type room struct {
	Type int
	AgentMap map[uint32] gate.Agent
}

func (r *room) AddAgent(userId uint32,a gate.Agent){
	log.T("userId%v的agent放在CachOutRoom中管理\n",userId)
	r.AgentMap[userId] = a
}

func (r *room) RemoveAgent(userId uint32){
	delete(r.AgentMap,userId);
}

/**
	发送信息
 */

func (r *room) BroadcastMsg(roomId int32,msg string){
	log.Normal("给房间号%v发送信息%v",roomId,msg)
	/* 使用 key 输出 map 值 */
	for key := range r.AgentMap {
		log.Normal("开始给%v发送消息",key)
		a :=r.AgentMap[key]
		result := bbproto.RoomMsg{}
		result.RoomId = &roomId
		a.WriteMsg(&result)
	}
}