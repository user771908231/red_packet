package room

import (
	"sync"
	"casino_server/msg/bbprotogo"
	"casino_server/service/pokerService"
	"time"
	"casino_server/service/userService"
	"github.com/nu7hatch/gouuid"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"errors"
	"sync/atomic"
	"github.com/golang/protobuf/proto"
	"casino_server/utils/jobUtils"
	"casino_server/common/Error"
)




//德州扑克 玩家的状态

var TH_USER_STATUS_NOGAME int32 = 0           //刚上桌子 没有游戏的玩家
var TH_USER_STATUS_WAITSEAT int32 = 1         //刚上桌子 等待入座的玩家
var TH_USER_STATUS_SEATED int32 = 2           //刚上桌子 已经入座没有准备的玩家
var TH_USER_STATUS_READY int32 = 3              //已经准备的玩家
var TH_USER_STATUS_BETING int32 = 4           //押注中
var TH_USER_STATUS_ALLINING int32 = 5         //allIn
var TH_USER_STATUS_FOLDED int32 = 6           //弃牌
var TH_USER_STATUS_WAIT_CLOSED int32 = 7      //等待结算
var TH_USER_STATUS_CLOSED int32 = 8           //已经结算


var TH_USER_GAME_STATUS_NOGAME int32 = 0      //不在游戏中
var TH_USER_GAME_STATUS_FRIEND int32 = 1      //朋友桌
var TH_USER_GAME_STATUS_CHAMPIONSHIP int32 = 2                //锦标赛


/**
	正在玩德州的人
 */
type ThUser struct {
	sync.Mutex
	UserId             uint32                //用户id
	NickName           string                //用户昵称

	deskId             int32                 //用户所在的桌子的编号
	Seat               int32                 //用户的座位号
	Agent              gate.Agent            //agent
	Status             int32                 //当前的状态,单局游戏的状态
	CSGamingStatus     bool                  //是否进行锦标赛,这个字段其实 是在服务器crash之后,恢复数据的时候可以用到
	GameStatus         int32                 //用户的游戏状态
	IsBreak            bool                  //用户断线的状态,这里判断用户是否断线
	IsLeave            bool                  //用户是否处于离开的状态
	HandCards          []*bbproto.Pai        //手牌
	thCards            *pokerService.ThCards //手牌加公共牌取出来的值,这个值可以实在结算的时候来取
	waiTime            time.Time             //等待时间
	waitUUID           string                //等待标志
	PreCoin            int64                 //前注
	TotalBet           int64                 //计算用户总共押注的多少钱
	TotalBet4calcAllin int64                 //押注总额 ***注意,目前这个值是用来计算all in 的
	winAmount          int64                 //总共赢了多少钱
	winAmountDetail    []int64               //赢钱的细节, 主要是每个记录每个奖池赢了多少钱
	TurnCoin           int64                 //单轮押注(总共四轮)的金额
	HandCoin           int64                 //用户下注多少钱、指单局
	RoomCoin           int64                 //用户上分的金额
	TotalRoomCoin      int64                 //用户总的带入金额是多少钱
	RebuyCount         int32                 //重购的次数
	LotteryCheck       bool                  //这个字段用于判断是否可以开奖,默认是false:   1,如果用户操作弃牌,则直接设置为true,2,如果本局是all in,那么要到本轮次押注完成之后,才能设置为true
	IsShowCard         bool                  //是否亮牌
}

func (t *ThUser) GetCoin() int64 {
	redu := userService.GetUserById(t.UserId)
	if redu == nil {
		return -1
	} else {
		return redu.GetCoin()
	}
}

func (t *ThUser) GetRoomCoin() int64 {
	return t.RoomCoin
}

func (t *ThUser) GetStatusDes() string {
	des := ""
	switch t.Status {
	case 0:
		des = "刚上桌子"
	case 1:
		des = "等待入座"
	case 2:
		des = "已经坐下"
	case 3:
		des = "已经准备"
	case 4:
		des = "押注中"
	case 5:
		des = "all"
	case 6:
		des = "弃牌"
	case 7:
		des = "等待结算"
	case 8:
		des = "已经结算"
	}

	return des

}
//
func (t *ThUser) trans2bbprotoThuser() *bbproto.THUser {
	thuserTemp := &bbproto.THUser{}
	thuserTemp.Status = &(t.Status)        //已经就做
	thuserTemp.User = userService.GetUserById(t.UserId)        //得到user
	thuserTemp.HandPais = t.HandCards
	thuserTemp.SeatNumber = new(int32)
	return thuserTemp
}

