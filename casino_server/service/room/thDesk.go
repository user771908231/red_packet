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
var TH_DESK_STATUS_LOTTERY int32 = 4             //正在开奖
var TH_DESK_STATUS_GAMEOVER int32 = 5             //已经开始的状态


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

	TH_DONGHUA_DURATION      time.Duration //动画的延时

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
	ThdeskConfig.TH_DONGHUA_DURATION = time.Second * 1        //发牌之后的动画,延时2秒
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
	blindLevel           int32                        //当前盲注的级别
	BeginTime            time.Time                    //游戏开始时间
	EndTime              time.Time                    //游戏结束时间
	RebuyCountLimit      int32                        //重购的次数
	RebuyBlindLevelLimit int32                        //rebuy盲注的限制

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
	Status               int32                        //牌桌的状态
	BetAmountNow         int64                        //当前的押注金额是多少
	RoundCount           int32                        //第几轮
	Jackpot              int64                        //奖金池
	EdgeJackpot          int64                        //边池
	MinRaise             int64                        //最低加注金额
	AllInJackpot         []*pokerService.AllInJackpot //allin的标记

	SendFlop             bool                         //公共底牌
	SendTurn             bool                         //第四章牌
	SendRive             bool                         //第五章牌
}

/**
	新生成一个德州的桌子
 */
func NewThDesk() *ThDesk {
	result := new(ThDesk)
	result.Id = newThDeskId()
	result.UserCount = 0
	result.Dealer = 0                //不需要创建  默认就是为空
	result.BetUserNow = 0
	result.BigBlind = 0
	result.SmallBlind = 0
	result.Users = make([]*ThUser, ThdeskConfig.TH_DESK_MAX_START_USER)
	result.RaiseUserId = 0
	result.RoundCount = 0
	result.NewRoundFirstBetUser = 0
	result.EdgeJackpot = 0
	result.SmallBlindCoin = ThGameRoomIns.SmallBlindCoin
	result.BigBlindCoin = 2 * ThGameRoomIns.SmallBlindCoin
	result.Status = TH_DESK_STATUS_STOP        //游戏还没有开始的状态
	result.GameType = intCons.GAME_TYPE_TH_CS                          //游戏桌子的类型
	result.JuCount = 0
	result.JuCountNow = 1                //
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
	//log.T("当前desk[%v]的信息:-----------------------------------begin----------------------------------", t.Id)
	//log.T("当前desk[%v]的信息的状态status[%v]", t.Id, t.Status)
	//for i := 0; i < len(t.Users); i++ {
	//	u := t.Users[i]
	//	if u != nil {
	//		//log.T("当前desk[%v]的user[%v]的状态status[%v],牌的信息[%v]", t.Id, u.UserId, u.Status, u.HandCards)
	//		log.T("当前desk[%v]的user[%v]的状态status[%v],HandCoin[%v],TurnCoin[%v],RoomCoin[%v]", t.Id, u.UserId, u.Status, u.HandCoin, u.TurnCoin, u.RoomCoin)
	//	}
	//}
	//log.T("当前desk[%v]的信息的状态,MatchId[%v]", t.Id, t.MatchId)
	//log.T("当前desk[%v]的信息的状态,RoomKey[%v]", t.Id, t.RoomKey)
	//log.T("当前desk[%v]的信息的状态,CreateFee[%v]", t.Id, t.CreateFee)
	//log.T("当前desk[%v]的信息的状态,GameType[%v]", t.Id, t.GameType)
	//log.T("当前desk[%v]的信息的状态,InitRoomCoin[%v]", t.Id, t.InitRoomCoin)
	//log.T("当前desk[%v]的信息的状态,JuCount[%v]", t.Id, t.JuCount)
	//log.T("当前desk[%v]的信息的状态,JuCountNow[%v]", t.Id, t.JuCountNow)
	//log.T("当前desk[%v]的信息的状态,PreCoin[%v]", t.Id, t.PreCoin)
	//log.T("当前desk[%v]的信息的状态,SmallBlindCoin[%v]", t.Id, t.SmallBlindCoin)
	//log.T("当前desk[%v]的信息的状态,BigBlindCoin[%v]", t.Id, t.BigBlindCoin)
	//log.T("当前desk[%v]的信息的状态,小盲注[%v]", t.Id, t.SmallBlind)
	//log.T("当前desk[%v]的信息的状态,大盲注[%v]", t.Id, t.BigBlind)
	//log.T("当前desk[%v]的信息的状态,RaiseUserId[%v]", t.Id, t.RaiseUserId)
	//log.T("当前desk[%v]的信息的状态,MinRaise[%v]", t.Id, t.MinRaise)
	//log.T("当前desk[%v]的信息的状态,BetUserNow[%v]", t.Id, t.BetUserNow)
	//log.T("当前desk[%v]的信息的状态,GameNumber[%v]", t.Id, t.GameNumber)
	//log.T("当前desk[%v]的信息的状态,压注人[%v]", t.Id, t.BetUserNow)
	//log.T("当前desk[%v]的信息的状态,压注轮次[%v]", t.Id, t.RoundCount)
	//log.T("当前desk[%v]的信息的状态,NewRoundFirstBetUser[%v]", t.Id, t.NewRoundFirstBetUser)
	//log.T("当前desk[%v]的信息的状态,总共押注Jackpot[%v]", t.Id, t.Jackpot)
	//log.T("当前desk[%v]的信息的状态,edgeJackpot[%v]", t.Id, t.EdgeJackpot)
	//log.T("当前desk[%v]的信息的状态,当前加注的人BetUserRaiseUserId[%v]", t.Id, t.RaiseUserId)
	//log.T("当前desk[%v]的信息:-----------------------------------end----------------------------------", t.Id)
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
			u.Agent = a                                                //设置用户的连接
			u.IsBreak = false               //设置用户的离线状态
			u.IsLeave = false
			u.UpdateAgentUserData()         //更新回话信息
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
	thUser.Agent = a
	thUser.Status = userStatus        //刚进房间的玩家
	thUser.deskId = t.Id                //桌子的id
	thUser.NickName = *redisUser.NickName                //todo 测试阶段,把nickName显示成用户id
	thUser.RoomCoin = t.InitRoomCoin
	thUser.IsBreak = false
	thUser.IsLeave = false
	thUser.LotteryCheck = true

	//根据桌子的状态 设置用户的游戏状态
	if t.IsChampionship() {
		thUser.GameStatus = TH_USER_GAME_STATUS_CHAMPIONSHIP        //新加入游戏的时候设置用户的游戏状态为锦标赛
		thUser.CSGamingStatus = true
	} else if t.IsFriend() {
		thUser.GameStatus = TH_USER_GAME_STATUS_FRIEND
	} else {
		thUser.GameStatus = TH_USER_GAME_STATUS_NOGAME
	}

	//3,添加thuser
	err := t.AddThuserBean(thUser)
	if err != nil {
		log.E("增加user【%v】到desk【%v】失败", thUser.UserId, t.Id)
		return nil, errors.New("增加user失败")
	}

	//4, 把用户的信息绑定到agent上
	thUser.UpdateAgentUserData()

	//5,等待的用户加1
	t.AddUserCount()
	return thUser, nil
}


