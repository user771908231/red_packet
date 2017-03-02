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
	"casino_common/proto/ddproto"
	"casino_common/common/sessionService"
	"casino_common/common/userService"
	"github.com/name5566/leaf/timer"
)

type SkeletonMJUser struct {
	desk            api.MjDesk
	userId          uint32
	UserStatus      *data.MjUserStatus //状态
	Coin            int64              //金币
	NickName        string             //昵称
	Sex             int32              //性别
	RobotType       int32              //机器人类型
	ReadyTimer      *timer.Timer
	Bill            *majiang.Bill
	GameData        *data.MJUserGameData
	Statisc         *majiang.MjUserStatisc //统计信息
	a               gate.Agent
	ActTimeoutCount int32                    //
	Log             *countService.T_game_log //任务统计信息
}

//初始化一个user骨架
func NewSkeletonMJUser(desk api.MjDesk, userId uint32, a gate.Agent) *SkeletonMJUser {
	//清空agent 的userData
	redisUser := userService.GetUserById(userId)
	if redisUser == nil {
		log.E("系统中找不到用户%v", userId)
		return nil
	}
	//返回用户信息
	return &SkeletonMJUser{
		desk:      desk,
		userId:    userId,
		a:         a,
		NickName:  redisUser.GetNickName(),
		Sex:       redisUser.GetSex(),
		RobotType: redisUser.GetRobotType(),
		UserStatus: &data.MjUserStatus{
			Status:   1,
			IsBanker: desk.GetMJConfig().Banker == userId,
			IsLeave:  false,
			IsBreak:  false,

		},
		GameData: &data.MJUserGameData{
			PlayerGameData: majiang.NewPlayerGameData(),
		},
		Statisc: majiang.NewMjUserStatisc(),
		Log:     &countService.T_game_log{},
	}
}

