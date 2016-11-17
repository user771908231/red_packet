package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_doudizhu/msg/protogo"
	"casino_server/msg/bbprotogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&bbproto.NullMsg{})        //0连接服务器
	Processor.Register(&ddzproto.Game_QuickConn{})        //1连接服务器
	Processor.Register(&ddzproto.Game_AckQuickConn{})//2登录游戏

	Processor.Register(&ddzproto.Game_Login{})//3
	Processor.Register(&ddzproto.Game_AckLogin{})//4
	Processor.Register(&ddzproto.Game_CreateRoom{})//5
	Processor.Register(&ddzproto.Game_AckCreateRoom{})//6
	Processor.Register(&ddzproto.Game_EnterRoom{})//7
	Processor.Register(&ddzproto.Game_AckEnterRoom{})//8
	Processor.Register(&ddzproto.Game_SendGameInfo{})//9
	Processor.Register(&ddzproto.Game_Ready{})//10
	Processor.Register(&ddzproto.Game_AckReady{})//11
	Processor.Register(&ddzproto.Game_Opening{})//12
	Processor.Register(&ddzproto.Game_DealCards{})//13
	Processor.Register(&ddzproto.Game_JiaoDiZhu{})//14	PID_JIAO_DIZHU = 14; //叫地主
	Processor.Register(&ddzproto.Game_JiaoDiZhuAck{})//15	PID_JIAO_DIZHU_ACK = 15; //叫地主-ack


	//欢乐斗地主
	Processor.Register(&ddzproto.Game_RobDiZhu{})//16	PID_ROB_DIZHU = 16; //抢地主
	Processor.Register(&ddzproto.Game_RobDiZhuAck{})//17	PID_ROB_DIZHU_ACK = 17; //抢地主-ack
	Processor.Register(&ddzproto.Game_Double{})//18
	Processor.Register(&ddzproto.Game_DoubleAck{})//19
	Processor.Register(&ddzproto.Game_ShowHandPokers{})//20	//明牌
	Processor.Register(&ddzproto.Game_ShowHandPokersAck{})//21

	//四川斗地主
	Processor.Register(&ddzproto.Game_MenuZhua{})//22
	Processor.Register(&ddzproto.Game_MenuZhuaAck{})//23
	Processor.Register(&ddzproto.Game_SeeCards{})//24
	Processor.Register(&ddzproto.Game_SeeCardsAck{})//25
	Processor.Register(&ddzproto.Game_Pull{})//26
	Processor.Register(&ddzproto.Game_PullAck{})//27
	Processor.Register(&ddzproto.Game_OutCards{})//28
	Processor.Register(&ddzproto.Game_OutCardsAck{})//29
	Processor.Register(&ddzproto.Game_ActGuo{})//30
	Processor.Register(&ddzproto.Game_ActGuoAck{})//31
	Processor.Register(&ddzproto.Game_StartPlay{})//32
	Processor.Register(&ddzproto.Game_OverTurn{})//33
	Processor.Register(&ddzproto.Game_SendCurrentResult{})//34 本局结束
	Processor.Register(&ddzproto.Game_SendEndLottery{})//35牌局结束
	Processor.Register(&ddzproto.Game_DissolveDesk{})//36
	Processor.Register(&ddzproto.Game_AckDissolveDesk{})//37

	Processor.Register(&ddzproto.Game_LeaveDesk{})//38
	Processor.Register(&ddzproto.Game_AckLeaveDesk{})//39
	Processor.Register(&ddzproto.Game_Message{})//40
	Processor.Register(&ddzproto.Game_SendMessage{})//41
	Processor.Register(&ddzproto.Game_GameRecord{})//42
	Processor.Register(&ddzproto.Game_AckGameRecord{})//43
}
