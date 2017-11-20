package gate

import (
	"casino_paoyao/msg"
	"casino_common/proto/ddproto"
	"casino_common/common/handlers"
	"casino_paoyao/game/game_handler"
)

func init() {
	msg.Processor.SetHandler(&ddproto.Heartbeat{}, game_handler.HandlerHeartbeat) //心跳
	msg.Processor.SetHandler(&ddproto.CommonReqGameLogin{}, handlers.HandlerGame_Login)          //登录

	msg.Processor.SetHandler(&ddproto.PaoyaoCreateDeskReq{}, game_handler.HandlerCreateDesk) //创建房间
	msg.Processor.SetHandler(&ddproto.PaoyaoEnterDeskReq{}, game_handler.HandlerEnterDesk) //进入房间

	msg.Processor.SetHandler(&ddproto.PaoyaoSwitchReadyReq{}, game_handler.HandlerReadyReq)  //准备


	msg.Processor.SetHandler(&ddproto.CommonReqApplyDissolve{}, game_handler.HandlerApplyDissolveReq)  //解散房间申请
	msg.Processor.SetHandler(&ddproto.CommonReqApplyDissolveBack{}, game_handler.HandlerDissolveBackReq)  //确定、取消解散房间请求

	msg.Processor.SetHandler(&ddproto.CommonReqMessage{}, game_handler.HandlerMessageReq)  //聊天请求

	msg.Processor.SetHandler(&ddproto.CommonReqLeaveDesk{}, game_handler.HandlerLeaveDesk)  //离开房间

	msg.Processor.SetHandler(&ddproto.PaoyaoCoinRoomListReq{}, game_handler.HandlerCoinRoomList)  //金币场进房列表
}
