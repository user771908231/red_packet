package zjhService

import (
	"github.com/name5566/leaf/gate"
	"casino_server/msg/bbprotogo"
	"casino_server/conf/intCons"
	"casino_server/service/room"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/service/userService"
)


func init(){

}


/**
处理关于游戏房间请求的action
 */
func HandlerZjhRoom(m *bbproto.ZjhRoom,a gate.Agent)(*bbproto.ZjhRoom,error){
	reqType := m.GetReqType()
	if reqType == intCons.REQ_TYPE_IN{
		//进入房间的请求
		getIntoRoom(m,a)
	}else{
		//退出房间的请求
		outRoom(m,a)
	}
	return nil,nil
}

/**
处理扎金花押注的请求,请求流程如下:
	1,判断是否属于押注的状态
	2,判断用户余额是否充足
	3,押注成功需要广播给每个人

 */
func HandlerZjhBet(m *bbproto.ZjhBet,a gate.Agent)(*bbproto.ZjhBet,error){
	//1,判断是否属于押注的状态
	//if !service.ZJHroom.Betable() {
	//	log.E("现在不能下注了")
	//	return nil,errors.New("现在不能下注了")
	//}

	//2,开始押注,判断用户资金是否足够,等

	//3,修改放房间的押注金额
	room.ZJHroom.AddZoneAmount(m.Betzone)

	//为了测试方便 随意返回数据
	result := &bbproto.ZjhBet{}
	header := &bbproto.ProtoHeader{}
	header.UserId = m.GetHeader().UserId
	header.Code = &intCons.CODE_SUCC		//表示请求成功
	a.WriteMsg(result)

	//广播发送押注信息
	room.ZJHroom.BroadcastProto(result,0)

	return  nil,nil
}

/**
请求进入房间
 */

func getIntoRoom(m *bbproto.ZjhRoom,a gate.Agent)(*bbproto.ZjhRoom,error){
	room.ZJHroom.AddAgent(m.GetHeader().GetUserId(),a)
	//这里给客户端返回信息,包括:押注中(剩余time）、开奖中（剩余time）、jackpot奖池金额、balance、庄家信息、在座玩家
	result := &bbproto.ZjhRoom{}
	result.Header		=	porotoUtils.GetSuccHeader()		//header
	result.Jackpot		=	&(room.ZJHroom.Jackpot)
	result.BetTime		= 	room.ZJHroom.GetBetRemainTime()		//剩余的押注时间
	result.LotteryTime	=	room.ZJHroom.GetLotteryRemainTime()	//剩余的开奖时间
	result.RoomStatus	=	&room.ZJHroom.Status			//房间的当前状态

	//个人,庄家的信息信息
	result.Banker =  userService.GetUserById(room.ZJHroom.BankerUserId)
	result.Me = userService.GetUserById(m.GetHeader().GetUserId())
	a.WriteMsg(result)
	return result,nil
}


/**
请求退出房间
 */

func outRoom(m *bbproto.ZjhRoom,a gate.Agent)(*bbproto.ZjhRoom,error){
	room.ZJHroom.RemoveAgent(m.GetHeader().GetUserId())
	return nil,nil
}