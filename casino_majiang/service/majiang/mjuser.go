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
	"fmt"
	"casino_server/utils/numUtils"
)

var MJUSER_STATUS_INTOROOM int32 = 1; ///刚进入游戏
var MJUSER_STATUS_SEATED int32 = 2; //坐下游戏
var MJUSER_STATUS_READY int32 = 3; ///准备游戏
var MJUSER_STATUS_DINGQUE int32 = 4; ///准备游戏
var MJUSER_STATUS_GAMING int32 = 5; ///正在游戏，这里的正在游戏，表示还没有胡牌..
var MJUSER_STATUS_HUPAI int32 = 6; ///准备游戏


var MJUSER_BILL_TYPE_YING_GNAG int32 = 1 //杠赢钱
var MJUSER_BILL_TYPE_SHU_GNAG int32 = 2 //杠输钱

var MJUSER_BILL_TYPE_YING_AN_GNAG int32 = 3 //暗杠赢钱
var MJUSER_BILL_TYPE_SHU_AN_GNAG int32 = 4 //被暗杠输钱

var MJUSER_BILL_TYPE_YING_BA_GANG int32 = 5 //巴杠赢钱
var MJUSER_BILL_TYPE_SHU_BA_GANG int32 = 6 //被巴杠输钱

var MJUSER_BILL_TYPE_YING_HU int32 = 7 //自摸赢钱
var MJUSER_BILL_TYPE_SHU_ZIMO int32 = 8 //自摸输钱

var MJUSER_BILL_TYPE_SHU_DIANPAO int32 = 9 //点炮输钱

var MJUSER_BILL_TYPE_SHU_DAJIAO int32 = 10 //被查叫输钱
var MJUSER_BILL_TYPE_YING_DAJIAO int32 = 11 //查大叫赢钱

var MJUSER_BILL_TYPE_SHU_CHAHUAZHU int32 = 12 //被查花猪输钱
var MJUSER_BILL_TYPE_YING_CHAHUAZHU int32 = 13 //查花猪赢钱


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

