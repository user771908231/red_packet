package coinPlay

import (
	"casino_majianagv2/core/ins/skeleton"
	"casino_common/common/log"
)

type CMJUser struct {
	*skeleton.SkeletonMJUser
}

//todo
func (u *CMJUser) GetCoin() int64 {
	return 0
}

func (u *CMJUser) Ready() error {
	//判断金币是否足够
	if u.GetCoin() < u.GetDesk().GetMJConfig().CreateFee {
		//金币不够则准备失败，返回金币不足的错误
		log.E("玩家[%v]准备失败，金币[%v]不够..", u.GetUserId(), u.GetCoin())
		return skeleton.ERR_READY_COIN_INSUFFICIENT
	}
	//设置为准备的状态,并且停止准备计时器
	u.SkeletonMJUser.Ready()
}