func (user *SkeletonMJUser) ActReady() {
	//设置为准备的状态,并且停止准备计时器
	user.UserStatus.SetStatus(majiang.MJUSER_STATUS_READY) //设置为准备的状态
	user.UserStatus.Ready = true
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
func (u *SkeletonMJUser) BeginInit(round int32, banker uint32) error {
	//1,游戏开始时候的初始化...
	u.GameData = &data.MJUserGameData{
		PlayerGameData: majiang.NewPlayerGameData(),
	} //初始化一个空的麻将牌
	u.GetStatus().DingQue = false
	u.GetStatus().Exchange = false
	if u.GetUserId() == banker {
		u.GetStatus().IsBanker = true
	} else {
		u.GetStatus().IsBanker = false
	}
	//杠牌信息
	u.GetGameData().DelPreMoGangInfo()
	//初始化账单
	u.Bill = majiang.NewBill()
	//2,初始化统计bean
	statisticsRoundBean := majiang.NewStatiscRound()
	*statisticsRoundBean.Round = round
	u.Statisc.RoundBean = append(u.Statisc.RoundBean, statisticsRoundBean)
	u.ActTimeoutCount = 0 //初始化超时的次数
	u.initTaskLog()

	return nil
}

func (u *SkeletonMJUser) initTaskLog() {
	u.Log.UserId = u.GetUserId()
	u.Log.Bill = 0
	u.Log.GameId = ddproto.CommonEnumGame_GID_MAHJONG
	u.Log.GameNumber = u.GetDesk().GetMJConfig().GameNumber
	u.Log.RoomType = ddproto.COMMON_ENUM_ROOMTYPE(u.GetDesk().GetMJConfig().RoomType)
	u.Log.RoomLevel = u.GetDesk().GetMJConfig().RoomLevel
	u.Log.StartTime = time.Now().Unix()
	u.Log.IsWine = false
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
	u.GetStatus().SetStatus(majiang.MJUSER_STATUS_SEATED) //设置为坐下的状态
	u.GetStatus().Ready = false                           //设置为没有准备的状态...
	u.GetStatus().DingQue = false                         //设置为没有定缺的状态...
	u.GetStatus().AgentMode = false                       //第二局开始默认不准备
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

//进入房间
func (user *SkeletonMJUser) ReEnterDesk(a gate.Agent) error {
	user.UpdateAgent(a)
	user.UpdateSession(int32(ddproto.COMMON_ENUM_GAMESTATUS_GAMING))
	return nil
}

func (u *SkeletonMJUser) UpdateAgent(a gate.Agent) error {
	log.T("玩家[%v]断线重连，进入房间UpdateAgent(a)", u.GetUserId())
	oldAgent := u.a
	if oldAgent != nil {
		//需要做处理
		log.T("断线重连，强制断开玩家[%v]老的链接", u.GetUserId())
		oldAgent.SetUserData(nil) //清空清空会话状态
		//oldAgent.Close()
	}
	//设置为没有断开链接
	u.GetStatus().IsBreak = false
	u.GetStatus().IsLeave = false
	u.a = a
	return nil

}

func (u *SkeletonMJUser) UpdateSession(gameStatus int32) error {
	//3,更新userSession,返回desk 的信息
	s, _ := sessionService.UpdateSession(
		u.GetUserId(),
		gameStatus,
		int32(ddproto.CommonEnumGame_GID_MAHJONG),
		u.GetDesk().GetMJConfig().GameNumber,
		u.GetDesk().GetMJConfig().RoomId,
		u.GetDesk().GetMJConfig().DeskId,
		u.GetDesk().GetMJConfig().Status,
		u.GetStatus().IsBreak,
		u.GetStatus().IsLeave,
		u.GetDesk().GetMJConfig().RoomType,
		u.GetDesk().GetMJConfig().Password)
	if s != nil {
		//给agent设置session
		agent := u.a
		if agent != nil {
			agent.SetUserData(s)
		}
	}
	return nil
}

//得到判定bean
func (u *SkeletonMJUser) GetCheckBean(p *majiang.MJPai, remainPaiCoun int32) *majiang.CheckBean {
	bean := majiang.NewCheckBean()

	*bean.CheckStatus = majiang.CHECK_CASE_BEAN_STATUS_CHECKING
	*bean.UserId = u.GetUserId()
	bean.CheckPai = p
	var fan int32 = 0

	//是否可以胡牌
	if u.IsCanInitCheckCaseHu() {
		*bean.CanHu, fan, _, _, _, _ = u.GetDesk().GetHuParser().GetCanHu(u.GameData.HandPai, p, false, 0)
	}
	//是否可以杠
	if u.IsCanInitCheckCaseGang() {
		*bean.CanGang, _ = u.GameData.HandPai.GetCanGang(p, remainPaiCoun)
	}
	//是否可以碰
	if u.IsCanInitCheckCasePeng() {
		*bean.CanPeng = u.GameData.HandPai.GetCanPeng(p)
	}

	log.T("得到用户[%v]对牌[%v]的check , bean[%v]", u.GetUserId(), p.LogDes(), bean)
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
func (u *SkeletonMJUser) IsCanInitCheckCaseGang() bool {
	//这里需要判断是否是 血流成河，目前暂时不判断...

	//1,普通规则
	if u.GetStatus().IsNotHu() {
		return true
	}

	//2,血流成河
	if u.GetStatus().IsHu() && u.GetDesk().GetMJConfig().XueLiuChengHe {
		return true
	}

	//其他情况返回false
	return false
}

func (u *SkeletonMJUser) IsCanInitCheckCasePeng() bool {
	//1,普通规则
	if u.GetStatus().IsNotHu() {
		return true
	} else {
		return false;
	}
}

//判断用户是否可以杠
func (u *SkeletonMJUser) IsCanInitCheckCaseHu() bool {
	return u.IsCanInitCheckCaseGang()
}

//是否已经有过胡了
func (u *SkeletonMJUser) HadGuoHuInfo(fan int32) bool {
	//通用的过胡判断
	if u.GameData.GuoHuInfo == nil || len(u.GameData.GuoHuInfo) <= 0 {
		return false
	}
	return true
}
