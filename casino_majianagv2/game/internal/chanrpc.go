package internal

import (
	"github.com/name5566/leaf/gate"
	"casino_common/common/log"
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

	}
}
