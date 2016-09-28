package AgentService

import (
	"github.com/name5566/leaf/gate"
	"casino_majiang/service/majiang"
	"casino_server/utils/numUtils"
	"strings"
	"casino_server/utils/redisUtils"
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

// session相关的...

var MJSESSION_KEY_PRE = "redis_majiang_session"

func getSessionKey(userId uint32) string {
	idstr, _ := numUtils.Uint2String(userId)
	ret := strings.Join([]string{MJSESSION_KEY_PRE, idstr}, "_")
	return ret
}

func GetSession(userId uint32) *majiang.MjSession {
	s := redisUtils.GetObj(getSessionKey(userId), &majiang.MjSession{})
	if s != nil {
		return s.(*majiang.MjSession)
	} else {
		return nil
	}
	return nil
}