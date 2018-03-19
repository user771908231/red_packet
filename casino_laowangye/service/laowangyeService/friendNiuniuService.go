package laowangyeService

import (
	"casino_common/proto/ddproto"
	"github.com/name5566/leaf/gate"
	"casino_common/proto/funcsInit"
	"casino_laowangye/service/laowangye"
	"github.com/golang/protobuf/proto"
	"casino_common/common/userService"
	"casino_common/common/log"
	"fmt"
	"casino_common/common/sessionService"
	"time"
	"casino_common/common/service/roomAgent"
	"casino_common/common/Error"
	"casino_common/gameManager/roomService"
)

//创建房间
func CreateDeskHandler(req *ddproto.LwyCreateDeskReq, agent gate.Agent) *ddproto.LwyEnterDeskAck {
	msg := &ddproto.LwyEnterDeskAck{
		Header: commonNewPorot.NewHeader(),
		DeskState: nil,
	}

	if req.Option == nil {
		log.E("create req option is nil.")
		*msg.Header.Code = -1
		*msg.Header.Error = "create req option is nil."
		return msg
	}
	//初始化参数
	req.Option.MinUser = proto.Int32(2)
	req.Option.MaxUser = proto.Int32(6)
	req.Option.BaseScore = proto.Int64(0)
	req.Option.IsCoinRoom = proto.Bool(false)
	req.Option.DenyHalfJoin = proto.Bool(true)
	if req.Option.GetBoardsCout() < 4 {
		req.Option.BoardsCout = proto.Int32(10)
	}

	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err == nil {
		*msg.Header.Code = -5
		*msg.Header.Error = fmt.Sprintf("您还在老王爷游戏中（房号%s DeskId%d），暂时无法创建房间！", user.Desk.GetPassword(), user.Desk.GetDeskId())
		return msg
	}

	//是否维护中
	isOnMaintain, maintainMsg := sessionService.IsOnMaintain(int32(ddproto.CommonEnumGame_GID_LAOWANGYE))
	if isOnMaintain == true {
		log.E("用户%d创建房间失败，原因：%s", req.Header.GetUserId(), maintainMsg)
		*msg.Header.Code = -6
		*msg.Header.Error = maintainMsg
		return msg
	}

	//判断是否在其他游戏中
	_,err = sessionService.CanEnter(req.Header.GetUserId(), int32(ddproto.CommonEnumGame_GID_LAOWANGYE), int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND))
	if err != nil {
		*msg.Header.Code = -7
		*msg.Header.Error = err.Error()
		agent.WriteMsg(msg)
		return msg
	}

	//房费
	var ownerFee int64 = 0
	//扣除房卡
	if !req.GetOption().GetIsCoinRoom() {
		//朋友桌扣房卡
		err = roomService.DecUserRoomcardToCreateDesk(req.Option.GetRoomCardBillType(), ddproto.CommonEnumGame_GID_LAOWANGYE, req.Option.GetBoardsCout(), req.Option.GetMaxUser(), req.Option.GetChanelId(), req.Header.GetUserId())

		if err != nil {
			log.E("用户%d创建房间失败，原因：%v", req.Header.GetUserId(), err)
			*msg.Header.Code = -4
			*msg.Header.Error = err.Error()
			return msg
		}
	}else {
		//检测最小金币数
		//if userService.GetUserCoin(req.Header.GetUserId()) < req.Option.GetMinEnterScore() {
		//	err_str := fmt.Sprintf("金币小于%d，建房失败！", req.Option.GetMinEnterScore())
		//	log.E("用户%d创建房间失败，原因：%v", err_str)
		//	*msg.Header.Code = -4
		//	*msg.Header.Error = err_str
		//	return msg
		//}
	}

	room, err := laowangye.Rooms.GetRoomById(0)
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
	if !req.Option.GetIsCoinRoom() && req.GetIsDaikai() == true {
		daikai_user := req.Header.GetUserId()
		*desk.IsDaikai = true
		*desk.DaikaiUser = daikai_user
		//更新代开列表
		err := roomAgent.CreateDesk(int32(ddproto.CommonEnumGame_GID_LAOWANGYE), desk.GetPassword(), desk.GetDeskId(), daikai_user, desk.GetTips(), time.Now().Unix(), desk.DeskOption.GetBoardsCout(), desk.DeskOption.GetMaxUser(), 0)
		if err != nil && req.GetOption().GetRoomCardBillType() == ddproto.COMMON_ENUM_ROOMCARD_BILL_TYPE_OWNER_PAY {
			//返还房费
			userService.INCRUserRoomcard(daikai_user, ownerFee, int32(ddproto.CommonEnumGame_GID_LAOWANGYE), "老王爷朋友桌，代开失败房费返还")
			*msg.Header.Code = -8
			*msg.Header.Error = err.Error()
			return msg
		}
		msg.DeskState = desk.GetClientDesk()
		log.T("用户%d代开房间%s成功！", daikai_user, desk.GetPassword())
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
	*msg.Header.Error = "创建房间成功"
	return msg
}

