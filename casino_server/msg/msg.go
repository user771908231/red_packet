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
	ProtoProcessor.Register(&bbproto.User{})                        //0
	ProtoProcessor.Register(&bbproto.GetIntoRoom{})                //1	进入房间时候的请求
	ProtoProcessor.Register(&bbproto.Shuiguoji{})                //2	水果机
	ProtoProcessor.Register(&bbproto.ShuiguojiHilomp{})        //3	水果机比大小
	ProtoProcessor.Register(&bbproto.ShuiguojiRes{})        //4	水果机的回应包

	//扎金花
	ProtoProcessor.Register(&bbproto.ZjhRoom{})                //5	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhBet{})                //6	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhMsg{})                //7	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhQueryNoSeatUser{})        //8	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhReqSeat{})                //9	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhLottery{})                //10	进入扎金花的房间
	ProtoProcessor.Register(&bbproto.ZjhBroadcastBeginBet{})//11 广播可以押注了


	//用户奖励相关的
	ProtoProcessor.Register(&bbproto.LoginSignInBonus{})        //12	登录签到奖励
	ProtoProcessor.Register(&bbproto.LoginTurntableBonus{})        //13	登录转盘奖励
	ProtoProcessor.Register(&bbproto.OlineBonus{})                //14	在线奖励
	ProtoProcessor.Register(&bbproto.TimingBonus{})                //15	定时奖励

	//德州扑克
	ProtoProcessor.Register(&bbproto.ThRoom{})                //16	德州扑克
	ProtoProcessor.Register(&bbproto.THBet{})                //17	德州扑克押注
	ProtoProcessor.Register(&bbproto.THBetBegin{})                //18	开始德州扑克的广播
	ProtoProcessor.Register(&bbproto.THBetBroadcast{})        //19	德州扑克,押注之后的广播
	ProtoProcessor.Register(&bbproto.THRoomAddUserBroadcast{})//20	房间增加一个玩家


	//联众游戏
	ProtoProcessor.Register(&bbproto.REQQuickConn{})        //21	//登陆
	ProtoProcessor.Register(&bbproto.ACKQuickConn{})        //22	登陆回复
	ProtoProcessor.Register(&bbproto.NullMsg{})                //23	空消息
	ProtoProcessor.Register(&bbproto.MatchList_ReqMobileMatchList{})        //24	快速开始游戏
	ProtoProcessor.Register(&bbproto.Game_LoginGame{})        //25	登陆游戏
	ProtoProcessor.Register(&bbproto.Game_EnterMatch{})        //26	进入房间
	ProtoProcessor.Register(&bbproto.Game_AckEnterMatch{})        //27
	ProtoProcessor.Register(&bbproto.Game_SendGameInfo{})        //28

	//开始游戏
	ProtoProcessor.Register(&bbproto.Game_BlindCoin{})        //29	盲注
	ProtoProcessor.Register(&bbproto.Game_InitCard{})        //30	手牌
	ProtoProcessor.Register(&bbproto.Game_SendFlopCard{})        //31	3张公共牌
	ProtoProcessor.Register(&bbproto.Game_SendTurnCard{})        //32	4张牌
	ProtoProcessor.Register(&bbproto.Game_SendRiverCard{})        //33	5张牌

	ProtoProcessor.Register(&bbproto.Game_RaiseBet{})        //34	加注
	ProtoProcessor.Register(&bbproto.Game_AckRaiseBet{})        //35	加注回复
	ProtoProcessor.Register(&bbproto.Game_FollowBet{})        //36	跟注
	ProtoProcessor.Register(&bbproto.Game_AckFollowBet{})        //37	跟注回复
	ProtoProcessor.Register(&bbproto.Game_FoldBet{})        //38	弃牌
	ProtoProcessor.Register(&bbproto.Game_AckFoldBet{})        //39	弃牌回复
	ProtoProcessor.Register(&bbproto.Game_CheckBet{})        //40	让牌
	ProtoProcessor.Register(&bbproto.Game_AckCheckBet{})        //41	让牌回复
	ProtoProcessor.Register(&bbproto.Game_SendOverTurn{})        //42	下一轮
	ProtoProcessor.Register(&bbproto.Game_SendAddUser{})        //43	新增用户
	ProtoProcessor.Register(&bbproto.Game_ShowCard{})        //44	请求开牌
	ProtoProcessor.Register(&bbproto.Game_AckShowCard{})        //45	回复 亮自己的手牌
	ProtoProcessor.Register(&bbproto.Game_TestResult{})        //46	一局结束之后,返回结果
	ProtoProcessor.Register(&bbproto.Game_PreCoin{})        //47	前注的协议号码
	ProtoProcessor.Register(&bbproto.Game_Notice{})                //48	请求公告
	ProtoProcessor.Register(&bbproto.Game_AckNotice{})        //49	回复请求的公告的协议
	ProtoProcessor.Register(&bbproto.Game_CreateDesk{})        //50	创建德州的游戏房间
	ProtoProcessor.Register(&bbproto.Game_AckCreateDesk{})        //51	创建德州的游戏房间

	ProtoProcessor.Register(&bbproto.Game_Ready{})                //52	准备游戏
	ProtoProcessor.Register(&bbproto.Game_AckReady{})        //53	//
	ProtoProcessor.Register(&bbproto.Game_Begin{})                //54 	开始游戏
	ProtoProcessor.Register(&bbproto.Game_GameRecord{})     //55 	查询战绩
	ProtoProcessor.Register(&bbproto.Game_AckGameRecord{})  //56 	查询战绩
	ProtoProcessor.Register(&bbproto.Game_BeanGameRecord{})        //57	战绩bean

	ProtoProcessor.Register(&bbproto.Game_DissolveDesk{})        //58	//解散房间
	ProtoProcessor.Register(&bbproto.Game_AckDissolveDesk{})//59	//解散房间回复

	ProtoProcessor.Register(&bbproto.Game_LeaveDesk{})        //60	//离开房间
	ProtoProcessor.Register(&bbproto.Game_ACKLeaveDesk{})        //61	//离开房间回复

	ProtoProcessor.Register(&bbproto.Game_SendDeskEndLottery{})        //62最终开奖的节奏
	ProtoProcessor.Register(&bbproto.Game_Message{})        //63	//发送信息
	ProtoProcessor.Register(&bbproto.Game_SendMessage{})        //64	//发送消息广播

	ProtoProcessor.Register(&bbproto.Game_TounamentBlind{})        //65
	ProtoProcessor.Register(&bbproto.Game_TounamentRewards{})//66
	ProtoProcessor.Register(&bbproto.Game_TounamentRank{})        //67
	ProtoProcessor.Register(&bbproto.Game_TounamentSummary{})        //68	//描述
	ProtoProcessor.Register(&bbproto.Game_MatchList{})        //69	竞标赛列表
	ProtoProcessor.Register(&bbproto.Game_TounamentPlayerRank{})        //70 每一局锦标赛完成之后的排名
	ProtoProcessor.Register(&bbproto.Game_Rebuy{})                //71
	ProtoProcessor.Register(&bbproto.Game_AckRebuy{})        //72
	ProtoProcessor.Register(&bbproto.Game_Login{})                //73	登陆大厅的协议
	ProtoProcessor.Register(&bbproto.Game_AckLogin{})        //74	回复的协议
	ProtoProcessor.Register(&bbproto.Game_Feedback{})        //75	反馈

	ProtoProcessor.Register(&bbproto.Game_NotRebuy{})        //76	反馈
	ProtoProcessor.Register(&bbproto.Game_AckNotRebuy{})        //77	反馈
	ProtoProcessor.Register(&bbproto.Game_SendChangeDeskOwner{})        //78	反馈
	ProtoProcessor.Register(&bbproto.Game_ChampionshipGameOver{})        //79	反馈
}
