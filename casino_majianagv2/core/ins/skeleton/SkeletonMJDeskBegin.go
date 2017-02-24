package skeleton

import (
	"casino_common/common/consts/tableName"
	"casino_common/common/log"
	"time"
	"casino_majiang/service/majiang"
	"errors"
	"casino_majiang/msg/funcsInit"
	"casino_common/utils/db"
	"casino_common/utils/timeUtils"
	"casino_majianagv2/core/majiangv2"
	"casino_majiang/msg/protogo"
	"casino_majianagv2/core/api"
)

//是否可以开始
func (d *SkeletonMJDesk) Time2begin() error {
	log.T("%v检测游戏是否可以开始,IsPlayerEnough:%v ", d.DlogDes(), d.IsPlayerEnough())
	log.T("%v检测游戏是否可以开始,IsAllReady:%v ", d.DlogDes(), d.IsAllReady())
	log.T("%v检测游戏是否可以开始,IsNotDingQue:%v ", d.DlogDes(), d.IsNotDingQue())
	if d.IsAllReady() && d.IsPlayerEnough() && d.IsNotDingQue() {
		return nil
	} else {
		return errors.New("开始游戏失败，因为还有人没有准备")
	}
	return nil
}

func (d *SkeletonMJDesk) BeginInit() error {
	log.T("desk[%v]round[%v]beginIint()...", d.GetMJConfig().DeskId, d.GetMJConfig().CurrPlayCount)

	//初始化桌子的信息
	//1,初始化庄的信息,如果目前没有庄，则设置房主为庄,如果有庄，则不用管，每局游戏借宿的时候，会设置下一局的庄
	if d.GetMJConfig().NextBanker != 0 {
		d.GetMJConfig().Banker = d.GetMJConfig().NextBanker
	}
	//金币场可能会出现没有banker的情况
	if d.GetBankerUser() == nil {
		d.GetMJConfig().Banker = d.GetUsers()[0].GetUserId()
	}
	//增加局数的统计
	d.GetMJConfig().CurrPlayCount ++
	//func (d *MjDesk) AddCurrPlayCount() { atomic.AddInt32(d.CurrPlayCount, 1) }
	//d.AddCurrPlayCount()                                                //场次数目加一
	d.GetMJConfig().GameNumber, _ = db.GetNextSeq(tableName.DBT_T_TH_GAMENUMBER_SEQ) //设置游戏编号
	d.GetMJConfig().BeginTime = timeUtils.Format(time.Now())
	d.GetMJConfig().NextBanker = 0 //设置下一次的为0

	//初始化每个玩家的信息
	for _, user := range d.GetUsers() {
		if user != nil && user.GetStatus().IsReady() {
			user.BeginInit(d.GetMJConfig().CurrPlayCount, d.GetMJConfig().Banker)
		}
	}
	//发送游戏开始的协议...
	d.SendNewGame_Opening()
	time.Sleep(majiang.SHAIZI_SLEEP_DURATION)
	return nil
}

/**
	初始化牌相关的信息
 */
