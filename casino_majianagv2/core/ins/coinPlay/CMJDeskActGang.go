package coinPlay

import "casino_server/common/log"

//成都金币场安杠
func (d *CMJDesk) ActGang(userId uint32, paiId int32, bu bool) error {
	log.T("锁日志: %v CMJDesk.ActGang(%v,%v)的时候等待锁", d.DlogDes(), userId, paiId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v CMJDesk.ActGang(%v,%v)的时候释放锁", d.DlogDes(), userId, paiId)
	}()
	//骨架的杠方法
	err := d.SkeletonChengDuMJDesk.ActGang(userId, paiId, bu)
	if err != nil {
		log.E("杠牌出错...err %v ", err)
		return err
	}
	//杠牌之后处理下一个CheckCase
	err = d.DoCheckCase() //杠牌之后，处理下一个判定牌
	if err != nil {
		log.E("杠牌之后DoCheckCase,出错...err %v ", err)
		return err
	}
	return nil
}
