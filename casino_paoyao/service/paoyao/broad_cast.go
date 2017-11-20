package paoyao

import (
	"casino_common/proto/ddproto"
	"casino_common/proto/funcsInit"
	"github.com/golang/protobuf/proto"
	"errors"
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
	msg := &ddproto.PaoyaoEnterDeskAck{
		Header: commonNewPorot.NewHeader(),
		DeskState: user.GetClientDesk(),
	}

	return user.WriteMsg(msg)
}

//进房广播
func (user *User) SendEnterDeskBC() error {
	msg := &ddproto.PaoyaoEnterDeskBc{
		Header: commonNewPorot.NewHeader(),
		User: user.GetClientUser(),
	}

	return user.BroadExclude(msg, user.GetUserId())
}

//准备ack
func (user *User) SendReadyAck(code int32, err string) error {
	msg := &ddproto.PaoyaoSwitchReadyBc{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//准备广播
func (user *User) SendReadyBC() error {
	msg := &ddproto.PaoyaoSwitchReadyBc{
		Header: commonNewPorot.NewHeader(),
		UserId: proto.Uint32(user.GetUserId()),
		IsReady: proto.Bool(true),
	}
	return user.BroadCast(msg)
}

//游戏结束统计数据广播
func (desk *Desk) SendGameEndResultBc() {
	return
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

//发送离线广播
func (user *User) SendOffineBc() error {
	isOffline := !user.GetIsOnline()
	msg := &ddproto.PaoyaoOfflineBc{
		UserId: user.UserId,
		IsOffline: &isOffline,
	}
	return user.Desk.BroadCast(msg)
}
