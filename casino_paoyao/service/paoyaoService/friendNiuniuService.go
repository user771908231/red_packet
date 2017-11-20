package paoyaoService

import (
	"casino_common/proto/ddproto"
	"github.com/name5566/leaf/gate"
	"casino_common/proto/funcsInit"
	"casino_paoyao/service/paoyao"
	"github.com/golang/protobuf/proto"
	"casino_common/common/userService"
	"casino_common/common/log"
	"fmt"
	"casino_common/common/sessionService"
	"time"
	"casino_common/common/service/roomAgent"
	"casino_common/common/Error"
)

//创建房间
func CreateDeskHandler(req *ddproto.PaoyaoCreateDeskReq, agent gate.Agent) *ddproto.PaoyaoEnterDeskAck {
	msg := &ddproto.PaoyaoEnterDeskAck{
		Header: commonNewPorot.NewHeader(),
		DeskState: nil,
		IsReconnect:proto.Bool(false),
	}
	defer agent.WriteMsg(msg)

	user, err := paoyao.FindUserById(req.Header.GetUserId())
	if err == nil {
		*msg.Header.Code = -5
		*msg.Header.Error = fmt.Sprintf("您还在刨幺游戏中（房号%s DeskId%d），暂时无法创建房间！", user.Desk.GetPwd(), user.Desk.GetDeskId())
		return msg
	}

	//是否维护中
	isOnMaintain, maintainMsg := sessionService.IsOnMaintain(int32(ddproto.CommonEnumGame_GID_PAOYAO))
	if isOnMaintain == true {
		log.E("用户%d创建房间失败，原因：%s", req.Header.GetUserId(), maintainMsg)
		*msg.Header.Code = -6
		*msg.Header.Error = maintainMsg
		return msg
	}

	//判断是否在其他游戏中
	_,err = sessionService.CanEnter(req.Header.GetUserId(), int32(ddproto.CommonEnumGame_GID_PAOYAO), int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND))
	if err != nil {
		*msg.Header.Code = -7
		*msg.Header.Error = err.Error()
		agent.WriteMsg(msg)
		return msg
	}

	//房费
	var ownerFee int64 = 0
	//朋友桌扣除房卡
	roomCard := userService.GetUserRoomCard(req.Header.GetUserId())
	ownerFee = int64(paoyao.GetOwnerFee(req.GetOption().GetMaxCircle()))
	if roomCard < ownerFee {
		log.E("用户%d创建房间失败，原因：房卡不足", req.Header.GetUserId())
		*msg.Header.Code = -4
		*msg.Header.Error = "房卡余额不足！"
		return msg
	}else {
		userService.DECRUserRoomcard(req.Header.GetUserId(), ownerFee, int32(ddproto.CommonEnumGame_GID_PAOYAO), "刨幺朋友桌，创建房间扣除房卡")
	}


	room, err := paoyao.Rooms.GetRoomById(0)
	if err != nil {
		log.E("用户%d创建房间失败，原因：%s", req.Header.GetUserId(), err.Error())
		*msg.Header.Code = -1
		*msg.Header.Error = "创建房间失败"
		return msg
	}
	desk, err := room.CreateFriendDesk(req.Option, req.Header.GetUserId())
	if err != nil {
		log.E("用户%d创建房间失败，原因：%s", req.Header.GetUserId(), err.Error())
		*msg.Header.Code = -2
		*msg.Header.Error = "创建房间失败"
		return msg
	}
	//如果是朋友桌代开
	if req.GetIsDaikai() == true {
		daikai_user := req.Header.GetUserId()
		*desk.IsDaikai = true
		*desk.DaikaiUser = daikai_user
		//更新代开列表
		err := roomAgent.CreateDesk(int32(ddproto.CommonEnumGame_GID_PAOYAO), desk.GetPwd(), desk.GetDeskId(), daikai_user, desk.GetTips(), time.Now().Unix(), desk.DeskOption.GetMaxCircle(), desk.DeskOption.GetMaxUser(), 0)
		if err != nil {
			//返还房费
			userService.INCRUserRoomcard(daikai_user, ownerFee, int32(ddproto.CommonEnumGame_GID_PAOYAO), "刨幺朋友桌，代开失败房费返还")
			*msg.Header.Code = -8
			*msg.Header.Error = err.Error()
			return msg
		}
		//msg.DeskState = desk.GetClientDesk()
		log.T("用户%d代开房间%s成功！", daikai_user, desk.GetPwd())
		*msg.Header.Code = 2
		*msg.Header.Error = "代开房间成功！"
		return msg
	}
	_,err = desk.AddUser(req.Header.GetUserId(), agent)
	if err != nil {
		log.E("用户%d创建房间失败，原因：%s", req.Header.GetUserId(), err.Error())
		*msg.Header.Code = -3
		*msg.Header.Error = "创建房间失败"
		return msg
	}
	msg.DeskState = desk.GetClientDesk()
	*msg.Header.Code = 1
	return msg
}

