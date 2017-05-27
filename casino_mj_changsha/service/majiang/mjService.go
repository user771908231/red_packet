package majiang

import (
	mjProto "casino_mj_changsha/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_mj_changsha/msg/funcsInit"
	"errors"
	"casino_mj_changsha/gamedata/dao"
	"casino_common/common/log"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
)

//service的作用就是handler的具体实现
/*
	创建room
	用户创建房间的逻辑
	1,如果用户之前已经创建了房间，怎么处理？
	2,余额不足怎么处理
	3,创建成功之后
 */

var ERR_DESK_NOT_FOUND = Error.NewError(consts.ACK_RESULT_ERROR, "房间没有找到")

func HandlerGame_CreateDesk(m *mjProto.Game_CreateRoom, a gate.Agent) {
	if m.GetHeader().GetUserId() == 0 {
		return
	}
	//log.T("收到请求，HandlerGame_CreateRoom(m[%v])", m)
	//1,查询用户是否已经创建了房间...

	//2,开始创建房间
	var UserCountLimit int32 = 4
	var FangCountLimit int32 = 3

	if m.GetRoomTypeInfo().GetMjRoomType() == mjProto.MJRoomType_roomType_sanRenLiangFang {
		UserCountLimit = 3
		FangCountLimit = 2
	}

	if m.GetRoomTypeInfo().GetMjRoomType() == mjProto.MJRoomType_roomType_siRenLiangFang {
		UserCountLimit = 4
		FangCountLimit = 2
	}

	if m.GetRoomTypeInfo().GetMjRoomType() == mjProto.MJRoomType_roomType_liangRenLiangFang {
		UserCountLimit = 2
		FangCountLimit = 2
	}

	if m.GetRoomTypeInfo().GetMjRoomType() == mjProto.MJRoomType_roomType_sanRenSanFang {
		UserCountLimit = 3
		FangCountLimit = 3
	}

	//如果是长沙麻将，根据长沙麻将的配置，设置人的数量
	if csconfig := m.GetRoomTypeInfo().GetChangShaPlayOptions(); csconfig != nil {
		UserCountLimit = csconfig.GetPlayerCount()
	}

	desk, err := MjroomManagerIns.FMJRoomIns.CreateFriendDesk(m.GetHeader().GetUserId(),
		int32(m.GetRoomTypeInfo().GetMjRoomType()), m.GetRoomTypeInfo().GetBoardsCout(), m.GetRoomTypeInfo().GetCapMax(),
		m.GetRoomTypeInfo().GetCardsNum(), m.GetRoomTypeInfo().GetSettlement(), m.GetRoomTypeInfo().GetBaseValue(),
		m.GetRoomTypeInfo().GetPlayOptions().GetZiMoRadio(), m.GetRoomTypeInfo().GetPlayOptions().GetOthersCheckBox(),
		m.GetRoomTypeInfo().GetPlayOptions().GetHuRadio(), m.GetRoomTypeInfo().GetPlayOptions().GetDianGangHuaRadio(),
		m.GetRoomTypeInfo().GetBoardsCout(), m.GetRoomTypeInfo().GetChangShaPlayOptions(), UserCountLimit, FangCountLimit)
	//3,返回数据

	if desk == nil {
		result := newProto.NewGame_AckCreateRoom()
		*result.Header.Code = consts.ACK_RESULT_ERROR
		*result.Header.Error = Error.GetErrorMsg(err)
		log.W("用户[%v]创建房间失败...err[%v]", m.GetHeader().GetUserId(), err)
		a.WriteMsg(result)

	} else {
		log.T("用户[%v]创建房间成功，deskKey[%v],desk[%v]", desk.GetOwner(), desk.GetPassword(), desk.GetDeskId())
		//创建成功之后，用户自动进入房间...
		HandlerGame_EnterDesk(m.GetHeader().GetUserId(), desk.GetPassword(), desk.GetRoomType(), 0, ENTERTYPE_NORMAL, a)
	}

}

/**
进入房间的逻辑
1，判断是否是重新进入房间：离开之后进入房间，掉线之后进入房间
2，进入成功【只】返回gameinfo
3，进入失败【只】返回AckEnterRoom
 */

var ENTERTYPE_NORMAL int32 = 0
var ENTERTYPE_AUTO int32 = 1

