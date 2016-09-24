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
