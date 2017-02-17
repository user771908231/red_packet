package skeleton

import (
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"casino_majiang/service/majiang"
	"sync/atomic"
	"casino_common/common/log"
	"errors"
	"casino_common/common/userService"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"github.com/name5566/leaf/gate"
)

func (u *SkeletonMJUser) GetUserId() uint32 {
	return u.userId
}

func (u *SkeletonMJUser) GetNickName() string {
	return u.NickName

}

func (u *SkeletonMJUser) GetStatus() *data.MjUserStatus {
	return u.UserStatus
}

//todo
func (u *SkeletonMJUser) GetDesk() api.MjDesk {
	return u.desk
}

/***************************************账单相关***************************************/
func (u *SkeletonMJUser) GetBill() *majiang.Bill {
	return u.Bill
}

func (u *SkeletonMJUser) SubBillAmount(amount int64) {
	atomic.AddInt64(u.GetBill().WinAmount, -amount)
}

func (u *SkeletonMJUser) AddBillAmount(amount int64) {
	atomic.AddInt64(u.GetBill().WinAmount, amount)
}

//删除账单
func (u *SkeletonMJUser) DelBillBean(pai *majiang.MJPai) (error, *majiang.BillBean) {
	var bean *majiang.BillBean
	index := -1
	for i, info := range u.GetBill().GetBills() {
		if info != nil && info.GetPai().GetIndex() == pai.GetIndex() {
			index = i
			bean = info
			break
		}
	}

	if index > -1 {
		u.Bill.Bills = append(u.Bill.Bills[:index], u.Bill.Bills[index+1:]...)
		u.SubBillAmount(bean.GetAmount()) //减去
		return nil, bean
	} else {
		log.E("服务器错误：删除账单 billBean的时候出错，没有找到对应的杠牌[%v]", pai)
		return errors.New("删除手牌时出错，没有找到对应的手牌..."), nil
	}
}

func (u *SkeletonMJUser) GetGameData() *data.MJUserGameData {
	return u.GameData
}

//增加一条账单
func (u *SkeletonMJUser) AddBillBean(bean *majiang.BillBean) error {
	u.Bill.Bills = append(u.Bill.Bills, bean)
	u.AddBillAmount(bean.GetAmount())
	return nil
}

//增加账单
func (u *SkeletonMJUser) AddBill(relationUserid uint32, billType int32, des string, score int64, pai *majiang.MJPai, roomType int32) error {
	//用户赢钱的账户,赢钱的账单
	bill := majiang.NewBillBean()
	*bill.UserId = u.GetUserId()
	*bill.OutUserId = relationUserid
	//*bill.Type = MJUSER_BILL_TYPE_YING_HU
	*bill.Type = billType
	*bill.Des = des
	*bill.Amount = score //杠牌的收入金额
	bill.Pai = pai
	u.AddBillBean(bill)

	//计算账单的地方 来加减用户的coin
	u.AddCoin(score, roomType)    //统计用户剩余多少钱
	u.AddStatisticsWinCoin(score) //统计用户输赢多少钱
	return nil
}

func (u *SkeletonMJUser) AddStatisticsWinCoin(coin int64) {
	atomic.AddInt64(u.Statisc.WinCoin, coin)
}

//更新用户金币
func (u *SkeletonMJUser) AddCoin(coin int64, roomType int32) {
	//减少玩家金额
	atomic.AddInt64(&u.Coin, coin) //更新账户余额
	//如果是金币场。需要更新用户的金币余额
	if roomType == majiang.ROOMTYPE_COINPLAY {
		remainCoin, _ := userService.INCRUserCOIN(u.GetUserId(), coin)
		u.Coin = remainCoin //增加用户的金币
	}
}

func (u *SkeletonMJUser) GetSkeletonMJDesk() *SkeletonMJDesk {
	return u.GetDesk().(*SkeletonMJDesk)
}

