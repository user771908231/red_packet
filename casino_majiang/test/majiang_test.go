package test

import (
	"testing"
	"reflect"
	"casino_server/common/log"
	"casino_majiang/conf/config"
	"casino_majiang/msg/funcsInit"
	"time"
)

func Test(t *testing.T) {
	config.InitConfig(false)
	m1 := newProto.NewGame_SendGameInfo()
	printP(m1)
	m2 := newProto.NewGame_AckActHu()
	printP(m2)

	time.Sleep(time.Second * 3)
}

func printP(msg interface{}) {
	log.T("agent发送的信息 type[%v],\t content[%v]", reflect.TypeOf(msg).String(), msg)
}
