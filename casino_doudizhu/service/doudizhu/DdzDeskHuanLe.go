package doudizhu

import (
	"casino_server/common/log"
	"casino_server/common/Error"
	"casino_doudizhu/msg/funcsInit"
)

//欢乐斗地主
func (d *DdzDesk) HLBegin() error {
	//判断是否可以开始
	err := d.IsTime2begin()
	if err != nil {
		log.E("开始斗地主的时候失败,不满足开始的条件err[%v]", err)
		return err
	}

	//开始叫地主
	err = d.sendJiaoDiZhuOverTurn(d.GetDizhuPaiUser())        //欢乐斗地主开始叫地主
	if err != nil {
		log.E("发送开始叫地主的时候出错，游戏begin失败")
		return err
	}

	return nil
}

//叫地主
func (d *DdzDesk) HLJiaoDiZhu(userId uint32) error {
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

	//k碍事叫地主
	user.SetQiangDiZhuStatus(DDZUSER_QIANGDIZHU_STATUS_JIAO)        //设置用户已经叫地主
	//叫地主成功，广播信息
	ack := newProto.NewDdzJiaoDiZhuAck()
	*ack.UserId = userId
	*ack.Jiao = true
	d.BroadCastProto(ack)

	//查找下一个没有操作过的人来抢地主
	nextUser := d.GetNextUserByPros(userId, func(u *DdzUser) bool {
		return u != nil && u.IsQiangDiZhuNoAct()
	})

	//表示没有下一家可以抢地主
	if nextUser == nil {
		d.HLAfterQiangDizhu()        //换掉了斗地主，叫地主
	} else {
		d.sendQiangDiZhuOverTurn(nextUser.GetUserId())        //欢乐斗地主，叫地主之后开始抢地主
	}

	return nil
}


//欢乐斗地主，抢地主
func (d *DdzDesk) HLQiangDiZhu(userId uint32, qiang bool) error {

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
	//抢地主和
	if qiang {
		//开始抢地主的逻辑
		d.SetDizhu(user.GetUserId())
		d.AddCountQiangDiZhu()        //增加抢地主的次数
		d.setQingDizhuValue(d.GetQingDizhuValue() * 2)//这里还需要计算低分
		user.SetQiangDiZhuStatus(DDZUSER_QIANGDIZHU_STATUS_QIANG)
	} else {
		user.SetQiangDiZhuStatus(DDZUSER_QIANGDIZHU_STATUS_BUQIANG)
	}


	//广播抢地主
	ack := newProto.NewDdzRobDiZhuAck()
	*ack.UserId = userId
	*ack.Rob = qiang
	d.BroadCastProto(ack)


	//表示地主都抢过来，又轮到第一家，抢地主的逻辑结束
	if user.IsQiangDiZhuQiang() && d.GetDizhuPaiUser() == user.GetUserId() {
		d.HLAfterQiangDizhu()        //欢乐斗地主，抢地主之后，如果已经抢完了，那么抢地主之后的逻辑
		return nil
	}

	//查找下一家抢地主的人
	nextUser := d.GetNextUserByPros(user.GetUserId(), func(u *DdzUser) bool {
		return u != nil && u.IsQiangDiZhuNoAct()
	})

	//表示没有下一家可以抢地主
	if nextUser == nil {
		d.HLAfterQiangDizhu()        //欢乐斗地主，已经没有合适的人抢地主的时候，进行抢地主之后的逻辑
	} else {
		d.sendQiangDiZhuOverTurn(nextUser.GetUserId())
	}

	return nil
}


//叫地主
func (d *DdzDesk) HLBuJiaoDiZhu(userId uint32) error {
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

	//设置用户为不叫地主
	user.SetQiangDiZhuStatus(DDZUSER_QIANGDIZHU_STATUS_BUJIAO)

	//得到下一个可以抢地主的人
	var nextUser = d.GetNextUserByPros(user.GetUserId(), func(u *DdzUser) bool {
		return u != nil && u.IsQiangDiZhuNoAct()
	})


	//表示没有下一家可以抢地主
	if nextUser == nil {
		d.HLAfterQiangDizhu()        //欢乐斗地主，不叫地址，并且没有找到下一个需要抢地主的人，进行抢地主之后的逻辑..
	} else {
		//判断是否有地主
		if d.IsHadDiZhuUser() {
			//发送抢地主的协议
			d.sendQiangDiZhuOverTurn(nextUser.GetUserId())
		} else {
			//发送叫地主的协议
			d.sendJiaoDiZhuOverTurn(nextUser.GetUserId())
		}

	}

	return nil
}

//欢乐斗地主开始加倍的逻辑
func (d *DdzDesk) ActJiaBei(userId uint32) error {
	//判断activeUser
	err := d.CheckActiveUser(userId)
	if err != nil {
		return err
	}
	//低分加倍的逻辑

	//需要返回加倍的ack
	ack := newProto.NewDdzDoubleAck()
	*ack.UserId = userId
	*ack.Double = 0
	d.BroadCastProto(ack)

	//下一个人加倍
	nextUser := d.GetNextUserByPros(userId, func(u *DdzUser) bool {
		return u != nil && u.IsJiaBeiNoAct()
	})

	if nextUser != nil {
		//发送加倍的overTurn
	} else {
		//加倍已经完成了，开始游戏


	}
	return nil
}

//处理加倍之后的操作
func (d *DdzDesk)  AfterJiaBei() error {
	//判断是否都已经加倍了
	if d.IsAllActJiaBei() {
		//进行下一步操作...地主出牌
		err := d.SendChuPaiOverTurn(d.GetDiZhuUserId())
		if err != nil {
			log.E("发送出牌overTurn 的时候出错")
			return err
		}
	} else {
		//返回错误
		return Error.NewFailError("还有玩家没有确认加倍的操作..")
	}
	return nil
}