package majiang

import (
	"casino_common/common/log"
	"errors"
	"casino_common/common/consts/tableName"
	mjproto        "casino_mj_changsha/msg/protogo"
	"casino_common/utils/db"
	"casino_common/utils/timeUtils"
	"time"
	"casino_mj_changsha/msg/funcsInit"
	"casino_common/common/userService"
	"github.com/golang/protobuf/proto"
)

//是否可以开始
func (d *MjDesk) time2begin() error {
	log.T("%v检测游戏是否可以开始,IsPlayerEnough:%v ,IsAllReady:%v,IsNotDingQue:%v ",
		d.DlogDes(), d.IsPlayerEnough(), d.IsAllReady(), d.IsNotDingQue())
	if d.IsAllReady() && d.IsPlayerEnough() && d.IsNotDingQue() {
		return nil
	} else {
		return errors.New("开始游戏失败，因为还有人没有准备")
	}
	return nil
}

//这里是开始麻将的路由
func (d *MjDesk) begin() error {
	err := d.time2begin()
	if err != nil {
		log.T("无法开始游戏:err[%v]", err)
		return err
	}
	return d.beginChangSha()
}

//扣除房卡
func (d *MjDesk) deductCreateFeeRoomCard() error {
	//计算朋友桌的房费
	if d.IsFriend() {
		log.T("%v 玩家[%v]开始扣除房卡，消耗房卡:%v", d.DlogDes(), d.GetOwner(), d.GetCreateFee())
		remain, err := userService.DECRUserRoomcard(d.GetOwner(), d.GetCreateFee()) //减少房主的房卡
		if err != nil {
			//扣费失败，创建房间失败
			log.E("玩家[%v]创建房间的时候扣房卡[%v]失败，创建房间失败...", d.GetOwner(), d.GetCreateFee())
			return ERR_ROOMCARD_INSUFFICIENT
		}
		log.T("%v 玩家[%v]扣除房卡成功，消耗房卡:%v,剩余房卡:%v", d.DlogDes(), d.GetOwner(), d.GetCreateFee(), remain)
	}

	return nil
}

//解散房间的时候，退还房卡
func (d *MjDesk) sendBackRoomCrad() error {
	if d.IsFriend() && d.GetCurrPlayCount() == 0 {
		remain, err := userService.INCRUserRoomcard(d.GetOwner(), d.GetCreateFee()) //退换房卡
		if err != nil {
			log.E("玩家[%v]剑三房间的时候退还房卡[%v]失败...remain[%v]", d.GetOwner(), d.GetCreateFee(), remain)
			return ERR_ROOMCARD_INSUFFICIENT
		}
	}
	return nil
}

/**
1,初始化desk
2,初始化user
 */
func (d *MjDesk) beginInit() error {
	log.T("%vbeginIint()...nextBankUser:%v", d.DlogDes(), d.GetNextBanker())

	//初始化桌子的信息
	d.InitBanker()                                                      //初始化庄
	d.AddCurrPlayCount()                                                //场次数目加一
	*d.GameNumber, _ = db.GetNextSeq(tableName.DBT_T_TH_GAMENUMBER_SEQ) //设置游戏编号
	*d.BeginTime = timeUtils.Format(time.Now())
	d.BirdInfo = []*mjproto.BirdInfo{}  //清空鸟牌信息
	d.CheckCase = nil                   ///初始化为nil,（不然第二句开始，如果庄起手能杠，点击杠的时候就会出错...）
	d.SetStatus(MJDESK_STATUS_OPENNING) //游戏中
	//初始化每个玩家的信息
	for _, user := range d.GetUsers() {
		if user != nil && user.CanBegin() {
			user.BeginInit(d.GetCurrPlayCount(), d.GetBanker())
		}
	}
	//发送游戏开始的协议...
	d.SendNewGame_Opening()
	time.Sleep(SHAIZI_SLEEP_DURATION)
	return nil
}

//初始化庄
func (d *MjDesk) InitBanker() {
	//1,初始化庄的信息,如果目前没有庄，则设置房主为庄,如果有庄，则不用管，每局游戏借宿的时候，会设置下一局的庄
	if d.GetNextBanker() != 0 {
		d.SetBanker(d.GetNextBanker())
	}
	//金币场可能会出现没有banker的情况
	if d.GetBankerUser() == nil {
		d.SetBanker(d.GetUsers()[0].GetUserId())
	}

	d.NextBanker = proto.Uint32(0) //设置下一次的为0
}

