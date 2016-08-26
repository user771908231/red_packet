package room

import (
	"casino_server/msg/bbprotogo"
	"sync"
	"time"
	"casino_server/conf/intCons"
	"casino_server/service/pokerService"
	"casino_server/conf/casinoConf"
	"casino_server/mode"
	"casino_server/service/userService"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"sort"
	"casino_server/common/log"
	"errors"
	"github.com/name5566/leaf/gate"
	"github.com/golang/protobuf/proto"
	"casino_server/utils/db"
	"sync/atomic"
	"casino_server/common/Error"
	"casino_server/utils/numUtils"
)


//联众德州,桌子状态
var (
	GAME_STATUS_READY int32 = 0        //准备
	GAME_STATUS_DEAL_CARDS int32 = 1   //发牌
	GAME_STATUS_PRECHIP int32 = 2        //盲注
	GAME_STATUS_FIRST_TURN int32 = 3        //第一轮
	GAME_STATUS_SECOND_TURN int32 = 4        //第二轮
	GAME_STATUS_THIRD_TURN int32 = 5        //第三轮
	GAME_STATUS_LAST_TURN int32 = 6        //第四轮
	GAME_STATUS_SHOW_RESULT int32 = 7        //完成
)

//德州扑克,牌桌的状态
var TH_DESK_STATUS_STOP int32 = 1                //没有开始的状态
var TH_DESK_STATUS_READY int32 = 2                 //游戏处于准备的状态
var TH_DESK_STATUS_RUN int32 = 3                //已经开始的状态
var TH_DESK_STATUS_LOTTERY int32 = 4             //已经开始的状态

var TH_DESK_ROUND1 int32 = 1                //第一轮押注
var TH_DESK_ROUND2 int32 = 2                //第二轮押注
var TH_DESK_ROUND3 int32 = 3                //第三轮押注
var TH_DESK_ROUND4 int32 = 4                //第四轮押注
var TH_DESK_ROUND_END int32 = 5             //完成押注


//押注的类型
var TH_DESK_BET_TYPE_BET int32 = 1                //押注
var TH_DESK_BET_TYPE_CALL int32 = 2        //跟注,和别人下相同的筹码
var TH_DESK_BET_TYPE_FOLD int32 = 3        //弃牌
var TH_DESK_BET_TYPE_CHECK int32 = 4                //让牌
var TH_DESK_BET_TYPE_RAISE int32 = 5                //加注
var TH_DESK_BET_TYPE_RERRAISE int32 = 6                //再加注
var TH_DESK_BET_TYPE_ALLIN int32 = 7                //全下


func init() {
	InitDeskConfig()
}


//桌子的一些配置
var ThdeskConfig struct {
					       //创建房间需要消耗的砖石
	CreateJuCountUnit        int32         //没多少局消耗多少钻石
	CreateFee                int64         //每多少局消耗多少砖石
	TH_TIMEOUT_DURATION      time.Duration //德州出牌的超时时间
	TH_TIMEOUT_DURATION_INT  int32         //德州出牌的超时时间
	TH_LOTTERY_DURATION      time.Duration //开奖的时间
	TH_DESK_LEAST_START_USER int32         //每局开始的最低人数
	TH_DESK_MAX_START_USER   int32         //桌子上最多坐多少人
	TH_GAME_SMALL_BLIND      int64         //默认小盲注是多少
}

//初始化桌子的配置参数
func InitDeskConfig() {

	//每4局消耗1颗砖石
	ThdeskConfig.CreateJuCountUnit = 4
	ThdeskConfig.CreateFee = 1

	//每个人出牌的等待时间
	ThdeskConfig.TH_TIMEOUT_DURATION = time.Second * 200
	ThdeskConfig.TH_TIMEOUT_DURATION_INT = 200

	//德州开奖的时间
	ThdeskConfig.TH_LOTTERY_DURATION = time.Second * 5

	ThdeskConfig.TH_DESK_LEAST_START_USER = 2
	ThdeskConfig.TH_DESK_MAX_START_USER = 9

	ThdeskConfig.TH_GAME_SMALL_BLIND = 10
}


/**
	一个德州扑克的房间
 */
type ThDesk struct {
	sync.Mutex
	Id                   int32                        //roomid
	MatchId              int32                        //matchId
	DeskOwner            uint32                       //房主的id
	RoomKey              string                       //room 自定义房间的钥匙
	CreateFee            int64                        //创建房间的费用
	GameType             int32                        //桌子的类型,1,表示自定义房间,2表示锦标赛的
	InitRoomCoin         int64                        //进入这个房间的roomCoin 带入金额标准是多少
	JuCount              int32                        //这个桌子最多能打多少局
	JuCountNow           int32                        //这个桌子已经玩了多少局
	PreCoin              int64                        //前注的金额
	SmallBlindCoin       int64                        //小盲注的押注金额
	BigBlindCoin         int64                        //大盲注的押注金额
	BeginTime            time.Time                    //游戏开始时间
	EndTime              time.Time                    //游戏结束时间

	Dealer               uint32                       //庄家
	BigBlind             uint32                       //大盲注
	SmallBlind           uint32                       //小盲注
	RaiseUserId          uint32                       //加注的人的Id,一轮结束的判断需要按照这个人为准
	NewRoundFirstBetUser uint32                       //新一轮,开始押注的第一个人//第一轮默认是小盲注,但是当小盲注弃牌之后,这个人要滑倒下一家去
	BetUserNow           uint32                       //当前押注人的Id

	GameNumber           int32                        //每一局游戏的游戏编号
	Users                []*ThUser                    //坐下的人
	PublicPai            []*bbproto.Pai               //公共牌的部分
	UserCount            int32                        //玩游戏的总人数
	UserCountOnline      int32                        //在线的人数
	ReadyCount           int32                        //已经准备的用户数
	Status               int32                        //牌桌的状态
	BetAmountNow         int64                        //当前的押注金额是多少
	RoundCount           int32                        //第几轮
	Jackpot              int64                        //奖金池
	edgeJackpot          int64                        //边池
	MinRaise             int64                        //最低加注金额
	AllInJackpot         []*pokerService.AllInJackpot //allin的标记
}

/**
	新生成一个德州的桌子
 */
func NewThDesk() *ThDesk {
	result := new(ThDesk)
	result.Id = newThDeskId()
	result.UserCount = 0
	result.UserCountOnline = 0
	result.Dealer = 0                //不需要创建  默认就是为空
	result.Status = TH_DESK_STATUS_STOP
	result.BetUserNow = 0
	result.BigBlind = 0
	result.SmallBlind = 0
	result.Users = make([]*ThUser, ThdeskConfig.TH_DESK_MAX_START_USER)
	result.RaiseUserId = 0
	result.RoundCount = 0
	result.NewRoundFirstBetUser = 0
	result.edgeJackpot = 0
	result.SmallBlindCoin = ThGameRoomIns.SmallBlindCoin
	result.BigBlindCoin = 2 * ThGameRoomIns.SmallBlindCoin
	result.Status = TH_DESK_STATUS_STOP        //游戏还没有开始的状态
	result.GameType = intCons.GAME_TYPE_TH_CS                          //游戏桌子的类型
	result.JuCount = 0
	result.JuCountNow = 1                //默认从第一局开始
	return result
}


//获取数据库的id
func newThDeskId() int32 {
	id, err := db.GetNextSeq(casinoConf.DBT_T_TH_DESK)
	if err != nil {
		log.E("new 一个desk的id 的时候出错!")
		return int32(0)
	}
	return int32(id)
}

