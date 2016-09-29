package majiang

import (
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/AgentService"
	"github.com/name5566/leaf/log"
)

var MJUSER_STATUS_INTOROOM int32 = 1; ///刚进入游戏

//麻将玩家

//发送接口
func (u *MjUser)WriteMsg(p proto.Message) error {
	agent := AgentService.GetAgent(u.GetUserId())
	if agent != nil {
		agent.WriteMsg(p)
	} else {
		log.Fatal("给用户[%v]发送proto[%v]失败，因为没有找到用户的agent。", u.UserId, p)
	}
	return nil
}


