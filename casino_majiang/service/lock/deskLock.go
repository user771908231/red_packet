package lock

import "sync"

var deskLockMap map[int32]*sync.Mutex = make(map[int32]*sync.Mutex)

//新建一个锁
func NewDeskLock(deskId int32) {
	lock := &sync.Mutex{}
	deskLockMap[deskId] = lock
}

//得到锁
func GetDeskLock(deskId int32) *sync.Mutex {
	return deskLockMap[deskId]
}

//删除锁
func delDeskLock(deskId int32) {
	delete(deskLockMap, deskId)
}