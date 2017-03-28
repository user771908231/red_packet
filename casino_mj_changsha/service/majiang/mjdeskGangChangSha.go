package majiang

import (
	"casino_common/common/log"
	"casino_common/common/Error"
	"errors"
	"casino_mj_changsha/msg/funcsInit"
	"github.com/golang/protobuf/proto"
)

//杠牌   怎么判断是明杠，暗杠，巴杠...
func (d *MjDesk) ActGang(userId uint32, paiId int32, buPai bool) error {
	log.T("锁日志: %v ActGang(%v,%v,%v)的时候等待锁", d.DlogDes(), userId, paiId, buPai)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActGang(%v,%v,%v)的时候释放锁", d.DlogDes(), userId, paiId, buPai)
	}()

	err := d.CheckActUser(userId, ACTTYPE_GANG)
	if err != nil {
		log.E("非法操作，没有轮到玩家[%v]操作杠牌...", userId)
		return Error.NewFailError("暂时没有轮到玩家操作")
	}

	if d.overTurnTimer != nil {
		d.overTurnTimer.Stop()
	}
	//检测参数是否正确
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("用户[%v]没有找到杠牌失败...", userId)
		return nil
	}

	gangPai := d.getPaiById(paiId)
	if gangPai == nil {
		log.E("用户[%v]没有找到杠牌,id[%v]，杠牌失败...", userId, paiId)
		return errors.New("服务器错误,杠牌失败..")
	}

	//判断是否可以杠牌

	var gangType int32 = 0
	var sendUserId uint32 = 0 //打出牌的人，暗杠的话 就表示是自己..
	var canGang bool = false

	if d.CheckCase != nil {
		user.GameData.HandPai.InPai = gangPai //把杠牌放进手里
		gangType = GANG_TYPE_DIAN             //明杠
		sendUserId = d.CheckCase.GetUserIdOut()
		canGang, _ = user.GameData.HandPai.GetCanGang(gangPai, d.GetRemainPaiCount())
	} else {
		log.T("user[%v]长沙杠的时候判断in[%v]牌", userId, user.GetGameData().GetHandPai().GetInPai().LogDes())
		canGang, _ = user.GameData.HandPai.GetCanGang(nil, d.GetRemainPaiCount())
		//如果碰牌中有这张牌表示是巴杠 //如果碰牌中没有这张牌，表示是暗杠
		isBaGang := user.GameData.HandPai.IsExistPengPai(gangPai)
		if isBaGang {
			gangType = GANG_TYPE_BA //巴杠
		} else {
			gangType = GANG_TYPE_AN //暗杠
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
		3,如果是长沙麻将，inpai可能在为空的情况下也能补牌，所以这里需要判断 inpai为空的情况
	 */
	if user.GameData.HandPai.InPai != nil {
		user.GameData.HandPai.Pais = append(user.GameData.HandPai.Pais, user.GameData.HandPai.InPai) //把inpai放进手里
	}

	if gangType == GANG_TYPE_BA {
		log.T("用户[%v]杠牌是巴杠,现在处理巴杠...", userId)
		//循环碰牌来处理
		user.GameData.HandPai.GangPais = append(user.GameData.HandPai.GangPais, gangPai)

		var pengKeys []int32
		for _, pengPai := range user.GameData.HandPai.PengPais {
			if pengPai != nil && pengPai.GetClientId() == gangPai.GetClientId() {
				//增加杠牌
				user.GameData.HandPai.GangPais = append(user.GameData.HandPai.GangPais, pengPai)
				pengKeys = append(pengKeys, pengPai.GetIndex())
			}
		}

		//删除碰牌,手中的杠牌
		for _, key := range pengKeys {
			log.T("巴杠删除碰牌牌..index[%v]", key)
			user.GameData.HandPai.DelPengPai(key)
		}

		//删除手牌
		user.GameData.HandPai.DelHandlPai(gangPai.GetIndex()) //

	} else if gangType == GANG_TYPE_DIAN || gangType == GANG_TYPE_AN {
		log.T("用户[%v]杠牌不是巴杠 是 gangType[%v]...", userId, gangType)

		//杠牌的类型
		var gangKey []int32
		//增加杠牌
		//如果不是摸的牌，而是手中本来就有的牌，那么需要把他移除
		for _, pai := range user.GameData.HandPai.Pais {
			if pai.GetClientId() == gangPai.GetClientId() {
				//增加杠牌
				user.GameData.HandPai.GangPais = append(user.GameData.HandPai.GangPais, pai)
				gangKey = append(gangKey, pai.GetIndex())
			}
		}

		log.T("用户杠牌[%v]之后移除需要移除的手牌id数组[%v]", userId, gangKey)
		//减少手中的杠牌
		for _, key := range gangKey {
			user.GameData.HandPai.DelHandlPai(key)
		}

		//如果是点杠，需要删除别人打牌玩家出牌列表里面的这张牌
		if gangType == GANG_TYPE_DIAN {
			errDelOut := outUser.GameData.HandPai.DelOutPai(gangPai.GetIndex())
			if errDelOut != nil {
				log.E("杠牌的时候，删除打牌玩家[%v]的out牌[%v]...", outUser.GetUserId(), gangPai.GetIndex())
			}
		}
	}

	//增加杠牌info
	info := NewGangPaiInfo()
	*info.GetUserId = user.GetUserId()
	*info.SendUserId = sendUserId
	info.Bu = proto.Bool(buPai)
	*info.GangType = gangType
	info.Pai = gangPai
	user.GameData.GangInfo = append(user.GameData.GangInfo, info)
	user.PreMoGangInfo = info         //增加杠牌状态
	user.GameData.HandPai.InPai = nil //1,设置inpai为nil
	user.DelGuoHuInfo()               //2,长沙麻将，杠牌之后 删除过胡的信息
	user.changshaGang = !buPai        //是否是杠牌
	log.T("%v玩家%v杠牌之后..user.changshaGang: %v", d.DlogDes(), userId, user.GetChangShaGangStatus())

	d.InitCheckCaseAfterGang(gangType, gangPai, user) //长沙杠牌之后 初始化checkCase

	//杠牌之后的逻辑
	result := newProto.NewGame_AckActGang()
	*result.GangType = user.PreMoGangInfo.GetGangType()
	*result.UserIdOut = user.PreMoGangInfo.GetSendUserId()
	*result.UserIdIn = user.GetUserId()
	result.ChangShaBu = proto.Bool(buPai) //判断长沙麻将是否补牌

	//组装杠牌的信息
	for _, ackpai := range user.GameData.HandPai.GangPais {
		if ackpai != nil && ackpai.GetClientId() == gangPai.GetClientId() {
			result.GangCard = append(result.GangCard, ackpai.GetCardInfo())
		}
	}
	log.T("广播玩家[%v]杠牌[%v]之后的ack[%v]", user.GetUserId(), gangPai, result)
	d.BroadCastProto(result)

	//长沙麻将特殊处理
	/**
		1，如果有抢杠的操作，可以直接操作抢杠
		2，如果没有抢杠，补的情况和以前一样，杠牌的情况 需要摸两张牌
	 */
	d.DoCheckCase() //长沙杠之后处理checkCase
	return nil
}