func (u *SkeletonMJUser) GetSkeletonMJUser() *SkeletonMJUser {
	return u
}
func (u *SkeletonMJUser) GetPlayerCard(showHand bool, needInpai bool) *mjproto.PlayerCard {
	//log.T("得到玩家[%v]的手牌,showHand[%v],needInpai[%v]", u.GetUserId(), showHand, needInpai)
	playerCard := newProto.NewPlayerCard()
	*playerCard.UserId = u.GetUserId()

	//得到inpai
	if needInpai {
		playerCard.HandCardCount = proto.Int32(1)
		if showHand {
			playerCard.HandCard = append(playerCard.HandCard, u.GameData.HandPai.InPai.GetCardInfo())
		} else {
			//直接不返回
			playerCard.HandCard = append(playerCard.HandCard, u.GameData.HandPai.InPai.GetBackPai())
		}
	}

	//手牌的长度
	playerCard.HandCardCount = proto.Int32(playerCard.GetHandCardCount() + int32(len(u.GameData.HandPai.Pais)))
	//得到手牌
	for _, pai := range u.GameData.HandPai.GetPais() {
		if pai != nil {
			if showHand {
				playerCard.HandCard = append(playerCard.HandCard, pai.GetCardInfo())
			} else {
				playerCard.HandCard = append(playerCard.HandCard, pai.GetBackPai())
			}
		}
	}

	//是否显示换三张的牌

	//类型（1,碰，2,明杠，3,暗杠）
	//得到碰牌
	for i, pai := range u.GameData.HandPai.GetPengPais() {
		if pai != nil && i%3 == 0 {
			com := newProto.NewComposeCard()
			*com.Value = pai.GetClientId()
			*com.Type = int32(mjproto.ComposeCardType_C_PENG) //todo ,需要把type 放置在常量里面 这里代表的是碰牌
			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}

	//得到吃的牌
	for i := 0; i < len(u.GetGameData().GetHandPai().GetChiPais()); i += 3 {
		com := newProto.NewComposeCard()
		*com.Type = int32(mjproto.ComposeCardType_C_CHI) //todo ,需要把type 放置在常量里面 这里代表的是碰牌
		com.ChiValue = append(com.ChiValue, u.GetGameData().GetHandPai().GetChiPais()[i].GetClientId())
		com.ChiValue = append(com.ChiValue, u.GetGameData().GetHandPai().GetChiPais()[i+1].GetClientId())
		com.ChiValue = append(com.ChiValue, u.GetGameData().GetHandPai().GetChiPais()[i+2].GetClientId())
		playerCard.ComposeCard = append(playerCard.ComposeCard, com)
	}

	//得到杠牌
	for _, info := range u.GameData.GangInfo {
		if info != nil {
			com := newProto.NewComposeCard()
			*com.Value = info.GetPai().GetClientId()

			if info.GetGangType() == majiang.GANG_TYPE_DIAN {
				*com.Type = int32(mjproto.ComposeCardType_C_MINGGANG) // 明杠
			} else if info.GetGangType() == majiang.GANG_TYPE_AN {
				*com.Type = int32(mjproto.ComposeCardType_C_ANGANG) // 暗杠
			} else if info.GetGangType() == majiang.GANG_TYPE_BA {
				*com.Type = int32(mjproto.ComposeCardType_C_BAGANG) // 巴杠
			}
			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}

	//得到胡牌
	for _, pai := range u.GameData.HandPai.GetHuPais() {
		if pai != nil {
			playerCard.HuCard = append(playerCard.HuCard, pai.GetCardInfo())
		}
	}

	//打出去的牌
	for _, pai := range u.GameData.HandPai.GetOutPais() {
		if pai != nil {
			playerCard.OutCard = append(playerCard.OutCard, pai.GetCardInfo())
		}
	}

	return playerCard
}

func (u *SkeletonMJUser) WriteMsg(p proto.Message) error {
	u.a.WriteMsg(p)
	return nil
}

func (u *SkeletonMJUser) GetWaitTime() int64 {
	return 30
}

func (u *SkeletonMJUser) SetCoin(c int64) {
	u.Coin = c
}

func (r *SkeletonMJUser) GetAgent() gate.Agent {
	return r.a
}

//返回一个用户信息
func (u *SkeletonMJUser) GetPlayerInfo(showHand bool, needInpai bool) *mjproto.PlayerInfo {
	info := newProto.NewPlayerInfo()
	*info.NHuPai = u.GetNHuPai()
	*info.BDingQue = u.GetBDingQue()
	*info.BExchanged = u.GetBExchanged()
	*info.BReady = u.getBReady()
	*info.Coin = u.Coin
	*info.IsBanker = u.GetStatus().IsBanker
	info.PlayerCard = u.GetPlayerCard(showHand, needInpai)
	*info.NickName = u.GetNickName()
	*info.UserId = u.GetUserId()
	info.WxInfo = u.GetWxInfo()
	*info.QuePai = u.GameData.HandPai.GetQueFlower()
	*info.Sex = u.Sex
	info.AgentMode = proto.Bool(u.GetStatus().AgentMode)
	return info
}

//是否胡牌
func (u *SkeletonMJUser) GetNHuPai() int32 {
	return 0
}

func (u *SkeletonMJUser) GetBDingQue() int32 {
	if u.GetStatus().DingQue {
		return 1
	} else {
		return 0
	}
}

func (u *SkeletonMJUser) GetBExchanged() int32 {
	if u.GetStatus().Exchange {
		return 1
	} else {
		return 0

	}
}

func (u *SkeletonMJUser) getBReady() int32 {
	if u.GetStatus().IsReady() {
		return 1
	} else {
		return 0
	}
}

func (u *SkeletonMJUser) GetWxInfo() *mjproto.WeixinInfo {
	user := userService.GetUserById(u.GetUserId())
	weixinInfo := newProto.NewWeixinInfo()
	*weixinInfo.NickName = user.GetNickName()
	*weixinInfo.HeadUrl = user.GetHeadUrl()
	*weixinInfo.OpenId = user.GetOpenId()
	return weixinInfo
}
