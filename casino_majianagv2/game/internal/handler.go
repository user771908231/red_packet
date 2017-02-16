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
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&ddproto.DdzReqCreateDesk{}, handlerCreateDesk)                   //创建房间
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

	//get mjconfig by m
	log.T("麻将的配置 %v ", m)
	config := data.SkeletonMJConfig{
		Owner: m.GetHeader().GetUserId(),
		//Password         string
		//DeskId           int32
		RoomType:   majiang.ROOMTYPE_FRIEND,
		Status:     majiang.MJDESK_STATUS_READY,
		MjRoomType: int32(m.GetRoomTypeInfo().GetMjRoomType()),
		//RoomId           int32
		//CreateFee        int64
		//BoardsCout       int32 //局数，如：4局（房卡 × 2）、8局（房卡 × 3）
		//CapMax           int32
		//CardsNum         int32
		//Settlement       int32
		//BaseValue        int64
		//ZiMoRadio        int32
		//OthersCheckBox   []int32
		//HuRadio          int32
		//DianGangHuaRadio int32
		//MJPaiCursor      int32
		//TotalPlayCount   int32
		//CurrPlayCount    int32
		//Banker           uint32
		//NextBanker       uint32
		//CheckCase        *CheckCase
		//ActiveUser       uint32
		//GameNumber       int32
		//ActUser          uint32
		//ActType          int32
		//NInitActionTime  int32
		//RoomLevel        int32
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
		room.EnterUser(m.GetHeader().GetUserId(), desk.GetMJConfig().Password)
	}
}

func handlerGame_EnterRoom(args []interface{}) {
	m := args[0].(*mjproto.Game_EnterRoom) //创建房间时候的配置
	a := args[1].(gate.Agent)              //连接

	//1,找到room
	room := roomMgr.GetRoom(m.GetRoomType(), m.GetRoomLevel())
	if room == nil {
		log.E("room没有找到...")
		a.WriteMsg(nil)
	}

	//2,进入房间
	err := room.EnterUser(m.GetUserId(), m.GetPassWord())
	if err != nil {
		a.WriteMsg(nil)
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
	log.T("收到请求，game_Ready([%v])", userId)
	desk := roomMgr.GetMjDeskBySession(userId)
	if desk == nil {
		// 准备失败
		log.E("用户[%v]准备失败.因为没有找到对应的desk", userId)
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
		*result.Header.Code = consts.ACK_RESULT_ERROR
		*result.Header.Error = Error.GetErrorMsg(err)
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
	err := desk.ActOut(userId, cardId) //普通玩家打牌
	if err != nil { //打牌失败
		result := newProto.NewGame_AckSendOutCard()
		*result.Header.Code = Error.GetErrorCode(err)
		*result.Header.Error = Error.GetErrorMsg(err)
		a.WriteMsg(result)
		return
	}
}

//碰牌
func handlerGame_ActPeng(args []interface{}) {
}

//杠牌
func handlerGame_ActGang(args []interface{}) {
}

//过牌
func handlerGame_ActGuo(args []interface{}) {
}

//胡牌
func handlerGame_ActHu(args []interface{}) {
}

//游戏记录
func handlerGame_GameRecord(args []interface{}) {
}

//发信息
func handlerGame_Message(args []interface{}) {
}

//强制退出
func HandlerCommonAckLogout(args []interface{}) {
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
