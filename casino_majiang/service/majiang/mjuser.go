package majiang

import (
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/AgentService"
	"github.com/name5566/leaf/log"
)


//麻将玩家
//type MjUser struct {
//
//}

//发送接口
func (u *MjUser)WriteMsg(p proto.Message) error {
	agent := AgentService.GetAgent(u.UserId)
	if agent != nil {
		agent.WriteMsg(p)
	} else {
		log.Fatal("给用户[%v]发送proto[%v]失败，因为没有找到用户的agent。", u.UserId, p)
	}
}