//增加一个user实体
func (t *ThDesk) AddThuserBean(user *ThUser) error {
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
	if user == nil {
		return errors.New("没有找到对应的用户...")
	}

	//1,如果是在游戏的过程中,则准备失败
	if t.Status == TH_DESK_STATUS_RUN {
		log.E("desk[%v]已经在游戏中了,user[%v]不能准备", t.Id, userId)
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_GAME_READY_REPEAT), "已经在游戏中,不能准备")
	}

	//1,如果用户已经准备,则返回重复准备
	if user.IsReady() {
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_GAME_READY_REPEAT), "已经准备了")
	}

	//2,如果用户余额不足,则准备失败
	if !t.IsUserRoomCoinEnough(user) {
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_GAME_READY_CHIP_NOT_ENOUGH), "筹码不足")
	}


	//3,准备成功
	user.Status = TH_USER_STATUS_READY
	return nil
}

//查看是否全部准备好了
func (t *ThDesk) IsAllReady() bool {
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			//用户既没有掉线,也没有离开的情况下,并且钱是充足的情况下,如果没有准备,那么返回false
			if !u.IsReady() && t.IsUserRoomCoinEnough(u) && !u.IsLeave && !u.IsBreak {
				return false
			}
		}
	}

	return true
}

//  用户退出德州游戏的房间,rmUser 需要在事物中进行,离开德州需要更具桌子的类型来做不同的处理
func (t *ThDesk) LeaveThuser(userId uint32) error {
	if t.IsFriend() {
		return t.FLeaveThuser(userId)
	} else if t.IsChampionship() {
		return t.CSLeaveThuser(userId)
	}

	return nil
}

func (t *ThDesk) FLeaveThuser(userId uint32) error {

	//1,离开之后,设置用户的信息
	user := t.GetUserByUserId(userId)
	user.IsLeave = true     //设置状态为离开
	user.GameStatus = TH_USER_GAME_STATUS_NOGAME        //用户离开之后,设置用户的游戏状态为没有游戏中
	user.UpdateAgentUserData()

	//2,自定义房间,如果其他人都准备了,那么开始游戏,离开房间和准备的处理是一样的
	//这样处理的作用是,防止最后一个未准备的人,离开房间以后游戏不能开始
	if t.JuCountNow > 1 && t.IsAllReady() {
		//用户离开之后,判断游戏是否开始
		go t.Run()
	}

	//3,返回离开房间之后的信息
	ret := bbproto.NewGame_ACKLeaveDesk()
	*ret.Result = intCons.ACK_RESULT_SUCC
	user.WriteMsg(ret)

	//离开之后,需要广播一次sendGameInfo,这里人还没有真正的离开
	t.BroadGameInfo(userId)

	//保存数据到redis
	t.UpdateThdeskAndUser2redis(user)

	return nil
}

func (t *ThDesk) CSLeaveThuser(userId uint32) error {

	//1,离开之后,设置用户的信息
	user := t.GetUserByUserId(userId)
	user.IsLeave = true     //设置状态为离开
	user.GameStatus = TH_USER_GAME_STATUS_NOGAME        //用户离开之后,设置用户的游戏状态为没有游戏中
	user.CSGamingStatus = false
	user.RoomCoin = 0
	user.UpdateAgentUserData()

	//2,更新锦标赛的数据
	ChampionshipRoom.UpdateUserRankInfo(user.UserId, t.MatchId, user.RoomCoin)
	ChampionshipRoom.SubOnlineCount()        //竞标赛的在线人数-1

	//3,返回离开房间之后的信息
	ret := bbproto.NewGame_ACKLeaveDesk()
	*ret.Result = intCons.ACK_RESULT_SUCC
	user.WriteMsg(ret)

	//4,删除用户
	t.RmUser(user.UserId)                         //删除用户,并且发送广播
	//离开之后,需要广播一次sendGameInfo,这里人还没有真正的离开
	t.BroadGameInfo(userId)

	//保存数据到redis
	t.UpdateThdeskAndUser2redis(user)

	return nil
}

func (r *ThDesk) RmUser(userId uint32) {
	for _, user := range r.Users {
		if user != nil && user.UserId == userId {
			//更新回话信息
			user.IsLeave = true     //设置状态为离开
			user.GameStatus = TH_USER_GAME_STATUS_NOGAME        //用户离开之后,设置用户的游戏状态为没有游戏中
			user.CSGamingStatus = false
			user.RoomCoin = 0
			user.deskId = 0
			user.IsLeave = true
			user.IsBreak = true
			user.UpdateAgentUserData()

			//设置为nil
			user = nil
		}
	}


	//删除redis中的数据
	DelRedisThUser(r.Id, r.GameNumber, userId)
}

