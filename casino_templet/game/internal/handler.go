package internal

import (
	"reflect"
	"casino_common/proto/ddproto"
	"casino_common/common/log"
	"casino_majiang/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_templet/core/data"
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
	config := data.DeskConfig{}
	err, _ := room.CreateDesk(config)
	if err != nil {
		log.E("创建房间失败 err %v", err)
		//ack
		a.WriteMsg(nil)
		return
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

func handlerGame_Ready(args []interface{}) {

}

func handlerGame_DingQue(args []interface{}) {

}
func handlerGame_ExchangeCards(args []interface{}) {
}

func handlerGame_SendOutCard(args []interface{}) {
}

func handlerGame_ActPeng(args []interface{}) {
}
func handlerGame_ActGang(args []interface{}) {
}
func handlerGame_ActGuo(args []interface{}) {
}

func handlerGame_ActHu(args []interface{}) {
}

func handlerGame_GameRecord(args []interface{}) {
}

func handlerGame_Message(args []interface{}) {
}

func HandlerCommonAckLogout(args []interface{}) {
}
func HandlerCommonReqLeaveDesk(args []interface{}) {
}
func handlerCommonReqApplyDissolve(args []interface{}) {
}
func handlerApplyDissolveBack(args []interface{}) {
}
func handlerEnterAgentMode(args []interface{}) {
}
func handlerQuitAgentMode(args []interface{}) {
}
