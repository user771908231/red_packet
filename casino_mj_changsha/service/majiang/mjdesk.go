package majiang

import (
	"errors"
	mjproto "casino_mj_changsha/msg/protogo"
	"casino_mj_changsha/msg/funcsInit"
	"time"
	"casino_mj_changsha/gamedata/model"
	"casino_common/common/log"
	"casino_common/utils/db"
	"casino_common/utils/timeUtils"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_common/common/consts/tableName"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/util"
	"casino_common/proto/ddproto"
	"github.com/name5566/leaf/gate"
	"math"
)

//状态表示的是当前状态.
var MJDESK_STATUS_READY int32 = 2    //正在准备
var MJDESK_STATUS_OPENNING int32 = 3 //正在准备
var MJDESK_STATUS_QISHOUHU int32 = 4 //起手胡牌增加工作量
var MJDESK_STATUS_EXCHANGE int32 = 5 //desk初始化完成之后，告诉玩家可以开始换牌
var MJDESK_STATUS_DINGQUE int32 = 6  //换牌结束之后，告诉玩家可以开始定缺
var MJDESK_STATUS_RUNNING int32 = 7  //定缺之后，开始打牌

//over turn type

var OVER_TURN_ACTTYPE_MOPAI int32 = 1 //摸牌的类型...
var OVER_TURN_ACTTYPE_OTHER int32 = 2 //碰OTHER

//act type
var MJDESK_ACT_TYPE_MOPAI int32 = 1                   //摸牌
var MJDESK_ACT_TYPE_DAPAI int32 = 2                   //打牌
var MJDESK_ACT_TYPE_WAIT_CHECK int32 = 3              //等待check
var MJDESK_ACT_TYPE_WAIT_HAIDI int32 = 4              //回复是否需要海底牌
var MJDESK_ACT_TYPE_BANK_FIRST_MOPAI int32 = 5        //庄第一次打牌
var MJDESK_ACT_TYPE_WAIT_CHECK_CHANGSHAGANG int32 = 6 //对长沙麻将杠需要做单独的处理

var SHAIZI_SLEEP_DURATION time.Duration = time.Second * 4       //定缺的延迟
var APPLYDISSOLVE_DURATION time.Duration = time.Second * 60 * 5 //申请房间 别人拒绝或者同意的超时时间
var COIN_RAEDY_DURATION time.Duration = time.Second * 60        //金币场，一局结束之后，需要准备，超时不准备需要提出房间

//所有的error 都定义在这里
var ERR_SYS = Error.NewError(consts.ACK_RESULT_FAIL, "系统错误")
var ERR_REQ_REPETITION error = Error.NewError(consts.ACK_RESULT_FAIL, "重复请求")
var ERR_ENTER_DESK error = Error.NewError(consts.ACK_RESULT_FAIL, "进入房间失败,人数已经满了")

//小于最低倍场次

var ERR_LEAVE_RUNNING error = Error.NewError(consts.ACK_RESULT_FAIL, "现在不能离开")
var ERR_LEAVE_ERROR error = Error.NewError(consts.ACK_RESULT_FAIL, "出现错误，离开失败")

//初始化checkCase
//如果出错 设置checkCase为nil
func (d *MjDesk) InitCheckCase(p *MJPai, outUser *MjUser, isBagang bool) error {
	if isBagang {
		log.T("%v 开始初始化outUser[%v] pai[%v]抢杠的iniCheckCase ", d.DlogDes(), outUser.GetUserId(), p.LogDes())
	}
	checkCase := NewCheckCase()
	checkCase.QiangGang = proto.Bool(isBagang)
	checkCase.DianPaoCount = proto.Int32(0) //设置点炮的次数为0
	*checkCase.UserIdOut = outUser.GetUserId()
	*checkCase.CheckStatus = CHECK_CASE_STATUS_CHECKING //正在判定
	checkCase.CheckMJPai = p
	checkCase.PreOutGangInfo = outUser.GetPreMoGangInfo()
	d.CheckCase = checkCase

	//这里需要对checkCase排序
	outUserIndex := d.getIndexByUserId(outUser.GetUserId())
	for i := (outUserIndex + 1); i < (len(d.GetUsers()) + outUserIndex); i++ {
		checkUser := d.GetUsers()[i%(len(d.GetUsers()))]
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() {
			log.T("%v用户[%v]打牌[%v],巴杠[%v]，判断user[%v]是否可以碰杠胡.手牌[%v]", checkUser.d.DlogDes(), outUser.GetUserId(), !isBagang, isBagang, checkUser.GetUserId(), checkUser.UserPai2String())
			//添加checkBean
			bean := checkUser.GetCheckBean(p, d.IsXueLiuChengHe(), d.GetRemainPaiCount(), isBagang)
			if bean != nil {
				checkCase.CheckB = append(checkCase.CheckB, bean)
			}
		}
	}

	log.T("%v判断最终的checkCase[%v]", d.DlogDes(), checkCase)
	if checkCase.CheckB == nil || len(checkCase.CheckB) == 0 {
		d.CheckCase = nil
	}

	return nil
}

