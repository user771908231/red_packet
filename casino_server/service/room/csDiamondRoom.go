package room

import (
	"casino_server/common/log"
)

type CsDiamondRoom struct {
	*CsRoomSkeleton
}

func (r *CsDiamondRoom) init(g CSGame) error {
	log.T("这里是CsDiamondRoom.init()")
	return nil
}

//是否可以开始游戏...
func (r *CsDiamondRoom) time2start() bool {
	return true
}

func (r *CsDiamondRoom) start() error {
	log.T("这里是CsDiamondRoom.start()")
	return nil
}

//是否可以开始游戏...
func (r *CsDiamondRoom) time2stop() bool {
	return true
}
func (r *CsDiamondRoom) stop() error {
	log.T("这里是CsDiamondRoom.stop()")
	return nil
}




