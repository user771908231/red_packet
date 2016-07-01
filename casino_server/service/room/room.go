package room

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"github.com/golang/protobuf/proto"
	"sync"
)


/**
游戏房间
 */
type room struct {
	sync.Mutex
	Type int
	RoomId	int32				//房间号
	AgentMap map[uint32] gate.Agent
}

func (r *room) AddAgent(userId uint32,a gate.Agent){
	log.T("userId%v的agent放在CachOutRoom中管理\n",userId)
	r.AgentMap[userId] = a

	//打印出 增加连接之后,但当前房间里的连接
	for key := range r.AgentMap {
		log.Normal("当前存在的连接%v",key)
	}
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

		//首先判断连接是否有断开
		a :=r.AgentMap[key]

		m := "服务器的消息"
		data := bbproto.RoomMsg{}
		data.RoomId = &roomId
		data.Msg    = &m
		a.WriteMsg(&data)
		log.Normal("给%v发送消息,发送完毕",key)
	}
}

/**
	给所有的人广播消息,ignoreUserId 的除外
		目前暂时没有实现这个功能
 */
func (r *room) BroadcastProto(p proto.Message,ignoreUserId int32){
	log.Normal("给每个房间发送proto 消息%v",p)
	for key := range r.AgentMap {
		log.Normal("开始给%v发送消息",key)
		//首先判断连接是否有断开
		a :=r.AgentMap[key]
		a.WriteMsg(p)
		log.Normal("给%v发送消息,发送完毕",key)
	}
}
