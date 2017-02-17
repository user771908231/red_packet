package skeleton

import (
	"casino_majiang/service/majiang"
	"casino_majianagv2/core/data"
	"time"
	"casino_majianagv2/core/api"
	"fmt"
	"casino_majianagv2/core/majiangv2"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/gate"
	"sync/atomic"
	"casino_common/common/log"
	"casino_common/common/Error"
	"casino_common/common/service/countService"
)

type SkeletonMJUser struct {
	desk            api.MjDesk
	status          *data.MjUserStatus
	userId          uint32
	Coin            int64  //金币
	NickName        string //昵称
	ReadyTimer      *time.Timer
	Bill            *majiang.Bill
	GameData        *data.MJUserGameData
	Statisc         *majiang.MjUserStatisc //统计信息
	a               gate.Agent
	ActTimeoutCount int32                    //
	Log             *countService.T_game_log //任务统计信息
}

//初始化一个user骨架
func NewSkeleconMJUser(userId uint32) *SkeletonMJUser {
	return nil
}

func (user *SkeletonMJUser) Ready() {
	//设置为准备的状态,并且停止准备计时器
	user.status.SetStatus(majiang.MJUSER_STATUS_READY)
	user.status.Ready = true
	if user.ReadyTimer != nil {
		user.ReadyTimer.Stop()
		user.ReadyTimer = nil
	}

}

func (user *SkeletonMJUser) UserPai2String() string {
	result := "玩家[%v]牌的信息,handPais[%v],inpai[%v],pengpais[%v],gangpai[%v]"
	result = fmt.Sprintf(result, user.GetUserId(),
		majiangv2.ServerPais2string(user.GetGameData().HandPai.Pais), user.GetGameData().HandPai.InPai.LogDes(),
		majiangv2.ServerPais2string(user.GetGameData().HandPai.PengPais), majiangv2.ServerPais2string(user.GetGameData().HandPai.GangPais))
	return result
}

//比较杠牌之后的叫牌和杠牌之前的叫牌的信息是否一样
func (u *SkeletonMJUser) AfterGangEqualJiaoPai(beforJiaoPais []*majiang.MJPai, gangPai *majiang.MJPai) bool {

	//1，获得杠牌之后的手牌
	var afterPais []*majiang.MJPai
	for _, p := range u.GameData.HandPai.Pais {
		if p.GetClientId() != gangPai.GetClientId() {
			afterPais = append(afterPais, p)
		}
	}

	//2，通过杠牌之后的手牌 获得此时的叫牌
	afterJiaoPais := u.GetDesk().GetHuParser().GetJiaoPais(afterPais)

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

//初始化方法
func (u *SkeletonMJUser) BeginInit(CurrPlayCount int32, banker uint32) {

}

//发送overTrun
/**
	这里需要区分有托管 和没有托管的状态：
	1，有托管的时候，给玩家发送
 */
func (u *SkeletonMJUser) SendOverTurn(p proto.Message) error {
	//如果是金币场有超时的处理...
	u.WriteMsg(p)
	return nil
}
func (u *SkeletonMJUser) printStatiscLog(round int32) {
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

/****统计信息相关方法****/

//得到每一局的统计bean...
func (u *SkeletonMJUser) GetStatisticsRoundBean(round int32) *majiang.StatiscRound {
	for _, bean := range u.Statisc.RoundBean {
		if bean != nil && bean.GetRound() == round {
			return bean
		}
	}
	//如果没有找到返回nil
	return nil
}

//增加用户被巴杠的统计记录
func (u *SkeletonMJUser) AddStatisticsCountBeiBaGang(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountBaGang(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountBeiAnGang(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountAnGang(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountMingGang(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountDianGang(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountHu(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountDianPao(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountBeiZiMo(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountZiMo(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountChaDaJiao(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountBeiChaJiao(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountChaHuaZhu(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountBeiChaHuaZhu(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountCatchBird(round int32) error {
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
func (u *SkeletonMJUser) AddStatisticsCountCaughtBird(round int32) error {
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

//判断是否是花猪
func (u *SkeletonMJUser) IsHuaZhu() bool {
	for _, pai := range u.GetGameData().HandPai.Pais {
		if pai != nil && pai.GetFlower() == u.GameData.HandPai.GetQueFlower() {
			//是花猪
			return true
		}
	}
	//不是花猪
	return false
}

//lottery之后，设置user为没有准备
func (u *SkeletonMJUser) AfterLottery() error {
	//准备状态
	u.GetStatus().SetStatus(majiang.MJUSER_STATUS_SEATED)
	u.GetStatus().Ready = false     //设置为没有准备的状态...
	u.GetStatus().DingQue = false   //设置为没有定缺的状态...
	u.GetStatus().AgentMode = false //第二局开始默认不准备
	u.UpdateTaskLog()
	return nil
}

//更新玩家任务系统用到的统计信息
func (u *SkeletonMJUser) UpdateTaskLog() {
	u.Log.EndTime = time.Now().Unix()
	u.Log.Bill = float64(u.GetBill().GetWinAmount())
	if u.GetBill().GetWinAmount() > 0 {
		u.Log.IsWine = true
	}
	u.Log.Insert()
}

func (u *SkeletonMJUser) AddGuoHuInfo(checkCase *majiang.CheckCase) {
	if checkCase == nil {
		return
	}
	checkBean := checkCase.GetBeanByUserIdAndStatus(u.GetUserId(), majiang.CHECK_CASE_BEAN_STATUS_CHECKING)
	if checkBean != nil && checkBean.GetCanHu() {
		guoHuInfo := majiang.NewGuoHuInfo()
		*guoHuInfo.SendUserId = checkCase.GetUserIdOut()
		guoHuInfo.Pai = checkCase.CheckMJPai
		*guoHuInfo.FanShu = 0 //现在都设置为0翻
		u.GameData.GuoHuInfo = append(u.GameData.GuoHuInfo, guoHuInfo)
	}

}
