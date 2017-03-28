package majiang

import (
	"fmt"
	"strings"
	mjproto        "casino_mj_changsha/msg/protogo"
	"casino_mj_changsha/msg/funcsInit"
	"casino_common/common/log"
	"casino_common/common/userService"
)

/**
	1，只剩一个玩家没有胡牌
	2, 已经没有牌了...
 */

func (d *MjDesk) Time2Lottery() bool {
	//游戏中的玩家只剩下一个人，表示游戏结束...
	gamingCount := d.GetGamingCount() //正在游戏中的玩家数量

	log.T("%v判断是否Time2Lottery...当前的gamingCount[%v],当前的PaiCursor[%v]", d.DlogDes(), gamingCount, d.GetMJPaiCursor())

	//1,只剩下一个人的时候. 表示游戏结束
	if gamingCount == 1 {
		return true
	}

	log.T("%v判断是否Time2Lottery...HandPaiCanMo[%v]", d.DlogDes(), d.HandPaiCanMo())
	//2，没有牌可以摸的时候，返回可以lottery了
	if !d.HandPaiCanMo() {
		return true
	}

	//如果是倒倒胡并且nextCheckCase为空
	if d.IsDaodaohu() && d.GetCheckCase().GetNextBean() == nil {
		for _, user := range d.Users {
			if user != nil && user.IsHu() {
				return true
			}
		}
	}
	return false
}

//成都麻将的 结算方式
func (d *MjDesk) LotteryChengDu() error {
	//结账需要分两中情况
	/**
		1，只剩一个玩家没有胡牌的时候
		2，没有生育麻将的时候.需要分别做处理...
	 */

	//判断是否可以胡牌
	log.T("现在开始处理lottery()的逻辑....")

	//查花猪
	d.ChaHuaZhu()

	//查大叫
	d.ChaDaJiao()

	//1，处理开奖的数据,
	d.DoLottery()

	//发送结束的广播
	d.SendLotteryData()

	//开奖之后 desk需要处理
	d.AfterLottery()

	//判断牌局结束(整场游戏结束)
	if !d.End() {
		//go d.begin()
	}
	return nil
}

// 一盘麻将结束....这里需要针对每个人结账...并且对desk和user的数据做清楚...
func (d *MjDesk) Lottery() error {
	//长沙麻将
	if d.IsChangShaMaJiang() {
		return d.LotteryChangSha()
	}
	//默认 成都麻将
	return d.LotteryChengDu()
}

//查花猪
/**
	查花猪是查用户是否没有缺
 */
func (d *MjDesk) ChaHuaZhu() error {
	for _, u := range d.GetUsers() {
		if !d.IsXueLiuChengHe() && u.IsHu() {
			//如果不是血流成河且用户已胡 则不能查u的花猪
			log.T("不查[%v]的花猪", u.GetUserId())
			continue
		}
		if u != nil && u.IsNotHu() {
			if u.IsHuaZhu() {
				log.T("玩家[%v]是花猪", u.GetUserId())
				d.DoHuaZhu(u)
			}
		}
	}
	return nil
}

//花猪玩家需要给封顶分数
func (d *MjDesk) DoHuaZhu(huazhu *MjUser) error {
	log.T("开始处理花猪[%v]", huazhu.GetUserId())
	fanTop := d.GetRoomTypeInfo().GetCapMax()
	score := d.GetBaseValue() * fanTop
	for _, user := range d.GetUsers() {
		if user != nil && user.IsNotHuaZhu() {
			//不是花猪的用户都可以收钱
			//判断不是花猪，可以赢钱...
			log.T("DoHuaZhu: 查[%v]的花猪", user.GetUserId())
			user.AddBill(huazhu.GetUserId(), MJUSER_BILL_TYPE_YING_CHAHUAZHU, "用户查花猪，赢钱", score, nil, d.GetRoomType())
			user.AddStatisticsCountChaHuaZhu(d.GetCurrPlayCount())

			huazhu.AddBill(user.GetUserId(), MJUSER_BILL_TYPE_SHU_CHAHUAZHU, "用户查花猪，输钱", -score, nil, d.GetRoomType())
			huazhu.AddStatisticsCountBeiChaHuaZhu(d.GetCurrPlayCount())
		}
	}

	return nil
}

