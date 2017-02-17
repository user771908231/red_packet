package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"sync"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_majiang/msg/funcsInit"
	"casino_majiang/service/majiang"
	"casino_common/utils/rand"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majianagv2/core/majiangv2"
	"time"
	"fmt"
	"strings"
	"github.com/name5566/leaf/module"
	"casino_common/common/sessionService"
	"github.com/name5566/leaf/gate"
)

var ERR_SYS = Error.NewError(consts.ACK_RESULT_FAIL, "系统错误")
var ERR_REQ_REPETITION error = Error.NewError(consts.ACK_RESULT_FAIL, "重复请求")
var ERR_ENTER_DESK error = Error.NewError(consts.ACK_RESULT_FAIL, "进入房间失败")
var ERR_COIN_INSUFFICIENT error = Error.NewError(consts.ACK_RESULT_FAIL, "进入金币场失败，金币不足")

var ERR_LEAVE_RUNNING error = Error.NewError(consts.ACK_RESULT_FAIL, "现在不能离开")
var ERR_LEAVE_ERROR error = Error.NewError(consts.ACK_RESULT_FAIL, "出现错误，离开失败")

var ERR_READY_STATE = Error.NewError(consts.ACK_RESULT_FAIL, "准备失败，不在准备阶段")
var ERR_READY_REPETITION = Error.NewError(consts.ACK_RESULT_FAIL, "准备失败，不在准备阶段")
var ERR_READY_COIN_INSUFFICIENT = Error.NewError(consts.ACK_RESULT_FAIL, "准备失败，金币不足")
var ERR_READY_state = Error.NewError(consts.ACK_RESULT_FAIL, "准备失败，不在准备阶段")

//desk 的骨架,业务逻辑的方法 放置在这里
type SkeletonMJDesk struct {
	*sync.Mutex
	*module.Skeleton
	config        *data.SkeletonMJConfig //这里不用使用指针，此配置创建之后不会再改变
	status        *data.MjDeskStatus     //桌子的所有状态都在这里
	HuParser      api.HuPaerApi          //胡牌解析器
	CheckCase     *majiang.CheckCase     //麻将的判定器
	Users         []api.MjUser           //所有的玩家
	AllMJPais     []*majiang.MJPai       //所有的麻将牌
	OverTurnTimer *time.Timer            //定时器
	BirdInfo      []*mjproto.BirdInfo    //抓鸟的信息
	Room          api.MjRoom             //room 的信息
}

func NewSkeletonMJDesk(config *data.SkeletonMJConfig) *SkeletonMJDesk {
	desk := &SkeletonMJDesk{
		config: config,
	}
	return desk
}

func (f *SkeletonMJDesk) EnterUser(userId uint32, a gate.Agent) error {
	log.Debug("玩家[%v]进入fdesk")
	return nil
}

//准备
func (d *SkeletonMJDesk) Ready(userId uint32) error {
	//判断desk状态
	if d.GetStatus().IsNotPreparing() {
		// 准备失败
		log.E("用户[%v]准备失败.desk[%v]不在准备的状态...", userId, d.GetMJConfig().DeskId)
		return ERR_READY_STATE
	}

	//找到需要准备的user
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("用户[%v]在desk[%v]准备的时候失败,没有找到对应的玩家", userId, d.GetMJConfig().DeskId)
		return ERR_SYS
	}

	if user.GetStatus().IsReady() {
		log.E("玩家[%v]已经准备好了...请不要重新准备...", userId)
		return ERR_READY_REPETITION
	}

	//如果是金币场，需要判断玩家的金币是否足够
	//判断金币是否足够,准备的阶段不会扣除房费，房费是在开始的时候扣除

	user.Ready()

	//准备成功,发送准备成功的广播
	result := newProto.NewGame_AckReady()
	*result.Header.Code = consts.ACK_RESULT_SUCC
	*result.Header.Error = "准备成功"
	*result.UserId = userId
	log.T("广播user[%v]在desk[%v]准备成功的广播..string(%v)", userId, d.GetMJConfig().DeskId, result.String())
	d.BroadCastProto(result)

	return nil
}

