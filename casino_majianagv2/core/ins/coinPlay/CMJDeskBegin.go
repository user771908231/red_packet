package coinPlay

import (
	"casino_common/common/service/awardService"
	"casino_common/common/log"
	"casino_common/common/userService"
)

//金币场beign
func (d *CMJDesk) begin() error {
	//因为begin是在准备之后才能准备，所以不用上锁，redy()上锁就可以了
	//1，检查是否可以开始游戏
	err := d.Time2begin()
	if err != nil {
		log.T("无法开始游戏:err[%v]", err)
		return err
	}

	log.T("[%v]begin()...", d.DlogDes())

	//2，初始化桌子的状态
	err = d.BeginInit()
	if err != nil {
		log.E("初始化牌的时候出错err[%v]", err)
		return err
	}
	//扣除房卡
	err = d.DeductCreateFee() //扣除房卡
	if err != nil {
		log.E("扣除房卡的时候出错[%v]", err)
		return err
	}
	//3，根据playoptions 发cardsNum张牌
	err = d.InitCards()
	if err != nil {
		log.E("初始化牌的时候出错err[%v]", err)
		return err
	}

	//这里需要判断//开始换三张
	if d.IsNeedExchange3zhang() {
		err = d.BeginExchange()
		if err != nil {
			log.E("发送开始换三张的广播的时候出错err[%v]", err)
			return err
		}
		//如果是换三张，开始换三张
		return nil
	}

	//如果是三人两房 or 四人两房 or 两人两房，游戏直接开始
	if d.GetMJConfig().FangCount == 2 {
		err = d.BeginStart()
		if err != nil {
			log.E("发送游戏开始的广播的时候出错err[%v]", err)
			return err
		}
		return nil
	}

	//开始定缺
	err = d.BeginDingQue()
	if err != nil {
		log.E("开始发送定缺广播的时候出错err[%v]", err)
		return err
	}

	return nil
}

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