//等待用户出牌
func (t *ThUser) wait() error {
	//如果不是押注中的状态,不用wait任务
	//这一步其实可以不用验证,因为出牌的游标滑动到这里的时候,已经验证过了
	//if !t.IsBetting() {
	//	return nil
	//}

	ticker := time.NewTicker(time.Second * 1)
	t.waiTime = time.Now().Add(ThdeskConfig.TH_TIMEOUT_DURATION)
	uuid, _ := uuid.NewV4()
	t.waitUUID = uuid.String()                //设置出牌等待的标志
	go func() {
		//tod 目前是测试性质的代码
		defer Error.ErrorRecovery("user.timeout()")

		for timeNow := range ticker.C {
			//表示已经过期了
			bool, err := t.TimeOut(timeNow)
			if err != nil {
				log.E("处理超时的逻辑出现错误,errMsg[%v]", err.Error())
				return
			}

			//判断是否已经超时
			if bool {
				log.E("user[%v]已经超时,结束等待任务", t.UserId)
				return
			}
		}
	}()
	return nil

}


//等待用户rebuy,时间过了之后,用户余额还是不够,那么游戏结束
func (t *ThUser) waitCsRebuy() {
	log.T("user【%v】开始等待重购买", t.UserId)
	timeEnd := time.Now().Add(time.Second * 10)        //5秒之后
	jobUtils.DoAsynJob(time.Second * 5, func() bool {
		if (time.Now().After(timeEnd)) {
			desk := t.GetDesk()
			if desk != nil && !desk.IsUserRoomCoinEnough(t) {
				log.T("user【%v】等待重购超时,系统自动notRebuy...", t.UserId)
				desk.CSNotRebuy(t.UserId)
			}
			return true
		} else {
			log.T("user【%v】等待重购ing...", t.UserId)
			return false
		}
	})

}

//返回自己所在的桌子
func (t *ThUser) GetDesk() *ThDesk {
	var desk *ThDesk
	if t.CSGamingStatus {
		desk = ChampionshipRoom.GetDeskById(t.deskId)
	} else {
		desk = ThGameRoomIns.GetDeskById(t.deskId)
	}
	return desk
}

//用户超时,做处理
func (t *ThUser) TimeOut(timeNow time.Time) (bool, error) {
	t.Lock()
	defer t.Unlock()

	//没有超时标志,直接返回
	if !t.IsWaiting() {
		log.T("用户[%v]的waitUUID==空,不用超时", t.UserId)
		return true, nil
	}

	//如果用户超时,或者用户已经离开,那么直接做弃牌的操作
	if t.waiTime.Before(timeNow) || t.IsLeave {
		log.T("玩家[%v]超时,现在做超时的处理", t.UserId)
		desk := t.GetDesk()
		if desk == nil {
			log.T("用户等待超时,自动弃牌的时候,desk为空,导致弃牌失败.")
			return true, nil                //表示游戏结束
		}
		log.T("玩家[%v]超时,现在做超时的处理:开始系统自动弃牌", t.UserId)
		err := desk.DDBet(t.Seat, TH_DESK_BET_TYPE_FOLD, 0)
		if err != nil {
			log.E("用户[%v]弃牌失败", t.UserId)
		}
		//这里需要设置为弃牌的状态
		log.T("玩家[%v]超时,现在做超时的处理,处理完毕", t.UserId)
		return true, err
	} else {
		//没有超时,继续等待
		log.T("desk[%v]玩家[%v]nickname[%v]出牌中还没有超时", t.deskId, t.UserId, t.NickName)
		return false, nil
	}
}

//用户结束等待
func (t *ThUser) FinishtWait() {
	t.waitUUID = ""
}

//判断用户是否正在等待出牌
func (t *ThUser) IsWaiting() bool {
	return t.waitUUID != ""
}


//操作押注时的waiting 状态
func (t *ThUser) CheckBetWaitStatus() error {
	if t.IsWaiting() {
		return nil
	} else {
		return errors.New("用户状态错误")
	}
}

