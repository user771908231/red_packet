package friendPlay

import (
	"casino_common/common/log"
	"casino_common/common/userService"
	"casino_majiang/service/majiang"
)

//开始游戏
func (d *FMJDesk) begin() error {
	//因为begin是在准备之后才能准备，所以不用上锁，redy()上锁就可以了
	//1，检查是否可以开始游戏
	err := d.Time2begin()
	if err != nil {
		log.T("无法开始游戏:err[%v]", err)
		return err
	}

	log.T("[%v]begin()...", d.DlogDes())

	//2，初始化桌子的状态
	d.BeginInit()
	d.DeductCreateFee() //扣除房卡
	//3，根据playoptions 发cardsNum张牌
	err = d.InitCards()
	if err != nil {
		log.E("初始化牌的时候出错err[%v]", err)
		return err
	}

	//这里需要判断
	if d.IsNeedExchange3zhang() {
		//开始换三张
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

//扣除房费
func (d *FMJDesk) DeductCreateFee() error {
	//计算朋友桌的房费
	if d.GetMJConfig().CurrPlayCount != 1 {
		return nil
	}

	log.T("%v 玩家[%v]开始扣除房卡，消耗房卡:%v", d.DlogDes(), d.GetMJConfig().Owner, d.GetMJConfig().CreateFee)
	remain, err := userService.DECRUserRoomcard(d.GetMJConfig().Owner, d.GetMJConfig().CreateFee)
	if err != nil {
		//扣费失败，创建房间失败
		log.E("玩家[%v]创建房间的时候扣房卡[%v]失败，创建房间失败...", d.GetMJConfig().Owner, d.GetMJConfig().CreateFee)
		return majiang.ERR_ROOMCARD_INSUFFICIENT
	}
	log.T("%v 玩家[%v]扣除房卡成功，消耗房卡:%v,剩余房卡:%v", d.DlogDes(), d.GetMJConfig().Owner, d.GetMJConfig().CreateFee, remain)
	return nil
}
