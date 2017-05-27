package internal

import (
	"reflect"
	"casino_common/proto/ddproto"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&ddproto.PdkReqCreateDesk{}, createDesk) //创建房间
}

//创建一个Desk
func createDesk(args []interface{}) {

}

func actReady(args []interface{}) {

}

//打牌
func actOut(args []interface{}) {

}
