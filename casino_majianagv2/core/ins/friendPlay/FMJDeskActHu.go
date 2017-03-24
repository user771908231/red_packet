package friendPlay

func (d *FMJDesk) ActHu(userId uint32) error {
	d.SkeletonMJDesk.ActHu(userId)
	//胡牌之后，需要判断游戏是否结束...
	if d.Time2Lottery() {
		//倒倒胡 某一玩家胡牌即结束
		return d.LotteryChengDu() //成都 胡牌之后判断是否lottery
	} else {
		//处理下一个
		return d.DoCheckCase() //胡牌之后，处理下一个判定牌
	}
}