//检测是否轮到当前玩家打牌...
func (d *SkeletonMJDesk) CheckActUser(userId uint32) bool {
	if d.GetMJConfig().ActUser == userId {
		return true //检测通过
	} else {
		//没有轮到当前玩家
		log.E("[%v]非法请求，没有轮到当前玩家打牌..应该是[%v]", d.DlogDes(), d.GetMJConfig().ActUser)
		return false
	}
}

//检测是否轮到操作
func (d *SkeletonMJDesk) CheckNotActUser(userId uint32) bool {
	return !d.CheckActUser(userId)
}

//定缺
func (f *SkeletonMJDesk) DingQue(userId uint32, color int32) error {
	return nil
}

func (d *SkeletonMJDesk) InitCheckCase(p *majiang.MJPai, outUser api.MjUser) error {
	return nil
}

//处理下一个checkCase
func (d *SkeletonMJDesk) DoCheckCase() error {
	return nil
}

func (d *SkeletonMJDesk) ActGang(userId uint32, c int32, bu bool) error {
	return nil
}

func (d *SkeletonMJDesk) ActHu(userId uint32) error {
	return nil
}

//胡牌解析器
func (d *SkeletonMJDesk) GetHuParser() api.HuPaerApi {
	return d.HuParser
}

func (d *SkeletonMJDesk) Leave(userId uint32) error {
	return nil
}

//指针指向的玩家
func (d *SkeletonMJDesk) SetActiveUser(userId uint32) error {
	d.GetMJConfig().ActiveUser = userId
	return nil
}

func (d *SkeletonMJDesk) Time2Lottery() bool {
	//游戏中的玩家只剩下一个人，表示游戏结束...
	gamingCount := d.GetGamingCount() //正在游戏中的玩家数量

	log.T("判断是否胡牌...当前的gamingCount[%v],当前的PaiCursor[%v]", gamingCount, d.GetMJConfig().MJPaiCursor)
	//1,只剩下一个人的时候. 表示游戏结束
	if gamingCount == 1 {
		return true
	}

	log.T("HandPaiCanMo[%v]", d.HandPaiCanMo())
	//2，没有牌可以摸的时候，返回可以lottery了
	if !d.HandPaiCanMo() {
		return true
	}

	//如果是倒倒胡并且nextCheckCase为空
	if d.IsDaodaohu() && d.GetCheckCase().GetNextBean() == nil {
		for _, user := range d.Users {
			if user != nil && user.GetStatus().IsHu() {
				return true
			}
		}
	}
	return false
}

//发送game_opening 的协议
func (d *SkeletonMJDesk) SendNewGame_Opening() {
	log.T("发送游戏开始的协议..")
	log.T("当前桌子共[%v]把，现在是第[%v]把游戏开始", d.GetMJConfig().TotalPlayCount, d.GetMJConfig().CurrPlayCount)

	open := newProto.NewGame_Opening()
	*open.CurrPlayCount = d.GetMJConfig().CurrPlayCount
	*open.Dice1 = d.GetDice1() //骰子
	*open.Dice2 = d.GetDice2() //骰子
	d.BroadCastProto(open)
}

//通过庄来判断骰子的数目
func (d *SkeletonMJDesk) GetDice1() int32 {
	return rand.Rand(1, 7)
}

func (d *SkeletonMJDesk) GetDice2() int32 {
	return rand.Rand(1, 7)
}

func GettPaiValueByCountPos(countPos int) int32 {
	return int32(countPos%9 + 1)
}

//从pais数组里删除一张pos位置的pai 注 pos是索引值 使用覆盖的方式
func removeFromPais(pais []*majiang.MJPai, pos int) []*majiang.MJPai {
	pais[pos] = pais[len(pais)-1]
	return pais[:len(pais)-1]
}

//将一张pai插入到指定pos的pais数组里去
func addPaiIntoPais(pai *majiang.MJPai, pais []*majiang.MJPai, pos int) []*majiang.MJPai {
	tempPais := make([]*majiang.MJPai, len(pais)+1)
	copy(tempPais[:pos], pais[:pos])
	tempPais[pos] = pai
	copy(tempPais[pos+1:], pais[pos:])
	return tempPais
}

