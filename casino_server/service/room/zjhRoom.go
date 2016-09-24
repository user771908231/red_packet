package room

import (
	"time"
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/service/userService"
	"strings"
	"casino_server/utils/redis"
	"casino_server/utils/numUtils"
	"casino_server/utils/timeUtils"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/conf/casinoConf"
	"casino_server/utils/db"
)

func init() {
	//OninitZjgRoom()
}


//config

var ZJH_BET_DURATION = time.Second * 10
var ZJH_LOTTERY_DURATION = time.Second * 1
var ZJH_AFTER_LOTTERY_DURATION = time.Second * 10
var ZJH_ROUND_NUMBER_PRE string = "zjh_round_pre"


//-------------------------------------------------扎金花房间的状态----------------------------------------------
var ZJH_STATUS_BETING int32 = 1                //押注中
var ZJH_STATUS_LOTTERING int32 = 2                //等待开奖
var ZJH_STATUS_LOTTERIED int32 = 3                //已经开奖

var ZJHroom zjhRoom        //扎金花的房间

/**
扎金花的room
todo 可以把游戏相关的信息放置在邮箱相关的结构体中
 */
type zjhRoom struct {
	room                              //继承基本的room的方法
	id             uint32             //数据库所以呢id
	BetStartTime   time.Time          //首次投注时间
	BetEndTime     time.Time          //结束投注时间
	LotteryTime    time.Time          // 开奖时间
	NextStartTime  time.Time          //下次开始时间
	Status         int32              //房间状态:押注中,开奖中,
	Jackpot        int64              //奖池大小
	zoneAmount     []int32            //押注区域的大小
	zoneWinAmount  []int32            //押注区的输赢分数
	BankerUserId   uint32             //庄家的信息
	Zjhpai         bbproto.ZjhPaiList //本轮游戏的纸牌
	ZjhRoundNumber string             //本局游戏的编号
}

/**
初始化扎金花的房间
 */
func OninitZjgRoom() {
	log.T("初始化扎金花的房间")
	//初始化参数 agent 的集合
	ZJHroom.AgentMap = make(map[uint32]gate.Agent)                //初始化agent缓冲区
	ZJHroom.zoneAmount = make([]int32, 4)                        //初始化押注区
	ZJHroom.zoneWinAmount = make([]int32, 4)                        //初始化押注区的输赢分数
	ZJHroom.Oninit(time.Now())                                //初始化开始押注时间,押注结束时间,开奖时间

	if casinoConf.SWITCH_ZJH {
		ZJHroom.run()                                                //启动扎金花房间的run任务
	}

}


/**
	初始化room
 */
func (r *zjhRoom) Oninit(t time.Time) {
	r.iniTime(t)                        //初始化时间
	r.OnInitZjhpai()                //初始化纸牌
	r.Status = ZJH_STATUS_BETING        //初始化状态
	r.Jackpot = 99909                //初始化奖池的大小
	r.zoneWinAmount = []int32{0, 0, 0, 0}
	//r.BroadcastBeginBet()		//广播可以开始押注了--这一步目前是放在next 函数中使用的...
	r.saveZjhRound()                //保存本局游戏到数据库

}

/**
	初始化本局编号
 */
func (r *zjhRoom) OnInitRoundNumber() {
	idStr, _ := numUtils.Uint2String(r.id)
	log.T("新的一轮扎金花游戏的-number-idstr[%v],", idStr)
	rtStr := strings.Join([]string{ZJH_ROUND_NUMBER_PRE, idStr}, "_")
	log.T("新的一轮扎金花游戏的-number[%v],", rtStr)
	r.ZjhRoundNumber = rtStr
}

/**
	保存本局扎金花的数据到数据库
	1.保存的过程中得到 id之后可以初始化本局游戏编号
 */