//执行判断事件
/**
	这里仅仅是用于判断打牌之后别人的碰杠胡
	这里仅仅是用于判断打牌之后别人的碰杠胡
	1,首先询问胡牌的人,如果有人胡，再询问下一个要胡的人
	2,再[依次序]询问碰杠的人，如果没有人碰杠，再询问吃的人，注意：
		在长沙麻将中，有可能是两个人都可以碰杠，所以需要依次序询问

 */
func (d *MjDesk) DoCheckCase() error {
	//检测参数
	if d.CheckCase.GetNextBean() == nil {
		log.T("[%v]已经没有需要处理的CheckCase,下一个玩家摸牌...", d.DlogDes())
		//直接跳转到下一个操作的玩家...,这里表示判断已经玩了...
		d.CheckCase = nil
		//在这之前需要保证 activeUser 是正确的...
		d.SendMopaiOverTurn()
		return nil
	} else {
		//1,找到胡牌的人来进行处理
		caseBean := d.CheckCase.GetNextBean()
		log.T("继续处理CheckCase,开处理下一个checkBean：%v", caseBean)
		//找到需要判断bean之后，发送给判断人	//send overTurn
		overTurn := d.GetOverTurnByCaseBean(d.CheckCase.CheckMJPai, caseBean, OVER_TURN_ACTTYPE_OTHER) //别人打牌，判断是否可以碰杠胡

		///发送overTurn 的信息
		//log.T("%v 开始发送overTurn[%v]", d.DlogDes(), overTurn)
		d.GetUserByUserId(caseBean.GetUserId()).SendOverTurn(overTurn)        //DoCheckCase
		d.SetActUserAndType(caseBean.GetUserId(), MJDESK_ACT_TYPE_WAIT_CHECK) //长沙麻将 DoCheckCase 设置当前活动的玩家
		return nil
	}
}

//得到一个canhuinfos
/**
	一次判断打出每一张牌的时候，有哪些牌可以胡，可以胡的翻数是多少
 */