//剩下的牌的数量
func (d *SkeletonMJDesk) GetLeftPaiCount(user api.MjUser, mjPai *majiang.MJPai) int {
	var count int = 0
	displayPais := d.GetDisplayPais(user)
	for i := 0; i < len(displayPais); i++ {
		if (displayPais[i].GetValue() == mjPai.GetValue()) && (displayPais[i].GetFlower() == mjPai.GetFlower()) {
			count++
		}
	}
	count = 4 - count
	if count < 0 {
		count = 0
	}
	//log.T("leftPai is %v Count is : %v", mjPai.GetDes(), count)
	return count
}

//获取用户已知亮出台面的牌 包括自己手牌、自己和其他玩家碰杠牌、其他玩家outPais
func (d *SkeletonMJDesk) GetDisplayPais(user api.MjUser) []*majiang.MJPai {
	//获取所有玩家的亮出台面的牌 outPais + pengPais + gangPais

	displayPais := []*majiang.MJPai{}
	for _, user := range d.GetUsers() {
		userHandPai := user.GetGameData().GetHandPai()

		if userHandPai.GetGangPais() != nil {
			displayPais = append(displayPais, userHandPai.GangPais...) //杠的牌
		}
		if userHandPai.GetPengPais() != nil {
			displayPais = append(displayPais, userHandPai.PengPais...) //碰的牌
		}
		if userHandPai.GetOutPais() != nil {
			displayPais = append(displayPais, userHandPai.OutPais...) //打出去的牌
		}
	}

	//在亮出台面的牌中加入用户自己的手牌
	userHandPai := user.GetGameData().GetHandPai()
	displayPais = append(displayPais, userHandPai.InPai)
	displayPais = append(displayPais, userHandPai.Pais...)
	return displayPais
}

//得到一个canhuinfos
/**
	一次判断打出每一张牌的时候，有哪些牌可以胡，可以胡的翻数是多少
 */
func (d *SkeletonMJDesk) GetJiaoInfos(user api.MjUser) []*mjproto.JiaoInfo {
	log.T("[%v]开始判断玩家[%v]的叫牌...GetJiaoInfos()", d.DlogDes(), user.GetUserId())
	if user == nil ||
		user.GetGameData() == nil ||
		user.GetGameData().HandPai == nil {
		log.E("[%v]开始判断玩家[%v]的叫牌...GetJiaoInfos()失败...因为手牌为nil", d.DlogDes(), user.GetUserId())
		return nil
	}

	jiaoInfos := []*mjproto.JiaoInfo{}

	//获取用户手牌 包括inPai
	userHandPai := majiang.NewMJHandPai()
	*userHandPai = *user.GetGameData().HandPai                //手牌
	userPais := make([]*majiang.MJPai, len(userHandPai.Pais)) //需要改变的牌
	copy(userPais, userHandPai.Pais)
	if userHandPai.InPai != nil {
		//碰牌 无inPai的情况
		userPais = append(userPais, userHandPai.InPai)
	}

	lenth := len(userPais)
	for i := 0; i < lenth; i++ {
		//从用户手牌中移除当前遍历的元素
		removedPai := userPais[i]
		userPais = removeFromPais(userPais, i)
		userHandPai.Pais = userPais
		jiaoInfo := majiang.NewJiaoInfo()

		//遍历麻将牌,看哪一张能胡牌
		for l := 0; l < len(majiangv2.MjpaiMap); l += 4 {

			//遍历未知牌
			//将遍历到的未知牌与用户手牌组合成handPai 去canhu
			mjPai := majiangv2.InitMjPaiByIndex(l)
			canHu, fan, _, _, _, _ := d.HuParser.GetCanHu(userHandPai, mjPai, true, 0)
			if canHu {
				//叫牌的信息
				mjPaiLeftCount := int32(d.GetLeftPaiCount(user, mjPai)) //该可胡牌在桌面中的剩余数量 注 对于自己而言的剩余
				jiaoPaiInfo := majiang.NewJiaoPaiInfo()
				jiaoPaiInfo.HuCard = mjPai.GetCardInfo()
				*jiaoPaiInfo.Fan = fan //可胡番数
				*jiaoPaiInfo.Count = mjPaiLeftCount
				//log.T("[%v],玩家[%v]打牌判断jiaoPaiInfo结果[%v]", d.DlogDes(), user.GetUserId(), jiaoPaiInfo)

				//增加到jiao info
				jiaoInfo.OutCard = removedPai.GetCardInfo() //当前打出去的牌
				jiaoInfo.PaiInfos = append(jiaoInfo.PaiInfos, jiaoPaiInfo)
			}

		}

		//回复手牌
		userPais = addPaiIntoPais(removedPai, userPais, i) //将移除的牌添加回原位置继续遍历
		///如果有叫牌，加入jiaoinfoS
		if jiaoInfo.PaiInfos != nil && len(jiaoInfo.PaiInfos) > 0 {
			jiaoInfos = append(jiaoInfos, jiaoInfo)
		} else {

		}
	}

	log.T("[%v],玩家[%v]判断jiaoInfo结果[%v]", d.DlogDes(), user.GetUserId(), jiaoInfos)
	return jiaoInfos
}

