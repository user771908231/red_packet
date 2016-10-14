package internal

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_majiang/service/majiang"
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
	log.T("agent 断开连接...")

	agentData := a.UserData()
	if agentData != nil {
		userData := agentData.(*majiang.MjSession)
		desk := majiang.GetMjDeskBySession(userData.GetUserId())
		if desk != nil {
			//这里一般不存在desk==nil的情况
			desk.SetOfflineStatus(userData.GetUserId())
		}
	}

}
