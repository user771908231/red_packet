package gamedata

import "sync"



/**
次数据和 agent 绑定在一起
 */
type AgentUserData struct {
	sync.Mutex
	userId int32
	fruitBlance int32
	fruitWinScoPre	int32		//上一次得分
}

/**
水果机上分
 */

func (user *AgentUserData) FriutshangFen(){
	user.Lock()
	defer user.Unlock()

}


/**
水果机 下分
 */

func (user *AgentUserData) FriutXiaFen(){

}

