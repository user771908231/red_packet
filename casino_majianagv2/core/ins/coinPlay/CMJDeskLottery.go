package coinPlay

import (
	"casino_common/common/userService"
	"casino_majiang/service/majiang"
	"casino_common/common/log"
	"casino_majianagv2/core/api"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
)

//成都麻将的 结算方式
func (d *CMJDesk) LotteryChengDu() error {
	//结账需要分两中情况
	/**
		1，只剩一个玩家没有胡牌的时候
		2，没有生育麻将的时候.需要分别做处理...
	 */

	//判断是否可以胡牌
	log.T("现在开始处理lottery()的逻辑....")

	//查花猪
	d.ChaHuaZhu()

	//查大叫
	d.ChaDaJiao()

	//1，处理开奖的数据,
	d.DoLottery()

	//发送结束的广播
	d.SendLotteryData()

	//开奖之后 desk需要处理
	d.AfterLottery()
	return nil
}

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
	d.Room.GetRoomMgr().GetRobotManger().ReleaseRobots(user.GetUserId())
	return nil
}

func (d *CMJDesk) beginEnter() error {
	d.GetStatus().SetStatus(majiang.MJDESK_STATUS_READY) //桌子开始ready
	d.initEnterTimer()                                   //beginEnter的时候 初始化 目前只有机器人才有
	d.initReadyTimer()                                   //已经再房间立的人，需要为准备做倒计时
	return nil
}

func (d *CMJDesk) ForceOutReadyTimeOutUser(userId uint32) {
	//因为此方法是在afterFun执行，所以需要上锁
	log.T("锁日志: %v ForceOutReadyTimeOutUser(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ForceOutReadyTimeOutUser(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//如果玩家没有准备，强制玩家退出
	readyUser := d.GetUserByUserId(userId)
	if readyUser != nil && !readyUser.GetStatus().IsReady() {
		log.T("%s,玩家[%v]超时没有准备，强制退出", d.DlogDes(), userId)
		err := d.rmUser(readyUser) //强制退出一个玩家
		if err != nil {
			log.E("%v玩家[%v]准备超时，强制退出的时候出错", d.DlogDes(), userId)
		} else {
			//发送强制离开的广播
			ack := new(ddproto.CommonBcKickout)
			ack.UserId = proto.Uint32(readyUser.GetUserId())
			ack.Type = ddproto.COMMON_ENUM_KICKOUT_K_TIMEOUT_NOTREADY_ENTERDESK.Enum()
			ack.Msg = proto.String("准备超时,退出房间")
			readyUser.WriteMsg(ack)
		}
	}
}

//设置准备倒计时
func (d *CMJDesk) initReadyTimer() {
	for _, u := range d.Users {
		if u != nil {
			log.T("%s 开始给玩家[%v]设置准备倒计时", d.DlogDes(), u.GetUserId())
			userId := u.GetUserId()
			d.GetSkeletonMJUserById(userId).ReadyTimer = d.AfterFunc(majiang.COIN_RAEDY_DURATION, func() {
				d.ForceOutReadyTimeOutUser(userId)
			})
		}
	}
}
