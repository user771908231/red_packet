package majiang

import (
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_common/common/game"
	"casino_common/common/log"
	"casino_common/common/service/countService"
	"casino_common/common/sessionService"
	"casino_common/common/userService"
	"casino_common/proto/ddproto"
	"casino_mj_changsha/msg/funcsInit"
	mjproto        "casino_mj_changsha/msg/protogo"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/gate"
	ltimer "github.com/name5566/leaf/timer"
	"github.com/name5566/leaf/util"
	"reflect"
	"sort"
	"sync/atomic"
	"time"
)

//玩家的游戏状态
var MJUSER_STATUS_INTOROOM int32 = 1 ///刚进入游戏
var MJUSER_STATUS_SEATED int32 = 2   //坐下游戏
var MJUSER_STATUS_READY int32 = 3    //准备游戏
var MJUSER_STATUS_DINGQUE int32 = 4  //准备定缺
var MJUSER_STATUS_GAMING int32 = 5   //正在游戏，这里的正在游戏，表示还没有胡牌..
var MJUSER_STATUS_HUPAI int32 = 6    //已经胡牌了

//账单的类型
var MJUSER_BILL_TYPE_YING_GNAG int32 = 1         //杠赢钱
var MJUSER_BILL_TYPE_SHU_GNAG int32 = 2          //杠输钱
var MJUSER_BILL_TYPE_YING_AN_GNAG int32 = 3      //暗杠赢钱
var MJUSER_BILL_TYPE_SHU_AN_GNAG int32 = 4       //被暗杠输钱
var MJUSER_BILL_TYPE_YING_BA_GANG int32 = 5      //巴杠赢钱
var MJUSER_BILL_TYPE_SHU_BA_GANG int32 = 6       //被巴杠输钱
var MJUSER_BILL_TYPE_YING_HU int32 = 7           //自摸赢钱
var MJUSER_BILL_TYPE_SHU_ZIMO int32 = 8          //自摸输钱
var MJUSER_BILL_TYPE_SHU_DIANPAO int32 = 9       //点炮输钱
var MJUSER_BILL_TYPE_SHU_DAJIAO int32 = 10       //被查叫输钱
var MJUSER_BILL_TYPE_YING_DAJIAO int32 = 11      //查大叫赢钱
var MJUSER_BILL_TYPE_SHU_CHAHUAZHU int32 = 12    //被查花猪输钱
var MJUSER_BILL_TYPE_YING_CHAHUAZHU int32 = 13   //查花猪赢钱
var MJUSER_BILL_TYPE_SHU_TUIGANGQIAN int32 = 14  //没叫退杠钱
var MJUSER_BILL_TYPE_YING_TUIGANGQIAN int32 = 15 //没叫退杠钱

//玩家解散房间的游戏状态
var MJUER_APPLYDISSOLVE_S_REFUSE int32 = -1 //拒绝解散
var MJUER_APPLYDISSOLVE_S_DEFAULT int32 = 0 //没有处理
var MJUER_APPLYDISSOLVE_S_AGREE int32 = 1   //同意解散

var MJUSER_NEEDHAIDI_STATUS_REFUSE int32 = -1 //不需要
var MJUSER_NEEDHAIDI_STATUS_DEFAULT int32 = 0 //默认没有判断
var MJUSER_NEEDHAIDI_STATUS_NEED int32 = 1    //需要

//
var ACTTYPE_PENG int32 = 1
var ACTTYPE_BU int32 = 2
var ACTTYPE_GANG int32 = 3
var ACTTYPE_CHI int32 = 4
var ACTTYPE_GUO int32 = 5
var ACTTYPE_HU int32 = 6
var ACTTYPE_OUT int32 = 7 //打牌

//麻将玩家

type ActCheck struct {
	canPeng bool
}

type MjUser struct {
	*PMjUser
	*game.GameUser
	dissolveTimer *ltimer.Timer            //申请解散房间的timer
	readyTimer    *ltimer.Timer            //到了时间没有准备也需要做处理
	dingQueTimer  *ltimer.Timer            //到了时间没有定缺 需要做处理
	d             *MjDesk                  //桌子
	Log           *countService.T_game_log //任务统计信息
	*MjuserChangShaConfig                  //长沙麻将的配置
	ip            string                   //玩家的Ip
	headUrl       string                   //头像headurl
	openId        string                   //opendId
	sex           int32                    //性别
	ActCheck      *ActCheck
}

