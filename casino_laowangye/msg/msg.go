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

	Processor.Register(&ddproto.LwyChizhuDetailReq{})  //15 吃注详情req
	Processor.Register(&ddproto.LwyChizhuDetailAck{})  //16 吃注详情ack

	Processor.Register(&ddproto.LwyYaoshaiziOt{}) //17 摇色子ot
	Processor.Register(&ddproto.LwyYaoshaiziReq{}) //18 摇色子req

	Processor.Register(&ddproto.LwyGameEndOne{}) //19 单局结束-摇色子结果

	Processor.Register(&ddproto.LwyGameEndAll{})  //20 全局结束

	Processor.Register(&ddproto.CommonReqApplyDissolve{})  //21 申请解散房间
	Processor.Register(&ddproto.CommonBcApplyDissolve{})  //22 申请解散房间广播
	Processor.Register(&ddproto.CommonReqApplyDissolveBack{})  //23 确定、拒绝解散房间请求
	Processor.Register(&ddproto.CommonAckApplyDissolveBack{})  //24 确定、拒绝解散房间广播

	Processor.Register(&ddproto.LwyDeskDissolveDoneBc{})  //25 解散成功广播

	Processor.Register(&ddproto.LwyOwnerDissolveReq{})  //26 房主解散房间不扣房卡req
	Processor.Register(&ddproto.LwyOwnerDissolveAck{})  //27 房主解散房间不扣房卡ack

	Processor.Register(&ddproto.CommonReqMessage{})  //28 聊天请求
	Processor.Register(&ddproto.CommonBcMessage{})  //29 聊天广播

	Processor.Register(&ddproto.CommonReqLeaveDesk{})  //30 退出房间req
	Processor.Register(&ddproto.CommonAckLeaveDesk{})  //31 退出房间ack、bc

	Processor.Register(&ddproto.LwyOfflineBc{})  //32 离线广播

	//Processor.Register(&ddproto.CommonReqListCoinInfo{})  //33 金币场牌桌列表req
	//Processor.Register(&ddproto.CommonAckListCoinInfo{})  //34 金币场牌桌列表ack
	//
	//Processor.Register(&ddproto.LwyDeskDissolveDoneBc{})  //35 老王爷解散房间广播
	//
	//
	//Processor.Register(&ddproto.LwyCoinRoomListReq{})  //36房间列表req
	//Processor.Register(&ddproto.LwyCoinRoomListAck{})  //37房间列表ack

}
