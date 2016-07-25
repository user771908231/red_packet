package internal

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/gamedata"
	"casino_server/service/room"
)



func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)

	log.Normal("调用了.....game.chanrpc.rpcCloseAgent")
	//对数据做异常退出时的保存工作
	//1,保存游戏数据
	//2,删除连接中管理的agent

	//测试代码--------------------------------------begin-----------------------------------
	//,如果UserId是10006的话,连接断开的时候,desk 删除user
	agentData := a.UserData()
	if agentData == nil {
		log.E("通过agent[%v]取出来的userData 是nil",&a)
	}else{
		log.T("开始出连接断开的处理工作")
		userData := agentData.(*gamedata.AgentUserData)
		if userData.UserId == 10006 {
			desk := room.ThGameRoomIns.GetDeskByUserId(userData.UserId)
			desk.RmThuser(userData.UserId)
		}
	}
	//测试代码--------------------------------------end--------------------------------------

}