func (u *MjUser) IsNotReady() bool {
	return !u.IsReady()
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

//判断用户有没有定缺
func (u *MjUser) IsNotDingQue() bool {
	return !u.IsDingQue()
}

//返回一个用户信息
func ( u *MjUser) GetPlayerInfo(showHand bool, needInpai bool) *mjproto.PlayerInfo {
	info := newProto.NewPlayerInfo()
	*info.NHuPai = u.GetNHuPai()
	*info.BDingQue = u.GetBDingQue()
	*info.BExchanged = u.GetBExchanged()
	*info.BReady = u.getBReady()
	*info.Coin = u.GetCoin()
	*info.IsBanker = u.GetIsBanker()
	info.PlayerCard = u.GetPlayerCard(showHand, needInpai)
	*info.NickName = u.GetNickName()
	*info.UserId = u.GetUserId()
	info.WxInfo = u.GetWxInfo()
	*info.QuePai = u.GameData.HandPai.GetQueFlower()
	*info.Sex = u.GetSex()
	return info
}

func (u *MjUser) GetSex() int32 {
	user := userService.GetUserById(u.GetUserId())
	return user.GetSex()
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
func (u *MjUser) GetPlayerCard(showHand bool, needInpai bool) *mjproto.PlayerCard {
	playerCard := newProto.NewPlayerCard()
	*playerCard.UserId = u.GetUserId()

	//得到inpai
	if needInpai {
		if showHand {
			playerCard.HandCard = append(playerCard.HandCard, u.GameData.HandPai.InPai.GetCardInfo())
		} else {
			playerCard.HandCard = append(playerCard.HandCard, u.GameData.HandPai.InPai.GetBackPai())
		}
	}

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
			*com.Type = int32(mjproto.ComposeCardType_C_PENG)        //todo ,需要把type 放置在常量里面 这里代表的是碰牌
			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}

	//得到杠牌
	for _, info := range u.GameData.GangInfo {
		if info != nil {
			com := newProto.NewComposeCard()
			*com.Value = info.GetPai().GetClientId()

			if info.GetGangType() == GANG_TYPE_DIAN {
				*com.Type = int32(mjproto.ComposeCardType_C_MINGGANG)   // 明杠
			} else if info.GetGangType() == GANG_TYPE_AN {
				*com.Type = int32(mjproto.ComposeCardType_C_ANGANG)       // 暗杠
			} else if info.GetGangType() == GANG_TYPE_BA {
				*com.Type = int32(mjproto.ComposeCardType_C_BAGANG)      // 巴杠
			}

			playerCard.ComposeCard = append(playerCard.ComposeCard, com)
		}
	}

	//得到胡牌
	for _, pai := range u.GameData.HandPai.GetHuPais() {
		if pai != nil {
			//*playerCard.HuCard = pai.GetClientId()
			playerCard.HuCard = append(playerCard.HuCard, pai.GetClientId())
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
	*u.WaitTime = time.Now().Add(time.Second * 30).Unix()
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
func (u *MjUser) GetCheckBean(p *MJPai, xueliuchenghe bool) *CheckBean {
	bean := NewCheckBean()

	*bean.CheckStatus = CHECK_CASE_BEAN_STATUS_CHECKING
	*bean.UserId = u.GetUserId()

	//是否可以胡牌
	if u.IsCanInitCheckCaseHu(xueliuchenghe) {
		*bean.CanHu, _ = u.GameData.HandPai.GetCanHu()
	}
	//是否可以gang
	if u.IsCanInitCheckCaseGang(xueliuchenghe) {
		*bean.CanGang, _ = u.GameData.HandPai.GetCanGang(p)
	}
	//是否可以碰
	if u.IsCanInitCheckCasePeng() {
		*bean.CanPeng = u.GameData.HandPai.GetCanPeng(p)
	}

	log.T("得到用户[%v]对牌[%v]的check , bean[%v]", u.GetUserId(), p.LogDes(), bean)
	//判断过胡.如果有过胡，那么就不能再胡了
	if u.HadGuoHuInfo(p) {
		*bean.CanHu = false
	}

	if bean.GetCanGang() || bean.GetCanHu() || bean.GetCanPeng() {
		return bean
	} else {
		return nil
	}
}

//判断用户是否可以杠
func (u *MjUser) IsCanInitCheckCaseGang(xueliuchenghe bool) bool {
	//这里需要判断是否是 血流成河，目前暂时不判断...

	//1,普通规则
	if u.IsNotHu() {
		return true
	}

	//2,血流成河
	if u.IsHu() && xueliuchenghe {
		return true
	}

	//其他情况返回false
	return false
}

func (u *MjUser) IsCanInitCheckCasePeng() bool {
	//1,普通规则
	if u.IsNotHu() {
		return true
	} else {
		return false;
	}
}

//判断用户是否可以杠
func (u *MjUser) IsCanInitCheckCaseHu(xueliuchenghe bool) bool {
	return u.IsCanInitCheckCaseGang(xueliuchenghe)
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
func (u *MjUser) CanMoPai(xueliuchenghe bool) bool {
	if u.IsNotHu() {
		return true
	}

	if u.IsHu() && xueliuchenghe {
		return true
	}
	return false
}


//判断用户是否可以开始游戏...
func (u *MjUser) CanBegin() bool {
	if u.IsReady() {
		return true
	} else {
		return false
	}
}


//开始游戏的时候 初始化user的信息..
func (u *MjUser) BeginInit(round int32, banker uint32) error {
	//1,游戏开始时候的初始化...
	u.GameData = NewPlayerGameData()        //初始化一个空的麻将牌
	*u.DingQue = false
	*u.Exchanged = false
	if u.GetUserId() == banker {
		*u.IsBanker = true
	} else {
		*u.IsBanker = false
	}

	//杠牌信息
	u.PreMoGangInfo = nil
	//初始化账单
	u.Bill = NewBill()
	//2,初始化统计bean
	statisticsRoundBean := NewStatiscRound()
	*statisticsRoundBean.Round = round
	u.Statisc.RoundBean = append(u.Statisc.RoundBean, statisticsRoundBean)

	return nil
}

//lottery之后，设置user为没有准备
func (u *MjUser) AfterLottery() error {
	*u.Ready = false        //设置为没有准备的状态...
	return nil
}
//得到用户的昵称
func (u *MjUser) GetNickName() string {
	user := userService.GetUserById(u.GetUserId())
	if user == nil {
		return "玩家不存在"
	} else {
		return user.GetNickName()
	}
}

func (u *MjUser) AddBillAmount(amount int64) {
	atomic.AddInt64(u.Bill.WinAmount, amount)
}

func (u *MjUser) SubBillAmount(amount int64) {
	atomic.AddInt64(u.Bill.WinAmount, -amount)
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
		u.SubBillAmount(bean.GetAmount())        //减去
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

//
func (u *MjUser) AddStatisticsWinCoin(coin int64) {
	atomic.AddInt64(u.Statisc.WinCoin, coin)
}

func (u *MjUser) AddCoin(coin int64) {
	atomic.AddInt64(u.Coin, coin)        //更新账户余额
}

func (u *MjUser) ADDCountBaGang() {
	atomic.AddInt32(u.Statisc.CountMingGang, 1)
}

func (u *MjUser) SubCountBaGang() {
	atomic.AddInt32(u.Statisc.CountMingGang, -1)
}

//当删除桌子的时候，需要清除掉用户的会话信息,设置deskId = 0 ，roomId = 0
func (u *MjUser)ClearAgentGameData() {
	log.T("清楚用户[%v]的session信息为默认状态....", u.GetUserId())
	UpdateSession(u.GetUserId(), MJUSER_SESSION_GAMESTATUS_NOGAME, 0, 0, "")
}

func (u *MjUser) BillToString() string {
	result := ""
	if u.Bill != nil {
		s1 := "玩家[%v]总共输赢[%v],下边是细节:\n"
		s1 = fmt.Sprintf(s1, u.GetUserId(), u.Bill.GetWinAmount())
		result += s1
		for _, bb := range u.Bill.GetBills() {
			if bb != nil {
				s2 := "玩家[%v] 关联玩家[%v],牌[%v],分数[%v],类型[%v],描述[%v]\n"
				s2 = fmt.Sprintf(s2, bb.GetUserId(), bb.GetOutUserId(), bb.GetPai().LogDes(), bb.GetAmount(), bb.GetType(), bb.GetDes())
				result += s2
			}
		}

	}
	return result
}

//增加用户被巴杠的统计记录
func (u *MjUser) AddStatisticsCountBeiBaGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountBeiBaGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiBaGang, 1) //汇总
	return nil
}

//增加用户巴杠的统计记录
func (u *MjUser) AddStatisticsCountBaGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountBaGnag, 1) //单局

	atomic.AddInt32(u.Statisc.CountBaGang, 1) //汇总
	return nil
}

//增加用户被暗杠的统计记录
func (u *MjUser) AddStatisticsCountBeiAnGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountBeiAnGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiAnGang, 1) //汇总
	return nil
}

