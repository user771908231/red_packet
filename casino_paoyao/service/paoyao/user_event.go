package paoyao

import (
	"casino_common/common/consts"
	"casino_common/common/log"
	"casino_common/proto/ddproto"
	"casino_common/proto/funcsInit"
	"errors"
	"github.com/golang/protobuf/proto"
	"time"
	"casino_common/common/Error"
	//"casino_common/common/service/countService"
)

//朋友桌准备
func (u *User) DoReadyFriend() {
	log.T("用户%d开始准备。", u.GetUserId())
	defer u.Desk.WipeSnapShot()

	if u.Desk.GetStatus() == ddproto.PaoyaoEnumDeskStatus_PAOYAO_DESK_STATUS_WAIT_READY {
		if u.GetIsReady() == true {
			log.E("用户%d在第%d圈，重复准备!", u.GetUserId(), u.Desk.GetCircleNo())
			u.SendReadyAck(1, "请不要重复准备！")
			return
		} else {
			log.T("用户%d在第%d圈,准备成功！", u.GetUserId(), u.Desk.GetCircleNo())
			*u.IsReady = true
			u.SendReadyAck(1, "准备成功！")
			u.SendReadyBC()
			//尝试开局
			if err := u.Desk.IsAllReady(); err == nil {
				u.Desk.DoStart()
			}
			return
		}
	}else {
		//中途加入，但不在准备阶段
		u.SendReadyAck(-3, "当前不在准备阶段。")
	}
}

//离开房间
func (user *User) DoLeaveDesk() error {
	log.T("用户%d开始离开房间%d。", user.GetUserId(), user.Desk.GetPwd())
	msg := &ddproto.CommonAckLeaveDesk{
		Header:     commonNewPorot.NewHeader(),
		UserId:     user.UserId,
		IsExchange: proto.Bool(false),
	}
	if !user.Desk.GetIsCoinRoom() {
		//朋友桌
		log.E("用户%d退出房间%d失败，原因：朋友桌不能离开。", user.GetUserId(), user.Desk.GetPwd())
		*msg.Header.Code = -1
		*msg.Header.Error = "朋友桌不能离开房间。"
		user.WriteMsg(msg)
		return errors.New("朋友桌不能离开房间。")
	}

	//金币场，未开局直接退出
	user.Desk.RemoveUser(user.GetUserId())

	defer user.Desk.WipeSnapShot()
	return nil
}

//开始比赛
func (desk *Desk) DoStart() error {


	return nil
}

//牌桌--10圈比赛打完初始化
func (desk *Desk) DoEnd() error {
	//每一圈结束
	//所有用户标记为处于游戏结束
	for _,u := range desk.Users{
		*u.IsReady = false
	}

	//朋友桌
	if !desk.GetIsCoinRoom() {

		//结束统计
		desk.DoCountEnd()

	}else {

	}

	return nil
}

//发起解散房间
func (user *User) DoApplyDissolve() error {
	log.T("用户%d发起解散房间", user.GetUserId())

	if user.Desk.GetIsOnDissolve() == true {
		log.T("用户%d发起解散房间失败，原因：%s", user.GetUserId(), "正处于解散投票阶段，不要重复申请！")
		return errors.New("正处于解散投票阶段，不要重复申请！")
	}

	//更改状态
	*user.Desk.IsOnDissolve = true
	user.Desk.DissolveUser = proto.Uint32(user.GetUserId())
	*user.Desk.DissolveTime = time.Now().Unix()

	//5分钟后强制解散房间
	if user.Desk.DissolveTimer != nil {
		user.Desk.DissolveTimer.Stop()
		user.Desk.DissolveTimer = nil
	}
	user.Desk.DissolveTimer = time.AfterFunc(consts.APPLYDISSOLVE_DURATION, func() {
		if user.Desk.GetIsOnDissolve() == true {
			for _, u := range user.Desk.Users {
				if u != nil && u.GetDissolveState() == 0 {
					u.DoDissolveBack(true)
				}
			}
		}
	})

	//发起解散房间广播
	user.SendApplyDissolveBc()

	//自动投票
	user.DoDissolveBack(true)

	return nil
}

//同意、拒绝解散房间
func (user *User) DoDissolveBack(isAgree bool) error {
	log.T("用户%d同意或拒绝解散房间,是否同意：%v", user.GetUserId(), isAgree)
	if user.Desk.GetIsOnDissolve() == false {
		log.E("用户%d解散房间投票失败，原因：当前房间未处于解散投票阶段，投票失败！", user.GetUserId())
		return errors.New("当前房间未处于解散投票阶段，投票失败！")
	}

	if user.GetDissolveState() != 0 {
		log.E("用户%d解散房间投票失败，原因：请不要重复投票！", user.GetUserId())
		return errors.New("请不要重复投票！")
	}

	//更新用户解散投票状态
	if isAgree == true {
		*user.DissolveState = 1
	} else {
		*user.DissolveState = -1
	}

	//发送投票广播
	user.SendDissolveBackBc(isAgree)

	if isAgree == false {
		//如果有人拒绝解散，则初始化房间和用户状态
		//todo  user.Desk.SendDissolveDoneBc(isAgree)
		if user.Desk.DissolveTimer != nil {
			user.Desk.DissolveTimer.Stop()
			user.Desk.DissolveTimer = nil
		}
		//初始化桌面解散状态
		*user.Desk.IsOnDissolve = false
		*user.Desk.DissolveTime = 0
		for _, u := range user.Desk.Users {
			*u.DissolveState = 0
		}
		return nil
	}

	//确认所有人都已投票
	for _, u := range user.Desk.Users {
		if u != nil {
			if u.GetDissolveState() == 0 {
				//如果离线，则自动同意
				if u.GetIsOnline() == false {
					*u.DissolveState = 1
				} else {
					log.T("房间%d解散房间失败，原因：用户%d未投票", user.Desk.GetPwd(), user.GetUserId())
					//有人未投票
					return nil
				}
			}
		}
	}

	// 当所有人都确认解散,发送解散成功或失败广播:
	// 如果解散失败则发协议，解散成功则直接发送牌局结束协议。
	if isAgree == false {
		//todo
		//user.Desk.SendDissolveDoneBc(isAgree)
	}

	if user.Desk.DissolveTimer != nil {
		user.Desk.DissolveTimer.Stop()
		user.Desk.DissolveTimer = nil
	}
	//初始化桌面解散状态
	*user.Desk.IsOnDissolve = false
	*user.Desk.DissolveTime = 0
	for _, u := range user.Desk.Users {
		if u != nil {
			*u.DissolveState = 0
		}
	}

	if isAgree == true {
		//删除牌桌状态
		user.Desk.Room.RemoveFriendDesk(user.Desk.GetDeskId())
		log.T("房间%d解散房间成功", user.Desk.GetPwd())
	}

	return nil
}

//统计开局
func (desk *Desk) DoCountStart() {
	//更新统计时间
	*desk.OneStartTime = time.Now().Unix()

	if desk.GetCircleNo() == 1 {
		*desk.AllStartTime = time.Now().Unix()
	}
}

//统计结束
func (desk *Desk) DoCountEnd() {
	if !desk.GetIsCoinRoom() {
		//朋友桌
		//插入10局记录
		if desk.GetCircleNo() == desk.DeskOption.GetMaxCircle() {
			go func(){
				defer Error.ErrorRecovery("DoCountEnd()->all")
				desk.InsertAllCounter()
			}()
		}

		//插入1局记录
		go func() {
			defer Error.ErrorRecovery("DoCountEnd()->one")
			desk.InsertOneCounter()
		}()
	}else {
		//金币场

	}
}
