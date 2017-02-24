package coinPlay

import "casino_common/common/log"

//个人开始定缺
func (d *CMJDesk) DingQue(userId uint32, color int32) error {
	log.T("锁日志: %v CMJDesk.DingQue(%v,%v)的时候等待锁", d.DlogDes(), userId, color)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v CMJDesk.DingQue(%v,%v)的时候释放锁", d.DlogDes(), userId, color, )
	}()

	err := d.SkeletonMJDesk.DingQue(userId, color)
	if err != nil {
		log.E("% 玩家 %v 定缺的时候出错:%v", d.DlogDes(), userId, err)
		return err
	}
	return nil
}
