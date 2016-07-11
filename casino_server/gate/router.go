package gate

import (
)
import (
	"casino_server/msg"
	"casino_server/msg/bbprotogo"
	"casino_server/game"
	"casino_server/login"
	"casino_server/system"
	"casino_server/bonus"
)

func init() {
	//指定protobuf格式的路由
	msg.PortoProcessor.SetRouter(&bbproto.TestP1{},game.ChanRPC)
	msg.PortoProcessor.SetRouter(&bbproto.Reg{},login.ChanRPC)
	msg.PortoProcessor.SetRouter(&bbproto.ReqAuthUser{},login.ChanRPC)
	msg.PortoProcessor.SetRouter(&bbproto.HeatBeat{},system.ChanRPC)

	//水果机
	msg.PortoProcessor.SetRouter(&bbproto.GetIntoRoom{},game.ChanRPC)
	msg.PortoProcessor.SetRouter(&bbproto.RoomMsg{},game.ChanRPC)		//给指定房间发送信息
	msg.PortoProcessor.SetRouter(&bbproto.GetRewards{},game.ChanRPC)	//获取奖励
	msg.PortoProcessor.SetRouter(&bbproto.Shuiguoji{},game.ChanRPC)		//水果机
	msg.PortoProcessor.SetRouter(&bbproto.ShuiguojiHilomp{},game.ChanRPC)	//水果机比大小

	//扎金花
	msg.PortoProcessor.SetRouter(&bbproto.ZjhRoom{},game.ChanRPC)		//进入扎金花的房间
	msg.PortoProcessor.SetRouter(&bbproto.ZjhLottery{},game.ChanRPC)		//进入扎金花的房间
	msg.PortoProcessor.SetRouter(&bbproto.ZjhMsg{},game.ChanRPC)		//进入扎金花的房间
	msg.PortoProcessor.SetRouter(&bbproto.ZjhBet{},game.ChanRPC)		//进入扎金花的房间
	msg.PortoProcessor.SetRouter(&bbproto.ZjhReqSeat{},game.ChanRPC)		//进入扎金花的房间
	msg.PortoProcessor.SetRouter(&bbproto.ZjhQueryNoSeatUser{},game.ChanRPC)		//进入扎金花的房间

	//奖励
	msg.PortoProcessor.SetRouter(&bbproto.LoginTurntableBonus{},bonus.ChanRPC)	//登陆奖励
	msg.PortoProcessor.SetRouter(&bbproto.LoginSignInBonus{},bonus.ChanRPC)		//连续签到的奖励

}