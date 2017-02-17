package coinPlay

import (
	"casino_common/common/userService"
	"casino_majiang/service/majiang"
	"casino_common/common/log"
	"casino_majianagv2/core/api"
)

func (d *CMJDesk) AfterLottery() error {
	//开奖完成之后的一些处理
	if d.OverTurnTimer != nil {
		d.OverTurnTimer.Stop()
	}

	//把信息更新到mgo
	for _, u := range d.GetUsers() {
		if u != nil {
			userService.UpdateUser2MgoById(u.GetUserId())
		}
	}

	d.ClearBreakUser()
	d.ClearLeaveUser()
	d.ClearRobotUser()
	//d.ClearCoinInsufficient() //踢掉金币不足的人

	//desk lottery 处理之后，开始等待新的玩家进入
	d.beginEnter() //一局结束之后的处理

	return nil

}

func (d *CMJDesk) ClearBreakUser() error {
	for _, u := range d.Users {
		if u != nil && u.GetStatus().IsBreak {
			log.T("%v 强制踢掉短线的玩家%v", d.DlogDes(), u.GetUserId())
			d.rmUser(u)
		}
	}
	return nil
}

func (d *CMJDesk) ClearLeaveUser() error {
	for _, u := range d.Users {
		if u != nil && u.GetStatus().IsLeave {
			log.T("%v 强制踢掉离开的玩家%v", d.DlogDes(), u.GetUserId())
			d.rmUser(u) //踢掉离开的人
		}
	}
	return nil
}

func (d *CMJDesk) ClearRobotUser() error {
	for _, u := range d.Users {
		if u != nil && u.GetStatus().IsRobot {
			log.T("%v 强制踢掉机器人的玩家%v", d.DlogDes(), u.GetUserId())
			d.rmUser(u) //踢掉机器人
		}
	}
	return nil
}

//删除一个user
func (d *CMJDesk) rmUser(user api.MjUser) error {
	//删除
	d.GetSkeletonMJDesk().RmUser(user)
	MjroomManagerIns.RobotManger.ReleaseRobots(user.GetUserId())
	return nil
}

func (d *CMJDesk) beginEnter() error {
	d.GetStatus().SetStatus(majiang.MJDESK_STATUS_READY) //桌子开始ready
	d.initEnterTimer()                                   //beginEnter的时候 初始化 目前只有机器人才有
	d.initReadyTimer()                                   //已经再房间立的人，需要为准备做倒计时
	return nil
}

//设置准备倒计时
func (d *CMJDesk) initReadyTimer() {
	for _, u := range d.Users {
		if u != nil {
			log.T("%s 开始给玩家[%v]设置准备倒计时", d.DlogDes(), u.GetUserId())
			userId := u.GetUserId()
			d.getCMJUser(u).ReadyTimer = d.AfterFunc(COIN_RAEDY_DURATION, func() {
				d.ForceOutReadyTimeOutUser(userId)
			})
		}
	}
}