//进入房间
func EnterDeskHandler(req *ddproto.PaoyaoEnterDeskReq, agent gate.Agent) {
	msg := &ddproto.PaoyaoEnterDeskAck{
		Header: commonNewPorot.NewHeader(),
		DeskState: nil,
		IsReconnect: proto.Bool(false),
	}

	//处理断  线重连
	user, err := paoyao.FindUserById(req.Header.GetUserId())
	if err == nil {
		log.T("用户%d断线重连成功。", user.GetUserId())
		*msg.Header.Code = 3
		*msg.Header.Error = "断线重连成功！"
		*msg.IsReconnect = true
		//切换为在线状态
		*user.IsOnline = true
		msg.DeskState = user.GetClientDesk()
		//更新断线用户的agent
		user.Agent = agent
		agent.WriteMsg(msg)
		//断线重连后的重发overturn等操作
		ReconnectProcess(user)
		return
	}else {
		//如果没找到用户，则强行删除该用户在刨幺的session
		session := sessionService.GetSessionAuto(req.Header.GetUserId())
		if session != nil && session.GetGameId() == int32(ddproto.CommonEnumGame_GID_PAOYAO) {
			sessionService.DelSessionByKey(session.GetUserId(),session.GetRoomType(),session.GetGameId(),session.GetDeskId())
		}
	}

	//判断是否在其他游戏中
	roomType := ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND
	if req.GetRoomId() > 0 {
		roomType = ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN
	}
	_,err = sessionService.CanEnter(req.Header.GetUserId(), int32(ddproto.CommonEnumGame_GID_PAOYAO), int32(roomType))
	if err != nil {
		*msg.Header.Code = -4
		*msg.Header.Error = err.Error()
		agent.WriteMsg(msg)
		return
	}

	//先找到room
	var room *paoyao.Room
	var desk *paoyao.Desk
	room, err = paoyao.Rooms.GetRoomById(req.GetRoomId())
	if err != nil {
		log.E("用户%d进房pwd:%v roomId:%v失败，原因：%s", req.Header.GetUserId(), req.GetDeskPwd(), req.GetRoomId(), err.Error())
		*msg.Header.Code = -1
		*msg.Header.Error = "进入房间失败！"
		agent.WriteMsg(msg)
		return
	}

	//朋友桌进牌桌
	if req.GetRoomId() == 0 {
		desk, err = room.GetDeskByPassword(req.GetDeskPwd())
		//清理session,以便容错
		if len(req.GetDeskPwd()) == 6 && req.GetDeskPwd()[2] == '2' && req.GetDeskPwd()[3] == '5' {
			go func() {
				defer Error.ErrorRecovery("EnterDeskHandler->DelSessionByKey()")
				sessionService.DelSessionByKey(req.Header.GetUserId(), int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND), int32(ddproto.CommonEnumGame_GID_PAOYAO), 0)
			}()
		}
		if err != nil {
			log.E("用户%d进房%v失败，原因：%s", req.Header.GetUserId(), req.GetDeskPwd(), err.Error())
			*msg.Header.Code = -2
			*msg.Header.Error = "进入房间失败！未找到该牌桌！"
			agent.WriteMsg(msg)
			return
		}
	}else {
		desk, err = room.GetFreeCoinDesk()
		if err != nil {
			log.E("用户%d进房pwd:%v roomId:%v失败，原因：%s", req.Header.GetUserId(), req.GetDeskPwd(), req.GetRoomId(), err.Error())
			*msg.Header.Code = -2
			*msg.Header.Error = "进入房间失败！未找到该房间！"
			agent.WriteMsg(msg)
			return
		}
	}

	_,err = desk.AddUser(req.Header.GetUserId(), agent)
	if err != nil {
		log.E("用户%d进房pwd:%v roomId:%v失败，原因：%s", req.Header.GetUserId(), req.GetDeskPwd(), req.GetRoomId(), err.Error())
		*msg.Header.Code = -3
		*msg.Header.Error = err.Error()
		agent.WriteMsg(msg)
		return
	}

	log.E("用户%d进房pwd:%v roomId:%v成功。", req.Header.GetUserId(), req.GetDeskPwd(), req.GetRoomId())
	msg.DeskState = desk.GetClientDesk()
	*msg.Header.Code = 1
	*msg.Header.Error = "新玩家加入房间成功！"
	agent.WriteMsg(msg)

	//真人加入后调度机器人加入
	if desk.GetIsCoinRoom() {
		desk.AutoJoinRobot()
	}
}

