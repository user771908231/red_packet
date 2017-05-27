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
	log.T("建立一个新的连接, 来自[%v]的请求", a.RemoteAddr())
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.T("agent 断开连接...")

	agentData := a.UserData()
	if agentData != nil {

	}
}
