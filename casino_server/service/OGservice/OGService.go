package OGservice

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/service/room"
	"casino_server/common/log"
	"casino_server/gamedata"
)

//联众德州,桌子状态
var (
	GAME_STATUS_READY int32 = 0        //准备
	GAME_STATUS_DEAL_CARDS int32 = 1        //发牌
	GAME_STATUS_PRECHIP int32 = 2        //盲注
	GAME_STATUS_FIRST_TURN int32 = 3        //第一轮
	GAME_STATUS_SECOND_TURN int32 = 4        //第二轮
	GAME_STATUS_THIRD_TURN int32 = 5        //第三轮
	GAME_STATUS_LAST_TURN int32 = 6        //第四轮
	GAME_STATUS_SHOW_RESULT int32 = 7        //完成
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



//处理登录游戏的协议
func HandlerGameEnterMatch(m *bbproto.Game_EnterMatch, a gate.Agent) error {
	//定义需要的参数
	var userId uint32 = uint32(m.GetTableid())                             //获取到用户的userId
	if userId == 0 {
		userId = 10006
	}
	result := newGame_SendGameInfo()                //需要返回的信息

	//1,进入房间,放回房间和错误信息
	mydesk, err := room.ThGameRoomIns.AddUser(userId, a)
	if err != nil || mydesk == nil {
		errMsg := err.Error()
		log.E("进入房间失败,errMsg[%v]", errMsg)
		//这里需要给客户端返回失败的信息
		a.WriteMsg(result)
		return err
	}

	//2 构造信息并且返回
	initGameSendgameInfoByDesk(mydesk, result)
	log.T("给请求登陆房间的人回复信息[%v]",result)
	a.WriteMsg(result)

	//3,进入房间的广播,告诉其他人有新的玩家进来了
	// todo 新增加og add user 广播....
	mydesk.THBroadcastAddUser(userId, userId)

	//4,最后:确定是否开始游戏, 上了牌桌之后,如果玩家人数大于1,并且游戏处于stop的状态,则直接开始游戏
	//这是游戏刚开始,的处理方式
	if mydesk.UserCount >= room.TH_DESK_LEAST_START_USER  && mydesk.Status == room.TH_DESK_STATUS_STOP {
		err = run(mydesk)
		if err != nil {
			log.E("开始德州扑克游戏的时候失败")
			return nil
		}
	}
	return nil
}


//初始化一个Game_SendGameInfo
func newGame_SendGameInfo() *bbproto.Game_SendGameInfo {
	result := &bbproto.Game_SendGameInfo{}
	result.TablePlayer = new(int32)
	result.Tableid = new(int32)
	result.BankSeat = new(int32)
	result.ChipSeat = new(int32)
	result.ActionTime = new(int32)
	result.DelayTime = new(int32)
	result.GameStatus = new(int32)
	//result.NRebuyCount
	//result.NAddonCount
	result.Pool = new(int64)
	result.MinRaise = new(int64)
	result.NInitActionTime = new(int32)
	result.NInitDelayTime = new(int32)
	result.Handcard = make([]*bbproto.Game_CardInfo, 0)
	result.TurnCoin = make([]int64, 0)
	result.BFold = make([]int32, 0)
	result.BAllIn = make([]int32, 0)
	result.BBreak = make([]int32, 0)
	result.BLeave = make([]int32, 0)

	return result
}

//返回房间的信息 todo 登陆成功的处理,给请求登陆的玩家返回登陆结果的消息
func initGameSendgameInfoByDesk(mydesk *room.ThDesk, result *bbproto.Game_SendGameInfo) error {
	//初始化桌子相关的信息
	//*result.Matchid		= 0
	*result.Tableid = int32(mydesk.Id)        //桌子的Id
	*result.TablePlayer = mydesk.UserCount
	*result.BankSeat = int32(mydesk.Dealer)        //庄家
	*result.ChipSeat = int32(mydesk.GetUserIndex(mydesk.BetUserNow))//当前活动玩家
	*result.ActionTime = int32(room.TH_TIMEOUT_DURATION)        //当前操作时间,服务器当前的时间
	*result.DelayTime = int32(0)        //当前延时时间
	*result.GameStatus = deskStatus2OG(mydesk)
	//result.NRebuyCount
	//result.NAddonCount
	*result.Pool = int64(mydesk.Jackpot)                //奖池
	result.Publiccard = thPublicCard2OGC(mydesk)        //公共牌...
	*result.MinRaise = int64(mydesk.MinRaise)        //最低加注金额
	*result.NInitActionTime = int32(room.TH_TIMEOUT_DURATION)
	*result.NInitDelayTime = int32(room.TH_TIMEOUT_DURATION)
	result.Handcard = getHandCard(mydesk)		//用户手牌

	//循环User来处理
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户当前状态
			result.BAllIn = append(result.BAllIn, isAllIn(u))                        //是否已经全下了
			result.BBreak = append(result.BBreak, isBreak(u))                        //是否
			result.BLeave = append(result.BLeave, isLeave(u))
			result.BEnable = append(result.BEnable, isEnable(u))                        //用户是否可以操作,0表示不能操作,1表示可以操作
			result.BFold = append(result.BFold, isFold(u))        //是否弃牌

			//用户的余额
			result.TurnCoin = append(result.HandCoin, int64(u.HandCoin))
			result.HandCoin = append(result.HandCoin, int64(u.RoundBet))

			//nickName			//seatId
			result.NickName = append(result.NickName, u.NickName)
			result.SeatId = append(result.SeatId, int32(i))

		} else {

		}

	}

	return nil

}

