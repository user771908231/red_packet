package room

import (
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"errors"
	"casino_server/service/pokerService"
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


//押注的通用接口
func (t *ThDesk) OGBet(seatId int32,betType int32,coin int64) error{
	t.Lock()
	defer t.Unlock()

	user := t.getUserBySeat(seatId)
	//1,得到跟注的用户
	if !t.CheckBetUserBySeat(user) {
		log.E("押注人的状态不正确")
		return errors.New("押注人的状态不正确")
	}


	switch betType {
	case TH_DESK_BET_TYPE_BET:
	//押注:大盲注之后的第一个人,这个类型不会使用到

	case TH_DESK_BET_TYPE_CALL:
		t.OgFollowBet(user)

	case TH_DESK_BET_TYPE_FOLD:
		t.OgFoldBet(user)

	case TH_DESK_BET_TYPE_CHECK:        //让牌
		t.OGCheckBet(user)

	case TH_DESK_BET_TYPE_RAISE:        //加注
		t.OGRaiseBet(user,coin)

	case TH_DESK_BET_TYPE_RERRAISE:        //再加注

	case TH_DESK_BET_TYPE_ALLIN:        //全部

	}

	t.nextRoundInfo()	//广播新一局的信息
	t.Lottery()		//开奖
	return nil
}


//这里只处理逻辑
func (t *ThDesk) OgFollowBet(user *ThUser) error {

	//1,押注
	log.T("用户[%v]开始押注",user.UserId)
	err := t.BetUserCall(user.UserId)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//2,初始化下一个押注的人
	t.NextBetUser()

	//3,返回信息
	log.T("返回用户[%v]押注的结果:",user.UserId)
	result := &bbproto.Game_AckFollowBet{}
	result.NextSeat = new(int32)
	result.Coin = new(int64)
	result.Seat = new(int32)
	result.Tableid =new(int32)
	result.CanRaise = new(int32)
	result.MinRaise = new(int64)
	result.Pool = new(int64)
	result.HandCoin = new(int64)

	*result.Coin		= user.Coin
	*result.Seat		= user.Seat                		//座位id
	*result.Tableid		= t.Id
	//*result.CanRaise	= t.CanRaise		     		//是否能加注
	*result.CanRaise	= int32(1)		     		//是否能加注
	*result.MinRaise	= t.MinRaise
	*result.Pool		= t.Jackpot
	*result.NextSeat	= int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	*result.HandCoin 	= user.HandCoin

	t.THBroadcastProtoAll(result)
	return nil
}

//这里只处理逻辑
func (t *ThDesk) OgFoldBet(user  *ThUser) error {

	log.T("用户[%v]开始弃牌",user.UserId)

	//1,弃牌
	//如果弃牌的人是 t.NewRoundBetUser ,需要重新设置值
	if t.NewRoundBetUser == user.UserId {
		t.NextNewRoundBetUser()
	}
	user.Status = TH_USER_STATUS_FOLDED


	//2,初始化下一个押注的人
	t.NextBetUser()


	//3,返回信息
	log.T("返回用户[%v]弃牌的结果:",user.UserId)
	result := &bbproto.Game_AckFoldBet{}
	result.NextSeat = new(int32)
	result.Pool = new(int64)
	result.MinRaise = new(int64)
	result.HandCoin = new(int64)

	*result.Pool = t.Jackpot
	*result.HandCoin = user.HandCoin
	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &user.Seat                			//座位id
	result.Tableid = &t.Id
	*result.MinRaise = t.MinRaise
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人

	t.THBroadcastProto(result,0)
	return nil
}



//联众德州 加注
func (t *ThDesk) OGRaiseBet(user *ThUser,coin int64) error{
	//1,得到跟注的用户,检测用户

	//2,开始处理加注
	err := t.BetUserRaise(user.UserId,coin)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	t.NextBetUser()

	log.T("准备给其他人发送加注的广播")
	//押注成功返回要住成功的消息
	result := &bbproto.Game_AckRaiseBet{}
	result.NextSeat = new(int32)
	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &user.Seat                			//座位id
	result.Tableid = &t.Id
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	result.HandCoin = &user.HandCoin			//表示需要加注多少

	//给所有人广播信息
	t.THBroadcastProto(result,0)
	log.T("开始处理seat[%v]加注的逻辑,t,OgFollowBet()...end",user.Seat)
	return nil
}


//联众德州 让牌
func (t *ThDesk) OGCheckBet(user *ThUser) error{

	//1,让牌
	err := t.BetUserCheck(user.UserId)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}


	//2,计算洗一个押注的人
	t.NextBetUser()


	//3押注成功返回要住成功的消息
	log.T("打印user[%v]让牌的结果:",user.UserId)
	result := &bbproto.Game_AckCheckBet{}
	result.NextSeat = new(int32)
	result.MinRaise = new(int64)
	result.CanRaise = new(int32)

	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &user.Seat                			//座位id
	result.Tableid = &t.Id
	*result.CanRaise = t.CanRaise		     		//是否能加注
	*result.MinRaise = t.MinRaise				//最低加注金额
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	t.THBroadcastProtoAll(result)

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
	var result []int64
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			result = append(result,int64(u.Coin))
		}
	}
	return result
}


