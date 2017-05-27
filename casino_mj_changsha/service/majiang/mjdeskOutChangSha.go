package majiang

import (
	"casino_common/common/log"
	"casino_mj_changsha/msg/funcsInit"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/util"
	mjproto        "casino_mj_changsha/msg/protogo"
)

//用户打一张牌出来
func (d *MjDesk) ActOutChangSha(userId uint32, paiKey int32, isAuto bool) error {
	outUser := d.GetUserByUserId(userId)
	if outUser == nil {
		log.E("[%v]打牌失败，没有找到玩家[%v]", d.DlogDes(), userId)
		return ERR_OUTPAI
	}

	//判断是否轮到当前玩家打牌了...
	err := d.CheckActUser(userId, ACTTYPE_OUT)
	if err != nil {
		log.E("%v没有轮到当前玩家%v打牌", d.DlogDes(), userId)
		return ERR_OUTPAI
	}

	//判断是不是在游戏中的状态
	if d.IsNotGaming() {
		log.E("%v玩家%v打牌失败，desk不在游戏状态[%v]", d.DlogDes(), userId, d.IsNotGaming())
		return ERR_OUTPAI
	}

	//停止定时器
	if d.overTurnTimer != nil {
		d.overTurnTimer.Stop()
	}

	//首先删除过胡的信息
	outUser.DelGuoHuInfo() //长沙麻将 打牌之后 删除过胡的信息

	/**
		长沙麻将打牌分为两种，普通打牌，杠牌打牌
		1，普通打牌和其他的打牌一致
		2，杠牌打牌只能打摸的两张牌
	 */

	if outUser.GameData.HandPai.GetInPai2() != nil {
		//杠牌打牌
		outPai1 := outUser.GameData.HandPai.GetInPai()
		outPai2 := outUser.GameData.HandPai.GetInPai2()

		//判断过胡
		canhu1, _, _, _, _, _ := d.HuParser.GetCanHu(outUser.GetGameData().GetHandPai(), outPai1, true, mjproto.HuType_H_NORMAL, d.IsBanker(outUser))
		canhu2, _, _, _, _, _ := d.HuParser.GetCanHu(outUser.GetGameData().GetHandPai(), outPai2, true, mjproto.HuType_H_NORMAL, d.IsBanker(outUser))
		if canhu1 || canhu2 {
			//增加过胡的信息
			guoHuInfo := NewGuoHuInfo()
			guoHuInfo.SendUserId = proto.Uint32(outUser.GetUserId())
			guoHuInfo.GetUserId = proto.Uint32(outUser.GetUserId())
			guoHuInfo.Pai = outUser.GetGameData().GetHandPai().GetInPai()
			guoHuInfo.FanShu = proto.Int32(0) //现在都设置为0翻
			outUser.GetGameData().GuoHuInfo = append(outUser.GetGameData().GuoHuInfo, guoHuInfo)
			if canhu1 {
				guoHuInfo.Pai = outPai1
			}
			if canhu2 {
				guoHuInfo.Pai = outPai2
			}
		}

		//打牌
		outUser.GameData.HandPai.OutPais = append(outUser.GameData.HandPai.OutPais, outPai1, outPai2) //自己桌子前面打出的牌，如果其他人碰杠胡了之后，需要把牌删除掉...
		outUser.GameData.HandPai.InPai = nil                                                          //打牌之后需要把自己的  inpai给移除掉...
		outUser.GameData.HandPai.InPai2 = nil                                                         //打牌之后需要把自己的  inpai2给移除掉...
		//初始化Init
		err := d.InitCheckCaseAfterChangShaGang(outPai1, outUser) //打牌之后
		if err != nil {
			//表示无人需要，直接给用户返回无人需要
			//给下一个人摸排，并且移动指针
			log.E("%v服务器错误，初始化判定牌的时候出错err[%v]", d.DlogDes(), err)
			return ERR_OUTPAI
		}
		checkCase1 := d.GetCheckCase()

		err = d.InitCheckCaseAfterChangShaGang(outPai2, outUser) //打牌之后
		if err != nil {
			//表示无人需要，直接给用户返回无人需要
			//给下一个人摸排，并且移动指针
			log.E("%v服务器错误，初始化判定牌的时候出错err[%v]", d.DlogDes(), err)
			return ERR_OUTPAI
		}

		checkCase2 := d.GetCheckCase()
		//合并两个checkCase
		var checkCaseEnd *CheckCase
		if checkCase1 == nil && checkCase2 != nil {
			checkCaseEnd = checkCase2
		} else if checkCase1 != nil && checkCase2 == nil {
			checkCaseEnd = checkCase1
		} else if checkCase1 != nil && checkCase2 != nil {
			//合并两个checkCase
			checkCaseEnd = checkCase1
			checkCaseEnd.CheckMJPai2 = checkCase2.CheckMJPai2 //这里报错了
			checkCaseEnd.CheckB = append(checkCaseEnd.CheckB, checkCase2.CheckB...)
		}
		d.CheckCase = checkCaseEnd  //初始化checkCase完成
		outUser.PreMoGangInfo = nil //清楚摸牌前的杠牌info

		result := newProto.NewGame_AckSendOutCard() //长沙麻将，杠牌之后打牌
		*result.UserId = userId
		result.Card = outPai1.GetCardInfo()
		result.Card2 = outPai2.GetCardInfo()
		result.IsAuto = proto.Bool(isAuto)
		d.BroadCastProto(result)
		log.T("[%v]用户[%v]已经打牌结束，开始处理下一个checkCase", d.DlogDes(), userId)
		d.DoCheckCase() //打牌之后，别人判定牌
		return nil

	} else {
		//得到参数普通打牌
		outPai := InitMjPaiByIndex(int(paiKey))

		/**
			1,如果是碰牌打牌的时候,inpai为nil，不需要增加
			2,如果是摸牌打牌（杠之后也是摸牌，需要增加in牌...）
		 */
		if outUser.GameData.HandPai.InPai != nil {
			outUser.GameData.HandPai.AddPai(outUser.GameData.HandPai.InPai) //把inpai放置到手牌上
		}

		//打牌之前需要判断当前是否可以胡牌了...如果可以胡牌的话，那么添加过胡的信息
		tempMjHand := &MJHandPai{}
		util.DeepCopy(tempMjHand, outUser.GetGameData().GetHandPai())
		tempMjHand.DelHandlPai(paiKey)
		if can, fan, _, _, _, _ := d.HuParser.GetCanHu(tempMjHand, outPai, true, mjproto.HuType_H_NORMAL, d.IsBanker(outUser)); can == true {
			log.T("玩家[%v]可以胡牌[%v]，但是没有胡，增加过胡,[%v]翻的信息", outUser.GetUserId(), outPai.LogDes(), fan)
			//如果能胡牌的话,增加过胡
			guoHuInfo := NewGuoHuInfo()
			*guoHuInfo.SendUserId = outUser.GetUserId()
			guoHuInfo.Pai = outPai
			guoHuInfo.FanShu = proto.Int32(fan) //长沙麻将 过胡的时候需要判断翻数
			outUser.GetGameData().GuoHuInfo = append(outUser.GetGameData().GuoHuInfo, guoHuInfo)
		}

		log.T("%v玩家[%v]打牌之前的手牌:%v", d.DlogDes(), outUser.GetUserId(), ServerPais2string(outUser.GameData.HandPai.Pais))
		errDapai := outUser.GameData.HandPai.DelHandlPai(outPai.GetIndex())
		if errDapai != nil {
			log.E("[%v]打牌的时候出现错误，没有找到要到的牌,id[%v]", d.DlogDes(), outPai.GetIndex())
			return ERR_OUTPAI
		}
		log.T("%v玩家[%v]打牌之后的手牌:%v", d.DlogDes(), outUser.GetUserId(), ServerPais2string(outUser.GameData.HandPai.Pais))

		outUser.GameData.HandPai.OutPais = append(outUser.GameData.HandPai.OutPais, outPai) //自己桌子前面打出的牌，如果其他人碰杠胡了之后，需要把牌删除掉...
		outUser.GameData.HandPai.InPai = nil                                                //打牌之后需要把自己的  inpai给移除掉...

		//打牌之后的逻辑,初始化判定事件
		err := d.InitCheckCase(outPai, outUser, false) //打牌之后
		if err != nil {
			//表示无人需要，直接给用户返回无人需要
			//给下一个人摸排，并且移动指针
			log.E("%v服务器错误，初始化判定牌的时候出错err[%v]", d.DlogDes(), err)
			return ERR_OUTPAI
		}

		log.T("[%v]玩家[%v]打牌之后InitCheckCase之后的checkCase[%v]", d.DlogDes(), userId, d.CheckCase)
		//回复消息,打牌之后，广播打牌的信息...s
		outUser.PreMoGangInfo = nil //清楚摸牌前的杠牌info
		result := newProto.NewGame_AckSendOutCard()
		*result.UserId = userId
		result.Card = outPai.GetCardInfo()
		result.IsAuto = proto.Bool(isAuto) //是否自动打牌
		d.BroadCastProto(result)

		log.T("[%v]用户[%v]已经打牌结束，开始处理下一个checkCase", d.DlogDes(), userId)
		d.DoCheckCase() //打牌之后，别人判定牌
		return nil
	}
	return nil
}

