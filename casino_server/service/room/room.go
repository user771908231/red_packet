package room

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/msg/bbprotogo"
	"time"
	"github.com/golang/protobuf/proto"
	"sync"
	"casino_server/service/porkService"
)

func init() {
	OninitSgjRoom()
	OninitZjgRoom()
}


//confi

var ZJH_BET_DURATION = time.Second * 10
var ZJH_LOTTERY_DURATION = time.Second *10



var SGJRoom sgjRoom        //水果机的房间
var ZJHroom zjhRoom        //扎金花的房间


//-------------------------------------------------扎金花房间的状态----------------------------------------------
var ZJH_STATUS_BETING		int32	=		1		//押注中
var ZJH_STATUS_LOTTERING	 int32	=		2		//等待开奖
var ZJH_STATUS_LOTTERIED 	int32	=		3		//已经开奖


/**
游戏房间
 */
type room struct {
	sync.Mutex
	Type int
	RoomId	int32				//房间号
	AgentMap map[uint32] gate.Agent
}

func (r *room) AddAgent(userId uint32,a gate.Agent){
	log.T("userId%v的agent放在CachOutRoom中管理\n",userId)
	r.AgentMap[userId] = a

	//打印出 增加连接之后,但当前房间里的连接
	for key := range r.AgentMap {
		log.Normal("当前存在的连接%v",key)
	}
}

func (r *room) RemoveAgent(userId uint32){
	delete(r.AgentMap,userId);
}

/**
	发送信息
 */

func (r *room) BroadcastMsg(roomId int32,msg string){
	log.Normal("给房间号%v发送信息%v",roomId,msg)
	/* 使用 key 输出 map 值 */
	for key := range r.AgentMap {
		log.Normal("开始给%v发送消息",key)

		//首先判断连接是否有断开
		a :=r.AgentMap[key]

		m := "服务器的消息"
		data := bbproto.RoomMsg{}
		data.RoomId = &roomId
		data.Msg    = &m
		a.WriteMsg(&data)
		log.Normal("给%v发送消息,发送完毕",key)
	}
}

/**
	给所有的人广播消息,ignoreUserId 的除外
		目前暂时没有实现这个功能
 */
func (r *room) BroadcastProto(p proto.Message,ignoreUserId int32){
	log.Normal("给每个房间发送proto 消息%v",p)
	for key := range r.AgentMap {
		log.Normal("开始给%v发送消息",key)
		//首先判断连接是否有断开
		a :=r.AgentMap[key]
		a.WriteMsg(p)
		log.Normal("给%v发送消息,发送完毕",key)
	}
}

/**

水果机的room
 */
type sgjRoom struct {
	room
}

/**
初始化水果机的房间
 */
func OninitSgjRoom(){
	log.T("初始化水果机的房间")
	SGJRoom.AgentMap = make(map[uint32] gate.Agent)
}



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
	Jackpot			int32		//奖池大小
}

/**
初始化扎金花的房间
 */
func OninitZjgRoom(){
	log.T("初始化扎金花的房间")

	//初始化参数 agent 的集合
	ZJHroom.AgentMap = make(map[uint32] gate.Agent)

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
			log.T("正在执行扎金花的逻辑...")
			if len(r.AgentMap) < 1 {
				log.T("没有玩家进入游戏...continue")
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


	//打印测试信息
	log.T("betEndTime",r.BetEndTime.String())
	log.T("lotteryTime",r.LotteryTime.String())
	log.T("now",t.String())
	log.T("",r.LotteryTime.Before(t))
	log.T("status",r.Status)

	//如果当前时间已经过了开奖时间,并且现在的状态是开奖中,则重新设置状态,并且开奖
	if r.LotteryTime.Before(t) && r.Status == ZJH_STATUS_LOTTERING {
		log.T("开奖。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。")

		r.Status = ZJH_STATUS_LOTTERIED                //设置状态已经开奖
		//需要重新设置下一轮的时间
		r.iniTime(t)

		//得到5副牌
		list := porkService.CreateZjhList()

		//需要伪造数据,并且发送给每个人
		var balance1 int32 =  77878
		var winAmount int32 =  666
		result := &bbproto.ZjhLottery{}
		result.Pbank 	= porkService.ConstructZjhPai(list[0])
		result.Pa 	= porkService.ConstructZjhPai(list[0])
		result.Pb 	= porkService.ConstructZjhPai(list[0])
		result.Pc 	= porkService.ConstructZjhPai(list[0])
		result.Pd 	= porkService.ConstructZjhPai(list[0])

		result.Balance = &balance1
		result.WinAmount = &winAmount
		//开始广播消息
		r.BroadcastProto(result,0)
	}

}


func (r *zjhRoom) betEnd(t time.Time){
	r.Lock()
	defer r.Unlock()

	if r.BetEndTime.Before(t) && r.Status == ZJH_STATUS_BETING {
		log.T("投注结束。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。")
		r.Status = ZJH_STATUS_LOTTERING
	}

}