func HandlerGame_EnterDesk(userId uint32, key string, roomType int32, roomLevel int32, enterType int32, a gate.Agent) error {
	log.T("收到请求，HandlerGame_EnterRoom(userId[%v],key[%v]),remoteAddr[%v],roomType[%v],roomLevel[%v]",
		userId, key, a.RemoteAddr(), roomType, roomLevel)

	//兼容老版本 直接进朋友桌
	if roomType == 0 && key != "" {
		roomType = ROOMTYPE_FRIEND
	}

	//1,找到合适的room
	room := MjroomManagerIns.GetFMJRoom()
	if room == nil {
		//没有找到room，进入房间失败
		log.T("用户[%v]进入房间失败，没有找到对应的room", userId)
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = consts.ACK_RESULT_ERROR
		*ack.Header.Error = "房间号输入错误"
		a.WriteMsg(ack)
		return errors.New("进入房间失败")
	}

	//2,返回进入的desk,如果进入房间失败，则返回进入房间失败...
	err := room.EnterRoom(key, userId, a)
	if err != nil {
		//进入房间失败
		if enterType == ENTERTYPE_NORMAL {
			log.T("用户[%v] ENTERTYPE_NORMAL (正常进入房间)进入房间失败，没有找到对应的room,返回进入失败的错误信息", userId)
			ack := newProto.NewGame_AckEnterRoom()
			*ack.Header.Code = Error.GetErrorCode(err)
			*ack.Header.Error = Error.GetErrorMsg(err)
			a.WriteMsg(ack)
		} else if enterType == ENTERTYPE_AUTO {
			log.T("用户[%v] ENTERTYPE_AUTO (自动进入房间)进入房间失败，没有找到对应的room", userId)
			ack := newProto.NewGame_AckEnterRoom()
			*ack.Header.Code = consts.ACK_RESULT_ERROR
			*ack.Header.Error = ""
			a.WriteMsg(ack)
		}
		return errors.New("进入房间失败")
	}
	return nil
}

//解散房间
func HandlerDissolveDesk(owner uint32) error {
	//1,通过房主找到房间
	desk := MjroomManagerIns.GetMjDeskBySession(owner) //解散房间
	if desk == nil {
		log.T("没有找到user[%v]对应的desk，解散房间失败", owner)
		return errors.New("解散房间失败...")
	}

	if desk.GetOwner() != owner {
		log.T("通过owner[%v]找到的desk的owner[%v]  不正确..", owner, desk.GetOwner())
		return errors.New("房间的房主不正确")
	}

	//开始解散房间
	err := MjroomManagerIns.GetFMJRoom().DissolveDesk(desk, true);
	if err != nil {
		return errors.New("解散朋友桌子的desk 失败...")
	}

	return nil
}

//用户开始准备游戏
func HandlerGame_Ready(userId uint32, a gate.Agent) {
	//log.T("收到请求，game_Ready([%v])", userId)
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //准备
	if desk == nil {
		// 准备失败
		log.W("用户[%v]准备失败.因为没有找到对应的desk", userId)
		if a == nil {
			return
		}
		result := newProto.NewGame_AckReady()
		*result.Header.Code = consts.ACK_RESULT_ERROR
		*result.Header.Error = "准备失败"
		result.UserId = proto.Uint32(userId)
		a.WriteMsg(result)
		return
	}

	//开始准备
	err := desk.Ready(userId)
	if err != nil {
		log.E("用户[%v]准备失败.err %v", userId, err)
		result := newProto.NewGame_AckReady()
		*result.Header.Code = Error.GetErrorCode(err)
		*result.Header.Error = Error.GetErrorMsg(err)
		result.UserId = proto.Uint32(userId)
		a.WriteMsg(result)
		return
	}
}
//出牌

/**
	打牌的协议
	1,接收到用户打出的牌
	2,判断其他人是否可以需要这张牌,以用户为单位
	3,一次让每个人判断牌是否需要...
 */

func HandlerGame_SendOutCard(userId uint32, cardId int32, a gate.Agent) {
	//检测参数
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //出牌
	if desk == nil {
		//打牌失败，因为没有找到对应的麻将桌子
		log.E("用户[%v]打牌失败，没有找到对应的desk.", userId)
		return
	}

	//开始打牌
	err := desk.ActOut(userId, cardId, false) //普通玩家打牌
	if err != nil { //打牌失败
		result := newProto.NewGame_AckSendOutCard() //打牌失败返回信息
		*result.Header.Code = Error.GetErrorCode(err)
		*result.Header.Error = Error.GetErrorMsg(err)
		a.WriteMsg(result)
		return
	}
}

