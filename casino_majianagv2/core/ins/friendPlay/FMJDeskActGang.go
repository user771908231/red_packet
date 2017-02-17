package friendPlay

import (
	"casino_common/common/log"
	"casino_common/common/Error"
	"errors"
	"casino_majiang/service/majiang"
	"casino_majiang/msg/funcsInit"
	"casino_majianagv2/core/api"
	"github.com/golang/protobuf/proto"
)

func (d *FMJDesk) ActGang(userId uint32, paiId int32, bu bool) error {
	log.T("锁日志: %v ActGang(%v,%v)的时候等待锁", d.DlogDes(), userId, paiId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActGang(%v,%v)的时候释放锁", d.DlogDes(), userId, paiId)
	}()

	if d.CheckNotActUser(userId) { //杠牌
		log.E("非法操作，没有轮到玩家[%v]操作杠牌...", userId)
		return Error.NewFailError("暂时没有轮到玩家操作")
	}

	if d.OverTurnTimer != nil {
		d.OverTurnTimer.Stop()
	}

	//检测参数是否正确
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("用户[%v]没有找到杠牌失败...", userId)
		return nil
	}

	gangPai := majiang.InitMjPaiByIndex(int(paiId))
	if gangPai == nil {
		log.E("用户[%v]没有找到杠牌,id[%v]，杠牌失败...", userId, paiId)
		return errors.New("服务器错误,杠牌失败..")
	}

	//判断是否可以杠牌

	var gangType int32 = 0
	var sendUserId uint32 = 0 //打出牌的人，暗杠的话 就表示是自己..
	var canGang bool = false

	if d.CheckCase != nil {
		user.GetGameData().HandPai.InPai = gangPai //把杠牌放进手里
		gangType = majiang.GANG_TYPE_DIAN          //明杠
		sendUserId = d.CheckCase.GetUserIdOut()
		canGang, _ = user.GetGameData().HandPai.GetCanGang(gangPai, d.GetRemainPaiCount())
	} else {
		canGang, _ = user.GetGameData().HandPai.GetCanGang(nil, d.GetRemainPaiCount())
		//如果碰牌中有这张牌表示是巴杠 //如果碰牌中没有这张牌，表示是暗杠
		isBaGang := user.GetGameData().HandPai.IsExistPengPai(gangPai)
		if isBaGang {
			gangType = majiang.GANG_TYPE_BA //巴杠
		} else {
			gangType = majiang.GANG_TYPE_AN //暗杠
		}
		sendUserId = userId
	}
	outUser := d.GetUserByUserId(sendUserId)

	//判断是否可以杠牌
	if !canGang {
		log.E("玩家[%v]杠牌[%v],牌id[%v]失败", userId, gangPai.LogDes(), paiId)
		return errors.New("用户杠牌失败..")
	}

	/**
		根据杠牌的类型不同，处理的方式也不同
		1,如果是巴杠，移除碰牌中的牌和碰牌info , 并且生成杠牌，和杠牌info
		2,如果是明杠或者暗杠，需要把所有的牌放在杠牌中，不用处理碰牌
	 */

	user.GetGameData().HandPai.Pais = append(user.GetGameData().HandPai.Pais, user.GetGameData().HandPai.InPai) //把inpai放进手里
	if gangType == majiang.GANG_TYPE_BA {
		log.T("用户[%v]杠牌是巴杠,现在处理巴杠...", userId)
		//循环碰牌来处理
		user.GetGameData().HandPai.GangPais = append(user.GetGameData().HandPai.GangPais, gangPai)

		var pengKeys []int32
		for _, pengPai := range user.GetGameData().HandPai.PengPais {
			if pengPai != nil && pengPai.GetClientId() == gangPai.GetClientId() {
				//增加杠牌
				user.GetGameData().HandPai.GangPais = append(user.GetGameData().HandPai.GangPais, pengPai)
				pengKeys = append(pengKeys, pengPai.GetIndex())
			}
		}

		//删除碰牌,手中的杠牌
		for _, key := range pengKeys {
			log.T("巴杠删除碰牌牌..index[%v]", key)
			user.GetGameData().HandPai.DelPengPai(key)
		}

		//删除手牌
		user.GetGameData().HandPai.DelHandlPai(gangPai.GetIndex()) //

	} else if gangType == majiang.GANG_TYPE_DIAN || gangType == majiang.GANG_TYPE_AN {
		log.T("用户[%v]杠牌不是巴杠 是 gangType[%v]...", userId, gangType)

		//杠牌的类型
		var gangKey []int32
		//增加杠牌
		//如果不是摸的牌，而是手中本来就有的牌，那么需要把他移除
		for _, pai := range user.GetGameData().HandPai.Pais {
			if pai.GetClientId() == gangPai.GetClientId() {
				//增加杠牌
				user.GetGameData().HandPai.GangPais = append(user.GetGameData().HandPai.GangPais, pai)
				gangKey = append(gangKey, pai.GetIndex())
			}
		}

		log.T("用户杠牌[%v]之后移除需要移除的手牌id数组[%v]", userId, gangKey)
		//减少手中的杠牌
		for _, key := range gangKey {
			log.T("用户杠牌[%v]之后移除手牌id[%v]", userId, key)
			user.GetGameData().HandPai.DelHandlPai(key)

		}

		//如果是点杠，需要删除别人打牌玩家出牌列表里面的这张牌
		if gangType == majiang.GANG_TYPE_DIAN {
			errDelOut := outUser.GetGameData().HandPai.DelOutPai(gangPai.GetIndex())
			if errDelOut != nil {
				log.E("杠牌的时候，删除打牌玩家[%v]的out牌[%v]...", outUser.GetUserId(), gangPai.GetIndex())
			}
		}
	}

	//增加杠牌info
	info := majiang.NewGangPaiInfo()
	*info.GetUserId = user.GetUserId()
	*info.SendUserId = sendUserId
	*info.GangType = gangType
	info.Pai = gangPai
	user.GetGameData().GangInfo = append(user.GetGameData().GangInfo, info)
	user.GetGameData().SetPreMoGangInfo(info) //增加杠牌状态
	user.GetGameData().HandPai.InPai = nil    //1,设置inpai为nil
	//user.StatisticsGangCount(d.GetCurrPlayCount(), gangType)        //处理杠牌的账单
	user.GetGameData().DelGuoHuInfo() //删除过胡的信息

	d.DoGangBill(info)
	d.InitCheckCaseAfterGang(gangType, gangPai, user)

	//杠牌之后的逻辑
	result := newProto.NewGame_AckActGang()
	*result.GangType = user.GetGameData().GetPreMoGangInfo().GetGangType()
	*result.UserIdOut = user.GetGameData().GetPreMoGangInfo().GetSendUserId()
	*result.UserIdIn = user.GetUserId()

	//组装杠牌的信息
	for _, ackpai := range user.GetGameData().HandPai.GangPais {
		if ackpai != nil && ackpai.GetClientId() == gangPai.GetClientId() {
			result.GangCard = append(result.GangCard, ackpai.GetCardInfo())
		}
	}
	log.T("广播玩家[%v]杠牌[%v]之后的ack[%v]", user.GetUserId(), gangPai, result)

	d.BroadCastProto(result)

	//
	d.DoCheckCase() //杠牌之后，处理下一个判定牌

	return nil
}

