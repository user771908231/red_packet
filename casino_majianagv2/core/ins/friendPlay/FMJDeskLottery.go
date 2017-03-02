package friendPlay

import (
	"casino_common/common/log"
	"casino_majiang/service/majiang"
	"casino_common/common/userService"
)

//成都麻将的 结算方式
func (d *FMJDesk) LotteryChengDu() error {
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

	//判断牌局结束(整场游戏结束)
	if !d.End() {
		//go d.begin()
	}
	return nil
}

func (d *FMJDesk) AfterLottery() error {
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
	//如果是金币场，需要把短线的，离开的，机器人都踢走
	d.GetStatus().SetStatus(majiang.MJDESK_STATUS_READY) //桌子开始ready
	return nil

}

func (d *FMJDesk) End() bool {
	//判断结束的条件,目前只有局数能判断
	log.T("%v游戏是否End() CurrPlayCount[%v], TotalPlayCount[%v]",
		d.DlogDes(), d.GetMJConfig().CurrPlayCount, d.GetMJConfig().TotalPlayCount)
	//朋友桌有整场结束的概念
	if d.GetMJConfig().CurrPlayCount < d.GetMJConfig().TotalPlayCount {
		//表示游戏还没有结束。。。.
		return false
	} else {
		d.DoEnd()
		return true
	}

}
