package internal

import (
	"reflect"
	"casino_common/proto/ddproto"
	"casino_common/common/log"
	"casino_majiang/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_majianagv2/core/data"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/funcsInit"
	"casino_common/common/consts"
	"casino_common/common/Error"
	"casino_majiang/service/majiang"
	"casino_common/common/userService"
	"casino_majiang/gamedata/dao"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&mjproto.Game_CreateRoom{}, handlerCreateDesk)                    //创建房间
	handler(&mjproto.Game_EnterRoom{}, handlerGame_EnterRoom)                 //进入房间
	handler(&mjproto.Game_DissolveDesk{}, handlerDissolveDesk)                //解散房间
	handler(&mjproto.Game_Ready{}, handlerGame_Ready)                         //准备
	handler(&mjproto.Game_DingQue{}, handlerGame_DingQue)                     //定缺
	handler(&mjproto.Game_ExchangeCards{}, handlerGame_ExchangeCards)         //换3张
	handler(&mjproto.Game_SendOutCard{}, handlerGame_SendOutCard)             //出牌
	handler(&mjproto.Game_ActPeng{}, handlerGame_ActPeng)                     //碰
	handler(&mjproto.Game_ActGang{}, handlerGame_ActGang)                     //杠
	handler(&mjproto.Game_ActGuo{}, handlerGame_ActGuo)                       //过
	handler(&mjproto.Game_ActHu{}, handlerGame_ActHu)                         //胡
	handler(&mjproto.Game_GameRecord{}, handlerGame_GameRecord)               //战绩相关
	handler(&ddproto.CommonReqMessage{}, handlerGame_Message)                 //聊天
	handler(&ddproto.CommonAckLogout{}, HandlerCommonAckLogout)               //退出游戏
	handler(&ddproto.CommonReqLeaveDesk{}, HandlerCommonReqLeaveDesk)         //退出游戏
	handler(&ddproto.CommonReqApplyDissolve{}, handlerCommonReqApplyDissolve) //申请解散房间
	handler(&ddproto.CommonReqApplyDissolveBack{}, handlerApplyDissolveBack)  //申请解散房间
	handler(&ddproto.CommonReqEnterAgentMode{}, handlerEnterAgentMode)        //申请进入托管
	handler(&ddproto.CommonReqQuitAgentMode{}, handlerQuitAgentMode)          //申请退出托管
	handler(&mjproto.Game_ActChi{}, handlerActChi)                            //吃牌
	handler(&mjproto.Game_ActChangShaQiShouHu{}, handlerQiShouHu)             //长沙麻将起手胡牌
	handler(&mjproto.Game_ReqDealHaiDiCards{}, handlerNeedHaidi)              //询问是否需要海底牌

}

