package internal

import (
	"reflect"
	"casino_common/proto/ddproto"
	"casino_common/common/handlers"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&ddproto.CommonReqGameLogin{}, handlers.HandlerGame_Login)
}
