package majiang

import (
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_common/common/log"
	"casino_mj_changsha/msg/funcsInit"
	mjproto	"casino_mj_changsha/msg/protogo"
	"errors"
	"github.com/golang/protobuf/proto"
	"strings"
)

func (d *MjDesk) ActQiShouHu(userId uint32, chooseHu bool) error {
	log.T("锁日志: %v ActQiShouHu(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActQiShouHu(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	huUser := d.GetUserByUserId(userId)
	if huUser == nil {
		return Error.NewError(consts.ACK_RESULT_ERROR, "胡牌失败")
	}
	var bean *CheckBean
	for i := 0; i < len(d.CheckCase.GetCheckB()); i++ {
		bean = d.CheckCase.GetCheckB()[i]
		if bean != nil && bean.GetCanHu() && bean.GetUserId() == userId {
			break
		}
	}

	//检测起手胡
	if bean == nil {
		log.E("没有找到起手胡")
		return Error.NewError(consts.ACK_RESULT_ERROR, "胡牌失败")
	}

	//选择过，直接处理下一个
	bean.CheckStatus = proto.Int32(CHECK_CASE_BEAN_STATUS_CHECKED) //设置为已经检测过了
	if !chooseHu {
		d.doQishouHuCheckCase()
		return nil
	}

	//开始处理起手胡牌:记录账单发送胡牌信息

	//放置在 h 测试牌
	canHu, fan, score, huCardStr, paiType, _ := d.HuParser.GetCanHu(huUser.GetGameData().GetHandPai(), nil, true, mjproto.HuType_H_changsha_qishouhu, d.IsBanker(huUser))
	score *= d.GetBaseValue()
	if !canHu {
		log.E("玩家[%v]不可以胡", userId)
		return errors.New("不可以胡牌...")
	}

	//胡牌之后的信息
	hu := NewHuPaiInfo()
	*hu.GetUserId = huUser.GetUserId()
	*hu.SendUserId = userId
	*hu.HuType = int32(mjproto.HuType_H_changsha_qishouhu) ////杠上炮 杠上花 抢杠 海底捞 海底炮 天胡 地胡
	*hu.HuDesc = strings.Join(huCardStr, " ")
	*hu.Fan = fan
	hu.PaiType = PaiTypeEnum2IntArry(paiType)
	*hu.Score = score //只是胡牌的分数，不是赢了多少钱
	huUser.GameData.HuInfo = append(huUser.GameData.HuInfo, hu)
	huUser.SetStatus(MJUSER_STATUS_HUPAI)
	d.SetActiveUser(huUser.GetUserId()) //胡牌之后，设置胡牌的人为activeuser，计算下一家玩家的时候，有用

	//胡牌之后计算账单
	d.doQIshouHuBill(hu) //计算起手胡牌的账单

	//发送起手胡结果
	ack := &mjproto.Game_AckActChangShaQiShouHu{}
	ack.Header = newProto.NewHeader()
	ack.Header.UserId = proto.Uint32(userId)
	ack.HuUserId = proto.Uint32(userId)
	ack.HuType = hu.HuType
	ack.PaiType = paiType
	for _, up := range huUser.GetGameData().GetHandPai().GetPais() {
		if up != nil {
			ack.HandPais = append(ack.HandPais, up.GetCardInfo())
		}
	}
	//如果是庄，需要加上in牌
	if huUser.GetGameData().GetHandPai().GetInPai() != nil {
		ack.HandPais = append(ack.HandPais, huUser.GetGameData().GetHandPai().GetInPai().GetCardInfo())
	}

	d.BroadCastProto(ack)
	//处理下一个起手胡
	d.doQishouHuCheckCase()
	return nil
}

//长沙麻将起手胡牌，单独来操作
func (d *MjDesk) doQIshouHuBill(hu *HuPaiInfo) {
	huUser := d.GetUserByUserId(hu.GetGetUserId())
	log.T("%v玩家[%v]胡牌，开始处理计算起手胡牌分数的逻辑...", d.DlogDes(), huUser.GetUserId())
	//如果是自摸的话，三家都需要给钱
	huUser.AddStatisticsCountZiMo(d.GetCurrPlayCount())
	for _, shuUser := range d.GetUsers() {
		log.T("判断是否计算输的玩家[%v]账单 输的玩家是否为空[%v] 是否在游戏中[%v] 是否不是胡牌的玩家[%v] 是否没有胡牌[%v]", shuUser.GetUserId(), shuUser != nil, shuUser.IsGaming(), shuUser.GetUserId() != huUser.GetUserId(), shuUser.IsNotHu())
		if shuUser != nil && shuUser.IsGaming() && (shuUser.GetUserId() != huUser.GetUserId()) && shuUser.IsNotHu() {
			log.T("开始计算输的玩家[%v]账单", shuUser.GetUserId())
			if banker := d.GetBankerUser(); banker != nil && banker.GetUserId() == huUser.GetUserId() {
				//庄家自摸
				//赢钱的账单
				huUser.AddCSQiShouBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_HU, "用户自摸，获得收入", d.GetYingScore(huUser, hu), hu.Pai, d.GetRoomType())

				//输钱的账单
				shuUser.AddCSQiShouBill(huUser.GetUserId(), MJUSER_BILL_TYPE_SHU_ZIMO, "用户自摸，输钱", -d.GetYingScore(huUser, hu), hu.Pai, d.GetRoomType())
			} else {
				//闲家自摸
				//赢钱的账单
				huUser.AddCSQiShouBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_HU, "用户自摸，获得收入", d.GetYingScore(shuUser, hu), hu.Pai, d.GetRoomType())

				//输钱的账单
				shuUser.AddCSQiShouBill(huUser.GetUserId(), MJUSER_BILL_TYPE_SHU_ZIMO, "用户自摸，输钱", -d.GetYingScore(shuUser, hu), hu.Pai, d.GetRoomType())
			}
			//被自摸的用户没人统计一次
			shuUser.AddStatisticsCountBeiZiMo(d.GetCurrPlayCount())
		}
	}
}
