package room

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"time"
	"casino_server/utils/db"
	"casino_server/conf/casinoConf"
	"casino_server/mode"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"casino_server/service/CSTHService"
	"sync/atomic"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/jobUtils"
	"sort"
	"casino_server/service/userService"
	"casino_server/conf/intCons"
)

var ChampionshipRoom CSThGameRoom        //锦标赛的房间


func init() {
	ChampionshipRoom.OnInitConfig()
	ChampionshipRoom.OnInit()        //初始化,开始运行
	//ChampionshipRoom.begin()
}

//对竞标赛的配置
var CSTHGameRoomConfig struct {
	gameDuration      time.Duration //一场比赛的时间周期
	checkDuration     time.Duration //检测的时间周期
	leastCount        int32         //游戏开始的最少人数
	nextRunDuration   time.Duration //开始下一场的间隔
	riseBlindDuration time.Duration //生盲的时间间隔
	blinds            []int64       //盲注
	initRoomCoin      int64         //初始的带入金额
	deskMaxUserCount  int32         //最多多少人
	roomMaxUserCount  int32         //room里最多能有多少人
	RebuyCountLimit   int32         //重构次数的限制
	quotaLimit        int32         //名额的限制
}

//对配置对象进行配置,以后可以从配置文件读取
func (r *CSThGameRoom) OnInitConfig() {
	log.T("初始化csthgameroom.config")
	CSTHGameRoomConfig.gameDuration = time.Second * 60 * 20                //游戏是20分钟异常
	CSTHGameRoomConfig.checkDuration = time.Second * 10
	CSTHGameRoomConfig.leastCount = 3; //最少要20人才可以开始游戏
	CSTHGameRoomConfig.nextRunDuration = time.Second * 60 * 1        //1 分钟之后开始下一场
	CSTHGameRoomConfig.riseBlindDuration = time.Second * 150        //每150秒生一次忙
	CSTHGameRoomConfig.blinds = []int64{5, 10, 20, 40, 80, 160, 320, 640, 1280, 2000, 10000, 100000, 1000000}
	CSTHGameRoomConfig.initRoomCoin = 1000;
	CSTHGameRoomConfig.deskMaxUserCount = 9
	CSTHGameRoomConfig.RebuyCountLimit = 5                //最多重构5次
	CSTHGameRoomConfig.quotaLimit = 2                        //能得到奖励的人
}


//锦标赛的状态
var CSTHGAMEROOM_STATUS_STOP int32 = 1;

var CSTHGAMEROOM_STATUS_READY int32 = 2;

var CSTHGAMEROOM_STATUS_RUN int32 = 3;

var CSTHGAMEROOM_STATUS_LOTTERY int32 = 4;


//锦标赛
type CSThGameRoom struct {
	ThGameRoom
	matchId         int32                   //比赛内容
	beginTime       time.Time               //游戏开始的时间
	endTime         time.Time               //游戏结束的时间
	gameDuration    time.Duration           //游戏的时长
	rankUserCount   int32                   //游戏总人数
	onlineCount     int32                   //总的在线人数
	gamingUserCount int32                   //游戏总的人数--正在玩,没有输掉比赛的
	status          int32                   //锦标赛的状态
	rankInfo        []*bbproto.CsThRankInfo //排名信息
	blindLevel      int32                   //盲注的等级
	initRoomCoin    int64                   //房间默认的带入金额
	UsersCopy       map[uint32]*ThUser      //这里是所有玩家信息的一份拷贝,只有当用户放弃之后,才会删除用户
}

func (r *CSThGameRoom) OnInit() {
	log.T("初始化锦标赛的房间.oninit()")
	r.SmallBlindCoin = CSTHGameRoomConfig.blinds[r.blindLevel];
	r.initRoomCoin = CSTHGameRoomConfig.initRoomCoin
	r.ThRoomSeatMax = CSTHGameRoomConfig.deskMaxUserCount
	r.RebuyCountLimit = CSTHGameRoomConfig.RebuyCountLimit                        //重购的次数
	r.UsersCopy = make(map[uint32]*ThUser, )

}

