package MJService

import (
	mjProto "casino_majiang/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_majiang/msg/funcsInit"
	"casino_majiang/service/majiang"
	"casino_server/conf/intCons"
	"casino_server/service/userService"
)


//service的作用就是handler的具体实现


/*
	创建room
	用户创建房间的逻辑
	1,如果用户之前已经创建了房间，怎么处理？
	2,余额不足怎么处理
	3,创建成功之后

 */
func HandlerGame_CreateRoom(m *mjProto.Game_CreateRoom, a gate.Agent) {
	log.T("收到请求，HandlerGame_CreateRoom(m[%v],a[%v])", m, a)
	//1,查询用户是否已经创建了房间...

	//2,开始创建房间
	desk := majiang.FMJRoomIns.CreateDesk(m)

	//3,返回数据
	result := newProto.NewGame_AckCreateRoom()

	if desk == nil {
		*result.Header.Code = intCons.ACK_RESULT_ERROR
		log.Error("用户[%v]创建房间失败...")
	} else {
		*result.Header.Code = intCons.ACK_RESULT_SUCC
		*result.Password = desk.GetPassword()
		*result.DeskId = desk.GetDeskId()
		*result.CreateFee = desk.GetCreateFee()
		result.RoomTypeInfo = desk.GetRoomTypeInfo()
		*result.UserBalance = userService.GetUserDiamond(m.GetHeader().GetUserId())

		//创建成功之后，用户自动进入房间...
		HandlerGame_EnterRoom(m.GetHeader().GetUserId(), desk.GetPassword(), a)
	}

	a.WriteMsg(result)

}

/**

进入房间的逻辑
1，判断是否是重新进入房间：离开之后进入房间，掉线之后进入房间
2，进入成功【只】返回gameinfo
3，进入失败【只】返回AckEnterRoom
 */
func HandlerGame_EnterRoom(userId uint32, key string, a gate.Agent) {
	log.T("收到请求，HandlerGame_EnterRoom(m[%v],a[%v])", userId, a)

	//1,找到合适的room
	room := majiang.GetMJRoom()
	if room == nil {
		//没有找到room，进入房间失败
		log.T("用户[%v]进入房间失败，没有找到对应的room", userId)
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = intCons.ACK_RESULT_ERROR
		a.WriteMsg(ack)
		return
	}

	//2,返回进入的desk
	desk, err := room.EnterRoom(key, userId, a)
	if err != nil || desk == nil {
		//进入房间失败
		log.E("用户[%v]进入房间,key[%v]失败err[%v]", userId, key, err)
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = intCons.ACK_RESULT_ERROR
		a.WriteMsg(ack)
	} else {
		//3,更新userSession,返回desk 的信息
		majiang.UpdateSession(userId, majiang.MJUSER_SESSION_GAMESTATUS_FRIEND, desk.GetRoomId(), desk.GetDeskId(), desk.GetPassword())
		gameinfo := desk.GetGame_SendGameInfo()
		*gameinfo.SenderUserId = userId
		a.WriteMsg(gameinfo)
	}
}

//用户开始准备游戏
func HandlerGame_Ready(m *mjProto.Game_Ready, a gate.Agent) {
	log.T("收到请求，game_Ready(m[%v],a[%v])", m, a)
	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId())
	if desk == nil {
		// 准备失败
		log.E("用户[%v]准备失败.因为没有找到对应的desk", m.GetHeader().GetUserId())
		result := newProto.NewGame_AckReady()
		*result.Header.Code = intCons.ACK_RESULT_ERROR
		*result.Header.Error = "准备失败"
		a.WriteMsg(result)
	} else {
		err := desk.Ready(m.GetHeader().GetUserId())
		if err != nil {
			//准备失败
			result := newProto.NewGame_AckReady()
			*result.Header.Code = intCons.ACK_RESULT_ERROR
			*result.Header.Error = "准备失败"
			a.WriteMsg(result)
		} else {
			//准备成功,发送准备成功的广播
			result := newProto.NewGame_AckReady()
			*result.Header.Code = intCons.ACK_RESULT_SUCC
			*result.Header.Error = "准备成功"
			*result.UserId = m.GetHeader().GetUserId()
			log.T("广播user[%v]在desk[%v]准备成功的广播..", m.GetHeader().GetUserId(), desk.GetDeskId())
			desk.BroadCastProto(result)

			//准备成功之后，是否需要开始游戏...
			desk.AfterReady()

		}
	}
}



