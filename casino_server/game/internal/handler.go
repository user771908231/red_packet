package internal

import (
	"github.com/name5566/leaf/gate"
	"reflect"
	"casino_server/msg/bbproto"
	"casino_server/common/log"
	"casino_server/gamedata"
)

func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&bbproto.TestP1{},handleTestP1)
	handler(&bbproto.Reg{},handleProtHello)
	handler(&bbproto.GetIntoRoom{},handlerGetIntoRoom)
	handler(&bbproto.RoomMsg{},handlerRoomMsg)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}


func handleProtHello(args []interface{}){
	log.T("进入handleProtHello()")
	// 收到的 Hello 消息
	m := args[0].(*bbproto.Reg)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("接收到的name %v", *m.Name)
	 //给发送者回应一个 Hello 消息
	var data bbproto.Reg
	//var n string = "a"+ time.Now().String()
	var n string = "hi leaf"
	data.Name = &n
	a.WriteMsg(&data)
}

func handleTestP1(args[]interface{}){
	log.Debug("进入handleTestP1()")
	// 收到的 Hello 消息
	m := args[0].(*bbproto.TestP1)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("接收到的name %v", *m.Name2)
	//给发送者回应一个 Hello 消息
	var data bbproto.TestP1
	var n string = "hi leaf testp2"
	data.Name2 = &n
	a.WriteMsg(&data)
}


/**
	请求进入游戏房间
	1,分配房间(根据游戏类型)
	2,proto中的标志 如果in=1表示进入房间,其他则表示退出房间

 */
func handlerGetIntoRoom(args []interface{}){
	log.T("进入到 game.handlerGetIntoRoom()")
	m := args[0].(*bbproto.GetIntoRoom)		//请求体
	a := args[1].(gate.Agent)		//连接
	log.T("请求进入房间的user %v ,in:%v\n",m.GetUserId(),m.GetIn())
	if m.GetIn() == 1 {
		gamedata.CashOutRoom.AddAgent(m.GetUserId(),a)
	}else{
		gamedata.CashOutRoom.RemoveAgent(m.GetUserId())
	}
}

/**

处理roomMsg

 */
func handlerRoomMsg(args []interface{}){
	log.T("进入到 game.handlerRoomMsg()")
	m := args[0].(*bbproto.RoomMsg)		//请求体
	gamedata.CashOutRoom.BroadcastMsg(m.GetRoomId(),m.GetMsg())
}