package doudizhu

import (
	"github.com/name5566/leaf/gate"
	"casino_doudizhu/msg/funcsInit"
	"casino_doudizhu/msg/protogo"
	"time"
	"casino_common/common/log"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_common/utils/rand"
)

//这里主要存放 玩斗地主的一些多逻辑....其他的基本方法都放在DdzDesk中
/**
	1，通用的地主逻辑
	2，主要的地主逻辑
	3欢乐，四川等逻辑经分开在两个文件
 */

var (
	DEALCARDS_DRUATION time.Duration = time.Second * 3        //发牌之后延时3秒
)

func (d *DdzDesk) EnterUser(userId uint32, a gate.Agent) (error, int32) {
	log.T("玩家[%v]开始进入房间[%v]", userId, d.GetDeskId())
	//判断是否是重复进入
	olduser := d.GetUserByUserId(userId)
	if olduser != nil {
		olduser.setAgent(a)
		//这里需要判断是否是短线重连
		olduser.SetOnline()
		olduser.UpdateSession()       //更新session 信息，这里可以更具需求来保存对应的属性...
		ret := newProto.NewGame_AckEnterRoom()
		a.WriteMsg(ret)
		return nil, 1
	}

	//新进入
	errAddNew := d.AddUser(userId, a)
	if errAddNew != nil {
		log.T("玩家[%v]进入房间[%v]失败", userId, d.GetDeskId())
		//进入失败 //返回失败的信息
		ret := newProto.NewGame_AckEnterRoom()
		*ret.Header.Code = consts.ACK_RESULT_ERROR
		*ret.Header.Error = "进入房间失败"
		a.WriteMsg(ret)
		return Error.NewError(-1, "进入房间失败"), 0
	} else {
		//进入成功,返回进入成功的信息
		log.T("玩家[%v]进入房间[%v]成功", userId, d.GetDeskId())
		ret := newProto.NewGame_AckEnterRoom()
		a.WriteMsg(ret)
		return nil, 0
	}

	return nil, 0
}

//开始准备
func (d *DdzDesk) Ready(userId uint32) error {
	d.Lock()
	defer d.Unlock()

	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("玩家[%v]准备游戏的时候失败，deks[%v]中没有找到对应的玩家[%v]", userId, d.GetDeskId(), userId)
		return Error.NewError(-1, "没有找到对应的玩家")
	}

	//设置状态为准备的状态...
	user.SetStatus(DDZUSER_STATUS_READY)
	ack := newProto.NewDdzAckReady()
	*ack.UserId = userId

	//广播ack
	d.BroadCastProto(ack)

	//准备之后的处理
	d.AfterReady()

	return nil

}


//都准备好了之后，开始游戏
func (d *DdzDesk) AfterReady() error {
	//如果都准备好了，那么可以开始抢地主或者发牌了..
	if d.IsAllReady() && d.IsEnoughUser() {
		d.Begin()
	}
	return nil
}

//开始游戏
func (d *DdzDesk) Begin() error {

	//common初始化
	err := d.BeginCommon()
	if err != nil {
		log.E("初始化桌子的时候出错...")
		return err
	}

	//欢乐斗地主
	if d.IsHuanLeDoudDiZhu() {
		return d.HLBegin()
	}

	return nil
}

func (d *DdzDesk) BeginCommon() error {

	//通用的一些初始化方法
	log.T("d.CommonBeginInit()")
	d.CommonBeginInit()

	log.T("d.CommonInitCards()")
	//给玩家发牌
	d.CommonInitCards()
	time.Sleep(DEALCARDS_DRUATION)

	return nil
}

func (d *DdzDesk) IsTime2begin() error {
	return nil
}


//开始抢地主的逻辑 目前只有欢乐斗地主才会使用
func (d *DdzDesk) sendQiangDiZhuOverTurn(userId uint32) error {
	d.SetActiveUser(userId)
	//发送开始抢地主的广播...
	overTurn := newProto.NewDdzOverTurn()
	*overTurn.ActType = ddzproto.ActType_T_ROB_DIZHU
	*overTurn.UserId = userId
	//开发发送广播
	d.BroadCastProto(overTurn)
	return nil
}