//判断当前时间是否已经超过了endtime
func (r *CSThGameRoom) IsOutofEndTime() bool {
	return r.endTime.Before(time.Now())
}

//只有开始之后才能进入游戏房间
func (r *CSThGameRoom) IsCanIntoRoom() bool {
	//时间过了不能进入
	if r.status == CSTHGAMEROOM_STATUS_RUN && r.IsOutofEndTime() {
		log.T("进入锦标赛的游戏房间失败,因为time.mow[].after (r.endTime[%v])", r.endTime)
		return false
	}

	//游戏开始之后,用户只剩10人不能进入游戏 todo 这里的10人需要放置在配置文件中
	if r.status == CSTHGAMEROOM_STATUS_RUN && r.gamingUserCount <= CSTHGameRoomConfig.quotaLimit {
		log.T("因为竞标赛已经是run的状态,并且游戏中的人数小于10,所以不能开始")
		return false
	}

	//以上情况都不满足的时候,表示可以进入游戏房间
	return true
}


//开始游戏
/**
	锦标赛的逻辑
	1,开始场次,这里的开始只是有这个场次,但是游戏还没有真正的开始,只有满足(人数足够的时候)游戏才真正的开始
	1,开始游戏,通过每个房间,游戏可以开始了,进行前注,盲注,发牌...
 */

func (r *CSThGameRoom) Begin() {
	//开始一局游戏的时候,生成一个matchId
	r.matchId, _ = db.GetNextSeq(casinoConf.DBT_T_CS_TH_RECORD)        //生成游戏的matchId
	r.status = CSTHGAMEROOM_STATUS_READY

	log.T("开始锦标赛的游戏matchId[%v]", r.matchId)
	//判断是否可以开始run

	jobUtils.DoAsynJob(CSTHGameRoomConfig.checkDuration, func() bool {
		//判断人数是否足够
		if r.gamingUserCount >= CSTHGameRoomConfig.leastCount {
			//开始游戏
			r.Run()
			r.BroadCastDeskRunGame()
			return true        //表示终止任务
		} else {
			log.T("锦标赛[%v]玩家数量[%v]不够[%v],暂时不开始游戏.", r.matchId, r.gamingUserCount, CSTHGameRoomConfig.leastCount)
			return false
		}
	})
}

//run游戏房间
func (r *CSThGameRoom) Run() {
	log.T("锦标赛游戏开始...")

	//设置room属性
	r.beginTime = time.Now()
	r.endTime = r.beginTime.Add(CSTHGameRoomConfig.gameDuration)       //一局游戏的时间是20分钟
	r.status = CSTHGAMEROOM_STATUS_RUN        //竞标赛当前的状态

	//通知desk开始游戏

	//保存游戏数据,1,保存数据到mongo,2,刷新redis中的信息
	saveData := &mode.T_cs_th_record{}
	saveData.Mid = bson.NewObjectId()
	saveData.Id = r.matchId
	saveData.BeginTime = r.beginTime
	saveData.EndTime = r.endTime
	db.InsertMgoData(casinoConf.DBT_T_CS_TH_RECORD, saveData)
	CSTHService.RefreshRedisMatchList()        //这里刷新redis中的锦标赛数据

	//通知desk游戏开始


	//这里定义一个计时器,每十秒钟检测一次游戏
	jobUtils.DoAsynJob(CSTHGameRoomConfig.checkDuration, func() bool {
		//log.T("开始time[%v]检测锦标赛matchId[%v]有没有结束...", timeNow, r.matchId)
		if r.checkEnd() {
			//重新开始
			time.Sleep(CSTHGameRoomConfig.nextRunDuration)        //开始下一场的延时
			go r.Run()
			return true
		}
		return false
	})

	//这里需要做生盲的逻辑
	jobUtils.DoAsynJob(CSTHGameRoomConfig.riseBlindDuration, func() bool {
		//开始生盲注
		log.T("锦标赛[%v]开始生盲", r.matchId)
		if r.blindLevel == int32(len(CSTHGameRoomConfig.blinds) - 1) {
			log.T("由于锦标赛[%v]的盲注已经达到了最大的级别[%v],所以不生了", r.matchId, r.SmallBlindCoin)
			return true
		} else {
			r.blindLevel ++
			r.SmallBlindCoin = CSTHGameRoomConfig.blinds[r.blindLevel]
			log.T("由于锦标赛[%v]的盲注达到了级别[%v],盲注【%v】", r.matchId, r.blindLevel, r.SmallBlindCoin)

			return false
		}
	})
}

