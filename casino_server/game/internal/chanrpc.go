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
		log.T("通过agent[%v]取出来的userData 是nil",a)
	}else{

		log.T("开始出a[%v]连接断开的处理工作agentData[%v]",a,agentData)
		userData := agentData.(*gamedata.AgentUserData)
		if userData.UserId == 10006 {
			log.T("因为用户的10006,所以退出的时候,游戏也跟着退出去....")
			desk := room.ThGameRoomIns.GetDeskByUserId(userData.UserId)
			desk.RmThuser(userData.UserId)
		}
	}
	//测试代码--------------------------------------end--------------------------------------

}
