package zjhService

import (
	"github.com/name5566/leaf/gate"
	"casino_server/msg/bbprotogo"
	"casino_server/conf/intCons"
	"casino_server/service/room"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/service/userService"
	"casino_server/common/log"
	"errors"
	"casino_server/mode"
)

func init(){

}


/**
	处理关于游戏房间请求的action
 */
func HandlerZjhRoom(m *bbproto.ZjhRoom,a gate.Agent)(*bbproto.ZjhRoom,error){
	reqType := m.GetReqType()
	log.T("进入房间的type:",reqType)
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
	if !room.ZJHroom.Betable() {
		log.E("现在不能下注了")
		return nil,errors.New("现在不能下注了")
	}

	//2,开始押注,判断用户资金是否足够,等
	//todo 1,判断用户资金是否足够,2,判断各种权限
	betrecord := room.GetTBetRecordByUserIdAnd(room.ZJHroom.ZjhRoundNumber,m.GetHeader().GetUserId())
	if betrecord == nil {
		log.T("玩家[%v]第一次押注[%v]",m.GetHeader().GetUserId(),m.GetBetzone())
		betrecord = &bbproto.TBetRecord{}
	}
	betrecord.Betzone = m.GetBetzone()
	betrecord.ZjhRoundNumber = &room.ZJHroom.ZjhRoundNumber
	betrecord.UserId = m.GetHeader().UserId
	room.SaveBetRecord(betrecord)		//保存数据到redis中去


	//3,修改放房间的押注金额
	room.ZJHroom.AddZoneAmount(m.Betzone)

	//为了测试方便 随意返回数据
	result := &bbproto.ZjhBet{}
	result.Header = protoUtils.GetSuccHeader()
	a.WriteMsg(result)

	//广播发送押注信息
	//room.ZJHroom.BroadcastProto(result,0)

	return  nil,nil
}

/**
请求进入房间
 */

func getIntoRoom(m *bbproto.ZjhRoom,a gate.Agent)(*bbproto.ZjhRoom,error){
	//设置用户锁
	l := &mode.LockUser{
		UserId:m.GetHeader().GetUserId(),
	}
	a.SetUserData(l)

	room.ZJHroom.AddAgent(m.GetHeader().GetUserId(),a)
	//这里给客户端返回信息,包括:押注中(剩余time）、开奖中（剩余time）、jackpot奖池金额、balance、庄家信息、在座玩家
	var retErr error = nil							//需要返回的错误信息
	result := &bbproto.ZjhRoom{}						//需要返回的数据
	result.Header		=	protoUtils.GetSuccHeader()		//header
	result.Jackpot		=	&(room.ZJHroom.Jackpot)			//奖池的大小
	result.BetTime		= 	room.ZJHroom.GetBetRemainTime()		//剩余的押注时间
	result.LotteryTime	=	room.ZJHroom.GetLotteryRemainTime()	//剩余的开奖时间
	result.RoomStatus	=	&room.ZJHroom.Status			//房间的当前状态
	result.Zjhpai		=	room.ZJHroom.Zjhpai			//当前的牌信息

	//个人,庄家的信息信息
	result.Banker =  userService.GetUserById(room.ZJHroom.BankerUserId)
	if result.Banker == nil {
		//没有查询到banker的信息,进入房间失败,返回到登录界面
		result.Header.Code = &intCons.CODE_FAIL
		retErr = errors.New("没有找到庄家信息")
	}
	result.Me = userService.GetUserById(m.GetHeader().GetUserId())
	if result.Me == nil {
		result.Header.Code = &intCons.CODE_FAIL
		retErr = errors.New("没有找到用户信息")
	}

	//给客户端返回信息
	log.T("进入扎进话房间之后返回的数据:",result)

	return result,retErr
}


/**
	请求退出房间,主要逻辑有:
	1,保存数据:redis中的数据需要同步到mongodb中
	2,删除游戏房间中管理的连接

 */

func outRoom(m *bbproto.ZjhRoom,a gate.Agent)(*bbproto.ZjhRoom,error){
	//todo 退出房间的时候,需要先保存数据
	//todo 1,用户余额的信息,2,用户游戏的信息,这里表示的扎金花游戏的数据,
	room.ZJHroom.RemoveAgent(m.GetHeader().GetUserId())
	return nil,nil
}

