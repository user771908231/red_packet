package room

import (
	"casino_server/common/log"
	"github.com/name5566/leaf/gate"
	"casino_server/msg/bbprotogo"
)


//通过座位号来找到user
func (t *ThDesk) getUserBySeat(seatId int32) *ThUser{
	return t.Users[seatId]
}

//这里只处理逻辑
func (t *ThDesk) OgFollowBet(seatId int32,a gate.Agent) error{
	log.T("开始处理og 跟注的逻辑,t,OgFollowBet()...")
	t.Lock()
	defer t.Unlock()

	user := t.getUserBySeat(seatId)
	err := t.BetUserCall(user.userId,t.BetAmountNow)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],",err.Error())
	}else{
		//押注成功返回要住成功的消息

		//初始化
		result := &bbproto.Game_AckFollowBet{}
		result.Coin = new(int64)

		//负值
		*result.Coin = t.BetAmountNow
		a.WriteMsg(result)
	}

	//判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		t.NextBetUser()
		//广播给下一个人押注
	}
	return nil
}