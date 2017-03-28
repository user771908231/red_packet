package majiang

import (
	"casino_common/common/Error"
	"casino_common/common/log"
	mjproto	"casino_mj_changsha/msg/protogo"
	"github.com/golang/protobuf/proto"
)

/**
长沙麻将的结算方式..结算之前开始抓鸟
*/
func (d *MjDesk) LotteryChangSha() error {
	log.T("现在开始处理lottery()的逻辑....")
	//1,抓鸟不用玩家点击，系统自动生成就行了...
	err := d.doChangShaBird()
	if err != nil {
		return err
	}
	//2，处理开奖的数据,
	err = d.DoLottery()
	if err != nil {
		return err
	}
	//3,发送结束的广播
	err = d.SendLotteryData()
	if err != nil {
		return err
	}

	//4,开奖之后 desk需要处理
	err = d.AfterLottery()
	if err != nil {
		return err
	}
	//5,判断牌局结束(整场游戏结束)
	d.End()
	return nil
}

//长沙麻将，抓鸟的逻辑
func (d *MjDesk) doChangShaBird() error {
	log.T("桌子[%v] 开始抓鸟", d.GetDeskId())

	//抓鸟数组
	if d.ChangShaPlayOptions == nil {
		log.E("桌子[%v]的房间信息丢失, 抓鸟失败", d.GetDeskId())
		return Error.NewError(-1, "房间信息错误")
	}
	birdPais := d.getBirdPais(int(d.GetRoomTypeInfo().GetChangShaPlayOptions().GetBirdCount()))

	//得到发起抓鸟的玩家
	u := d.getBeginBirdUser()
	if u == nil {
		log.E("桌子[%v] 未找到发起抓鸟的玩家 扎鸟失败", d.GetDeskId())
		return nil
	}

	//遍历鸟牌
	for _, birdPai := range birdPais {
		cpBird := birdPai
		target := d.releaseABird(u, cpBird)
		if target == nil {
			log.E("桌子[%v]扎到的玩家不存在 扎鸟出错", d.GetDeskId())
			continue
		}

		//u.AddStatisticsCountCatchBird(d.GetCurrPlayCount())
		//target.AddStatisticsCountCaughtBird(d.GetCurrPlayCount())

		isCatched := d.doBirdBill(target, cpBird)

		birdInfo := new(mjproto.BirdInfo)
		birdInfo.BirdPai = cpBird.GetCardInfo()
		birdInfo.ZhuaUser = proto.Uint32(u.GetUserId())
		birdInfo.BirdUser = proto.Uint32(target.GetUserId())
		if !isCatched {
			birdInfo.BirdUser = proto.Uint32(0)
		}
		d.BirdInfo = append(d.BirdInfo, birdInfo)
		log.T("桌子[%v] 得到抓鸟的信息birdInfo[%v]", d.GetDeskId(), birdInfo)
	}
	log.T("桌子[%v] 抓鸟结束", d.GetDeskId())
	return nil
}

func (d *MjDesk) getBeginBirdUser() *MjUser {
	if d.GetRoomTypeInfo().GetChangShaPlayOptions().GetIgnoreBank() {
		//不区分庄闲 则取下一局的庄家为开始抓鸟的玩家
		return d.GetUserByUserId(d.GetNextBanker())
	}
	//区分庄家 则取本局的庄家为开始抓鸟的玩家
	return d.GetBankerUser()
}

//获取几只鸟牌
func (d *MjDesk) getBirdPais(num int) (birdPais []*MJPai) {
	log.T("桌子[%v] 开始获取鸟牌", d.GetDeskId())
	//这里注意 牌堆中牌不够时 取最后一张牌为鸟

	for i := 0; i < num; i++ {
		birdPai := d.GetNextPai() //抓鸟牌
		if birdPai != nil {
			log.T("桌子[%v] 抓到鸟牌[%v]", d.GetDeskId(), birdPai.GetDes())
			birdPais = append(birdPais, birdPai)
			continue
		}
		birdPai = d.AllMJPai[len(d.AllMJPai)-1]
		log.T("桌子[%v] 牌堆空了 将最后一张作为鸟牌[%v]", birdPai.GetDes())
		birdPais = append(birdPais, birdPai)
	}
	log.T("桌子[%v] 成功获取鸟牌[%v]", d.GetDeskId(), birdPais)
	return birdPais
}

//扎鸟 由玩家u释放一个鸟牌 根据鸟牌的牌值得到扎到的玩家
func (d *MjDesk) releaseABird(u *MjUser, bird *MJPai) (target *MjUser) {
	log.T("桌子[%v]开始扎鸟 鸟牌为[%v]", d.GetDeskId(), bird.GetDes())

	//根据鸟牌的牌值得到扎到的玩家
	uIndex := d.getIndexByUserId(u.GetUserId())
	if uIndex == -1 {
		log.E("桌子[%v] 未找到发起抓鸟的玩家[%v]索引 扎鸟失败", d.GetDeskId(), u.GetUserId())
		return nil
	}
	targetIndex := (bird.GetValue() - 1 + int32(uIndex)) % d.GetUserCountLimit() //计算扎到的玩家索引 (鸟牌值 - 1 + 发起人的索引) % 玩家人数
	target = d.GetUserByIndex(int(targetIndex))
	if target == nil {
		log.E("桌子[%v] 未找到抓鸟玩家 玩家索引[%v] 扎鸟失败", d.GetDeskId(), targetIndex)
		return nil
	}
	log.T("桌子[%v]结束扎鸟 鸟牌为[%v] 扎到的玩家为[%v]", d.GetDeskId(), bird.GetDes(), target.GetUserId())
	return target
}

