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
	"sync/atomic"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/jobUtils"
	"sort"
	"casino_server/conf/intCons"
	"casino_server/common/Error"
	"gopkg.in/mgo.v2"
	"casino_server/utils/numUtils"
	"casino_server/service/userService"
	"casino_server/utils/redisUtils"
	"casino_server/utils/timeUtils"
)

var ChampionshipRoom CSThGameRoom        //锦标赛的房间


func init() {
	ChampionshipRoom.OnInitConfig()
	//ChampionshipRoom.begin()
}

//对竞标赛的配置
var CSTHGameRoomConfig struct {
	gameDuration         time.Duration //一场比赛的时间周期
	checkDuration        time.Duration //检测的时间周期
	leastCount           int32         //游戏开始的最少人数
	nextRunDuration      time.Duration //开始下一场的间隔
	riseBlindDuration    time.Duration //生盲的时间间隔
	Blinds               []int64       //盲注
	initRoomCoin         int64         //初始的带入金额
	deskMaxUserCount     int32         //最多多少人
	roomMaxUserCount     int32         //room里最多能有多少人
	RebuyCountLimit      int32         //重构次数的限制
	RebuyBlindLevelLimit int32         //可以购买的盲注级别
	quotaLimit           int32         //名额的限制
}

//对配置对象进行配置,以后可以从配置文件读取
func (r *CSThGameRoom) OnInitConfig() {
	log.T("初始化csthgameroom.config")
	CSTHGameRoomConfig.gameDuration = time.Second * 60 * 20                //游戏是20分钟异常
	CSTHGameRoomConfig.checkDuration = time.Second * 10
	CSTHGameRoomConfig.leastCount = 3; //最少要20人才可以开始游戏
	CSTHGameRoomConfig.nextRunDuration = time.Second * 60 * 1        //1 分钟之后开始下一场
	CSTHGameRoomConfig.riseBlindDuration = time.Second * 150        //每150秒生一次忙
	CSTHGameRoomConfig.Blinds = []int64{0, 25, 50, 75, 100, 200, 300, 400, 500, 1000, 2000, 3000, 4000, 5000, 10000 }
	CSTHGameRoomConfig.initRoomCoin = 1000;
	CSTHGameRoomConfig.deskMaxUserCount = 9
	CSTHGameRoomConfig.RebuyCountLimit = 5           //最多重构5次
	CSTHGameRoomConfig.quotaLimit = 2                //能得到奖励的人
	CSTHGameRoomConfig.RebuyBlindLevelLimit = 7      //7级盲注以前可以购买
	CSTHGameRoomConfig.roomMaxUserCount = 500        //最多500人玩
}


//锦标赛的状态
var CSTHGAMEROOM_STATUS_STOP int32 = 1; //竞标赛已经停止

var CSTHGAMEROOM_STATUS_READY int32 = 2;

var CSTHGAMEROOM_STATUS_RUN int32 = 3;

var CSTHGAMEROOM_STATUS_LOTTERY int32 = 4;


//锦标赛
type CSThGameRoom struct {
	ThGameRoom
	MatchId              int32                   //比赛内容
	ReadyTime            time.Time               //游戏开始准备的时间
	BeginTime            time.Time               //游戏开始的时间
	EndTime              time.Time               //游戏结束的时间
	gameDuration         time.Duration           //游戏的时长
	rankUserCount        int32                   //游戏总人数
	onlineCount          int32                   //总的在线人数
	Status               int32                   //锦标赛的状态
	RankInfo             []*bbproto.CsThRankInfo //排名信息
	BlindLevel           int32                   //盲注的等级
	initRoomCoin         int64                   //房间默认的带入金额
	UsersCopy            map[uint32]*ThUser      //这里是所有玩家信息的一份拷贝,只有当用户放弃之后,才会删除用户
	RebuyBlindLevelLimit int32                   //盲注可以购买的级别
}

func (r *CSThGameRoom) OnInit() {
	log.T("初始化锦标赛的房间.oninit()")
	r.BlindLevel = 0
	r.ThDeskBuf = nil
	r.SmallBlindCoin = CSTHGameRoomConfig.Blinds[r.BlindLevel];
	r.initRoomCoin = CSTHGameRoomConfig.initRoomCoin
	r.ThRoomSeatMax = CSTHGameRoomConfig.deskMaxUserCount
	r.RebuyCountLimit = CSTHGameRoomConfig.RebuyCountLimit                        //重购的次数
	r.RebuyBlindLevelLimit = CSTHGameRoomConfig.RebuyBlindLevelLimit
	r.UsersCopy = make(map[uint32]*ThUser, CSTHGameRoomConfig.roomMaxUserCount)
	r.MatchId, _ = db.GetNextSeq(casinoConf.DBT_T_CS_TH_RECORD)        //生成游戏的matchId
	r.Status = CSTHGAMEROOM_STATUS_READY
	r.ReadyTime = time.Now()
	r.RankInfo = make([]*bbproto.CsThRankInfo, 0)
}


