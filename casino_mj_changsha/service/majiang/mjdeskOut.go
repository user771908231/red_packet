package majiang

import (
	"casino_common/common/log"
	"casino_common/common/consts"
	"casino_common/common/Error"
)

var ERR_OUTPAI error = Error.NewError(consts.ACK_RESULT_ERROR, "打牌失败")

func (d *MjDesk) ActOut(userId uint32, paiKey int32, isAuto bool) error {
	defer Error.ErrorRecovery("actOut")
	log.T("锁日志: %v ActOut(%v,%v)的时候等待锁", d.DlogDes(), userId, paiKey)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActOut(%v,%v)的时候释放锁", d.DlogDes(), userId, paiKey, )
	}()

	return d.ActOutChangSha(userId, paiKey, isAuto)

}