//定缺
/**
定缺之后 需要判断，如果所有人都已经定缺了，那么庄开始发牌，

 */
func HandlerGame_DingQue(m *mjProto.Game_DingQue, a gate.Agent) {
	log.T("收到请求，HandlerGame_DingQue(m[%v],a[%v])", m, a)

	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId())
	err := desk.DingQue(m.GetHeader().GetUserId(), m.GetColor())
	if err != nil {
		log.E("定缺失败...")
		return
	}

	//发送成功的回复
	a.WriteMsg(m)

	//如果所有人都定缺了，那么可以通知庄打牌了..
	if desk.AllDingQue() {

		//通知庄家打一张牌,这里初始化信息，这里应该是广播的..

		//注意是否可以碰，可以杠牌，可以胡牌，只有当时人才能看到，所以广播的和当事人的收到的数据不一样...
		overTurn := newProto.NewGame_OverTurn()
		overTurn.ActCard = nil
		*overTurn.ActType = 1
		*overTurn.CanHu = false
		*overTurn.CanPeng = false
		*overTurn.CanGang = false
		desk.BroadCastProtoExclusive(overTurn, desk.GetBanker())

		//发送给当事人
		bankUser := desk.GetBankerUser()

		overTurn.ActCard = nil
		*overTurn.ActType = 1
		*overTurn.CanHu = bankUser.MJHandPai.GetCanHu()
		*overTurn.CanGang = bankUser.MJHandPai.GetCanGang()
		*overTurn.CanPeng = bankUser.MJHandPai.GetCanPeng()
		bankUser.SendOverTurn(overTurn)
	}

}

//换3张
func HandlerGame_ExchangeCards(m *mjProto.Game_ExchangeCards, a gate.Agent) {
	log.Debug("收到请求，HandlerGame_ExchangeCards(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckExchangeCards()
	//a.WriteMsg(result)
}

//出牌

/**
	打牌的协议
	1,接收到用户打出的牌
	2,判断其他人是否可以需要这张牌,以用户为单位
	3,一次让每个人判断牌是否需要...
 */

func HandlerGame_SendOutCard(m *mjProto.Game_SendOutCard, a gate.Agent) {
	log.Debug("收到请求，HandlerGame_SendOutCard(m[%v],a[%v])", m, a)

	//检测参数
	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId())
	if desk == nil {
		//打牌失败，因为没有找到对应的麻将桌子

		return
	}

	//打牌之后的逻辑,初始化判定事件
	desk.InitCheckCase(majiang.InitMjPaiByIndex(int(m.GetCardId())))
	if desk.GetCheckCase() == nil {
		//表示无人需要，直接给用户返回无人需要
		//给下一个人摸排，并且移动指针
	} else {
		//如果判定事件不为空，那么开始执行判断事件
		desk.DoCheckCase()
		//并且给用户返回，牌需要判定，牌是正在打出的状态，并没有落到桌子上
	}
}

//碰
func HandlerGame_ActPeng(m *mjProto.Game_ActPeng, a gate.Agent) {
	log.T("收到请求，HandlerGame_ActPeng(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActPeng{}
	result.Header = newProto.SuccessHeader()

	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId()) //通过userId 的session 得到对应的desk
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result.Header = newProto.ErrorHeader()
		return
	}

	//找到玩家
	user := desk.GetUserByUserId(m.GetHeader().GetUserId())
	if user == nil {
		return
	}

	user.ActHu()        //这里碰牌的逻辑

	//todo 设置checkCase 为已经验证过了

	//设置下一个人摸牌
	desk.SendMopaiOverTurn(nil)

	a.WriteMsg(result)
}

//杠
func HandlerGame_ActGang(m *mjProto.Game_ActGang, a gate.Agent) {
	log.Debug("收到请求，game_ActGang(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActGang{}
	result.Header = newProto.SuccessHeader()

	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId()) //通过userId 的session 得到对应的desk
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result.Header = newProto.ErrorHeader()
		return
	}

	//找到玩家
	user := desk.GetUserByUserId(m.GetHeader().GetUserId())
	if user == nil {
		return
	}

	user.ActHu()        //这里碰牌的逻辑

	//todo 设置checkCase 为已经验证过了

	//设置下一个人摸牌
	desk.SendMopaiOverTurn(user)

	//杠牌之后 自己摸牌
	a.WriteMsg(result)
}