//断线重连处理
func ReconnectProcess(user *paoyao.User) {
	//更新离开状态
	user.IsLeave = proto.Bool(false)
	//刷新白名单
	user.CheckWhiteList()
	//更新离线状态广播
	user.SendOffineBc()
	//由客户端自己判断状态
	if user.Desk.GetIsOnDissolve() && user.GetDissolveState() == 0 {
		//如果牌桌正在解散牌桌中,且未投票,则给他发送发起解散房间广播
		msg := &ddproto.CommonBcApplyDissolve{
			Header: commonNewPorot.NewHeader(),
			UserId: user.Desk.DissolveUser,
		}
		time.AfterFunc(500 * time.Millisecond, func() {
			user.WriteMsg(msg)
			log.T("重发解散申请：%v", msg, user.Desk.DissolveUser)
		})
	}

	//金币场开局
	if user.Desk.GetIsCoinRoom() {
		user.Desk.DoStart()
	}
}

//准备
func ReadyHandler(req *ddproto.PaoyaoSwitchReadyReq, agent gate.Agent) {
	user, err := paoyao.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		if !user.GetIsCoinRoom() {
			//朋友桌准备
			user.DoReadyFriend()
		}else {
			//金币场准备
			user.DoReadyCoin()
		}
	}else {
		user = &paoyao.User{
			Agent: agent,
			PaoyaoSrvUser: nil,
		}
		user.SendReadyAck(-1, "您当前未在房间中")
	}
}

//申请解散房间
func ApplyDissolveReqHandler(req *ddproto.CommonReqApplyDissolve, agent gate.Agent)  {
	user, err := paoyao.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		user.DoApplyDissolve()
	}
}

//确定、取消解散房间
func DissolveBackReqHandler(req *ddproto.CommonReqApplyDissolveBack, agent gate.Agent) {
	user, err := paoyao.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		user.DoDissolveBack(req.GetAgree())
	}
}

//聊天请求
func MessageReqHandler(req *ddproto.CommonReqMessage, agent gate.Agent) {
	user, err := paoyao.FindUserById(req.Header.GetUserId())
	if err != nil {
		return
	}
	user.Desk.ReqLock.Lock()
	defer user.Desk.ReqLock.Unlock()
	user.UpdateAgent(agent)

	//将聊天请求广播出去
	msg := &ddproto.CommonBcMessage{
		UserId: req.Header.UserId,
		Id: req.Id,
		Msg: req.Msg,
		MsgType: req.MsgType,
		ToUserId: req.ToUserId,
	}
	//广播
	user.Desk.BroadCast(msg)
}

//离开房间
func LeaveDeskReqHandler(req *ddproto.CommonReqLeaveDesk, agent gate.Agent) {
	user, err := paoyao.FindUserById(req.Header.GetUserId())
	if err != nil {
		return
	}
	user.Desk.ReqLock.Lock()
	defer user.Desk.ReqLock.Unlock()
	user.UpdateAgent(agent)

	//离开房间
	user.DoLeaveDesk()
}

//金币场房间列表
func CoinRoomListHandler(req *ddproto.PaoyaoCoinRoomListReq, agent gate.Agent) {
	ack := &ddproto.PaoyaoCoinRoomListAck{
		Header: commonNewPorot.NewHeader(),
		Rooms: []*ddproto.PaoyaoSrvRoom{},
	}

	for _,room := range paoyao.Rooms{
		if room.GetRoomType() == int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN) {
			ack.Rooms = append(ack.Rooms, room.PaoyaoSrvRoom)
		}
	}

	agent.WriteMsg(ack)
}