//可以把overturn放在一个地方,目前都是摸牌的时候在用
func (d *SkeletonMJDesk) GetMoPaiOverTurn(userApi api.MjUser, isOpen bool) *mjproto.Game_OverTurn {
	user := d.GetSkeletonMJUser(userApi)
	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = user.GetUserId()                 //这个是摸牌的，所以是广播...
	*overTurn.PaiCount = d.GetRemainPaiCount()          //桌子剩余多少牌
	*overTurn.ActType = majiang.OVER_TURN_ACTTYPE_MOPAI //摸牌
	*overTurn.Time = 30
	if isOpen {
		overTurn.ActCard = user.GetGameData().HandPai.InPai.GetBackPai()
	} else {
		overTurn.ActCard = user.GetGameData().HandPai.InPai.GetCardInfo()
	}

	log.T("[%v]摸牌的时候牌:%v", d.DlogDes(), user.GetSkeletonMJUser().UserPai2String())
	*overTurn.CanHu, _, _, _, _, _ = d.HuParser.GetCanHu(user.GetGameData().GetHandPai(), user.GetGameData().GetHandPai().GetInPai(), true, 0) //是否可以胡牌
	*overTurn.CanPeng = false                                                                                                                  //是否可以碰牌

	//处理杠牌的时候
	/**
		1，血战到底：用户胡牌之后是不会进入到这个方法的
		2，血流成河：用户已经胡牌，那么杠牌之后，胡牌不会改变的情况下，才可以杠 // todo
	 */
	canGangBool, gangPais := user.GetGameData().HandPai.GetCanGang(nil, d.GetRemainPaiCount()) //是否可以杠牌
	*overTurn.CanGang = canGangBool
	if canGangBool && gangPais != nil {
		if user.GetStatus().IsHu() && d.IsXueLiuChengHe() {
			//血流成河，胡牌之后 杠牌的逻辑
			//jiaoPais := user.GetJiaoPaisByHandPais(); //得到杠牌之前的可以胡的叫牌
			jiaoPais := d.HuParser.GetJiaoPais(user.GetGameData().HandPai.Pais)
			for _, g := range gangPais {
				//判断杠牌之后的叫牌是否和杠牌之前一样
				if user.GetSkeletonMJUser().AfterGangEqualJiaoPai(jiaoPais, g) {
					overTurn.GangCards = append(overTurn.GangCards, g.GetCardInfo())
				}
			}
		} else {
			//没有胡牌之前，杠牌的逻辑....
			for _, g := range gangPais {
				overTurn.GangCards = append(overTurn.GangCards, g.GetCardInfo())
			}
		}
	}

	//最后判断是否可以杠牌
	if overTurn.GangCards == nil || len(overTurn.GangCards) <= 0 {
		*overTurn.CanGang = false
	}

	//最后判断是否需要增加过(可以杠，可以胡的时候需要增加可以过的按钮)
	if overTurn.GetCanGang() || overTurn.GetCanHu() {
		overTurn.CanGuo = proto.Bool(true)
	}

	//对长沙麻将做特殊处理
	overTurn.JiaoInfos = d.GetJiaoInfos(userApi)
	return overTurn
}

