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



//处理登录游戏的协议
/**
	1,判断用户是否已经登陆了游戏
	2,如果已经登陆了游戏,替换现有的agent
	3,如果没有登陆游戏,走正常的流程
 */
func HandlerGameEnterMatch(m *bbproto.Game_EnterMatch, a gate.Agent) error {

	userId := m.GetUserId()
	//定义需要的参数
	result := newGame_SendGameInfo()                //需要返回的信息

	//1,进入房间,返回房间和错误信息
	mydesk, err := room.ThGameRoomIns.AddUser(userId, a)
	if err != nil || mydesk == nil {
		errMsg := err.Error()
		log.E("用户[%v]进入房间失败,errMsg[%v]",userId,errMsg)
		a.WriteMsg(result)
		return err
	}

	//2 构造信息并且返回
	initGameSendgameInfoByDesk(mydesk, result,userId)
	log.T("给请求登陆房间的人[%v]回复信息[%v]",userId,result)
	a.WriteMsg(result)

	//3 发送进入游戏房间的广播
	mydesk.OGTHBroadAddUser(userId)

	//4,最后:确定是否开始游戏, 上了牌桌之后,如果玩家人数大于1,并且游戏处于stop的状态,则直接开始游戏
	if mydesk.IsTime2begin() {
		err = run(mydesk)
		if err != nil {
			log.E("开始德州扑克游戏的时候失败!")
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
	result.Seat = new(int32)

	return result
}

//返回房间的信息 todo 登陆成功的处理,给请求登陆的玩家返回登陆结果的消息
func initGameSendgameInfoByDesk(mydesk *room.ThDesk, result *bbproto.Game_SendGameInfo,myUserId uint32) error {
	//初始化桌子相关的信息
	*result.Tableid = int32(mydesk.Id)        //桌子的Id
	*result.TablePlayer = mydesk.UserCount
	*result.BankSeat = int32(mydesk.Dealer)        //庄家
	*result.ChipSeat = int32(mydesk.GetUserIndex(mydesk.BetUserNow))//当前活动玩家
	*result.ActionTime = int32(room.TH_TIMEOUT_DURATION_INT)        //当前操作时间,服务器当前的时间
	*result.DelayTime = int32(0)        //当前延时时间
	*result.GameStatus = deskStatus2OG(mydesk)
	*result.Pool = int64(mydesk.Jackpot)                //奖池
	result.Publiccard = thPublicCard2OGC(mydesk)        //公共牌...
	*result.MinRaise = int64(mydesk.MinRaise)        //最低加注金额
	*result.NInitActionTime = int32(room.TH_TIMEOUT_DURATION_INT)
	*result.NInitDelayTime = int32(room.TH_TIMEOUT_DURATION)
	result.Handcard = getHandCard(mydesk)		//用户手牌
	result.HandCoin = mydesk.GetCoin()	//下注的金额
	result.TurnCoin = getTurnCoin(mydesk)
	*result.Seat	= int32(mydesk.GetUserIndex(myUserId))	//我

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
	//log.T("把desk的手牌,转化为og的手牌")
	var handCard []*bbproto.Game_CardInfo
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户手牌
			result := make([]*bbproto.Game_CardInfo, 0)
			for i := 0; i < len(u.Cards); i++ {
				c := u.Cards[i]
				gc := room.ThCard2OGCard(c)
				result = append(result, gc)
			}
			handCard = append(handCard, result...)
		} else {

		}

	}
	return handCard
}


//解析每个人下注的金额
func getHandCoin(mydesk *room.ThDesk) []int64{
	result := make([]int64,len(mydesk.Users))
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户手牌
			result[i] = int64(u.HandCoin)
		} else {
			result[i]= int64(0)
		}
	}
	return result
}


//解析TurnCoin
func getTurnCoin(mydesk *room.ThDesk) []int64{
	result := make([]int64,len(mydesk.Users))
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户手牌
			result[i] = int64(u.TurnCoin)
		} else {
			result[i]= int64(0)
		}
	}
	return result
}




