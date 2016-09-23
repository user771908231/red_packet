package room

import (
	"casino_server/common/log"
	"casino_server/mode"
	"gopkg.in/mgo.v2/bson"
	"casino_server/conf/casinoConf"
	"casino_server/utils/jobUtils"
	"time"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/db"
	"gopkg.in/mgo.v2"
	"casino_server/common/Error"
	"github.com/name5566/leaf/gate"
	"errors"
	"sync/atomic"
)

type CsDiamondRoom struct {
	*CsRoomSkeleton
}

func NewCsDiamondRoom() CSRoom {
	room := new(CsDiamondRoom)
	room.CsRoomSkeleton = new(CsRoomSkeleton)
	return room
}

func (r *CsDiamondRoom) init(g CSGame) error {
	log.T("这里是CsDiamondRoom.init()")

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

	//3,保存游戏数据,1,保存数据到mongo,2,刷新redis中的信息
	saveData := &mode.T_cs_th_record{}
	saveData.Mid = bson.NewObjectId()
	saveData.Id = r.MatchId
	saveData.ReadyTime = r.ReadyTime
	saveData.BeginTime = r.BeginTime
	saveData.EndTime = r.EndTime
	saveData.Status = r.Status
	db.InsertMgoData(casinoConf.DBT_T_CS_TH_RECORD, saveData)

	RefreshRedisMatchList()        //这里刷新redis中的锦标赛数据

	return nil
}



//是否可以开始游戏...
func (r *CsDiamondRoom) time2start() bool {
	if r.GetGamingCount() >= CSTHGameRoomConfig.leastCount {
		log.T("游戏人书已经足够了，可以开始游戏了...")
		return true        //表示终止任务
	} else {
		log.T("锦标赛[%v]玩家数量[%v]不够[%v],暂时不开始游戏.", r.MatchId, r.GetGamingCount(), CSTHGameRoomConfig.leastCount)
		return false
	}

	return true
}

func (r *CsDiamondRoom) start() error {
	log.T("这里是CsDiamondRoom.start()")

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

	//通知desk开始desk.run
	r.BroadCastDeskRunGame()

	return nil
}

//是否可以开始游戏...
func (r *CsDiamondRoom) time2stop() bool {
	if r.IsOutofEndTime() || r.GetGamingCount() <= 1 {
		return true
	} else {
		log.T("锦标赛matchId[%v]没有结束比赛IsOutofEndTime[%v],GetGamingCount[%v]:", r.MatchId, r.IsOutofEndTime(), r.GetGamingCount())
		return false
	}
}
func (r *CsDiamondRoom) stop() error {
	log.T("这里是CsDiamondRoom.stop()")
	return nil
}

func (r *CsDiamondRoom) GetGamingCount() int32 {
	var count int32 = 0
	desks := r.ThDeskBuf
	for _, desk := range desks {
		if desk != nil {
			for _, user := range desk.Users {
				if user != nil && user.CSGamingStatus && !user.IsLeave && !user.IsBreak {
					count ++
				}
			}
		}
	}
	//log.T("获取锦标赛当前的游戏(user != nil && user.CSGamingStatus && !user.IsLeave && !user.IsBreak)人数[%v]", count)
	return count

}

func (r *CsDiamondRoom) IsOutofEndTime() bool {
	return r.EndTime.Before(time.Now())
}


//通过所有的desk可以开始游戏了
func (r *CsDiamondRoom) BroadCastDeskRunGame() {
	log.T("通知所有的desk开始游戏")
	for i := 0; i < len(r.ThDeskBuf); i++ {
		desk := r.ThDeskBuf[i]
		if desk != nil {
			go desk.Run()
		}
	}
}

//游戏大厅增加一个玩家
func (r *CsDiamondRoom) AddUser(userId uint32, a gate.Agent) (*ThDesk, error) {
	r.Lock()
	defer r.Unlock()
	log.T("userid【%v】进入锦标赛的游戏房间", userId)

	//这里需要判断锦标赛是否可以开始游戏
	e := r.CheckIntoRoom(r.MatchId)
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

func (t *CsDiamondRoom) AddCopyUser(user *ThUser) {
	t.UsersCopy[user.UserId] = user
}

func (r *CsDiamondRoom) AddUserRankInfo(userId uint32, matchId int32, balance int64) {
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

func (r *CsDiamondRoom) AddrankUserCount() {
	atomic.AddInt32(&r.rankUserCount, 1)        //在线人数增加一人
}

func (r *CsDiamondRoom) SubrankUserCount() {
	atomic.AddInt32(&r.rankUserCount, -1)        //在线人数减少一人

}

func (r *CsDiamondRoom) AddOnlineCount() {
	atomic.AddInt32(&r.onlineCount, 1)        //在线人数增加一人
}

func (r *CsDiamondRoom) SubOnlineCount() {
	atomic.AddInt32(&r.onlineCount, -1)        //在线人数减少一人

}

//新建一个desk并且加入 到锦标赛的房间
func (r *CsDiamondRoom) NewThdeskAndSave() *ThDesk {
	mydesk := NewThDesk()
	mydesk.MatchId = r.MatchId
	mydesk.InitRoomCoin = r.initRoomCoin
	mydesk.RebuyCountLimit = r.RebuyCountLimit
	mydesk.blindLevel = r.BlindLevel
	mydesk.RebuyBlindLevelLimit = r.RebuyBlindLevelLimit
	r.AddThDesk(mydesk)
	return mydesk
}

//增加一个thRoom
func (r *CsDiamondRoom) AddThDesk(throom *ThDesk) error {
	r.ThDeskBuf = append(r.ThDeskBuf, throom)
	return nil
}


//只有开始之后才能进入游戏房间
func (r *CsDiamondRoom) CheckIntoRoom(matchId int32) error {

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

func (r *CsDiamondRoom) IsRepeatIntoRoom(userId uint32, a gate.Agent) (*ThDesk, error) {
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

func (r *CsDiamondRoom) GetAbleIntoDesk() *ThDesk {
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

func (t *CsDiamondRoom) GetCopyUserById(userId uint32) *ThUser {
	return t.UsersCopy[userId]
}

func (r *CsDiamondRoom) GetDeskById(id int32) *ThDesk {
	var result *ThDesk = nil
	for i := 0; i < len(r.ThDeskBuf); i++ {
		if r.ThDeskBuf[i] != nil && r.ThDeskBuf[i].Id == id {
			result = r.ThDeskBuf[i]
			break
		}
	}
	return result
}