//设置用户为掉线的状态
func (t *ThDesk) SetOfflineStatus(userId uint32) error {

	u := t.GetUserByUserId(userId)
	//1,设置状态为断线
	//这里需要保存掉线钱的状态
	u.IsBreak = true        //设置为掉线的状态
	u.UpdateAgentUserData()
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
	*blindB.Banker = t.GetUserSeatByUserId(t.Dealer)        //庄
	blindB.Bigblind = &t.BigBlindCoin        //大盲注
	blindB.Smallblind = &t.SmallBlindCoin        //小盲注
	*blindB.Bigblindseat = t.GetUserSeatByUserId(t.BigBlind)        //      //大盲注座位号
	*blindB.Smallblindseat = t.GetUserSeatByUserId(t.SmallBlind)        //     //小盲注座位号
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
func (t *ThDesk) InitUserStatus() error {
	log.T("开始一局新的游戏,开始初始化用户的状态")

	//清空状态leave 的玩家
	t.RmLeaveUser()

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
		u.LotteryCheck = true                                //游戏开始的时候设置为false

		//如果用户的余额不足或者用户的状态是属于断线的状态,则设置用户为等待入座
		if !t.IsUserRoomCoinEnough(u) {
			log.T("由于用户[%v] status[%v],的roomCoin[%v] <= desk.BigBlindCoin 所以设置用户为TH_USER_STATUS_WAITSEAT", u.UserId, u.IsBreak, u.RoomCoin, t.BigBlindCoin)
			u.Status = TH_USER_STATUS_WAITSEAT        //只是坐下,没有游戏中
			u.LotteryCheck = true                          //计算是否开奖的时候需要用到
			continue
		}

		//用户不是离线的状态,并且,用户已经准备好了,则可以开始游戏
		if !u.IsBreak && !u.IsLeave && u.IsReady() {
			log.T("由于用户[%v]的status[%v]BreakStatus[%v],所以设置状态为TH_USER_STATUS_BETING", u.UserId, u.Status, u.IsBreak)
			u.Status = TH_USER_STATUS_BETING
			u.LotteryCheck = false
		}
	}

	//------------------------------------由于联众前端设计的问题...这里的user需要重新排列user的顺序------------------------------------

	log.T("开始一局新的游戏,初始化用户的状态完毕")
	return nil
}

//清除离开房间的人
func (t *ThDesk) RmLeaveUser() {
	//清空状态为离开的玩家
	for _, user := range t.Users {
		if user != nil {
			if user.IsLeave {
				t.RmUser(user.UserId)
				continue
			} else if user.IsBreak && user.GetDesk().IsChampionship() {
				t.RmUser(user.UserId)
			}
		}
	}

	//重新对玩家进行赋值
	usersTemp := make([]*ThUser, len(t.Users))
	copy(usersTemp, t.Users)
	log.T("原来的thsuers:[%v]", t.Users)

	//初始化之前的用户为nil
	for i := 0; i < len(t.Users); i++ {
		t.Users[i] = nil
	}

	//排序游戏中的玩家
	for i := 0; i < len(usersTemp); i++ {
		u := usersTemp[i]
		if u != nil {
			t.AddThuserBean(u)
		}
	}

	//重新设置座位号
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			u.Seat = int32(i)
		}
	}

	//发送gameInfo的信息
	t.BroadGameInfo(0)
	time.Sleep(time.Second * 1)        //1秒时候发送其他的信息
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
			//log.T("用户[%v]的手牌[%v]", t.Users[i].UserId, t.Users[i].HandCards)
			t.Users[i].thCards = pokerService.GetTHPoker(t.Users[i].HandCards, t.PublicPai, 5)
			//log.T("用户[%v]的:拍类型,所有牌[%v],th[%v]", t.Users[i].UserId, t.Users[i].thCards.ThType, t.Users[i].thCards.Cards, t.Users[i].thCards)
		}
	}

	log.T("用户开始等待")
	//本次押注的热开始等待
	waitUser := t.GetUserByUserId(t.BetUserNow)
	log.T("用户开始等待,等待之前桌子的状态[%v]:", t.Status)
	waitUser.wait()

	log.T("广播Game_InitCard的信息")
	initCardB := bbproto.NewGame_InitCard()
	*initCardB.Tableid = t.Id
	initCardB.HandCard = t.GetHandCard()
	initCardB.PublicCard = t.ThPublicCard2OGC()
	*initCardB.MinRaise = t.GetMinRaise()
	*initCardB.NextUser = t.GetUserSeatByUserId(t.BetUserNow)
	*initCardB.ActionTime = ThdeskConfig.TH_TIMEOUT_DURATION_INT
	*initCardB.CurrPlayCount = t.JuCountNow
	*initCardB.TotalPlayCount = t.JuCount

	t.OGTHBroadInitCard(initCardB)
	log.T("开始一局新的游戏,初始化牌的信息完毕...")

	return nil
}

//广播porto消息的通用方法
func (t *ThDesk) THBroadcastProto(p proto.Message, ignoreUserId uint32) error {
	log.Normal("开始广播proto消息【%v】", p)
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]                //给这个玩家发送广播信息
		if u != nil && u.UserId != ignoreUserId && u.IsLeave == false && u.IsBreak == false {
			u.WriteMsg(p)
		}
	}
	return nil
}


//发送本局牌局的结果
func (t *ThDesk) BroadcastTestResult(p *bbproto.Game_TestResult) error {
	//发送的时候,初始化自己的排名
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]                //给这个玩家发送广播信息
		if u != nil && !u.IsLeave  && !u.IsBreak && u.IsClose() {
			*p.Rank = t.getRankByUserId(u)        //获取用户的排名
			*p.CanRebuy = t.getCanRebuyByUserId(u)        //是否可以重构
			*p.RebuyCount = u.RebuyCount        //重购的次数

			//判断是否可以
			t.Users[i].WriteMsg(p)
		}
	}
	return nil
}

//获取用户的排名
func (t *ThDesk) getRankByUserId(user *ThUser) int32 {
	if t.IsChampionship() {
		return ChampionshipRoom.GetRankByuserId(user.UserId)
	} else {
		return 0
	}
}

func (t *ThDesk) getRankUserCount() int32 {
	return ChampionshipRoom.onlineCount
}

