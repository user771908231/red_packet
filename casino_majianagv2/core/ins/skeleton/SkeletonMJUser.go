package skeleton

import (
	"casino_majiang/service/majiang"
	"casino_majianagv2/core/data"
	"time"
	"casino_majianagv2/core/api"
	"casino_common/common/log"
	"errors"
	"sync/atomic"
)

type SkeletonMJUser struct {
	desk       api.MjDesk
	status     *data.MjUserStatus
	userId     uint32
	readyTimer *time.Timer
	Bill       *majiang.Bill
	UserData   *data.MJUserData
}

//初始化一个user骨架
func NewSkeleconMJUser(userId uint32) *SkeletonMJUser {
	return nil
}

func (user *SkeletonMJUser) Ready() {
	//设置为准备的状态,并且停止准备计时器
	user.status.SetStatus(majiang.MJUSER_STATUS_READY)
	user.status.Ready = true
	if user.readyTimer != nil {
		user.readyTimer.Stop()
		user.readyTimer = nil
	}

}