//手牌转换为OG可以使用的牌
func thPublicCard2OGC(desk *room.ThDesk) []*bbproto.Game_CardInfo {
	result := make([]*bbproto.Game_CardInfo, len(desk.PublicPai))
	for i := 0; i < len(desk.PublicPai); i++ {
		result[i] = room.ThCard2OGCard(desk.PublicPai[i])
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

	//1,初始化默认信息
	err := mydesk.Run()
	if err != nil {
		log.E("开始游戏的时候出错errMsg[%v]",err.Error())
		return err
	}

	//2.1 盲注开始押注
	err = mydesk.BlindBet()
	if err != nil {
		log.E("盲注下注的时候出错errMsg[%v]",err.Error())
		return err
	}

	//2.2,------------------------------------发送盲注的广播------------------------------------
	log.T("开始广播盲注的信息")
	blindB := &bbproto.Game_BlindCoin{}

	//初始化指针地址
	blindB.Banker = new(int32)
	blindB.Bigblindseat = new(int32)
	blindB.Smallblindseat = new(int32)

	//初始化默认值
	blindB.Tableid = &mydesk.Id	//deskid
	//blindB.Matchid = &room.ThGameRoomIns.Id //roomId
	*blindB.Banker = int32(mydesk.GetUserIndex(mydesk.Dealer))	//庄
	blindB.Bigblind = &mydesk.BigBlindCoin	//大盲注
	blindB.Smallblind = &mydesk.SmallBlindCoin	//小盲注
	*blindB.Bigblindseat = int32(mydesk.GetUserIndex(mydesk.BigBlind))	//大盲注座位号
	*blindB.Smallblindseat = int32(mydesk.GetUserIndex(mydesk.SmallBlind))	//小盲注座位号
	blindB.Coin	= mydesk.GetCoin()	//每个人手中的coin
	blindB.Handcoin = getHandCoin(mydesk)	//每个人下注的coin
	blindB.Pool	= &mydesk.Jackpot	//奖池

	mydesk.THBroadcastProto(blindB,0)
	//mydesk.Testb(blindB)
	log.T("广播盲注的信息完毕")

	//3,------------------------------------发送手牌的广播------------------------------------
	log.T("广播initCard的信息")
	initCardB := &bbproto.Game_InitCard{}

	//设置默认值
	initCardB.Tableid = new(int32)
	initCardB.ActionTime = new(int32)
	initCardB.DelayTime  = new(int32)
	initCardB.NextUser = new(int32)

	//设置初始化值
	*initCardB.Tableid = int32(mydesk.Id)
	initCardB.HandCard = getHandCard(mydesk)
	initCardB.PublicCard = thPublicCard2OGC(mydesk)
	initCardB.MinRaise = &mydesk.MinRaise
	//initCardB.NextUser = mydesk.GetResUserModelById(mydesk.BetUserNow).SeatNumber
	*initCardB.NextUser = int32(mydesk.GetUserIndex(mydesk.BetUserNow))
	initCardB.Seat = &mydesk.UserCount
	mydesk.THBroadcastProto(initCardB,0)
	//mydesk.Testb(initCardB)

	log.T("广播initCard的信息完毕")
	return nil
}

//处理押注的请求
func HandlerFollowBet(m  *bbproto.Game_FollowBet,a gate.Agent) error{
	log.T("处理用户押注的请求")
	seatId := m.GetSeat()
	desk := room.ThGameRoomIns.GetDeskById(m.GetTableid())
	desk.OgFollowBet(seatId)
	return nil
}


//处理加注
func HandlerRaiseBet(m *bbproto.Game_RaiseBet,a gate.Agent) error{
	seatId := m.GetSeat()
	desk := room.ThGameRoomIns.GetDeskById(m.GetTableid())
	desk.OGRaiseBet(seatId,m.GetCoin())
	return nil
}

//处理让牌
func HandlerCheckBet(m *bbproto.Game_CheckBet,a gate.Agent) error{
	seatId := m.GetSeat()
	desk := room.ThGameRoomIns.GetDeskById(m.GetTableid())
	desk.OGCheckBet(seatId)
	return nil
}


//处理让牌
func HandlerFoldBet(m *bbproto.Game_FoldBet,a gate.Agent) error{
	seatId := m.GetSeat()
	desk := room.ThGameRoomIns.GetDeskById(m.GetTableid())
	desk.OgFoldBet(seatId)
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