func (r *zjhRoom) saveZjhRound() error {
	bst := r.BetStartTime.UnixNano()        //开始押注时间毫秒数
	bnt := r.BetEndTime.UnixNano()                //结束押注时间
	ltt := r.LotteryTime.UnixNano()                //开奖时间

	//得到需要保存的数据
	d := &bbproto.TZjhRound{}
	d.BeginTime = &bst
	d.BetEndTime = &bnt
	d.LotteryTime = &ltt
	d.BankerUserId = &r.BankerUserId
	d.ZjhPaiList = r.Zjhpai
	d.ZoneAmount = r.zoneAmount
	d.ZoneWinAmount = r.zoneWinAmount

	// 获取连接 connection

	id, _ := db.GetNextSeq(casinoConf.DBT_T_ZJH_ROUND)
	log.T("通过数据库获取到的 t_zjh_round seq[%v]", id)
	r.id = uint32(id)
	d.Id = &r.id
	r.OnInitRoundNumber()                //初始化编号
	d.Number = &r.ZjhRoundNumber        //编号也要保存到数据库
	e := db.InsertMgoData(casinoConf.DBT_T_ZJH_ROUND, d)
	if e != nil {
		log.Error(e.Error())
		return e
	}
	return nil

}

func (r *zjhRoom) iniTime(t time.Time) {
	//r.Lock()
	//defer r.Unlock()
	r.BetStartTime = t                                        //首次开始押注的时间
	r.BetEndTime = r.BetStartTime.Add(ZJH_BET_DURATION)        //首次押注结束的时间
	r.LotteryTime = r.BetEndTime.Add(ZJH_LOTTERY_DURATION)        //首次开奖的时间
	r.NextStartTime = r.LotteryTime.Add(ZJH_AFTER_LOTTERY_DURATION)        //下次开始的时间
	//下次押注的时间
}

//初始化本轮的牌
func (r *zjhRoom) OnInitZjhpai() {
	//list := porkService.CreateZjhList()
	//zs := make([]*bbproto.ZjhPai,5)
	//for i := 0; i < 5; i++ {
	//	zs[i] 	= porkService.ConstructZjhPai(list[i])
	//}
	//r.Zjhpai = zs

	r.Zjhpai = bbproto.CreateZjhList()

}


/**
扎金花需要的启动函数
启动函数实际上是一个定时任务
	1,判断当前是否为投注状态,如果是投注的状态则可以投注
	2,判断当前是否为开奖的时间,如果是则不能投注,等待开奖
 */
func (r *zjhRoom) run() {
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for t := range ticker.C {
			log.T("正在执行扎金花的逻辑当前状态:\n status[%v],\n betEndTime[%v],\nlotterTime[%v],\nnextStartTime[%v],\nbankerUserId[%v],\njackPort[%v],\nnow[%v]",
				r.Status, r.BetEndTime.Format(timeUtils.TIME_LAYOUT), r.LotteryTime.Format(timeUtils.TIME_LAYOUT), r.NextStartTime.Format(timeUtils.TIME_LAYOUT), r.BankerUserId, r.Jackpot, t.Format(timeUtils.TIME_LAYOUT))
			if len(r.AgentMap) < 1 {
				log.T("没有玩家进入游戏...continue\n\n\n")
				//如果没有玩家进入 需要重置房间的状态和时间
				r.iniTime(t)
				continue
			}

			//这里测试代码,打印当前正在游戏中的玩家
			log.T("当前玩家正在房间中")
			/* 使用 key 输出 map 值 */
			for a := range r.AgentMap {
				log.T("agent-userId:%v", a)
			}

			//投注结束的广播
			r.betEnd(t)
			//开奖的广播
			r.lottery(t)
			//进行下一轮
			r.next(t)

		}

	}()
}

/**
判断当前房间是否可以扎金花
 */
func (r *zjhRoom) Betable() bool {
	if r.Status == ZJH_STATUS_BETING && r.BetEndTime.After(time.Now()) {
		return true
	} else {
		return false
	}
}


/**
	为每个人清算本局输赢
	1,得到当前的押注记录
	2,比较大小

 */


