package coinPlay

import (
	"casino_common/common/service/awardService"
	"casino_common/common/log"
	"casino_common/common/userService"
)

func (d *CMJDesk) DeductCreateFee() error {
	//扣除房费,每个人都要扣除
	for _, u := range d.GetUsers() {
		if u != nil {
			log.T("玩家[%v]进入麻将金币场，准备之后扣除房费[%v]", u.GetUserId(), d.GetMJConfig().CreateFee)
			remain, err := userService.DECRUserCOIN(u.GetUserId(), d.GetMJConfig().CreateFee) //扣除进入房间的房费并更新redisUser的金币数
			if err != nil {
				log.E("%v 金币场开始游戏，扣除玩家[%v]房费[%v]的时候出错,err %v ", d.DlogDes(), u.GetUserId(), d.GetMJConfig().CreateFee, err)
			} else {
				u.SetCoin(remain)                                                       //玩家的金币
				awardService.AddOnlineCapital(u.GetUserId(), d.GetMJConfig().CreateFee) //增加在线奖励的提成本金
			}
		}
	}
	return nil
}
