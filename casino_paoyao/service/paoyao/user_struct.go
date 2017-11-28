package paoyao

import (
	"github.com/name5566/leaf/gate"
	"casino_common/proto/ddproto"
	"casino_common/common/sessionService"
	"github.com/golang/protobuf/proto"
	"errors"
	"casino_common/common/log"
	"casino_common/common/service/whiteListService"
	"casino_common/common/Error"
)

//用户
type User struct {
	Agent gate.Agent
	*ddproto.PaoyaoSrvUser
	*Desk
}

//更新Session
func (user *User) UpdateSession() error {
	var err error
	if !user.Desk.GetIsCoinRoom() {
		//朋友桌
		go func() {
			defer Error.ErrorRecovery("UpdateSession()->friend")
			sessionService.UpdateSession(user.GetUserId(), int32(ddproto.COMMON_ENUM_GAMESTATUS_GAMING), int32(ddproto.CommonEnumGame_GID_PAOYAO), 0, user.GetRoomId(), user.GetDeskId(), 0, false, false,int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND), user.Desk.GetPwd())
		}()
	}else {
		//金币场
		go func() {
			defer Error.ErrorRecovery("UpdateSession()->coin")
			sessionService.UpdateSession(user.GetUserId(), int32(ddproto.COMMON_ENUM_GAMESTATUS_GAMING), int32(ddproto.CommonEnumGame_GID_PAOYAO), 0, user.GetRoomId(), user.GetDeskId(), 0, false, false,int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN), user.Desk.GetPwd())
		}()
	}

	if user.Agent != nil {
		user.Agent.SetUserData(user.GetUserId())
	}

	return err
}

//清理Session
func (user *User) ClearSession() error {
	log.T("开始清理用户%d的session", user.GetUserId())
	if !user.Desk.GetIsCoinRoom() {
		//朋友桌
		go func() {
			defer Error.ErrorRecovery("ClearSession()->friend")
			sessionService.DelSessionByKey(user.GetUserId(), int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND), int32(ddproto.CommonEnumGame_GID_PAOYAO), user.GetDeskId())
		}()
	}else {
		go func() {
			defer Error.ErrorRecovery("ClearSession()->coin")
			sessionService.DelSessionByKey(user.GetUserId(), int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN), int32(ddproto.CommonEnumGame_GID_PAOYAO), user.GetDeskId())
		}()
	}

	if user.Agent != nil {
		user.Agent.SetUserData(nil)
	}
	//更新user缓存
	delete(UserDeskMap, user.GetUserId())
	return nil
}

//获取客户端user
func (user *User) GetClientUser() *ddproto.PaoyaoClientUser {
	client_user := &ddproto.PaoyaoClientUser{
		UserId: proto.Uint32(user.GetUserId()),
		IsOnline: proto.Bool(user.GetIsOnline()),
		Index: proto.Int32(user.GetIndex()),
		Pokers: nil,  //牌默认隐藏
		OutPai: GetClientPoker(user.GetOutPai()),
		SurplusPokerNum: proto.Int32(int32(len(user.GetPokers().GetPais()))),
		IsPass: user.IsPass,
		IsReady: user.IsReady,
		LastScore: user.LastScore,
		AllScore: user.Bill.Score,
		WxInfo: user.WxInfo,
		DissolveState: user.DissolveState,
	}
	return client_user
}

//是否为房主
func (user *User) IsOwner() bool {
	if user.Desk.GetOwner() == user.GetUserId() {
		return true
	}
	return false
}

//获取队友
func (user *User) GetTeamMateUser() *User {
	mate,_ := user.Desk.GetUserByUid(user.GetTeamMate())
	return mate
}

//更新链接
func (user *User) UpdateAgent(agent gate.Agent) error {
	if user != nil {
		user.Agent = agent
		return nil
	}
	return errors.New("user is nil")
}

//刷新白名单
func (user *User) CheckWhiteList() {
	//刷新白名单
	whiteListService.RefreshWhiteList(int32(ddproto.CommonEnumGame_GID_PAOYAO))
	//是否在白名单中
	whiteUser := whiteListService.GetWhiteUser(int32(ddproto.CommonEnumGame_GID_PAOYAO), user.GetUserId())
	if whiteUser != nil {
		*user.IsOnWhiteList = true
		*user.WhiteWinRate = whiteUser.WinRate
	}
}

//获取对方队伍分数
func (user *User) GetTeamScore() (ourside_score, oppsite_score int32) {
	for _,u := range user.Desk.Users {
		switch u.GetUserId() {
		case user.GetUserId(), user.GetTeamMate():
			ourside_score += u.GetDeskScore()
		default:
			oppsite_score += u.GetDeskScore()
		}
	}
	return
}

//对方已出完牌
func (user *User) IsOppoSideUserChupaiDone() bool {
	for _,u := range user.Desk.Users {
		switch u.GetUserId() {
		case user.GetUserId(), user.GetTeamMate():
		default:
			if u.Pokers != nil && len(u.Pokers.Pais) > 0 {
				return false
			}
		}
	}
	return true
}

//我方已出完牌
func (user *User) IsOurSideUserChupaiDone() bool {
	for _,u := range user.Desk.Users {
		switch u.GetUserId() {
		case user.GetUserId(), user.GetTeamMate():
			if u.Pokers != nil && len(u.Pokers.Pais) > 0 {
				return false
			}
		}
	}
	return true
}

//我方已扛旗
func (user *User) IsOurSideKangQi() bool {
	if user.GetIsKangQi() || user.GetTeamMateUser().GetIsKangQi() {
		return true
	}
	return false
}