func (t *ThDesk) LogString() {
	log.T("当前desk[%v]的信息:-----------------------------------begin----------------------------------", t.Id)
	log.T("当前desk[%v]的信息的状态status[%v]", t.Id, t.Status)
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			//log.T("当前desk[%v]的user[%v]的状态status[%v],牌的信息[%v]", t.Id, u.UserId, u.Status, u.HandCards)
			log.T("当前desk[%v]的user[%v]的状态status[%v],HandCoin[%v],TurnCoin[%v],RoomCoin[%v]", t.Id, u.UserId, u.Status,u.HandCoin,u.TurnCoin,u.RoomCoin)
		}
	}
	log.T("当前desk[%v]的信息的状态,MatchId[%v]", t.Id, t.MatchId)
	log.T("当前desk[%v]的信息的状态,RoomKey[%v]", t.Id, t.RoomKey)
	log.T("当前desk[%v]的信息的状态,CreateFee[%v]", t.Id, t.CreateFee)
	log.T("当前desk[%v]的信息的状态,GameType[%v]", t.Id, t.GameType)
	log.T("当前desk[%v]的信息的状态,InitRoomCoin[%v]", t.Id, t.InitRoomCoin)
	log.T("当前desk[%v]的信息的状态,JuCount[%v]", t.Id, t.JuCount)
	log.T("当前desk[%v]的信息的状态,JuCountNow[%v]", t.Id, t.JuCountNow)
	log.T("当前desk[%v]的信息的状态,PreCoin[%v]", t.Id, t.PreCoin)
	log.T("当前desk[%v]的信息的状态,SmallBlindCoin[%v]", t.Id, t.SmallBlindCoin)
	log.T("当前desk[%v]的信息的状态,BigBlindCoin[%v]", t.Id, t.BigBlindCoin)
	log.T("当前desk[%v]的信息的状态,小盲注[%v]", t.Id, t.SmallBlind)
	log.T("当前desk[%v]的信息的状态,大盲注[%v]", t.Id, t.BigBlind)
	log.T("当前desk[%v]的信息的状态,BeginTime[%v]", t.Id, t.BeginTime)
	log.T("当前desk[%v]的信息的状态,EndTime[%v]", t.Id, t.EndTime)
	log.T("当前desk[%v]的信息的状态,RaiseUserId[%v]", t.Id, t.RaiseUserId)
	log.T("当前desk[%v]的信息的状态,MinRaise[%v]", t.Id, t.MinRaise)
	log.T("当前desk[%v]的信息的状态,BetUserNow[%v]", t.Id, t.BetUserNow)
	log.T("当前desk[%v]的信息的状态,GameNumber[%v]", t.Id, t.GameNumber)
	log.T("当前desk[%v]的信息的状态,ReadyCount[%v]", t.Id, t.ReadyCount)
	log.T("当前desk[%v]的信息的状态,总人数SeatedCount[%v],在线人数[%v]", t.Id, t.UserCount, t.UserCountOnline)
	log.T("当前desk[%v]的信息的状态,压注人[%v]", t.Id, t.BetUserNow)
	log.T("当前desk[%v]的信息的状态,压注轮次[%v]", t.Id, t.RoundCount)
	log.T("当前desk[%v]的信息的状态,NewRoundFirstBetUser[%v]", t.Id, t.NewRoundFirstBetUser)
	log.T("当前desk[%v]的信息的状态,总共押注Jackpot[%v]", t.Id, t.Jackpot)
	log.T("当前desk[%v]的信息的状态,edgeJackpot[%v]", t.Id, t.edgeJackpot)
	log.T("当前desk[%v]的信息的状态,当前加注的人BetUserRaiseUserId[%v]", t.Id, t.RaiseUserId)
	log.T("当前desk[%v]的信息:-----------------------------------end----------------------------------", t.Id)
}

//打印测试信息用
func (t *ThDesk) LogStirngWinCoin() {
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("用户[%v]的wincoin[%v],detail[%v]", u.UserId, u.winAmount, u.winAmountDetail)
		}
	}
}


//判断是否重新进入
func (t *ThDesk) IsrepeatIntoWithRoomKey(userId uint32, a gate.Agent) bool {
	//1,先判断用户是否已经在房间了
	for i := 0; i < len(t.Users); i ++ {
		u := t.Users[i]
		if u != nil && u.UserId == userId {
			//如果u!=nil 那么
			log.T("用户[%v]断线重连", userId)
			u.agent = a                                                //设置用户的连接
			u.IsBreak = false                //设置用户的离线状态
			u.IsLeave = false
			u.UpdateAgentUserData(a, t.Id, t.MatchId)
			u.Update2redis()                //更新用户数据到redis
			t.AddUserCountOnline()
			return true
		}
	}
	return false
}

/**
	为桌子增加一个人
 */
func (t *ThDesk) AddThUser(userId uint32, userStatus int32, a gate.Agent) (*ThUser, error) {
	//1,从redis得到redisUser
	redisUser := userService.GetUserById(userId)
	//2,通过userId 和agent 够做一个thuser
	thUser := NewThUser()
	thUser.UserId = userId
	thUser.agent = a
	thUser.Status = userStatus        //刚进房间的玩家
	thUser.deskId = t.Id                //桌子的id
	thUser.MatchId = t.MatchId
	thUser.GameNumber = t.GameNumber
	thUser.NickName = *redisUser.NickName                //todo 测试阶段,把nickName显示成用户id
	thUser.RoomCoin = t.InitRoomCoin
	thUser.InitialRoomCoin = thUser.RoomCoin
	log.T("初始化thuser的时候coin[%v]:,roomCoin[%v]", thUser.GetCoin(), thUser.GetRoomCoin())

	//3,添加thuser
	err := t.addThuserBean(thUser)
	if err != nil {
		log.E("增加user【%v】到desk【%v】失败", thUser.UserId, t.Id)
		return nil, errors.New("增加user失败")
	}

	//4, 把用户的信息绑定到agent上
	thUser.UpdateAgentUserData(a, t.Id, t.MatchId)
	NewRedisThuser(thUser)

	//5,等待的用户加1
	t.AddUserCount()
	t.AddUserCountOnline()
	if userStatus == TH_USER_STATUS_READY {
		t.AddReadyCount()
	}
	return thUser, nil
}


//增加一个user实体
func (t *ThDesk) addThuserBean(user *ThUser) error {
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] == nil {
			user.Seat = int32(i)                //给用户设置位置编号
			t.Users[i] = user
			return nil
		}

		if (i + 1) == len(t.Users) {
			log.E("玩家加入桌子失败")
			return errors.New("加入房间失败")
		}
	}

	return nil
}

//用户准备的时候需要判断
func (t *ThDesk) Ready(userId uint32) error {
	user := t.GetUserByUserId(userId)

	//1,如果是在游戏的过程中,则准备失败
	if t.Status == TH_DESK_STATUS_RUN {
		log.E("desk[%v]已经在游戏中了,user[%v]不能准备", t.Id, userId)
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_GAME_READY_REPEAT), "已经在游戏中,不能准备")
	}

	//1,如果用户已经准备,则返回重复准备
	if user.Status == TH_USER_STATUS_READY {
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_GAME_READY_REPEAT), "已经准备了")
	}

	//2,如果用户余额不足,则准备失败
	if !t.IsUserCoinEnough(user) {
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_GAME_READY_CHIP_NOT_ENOUGH), "筹码不足")
	}


	//3,准备成功
	user.Status = TH_USER_STATUS_READY
	t.AddReadyCount()
	return nil
}

//查看是否全部准备好了
func (t *ThDesk) IsAllReady() bool {
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			if u.Status != TH_USER_STATUS_READY {
				return false
			}
		}
	}

	return true
}

//  用户退出德州游戏的房间,rmUser 需要在事物中进行,离开德州需要更具桌子的类型来做不同的处理
func (t *ThDesk) LeaveThuser(userId uint32) error {

	//离线的时候,桌子统一的处理
	user := t.GetUserByUserId(userId)
	user.IsLeave = true     //设置状态为离开
	t.SubUserCountOnline()        //房间的在线人数减一

	//根据不同的游戏类型做不同的处理
	if t.GameType == intCons.GAME_TYPE_TH {
		//自定义房间

	} else if t.GameType == intCons.GAME_TYPE_TH_CS {
		//用户直接放弃游戏,设置roomCoin=0,并且更新rankxin
		user.AddRoomCoin(user.RoomCoin)
		ChampionshipRoom.UpdateUserRankInfo(user.UserId, user.MatchId, user.RoomCoin)
		//竞标赛的房间
		ChampionshipRoom.SubOnlineCount()        //竞标赛的在线人数-1
	}

	//给离开的人发送信息
	ret := &bbproto.Game_ACKLeaveDesk{}
	ret.Result = &intCons.ACK_RESULT_SUCC
	user.WriteMsg(ret)

	//离开之后,需要广播一次sendGameInfo
	t.BroadGameInfo()
	return nil
}

//设置用户为掉线的状态
func (t *ThDesk) SetOfflineStatus(userId uint32) error {

	u := t.GetUserByUserId(userId)
	//1,设置状态为断线
	//这里需要保存掉线钱的状态
	u.IsBreak = true        //设置为掉线的状态
	t.SubUserCountOnline()
	return nil
}

//初始化前注的信息
func (t *ThDesk) OinitPreCoin() error {
	log.T("开始一局新的游戏,现在开始初始化前注的信息[%v]", t.PreCoin)

	//每个都减少前注的金额
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.IsBetting() {
			t.calcPreCoin(u.UserId, t.PreCoin)
		}
	}

	//发送前主的广播
	ret := bbproto.NewGame_PreCoin()
	ret.Coin = t.GetRoomCoin()
	*ret.Pool = t.Jackpot
	*ret.Tableid = t.Id
	ret.Precoin = t.GetPreCoin()

	t.THBroadcastProtoAll(ret)
	log.T("开始一局新的游戏,现在开始初始化前注的信息完毕....")
	return nil
}

// 盲注开始押注
func (t *ThDesk) InitBlindBet() error {
	log.T("开始一局新的游戏,现在开始初始化盲注的信息")
	//小盲注押注
	t.calcBetCoin(t.SmallBlind, t.SmallBlindCoin)

	//大盲注押注
	t.calcBetCoin(t.BigBlind, t.BigBlindCoin)

	//发送盲注的信息
	log.T("开始广播盲注的信息")
	blindB := bbproto.NewGame_BlindCoin()

	//初始化默认值
	blindB.Tableid = &t.Id        //deskid
	//blindB.Matchid = &room.ThGameRoomIns.Id //roomId
	*blindB.Banker = t.GetUserByUserId(t.Dealer).Seat        //庄
	blindB.Bigblind = &t.BigBlindCoin        //大盲注
	blindB.Smallblind = &t.SmallBlindCoin        //小盲注
	*blindB.Bigblindseat = t.GetUserByUserId(t.BigBlind).Seat        //      //大盲注座位号
	*blindB.Smallblindseat = t.GetUserByUserId(t.SmallBlind).Seat        //     //小盲注座位号
	//blindB.Coin = t.GetCoin()        //每个人手中的coin
	blindB.Coin = t.GetRoomCoin()
	blindB.Handcoin = t.GetHandCoin()        //每个人下注的coin
	blindB.Pool = &t.Jackpot        //奖池
	t.THBroadcastProto(blindB, 0)
	log.T("广播盲注的信息完毕")

	log.T("开始一局新的游戏,现在开始初始化盲注的信息完毕....")
	return nil
}