func (u *MjUser) GetIp() string {
	return u.ip
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

//是否换三张
func (u *MjUser) IsExchange() bool {
	return u.GetExchanged()
}

//是否已经换三张
func (u *MjUser) IsNotExchange() bool {
	return !u.IsExchange()
}

//返回一个用户信息
func (u *MjUser) GetPlayerInfo(showHand bool) *mjproto.PlayerInfo {
	info := newProto.NewPlayerInfo()
	*info.NHuPai = u.GetNHuPai()
	*info.BDingQue = u.GetBDingQue()
	*info.BExchanged = u.GetBExchanged()
	*info.BReady = u.getBReady()
	*info.Coin = u.GetCoin()
	*info.IsBanker = u.GetIsBanker()
	info.PlayerCard = u.GetPlayerCard(showHand)
	*info.NickName = u.GetNickName()
	*info.UserId = u.GetUserId()
	*info.QuePai = u.GameData.HandPai.GetQueFlower()
	*info.Sex = u.sex
	info.AgentMode = proto.Bool(u.GetAgentMode())
	info.Ip = proto.String(u.GetIp())
	info.WxInfo = &mjproto.WeixinInfo{
		NickName: proto.String(u.GetNickName()),
		HeadUrl:  proto.String(u.headUrl),
		OpenId:   proto.String(u.openId),
	}
	return info
}

//得到手牌
//showHand 是否显示手牌
func (u *MjUser) GetPlayerCard(showHand bool) *mjproto.PlayerCard {
	//log.T("得到玩家[%v]的手牌,showHand[%v],needInpai[%v]", u.GetUserId(), showHand, needInpai)
	playerCard := newProto.NewPlayerCard()
	*playerCard.UserId = u.GetUserId()
	playerCard.HandCardCount = proto.Int32(int32(len(u.GetGameData().GetHandPai().GetPais()))) //手牌的长度

	//得到inpai
	/*
		1，in1 !=nil
		2，in2 ==nil
		这样能保证在杠牌之后断线 不发送in牌
	*/
	if u.GetGameData().GetHandPai().GetInPai() != nil && u.GetGameData().GetHandPai().GetInPai2() == nil {
		*playerCard.HandCardCount++
		if showHand {
			playerCard.HandCard = append(playerCard.HandCard, u.GameData.HandPai.InPai.GetCardInfo())
		}
	}

	//得到所有的手牌手牌
	if showHand {
		for _, pai := range u.GameData.HandPai.GetPais() {
			if pai != nil {
				playerCard.HandCard = append(playerCard.HandCard, pai.GetCardInfo())
			}
		}
	}

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

			if info.GetGangType() == GANG_TYPE_DIAN {
				*com.Type = int32(mjproto.ComposeCardType_C_MINGGANG) // 明杠
			} else if info.GetGangType() == GANG_TYPE_AN {
				*com.Type = int32(mjproto.ComposeCardType_C_ANGANG) // 暗杠
			} else if info.GetGangType() == GANG_TYPE_BA {
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

//玩家发送消息
func (u *MjUser) WriteMsg(msgold proto.Message) {
	msg := util.DeepClone(msgold).(proto.Message) //防治出错
	//log.T("开始准备给用户发送信息:leave %v ,isBreadk %v", u.GetIsLeave(), u.GetIsBreak())
	typeString := reflect.TypeOf(msg).String()
	if u.Agent == nil || u.GetIsLeave() || u.GetIsBreak() {
		log.T("开始给真实user[%v],agent ==nil ? %v ,leave :%v break: %v 发送-----失败----[%v]: %v", u.GetUserId(), u.Agent == nil, u.GetIsLeave(), u.GetIsBreak(), typeString, msg)
		return
	}

	log.T("%v开始给真实user[%v]发送[%v]: %v", u.d.DlogDes(), u.GetUserId(), typeString, msg)
	u.Agent.WriteMsg(msg)
}

//发送overTrun
/**
这里需要区分有托管 和没有托管的状态：
1，有托管的时候，给玩家发送
*/
func (u *MjUser) SendOverTurn(p proto.Message) error {
	if p == nil {
		return Error.NewError(consts.ACK_RESULT_ERROR, "overturn == nil")
	}

	//长沙麻将
	//长沙麻将的特殊处理...
	/*
		如果：1，长沙麻将，2,有杠，3,overTurn,4,普通overTurn
			就自动打牌...
	*/
	u.WriteMsg(p)
	if u.d.IsChangShaMaJiang() && u.GetChangShaGangStatus() {
		t := reflect.TypeOf(p)
		if t == reflect.TypeOf(&mjproto.Game_OverTurn{}) {
			overTrun := p.(*mjproto.Game_OverTurn)
			if !overTrun.GetCanHu() {
				go func() {
					time.Sleep(time.Second * 1)
					log.T("[%v][%v]杠之后开摸的牌，然后自动打牌【%v】...", u.d.DlogDes(), u.GetUserId(), u.UserPai2String(), overTrun)
					u.d.ActOut(u.GetUserId(), overTrun.GetActCard().GetId(), true) //长沙杠之后自动打牌...
				}()
			}
		}
	}
	return nil
}

//得到判定bean
func (u *MjUser) GetCheckBean(p *MJPai, xueliuchenghe bool, remainPaiCoun int32, isBaGang bool) *CheckBean {
	bean := NewCheckBean()

	*bean.CheckStatus = CHECK_CASE_BEAN_STATUS_CHECKING
	*bean.UserId = u.GetUserId()
	bean.CheckPai = p
	var fan int32 = 0

	//是否可以胡牌
	if u.IsCanInitCheckCaseHu(xueliuchenghe) {
		*bean.CanHu, fan, _, _, _, _ = u.d.HuParser.GetCanHu(u.GameData.HandPai, p, false, 0, u.d.IsBanker(u))
	}
	//是否可以杠
	if u.IsCanInitCheckCaseGang(xueliuchenghe) { //initCheckCase的时候 得到一个ChechCaseGang
		*bean.CanGang, bean.GangCards = u.GameData.HandPai.GetCanGang(p, remainPaiCoun)
		//log.T("得到可以杠:%v,杠牌:%v", bean.GetCanGang(), bean.GetGangCards())
	}
	//是否可以碰
	if u.IsCanInitCheckCasePeng() {
		*bean.CanPeng = u.GameData.HandPai.GetCanPeng(p)
	}

	//是否可以吃牌
	if u.IsCanInitCheckCaseChi() && !isBaGang { //解决巴杠被吃的问题
		*bean.CanChi, bean.ChiCards = u.GameData.HandPai.GetCanChi(p)
	}

	log.T("%v得到用户[%v]对牌[%v]的check , bean[%v]", u.d.DlogDes(), u.GetUserId(), p.LogDes(), bean)
	//判断过胡.如果有过胡，那么就不能再胡了
	/**
	过胡要分两种：
	成都麻将：只有要过胡，那就不能胡
	长沙麻将：如果是自己的，那不胡，如果是别人点了，翻数<= 过户的时候，不能胡，翻数> 过户的时候可以胡
	*/
	if u.HadGuoHuInfo(fan) {
		*bean.CanHu = false
	}

	if bean.GetCanGang() || bean.GetCanHu() || bean.GetCanPeng() || bean.GetCanChi() {
		return bean
	} else {
		return nil
	}
}

//判断用户是否可以杠
func (u *MjUser) IsCanInitCheckCaseGang(xueliuchenghe bool) bool {

	//长沙麻将的判断规则
	if u.d.IsChangShaMaJiang() {
		//1,如果有杠，直接返回false
		if u.GetChangShaGangStatus() {
			return false
		}
		//2,如果没有胡牌，返回true
		if u.IsNotHu() {
			return true
		}

	} else {
		//成都麻将的判断规则 这个是默认规则
		//1,普通规则
		if u.IsNotHu() {
			return true
		}

		//2,血流成河
		if u.IsHu() && xueliuchenghe {
			return true
		}
	}

	//其他情况返回false
	return false
}

func (u *MjUser) IsCanInitCheckCasePeng() bool {
	if u.d.IsChangShaMaJiang() {
		//长沙麻将的就判断规则
		//1,如果有杠，直接返回false
		if u.GetChangShaGangStatus() {
			return false
		}
		//没有胡牌返回true
		if u.IsNotHu() {
			return true
		}

	} else {
		//1,默认成都的规则普通规则
		if u.IsNotHu() {
			return true
		} else {
			return false
		}
	}
	return false
}

//目前主要是长沙麻将使用
func (u *MjUser) IsCanInitCheckCaseChi() bool {
	//杠牌之后不允许吃
	if u.GetChangShaGangStatus() {
		return false
	}

	//判断顺序是否正确
	if u.d.IsChangShaMaJiang() {
		if (u.d.getIndexByUserId(u.d.GetCheckCase().GetUserIdOut())+1)%int(u.d.GetUserCountLimit()) == u.d.getIndexByUserId(u.GetUserId()) {
			return true
		}
	}
	return false
}

//判断用户是否可以杠
func (u *MjUser) IsCanInitCheckCaseHu(xueliuchenghe bool) bool {
	//长沙麻将的判断规则
	if u.d.IsChangShaMaJiang() {
		//2,如果没有胡牌，返回true
		if u.IsNotHu() {
			return true
		}

	} else {
		//成都麻将的判断规则 这个是默认规则
		//1,普通规则
		if u.IsNotHu() {
			return true
		}

		//2,血流成河
		if u.IsHu() && xueliuchenghe {
			return true
		}
	}
	//其他情况返回false
	return false
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

/**
补牌的逻辑问题
*/
func (u *MjUser) GetCanChangShaGang(pai *MJPai) bool {
	log.T("%v 玩家[%v]开始判断是否能长沙杠pai[%v]", u.d.DlogDes(), u.GetUserId(), pai.LogDes())
	if u.GameData == nil || u.GameData.HandPai == nil || u.GameData.HandPai.Pais == nil || len(u.GameData.HandPai.Pais) <= 0 {
		log.T("%v 玩家[%v]判断是否能杠错误 玩家手牌为空", u.d.DlogDes(), u.GetUserId())
		return false
	}
	//1 手牌中和pai 花色相同的去掉 by 彬哥

	var newHandPai *MJHandPai
	if newHandi := util.DeepClone(u.GetGameData().GetHandPai()); newHandi != nil {
		newHandPai = newHandi.(*MJHandPai)
	}

	//过滤到需要杠的牌
	newHandPai.Pais = make([]*MJPai, 0)
	pais := u.GameData.HandPai.Pais
	for i := 0; i < len(pais); i++ {
		fp := pais[i]
		if fp != nil && fp.GetClientId() != pai.GetClientId() {
			newHandPai.Pais = append(newHandPai.Pais, fp)
		}
	}

	//添加in牌
	inpai := u.GetGameData().GetHandPai().GetInPai()
	if inpai != nil && inpai.GetClientId() != pai.GetClientId() {
		newHandPai.Pais = append(newHandPai.Pais, inpai)
	}

	//2 看剩下的牌有无听 by 彬哥
	jiaoPais := u.d.HuParser.GetJiaoPais(newHandPai) //GetCanChangShaGang
	if jiaoPais != nil && len(jiaoPais) > 0 {
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
	u.GameData = NewPlayerGameData() //初始化一个空的麻将牌
	*u.DingQue = false
	*u.Exchanged = false
	u.changshaGang = false
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
	u.ActTimeoutCount = proto.Int32(0) //初始化超时的次数
	u.NeedHaidiStatus = proto.Int32(MJUSER_NEEDHAIDI_STATUS_DEFAULT)
	u.AgentMode = proto.Bool(false) //第二局开始默认不托管
	u.initTaskLog()

	return nil
}

//初始化玩家任务统计信息
func (u *MjUser) initTaskLog() {
	u.Log.UserId = u.GetUserId()
	u.Log.Bill = 0
	u.Log.GameId = ddproto.CommonEnumGame_GID_MAHJONG
	u.Log.GameNumber = u.GetGameNumber()
	u.Log.RoomType = ddproto.COMMON_ENUM_ROOMTYPE(u.d.GetRoomType())
	u.Log.RoomLevel = u.d.GetRoomLevel()
	u.Log.StartTime = time.Now().Unix()
	u.Log.IsWine = false
}

//更新玩家任务系统用到的统计信息
func (u *MjUser) UpdateTaskLog() {
	u.Log.EndTime = time.Now().Unix()
	u.Log.Bill = float64(u.GetBill().GetWinAmount())
	if u.GetBill().GetWinAmount() > 0 {
		u.Log.IsWine = true
	}

	//修成改异步的插入数据
	go func(l *countService.T_game_log) {
		defer Error.ErrorRecovery("插入统计数据")
		l.Insert()
	}(u.Log)
}

//lottery之后，设置user为没有准备
func (u *MjUser) AfterLottery() error {
	//准备状态
	u.SetStatus(MJUSER_STATUS_SEATED)
	u.Ready = proto.Bool(false)     //设置为没有准备的状态...
	u.DingQue = proto.Bool(false)   //设置为没有定缺的状态...
	u.AgentMode = proto.Bool(false) //第二局开始默认不托管

	log.T("%v玩家[%v]开始开始插入统计数据", u.d.DlogDes(), u.GetUserId())
	u.UpdateTaskLog()
	return nil
}

func (u *MjUser) SubBillAmount(amount int64) {
	atomic.AddInt64(u.Bill.WinAmount, -amount)
}

//删除账单
func (u *MjUser) DelBillBean(pai *MJPai) (error, *BillBean) {
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
		u.Bill.Bills = append(u.Bill.Bills[:index], u.Bill.Bills[index+1:]...)
		u.SubBillAmount(bean.GetAmount()) //减去
		return nil, bean
	} else {
		log.E("服务器错误：删除账单 billBean的时候出错，没有找到对应的杠牌[%v]的账单", pai)
		return errors.New("删除手牌时出错，没有找到对应的手牌..."), nil
	}

}

//增加一条账单
func (u *MjUser) AddBillBean(bean *BillBean) error {
	u.Bill.Bills = append(u.Bill.Bills, bean)
	atomic.AddInt64(u.Bill.WinAmount, bean.GetAmount()) //Bill里面增加账单
	return nil
}

//
func (u *MjUser) AddStatisticsWinCoin(coin int64) {
	atomic.AddInt64(u.Statisc.WinCoin, coin)
}

//更新用户金币
func (u *MjUser) AddCoin(coin int64, roomType int32) {
	//减少玩家金额
	atomic.AddInt64(u.Coin, coin) //更新账户余额
	//如果是金币场。需要更新用户的金币余额
	if roomType == ROOMTYPE_COINPLAY {
		remainCoin, _ := userService.INCRUserCOIN(u.GetUserId(), coin) //增加用户的记录
		u.Coin = proto.Int64(remainCoin)                               //增加用户的金币
	}
}

func (u *MjUser) ADDCountBaGang() {
	atomic.AddInt32(u.Statisc.CountMingGang, 1)
}

func (u *MjUser) SubCountBaGang() {
	atomic.AddInt32(u.Statisc.CountMingGang, -1)
}

//当删除桌子的时候，需要清除掉用户的会话信息,设置deskId = 0 ，roomId = 0
//直接删除玩家的绘画信息
func (u *MjUser) ClearAgentGameData() {
	log.T("清除用户[%v]的session信息为默认状态....", u.GetUserId())
	*u.RoomId = 0
	*u.GameNumber = 0
	agent := u.Agent
	if agent != nil {
		agent.SetUserData(nil)
	}
}

func (u *MjUser) WinAmount2String() string {
	ret := ""
	if u.Bill != nil {
		s := "玩家[%v]总共输赢[%v],下边是细节:"
		s = fmt.Sprintf(s, u.GetUserId(), u.Bill.GetWinAmount())
		ret += s
	}
	return ret
}

//统计玩家某类型账单的条数
func (u *MjUser) getBillBeanTypeCount(billType int32) (count int) {
	if u.GetBill() == nil || u.GetBill().GetBills() == nil || len(u.GetBill().GetBills()) <= 0 {
		return -1
	}
	for _, billBean := range u.GetBill().GetBills() {
		if billBean.GetType() == billType {
			count++
		}
	}
	return count
}

func (u *MjUser) Bill2String(bb *BillBean) string {
	ret := ""
	if bb != nil {
		ret = "玩家[%v] 关联玩家[%v],牌[%v],分数[%v],类型[%v],描述[%v]"
		ret = fmt.Sprintf(ret, bb.GetUserId(), bb.GetOutUserId(), bb.GetPai().LogDes(), bb.GetAmount(), bb.GetType(), bb.GetDes())
	}
	return ret
}

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
		*guoHuInfo.FanShu = 0 //现在都设置为0翻
		u.GameData.GuoHuInfo = append(u.GameData.GuoHuInfo, guoHuInfo)
	}

}

//增加胡牌的信息
func (u *MjUser) AddHuPaiInfo(hu *HuPaiInfo) {
	u.GameData.HuInfo = append(u.GameData.HuInfo, hu)
	u.GameData.HandPai.HuPais = append(u.GameData.HandPai.HuPais, hu.Pai) //增加胡牌
	u.SetStatus(MJUSER_STATUS_HUPAI)
}

//删除过胡的信息
func (u *MjUser) DelGuoHuInfo() error {
	u.GameData.GuoHuInfo = nil
	return nil
}

//是否已经有过胡了
func (u *MjUser) HadGuoHuInfo(fan int32) bool {
	//通用的过胡判断
	if u.GameData.GuoHuInfo == nil || len(u.GameData.GuoHuInfo) <= 0 {
		return false
	}
	//通用1，如果自摸不胡，那么别人点的都不能要
	guoInfos := u.GetGameData().GetGuoHuInfo()
	//长沙麻将的过胡判断
	/**

	 */
	if u.d.IsChangShaMaJiang() {
		//成都麻将的过胡判断
		for _, info := range guoInfos {
			if info.GetSendUserId() == u.GetUserId() {
				log.T("长沙麻将 user:%v 可以胡牌但是腿牌之后，返回有过胡[%v]的信息", u.GetUserId(), guoInfos)
				return true
			}

			if info.GetFanShu() >= fan {
				log.T("长沙麻将 user:%v 可过胡[%v]的翻数大于现在的翻数[%v]，有过胡，不能胡", u.GetUserId(), guoInfos, fan)

				return true
			}
		}

	} else {
		//成都麻将 有过胡直接返回不能胡牌
		log.T("成都麻将 user:%v 有过胡[%v]的信息,直接返回不能胡牌", u.GetUserId(), guoInfos)
		return true
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
		//log.T("user[%v] GetGameData:%v", u.GetUserId(), u.GetGameData())
		//log.T("user[%v] GetHandPai:%v", u.GetUserId(), u.GetGameData().GetHandPai())
		//log.T("user[%v] GetQueFlower:%v", u, u.GetGameData().GetHandPai().GetQueFlower())
		if pai.GetFlower() == u.GameData.HandPai.GetQueFlower() {
			continue
		}

		canhu, _, _, _, _, _ := u.d.HuParser.GetCanHu(u.GameData.HandPai, pai, false, 0, u.d.IsBanker(u))
		if canhu {
			//
			log.T("玩家查叫的时候，查到可以胡牌[%v]", pai.LogDes())
			return true
		}

	}
	return false
}

/****统计信息相关方法****/

//增加用户被巴杠的统计记录
func (u *MjUser) AddStatisticsCountBeiBaGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountBeiBaGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiBaGang, 1) //汇总

	//log.T("用户[%v] 被巴杠+1, 当局被巴杠[%v]次, 汇总被巴杠[%v]次", u.GetUserId(), roundBean.GetCountBeiBaGang(), u.Statisc.GetCountBeiBaGang())
	u.printStatiscLog(round)
	return nil
}

