package game_handler


import (
	"casino_common/proto/ddproto"
	"github.com/name5566/leaf/gate"
	"casino_paosangong/service/paosangongService"
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
	m := args[0].(*ddproto.NiuCreateDeskReq)
	a := args[1].(gate.Agent)
	paosangongService.CreateDeskHandler(m, a)
}
//进入房间
func HandlerEnterDesk(args []interface{}) {
	m := args[0].(*ddproto.NiuEnterDeskReq)
	a := args[1].(gate.Agent)
	paosangongService.EnterDeskHandler(m, a)
}

//金币场房间列表
func HandlerCoinDeskList(args []interface{}) {
	m := args[0].(*ddproto.CommonReqListCoinInfo)
	a := args[1].(gate.Agent)
	paosangongService.CoinDeskListHandler(m, a)
}

//入座
func HandlerSiteDownReq(args []interface{})  {
	m := args[0].(*ddproto.NiuSiteDownReq)
	a := args[1].(gate.Agent)
	paosangongService.SiteDownHandler(m, a)
}

//站起
func HandlerSiteUpReq(args []interface{})  {
	m := args[0].(*ddproto.NiuSiteUpReq)
	a := args[1].(gate.Agent)
	paosangongService.SiteUpHandler(m, a)
}

//准备
func HandlerReadyReq(args []interface{})  {
	m := args[0].(*ddproto.NiuSwitchReadyReq)
	a := args[1].(gate.Agent)
	paosangongService.ReadyHandler(m, a)
}

//抢庄
func HandlerQiangzhuangReq(args []interface{})  {
	m := args[0].(*ddproto.NiuQiangzhuangReq)
	a := args[1].(gate.Agent)
	paosangongService.QiangzhuangHandler(m, a)
}

//加倍
func HandlerJiabeiReq(args []interface{})  {
	m := args[0].(*ddproto.NiuJiabeiReq)
	a := args[1].(gate.Agent)
	paosangongService.JiabeiHandler(m, a)
}

//亮牌
func HandlerLiangpaiReq(args []interface{})  {
	m := args[0].(*ddproto.NiuLiangpaiReq)
	a := args[1].(gate.Agent)
	paosangongService.LiangpaiHandler(m, a)
}

//托管
func HandlerTuoguanReq(args []interface{})  {
	m := args[0].(*ddproto.NiuTuoguanReq)
	a := args[1].(gate.Agent)
	paosangongService.TuoguanHandler(m, a)
}


//申请解散房间
func HandlerApplyDissolveReq(args []interface{})  {
	m := args[0].(*ddproto.CommonReqApplyDissolve)
	a := args[1].(gate.Agent)
	paosangongService.ApplyDissolveReqHandler(m, a)
}

//确定、取消解散房间
func HandlerDissolveBackReq(args []interface{}) {
	m := args[0].(*ddproto.CommonReqApplyDissolveBack)
	a := args[1].(gate.Agent)
	paosangongService.DissolveBackReqHandler(m, a)
}

//聊天请求
func HandlerMessageReq(args []interface{})  {
	m := args[0].(*ddproto.CommonReqMessage)
	a := args[1].(gate.Agent)
	paosangongService.MessageReqHandler(m, a)
}

//离开房间
func HandlerLeaveDesk(args []interface{})  {
	m := args[0].(*ddproto.CommonReqLeaveDesk)
	a := args[1].(gate.Agent)
	paosangongService.LeaveDeskReqHandler(m, a)
}

//房主解散房间不扣房卡
func HandlerOwnerDissolveReq(args []interface{})  {
	m := args[0].(*ddproto.NiuOwnerDissolveReq)
	a := args[1].(gate.Agent)
	paosangongService.OwnerDissolveHandler(m, a)
}

//金币场房间列表
func HandlerCoinRoomList(args []interface{})  {
	m := args[0].(*ddproto.NiuCoinRoomListReq)
	a := args[1].(gate.Agent)
	paosangongService.CoinRoomListHandler(m, a)
}
