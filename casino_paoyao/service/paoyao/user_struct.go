package paoyao

import (
	"github.com/name5566/leaf/gate"
	"casino_common/proto/ddproto"
	"casino_common/common/sessionService"
	"github.com/golang/protobuf/proto"
	"casino_common/common/userService"
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
	user_info := userService.GetUserById(user.GetUserId())
	wx_info := &ddproto.WeixinInfo{}
	if user_info != nil {
		wx_info.City = user_info.City
		wx_info.HeadUrl = proto.String(userService.GetUserHeadImg(user_info))
		wx_info.NickName = user_info.NickName
		wx_info.OpenId = user_info.OpenId
		wx_info.Sex = user_info.Sex
		wx_info.UnionId = user_info.UnionId
	}
	client_user := &ddproto.PaoyaoClientUser{
		UserId: proto.Uint32(user.GetUserId()),
		IsOnline: proto.Bool(user.GetIsOnline()),
		Index: proto.Int32(user.GetIndex()),
		//Pokers: GetClientPoker(user.Pokers),
		IsReady: proto.Bool(user.GetIsReady()),
		WxInfo: wx_info,
		LastScore: user.LastScore,
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
