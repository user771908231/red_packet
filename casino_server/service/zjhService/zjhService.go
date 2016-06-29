package zjhService

import (
	"github.com/name5566/leaf/gate"
	"casino_server/msg/bbprotogo"
	"casino_server/conf/intCons"
	"casino_server/service"
)



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


	//模拟返回数据
	//1,庄家的信息2,自己的信息
	var bbalance int32 = 1000
	var bname string = "我庄稼"

	var mbalance int32 = 9999
	var uname string = "用户哈哈哈哈"

	userBanker := &bbproto.User{}
	userMe	:= &bbproto.User{}
	result := &bbproto.ZjhRoom{}

	result.Banker = userBanker
	result.Me = userMe

	userBanker.Name = &bname
	userBanker.Balance = &bbalance

	userMe.Name = &uname
	userMe.Balance = &mbalance

	//返回信息
	a.WriteMsg(result)
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

	//为了测试方便 随意返回数据

	var ba int32 = 98989

	result := &bbproto.ZjhBet{}
	header := &bbproto.ProtoHeader{}
	header.UserId = m.GetHeader().UserId
	result.BetAmount = &ba
	a.WriteMsg(result)

	//广播发送押注信息
	service.ZJHroom.BroadcastProto(result,0)

	return  nil,nil
}

/**
请求进入房间
 */

func getIntoRoom(m *bbproto.ZjhRoom,a gate.Agent)(*bbproto.ZjhRoom,error){
	service.ZJHroom.AddAgent(m.GetHeader().GetUserId(),a)

	//这里给客户端返回信息,包括:押注中(剩余time）、开奖中（剩余time）、jackpot奖池金额、balance、庄家信息、在座玩家


	return nil,nil
}


/**
请求退出房间
 */

func outRoom(m *bbproto.ZjhRoom,a gate.Agent)(*bbproto.ZjhRoom,error){
	service.ZJHroom.RemoveAgent(m.GetHeader().GetUserId())
	return nil,nil
}