//增加用户巴杠的统计记录
func (u *MjUser) AddStatisticsCountBaGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountBaGnag, 1) //单局

	atomic.AddInt32(u.Statisc.CountBaGang, 1) //汇总

	//log.T("用户[%v] 巴杠+1, 当局巴杠[%v]次, 汇总巴杠[%v]次", u.GetUserId(), roundBean.GetCountBaGnag(), u.Statisc.GetCountBaGang())
	u.printStatiscLog(round)
	return nil
}

//增加用户被暗杠的统计记录
func (u *MjUser) AddStatisticsCountBeiAnGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountBeiAnGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiAnGang, 1) //汇总

	//log.T("用户[%v] 被暗杠+1, 当局被暗杠[%v]次, 汇总被暗杠[%v]次", u.GetUserId(), roundBean.GetCountBeiAnGang(), u.Statisc.GetCountBeiAnGang())
	u.printStatiscLog(round)
	return nil
}

//增加用户暗杠的统计记录
func (u *MjUser) AddStatisticsCountAnGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountAnGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountAnGang, 1) //汇总

	//log.T("用户[%v] 暗杠+1, 当局暗杠[%v]次, 汇总暗杠[%v]次", u.GetUserId(), roundBean.GetCountAnGang(), u.Statisc.GetCountAnGang())
	u.printStatiscLog(round)
	return nil
}