//查大叫
/**
	查用户有没有叫
 */
func (d *MjDesk) ChaDaJiao() error {
	//循环判断谁可以被查叫
	for _, u := range d.GetUsers() {
		d.DoDaJiao(u)
	}
	return nil

}

//返回一个牌局结果
func (d *MjDesk) GetWinCoinInfo(user *MjUser) *mjproto.WinCoinInfo {
	win := newProto.NewWinCoinInfo()
	*win.NickName = user.GetNickName()
	*win.UserId = user.GetUserId()
	*win.WinCoin = user.Bill.GetWinAmount()           //本次输赢多少(负数表示输了)
	*win.Coin = user.GetCoin()                        // 输赢以后，当前筹码是多少
	*win.CardTitle = d.GetCardTitle4WinCoinInfo(user) // 赢牌牌型信息( 如:"点炮x2 明杠x2 根x2 自摸 3番" )
	log.T("%v用户[%v]的CardTitle is [%v]", d.DlogDes(), user.GetUserId(), *win.CardTitle)
	win.Cards = user.GetPlayerCard(true) //牌信息,true 表示要显示牌的信息...
	*win.IsDealer = user.GetIsBanker()   //是否是庄家
	//*win.HuCount = user.Statisc.GetCountHu()        //本局胡的次数(血流成河会多次胡)
	roundBean := user.GetStatisticsRoundBean(d.GetCurrPlayCount())
	*win.HuCount = roundBean.GetCountHu() + roundBean.GetCountZiMo() //本局胡的次数(血流成河会多次胡)
	log.T("%v用户[%v]的HuCount is [%v]", d.DlogDes(), user.GetUserId(), *win.HuCount)
	return win
}

//得到这个人的胡牌描述
func (d *MjDesk) GetCardTitle4WinCoinInfo(user *MjUser) string {
	var huDes []string

	//获取当局的统计信息
	roundBean := user.GetStatisticsRoundBean(d.GetCurrPlayCount())
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
	log.T("%v用户[%v]GameData.HuInfo[%v],huDes[%v]", d.DlogDes(), user.GetUserId(), user.GetGameData().GetHuInfo(), huDes)
	s := strings.Join(huDes, " ")
	return s
}

//得到EndLotteryInfo结果...
func (d *MjDesk) GetEndLotteryInfo(user *MjUser) *mjproto.EndLotteryInfo {
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

func (d *MjDesk) AfterLottery() error {
	//开奖完成之后的一些处理
	if d.overTurnTimer != nil {
		d.overTurnTimer.Stop()
	}

	//把信息更新到mgo
	for _, u := range d.GetUsers() {
		if u != nil {
			userService.UpdateUser2MgoById(u.GetUserId()) //afterlottery 之后更新玩家的信息
		}
	}

	//desk lottery 处理之后，开始等待新的玩家进入
	d.beginEnter() //一局结束之后的处理

	return nil

}

func (d *MjDesk) End() bool {
	//判断结束的条件,目前只有局数能判断
	log.T("判断desk[%v],round[%v]游戏是否End() CurrPlayCount[%v], TotalPlayCount[%v]", d.GetDeskId(), d.GetCurrPlayCount(), d.GetCurrPlayCount(), d.GetTotalPlayCount())
	//朋友桌有整场结束的概念
	if d.GetCurrPlayCount() < d.GetTotalPlayCount() {
		//表示游戏还没有结束。。。.
		return false
	} else {
		d.DoEnd()
		return true
	}
	return false

}

func (d *MjDesk) DoEnd() error {

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
	MjroomManagerIns.GetFMJRoom().DissolveDesk(d, false)

	return nil
}

func (d *MjDesk) SendLotteryData() error {
	//发送开奖的数据,需要得到每个人的winCoinInfo
	result := newProto.NewGame_SendCurrentResult()
	for _, user := range d.GetUsers() {
		if user != nil {
			result.WinCoinInfo = append(result.WinCoinInfo, d.GetWinCoinInfo(user))
		}
	}
	result.BridInfo = d.BirdInfo

	//开始发送开奖的广播
	d.BroadCastProto(result)
	log.T("desk(%v),gameNumber(%v)SendLotteryData()", d.GetDeskId(), d.GetGameNumber())
	return nil
}