/**
	游戏开始的时候,初始化玩家的信息
 */
func (t *ThDesk) InitUserBeginStatus() error {
	log.T("开始一局新的游戏,开始初始化用户的状态")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		//判断用户是否为空
		if u == nil {
			continue
		}

		log.T("用户[%v]的status[%v],BreakStatus[%v],:", u.UserId, u.Status, u.IsBreak)
		//新的一局开始,设置用户的状态
		u.HandCoin = 0
		u.TurnCoin = 0
		u.winAmount = 0
		u.TotalBet4calcAllin = 0
		u.TotalBet = 0                                //新的一局游戏开始,把总的押注金额设置为0
		u.winAmountDetail = nil

		//如果用户的余额不足或者用户的状态是属于断线的状态,则设置用户为等待入座
		//if u.RoomCoin <= (t.BigBlindCoin + t.PreCoin) || u.IsBreak == true {
		if u.RoomCoin <= (t.BigBlindCoin + t.PreCoin) {
			log.T("由于用户[%v] status[%v],的roomCoin[%v] <= desk.BigBlindCoin 所以设置用户为TH_USER_STATUS_WAITSEAT", u.UserId, u.IsBreak, u.RoomCoin, t.BigBlindCoin)
			u.Status = TH_USER_STATUS_WAITSEAT        //只是坐下,没有游戏中
			continue
		}

		//用户不是离线的状态,并且,用户已经准备好了,则可以开始游戏
		if u.IsBreak == false && u.Status == TH_USER_STATUS_READY {
			log.T("由于用户[%v]的status[%v]BreakStatus[%v],所以设置状态为TH_USER_STATUS_BETING", u.UserId, u.Status, u.IsBreak)
			u.Status = TH_USER_STATUS_BETING
		}
	}

	//------------------------------------由于联众前端设计的问题...这里的user需要重新排列user的顺序------------------------------------
	usersTemp := make([]*ThUser, len(t.Users))
	copy(usersTemp, t.Users)
	log.T("原来的thsuers:[%v]", t.Users)
	log.T("复制的thsuers:[%v]", usersTemp)
	log.T("排序之前的thusers【%v】", t.Users)

	//是原来的thusers置为nil
	for i := 0; i < len(t.Users); i++ {
		t.Users[i] = nil
	}

	//排序游戏中的玩家
	for i := 0; i < len(usersTemp); i++ {
		u := usersTemp[i]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.addThuserBean(u)
		}
	}

	for i := 0; i < len(usersTemp); i++ {
		u := usersTemp[i]
		if u != nil && u.Status != TH_USER_STATUS_BETING {
			t.addThuserBean(u)
		}
	}
	log.T("排序之后的thusers【%v】", t.Users)

	log.T("开始一局新的游戏,初始化用户的状态完毕")
	return nil
}

/**
	初始化纸牌的信息
 */
func (t *ThDesk) OnInitCards() error {
	log.T("开始一局新的游戏,初始化牌的信息")
	var total = int(2 * ThdeskConfig.TH_DESK_MAX_START_USER + 5); //人数*手牌+5张公共牌
	totalCards := pokerService.RandomTHPorkCards(total)        //得到牌

	//得到所有的牌,前五张为公共牌,后边的每两张为手牌
	t.PublicPai = totalCards[0:5]
	log.T("初始化得到的公共牌的信息:")
	//给每个人分配手牌
	for i := 0; i < len(t.Users); i++ {
		//只有当用不为空,并且是在游戏中的状态的时候,才可以发牌
		if t.Users[i] != nil && t.Users[i].Status == TH_USER_STATUS_BETING {
			begin := i * 2 + 5
			end := i * 2 + 5 + 2
			t.Users[i].HandCards = totalCards[begin:end]
			log.T("用户[%v]的手牌[%v]", t.Users[i].UserId, t.Users[i].HandCards)
			t.Users[i].thCards = pokerService.GetTHPoker(t.Users[i].HandCards, t.PublicPai, 5)
			log.T("用户[%v]的:拍类型,所有牌[%v],th[%v]", t.Users[i].UserId, t.Users[i].thCards.ThType, t.Users[i].thCards.Cards, t.Users[i].thCards)
		}
	}
	log.T("开始一局新的游戏,初始化牌的信息完毕...")
	return nil
}

//广播porto消息的通用方法
func (t *ThDesk) THBroadcastProto(p proto.Message, ignoreUserId uint32) error {
	log.Normal("开始广播proto消息【%v】", p)
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]                //给这个玩家发送广播信息
		if u != nil && u.UserId != ignoreUserId && u.IsLeave == false && u.IsBreak == false {
			a := t.Users[i].agent
			a.WriteMsg(p)
		}
	}
	return nil
}



//给全部人发送广播
func (t *ThDesk) BroadcastTestResult(p *bbproto.Game_TestResult) error {
	//发送的时候,初始化自己的排名
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]                //给这个玩家发送广播信息
		if u != nil && u.IsLeave != true {

			//只有锦标赛的时候才可以有排名
			if t.GameType == intCons.GAME_TYPE_TH_CS {
				*p.Rank = ChampionshipRoom.GetRankByuserId(u.UserId)
			}
			a := t.Users[i].agent
			a.WriteMsg(p)
		}
	}
	return nil
}

//给全部人发送广播
func (t *ThDesk) THBroadcastProtoAll(p proto.Message) error {
	return t.THBroadcastProto(p, 0)
}


/**
	返回res需要的User实体
 */
func (t *ThDesk) GetResUserModel() []*bbproto.THUser {
	result := make([]*bbproto.THUser, ThdeskConfig.TH_DESK_MAX_START_USER)
	//就坐的人
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil {
			result[i] = t.Users[i].trans2bbprotoThuser()
			*result[i].SeatNumber = int32(i)
		} else {
			result[i] = &bbproto.THUser{}
			result[i].SeatNumber = new(int32)        //设置为0
		}
	}

	return result
}

func (t *ThDesk) GetResUserModelById(userId uint32) *bbproto.THUser {
	var result *bbproto.THUser
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil &&  t.Users[i].UserId == userId {
			result = t.Users[i].trans2bbprotoThuser()
		}
	}
	log.T("通过userId得到的bbproto.THUser的情况【%v】", result)
	return result
}

// 返回res需要的User实体 并且排序,排序规则是,当前用户排在第一个
func (t *ThDesk) GetResUserModelClieSeq(userId uint32) []*bbproto.THUser {
	//需要根据当前用户的Userid来进行排序
	users := t.GetResUserModel()
	var userIndex int = 0
	for i := 0; i < len(users); i++ {
		if users[i] != nil && users[i].User != nil && *(users[i].User.Id) == userId {
			userIndex = i
			break
		}
	}

	for i := 0; i < len(users); i++ {
		u := users[(i + userIndex) % len(users)]
		if u != nil {
			*u.SeatNumber = int32(i)
		}
	}

	log.T("得到排序后的bbproto.THUser的情况的信息[%v]", users)
	return users
}




// 	初始化第一个押注的人,当前押注的人
func (t *ThDesk) OninitThDeskBeginStatus() error {
	log.T("开始一局游戏,现在初始化desk的信息")
	//设置德州desk状态//设置状态为开始游戏
	t.Status = TH_DESK_STATUS_RUN

	userTemp := make([]*ThUser, len(t.Users))
	copy(userTemp, t.Users)
	//这里需要定义一个庄家,todo 暂时默认为第一个,后边再修改
	var dealerIndex int = 0;
	if t.Dealer == 0 {
		log.T("游戏没有庄家,现在默认初始化第一个人为庄家")
		t.Dealer = t.Users[0].UserId
	} else {
		dealerIndex = t.GetUserIndex(t.Dealer)
		for i := dealerIndex; i < len(t.Users); i++ {
			u := t.Users[(i + 1) % len(t.Users)]
			if u != nil && u.Status == TH_USER_STATUS_BETING {
				t.Dealer = u.UserId
			}
		}
	}

	//设置小盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.SmallBlind = u.UserId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	//设置大盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.BigBlind = u.UserId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	/**
		设置第一个押注的人
	 */

	if t.UserCountOnline == int32(2) {
		//如果只有两个人,当前押注的人是小盲注
		t.BetUserNow = t.SmallBlind
	} else {
		//设置当前押注的人
		for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
			u := userTemp[(i + 1) % len(userTemp)]
			if u != nil && u.Status == TH_USER_STATUS_BETING {
				t.BetUserNow = u.UserId
				userTemp[(i + 1) % len(userTemp)] = nil
				break
			}
		}
	}

	t.RaiseUserId = t.BetUserNow        //第一个加注的人
	t.NewRoundFirstBetUser = t.SmallBlind           //新一轮开始默认第一个押注的人,第一轮默认是小盲注
	t.RoundCount = TH_DESK_ROUND1
	t.BetAmountNow = t.BigBlindCoin                   //设置第一次跟住时的跟注金额应该是多少
	t.MinRaise = t.BigBlindCoin
	t.Jackpot = 0
	t.edgeJackpot = 0
	t.AllInJackpot = nil                          // 初始化allInJackpot 为空

	t.LogString()
	log.T("开始一局游戏,现在初始化desk的信息完毕...")
	return nil
}