//发送game_opening 的协议
func (d *MjDesk) SendNewGame_Opening() {
	log.T("%v发送游戏开始的协议.当前桌子共[%v]把，现在是第[%v]把游戏开始", d.DlogDes())
	open := newProto.NewGame_Opening()
	*open.CurrPlayCount = d.GetCurrPlayCount()
	*open.Dice1 = d.GetDice1()                //骰子
	*open.Dice2 = d.GetDice2()                //骰子
	open.UserCoinBeans = d.GetUserCoinBeans() //玩家的金币信息
	d.BroadCastProto(open)
}

/**
	初始化牌相关的信息
 */
func (d *MjDesk) initCards() error {
	log.T("%v 开始initCards()...", d.DlogDes())
	d.AllMJPai = XiPai()
	//d.AllMJPai = XiPaiTestHu()
	if d.IsLiangFang() {
		//如果是两房 过滤掉万牌
		d.AllMJPai = IgnoreFlower(d.AllMJPai, W)
	}

	//cardsNum 发牌张数处理
	cardsNum := int(d.GetCardsNum())
	if cardsNum == 0 {
		cardsNum = 13 //默认13张
	}
	for i, u := range d.Users {
		if u != nil && u.IsReady() {
			//log.T("开始给你玩家[%v]初始化手牌...", u.GetUserId())
			ps := make([]*MJPai, cardsNum)
			copy(ps, d.AllMJPai[i*cardsNum: (i+1)*cardsNum]) //这里这样做的目的是不能更改base的值
			u.GameData.HandPai.Pais = ps
			*d.MJPaiCursor = int32((i+1)*cardsNum) - 1
		}
	}

	//庄需要多发一张牌
	bankUser := d.GetBankerUser()
	if bankUser == nil {
		return errors.New("系统错误")
	}
	bankUser.GameData.HandPai.InPai = d.GetNextPai()

	//发牌的协议game_DealCards  初始化完成之后，给每个人发送牌
	for _, user := range d.Users {
		if user != nil {
			dealCards := d.GetDealCards(user)
			if dealCards != nil {
				log.T("开始发送玩家[%v]的手牌[%v]...", user.GetUserId(), dealCards)
				user.WriteMsg(dealCards)
			} else {
				log.E("给user[%v]发牌的时候出现错误..", user.GetUserId())
			}
		}
	}

	//发送发牌的广播
	return nil
}

//发牌的协议
func (d *MjDesk) GetDealCards(user *MjUser) *mjproto.Game_DealCards {
	dealCards := newProto.NewGame_DealCards()
	*dealCards.DealerUserId = d.GetBanker()
	for _, u := range d.GetUsers() {
		if u != nil {
			if u.GetUserId() == user.GetUserId() {
				//表示是自己，可以看到手牌
				pc := u.GetPlayerCard(true)
				dealCards.PlayerCard = append(dealCards.PlayerCard, pc)
			} else {
				pc := u.GetPlayerCard(false) //表示不是自己，不能看到手牌

				if d.GetBanker() == u.GetUserId() {
					pc.HandCard = append(pc.HandCard, NewBackPai())
				}

				dealCards.PlayerCard = append(dealCards.PlayerCard, pc)
			}

		}

	}

	return dealCards
}

//游戏正式开始，庄家打牌的广播
func (d *MjDesk) BeginStart() error {
	log.T("%v beginStart()...", d.DlogDes())

	//设置游戏开始的状态
	d.SetStatus(MJDESK_STATUS_RUNNING)
	d.UpdateUserStatus(MJUSER_STATUS_GAMING)                      //设置为游戏中的状态
	d.SetAATUser(d.GetBanker(), MJDESK_ACT_TYPE_BANK_FIRST_MOPAI) //庄开始打牌

	//通知庄家打一张牌,这里初始化信息，这里应该是广播的..
	//注意是否可以碰，可以杠牌，可以胡牌，只有当时人才能看到，所以广播的和当事人的收到的数据不一样...
	bankUser := d.GetBankerUser()

	overTurn := d.GetMoPaiOverTurn(bankUser, true) //定缺完了之后，庄摸牌

	bankUser.SendOverTurn(overTurn) //庄家打牌

	//广播时候的信息
	overTurn.ActCard = nil
	*overTurn.CanHu = false
	*overTurn.CanGang = false
	*overTurn.CanPeng = false
	d.BroadCastProtoExclusive(overTurn, d.GetBanker())

	return nil
}
