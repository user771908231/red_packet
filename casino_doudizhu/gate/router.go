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
	msg.Processor.SetRouter(&ddzproto.Game_Login{}, login.ChanRPC)    //登录协议

	//游戏
	msg.Processor.SetRouter(&ddzproto.Game_CreateRoom{}, game.ChanRPC)    //创建房间
	msg.Processor.SetRouter(&ddzproto.Game_EnterRoom{}, game.ChanRPC)    //进入房间
	msg.Processor.SetRouter(&ddzproto.Game_Ready{}, game.ChanRPC)    //准备
	msg.Processor.SetRouter(&ddzproto.Game_JiaoDiZhu{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_RobDiZhu{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_Double{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_ShowHandPokers{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_MenuZhua{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_SeeCards{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_Pull{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_OutCards{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_ActGuo{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_DissolveDesk{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_LeaveDesk{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_Message{}, game.ChanRPC)    //
	msg.Processor.SetRouter(&ddzproto.Game_GameRecord{}, game.ChanRPC)    //
}
