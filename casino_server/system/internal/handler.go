package internal

import (
	"casino_server/msg/bbprotogo"
	"reflect"
	"casino_server/common/log"
	"github.com/name5566/leaf/gate"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}


func init() {
	handleMsg(&bbproto.HeatBeat{},handleHeatBeat)

}

/*
 */
func handleHeatBeat(args []interface{}){
	log.Debug("进入login.handler.handleHeatBeat()")

	// 收到的 Hello 消息
	//m := args[0].(*bbproto.HeatBeat)
	// 消息的发送者
	a := args[1].(gate.Agent)
	// 输出收到的消息的内容
	log.T("agent:",a)

}