//得到这个用户是可以重购,不同的桌子,来判断的逻辑不用
func (t *ThDesk) getCanRebuyByUserId(user *ThUser) bool {
	if t.IsFriend() {
		log.T("判断user【%v】是否能重购买,t.juCountNow[%v],t.juCount[%v]", user.UserId, t.JuCountNow, t.JuCount)
		if !t.IsUserRoomCoinEnough(user) && t.JuCountNow < t.JuCount {
			return true
		} else {
			return false
		}
	} else if t.IsChampionship() {
		//用户的金额不足的时候并且重构的次数小于desk的重购限制的时候
		if !t.IsUserRoomCoinEnough(user) &&
		user.RebuyCount < t.RebuyCountLimit &&
		t.blindLevel < t.RebuyBlindLevelLimit {
			return true
		} else {
			return false
		}
	}

	return false

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
func (t *ThDesk) OninitThDeskStatus() error {
	log.T("开始一局游戏,现在初始化desk的信息")

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
			if u != nil && u.IsBetting() {
				t.Dealer = u.UserId
			}
		}
	}

	//设置小盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil && u.IsBetting() {
			t.SmallBlind = u.UserId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	//设置大盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil && u.IsBetting() {
			t.BigBlind = u.UserId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	/**
		设置第一个押注的人
	 */

	if t.GetBettingUserCount() == int32(2) {
		//如果只有两个人,当前押注的人是小盲注
		log.T("由于当前游戏中的t.userCountOnline==2,所以设置betUserNow 是盲注", )
		t.BetUserNow = t.SmallBlind
	} else {
		//设置当前押注的人
		for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
			u := userTemp[(i + 1) % len(userTemp)]
			if u != nil && u.IsBetting() {
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
	t.EdgeJackpot = 0
	t.AllInJackpot = nil                          // 初始化allInJackpot 为空
	t.SendFlop = false        //是否已经发了三张底牌
	t.SendTurn = false        //是否已经发了第四张牌
	t.SendRive = false        //是否已经发了第五张牌
	t.GameNumber, _ = db.GetNextSeq(casinoConf.DBT_T_TH_GAMENUMBER_SEQ)
	t.Status = TH_DESK_STATUS_RUN                //设置德州desk状态//设置状态为开始游戏

	//如果是锦标赛,需要设置锦标赛的属性
	if t.IsChampionship() {
		t.MatchId = ChampionshipRoom.MatchId
	}

	t.LogString()
	log.T("开始一局游戏,现在初始化desk的信息完毕...")
	return nil
}

func (t *ThDesk) GetBettingUserCount() int32 {
	var count int32 = 0
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.IsBetting() {
			count ++
		}
	}

	return count
}

//判断是不是朋友桌
func (t *ThDesk) IsFriend() bool {
	return t.GameType == intCons.GAME_TYPE_TH
}

//判断是否是锦标赛
func (t *ThDesk) IsChampionship() bool {
	return t.GameType == intCons.GAME_TYPE_TH_CS
}


//判断lotteryCheck == false 的count
func (t *ThDesk) GetLotteryCheckFalseCount() int32 {
	var count int32 = 0
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && !u.LotteryCheck {
			count ++

		}
	}
	return count
}


//
func (t *ThDesk) GetBettingCount() int32 {
	var count int32 = 0
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.IsBetting() {
			count ++
		}
	}
	return count
}