//判断是否是开奖的时刻
/**
开奖的时候
1,目前已经是第四轮
2,计算出来的下一个押注这和当前押注的是同一个人
3,即使没有到第四轮,但是所有人都all in 了,那么还是要开奖
 */
func (t *ThDesk) Tiem2Lottery() bool {
	//如果处于押注状态的人只有一个人了,那么是开奖的时刻
	//

	log.T("判断是否应该开奖,打印每个人的信息://1,刚上桌子,2,坐下,3,ready 4 押注中,5,allin,6,弃牌,7,等待结算,8,已经结算,9,离开")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("[%v]判断是否应该开奖,打印user[%v]的状态[%v]:", i, u.UserId, u.Status)
		}
	}

	var betingCount int = 0
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil && t.Users[i].Status == TH_USER_STATUS_BETING {
			betingCount ++
		}
	}

	log.T("当前处于押注中的人数是[%v]", betingCount)
	//如果押注的人只有一个人了,那么是开奖的时刻
	if betingCount <= 1 {
		log.T("现在处于押注中(beting)状态的人,只剩下一个了,所以直接开奖")
		return true
	}

	//第四轮,并且计算出来的押注人和start是同一个人
	if t.RoundCount == TH_DESK_ROUND_END {
		log.T("现在处于第[%v]轮押注,所以可以直接开奖", t.RoundCount)
		return true
	}

	//如果只有一个人没有all in  或者全部都all in 了也要开牌

	return false
}



//计算牌面是否赢
func (t *ThDesk) CalcThcardsWin() error {
	log.T("开始计算谁的牌是赢牌:")

	log.T("打印每个人牌的信息:")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED {
			log.T("玩家[%v] ", u.UserId)
			log.T("牌的信息:牌类型[%v],", u.thCards.ThType)
			log.T("所有牌[%v]", u.thCards.Cards)
		}
	}

	userWin := t.Users[0]                //最大的牌的userId
	for i := 1; i < len(t.Users); i++ {
		if t.Less(userWin, t.Users[i]) {
			userWin = t.Users[i]
		}
	}


	//1.3打印得到那个人的信息
	if userWin == nil {
		log.E("服务器出错,没有找到赢牌的人...")
		return errors.New("没有找到赢牌的人")
	} else {
		log.T("得到的牌最大的user[%v]的信息[%v]", userWin.UserId, userWin)
	}

	//赢牌的人依次置为1,每个用户的牌都需要和最大的用户的牌来想比较
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil &&
		u.Status == TH_USER_STATUS_WAIT_CLOSED  &&
		pokerService.ThCompare(userWin.thCards, u.thCards) == pokerService.THPOKER_COMPARE_EQUALS {
			u.thCards.IsWin = true
		}
	}

	//------------------------------------------------打印测试信息
	log.T("开始计算谁的牌是赢牌,计算出来的结果:")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED {
			log.T("user[%v]的牌 isWin[%v]", u.UserId, u.thCards.IsWin)
		}
	}

	return nil
}

//比较两张牌的大小
func (t *ThDesk) Less(u1, u2 *ThUser) bool {

	if u1 == nil || u1.Status != TH_USER_STATUS_WAIT_CLOSED {
		return true
	}

	if u2 == nil || u2.Status != TH_USER_STATUS_WAIT_CLOSED {
		return false
	}

	//必将两个人的牌,u1的牌是否大于u2的牌
	ret := pokerService.ThCompare(u1.thCards, u2.thCards)
	if ret == pokerService.THPOKER_COMPARE_BIG {
		return false
	} else {
		return true
	}
}

//设置用户的状态喂等待开奖
func (t *ThDesk) SetStatusWaitClose() error {
	log.T("开奖之前答应每个人的状态,并且修改为等待结算")
	//设置用户的状态为等待开奖
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("用户[%v].nickname[%v]的status[%v]", u.UserId, u.NickName, u.Status)
			u.InitWait()        //不再等待
			if u.Status == TH_USER_STATUS_ALLINING || u.Status == TH_USER_STATUS_BETING {
				//如果用户当前的状态是押注中,或者all in,那么设置用户的状态喂等待结算
				u.Status = TH_USER_STATUS_WAIT_CLOSED
			} else {
				u.Status = TH_USER_STATUS_CLOSED
			}
		}
	}
	return nil
}

//计算用户输赢盈眶
func (t *ThDesk) calcUserWinAmount() error {

	//计算每个allin池子的将近分配
	for i := 0; i < len(t.AllInJackpot); i++ {
		log.T("现在开始开奖,计算allInJackpot[%v]的奖励....", i)
		a := t.AllInJackpot[i]
		if a != nil {
			//对这个奖池做计算
			/**
				1,几个人的牌是赢牌
				2,user的状态必须是没有结算的状态
			 */
			var winCount int = t.GetWinCount()
			bonus := a.Jackpopt / int64(winCount)        //每个人赢的奖金
			log.T("allInJackpot[%v]的总共有[%v]人获得奖励[%v],平均[%v]", i, winCount, a.Jackpopt, bonus)
			//这里吧奖金发放给每个人之后,需要把这局allin的人排除掉,再来计算剩下的人的将近
			//牌的方式只需要把这个人的状态设置为已经结清就行了
			for j := 0; j < len(t.Users); j++ {
				//todo 这里的将近可以选择使用一个数组存醋,方便clien制作动画
				//todo 目前只是计算总的金额
				u := t.Users[j]

				//如果用户nil,则直接开始下一次循环
				if u == nil {
					continue
				}

				//判断用户是否得奖
				if u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
					//可以发送奖金
					log.T("用户[%v].status[%v],iswin[%v]在allin.index[%v]活的奖金[%v]", u.UserId, u.Status, u.thCards.IsWin, i, bonus)
					u.AddWinAmount(bonus)
					u.AddRoomCoin(bonus)
					u.winAmountDetail = append(u.winAmountDetail, bonus)
					u.Update2redis()        //把数据保存到redis中
				}

				//如果用户是这个奖金池all in的用户,则此用户设置喂已经结清的状态
				if u.UserId == a.UserId {
					u.Status = TH_USER_STATUS_CLOSED
				}
			}
		}
	}

	//计算边池的奖金	t.bianJackpot,同样需要看是几个人赢,然后评分将近
	log.T("现在开始开奖,计算边池的奖励....")
	bwinCount := t.GetWinCount()

	log.T("获奖的人数:[%v]", bwinCount)
	if bwinCount != 0 {
		var bbonus int64 = t.edgeJackpot / int64(bwinCount)
		for i := 0; i < len(t.Users); i++ {
			u := t.Users[i]

			if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
				//
				//对这个用户做结算...
				log.T("现在开始开奖,计算边池的奖励,user[%v]得到[%v]....", u.UserId, bbonus)
				u.AddWinAmount(bbonus)
				u.AddRoomCoin(bbonus)
				u.winAmountDetail = append(u.winAmountDetail, bbonus)        //详细的奖励(边池主池分开)
				u.Update2redis()
			}

			//设置为结算完了的状态
			if u != nil {
				u.Status = TH_USER_STATUS_CLOSED        //结算完了之后需要,设置用户的状态为已经结算
			}

		}
		//log.T	("现在开始开奖,计算奖励之后t.getWinCoinInfo()[%v]", t.getWinCoinInfo())
	}

	return nil

}


//开奖
func (t *ThDesk) Lottery() error {
	log.T("现在开始开奖,并且发放奖励....")

	//todo 开奖之前 是否需要把剩下的牌 全部发完**** 目前是不可能
	t.Status = TH_DESK_STATUS_LOTTERY

	//设置用户的状态都为的等待开奖
	t.SetStatusWaitClose()

	//需要计算本局allin的奖金池
	t.CalcAllInJackpot()

	//计算用户输赢情况
	t.calcUserWinAmount()

	//todo 这里需要删除 打印测试信息的代码
	t.LogStirngWinCoin()

	//保存数据到数据库
	t.SaveLotteryData()

	//广播开奖的额结果
	t.broadLotteryResult()

	//开奖之后,设置状态为 没有开始游戏
	t.afterLottery()

	//判断游戏是否结束
	if t.isEnd() {
		t.End()
	} else {
		//表示不能继续开始游戏
		time.Sleep(ThdeskConfig.TH_LOTTERY_DURATION)        //开奖的延迟
		go t.Run()
	}
	return nil
}