func (d *SkeletonMJDesk) InitCards() error {
	log.T("%vinitCards()...", d.DlogDes())
	d.AllMJPais = majiangv2.XiPai() //真累需要穿参数 房.
	//d.AllMJPai = XiPaiTestHu()
	if d.GetMJConfig().FangCount == 2 {
		//如果是两房 过滤掉万牌
		d.AllMJPais = majiangv2.IgnoreFlower(d.AllMJPais, majiangv2.W)
	}

	//cardsNum 发牌张数处理
	cardsNum := int(d.GetMJConfig().CardsNum)
	if cardsNum == 0 {
		cardsNum = 13 //默认13张
	}
	for i, u := range d.Users {
		if u != nil && u.GetStatus().IsReady() {
			//log.T("开始给你玩家[%v]初始化手牌...", u.GetUserId())
			ps := make([]*majiang.MJPai, cardsNum)
			copy(ps, d.AllMJPais[i*cardsNum: (i+1)*cardsNum]) //这里这样做的目的是不能更改base的值
			u.GetGameData().HandPai.Pais = ps
			d.GetMJConfig().MJPaiCursor = int32((i+1)*cardsNum) - 1;
		}
	}

	//庄需要多发一张牌
	bankUser := d.GetBankerUser()
	if bankUser == nil {
		return errors.New("系统错误")
	}
	bankUser.GetGameData().HandPai.InPai = d.GetNextPai()

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
func (d *SkeletonMJDesk) GetDealCards(user api.MjUser) *mjproto.Game_DealCards {
	dealCards := newProto.NewGame_DealCards()
	*dealCards.DealerUserId = d.GetMJConfig().Banker
	for _, u := range d.GetSkeletonMJUsers() {
		if u != nil {
			if u.GetUserId() == user.GetUserId() {
				//表示是自己，可以看到手牌
				pc := u.GetPlayerCard(true)
				if d.GetMJConfig().Banker == u.GetUserId() {
					pc.HandCard = append(pc.HandCard, u.GameData.HandPai.InPai.GetCardInfo())
				}

				dealCards.PlayerCard = append(dealCards.PlayerCard, pc)
			} else {
				pc := u.GetPlayerCard(false) //表示不是自己，不能看到手牌

				if d.GetMJConfig().Banker == u.GetUserId() {
					pc.HandCard = append(pc.HandCard, majiang.NewBackPai())
				}
				dealCards.PlayerCard = append(dealCards.PlayerCard, pc)
			}

		}

	}
	return dealCards
}

//开始换三张
func (d *SkeletonMJDesk) BeginExchange() error {
	log.T("desk[%v]round[%v]beginExchange()...", d.GetMJConfig().DeskId, d.GetMJConfig().CurrPlayCount)

	d.GetStatus().SetStatus(majiang.MJDESK_STATUS_EXCHANGE)
	data := newProto.NewGame_BroadcastBeginExchange()
	*data.Reconnect = false
	d.BroadCastProto(data)
	return nil
}

//开始定缺
func (d *SkeletonMJDesk) BeginDingQue() error {
	log.T("desk[%v]round[%v]beginDingQue()...", d.GetMJConfig().DeskId, d.GetMJConfig().CurrPlayCount)

	//开始定缺，修改desk的状态
	d.GetStatus().SetStatus(majiang.MJDESK_STATUS_DINGQUE)

	//给每个人发送开始定缺的信息
	beginQue := newProto.NewGame_BroadcastBeginDingQue()
	*beginQue.Reconnect = false
	//
	d.BroadCastProto(beginQue)
	return nil
}

//游戏正式开始，庄家打牌的广播
func (d *SkeletonMJDesk) BeginStart() error {
	log.T("desk[%v]round[%v]beginStart()...", d.GetMJConfig().DeskId, d.GetMJConfig().CurrPlayCount)

	//设置游戏开始的状态
	d.GetStatus().SetStatus(majiang.MJDESK_STATUS_RUNNING)
	d.UpdateUserStatus(majiang.MJUSER_STATUS_GAMING)
	d.SetActiveUser(d.GetMJConfig().Banker) // d.beginStart
	d.SetActUserAndType(d.GetMJConfig().Banker, majiang.MJDESK_ACT_TYPE_MOPAI)

	//通知庄家打一张牌,这里初始化信息，这里应该是广播的..
	//注意是否可以碰，可以杠牌，可以胡牌，只有当时人才能看到，所以广播的和当事人的收到的数据不一样...
	bankUser := d.GetBankerUser()

	overTurn := d.GetMoPaiOverTurn(bankUser, true) //定缺完了之后，庄摸牌
	bankUser.SendOverTurn(overTurn)                //庄家打牌
	//广播时候的信息
	overTurn.ActCard = nil
	*overTurn.CanHu = false
	*overTurn.CanGang = false
	*overTurn.CanPeng = false
	d.BroadCastProtoExclusive(overTurn, d.GetMJConfig().Banker)
	return nil
}