//创建一个房间,请求创建房间必定是朋友桌的请求，金币场没有创建房间
func handlerCreateDesk(args []interface{}) {
	m := args[0].(*mjproto.Game_CreateRoom) //创建房间时候的配置
	a := args[1].(gate.Agent)               //连接

	room := roomMgr.GetRoom(int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND), 0)
	if room == nil {
		log.E("创建朋友桌房间失败")
		//ack
		return
	}

	var UserCountLimit int32 = 4 //人数
	var FangCountLimit int32 = 3 //麻将房数

	if m.GetRoomTypeInfo().GetMjRoomType() == mjproto.MJRoomType_roomType_sanRenLiangFang {
		UserCountLimit = 3
		FangCountLimit = 2
	}

	if m.GetRoomTypeInfo().GetMjRoomType() == mjproto.MJRoomType_roomType_siRenLiangFang {
		UserCountLimit = 4
		FangCountLimit = 2
	}

	if m.GetRoomTypeInfo().GetMjRoomType() == mjproto.MJRoomType_roomType_liangRenLiangFang {
		UserCountLimit = 2
		FangCountLimit = 2
	}

	if m.GetRoomTypeInfo().GetMjRoomType() == mjproto.MJRoomType_roomType_sanRenSanFang {
		UserCountLimit = 3
		FangCountLimit = 3
	}

	config := &data.SkeletonMJConfig{
		Owner:            m.GetHeader().GetUserId(),
		RoomType:         majiang.ROOMTYPE_FRIEND,
		Status:           majiang.MJDESK_STATUS_READY,
		MjRoomType:       int32(m.GetRoomTypeInfo().GetMjRoomType()),
		BoardsCout:       m.GetRoomTypeInfo().GetBoardsCout(), //局数，如：4局（房卡 × 2）、8局（房卡 × 3）
		CreateFee:        room.CalcCreateFee(m.GetRoomTypeInfo().GetBoardsCout()),
		CapMax:           m.GetRoomTypeInfo().GetCapMax(),
		CardsNum:         m.GetRoomTypeInfo().GetCardsNum(),
		Settlement:       m.GetRoomTypeInfo().GetSettlement(),
		BaseValue:        m.GetRoomTypeInfo().GetBaseValue(),
		ZiMoRadio:        m.GetRoomTypeInfo().GetPlayOptions().GetZiMoRadio(),
		OthersCheckBox:   m.GetRoomTypeInfo().GetPlayOptions().GetOthersCheckBox(),
		HuRadio:          m.GetRoomTypeInfo().GetPlayOptions().GetHuRadio(),
		DianGangHuaRadio: m.GetRoomTypeInfo().GetPlayOptions().GetDianGangHuaRadio(),
		MJPaiCursor:      0,
		TotalPlayCount:   m.GetRoomTypeInfo().GetBoardsCout(),
		CurrPlayCount:    0,
		Banker:           0,
		NextBanker:       0,
		ActiveUser:       0,
		ActUser:          0,
		ActType:          0,
		NInitActionTime:  30,
		RoomLevel:        0,
		PlayerCountLimit: UserCountLimit,
		FangCount:        FangCountLimit,
		XueLiuChengHe:    m.GetRoomTypeInfo().GetMjRoomType() == mjproto.MJRoomType_roomType_xueLiuChengHe, //是否是血流成河
	}

	desk, err := room.CreateDesk(config)
	if desk == nil {
		result := newProto.NewGame_AckCreateRoom()
		*result.Header.Code = consts.ACK_RESULT_ERROR
		*result.Header.Error = Error.GetErrorMsg(err)
		log.Error("用户[%v]创建房间失败...err[%v]", m.GetHeader().GetUserId(), err)
		a.WriteMsg(result)
	} else {
		//创建desk成功...设置为开始准备的状态
		//desk.SetStatus(MJDESK_STATUS_READY) //设置为开始准备的状态
		log.T("用户[%v]创建房间成功，roomKey[%v]", desk.GetMJConfig().Owner, desk.GetMJConfig().Password)
		result := newProto.NewGame_AckCreateRoom()
		*result.Header.Code = consts.ACK_RESULT_SUCC
		*result.Password = desk.GetMJConfig().Password
		*result.DeskId = desk.GetMJConfig().DeskId
		*result.CreateFee = desk.GetMJConfig().CreateFee
		result.RoomTypeInfo = &mjproto.RoomTypeInfo{
			MjRoomType: mjproto.MJRoomType(desk.GetMJConfig().MjRoomType).Enum(),
			BoardsCout: proto.Int32(desk.GetMJConfig().BoardsCout),
			CapMax:     proto.Int64(desk.GetMJConfig().CapMax),
			PlayOptions: &mjproto.PlayOptions{
				ZiMoRadio:        proto.Int32(desk.GetMJConfig().ZiMoRadio),
				DianGangHuaRadio: proto.Int32(desk.GetMJConfig().DianGangHuaRadio),
				OthersCheckBox:   desk.GetMJConfig().OthersCheckBox,
				HuRadio:          proto.Int32(desk.GetMJConfig().HuRadio),
			},
			CardsNum:   proto.Int32(desk.GetMJConfig().CardsNum),
			Settlement: proto.Int32(desk.GetMJConfig().Settlement),
			BaseValue:  proto.Int64(desk.GetMJConfig().BaseValue),
			ChangShaPlayOptions: &mjproto.ChangShaPlayOptions{
				//PlayerCount:
				//IgnoreBank:
				//BirdCount:
				//BirdMultiple:
			},
		}
		*result.UserBalance = userService.GetUserRoomCard(m.GetHeader().GetUserId())
		a.WriteMsg(result)

		//创建成功之后，用户自动进入房间...
		err := desk.EnterUser(m.GetHeader().GetUserId(), a)
		if err != nil {
			log.E("进入房间失败...")
		}
	}
}

