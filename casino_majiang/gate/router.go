package gate

import (
	"casino_majiang/msg"
	"casino_server/msg/bbprotogo"
	"casino_majiang/game"
	"casino_majiang/msg/protogo"
	"casino_majiang/login"
)

func init() {
	msg.Processor.SetRouter(&bbproto.NullMsg{}, game.ChanRPC)
	msg.Processor.SetRouter(&bbproto.REQQuickConn{}, login.ChanRPC)
	msg.Processor.SetRouter(&mjproto.Game_CreateRoom{}, game.ChanRPC)        //创建房间
	msg.Processor.SetRouter(&mjproto.Game_EnterRoom{}, game.ChanRPC)                //进入房间
	msg.Processor.SetRouter(&mjproto.Game_Ready{}, game.ChanRPC)                        //准备

	msg.Processor.SetRouter(&mjproto.Game_DingQue{}, game.ChanRPC)                //定缺
	msg.Processor.SetRouter(&mjproto.Game_ExchangeCards{}, game.ChanRPC)                //换3张

	msg.Processor.SetRouter(&mjproto.Game_SendOutCard{}, game.ChanRPC)  //出牌

	msg.Processor.SetRouter(&mjproto.Game_ActPeng{}, game.ChanRPC)                //碰
	msg.Processor.SetRouter(&mjproto.Game_ActGang{}, game.ChanRPC)                //杠
	msg.Processor.SetRouter(&mjproto.Game_ActGuo{}, game.ChanRPC)                //过
	msg.Processor.SetRouter(&mjproto.Game_ActHu{}, game.ChanRPC)                        //胡


}