//初始化checkCase
//如果出错 设置checkCase为nil
func (d *FMJDesk) InitCheckCase(p *majiang.MJPai, outUser api.MjUser) error {

	checkCase := majiang.NewCheckCase()
	checkCase.DianPaoCount = proto.Int32(0) //设置点炮的次数为0
	*checkCase.UserIdOut = outUser.GetUserId()
	*checkCase.CheckStatus = majiang.CHECK_CASE_STATUS_CHECKING //正在判定
	checkCase.CheckMJPai = p
	checkCase.PreOutGangInfo = outUser.GetGameData().GetPreMoGangInfo()
	d.CheckCase = checkCase

	//初始化checkbean
	for _, checkUser := range d.GetFMJUsers() {
		//这里要判断用户是不是已经胡牌
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() {
			log.T("用户[%v]打牌，判断user[%v]是否可以碰杠胡.手牌[%v]", outUser.GetUserId(), checkUser.GetUserId(), checkUser.GameData.HandPai.GetDes())
			//添加checkBean
			bean := checkUser.GetCheckBean(p, d.IsXueLiuChengHe(), d.GetRemainPaiCount())
			if bean != nil {
				checkCase.CheckB = append(checkCase.CheckB, bean)
			}
		}
	}

	log.T("判断最终的checkCase[%v]", checkCase)
	if checkCase.CheckB == nil || len(checkCase.CheckB) == 0 {
		d.CheckCase = nil
	}

	return nil
}