func (t *ThDesk) GetBettingOrAllinCount() int32 {
	var count int32 = 0
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && (u.IsBetting() || u.IsAllIn()) {
			count ++
		}
	}
	return count
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

	log.T("判断是否应该开奖,打印每个人的信息://1,刚上桌子,2,坐下,3,ready 4 押注中,5,allin,6,弃牌,7,等待结算,8,已经结算,9,离开")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("[%v]判断是否应该开奖,打印user[%v]的状态[%v]:", i, u.UserId, u.Status)
		}
	}


	//
	lotteryCheckFalseCount := t.GetLotteryCheckFalseCount()
	log.T("t.getlotteryCheckFalseCount[%v]", lotteryCheckFalseCount)

	if lotteryCheckFalseCount <= 1 {
		log.T("因为getlotteryCheckFalseCount == 1  ,所以现在开始开奖...")
		return true
	}

	//
	//var betingCount int = 0
	//for i := 0; i < len(t.Users); i++ {
	//	if t.Users[i] != nil && t.Users[i].Status == TH_USER_STATUS_BETING {
	//		betingCount ++
	//	}
	//}
	//
	//log.T("当前处于押注中的人数是[%v]", betingCount)
	////如果押注的人只有一个人了,那么是开奖的时刻
	//if betingCount <= 1 {
	//	log.T("现在处于押注中(beting)状态的人,只剩下一个了,所以直接开奖")
	//	return true
	//}

	//第四轮,并且计算出来的押注人和start是同一个人
	if t.RoundCount == TH_DESK_ROUND_END {
		log.T("现在处于第[%v]轮押注,所以可以直接开奖", t.RoundCount)
		return true
	}

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
func (t *ThDesk) InitLotteryStatus() error {

	//1,设置desk当前的状态为开奖中
	t.Status = TH_DESK_STATUS_LOTTERY

	log.T("开奖之前答应每个人的状态,并且修改为等待结算")
	//2,设置用户的状态为等待开奖
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("用户[%v].nickname[%v]的status[%v]", u.UserId, u.NickName, u.Status)
			u.FinishtWait()        //不再等待
			if u.IsAllIn() || u.IsBetting() {
				//如果用户当前的状态是押注中,或者all in,那么设置用户的状态喂等待结算
				u.Status = TH_USER_STATUS_WAIT_CLOSED
			} else if u.IsFold() {
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
		var bbonus int64 = t.EdgeJackpot / int64(bwinCount)
		for i := 0; i < len(t.Users); i++ {
			u := t.Users[i]

			if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
				//
				//对这个用户做结算...
				log.T("现在开始开奖,计算边池的奖励,user[%v]得到[%v]....", u.UserId, bbonus)
				u.AddWinAmount(bbonus)
				u.AddRoomCoin(bbonus)
				u.winAmountDetail = append(u.winAmountDetail, bbonus)        //详细的奖励(边池主池分开)
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

	//发送还没有发送的牌
	t.sendReaminCard()

	//设置用户的状态都为的等待开奖
	t.InitLotteryStatus()

	//需要计算本局allin的奖金池
	t.CalcAllInJackpot()

	//计算用户输赢情况
	t.calcUserWinAmount()

	//保存数据到数据库
	t.SaveLotteryData()

	//广播开奖的额结果
	t.broadLotteryResult()

	//开奖之后,设置状态为 没有开始游戏
	t.afterLottery()

	//判断游戏是否结束
	if !t.end() {
		//表示不能继续开始游戏
		go t.Run()
	}

	//备份数据到redis
	t.UpdateThdeskAndAllUser2redis()
	return nil
}


//判断是否可以开始下一句游戏
func (t *ThDesk) end() bool {
	if t.IsFriend() {
		//朋友桌是否结束游戏
		return t.EndTh()
	} else if t.IsChampionship() {
		return t.EndCsTh()        //锦标赛进入游戏
	} else {
		return true
	}
}


//广播开奖的结果
func (t *ThDesk) broadLotteryResult() error {
	//发送是否需要加注的时候,需要升盲之后才能确定
	if t.IsChampionship() {
		t.SmallBlindCoin = ChampionshipRoom.SmallBlindCoin
		t.BigBlindCoin = ChampionshipRoom.SmallBlindCoin * 2
	}

	//1.发送输赢结果
	result := bbproto.NewGame_TestResult()
	*result.Tableid = t.Id                          //桌子
	result.BCanShowCard = t.GetBshowCard()          //
	result.BShowCard = t.GetBshowCard()             //亮牌
	result.Handcard = t.GetHandCard()               //手牌
	result.WinCoinInfo = t.getWinCoinInfo()
	result.HandCoin = t.GetRoomCoin()                //现实用户的余额
	result.CoinInfo = t.getCoinInfo()               //每个人的输赢情况
	*result.RankUserCount = t.getRankUserCount()
	t.BroadcastTestResult(result)
	return nil
}

//开奖之后的处理
func (t *ThDesk) afterLottery() error {
	//1,设置游戏桌子的状态
	log.T("开奖结束,设置desk的状态为stop")
	t.Status = TH_DESK_STATUS_STOP                //设置为没有开始开始游戏
	t.Jackpot = 0; //主池设置为0
	t.EdgeJackpot = 0; //边池设置为0
	t.AllInJackpot = nil;
	t.blindLevel = 0;
	t.JuCountNow ++

	//2,设置用户的状态
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && !u.IsBreak {
			if t.IsFriend() {
				//如果是自定义的房间,设置每个人都是坐下的状态
				u.Status = TH_USER_STATUS_SEATED
			} else if t.IsChampionship() {
				//这里不应该这么判断,因为升盲之后,钱已经变了...
				if t.IsUserRoomCoinEnough(u) {
					u.Status = TH_USER_STATUS_READY   //如果是锦标赛的房间,用户的钱足够
				} else {
					//如果是锦标赛的房间,用户的钱不够
					u.Status = TH_USER_STATUS_SEATED
					u.waitCsRebuy()//等待重购
				}
			}
		}
	}

	log.T("lottery 之后保存用户和desk的数据到redis")
	t.UpdateThdeskAndAllUser2redis()

	return nil
}

//判断用户的余额是否足够开始下一场游戏
func (t *ThDesk) IsUserRoomCoinEnough(u *ThUser) bool {
	if u.RoomCoin < t.PreCoin + t.BigBlindCoin {
		log.T("用户[%v]的roomCoin[%v] < preCoin[%v]+bigBlindCoin[%v]", u.UserId, u.RoomCoin, t.PreCoin, t.BigBlindCoin)
		return false
	} else {
		return true
	}
}


//发送没有发送玩的牌
func (t *ThDesk) sendReaminCard() error {
	log.T("判断是否需要发送剩余的牌GetBettingOrAllinCount[%v]", t.GetBettingOrAllinCount())
	if t.GetBettingOrAllinCount() > 1 {
		//大于一个人则需要发送剩余的牌
		if !t.SendFlop {
			t.sendFlopCard()
		}

		if !t.SendTurn {
			t.sendTurnCard()
		}

		if !t.SendRive {
			t.sendRiverCard()
		}
	}
	return nil
}



//保存数据到数据库
//这里需要根据游戏类型的不同来分别存醋

func (t *ThDesk)  SaveLotteryData() error {

	if t.GameType == intCons.GAME_TYPE_TH {
		//自定义房间
		return t.SaveLotteryDatath()
	} else if t.IsChampionship() {
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
		ChampionshipRoom.UpdateUserRankInfo(u.UserId, t.MatchId, u.RoomCoin)
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
	log.T("开始[%v],nickname[%v]操作用户加注[%v]的操作,user.RoomCoin[%v],user.handcoin[%v].,user.turnCoin[%v].,t.BetAmountNow[%v]", user.UserId, user.NickName, coin, user.RoomCoin, user.HandCoin, user.TurnCoin, t.BetAmountNow)

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

//通过UserId找到User//默认返回第一个人
func (t *ThDesk) GetUserByUserId(userId uint32) *ThUser {
	for _, u := range t.Users {
		if u != nil && u.UserId == userId {
			return u
		}
	}
	return nil
}

//返回用户的seat ,如果找不到userId ,则返回-1
func (t *ThDesk) GetUserSeatByUserId(userId uint32) int32 {
	u := t.GetUserByUserId(userId)
	if u == nil {
		return -1
	} else {
		return u.Seat
	}
}

// 用户加注,跟住,allin 之后对他的各种余额属性进行计算
func (t *ThDesk) calcPreCoin(userId uint32, coin int64) error {
	user := t.GetUserByUserId(userId)

	//这里是前注的信息,handCoin 是否还需要?
	user.AddPreCoin(coin)
	user.AddTotalBet4calcAllin(coin)
	user.AddTotalBet(coin)
	user.AddRoomCoin(-coin)
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
	t.AddJackpot(coin)                   //底池 增加,
	t.AddedgeJackpot(coin)
	return nil
}

/**
	初始化下一个押注的人
	初始化下一个人的时候需要一个超时的处理
 */
func (t *ThDesk) NextBetUser() error {

	log.T("开始计算下一个押注的人是谁,当前押注的人是userId[%v]", t.BetUserNow)
	index := t.GetUserIndex(t.BetUserNow)
	t.BetUserNow = 0        //这里设置为-1是为了方便判断找不到下一个人的时候,设置为新的一局
	for i := index; i < len(t.Users) + index; i++ {
		u := t.Users[(i + 1) % len(t.Users)]
		if u != nil && (u.IsBetting() || (u.IsAllIn() && t.RaiseUserId == u.UserId)) {
			log.T("计算出下一个押注的人,设置betUserNow 为[%v]", u.UserId)
			t.BetUserNow = u.UserId
			break
		}
	}

	log.T("判断是否是下一轮,计算出来的t.BetUserNow[%v],t.RaiseUserId[%v],t.NextNewRoundBetUser[%v]", t.BetUserNow, t.RaiseUserId, t.NewRoundFirstBetUser)
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
	*sendData.NextSeat = t.GetUserSeatByUserId(t.BetUserNow)        //
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
			//如果用户是allin的状态,那么需要设置
			if u.IsAllIn() {
				u.LotteryCheck = true
			}
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
	log.T("判断是否是新的一轮t.BetUserRaiseUserId[%v],t.BetUserNow(%v),t.status[%v].", t.RaiseUserId, t.BetUserNow, t.Status)
	if t.RaiseUserId == t.BetUserNow &&  t.IsRun() {
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
	t.SendFlop = true

	//广播消息
	t.THBroadcastProto(result, 0)
	time.Sleep(ThdeskConfig.TH_DONGHUA_DURATION)
	log.T("发送三张公共牌end")

	return nil
}


//发送第四张牌
func (t *ThDesk) sendTurnCard() error {
	log.T("发送第四张公共牌begin")

	result := &bbproto.Game_SendTurnCard{}
	result.Tableid = &t.Id
	result.Card = ThCard2OGCard(t.PublicPai[3])
	t.SendTurn = true

	t.THBroadcastProto(result, 0)
	time.Sleep(ThdeskConfig.TH_DONGHUA_DURATION)

	log.T("发送第四张公共牌end")

	return nil
}

//发送低五张牌
func (t *ThDesk) sendRiverCard() error {
	log.T("发送第五张公共牌begin")

	result := &bbproto.Game_SendRiverCard{}
	result.Tableid = &t.Id
	result.Card = ThCard2OGCard(t.PublicPai[4])
	t.SendRive = true

	t.THBroadcastProto(result, 0)
	time.Sleep(ThdeskConfig.TH_DONGHUA_DURATION)

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
		t.EdgeJackpot -= all.Jackpopt
		log.T("开始给allIn[%v]计算all in 时的池子金额---------------------------------end---------------", i)
		log.T("目前t.bianJackPot 的剩余值是[%v]", t.EdgeJackpot)
	}
	log.T("计算出来的allIn:【%v】", t.AllInJackpot)
	log.T("开始计算allin将近池子end")
	return nil

}

func (t *ThDesk) CheckBetUserBySeat(user *ThUser) bool {
	//2,判断押注的用户是否是当前的用户
	if t.BetUserNow != user.UserId {
		log.T("user[%v]不是desk的betUserNow[%v],押注状态不正确", user.UserId, t.BetUserNow)
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
	user.FinishtWait()
	return true
}


//获取游戏准备中的数量
func (t *ThDesk) GetGameReadyCount() int32 {
	var count int32 = 0
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.IsReady() && !u.IsBreak && !u.IsLeave {
			count ++
		}
	}
	return count
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

	//1.1,判断桌子是否是正在进行中的状态...
	if !t.IsStop() {
		log.T("desk[%v]的状态是stop[%v]的状态,所以不能开始游戏", t.Id, t.Status)
		return false
	}

	log.T("当前准备的玩家数量readyCount[%v]", t.GetGameReadyCount())
	//1.2,判断在线的用户数是否达到标准
	if t.GetGameReadyCount() < ThdeskConfig.TH_DESK_LEAST_START_USER {
		log.T("desk[%v]的准备用户数量[%v]不够", t.Id, t.GetGameReadyCount())
		return false
	}

	log.T("现在开始判断是否可以开始一局新的游戏,2,判断锦标赛的逻辑:")
	if t.IsChampionship() {
		//锦标赛游戏房间的状态
		if !ChampionshipRoom.CanNextDeskRun() {
			log.T("锦标赛的逻辑,判断不能开始下一局")
			return false
		}
	}

	log.T("现在开始判断是否可以开始一局新的游戏,1,判断自定义的逻辑:")
	if t.IsFriend() {
		//自定义房间的标准
	}

	return true

}

//判断桌子是否是stop的状态
func (t *ThDesk) IsStop() bool {
	return t.Status == TH_DESK_STATUS_STOP
}

func (t *ThDesk) IsReady() bool {
	return t.Status == TH_DESK_STATUS_READY
}

func (t *ThDesk) IsRun() bool {
	return t.Status == TH_DESK_STATUS_RUN
}

func (t *ThDesk) IsLottery() bool {
	return t.Status == TH_DESK_STATUS_LOTTERY
}

func (t *ThDesk) IsOver() bool {
	return t.Status == TH_DESK_STATUS_GAMEOVER
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

	//2,初始化玩家的信息,是否可以开始游戏,
	err := mydesk.InitUserStatus()
	if err != nil {
		log.E("开始游戏失败,errMsg[%v]", err.Error())
		return err
	}

	//3,初始化游戏房间的状态
	err = mydesk.OninitThDeskStatus()
	if err != nil {
		log.E("开始德州扑克游戏,初始化房间的状态的时候报错")
		return err
	}

	//4, 初始化前注的信息
	err = mydesk.OinitPreCoin()
	if err != nil {
		log.E("开始德州扑克游戏,初始化房间的状态的时候报错")
		return err
	}

	//5 初始化盲注开始押注
	err = mydesk.InitBlindBet()
	if err != nil {
		log.E("盲注下注的时候出错errMsg[%v]", err.Error())
		return err
	}

	//6,初始化牌的信息
	err = mydesk.OnInitCards()
	if err != nil {
		log.E("开始德州扑克游戏,初始化扑克牌的时候出错")
		return err
	}


	//7,保存用户和desk的信息到redis
	mydesk.UpdateThdeskAndAllUser2redis()

	//8,初始化玩比之后,redis中放置一份运行中的
	AddRunningDesk(mydesk)

	log.T("\n\n开始一局新的游戏,初始化完毕\n\n")
	return nil
}

//得到所有的id
func (t *ThDesk) GetuserIds() []uint32 {
	var ids []uint32
	for _, u := range t.Users {
		if u != nil {
			ids = append(ids, u.UserId)
		}
	}
	return ids
}

//锦标赛结束
/**
	关于本桌子的是否还能进入下一轮游戏
	1,如果锦标赛的游戏时间已经到了,表示锦标赛结束
	2,如果本次桌子只剩下一个人可以继续游戏,则解散
 */
func (t *ThDesk) EndCsTh() bool {

	//2,判断人数是否符合规则
	//找到符合下一局游戏的user
	var nextUsers []*ThUser
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && t.IsUserRoomCoinEnough(u) {
			nextUsers = append(nextUsers, u)
		}
	}

	//判断人数:如果只有一个人满足下一次游戏的结果,那解散这个房间,并且把user安插到其他的房间去
	if len(nextUsers) == 1 {
		//如果不是锦标赛的第一个房间,那么解散这个房间
		if ChampionshipRoom.ThDeskBuf[0] != t {
			err := ChampionshipRoom.DissolveDesk(t, nextUsers[0])
			if err != nil {
				log.E("解散房间失败")
			}
			ChampionshipRoom.Join(nextUsers[0])
		}

		return true
	}

	return false
}

//表示游戏结束,自定义房间结束
func (t *ThDesk) EndTh() bool {
	//如果是自定义房间
	log.T("判断自定义的desk是否结束游戏t.jucount[%v],t.jucountnow[%v],", t.JuCount, t.JuCountNow)
	if t.JuCountNow <= t.JuCount {
		return false
	}

	log.T("整局(多场游戏)已经结束...)")

	//设置结束时候的状态
	t.Status = TH_DESK_STATUS_GAMEOVER                //设置为没有开始开始游戏

	//广播结算的信息
	result := &bbproto.Game_SendDeskEndLottery{}
	result.Result = &intCons.ACK_RESULT_SUCC

	maxWin := int64(0)
	maxUserid := uint32(0)
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			gel := bbproto.NewGame_EndLottery()
			*gel.Coin = u.RoomCoin - t.InitRoomCoin        //这里不是u.winamount,u.winamount  表示本局赢得底池的金额
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

	//整局游戏结束之后,解散游戏房间,并且更新每个人的agent信息
	t.clearAgentData(0)
	ThGameRoomIns.RmThDesk(t) //自定义房间解散
	return true        //已经结束了本场游戏
}

//清空user的agentUserData
func (t *ThDesk) clearAgentData(ignoreUserId uint32) {
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.UserId != ignoreUserId {
			u.deskId = 0
			u.GameStatus = TH_USER_GAME_STATUS_NOGAME
			u.UpdateAgentUserData()
		}
	}
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
	//如果bettingcount>1则可以加注,如果 否则,不可以加注
	if t.GetBettingCount() > 1 {
		return 1
	} else {
		return 0
	}
}


//通过座位号来找到user
func (t *ThDesk) getUserBySeat(seatId int32) *ThUser {
	return t.Users[seatId]
}


//这个加注
func (t *ThDesk) Rebuy(userId uint32) error {

	if t.IsFriend() {
		return t.FRebuy(userId)
	} else if t.IsChampionship() {
		return t.CSRebuy(userId)
	} else {
		return nil
	}
}


//朋友桌rebuy
func (t *ThDesk) FRebuy(userId uint32) error {
	ret := bbproto.NewGame_AckRebuy()        //返回的解雇
	user := t.GetUserByUserId(userId)        //要操作的用户
	//1,为用户增加金额
	user.AddRoomCoin(t.InitRoomCoin)
	user.Update2redis()                         //rebuy需要更新redis中的缓存

	//得到需要扣除的砖石
	var feeDiamond int64 = -1
	banlance, err := userService.UpdateUserDiamond(userId, feeDiamond)
	if err != nil {
		log.E("rebuy的时候出错,error", err.Error())
		*ret.Result = intCons.ACK_RESULT_ERROR                                //错误码
		user.WriteMsg(ret)
		return err
	}
	//2,生成一条交易记录
	err = userService.CreateDiamonDetail(userId, mode.T_USER_DIAMOND_DETAILS_TYPE_REBUY, feeDiamond, banlance, "rebuy消耗钻石");
	if err != nil {
		log.E("创建用户的钻石交易记录(rebuy)失败")
		*ret.Result = intCons.ACK_RESULT_ERROR                                //错误码
		user.WriteMsg(ret)
		return err
	}
	//3,更新状态,并且返回交易结果

	*ret.Result = intCons.ACK_RESULT_SUCC                        //错误码
	*ret.CurrChip = user.RoomCoin
	*ret.RemainCount = -1        //朋友桌可以无限购买
	user.WriteMsg(ret)
	return nil
}

//锦标赛rebuy
func (t *ThDesk) CSRebuy(userId uint32) error {
	ret := bbproto.NewGame_AckRebuy()        //返回的解雇
	user := t.GetUserByUserId(userId)        //要操作的用户


	//0,购买之前需要判断锦标赛的状态
	if !user.CSGamingStatus {
		log.E("用户已经放弃了比赛,不能在rebuy了")
		*ret.Result = intCons.ACK_RESULT_ERROR                                //错误码
		user.WriteMsg(ret)
		return errors.New("用户已经发起比赛了,不能在rebuy")
	}


	//先扣除钻石,扣除成功之后,再给用户加上roomCoin
	var feeDiamond int64 = -1
	banlance, err := userService.UpdateUserDiamond(userId, feeDiamond)
	if err != nil {
		log.E("rebuy的时候出错,error", err.Error())
		*ret.Result = Error.GetErrorCode(err)                               //错误码
		user.WriteMsg(ret)
		return err
	}

	//如果用户的钱足够,那么设置喂准备的状态
	if t.IsUserRoomCoinEnough(user) {
		user.Status = TH_USER_STATUS_READY
	}

	//1,为用户增加金额
	//rebuy需要更新redis中的缓存
	user.AddRoomCoin(t.InitRoomCoin)
	user.Update2redis()


	//2,生成一条交易记录
	err = userService.CreateDiamonDetail(userId, mode.T_USER_DIAMOND_DETAILS_TYPE_REBUY, feeDiamond, banlance, "rebuy消耗钻石");
	if err != nil {
		log.E("创建用户的钻石交易记录(rebuy)失败")
		*ret.Result = intCons.ACK_RESULT_ERROR                                //错误码
		user.WriteMsg(ret)
		return err
	}
	//3,更新状态,并且返回交易结果

	*ret.Result = intCons.ACK_RESULT_SUCC                        //错误码
	*ret.CurrChip = user.RoomCoin
	*ret.RemainCount = t.RebuyCountLimit - user.RebuyCount        //锦标赛的购买次数限制
	user.WriteMsg(ret)
	return nil
}


//用户放弃购买
func (t *ThDesk) NotRebuy(userId uint32) error {
	//1,区分锦标赛还是朋友桌不重够的种情况
	if t.IsFriend() {
		//朋友桌
		t.FNotRebuy(userId)
	} else if t.IsChampionship() {
		//锦标赛
		t.CSNotRebuy(userId)
	}


	//2,不重购之后,需要在redis中保存用户的信息
	t.UpdateThdeskAndUser2redis(t.GetUserByUserId(userId))

	return nil
}

func (t *ThDesk) FNotRebuy(userId uint32) {
	//房主和非房主之间的处理方式是不同的
	if t.IsOwner(userId) {
		//是房主的处理方式
		//1,房主易主
		t.ChangeOwner()
	} else {
		//不是房主的处理方式
		//1,检测是否可以开始
		if t.JuCountNow > 1 && t.IsAllReady() {
			//用户离开之后,判断游戏是否开始
			go t.Run()
		}
	}
	//2,返回信息

}

func (t *ThDesk) ChangeOwner() error {

	var ouId uint32 = 0
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && t.IsUserRoomCoinEnough(u) {
			ouId = u.UserId
			break
		}
	}

	//没有找到合适房主的请求
	if ouId == 0 {
		return errors.New("没有找到合适的房主")
	}

	oldOwnerUserId := t.DeskOwner        //旧的房主
	t.DeskOwner = ouId        //设置新的房主

	result := bbproto.NewGame_SendChangeDeskOwner()
	*result.DeskId = t.Id
	*result.OldOwner = oldOwnerUserId
	*result.OldOwnerSeat = t.GetUserSeatByUserId(oldOwnerUserId)

	*result.NewOwner = t.DeskOwner
	*result.NewOwnerSeat = t.GetUserSeatByUserId(t.DeskOwner)

	//发送房主变更的广播
	t.THBroadcastProtoAll(result)
	return nil

}


