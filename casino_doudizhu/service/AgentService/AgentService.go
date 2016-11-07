package AgentService

import (
	"github.com/name5566/leaf/gate"
)

func init() {
	AgentService.agents = make(map[uint32]gate.Agent)        //初始化map
}

//存放全局的agent
var AgentService struct {
	agents map[uint32]gate.Agent
}

//得到一个agent
func GetAgent(userId uint32) gate.Agent {
	return AgentService.agents[userId]
}

//设置一个agent
func SetAgent(userId uint32, a gate.Agent) {
	AgentService.agents[userId] = a
}