//通过checkCase 得到一个OverTurn
func (d *SkeletonMJDesk) GetOverTurnByCaseBean(checkPai *majiang.MJPai, caseBean *majiang.CheckBean, actType int32) *mjproto.Game_OverTurn {
	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = caseBean.GetUserId()
	*overTurn.CanGang = caseBean.GetCanGang()
	*overTurn.CanPeng = caseBean.GetCanPeng()
	*overTurn.CanHu = caseBean.GetCanHu()
	*overTurn.PaiCount = d.GetRemainPaiCount() //剩余多少钱
	overTurn.ActCard = checkPai.GetCardInfo()  //
	*overTurn.ActType = actType
	*overTurn.Time = 30
	overTurn.CanGuo = caseBean.CanGuo //目前默认是能过的
	overTurn.CanGuo = proto.Bool(true)
	overTurn.CanChi = caseBean.CanChi
	for i := 0; i < len(caseBean.ChiCards); i += 3 {
		c := &mjproto.ChiOverTurn{}
		c.ChiCard = append(c.ChiCard, caseBean.ChiCards[i].GetCardInfo())
		c.ChiCard = append(c.ChiCard, caseBean.ChiCards[i+1].GetCardInfo())
		c.ChiCard = append(c.ChiCard, caseBean.ChiCards[i+2].GetCardInfo())
		overTurn.ChiInfo = append(overTurn.ChiInfo, c)
	}

	return overTurn
}

//返回一个牌局结果
func (d *SkeletonMJDesk) GetWinCoinInfo(user api.MjUser) *mjproto.WinCoinInfo {
	suser := d.GetSkeletonMJUser(user)
	win := newProto.NewWinCoinInfo()
	*win.NickName = suser.GetNickName()
	*win.UserId = user.GetUserId()
	*win.WinCoin = suser.Bill.GetWinAmount()          //本次输赢多少(负数表示输了)
	*win.Coin = suser.Coin                            // 输赢以后，当前筹码是多少
	*win.CardTitle = d.GetCardTitle4WinCoinInfo(user) // 赢牌牌型信息( 如:"点炮x2 明杠x2 根x2 自摸 3番" )
	log.T("用户[%v]的CardTitle is [%v]", user.GetUserId(), *win.CardTitle)
	win.Cards = suser.GetPlayerCard(true, false) //牌信息,true 表示要显示牌的信息...
	*win.IsDealer = suser.GetStatus().IsBanker   //是否是庄家
	//*win.HuCount = user.Statisc.GetCountHu()        //本局胡的次数(血流成河会多次胡)
	roundBean := suser.GetStatisticsRoundBean(d.GetMJConfig().CurrPlayCount)
	*win.HuCount = roundBean.GetCountHu() + roundBean.GetCountZiMo() //本局胡的次数(血流成河会多次胡)
	log.T("用户[%v]的HuCount is [%v]", user.GetUserId(), *win.HuCount)
	return win
}

