package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_common/proto/ddproto"
	"casino_majiang/msg/protogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&ddproto.Heartbeat{})                    //0
	Processor.Register(&mjproto.Game_QuickConn{})               //1 接入服务器
	Processor.Register(&mjproto.Game_AckQuickConn{})            //2
	Processor.Register(&ddproto.CommonReqGameLogin{})           //3 登录游戏
	Processor.Register(&ddproto.CommonAckGameLogin{})           //4
	Processor.Register(&mjproto.Game_CreateRoom{})              //5 创建房间
	Processor.Register(&mjproto.Game_AckCreateRoom{})           //6
	Processor.Register(&mjproto.Game_EnterRoom{})               //7 进入房间
	Processor.Register(&mjproto.Game_AckEnterRoom{})            //8
	Processor.Register(&mjproto.Game_SendGameInfo{})            //9 卓内游戏数据
	Processor.Register(&mjproto.Game_Ready{})                   //10 准备
	Processor.Register(&mjproto.Game_AckReady{})                //11
	Processor.Register(&mjproto.Game_ExchangeCards{})           //12 换3张
	Processor.Register(&mjproto.Game_AckExchangeCards{})        //13 换3张-回复
	Processor.Register(&mjproto.Game_DingQue{})                 //14 定缺
	Processor.Register(&mjproto.Game_Opening{})                 //15 开始(表示都已经准备完了)
	Processor.Register(&mjproto.Game_DealCards{})               //16 发牌
	Processor.Register(&mjproto.Game_GetInCard{})               //17 摸牌
	Processor.Register(&mjproto.Game_SendOutCard{})             //18 出牌
	Processor.Register(&mjproto.Game_AckSendOutCard{})          //19 出牌-ack
	Processor.Register(&mjproto.Game_ActPeng{})                 //20 碰
	Processor.Register(&mjproto.Game_AckActPeng{})              //21
	Processor.Register(&mjproto.Game_ActGang{})                 //22 杠
	Processor.Register(&mjproto.Game_AckActGang{})              //23
	Processor.Register(&mjproto.Game_ActGuo{})                  //24 过
	Processor.Register(&mjproto.Game_AckActGuo{})               //25
	Processor.Register(&mjproto.Game_ActHu{})                   //26 胡
	Processor.Register(&mjproto.Game_AckActHu{})                //27
	Processor.Register(&mjproto.Game_BroadcastBeginDingQue{})   //28 开始定缺(广播)
	Processor.Register(&mjproto.Game_BroadcastBeginExchange{})  //29 开始换牌(广播)
	Processor.Register(&mjproto.Game_OverTurn{})                //30 轮到下一人
	Processor.Register(&mjproto.Game_SendCurrentResult{})       //31 本局结果
	Processor.Register(&mjproto.Game_SendEndLottery{})          //32 牌局结束
	Processor.Register(&mjproto.Game_DissolveDesk{})            //33 解散房间
	Processor.Register(&mjproto.Game_AckDissolveDesk{})         //34
	Processor.Register(&ddproto.CommonReqLeaveDesk{})           //35 离开房间
	Processor.Register(&ddproto.CommonAckLeaveDesk{})           //36
	Processor.Register(&ddproto.CommonReqMessage{})             //37 发送聊天消息
	Processor.Register(&ddproto.CommonBcMessage{})              //38 广播聊天
	Processor.Register(&mjproto.Game_DingQueEnd{})              //39 定缺结束
	Processor.Register(&mjproto.Game_GameRecord{})              //40 查询战绩
	Processor.Register(&mjproto.Game_AckGameRecord{})           //41 战绩回复
	Processor.Register(&mjproto.Game_ExchangeCardsEnd{})        //42 换三张 结束之后的广播
	Processor.Register(&ddproto.CommonReqNotice{})              //43通知信息
	Processor.Register(&ddproto.CommonAckNotice{})              //44通知信息回复
	Processor.Register(&ddproto.CommonReqLogout{})              //45请求推出
	Processor.Register(&ddproto.CommonAckLogout{})              //46回复请求推出
	Processor.Register(&ddproto.AwardReqOnline{})               //47 在线奖励
	Processor.Register(&ddproto.WardAckOnline{})                //48 在线奖励回复
	Processor.Register(&ddproto.HallReqTask{})                  //49 任务
	Processor.Register(&ddproto.HallAckTask{})                  //50 任务回复
	Processor.Register(&ddproto.CommonReqEnterAgentMode{})      //51 进入托管
	Processor.Register(&ddproto.CommonAckEnterAgentMode{})      //52 进入托管回复
	Processor.Register(&ddproto.CommonReqQuitAgentMode{})       //53 退出托管
	Processor.Register(&ddproto.CommonAckQuitAgentMode{})       //54 退出托管回复
	Processor.Register(&ddproto.CommonReqReg{})                 //55 注册
	Processor.Register(&ddproto.CommonAckReg{})                 //56 注册回复
	Processor.Register(&ddproto.CommonReqGameState{})           //57 玩家游戏状态
	Processor.Register(&ddproto.CommonAckGameState{})           //58 玩家游戏状态回复
	Processor.Register(&ddproto.CommonReqFeedback{})            //59 反馈
	Processor.Register(&ddproto.CommonReqApplyDissolve{})       //60 申请解散房间
	Processor.Register(&ddproto.CommonBcApplyDissolve{})        //61 申请解散房间回复
	Processor.Register(&ddproto.CommonReqApplyDissolveBack{})   //62 同意或拒绝解散房间回复
	Processor.Register(&ddproto.CommonAckApplyDissolveBack{})   //63 同意或拒绝解散房间回复
	Processor.Register(&ddproto.CommonBcKickout{})              //64 强制退出
	Processor.Register(&mjproto.Game_ActChi{})                  //65吃牌
	Processor.Register(&mjproto.Game_AckActChi{})               //66吃牌回复
	Processor.Register(&mjproto.Game_ChangShaAckActGang{})      //67长沙杠回复
	Processor.Register(&mjproto.Game_ActChangShaQiShouHu{})     //68长沙起手胡
	Processor.Register(&mjproto.Game_AckActChangShaQiShouHu{})  //69回复长沙起手胡
	Processor.Register(&mjproto.Game_ChangshQiShouHuOverTurn{}) //70回复长沙起手胡
}
