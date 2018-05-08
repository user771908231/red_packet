package gate

import (
	"casino_paosangong/msg"
	"casino_common/proto/ddproto"
	"casino_common/common/handlers"
	"casino_paosangong/game/game_handler"
)

func init() {
	msg.Processor.SetHandler(&ddproto.Heartbeat{}, game_handler.HandlerHeartbeat) //心跳
	msg.Processor.SetHandler(&ddproto.CommonReqGameLogin{}, handlers.HandlerGame_Login)          //登录

	msg.Processor.SetHandler(&ddproto.NiuCreateDeskReq{}, game_handler.HandlerCreateDesk) //创建房间
	msg.Processor.SetHandler(&ddproto.NiuEnterDeskReq{}, game_handler.HandlerEnterDesk) //进入房间

	msg.Processor.SetHandler(&ddproto.NiuSwitchReadyReq{}, game_handler.HandlerReadyReq)  //准备

	msg.Processor.SetHandler(&ddproto.NiuQiangzhuangReq{}, game_handler.HandlerQiangzhuangReq)  //抢庄
	msg.Processor.SetHandler(&ddproto.NiuJiabeiReq{}, game_handler.HandlerJiabeiReq)  //加倍

	msg.Processor.SetHandler(&ddproto.CommonReqApplyDissolve{}, game_handler.HandlerApplyDissolveReq)  //解散房间申请
	msg.Processor.SetHandler(&ddproto.CommonReqApplyDissolveBack{}, game_handler.HandlerDissolveBackReq)  //确定、取消解散房间请求

	msg.Processor.SetHandler(&ddproto.CommonReqMessage{}, game_handler.HandlerMessageReq)  //聊天请求

	msg.Processor.SetHandler(&ddproto.CommonReqLeaveDesk{}, game_handler.HandlerLeaveDesk)  //离开房间

	msg.Processor.SetHandler(&ddproto.NiuOwnerDissolveReq{}, game_handler.HandlerOwnerDissolveReq)  //房主解散房间不扣房卡

	msg.Processor.SetHandler(&ddproto.NiuCoinRoomListReq{}, game_handler.HandlerCoinRoomList)  //金币场进房列表

	msg.Processor.SetHandler(&ddproto.NiuSiteDownReq{}, game_handler.HandlerSiteDownReq)  //入座
	msg.Processor.SetHandler(&ddproto.NiuSiteUpReq{}, game_handler.HandlerSiteUpReq)  //站起

	msg.Processor.SetHandler(&ddproto.CommonReqListCoinInfo{}, game_handler.HandlerCoinDeskList)  //金币场房间列表

	msg.Processor.SetHandler(&ddproto.NiuLiangpaiReq{}, game_handler.HandlerLiangpaiReq)  //亮牌

	msg.Processor.SetHandler(&ddproto.NiuTuoguanReq{}, game_handler.HandlerTuoguanReq)  //托管
}