//解析每个人下注的金额
func (mydesk *ThDesk) GetHandCoin() []int64{
	var result []int64
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户手牌
			result = append(result,int64(u.HandCoin))
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

	//边池
	ret = append(ret,t.bianJackpot)

	log.T("下一句开始,返回的secondPool【%v】",ret)
	return ret
}

//发送新增用户的广播
func (t *ThDesk) OGTHBroadAddUser(msg *bbproto.Game_SendGameInfo){
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil {

			//给用户发送广播的时候需要判断自己的座位号是多少
			a := t.Users[i].agent
			*msg.Seat = t.Users[i].Seat
			a.WriteMsg(msg)
		}
	}
}



//解析手牌
func (mydesk *ThDesk ) GetHandCard() []*bbproto.Game_CardInfo {
	//log.T("把desk的手牌,转化为og的手牌")
	var handCard []*bbproto.Game_CardInfo
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			log.T("开始给玩家[%v]解析手牌",u.UserId)
			result := make([]*bbproto.Game_CardInfo, 0)
			//用户手牌
			if len(u.Cards) == 2 {
				for i := 0; i < len(u.Cards); i++ {
					c := u.Cards[i]
					gc := ThCard2OGCard(c)
					//增加到数组中
					result = append(result, gc)
				}
			}else{
				log.T("玩家[%v]刚刚进房间,等待别人完成游戏,所以手牌为空",u.UserId)
				gc := &bbproto.Game_CardInfo{}
				gc.Color = new(int32)
				gc.Value = new(int32)
				*gc.Color = POKER_COLOR_COUNT
				*gc.Value = POKER_VALUE_EMPTY
				result = append(result, gc)
				result = append(result, gc)
			}
			handCard = append(handCard, result...)
		} else {

		}

	}
	return handCard
}


func NewGame_WinCoin() *bbproto.Game_WinCoin{
	gwc := &bbproto.Game_WinCoin{}
	gwc.Card1 = new(int32)
	gwc.Card2 = new(int32)
	gwc.Card3 = new(int32)
	gwc.Card4 = new(int32)
	gwc.Card5 = new(int32)
	gwc.Cardtype = new(int32)
	gwc.Coin = new(int64)
	gwc.Seat = new(int32)
	gwc.PoolIndex = new(int32)
	//gwc.Rolename = new(int32)

	return gwc

}

//
func (t *ThDesk) getWinCoinInfo() []*bbproto.Game_WinCoin{
	var ret []*bbproto.Game_WinCoin

	//对每个人做计算
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.winAmount > 0{
			//开是对这个人计算
			gwc := NewGame_WinCoin()

			//这里需要先对牌进行排序
			paicards := OGTHCardPaixu(u.thCards)

			*gwc.Card1 = paicards[0].GetMapKey()
			*gwc.Card2 = paicards[1].GetMapKey()
			*gwc.Card3 = paicards[2].GetMapKey()
			*gwc.Card4 = paicards[3].GetMapKey()
			*gwc.Card5 = paicards[4].GetMapKey()
			*gwc.Cardtype = u.thCards.GetOGCardType()
			*gwc.Coin = u.winAmount
			*gwc.Seat = u.Seat
			*gwc.PoolIndex = int32(0)
			ret = append(ret,gwc)
		}
	}
	return ret
}

func (t *ThDesk) GetBshowCard() []int32{
	var ret []int32
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			ret = append(ret,int32(1))
		}
	}

	return ret
}

