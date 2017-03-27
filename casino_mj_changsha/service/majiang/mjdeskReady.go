package majiang

import (
	"casino_common/common/log"
	"casino_mj_changsha/msg/funcsInit"
	"casino_common/common/consts"
)

//用户准备
func (d *MjDesk) Ready(userId uint32) error {
	log.T("锁日志: %v ready(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ready(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//判断desk状态
	if d.IsNotPreparing() {
		// 准备失败
		log.E("用户[%v]准备失败.desk[%v]不在准备的状态...", userId, d.GetDeskId())
		return nil
	}

	//找到需要准备的user
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.W("用户[%v]在desk[%v]准备的时候失败,没有找到对应的玩家", userId, d.GetDeskId())
		return ERR_SYS
	}

	if user.IsReady() {
		log.W("玩家[%v]已经准备好了...请不要重新准备...", userId)
		return nil
	}

	//设置为准备的状态,并且停止准备计时器
	user.SetStatus(MJUSER_STATUS_READY)
	*user.Ready = true
	if user.readyTimer != nil {
		user.readyTimer.Stop()
		user.readyTimer = nil
	}

	//准备成功,发送准备成功的广播
	result := newProto.NewGame_AckReady()
	*result.Header.Code = consts.ACK_RESULT_SUCC
	*result.Header.Error = "准备成功"
	*result.UserId = userId
	log.T("广播user[%v]在desk[%v]准备成功的广播..string(%v)", userId, d.GetDeskId(), result.String())
	d.BroadCastProto(result)

	//如果人数还是不够，就需要在计划增加机器人
	d.begin()
	return nil
}