//判断开奖之后是否可以继续游戏
func (t *ThDesk) isEnd() bool {
	if t.GameType == intCons.GAME_TYPE_TH {
		//如果是自定义房间
		//1,局数
		//2,人数
		log.T("判断自定义的desk是否结束游戏t.jucount[%v],t.jucountnow[%v],", t.JuCount, t.JuCountNow)
		if t.JuCountNow <= t.JuCount {
			return false
		} else {
			return true
		}
	} else if t.GameType == intCons.GAME_TYPE_TH_CS {
		//判断锦标赛有没有结束,如果所有的desk都已经stop了,则表示游戏结束
		if ChampionshipRoom.allStop() {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}


//广播开奖的结果
func (t *ThDesk) broadLotteryResult() error {

	//1.发送输赢结果
	result := bbproto.NewGame_TestResult()
	*result.Tableid = t.Id                           //桌子
	result.BCanShowCard = t.GetBshowCard()           //
	result.BShowCard = t.GetBshowCard()              //亮牌
	result.Handcard = t.GetHandCard()                //手牌
	result.WinCoinInfo = t.getWinCoinInfo()
	result.HandCoin = t.GetHandCoin()
	result.CoinInfo = t.getCoinInfo()                //每个人的输赢情况
	t.BroadcastTestResult(result)

	//2.发送每个人在锦标赛中目前的排名
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			//开始风别给每个人发送自己的排名信息

		}
	}

	return nil

}

//开奖之后的处理
func (t *ThDesk) afterLottery() error {
	//1,设置游戏桌子的状态
	log.T("开奖结束,设置desk的状态为stop")
	t.Status = TH_DESK_STATUS_STOP        //设置为没有开始开始游戏
	t.ReadyCount = 0; //准备的人数为0
	t.JuCountNow ++

	//2,设置用户的状态
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.IsBreak == false {
			if t.GameType == intCons.GAME_TYPE_TH {
				//如果是自定义的房间,设置每个人都是坐下的状态
				u.Status = TH_USER_STATUS_SEATED
			} else if t.GameType == intCons.GAME_TYPE_TH_CS {
				if t.IsUserCoinEnough(u) {
					//如果是锦标赛的房间,用户的钱足够
					u.Status = TH_USER_STATUS_READY
					t.AddReadyCount()
				} else {
					//如果是锦标赛的房间,用户的钱不够
					u.Status = TH_USER_STATUS_SEATED
				}
			}
		}
	}

	//3,

	return nil
}

//判断用户的余额是否足够开始下一场游戏
func (t *ThDesk) IsUserCoinEnough(u *ThUser) bool {
	if u.RoomCoin < t.PreCoin + t.BigBlindCoin {
		return false
	} else {
		return true
	}
}

//保存数据到数据库

//这里需要根据游戏类型的不同来分别存醋

func (t *ThDesk)  SaveLotteryData() error {

	if t.GameType == intCons.GAME_TYPE_TH {
		//自定义房间
		return t.SaveLotteryDatath()
	} else if t.GameType == intCons.GAME_TYPE_TH_CS {
		//锦标赛
		return t.SaveLotteryDatacsth()
	}

	return nil

}

func (t *ThDesk) SaveLotteryDatath() error {
	log.T("一局游戏结束,开始保存游戏的数据到数据库")

	//为每一局保存一组数据
	deskRecord := &mode.T_th_desk_record{}
	deskRecord.DeskId = t.Id
	deskRecord.BeginTime = t.BeginTime
	deskRecord.UserIds = ""

	//循环对每个人做处理
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u == nil {
			continue
		}

		if u.Status != TH_USER_STATUS_CLOSED {
			log.E("保存用户[%v]信息到数据库的时候出错,状态[%v]不正确", u.UserId, u.Status)
			continue
		}

		//1,修改user在redis中的数据
		userService.FlashUser2Mongo(u.UserId)                        //刷新redis中的数据到mongo
		//2,保存游戏相关的数据
		//todo  游戏相关的数据结构 还没有建立,
		thData := &mode.T_th_record{}
		thData.Mid = bson.NewObjectId()
		thData.BetAmount = u.TotalBet
		thData.UserId = u.UserId
		thData.DeskId = u.deskId
		thData.WinAmount = u.winAmount - u.TotalBet
		thData.Blance = *(userService.GetUserById(u.UserId).Coin)
		thData.BeginTime = time.Now()
		thData.GameNumber = t.GameNumber
		db.InsertMgoData(casinoConf.DBT_T_TH_RECORD, thData)

		//获取游戏数据
		userRecord := mode.BeanRecord{}
		userRecord.UserId = u.UserId
		userRecord.NickName = u.NickName
		userRecord.WinAmount = u.winAmount - u.TotalBet

		deskRecord.Records = append(deskRecord.Records, userRecord)
		idStr, _ := numUtils.Uint2String(u.UserId)
		deskRecord.UserIds = strings.Join([]string{deskRecord.UserIds, idStr}, ",")
	}

	log.T("开始保存DBT_T_TH_DESK_RECORD的信息")
	//保存桌子的用户信息
	db.InsertMgoData(casinoConf.DBT_T_TH_DESK_RECORD, deskRecord)
	return nil

}

func (t *ThDesk) SaveLotteryDatacsth() error {
	log.T("一局游戏结束,开始保存游戏的数据到数据库")

	//为每一局保存一组数据
	deskRecord := &mode.T_th_desk_record{}
	deskRecord.DeskId = t.Id
	deskRecord.BeginTime = t.BeginTime
	deskRecord.UserIds = ""

	//循环对每个人做处理
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u == nil {
			continue
		}

		if !u.IsClose() {
			log.E("保存用户[%v]信息到数据库的时候出错,状态[%v]不正确", u.UserId, u.Status)
			continue
		}

		//1,修改user在redis中的数据
		userService.FlashUser2Mongo(u.UserId)                        //刷新redis中的数据到mongo
		//2,保存游戏相关的数据
		//todo  游戏相关的数据结构 还没有建立,
		thData := &mode.T_th_record{}
		thData.Mid = bson.NewObjectId()
		thData.BetAmount = u.TotalBet
		thData.UserId = u.UserId
		thData.DeskId = u.deskId
		thData.WinAmount = u.winAmount - u.TotalBet
		thData.Blance = *(userService.GetUserById(u.UserId).Coin)
		thData.BeginTime = time.Now()
		thData.GameNumber = t.GameNumber
		db.InsertMgoData(casinoConf.DBT_T_TH_RECORD, thData)

		//获取游戏数据
		userRecord := mode.BeanRecord{}
		userRecord.UserId = u.UserId
		userRecord.NickName = u.NickName
		userRecord.WinAmount = u.winAmount - u.TotalBet

		deskRecord.Records = append(deskRecord.Records, userRecord)
		deskRecord.UserIds = strings.Join([]string{deskRecord.UserIds, u.NickName}, ",")

		//保存锦标赛用户的排名信息
		ChampionshipRoom.UpdateUserRankInfo(u.UserId, u.MatchId, u.RoomCoin)
	}

	log.T("开始保存DBT_T_TH_DESK_RECORD的信息")
	//保存桌子的用户信息
	db.InsertMgoData(casinoConf.DBT_T_CS_TH_DESK_RECORD, deskRecord)
	return nil

}

//得到这局胜利的人有几个
func (t *ThDesk) GetWinCount() int {
	t.CalcThcardsWin()        //先计算牌的局面

	var result int = 0
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
			//如果用户不为空,并且状态是等待结算,牌的信息现实的是win 那么,表示一个赢的人
			result ++
		}
	}
	log.T("本局总共有[%v]人是赢牌,", result)
	return result
}


//新用户进入的时候,返回的信息
func (t *ThDesk) GetWeiXinInfos() []*bbproto.WeixinInfo {
	var result []*bbproto.WeixinInfo
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			ru := userService.GetUserById(u.UserId)
			wxi := &bbproto.WeixinInfo{}
			wxi.HeadUrl = ru.HeadUrl
			wxi.NickName = ru.NickName
			wxi.OpenId = ru.OpenId

			//放在列表中
			result = append(result, wxi)
		}
	}
	return result
}

//跟注:跟注的时候 不需要重新设置押注的人
//只是跟注,需要减少用户的资产,增加奖池的金额
func (t *ThDesk) BetUserCall(user  *ThUser) error {
	log.T("用户[%v],nikename[%v],t.BetAmountNow[%v],user.HandCoin[%v]", user.UserId, user.NickName, t.BetAmountNow, user.HandCoin)
	followCoin := t.BetAmountNow - user.HandCoin
	if user.RoomCoin <= followCoin {
		//allin
		t.BetUserAllIn(user.UserId, user.RoomCoin)
	} else {
		//1,增加奖池的金额
		t.calcBetCoin(user.UserId, followCoin)
	}
	return nil
}

//如果弃牌的人是 t.NewRoundBetUser ,需要重新设置值
func (t *ThDesk) NextNewRoundBetUser() error {
	index := t.GetUserIndex(t.NewRoundFirstBetUser)
	for i := index + 1; i < len(t.Users) + index; i++ {
		u := t.Users[i % len(t.Users)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.NewRoundFirstBetUser = u.UserId
			break
		}
		//如果没有找到,那么返回失败
		if i == (len(t.Users) + index - 1) {
			log.T("没有找到下一个默认开始押注的人")
			return errors.New("没有找到下一个默认开始的押注的人")
		}
	}

	return nil

}

//让牌:只有第一个人才可以让牌
func (t *ThDesk) BetUserCheck(userId uint32) error {
	return nil
}

