package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_common/proto/ddproto"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&ddproto.Heartbeat{}) //0-连接服务器

	Processor.Register(&ddproto.LwyCreateDeskReq{}) //1 创建房间

	Processor.Register(&ddproto.LwyEnterDeskReq{}) //2 进入房间req
	Processor.Register(&ddproto.LwyEnterDeskAck{})  //3 进入房间ack

	Processor.Register(&ddproto.LwySiteDownReq{}) //4 入座
	Processor.Register(&ddproto.LwySiteDownBc{}) //5 入座bc

	Processor.Register(&ddproto.LwySwitchReadyReq{})  //6 准备req
	Processor.Register(&ddproto.LwySwitchReadyBc{})  //7 准备bc

	Processor.Register(&ddproto.LwyStartGameOt{})  //8 房主开局ot

	Processor.Register(&ddproto.LwyQiangzhuangOt{})  //9 抢庄ot
	Processor.Register(&ddproto.LwyQiangzhuangReq{})  //10 抢庄req
	Processor.Register(&ddproto.LwyQiangzhuangBc{})  //11 抢庄ack

	Processor.Register(&ddproto.LwyYazhuOt{})  //12 押注ot
	Processor.Register(&ddproto.LwyYazhuReq{})  //13 押注req
	Processor.Register(&ddproto.LwyYazhuBc{})  //14 押注bc

	Processor.Register(&ddproto.LwyYaoshaiziOt{}) //15 摇色子ot
	Processor.Register(&ddproto.LwyYaoshaiziReq{}) //16 摇色子req

	Processor.Register(&ddproto.LwyGameEndOne{}) //17 单局结束-摇色子结果

	Processor.Register(&ddproto.LwyGameEndAll{})  //18 全局结束

	Processor.Register(&ddproto.CommonReqApplyDissolve{})  //19 申请解散房间
	Processor.Register(&ddproto.CommonBcApplyDissolve{})  //20 申请解散房间广播
	Processor.Register(&ddproto.CommonReqApplyDissolveBack{})  //21 确定、拒绝解散房间请求
	Processor.Register(&ddproto.CommonAckApplyDissolveBack{})  //22 确定、拒绝解散房间广播

	Processor.Register(&ddproto.LwyOwnerDissolveReq{})  //23 房主解散房间不扣房卡req
	Processor.Register(&ddproto.LwyOwnerDissolveAck{})  //24 房主解散房间不扣房卡ack

	Processor.Register(&ddproto.CommonReqMessage{})  //25 聊天请求
	Processor.Register(&ddproto.CommonBcMessage{})  //26 聊天广播

	Processor.Register(&ddproto.CommonReqLeaveDesk{})  //27 退出房间req
	Processor.Register(&ddproto.CommonAckLeaveDesk{})  //28 退出房间ack、bc

	Processor.Register(&ddproto.LwyOfflineBc{})  //29 离线广播

	//Processor.Register(&ddproto.CommonReqListCoinInfo{})  //39 金币场牌桌列表req
	//Processor.Register(&ddproto.CommonAckListCoinInfo{})  //40 金币场牌桌列表ack
	//
	//Processor.Register(&ddproto.LwyDeskDissolveDoneBc{})  //22 老王爷解散房间广播
	//
	//
	//Processor.Register(&ddproto.LwyCoinRoomListReq{})  //35房间列表req
	//Processor.Register(&ddproto.LwyCoinRoomListAck{})  //36房间列表ack

}