//进入房间
func EnterDeskHandler(req *ddproto.LwyEnterDeskReq, agent gate.Agent) {
	msg := &ddproto.LwyEnterDeskAck{
		Header: commonNewPorot.NewHeader(),
		DeskState: nil,
		IsReconnect: proto.Bool(false),
	}

	//处理断  线重连
	user, err := laowangye.FindUserById(req.Header.GetUserId())
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
		//如果没找到用户，则强行删除该用户在老王爷的session
		session := sessionService.GetSessionAuto(req.Header.GetUserId())
		if session != nil && session.GetGameId() == int32(ddproto.CommonEnumGame_GID_LAOWANGYE) {
			sessionService.DelSessionByKey(session.GetUserId(),session.GetRoomType(),session.GetGameId(),session.GetDeskId())
		}
	}

	//判断是否在其他游戏中
	roomType := ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND
	if req.GetRoomId() > 0 {
		roomType = ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN
	}
	_,err = sessionService.CanEnter(req.Header.GetUserId(), int32(ddproto.CommonEnumGame_GID_LAOWANGYE), int32(roomType))
	if err != nil {
		*msg.Header.Code = -4
		*msg.Header.Error = err.Error()
		agent.WriteMsg(msg)
		return
	}

	//先找到room
	var room *laowangye.Room
	var desk *laowangye.Desk
	room, err = laowangye.Rooms.GetRoomById(req.GetRoomId())
	if err != nil {
		log.E("用户%d进房pwd:%v roomId:%v失败，原因：%s", req.Header.GetUserId(), req.GetPassword(), req.GetRoomId(), err.Error())
		*msg.Header.Code = -1
		*msg.Header.Error = "进入房间失败！"
		agent.WriteMsg(msg)
		return
	}

	//朋友桌进牌桌
	if !req.GetIsCoinRoom() {
		desk, err = room.GetDeskByPassword(req.GetPassword())
		//清理session,以便容错
		if len(req.GetPassword()) == 6 && req.GetPassword()[2] == '1' && req.GetPassword()[3] == '1' {
			go func() {
				defer Error.ErrorRecovery("EnterDeskHandler->DelSessionByKey()")
				sessionService.DelSessionByKey(req.Header.GetUserId(), int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND), int32(ddproto.CommonEnumGame_GID_LAOWANGYE), 0)
			}()
		}
		if err != nil {
			log.E("用户%d进房%v失败，原因：%s", req.Header.GetUserId(), req.GetPassword(), err.Error())
			*msg.Header.Code = -2
			*msg.Header.Error = "进入房间失败！未找到该牌桌！"
			agent.WriteMsg(msg)
			return
		}

		//房卡准入检测
		err = roomService.HasEnoughRoomcardToEnterDesk(desk.DeskOption.GetRoomCardBillType(), ddproto.CommonEnumGame_GID_LAOWANGYE, desk.DeskOption.GetBoardsCout(), desk.DeskOption.GetMaxUser(), desk.DeskOption.GetChanelId(), req.Header.GetUserId())
		if err != nil {
			log.E("用户%d进房%v失败，原因：%s", req.Header.GetUserId(), req.GetPassword(), err.Error())
			*msg.Header.Code = -7
			*msg.Header.Error = err.Error()
			agent.WriteMsg(msg)
			return
		}
	}else {
		//判断准入
		if userService.GetUserCoin(req.Header.GetUserId()) < int64(room.GetEnterCoin()) {
			log.E("用户%d进房pwd:%v roomId:%v失败，原因：%s", req.Header.GetUserId(), req.GetPassword(), req.GetRoomId(), "小于进房最小金币数")
			*msg.Header.Code = -5
			*msg.Header.Error = fmt.Sprintf("您的金币小于%d,进房失败。", room.GetEnterCoin())
			agent.WriteMsg(msg)
			return
		}
		if req.GetPassword() != "" || req.GetRoomId() == 0 {
			//带密码的金币场
			desk,err = room.GetDeskByPassword(req.GetPassword())
			if err != nil {
				log.E("用户%d进房pwd:%v roomId:%v失败，原因：%s", req.Header.GetUserId(), req.GetPassword(), req.GetRoomId(), err.Error())
				*msg.Header.Code = -2
				*msg.Header.Error = "进入房间失败！未找到该房间！"
				agent.WriteMsg(msg)
				return
			}
		}else {
			//传统金币场
			desk, err = room.GetFreeCoinDesk()
			if err != nil {
				log.E("用户%d进房pwd:%v roomId:%v失败，原因：%s", req.Header.GetUserId(), req.GetPassword(), req.GetRoomId(), err.Error())
				*msg.Header.Code = -2
				*msg.Header.Error = "进入房间失败！"
				agent.WriteMsg(msg)
				return
			}
		}
	}

	_,err = desk.AddUser(req.Header.GetUserId(), agent)
	if err != nil {
		log.E("用户%d进房pwd:%v roomId:%v失败，原因：%s", req.Header.GetUserId(), req.GetPassword(), req.GetRoomId(), err.Error())
		*msg.Header.Code = -3
		*msg.Header.Error = err.Error()
		agent.WriteMsg(msg)
		return
	}

	log.E("用户%d进房pwd:%v roomId:%v成功。", req.Header.GetUserId(), req.GetPassword(), req.GetRoomId())
	msg.DeskState = desk.GetClientDesk()
	*msg.Header.Code = 1
	*msg.Header.Error = "新玩家加入房间成功！"
	agent.WriteMsg(msg)

	//传统金币场，真人加入后调度机器人加入
	if desk.Room.GetRoomId() > 0 {
		desk.AutoJoinRobot()
	}
}

