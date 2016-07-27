package room

import (
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"errors"
)

//牌的花色 和值
var (
	POKER_COLOR_DIAMOND int32 = 0    //方片
	POKER_COLOR_CLUB int32 = 1    //梅花
	POKER_COLOR_HEARTS int32 = 2    //红桃
	POKER_COLOR_SPADE int32 = 3    //黑桃
	POKER_COLOR_COUNT int32 = 4
)


var (
	POKER_VALUE_A int32 = 0
	POKER_VALUE_2 int32 = 1
	POKER_VALUE_3 int32 = 2
	POKER_VALUE_4 int32 = 3
	POKER_VALUE_5 int32 = 4
	POKER_VALUE_6 int32 = 5
	POKER_VALUE_7 int32 = 6
	POKER_VALUE_8 int32 = 7
	POKER_VALUE_9 int32 = 8
	POKER_VALUE_10 int32 = 9
	POKER_VALUE_J int32 = 10
	POKER_VALUE_Q int32 = 11
	POKER_VALUE_K int32 = 12
	POKER_VALUE_COUNT int32 = 13
	POKER_VALUE_BACK int32 = 52    // 牌背
	POKER_VALUE_EMPTY int32 = 53    // 没牌
)




//通过座位号来找到user
func (t *ThDesk) getUserBySeat(seatId int32) *ThUser {
	return t.Users[seatId]
}

//这里只处理逻辑
func (t *ThDesk) OgFollowBet(seatId int32) error {
	log.T("开始处理seat[%v]跟注的逻辑,t,OgFollowBet()...",seatId)
	t.Lock()
	defer t.Unlock()

	user := t.getUserBySeat(seatId)
	if !t.CheckBetUser(user.UserId) {
		log.T("押注人是[%v]而不是[%v]",t.BetUserNow,user.UserId)
		return errors.New("服务器错误")
	}

	user.InitWait()
	err := t.BetUserCall(user.UserId, t.BetAmountNow)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		t.NextBetUser()
		log.T("准备给其他人发送押注的广播")
		//广播给下一个人押注
	}

	//押注成功返回要住成功的消息
	//初始化
	result := &bbproto.Game_AckFollowBet{}
	result.NextSeat = new(int32)

	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &seatId                			//座位id
	result.Tableid = &t.Id
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat = int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	//a.WriteMsg(result)

	//给所有人广播信息
	t.THBroadcastProto(result,0)
	log.T("开始处理seat{%v}跟注的逻辑,t,OgFollowBet()...end",seatId)
	return nil
}

//这里只处理逻辑
func (t *ThDesk) OgFoldBet(seatId int32) error {
	log.T("开始处理seat[%v]弃牌的逻辑,t,OgFollowBet()...",seatId)
	t.Lock()
	defer t.Unlock()

	user := t.getUserBySeat(seatId)
	user.InitWait()
	user.waitUUID = ""
	err := t.BetUserFold(user.UserId)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		t.NextBetUser()
		log.T("准备给其他人发送弃牌的广播")
	}

	//押注成功返回要住成功的消息
	//初始化
	result := &bbproto.Game_AckFoldBet{}
	result.NextSeat = new(int32)
	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &seatId                			//座位id
	result.Tableid = &t.Id
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	//给所有人广播信息
	t.THBroadcastProto(result,0)
	log.T("开始处理seat[%v]弃牌的逻辑,t,OgFollowBet()...end",seatId)
	return nil
}



//联众德州 加注
func (t *ThDesk) OGRaiseBet(seatId int32,coin int64) error{
	log.T("开始处理seat[%v]弃牌的逻辑,t,OgFollowBet()...",seatId)
	t.Lock()
	defer t.Unlock()

	user := t.getUserBySeat(seatId)
	user.InitWait()
	user.waitUUID = ""
	err := t.BetUserRaise(user.UserId,coin)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		t.NextBetUser()
		log.T("准备给其他人发送弃牌的广播")
	}

	//押注成功返回要住成功的消息
	//初始化
	result := &bbproto.Game_AckRaiseBet{}
	result.NextSeat = new(int32)
	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &seatId                			//座位id
	result.Tableid = &t.Id
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	//给所有人广播信息
	t.THBroadcastProto(result,0)
	log.T("开始处理seat[%v]弃牌的逻辑,t,OgFollowBet()...end",seatId)
	return nil

}


//联众德州 让牌
func (t *ThDesk) OGCheckBet(seatId int32) error{
	log.T("开始处理seat[%v]弃牌的逻辑,t,OgFollowBet()...",seatId)
	t.Lock()
	defer t.Unlock()

	user := t.getUserBySeat(seatId)
	user.InitWait()
	user.waitUUID = ""
	err := t.BetUserCheck(user.UserId)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		t.NextBetUser()
		log.T("准备给其他人发送弃牌的广播")
	}

	//押注成功返回要住成功的消息
	//初始化
	result := &bbproto.Game_AckCheckBet{}
	result.NextSeat = new(int32)
	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &seatId                			//座位id
	result.Tableid = &t.Id
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	//给所有人广播信息
	t.THBroadcastProto(result,0)
	log.T("开始处理seat[%v]弃牌的逻辑,t,OgFollowBet()...end",seatId)
	return nil
}



//th的牌转换成OG的牌
func ThCard2OGCard(pai *bbproto.Pai) *bbproto.Game_CardInfo {
	result := &bbproto.Game_CardInfo{}
	result.Value = new(int32)
	result.Color = new(int32)

	//log.T("初始化牌的花色:*pai.flower[%v]",*pai.Flower)

	//初始化花色
	switch *pai.Flower {
	case "heart" :
		*result.Color = POKER_COLOR_HEARTS
	case "diamond" :
		*result.Color = POKER_COLOR_DIAMOND
	case "club" :
		*result.Color = POKER_COLOR_CLUB
	case "spade" :
		*result.Color = POKER_COLOR_SPADE
	}

	//初始化值
	switch *pai.Value {
	case 2:
		*result.Value = POKER_VALUE_2
	case 3:
		*result.Value = POKER_VALUE_3
	case 4:
		*result.Value = POKER_VALUE_4
	case 5:
		*result.Value = POKER_VALUE_5
	case 6:
		*result.Value = POKER_VALUE_6
	case 7:
		*result.Value = POKER_VALUE_7
	case 8:
		*result.Value = POKER_VALUE_8
	case 9:
		*result.Value = POKER_VALUE_9
	case 10:
		*result.Value = POKER_VALUE_10
	case 11:                                //j
		*result.Value = POKER_VALUE_J
	case 12:                                //q
		*result.Value = POKER_VALUE_Q
	case 13:                                //k
		*result.Value = POKER_VALUE_K
	case 14:                                //a
		*result.Value = POKER_VALUE_A

	}

	return result
}

func (mydesk *ThDesk) GetCoin() []int64{
	result := make([]int64,len(mydesk.Users))
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户手牌
			result[i] = int64(u.Coin)
		} else {
			result[i]= int64(0)
		}
	}
	return result
}

//
func NewGame_SendOverTurn() *bbproto.Game_SendOverTurn{
	ret := &bbproto.Game_SendOverTurn{}
	ret.NextSeat = new(int32)
	ret.Tableid  = new(int32)
	ret.Pool = new(int64)
	ret.MinRaise = new(int64)
	return ret
}

//取边池的大小
func (t *ThDesk) GetSecondPool() []int64{
	var ret []int64
	for i := 0; i < len(t.AllInJackpot); i++ {
		a := t.AllInJackpot[i]
		if a != nil {
			ret = append(ret,a.Jackpopt)
		}
	}
	return ret
}