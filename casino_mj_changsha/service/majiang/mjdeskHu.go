package majiang

import (
	mjproto       "casino_mj_changsha/msg/protogo"
	"casino_common/common/log"
	"casino_mj_changsha/msg/funcsInit"
	"sync/atomic"
)

//处理抢杠的逻辑
/**
处理抢杠的逻辑，抢杠的逻辑需要特殊处理...
1,首先是清楚杠牌的info
2,增加碰牌
3,删除杠牌的账单
*/
func (d *MjDesk) DoQiangGang(hu *HuPaiInfo) error {
	if hu.GetHuType() != int32(mjproto.HuType_H_QiangGang) {
		return nil
	}

	log.T("%v 由于CheckCase :%v 开始处理抢杠的逻辑....", d.DlogDes(), d.GetCheckCase())
	dianUser := d.GetUserByUserId(hu.GetSendUserId())
	//1,首先是清楚杠牌的info
	var gangKeys []int32
	for _, pai := range dianUser.GameData.HandPai.GangPais {
		//需要删除的杠牌
		if pai != nil && pai.GetClientId() == hu.Pai.GetClientId() {
			gangKeys = append(gangKeys, pai.GetIndex()) //需要删除的杠牌
			if pai.GetIndex() != hu.Pai.GetIndex() {
				dianUser.GameData.HandPai.PengPais = append(dianUser.GameData.HandPai.PengPais, pai) //碰牌
			}
		}
	}

	//删除杠牌的信息
	dianUser.GameData.DelGangInfo(hu.Pai)
	dianUser.GetGameData().GetHandPai().DelGangInfo(hu.Pai)

	dianUser.PreMoGangInfo = nil //摸牌前的杠牌处理也要删除
	//删除杠牌的账单
	for _, billUser := range d.GetUsers() {
		//处理每一个人的账单,并且减去amount
		billUser.DelBillBean(hu.Pai)
	}

	return nil
}

//计算胡牌的账单
/**
	增加账单
	//todo 这里需要完善算账的逻辑逻辑,目前就自摸和点炮来做
 */
func (d *MjDesk) DoHuBill(hu *HuPaiInfo) {
	isZimo := (hu.GetGetUserId() == hu.GetSendUserId())
	outUser := hu.GetSendUserId()
	huUser := d.GetUserByUserId(hu.GetGetUserId())

	log.T("%v玩家[%v]胡牌，开始处理计算分数的逻辑...", d.DlogDes(), huUser.GetUserId())
	if isZimo {
		//如果是自摸的话，三家都需要给钱
		huUser.AddStatisticsCountZiMo(d.GetCurrPlayCount())
		for _, shuUser := range d.GetUsers() {
			log.T("判断是否计算输的玩家[%v]账单 输的玩家是否为空[%v] 是否在游戏中[%v] 是否不是胡牌的玩家[%v] 是否没有胡牌[%v]", shuUser.GetUserId(), shuUser != nil, shuUser.IsGaming(), shuUser.GetUserId() != huUser.GetUserId(), shuUser.IsNotHu())
			if shuUser != nil && shuUser.IsGaming() && (shuUser.GetUserId() != huUser.GetUserId()) && shuUser.IsNotHu() {
				log.T("开始计算输的玩家[%v]账单", shuUser.GetUserId())
				if banker := d.GetBankerUser(); banker != nil && banker.GetUserId() == huUser.GetUserId() {
					//庄家自摸
					//赢钱的账单
					huUser.AddBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_HU, "用户自摸，获得收入", d.GetYingScore(huUser, hu), hu.Pai, d.GetRoomType())

					//输钱的账单
					shuUser.AddBill(huUser.GetUserId(), MJUSER_BILL_TYPE_SHU_ZIMO, "用户自摸，输钱", -d.GetYingScore(huUser, hu), hu.Pai, d.GetRoomType())
				} else {
					//闲家自摸
					//赢钱的账单
					huUser.AddBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_HU, "用户自摸，获得收入", d.GetYingScore(shuUser, hu), hu.Pai, d.GetRoomType())

					//输钱的账单
					shuUser.AddBill(huUser.GetUserId(), MJUSER_BILL_TYPE_SHU_ZIMO, "用户自摸，输钱", -d.GetYingScore(shuUser, hu), hu.Pai, d.GetRoomType())
				}

				//被自摸的用户没人统计一次
				shuUser.AddStatisticsCountBeiZiMo(d.GetCurrPlayCount())
			}
		}
	} else {

		//如果是点炮的话，只有一家需要给钱...
		shuUser := d.GetUserByUserId(outUser)

		if banker := d.GetBankerUser(); banker != nil && banker.GetUserId() == huUser.GetUserId() {
			//闲家点炮
			//赢钱的账单
			huUser.AddBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_HU, "用户接炮，获得收入", d.GetYingScore(huUser, hu), hu.Pai, d.GetRoomType())

			//输钱的账单
			shuUser.AddBill(huUser.GetUserId(), MJUSER_BILL_TYPE_SHU_DIANPAO, "闲家点炮，输钱", -d.GetYingScore(huUser, hu), hu.Pai, d.GetRoomType())
		} else {
			//庄家点炮
			//赢钱的账单
			huUser.AddBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_HU, "用户接炮，获得收入", d.GetYingScore(shuUser, hu), hu.Pai, d.GetRoomType())

			//输钱的账单
			shuUser.AddBill(huUser.GetUserId(), MJUSER_BILL_TYPE_SHU_DIANPAO, "庄家点炮，输钱", -d.GetYingScore(shuUser, hu), hu.Pai, d.GetRoomType())
		}

		huUser.AddStatisticsCountHu(d.GetCurrPlayCount())       //胡的用户统计信息
		shuUser.AddStatisticsCountDianPao(d.GetCurrPlayCount()) //点炮的用户统计信息
	}
}