//金币场牌桌列表
func CoinDeskListHandler(req *ddproto.CommonReqListCoinInfo, a gate.Agent) {
	//ack := &ddproto.CommonAckListCoinInfo{
	//	Header: commonNewPorot.NewHeader(),
	//	GameId: req.GameId,
	//	RoomType: req.RoomType,
	//	ListCoinLevel: []*ddproto.CommonCoinLevelInfo{},
	//	ListCoinDesk: []*ddproto.CommonCoinDeskInfo{},
	//}
	//
	//if req.GetGameId() != int32(ddproto.CommonEnumGame_GID_LAOWANGYE) {
	//	*ack.Header.Code = 0
	//	*ack.Header.Error = "参数异常！"
	//	a.WriteMsg(ack)
	//	return
	//}
	//
	//bankRule := ddproto.LwyEnumBankerRule(req.GetRoomType())
	//
	//for _,room := range laowangye.Rooms {
	//	if room.GetRoomId() == 0 {
	//		for _, desk := range room.Desks {
	//			if !desk.GetIsCoinRoom() || desk.DeskOption.GetIsPrivate() || desk.DeskOption.GetBankRule() != bankRule {
	//				continue
	//			}
	//			new_coin_desk := &ddproto.CommonCoinDeskInfo{
	//				DeskId: proto.Int32(desk.GetDeskId()),
	//				Password: proto.String(desk.GetPassword()),
	//				CfgStr: proto.String(desk.GetTips()),
	//				CurrPlayerNum: proto.Int32(desk.GetSitedGammerNum()),
	//				TotalPlayerNum: proto.Int32(desk.DeskOption.GetMaxUser()),
	//				BaseValue: proto.Int64(desk.DeskOption.GetBaseScore()),
	//				CoinLimitEnter: proto.Int64(desk.Room.GetEnterCoin()),
	//				CoinLimitLeave: proto.Int64(desk.Room.GetEnterCoin()),
	//			}
	//			ack.ListCoinDesk = append(ack.ListCoinDesk, new_coin_desk)
	//		}
	//	}else {
	//		if int32(room.GetBankRule()) == req.GetRoomType() {
	//			new_coin_level := &ddproto.CommonCoinLevelInfo{
	//				Id: proto.Int32(room.GetRoomId()),
	//				Name: proto.String(room.GetRoomTitle()),
	//				BaseValue: proto.Int64(room.GetBaseChip()),
	//				CoinLimitEnter: proto.Int64(room.GetEnterCoin()),
	//				CoinLimitLeave: proto.Int64(room.GetEnterCoin()),
	//				PlayingUserCount: proto.Int32(room.GetPlayerNum(bankRule)),
	//			}
	//			ack.ListCoinLevel = append(ack.ListCoinLevel, new_coin_level)
	//		}
	//	}
	//}
	//
	//a.WriteMsg(ack)
}