/**
	这个是开奖的广播,开奖的流程如下
	1,清算本局扎金花的输赢结果
	2,分别给每个人发送自己的得分,牌的结果
	3,设置房间的状态准备下一轮游戏开始
 */
func (r *zjhRoom) lottery(t time.Time) {
	r.Lock()
	defer r.Unlock()
	//如果当前时间已经过了开奖时间,并且现在的状态是开奖中,则重新设置状态,并且开奖
	if r.LotteryTime.Before(t) && r.Status == ZJH_STATUS_LOTTERING {
		log.T("-----------------------------------------开奖-----------------------------------------")

		//1,开奖之前,首先为每个人结账,清算
		//如果游戏房间里没有人,直接返回不用清算
		if len(r.AgentMap) < 1 {
			log.E("游戏房间中没有玩家,不用结算得分")
			return
		}

		//一次对每个人对清算,并且统计m,key 表示userId
		for key := range r.AgentMap {
			log.T("开始给玩家userId[%v]结算本局得分", key)
			//这里通过userId (key)取到用户的押注记录和用户信息,再更具牌面的大小来做结算
			betRecode := GetTBetRecordByUserIdAnd(r.ZjhRoundNumber, key)
			if betRecode.GetUserId() == 0 {
				log.T("玩家userId[%v]没有押注", key)
				continue
			}

			//打印目前的押注
			log.T("玩家[%v]的押注信息:[%v]", key, betRecode)

			//开始计算得分
			for i, d := range betRecode.Betzone {
				//比较牌的大小,从第二组牌开始算.如果庄家的牌大,则赢钱
				if r.Zjhpai.Less(0, (i + 1)) {
					log.T("押注[%v]输了", i)
					var btwa int32 = betRecode.GetWinAmount() - d
					betRecode.WinAmount = &btwa
					r.zoneWinAmount[i] += d                //游戏房间押注区域分数计算,庄家赢,庄家加分
				} else {
					log.T("押注[%v]赢了", i)
					var btwa int32 = betRecode.GetWinAmount() + d
					betRecode.WinAmount = &btwa                //庄家赢,个人加分
					r.zoneWinAmount[i] -= d                //游戏房间押注区域分数计算,庄家赢,庄家减分
				}
			}

			log.T("给玩家userId[%v]结算本局得分结束,得分【%v】", key, betRecode.GetWinAmount())

			//2,真是数据给每个人发送开奖信息
			result := &bbproto.ZjhLottery{}
			result.Header = protoUtils.GetSuccHeader()                                        //包头,返回结果
			result.Zjhpai = r.Zjhpai                                                        //纸牌中,第一幅牌是庄家的牌
			userMe := userService.GetUserById(key)
			if userMe == nil {
				log.E("用户userId[%v]没有找到.", key)
				return
			}

			result.Balance = new(int32)
			result.WinAmount = new(int32)                                        //本局的输赢分数
			*result.Balance = int32(*userService.GetUserById(key).Coin)                                //自己的余额
			*result.WinAmount = int32(*betRecode.WinAmount)                                              //本局的输赢分数

			//更新用户的信息,保存用户信息到redis
			SaveBetRecord(betRecode)

			//todo 不需要保存用户信息到数据库
			//userService.UpUserBalance(key, int64(betRecode.GetWinAmount()), userService.UPDATE_TYPE_ONLY_REDIS)

			//利用agent发送数据

			log.Normal("开始给%d发送消息", key)
			a := r.AgentMap[key]
			a.WriteMsg(result)
			log.Normal("给%v发送消息完毕", key)
		}

		//3,设置房间状态
		r.Status = ZJH_STATUS_LOTTERIED                                      //设置状态已经开奖
		log.T("---------------------------------------开奖结束-----------------------------------------")
	}
}

/**
	开始下一轮的游戏
 */
