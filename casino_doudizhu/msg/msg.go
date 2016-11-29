package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_doudizhu/msg/protogo"
	"casino_common/proto"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&casinoCommonProto.Heartbeat{})        //0连接服务器 使用公共库的东西
	Processor.Register(&ddzproto.DdzQuickConn{})        //1连接服务器
	Processor.Register(&ddzproto.DdzAckQuickConn{})//2登录游戏

	Processor.Register(&ddzproto.DdzLogin{})//3
	Processor.Register(&ddzproto.DdzAckLogin{})//4
	Processor.Register(&ddzproto.DdzCreateRoom{})//5
	Processor.Register(&ddzproto.DdzAckCreateRoom{})//6
	Processor.Register(&ddzproto.DdzEnterRoom{})//7
	Processor.Register(&ddzproto.DdzAckEnterRoom{})//8
	Processor.Register(&ddzproto.DdzSendGameInfo{})//9
	Processor.Register(&ddzproto.DdzReady{})//10
	Processor.Register(&ddzproto.DdzAckReady{})//11
	Processor.Register(&ddzproto.DdzOpening{})//12
	Processor.Register(&ddzproto.DdzDealCards{})//13
	Processor.Register(&ddzproto.DdzJiaoDiZhu{})//14	PID_JIAO_DIZHU = 14; //叫地主
	Processor.Register(&ddzproto.DdzJiaoDiZhuAck{})//15	PID_JIAO_DIZHU_ACK = 15; //叫地主-ack


	//欢乐斗地主
	Processor.Register(&ddzproto.DdzRobDiZhu{})//16	PID_ROB_DIZHU = 16; //抢地主
	Processor.Register(&ddzproto.DdzRobDiZhuAck{})//17	PID_ROB_DIZHU_ACK = 17; //抢地主-ack
	Processor.Register(&ddzproto.DdzDouble{})//18
	Processor.Register(&ddzproto.DdzDoubleAck{})//19
	Processor.Register(&ddzproto.DdzShowHandPokers{})//20	//明牌
	Processor.Register(&ddzproto.DdzShowHandPokersAck{})//21

	//四川斗地主
	Processor.Register(&ddzproto.DdzMenuZhua{})//22
	Processor.Register(&ddzproto.DdzMenuZhuaAck{})//23
	Processor.Register(&ddzproto.DdzSeeCards{})//24
	Processor.Register(&ddzproto.DdzSeeCardsAck{})//25
	Processor.Register(&ddzproto.DdzPull{})//26
	Processor.Register(&ddzproto.DdzPullAck{})//27
	Processor.Register(&ddzproto.DdzOutCards{})//28
	Processor.Register(&ddzproto.DdzOutCardsAck{})//29
	Processor.Register(&ddzproto.DdzActGuo{})//30
	Processor.Register(&ddzproto.DdzActGuoAck{})//31
	Processor.Register(&ddzproto.DdzStartPlay{})//32
	Processor.Register(&ddzproto.DdzOverTurn{})//33
	Processor.Register(&ddzproto.DdzSendCurrentResult{})//34 本局结束
	Processor.Register(&ddzproto.DdzSendEndLottery{})//35牌局结束
	Processor.Register(&ddzproto.DdzDissolveDesk{})//36
	Processor.Register(&ddzproto.DdzAckDissolveDesk{})//37

	Processor.Register(&ddzproto.DdzLeaveDesk{})//38
	Processor.Register(&ddzproto.DdzAckLeaveDesk{})//39
	Processor.Register(&ddzproto.DdzMessage{})//40
	Processor.Register(&ddzproto.DdzSendMessage{})//41
	Processor.Register(&ddzproto.DdzGameRecord{})//42
	Processor.Register(&ddzproto.DdzAckGameRecord{})//43

}