func handlerGame_EnterRoom(args []interface{}) {
	m := args[0].(*mjproto.Game_EnterRoom) //创建房间时候的配置
	a := args[1].(gate.Agent)              //连接

	userId := m.GetHeader().GetUserId()
	passWord := m.GetPassWord()
	//1,找到room
	room := roomMgr.GetRoom(m.GetRoomType(), m.GetRoomLevel())
	if room == nil {
		//没有找到room，进入房间失败
		log.T("用户[%v]进入房间失败，没有找到对应的room", userId)
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = consts.ACK_RESULT_ERROR
		*ack.Header.Error = "房间号输入错误"
		a.WriteMsg(ack)
	}

	//2,进入房间
	err := room.EnterUser(userId, passWord, a)
	if err != nil {
		//进入房间失败
		ack := newProto.NewGame_AckEnterRoom()
		*ack.Header.Code = Error.GetErrorCode(err)
		*ack.Header.Error = Error.GetErrorMsg(err)
		a.WriteMsg(ack)
	}
}

func handlerDissolveDesk(args []interface{}) {
	desk := roomMgr.GetDesk()
	err := desk.Leave(0)
	if err != nil {
		log.E("离开失败")
		//return errenterdesk
	}
}

//准备的协议
func handlerGame_Ready(args []interface{}) {
	m := args[0].(*mjproto.Game_Ready)
	a := args[1].(gate.Agent)

	userId := m.GetUserId()
	desk := roomMgr.GetMjDeskBySession(userId)
	if desk == nil {
		// 准备失败
		log.E("用户[%v]准备失败.因为没有找到对应的desk", userId)
		result := newProto.NewGame_AckReady()
		*result.Header.Code = consts.ACK_RESULT_SUCC
		*result.Header.Error = ""
		result.UserId = proto.Uint32(userId)
		a.WriteMsg(result)
		return
	}

	//开始准备
	err := desk.Ready(userId)
	if err != nil {
		log.E("用户[%v]准备失败.err %v", userId, err)
		result := newProto.NewGame_AckReady()
		*result.Header.Code = consts.ACK_RESULT_SUCC
		*result.Header.Error = ""
		result.UserId = proto.Uint32(userId)
		a.WriteMsg(result)
		return
	}
}

func handlerGame_DingQue(args []interface{}) {
	m := args[0].(*mjproto.Game_DingQue)
	userId := m.GetUserId()
	color := m.GetColor()
	desk := roomMgr.GetMjDeskBySession(userId)
	if desk == nil {
		log.E("玩家[%v]定缺[%v]失败,desk 没有找到..", userId, color)
		return
	}
	err := desk.DingQue(userId, color) //普通玩家定缺
	if err != nil {
		log.E("用户[%v]定缺失败...", userId)
	}
}

//换三张
func handlerGame_ExchangeCards(args []interface{}) {

}

//打牌
func handlerGame_SendOutCard(args []interface{}) {
	m := args[0].(*mjproto.Game_SendOutCard)
	a := args[1].(gate.Agent)
	userId := m.GetHeader().GetUserId()
	cardId := m.GetCardId()

	//检测参数
	desk := roomMgr.GetMjDeskBySession(userId)
	if desk == nil {
		//打牌失败，因为没有找到对应的麻将桌子
		log.E("用户[%v]打牌", userId)
		return
	}

	//开始打牌
	err := desk.ActOut(userId, cardId, false) //普通玩家打牌
	if err != nil {
		//打牌失败
		result := newProto.NewGame_AckSendOutCard()
		*result.Header.Code = Error.GetErrorCode(err)
		*result.Header.Error = Error.GetErrorMsg(err)
		a.WriteMsg(result)
		return
	}
}

