package gamedata


//const

var(
	AGENT_USER_STATUS_ONLINE	int32	=	1	//只是在先 没有在游戏中
	AGENT_USER_STATUS_GAMING	int32	=	2	//游戏中


)

/**
次数据和 agent 绑定在一起
 */
type AgentUserData struct {
	UserId   uint32
	Status   int32  //当前状态
	ThDeskId uint32 //德州扑克deskId
}

//初始化一个agentUser,这个agentUser 在登录的时候使用
func NewAgentUser(userId uint32) *AgentUserData{
	result := &AgentUserData{}
	result.UserId = userId
	result.Status = AGENT_USER_STATUS_ONLINE
	return result
}