// 开始叫地主
func (d *DdzDesk) sendJiaoDiZhuOverTurn(userId uint32) error {
	d.SetActiveUser(userId)        //设置为当前玩家

	//开始叫地主
	overTurn := newProto.NewDdzOverTurn()
	*overTurn.ActType = ddzproto.ActType_T_JIAO_DIZHU      //类型是开始叫地主
	*overTurn.UserId = userId
	*overTurn.CanOutCards = false
	*overTurn.CanDouble = false

	//广播开始叫地主的协议
	d.BroadCastProto(overTurn)
	return nil
}

//发送加倍
func (d *DdzDesk) sendJiaBeiOverTurn(userId uint32) error {
	d.SetActiveUser(userId)        //设置当前活动的玩家
	overTurn := newProto.NewDdzOverTurn()
	*overTurn.ActType = ddzproto.ActType_T_DOUBLE        ///加倍不加倍
	*overTurn.UserId = userId
	d.BroadCastProto(overTurn)
	return nil
}

//一场结束
func (d *DdzDesk) Lottery(user *DdzUser) {
	//开始结算....
	//1,计算炸弹的个数，计算分数,这里需要判断user的身份是地主还是平民
	if d.IsDiZhuRole(user) {
		//地主赢了,增加账单
		for _, loseUser := range d.Users {
			if loseUser.GetUserId() != user.GetUserId() {
				user.AddNewBill(d.GetWinValue(), user.GetUserId(), loseUser.GetUserId(), "地主赢了")
				loseUser.AddNewBill(-d.GetWinValue(), user.GetUserId(), loseUser.GetUserId(), "平明输了")
			}
		}

	} else {
		//地主输了,增加账单
		dizhuUser := d.GetUserByUserId(d.GetDiZhuUserId())
		for _, winUser := range d.Users {
			if winUser.GetUserId() != dizhuUser.GetUserId() {
				user.AddNewBill(d.GetWinValue(), user.GetUserId(), winUser.GetUserId(), "平明赢了")
				dizhuUser.AddNewBill(-d.GetWinValue(), user.GetUserId(), winUser.GetUserId(), "地主输了")
			}
		}
	}

	//2,发送结算的通知


}

//牌局结束
func (d *DdzDesk) DoEnd() {

}

func (d *DdzDesk) CommonBeginInit() error {
	//desk.init
	*d.DizhuPaiUser = d.Users[rand.Rand(0, 3)].GetUserId()        //随机一个 第一个叫地主的人
	d.SetActiveUser(0)        //设置当前激活的玩家为0

	//userInit
	for _, user := range d.Users {
		if user != nil {
			//初始化每个用户
			user.beginInit()
		}
	}

	return nil
}

//初始化每个人的牌
func (d *DdzDesk) CommonInitCards() {
	//获得一副洗好的牌
	d.AllPokerPai = XiPai()                //获得洗好的一副扑克牌

	//为每个人分配牌
	for i, user := range d.Users {
		user.GameData.HandPokers = make([]*PPokerPai, 17)        //这里的17暂时写死...
		copy(user.GameData.HandPokers, d.AllPokerPai[i * 17:(i + 1) * 17])
	}

	//底牌
	d.DiPokerPai = make([]*PPokerPai, 3)
	copy(d.DiPokerPai, d.AllPokerPai[54 - 3:54])

	//开始发牌
	d.EveryUserDoSomething(func(user *DdzUser) error {
		//开始发牌
		cardsInfo := newProto.NewDdzDealCards()
		cardsInfo.PlayerPokers = user.GetPlayerPokers()
		user.WriteMsg(cardsInfo)        //给用户发送牌的信息
		return nil
	})
}

//判断出牌的用户是否合法
func (d *DdzDesk) CheckOutUser(userId uint32) error {
	if d.GetActiveUserId() == userId {
		return nil
	} else {
		return Error.NewError(-1, "activeUser 不正确...")
	}
}