//判断是否allIn
func isAllIn(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_ALLINING {
		return 1
	} else {
		return 0
	}
}


//判断是否allIn
func isBreak(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_BREAK {
		return 1
	} else {
		return 0
	}
}

//判断是否allIn
func isLeave(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_LEAVE {
		return 1
	} else {
		return 0
	}
}

//判断是否allIn
func isFinalAddon(u *room.ThUser) int32 {
	return 0
}

//判断是否allIn
func isEnable(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_BETING {
		//游戏中
		return 1
	} else {
		return 0
	}
}

//判断是否allIn
func isFold(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_FOLDED {
		//游戏中
		return 1
	} else {
		return 0
	}
}



//解析手牌
func getHandCard(mydesk *room.ThDesk) []*bbproto.Game_CardInfo {
	var handCard []*bbproto.Game_CardInfo
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户手牌
			result := make([]*bbproto.Game_CardInfo, 0)
			for i := 0; i < len(u.Cards); i++ {
				c := u.Cards[i]
				gc := thCard2OGCard(c)
				result = append(result, gc)
			}
			handCard = append(handCard, )
		} else {

		}

	}
	return handCard
}

//th的牌转换成OG的牌
func thCard2OGCard(pai *bbproto.Pai) *bbproto.Game_CardInfo {
	result := &bbproto.Game_CardInfo{}
	result.Value = new(int32)
	result.Color = new(int32)

	//初始化花色
	switch *pai.Flower {
	case "HEART" :
		*result.Color = POKER_COLOR_HEARTS
	case "DIAMOND" :
		*result.Color = POKER_COLOR_DIAMOND
	case "CLUB" :
		*result.Color = POKER_COLOR_CLUB
	case "SPADE" :
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

//手牌转换为OG可以使用的牌
func thPublicCard2OGC(desk *room.ThDesk) []*bbproto.Game_CardInfo {
	result := make([]*bbproto.Game_CardInfo, len(desk.PublicPai))
	for i := 0; i < len(desk.PublicPai); i++ {
		result[i] = thCard2OGCard(desk.PublicPai[i])
	}
	return result
}


//游戏状态的转换
func deskStatus2OG(desk *room.ThDesk) int32 {
	var result int32 = GAME_STATUS_READY
	status := desk.Status
	round := desk.RoundCount

	if status == room.TH_DESK_STATUS_STOP {
		//没有开始
		result = GAME_STATUS_READY
	} else if status == room.TH_DESK_STATUS_SART {
		switch round {
		case room.TH_DESK_ROUND1:
			result = GAME_STATUS_FIRST_TURN
		case room.TH_DESK_ROUND2:
			result = GAME_STATUS_SECOND_TURN
		case room.TH_DESK_ROUND3:
			result = GAME_STATUS_THIRD_TURN
		case room.TH_DESK_ROUND4:
			result = GAME_STATUS_LAST_TURN
		default:
			result = GAME_STATUS_DEAL_CARDS
		}
	} else if status == room.TH_DESK_STATUS_LOTTERY {
		result = GAME_STATUS_SHOW_RESULT
	}

	return result
}


//开始游戏
func  run(mydesk *room.ThDesk)error{
	mydesk.Run()

	//2,发送盲注的广播
	blindB := &bbproto.Game_BlindCoin{}
	//blindB.Tableid	//deskid
	//blindB.Matchid  //roomId
	//blindB.Banker	//庄
	//blindB.Bigblind	//大盲注
	//blindB.Bigblindseat	//大盲注座位号
	//
	//blindB.Smallblind
	//blindB.Smallblindseat
	//
	//blindB.Coin	//
	//blindB.Pool

	mydesk.THBroadcastProto(blindB,0)

	//3,发送手牌的广播
	initCardB := &bbproto.Game_InitCard{}

	initCardB.Tableid = new(int32)
	*initCardB.Tableid = int32(mydesk.Id)

	initCardB.ActionTime = new(int32)
	initCardB.DelayTime  = new(int32)
	initCardB.HandCard = getHandCard(mydesk)
	initCardB.PublicCard = thPublicCard2OGC(mydesk)
	initCardB.MinRaise = &mydesk.MinRaise
	initCardB.NextUser = mydesk.GetResUserModelById(mydesk.BetUserNow).SeatNumber
	initCardB.Seat = &mydesk.UserCount
	mydesk.THBroadcastProto(initCardB,0)

	return nil

}

//处理押注的请求
func HandlerFollowBet(m  *bbproto.Game_FollowBet,a gate.Agent) error{
	log.T("处理用户押注的请求")
	seatId := m.GetSeat()
	desk := room.ThGameRoomIns.GetDeskById(m.GetTableid())
	desk.OgFollowBet(seatId,a)
	return nil
}


//处理加注
func HandlerRaiseBet(m *bbproto.Game_RaiseBet,a gate.Agent) error{
	return nil
}

//处理让牌
func HandlerCheckBet(m *bbproto.Game_CheckBet,a gate.Agent) error{
	return nil
}


//处理让牌
func HandlerFoldBet(m *bbproto.Game_FoldBet,a gate.Agent) error{
	return nil
}


//通过agent返回UserId
func getUserIdByAgent( a gate.Agent) uint32{
	//获取agent中的userData
	ad := a.UserData()
	if ad == nil {
		log.E("agent中的userData为nil")
		return uint32(0)

	}

	userData := ad.(*gamedata.AgentUserData)
	log.T("得到的UserAgent中的userId是[%v]",userData.UserId)
	//return userData.UserId

	//测试代码,返回10006
	return 10006
}