//增加用户暗杠的统计记录
func (u *MjUser) AddStatisticsCountAnGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountAnGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountAnGang, 1) //汇总
	return nil
}

//增加用户明杠的统计记录
func (u *MjUser) AddStatisticsCountMingGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountMingGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountMingGang, 1) //汇总
	return nil
}

//增加用户点杠的统计记录
func (u *MjUser) AddStatisticsCountDianGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountDianGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountDianGang, 1) //汇总
	return nil
}

//增加用户胡牌的统计记录
func (u *MjUser) AddStatisticsCountHu(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountHu, 1) //单局

	atomic.AddInt32(u.Statisc.CountHu, 1) //汇总
	return nil
}

//增加用户点炮的统计记录
func (u *MjUser) AddStatisticsCountDianPao(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountDianPao, 1) //单局

	atomic.AddInt32(u.Statisc.CountDianPao, 1) //汇总
	return nil
}

//增加用户被自摸的统计记录
func (u *MjUser) AddStatisticsCountBeiZiMo(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountBeiZiMo, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiZiMo, 1) //汇总
	return nil
}

//增加用户自摸的统计记录
func (u *MjUser) AddStatisticsCountZiMo(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountZiMo, 1) //单局

	atomic.AddInt32(u.Statisc.CountZiMo, 1) //汇总
	return nil
}

//增加用户查大叫的统计记录
func (u *MjUser) AddStatisticsCountChaDaJiao(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountChaDaJiao, 1) //单局

	atomic.AddInt32(u.Statisc.CountChaDaJiao, 1) //汇总
	return nil
}

//增加用户被查叫的统计记录
func (u *MjUser) AddStatisticsCountBeiChaJiao(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountBeiChaJiao, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiChaJiao, 1) //汇总
	return nil
}

//增加用户查花猪的统计记录
func (u *MjUser) AddStatisticsCountChaHuaZhu(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountChaDaJiao, 1) //单局

	atomic.AddInt32(u.Statisc.CountChaDaJiao, 1) //汇总
	return nil
}

