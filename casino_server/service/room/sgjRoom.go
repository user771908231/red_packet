package room

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
)

var SGJRoom sgjRoom        //水果机的房间


func init(){
	OninitSgjRoom()
}

/**

水果机的room
 */
type sgjRoom struct {
	room
}

/**
初始化水果机的房间
 */
func OninitSgjRoom(){
	log.T("初始化水果机的房间")
	SGJRoom.AgentMap = make(map[uint32] gate.Agent)
}


