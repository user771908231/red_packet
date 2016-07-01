package room

import (
	"time"
	"casino_server/service/porkService"
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/utils/time"
)

func init() {
	OninitZjgRoom()
}


//config

var ZJH_BET_DURATION = time.Second * 10
var ZJH_LOTTERY_DURATION = time.Second *10


//-------------------------------------------------扎金花房间的状态----------------------------------------------
var ZJH_STATUS_BETING		int32	=		1		//押注中
var ZJH_STATUS_LOTTERING	int32	=		2		//等待开奖
var ZJH_STATUS_LOTTERIED 	int32	=		3		//已经开奖

var ZJHroom zjhRoom        //扎金花的房间

/**
扎金花的room
 */
type zjhRoom struct {
	room					//继承基本的room的方法
	BetStartTime		time.Time	//首次投注时间
	BetEndTime		time.Time	//结束投注时间
	LotteryTime		time.Time	// 开奖时间
	NextStartTime		time.Time	//下次开始时间
	Status			int32		//房间状态:押注中,开奖中,
	Jackpot			int64		//奖池大小
	zoneAmount		[]int32		//押注区A的押注
	BankerUserId		uint32		//庄家的信息
}

/**
初始化扎金花的房间
 */
func OninitZjgRoom(){
	log.T("初始化扎金花的房间")

	//初始化参数 agent 的集合
	ZJHroom.AgentMap = make(map[uint32] gate.Agent)
	ZJHroom.zoneAmount = make([]int32,4)

	//初始化开始押注时间,押注结束时间,开奖时间
	ZJHroom.iniTime(time.Now())
	ZJHroom.Status = ZJH_STATUS_BETING
	//启动扎金花房间的run任务
	ZJHroom.run()
}


func (r *zjhRoom) iniTime(t time.Time){
	//r.Lock()
	//defer r.Unlock()
	r.BetStartTime	=	t					//首次开始押注的时间
	r.BetEndTime	=	r.BetStartTime.Add(ZJH_BET_DURATION)	//首次押注结束的时间
	r.LotteryTime	=	r.BetEndTime.Add(ZJH_LOTTERY_DURATION)	//首次开奖的时间
	r.NextStartTime = 	r.LotteryTime.Add(ZJH_LOTTERY_DURATION)	//下次开始的时间
	//下次押注的时间

}

/**
扎金花需要的启动函数
启动函数实际上是一个定时任务
	1,判断当前是否为投注状态,如果是投注的状态则可以投注
	2,判断当前是否为开奖的时间,如果是则不能投注,等待开奖
 */
func (r *zjhRoom) run(){
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for t := range ticker.C {
			log.T("正在执行扎金花的逻辑当前状态:\n status[%v],\n betEndTime[%v],\nlotterTime[%v],\nnextStartTime[%v],\nbankerUserId[%v],\njackPort[%v],\nnow[%v]",r.Status,r.BetEndTime.Format(timeUtils.TIME_LAYOUT),r.LotteryTime.Format(timeUtils.TIME_LAYOUT),r.NextStartTime.Format(timeUtils.TIME_LAYOUT),r.BankerUserId,r.Jackpot,t.Format(timeUtils.TIME_LAYOUT))
			if len(r.AgentMap) < 1 {
				log.T("没有玩家进入游戏...continue")
				//如果没有玩家进入 需要重置房间的状态和时间
				r.iniTime(t)
				continue
			}

			//这里测试代码,打印当前正在游戏中的玩家
			log.T("当前玩家正在房间中")
			/* 使用 key 输出 map 值 */
			for a := range r.AgentMap {
				log.T("agent-userId:%v",a)
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
func (r *zjhRoom) Betable()bool{
	if r.Status == ZJH_STATUS_BETING && r.BetEndTime.After(time.Now()){
		return true
	}else{
		return false
	}
}

/**
这个是开奖的广播
 */
func (r *zjhRoom) lottery(t time.Time){
	r.Lock()
	defer r.Unlock()
	//如果当前时间已经过了开奖时间,并且现在的状态是开奖中,则重新设置状态,并且开奖
	if r.LotteryTime.Before(t) && r.Status == ZJH_STATUS_LOTTERING {
		log.T("-----------------------------------------开奖-----------------------------------------")
		//得到5副牌
		list := porkService.CreateZjhList()
		//需要伪造数据,并且发送给每个人
		var balance1 int32 =  77878
		var winAmount int32 =  666
		result := &bbproto.ZjhLottery{}
		zs := make([]*bbproto.ZjhPai,5)
		for i := 0; i < 5; i++ {
			zs[i] 	= porkService.ConstructZjhPai(list[i])
		}
		result.Zjhpai = zs
		result.Balance = &balance1
		result.WinAmount = &winAmount
		//开始广播消息
		r.BroadcastProto(result,0)
		r.Status = ZJH_STATUS_LOTTERIED			              //设置状态已经开奖
		log.T("---------------------------------------开奖结束-----------------------------------------")
	}
}

/**
	开始下一轮的游戏
 */
func (r *zjhRoom) next(t time.Time){
	r.Lock()
	defer r.Unlock()
	if t.After(r.NextStartTime) && r.Status == ZJH_STATUS_LOTTERIED {
		log.T("---------------------------------------初始化下一轮-----------------------------------------")
		//开奖已经结束了..可以重新开始
		r.iniTime(t)
		r.Status = ZJH_STATUS_BETING
		log.T("--------------------------------------初始化下一轮结束---------------------------------------")

	}

}


func (r *zjhRoom) betEnd(t time.Time){
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
func (r *zjhRoom) AddZoneAmount(d []int32){
	r.Lock()
	defer r.Unlock()
	for i := 0; i < len(d); i++ {
		r.zoneAmount[i] = d[i]
	}
	log.T("增加积分之后,目前房间的zoneAmount:%v",r.zoneAmount)

}

/**
	返回押注剩余的时间
 */
func (r *zjhRoom) GetBetRemainTime() *int32{
	var result int32
	if r.Status == ZJH_STATUS_BETING {
		now := time.Now()
		diff := r.BetEndTime.Sub(now)
		result = int32(diff.Seconds())
	}else{
		result = 0
	}
	return &result
}

/**
	返回开奖剩余的时间
 */
func (r *zjhRoom) GetLotteryRemainTime() *int32{
	var result int32
	if r.Status == ZJH_STATUS_BETING || r.Status == ZJH_STATUS_LOTTERING {
		//押注中的状态或者等待开奖的状态的时候,可以返回等待开奖的时间
		now := time.Now()
		diff := r.LotteryTime.Sub(now)
		result = int32(diff.Seconds())
	}else{
		result = 0
	}
	return &result

}