//增加用户明杠的统计记录
func (u *MjUser) AddStatisticsCountMingGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountMingGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountMingGang, 1) //汇总

	//log.T("用户[%v] 明杠+1, 当局明杠[%v]次, 汇总明杠[%v]次", u.GetUserId(), roundBean.GetCountMingGang(), u.Statisc.GetCountMingGang())
	u.printStatiscLog(round)
	return nil
}

//增加用户点杠的统计记录
func (u *MjUser) AddStatisticsCountDianGang(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountDianGang, 1) //单局

	atomic.AddInt32(u.Statisc.CountDianGang, 1) //汇总

	//log.T("用户[%v] 点杠+1, 当局点杠[%v]次, 汇总点杠[%v]次", u.GetUserId(), roundBean.GetCountDianGang(), u.Statisc.GetCountDianGang())
	u.printStatiscLog(round)
	return nil
}

//增加用户胡牌的统计记录
func (u *MjUser) AddStatisticsCountHu(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountHu, 1) //单局

	atomic.AddInt32(u.Statisc.CountHu, 1) //汇总

	//log.T("用户[%v] 胡+1, 当局胡[%v]次, 汇总胡[%v]次", u.GetUserId(), roundBean.GetCountHu(), u.Statisc.GetCountHu())
	u.printStatiscLog(round)
	return nil
}

