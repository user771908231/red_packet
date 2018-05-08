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

	Processor.Register(&ddproto.NiuCreateDeskReq{}) //5创建房间

	Processor.Register(&ddproto.NiuEnterDeskReq{}) //6进入房间
	Processor.Register(&ddproto.NiuEnterDeskAck{})  //7 进入房间ack
	Processor.Register(&ddproto.NiuEnterDeskBc{}) //8 进入房间广播

	Processor.Register(&ddproto.NiuSwitchReadyReq{})  //9 准备req
	Processor.Register(&ddproto.NiuSwitchReadyAck{})  //10 准备ack
	Processor.Register(&ddproto.NiuSwitchReadyBc{})  //11 准备bc

	Processor.Register(&ddproto.NiuStartGameOt{})  //12 通知房主开局

	Processor.Register(&ddproto.NiuQiangzhuangOt{})  //13 抢庄ot
	Processor.Register(&ddproto.NiuQiangzhuangReq{})  //14 抢庄req
	Processor.Register(&ddproto.NiuQiangzhuangAck{})  //15 抢庄ack
	Processor.Register(&ddproto.NiuQiangzhuangResBc{})  //16 抢庄bc

	Processor.Register(&ddproto.NiuJiabeiOt{})  //17 加倍ot
	Processor.Register(&ddproto.NiuJiabeiReq{})  //18 加倍req
	Processor.Register(&ddproto.NiuJiabeiAck{})  //19 加倍ack
	Processor.Register(&ddproto.NiuJiabeiBc{})  //20 加倍bc

	Processor.Register(&ddproto.NiuBipaiResultBc{}) //21 比牌结果

	Processor.Register(&ddproto.NiuGameEnd{})  //22 游戏结束，统计数据

	Processor.Register(&ddproto.CommonReqApplyDissolve{})  //23 申请解散房间
	Processor.Register(&ddproto.CommonBcApplyDissolve{})  //24 申请解散房间广播
	Processor.Register(&ddproto.CommonReqApplyDissolveBack{})  //25 确定、拒绝解散房间请求
	Processor.Register(&ddproto.CommonAckApplyDissolveBack{})  //26 确定、拒绝解散房间广播
	Processor.Register(&ddproto.NiuDeskDissolveDoneBc{})  //27 牛牛解散房间广播

	Processor.Register(&ddproto.CommonReqMessage{})  //28 聊天请求
	Processor.Register(&ddproto.CommonBcMessage{})  //29 聊天广播

	Processor.Register(&ddproto.CommonReqLeaveDesk{})  //30 退出房间req
	Processor.Register(&ddproto.CommonAckLeaveDesk{})  //31 退出房间ack、bc

	Processor.Register(&ddproto.NiuOwnerDissolveReq{})  //32 房主解散房间不扣房卡req
	Processor.Register(&ddproto.NiuOwnerDissolveAck{})  //33 房主解散房间不扣房卡ack

	Processor.Register(&ddproto.NiuOfflineBc{})  //34 离线广播

	Processor.Register(&ddproto.NiuCoinRoomListReq{})  //35房间列表req
	Processor.Register(&ddproto.NiuCoinRoomListAck{})  //36房间列表ack

	Processor.Register(&ddproto.NiuSiteDownReq{})  //37 入座
	Processor.Register(&ddproto.NiuSiteUpReq{})  //38 站起

	Processor.Register(&ddproto.CommonReqListCoinInfo{})  //39 金币场牌桌列表req
	Processor.Register(&ddproto.CommonAckListCoinInfo{})  //40 金币场牌桌列表ack

	Processor.Register(&ddproto.NiuLiangpaiOt{})  //41 亮牌ot
	Processor.Register(&ddproto.NiuLiangpaiReq{})  //42 亮牌req
	Processor.Register(&ddproto.NiuLiangpaiAck{})  //43 亮牌ack
	Processor.Register(&ddproto.NiuLiangpaiBc{})  //44 亮牌bc

	Processor.Register(&ddproto.NiuTuoguanReq{})  //45 托管req
	Processor.Register(&ddproto.NiuTuoguanBc{})  //46 托管ack/bc
}