//增加用户被查花猪的统计记录
func (u *MjUser) AddStatisticsCountBeiChaHuaZhu(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return errors.New("没有找到统计的roundBean，无法统计")
	}
	atomic.AddInt32(roundBean.CountBeiChaHuaZhu, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiChaHuaZhu, 1) //汇总
	return nil
}


//统计杠的数量
//func (u *MjUser) StatisticsGangCount(round int32, gangType int32) error {
//
//	//统计每一轮和总的次数
//	bean := u.GetStatisticsRoundBean(round)
//	if bean == nil {
//		log.E("统计的时候出错...")
//		return errors.New("没有找到统计的roundBean，无法统计")
//	}
//
//	if gangType == GANG_TYPE_DIAN {
//		atomic.AddInt32(bean.CountMingGang, 1)
//		atomic.AddInt32(u.Statisc.CountMingGang, 1)
//
//	} else if gangType == GANG_TYPE_BA {
//		atomic.AddInt32(bean.CountBaGnag, 1)
//		atomic.AddInt32(u.Statisc.CountMingGang, 1)
//
//		//总的统计 + 1
//		u.ADDCountBaGang()
//
//	} else if gangType == GANG_TYPE_AN {
//		atomic.AddInt32(bean.CountAnGang, 1)
//		atomic.AddInt32(u.Statisc.CountAnGang, 1)
//	}
//
//	return nil
//}

//计算胡牌的次数,初步是胡了几次可以统计，详细的各种类型 还没有统计...
//func (u *MjUser) StatisticsHuCount(round int32, huUserId uint32, huType int32) error {
//
//	//计算总的
//	atomic.AddInt32(u.Statisc.CountHu, 1)
//
//	//计算每一句的胡牌
//	bean := u.GetStatisticsRoundBean(round)
//	if bean == nil {
//		log.E("统计的时候出错...")
//		return errors.New("没有找到统计的roundBean，无法统计")
//	}
//	atomic.AddInt32(bean.CountHu, 1)
//	return nil
//}


//统计点炮的次数
func (u *MjUser) StatisticsDianCount(dianUserId uint32, dianType int32) {
	atomic.AddInt32(u.Statisc.CountDianPao, 1)
}


//得到每一局的统计bean...
func (u *MjUser) GetStatisticsRoundBean(round int32) *StatiscRound {
	for _, bean := range u.Statisc.RoundBean {
		if bean != nil && bean.GetRound() == round {
			return bean
		}
	}
	//如果没有找到返回nil
	return nil
}

//增加一个过胡的info
/**
	用户操作1，摸牌，2，碰牌，3，杠牌 之后，需要删除过胡的信息
 */


func (u *MjUser) AddGuoHuInfo(checkCase *CheckCase) {
	if checkCase == nil {
		return
	}

	checkBean := checkCase.GetBeanByUserIdAndStatus(u.GetUserId(), CHECK_CASE_BEAN_STATUS_CHECKING)
	if checkBean != nil && checkBean.GetCanHu() {
		guoHuInfo := NewGuoHuInfo()
		*guoHuInfo.SendUserId = checkCase.GetUserIdOut()
		guoHuInfo.Pai = checkCase.CheckMJPai
		*guoHuInfo.FanShu = 0        //现在都设置为0翻
		u.GameData.GuoHuInfo = append(u.GameData.GuoHuInfo, guoHuInfo)
	}

}

//增加胡牌的信息
func (u *MjUser) AddHuPaiInfo(hu *HuPaiInfo) {
	u.GameData.HuInfo = append(u.GameData.HuInfo, hu)
	u.GameData.HandPai.HuPais = append(u.GameData.HandPai.HuPais, hu.Pai)        //增加胡牌
	u.SetStatus(MJUSER_STATUS_HUPAI)
}

//删除过胡的信息
func (u *MjUser)DelGuoHuInfo() error {
	u.GameData.GuoHuInfo = nil
	return nil
}

//是否已经有过胡了
func (u *MjUser) HadGuoHuInfo(pai *MJPai) bool {
	if u.GameData.GuoHuInfo == nil || len(u.GameData.GuoHuInfo) <= 0 {
		return false
	}

	//目前只做成  牌一样的时候再判断

	//如果huinfo的牌和pai 一样，表示有guohu的info
	for _, info := range u.GameData.GuoHuInfo {
		if pai.GetClientId() == info.GetPai().GetClientId() {
			return true
		}
	}

	//没有过胡的信息
	return false
}

