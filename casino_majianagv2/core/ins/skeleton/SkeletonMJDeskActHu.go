package skeleton

import (
	"casino_majiang/service/majiang"
	"strings"
	"casino_majiang/msg/protogo"
	"casino_common/common/log"
	"casino_common/common/Error"
	"errors"
	"sync/atomic"
	"casino_majiang/msg/funcsInit"
)

func (d *SkeletonMJDesk) ActHu(userId uint32) error {
	log.T("锁日志: %v ActHuChangsha(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActHuChangsha(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//对于杠，有摸牌前 杠的状态，有打牌前杠的状态

	if d.CheckNotActUser(userId) { //胡牌
		return Error.NewFailError("不该当前用户操作..")
	}
	if d.OverTurnTimer != nil {
		d.OverTurnTimer.Stop()
	}

	//玩家胡牌
	huUser := d.GetUserByUserId(userId)
	if huUser == nil {
		log.E("服务器错误：没有找到胡牌的user[%v]", userId)
		return errors.New("服务器错误，没有找到胡牌的user")
	}

	//1,胡的牌是当前check里面的牌，如果没有check，则表示是自摸
	//得到胡牌的信息
	var isZimo bool = false                //是否是自摸
	var isLastPai bool = !d.HandPaiCanMo() //是否是最后一张牌
	var isPreMoGang bool = false           //之前是否有摸杠
	var extraAct mjproto.HuType = 0        //杠上花，
	var outUserId uint32
	var huPai *majiang.MJPai

	//判断胡的牌是别人打的还是自己摸的...
	if d.GetCheckCase() != nil {
		huPai = d.GetCheckCase().CheckMJPai
		outUserId = d.GetCheckCase().GetUserIdOut()
	} else {
		isZimo = true //自摸
		huPai = huUser.GetGameData().GetHandPai().GetInPai()
		outUserId = userId
	}

	//胡牌之前的杠牌信息
	if huUser.GetGameData().GetPreMoGangInfo() != nil {
		isPreMoGang = true
	}

	switch {
	case isZimo && isPreMoGang && !isLastPai: // 杠上花
		extraAct = mjproto.HuType_H_GangShangHua
	case !isZimo && isPreMoGang && !isLastPai: //杠上炮
		extraAct = mjproto.HuType_H_GangShangPao
	case !isZimo && isPreMoGang && isLastPai: //海底杠上炮
		extraAct = mjproto.HuType_H_HaidiGangShangPao
	case isZimo && isPreMoGang && isLastPai: //海底杠上花
		extraAct = mjproto.HuType_H_HaidiGangShangHua
	case isZimo && !isPreMoGang && isLastPai: //海底捞
		//isHaiDiLao = true
		extraAct = mjproto.HuType_H_HaiDiLao
	default:
		extraAct = 0
	}

	//得到inpai 之后，在获取枫树之前需要先判断是否可以胡牌，如果不能胡牌直接返回
	canHu, fan, score, huCardStr, paiType, _ := d.HuParser.GetCanHu(huUser.GetGameData().GetHandPai(), huPai, isZimo, extraAct)
	if !canHu {
		log.E("玩家[%v]不可以胡", userId)
		return errors.New("不可以胡牌...")
	}
	//二人两房 平胡不能被点炮
	if d.GetMJConfig().TotalPlayCount == 2 && d.GetMJConfig().FangCount == 2 && !isZimo {
		//是两人两房 并且 不是自摸
		if score <= 0 {
			//得分为零
			log.E("玩家[%v]不可以胡", userId)
			return errors.New("二人两房 平胡不能被点炮...")
		}
	}

	//胡牌之后的信息
	hu := majiang.NewHuPaiInfo()
	*hu.GetUserId = huUser.GetUserId()
	*hu.SendUserId = outUserId
	*hu.HuType = int32(extraAct) ////杠上炮 杠上花 抢杠 海底捞 海底炮 天胡 地胡
	*hu.HuDesc = strings.Join(huCardStr, " ");
	*hu.Fan = fan
	hu.PaiType = majiang.PaiTypeEnum2IntArry(paiType)
	*hu.Score = score //只是胡牌的分数，不是赢了多少钱
	hu.Pai = huPai
	huUser.GetGameData().AddHuPaiInfo(hu) //胡牌之后，设置用户的数据
	huUser.GetStatus().SetStatus(majiang.MJUSER_STATUS_HUPAI)

	/**
		 //胡牌之后清湖PreMoGang的信息这里这样做的作用是保证计算下一摸牌的人时候有作用
		 还需要考虑一种情况就是血流成和，胡牌之后就不能继续有preMoGang了，否则下一次胡牌还是 杠上花
	 */
	huUser.GetGameData().DelPreMoGangInfo()
	d.SetActiveUser(huUser.GetUserId()) //胡牌之后，设置胡牌的人为activeuser，计算下一家玩家的时候，有用

	//处理抢杠的逻辑
	d.DoQiangGang(hu)

	//胡牌之后计算账单
	d.DoHuBill(hu)

	//点炮之后设置checkCase的状态
	d.DoAfterDianPao(hu)

	//设置下一次的庄
	d.InitNextBanker(hu)

	//发送胡牌成功的回复
	d.SendAckActHu(hu)


	return nil
}

//处理抢杠的逻辑
/**
处理抢杠的逻辑，抢杠的逻辑需要特殊处理...
1,首先是清楚杠牌的info
2,增加碰牌
3,删除杠牌的账单
*/
func (d *SkeletonMJDesk) DoQiangGang(hu *majiang.HuPaiInfo) error {
	if d.CheckCase != nil && d.CheckCase.PreOutGangInfo != nil && d.CheckCase.PreOutGangInfo.GetGangType() == majiang.GANG_TYPE_BA {
		log.T("开始处理抢杠的逻辑....")
		dianUser := d.GetUserByUserId(hu.GetSendUserId())

		//1,首先是清楚杠牌的info
		var gangKeys []int32
		for _, pai := range dianUser.GetGameData().GetHandPai().GetGangPais() {
			if pai == nil {
				continue
			}
			//需要删除的杠牌
			if pai != nil && pai.GetClientId() == hu.Pai.GetClientId() {
				gangKeys = append(gangKeys, pai.GetIndex()) //需要删除的杠牌
				if pai.GetIndex() != hu.Pai.GetIndex() {
					dianUser.GetGameData().GetHandPai().PengPais = append(dianUser.GetGameData().GetHandPai().PengPais, pai) //碰牌
				}
			}
		}

		//删除杠牌的信息
		dianUser.GetGameData().DelGangInfo(hu.Pai)
		dianUser.GetGameData().DelPreMoGangInfo() //摸牌前的杠牌处理也要删除
		//删除杠牌的账单
		for _, billUser := range d.GetUsers() {
			//处理每一个人的账单,并且减去amount
			billUser.DelBillBean(hu.Pai)
		}
	}

	return nil
}

//计算胡牌的账单
/**
	增加账单
	//todo 这里需要完善算账的逻辑逻辑,目前就自摸和点炮来做
 */
func (d *SkeletonMJDesk) DoHuBill(hu *majiang.HuPaiInfo) {
	isZimo := (hu.GetGetUserId() == hu.GetSendUserId())
	outUser := hu.GetSendUserId()
	huUser := d.GetSkeletonMJUserById(hu.GetGetUserId())

	log.T("玩家[%v]胡牌，开始处理计算分数的逻辑...", huUser.GetUserId())
	if isZimo {
		//如果是自摸的话，三家都需要给钱
		huUser.AddStatisticsCountZiMo(d.GetMJConfig().CurrPlayCount)
		for _, shuUser := range d.GetSkeletonMJUsers() {
			if shuUser != nil && shuUser.GetStatus().IsGaming() && (shuUser.GetUserId() != huUser.GetUserId()) && shuUser.GetStatus().IsNotHu() {

				//赢钱的账单
				huUser.AddBill(shuUser.GetUserId(), majiang.MJUSER_BILL_TYPE_YING_HU, "用户自摸，获得收入", d.GetYingScore(huUser, hu.GetScore()), hu.Pai, d.GetMJConfig().RoomType)

				//输钱的账单
				shuUser.AddBill(huUser.GetUserId(), majiang.MJUSER_BILL_TYPE_SHU_ZIMO, "用户自摸，输钱", -d.GetYingScore(shuUser, hu.GetScore()), hu.Pai, d.GetMJConfig().RoomType)

				//被自摸的用户没人统计一次
				shuUser.AddStatisticsCountBeiZiMo(d.GetMJConfig().CurrPlayCount)
			}
		}
	} else {

		//如果是点炮的话，只有一家需要给钱...
		shuUser := d.GetSkeletonMJUserById(outUser)
		//赢钱的账单
		huUser.AddBill(shuUser.GetUserId(), majiang.MJUSER_BILL_TYPE_YING_HU, "点炮胡牌，获得收入", d.GetYingScore(huUser, hu.GetScore()), hu.Pai, d.GetMJConfig().RoomType)

		//输钱的账单
		shuUser.AddBill(huUser.GetUserId(), majiang.MJUSER_BILL_TYPE_SHU_DIANPAO, "用户点炮，输钱", -d.GetYingScore(huUser, hu.GetScore()), hu.Pai, d.GetMJConfig().RoomType)

		huUser.AddStatisticsCountHu(d.GetMJConfig().CurrPlayCount)       //胡的用户统计信息
		shuUser.AddStatisticsCountDianPao(d.GetMJConfig().CurrPlayCount) //点炮的用户统计信息
	}
}

//获取玩家赢分 为长沙麻将添加的方法
func (d *SkeletonMJDesk) GetYingScore(yingUser *SkeletonMJUser, score int64) int64 {
	return score
}

func (d *SkeletonMJDesk) DoAfterDianPao(hu *majiang.HuPaiInfo) {
	//自摸的不用关
	if hu.GetGetUserId() == hu.GetSendUserId() {
		return
	}
	//点炮胡牌成功之后的处理... 处理checkCase
	d.CheckCase.UpdateCheckBeanStatus(hu.GetGetUserId(), majiang.CHECK_CASE_BEAN_STATUS_CHECKED) // update checkCase...
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
	errDelOut := outUser.GetGameData().GetHandPai().DelOutPai(hu.Pai.GetIndex())
	if errDelOut != nil {
		log.E("胡牌的时候，删除打牌玩家的out牌[%v]...注意:[这里有可能不是错误，一炮多响的情况会多次删除手牌，就可能出现这种情况，等待解决] err[%v]", hu.Pai.GetIndex(), errDelOut)
	}

}

//设置下一次的庄
/**
	1，如果当前的nextBanker 没有值(nextBanker==0)，那代表此人是第一个胡牌的，设置为nextBanekr
	2,如果当前的nextBanker有值(nextBanker > 0 ),那需要判断是不是当前的点炮的人一炮双向
 */
func (d *SkeletonMJDesk) InitNextBanker(hu *majiang.HuPaiInfo) {
	log.T("%v 通过huinfo %v 开始设置下一次的庄的玩家", d.DlogDes(), hu)
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

func (d *SkeletonMJDesk) SendAckActHu(hu *majiang.HuPaiInfo) {
	ack := newProto.NewGame_AckActHu()
	*ack.HuType = hu.GetHuType() //这里需要判断是自摸还是点炮
	*ack.UserIdIn = hu.GetGetUserId()
	*ack.UserIdOut = hu.GetSendUserId()
	ack.HuCard = hu.Pai.GetCardInfo()
	*ack.IsZiMo = (hu.GetGetUserId() == hu.GetSendUserId())
	ack.PaiType = d.IntArry2PaiTypeEnum(hu.PaiType)
	log.T("给用户[%v]广播胡牌的ack[%v]", hu.GetGetUserId(), ack)
	d.BroadCastProto(ack)
}
