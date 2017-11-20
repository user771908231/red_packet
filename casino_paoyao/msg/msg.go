package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_common/proto/ddproto"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&ddproto.Heartbeat{}) //0-连接服务器
	Processor.Register(&ddproto.QuickConn{})       //1-连接服务器
	Processor.Register(&ddproto.AckQuickConn{})    //2-连接服务器

	Processor.Register(&ddproto.CommonReqGameLogin{})           //3-连接服务器
	Processor.Register(&ddproto.CommonAckGameLogin{})        //4登录游戏,登录返回

	Processor.Register(&ddproto.PaoyaoCreateDeskReq{}) //5创建房间

	Processor.Register(&ddproto.PaoyaoEnterDeskReq{}) //6进入房间
	Processor.Register(&ddproto.PaoyaoEnterDeskAck{})  //7 进入房间ack
	Processor.Register(&ddproto.PaoyaoEnterDeskBc{}) //8 进入房间广播

	Processor.Register(&ddproto.PaoyaoSwitchReadyReq{})  //9 准备req
	Processor.Register(&ddproto.PaoyaoSwitchReadyBc{})  //11 准备bc

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
