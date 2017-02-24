package coinPlay

import (
	"casino_common/common/log"
)

//准备
func (d *CMJDesk) Ready(userId uint32) error {
	log.T("锁日志: %v CMJDesk.Ready(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v CMJDesk.Ready(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	err := d.SkeletonMJDesk.Ready(userId) //金币场调用骨架的准备
	if err != nil {
		log.E("%v玩家%v准备失败...", d.DlogDes(), userId)
		return err
	}

	d.initEnterTimer() //房间进入一个人之后开始计划添加机器人
	d.begin()
	return nil
}
