package doudizhu

import "sync"

type DdzUser struct {
	sync.Mutex
	*DdzUser
}

//清楚session
func (u *DdzUser)ClearAgentGameData() {

}
