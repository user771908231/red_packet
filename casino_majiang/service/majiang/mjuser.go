package majiang

import (
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/AgentService"
	"github.com/name5566/leaf/log"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
)

var MJUSER_STATUS_INTOROOM int32 = 1; ///刚进入游戏
var MJUSER_STATUS_SEATED int32 = 2; ///刚进入游戏
var MJUSER_STATUS_READY int32 = 3; ///刚进入游戏


//麻将玩家

//发送接口
func (u *MjUser)WriteMsg(p proto.Message) error {
	agent := AgentService.GetAgent(u.GetUserId())
	if agent != nil {
		agent.WriteMsg(p)
	} else {
		log.Error("给用户[%v]发送proto[%v]失败，因为没有找到用户的agent。", u.GetUserId(), p)
	}
	return nil
}

//是否是准备中...
func (u *MjUser) IsReady() bool {
	return u.GetStatus() == MJUSER_STATUS_READY
}

//玩家是否在游戏状态中
func (u *MjUser) IsGaming() bool {
	return true

}

//返回一个用户信息
func ( u *MjUser) GetPlayerInfo() *mjproto.PlayerInfo {
	info := newProto.NewPlayerInfo()
	*info.NHuPai = u.GetNHuPai()
	*info.BDingQue = u.GetBDingQue()
	*info.BExchanged = u.GetBExchanged()
	*info.BReady = u.getBReady()
	*info.Coin = u.GetCoin()
	*info.IsBanker = u.GetIsBanker()
	info.PlayerCard = u.GetPlayerCard()
	*info.NickName = "测试nickName"
	//info.SeatId = u
	return info
}

//得到手牌
func (u *MjUser) GetPlayerCard() *mjproto.PlayerCard {
	playerCard := newProto.NewPlayerCard()

	//得到手牌
	for _, pai := range u.MJHandPai.GetPais() {
		if pai != nil {
			playerCard.HandCard = append(playerCard.HandCard, pai.GetCardInfo())
		}
	}


	//得到碰牌
	for i, pai := range u.MJHandPai.GetPengPais() {
		if pai != nil && i % 3 == 0 {
			com := newProto.NewComposeCard()
			*com.Value = pai.GetClientId()
			//com.Type =	这里代表的是碰牌
			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}


	//得到杠牌
	for i, pai := range u.MJHandPai.GetGangPais() {
		if pai != nil && i % 4 == 0 {
			com := newProto.NewComposeCard()
			*com.Value = pai.GetClientId()
			//com.Type =	这里代表的是杠牌
			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}

	//得到胡牌
	for _, pai := range u.MJHandPai.GetPais() {
		if pai != nil {
			*playerCard.HuCard = pai.GetClientId()
		}
	}


	//打出去的牌
	for _, pai := range u.MJHandPai.GetPais() {
		if pai != nil {
			playerCard.OutCard = append(playerCard.OutCard, pai.GetClientId())
		}
	}

	return playerCard
}


//是否胡牌
func (u *MjUser) GetNHuPai() int32 {
	return 0
}

func (u *MjUser) GetBDingQue() int32 {
	if u.GetDingQue() {
		return 1
	} else {
		return 0
	}
}

func (u *MjUser) GetBExchanged() int32 {
	if u.GetExchanged() {
		return 1
	} else {
		return 0

	}
}

func (u *MjUser) getBReady() int32 {
	if u.GetReady() {
		return 1
	} else {
		return 0

	}
}

//发牌
func (u *MjUser) GetDealCards() *mjproto.Game_DealCards {
	dealCards := newProto.NewGame_DealCards()
	*dealCards.Header.UserId = u.GetUserId()
	dealCards.PlayerCard = u.GetPlayerCard()
	return dealCards
}

//发送overTrun
func (u *MjUser) SendOverTurn(p *mjproto.Game_OverTurn) error {

	go u.Wait()

	u.WriteMsg(p)

	return nil
}


//等待超时
func (u *MjUser) Wait() error {
	return nil

}

//用户胡牌
func (u *MjUser) ActHu() error {

	//判断自摸

	//判断点炮

	return nil
}