//判断这个userId 是不是房主
func (t *ThDesk) IsOwner(userId uint32) bool {
	return t.DeskOwner == userId
}

//锦标赛不重新够买的处理
func (t *ThDesk) CSNotRebuy(userId uint32) {
	log.T("锦标赛user[%v]notrebuy的请求", userId)
	//1,设置当前用户的锦标赛状态 为结束
	user := t.GetUserByUserId(userId)
	if !user.CSGamingStatus {
		log.E("用户已经放弃了比赛,不能重新notRebuy了")
		return
	}

	user.RebuyCount = ChampionshipRoom.RebuyCountLimit                //取消重构之后,下一局就不能重购买了
	user.CSGamingStatus = false;

	//2,取消之后,现实最终的排名

	log.T("用户notRebuy的时候,发送User的最终排名...")
	ret := bbproto.NewGame_TounamentPlayerRank()
	*ret.Message = "测试最终排名的信息"
	*ret.PlayerRank = ChampionshipRoom.GetRankByuserId(user.UserId)
	user.WriteMsg(ret)

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

	case TH_DESK_BET_TYPE_CHECK:
		//让牌
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

		//保存用户和desk的数据到redis
		t.UpdateThdeskAndUser2redis(user)
		return nil
	}
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
	*result.NextSeat = t.GetUserSeatByUserId(t.BetUserNow) //		//下一个押注的人
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
	user.LotteryCheck = true

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
	*result.NextSeat = t.GetUserSeatByUserId(t.BetUserNow)        //		//下一个押注的人

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
	*result.NextSeat = t.GetUserSeatByUserId(t.BetUserNow)        //		//下一个押注的人
	*result.HandCoin = user.HandCoin                        //表示需要加注多少

	//给所有人广播信息
	t.THBroadcastProto(result, 0)
	log.T("开始处理seat[%v]加注的逻辑,t,OGRaiseBet()...end", user.Seat)
	return nil
}


