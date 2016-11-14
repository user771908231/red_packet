package doudizhu

import "casino_server/common/log"

//这里主要存放 玩斗地主的一些多逻辑....其他的基本方法都放在DdzDesk中

//开始游戏
func (d *DdzDesk) Begin() error {

	//
	err := d.IsTime2begin()
	if err != nil {
		log.E("开始斗地主的时候失败,不满足开始的条件err[%v]", err)
		return err
	}


	//初始化，这里着重初始化 默认值，状态等...
	err = d.BeginInit()
	if err != nil {
		log.E("开始斗地主的时候,beginInit()失败..err[%v]", err)
		return err
	}

	//开始抢地主


	return nil
}

func (d *DdzDesk) IsTime2begin() error {
	return nil

}


//开始时候的初始化
func (d *DdzDesk) BeginInit() error {
	return nil
}

//一场结束
func (d *DdzDesk) Lottery() {

}

//牌局结束
func (d *DdzDesk) DoEnd() {

}

//初始化每个人的牌
func (d *DdzDesk) InitCards() {
	//获得一副洗好的牌
	d.AllPokerPai = XiPai()                //获得洗好的一副扑克牌

	//为每个人分配牌
	for i, user := range d.Users {
		user.GameData.HandPokers = make([]*PPokerPai, 17)        //这里的17暂时写死...
		copy(user.GameData.HandPokers, d.AllPokerPai[(i - 1) * 17:i * 17])
	}

	//底牌
	d.DiPokerPai = make([]*PPokerPai, 3)
	copy(d.DiPokerPai, d.AllPokerPai[54 - 3:54])

}