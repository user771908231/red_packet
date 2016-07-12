package internal

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
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
	_ = a

	log.Normal("调用了.....game.chanrpc.rpcCloseAgent")
	//对数据做异常退出时的保存工作
	//1,保存游戏数据
	//2,删除连接中管理的agent

}
