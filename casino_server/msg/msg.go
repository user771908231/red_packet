package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	"casino_server/msg/bbprotogo"
)

// protobuf 消息处理器
var PortoProcessor = protobuf.NewProcessor()

func init() {
	// 这里我们注册了一个 JSON 消息 Hello

	//次处注册proto 的消息
	PortoProcessor.Register(&bbproto.TestP1{})		//0	测试用
	PortoProcessor.Register(&bbproto.Reg{})			//1	注册协议(已经废弃)
	PortoProcessor.Register(&bbproto.ReqAuthUser{})		//2	登陆、注册的协议
	PortoProcessor.Register(&bbproto.HeatBeat{})		//3	心跳协议,检测网络是否联通
	PortoProcessor.Register(&bbproto.GetIntoRoom{})		//4	进入房间时候的请求
	PortoProcessor.Register(&bbproto.RoomMsg{})		//5	给指定房间发送信息
	PortoProcessor.Register(&bbproto.GetRewards{})		//6	各种奖励
	PortoProcessor.Register(&bbproto.Shuiguoji{})		//7	水果机
	PortoProcessor.Register(&bbproto.ShuiguojiHilomp{})	//8	水果机比大小
	PortoProcessor.Register(&bbproto.ShuiguojiRes{})	//9	水果机的回应包

	//扎金花
	PortoProcessor.Register(&bbproto.ZjhRoom{})		//10	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhBet{})		//11	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhMsg{})		//12	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhQueryNoSeatUser{})	//13	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhReqSeat{})		//14	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhLottery{})		//15	进入扎金花的房间
	PortoProcessor.Register(&bbproto.ZjhBroadcastBeginBet{})	//16 广播可以押注了


	//用户奖励相关的
	PortoProcessor.Register(&bbproto.LoginSignInBonus{})	//17	登录签到奖励
	PortoProcessor.Register(&bbproto.LoginTurntableBonus{})	//18	登录转盘奖励
	PortoProcessor.Register(&bbproto.OlineBonus{})		//19	在线奖励
	PortoProcessor.Register(&bbproto.TimingBonus{})		//20	定时奖励


	//德州扑克
	PortoProcessor.Register(&bbproto.ThRoom{})		//21	德州扑克
	PortoProcessor.Register(&bbproto.THBet{})		//22	德州扑克押注
	PortoProcessor.Register(&bbproto.THBetBegin{})		//23	开始德州扑克的广播
	PortoProcessor.Register(&bbproto.THBetBroadcast{})	//24	德州扑克,押注之后的广播


}