//获取玩家赢分 为长沙麻将添加的方法
func (d *MjDesk) GetYingScore(yingUser *MjUser, hu *HuPaiInfo) (score int64) {
	bankAddScore := d.GetBaseValue()
	score = hu.GetScore()
	if !d.IsChangShaMaJiang() {
		return
	}
	if banker := d.GetBankerUser(); banker != nil && banker.GetUserId() == yingUser.GetUserId() {
		//长沙麻将 庄家加分的处理

		//只要包含特殊的胡牌方式 就加分
		isContainHuType := false
		isQishouhu := false
		if hu.GetHuType() > 0 {
			isContainHuType = true
		}
		if hu.GetHuType() == int32(mjproto.HuType_H_changsha_qishouhu) {
			isQishouhu = true
		}

		for _, pt := range hu.GetPaiType() {
			//遍历胡牌类型 给每个大胡和平胡的胡牌类型都加分
			switch mjproto.PaiType(pt) {
			case mjproto.PaiType_H_CHANGSHA_QIXIAODUI_HAOHUA_DOUBLE:
				score += bankAddScore
				score += bankAddScore
				score += bankAddScore
				score += bankAddScore
			case mjproto.PaiType_H_CHANGSHA_QIXIAODUI_HAOHUA:
				score += bankAddScore
				score += bankAddScore
			case mjproto.PaiType_H_CHANGSHA_PINGHU: //平胡
				if !isContainHuType || isQishouhu {
					//没有特殊胡牌方式或者是起手胡才累加平胡的得分
					score += bankAddScore
				}
			default:
				score += bankAddScore
			}
		}
		//有特殊胡牌方式 且 不是起手胡 才累加庄家得分
		if isContainHuType && !isQishouhu {
			score += bankAddScore
		}

	}
	return score
}

func (d *MjDesk) DoAfterDianPao(hu *HuPaiInfo) {
	//自摸的不用关
	if hu.GetGetUserId() == hu.GetSendUserId() {
		return
	}
	//点炮胡牌成功之后的处理... 处理checkCase
	d.CheckCase.UpdateCheckBeanStatus(hu.GetGetUserId(), CHECK_CASE_BEAN_STATUS_CHECKED) // update checkCase...
	atomic.AddInt32(d.CheckCase.DianPaoCount, 1)
	//胡牌之后，更新canPeng 或者 canGang 的checkCase
	for _, bean := range d.CheckCase.CheckB {
		if bean != nil && (bean.GetCanPeng() || bean.GetCanGang() || bean.GetCanChi()) {
			//*bean.CheckStatus = CHECK_CASE_BEAN_STATUS_CHECKED
			*bean.CanPeng = false
			*bean.CanGang = false
			*bean.CanChi = false
		}
	}

	// 删除点炮者的out牌 //todo 如果有多个人胡牌，那么有可能重复删除？这里怎么处理
	outUser := d.GetUserByUserId(hu.GetSendUserId())
	errDelOut := outUser.GameData.HandPai.DelOutPai(hu.Pai.GetIndex())
	if errDelOut != nil {
		log.E("胡牌的时候，删除打牌玩家的out牌[%v]...注意:[这里有可能不是错误，一炮多响的情况会多次删除手牌，就可能出现这种情况，等待解决] err[%v]", hu.Pai.GetIndex(), errDelOut)
	}

}

//设置下一次的庄
/**
	1，如果当前的nextBanker 没有值(nextBanker==0)，那代表此人是第一个胡牌的，设置为nextBanekr
	2,如果当前的nextBanker有值(nextBanker > 0 ),那需要判断是不是当前的点炮的人一炮双向
 */
func (d *MjDesk) InitNextBanker(hu *HuPaiInfo) {
	log.T("%v  开始设置下一次[%v]的庄的玩家通过huinfo %v", d.DlogDes(), d.GetNextBanker(), hu)
	if d.IsNextBankerExist() {
		//已经存在的情况 //有双响就双响点炮的人做庄，不论之前是否有人胡牌  by 亮哥
		//这里可以用过pai 查询点炮账单的个数

		//判断是否是自摸
		isZimo := (hu.GetGetUserId() == hu.GetSendUserId())
		if isZimo {
			//如果是自摸，并且nextBanker已经有值了,那么直接返回不用设置
			return
		}
		//表示多响
		if d.GetCheckCase().GetDianPaoCount() > 1 {
			//设置一炮多响的人为庄
			d.SetNextBanker(hu.GetSendUserId())
		}
	} else {
		d.SetNextBanker(hu.GetGetUserId())
	}
}

func (d *MjDesk) SendAckActHu(hu *HuPaiInfo) {
	ack := newProto.NewGame_AckActHu()
	*ack.HuType = hu.GetHuType() //这里需要判断是自摸还是点炮
	*ack.UserIdIn = hu.GetGetUserId()
	*ack.UserIdOut = hu.GetSendUserId()
	ack.HuCard = hu.Pai.GetCardInfo()
	*ack.IsZiMo = (hu.GetGetUserId() == hu.GetSendUserId())
	ack.PaiType = IntArry2PaiTypeEnum(hu.PaiType)
	ack.UserCoinBeans = d.GetUserCoinBeans()
	log.T("给用户[%v]广播胡牌的ack[%v]", hu.GetGetUserId(), ack)
	d.BroadCastProto(ack)
}
