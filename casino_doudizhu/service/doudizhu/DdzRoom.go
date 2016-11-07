package doudizhu

import "sync"


//斗地主的房间
type  DdzRoom struct {
	sync.Mutex
}

func (room *DdzRoom) CreateDesk() *DdzDesk {
	//创建一个desk

	//1,得到一个key
	key := room.NewRoomKey()

	//2, newDesk and 赋值
	desk := NewDdzDesk()
	desk.key = key

	return desk

}

//得到一个roomKey
func (room *DdzRoom) NewRoomKey() string {
	return ""
}
