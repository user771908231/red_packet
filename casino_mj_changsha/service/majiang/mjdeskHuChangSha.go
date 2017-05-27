package majiang

import (
	"strings"
	mjproto "casino_mj_changsha/msg/protogo"
	"casino_common/common/Error"
	"casino_common/common/log"
	"errors"
	"casino_mj_changsha/msg/funcsInit"
	"github.com/golang/protobuf/proto"
)

//长沙麻将的杠上花单独做
func (d *MjDesk) ActHu(userId uint32) error {
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActHu(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//对于杠，有摸牌前 杠的状态，有打牌前杠的状态

	err := d.CheckActUser(userId, ACTTYPE_HU)
	if err != nil {
		return Error.NewFailError("不该当前用户操作..")
	}
	if d.overTurnTimer != nil {
		d.overTurnTimer.Stop()
	}

	//玩家胡牌
	huUser := d.GetUserByUserId(userId)
	if huUser == nil {
		log.E("服务器错误：没有找到胡牌的user[%v]", userId)
		return errors.New("服务器错误，没有找到胡牌的user")
	}

	//1,胡的牌是当前check里面的牌，如果没有check，则表示是自摸
	//得到胡牌的信息
	var isZimo bool = false                                                                //是否是自摸
	var isHaidilaoyue bool = (huUser.GetNeedHaidiStatus() == MJUSER_NEEDHAIDI_STATUS_NEED) //是否是最后一张牌
	var isHaiDiPao bool                                                                    //海底炮
	var extraAct mjproto.HuType = 0                                                        //杠上花，
	var outUserId uint32                                                                   //打牌的人
	var huPai []*MJPai = make([]*MJPai, 2)                                                 //胡牌的麻将
	var huInfos []*HuPaiInfo                                                               //胡牌的信息
	var gangshanghua bool = false                                                          //是否是杠上花
	var qiangGnag = false                                                                  //判断是否是抢杠

	//判断胡的牌是别人打的还是自己摸的...
	if d.GetCheckCase() != nil {
		outUserId = d.GetCheckCase().GetUserIdOut()
		huPai[0] = d.GetCheckCase().GetCheckMJPai()  //胡牌1
		huPai[1] = d.GetCheckCase().GetCheckMJPai2() //胡牌2
		qiangGnag = d.GetCheckCase().GetQiangGang()  //是否是抢杠
		if d.GetUserByUserId(outUserId).GetNeedHaidiStatus() == MJUSER_NEEDHAIDI_STATUS_NEED {
			isHaiDiPao = true
		}
	} else {
		isZimo = true                                            //自摸
		huPai[0] = huUser.GetGameData().GetHandPai().GetInPai()  //胡牌1
		huPai[1] = huUser.GetGameData().GetHandPai().GetInPai2() //胡牌2
		outUserId = userId
	}

	//胡牌之前的杠牌信息

	/**
		是否是杠上炮的逻辑:
		1,不是自摸
		2，checkCase 需要有preOutGangInfo
		3, preOutGangInfo 里不是补牌，而是杠牌的类型
	 */

	log.T("%v 玩家%v开始胡牌，检测胡牌方式::海底捞月:%v,海底炮:%v,自摸:%v,自摸前杠：%v,点炮点杠:%v",
		d.DlogDes(), userId, isHaidilaoyue, isHaiDiPao, isZimo, huUser.GetPreMoGangInfo(), d.GetCheckCase().GetPreOutGangInfo())

	if isHaidilaoyue || isHaiDiPao {
		extraAct = mjproto.HuType_H_HaiDiLao //海底捞
		//这里判断海底捞月还是海底炮
		//if isHaidilaoyue ==> 海底捞月;if
	} else if qiangGnag {
		extraAct = mjproto.HuType_H_QiangGang //长沙麻将需要增加抢杠的逻辑	// 注意这里必须要优先判断抢杠(不抢杠的时候，也能识别杠上炮)
	} else if isZimo && huUser.GetPreMoGangInfo() != nil {
		gangshanghua = true
		extraAct = mjproto.HuType_H_GangShangHua
	} else if !isZimo && d.GetCheckCase().GetPreOutGangInfo() != nil && !d.GetCheckCase().GetPreOutGangInfo().GetBu() {
		//如果不是自摸，并且打牌前拥有杠牌的info，并且不是补牌，那么就是杠上炮
		extraAct = mjproto.HuType_H_GangShangPao
	}

	//开始胡
	for _, h := range huPai {
		if h == nil {
			continue
		}
		canHu, fan, score, huCardStr, paiType, _ := d.HuParser.GetCanHu(huUser.GetGameData().GetHandPai(), h, isZimo, extraAct, d.IsBanker(huUser))
		score *= d.GetBaseValue()
		if !canHu {
			log.E("%d玩家[%v]不可以胡牌[%v],胡牌失败...", userId, h.LogDes())
			continue
		} else {
			//如果客户胡牌，处理胡牌的信息
			if huUser.GetGameData().GetHandPai().GetInPai().GetIndex() == h.GetIndex() {
				huUser.GetGameData().GetHandPai().InPai = nil
			}
			if huUser.GetGameData().GetHandPai().GetInPai2().GetIndex() == h.GetIndex() {
				huUser.GetGameData().GetHandPai().InPai2 = nil
			}

			hu := NewHuPaiInfo()
			*hu.GetUserId = huUser.GetUserId()
			*hu.SendUserId = outUserId
			//*hu.ByWho = 打牌的方位，对家，上家，下家？
			*hu.HuType = int32(extraAct) ////杠上炮 杠上花 抢杠 海底捞 海底炮 天胡 地胡
			*hu.HuDesc = strings.Join(huCardStr, " ");
			*hu.Fan = fan
			*hu.Score = score //只是胡牌的分数，不是赢了多少钱
			hu.Pai = h
			hu.PaiType = PaiTypeEnum2IntArry(paiType)
			huUser.AddHuPaiInfo(hu)             //胡牌之后，设置用户的数据
			d.SetActiveUser(huUser.GetUserId()) //长沙麻将胡牌之后，设置胡牌的人为activeuser，计算下一家玩家的时候，有用

			//处理抢杠的逻辑
			d.changShaDoQiangGang(hu)

			//胡牌之后计算账单
			d.DoHuBill(hu)

			//点炮之后设置checkCase的状态
			d.DoAfterDianPao(hu)

			//设置下一次的庄
			d.InitNextBanker(hu)

			//发送胡牌成功的回复
			huInfos = append(huInfos, hu)
		}
	}

	//胡牌结果
	if gangshanghua && len(huInfos) > 0 {
		//发送结果 有可能胡两张牌
		ack := &mjproto.Game_AckActHuChangSha{}
		for _, huinfo := range huInfos {
			if huinfo != nil {
				a := newProto.NewGame_AckActHu()
				*a.HuType = huinfo.GetHuType()        //这里需要判断是自摸还是点炮
				*a.UserIdIn = huinfo.GetGetUserId()   //胡牌的人
				*a.UserIdOut = huinfo.GetSendUserId() //打牌的人
				a.HuCard = huinfo.Pai.GetCardInfo()   //胡的牌
				a.IsZiMo = proto.Bool(true)           //杠上花都是自摸
				ack.Hus = append(ack.Hus, a)
				log.T("给用户[%v]广播胡牌的ack[%v]", huinfo.GetGetUserId(), a)
			}
		}
		d.BroadCastProto(ack)
	} else {
		d.SendAckActHu(huInfos[0])
	}

	//胡牌之后，需要判断游戏是否结束...
	if d.Time2Lottery() {
		//倒倒胡 某一玩家胡牌即结束
		return d.Lottery() //长沙胡牌之后判断是否lottery
	} else {
		//处理下一个

		return d.DoCheckCase() //胡牌之后，处理下一个判定牌
	}
	return nil
}

//处理长沙麻将抢杠的逻辑
func (d *MjDesk) changShaDoQiangGang(hu *HuPaiInfo) error {
	//不是抢杠的胡法直接返回
	if hu.GetHuType() != int32(mjproto.HuType_H_QiangGang) {
		return nil
	}

	//开始处理抢杠的逻辑
	log.T("%v 由于CheckCase :%v 开始处理长沙麻将抢杠的逻辑....", d.DlogDes(), d.GetCheckCase())
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
	dianUser.PreMoGangInfo = nil //摸牌前的杠牌处理也要删除
	return nil
}