//判断是否是花猪
func (u *MjUser) IsHuaZhu() bool {
	for _, pai := range u.GameData.HandPai.Pais {
		if pai != nil && pai.GetFlower() == u.GameData.HandPai.GetQueFlower() {
			//是花猪
			return true
		}
	}
	//不是花猪
	return false
}

//判断不是花猪
func (u *MjUser) IsNotHuaZhu() bool {
	return !u.IsHuaZhu()
}


//判断用户是不是没有叫 有叫为true
func (u *MjUser) IsYouJiao() bool {
	//
	//u.GameData.HandPai.GetCanHu()
	for i := 0; i < 108; i++ {
		pai := InitMjPaiByIndex(i)
		if pai.GetFlower() == u.GameData.HandPai.GetQueFlower() {
			continue
		}

		u.GameData.HandPai.InPai = pai
		canhu, _ := u.GameData.HandPai.GetCanHu()
		if canhu {
			//
			log.T("玩家查叫的时候，查到可以胡牌[%v]", pai.LogDes())
			return true
		}

	}
	return false
}

//增加账单
func (u *MjUser) AddBill(relationUserid uint32, billType int32, des string, score int64, pai *MJPai) error {
	//用户赢钱的账户,赢钱的账单
	bill := NewBillBean()
	*bill.UserId = u.GetUserId()
	*bill.OutUserId = relationUserid
	//*bill.Type = MJUSER_BILL_TYPE_YING_HU
	*bill.Type = billType
	*bill.Des = des
	*bill.Amount = score        //杠牌的收入金额
	bill.Pai = pai
	u.AddBillBean(bill)

	//计算账单的地方 来加减用户的coin
	u.AddCoin(score)                //统计用户剩余多少钱
	u.AddStatisticsWinCoin(score)        //统计用户输赢多少钱
	return nil
}

//通过手牌，得到杠牌的信息
func (u *MjUser) GetJiaoPaisByHandPais() []*MJPai {
	return GetJiaoPais(u.GameData.HandPai.Pais)
}


//比较杠牌之后的叫牌和杠牌之前的叫牌的信息是否一样
func (u *MjUser) AfterGangEqualJiaoPai(beforJiaoPais []*MJPai, gangPai *MJPai) bool {

	//1，获得杠牌之后的手牌
	var afterPais []*MJPai
	for _, p := range u.GameData.HandPai.Pais {
		if p.GetClientId() != gangPai.GetClientId() {
			afterPais = append(afterPais, p)
		}
	}

	//2，通过杠牌之后的手牌 获得此时的叫牌
	afterJiaoPais := GetJiaoPais(afterPais)

	//2,比较beforJiaoPais 和 afterJiaoPais
	if len(afterPais) != len(beforJiaoPais) {
		return false
	}

	for _, aj := range afterJiaoPais {

		forbool := false
		for _, bj := range beforJiaoPais {
			if aj.GetClientId() == bj.GetClientId() {
				forbool = true
				break
			}
		}

		if !forbool {
			return false
		}
	}

	return true;
}

func (u *MjUser) IsQuePai(mjPai *MJPai) bool {
	if u.GetGameData().GetHandPai().GetQueFlower() == mjPai.GetFlower() {
		return true
	}
	return false
}

func (u *MjUser) GetUserPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	} else {
		return u.GameData.HandPai.GetDes()
	}
}

func (u *MjUser) GetUserPengPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range u.GameData.HandPai.PengPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}

func (u *MjUser) GetUserGnagPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range u.GameData.HandPai.GangPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}

func (u *MjUser) GetUserHuPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range u.GameData.HandPai.HuPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}
func (u *MjUser) GetUserInPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	return u.GameData.HandPai.InPai.LogDes()

}

func (u *MjUser) GetTransferredStatus() string {
	ret := ""
	switch u.GetStatus() {
	case MJUSER_STATUS_DINGQUE:
		ret = "已定缺"
	case MJUSER_STATUS_GAMING:
		ret = "游戏中"
	case MJUSER_STATUS_HUPAI:
		ret = "已胡牌"
	case MJUSER_STATUS_INTOROOM:
		ret = "已进入房间"
	case MJUSER_STATUS_READY:
		ret = "已准备"
	case MJUSER_STATUS_SEATED:
		ret = ""
	default:
	}
	return ret
}