//碰牌 todo
func handlerGame_ActPeng(args []interface{}) {
	m := args[0].(*mjproto.Game_ActPeng)
	a := args[1].(gate.Agent)
	userId := m.GetUserId()

	//找到桌子
	desk := roomMgr.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk
	if desk == nil {
		//这里属于服务器错误... 是否需要给客户端返回信息？
		log.E("没有找到对应的desk ..")
		result := &mjproto.Game_AckActPeng{}
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

//杠牌 todo
func handlerGame_ActGang(args []interface{}) {
	m := args[0].(*mjproto.Game_ActGang)

	userId := m.GetUserId()
	cardId := m.GetGangCard().GetId()
	bu := m.GetBu()

	result := &mjproto.Game_AckActGang{}
	result.Header = newProto.SuccessHeader()

	desk := roomMgr.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk
	//通过userId 的session 得到对应的desk
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

//过牌 todo
func handlerGame_ActGuo(args []interface{}) {
	m := args[0].(*mjproto.Game_ActGuo)
	a := args[1].(gate.Agent)
	userId := m.GetUserId()
	desk := roomMgr.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk
	if desk == nil {
		log.E("玩家%v过的时候失败,没有找到对应的desk", userId)
		return
	}
	err := desk.ActGuo(userId)
	if err != nil {
		result := newProto.NewGame_AckActGuo()
		result.Header = newProto.ErrorHeader()
		a.WriteMsg(result)
		return
	}
}

//胡牌 todo
func handlerGame_ActHu(args []interface{}) {
	m := args[0].(*mjproto.Game_ActHu)
	userId := m.GetUserId()

	//区分自摸点炮:1,如果自己的手牌就已经糊了（或者如果自己自己的牌是14，11，8，5，2 张的时候），那么就自摸，如果需要加上判定牌，那就是点炮
	desk := roomMgr.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk

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

//游戏记录 todo
func handlerGame_GameRecord(args []interface{}) {
	m := args[0].(*mjproto.Game_GameRecord)
	a := args[1].(gate.Agent)
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
}

//发信息 todo
func handlerGame_Message(args []interface{}) {
	m := args[0].(*ddproto.CommonReqMessage)
	log.T("请求发送信息[%v]", m)
	userId := m.GetHeader().GetUserId()
	desk := roomMgr.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk
	if desk != nil {
		desk.SendMessage(m)
	} else {
		log.E("玩家发送消息失败....因为没有找到玩家【%v】所在的desk", userId)
	}
}

//强制退出
func HandlerCommonAckLogout(args []interface{}) {
	m := args[0].(*ddproto.CommonReqLeaveDesk)
	a := args[1].(gate.Agent)
	userId := m.GetHeader().GetUserId()

	desk := roomMgr.GetMjDeskBySession(userId) //通过userId 的session 得到对应的desk
	if desk == nil {
		return
	}

	if m.GetIsExchange() {
		//换房间
		desk.ExchangeRoom(userId, a)
	} else {
		//离开房间
		log.T("玩家[%v]开始离开房间...", userId)
		err := desk.Leave(userId)
		if err != nil {
			log.E("玩家[%v]离开房间的时候出错", userId)
		}
	}
}

//离开房间
func HandlerCommonReqLeaveDesk(args []interface{}) {

}

//申请解散房间
func handlerCommonReqApplyDissolve(args []interface{}) {

}

//申请解散 回复
func handlerApplyDissolveBack(args []interface{}) {

}

//托管模式
func handlerEnterAgentMode(args []interface{}) {

}

//拖出托管模式
func handlerQuitAgentMode(args []interface{}) {

}

//吃
func handlerActChi(args []interface{}) {
}

//长沙麻将起手胡牌
func handlerQiShouHu(args []interface{}) {
}

func handlerNeedHaidi(args []interface{}) {
	m := args[0].(*mjproto.Game_ReqDealHaiDiCards)
	//a := args[1].(gate.Agent)

	userId := m.GetUserId()
	need := m.GetNeed()

	desk := roomMgr.GetMjDeskBySession(userId)
	if desk == nil {
		log.E("玩家%v need:%v 海底牌的时候出错...没有找到desk", userId, need)
		//这里是否需要返回错误信息
	}

	err := desk.NeedHaidi(userId, need) //普通玩家开始起手胡牌
	if err != nil {
		log.E("服务器错误，胡牌失败..err[%v]", err)
	}
}
