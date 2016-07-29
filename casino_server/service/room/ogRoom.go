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

	t.nextRoundInfo()
	t.Lottery()//

	return nil
}


//这里只处理逻辑
func (t *ThDesk) OgFollowBet(user *ThUser) error {
	log.T("开始处理seat[%v]跟注的逻辑,t,OgFollowBet()...",user.Seat)

	//2,开始押注
	err := t.BetUserCall(user.UserId)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//3,判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人

	t.NextBetUser()

	//4,押注成功返回要住成功的消息
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

	//给所有人广播信息
	t.THBroadcastProto(result,0)
	log.T("开始处理seat{%v}跟注的逻辑,t,OgFollowBet()...end",user.Seat)
	return nil
}

//这里只处理逻辑
func (t *ThDesk) OgFoldBet(user  *ThUser) error {

	//2,开始弃牌
	err := t.BetUserFold(user.UserId)
	if err != nil {
		log.E("弃牌的时候出错了.errMsg[%v],", err.Error())
		return err
	}

	//3,判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		t.NextBetUser()
		log.T("准备给其他人发送弃牌的广播")
	}

	//押注成功返回要住成功的消息
	result := &bbproto.Game_AckFoldBet{}
	result.NextSeat = new(int32)

	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &user.Seat                			//座位id
	result.Tableid = &t.Id
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	//给所有人广播信息
	t.THBroadcastProto(result,0)
	log.T("开始处理seat[%v]弃牌的逻辑,t,OgFollowBet()...end",user.Seat)
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
	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		t.NextBetUser()
	}

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

	//2,开始弃牌
	err := t.BetUserCheck(user.UserId)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//3,判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		t.NextBetUser()
		log.T("准备给其他人发送弃牌的广播")
	}

	//4押注成功返回要住成功的消息
	result := &bbproto.Game_AckCheckBet{}
	result.NextSeat = new(int32)
	result.Coin = &t.BetAmountNow        			//本轮压了多少钱
	result.Seat = &user.Seat                			//座位id
	result.Tableid = &t.Id
	result.CanRaise	= &t.CanRaise		     		//是否能加注
	*result.NextSeat =int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	t.THBroadcastProto(result,0)

	log.T("开始处理seat[%v]让牌的逻辑,t,OgFollowBet()...end",user.Seat)
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


//发送新增玩家的广播
func (t *ThDesk) OGTHBroadAddUser(newUserId uint32) error{
	newUser := t.GetUserByUserId(newUserId)

	//生成广播的信息
	broadUser := &bbproto.Game_SendAddUser{}
	broadUser.Coin = new(int64)
	broadUser.Seat = new(int32)
	broadUser.Tableid = new(int32)
	broadUser.NickName= new(string)

	*broadUser.Coin = newUser.Coin
	*broadUser.Seat = newUser.Seat
	*broadUser.NickName = newUser.NickName
	*broadUser.Tableid = t.Id

	log.T("开始广播新增用户的proto消息[%v]",broadUser)
	t.THBroadcastProto(broadUser,newUserId)
	return nil
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
		if u != nil {
			//开是对这个人计算
			gwc := NewGame_WinCoin()
			*gwc.Card1 = GetOGCardIndex(u.thCards.Cards[0])
			*gwc.Card2 = GetOGCardIndex(u.thCards.Cards[1])
			*gwc.Card3 = GetOGCardIndex(u.thCards.Cards[2])
			*gwc.Card4 = GetOGCardIndex(u.thCards.Cards[3])
			*gwc.Card5 = GetOGCardIndex(u.thCards.Cards[4])
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

//通过牌的信息返回index
func GetOGCardIndex(p *bbproto.Pai) int32{
	return p.GetMapKey()
}