//增加用户点炮的统计记录
func (u *MjUser) AddStatisticsCountDianPao(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountDianPao, 1) //单局

	atomic.AddInt32(u.Statisc.CountDianPao, 1) //汇总

	//log.T("用户[%v] 点炮+1, 当局点炮[%v]次, 汇总点炮[%v]次", u.GetUserId(), roundBean.GetCountDianPao(), u.Statisc.GetCountDianPao())
	u.printStatiscLog(round)
	return nil
}

//增加用户被自摸的统计记录
func (u *MjUser) AddStatisticsCountBeiZiMo(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountBeiZiMo, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiZiMo, 1) //汇总

	//log.T("用户[%v] 被自摸+1, 当局被自摸[%v]次, 汇总被自摸[%v]次", u.GetUserId(), roundBean.GetCountBeiZiMo(), u.Statisc.GetCountBeiZiMo())
	u.printStatiscLog(round)
	return nil
}

//增加用户自摸的统计记录
func (u *MjUser) AddStatisticsCountZiMo(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountZiMo, 1) //单局

	atomic.AddInt32(u.Statisc.CountZiMo, 1) //汇总

	//log.T("用户[%v] 自摸+1, 当局自摸[%v]次, 汇总自摸[%v]次", u.GetUserId(), roundBean.GetCountZiMo(), u.Statisc.GetCountZiMo())
	u.printStatiscLog(round)
	return nil
}

