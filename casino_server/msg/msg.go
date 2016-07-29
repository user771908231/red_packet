package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_server/msg/bbprotogo"
)

// protobuf 消息处理器
var ProtoProcessor = protobuf.NewProcessor()

func init() {
	// 这里我们注册了一个 JSON 消息 Hello

	//次处注册proto 的消息
	ProtoProcessor.Register(&bbproto.TestP1{})		//0	测试用
	ProtoProcessor.Register(&bbproto.Reg{})			//1	注册协议(已经废弃)
	ProtoProcessor.Register(&bbproto.ReqAuthUser{})		//2	登陆、注册的协议
	ProtoProcessor.Register(&bbproto.HeatBeat{})		//3	心跳协议,检测网络是否联通
	ProtoProcessor.Register(&bbproto.GetIntoRoom{})		//4	进入房间时候的请求
	ProtoProcessor.Register(&bbproto.RoomMsg{})		//5	给指定房间发送信息
	ProtoProcessor.Register(&bbproto.GetRewards{})		//6	各种奖励
	ProtoProcessor.Register(&bbproto.Shuiguoji{})		//7	水果机
	ProtoProcessor.Register(&bbproto.ShuiguojiHilomp{})	//8	水果机比大小
	ProtoProcessor.Register(&bbproto.ShuiguojiRes{})	//9	水果机的回应包

	//扎金花
	ProtoProcessor.Register(&bbproto.ZjhRoom{})		//10	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhBet{})		//11	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhMsg{})		//12	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhQueryNoSeatUser{})	//13	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhReqSeat{})		//14	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhLottery{})		//15	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhBroadcastBeginBet{})	//16 广播可以押注了


	//用户奖励相关的
	ProtoProcessor.Register(&bbproto.LoginSignInBonus{})	//17	登录签到奖励
	ProtoProcessor.Register(&bbproto.LoginTurntableBonus{})	//18	登录转盘奖励
	ProtoProcessor.Register(&bbproto.OlineBonus{})		//19	在线奖励
	ProtoProcessor.Register(&bbproto.TimingBonus{})		//20	定时奖励


	//德州扑克
	ProtoProcessor.Register(&bbproto.ThRoom{})		//21	德州扑克
	ProtoProcessor.Register(&bbproto.THBet{})		//22	德州扑克押注
	ProtoProcessor.Register(&bbproto.THBetBegin{})		//23	开始德州扑克的广播
	ProtoProcessor.Register(&bbproto.THBetBroadcast{})	//24	德州扑克,押注之后的广播
	ProtoProcessor.Register(&bbproto.THRoomAddUserBroadcast{})	//25	房间增加一个玩家


	//联众游戏
	ProtoProcessor.Register(&bbproto.REQQuickConn{})	//26	//登陆
	ProtoProcessor.Register(&bbproto.ACKQuickConn{})	//27	登陆回复
	ProtoProcessor.Register(&bbproto.NullMsg{})		//28	空消息
	ProtoProcessor.Register(&bbproto.MatchList_ReqMobileMatchList{})	//29	快速开始游戏
	ProtoProcessor.Register(&bbproto.Game_LoginGame{})	//30	登陆游戏
	ProtoProcessor.Register(&bbproto.Game_EnterMatch{})	//31	进入房间
	ProtoProcessor.Register(&bbproto.Game_AckEnterMatch{})	//32
	ProtoProcessor.Register(&bbproto.Game_SendGameInfo{})	//33

	//开始游戏
	ProtoProcessor.Register(&bbproto.Game_BlindCoin{})	//34	盲注
	ProtoProcessor.Register(&bbproto.Game_InitCard{})	//35	手牌
	ProtoProcessor.Register(&bbproto.Game_SendFlopCard{})	//36	3张公共牌
	ProtoProcessor.Register(&bbproto.Game_SendTurnCard{})	//37	4张牌
	ProtoProcessor.Register(&bbproto.Game_SendRiverCard{})	//38	5张牌

	ProtoProcessor.Register(&bbproto.Game_RaiseBet{})	//39	加注
	ProtoProcessor.Register(&bbproto.Game_AckRaiseBet{})	//40	加注回复
	ProtoProcessor.Register(&bbproto.Game_FollowBet{})	//41	跟注
	ProtoProcessor.Register(&bbproto.Game_AckFollowBet{})	//42	跟注回复
	ProtoProcessor.Register(&bbproto.Game_FoldBet{})	//43	弃牌
	ProtoProcessor.Register(&bbproto.Game_AckFoldBet{})	//44	弃牌回复
	ProtoProcessor.Register(&bbproto.Game_CheckBet{})	//45	让牌
	ProtoProcessor.Register(&bbproto.Game_AckCheckBet{})	//46	让牌回复
	ProtoProcessor.Register(&bbproto.Game_SendOverTurn{})	//47	下一轮
	ProtoProcessor.Register(&bbproto.Game_SendAddUser{})	//48	新增用户
	ProtoProcessor.Register(&bbproto.Game_ShowCard{})	//49	请求开牌
	ProtoProcessor.Register(&bbproto.Game_AckShowCard{})	//50	回复 亮自己的手牌
	ProtoProcessor.Register(&bbproto.Game_TestResult{})	//51	一局结束之后,返回结果
}