func (d *MjDesk) GetJiaoInfos(user *MjUser) []*mjproto.JiaoInfo {
	log.T("[%v]开始判断玩家[%v]的叫牌...GetJiaoInfos()", d.DlogDes(), user.GetUserId())
	if user == nil ||
		user.GameData == nil ||
		user.GameData.HandPai == nil {
		log.E("[%v]开始判断玩家[%v]的叫牌...GetJiaoInfos()失败...因为手牌为nil", d.DlogDes(), user.GetUserId())
		return nil
	}

	jiaoInfos := []*mjproto.JiaoInfo{}

	//获取用户手牌 包括inPai
	userHandPai := NewMJHandPai()
	*userHandPai = *user.GetGameData().HandPai        //手牌
	userPais := make([]*MJPai, len(userHandPai.Pais)) //需要改变的牌
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
		jiaoInfo := NewJiaoInfo()

		//遍历麻将牌,看哪一张能胡牌
		for l := 0; l < len(mjpaiMap); l += 4 {

			//遍历未知牌
			//将遍历到的未知牌与用户手牌组合成handPai 去canhu
			mjPai := InitMjPaiByIndex(l)
			canHu, fan, _, _, _, _ := d.HuParser.GetCanHu(userHandPai, mjPai, true, 0, d.IsBanker(user))
			if canHu {
				//叫牌的信息
				mjPaiLeftCount := int32(d.GetLeftPaiCount(user, mjPai)) //该可胡牌在桌面中的剩余数量 注 对于自己而言的剩余
				jiaoPaiInfo := NewJiaoPaiInfo()
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

//获取用户听张的最大番数
func (d *MjDesk) GetJiaoMaxFan(u *MjUser) int32 {
	maxFan := 0

	//取得可以叫的infos
	jiaoInfos := d.GetJiaoInfos(u)

	if jiaoInfos == nil {
		return 0
	}
	for _, jiaoInfo := range jiaoInfos {
		paiInfos := jiaoInfo.PaiInfos
		if paiInfos != nil {
			for _, paiInfo := range paiInfos {
				fan := paiInfo.Fan
				if int32(maxFan) < *fan {
					//如果最大番小于当前可叫的番数
					maxFan = int(*fan)
				}
			}
		}
	}
	return int32(maxFan)
}

//用户没有叫的处理了
func (d *MjDesk) DoDaJiao(u *MjUser) {

	//user == nil 不查  直接返回
	if u == nil {
		//log
		return
	}

	//胡牌了 不用查大叫
	if u.IsHu() {
		return
	}

	//花猪不用查大叫
	if u.IsHuaZhu() {
		return
	}

	//有叫的时候不用查大叫
	if u.IsYouJiao() {
		return
	}

	//处理u 有杠，但是没有听的情况
	d.backGangMoney(u)
	//判断谁可以查u的大叫
	//没听牌的玩家(花猪除外)赔给听牌的玩家 按听牌的最大番型给
	//log.T("开始处理玩家[%v]没叫,开始处理查大叫...", u.GetUserId())

	for _, user := range d.GetUsers() {
		//如果是自己，不用管
		if user.GetUserId() == u.GetUserId() {
			continue
		}

		//没有叫的玩家不能去查别人
		if !user.IsYouJiao() {
			continue
		}

		if !d.IsXueLiuChengHe() && user.IsHu() {
			//如果不是血流成河且用户已胡 则不能查u的大叫
			continue
		}

		//获得听牌的最大番数.
		fan := d.GetJiaoMaxFan(user)
		score := d.GetBaseValue() * int64(math.Pow(2, float64(fan)))

		//如果looper不是被查大叫的玩家 且 该looper有听牌 且 该looper没有胡 为该玩家增加赢钱的bill
		user.AddBill(u.GetUserId(), MJUSER_BILL_TYPE_YING_DAJIAO, "用户查大叫，赢钱", score, nil, d.GetRoomType())
		user.AddStatisticsCountChaDaJiao(d.GetCurrPlayCount())

		u.AddBill(user.GetUserId(), MJUSER_BILL_TYPE_SHU_DAJIAO, "用户被查叫，输钱", -score, nil, d.GetRoomType())
		u.AddStatisticsCountBeiChaJiao(d.GetCurrPlayCount())
	}
}

//退换杠钱
func (d *MjDesk) backGangMoney(u *MjUser) error {
	if u == nil {
		return nil
	}

	if len(u.GetGameData().GetGangInfo()) <= 0 {
		return nil
	}

	//循环处理每一个杠
	for _, g := range u.GetGameData().GetGangInfo() {
		if g == nil {
			continue
		}

		log.T("%v 开始退换user[%v]杠[%v]收的钱：", d.DlogDes(), u.GetUserId(), g)
		//开始处理每一个杠，通过循环自己的账单，来退钱
		bills := u.GetBill().GetBills()
		for _, b := range bills {
			//首先确定是和这个杠info相关的账单
			if b == nil || b.GetPai().GetIndex() != g.GetPai().GetIndex() {
				continue
			}
			//确定关联账单的人
			ru := d.GetUserByUserId(b.GetOutUserId())
			log.T("%v 开始退换user[%v]杠收的钱[%v]给user[%v]：", d.DlogDes(), u.GetUserId(), b.GetAmount(), ru.GetUserId())
			ru.AddBill(u.GetUserId(), MJUSER_BILL_TYPE_YING_TUIGANGQIAN, "查大叫退杠钱", int64(math.Abs(float64(b.GetAmount()))), b.GetPai(), d.GetRoomType())   //加钱
			u.AddBill(ru.GetUserId(), MJUSER_BILL_TYPE_SHU_TUIGANGQIAN, "查大叫收退的杠钱", -int64(math.Abs(float64(b.GetAmount()))), b.GetPai(), d.GetRoomType()) //退钱
		}
	}
	return nil
}

//处理lottery的数据

//需要保存到 ..T_mj_desk_round   ...这里设计到保存数据，战绩相关的查询都要从这里查询
func (d *MjDesk) DoLottery() error {
	log.T("desk(%v),gameNumber(%v)处理DoLottery()", d.GetDeskId(), d.GetGameNumber())
	data := model.T_mj_desk_round{}
	data.DeskId = d.GetDeskId()
	data.GameNumber = d.GetGameNumber()
	data.BeginTime = timeUtils.String2YYYYMMDDHHMMSS(d.GetBeginTime())
	data.EndTime = time.Now()
	data.UserIds = d.GetUserIds()

	//一次处理每个胡牌的人
	for _, user := range d.GetUsers() {
		//这里不应该是胡牌的人才有记录...而是应该记录每一个人...
		if user != nil {
			//处理胡牌之后，分数相关的逻辑.
			//这里有一个统计...实在杠牌，或者胡牌之后会更新的数据...结算的时候，数据落地可以使用这个...
			//user.Statisc
			bean := model.MjRecordBean{}
			bean.UserId = user.GetUserId()
			bean.NickName = user.GetNickName()
			bean.WinAmount = user.Bill.GetWinAmount() //	赢了多少钱...

			//添加到record
			data.Records = append(data.Records, bean)
			//开奖之后，玩家的状态修改
			user.AfterLottery()
		}
	}

	log.T("%v 一局游戏结束 开始插入数据[%v]到mongo ", d.DlogDes(), data)
	//保存数据
	go func(d interface{}) {
		defer Error.ErrorRecovery("保存游戏数据到mgo")
		db.InsertMgoData(tableName.DBT_MJ_DESK_ROUND, d)
	}(&data)

	log.T("desk(%v),gameNumber(%v)处理DoLottery(),处理完毕,保存数据data[%v]", d.GetDeskId(), d.GetGameNumber(), data)

	return nil

}

//得到下一个摸牌的人...
func (d *MjDesk) GetNextMoPaiUser() *MjUser {

	//首先找，刚刚杠过牌的User，否则找下一个User
	for _, u := range d.GetUsers() {
		if u != nil && u.GetPreMoGangInfo() != nil {
			return u
		}
	}

	//log.T("查询下一个玩家...当前的activeUser[%v]", d.GetActiveUser())
	var activeUser *MjUser = nil
	activeIndex := -1
	for i, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == d.GetActiveUser() {
			activeIndex = i
			break
		}
	}
	//log.T("查询下一个玩家...当前的activeUser[%v],activeIndex[%v]", d.GetActiveUser(), activeIndex)
	if activeIndex == -1 {
		return nil
	}

	for i := activeIndex + 1; i < activeIndex+int(d.GetUserCount()); i++ {
		user := d.GetUsers()[i%int(d.GetUserCount())]
		//log.T("查询下一个玩家...当前的activeUser[%v],activeIndex[%v],循环检测index[%v],user.IsNotHu(%v),user.CanMoPai[%v]", d.GetActiveUser(), activeIndex, i, user.IsNotHu(), user.CanMoPai(d.IsXueLiuChengHe()))
		if user != nil && user.CanMoPai(d.IsXueLiuChengHe()) {
			activeUser = user
			break
		}
	}

	//找到下一个操作的user
	return activeUser

}

//得到下一张牌...
func (d *MjDesk) GetNextPai() *MJPai {
	*d.MJPaiCursor ++
	if d.GetMJPaiCursor() >= d.GetTotalMjPaiCount() {
		log.T("牌已经摸完了:要找的牌的坐标[%v]已经超过整副麻将的坐标了... ", d.GetMJPaiCursor())
		*d.MJPaiCursor --
		return nil
	} else {
		p := d.AllMJPai[d.GetMJPaiCursor()]
		pai := NewMjpai()
		*pai.Des = p.GetDes()
		*pai.Flower = p.GetFlower()
		*pai.Index = p.GetIndex()
		*pai.Value = p.GetValue()
		return pai
	}
}

func (d *MjDesk) SendMopaiOverTurn() error {
	if d.IsChangShaMaJiang() {
		return d.SendMopaiOverTurnChangSha()
	} else {
		return d.SendMopaiOverTurnChengDu()
	}
}

//长沙的摸牌
func (d *MjDesk) SendMopaiOverTurnChangSha() error {
	//首先判断是否可以lottery(),如果可以那么直接开奖
	if d.Time2Lottery() {
		d.Lottery() //摸牌的时候判断可以lottery了
		return nil
	}

	//开始摸牌的逻辑
	user := d.GetNextMoPaiUser()
	if user == nil {
		log.E("服务器出现错误..没有找到下一个摸牌的玩家...")
		return errors.New("没有找到下一家")
	}
	d.SetAATUser(user.GetUserId(), MJDESK_ACT_TYPE_MOPAI) //用户摸牌之后，设置前端指针指向的玩家//长沙麻将,用户摸牌之后，设置当前活动的玩家

	//这里需要判断特殊情况
	if d.IsChangShaMaJiang() && user.GetPreMoGangInfo() != nil && !user.GetPreMoGangInfo().GetBu() {
		log.T("%v 玩家[%v]长沙杠之后开始处理摸两张牌", d.DlogDes(), user.GetUserId())
		//如果又是长沙麻将，又是杠那么需要摸两张牌,否则就是普通麻将 摸一张
		//杠牌之后一次性摸两张牌
		inpai1 := d.GetNextPai()
		inpai2 := d.GetNextPai()

		user.GameData.HandPai.InPai = inpai1
		overTrun1 := d.GetMoPaiOverTurn(user, false) //长沙玩家杠牌之后，用户摸牌的时候,发送一个用户摸牌的overturn
		overTrun1.CanGang = proto.Bool(false)
		overTrun1.GangCards = nil
		overTrun1.CanBu = proto.Bool(false)
		overTrun1.BuCards = nil

		user.GameData.HandPai.InPai2 = inpai1
		user.GameData.HandPai.InPai = inpai2
		overTrun2 := d.GetMoPaiOverTurn(user, false) //长沙玩家杠牌之后，用户摸牌的时候,发送一个用户摸牌的overturn
		overTrun2.CanGang = proto.Bool(false)
		overTrun2.GangCards = nil
		overTrun2.CanBu = proto.Bool(false)
		overTrun2.BuCards = nil

		/**
			注意：杠牌和补牌的差别...如果是杠牌，摸两张，然后自动打出去...如果补，那么和普通麻将的杠是一样的...
		 */
		ack := &mjproto.Game_ChangShaOverTurnAfterGang{}
		ack.GangPai = append(ack.GangPai, user.GetGameData().GetHandPai().GetInPai().GetCardInfo(), user.GetGameData().GetHandPai().GetInPai2().GetCardInfo())
		ack.Header = newProto.NewHeader()
		ack.Header.UserId = proto.Uint32(user.GetUserId())
		if overTrun1.GetCanHu() {
			ack.CanHu = proto.Bool(true)
			ack.CanGuo = proto.Bool(true)
			ack.HuCards = append(ack.HuCards, overTrun1.ActCard)
		}

		if overTrun2.GetCanHu() {
			ack.CanHu = proto.Bool(true)
			ack.CanGuo = proto.Bool(true)
			ack.HuCards = append(ack.HuCards, overTrun2.ActCard)
		}

		log.T("[%v][%v]开摸的牌【%v】...", d.DlogDes(), user.GetUserId(), user.UserPai2String(), ack)
		d.BroadCastProto(ack)
		//杠起来的两张牌，如果可以胡牌，需要玩家选择是否胡牌，如果不能胡牌，直接系统把这两张牌打出去
		if !ack.GetCanHu() {
			//不能胡牌，直接把两张牌打出去,第二个参数传0
			go d.ActOutChangSha(user.GetUserId(), 0, true)
		} else {
			d.SetAATUser(user.GetUserId(), MJDESK_ACT_TYPE_WAIT_CHECK_CHANGSHAGANG) //长沙麻将,用户摸牌两张牌之后，设置当前活动的玩家
		}

	} else {
		//发送给当事人时候的信息
		user.GameData.HandPai.InPai = d.GetNextPai()
		//判断是否是海底牌
		if d.IsChangShaMaJiang() && user.GameData.HandPai.InPai.GetIndex() == d.GetLastMjPai().GetIndex() {
			if user.GetNeedHaidiStatus() == MJUSER_NEEDHAIDI_STATUS_DEFAULT {
				//询问是否需要海底牌
				d.enquireHaiDi(user)
				return nil
			}
		}

		//长沙麻将杠后摸牌，自动打牌
		overTrun := d.GetMoPaiOverTurn(user, false) //普通摸牌，用户摸牌的时候,发送一个用户摸牌的overturn
		if d.IsChangShaMaJiang() && user.GetChangShaGangStatus() {
			overTrun.CanGang = proto.Bool(false) //设置不能杠
			overTrun.CanPeng = proto.Bool(false) //设置不能碰
			overTrun.CanBu = proto.Bool(false)   //设置不能补
			overTrun.CanGuo = overTrun.CanHu
			overTrun.BuCards = nil
			overTrun.GangCards = nil
		}

		//普通打牌
		log.T("[%v][%v]长沙麻将的牌【%v】...", d.DlogDes(), user.GetUserId(), user.UserPai2String(), overTrun)
		user.SendOverTurn(overTrun) //玩家摸排之后发送overturn

		//给其他人广播协议
		bcOverTurn := &mjproto.Game_OverTurn{}
		util.DeepCopy(bcOverTurn, overTrun)
		bcOverTurn.CanHu = proto.Bool(false)
		bcOverTurn.ActCard = NewBackPai()
		d.BroadCastProtoExclusive(bcOverTurn, user.GetUserId())
	}

	return nil
}

//发送摸牌的广播
//指定一个摸牌，如果没有指定，则系统通过游标来判断
func (d *MjDesk) SendMopaiOverTurnChengDu() error {
	//首先判断是否可以lottery(),如果可以那么直接开奖
	if d.Time2Lottery() {
		d.Lottery() //摸牌的时候判断可以lottery了
		return nil
	}

	//开始摸牌的逻辑
	user := d.GetNextMoPaiUser()
	if user == nil {
		log.E("服务器出现错误..没有找到下一个摸牌的玩家...")
		return errors.New("没有找到下一家")
	}
	d.SetAATUser(user.GetUserId(), MJDESK_ACT_TYPE_MOPAI) //长度麻将 用户摸牌之后，设置当前活动的玩家

	//发送摸牌的OverTrun
	user.GameData.HandPai.InPai = d.GetNextPai()
	overTrun := d.GetMoPaiOverTurn(user, false) //普通摸牌，用户摸牌的时候,发送一个用户摸牌的overturn
	user.SendOverTurn(overTrun)                 //玩家摸排之后发送overturn
	log.T("[%v][%v]开摸的牌【%v】...", d.DlogDes(), user.GetUserId(), user.UserPai2String(), overTrun)
	//给其他人广播协议
	*overTrun.CanHu = false
	*overTrun.CanGang = false
	overTrun.ActCard = NewBackPai()
	d.BroadCastProtoExclusive(overTrun, user.GetUserId())

	return nil
}

//设置用户的状态为离线
func (d *MjDesk) SetOfflineStatus(userId uint32) {
	log.T("玩家[%v]断开连接，设置当前状态为离线的状态...", userId)
	user := d.GetUserByUserId(userId)
	*user.IsBreak = true

	//离线之后，给其他的玩家发送离线的广播
	bc := &ddproto.CommonBcUserBreak{
		UserId: proto.Uint32(userId),
	}
	d.BroadCastProto(bc)
}

//设置用户的状态为离线
func (d *MjDesk) SetReconnectStatus(userId uint32, a gate.Agent) {
	log.T("玩家[%v]重新链接，设置当前状态为在线的状态的状态...", userId)
	user := d.GetUserByUserId(userId)
	*user.IsBreak = false
	user.Agent = a

	//离线之后，给其他的玩家发送离线的广播
	bc := &ddproto.CommonAckReconnect{
		UserId: proto.Uint32(userId),
	}
	//
	d.BroadCastProto(bc)
}

//剩余牌的数量
func (d *MjDesk) GetRemainPaiCount() int32 {
	remainCount := d.GetTotalMjPaiCount() - d.GetMJPaiCursor() - 1
	//log.T("GetRemainPaiCount[%v], GetTotalMjPaiCount[%v], GetMJPaiCursor[%v]", remainCount, d.GetTotalMjPaiCount(), d.GetMJPaiCursor())
	return remainCount
}

//可以把overturn放在一个地方,目前都是摸牌的时候在用
func (d *MjDesk) GetMoPaiOverTurn(user *MjUser, isOpen bool) *mjproto.Game_OverTurn {

	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = user.GetUserId()         //这个是摸牌的，所以是广播...
	*overTurn.PaiCount = d.GetRemainPaiCount()  //桌子剩余多少牌
	*overTurn.ActType = OVER_TURN_ACTTYPE_MOPAI //摸牌
	*overTurn.Time = 30
	if isOpen {
		overTurn.ActCard = user.GameData.HandPai.InPai.GetBackPai()
	} else {
		overTurn.ActCard = user.GameData.HandPai.InPai.GetCardInfo()
	}

	log.T("[%v]开始摸牌:%v", d.DlogDes(), user.UserPai2String())
	*overTurn.CanHu, _, _, _, _, _ = d.HuParser.GetCanHu(user.GetGameData().GetHandPai(), user.GetGameData().GetHandPai().GetInPai(), true, 0, d.IsBanker(user)) //是否可以胡牌
	*overTurn.CanPeng = false                                                                                                                                    //是否可以碰牌

	//处理杠牌的时候
	/**
		1，血战到底：用户胡牌之后是不会进入到这个方法的
		2，血流成河：用户已经胡牌，那么杠牌之后，胡牌不会改变的情况下，才可以杠 // todo
	 */
	canGangBool, gangPais := user.GameData.HandPai.GetCanGang(nil, d.GetRemainPaiCount()) //是否可以杠牌
	log.T("%v 判断玩家的牌受否可以杠%v,handpai:%v,inpai%v", d.DlogDes(), canGangBool, user.GetGameData().GetHandPai().GetPais(), user.GetGameData().GetHandPai().GetInPai())
	*overTurn.CanGang = canGangBool
	if canGangBool && gangPais != nil {
		if user.IsHu() && d.IsXueLiuChengHe() {
			//血流成河，胡牌之后 杠牌的逻辑
			jiaoPais := d.HuParser.GetJiaoPais(user.GetGameData().GetHandPai()) //GetMoPaiOverTurn
			for _, g := range gangPais {
				//判断杠牌之后的叫牌是否和杠牌之前一样
				if user.AfterGangEqualJiaoPai(jiaoPais, g) {
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
	overTurn.JiaoInfos = d.GetJiaoInfos(user)

	//这里需要对长沙麻将做特殊处理(主要是杠，补的处理)
	if d.IsChangShaMaJiang() {
		if overTurn.GetCanGang() {
			overTurn.CanBu = proto.Bool(true)
			overTurn.CanGang = proto.Bool(false)
			overTurn.BuCards = overTurn.GangCards
			overTurn.GangCards = nil
			//判断长沙麻将能不能杠
			for _, g := range overTurn.BuCards {
				cang := user.GetCanChangShaGang(InitMjPaiByIndex(int(g.GetId()))) // 摸排的时候 判断能否gang
				log.T("判断玩家[%v]对牌[%v]是否可以长沙杠[%v]", user.GetUserId(), g.GetId(), cang)
				if cang {
					overTurn.CanGang = proto.Bool(true)
					overTurn.GangCards = append(overTurn.GangCards, g)
				}
			}
		}
	}
	return overTurn
}

//通过checkCase 得到一个OverTurn
func (d *MjDesk) GetOverTurnByCaseBean(checkPai *MJPai, caseBean *CheckBean, actType int32) *mjproto.Game_OverTurn {
	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = caseBean.GetUserId()
	*overTurn.CanGang = caseBean.GetCanGang()
	for _, p := range caseBean.GetGangCards() {
		overTurn.GangCards = append(overTurn.GangCards, p.GetCardInfo())
	}
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

	//这里需要对长沙麻将做特殊处理(主要是杠，补的处理)
	if d.IsChangShaMaJiang() && overTurn.GetCanGang() {
		user := d.GetUserByUserId(caseBean.GetUserId()) //判断玩家是否可以长沙杠
		overTurn.CanBu = proto.Bool(true)
		overTurn.CanGang = proto.Bool(false)
		overTurn.BuCards = overTurn.GangCards
		overTurn.GangCards = nil
		//判断长沙麻将能不能杠
		for _, g := range overTurn.BuCards {
			cang := user.GetCanChangShaGang(InitMjPaiByIndex(int(g.GetId()))) //overTurn by caseBean
			if cang {
				overTurn.CanGang = proto.Bool(true)
				overTurn.GangCards = append(overTurn.GangCards, g)
			}

		}
	}

	return overTurn
}