//增加用户查大叫的统计记录
func (u *MjUser) AddStatisticsCountChaDaJiao(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountChaDaJiao, 1) //单局

	atomic.AddInt32(u.Statisc.CountChaDaJiao, 1) //汇总

	//log.T("用户[%v] 查大叫+1, 当局查大叫[%v]次, 汇总查大叫[%v]次", u.GetUserId(), roundBean.GetCountChaDaJiao(), u.Statisc.GetCountChaDaJiao())
	u.printStatiscLog(round)
	return nil
}

//增加用户被查叫的统计记录
func (u *MjUser) AddStatisticsCountBeiChaJiao(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountBeiChaJiao, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiChaJiao, 1) //汇总

	//log.T("用户[%v] 被查大叫+1, 当局被查叫[%v]次, 汇总被查叫[%v]次", u.GetUserId(), roundBean.GetCountBeiChaJiao(), u.Statisc.GetCountBeiChaJiao())
	u.printStatiscLog(round)
	return nil
}

//增加用户查花猪的统计记录
func (u *MjUser) AddStatisticsCountChaHuaZhu(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountChaHuaZhu, 1) //单局

	atomic.AddInt32(u.Statisc.CountChaHuaZhu, 1) //汇总

	//log.T("用户[%v] 查花猪+1, 当局查花猪[%v]次, 汇总查花猪[%v]次", u.GetUserId(), roundBean.GetCountChaHuaZhu(), u.Statisc.GetCountChaHuaZhu())
	u.printStatiscLog(round)
	return nil
}

//增加用户被查花猪的统计记录
func (u *MjUser) AddStatisticsCountBeiChaHuaZhu(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountBeiChaHuaZhu, 1) //单局

	atomic.AddInt32(u.Statisc.CountBeiChaHuaZhu, 1) //汇总

	//log.T("用户[%v] 被查花猪+1, 当局被查花猪[%v]次, 汇总被查花猪[%v]次", u.GetUserId(), roundBean.GetCountBeiChaHuaZhu(), u.Statisc.GetCountBeiChaHuaZhu())
	u.printStatiscLog(round)
	return nil
}

//增加用户抓鸟的统计记录
func (u *MjUser) AddStatisticsCountCatchBird(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountCatchBird, 1) //单局

	atomic.AddInt32(u.Statisc.CountCatchBird, 1) //汇总

	//log.T("用户[%v] 被自摸+1, 当局被自摸[%v]次, 汇总被自摸[%v]次", u.GetUserId(), roundBean.GetCountBeiZiMo(), u.Statisc.GetCountBeiZiMo())
	u.printStatiscLog(round)
	return nil
}