//用户加注
func (t *ThDesk) BetUserRaise(user *ThUser, coin int64) error {
	log.T("开始[%v],nickname[%v]操作用户加注[%v]的操作,user.RoomCoin[%v],user.handcoin[%v].,user.turnCoin[%v].,t.BetAmountNow[%v]", user.UserId, user.NickName, coin, user.RoomCoin, user.HandCoin,user.TurnCoin, t.BetAmountNow)

	//1,加注的时候判断用户金额是否足够
	if (coin > user.RoomCoin) {
		return errors.New("用户加注的金额,比自己手中的金额还要多,出错!")
	}

	//all-in的情况
	if coin == user.RoomCoin {
		//如果加注的金额和手上的余额一样大的话,all in
		if (coin + user.HandCoin) < t.BetAmountNow {
			//表示钱不够,被迫all-in,此时的betamountNow 不会改变
		} else {
			//主动设置下家需要跟注的金额
			t.MinRaise = user.HandCoin + coin - t.BetAmountNow
			t.BetAmountNow = user.HandCoin + coin
			t.RaiseUserId = user.UserId
		}
		//开始allin
		t.BetUserAllIn(user.UserId, user.RoomCoin)

	} else {
		t.MinRaise = user.HandCoin + coin - t.BetAmountNow
		t.BetAmountNow = user.HandCoin + coin
		t.calcBetCoin(user.UserId, coin)                        //thuser
		//3,设置状态:设置为第一个加注的人,如果后边所有人都是跟注,可由这个人判断一轮是否结束
		t.RaiseUserId = user.UserId
	}

	return nil
}

//用户AllIn
func (t *ThDesk) BetUserAllIn(userId uint32, coin int64) error {
	log.T("用户[%v]开始allin[%v]", userId, coin)
	//2,减少用户的金额
	t.calcBetCoin(userId, coin)

	//3,设置用户的状态
	t.GetUserByUserId(userId).Status = TH_USER_STATUS_ALLINING        //设置用户的状态为all-in
	if t.NewRoundFirstBetUser == userId {
		//如果用户是第一个押注的人,all-in之后第一个押注的是往下滑
		t.NextNewRoundBetUser()
		log.T("重新设置了t.NewRoundFirstBetUser[%v]", t.NewRoundFirstBetUser)
	}

	//4,增加all in的状态
	allinpot := &pokerService.AllInJackpot{}
	allinpot.UserId = userId
	allinpot.Jackpopt = 0
	allinpot.ThroundCount = t.RoundCount
	allinpot.AllInAmount = t.GetUserByUserId(userId).TotalBet4calcAllin
	t.AllInJackpot = append(t.AllInJackpot, allinpot)        //增加一个池子
	log.T("用户[%v] all in 的时候,allin的值是[%v]", allinpot.UserId, allinpot.AllInAmount)
	return nil
}

/**
	根据userId 找到在桌子上的index
 */
func (t *ThDesk) GetUserIndex(userId uint32) int {
	var result int = 0
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil && t.Users[i].UserId == userId {
			result = i
			break
		}
	}
	return result
}

//通过UserId找到User
func (t *ThDesk) GetUserByUserId(userId uint32) *ThUser {
	index := t.GetUserIndex(userId)
	return t.Users[index]
}

// 用户加注,跟住,allin 之后对他的各种余额属性进行计算
func (t *ThDesk) calcPreCoin(userId uint32, coin int64) error {
	user := t.GetUserByUserId(userId)

	//这里是前注的信息,handCoin 是否还需要?
	user.AddPreCoin(coin)
	user.AddTotalBet4calcAllin(coin)
	user.AddTotalBet(coin)
	user.AddRoomCoin(-coin)
	user.Update2redis()                //用户信息更新之后,保存到数据库
	t.AddJackpot(coin)                   //底池 增加
	t.AddedgeJackpot(coin)
	return nil
}



// 用户加注,跟住,allin 之后对他的各种余额属性进行计算
func (t *ThDesk) calcBetCoin(userId uint32, coin int64) error {
	user := t.GetUserByUserId(userId)

	user.AddTurnCoin(coin)
	user.AddHandCoin(coin)
	user.AddTotalBet4calcAllin(coin)
	user.AddTotalBet(coin)
	user.AddRoomCoin(-coin)
	user.Update2redis()                //用户信息更新之后,保存到数据库
	t.AddJackpot(coin)                   //底池 增加
	t.AddedgeJackpot(coin)
	return nil
}

/**
	初始化下一个押注的人
	初始化下一个人的时候需要一个超时的处理
 */
func (t *ThDesk) NextBetUser() error {

	log.T("开始计算下一个押注的人是睡,当前押注的人是userId[%v]", t.BetUserNow)
	index := t.GetUserIndex(t.BetUserNow)
	t.BetUserNow = 0        //这里设置为-1是为了方便判断找不到下一个人的时候,设置为新的一局
	for i := index; i < len(t.Users) + index; i++ {
		u := t.Users[(i + 1) % len(t.Users)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			log.T("计算出下一个押注的人,设置betUserNow 为[%v]", u.UserId)
			t.BetUserNow = u.UserId
			break
		}
	}

	log.T("判断是否是下一轮,计算出来的t.BetUserNow[%v],t.BetUserRaiseUserId[%v],t.NextNewRoundBetUser[%v]", t.BetUserNow, t.RaiseUserId, t.NewRoundFirstBetUser)
	//设置新一轮
	if t.BetUserNow == 0 || t.RaiseUserId == t.BetUserNow {
		//处理allin 奖金池分割的问题
		t.CalcAllInJackpot()
		t.RaiseUserId = t.NewRoundFirstBetUser
		t.BetUserNow = t.NewRoundFirstBetUser
		t.BetAmountNow = 0
		t.RoundCount ++
		t.MinRaise = t.BigBlindCoin        //第二轮开始的时候,最低加注金额设置喂大盲注

		////这里吧TurnCoin设置0 是为了是getMinRaise的值正确,目前turnCoin的作用只是在这里
		for i := 0; i < len(t.Users); i++ {
			u := t.Users[i]
			if u != nil {
				u.TurnCoin = 0
			}
		}

		log.T("设置下次押注的人是小盲注,下轮次[%v]", t.RoundCount)
	}

	//如果第一个押注的人,已经弃牌了,那么BetUserRaiseUserId 需要滑向下一个人
	if t.GetUserByUserId(t.RaiseUserId).Status == TH_USER_STATUS_FOLDED {
		log.T("第一个押注的人,弃牌了,需要把t.BetUserRaiseUserId 设置为下一个人")
		for i := t.GetUserIndex(t.RaiseUserId); i < len(t.Users) + index; i++ {
			u := t.Users[(i + 1) % len(t.Users)]
			if u != nil && u.Status == TH_USER_STATUS_BETING {
				log.T("设置betUserNow 为[%v]", u.UserId)
				t.RaiseUserId = u.UserId
				break
			}
		}
	}

	//打印当前桌子的信息
	t.LogString()
	return nil

}

//下一轮
func (t *ThDesk) nextRoundInfo() {
	if !t.isNewRound() {
		return
	}

	log.T("本次设置的押注人和之前的是同一个人,所以开始第[%v]轮的游戏", t.RoundCount)
	//一轮完之后需要发送完成的消息
	sendData := NewGame_SendOverTurn()
	*sendData.Tableid = t.Id
	*sendData.MinRaise = t.GetMinRaise()
	*sendData.NextSeat = t.GetUserByUserId(t.BetUserNow).Seat        //int32(t.GetUserIndex(t.BetUserNow))
	sendData.Handcoin = t.GetHandCoin()
	sendData.Coin = t.GetRoomCoin()
	*sendData.Pool = t.Jackpot
	sendData.SecondPool = t.GetSecondPool()
	t.THBroadcastProto(sendData, 0)

	//清空每个人的handCoin
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			u.HandCoin = 0
			u.TurnCoin = 0
		}
	}

	switch t.RoundCount {
	case TH_DESK_ROUND2:
		//发三张公共牌
		t.sendFlopCard()
	case TH_DESK_ROUND3:
		//发第四章牌
		t.sendTurnCard()
	case TH_DESK_ROUND4:
		//发第五章牌
		t.sendRiverCard()
	}

	//做redis的持久化
	/**
		需要持久化的数据
		1,当前桌子的状态
		2,当鸭玩家的状态
	 */

}




//判断是否是新的一局
func (t *ThDesk) isNewRound() bool {
	//
	log.T("判断是否是新的一轮t.BetUserRaiseUserId[%v],t.BetUserNow(%v),t.status[%v].//status:1,stop,2,start,3,lottery", t.RaiseUserId, t.BetUserNow, t.Status)
	if t.RaiseUserId == t.BetUserNow &&  t.Status == TH_DESK_STATUS_RUN {
		log.T("t.BetUserRaiseUserId[%v] == t.BetUserNow[%v],新的一局开始", t.RaiseUserId, t.BetUserNow)
		return true
	} else {
		return false
	}
}


//发三张公共牌
func (t *ThDesk) sendFlopCard() error {
	log.T("发送三张公共牌begin")
	//申明
	result := &bbproto.Game_SendFlopCard{}
	result.Tableid = &t.Id
	result.Card0 = ThCard2OGCard(t.PublicPai[0])
	result.Card1 = ThCard2OGCard(t.PublicPai[1])
	result.Card2 = ThCard2OGCard(t.PublicPai[2])

	//广播消息
	t.THBroadcastProto(result, 0)
	log.T("发送三张公共牌end")

	return nil
}


//发送第四张牌
func (t *ThDesk) sendTurnCard() error {
	log.T("发送第四张公共牌begin")

	result := &bbproto.Game_SendTurnCard{}
	result.Tableid = &t.Id
	result.Card = ThCard2OGCard(t.PublicPai[3])

	t.THBroadcastProto(result, 0)
	log.T("发送第四张公共牌end")

	return nil
}

