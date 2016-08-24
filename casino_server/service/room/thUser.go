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
)




//德州扑克 玩家的状态
var TH_USER_STATUS_WAITSEAT int32 = 1           //刚上桌子 等待开始的玩家
var TH_USER_STATUS_SEATED int32 = 2                //刚上桌子 但是没有在游戏中
var TH_USER_STATUS_READY int32 = 3
var TH_USER_STATUS_BETING int32 = 4                //押注中
var TH_USER_STATUS_ALLINING int32 = 5                //allIn
var TH_USER_STATUS_FOLDED int32 = 6                //弃牌
var TH_USER_STATUS_WAIT_CLOSED int32 = 7                 //等待结算
var TH_USER_STATUS_CLOSED int32 = 8                //已经结算


/**
	正在玩德州的人
 */
type ThUser struct {
	sync.Mutex
	UserId             uint32                //用户id
	NickName           string                //用户昵称

	deskId             int32                 //用户所在的桌子的编号
	MatchId            int32                 //matchId
	GameNumber         int32                 //游戏编号
	Seat               int32                 //用户的座位号
	agent              gate.Agent            //agent
	Status             int32                 //当前的状态
	IsBreak            bool                  //用户断线的状态,这里判断用户是否断线
	IsLeave            bool                  //用户是否处于离开的状态
	HandCards          []*bbproto.Pai        //手牌
	thCards            *pokerService.ThCards //手牌加公共牌取出来的值,这个值可以实在结算的时候来取
	waiTime            time.Time             //等待时间
	waitUUID           string                //等待标志
	TotalBet           int64                 //计算用户总共押注的多少钱
	TotalBet4calcAllin int64                 //押注总额 ***注意,目前这个值是用来计算all in 的
	winAmount          int64                 //总共赢了多少钱
	winAmountDetail    []int64               //赢钱的细节, 主要是每个记录每个奖池赢了多少钱
	TurnCoin           int64                 //单轮押注(总共四轮)的金额
	HandCoin           int64                 //用户下注多少钱、指单局
	RoomCoin           int64                 //用户上分的金额
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
	log.T("用户当前的状态[%v]", t.Status)
	if t.Status != TH_USER_STATUS_BETING {
		return nil
	}

	ticker := time.NewTicker(time.Second * 1)
	t.waiTime = time.Now().Add(ThdeskConfig.TH_TIMEOUT_DURATION)
	uuid, _ := uuid.NewV4()
	t.waitUUID = uuid.String()                //设置出牌等待的标志
	go func() {
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

//返回自己所在的桌子
func (t *ThUser) GetDesk() *ThDesk {
	desk := GetDeskByAgent(t.agent)
	return desk
}

//用户超时,做处理
func (t *ThUser) TimeOut(timeNow time.Time) (bool, error) {
	t.Lock()
	defer t.Unlock()

	//没有超时标志,直接返回
	if t.waitUUID == "" {
		//不需要等待
		log.T("用户[%v]的waitUUID==空,不用超时", t.UserId)
		return true, nil
	}

	//如果用户超市,或者用户选择离线,那么直接做弃牌的操作
	if t.waiTime.Before(timeNow) || t.IsLeave {
		log.T("玩家[%v]超时,现在做超时的处理", t.UserId)
		err := t.GetDesk().DDBet(t.Seat, TH_DESK_BET_TYPE_FOLD, 0)
		if err != nil {
			log.E("用户[%v]弃牌失败", t.UserId)
		}
		//这里需要设置为弃牌的状态
		log.T("玩家[%v]超时,现在做超时的处理,处理完毕", t.UserId)
		return true, err
	} else {
		//没有超时,继续等待
		log.T("玩家[%v]nickname[%v]出牌中还没有超时", t.UserId, t.NickName)
		return false, nil
	}
}

func (t *ThUser) InitWait() {
	t.waitUUID = ""
}

//判断用户是否正在等待出牌
func (t *ThUser) IsWaiting() bool {
	return t.waitUUID != ""
}


//操作押注时的waiting 状态
func (t *ThUser) CheckBetWaitStatus() error {
	if t.IsWaiting() {
		t.InitWait()
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
	return result
}

//更新用户的agentUserData数据
func (t *ThUser) UpdateAgentUserData(a gate.Agent, deskId int32, matchId int32) {

	//保存回话信息
	userAgentData := bbproto.NewThServerUserSession()                //绑定参数
	*userAgentData.UserId = t.UserId
	*userAgentData.DeskId = deskId
	*userAgentData.MatchId = matchId
	a.SetUserData(userAgentData)
	t.agent = a
	t.IsBreak = false

	//回话信息保存到redis
	userService.SaveUserSession(userAgentData)
}

//把用户数据保存到redis中
func (u *ThUser) Update2redis() {
	UpdateRedisThuser(u)
}

//计算用户的各种金额

func (t *ThUser) AddRoomCoin(coin int64) {
	atomic.AddInt64(&t.RoomCoin, coin)
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

//得到自己在当前锦标赛中的名次
func (t *ThUser) GetCsRank() int64 {
	//todo 获取名词之前,需要判断用户正在玩的游戏的类型
	return GetCSTHuserRank(t.MatchId,t.UserId)
}