//这里需要分开计算，因为需要分开询问 胡，碰/杠，吃
func (d *MjDesk) InitCheckCaseAfterChangShaGang(p *MJPai, outUser *MjUser) error {

	checkCase := NewCheckCase()
	checkCase.DianPaoCount = proto.Int32(0) //InitCheckCaseAfterChangShaGang
	*checkCase.UserIdOut = outUser.GetUserId()
	*checkCase.CheckStatus = CHECK_CASE_STATUS_CHECKING //正在判定
	checkCase.CheckMJPai = p
	checkCase.PreOutGangInfo = outUser.GetPreMoGangInfo()
	d.CheckCase = checkCase

	//更具长沙麻将的规则，先初始化可以胡的bean
	for _, checkUser := range d.GetUsers() {
		//这里要判断用户是不是已经胡牌
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() {
			log.T("用户[%v]打牌，判断user[%v]是否可以胡.手牌[%v]", outUser.GetUserId(), checkUser.GetUserId(), checkUser.GameData.HandPai.GetDes())
			//添加checkBean
			bean := NewCheckBean()
			*bean.CheckStatus = CHECK_CASE_BEAN_STATUS_CHECKING
			*bean.UserId = checkUser.GetUserId()
			bean.CheckPai = p
			fan := int32(0)
			//是否可以胡牌
			if checkUser.IsCanInitCheckCaseHu(d.IsXueLiuChengHe()) {
				//*bean.CanHu, _ = checkUser.GameData.HandPai.GetCanHu()
				*bean.CanHu, fan, _, _, _, _ = d.HuParser.GetCanHu(checkUser.GetGameData().GetHandPai(), p, false, 0, d.IsBanker(checkUser))
			}

			log.T("得到用户[%v]对牌[%v]的check , bean[%v]", checkUser.GetUserId(), p.LogDes(), bean)
			//判断过胡.如果有过胡，那么就不能再胡了
			if checkUser.HadGuoHuInfo(fan) {
				*bean.CanHu = false
			}

			if bean.GetCanHu() {
				checkCase.CheckB = append(checkCase.CheckB, bean)
			}
		}
	}

	//再初始化可以碰，杠的bean
	for _, checkUser := range d.GetUsers() {
		//这里要判断用户是不是已经胡牌
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() {

			log.T("用户[%v]打牌，判断user[%v]是否可以碰杠胡.手牌[%v]", outUser.GetUserId(), checkUser.GetUserId(), checkUser.GameData.HandPai.GetDes())
			checkUser.GameData.HandPai.InPai = p
			//添加checkBean

			bean := NewCheckBean()
			*bean.CheckStatus = CHECK_CASE_BEAN_STATUS_CHECKING
			*bean.UserId = checkUser.GetUserId()
			bean.CheckPai = p

			//是否可以杠
			if checkUser.IsCanInitCheckCaseGang(d.IsXueLiuChengHe()) {
				*bean.CanGang, _ = checkUser.GameData.HandPai.GetCanGang(p, d.GetRemainPaiCount())
			}

			//是否可以碰
			if checkUser.IsCanInitCheckCasePeng() {
				*bean.CanPeng = checkUser.GameData.HandPai.GetCanPeng(p)
			}

			if bean.GetCanGang() || bean.GetCanPeng() {
				checkCase.CheckB = append(checkCase.CheckB, bean)
			}

			checkUser.GameData.HandPai.InPai = nil //判定之后把inpai设置为空...

		}
	}
	//再初始化可以吃的bean

	//初始化checkbean
	for _, checkUser := range d.GetUsers() {
		//这里要判断用户是不是已经胡牌
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() {
			log.T("用户[%v]打牌，判断user[%v]是否可以碰杠胡.手牌[%v]", outUser.GetUserId(), checkUser.GetUserId(), checkUser.GameData.HandPai.GetDes())
			checkUser.GameData.HandPai.InPai = p
			//添加checkBean
			bean := NewCheckBean()
			*bean.CheckStatus = CHECK_CASE_BEAN_STATUS_CHECKING
			*bean.UserId = checkUser.GetUserId()
			bean.CheckPai = p

			//是否可以吃牌
			if checkUser.IsCanInitCheckCaseChi() {
				*bean.CanChi, bean.ChiCards = checkUser.GameData.HandPai.GetCanChi(p)
			}
			if bean.GetCanChi() {
				checkCase.CheckB = append(checkCase.CheckB, bean)
			}

			checkUser.GameData.HandPai.InPai = nil //判定之后把inpai设置为空...
		}
	}

	log.T("判断最终的checkCase[%v]", checkCase)
	if checkCase.CheckB == nil || len(checkCase.CheckB) == 0 {
		d.CheckCase = nil
	}
	return nil
}