//把结果牌 针对og来排序
func OGTHCardPaixu(s *pokerService.ThCards) []*bbproto.Pai{
	lensCards := len(s.Cards)
	ret := make([]*bbproto.Pai,lensCards)

	for i := 0; i <lensCards;i++  {
		log.T("og 排序之前[%v]----[%v]:",i,s.Cards[i].GetValue())
	}

	if s.ThType == pokerService.THPOKER_TYPE_SITIAO {
		//如果是四条的话,前四张是四条
		if s.KeyValue[0] == s.Cards[0].GetValue() {
			copy(ret,s.Cards)
		}else{
			copy(ret[0:4], s.Cards[1:5])
			copy(ret[4:5], s.Cards[0:1])
		}
	}else if s.ThType == pokerService.THPOKER_TYPE_HULU{
		//如果是葫芦的话,前三张是三条
		if s.KeyValue[0] == s.Cards[0].GetValue() {
			copy(ret,s.Cards)
		}else{
			copy(ret[0:3], s.Cards[2:5])
			copy(ret[3:5], s.Cards[0:2])
		}
	}else if s.ThType == pokerService.THPOKER_TYPE_SANTIAO{
		//如果是三条的话,前三张是三条
		santiaoIndex := 0
		for i := 0; i < lensCards; i++ {
			if s.Cards[i].GetValue() == s.KeyValue[0] {
				santiaoIndex = i
				break
			}
		}
		ret[0] = s.Cards[santiaoIndex]
		ret[1] = s.Cards[santiaoIndex +1]
		ret[2] = s.Cards[santiaoIndex +2]

		log.T("排序葫芦1:[%v]",ret)
		tempS := make([]*bbproto.Pai,lensCards)
		copy(tempS,s.Cards)
		for i := 3; i < 5; i++ {
			for j := 0; j < lensCards; j++ {
				ts := tempS[j]
				if ts!=nil && ts.GetValue() != ret[0].GetValue() {
					ret[i] = ts
					log.T("排序葫芦[%v]:[%v]",i,ret[i])
					tempS[j] = nil
					break
				}
			}
		}
	}else if s.ThType == pokerService.THPOKER_TYPE_LIANGDUI{
		//如果是两队的话,

		yiduiIndex,liangduiIndex,danpai := 0,0,0

		for i := 0; i < lensCards; i++ {
			if s.Cards[i].GetValue() == s.KeyValue[0] {
				yiduiIndex = i
				break
			}
		}

		for i := 0; i < lensCards; i++ {
			if s.Cards[i].GetValue() == s.KeyValue[1] {
				liangduiIndex = i
				break
			}
		}

		for i := 0; i < lensCards; i++ {
			if s.Cards[i].GetValue() == s.KeyValue[2] {
				danpai = i
				break
			}
		}

		copy(ret[0:2],s.Cards[yiduiIndex:yiduiIndex+2])
		copy(ret[2:4],s.Cards[liangduiIndex:liangduiIndex+2])
		copy(ret[4:lensCards],s.Cards[danpai:danpai+1])

	}else if s.ThType == pokerService.THPOKER_TYPE_YIDUI{
		//如果是一对的话,
		//1,找到一对的坐标
		tempS := make([]*bbproto.Pai,lensCards)
		copy(tempS,s.Cards)
		log.T("temps:[%v]",tempS)

		yiduiIndex := 0
		for i := 0; i < lensCards; i++ {
			if s.Cards[i].GetValue() == s.KeyValue[0] {
				yiduiIndex = i
				break
			}
		}
		ret[0] = tempS[yiduiIndex]
		tempS[yiduiIndex] = nil
		ret[1] = tempS[yiduiIndex+1]
		tempS[yiduiIndex+1] = nil

		for i := 2; i < 5; i++ {
			for j := 0; j < lensCards; j++ {
				ts := tempS[j]
				if ts!=nil {
					ret[i] = ts
					tempS[j] = nil
					log.T("排序一对[%v],[%v]",i,ret[i])
					break
				}
			}
		}

	}else{
		copy(ret,s.Cards)
	}

	for i := 0; i <lensCards;i++  {
		log.T("og 排序之后[%v]----[%v]:",i,ret[i].GetValue())
	}
	return ret
}