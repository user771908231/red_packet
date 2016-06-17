package gamedata

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"fmt"
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

func (r *room) BroadcastMsg(){
	/* 使用 key 输出 map 值 */
	for key := range r.AgentMap {
		fmt.Println("广播测试 of",key)
		a :=r.AgentMap[key]
		result := bbproto.GetIntoRoom{}
		key2 := int32(key+1)
		result.RoomId = &key2
		a.WriteMsg(&result)
	}
}