//过

/**

	设置checkCaseBean为已经check过就行了，不做其他的处理...
 */
func HandlerGame_ActGuo(m *mjProto.Game_ActGuo, a gate.Agent) {
	log.Debug("收到请求，game_ActGuo(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActGuo{}
	result.Header = newProto.SuccessHeader()
	// 设置当前CheckBean 为已经check ，处理下一个checkBean


	a.WriteMsg(result)
}




//胡

/**
	胡牌需要注意的是：
	1,如何区分 只自摸还是点炮...
	2,点炮的时候需要注意区分  抢杠，杠上炮，普通点炮
 */
func HandlerGame_ActHu(m *mjProto.Game_ActHu) {
	log.Debug("收到请求，game_ActHu(m[%v])", m)

	//需要返回的数据
	result := &mjProto.Game_AckActHu{}

	//区分自摸点炮:1,如果自己的手牌就已经糊了（或者如果自己自己的牌是14，11，8，5，2 张的时候），那么就自摸，如果需要加上判定牌，那就是点炮
	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId()) //通过userId 的session 得到对应的desk
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result.Header = newProto.ErrorHeader()
		return
	}

	//找到玩家
	user := desk.GetUserByUserId(m.GetHeader().GetUserId())
	if user == nil {
		return
	}

	//玩家胡牌
	err := user.ActHu()
	if err != nil {
		//如果这里出现胡牌失败，证明是系统有问题...
		log.E("用户[%v]胡牌失败...", m.GetHeader().GetUserId())
		result.Header = newProto.ErrorHeader()
		user.WriteMsg(result)        //返回失败的信息
		return
	}

	//胡牌成功之后的处理...
	desk.SetNestUserCursor(user.GetUserId())        // 胡牌之后 设置当前操作的用户为当前胡牌的人...
	desk.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), majiang.CHECK_CASE_bean_STATUS_CHECKED)        // update checkCase...

	//todo 胡牌之后，如果只剩下一个人..那么这句游戏结束...
	if desk.Time2Lottery() {
		desk.Lottery()
		//因为可以开奖了，所以不操作后边的，直接返回
		return
	}

	//还没有到牌局结束的时候，轮到下一个人...
	mjHandPai := user.GetMJHandPai()
	if mjHandPai == nil {
		//服务器出错
		return
	}

	//胡牌之后的处理
	if mjHandPai.GetZiMo() {
		// 如果是自摸，则轮到下一个人摸排，轮到胡牌的下一个人...
		desk.SendMopaiOverTurn(nil)                //给下一个人发送overTurn 发牌的类型...
	} else {
		//todo 如果是点炮,那么计算判断其他人是否需要继续胡牌,有的话继续胡牌，没有的话设置下一个人摸牌...

		//1,找到checkCase 是否有下一个人胡牌，如果有，那么让下一个人验证，如果没有，下一个人摸牌。。
		nextBean := desk.CheckCase.GetBuBean(majiang.CHECK_CASE_bean_STATUS_CHECKED)
		if nextBean == nil {
			//表示没有下一个胡牌的人了,和自摸同样的处理
			desk.CheckCase.UpdateChecStatus(majiang.CHECK_CASE_STATUS_CHECKED)
			desk.SendMopaiOverTurn(nil)                //给下一个人发送overTurn 发牌的类型...
		} else {
			//发送overTurn 给下一个判定的人...

			overTurn := newProto.NewGame_OverTurn()
			*overTurn.UserId = nextBean.GetUserId()
			*overTurn.CanGang = nextBean.GetCanGang()
			*overTurn.CanPeng = nextBean.GetCanPeng()
			*overTurn.CanHu = nextBean.GetCanHu()
			overTurn.ActCard = desk.CheckCase.CheckMJPai.GetCardInfo()        //
			*overTurn.ActType = majiang.OVER_TURN_ACTTYPE_OTHER

			///发送overTurn 的信息
			desk.GetUserByUserId(nextBean.GetUserId()).SendOverTurn(overTurn)

		}
	}
}


