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
		log.Fatal("给用户[%v]发送proto[%v]失败，因为没有找到用户的agent。", u.UserId, p)
	}
	return nil
}

//是否是准备中...
func (u *MjUser) IsReady() bool {
	return u.GetStatus() == MJUSER_STATUS_SEATED
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
	//info.SeatId = u

	return info
}

//得到手牌
func (u *MjUser) GetPlayerCard() *mjproto.PlayerCard {
	playerCard := newProto.NewPlayerCard()

	for _, pai := range u.MJHandPai.GetPais() {
		if pai != nil {
			playerCard.HandCard = append(playerCard.HandCard, pai.GetCardInfo())
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