//增加用户被抓鸟的统计记录
func (u *MjUser) AddStatisticsCountCaughtBird(round int32) error {
	roundBean := u.GetStatisticsRoundBean(round)
	if roundBean == nil {
		log.E("统计的时候出错...")
		return Error.NewFailError("没有找到统计的roundBean，无法统计")
	}

	atomic.AddInt32(roundBean.CountCaughtBird, 1) //单局

	atomic.AddInt32(u.Statisc.CountCaughtBird, 1) //汇总

	//log.T("用户[%v] 被自摸+1, 当局被自摸[%v]次, 汇总被自摸[%v]次", u.GetUserId(), roundBean.GetCountBeiZiMo(), u.Statisc.GetCountBeiZiMo())
	u.printStatiscLog(round)
	return nil
}

//增加账单
func (u *MjUser) AddBill(relationUserid uint32, billType int32, des string, score int64, pai *MJPai, roomType int32) error {
	//用户赢钱的账户,赢钱的账单
	bill := NewBillBean()
	*bill.UserId = u.GetUserId()
	*bill.OutUserId = relationUserid
	*bill.Type = billType
	*bill.Des = des
	*bill.Amount = score //杠牌的收入金额
	bill.Pai = pai
	bill.IsQiShouHu = proto.Bool(false) ///普通胡牌，不是起手胡牌
	u.AddBillBean(bill)

	//计算账单的地方 来加减用户的coin
	u.AddCoin(score, roomType)    //统计用户剩余多少钱
	u.AddStatisticsWinCoin(score) //统计用户输赢多少钱
	return nil
}

//增加账单
func (u *MjUser) AddCSQiShouBill(relationUserid uint32, billType int32, des string, score int64, pai *MJPai, roomType int32) error {
	//用户赢钱的账户,赢钱的账单
	bill := NewBillBean()
	*bill.UserId = u.GetUserId()
	*bill.OutUserId = relationUserid
	*bill.Type = billType
	*bill.Des = des
	*bill.Amount = score //杠牌的收入金额
	bill.Pai = pai
	bill.IsQiShouHu = proto.Bool(true)
	u.AddBillBean(bill)

	//计算账单的地方 来加减用户的coin
	u.AddCoin(score, roomType)    //统计用户剩余多少钱
	u.AddStatisticsWinCoin(score) //统计用户输赢多少钱
	return nil
}

//长沙麻将addBill
func (u *MjUser) AddChangShaBill(relationUserid uint32, billType int32, des string, score int64, pai *MJPai, roomType int32, isBird bool) error {
	//用户赢钱的账户,赢钱的账单
	bill := NewBillBean()
	*bill.UserId = u.GetUserId()
	*bill.OutUserId = relationUserid
	//*bill.Type = MJUSER_BILL_TYPE_YING_HU
	*bill.Type = billType
	*bill.Des = des
	*bill.Amount = score //杠牌的收入金额
	bill.Pai = pai

	*bill.IsBird = isBird //是否是抓鸟

	u.AddBillBean(bill)

	//计算账单的地方 来加减用户的coin
	u.AddCoin(score, roomType)    //统计用户剩余多少钱
	u.AddStatisticsWinCoin(score) //统计用户输赢多少钱
	return nil
}

//比较杠牌之后的叫牌和杠牌之前的叫牌的信息是否一样
func (u *MjUser) AfterGangEqualJiaoPai(beforJiaoPais []*MJPai, gangPai *MJPai) bool {

	//1，获得杠牌之后的手牌
	var chandPai *MJHandPai = util.DeepClone(u.GetGameData().GetHandPai()).(*MJHandPai)
	chandPai.Pais = make([]*MJPai, 0)
	for _, p := range u.GameData.HandPai.Pais {
		if p.GetClientId() != gangPai.GetClientId() {
			chandPai.Pais = append(chandPai.Pais, p)
		}
	}

	//2，通过杠牌之后的手牌 获得此时的叫牌
	afterJiaoPais := u.d.HuParser.GetJiaoPais(chandPai) //AfterGangEqualJiaoPai

	//2,比较beforJiaoPais 和 afterJiaoPais
	if len(chandPai.Pais) != len(beforJiaoPais) {
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

	return true
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

func (u *MjUser) GetUserSortedPais() []*MJPai {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return nil
	} else {
		pais := make([]*MJPai, len(u.GameData.HandPai.Pais))
		copy(pais, u.GameData.HandPai.Pais)
		var list MjPaiList = pais
		sort.Sort(list)
		return pais
	}
}

func (u *MjUser) GetUserPengPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	}
	return ServerPais2string(u.GameData.HandPai.PengPais)

}

func (u *MjUser) GetUserGnagPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	return ServerPais2string(u.GameData.HandPai.GangPais)
}

func (u *MjUser) GetUserHuPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	return ServerPais2string(u.GameData.HandPai.HuPais)
}

func (u *MjUser) GetUserInPaiInfo() string {
	if u.GameData == nil || u.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	return u.GameData.HandPai.InPai.LogDes()

}

func (u *MjUser) GetExchangedCardsInfo() string {
	if u == nil || u.GameData.ExchangeCardsOut == nil {
		return "无"
	}

	return ServerPais2string(u.GameData.ExchangeCardsOut)
}