//活的当前正在游戏中的人数
func (r *CSThGameRoom) GetGamingCount() int32 {
	var count int32 = 0
	desks := r.ThDeskBuf
	for _, desk := range desks {
		if desk != nil {
			for _, user := range desk.Users {
				if user != nil && user.CSGamingStatus {
					count ++
				}
			}

		}
	}
	return count

}

//判断当前时间是否已经超过了endtime
func (r *CSThGameRoom) IsOutofEndTime() bool {
	return r.EndTime.Before(time.Now())
}

//只有开始之后才能进入游戏房间
func (r *CSThGameRoom) CheckIntoRoom(matchId int32) error {

	if r.MatchId != matchId {
		log.T("进入锦标赛失败,游戏场次matchId[%v]不正确", matchId)
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_INTO_DESK_NOTFOUND), "游戏已经过期")
	}

	//时间过了不能进入
	if r.Status == CSTHGAMEROOM_STATUS_RUN && r.IsOutofEndTime() {
		log.T("进入锦标赛的游戏房间失败,因为time.mow[].after (r.endTime[%v])", r.EndTime)
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_INTO_DESK_NOTFOUND), "游戏已经过期")

	}

	//游戏开始之后,用户只剩10人不能进入游戏 todo 这里的10人需要放置在配置文件中
	if r.Status == CSTHGAMEROOM_STATUS_RUN && r.GetGamingCount() <= CSTHGameRoomConfig.quotaLimit {
		log.T("因为竞标赛已经是run的状态,并且游戏中的人数小于10,所以不能开始")
		return Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_INTO_DESK_NOTFOUND), "游戏已经过期")
	}

	//以上情况都不满足的时候,表示可以进入游戏房间
	return nil
}


//开始游戏
/**
	锦标赛的逻辑
	1,开始场次,这里的开始只是有这个场次,但是游戏还没有真正的开始,只有满足(人数足够的时候)游戏才真正的开始
	1,开始游戏,通过每个房间,游戏可以开始了,进行前注,盲注,发牌...
 */

func (r *CSThGameRoom) Begin() {
	r.OnInit()                //每次开始的时候做初始化
	//保存游戏数据,1,保存数据到mongo,2,刷新redis中的信息
	saveData := &mode.T_cs_th_record{}
	saveData.Mid = bson.NewObjectId()
	saveData.Id = r.MatchId
	saveData.ReadyTime = r.ReadyTime
	saveData.BeginTime = r.BeginTime
	saveData.EndTime = r.EndTime
	saveData.Status = r.Status
	db.InsertMgoData(casinoConf.DBT_T_CS_TH_RECORD, saveData)
	RefreshRedisMatchList()        //这里刷新redis中的锦标赛数据

	log.T("开始锦标赛的游戏matchId[%v]", r.MatchId)

	//判断是否可以开始run
	jobUtils.DoAsynJob(CSTHGameRoomConfig.checkDuration, func() bool {
		//判断人数是否足够
		if r.GetGamingCount() >= CSTHGameRoomConfig.leastCount {
			//开始游戏
			r.Run()

			//通知desk开始desk.run
			r.BroadCastDeskRunGame()
			return true        //表示终止任务
		} else {
			log.T("锦标赛[%v]玩家数量[%v]不够[%v],暂时不开始游戏.", r.MatchId, r.GetGamingCount(), CSTHGameRoomConfig.leastCount)
			return false
		}
	})
}


