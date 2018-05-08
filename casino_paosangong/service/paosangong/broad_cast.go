package paosangong

import (
	"casino_common/proto/ddproto"
	"casino_common/proto/funcsInit"
	"github.com/golang/protobuf/proto"
	"time"
	"errors"
	"casino_common/common/userService"
	"casino_common/gameManager/roomService"
	"casino_common/common/Error"
)

//发送消息
func (user *User) WriteMsg(msg proto.Message) error {
	if user == nil {
		return errors.New("user is nil.")
	}
	if user.Agent == nil {
		return errors.New("user.agent is nil.")
	}
	user.Agent.WriteMsg(msg)
	return nil
}

//牌桌广播
func (desk *Desk) BroadCast(msg proto.Message) error {
	for _,u := range desk.Users {
		if u != nil {
			u.WriteMsg(msg)
		}
	}
	return nil
}

//广播-排除某用户
func (desk *Desk) BroadExclude(msg proto.Message, exclude_user_id uint32) error {
	for _,u := range desk.Users {
		if u != nil && u.GetUserId() != exclude_user_id {
			u.WriteMsg(msg)
		}
	}
	return nil
}

//进房ack
func (user *User) SendEnterDeskAck() error {
	msg := &ddproto.NiuEnterDeskAck{
		Header: commonNewPorot.NewHeader(),
		DeskState: user.GetClientDesk(),
	}

	return user.WriteMsg(msg)
}

//进房广播
func (user *User) SendEnterDeskBC() error {
	msg := &ddproto.NiuEnterDeskBc{
		Header: commonNewPorot.NewHeader(),
		User: user.GetClientUser(),
		Owner: proto.Uint32(user.Desk.GetOwner()),
	}

	return user.BroadCast(msg)
}

