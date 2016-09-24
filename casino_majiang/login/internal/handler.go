package internal

import (
	"reflect"
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/gate"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&bbproto.REQQuickConn{}, handlerREQQuickConn)
}

//处理登陆
func handlerREQQuickConn(args []interface{}) {
	m := args[0].(*bbproto.REQQuickConn)
	log.Debug("进入到 game.handlerREQQuickConn()m[%v]", m)
	a := args[1].(gate.Agent)

	m.ChannelID = new(string)
	*m.ChannelID = "测试回复..."
	a.WriteMsg(m)
}
