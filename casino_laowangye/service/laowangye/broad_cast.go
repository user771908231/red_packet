package laowangye

import (
	"casino_common/proto/ddproto"
	"casino_common/proto/funcsInit"
	"github.com/golang/protobuf/proto"
	"time"
	"errors"
	"casino_common/common/userService"
	//"casino_common/gameManager/roomService"
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
	msg := &ddproto.LwyEnterDeskAck{
		Header: commonNewPorot.NewHeader(),
		DeskState: user.GetClientDesk(),
	}

	return user.WriteMsg(msg)
}

//进房广播
func (user *User) SendSiteDownBC() error {
	msg := &ddproto.LwySiteDownBc{
		Header: commonNewPorot.NewHeader(),
		User: user.GetClientUser(),
	}

	return user.BroadCast(msg)
}

//入座、离座 ack
func (user *User) SendSiteDownACK(code int32, err string) error {
	msg := &ddproto.LwySiteDownBc{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//准备ack
func (user *User) SendReadyAck(code int32, err string) error {
	msg := &ddproto.LwySwitchReadyBc{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//准备广播
func (user *User) SendReadyBC() error {
	msg := &ddproto.LwySwitchReadyBc{
		Header: commonNewPorot.NewHeader(),
		User: proto.Uint32(user.GetUserId()),
		IsReady: proto.Bool(true),
	}
	return user.BroadCast(msg)
}

//房主开局overturn
func (desk *Desk) SendStartOt() error {
	msg := &ddproto.LwyStartGameOt{}
	owner,err := desk.GetUserByUid(desk.GetOwner())
	if err == nil {
		return owner.WriteMsg(msg)
	}
	return err
}

//开始抢庄timer
func (desk *Desk) StartQiangzhuangTimer() {
	//设置抢庄倒计时，金币场独有
	if desk.GetIsCoinRoom() {
		if desk.QiangzhuangTimer != nil {
			desk.QiangzhuangTimer.Stop()
		}
		desk.QiangzhuangTimer = time.AfterFunc(10*time.Second, func() {
			for _, u := range desk.Users {
				if u != nil && u.GetIsOnGamming() && u.GetBankerScore() == 0 {
					u.DoQiangzhuang(-1)
				}
			}
		})
	}
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
				msg := &ddproto.LwyQiangzhuangOt{
					Header: commonNewPorot.NewHeader(),
					CurrCircle: proto.Int32(desk.GetCircleNo()),
					IsOnGaming: proto.Bool(u.GetIsOnGamming()),
					SurplusCoin:proto.Int64(userService.GetUserCoin(u.GetUserId())),
				}
				u.WriteMsg(msg)
			}else {
				//机器人
				u.DoRobotQiangzhuang()
			}
		}
	}
	return nil
}

//抢庄结果
//func (desk *Desk) SendQiangzhuangResBc() error {
	//res := []*ddproto.LwyQiangzhuangResItem{}
	//
	//for _,u := range desk.Users {
	//	if u != nil {
	//		res = append(res, &ddproto.LwyQiangzhuangResItem{
	//			User: u.UserId,
	//			Score: u.BankerScore,
	//			IsBanker: proto.Bool(u.IsBanker()),
	//		})
	//	}
	//}
	//
	//msg := &ddproto.LwyQiangzhuangResBc{
	//	Result: res,
	//}
	//
	//return desk.BroadCast(msg)
	//return nil
//}

//发起加倍overturn
func (desk *Desk) SendYazhuOt() error {
	//发广播
	for _,u := range desk.Users {
		if u != nil {
			if !u.GetIsRobot() {
				//真人
				msg := &ddproto.LwyYazhuOt{
					Header: commonNewPorot.NewHeader(),
					Banker: proto.Uint32(desk.GetCurrBanker()),
				}
				u.WriteMsg(msg)
			}else {
				//机器人
				u.DoRobotJiabei()
			}
		}
	}
	return nil
}

//押注ack
func (user *User) SendYazhuAck(code int32, err string) error {
	msg := &ddproto.LwyYazhuBc{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//加倍广播
//func (user *User) SendJiabeiBc() error {
//	msg := &ddproto.LwyJiabeiBc{
//		Header: commonNewPorot.NewHeader(),
//		Score: user.DoubleScore,
//		UserId: user.UserId,
//	}
//	*msg.Header.Code = 1
//	*msg.Header.Error = "加倍成功！"
//	*msg.Header.UserId = user.GetUserId()
//	return user.Desk.BroadCast(msg)
//}

//抢庄ack
func (user *User) SendQiangzhuangAck(code int32, err string) error {
	msg := &ddproto.LwyQiangzhuangBc{
		Header: commonNewPorot.NewHeader(),
		Score: proto.Int64(user.GetBankerScore()),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//游戏结束统计数据广播
func (desk *Desk) SendGameEndResultBc() {
	//item_list := []*ddproto.LwyUserBill{}
	//var max_score int64 = 0
	//var max_user uint32 = 0
	//for _,u := range desk.Users {
	//	if u != nil {
			//item_list = append(item_list, u.Bill)
			//if max_score == 0 {
			//	max_score = u.Bill.GetScore()
			//	max_user = u.GetUserId()
			//}
			//if u.Bill.GetScore() > max_score {
			//	max_score = u.Bill.GetScore()
			//	max_user = u.GetUserId()
			//}
		//}
	//}

	//allUsers, winUsers := []uint32{}, []uint32{}
	//for _,u := range desk.Users {
	//	if u == nil {
	//		continue
	//	}

		//allUsers = append(allUsers, u.GetUserId())
		//if u.Bill.GetScore() == max_score {
		//	winUsers = append(winUsers, u.GetUserId())
		//}
	//}

	//create_user_id := desk.GetOwner()
	////如果是代开
	//if desk.GetIsDaikai() {
	//	create_user_id = desk.GetDaikaiUser()
	//}
	////AA扣房卡
	//roomService.DoDecUsersRoomcard(desk.DeskOption.GetRoomCardBillType(), ddproto.CommonEnumGame_GID_LAOWANGYE, desk.DeskOption.GetBoardsCout(), desk.DeskOption.GetMaxUser(), desk.DeskOption.GetChanelId(), allUsers, winUsers, desk.GetCircleNo(), create_user_id)
	//
	////发送10局牌局结束后的统计数据
	//msg := &ddproto.LwyGameEnd{
	//	Header: commonNewPorot.NewHeader(),
	//	Data: item_list,
	//	EndTime: proto.Int64(time.Now().Unix()),
	//	BigWiner: &max_user,
	//}
	//desk.BroadCast(msg)
}


//申请解散房间广播
func (user *User) SendApplyDissolveBc() error {
	msg := &ddproto.CommonBcApplyDissolve{
		Header: commonNewPorot.NewHeader(),
		UserId: user.UserId,
	}
	return user.Desk.BroadCast(msg)
}

//确定、拒绝解散房间BC
func (user *User) SendDissolveBackBc(isAgree bool) error {
	msg := &ddproto.CommonAckApplyDissolveBack{
		Header: commonNewPorot.NewHeader(),
		UserId: user.UserId,
		Agree: &isAgree,
	}

	return user.Desk.BroadCast(msg)
}

//确认解散房间广播
func (desk *Desk) SendDissolveDoneBc(isAllAgree bool) error {
	msg := &ddproto.LwyDeskDissolveDoneBc{
		Header: commonNewPorot.NewHeader(),
		IsDissolve: &isAllAgree,
	}

	return desk.BroadCast(msg)
}

//房主解散房间不扣房卡ack
func (user *User) SendOwnerDissolveAck(code int32, err string) error {
	msg := &ddproto.LwyOwnerDissolveAck{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//发送离线广播
func (user *User) SendOffineBc() error {
	isOffline := !user.GetIsOnline()
	msg := &ddproto.LwyOfflineBc{
		UserId: user.UserId,
		IsOffline: &isOffline,
	}
	return user.Desk.BroadCast(msg)
}