//断线重连处理
func ReconnectProcess(user *laowangye.User) {
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

//入座
func SiteDownHandler(req *ddproto.LwySiteDownReq, agent gate.Agent) {
	return
	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		user.DoSiteDown()
	}else {
		user = &laowangye.User{
			Agent: agent,
			LwySrvUser: nil,
		}
		user.SendSiteDownACK(-1, "您当前未在房间中,入座失败！")
	}
}

//离坐
//func SiteUpHandler(req *ddproto.LwySiteUpReq, agent gate.Agent) {
//	return
//	user, err := laowangye.FindUserById(req.Header.GetUserId())
//	if err == nil {
//		user.Desk.ReqLock.Lock()
//		defer user.Desk.ReqLock.Unlock()
//		user.UpdateAgent(agent)
//
//		user.DoSiteUp()
//	}else {
//		user = &laowangye.User{
//			Agent: agent,
//			LwySrvUser: nil,
//		}
//		user.SendEnterDeskBcACK(-1, "您当前未在房间中,离座失败！")
//	}
//}

//准备
func ReadyHandler(req *ddproto.LwySwitchReadyReq, agent gate.Agent) {
	user, err := laowangye.FindUserById(req.Header.GetUserId())
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
		user = &laowangye.User{
			Agent: agent,
			LwySrvUser: nil,
		}
		user.SendReadyAck(-1, "您当前未在房间中")
	}
}

//抢庄
func QiangzhuangHandler(req *ddproto.LwyQiangzhuangReq, agent gate.Agent) {
	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		user.DoQiangzhuang(req.GetScore())
	}else {
		user = &laowangye.User{
			Agent: agent,
			LwySrvUser: nil,
		}
		user.SendQiangzhuangAck(-1, "您当前未在房间中")
	}
}

//押注
func YazhuHandler(req *ddproto.LwyYazhuReq, agent gate.Agent) {
	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		user.DoYazhu(req.GetYazhuType(), req.GetYazhuScore())
	}else {
		user = &laowangye.User{
			Agent: agent,
			LwySrvUser: nil,
		}
		user.SendYazhuAck(-1, "您当前未在房间中")
	}
}

//摇色子
func YaoshaiziHandler(req *ddproto.LwyYaoshaiziReq, agent gate.Agent) {
	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		user.DoYaoshaizi()
	}else {
		user = &laowangye.User{
			Agent: agent,
			LwySrvUser: nil,
		}
		user.SendYaoshaiziAck(-1, "您当前未在房间中")
	}
}

