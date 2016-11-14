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
	"errors"
	"casino_majiang/gamedata/dao"
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
	log.T("收到请求，HandlerGame_CreateRoom(m[%v])", m)
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
		//创建desk成功...设置为开始准备的状态
		desk.SetStatus(majiang.MJDESK_STATUS_READY)        //设置为开始准备的状态

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
	room := majiang.GetFMJRoom()
	if room == nil {
		//没有找到room，进入房间失败
		log.T("用户[%v]进入房间失败，没有找到对应的room", userId)
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = intCons.ACK_RESULT_ERROR
		a.WriteMsg(ack)
		return
	}

	//2,返回进入的desk,如果进入房间失败，则返回进入房间失败...
	desk, reconnect, err := room.EnterRoom(key, userId, a)
	if err != nil || desk == nil {
		//进入房间失败
		log.E("用户[%v]进入房间,key[%v]失败err[%v]", userId, key, err)
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = intCons.ACK_RESULT_ERROR
		a.WriteMsg(ack)
		return
	}

	//3,更新userSession,返回desk 的信息
	s, _ := majiang.UpdateSession(userId, majiang.MJUSER_SESSION_GAMESTATUS_FRIEND, desk.GetRoomId(), desk.GetDeskId(), desk.GetPassword())
	if s != nil {
		//给agent设置session
		a.SetUserData(s)
	}

	gameinfo := desk.GetGame_SendGameInfo(userId)
	log.T("用户[%v]进入房间,reconnect[%v]之后，返回的数据gameInfo[%v]", userId, reconnect, gameinfo)
	desk.BroadCastProto(gameinfo)


	//如果是重新进入房间，需要发送重近之后的处理
	if reconnect {
		time.Sleep(time.Second * 5)
		desk.SendReconnectOverTurn(userId)
	}
}

//解散房间
func HandlerDissolveDesk(owner uint32) error {
	//1,通过房主找到房间
	desk := majiang.GetMjDeskBySession(owner)
	if desk == nil {
		log.T("没有找到user[%v]对应的desk，解散房间失败", owner)
		return errors.New("解散房间失败...")
	}

	if desk.GetOwner() != owner {
		log.T("通过owner[%v]找到的desk的owner  不正确..", owner, desk.GetOwner())
		return errors.New("房间的房主不正确")
	}

	//开始解散房间
	err := majiang.GetFMJRoom().DissolveDesk(desk, true);
	if err != nil {
		return errors.New("解散朋友桌子的desk 失败...")
	}

	return nil
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
		return
	}

	//判断desk状态
	if desk.IsNotPreparing() {
		// 准备失败
		log.E("用户[%v]准备失败.desk[%v]不再准备的状态...", userId, desk.GetDeskId())
		result := newProto.NewGame_AckReady()
		*result.Header.Code = intCons.ACK_RESULT_ERROR
		*result.Header.Error = "准备失败"
		a.WriteMsg(result)
		return
	}

	//开始准备
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
		log.T("广播user[%v]在desk[%v]准备成功的广播..string(%v)", userId, desk.GetDeskId(), result.String())
		desk.BroadCastProto(result)

		//准备成功之后，是否需要开始游戏...
		desk.AfterReady()
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

		//游戏开始 庄家打牌
		desk.BeginStart()
	}

}




//换3张
func HandlerGame_ExchangeCards(m *mjProto.Game_ExchangeCards) {
	log.T("收到请求，HandlerGame_ExchangeCards(m[%v])", m)
	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId())        //得到desk
	//开始换三张
	desk.DoExchange(m.GetHeader().GetUserId(), m.GetExchangeNum(), m.GetExchangeOutCards())
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


	请求的场景:
	1,别人点炮，自己不糊的时候，需要请求过
	2,别人打牌，自己可以碰的时候，如果自己不碰，那么点过
	3,别人打牌，自己可以杠的时候，如果自己不杠，那么需要点过


	//过胡
	1，场景：
		别人打牌（checkBean不为nil） && 自己可以胡(canhu==true) && 自己点了过  确认为一个过胡的场景
	过胡...在下一次操作之前，如果还有其他人打牌，你将不能胡，除非翻📖比之前的要大.. 注意自摸的要 除开，之前点炮的才需要判断...



 */


func HandlerGame_ActGuo(m *mjProto.Game_ActGuo) {
	log.T("收到过牌的请求，game_ActGuo(m[%v])", m)

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

	//添加一个过hu的info,下次init的时候，需要判断是否有这个guohu
	user.AddGuoHuInfo(desk.CheckCase)

	err := desk.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), majiang.CHECK_CASE_BEAN_STATUS_PASS)        // update checkCase...
	if err != nil {
		log.T("过牌的时候失败，err[%v]", err)
	}
	//设置为过

	//返回信息,过 只返回给过的
	result := newProto.NewGame_AckActGuo()
	result.Header = newProto.SuccessHeader()
	*result.UserId = user.GetUserId()
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

//查询用户的战绩
func HandlerGame_GameRecord(userId uint32, a gate.Agent) error {
	log.T("用户[%v]请求战绩", userId)
	//战绩 mongoData
	data := dao.GetByUserId(userId)
	log.T("data[%v]", data)

	//返回数据到client
	result := newProto.NewGame_AckGameRecord()
	*result.UserId = userId
	//*result.Records
	//增加records
	for _, d := range data {
		bean := d.TransRecord()
		result.Records = append(result.Records, bean)
	}

	//发送战绩
	log.T("发送玩家[%v]的战绩[%v]", userId, result)
	a.WriteMsg(result)
	return nil
}

//聊天协议
func HandlerGame_Message(m *mjProto.Game_Message) {
	log.T("请求发送信息[%v]", m)
	userId := m.GetHeader().GetUserId()
	desk := majiang.GetMjDeskBySession(userId)
	if desk == nil {
		log.E("玩家[%v]聊天的时候没有找到desk", userId)
		return
	}
	result := newProto.NewGame_SendMessage()
	*result.UserId = m.GetHeader().GetUserId()
	*result.Id = m.GetId()
	*result.Msg = m.GetMsg()
	*result.MsgType = m.GetMsgType()
	desk.BroadCastProtoExclusive(result, result.GetUserId())
}