//通过所有的desk可以开始游戏了
func (r *CSThGameRoom) BroadCastDeskRunGame() {
	log.T("通知所有的desk开始游戏")
	for i := 0; i < len(r.ThDeskBuf); i++ {
		desk := r.ThDeskBuf[i]
		if desk != nil {
			go desk.Run()
		}
	}
}

func (r *CSThGameRoom) AddrankUserCount() {
	atomic.AddInt32(&r.rankUserCount, 1)        //在线人数增加一人
}

func (r *CSThGameRoom) SubrankUserCount() {
	atomic.AddInt32(&r.rankUserCount, -1)        //在线人数减少一人

}

func (r *CSThGameRoom) AddOnlineCount() {
	atomic.AddInt32(&r.onlineCount, 1)        //在线人数增加一人
}

func (r *CSThGameRoom) SubOnlineCount() {
	atomic.AddInt32(&r.onlineCount, -1)        //在线人数减少一人

}

//检测结束
func (r *CSThGameRoom) checkEnd() bool {
	//如果时间已经过了,并且所有桌子的状态都是已经停止游戏,那么表示这一局结束,为什么是所有的桌子?因为有可能时间到了,有很多桌子还在游戏中
	if r.IsOutofEndTime() && r.allStop() {
		//结算本局
		log.T("锦标赛matchid[%v]已经结束.现在开始保存数据", r.matchId)
		r.status = CSTHGAMEROOM_STATUS_STOP
		//这里需要保存每一个人锦标赛的结果信息
		log.T("保存每一个人竞标赛的信息")

		return true
	} else {
		return false
	}

}


//判断是否所有的desk停止游戏
//如果没有desk 是代表停止游戏还是游戏未开始?
func (r *CSThGameRoom) allStop() bool {
	result := true
	for i := 0; i < len(r.ThDeskBuf); i++ {
		desk := r.ThDeskBuf[i]
		if desk != nil && desk.CStatus != CSTH_DESK_STATUS_STOP {
			result = false
			break
		}
	}
	return result

}


//本场锦标赛 结束的处理
func (r *CSThGameRoom) End() {
	log.T("锦标赛游戏结束")
}


//锦标赛: 判断锦标赛是否重新进入房间
/**
	锦标赛的特殊性:锦标赛有可能没了桌子,但是人也是在游戏中的,所以不能通过agent来寻找桌子
	1,在csthroom的buf中来寻找thusers,找到之后看其状态

 */
func (r *CSThGameRoom) IsRepeatIntoRoom(userId uint32, a gate.Agent) (*ThDesk,error) {

	user := r.GetCopyUserById(userId)
	if user == nil {
		//表示没有进入过锦标赛
		return nil,nil
	}

	if user.IsLeave {
		//表示用户已经离开,不能进入游戏
		return nil,errors.New("用户已经离开了")
	}

	log.T("用户【%v】断线重连...",userId)
	return r.GetDeskById(user.deskId)
}