//run游戏房间
func (r *CSThGameRoom) Run() {
	log.T("锦标赛游戏开始...run()")

	//设置room属性
	r.Status = CSTHGAMEROOM_STATUS_RUN        //竞标赛当前的状态
	r.BeginTime = time.Now()
	r.EndTime = r.BeginTime.Add(CSTHGameRoomConfig.gameDuration)       //一局游戏的时间是20分钟

	//update 锦标赛的数据
	saveData := &mode.T_cs_th_record{}
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_CS_TH_RECORD).Find(bson.M{"id":r.MatchId}).One(saveData)
	})

	saveData.BeginTime = r.BeginTime
	saveData.EndTime = r.EndTime
	saveData.Status = r.Status
	db.UpdateMgoData(casinoConf.DBT_T_CS_TH_RECORD, saveData)


	//通知desk游戏开始
	//这里定义一个计时器,每十秒钟检测一次游戏
	jobUtils.DoAsynJob(CSTHGameRoomConfig.checkDuration, func() bool {
		//log.T("开始time[%v]检测锦标赛matchId[%v]有没有结束...", timeNow, r.matchId)
		if r.checkEnd() {
			//重新开始
			time.Sleep(CSTHGameRoomConfig.nextRunDuration)        //开始下一场的延时
			go r.Begin()
			return true
		}
		return false
	})

	//这里需要做生盲的逻辑
	jobUtils.DoAsynJob(CSTHGameRoomConfig.riseBlindDuration, func() bool {
		//1,如果游戏还没有开始,停止升盲注的任务
		if r.Status != CSTHGAMEROOM_STATUS_RUN {
			return true
		}

		//2,开始生盲注
		log.T("锦标赛[%v]开始生盲", r.MatchId)
		if r.BlindLevel == int32(len(CSTHGameRoomConfig.Blinds) - 1) {
			log.T("由于锦标赛[%v]的盲注已经达到了最大的级别[%v],所以不生了", r.MatchId, r.SmallBlindCoin)
			return true
		} else {
			r.BlindLevel ++
			r.SmallBlindCoin = CSTHGameRoomConfig.Blinds[r.BlindLevel]
			log.T("由于锦标赛[%v]的盲注达到了级别[%v],盲注【%v】", r.MatchId, r.BlindLevel, r.SmallBlindCoin)

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
	if r.IsOutofEndTime() || r.GetGamingCount() <= 1 {
		//结算本局
		log.T("锦标赛matchid[%v]已经结束.现在开始保存数据", r.MatchId)
		r.End()
		//这里需要保存每一个人锦标赛的结果信息
		log.T("保存每一个人竞标赛的信息")

		return true
	} else {
		return false
	}
}


//本场锦标赛 结束的处理
func (r *CSThGameRoom) End() {
	log.T("锦标赛游戏结束")
	//设置锦标赛的状态为结束,并且更新数据库数据
	//保存锦标赛的数据,玩家的游戏数据
	r.Status = CSTHGAMEROOM_STATUS_STOP
	r.RefreshRank()
	saveData := &mode.T_cs_th_record{}
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_CS_TH_RECORD).Find(bson.M{"Id":r.MatchId}).One(saveData)
	})
	//更新mongo中锦标赛的状态
	saveData.Status = r.Status

	//更新锦标赛的排名信息
	for _, rank := range r.RankInfo {
		bean := mode.T_cs_th_rank_bean{}
		bean.UserId = rank.GetUserId()
		bean.WinCoin = rank.GetBalance()
		saveData.Ranks = append(saveData.Ranks, bean)
	}
	//保存信息
	db.UpdateMgoData(casinoConf.DBT_T_CS_TH_RECORD, saveData)


	//给没有发送过游戏排名的玩家发送游戏排名
	for _, desk := range r.ThDeskBuf {
		if desk != nil {
			for _, user := range desk.Users {
				if user != nil && user.CSGamingStatus {
					user.CSGamingStatus = false; //设置游戏状态
					log.T("锦标赛[%v]结束,给用户[%v]发送游戏排名", r.MatchId, user.UserId)
					ret := bbproto.NewGame_TounamentPlayerRank()
					*ret.Message = "测试最终排名的信息"
					*ret.PlayerRank = ChampionshipRoom.GetRankByuserId(user.UserId)
					user.WriteMsg(ret)
				}
			}
		}
	}

	//给在desk上的人发送游戏结束的广播
	csUser := userService.GetUserById(r.RankInfo[0].GetUserId())
	gameOver := bbproto.NewGame_ChampionshipGameOver()
	*gameOver.Coin = r.RankInfo[0].GetBalance()
	*gameOver.UserName = csUser.GetNickName()
	*gameOver.HeadUrl = csUser.GetHeadUrl()
	for _, desk := range r.ThDeskBuf {
		if desk != nil {
			desk.THBroadcastProtoAll(csUser)
		}
	}

}

