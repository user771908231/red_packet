package majiang

import (
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/AgentService"
	"github.com/name5566/leaf/log"
)

var MJUSER_STATUS_INTOROOM int32 = 1; ///刚进入游戏
var MJUSER_STATUS_SEATED int32 = 2; ///刚进入游戏
var MJUSER_STATUS_READY int32 = 3; ///刚进入游戏


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

//是否是准备中...
func (u *MjUser) IsReady() bool {
	return u.GetStatus() == MJUSER_STATUS_SEATED
}

//玩家是否在游戏状态中
func (u *MjUser) IsGaming() bool {
	return true
	
}