//碰
func HandlerGame_ActChi(userId uint32, cards []*mjProto.CardInfo, a gate.Agent) {
	//找到桌子
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //吃牌
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result := &mjProto.Game_AckActPeng{}
		result.Header = newProto.ErrorHeader()
		a.WriteMsg(result)
		return
	}

	//开始碰牌
	err := desk.ActChi(userId, cards) //普通玩家碰牌
	if err != nil {
		log.E("服务器错误: 用户[%v]吃牌失败...", userId)
	}
}

func HandlerGame_ActPeng(userId uint32, a gate.Agent) {
	//找到桌子
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //碰牌
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result := &mjProto.Game_AckActPeng{}
		result.Header = newProto.ErrorHeader()
		a.WriteMsg(result)
		return
	}

	//开始碰牌
	err := desk.ActPeng(userId) //普通玩家碰牌
	if err != nil {
		log.E("服务器错误: 用户[%v]碰牌失败...", userId)
	}
}

//杠
func HandlerGame_ActGang(userId uint32, cardId int32, bu bool) {
	result := &mjProto.Game_AckActGang{}
	result.Header = newProto.SuccessHeader()

	desk := MjroomManagerIns.GetMjDeskBySession(userId) //杠牌
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result.Header = newProto.ErrorHeader()
		return
	}

	//先杠牌
	err := desk.ActGang(userId, cardId, bu) //普通玩家开始杠牌
	if err != nil {
		log.E("服务器错误：用户[%v]杠牌的时候出错err[%v]", userId, err)
	}
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

func HandlerGame_ActGuo(userId uint32, a gate.Agent) {
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //过牌
	if desk == nil {
		log.E("玩家%v过的时候失败,没有找到对应的desk", userId)

		return
	}
	err := desk.ActGuo(userId)
	if err != nil {
		log.T("%v玩家过的时候失败，err %v", desk.DlogDes(), err)
		result := newProto.NewGame_AckActGuo()
		result.Header = newProto.ErrorHeader()
		a.WriteMsg(result)
		return
	}
}

//胡

/**
	胡牌需要注意的是：
	1,如何区分 只自摸还是点炮...
	2,点炮的时候需要注意区分  抢杠，杠上炮，普通点炮
 */
func HandlerGame_ActHu(userId uint32) {
	log.T("收到胡牌请求，game_ActHu(userId [%v])", userId)

	//区分自摸点炮:1,如果自己的手牌就已经糊了（或者如果自己自己的牌是14，11，8，5，2 张的时候），那么就自摸，如果需要加上判定牌，那就是点炮
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //胡牌
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result := newProto.NewGame_AckActHu()
		result.Header = newProto.ErrorHeader()
		return
	}

	//开始胡牌...
	err := desk.ActHu(userId) //普通玩家开始胡牌
	if err != nil {
		log.E("服务器错误，胡牌失败..err[%v]", err)
		return
	}

}

//查询用户的战绩
func HandlerGame_GameRecord(m *mjProto.Game_GameRecord, a gate.Agent) error {
	log.T("用户[%v]请求战绩", m.GetUserId())

	//todo 复用
	switch m.GetGameId() {
	case int32(ddproto.CommonEnumGame_GID_DDZ):
		data := dao.GetDdzDeskRoundByUserId(m.GetUserId())
		//战绩 mongoData
		log.T("data[%v]", data)
		//返回数据到client
		result := newProto.NewGame_AckGameRecord()
		*result.UserId = m.GetUserId()
		//*result.Records
		//增加records
		for _, d := range data {
			bean := d.TransRecord()
			result.Records = append(result.Records, bean)
		}

		//发送战绩
		log.T("发送玩家[%v]的战绩[%v]", m.GetUserId(), result)
		a.WriteMsg(result)

	case int32(ddproto.CommonEnumGame_GID_MAHJONG):
		data := dao.GetMjDeskRoundByUserId(m.GetUserId())
		//战绩 mongoData
		log.T("data[%v]", data)
		//返回数据到client
		result := newProto.NewGame_AckGameRecord()
		*result.UserId = m.GetUserId()
		//*result.Records
		//增加records
		for _, d := range data {
			bean := d.TransRecord()
			result.Records = append(result.Records, bean)
		}

		//发送战绩
		log.T("发送玩家[%v]的战绩[%v]", m.GetUserId(), result)
		a.WriteMsg(result)
	default:
		data := dao.GetMjDeskRoundByUserId(m.GetUserId())
		//战绩 mongoData
		log.T("data[%v]", data)
		//返回数据到client
		result := newProto.NewGame_AckGameRecord()
		*result.UserId = m.GetUserId()
		//*result.Records
		//增加records
		for _, d := range data {
			bean := d.TransRecord()
			result.Records = append(result.Records, bean)
		}

		//发送战绩
		log.T("发送玩家[%v]的战绩[%v]", m.GetUserId(), result)
		a.WriteMsg(result)
	}
	return nil
}