//入座、离座 ack
func (user *User) SendEnterDeskBcACK(code int32, err string) error {
	msg := &ddproto.NiuEnterDeskBc{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//准备ack
func (user *User) SendReadyAck(code int32, err string) error {
	msg := &ddproto.NiuSwitchReadyAck{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//准备广播
func (user *User) SendReadyBC() error {
	msg := &ddproto.NiuSwitchReadyBc{
		Header: commonNewPorot.NewHeader(),
		User: proto.Uint32(user.GetUserId()),
		IsReady: proto.Bool(true),
	}
	return user.BroadCast(msg)
}

//房主开局overturn
func (desk *Desk) SendStartOt() error {
	msg := &ddproto.NiuStartGameOt{}
	owner,err := desk.GetUserByUid(desk.GetOwner())
	if err == nil {
		return owner.WriteMsg(msg)
	}

	//托管房主自动开局
	if owner,_ := desk.GetUserByUid(desk.GetOwner());owner.GetIsTuoguan() {
		owner.DoReady()
	}
	return err
}

//开始抢庄timer
func (desk *Desk) StartQiangzhuangTimer() {
	//设置抢庄倒计时，金币场独有
	if desk.QiangzhuangTimer != nil {
		desk.QiangzhuangTimer.Stop()
	}
	desk.StartTime = proto.Int64(time.Now().Unix())
	desk.QiangzhuangTimer = time.AfterFunc(7*time.Second, func() {
		defer Error.ErrorRecovery("desk.QiangzhuangTimer()")
		for _, u := range desk.Users {
			if u != nil && u.GetIsOnGamming() && u.GetBankerScore() == 0 {
				u.DoQiangzhuang(-1)
			}
		}
	})
}

//发起抢庄overturn
func (desk *Desk) SendQiangzhuangOt() error {
	//设置抢庄倒计时，金币场独有
	desk.StartQiangzhuangTimer()
	//发广播
	for _,u := range desk.Users {
		if u != nil {
			if !u.GetIsRobot() {
				//真人
				msg := &ddproto.NiuQiangzhuangOt{
					Header: commonNewPorot.NewHeader(),
					Pokers: GetClientPoker(u.Pokers),
					CurrCircle: proto.Int32(desk.GetCircleNo()),
					IsOnGaming: proto.Bool(u.GetIsOnGamming()),
					SurplusCoin:proto.Int64(userService.GetUserCoin(u.GetUserId())),
					DeskTime: proto.Int32(int32(desk.GetSurplusTime())),
				}
				u.WriteMsg(msg)
			}else {
				//机器人
				u.DoRobotQiangzhuang()
			}
		}
	}
	//托管操作
	desk.DoTuoguan()

	return nil
}

//抢庄结果
func (desk *Desk) SendQiangzhuangResBc() error {
	res := []*ddproto.NiuQiangzhuangResItem{}

	for _,u := range desk.Users {
		if u != nil {
			res = append(res, &ddproto.NiuQiangzhuangResItem{
				User: u.UserId,
				Score: u.BankerScore,
				IsBanker: proto.Bool(u.IsBanker()),
			})
		}
	}

	msg := &ddproto.NiuQiangzhuangResBc{
		Result: res,
	}

	return desk.BroadCast(msg)
}

//加倍timer
func (desk *Desk) StartJiabeiTimer() {
	//如果是金币场，则设置7秒超时
	if desk.JiaBeiTimer != nil {
		desk.JiaBeiTimer.Stop()
	}
	desk.StartTime = proto.Int64(time.Now().Unix())
	desk.JiaBeiTimer = time.AfterFunc(7 * time.Second, func() {
		defer Error.ErrorRecovery("desk.JiabeiTimer()")
		for _,u1 := range desk.Users {
			if u1 != nil && u1.GetIsOnGamming() && !u1.IsBanker() && u1.GetDoubleScore() == 0 {
				u1.DoJiabei(int64(desk.DeskOption.GetDiFen()))
			}
		}
		//将倒计时置空
		desk.JiaBeiTimer = nil
	})
}

//发起加倍overturn
func (desk *Desk) SendJiabeiOt() error {
	//开始加倍倒计时，金币场独有
	desk.StartJiabeiTimer()
	//发广播
	for _,u := range desk.Users {
		if u != nil {
			if !u.GetIsRobot() {
				//真人
				msg := &ddproto.NiuJiabeiOt{
					Header: commonNewPorot.NewHeader(),
					Pokers: GetClientPoker(u.Pokers),
					CuurBanker: proto.Int32(int32(desk.GetCurrBanker())),
					CuurCircle: proto.Int32(desk.GetCircleNo()),
					IsOnGaming: proto.Bool(u.GetIsOnGamming()),
					TuizhuScore: proto.Int32(u.GetTuizhuScore()),
					DeskTime: proto.Int32(int32(desk.GetSurplusTime())),
				}
				u.WriteMsg(msg)
			}else {
				//机器人
				u.DoRobotJiabei()
			}
		}
	}

	//托管操作
	desk.DoTuoguan()

	return nil
}

//加倍ack
func (user *User) SendJiabeiAck(code int32, err string) error {
	msg := &ddproto.NiuJiabeiAck{
		Header: commonNewPorot.NewHeader(),
		Score: proto.Int64(user.GetDoubleScore()),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//加倍广播
func (user *User) SendJiabeiBc() error {
	msg := &ddproto.NiuJiabeiBc{
		Header: commonNewPorot.NewHeader(),
		Score: user.DoubleScore,
		UserId: user.UserId,
	}
	*msg.Header.Code = 1
	*msg.Header.Error = "加倍成功！"
	*msg.Header.UserId = user.GetUserId()
	return user.Desk.BroadCast(msg)
}

//抢庄ack
func (user *User) SendQiangzhuangAck(code int32, err string) error {
	msg := &ddproto.NiuQiangzhuangAck{
		Header: commonNewPorot.NewHeader(),
		Score: proto.Int64(user.GetBankerScore()),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//抢庄bc
func (user *User) SendQiangzhuangBC() error {
	msg := &ddproto.NiuQiangzhuangAck{
		Header: commonNewPorot.NewHeader(),
		Score: proto.Int64(user.GetBankerScore()),
	}

	*msg.Header.UserId = user.GetUserId()
	*msg.Header.Code = 1
	*msg.Header.Error = "抢庄成功！"
	return user.Desk.BroadCast(msg)
}

//亮牌timer
func (desk *Desk) StartLiangpaiTimer() {
	//如果是金币场，则设置10秒超时
	if desk.LiangpaiTimer != nil {
		desk.LiangpaiTimer.Stop()
	}
	desk.StartTime = proto.Int64(time.Now().Unix())
	desk.LiangpaiTimer = time.AfterFunc(10 * time.Second, func() {
		defer Error.ErrorRecovery("desk.LiangpaiTimer()")
		for _,u1 := range desk.Users {
			if u1 != nil && u1.GetIsOnGamming() && !u1.GetIsLiangpai() {
				u1.DoLiangpai()
			}
		}
		//将倒计时置空
		desk.LiangpaiTimer = nil
	})
}

//发起加倍overturn
func (desk *Desk) SendLiangpaiOt() error {
	//开始亮牌倒计时
	desk.StartLiangpaiTimer()
	//发广播
	for _,u := range desk.Users {
		if u != nil {
			if !u.GetIsRobot() {
				//真人
				msg := &ddproto.NiuLiangpaiOt{
					Header: commonNewPorot.NewHeader(),
					Pokers: GetClientPoker(u.Pokers),
					IsOnGaming: proto.Bool(u.GetIsOnGamming()),
					DeskTime: proto.Int32(int32(desk.GetSurplusTime())),
				}
				u.WriteMsg(msg)
			}else {
				//机器人
				//u.DoRobotJiabei()
			}
		}
	}

	//托管操作
	desk.DoTuoguan()
	return nil
}

//托管ack
func (user *User) SendTuoguanAck(code int32, err string) error {
	msg := &ddproto.NiuTuoguanBc{
		Header: commonNewPorot.NewHeader(),
		IsTuoguan: proto.Bool(user.GetIsTuoguan()),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//托管bc
func (user *User) SendTuoguanBc() error {
	msg := &ddproto.NiuTuoguanBc{
		Header: commonNewPorot.NewHeader(),
		IsTuoguan: proto.Bool(user.GetIsTuoguan()),
	}
	*msg.Header.UserId = user.GetUserId()
	*msg.Header.Code = 1
	*msg.Header.Error = "[开启/关闭]托管成功！"
	return user.Desk.BroadCast(msg)
}

//亮牌ack
func (user *User) SendLiangpaiAck(code int32, err string) error {
	msg := &ddproto.NiuLiangpaiAck{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//亮牌bc
func (user *User) SendLiangpaiBc() error {
	msg := &ddproto.NiuLiangpaiBc{
		UserId: user.UserId,
		Pokers: GetClientPoker(user.Pokers),
	}
	return user.Desk.BroadCast(msg)
}

//游戏结束统计数据广播
func (desk *Desk) SendGameEndResultBc() {
	item_list := []*ddproto.NiuniuUserBill{}
	var max_score int64 = 0
	var max_user uint32 = 0
	for _,u := range desk.Users {
		if u != nil && u.GetIndex() != -1 {
			item_list = append(item_list, u.Bill)
			if max_score == 0 {
				max_score = u.Bill.GetScore()
				max_user = u.GetUserId()
			}
			if u.Bill.GetScore() > max_score {
				max_score = u.Bill.GetScore()
				max_user = u.GetUserId()
			}
		}
	}

	allUsers, winUsers := []uint32{}, []uint32{}
	for _,u := range desk.Users {
		if u == nil || u.GetIndex() == -1 {
			continue
		}

		allUsers = append(allUsers, u.GetUserId())
		if u.Bill.GetScore() == max_score {
			winUsers = append(winUsers, u.GetUserId())
		}
	}

	create_user_id := desk.GetOwner()
	//如果是代开
	if desk.GetIsDaikai() {
		create_user_id = desk.GetDaikaiUser()
	}
	//AA扣房卡
	roomService.DoDecUsersRoomcard(desk.DeskOption.GetRoomCardBillType(), ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN, desk.DeskOption.GetMaxCircle(), desk.DeskOption.GetMaxUser(), desk.DeskOption.GetChanelId(), allUsers, winUsers, desk.GetCircleNo(), create_user_id)

	//发送10局牌局结束后的统计数据
	msg := &ddproto.NiuGameEnd{
		Header: commonNewPorot.NewHeader(),
		Data: item_list,
		EndTime: proto.Int64(time.Now().Unix()),
		BigWiner: &max_user,
	}
	desk.BroadCast(msg)
}


//申请解散房间广播
func (user *User) SendApplyDissolveBc() error {
	msg := &ddproto.CommonBcApplyDissolve{
		Header: commonNewPorot.NewHeader(),
		UserId: user.UserId,
	}
	for _,u := range user.Desk.Users {
		if u == nil {
			continue
		}
		u.WriteMsg(msg)
	}
	return nil
}

//确定、拒绝解散房间BC
func (user *User) SendDissolveBackBc(isAgree bool) error {
	msg := &ddproto.CommonAckApplyDissolveBack{
		Header: commonNewPorot.NewHeader(),
		UserId: user.UserId,
		Agree: &isAgree,
	}

	for _,u := range user.Desk.Users {
		if u == nil {
			continue
		}
		u.WriteMsg(msg)
	}
	return nil
}

//确认解散房间广播
func (desk *Desk) SendDissolveDoneBc(isAllAgree bool) error {
	msg := &ddproto.NiuDeskDissolveDoneBc{
		Header: commonNewPorot.NewHeader(),
		IsDissolve: &isAllAgree,
	}

	for _,u := range desk.Users {
		if u == nil {
			continue
		}
		u.WriteMsg(msg)
	}
	return nil
}

//房主解散房间不扣房卡ack
func (user *User) SendOwnerDissolveAck(code int32, err string) error {
	msg := &ddproto.NiuOwnerDissolveAck{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//发送离线广播
func (user *User) SendOffineBc() error {
	isOffline := !user.GetIsOnline()
	msg := &ddproto.NiuOfflineBc{
		UserId: user.UserId,
		IsOffline: &isOffline,
	}
	return user.Desk.BroadCast(msg)
}