//得到当前用户需要加注的金额
func (t *ThDesk) GetMinRaise() int64 {

	//得到当前押注的人...
	betUser := t.GetUserByUserId(t.BetUserNow)
	var betUserTurnCoin int64 = 0
	if betUser != nil {
		betUserTurnCoin = betUser.TurnCoin
	}

	log.T("获取用户[%v]的最低加注金额t.TurnCoin[%v],t.MinRaise[%v],t.BetAmountNow[%v],", t.BetUserNow, betUserTurnCoin, t.MinRaise, t.BetAmountNow)

	result := t.MinRaise + t.BetAmountNow - betUserTurnCoin
	if result < 0 {
		//现在处理的有可能是新的一局开始
		result = t.BigBlindCoin
	}
	log.T("获取用户[%v]的最低加注金额是[%v]", t.BetUserNow, result)
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
	*result.NextSeat = t.GetUserSeatByUserId(t.BetUserNow)        //		//下一个押注的人
	*result.HandCoin = user.HandCoin

	t.THBroadcastProtoAll(result)

	return nil
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
	atomic.AddInt64(&t.EdgeJackpot, coin)
}

//把程序的allin 转化为proto的格式
func (t *ThDesk) GetServerProtoAllInJackPot() []*bbproto.ThServerAllInJackpot {
	var result []*bbproto.ThServerAllInJackpot
	for i := 0; i < len(t.AllInJackpot); i++ {
		a := t.AllInJackpot[i]
		if a != nil {
			sa := bbproto.NewThServerAllInJack()
			*sa.UserId = a.UserId
			*sa.AllInAmount = a.AllInAmount
			*sa.Jackpopt = a.Jackpopt
			*sa.ThroundCount = a.ThroundCount
			result = append(result, sa)
		}
	}
	return result
}