//聊天协议
func HandlerGame_Message(m *ddproto.CommonReqMessage) {
	log.T("请求发送信息[%v]", m)
	userId := m.GetHeader().GetUserId()
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //发信息
	if desk != nil {
		desk.SendMessage(m)
	} else {
		log.E("玩家发送消息失败....因为没有找到玩家【%v】所在的desk", userId)
	}
}

//这个功能目前只有金币场才会用到...
//1,更换房间，第一步是先离开房间
//2,自动寻找房间，然后进入房间
func HandlerExchangeRoom(userId uint32, a gate.Agent) error {
	log.T("玩家[%v]开始更换房间...", userId)
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //换桌子
	if desk == nil {
		ack := &ddproto.CommonAckLeaveDesk{
			UserId:     proto.Uint32(userId),
			IsExchange: proto.Bool(false)}
		a.WriteMsg(ack) //回复离开房间的回复
		return ERR_DESK_NOT_FOUND
	}

	desk.ExchangeRoom(userId, a)
	return nil
}

//离开房间
/**
	朋友做：
		1，游戏没有开始的时候可以离开房间

	金币场：
		2，托管的时候可以离开房间
 */
func HandlerLeaveRoom(userId uint32, a gate.Agent) error {
	log.T("玩家[%v]开始离开房间...", userId)
	ack := &ddproto.CommonAckLeaveDesk{
		UserId:     proto.Uint32(userId),
		IsExchange: proto.Bool(false)}
	//找到对应的房间
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //离开房间

	if desk == nil {
		log.T("通过 session 没有找到玩家[%v]的房间，因为session为nil ...", userId)
		a.WriteMsg(ack) //回复离开房间的回复
		return ERR_DESK_NOT_FOUND
	}

	//朋友桌离开
	return desk.Leave(userId)
}

//申请解散房间 这个方法只有朋友桌才会有
func HandlerApplyDissolve(userId uint32, a gate.Agent) error {
	desk := MjroomManagerIns.GetFMjDeskBySession(userId)
	if desk == nil {
		return ERR_DESK_NOT_FOUND
	}

	//申请
	err := desk.ApplyDissolve(userId)
	if err != nil {
		log.E("玩家申请解散房间失败%v", err)
		return err
	}
	return nil
}

//回复别人解散房间的申请 同意还是拒绝
func HndlerApplyDissolveBack(userId uint32, agree bool, a gate.Agent) error {
	desk := MjroomManagerIns.GetFMjDeskBySession(userId) //只有朋友桌才能由此协议
	if desk == nil {
		return ERR_DESK_NOT_FOUND
	}
	//申请
	err := desk.ApplyDissolveBack(userId, agree)
	if err != nil {
		log.E("玩家agree %v 申请解散房间失败%v", agree, err)
		return err
	}
	return nil
}

//长沙麻将起手胡牌
func HandlerQiShouHu(userId uint32, chooseHu bool, a gate.Agent) error {
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //起手胡牌
	if desk == nil {
		//这里是否需要返回错误信息
		return ERR_DESK_NOT_FOUND
	}

	err := desk.ActQiShouHu(userId, chooseHu) //普通玩家开始起手胡牌
	if err != nil {
		log.E("服务器错误，胡牌失败..err[%v]", err)
		return err
	}

	return nil
}

//是否需要海底牌
func HandlerNeedHaidi(userId uint32, need bool, a gate.Agent) error {
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //海底牌
	if desk == nil {
		log.E("玩家%v need:%v 海底牌的时候出错...没有找到desk", userId, need)
		//这里是否需要返回错误信息
		return ERR_DESK_NOT_FOUND
	}

	err := desk.NeedHaidi(userId, need) //普通玩家开始起手胡牌
	if err != nil {
		log.E("服务器错误，胡牌失败..err[%v]", err)
		return err
	}

	return nil
}

//断线重连的处理
func HandlerReconnect(userId uint32, a gate.Agent) error {
	desk := MjroomManagerIns.GetMjDeskBySession(userId) //重连
	if desk == nil {
		log.T("玩家%v没有在游戏状态，断线重连之后不需要做处理", userId)
		return nil
	}

	//设置断线处理
	desk.SetReconnectStatus(userId, a)
	return nil
}
