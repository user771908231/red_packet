package room

import (
	"sync"
	"time"
	"casino_server/msg/bbprotogo"
	"casino_server/common/Error"
	"sync/atomic"
	"casino_server/common/log"
	"github.com/name5566/leaf/gate"
	"errors"
)



//锦标赛的接口
type CSRoom interface {
	init(CSGame) error
	time2start() bool
	start() error
	time2stop() bool
	stop() error                //做停止的处理
}

//游戏的骨架
type CsRoomSkeleton struct {
	sync.Mutex
	RoomStatus           int32                   //游戏大厅的状态
	ThDeskBuf            []*ThDesk
	ThRoomSeatMax        int32                   //每个房间的座位数目
	ThRoomCount          int32                   //房间数目
	Id                   int32                   //房间的id
	SmallBlindCoin       int64                   //小盲注的金额
	RebuyCountLimit      int32                   //重购的次数限制
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

//备份的用户
func (t *CsRoomSkeleton) GetCopyUserById(userId uint32) *ThUser {
	return t.UsersCopy[userId]
}

//通过deskId 得到desk
func (r *CsRoomSkeleton) GetDeskById(id int32) *ThDesk {
	var result *ThDesk = nil
	for i := 0; i < len(r.ThDeskBuf); i++ {
		if r.ThDeskBuf[i] != nil && r.ThDeskBuf[i].Id == id {
			result = r.ThDeskBuf[i]
			break
		}
	}
	return result
}



//是否可以开始游戏...
func (r *CsRoomSkeleton) time2stop() bool {
	if r.IsOutofEndTime() || r.GetGamingCount() <= 1 {
		return true
	} else {
		log.T("锦标赛matchId[%v]没有结束比赛IsOutofEndTime[%v],GetGamingCount[%v]:", r.MatchId, r.IsOutofEndTime(), r.GetGamingCount())
		return false
	}
}
func (r *CsRoomSkeleton) stop() error {
	log.T("这里是CsDiamondRoom.stop()")
	return nil
}

func (r *CsRoomSkeleton) GetGamingCount() int32 {
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

func (r *CsRoomSkeleton) IsOutofEndTime() bool {
	return r.EndTime.Before(time.Now())
}


//通过所有的desk可以开始游戏了
func (r *CsRoomSkeleton) BroadCastDeskRunGame() {
	log.T("通知所有的desk开始游戏")
	for i := 0; i < len(r.ThDeskBuf); i++ {
		desk := r.ThDeskBuf[i]
		if desk != nil {
			go desk.Run()
		}
	}
}

//游戏大厅增加一个玩家
func (r *CsRoomSkeleton) AddUser(userId uint32, a gate.Agent) (*ThDesk, error) {
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

func (t *CsRoomSkeleton) AddCopyUser(user *ThUser) {
	t.UsersCopy[user.UserId] = user
}

func (r *CsRoomSkeleton) AddUserRankInfo(userId uint32, matchId int32, balance int64) {
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

func (r *CsRoomSkeleton) AddrankUserCount() {
	atomic.AddInt32(&r.rankUserCount, 1)        //在线人数增加一人
}

func (r *CsRoomSkeleton) SubrankUserCount() {
	atomic.AddInt32(&r.rankUserCount, -1)        //在线人数减少一人

}

func (r *CsRoomSkeleton) AddOnlineCount() {
	atomic.AddInt32(&r.onlineCount, 1)        //在线人数增加一人
}

func (r *CsRoomSkeleton) SubOnlineCount() {
	atomic.AddInt32(&r.onlineCount, -1)        //在线人数减少一人

}

//新建一个desk并且加入 到锦标赛的房间
func (r *CsRoomSkeleton) NewThdeskAndSave() *ThDesk {
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
func (r *CsRoomSkeleton) AddThDesk(throom *ThDesk) error {
	r.ThDeskBuf = append(r.ThDeskBuf, throom)
	return nil
}


//只有开始之后才能进入游戏房间
func (r *CsRoomSkeleton) CheckIntoRoom(matchId int32) error {

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

func (r *CsRoomSkeleton) IsRepeatIntoRoom(userId uint32, a gate.Agent) (*ThDesk, error) {
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

func (r *CsRoomSkeleton) GetAbleIntoDesk() *ThDesk {
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






