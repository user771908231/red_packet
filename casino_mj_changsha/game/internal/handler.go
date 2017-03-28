package internal

import (
	"reflect"
	"github.com/name5566/leaf/gate"
	mjProto "casino_mj_changsha/msg/protogo"
	MJService "casino_mj_changsha/service/majiang"
	"casino_common/proto/ddproto"
	"casino_common/common/log"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&mjProto.Game_CreateRoom{}, handlerGame_CreateDesk)
	handler(&mjProto.Game_EnterRoom{}, handlerGame_EnterRoom)
	handler(&mjProto.Game_DissolveDesk{}, handlerDissolveDesk)
	handler(&mjProto.Game_Ready{}, handlerGame_Ready)
	handler(&mjProto.Game_SendOutCard{}, handlerGame_SendOutCard)             //出牌
	handler(&mjProto.Game_ActPeng{}, handlerGame_ActPeng)                     //碰
	handler(&mjProto.Game_ActGang{}, handlerGame_ActGang)                     //杠
	handler(&mjProto.Game_ActGuo{}, handlerGame_ActGuo)                       //过
	handler(&mjProto.Game_ActHu{}, handlerGame_ActHu)                         //胡
	handler(&mjProto.Game_GameRecord{}, handlerGame_GameRecord)               //战绩相关
	handler(&ddproto.CommonReqMessage{}, handlerGame_Message)                 //聊天
	handler(&ddproto.CommonAckLogout{}, HandlerCommonAckLogout)               //退出游戏
	handler(&ddproto.CommonReqLeaveDesk{}, HandlerCommonReqLeaveDesk)         //退出游戏
	handler(&ddproto.CommonReqApplyDissolve{}, handlerCommonReqApplyDissolve) //申请解散房间
	handler(&ddproto.CommonReqApplyDissolveBack{}, handlerApplyDissolveBack)  //申请解散房间
	handler(&mjProto.Game_ActChi{}, handlerActChi)                            //吃牌
	handler(&mjProto.Game_ActChangShaQiShouHu{}, handlerQiShouHu)             //长沙麻将起手胡牌
	handler(&mjProto.Game_ReqDealHaiDiCards{}, handlerNeedHaidi)              //询问是否需要海底牌
	handler(&ddproto.CommonReqReconnect{}, handlerReconnect)                  //断线重连的处理

}

//处理创建房间
func handlerGame_CreateDesk(args []interface{}) {
	m := args[0].(*mjProto.Game_CreateRoom)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_CreateDesk(m, a)
}

//处理进入房间
func handlerGame_EnterRoom(args []interface{}) {
	m := args[0].(*mjProto.Game_EnterRoom)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_EnterDesk(m.GetHeader().GetUserId(), m.GetPassWord(), m.GetRoomType(), m.GetRoomLevel(), m.GetEnterType(), a)
}

//解散房间
func handlerDissolveDesk(args []interface{}) {
	m := args[0].(*mjProto.Game_DissolveDesk) //解散房间
	MJService.HandlerDissolveDesk(m.GetHeader().GetUserId())
}

//准备游戏
func handlerGame_Ready(args []interface{}) {
	m := args[0].(*mjProto.Game_Ready)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_Ready(m.GetHeader().GetUserId(), a)
}

//出牌
func handlerGame_SendOutCard(args []interface{}) {
	m := args[0].(*mjProto.Game_SendOutCard)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_SendOutCard(m.GetHeader().GetUserId(), m.GetCardId(), a)
}

//碰
func handlerGame_ActPeng(args []interface{}) {
	m := args[0].(*mjProto.Game_ActPeng)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_ActPeng(m.GetUserId(), a)
}

//杠
func handlerGame_ActGang(args []interface{}) {
	m := args[0].(*mjProto.Game_ActGang)
	MJService.HandlerGame_ActGang(m.GetUserId(), m.GetGangCard().GetId(), m.GetBu())
}

//过
func handlerGame_ActGuo(args []interface{}) {
	m := args[0].(*mjProto.Game_ActGuo)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_ActGuo(m.GetUserId(), a)
}

//胡
func handlerGame_ActHu(args []interface{}) {
	m := args[0].(*mjProto.Game_ActHu)
	MJService.HandlerGame_ActHu(m.GetUserId())
}

//吃
func handlerActChi(args []interface{}) {
	m := args[0].(*mjProto.Game_ActChi)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_ActChi(m.GetUserId(), m.GetChooseCards(), a)
}

//查询战绩
func handlerGame_GameRecord(args []interface{}) {
	m := args[0].(*mjProto.Game_GameRecord)
	a := args[1].(gate.Agent)
	MJService.HandlerGame_GameRecord(m, a)
}

//聊天的协议
func handlerGame_Message(args []interface{}) {
	m := args[0].(*ddproto.CommonReqMessage)
	MJService.HandlerGame_Message(m)
}

//推出的
func HandlerCommonAckLogout(args []interface{}) {

}

//离开房间的协议...
func HandlerCommonReqLeaveDesk(args []interface{}) {
	m := args[0].(*ddproto.CommonReqLeaveDesk)
	a := args[1].(gate.Agent)

	if m.GetIsExchange() {
		//换房间
		MJService.HandlerExchangeRoom(m.GetHeader().GetUserId(), a)
	} else {
		//离开房间
		err := MJService.HandlerLeaveRoom(m.GetHeader().GetUserId(), a)
		if err != nil {
			log.T("玩家退出晚间的时候失败(有可能是客户端已经被强制退出了，又请求了一个离开房间):err %v", err)
		}
	}
	//MJService.HandlerLeveav2(m.GetHeader().GetUserId(), m.GetIsExchange(), a)

}

//申请解散房间
func handlerCommonReqApplyDissolve(args []interface{}) {
	m := args[0].(*ddproto.CommonReqApplyDissolve)
	a := args[1].(gate.Agent)
	MJService.HandlerApplyDissolve(m.GetUserId(), a)
}

//回复别人解散房间的申请
func handlerApplyDissolveBack(args []interface{}) {
	m := args[0].(*ddproto.CommonReqApplyDissolveBack)
	a := args[1].(gate.Agent)
	MJService.HndlerApplyDissolveBack(m.GetUserId(), m.GetAgree(), a)
}

//长沙麻将起手胡牌
func handlerQiShouHu(args []interface{}) {
	m := args[0].(*mjProto.Game_ActChangShaQiShouHu)
	a := args[1].(gate.Agent)
	MJService.HandlerQiShouHu(m.GetUserId(), m.GetHu(), a)
}

func handlerNeedHaidi(args []interface{}) {
	m := args[0].(*mjProto.Game_ReqDealHaiDiCards)
	a := args[1].(gate.Agent)
	MJService.HandlerNeedHaidi(m.GetUserId(), m.GetNeed(), a)
}

//断线重连的处理
func handlerReconnect(args []interface{}) {
	m := args[0].(*ddproto.CommonReqReconnect)
	a := args[1].(gate.Agent)
	MJService.HandlerReconnect(m.GetUserId(), a)
}
