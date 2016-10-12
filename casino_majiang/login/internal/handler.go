package internal

import (
	"reflect"
	"casino_server/common/log"
	"casino_majiang/msg/protogo"
	"github.com/name5566/leaf/gate"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(mjproto.Game_QuickConn{}, handlerREQQuickConn)
}

//处理登陆
func handlerREQQuickConn(args []interface{}) {
	m := args[0].(*mjproto.Game_QuickConn)
	a := args[1].(gate.Agent)
	log.T("请求了handlerREQQuickConn m[%v]")
	a.WriteMsg(m)
}