//处理账单
/**
	没有胡牌的人，都需要给钱  ,目前不是承包的方式...
 */

func (d *FMJDesk) DoGangBill(info *majiang.GangPaiInfo) {
	gangType := info.GetGangType()
	gangUser := d.GetFMJUser(info.GetGetUserId())
	gangPai := info.GetPai()

	if gangType == majiang.GANG_TYPE_AN {
		//处理暗杠的账单
		score := d.GetMJConfig().BaseValue * 2 //暗杠的分数
		for _, ou := range d.GetFMJUsers() {
			//不为nil 并且不是本人，并且没有胡牌
			if ou != nil && (ou.GetUserId() != gangUser.GetUserId()) && ou.GetStatus().IsGaming() && ou.GetStatus().IsNotHu() {
				gangUser.AddBill(ou.GetUserId(), majiang.MJUSER_BILL_TYPE_YING_AN_GNAG, "用户暗杠，收入", score, gangPai, d.GetMJConfig().RoomType) //用户赢钱的账户
				ou.AddBill(gangUser.GetUserId(), majiang.MJUSER_BILL_TYPE_SHU_AN_GNAG, "用户暗杠，输钱", -score, gangPai, d.GetMJConfig().RoomType) //用户输钱的账单
				ou.AddStatisticsCountBeiAnGang(d.GetMJConfig().CurrPlayCount)                                                                //被暗杠用户的统计信息
			} else if ou != nil && (ou.GetUserId() == gangUser.GetUserId()) && ou.GetStatus().IsGaming() && ou.GetStatus().IsNotHu() {
				gangUser.AddStatisticsCountAnGang(d.GetMJConfig().CurrPlayCount) //暗杠用户的统计信息
			}
		}

	} else if gangType == majiang.GANG_TYPE_DIAN {
		//处理点杠点账单
		score := d.GetMJConfig().BaseValue * 2 //点杠的分数
		shuUser := d.GetFMJUser(info.GetSendUserId())

		gangUser.AddBill(shuUser.GetUserId(), majiang.MJUSER_BILL_TYPE_YING_GNAG, "用户点杠，收入", score, gangPai, d.GetMJConfig().RoomType) //用户赢钱的账户
		shuUser.AddBill(gangUser.GetUserId(), majiang.MJUSER_BILL_TYPE_SHU_GNAG, "用户点杠，输钱", -score, gangPai, d.GetMJConfig().RoomType) //用户输钱的账单

		gangUser.AddStatisticsCountMingGang(d.GetMJConfig().CurrPlayCount) //明杠用户的统计信息
		shuUser.AddStatisticsCountDianGang(d.GetMJConfig().CurrPlayCount)  //点杠用户的统计信息

	} else if gangType == majiang.GANG_TYPE_BA {
		//处理巴杠的账单
		score := d.GetMJConfig().BaseValue //巴杠的分数
		for _, ou := range d.GetFMJUsers() {
			if ou != nil && (ou.GetUserId() != gangUser.GetUserId()) && ou.GetStatus().IsGaming() && ou.GetStatus().IsNotHu() {
				//账单多次添加
				gangUser.AddBill(ou.GetUserId(), majiang.MJUSER_BILL_TYPE_YING_BA_GANG, "用户巴杠，收入", score, gangPai, d.GetMJConfig().RoomType) //用户赢钱的账户
				ou.AddBill(gangUser.GetUserId(), majiang.MJUSER_BILL_TYPE_SHU_BA_GANG, "用户巴杠，输钱", -score, gangPai, d.GetMJConfig().RoomType) //用户输钱的账单

				ou.AddStatisticsCountBeiBaGang(d.GetMJConfig().CurrPlayCount) //被巴杠用户的统计信息

			} else if ou != nil && (ou.GetUserId() == gangUser.GetUserId()) && ou.GetStatus().IsGaming() && ou.GetStatus().IsNotHu() {
				gangUser.AddStatisticsCountBaGang(d.GetMJConfig().CurrPlayCount) //巴杠用户的统计信息
			}
		}
	}
}