func NewThUser() *ThUser {
	result := &ThUser{}
	result.UserId = 0
	result.Status = 0
	result.TurnCoin = 0
	result.TotalBet4calcAllin = 0
	result.TotalBet = 0
	result.winAmount = 0
	result.RoomCoin = 0
	result.IsBreak = false
	result.IsLeave = false
	result.IsShowCard = false
	return result
}

//更新用户的agentUserData数据
func (u *ThUser) UpdateAgentUserData() {
	userAgentData := bbproto.NewThServerUserSession()                //绑定参数

	//保存回话信息
	*userAgentData.UserId = u.UserId
	*userAgentData.DeskId = u.deskId
	*userAgentData.GameStatus = u.GameStatus //返回用户当前的状态 0：未游戏  1：正在朋友桌  2：正在锦标赛
	*userAgentData.IsBreak = u.IsBreak
	*userAgentData.IsLeave = u.IsLeave
	agent := u.Agent
	//这里出现nil的情况是,回复数据的时候,用户没有连接的时候就有可能出现nil
	if agent != nil {
		agent.SetUserData(userAgentData)        //设置用户的agentData
	}

	desk := u.GetDesk()
	if desk == nil {
		*userAgentData.MatchId = 0
		*userAgentData.RoomKey = ""
	} else {
		*userAgentData.MatchId = desk.MatchId
		*userAgentData.RoomKey = desk.RoomKey
	}
	//回话信息保存到redis
	userService.SaveUserSession(userAgentData)
}

//把用户数据保存到redis中
func (u *ThUser) Update2redis() {
	//log.T("用户数据改变之后的值,需要保存在rendis中[%v]", *u)
	UpdateRedisThuser(u)
}

//计算用户的各种金额


func (t *ThUser) AddPreCoin(coin int64) {
	atomic.AddInt64(&t.PreCoin, coin)
}

func (t *ThUser) AddRoomCoin(coin int64) {
	atomic.AddInt64(&t.RoomCoin, coin)
}

func (t *ThUser) AddTotalRoomCoin(coin int64) {
	atomic.AddInt64(&t.TotalRoomCoin, coin)
}

func (t *ThUser) AddWinAmount(coin int64) {
	atomic.AddInt64(&t.winAmount, coin)
}

func (t *ThUser) AddTurnCoin(coin int64) {
	atomic.AddInt64(&t.TurnCoin, coin)
}

func (t *ThUser) AddHandCoin(coin int64) {
	atomic.AddInt64(&t.HandCoin, coin)
}

func (t *ThUser) AddTotalBet4calcAllin(coin int64) {
	atomic.AddInt64(&t.TotalBet4calcAllin, coin)
}

func (t *ThUser) AddTotalBet(coin int64) {
	atomic.AddInt64(&t.TotalBet, coin)
}

//判断用户是否正在押注中
func (t *ThUser) IsBetting() bool {
	//正在押注中 是否需要判断是否断线,是否离线?
	return t.Status == TH_USER_STATUS_BETING
}

func (t *ThUser) IsAllIn() bool {
	//正在押注中 是否需要判断是否断线,是否离线?
	return t.Status == TH_USER_STATUS_ALLINING
}

func (t *ThUser) IsFold() bool {
	//正在押注中 是否需要判断是否断线,是否离线?
	return t.Status == TH_USER_STATUS_FOLDED
}

func (t *ThUser) IsClose() bool {
	return t.Status == TH_USER_STATUS_CLOSED
}


//查看用户是否准备
func (t *ThUser) IsReady() bool {
	return t.Status == TH_USER_STATUS_READY
}

//得到自己在当前锦标赛中的名次
func (t *ThUser) GetCsRank() int64 {
	return int64(0)
}

func (t *ThUser) WriteMsg(p proto.Message) error {
	agent := t.Agent
	if agent != nil {
		log.T("agent.WriteMsg()[%v]:", t.UserId)
		agent.WriteMsg(p)
		return nil
	} else {
		log.T("用户[%v]的agent为nil,不能发送信息", t.UserId)
		return errors.New("用户的agent为nil,不能发送信息")
	}

}

//通过agent得到session
func GetUserDataByAgent(a gate.Agent) *bbproto.ThServerUserSession {
	//获取agent中的userData
	ad := a.UserData()
	if ad == nil {
		log.E("agent中的userData为nil")
		return nil
	}

	userData := ad.(*bbproto.ThServerUserSession)
	log.T("得到的UserAgent中的userId是[%v]", userData.UserId)
	return userData
}
