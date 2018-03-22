package gate

import (
	"casino_laowangye/msg"
	"casino_common/proto/ddproto"
	"casino_laowangye/game/game_handler"
)

func init() {
	msg.Processor.SetHandler(&ddproto.Heartbeat{}, game_handler.HandlerHeartbeat) //心跳
	msg.Processor.SetHandler(&ddproto.LwyCreateDeskReq{}, game_handler.HandlerCreateDesk) //创建房间

	msg.Processor.SetHandler(&ddproto.LwyEnterDeskReq{}, game_handler.HandlerEnterDesk) //进入房间

	msg.Processor.SetHandler(&ddproto.LwySwitchReadyReq{}, game_handler.HandlerReadyReq)  //准备

	msg.Processor.SetHandler(&ddproto.LwyQiangzhuangReq{}, game_handler.HandlerQiangzhuangReq)  //抢庄
	msg.Processor.SetHandler(&ddproto.LwyYazhuReq{}, game_handler.HandlerYazhuReq)  //押注

	msg.Processor.SetHandler(&ddproto.LwyChizhuDetailReq{}, game_handler.HandlerChizhuDetailReq)  //吃注详情

	msg.Processor.SetHandler(&ddproto.LwyYaoshaiziReq{}, game_handler.HandlerYaoshaiziReq)  //摇色子

	msg.Processor.SetHandler(&ddproto.CommonReqApplyDissolve{}, game_handler.HandlerApplyDissolveReq)  //解散房间申请
	msg.Processor.SetHandler(&ddproto.CommonReqApplyDissolveBack{}, game_handler.HandlerDissolveBackReq)  //确定、取消解散房间请求

	msg.Processor.SetHandler(&ddproto.CommonReqMessage{}, game_handler.HandlerMessageReq)  //聊天请求

	msg.Processor.SetHandler(&ddproto.CommonReqLeaveDesk{}, game_handler.HandlerLeaveDesk)  //离开房间

	msg.Processor.SetHandler(&ddproto.LwyOwnerDissolveReq{}, game_handler.HandlerOwnerDissolveReq)  //房主解散房间不扣房卡

	msg.Processor.SetHandler(&ddproto.LwySiteDownReq{}, game_handler.HandlerSiteDownReq)  //入座
}