//游戏大厅增加一个玩家
func (r *CSThGameRoom) AddUser(userId uint32, a gate.Agent) (*ThDesk, error) {
	r.Lock()
	defer r.Unlock()
	log.T("userid【%v】进入锦标赛的游戏房间", userId)

	//这里需要判断锦标赛是否可以开始游戏
	if !r.IsCanIntoRoom() {
		log.T("用户[%v]进入锦标赛的房间失败,因为游戏还没有开始", userId)
		return nil, errors.New("游戏还没有开始")
	}

	var mydesk *ThDesk = nil                //为用户找到的desk
	//1,判断用户是否已经在房间里了,如果是在房间里,那么替换现有的agent
	//重新进入房间,如果是锦标赛,那么只有断线重连,如果是
	mydesk = r.IsRepeatIntoRoom(userId, a)
	if mydesk != nil {
		return mydesk, nil
	}

	//2,找到可以进入游戏的桌子
	mydesk = r.GetAbleintoDesk()        //找到可以进入的桌子,如果没有找到合适的desk,则新生成一个并且返回
	if mydesk == nil {
		//如果没有找到合适的桌子,那么新生成一个并且返回
		mydesk = r.NewThdeskAndSave()
	}

	//3,进入房间,竞标赛进入房间的时候,默认就是准备的状态
	user, err := mydesk.AddThUser(userId, TH_USER_STATUS_READY, a)
	if err != nil {
		log.E("用户上德州扑克的桌子 失败...")
		return nil, err
	}

	r.AddOnlineCount()        //在线用户增加1
	r.AddrankUserCount()
	r.AddgamingUserCount()    //游戏玩家数量+1
	r.AddUserRankInfo(user.UserId, user.MatchId, user.RoomCoin)
	r.AddCopyUser(user)	//用户列表总增加一个用户

	mydesk.LogString()        //打印当前房间的信息
	return mydesk, nil
}

func (r *CSThGameRoom) NewThdeskAndSave() *ThDesk {
	mydesk := NewThDesk()
	mydesk.MatchId = r.matchId
	mydesk.InitRoomCoin = r.initRoomCoin
	mydesk.RebuyCountLimit = r.RebuyCountLimit
	r.AddThDesk(mydesk)
	return mydesk
}

func (r *CSThGameRoom) GetAbleintoDesk() *ThDesk {
	var mydesk *ThDesk = nil
	//2,查询哪个德州的房间缺人:循环每个德州的房间,然后查询哪个房间缺人
	for deskIndex := 0; deskIndex < len(r.ThDeskBuf); deskIndex++ {
		log.T("查找竞标赛index=[%v]的房间", deskIndex)
		tempDesk := r.ThDeskBuf[deskIndex]
		if tempDesk == nil {
			log.E("找到房间为nil,出错")
			break
		}
		log.T("查找竞标赛index=[%v]的房间:tempDesk.UserCount[%v],r.ThRoomSeatMax", tempDesk.UserCount, r.ThRoomSeatMax)
		if tempDesk.UserCount < r.ThRoomSeatMax {
			mydesk = tempDesk        //通过roomId找到德州的room
			break;
		}
	}
	return mydesk

}

//是否可以进行下把游戏
func (r *CSThGameRoom) CanNextDeskRun() bool {
	//如果当前时间已经在结束时间之后,那么本局游戏结束
	if r.IsOutofEndTime() {
		log.T("游戏时间已经到了,不能开始游戏")
		return false
	}

	//如果锦标赛不是在 run 的状态,则desk不能开始
	if r.status != CSTHGAMEROOM_STATUS_RUN {
		return false
	}

	return true

}

//查找一个人的赛况
func (r *CSThGameRoom) GetRankInfo(userId uint32) *bbproto.CsThRankInfo {
	for i := 0; i < len(r.rankInfo); i++ {
		r := r.rankInfo[i]
		if r != nil && r.GetUserId() == userId {
			return r
		}
	}
	return nil
}

//更新用户的排名信息
func (r *CSThGameRoom) UpdateUserRankInfo(userId uint32, matchId int32, balance int64) error {
	info := r.GetRankInfo(userId)
	if info == nil {
		return errors.New("没有找到排名信息")
	}

	*info.Balance = time.Now().UnixNano()        //结束的时间
	*info.Balance = balance                        //余额

	//更新之后,保存数据到redis
	return nil
}

