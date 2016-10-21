package MJService

import (
	mjProto "casino_majiang/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_majiang/msg/funcsInit"
	"casino_majiang/service/majiang"
	"casino_server/conf/intCons"
	"casino_server/service/userService"
	"time"
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

	if desk == nil {
		result := newProto.NewGame_AckCreateRoom()
		*result.Header.Code = intCons.ACK_RESULT_ERROR
		log.Error("用户[%v]创建房间失败...")
		a.WriteMsg(result)

	} else {
		log.T("用户[%v]创建房间成功，roomKey[%v]", desk.GetOwner(), desk.GetPassword())
		result := newProto.NewGame_AckCreateRoom()
		*result.Header.Code = intCons.ACK_RESULT_SUCC
		*result.Password = desk.GetPassword()
		*result.DeskId = desk.GetDeskId()
		*result.CreateFee = desk.GetCreateFee()
		result.RoomTypeInfo = desk.GetRoomTypeInfo()
		*result.UserBalance = userService.GetUserDiamond(m.GetHeader().GetUserId())
		a.WriteMsg(result)

		//创建成功之后，用户自动进入房间...
		HandlerGame_EnterRoom(m.GetHeader().GetUserId(), desk.GetPassword(), a)
	}

}

/**

进入房间的逻辑
1，判断是否是重新进入房间：离开之后进入房间，掉线之后进入房间
2，进入成功【只】返回gameinfo
3，进入失败【只】返回AckEnterRoom
 */
func HandlerGame_EnterRoom(userId uint32, key string, a gate.Agent) {
	log.T("收到请求，HandlerGame_EnterRoom(userId[%v],key[%v])", userId, key)

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
		s, _ := majiang.UpdateSession(userId, majiang.MJUSER_SESSION_GAMESTATUS_FRIEND, desk.GetRoomId(), desk.GetDeskId(), desk.GetPassword())
		if s != nil {
			//给agent设置session
			a.SetUserData(s)
		}

		gameinfo := desk.GetGame_SendGameInfo(userId)
		*gameinfo.SenderUserId = userId
		//a.WriteMsg(gameinfo)
		log.T("用户[%v]进入房间之后，返回的数据gameInfo[%v]", userId, gameinfo)
		desk.BroadCastProto(gameinfo)

	}
}

//用户开始准备游戏
func HandlerGame_Ready(m *mjProto.Game_Ready, a gate.Agent) {
	log.T("收到请求，game_Ready(m[%v])", m)
	userId := m.GetHeader().GetUserId()
	desk := majiang.GetMjDeskBySession(userId)
	if desk == nil {
		// 准备失败
		log.E("用户[%v]准备失败.因为没有找到对应的desk", userId)
		result := newProto.NewGame_AckReady()
		*result.Header.Code = intCons.ACK_RESULT_ERROR
		*result.Header.Error = "准备失败"
		a.WriteMsg(result)
	} else {
		err := desk.Ready(userId)
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
			*result.UserId = userId
			log.T("广播user[%v]在desk[%v]准备成功的广播..", userId, desk.GetDeskId())
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
		log.E("用户[%v]定缺失败...", m.GetHeader().GetUserId())
		return
	}

	//如果所有人都定缺了，那么可以通知庄打牌了..
	if desk.AllDingQue() {
		//首先发送定缺结束的广播，然后发送庄家出牌的广播...

		ques := desk.GetDingQueEndInfo()
		desk.BroadCastProto(ques)

		//通知庄家打一张牌,这里初始化信息，这里应该是广播的..

		//注意是否可以碰，可以杠牌，可以胡牌，只有当时人才能看到，所以广播的和当事人的收到的数据不一样...
		overTurn := newProto.NewGame_OverTurn()
		*overTurn.UserId = desk.GetBanker()
		*overTurn.ActType = majiang.OVER_TURN_ACTTYPE_MOPAI
		*overTurn.CanPeng = false        ///自己的手牌不能碰

		//广播时候的信息
		overTurn.ActCard = nil
		*overTurn.CanHu = false
		*overTurn.CanGang = false
		*overTurn.CanPeng = false
		desk.BroadCastProtoExclusive(overTurn, desk.GetBanker())

		//发送给当事人时候的信息
		bankUser := desk.GetBankerUser()
		*overTurn.CanHu = bankUser.GameData.HandPai.GetCanHu()

		//判断是否可以杠牌
		canGangBool, gangPais := bankUser.GameData.HandPai.GetCanGang(nil)//判断自己摸牌的情况，有可能有多个杠牌
		*overTurn.CanGang = canGangBool
		if canGangBool && gangPais != nil {
			for _, p := range gangPais {
				if p != nil {
					overTurn.GangCards = append(overTurn.GangCards, p.GetCardInfo())
				}
			}
		}

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
	log.T("收到请求，HandlerGame_SendOutCard(m[%v],a[%v])", m, a)
	userId := m.GetHeader().GetUserId()
	//检测参数
	desk := majiang.GetMjDeskBySession(userId)
	if desk == nil {
		//打牌失败，因为没有找到对应的麻将桌子
		log.E("用户[%v]打牌", userId)
		return
	}

	err := desk.ActOut(userId, m.GetCardId())
	if err != nil {
		//打牌失败
		log.E("打牌失败...errMsg[%v]", err)
		return
	}

	log.T("用户[%v]已经打牌，处理下一个checkCase", userId)
	desk.DoCheckCase(nil)        //打牌之后，别人判定牌

}

//碰
func HandlerGame_ActPeng(m *mjProto.Game_ActPeng, a gate.Agent) {
	log.T("收到请求，HandlerGame_ActPeng(m[%v],a[%v])", m, a)

	//找到桌子
	userId := m.GetHeader().GetUserId()
	desk := majiang.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result := &mjProto.Game_AckActPeng{}
		result.Header = newProto.ErrorHeader()
		a.WriteMsg(result)
		return
	}

	//开始碰牌
	err := desk.ActPeng(userId)
	if err != nil {
		log.E("服务器错误: 用户[%v]碰牌失败...", userId)
		//todo 需要做特殊处理
	}

	//操作下一个
	//desk.DoCheckCase(desk.GetUserByUserId(userId))        //碰牌之后，别人判定牌	//碰牌之后不需要处理desk.DoCheckCase
}


