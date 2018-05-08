package internal

import (
	"github.com/name5566/leaf/gate"
	"casino_paosangong/service/paosangong"
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

	//朋友卓牛牛断线处理
	paosangong.OnOffLine(a)
}
