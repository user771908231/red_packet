package room

import (
	"casino_server/common/log"
	"casino_server/utils/jobUtils"
	"time"
	"errors"
)



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