func (r *CSThGameRoom) AddUserRankInfo(userId uint32, matchId int32, balance int64) {
	//增加一个用户的rankinfo信息
	rank := &bbproto.CsThRankInfo{}
	rank.UserId = new(uint32)
	rank.MatchId = new(int32)
	rank.Balance = new(int64)
	rank.EndTime = new(int64)

	*rank.UserId = userId
	*rank.MatchId = matchId
	*rank.Balance = balance
	*rank.EndTime = time.Now().UnixNano()

}



//------------------------------------------------------关于排名的排序-------------------------------------------------
type RankList []*bbproto.CsThRankInfo

func ( list RankList) Len() int {
	return len(list)
}

//由大到小的排序
func ( list RankList) Less(i, j int) bool {
	if list[i].GetBalance() > list[j].GetBalance() {
		return false
	} else if list[i].GetBalance() == list[j].GetBalance() {
		if list[i].GetEndTime() > list[i].GetEndTime() {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

//交换函数
func ( list RankList) Swap(i, j int) {
	var temp *bbproto.CsThRankInfo = list[i]
	list[i] = list[j]
	list[j] = temp
}

//更具用户信息得到排名
func (r *CSThGameRoom) GetRankByuserId(userId uint32) int32 {
	var tempList RankList = make([]*bbproto.CsThRankInfo, len(r.rankInfo))
	copy(tempList, r.rankInfo)
	sort.Sort(tempList)                //开始排序

	index := 0
	for i := 0; i < len(tempList); i++ {
		info := tempList[i]
		if info != nil && info.GetUserId() == userId {
			index = i
			break
		}

	}

	rank := len(tempList) - index
	log.T("查询用户[%v]的锦标赛rank排名是[%v]", userId, rank)
	return int32(rank)
}
//------------------------------------------------------关于排名的排序-end---------------------------------------------

func (r *CSThGameRoom) AddgamingUserCount() {
	atomic.AddInt32(&r.gamingUserCount, 1)
}

func (r *CSThGameRoom) SubgamingUserCount() {
	atomic.AddInt32(&r.gamingUserCount, -1)

}

//解散锦标赛的房间,并且保留需要继续游戏的user
func (r *CSThGameRoom) DissolveDesk(desk *ThDesk, reserveUser *ThUser) error {
	//解散房间,给每个人发送解散房间的广播,并且删除房间
	log.T("锦标赛开始解散desk[%v]", desk.Id)
	result := &bbproto.Game_AckDissolveDesk{}
	result.Result = new(int32)
	result.UserId = new(uint32)
	result.DeskId = new(int32)
	result.PassWord = new(string)

	//1,找到桌子,并且判断是否能够解散
	if desk == nil {
		*result.Result = intCons.ACK_RESULT_ERROR
		reserveUser.WriteMsg(result)
		return errors.New("房间已经解散了")
	}

	if desk.Status != TH_DESK_STATUS_STOP {
		*result.Result = intCons.ACK_RESULT_ERROR
		reserveUser.WriteMsg(result)
		return errors.New("房间正在游戏中,不能解散")
	}

	//2,解散
	desk.clearAgentData(reserveUser.UserId)                //清空其他人的回话信息
	ChampionshipRoom.RmThDesk(desk)

	//3,发送解散的广播
	*result.Result = intCons.ACK_RESULT_SUCC
	*result.UserId = desk.DeskOwner
	*result.PassWord = desk.RoomKey
	desk.THBroadcastProtoAll(result)
	return nil
}

//用户加入到其他的房间
func (t *CSThGameRoom) Join(user *ThUser) error {
	log.T("用户[%v]加入到其他的房间", user.UserId)
	if user == nil {
		log.T("用户为nil,不能加入其他的房间")
		return errors.New("系统错误")
	}

	return nil
}

func (t *CSThGameRoom) RmCopyUser(userId uint32) {
	delete(t.UsersCopy, userId)
}

func (t *CSThGameRoom) AddCopyUser(user *ThUser) {
	t.UsersCopy[user.UserId] = user
}

//得到buf中的thusr
func (t *CSThGameRoom) GetCopyUserById(userId uint32) *ThUser{
	return  t.UsersCopy[userId]
}
