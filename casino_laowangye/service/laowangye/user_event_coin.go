package laowangye

import (
	"casino_common/common/log"
	"casino_common/proto/ddproto"
	"casino_common/common/userService"
)

//金币场准备
func (u *User) DoReadyCoin() {
	defer u.Desk.WipeSnapShot()

	if !u.Desk.GetIsCoinRoom() || u.Desk.GetStatus() != ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY {
		u.SendReadyAck(-5, "准备失败，当前牌桌未在准备状态。")
		return
	}
	log.T("用户%d开始准备。", u.GetUserId())
	if u.GetIsReady() == true {
		u.SendReadyAck(1, "已经准备了，请不要重复准备！")
		return
	}

	//检测金币
	if userService.GetUserCoin(u.GetUserId()) < int64(u.Room.GetEnterCoin()) {
		u.SendReadyAck(-6, "金币不足，请充值！")
		return
	}

	//准备成功
	*u.IsReady = true
	u.SendReadyAck(1, "准备成功！")
	u.SendReadyBC()

	//尝试开局
	u.Desk.DoStart()
}

//金币场是否可以开局
func (desk *Desk) CoinHasEoughtUserReady() bool {
	var i int32 = 0
	for _,u := range desk.Users{
		if u != nil && u.GetIsReady() && u.GetIsOnline() {
			i++
			if i >= desk.DeskOption.GetMinUser() {
				return true
			}
		}
	}

	return false
}