//发送低五张牌
func (t *ThDesk) sendRiverCard() error {
	log.T("发送第五张公共牌begin")

	result := &bbproto.Game_SendRiverCard{}
	result.Tableid = &t.Id
	result.Card = ThCard2OGCard(t.PublicPai[4])

	t.THBroadcastProto(result, 0)
	log.T("发送第五张公共牌end")

	return nil
}

//计算奖金池的划分的问题
func (t *ThDesk) CalcAllInJackpot() error {
	log.T("开始计算allin将近池子begin")
	//1,对allin 进行排序,排序之后,可以对奖金池做划分,得到当前未做处理的all和边池的值
	var list pokerService.AllInJackpotList = t.AllInJackpot
	sort.Sort(list)

	for i := 0; i < len(t.AllInJackpot); i++ {
		all := t.AllInJackpot[i]
		log.T("第[%v]次循环的时候,allinlist[%v]", i, all)

		//如果这个池子为nil ,则跳过这个循环,如果这个allIn 不是本轮的,则把之前的allin 的jackpot 累加起来
		if all == nil || all.ThroundCount != t.RoundCount {
			continue
		}

		log.T("开始计算用户[%v]allIn.index[%v] allin.amount[%v]计算all in 时的池子金额", all.UserId, i, all.AllInAmount)
		//每个allin计算金额
		for n := 0; n < len(t.Users); n++ {
			u := t.Users[n]
			if u != nil {
				log.T("用户[%v]押注的总金额是[%v]", u.UserId, u.TotalBet4calcAllin)
				if u.TotalBet4calcAllin > all.AllInAmount {
					all.Jackpopt += all.AllInAmount
					u.TotalBet4calcAllin -= all.AllInAmount
					log.T("用户[%v]押注加入all的金额是[%v]", u.UserId, all.AllInAmount)
				} else {
					all.Jackpopt += u.TotalBet4calcAllin
					u.TotalBet4calcAllin = 0
					log.T("用户[%v]押注加入all的金额是[%v]", u.UserId, u.TotalBet4calcAllin)
				}

			}
		}
		log.T("计算出来用户[%v]allIn.index[%v] allin.amount[%v]计算all in 的池子总金额[]", all.UserId, i, all.AllInAmount, all.Jackpopt)

		//之后的allinamount - 当前allin
		for k := i; k < len(t.AllInJackpot); k++ {
			allk := t.AllInJackpot[k]
			if allk != nil {
				allk.AllInAmount -= all.AllInAmount
			}
		}
		t.edgeJackpot -= all.Jackpopt
		log.T("开始给allIn[%v]计算all in 时的池子金额---------------------------------end---------------", i)
		log.T("目前t.bianJackPot 的剩余值是[%v]", t.edgeJackpot)
	}
	log.T("计算出来的allIn:【%v】", t.AllInJackpot)
	log.T("开始计算allin将近池子end")
	return nil

}

func (t *ThDesk) CheckBetUserBySeat(user *ThUser) bool {
	//2,判断押注的用户是否是当前的用户
	if t.BetUserNow != user.UserId {
		return false
	}

	//3, 判断用户的状态是否正确
	err := user.CheckBetWaitStatus()
	if err != nil {
		return false
	}

	if user.RoomCoin <= 0 {
		log.E("用户userId[%v]name[%v]的带入金额小于0,所以不能押注或者投注了", user.UserId, user.NickName)
		return false
	}

	//用户合法,设置等待状态
	user.InitWait()
	return true
}

//是不是可以开始游戏了
/**
	1,通用的判断
	2,如果是锦标赛,判断锦标赛的逻辑
	3,如果是自定义,判断自定义的逻辑

 */
func (t *ThDesk) IsTime2begin() bool {
	log.T("现在开始判断是否可以开始一局新的游戏:")

	log.T("现在开始判断是否可以开始一局新的游戏,1,判断通用的逻辑:")

	log.T("测试信息,打印每个玩家的状态://1,等待开始,2,坐下,3,已经准备,4,beting5,allin,6,弃牌,7,等待结算,8,已经结算,9,裂开,10,掉线")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("用户[%v].seat[%v]的状态是[%v]", u.UserId, u.Seat, u.Status)
		}
	}

	//1.1,判断桌子当前的状态
	if !t.IsStop() {
		log.T("desk[%v]的状态不是stop[%v]的状态,所以不能开始游戏", t.Id, t.Status)
		return false
	}

	//1.2,判断在线的用户数是否达到标准
	if t.ReadyCount < ThdeskConfig.TH_DESK_LEAST_START_USER {
		log.T("desk[%v]的准备用户数量[%v]不够", t.Id, t.ReadyCount)
		return false
	}

	log.T("现在开始判断是否可以开始一局新的游戏,2,判断锦标赛的逻辑:")
	if t.GameType == intCons.GAME_TYPE_TH_CS {
		//锦标赛游戏房间的状态
		if !ChampionshipRoom.CanNextDeskRun() {
			log.T("锦标赛的逻辑,判断不能开始下一局")
			return false
		}
	}

	log.T("现在开始判断是否可以开始一局新的游戏,1,判断自定义的逻辑:")
	if t.GameType == intCons.GAME_TYPE_TH {
		//自定义房间的标准
	}

	return true

}

//判断桌子是否是stop的状态
func (t *ThDesk) IsStop() bool {
	return t.Status == TH_DESK_STATUS_STOP
}


//开始游戏
func (mydesk *ThDesk) Run() error {
	mydesk.Lock()
	defer mydesk.Unlock()

	log.T("\n\n开始一局新的游戏\n\n")
	//1,判断是否可以开始游戏
	if !mydesk.IsTime2begin() {
		log.T("\n\n不能开始一局新的游戏\n\n")
		return nil
	}

	//2,初始化玩家的信息
	err := mydesk.InitUserBeginStatus()
	if err != nil {
		log.E("开始游戏失败,errMsg[%v]", err.Error())
		return err
	}

	//2,初始化牌的信息
	err = mydesk.OnInitCards()
	if err != nil {
		log.E("开始德州扑克游戏,初始化扑克牌的时候出错")
		return err
	}

	//3,初始化游戏房间的状态
	err = mydesk.OninitThDeskBeginStatus()
	if err != nil {
		log.E("开始德州扑克游戏,初始化房间的状态的时候报错")
		return err
	}

	//3, 初始化前注的信息
	err = mydesk.OinitPreCoin()
	if err != nil {
		log.E("开始德州扑克游戏,初始化房间的状态的时候报错")
		return err
	}

	//4 初始化盲注开始押注
	err = mydesk.InitBlindBet()
	if err != nil {
		log.E("盲注下注的时候出错errMsg[%v]", err.Error())
		return err
	}

	//
	log.T("用户开始等待")
	//本次押注的热开始等待
	waitUser := mydesk.GetUserByUserId(mydesk.BetUserNow)
	log.T("用户开始等待,等待之前桌子的状态[%v]:", mydesk.Status)
	waitUser.wait()

	log.T("广播Game_InitCard的信息")
	initCardB := bbproto.NewGame_InitCard()
	*initCardB.Tableid = int32(mydesk.Id)
	initCardB.HandCard = mydesk.GetHandCard()
	initCardB.PublicCard = mydesk.ThPublicCard2OGC()
	*initCardB.MinRaise = mydesk.GetMinRaise()
	*initCardB.NextUser = mydesk.GetUserByUserId(mydesk.BetUserNow).Seat                //	int32(mydesk.GetUserIndex(mydesk.BetUserNow))
	*initCardB.ActionTime = ThdeskConfig.TH_TIMEOUT_DURATION_INT
	//initCardB.Seat = &mydesk.UserCount
	mydesk.OGTHBroadInitCard(initCardB)
	log.T("广播Game_InitCard的信息完毕")

	log.T("\n\n开始一局新的游戏,初始化完毕\n\n")
	return nil
}

//表示游戏结束
func (t *ThDesk) End() {
	log.T("整局(多场游戏)已经结束...)")
	//广播结算的信息
	result := &bbproto.Game_SendDeskEndLottery{}
	result.Result = &intCons.ACK_RESULT_SUCC

	maxWin := int64(0)
	maxUserid := uint32(0)
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			gel := bbproto.NewGame_EndLottery()
			*gel.Coin = u.RoomCoin - u.InitialRoomCoin        //这里不是u.winamount,u.winamount  表示本局赢得底池的金额
			*gel.Seat = u.Seat

			if t.DeskOwner == u.UserId {
				*gel.Owner = true
			} else {
				*gel.Owner = false
			}

			*gel.Rolename = u.NickName
			*gel.UserId = u.UserId

			if *gel.Coin > maxWin {
				maxWin = *gel.Coin
				maxUserid = u.UserId
			}
			result.CoinInfo = append(result.CoinInfo, gel)
		}
	}

	//赋值大赢家
	for i := 0; i < len(result.CoinInfo); i++ {
		ci := result.CoinInfo[i]
		if ci != nil && ci.GetUserId() == maxUserid {
			*ci.BigWin = true                //设置大赢家
		}
	}

	//设置bigWin
	//广播消息
	t.THBroadcastProtoAll(result)

	//整局游戏结束之后,解散游戏房间
	RmThdesk(t)
}


