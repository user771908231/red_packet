package game_handler


import (
	"casino_common/proto/ddproto"
	"github.com/name5566/leaf/gate"
	"casino_paoyao/service/paoyaoService"
)

//心跳
func HandlerHeartbeat(args []interface{}) {
	//m := args[0].(*casinoCommonProto.Heartbeat)
	a := args[1].(gate.Agent)
	ack := new(ddproto.Heartbeat)
	a.WriteMsg(ack)
}

//创建房间
func HandlerCreateDesk(args []interface{}) {
	m := args[0].(*ddproto.PaoyaoCreateDeskReq)
	a := args[1].(gate.Agent)
	paoyaoService.CreateDeskHandler(m, a)
}
//进入房间
func HandlerEnterDesk(args []interface{}) {
	m := args[0].(*ddproto.PaoyaoEnterDeskReq)
	a := args[1].(gate.Agent)
	paoyaoService.EnterDeskHandler(m, a)
}

//准备
func HandlerReadyReq(args []interface{})  {
	m := args[0].(*ddproto.PaoyaoSwitchReadyReq)
	a := args[1].(gate.Agent)
	paoyaoService.ReadyHandler(m, a)
}

//加倍请求
func HandlerJiabeiReq(args []interface{})  {
	m := args[0].(*ddproto.PaoyaoJiabeiReq)
	a := args[1].(gate.Agent)
	paoyaoService.JiabeiHandler(m, a)
}

//出牌
func HandlerChupaiReq(args []interface{})  {
	m := args[0].(*ddproto.PaoyaoChupaiReq)
	a := args[1].(gate.Agent)
	paoyaoService.ChupaiHandler(m, a)
}

//过牌
func HandlerGuopaiReq(args []interface{})  {
	m := args[0].(*ddproto.PaoyaoGuopaiReq)
	a := args[1].(gate.Agent)
	paoyaoService.GuopaiHandler(m, a)
}





//============================非游戏流程协议========================

//申请解散房间
func HandlerApplyDissolveReq(args []interface{})  {
	m := args[0].(*ddproto.CommonReqApplyDissolve)
	a := args[1].(gate.Agent)
	paoyaoService.ApplyDissolveReqHandler(m, a)
}

//确定、取消解散房间
func HandlerDissolveBackReq(args []interface{}) {
	m := args[0].(*ddproto.CommonReqApplyDissolveBack)
	a := args[1].(gate.Agent)
	paoyaoService.DissolveBackReqHandler(m, a)
}

//聊天请求
func HandlerMessageReq(args []interface{})  {
	m := args[0].(*ddproto.CommonReqMessage)
	a := args[1].(gate.Agent)
	paoyaoService.MessageReqHandler(m, a)
}

//离开房间
func HandlerLeaveDesk(args []interface{})  {
	m := args[0].(*ddproto.CommonReqLeaveDesk)
	a := args[1].(gate.Agent)
	paoyaoService.LeaveDeskReqHandler(m, a)
}

//金币场房间列表
func HandlerCoinRoomList(args []interface{})  {
	m := args[0].(*ddproto.PaoyaoCoinRoomListReq)
	a := args[1].(gate.Agent)
	paoyaoService.CoinRoomListHandler(m, a)
}
