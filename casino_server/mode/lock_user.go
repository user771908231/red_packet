package mode

import "sync"

//User相关的锁
type LockUser struct {
	sync.Mutex
	UserId uint32
}
