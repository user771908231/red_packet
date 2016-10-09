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
	}

	a.WriteMsg(result)
}

/**

进入房间的逻辑
1，判断是否是重新进入房间：离开之后进入房间，掉线之后进入房间
2，进入成功【只】返回gameinfo
3，进入失败【只】返回AckEnterRoom


 */
func HandlerGame_EnterRoom(m *mjProto.Game_EnterRoom, a gate.Agent) {
	log.T("收到请求，HandlerGame_EnterRoom(m[%v],a[%v])", m, a)
	//todo 根据游戏类型不同，加入房间的方式也不同...这里可以通过agent 来区分

	//1,找到合适的room
	room := majiang.GetMJRoom()
	if room == nil {
		//没有找到room，进入房间失败
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = intCons.ACK_RESULT_ERROR
		a.WriteMsg(ack)
		return
	}

	//2,返回进入的desk
	desk, err := room.EnterRoom("", 0, a)
	if err != nil || desk == nil {
		//进入房间失败
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = intCons.ACK_RESULT_ERROR
		a.WriteMsg(ack)
	} else {
		//3,更新userSession,返回desk 的信息
		majiang.UpdateSession(m.GetHeader().GetUserId(), majiang.MJUSER_SESSION_GAMESTATUS_FRIEND, desk.GetRoomId(), desk.GetDeskId(), desk.GetPassword())
		gameinfo := desk.GetGame_SendGameInfo()
		*gameinfo.SenderUserId = m.GetHeader().GetUserId()
		a.WriteMsg(gameinfo)
	}
}

//用户开始准备游戏
func HandlerGame_Ready(m *mjProto.Game_Ready, a gate.Agent) {
	log.Debug("收到请求，game_Ready(m[%v],a[%v])", m, a)
	desk := majiang.GetMjDeskBySession(m.GetHeader().GetUserId())
	if desk == nil {
		// 准备失败
		result := newProto.NewGame_AckCreateRoom()
		*result.Header.Code = intCons.ACK_RESULT_ERROR
		*result.Header.Error = "准备失败"
		a.WriteMsg(result)
	} else {
		err := desk.Ready(m.GetHeader().GetUserId())
		if err != nil {
			//准备失败
			result := newProto.NewGame_AckCreateRoom()
			*result.Header.Code = intCons.ACK_RESULT_ERROR
			*result.Header.Error = "准备失败"
			a.WriteMsg(result)
		} else {
			//准备成功,发送准备成功的广播
			result := newProto.NewGame_AckCreateRoom()
			*result.Header.Code = intCons.ACK_RESULT_SUCC
			*result.Header.Error = "准备成功"
			log.T("广播user[%v]在desk[%v]准备成功的广播..", m.GetHeader().GetUserId(), desk.GetDeskId())
			desk.BroadCastProto(result)

			//准备成功之后，是否需要开始游戏...
			desk.AfterReady()

		}
	}
}



//定缺
func HandlerGame_DingQue(m *mjProto.Game_DingQue, a gate.Agent) {
	log.Debug("收到请求，HandlerGame_DingQue(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckEnterRoom()
	//a.WriteMsg(result)
}

//换3张
func HandlerGame_ExchangeCards(m *mjProto.Game_ExchangeCards, a gate.Agent) {
	log.Debug("收到请求，HandlerGame_ExchangeCards(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckExchangeCards()
	//a.WriteMsg(result)
}

//出牌
func HandlerGame_SendOutCard(m *mjProto.Game_SendOutCard, a gate.Agent) {
	log.Debug("收到请求，HandlerGame_SendOutCard(m[%v],a[%v])", m, a)
	//result := newProto.NewGame_AckEnterRoom()
	//a.WriteMsg(result)
}

//碰
func HandlerGame_ActPeng(m *mjProto.Game_ActPeng, a gate.Agent) {
	log.Debug("收到请求，HandlerGame_ActPeng(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActPeng{}
	result.Header = newProto.SuccessHeader()

	a.WriteMsg(result)
}

//杠
func HandlerGame_ActGang(m *mjProto.Game_ActGang, a gate.Agent) {
	log.Debug("收到请求，game_ActGang(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActGang{}
	result.Header = newProto.SuccessHeader()

	a.WriteMsg(result)
}

//过
func HandlerGame_ActGuo(m *mjProto.Game_ActGuo, a gate.Agent) {
	log.Debug("收到请求，game_ActGuo(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActGuo{}
	result.Header = newProto.SuccessHeader()

	a.WriteMsg(result)
}

//胡
func HandlerGame_ActHu(m *mjProto.Game_ActHu, a gate.Agent) {
	log.Debug("收到请求，game_ActHu(m[%v],a[%v])", m, a)

	result := &mjProto.Game_AckActHu{}
	result.Header = newProto.SuccessHeader()

	a.WriteMsg(result)
}
