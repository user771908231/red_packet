package doudizhu

import "sync"

type DdzUser struct {
	sync.Mutex
	*PDdzUser
}

//清楚session
func (u *DdzUser)ClearAgentGameData() {

}
