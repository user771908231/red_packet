package gate

import (
	"casino_doudizhu/login"
	"casino_doudizhu/msg"
	"casino_doudizhu/game"
	"casino_server/msg/bbprotogo"
	"casino_doudizhu/msg/protogo"
)

func init() {
	msg.Processor.SetRouter(&bbproto.NullMsg{}, game.ChanRPC)        //空协议
	msg.Processor.SetRouter(&ddzproto.DdzLogin{}, login.ChanRPC)    //登录协议

	//游戏
	msg.Processor.SetRouter(&ddzproto.DdzCreateRoom{}, game.ChanRPC)    //创建房间
	msg.Processor.SetRouter(&ddzproto.DdzEnterRoom{}, game.ChanRPC)    //进入房间
	msg.Processor.SetRouter(&ddzproto.DdzReady{}, game.ChanRPC)    //准备
	msg.Processor.SetRouter(&ddzproto.DdzJiaoDiZhu{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzRobDiZhu{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzDouble{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzShowHandPokers{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzMenuZhua{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzSeeCards{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzPull{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzOutCards{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzActGuo{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzDissolveDesk{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzLeaveDesk{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzMessage{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.DdzGameRecord{}, game.ChanRPC)    //
}