//得到这个人的胡牌描述
func (d *SkeletonMJDesk) GetCardTitle4WinCoinInfo(suser api.MjUser) string {
	user := d.GetSkeletonMJUser(suser)
	var huDes []string

	//获取当局的统计信息
	roundBean := user.GetStatisticsRoundBean(d.GetMJConfig().CurrPlayCount)
	var count int32 = 0 //统计数大于count才显示

	mingGangCount := roundBean.GetCountMingGang()         //明杠
	baGangCount := roundBean.GetCountBaGnag()             //巴杠
	beiBaGangCount := roundBean.GetCountBeiBaGang()       //被巴杠
	anGangCount := roundBean.GetCountAnGang()             //暗杠
	beiAnGangCount := roundBean.GetCountBeiAnGang()       //被暗杠
	dianPaoCount := roundBean.GetCountDianPao()           //点炮
	dianGangCount := roundBean.GetCountDianGang()         //点杠
	beiZimoCount := roundBean.GetCountBeiZiMo()           //被自摸
	zimoCount := roundBean.GetCountZiMo()                 //自摸
	beiChaHuaZhuCount := roundBean.GetCountBeiChaHuaZhu() //被查花猪
	beiChaJiaoCount := roundBean.GetCountBeiChaJiao()     //被查叫
	chaJiaoCount := roundBean.GetCountChaDaJiao()         //查大叫

	catchBirdCount := roundBean.GetCountCatchBird()   //抓鸟
	caughtBirdCount := roundBean.GetCountCaughtBird() //被抓鸟

	//log.T("user[%v] roundBean is %v", user.GetUserId(), roundBean)

	if zimoCount > count {
		//log.T("user[%v] zimoCount[%v]", user.GetUserId(), zimoCount)
		huDes = append(huDes, fmt.Sprintf("自摸x%d", zimoCount))
	}

	if mingGangCount > count {
		//log.T("user[%v] mingGangCount[%v]", user.GetUserId(), mingGangCount)
		huDes = append(huDes, fmt.Sprintf("明杠x%d", mingGangCount))
	}

	if baGangCount > count {
		//log.T("user[%v] baGangCount[%v]", user.GetUserId(), baGangCount)
		huDes = append(huDes, fmt.Sprintf("巴杠x%d", baGangCount))
	}

	if anGangCount > count {
		//log.T("user[%v] anGangCount[%v]", user.GetUserId(), anGangCount)
		huDes = append(huDes, fmt.Sprintf("暗杠x%d", anGangCount))
	}

	if chaJiaoCount > count {
		//log.T("user[%v] chaJiaoCount[%v]", user.GetUserId(), chaJiaoCount)
		huDes = append(huDes, fmt.Sprintf("查大叫x%d", chaJiaoCount))
	}

	if beiBaGangCount > count {
		//log.T("user[%v] beiBaGangCount[%v]", user.GetUserId(), anGangCount)
		huDes = append(huDes, fmt.Sprintf("被巴杠x%d", beiBaGangCount))
	}

	if beiAnGangCount > count {
		//log.T("user[%v] beiAnGangCount[%v]", user.GetUserId(), beiAnGangCount)
		huDes = append(huDes, fmt.Sprintf("被暗杠x%d", beiAnGangCount))
	}

	if dianPaoCount > count {
		//log.T("user[%v] dianPaoCount[%v]", user.GetUserId(), dianPaoCount)
		huDes = append(huDes, fmt.Sprintf("点炮x%d", dianPaoCount))
	}

	if dianGangCount > count {
		//log.T("user[%v] dianGangCount[%v]", user.GetUserId(), dianGangCount)
		huDes = append(huDes, fmt.Sprintf("点杠x%d", dianGangCount))
	}

	if beiZimoCount > count {
		//log.T("user[%v] beiZimoCount[%v]", user.GetUserId(), beiZimoCount)
		huDes = append(huDes, fmt.Sprintf("被自摸x%d", beiZimoCount))
	}

	if beiChaHuaZhuCount > count {
		//log.T("user[%v] beiChaHuaZhuCount[%v]", user.GetUserId(), beiChaHuaZhuCount)
		huDes = append(huDes, fmt.Sprintf("被查花猪x%d", beiChaHuaZhuCount))
	}

	if beiChaJiaoCount > count {
		//log.T("user[%v] beiChaJiaoCount[%v]", user.GetUserId(), beiChaJiaoCount)
		huDes = append(huDes, fmt.Sprintf("被查叫x%d", beiChaJiaoCount))
	}

	if catchBirdCount > count {
		//log.T("user[%v] beiChaJiaoCount[%v]", user.GetUserId(), beiChaJiaoCount)
		huDes = append(huDes, fmt.Sprintf("抓鸟x%d", catchBirdCount))
	}

	if caughtBirdCount > count {
		//log.T("user[%v] beiChaJiaoCount[%v]", user.GetUserId(), beiChaJiaoCount)
		huDes = append(huDes, fmt.Sprintf("被抓鸟x%d", caughtBirdCount))
	}

	//获取胡番的信息
	if user.GameData.HuInfo != nil && len(user.GameData.HuInfo) > 0 {
		huDes = append(huDes, user.GameData.HuInfo[0].GetHuDesc())
	}
	log.T("用户[%v]GameData.HuInfo is [%v]", user.GetUserId(), user.GetGameData().GetHuInfo())
	log.T("用户[%v] huDes is [%v]", user.GetUserId(), huDes)
	s := strings.Join(huDes, " ")
	log.T("用户[%v] DesStr is [%v]", user.GetUserId(), s)
	return s
}