//计算鸟扎中玩家的账单
func (d *MjDesk) doBirdBill(u *MjUser, birdPai *MJPai) (isCatched bool) {
	log.T("桌子[%v]开始计算玩家[%v]的抓鸟账单", d.GetDeskId(), u.GetUserId())
	//将玩家的账单翻倍
	if u.GetBill() == nil {
		log.E("计算抓鸟账单错误 玩家[%v]的账单数据异常", u.GetUserId())
		return
	}

	if u.GetBill().GetBills() == nil {
		log.T("玩家[%v]的没有账单数据, 无法计算抓鸟", u.GetUserId())
		return
	}

	userBills := u.GetBill().GetBills()
	if userBills != nil && len(userBills) <= 0 {
		log.T("计算抓鸟账单错误 玩家[%v]的账单为空", u.GetUserId())
		return
	}

	bills := make([]*BillBean, len(userBills))
	copy(bills, u.GetBill().GetBills())

	if bills[0].GetAmount() > 0 {
		//玩家是赢钱 增加抓鸟的统计数据
		u.AddStatisticsCountCatchBird(d.GetCurrPlayCount())
	} else {
		//玩家是输钱 增加被抓鸟的统计数据
		u.AddStatisticsCountCaughtBird(d.GetCurrPlayCount())
	}

	for _, billBean := range bills {

		//对已经是抓鸟账单的不作处理
		if billBean.GetIsBird() {
			continue
		}

		//如果是起手胡牌的账单，那么不做处理
		if billBean.GetIsQiShouHu() {
			continue
		}

		isCatched = true

		//给中鸟的玩家增加抓鸟账单
		log.T("桌子[%v] 给玩家[%v] 添加抓鸟账单", d.GetDeskId(), u.GetUserId())
		d.AddBirdBill(u, billBean, birdPai)

		//给关联玩家的非抓鸟账单加倍
		relationUser := d.GetUserByUserId(billBean.GetOutUserId())
		rBills := relationUser.GetBill().GetBills()

		if rBills[0].GetAmount() > 0 {
			//关联玩家是赢钱 增加抓鸟的统计数据
			relationUser.AddStatisticsCountCatchBird(d.GetCurrPlayCount())
		} else {
			//关联玩家是输钱 增加被抓鸟的统计数据
			relationUser.AddStatisticsCountCaughtBird(d.GetCurrPlayCount())
		}

		rBB := &BillBean{
			UserId:    proto.Uint32(billBean.GetOutUserId()),
			OutUserId: proto.Uint32(billBean.GetUserId()),
			Type:      proto.Int32(billBean.GetType()),
			Des:       proto.String(billBean.GetDes()),
			Amount:    proto.Int64(-billBean.GetAmount()),
			Pai:       billBean.GetPai(),
			IsBird:    proto.Bool(billBean.GetIsBird()),
		}
		log.T("桌子[%v] 给玩家[%v] 添加抓鸟账单", d.GetDeskId(), relationUser.GetUserId())
		d.AddBirdBill(relationUser, rBB, birdPai)
		//for _, rBillBean := range rBills {
		//	if rBillBean.GetOutUserId() == u.GetUserId() && !rBillBean.GetIsBird() {
		//		log.T("桌子[%v] 给玩家[%v] 添加抓鸟账单", d.GetDeskId(), relationUser.GetUserId())
		//		d.AddBirdBill(relationUser, rBillBean, birdPai)
		//	}
		//}
	}
	log.T("桌子[%v]玩家[%v]的抓鸟账单计算完成", d.GetDeskId(), u.GetUserId())
	return isCatched
}

//增加抓鸟账单
func (d *MjDesk) AddBirdBill(u *MjUser, bb *BillBean, birdPai *MJPai) {
	fenshu := int64(0)
	baseNiaoValue := d.GetBaseValue()
	switch d.GetRoomTypeInfo().GetChangShaPlayOptions().GetBirdMultiple() {
	case 2: //加倍
		fenshu = bb.GetAmount()
	default:
		//加底
		switch {
		case bb.GetAmount() > 0:
			fenshu = baseNiaoValue
		case bb.GetAmount() < 0:
			fenshu = -baseNiaoValue
		default:
			fenshu = 0
		}
	}
	u.AddChangShaBill(bb.GetOutUserId(), bb.GetType(), bb.GetDes()+" 抓鸟", fenshu, birdPai, d.GetRoomType(), true)
}