func (r *zjhRoom) next(t time.Time) {
	r.Lock()
	defer r.Unlock()
	if t.After(r.NextStartTime) && r.Status == ZJH_STATUS_LOTTERIED {
		log.T("---------------------------------------初始化下一轮-----------------------------------------")
		//开奖已经结束了..可以重新开始
		r.Oninit(t)
		r.BroadcastBeginBet()
		log.T("--------------------------------------初始化下一轮结束---------------------------------------")

	}

}

func (r *zjhRoom) betEnd(t time.Time) {
	r.Lock()
	defer r.Unlock()

	if r.BetEndTime.Before(t) && r.Status == ZJH_STATUS_BETING {
		r.Status = ZJH_STATUS_LOTTERING
		log.T("---------------------------------------押注结束-----------------------------------------")

	}

}


/**
	给对应的zone增加对应的积分
 */
func (r *zjhRoom) AddZoneAmount(d []int32) {
	r.Lock()
	defer r.Unlock()
	for i := 0; i < len(d); i++ {
		r.zoneAmount[i] = d[i]
	}
	log.T("增加积分之后,目前房间的zoneAmount:%v", r.zoneAmount)

}

/**
	返回押注剩余的时间
 */
func (r *zjhRoom) GetBetRemainTime() *int32 {
	var result int32
	if r.Status == ZJH_STATUS_BETING {
		now := time.Now()
		diff := r.BetEndTime.Sub(now)
		result = int32(diff.Seconds())
	} else {
		result = 0
	}
	return &result
}

/**
	返回开奖剩余的时间
 */
func (r *zjhRoom) GetLotteryRemainTime() *int32 {
	var result int32
	if r.Status == ZJH_STATUS_BETING || r.Status == ZJH_STATUS_LOTTERING {
		//押注中的状态或者等待开奖的状态的时候,可以返回等待开奖的时间
		now := time.Now()
		diff := r.LotteryTime.Sub(now)
		result = int32(diff.Seconds())
	} else {
		result = 0
	}
	return &result

}


/**
	广播消息,开始押注
 */
func (r *zjhRoom) BroadcastBeginBet() {
	log.T("所有人发送可以开始押注的广播:")
	//通知押注的信息
	result := &bbproto.ZjhBroadcastBeginBet{}
	result.Jackpot = &r.Jackpot
	result.Banker = userService.GetUserById(r.BankerUserId)
	result.Zjhpai = r.Zjhpai
	result.BetTime = r.GetLotteryRemainTime()
	result.LotteryTime = r.GetLotteryRemainTime()
	r.BroadcastProto(result, 0)
	log.T("所有人发送可以开始押注的广播借宿")

}

/**
	构造要住记录的key
	每一局游戏,玩家对应的押注的key是   zjhRoundNum-userId
 */
func GetReidsKeyBetRecord(zjhRoundNum string, userId uint32) string {
	userIdStr, _ := numUtils.Uint2String(userId)
	result := strings.Join([]string{zjhRoundNum, userIdStr}, "-")
	log.T("得到获取redis中用户的押注记录的key[%v]值:", result)
	return result
}

func GetTBetRecordByUserIdAnd(zjhRoundNum string, userId uint32) *bbproto.TBetRecord {
	conn := data.Data{}
	conn.Open("test")
	defer conn.Close()
	key := GetReidsKeyBetRecord(zjhRoundNum, userId)
	result := &bbproto.TBetRecord{}
	err := conn.GetObj(key, result)
	if err != nil {
		log.E("reids中没有找到user【%v】对应的押注记录", userId)
		return nil
	}

	if result.GetUserId() == 0 {
		return nil
	}
	return result
}

/**
	保存押注记录,目前是只保存到redis 还没有保存到数据库当中
 */
func SaveBetRecord(record *bbproto.TBetRecord) error {
	conn := data.Data{}
	conn.Open("test")
	defer conn.Close()
	key := GetReidsKeyBetRecord(record.GetZjhRoundNumber(), uint32(record.GetUserId()))
	conn.SetObj(key, record)
	return nil
}