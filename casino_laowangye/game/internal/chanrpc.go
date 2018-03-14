package internal

import (
	"github.com/name5566/leaf/gate"
	"casino_laowangye/service/laowangye"
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

	//朋友卓老王爷断线处理
	laowangye.OnOffLine(a)
}
