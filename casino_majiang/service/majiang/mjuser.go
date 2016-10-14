package majiang

import (
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/AgentService"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"casino_server/common/log"
	"casino_server/utils/jobUtils"
	"time"
)

var MJUSER_STATUS_INTOROOM int32 = 1; ///刚进入游戏
var MJUSER_STATUS_SEATED int32 = 2; //坐下游戏
var MJUSER_STATUS_READY int32 = 3; ///准备游戏
var MJUSER_STATUS_DINGQUE int32 = 4; ///准备游戏

//麻将玩家

//发送接口
func (u *MjUser)WriteMsg(p proto.Message) error {
	agent := AgentService.GetAgent(u.GetUserId())
	if agent != nil {
		agent.WriteMsg(p)
	} else {
		log.T("给用户[%v]发送proto[%v]失败，因为没有找到用户的agent。", u.GetUserId(), p)
	}
	return nil
}

//是否是准备中...
func (u *MjUser) IsReady() bool {
	return u.GetStatus() == MJUSER_STATUS_READY
}

//用户是否胡牌
func (u *MjUser) IsHu() bool {
	return true
}

func (u *MjUser) IsNotHu() bool {
	return !u.IsHu()
}


//todo 玩家是否在游戏状态中
func (u *MjUser) IsGaming() bool {
	return true

}

//判断用户是否已经定缺
func (u *MjUser) IsDingQue() bool {
	return u.GetDingQue()
}

//返回一个用户信息
func ( u *MjUser) GetPlayerInfo(showHand bool) *mjproto.PlayerInfo {
	info := newProto.NewPlayerInfo()
	*info.NHuPai = u.GetNHuPai()
	*info.BDingQue = u.GetBDingQue()
	*info.BExchanged = u.GetBExchanged()
	*info.BReady = u.getBReady()
	*info.Coin = u.GetCoin()
	*info.IsBanker = u.GetIsBanker()
	info.PlayerCard = u.GetPlayerCard(showHand)
	*info.NickName = "测试nickName"
	*info.UserId = u.GetUserId()
	return info
}

//得到手牌
//showHand 是否显示手牌
func (u *MjUser) GetPlayerCard(showHand bool) *mjproto.PlayerCard {
	playerCard := newProto.NewPlayerCard()
	*playerCard.UserId = u.GetUserId()

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


	//得到碰牌
	for i, pai := range u.GameData.HandPai.GetPengPais() {
		if pai != nil && i % 3 == 0 {
			com := newProto.NewComposeCard()
			*com.Value = pai.GetClientId()
			//com.Type =	这里代表的是碰牌
			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}


	//得到杠牌
	for i, pai := range u.GameData.HandPai.GetGangPais() {
		if pai != nil && i % 4 == 0 {
			com := newProto.NewComposeCard()
			*com.Value = pai.GetClientId()
			//com.Type =	这里代表的是杠牌
			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}

	//得到胡牌
	for _, pai := range u.GameData.HandPai.GetHuPais() {
		if pai != nil {
			*playerCard.HuCard = pai.GetClientId()
		}
	}


	//打出去的牌
	for _, pai := range u.GameData.HandPai.GetOutPais() {
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

//发送overTrun
func (u *MjUser) SendOverTurn(p *mjproto.Game_OverTurn) error {

	u.Wait()
	u.WriteMsg(p)

	return nil
}


//等待超时
func (u *MjUser) Wait() error {

	//这里需要根据等待的类型来做不同的处理...
	/**
		如果是摸牌之后等带,需要自动打一张牌
		如果是判定，系统自动过
	 */
	jobUtils.DoAsynJob(time.Second * 1, func() bool {
		//如果超时，那么系统自动打一张牌...
		return true
	})
	return nil

}


//用户胡牌
func (u *MjUser) ActHu(p *MJPai, sendUserId uint32) error {
	//判断能不能胡

	//得到胡牌的信息

	//胡牌之后的信息
	hu := NewHuPaiInfo()
	*hu.SendUserId = sendUserId
	hu.Pai = p
	//hu.Fan
	u.GameData.HuInfo = append(u.GameData.HuInfo, hu)

	//增加胡牌
	u.GameData.HandPai.HuPais = append(u.GameData.HandPai.HuPais, p)

	//发送胡牌成功的回复
	ack := newProto.NewGame_AckActHu()

	log.T("给用户[%v]发送胡牌的ack[%v]", u.GetUserId(), ack)
	u.WriteMsg(ack)
	return nil
}

//用户杠牌,主要是存储数据
func (u *MjUser) Gang(p *MJPai, sendUserId uint32) error {

	//杠牌的类型
	var gangType int32 = 0
	var gangKey []int32
	//增加杠牌
	u.GameData.HandPai.GangPais = append(u.GameData.HandPai.GangPais, p)
	for _, pai := range u.GameData.HandPai.Pais {
		if pai.GetFlower() == p.GetFlower() && pai.GetValue() == p.GetValue() {
			//增加杠牌
			u.GameData.HandPai.GangPais = append(u.GameData.HandPai.GangPais, pai)
			gangKey = append(gangKey, pai.GetIndex())
		}
	}

	//增加杠牌info
	info := NewGangPaiInfo()
	*info.SendUserId = sendUserId
	*info.GangType = gangType
	info.Pai = p
	u.GameData.GangInfo = append(u.GameData.GangInfo, info)

	//增加杠牌状态
	u.PreMoGangInfo = info

	//减少手中的杠牌
	for _, key := range gangKey {
		u.GameData.HandPai.DelPai(key)
	}

	return nil
}

//得到判定bean
func (u *MjUser) GetCheckBean(p *MJPai) *CheckBean {
	return nil
}

//玩家打一张牌
func (u *MjUser) DaPai(p *MJPai) error {
	u.GameData.HandPai.DelPai(p.GetIndex())
	return nil
}

//设置用户的状态
func (u *MjUser) SetStatus(s int32) error {
	*u.Status = s
	return nil
}