//打牌
func (d *DdzDesk) ActOut(userId uint32, out *POutPokerPais) error {

	err := d.CheckOutUser(userId)
	if err != nil {

	}

	//判断牌是否合法
	err = d.CheckOutPai(out)
	if err != nil {
		log.E("玩家[%v]出牌[%v]失败,desk.outpai[%v]", userId, out, d.OutPai)
		return Error.NewError(-1, "出牌失败，牌型有误")
	}

	user := d.GetUserByUserId(userId)
	if user == nil {
		return Error.NewError(-1, "没有找到玩家,出牌失败")
	}

	//牌型合法,1保存用户出的牌，2删除手里面的牌
	err = user.DOPoutPokerPais(out)
	if err != nil {
		log.E("玩家[%v]出牌的时候错误", userId)
		return Error.NewError(-1, "玩家出牌的时候出错.")
	}

	//成功之后需要把炸弹的信息保存下来
	if out.GetIsBomb() {
		d.addBombTongjiInfo(out)
		d.setWinValue(d.GetQingDizhuValue() * 2)
	}

	//initChechOut
	d.OutPai = out        //给desk赋值

	//返回成功的消息,广播用户出的牌
	ack := newProto.NewDdzOutCardsAck()
	*ack.UserId = user.GetUserId()
	*ack.CardType = ddzproto.DdzPaiType(out.GetType())
	*ack.RemainPokers = user.GetHandPaiCount()
	ack.OutCards = out.GetClientPokers()
	d.BroadCastProto(ack)


	//判断游戏是否结束
	if user.GetHandPaiCount() == 0 {
		//出牌的人 手牌为0，表示游戏结束
		d.Lottery(user)
		return nil
	}

	d.NextUserChuPai()        //出牌之后，下一个人出牌
	return nil
}

//用户过牌，不出牌
func (d *DdzDesk) ActPass(userId uint32) error {
	err := d.CheckOutUser(userId)
	if err != nil {

	}

	user := d.GetUserByUserId(userId)
	if user == nil {
		return Error.NewError(-1, "没有找到玩家,出牌失败")
	}

	//返回成功的消息 todo  返回ack
	ack := newProto.NewDdzActGuoAck()
	user.WriteMsg(ack)

	//轮到下一个玩家
	d.NextUserChuPai()        //玩家过牌，下一个人出牌
	return nil
}

//轮到下一个人出牌
func (d *DdzDesk) NextUserChuPai() error {
	index := d.GetUserIndexByUserId(d.GetActiveUserId())
	if index < 0 {
		log.E("轮到下一个玩家的时候出错,desk.activeUser[%v]", d.GetActiveUserId())
		return Error.NewError(-1, "轮到一下个玩家的时候出错.")
	}

	nextUser := d.Users[(index + 1) % int(d.GetUserCountLimit())]
	if nextUser == nil {
		log.E("轮到下一个玩家的时候出错,desk.activeUser[%v]", d.GetActiveUserId())
		return Error.NewError(-1, "轮到一下个玩家的时候出错.")
	} else {
		d.SendChuPaiOverTurn(nextUser.GetUserId())
	}
	return nil
}





//抢完地主之后的操作
func (d *DdzDesk) afterQiangDizhu() {

}

//欢乐斗地主抢完地主之后的操作
func (d *DdzDesk) HLAfterQiangDizhu() {
	log.T("抢地主结束，现在开始 加倍的逻辑....")

	//首先是判断是否有人当地主,如果没有人当地主，那么洗牌开始下一局
	if d.GetDiZhuUserId() == 0 {
		d.Begin()        ///重新开始
		return
	}

	//显示底牌
	showDiPai := newProto.NewDdzStartPlay()
	showDiPai.FootPokers = d.GetDiPaiClientPokers()
	*showDiPai.FootRate = d.GetFootRate()
	*showDiPai.Dizhu = d.GetDiZhuUserId()
	d.BroadCastProto(showDiPai)

	//有地主的情况,抢地主完成之后开会加倍不加倍的操作
	d.sendJiaBeiOverTurn(d.GetDiZhuUserId())        //发送加倍的广播
}


//明牌
func (d *DdzDesk) ShowHandPokers(userId uint32) error {
	d.Lock()
	defer d.Unlock()
	//处理明牌的逻辑
	return nil
}

//离开房间的协议
/**
	用户离开房间之后，需要设置用户的状态为离开的状态
 */
func (d *DdzDesk) LeaveDesk(userId uint32) error {

	//查找user是否存在
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("没有找到用户，用户离开失败...")
		return Error.NewFailError("没有找到用户")
	}

	//检测user的状态，已经里开的不用离开

	//
	*user.IsLeave = true        //设置用户离开...下句开始之前需要删除掉此用户

	ack := newProto.NewDdzAckLeaveDesk()
	user.WriteMsg(ack)        //todo 这里是需要广播还是单独回复?	//暂时没有做

	return nil
}

