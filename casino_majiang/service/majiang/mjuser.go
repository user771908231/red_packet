package majiang

import (
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/AgentService"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"casino_server/common/log"
	"casino_server/utils/jobUtils"
	"time"
	"casino_server/service/userService"
	"sync/atomic"
	"errors"
)

var MJUSER_STATUS_INTOROOM int32 = 1; ///刚进入游戏
var MJUSER_STATUS_SEATED int32 = 2; //坐下游戏
var MJUSER_STATUS_READY int32 = 3; ///准备游戏
var MJUSER_STATUS_DINGQUE int32 = 4; ///准备游戏
var MJUSER_STATUS_GAMING int32 = 5; ///正在游戏，这里的正在游戏，表示还没有胡牌..
var MJUSER_STATUS_HUPAI int32 = 6; ///准备游戏

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
	return u.GetStatus() == MJUSER_STATUS_HUPAI
}

func (u *MjUser) IsNotHu() bool {
	return !u.IsHu()
}


//玩家正在游戏中，
func (u *MjUser) IsGaming() bool {
	return u.GetStatus() == MJUSER_STATUS_GAMING
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
	info.WxInfo = u.GetWxInfo()
	return info
}

func (u *MjUser) GetWxInfo() *mjproto.WeixinInfo {
	user := userService.GetUserById(u.GetUserId())
	weixinInfo := newProto.NewWeixinInfo()
	*weixinInfo.NickName = user.GetNickName()
	*weixinInfo.HeadUrl = user.GetHeadUrl()
	*weixinInfo.OpenId = user.GetOpenId()
	return weixinInfo
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

	//类型（1,碰，2,明杠，3,暗杠）
	//得到碰牌
	for i, pai := range u.GameData.HandPai.GetPengPais() {
		if pai != nil && i % 3 == 0 {
			com := newProto.NewComposeCard()
			*com.Value = pai.GetClientId()
			*com.Type = 1        //todo ,需要把type 放置在常量里面 这里代表的是碰牌
			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}


	//得到杠牌
	for i, pai := range u.GameData.HandPai.GetGangPais() {
		if pai != nil && i % 4 == 0 {
			com := newProto.NewComposeCard()
			*com.Value = pai.GetClientId()
			*com.Type = 2       // 这里代表的是杠牌
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
func (u *MjUser) ActHu(p *MJPai, sendUserId uint32, desk *MjDesk) error {

	return nil
}

//用户杠牌,主要是存储数据
func (u *MjUser) Gang(p *MJPai, sendUserId uint32) error {

	return nil
}

//得到判定bean
func (u *MjUser) GetCheckBean(p *MJPai) *CheckBean {
	bean := NewCheckBean()
	*bean.CheckStatus = CHECK_CASE_BEAN_STATUS_CHECKING
	*bean.CanHu = u.GameData.HandPai.GetCanHu()
	*bean.CanPeng = u.GameData.HandPai.GetCanPeng(p)
	*bean.CanGang, _ = u.GameData.HandPai.GetCanGang(p)
	*bean.UserId = u.GetUserId()
	log.T("得到用户[%v]对牌[%v]的check , bean[%v]", u.GetUserId(), p.LogDes(), bean)

	if bean.GetCanGang() || bean.GetCanHu() || bean.GetCanPeng() {
		return bean
	} else {
		return nil
	}
}

//玩家打一张牌
func (u *MjUser) DaPai(p *MJPai) error {
	return u.GameData.HandPai.DelHandlPai(p.GetIndex())
}

//设置用户的状态
func (u *MjUser) SetStatus(s int32) error {
	*u.Status = s
	return nil
}

//判断用户是否可以摸牌
func (u *MjUser) CanMoPai() bool {
	if u.IsNotHu() {
		return true
	} else {
		return false
	}
}
//得到用户的昵称
func (u *MjUser) GetNickName() string {
	return "nickName"
}

func (u *MjUser) AddBillAmount(amount int64) {
	atomic.AddInt64(u.Bill.WinAmount, amount)
}

//删除账单
func (u *MjUser)DelBillBean(pai *MJPai) (error, *BillBean) {
	var bean *BillBean
	index := -1
	for i, info := range u.Bill.Bills {
		if info != nil && info.GetPai().GetIndex() == pai.GetIndex() {
			index = i
			bean = info
			break
		}
	}

	if index > -1 {
		u.Bill.Bills = append(u.Bill.Bills[:index], u.Bill.Bills[index + 1:]...)
		return nil, bean
	} else {
		log.E("服务器错误：删除账单 billBean的时候出错，没有找到对应的杠牌[%v]", pai)
		return errors.New("删除手牌时出错，没有找到对应的手牌..."), nil
	}

}

//增加一条账单
func (u *MjUser) AddBillBean(bean *BillBean) error {
	u.Bill.Bills = append(u.Bill.Bills, bean)
	u.AddBillAmount(bean.GetAmount())
	return nil
}

