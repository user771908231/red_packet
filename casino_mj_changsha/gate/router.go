package gate

import (
	"casino_mj_changsha/msg"
	"casino_mj_changsha/game"
	mjproto	"casino_mj_changsha/msg/protogo"
	"casino_mj_changsha/login"
	"casino_common/proto/ddproto"
	"casino_common/common/handlers"
)

func init() {
	msg.Processor.SetHandler(&ddproto.Heartbeat{}, handlers.HandlerHeartBeat)                   //心跳
	msg.Processor.SetRouter(&mjproto.Game_QuickConn{}, login.ChanRPC)                           //接入服务器
	msg.Processor.SetHandler(&ddproto.CommonAckReg{}, handlers.HandlerReg)                      //注册
	msg.Processor.SetRouter(&mjproto.Game_CreateRoom{}, game.ChanRPC)                           //创建房间
	msg.Processor.SetRouter(&mjproto.Game_EnterRoom{}, game.ChanRPC)                            //进入房间
	msg.Processor.SetRouter(&mjproto.Game_Ready{}, game.ChanRPC)                                //准备
	msg.Processor.SetRouter(&mjproto.Game_DingQue{}, game.ChanRPC)                              //定缺
	msg.Processor.SetRouter(&mjproto.Game_ExchangeCards{}, game.ChanRPC)                        //换3张
	msg.Processor.SetRouter(&mjproto.Game_SendOutCard{}, game.ChanRPC)                          //出牌
	msg.Processor.SetRouter(&mjproto.Game_ActPeng{}, game.ChanRPC)                              //碰
	msg.Processor.SetRouter(&mjproto.Game_ActGang{}, game.ChanRPC)                              //杠
	msg.Processor.SetRouter(&mjproto.Game_ActGuo{}, game.ChanRPC)                               //过
	msg.Processor.SetRouter(&mjproto.Game_ActHu{}, game.ChanRPC)                                //胡
	msg.Processor.SetRouter(&ddproto.CommonReqGameLogin{}, login.ChanRPC)                       //游戏登录
	msg.Processor.SetRouter(&mjproto.Game_DissolveDesk{}, game.ChanRPC)                         //解散房间
	msg.Processor.SetRouter(&mjproto.Game_GameRecord{}, game.ChanRPC)                           //游戏战绩
	msg.Processor.SetRouter(&ddproto.CommonReqMessage{}, game.ChanRPC)                          //发送聊天信息
	msg.Processor.SetHandler(&ddproto.CommonReqNotice{}, handlers.HandlerGetCommonAckNotice)    //麻将信息
	msg.Processor.SetRouter(&ddproto.CommonReqLogout{}, game.ChanRPC)                           //麻将信息
	msg.Processor.SetHandler(&ddproto.CommonReqGameState{}, handlers.HandlerCommonReqGameState) //游戏状态
	msg.Processor.SetHandler(&ddproto.CommonReqFeedback{}, handlers.HandlerFeedback)            //反馈
	msg.Processor.SetRouter(&ddproto.CommonReqLeaveDesk{}, game.ChanRPC)                        //换桌
	msg.Processor.SetRouter(&ddproto.CommonReqApplyDissolve{}, game.ChanRPC)                    //申请解散房间
	msg.Processor.SetRouter(&ddproto.CommonReqApplyDissolveBack{}, game.ChanRPC)                //申请解散房间
	msg.Processor.SetRouter(&ddproto.CommonReqEnterAgentMode{}, game.ChanRPC)                   //申请进入托管
	msg.Processor.SetRouter(&ddproto.CommonReqQuitAgentMode{}, game.ChanRPC)                    //申请退出托管模式
	msg.Processor.SetRouter(&mjproto.Game_ActChi{}, game.ChanRPC)                               //申请退出托管模式
	msg.Processor.SetRouter(&mjproto.Game_ActChangShaQiShouHu{}, game.ChanRPC)                  //长沙起手胡牌
	msg.Processor.SetRouter(&mjproto.Game_ReqDealHaiDiCards{}, game.ChanRPC)                    //是否需要海底牌
	msg.Processor.SetRouter(&ddproto.CommonReqReconnect{}, game.ChanRPC)                        //断线重连的处理
}
