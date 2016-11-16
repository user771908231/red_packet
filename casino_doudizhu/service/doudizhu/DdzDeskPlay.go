package doudizhu

import (
	"casino_server/common/log"
	"casino_server/common/Error"
)

//这里主要存放 玩斗地主的一些多逻辑....其他的基本方法都放在DdzDesk中


//
func (d *DdzDesk) EnterUser(userId uint32) error {
	//判断是否是重复进入
	olduser := d.GetUserByUserId(userId)
	if olduser != nil {
		//这里需要判断是否是短线重连
		olduser.SetOnline()
		//todo 返回信息
		olduser.UpdateSession()       //更新session 信息，这里可以更具需求来保存对应的属性...
		return
	}

	//新进入
	errAddNew := d.AddUser(userId)
	if errAddNew != nil {
		//进入失败 //返回失败的信息
		return Error.NewError(-1, "进入房间失败")
	} else {
		//进入成功
		//todo 返回进入房间成功的消息
		return nil
	}

	return nil
}

//开始准备
func (d *DdzDesk) Ready(userId uint32) error {
	d.Lock()
	defer d.Unlock()

	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("玩家[%v]准备游戏的时候失败，deks[%v]中没有找到对应的玩家[%v]", userId, d.GetDeskId(), userId)
		return Error.NewError(-1, "没有找到对应的玩家")
	}

	//设置状态为准备的状态...
	user.SetStatus(DDZUSER_STATUS_READY)

	//准备之后的处理
	d.AfterReady()

	return nil

}

func (d *DdzDesk) AfterReady() error {
	//如果都准备好了，那么可以开始抢地主或者发牌了..
	if d.IsAllReady() {
		//开始处理准备之后的事情...
		// todo 开始抢地主.

	}
	return nil
}




//开始游戏
func (d *DdzDesk) Begin() error {

	//判断是否可以开始
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
	err = d.beginQiangDiZhu()
	if err != nil {
		log.E("开始斗地主的时候,beginInit()失败..err[%v]", err)
		return err
	}

	return nil
}

func (d *DdzDesk) IsTime2begin() error {
	return nil

}


//开始时候的初始化
func (d *DdzDesk) BeginInit() error {
	//desk.init


	//userInit
	for _, user := range d.Users {
		if user != nil {
			//初始化每个用户
			user.beginInit()
		}
	}
	return nil
}

//开始抢地主的逻辑
func (d *DdzDesk) beginQiangDiZhu() error {
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

//判断出牌的用户是否合法
func (d *DdzDesk) CheckOutUser(userId uint32) error {
	if d.GetActiveUser() == userId {
		return nil
	} else {
		return Error.NewError(-1, "activeUser 不正确...")
	}
}

//打牌
func (d *DdzDesk) ActOut(userId uint32, out POutPokerPais) error {

	err := d.CheckOutUser(userId)
	if err != nil {

	}

	//判断牌是否合法
	err = d.CheckOutPai(out)
	if err != nil {
		log.E("玩家[%v]出牌[%v]失败,desk.outpai[%v]", userId, out, d.OutPai)
		return Error.NewError(-1, "出牌失败，牌型有误")
	}

	user := d.GetUserByUserId(userId)
	if user == nil {
		return Error.NewError(-1, "没有找到玩家,出牌失败")
	}

	//牌型合法,1保存用户出的牌，2删除手里面的牌
	err = user.DOPoutPokerPais(out)
	if err != nil {
		log.E("玩家[%v]出牌的时候错误", userId)
		return Error.NewError(-1, "玩家出牌的时候出错.")
	}
	//返回成功的消息 todo  返回ack

	//判断游戏是否结束
	if user.GetHandPaiCount() == 0 {
		//出牌的人 手牌为0，表示游戏结束
		d.Lottery()
		return
	}

	d.NextUser()
	return nil
}

//用户过牌，不出牌
func (d *DdzDesk) ActPass(userId uint32) {

}

//轮到下一个人出牌
func (d *DdzDesk) NextUser() error {
	index := d.GetUserIndexByUserId(d.GetActiveUser())
	if index < 0 {
		log.E("轮到下一个玩家的时候出错,desk.activeUser[%v]", d.GetActiveUser())
		return Error.NewError(-1, "轮到一下个玩家的时候出错.")
	}

	nextUser := d.Users[(index + 1) % d.UserCountLimit]
	if nextUser == nil {
		log.E("轮到下一个玩家的时候出错,desk.activeUser[%v]", d.GetActiveUser())
		return Error.NewError(-1, "轮到一下个玩家的时候出错.")
	} else {
		d.SetActiveUser(nextUser.GetUserId())
		//todo 发送下一个人开始出牌的广播
	}

	return nil
}

//抢地主
func (d *DdzDesk) QiangDiZhu(userId uint32, qiangType int32) error {

	//验证活动玩家
	err := d.CheckActiveUser(userId)
	if err != nil {
		return err
	}

	//验证用户是否为空
	user := d.GetUserByUserId(userId)
	if user == nil {
		return Error.NewFailError("玩家没找到，抢地主失败")
	}

	//抢地主
	if qiangType == 1 {
		//开始抢地主的逻辑
		d.SetDizhu(user.GetUserId())
		user.SetQiangDiZhuStatus(DDZUSER_QIANGDIZHU_STATUS_QIANG)
		//表示地主都抢过来，又轮到第一家，抢地主的逻辑结束
		if user.IsQiangDiZhu() && d.GetDizhuPaiUser() == user.GetUserId() {
			//todo 抢地主结束的操作
			return nil
		}

	} else if qiangType == 2 {
		//不叫地主,直接轮到下一个人抢地主
		user.SetQiangDiZhuStatus(DDZUSER_QIANGDIZHU_STATUS_PASS)
	}


	//查找下一家抢地主的人
	index := d.GetUserIndexByUserId(user.GetUserId())
	var nextUser *DdzUser
	for i := index + 1; i < len(d.Users) + index; i++ {
		u := d.Users[(i) / len(d.Users)]
		if u != nil && !u.IsBuJiao() {
			nextUser = u
		}
	}
	//表示没有下一家可以抢地主
	if nextUser == nil {
		//todo 抢地主结束的操作
		d.afterQiangDizhu()
	} else {
		//todo 给nextUser 发送抢地主的协议

	}

	return nil
}

//地主开始出牌
func (d *DdzDesk) afterQiangDizhu() {

}