/**
玩家开始进入房间
1,如果是朋友桌，不做操作
2,如果是金币场，到了时间之后增加机器人
 */
func (d *SkeletonMJDesk) beginEnter() error {
	d.GetStatus().SetStatus(majiang.MJDESK_STATUS_READY) //桌子开始ready
	return nil
}

//删除一个user
func (d *SkeletonMJDesk) RmUser(user api.MjUser) error {
	for i, u := range d.Users {
		if u != nil && u.GetUserId() == user.GetUserId() {
			//更新session
			d.Users[i] = nil
			log.T("%v d.rmUser(%v)", d.DlogDes(), u.GetUserId())
			sessionService.DelSessionByKey(u.GetUserId(), d.GetMJConfig().RoomType)

		}
	}

	return nil
}

//doend
func (d *SkeletonMJDesk) DoEnd() error {

	//1
	//首先发送游戏 结束的广播....game_SendEndLottery
	result := newProto.NewGame_SendEndLottery()
	for _, user := range d.GetUsers() {
		if user != nil {
			result.CoinInfo = append(result.CoinInfo, d.GetEndLotteryInfo(user))
		}
	}

	//发送之前需要判断谁是大赢家...这里暂时没有判断...

	//发送游戏结束的结果
	d.BroadCastProto(result)

	//2,清楚数据，解散房间....
	d.Room.DissolveDesk(d, false)
	return nil
}

//得到EndLotteryInfo结果...
func (d *SkeletonMJDesk) GetEndLotteryInfo(suser api.MjUser) *mjproto.EndLotteryInfo {
	user := d.GetSkeletonMJUser(suser)
	end := newProto.NewEndLotteryInfo()
	*end.UserId = user.GetUserId()
	*end.BigWin = false                                  //是否是大赢家...
	*end.CountAnGang = user.Statisc.GetCountAnGang()     //暗杠的次数
	*end.CountChaJiao = user.Statisc.GetCountChaDaJiao() //查叫的次数..
	*end.CountDianGang = user.Statisc.GetCountDianGang() // 点杠的次数
	*end.CountDianPao = user.Statisc.GetCountDianPao()   //点炮的次数
	*end.CountHu = user.Statisc.GetCountHu()             //胡牌的次数
	*end.CountZiMo = user.Statisc.GetCountZiMo()         //自摸的次数
	*end.WinCoin = user.Statisc.GetWinCoin()             //赢了多少钱
	*end.CountMingGang = user.Statisc.GetCountMingGang() //明杠
	return end
}

//巴杠之后需要初始化抢杠的CheckCase
func (d *SkeletonMJDesk) InitCheckCaseAfterGang(gangType int32, gangPai *majiang.MJPai, user api.MjUser) {
	d.CheckCase = nil //设置 判断的为nil
	///如果是巴杠，需要设置巴杠的判断  initCheckCase
	if gangType == majiang.GANG_TYPE_BA {
		d.InitCheckCase(gangPai, user) //杠牌之后 初始化checkCase
		if d.CheckCase == nil {
			log.T("巴杠没有人可以抢杠...")
		}
	}

}