//杠
func HandlerGame_ActGang(m *mjProto.Game_ActGang) {
	log.T("收到请求，game_ActGang(m[%v])", m)
	userId := m.GetHeader().GetUserId()

	result := &mjProto.Game_AckActGang{}
	result.Header = newProto.SuccessHeader()

	desk := majiang.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result.Header = newProto.ErrorHeader()
		return
	}

	//先杠牌
	err := desk.ActGang(m.GetHeader().GetUserId(), m.GetGangCard().GetId())
	if err != nil {
		log.E("服务器错误：用户[%v]杠牌的时候出错err[%v]", userId, err)
	}

	time.Sleep(time.Second * 1)        //间隔两秒 进行下一个动作
	//处理下一个人
	desk.DoCheckCase(desk.GetUserByUserId(userId))        //杠牌之后，处理下一个判定牌
}

//过

/**

	设置checkCaseBean为已经check过就行了，不做其他的处理...

	注意 *   本协议  只有判断别人出牌是否需要的时候，才会请求...
	胡牌的过，之后的人可以继续碰或者杠
 */
func HandlerGame_ActGuo(m *mjProto.Game_ActGuo) {
	log.T("收到杠牌的请求，game_ActGuo(m[%v])", m)

	userId := m.GetHeader().GetUserId()
	desk := majiang.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk
	user := desk.GetUserByUserId(userId)
	if desk.CheckCase == nil {
		/**
			只有判断别人打的牌的时候，需要过的时候才会请求这个协议，自己摸牌 需不需要过的时候不需要请求这个协议...
		 */
		log.E("玩家【%v】过牌的时候出错，因为checkCase为nil", userId)
		return

	}
	err := desk.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), majiang.CHECK_CASE_BEAN_STATUS_PASS)        // update checkCase...
	if err != nil {
		log.T("过牌的时候失败，err[%v]", err)
	}
	//设置为过

	//返回信息,过 只返回给过的
	result := &mjProto.Game_AckActGuo{}
	result.Header = newProto.SuccessHeader()
	// 设置当前CheckBean 为已经check ，处理下一个checkBean
	user.WriteMsg(result)

	//进行下一个判断
	desk.DoCheckCase(nil)        //过牌之后，处理下一个判定牌
}

//胡

/**
	胡牌需要注意的是：
	1,如何区分 只自摸还是点炮...
	2,点炮的时候需要注意区分  抢杠，杠上炮，普通点炮
 */
func HandlerGame_ActHu(m *mjProto.Game_ActHu) {
	log.T("收到胡牌请求，game_ActHu(m[%v])", m)

	//需要返回的数据
	userId := m.GetHeader().GetUserId()

	//区分自摸点炮:1,如果自己的手牌就已经糊了（或者如果自己自己的牌是14，11，8，5，2 张的时候），那么就自摸，如果需要加上判定牌，那就是点炮
	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId()) //通过userId 的session 得到对应的desk
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result := newProto.NewGame_AckActHu()
		result.Header = newProto.ErrorHeader()
		return
	}

	//开始胡牌...
	err := desk.ActHu(userId)
	if err != nil {
		log.E("服务器错误，胡牌失败..")
	}

	//这里是否需要广播胡牌的广播...

	//胡牌之后，需要判断游戏是否结束...
	if desk.Time2Lottery() {
		desk.Lottery()
		//因为可以开奖了，所以不操作后边的，直接返回
		return
	} else {
		//处理下一个
		desk.DoCheckCase(nil)        //胡牌之后，处理下一个判定牌
	}

}