//刷新排名
func (r *CSThGameRoom) RefreshRank() {
	var tempList RankList = make([]*bbproto.CsThRankInfo, len(r.RankInfo))
	copy(tempList, r.RankInfo)
	sort.Sort(tempList)                //开始排序
	r.RankInfo = tempList
}


//锦标赛: 判断锦标赛是否重新进入房间
/**
	锦标赛的特殊性:锦标赛有可能没了桌子,但是人也是在游戏中的,所以不能通过agent来寻找桌子
	1,在csthroom的buf中来寻找thusers,找到之后看其状态

 */
func (r *CSThGameRoom) IsRepeatIntoRoom(userId uint32, a gate.Agent) (*ThDesk, error) {
	user := r.GetCopyUserById(userId)
	if user == nil {
		//表示没有进入过锦标赛
		return nil, nil
	}

	//表示用户已经离开,并且是自己放弃的比赛,那么不能重新进入游戏...
	if user.IsLeave && !user.CSGamingStatus {
		return nil, errors.New("用户已经放弃比赛了")
	}

	//设置用户的状态
	user.IsBreak = false
	user.IsLeave = false
	user.Agent = a
	user.GameStatus = TH_USER_GAME_STATUS_CHAMPIONSHIP
	user.UpdateAgentUserData()        //更新用户的session信息

	log.T("用户【%v】断线重连...", userId)
	desk := r.GetDeskById(user.deskId)
	return desk, nil
}


//游戏大厅增加一个玩家
func (r *CSThGameRoom) AddUser(userId uint32, matchId int32, a gate.Agent) (*ThDesk, error) {
	r.Lock()
	defer r.Unlock()
	log.T("userid【%v】进入锦标赛的游戏房间", userId)

	//这里需要判断锦标赛是否可以开始游戏
	e := r.CheckIntoRoom(matchId)
	if e != nil {
		return nil, Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_TOURNAMENT_CANNOT_JOIN), "锦标赛入场时间已过,不能加入")
	}

	//1,判断用户是否已经在房间里了,如果是在房间里,那么替换现有的agent
	//重新进入房间,如果是锦标赛,那么只有断线重连,如果是
	mydesk, err := r.IsRepeatIntoRoom(userId, a)
	if err != nil {
		//直接不能进入游戏
		return nil, err
	}

	if mydesk != nil {
		return mydesk, nil
	}

	//2,找到可以进入游戏的桌子
	mydesk = r.GetAbleIntoDesk()        //找到可以进入的桌子,如果没有找到合适的desk,则新生成一个并且返回
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

	//更新room的信息
	r.AddOnlineCount()        //在线用户增加1
	r.AddrankUserCount()
	r.AddUserRankInfo(user.UserId, mydesk.MatchId, user.RoomCoin)
	r.AddCopyUser(user)       //用户列表总增加一个用户

	mydesk.LogString()        //打印当前房间的信息
	return mydesk, nil
}

//desk和user
func (t *ThDesk) UpdateThdeskAndAllUser2redis() error {
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			u.Update2redis()
			u.UpdateAgentUserData()
		}
	}
	t.Update2redis()
	return nil
}

//删除redis中的数据
func (t *ThDesk) DelThdeskAndAllUserRedisCache() error {
	//1,删除缓存中的desk
	DelRedisThdesk(t.Id)

	//2,删除缓存中desk,gamenumber 相关的user
	for _, u := range t.Users {
		if u != nil {
			DelRedisThUser(t.Id, u.UserId)
		}
	}

	//3,删除running 中的desk的key
	RmRunningDesk(t)
	return nil
}


//删除所有当前用户的信息
func (t *ThDesk) DelAllUsersRedisCache() {
	for _, u := range t.Users {
		if u != nil {
			DelRedisThUser(t.Id, u.UserId)
		}
	}

}


//删除离开的用户,离开的用户,设置gameNumber为0
func (t *ThDesk) DelAllLeaveUsersRedisCache() {
	for _, u := range t.Users {
		if u != nil {
			DelRedisThUser(t.Id, u.UserId)
		}
	}
}


//把desk的信息和指定的user备份到redis
func (t *ThDesk) UpdateThdeskAndUser2redis(user  *ThUser) error {
	if user != nil {
		user.Update2redis()
		t.Update2redis()
		return nil
	} else {
		return errors.New("用户为nil,无法保存数据")
	}
}


