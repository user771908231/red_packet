package internal

import (
	"reflect"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	//handler({}, handlerREQQuickConn)
}

//处理登陆
//func handlerREQQuickConn(args []interface{}) {
//	m := args[0].(*bbproto.REQQuickConn)
//	log.Debug("进入到 game.handlerREQQuickConn()m[%v]", m)
//	a := args[1].(gate.Agent)
//
//	m.ChannelID = new(string)
//	*m.ChannelID = "测试回复..."
//	a.WriteMsg(m)
//}
