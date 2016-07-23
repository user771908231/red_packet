package internal

import (
	"github.com/name5566/leaf/gate"
	"reflect"
	"casino_server/common/log"
	"casino_server/service/rewardService"
	"casino_server/service/fruitService"
	"casino_server/msg/bbprotogo"
	"casino_server/service/zjhService"
	"casino_server/conf/intCons"
	"casino_server/service/room"
	"casino_server/service/thService"
	"casino_server/service/OGservice"
)

func init() {
	//
	handler(&bbproto.TestP1{}, handleTestP1)
	handler(&bbproto.Reg{}, handleProtHello)

	//水果机相关的业务
	handler(&bbproto.GetIntoRoom{}, handlerGetIntoRoom)
	handler(&bbproto.RoomMsg{}, handlerRoomMsg)
	handler(&bbproto.GetRewards{}, handlerRewards)
	handler(&bbproto.Shuiguoji{}, handlerShuiguoji)
	handler(&bbproto.ShuiguojiHilomp{},handlerShuiguojiHilomp)

	//扎金花相关的业务
	handler(&bbproto.ZjhRoom{},handlerZjhRoom)
	handler(&bbproto.ZjhLottery{},handlerZjhLottery)
	handler(&bbproto.ZjhMsg{},handlerZjhMsg)
	handler(&bbproto.ZjhBet{},handlerZjhBet)
	handler(&bbproto.ZjhReqSeat{},handlerZjhReqSeat)
	handler(&bbproto.ZjhQueryNoSeatUser{},handlerZjhQueryNoSeatUser)

	//德州扑克相关的业务
	handler(&bbproto.ThRoom{},handlerThRoom)
	handler(&bbproto.THBet{},handlerThBet)		//押注的协议


	//联众的德州扑克
	handler(&bbproto.Game_LoginGame{},handlerGameLoginGame)	//登陆游戏
	handler(&bbproto.Game_EnterMatch{},handlerGameEnterMatch)//进入房间
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleProtHello(args []interface{}) {
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

func handleTestP1(args[]interface{}) {
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
func handlerGetIntoRoom(args []interface{}) {
	log.T("进入到 game.handlerGetIntoRoom()")
	m := args[0].(*bbproto.GetIntoRoom)                //请求体
	a := args[1].(gate.Agent)                //连接
	log.T("agent:", &a)
	log.T("请求进入房间的user %v ,in:%v\n", m.GetUserId(), m.GetIn())
	if m.GetIn() == intCons.REQ_TYPE_IN{
		room.SGJRoom.AddAgent(m.GetUserId(), a)
	} else {
		room.SGJRoom.RemoveAgent(m.GetUserId())
	}
}

/**

处理roomMsg

 */
func handlerRoomMsg(args []interface{}) {
	log.T("进入到 game.handlerRoomMsg()")
	m := args[0].(*bbproto.RoomMsg)                //请求体
	a := args[1].(gate.Agent)
	log.T("agent:", &a)
	room.SGJRoom.BroadcastMsg(m.GetRoomId(), m.GetMsg())
}


/**
领取奖励的入口都在这里
 */
func handlerRewards(args []interface{}) {
	log.T("进入到 game.handlerRewards()")
	//检测参数是否正确
	m := args[0].(*bbproto.GetRewards)                //请求体
	a := args[1].(gate.Agent)
	err := rewardService.HandlerRewards(m, a)                //调用处理函数来处理
	if err != nil {
		log.E(err.Error())
	}
}

/**
处理水果机的业务
 */
func handlerShuiguoji(args []interface{}) {
	log.T("进入到 game.handlerShuiguoji()")
	//检测参数是否正确
	m := args[0].(*bbproto.Shuiguoji)                //请求体
	a := args[1].(gate.Agent)
	result, err := fruitService.HandlerShuiguoji(m,a)
	if err != nil {
		log.E(err.Error())
	}

	//给客户端返回数据
	log.N("给客户端返回的数据%v",result)
	a.WriteMsg(result)
}


/**
	处理水果机比大小的业务
 */
func handlerShuiguojiHilomp(args []interface{}){
	log.T("进入到 game.handlerShuiguojiHilomp()")
	//检测参数是否正确
	m := args[0].(*bbproto.ShuiguojiHilomp)                //请求体
	a := args[1].(gate.Agent)
	result, err := fruitService.HandlerShuiguojiHilomp(m)
	if err != nil {
		log.E(err.Error())
	}
	a.WriteMsg(result)
	log.N("%v",result)

}

/**
	进入扎金花的房间
 */
func handlerZjhRoom(args []interface{}){
	log.T("进入到扎金花的房间 game.handlerZjhRoom()")
	//检测参数是否正确
	m := args[0].(*bbproto.ZjhRoom)                //请求体
	a := args[1].(gate.Agent)

	//通过serVice来处理
	result,err := zjhService.HandlerZjhRoom(m,a)
	if err != nil {
		log.E(err.Error())
	}

	//处理返回信息
	log.T("得到的结果:",result)

}


/**
	扎金花奖励

 */
func handlerZjhLottery(args []interface{}){
	log.T("进入到扎金花的房间 game.handlerZjhRoom()")
}


/**
	扎金花房间消息
 */
func handlerZjhMsg(args []interface{}){
	log.T("进入到扎金花的房间 game.handlerZjhMsg()")
}


/**
扎金花 押注
 */
func handlerZjhBet(args []interface{}){
	log.T("进入到扎金花的房间 game.handlerZjhBet()")
	//检测参数是否正确
	m := args[0].(*bbproto.ZjhBet)                //请求体
	a := args[1].(gate.Agent)
	zjhService.HandlerZjhBet(m,a)
}


/**
	扎金花请求座位
 */
func handlerZjhReqSeat(args []interface{}){
	log.T("进入到扎金花的房间 game.handlerZjhReqSeat()")
}


/*
	扎金花请求没有作为的玩家
 */


func handlerZjhQueryNoSeatUser(args []interface{}){
	log.T("进入到扎金花的房间 game.handlerZjhQueryNoSeatUser()")
}


/**
	退出或者进入德州扑克的房间
 */
func handlerThRoom(args []interface{}){
	log.T("handlerThRoom()")
	m := args[0].(*bbproto.ThRoom)                //请求体
	a := args[1].(gate.Agent)
	err := dzService.HandlerThRoom(m,a)
	if err != nil {
		log.E("请求进入或者退出德州扑克房间的时候出错了[%v]",err.Error())
	}
}


/**
	处理德州扑克押注
 */
func handlerThBet(args []interface{}){
	log.T("处理德州扑克押注")
	m := args[0].(*bbproto.THBet)                //请求体
	a := args[1].(gate.Agent)
	err := dzService.HandlerTHBet(m,a)
	if err != nil {
		log.E("请求进入或者退出德州扑克房间的时候出错了[%v]",err.Error())
	}

}

func handlerGameLoginGame(args []interface{}){
	log.T("快速登陆德州扑克游戏")
	m := args[0].(*bbproto.Game_LoginGame)
	log.T("收到的数据[%v]",m)

	a := args[1].(gate.Agent)
	result := &bbproto.Game_LoginGame{}
	result.Result = new(int32)		//默认是0表示成功
	a.WriteMsg(result)
}

// 处理请求进入游戏房间
func  handlerGameEnterMatch(args []interface{}){
	log.T("进入到-handlerGameEnterMatch()")
	m := args[0].(*bbproto.Game_EnterMatch)
	a := args[1].(gate.Agent)
	log.T("收到的数据-handlerGameEnterMatch()--[%v]",m)

	//返回房间的信息
	OGservice.HandlerGameEnterMatch(m,a)

}