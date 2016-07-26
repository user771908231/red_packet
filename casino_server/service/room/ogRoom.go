package room

import (
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
)


//通过座位号来找到user
func (t *ThDesk) getUserBySeat(seatId int32) *ThUser {
	return t.Users[seatId]
}

//这里只处理逻辑
func (t *ThDesk) OgFollowBet(seatId int32) error {
	log.T("开始处理og 跟注的逻辑,t,OgFollowBet()...")
	t.Lock()
	defer t.Unlock()

	user := t.getUserBySeat(seatId)
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

	return nil
}

//这里只处理逻辑
func (t *ThDesk) OgFoldBet(seatId int32) error {
	log.T("开始处理og 弃牌的逻辑,t,OgFollowBet()...")
	t.Lock()
	defer t.Unlock()

	user := t.getUserBySeat(seatId)
	err := t.BetUserFold(user.UserId)
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
	result := &bbproto.Game_AckFoldBet{}
	result.NextSeat = new(int32)
	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &seatId                			//座位id
	result.Tableid = &t.Id
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	//给所有人广播信息
	t.THBroadcastProto(result,0)

	return nil
}