func (u *MjUser) GetExchangedInCardsInfo() string {
	if u == nil || u.GameData.ExchangeCardsIn == nil {
		return "无"
	}

	s := ""
	for _, p := range u.GameData.ExchangeCardsIn {
		s = s + p.LogDes() + "\t "
	}

	return s
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
		ret = "已坐下"
	default:
	}
	return ret
}

func (u *MjUser) printStatiscLog(round int32) {
	//roundBean := u.GetStatisticsRoundBean(round)
	//log.T("用户[%v] 被巴杠+1, 当局被巴杠[%v]次, 汇总被巴杠[%v]次", u.GetUserId(), roundBean.GetCountBeiBaGang(), u.Statisc.GetCountBeiBaGang())
	//log.T("用户[%v] 巴杠+1, 当局巴杠[%v]次, 汇总巴杠[%v]次", u.GetUserId(), roundBean.GetCountBaGnag(), u.Statisc.GetCountBaGang())
	//log.T("用户[%v] 被暗杠+1, 当局被暗杠[%v]次, 汇总被暗杠[%v]次", u.GetUserId(), roundBean.GetCountBeiAnGang(), u.Statisc.GetCountBeiAnGang())
	//log.T("用户[%v] 暗杠+1, 当局暗杠[%v]次, 汇总暗杠[%v]次", u.GetUserId(), roundBean.GetCountAnGang(), u.Statisc.GetCountAnGang())
	//log.T("用户[%v] 明杠+1, 当局明杠[%v]次, 汇总明杠[%v]次", u.GetUserId(), roundBean.GetCountMingGang(), u.Statisc.GetCountMingGang())
	//log.T("用户[%v] 点杠+1, 当局点杠[%v]次, 汇总点杠[%v]次", u.GetUserId(), roundBean.GetCountDianGang(), u.Statisc.GetCountDianGang())
	//log.T("用户[%v] 胡+1, 当局胡[%v]次, 汇总胡[%v]次", u.GetUserId(), roundBean.GetCountHu(), u.Statisc.GetCountHu())
	//log.T("用户[%v] 点炮+1, 当局点炮[%v]次, 汇总点炮[%v]次", u.GetUserId(), roundBean.GetCountDianPao(), u.Statisc.GetCountDianPao())
	//log.T("用户[%v] 被自摸+1, 当局被自摸[%v]次, 汇总被自摸[%v]次", u.GetUserId(), roundBean.GetCountBeiZiMo(), u.Statisc.GetCountBeiZiMo())
	//log.T("用户[%v] 自摸+1, 当局自摸[%v]次, 汇总自摸[%v]次", u.GetUserId(), roundBean.GetCountZiMo(), u.Statisc.GetCountZiMo())
	//log.T("用户[%v] 查大叫+1, 当局查大叫[%v]次, 汇总查大叫[%v]次", u.GetUserId(), roundBean.GetCountChaDaJiao(), u.Statisc.GetCountChaDaJiao())
	//log.T("用户[%v] 被查大叫+1, 当局被查叫[%v]次, 汇总被查叫[%v]次", u.GetUserId(), roundBean.GetCountBeiChaJiao(), u.Statisc.GetCountBeiChaJiao())
	//log.T("用户[%v] 查花猪+1, 当局查花猪[%v]次, 汇总查花猪[%v]次", u.GetUserId(), roundBean.GetCountChaHuaZhu(), u.Statisc.GetCountChaHuaZhu())
	//log.T("用户[%v] 被查花猪+1, 当局被查花猪[%v]次, 汇总被查花猪[%v]次", u.GetUserId(), roundBean.GetCountBeiChaHuaZhu(), u.Statisc.GetCountBeiChaHuaZhu())

}

//更新用户的信息
func (u *MjUser) UpdateSession(gameStatus int32) {
	//3,更新userSession,返回desk 的信息
	s, _ := sessionService.UpdateSession(u.GetUserId(), gameStatus, int32(ddproto.CommonEnumGame_GID_MAHJONG), u.GetGameNumber(), u.GetRoomId(), u.d.GetDeskId(), u.GetStatus(), u.GetIsBreak(), u.GetIsLeave(), u.GetRoomType(), u.GetRoomPassword())
	if s != nil {
		//给agent设置session
		agent := u.Agent
		if agent != nil {
			agent.SetUserData(s)
		}
	}
}

//更新agent
func (u *MjUser) UpdateAgent(a gate.Agent) {
	//设置为没有断开链接
	*u.IsBreak = false
	*u.IsLeave = false
	u.Agent = a
}

func (user *MjUser) UserPai2String() string {
	result := "玩家[%v]牌的信息,handPais[%v],inpai[%v],pengpais[%v],gangpai[%v]，chipai[%v]"
	result = fmt.Sprintf(result, user.GetUserId(),
		ServerPais2string(user.GameData.HandPai.Pais), user.GameData.HandPai.InPai.LogDes(),
		ServerPais2string(user.GameData.HandPai.PengPais), ServerPais2string(user.GameData.HandPai.GangPais),
		ServerPais2string(user.GetGameData().GetHandPai().GetChiPais()))
	return result
}

//设置为托管模式
func (u *MjUser) setAgentMode(m bool) error {
	u.AgentMode = proto.Bool(m)
	return nil
}

//设置玩家为盛情解散
func (u *MjUser) applayDissolve(b int32) {
	u.ApplyDissolve = proto.Int32(b)
}