//申请解散房间
func ApplyDissolveReqHandler(req *ddproto.CommonReqApplyDissolve, agent gate.Agent)  {
	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		user.DoApplyDissolve()
	}
}

//确定、取消解散房间
func DissolveBackReqHandler(req *ddproto.CommonReqApplyDissolveBack, agent gate.Agent) {
	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		user.DoDissolveBack(req.GetAgree())
	}
}

//聊天请求
func MessageReqHandler(req *ddproto.CommonReqMessage, agent gate.Agent) {
	user, err := laowangye.FindUserById(req.Header.GetUserId())
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
	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err != nil {
		return
	}
	user.Desk.ReqLock.Lock()
	defer user.Desk.ReqLock.Unlock()
	user.UpdateAgent(agent)

	//离开房间
	user.DoLeaveDesk()
}

//房主解散房间不扣除房卡
func OwnerDissolveHandler(req *ddproto.LwyOwnerDissolveReq, agent gate.Agent) {
	msg := &ddproto.LwyOwnerDissolveAck{
		Header: commonNewPorot.NewHeader(),
	}
	//代开解散
	if req.GetDeskId() > 0 && req.Header.GetUserId() > 0 {
		room,_ := laowangye.Rooms.GetRoomById(0)
		desk,err := room.GetDeskById(req.GetDeskId())
		if err != nil {
			//如果没找到该房间，则强制从代开列表中清除这条记录
			err = roomAgent.DoDissolve(req.Header.GetUserId(), int32(ddproto.CommonEnumGame_GID_LAOWANGYE), req.GetDeskId())
			log.T("roomAgent.DoDissolve() err:%v", err)
			*msg.Header.Code = 0
			*msg.Header.Error = "代开房间解散成功！"
			agent.WriteMsg(msg)
			return
		}
		if desk.GetDaikaiUser() != req.GetHeader().GetUserId() {
			*msg.Header.Code = -3
			*msg.Header.Error = "您不是房主，代开解散失败！"
			agent.WriteMsg(msg)
			return
		}
		//清除牌桌状态
		desk.Room.RemoveFriendDesk(desk.GetDeskId())
		*msg.Header.Code = 0
		*msg.Header.Error = "代开房间解散成功！"
		agent.WriteMsg(msg)
		return
	}
	user, err := laowangye.FindUserById(req.Header.GetUserId())
	if err == nil {
		user.Desk.ReqLock.Lock()
		defer user.Desk.ReqLock.Unlock()
		user.UpdateAgent(agent)

		if user.IsOwner() && user.Desk.GetCircleNo() == 1 && user.Desk.GetStatus() == ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY {
			log.T("房主解散房间%d成功", user.Desk.GetPassword())
			user.SendOwnerDissolveAck(0, "解散房间成功！")
			user.Desk.SendDissolveDoneBc(true)

			//清除牌桌状态
			user.Desk.Room.RemoveFriendDesk(user.Desk.GetDeskId())
		}else {
			log.E("房主解散房间%d失败，原因：用户%d不是房主，或房间已开局（当前：%s）。", user.Desk.GetPassword(), user.GetUserId(), user.Desk.GetStatus().String())
			user.SendOwnerDissolveAck(-1, "解散房间失败，因为你不是房主或房间已开局！")
		}
	}
}

//金币场房间列表
func CoinRoomListHandler(req *ddproto.LwyCoinRoomListReq, agent gate.Agent) {
	//ack := &ddproto.LwyCoinRoomListAck{
	//	Header: commonNewPorot.NewHeader(),
	//	Rooms: []*ddproto.LwySrvRoom{},
	//}
	//
	//for _,room := range laowangye.Rooms{
	//	if room.GetRoomType() == int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN) && room.GetBankRule() == req.GetBankRule() {
	//		ack.Rooms = append(ack.Rooms, room.LwySrvRoom)
	//	}
	//}
	//
	//agent.WriteMsg(ack)
}
