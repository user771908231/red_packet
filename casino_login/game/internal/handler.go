package internal

import (
	"reflect"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {

}


//创建房间
func handlerCreateDesk(args []interface{}) {

}

//进入房间
func HandlerEnterRoom(args []interface{}) {

}

//准备
func gandlerReady(args []interface{}) {

}

//抢地主
func handlerQiangDiZhu(args []interface{}) {

}

//加倍
func handlerJiaBei(args []interface{}) {

}

//出牌
func handlerChuPai(args []interface{}) {

}







