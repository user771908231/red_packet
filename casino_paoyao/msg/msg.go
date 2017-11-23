package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_common/proto/ddproto"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&ddproto.Heartbeat{}) //0-连接服务器

	Processor.Register(&ddproto.PaoyaoCreateDeskReq{}) //1 创建房间

	Processor.Register(&ddproto.PaoyaoEnterDeskReq{}) //2进入房间
	Processor.Register(&ddproto.PaoyaoEnterDeskAck{})  //3 进入房间ack
	Processor.Register(&ddproto.PaoyaoEnterDeskBc{}) //4 进入房间广播

	Processor.Register(&ddproto.PaoyaoSwitchReadyReq{})  //5 准备req
	Processor.Register(&ddproto.PaoyaoSwitchReadyBc{})  //6 准备bc

	Processor.Register(&ddproto.PaoyaoFapaiOt{})  //7 开局发牌overturn

	Processor.Register(&ddproto.PaoyaoJiabeiReq{})  //8 加倍req
	Processor.Register(&ddproto.PaoyaoJiabeiBc{})  //9 加倍bc

	Processor.Register(&ddproto.PaoyaoChupaiOt{})  //10 出牌overturn
	Processor.Register(&ddproto.PaoyaoChupaiReq{})  //11 出牌req
	Processor.Register(&ddproto.PaoyaoChupaiBc{})  //12 出牌bc

	Processor.Register(&ddproto.PaoyaoGuopaiReq{})  //13 过牌req
	Processor.Register(&ddproto.PaoyaoGuopaiBc{})  //14 过牌bc

	Processor.Register(&ddproto.PaoyaoGameEndOneBc{})  //15 单局结束
	Processor.Register(&ddproto.PaoyaoGameEndAllBc{})  //16 全局结算

	//========================非游戏流程协议=======================

	Processor.Register(&ddproto.CommonReqApplyDissolve{})  //23 申请解散房间
	Processor.Register(&ddproto.CommonBcApplyDissolve{})  //24 申请解散房间广播
	Processor.Register(&ddproto.CommonReqApplyDissolveBack{})  //25 确定、拒绝解散房间请求
	Processor.Register(&ddproto.CommonAckApplyDissolveBack{})  //26 确定、拒绝解散房间广播

	Processor.Register(&ddproto.CommonReqMessage{})  //28 聊天请求
	Processor.Register(&ddproto.CommonBcMessage{})  //29 聊天广播

	Processor.Register(&ddproto.CommonReqLeaveDesk{})  //30 退出房间req
	Processor.Register(&ddproto.CommonAckLeaveDesk{})  //31 退出房间ack、bc

	Processor.Register(&ddproto.PaoyaoOfflineBc{})  //34 离线广播

	Processor.Register(&ddproto.PaoyaoCoinRoomListReq{})  //35房间列表req
	Processor.Register(&ddproto.PaoyaoCoinRoomListAck{})  //36房间列表ack

}