//手牌转换为OG可以使用的牌
func (t *ThDesk) ThPublicCard2OGC() []*bbproto.Game_CardInfo {
	result := make([]*bbproto.Game_CardInfo, len(t.PublicPai))
	for i := 0; i < len(t.PublicPai); i++ {
		result[i] = ThCard2OGCard(t.PublicPai[i])
	}
	return result
}

//canrase
func (t *ThDesk) GetCanRise() int32 {
	if t.Tiem2Lottery() {
		return 0
	} else {
		return 1
	}
}




//通过座位号来找到user
func (t *ThDesk) getUserBySeat(seatId int32) *ThUser {
	return t.Users[seatId]
}

//押注的通用接口
func (t *ThDesk) DDBet(seatId int32, betType int32, coin int64) error {
	t.Lock()
	defer t.Unlock()
	user := t.getUserBySeat(seatId)
	//1,得到跟注的用户
	if !t.CheckBetUserBySeat(user) {
		log.E("押注人的状态不正确")
		return errors.New("押注人的状态不正确")
	}

	switch betType {
	case TH_DESK_BET_TYPE_CALL:
		t.DDFollowBet(user)

	case TH_DESK_BET_TYPE_FOLD:
		t.DDFoldBet(user)

	case TH_DESK_BET_TYPE_CHECK:        //让牌
		t.DDCheckBet(user)

	case TH_DESK_BET_TYPE_RAISE:        //加注
		t.DDRaiseBet(user, coin)

	case TH_DESK_BET_TYPE_RERRAISE:        //再加注

	case TH_DESK_BET_TYPE_ALLIN:        //全部

	}

	t.nextRoundInfo()        //广播新一局的信息

	if t.Tiem2Lottery() {
		log.T("--------------------------------------现在可以开奖了---------------------------------------------")
		return t.Lottery()
	} else {
		//用户开始等待,如果超时,需要做超时的处理
		t.GetUserByUserId(t.BetUserNow).wait()                //当前押注的人开始等待
		return nil
	}                //开奖
	return nil
}


//这里只处理逻辑
func (t *ThDesk) DDFollowBet(user *ThUser) error {

	//1,押注
	log.T("用户[%v]开始押注", user.UserId)

	//这里需要判断是跟住还是全下
	err := t.BetUserCall(user)

	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//2,初始化下一个押注的人
	t.NextBetUser()

	//3,返回信息
	log.T("返回用户[%v]押注的结果:", user.UserId)
	result := &bbproto.Game_AckFollowBet{}
	result.NextSeat = new(int32)
	result.Coin = new(int64)
	result.Seat = new(int32)
	result.Tableid = new(int32)
	result.CanRaise = new(int32)
	result.MinRaise = new(int64)
	result.Pool = new(int64)
	result.HandCoin = new(int64)

	*result.Coin = user.GetRoomCoin()
	*result.Seat = user.Seat                                //座位id
	*result.Tableid = t.Id
	*result.CanRaise = t.GetCanRise()                               //是否能加注
	*result.MinRaise = t.GetMinRaise()
	*result.Pool = t.Jackpot
	*result.NextSeat = t.GetUserByUserId(t.BetUserNow).Seat //int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	*result.HandCoin = user.HandCoin

	t.THBroadcastProtoAll(result)
	return nil
}

//这里只处理逻辑
func (t *ThDesk) DDFoldBet(user  *ThUser) error {

	log.T("用户[%v]开始弃牌", user.UserId)

	//1,弃牌
	//如果弃牌的人是 t.NewRoundBetUser ,需要重新设置值
	if t.NewRoundFirstBetUser == user.UserId {
		t.NextNewRoundBetUser()
	}
	user.Status = TH_USER_STATUS_FOLDED

	//如果用户是离开的情况,设置用户已经离开
	if user.IsLeave {
		t.LeaveThuser(user.UserId)
	}

	//2,初始化下一个押注的人
	t.NextBetUser()


	//3,返回信息
	log.T("返回用户[%v]弃牌的结果:", user.UserId)
	result := &bbproto.Game_AckFoldBet{}
	result.NextSeat = new(int32)
	result.Pool = new(int64)
	result.MinRaise = new(int64)
	result.HandCoin = new(int64)
	result.CanRaise = new(int32)
	result.Coin = new(int64)

	*result.Pool = t.Jackpot
	*result.HandCoin = user.HandCoin
	*result.Coin = user.GetRoomCoin()                        //本轮压了多少钱

	result.Seat = &user.Seat                                        //座位id
	result.Tableid = &t.Id
	*result.MinRaise = t.GetMinRaise()
	*result.CanRaise = t.GetCanRise()                                //是否能加注
	*result.NextSeat = t.GetUserByUserId(t.BetUserNow).Seat        //int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人

	t.THBroadcastProto(result, 0)
	return nil
}



//联众德州 加注
func (t *ThDesk) DDRaiseBet(user *ThUser, coin int64) error {
	//1,得到跟注的用户,检测用户

	//2,开始处理加注

	//这里需要判断加注是否是all in

	err := t.BetUserRaise(user, coin)
	if err != nil {
		log.E("跟注的时候出错了.errMsg[%v],", err.Error())
	}

	//判断是否属于开奖的时候,如果是,那么开奖,如果不是,设置下一个押注的人
	t.NextBetUser()

	log.T("准备给其他人发送加注的广播")
	//押注成功返回要住成功的消息
	result := &bbproto.Game_AckRaiseBet{}
	result.NextSeat = new(int32)
	result.MinRaise = new(int64)
	result.CanRaise = new(int32)
	result.HandCoin = new(int64)
	result.Coin = new(int64)

	*result.Coin = user.GetRoomCoin()                                //本轮压了多少钱
	result.Seat = &user.Seat                                //座位id
	result.Tableid = &t.Id
	*result.CanRaise = t.GetCanRise()                //是否能加注
	*result.MinRaise = t.GetMinRaise()
	*result.NextSeat = t.GetUserByUserId(t.BetUserNow).Seat        //int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	*result.HandCoin = user.HandCoin                        //表示需要加注多少

	//给所有人广播信息
	t.THBroadcastProto(result, 0)
	log.T("开始处理seat[%v]加注的逻辑,t,OGRaiseBet()...end", user.Seat)
	return nil
}


//得到当前用户需要加注的金额
func (t *ThDesk) GetMinRaise() int64 {
	log.T("获取用户[%v]的最低加注金额,handCoin[%v],t.TurnCoin[%v],t.MinRaise[%v],t.BetAmountNow[%v],",
		t.BetUserNow, t.GetUserByUserId(t.BetUserNow).HandCoin, t.GetUserByUserId(t.BetUserNow).TurnCoin, t.MinRaise, t.BetAmountNow)
	result := t.MinRaise + t.BetAmountNow - t.GetUserByUserId(t.BetUserNow).TurnCoin
	if result < 0 {
		//现在处理的有可能是新的一局开始
		result = t.BigBlindCoin
	}
	log.T("获取用户[%v]的最低加注金额是[%v]", t.BetUserNow,result)
	return result
}


//联众德州 让牌
func (t *ThDesk) DDCheckBet(user *ThUser) error {

	log.T("用户[%v]开始让牌", user.UserId)
	//1,让牌
	err := t.BetUserCheck(user.UserId)
	if err != nil {
		log.E("用户[%v]让牌的时候出错了.errMsg[%v],", user.UserId, err.Error())
	}


	//2,计算洗一个押注的人
	t.NextBetUser()

	//3押注成功返回要住成功的消息
	log.T("打印user[%v]让牌的结果:", user.UserId)
	result := &bbproto.Game_AckCheckBet{}
	result.NextSeat = new(int32)
	result.MinRaise = new(int64)
	result.CanRaise = new(int32)
	result.HandCoin = new(int64)
	result.Coin = new(int64)

	*result.Coin = user.GetRoomCoin()                                //本轮压了多少钱
	result.Seat = &user.Seat                                        //座位id
	result.Tableid = &t.Id
	*result.CanRaise = t.GetCanRise()                                //是否能加注
	*result.MinRaise = t.GetMinRaise()                                //最低加注金额
	*result.NextSeat = t.GetUserByUserId(t.BetUserNow).Seat        //int32(t.GetUserIndex(t.BetUserNow))		//下一个押注的人
	*result.HandCoin = user.HandCoin

	t.THBroadcastProtoAll(result)

	return nil
}





func (t *ThDesk) AddUserCountOnline() {
	atomic.AddInt32(&t.UserCountOnline, 1)
}

func (t *ThDesk) SubUserCountOnline() {
	atomic.AddInt32(&t.UserCountOnline, -1)
}

func (t *ThDesk) AddReadyCount() {
	atomic.AddInt32(&t.ReadyCount, 1)
}

func (t *ThDesk) SubReadyCount() {
	atomic.AddInt32(&t.ReadyCount, -1)
}

func (t *ThDesk) AddUserCount() {
	atomic.AddInt32(&t.UserCount, 1)
}

func (t *ThDesk) SubUserCount() {
	atomic.AddInt32(&t.UserCount, -1)
}

func (t *ThDesk) AddJackpot(coin int64) {
	atomic.AddInt64(&t.Jackpot, coin)
}

func (t *ThDesk) AddedgeJackpot(coin int64) {
	atomic.AddInt64(&t.edgeJackpot, coin)
}


