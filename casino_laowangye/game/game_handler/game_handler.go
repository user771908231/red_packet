package game_handler


import (
	"casino_common/proto/ddproto"
	"github.com/name5566/leaf/gate"
	"casino_laowangye/service/laowangyeService"
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
	m := args[0].(*ddproto.LwyCreateDeskReq)
	a := args[1].(gate.Agent)
	a.WriteMsg(laowangyeService.CreateDeskHandler(m, a))
}
//进入房间
func HandlerEnterDesk(args []interface{}) {
	m := args[0].(*ddproto.LwyEnterDeskReq)
	a := args[1].(gate.Agent)
	laowangyeService.EnterDeskHandler(m, a)
}

//金币场房间列表
func HandlerCoinDeskList(args []interface{}) {
	m := args[0].(*ddproto.CommonReqListCoinInfo)
	a := args[1].(gate.Agent)
	laowangyeService.CoinDeskListHandler(m, a)
}

//入座
func HandlerSiteDownReq(args []interface{})  {
	m := args[0].(*ddproto.LwySiteDownReq)
	a := args[1].(gate.Agent)
	laowangyeService.SiteDownHandler(m, a)
}

//站起
func HandlerSiteUpReq(args []interface{})  {
	//m := args[0].(*ddproto.LwySiteUpReq)
	//a := args[1].(gate.Agent)
	//laowangyeService.SiteUpHandler(m, a)
}

//准备
func HandlerReadyReq(args []interface{})  {
	m := args[0].(*ddproto.LwySwitchReadyReq)
	a := args[1].(gate.Agent)
	laowangyeService.ReadyHandler(m, a)
}

//抢庄
func HandlerQiangzhuangReq(args []interface{})  {
	m := args[0].(*ddproto.LwyQiangzhuangReq)
	a := args[1].(gate.Agent)
	laowangyeService.QiangzhuangHandler(m, a)
}

//加倍
func HandlerYazhuReq(args []interface{})  {
	m := args[0].(*ddproto.LwyYazhuReq)
	a := args[1].(gate.Agent)
	laowangyeService.YazhuHandler(m, a)
}

//申请解散房间
func HandlerApplyDissolveReq(args []interface{})  {
	m := args[0].(*ddproto.CommonReqApplyDissolve)
	a := args[1].(gate.Agent)
	laowangyeService.ApplyDissolveReqHandler(m, a)
}

//确定、取消解散房间
func HandlerDissolveBackReq(args []interface{}) {
	m := args[0].(*ddproto.CommonReqApplyDissolveBack)
	a := args[1].(gate.Agent)
	laowangyeService.DissolveBackReqHandler(m, a)
}

//聊天请求
func HandlerMessageReq(args []interface{})  {
	m := args[0].(*ddproto.CommonReqMessage)
	a := args[1].(gate.Agent)
	laowangyeService.MessageReqHandler(m, a)
}

//离开房间
func HandlerLeaveDesk(args []interface{})  {
	m := args[0].(*ddproto.CommonReqLeaveDesk)
	a := args[1].(gate.Agent)
	laowangyeService.LeaveDeskReqHandler(m, a)
}

//房主解散房间不扣房卡
func HandlerOwnerDissolveReq(args []interface{})  {
	m := args[0].(*ddproto.LwyOwnerDissolveReq)
	a := args[1].(gate.Agent)
	laowangyeService.OwnerDissolveHandler(m, a)
}

//金币场房间列表
func HandlerCoinRoomList(args []interface{})  {
	m := args[0].(*ddproto.LwyCoinRoomListReq)
	a := args[1].(gate.Agent)
	laowangyeService.CoinRoomListHandler(m, a)
}