//当桌子的信息有所改变的时候,需要调用这个方法把桌子的数据保存在redis中
// 保存thDesk 的时候,同时需要保存room的信息到redis ,目前只需要保存锦标赛 room 的信息

func (t *ThDesk) Update2redis() {
	UpdateTedisThDesk(t)
}

//新建一个desk并且加入 到锦标赛的房间
func (r *CSThGameRoom) NewThdeskAndSave() *ThDesk {
	mydesk := NewThDesk()
	mydesk.MatchId = r.MatchId
	mydesk.InitRoomCoin = r.initRoomCoin
	mydesk.RebuyCountLimit = r.RebuyCountLimit
	mydesk.blindLevel = r.BlindLevel
	mydesk.RebuyBlindLevelLimit = r.RebuyBlindLevelLimit
	r.AddThDesk(mydesk)
	return mydesk
}

func (r *CSThGameRoom) GetAbleIntoDesk() *ThDesk {
	var mydesk *ThDesk = nil
	//2,查询哪个德州的房间缺人:循环每个德州的房间,然后查询哪个房间缺人
	for deskIndex := 0; deskIndex < len(r.ThDeskBuf); deskIndex++ {
		log.T("查找竞标赛index=[%v]的房间", deskIndex)
		tempDesk := r.ThDeskBuf[deskIndex]
		if tempDesk == nil {
			log.E("找到房间为nil,出错")
			break
		}
		log.T("查找竞标赛index=[%v]的房间:tempDesk.UserCount[%v],r.ThRoomSeatMax", tempDesk.GetUserCount(), r.ThRoomSeatMax)
		if tempDesk.GetUserCount() < r.ThRoomSeatMax {
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
	if r.Status != CSTHGAMEROOM_STATUS_RUN {
		return false
	}

	return true

}

//查找一个人的赛况
func (r *CSThGameRoom) GetRankInfo(userId uint32) *bbproto.CsThRankInfo {
	for i := 0; i < len(r.RankInfo); i++ {
		r := r.RankInfo[i]
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

	r.RankInfo = append(r.RankInfo, rank)        //保存到rankInfo中去
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
	r.RefreshRank()
	index := 0
	for i := 0; i < len(r.RankInfo); i++ {
		info := r.RankInfo[i]
		if info != nil && info.GetUserId() == userId {
			index = i
			break
		}

	}

	rank := len(r.RankInfo) - index
	log.T("查询用户[%v]的锦标赛rank排名是[%v]", userId, rank)
	return int32(rank)
}
//------------------------------------------------------关于排名的排序-end---------------------------------------------

//解散锦标赛的房间
func (r *CSThGameRoom) DissolveDesk(desk *ThDesk) error {
	//解散房间,给每个人发送解散房间的广播,并且删除房间
	log.T("锦标赛开始解散desk[%v]", desk.Id)
	result := &bbproto.Game_AckDissolveDesk{}
	result.Result = new(int32)
	result.UserId = new(uint32)
	result.DeskId = new(int32)
	result.PassWord = new(string)

	//1,找到桌子,并且判断是否能够解散
	if desk == nil {
		return errors.New("房间已经解散了")
	}

	if desk.IsRun() || desk.IsLottery() {
		return errors.New("房间正在游戏中,不能解散")
	}

	//2,发送解散的广播
	*result.Result = intCons.ACK_RESULT_SUCC
	*result.UserId = desk.DeskOwner
	*result.PassWord = desk.RoomKey
	desk.THBroadcastProtoAll(result)

	//3,解散桌子...
	ChampionshipRoom.RmThDesk(desk)                //删除buf中的桌子
	return nil
}

func (t *CSThGameRoom) RmCopyUser(userId uint32) {
	delete(t.UsersCopy, userId)
}

func (t *CSThGameRoom) AddCopyUser(user *ThUser) {
	t.UsersCopy[user.UserId] = user
}

//得到buf中的thusr
func (t *CSThGameRoom) GetCopyUserById(userId uint32) *ThUser {
	return t.UsersCopy[userId]
}

//得到本场锦标赛盲注的信息
func (r *CSThGameRoom) GetGame_TounamentBlind() *bbproto.Game_TounamentBlind {
	ret := bbproto.NewGame_TounamentBlind()
	//得到盲注信息
	blindLevel := int32(0)
	for _, b := range CSTHGameRoomConfig.Blinds {
		if blindLevel >= 15 {
			//暂时只提供15级盲注
			break
		}
		blindLevel += 1
		bean := bbproto.NewGame_TounamentBlindBean()
		*bean.BlindLevel, _ = numUtils.Int2String(blindLevel)
		*bean.Ante, _ = numUtils.Int642String(CSTHGameRoomConfig.Blinds[blindLevel - 1]) //"前注"
		*bean.SmallBlind, _ = numUtils.Int642String(b) //+ "/" + umUtils.Int642String(b*2)
		*bean.CanRebuy = (blindLevel <= 7 )
		*bean.RaiseTime = "150秒"
		ret.Data = append(ret.Data, bean)
		log.T("GetBlind  >>> bean[%v]: ante:%v  canRebuy: %v", *bean.BlindLevel, bean.GetAnte(), *bean.CanRebuy)
	}
	return ret
}



//合并桌子,并且返回,合并之后的桌子,合并失败
func (r *CSThGameRoom) MergeDesk(mt *ThDesk) (*ThDesk, error) {
	r.Lock()
	defer r.Unlock()

	//1,判断mt 是否需要重组
	if mt.GetCSGamingUserCount() > 4 {
		return nil, errors.New("不需要重组,人数大于4人")
	}

	if !mt.IsStop() {
		return nil, errors.New("不能重组,desk的状态不是stop的状态")
	}

	//2,找到可以重组的desk
	var dt *ThDesk = nil
	for _, desk := range r.ThDeskBuf {
		if desk != nil && desk.Id != mt.Id && desk.GetUserCount() <= 4 {
			dt = desk
			break
		}
	}

	if dt == nil {
		return nil, errors.New("没有找到合适的desk,重组失败...")
	}


	//3,开始重组
	for _, user := range mt.Users {
		if user != nil && user.CSGamingStatus {
			//用户不为空,并且是锦标赛游戏状态中的
			user.deskId = dt.Id
			addError := dt.AddThuserBean(user)
			if addError != nil {
				log.E("用户假如房间失败...")
				//这里的处理是让用户退出房间,自己重新加入
			}

		}
	}

	dt.UpdateThdeskAndAllUser2redis()        //更新缓存中的数据
	dt.BroadGameInfo(0)                        //发送desk的info信息


	//4,解散之前的房间
	ChampionshipRoom.DissolveDesk(mt)
	return dt, nil
}


//刷新锦标赛的列表信息
func RefreshRedisMatchList() {
	//1,获取数据库中的近20场次的信息(通过时间来排序)
	data := []mode.T_cs_th_record{}
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_CS_TH_RECORD).Find(bson.M{"id":bson.M{"$gt":0}}).Sort("-id").Limit(10).All(&data)
	})

	//把得到的数据保存在数据库中
	if len(data) > 0 {
		//表示有数据,需要存储在redis中
		sdata := bbproto.NewGame_MatchList()
		*sdata.HelpMessage = "1、竞技场玩法类似传统的德州扑克锦标赛,多人角逐"
		for i := 0; i < len(data); i++ {
			d := data[i]
			sd := bbproto.NewGame_MatchItem()
			*sd.CostFee = d.CostFee
			idStr, _ := numUtils.Int2String(d.Id)
			*sd.Title = "神经德州赢红包大赛" + idStr
			*sd.Status = d.Status
			*sd.Type = d.GameType
			*sd.Time = timeUtils.Format(d.BeginTime)
			*sd.MatchId = d.Id

			//如果是真在run或者ready的状态则表示为游戏中
			if d.Status == CSTHGAMEROOM_STATUS_RUN || d.Status == CSTHGAMEROOM_STATUS_READY {
				*sd.Status = 1
			} else {
				d.Status = 2
			}

			sdata.Items = append(sdata.Items, sd)
		}

		//存储
		redisUtils.SetObj(GetMatchListRedisKey(), sdata)
	}

}

//返回存放竞标赛列表的redis-key
func GetMatchListRedisKey() string {
	return "game_match_list_redis_key"
}

//竞标赛列表主要存放在标 t_cs_th_record中,但是一般都是在redis中取,如果redis中没有再从数据库中取
func GetGameMatchList() *bbproto.Game_MatchList {
	data := redisUtils.GetObj(GetMatchListRedisKey(), &bbproto.Game_MatchList{})
	if data == nil {
		log.T("redis中没有找到竞标赛列表,需要在数据库中查找")
		RefreshRedisMatchList()
		result := bbproto.NewGame_MatchList()
		return result
	} else {
		//更新其状态
		return data.(*bbproto.Game_MatchList)
	}
}



