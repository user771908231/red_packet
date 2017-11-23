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
		PlayerInfo: []*ddproto.PaoyaoClientUser{},
	}

	for _,u := range user.Desk.Users {
		msg.PlayerInfo = append(msg.PlayerInfo, u.GetClientUser())
	}

	return user.BroadExclude(msg, user.GetUserId())
}

//准备ack
func (user *User) SendReadyAck(code int32, err string) error {
	msg := &ddproto.PaoyaoSwitchReadyBc{
		Header: commonNewPorot.NewHeader(),
		UserId: proto.Uint32(user.GetUserId()),
		IsReady: proto.Bool(user.GetIsReady()),
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
		IsReady: proto.Bool(user.GetIsReady()),
	}
	*msg.Header.Code = 0
	*msg.Header.Error = "准备成功"
	return user.BroadCast(msg)
}

//发牌overturn
func (desk *Desk) SendFapaiOt() error {
	msg := &ddproto.PaoyaoFapaiOt{
		Header: commonNewPorot.NewHeader(),
		Pokers: nil,
		CurrCircle: desk.CircleNo,
	}

	for _,u := range desk.Users {
		if u == nil {
			continue
		}
		*msg.Header.UserId = u.GetUserId()
		msg.Pokers = GetClientPoker(u.Pokers)
		u.WriteMsg(msg)
	}

	return nil
}

//出牌ack
func (user *User) SendChupaiAck(code int32, err string) error {
	msg := &ddproto.PaoyaoChupaiBc{
		Header: commonNewPorot.NewHeader(),
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	return user.WriteMsg(msg)
}

//出牌广播
func (user *User) SendChupaiBc(add_score int32) error {
	surplus_poker_num := int32(len(user.Pokers.Pais))
	msg := &ddproto.PaoyaoChupaiBc{
		Header: commonNewPorot.NewHeader(),
		OutPai:GetClientPoker(user.OutPai),
		SurplusPokerNum:proto.Int32(surplus_poker_num),
		FriendPoker: GetClientPoker(nil),
		DeskScore: user.Desk.CurrDeskScore,
		OppositeScore: proto.Int32(0),
		OurSideScore: proto.Int32(0),
		AddScore: &add_score,
	}
	*msg.Header.Code = 0
	*msg.Header.Error = "出牌成功！"

	//查看队友余牌
	mate := user.GetTeamMateUser()
	if surplus_poker_num == 0 && mate != nil {
		msg.FriendPoker = GetClientPoker(mate.GetPokers())
	}
	//双方分数
	*msg.OppositeScore, *msg.OurSideScore = user.GetTeamScore()
	//先给自己发
	user.WriteMsg(msg)

	//再给其他人发
	for _,u := range user.Desk.Users {
		if u == nil {
			continue
		}
		//双方分数
		*msg.OppositeScore, *msg.OurSideScore = u.GetTeamScore()
		u.WriteMsg(msg)
	}

	return nil
}

//过牌ack
func (user *User) SendGuopaiAck(code int32, err string) error {
	msg := &ddproto.PaoyaoGuopaiBc{
		Header: commonNewPorot.NewHeader(),
		UserId: user.UserId,
	}
	*msg.Header.Code = code
	*msg.Header.Error = err
	*msg.Header.UserId = user.GetUserId()
	return user.WriteMsg(msg)
}

//过牌bc
func (user *User) SendGuopaiBc() error {
	msg := &ddproto.PaoyaoGuopaiBc{
		Header: commonNewPorot.NewHeader(),
		UserId: user.UserId,
	}
	*msg.Header.Code = 0
	*msg.Header.Error = "过牌成功！"

	//广播
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
