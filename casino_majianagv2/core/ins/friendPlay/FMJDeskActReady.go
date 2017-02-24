package friendPlay

import (
	"casino_common/common/log"
)

//成都麻将处理准备
func (d *FMJDesk) Ready(userId uint32) error {
	log.T("锁日志: %v FMJDesk.ready(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v FMJDesk.ready(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	d.SkeletonMJDesk.Ready(userId) //朋友桌 调用骨架的准备

	//如果人数还是不够，就需要在计划增加机器人
	d.begin()
	return nil
}
