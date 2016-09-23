package room

import (
	"casino_server/common/log"
	"casino_server/utils/jobUtils"
	"time"
	"errors"
	"sync"
	"casino_server/msg/bbprotogo"
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

////////////////////////////////////////////////////////game///////////////////////////////////////////////////////////

//游戏的接口
type CSGame interface {
	run(CSRoom) error
	stop() error
}

/*
	锦标赛:
		1,开始条件
		2,结束条件

*/

//游戏的骨架
type CsGameSkeleton struct {
	roomBuf map[int32]CSRoom //roomBuf
}

//实现csGame.run()方法
func (g *CsGameSkeleton) run(room CSRoom) error {
	log.T("这里是CsGameSkeleton.run()")
	//创建一个房间，并且开始游戏...
	if room == nil {
		log.E("创建钻石锦标赛的时候room==nil")
		return errors.New("创建钻石锦标赛的房间出错")
	} else {
		//初始化房间
		room.init(g)

		//判断是否开始
		jobUtils.DoAsynJob(time.Second * 2, func() bool {
			//判断人数是否足够
			if room.time2start() {
				room.start()
				return true
			} else {
				return false
			}
		})

		//判断是否结束
		jobUtils.DoAsynJob(time.Second * 2, func() bool {
			//判断人数是否足够
			if room.time2stop() {
				room.stop()
				log.T("重新开始一局游戏...")
				time.Sleep(CSTHGameRoomConfig.nextRunDuration)        //开始下一场的延时
				g.run(NewCsDiamondRoom())      //重新开始一局游戏...
				return true
			} else {
				return false
			}
		})

		//把房间增加到room buf中
		g.addRoom(room)
		return nil

	}
	return nil
}

//停止游戏
func (g *CsGameSkeleton) stop() error {
	for key := range g.roomBuf {
		room := g.roomBuf[key]
		room.stop()
	}
	return nil
}

//
func (g *CsGameSkeleton) addRoom(r CSRoom) error {
	return nil
}





