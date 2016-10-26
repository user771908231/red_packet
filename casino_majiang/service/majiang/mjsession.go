package majiang

import (
	"strings"
	"casino_server/utils/redisUtils"
	"casino_server/utils/numUtils"
	"casino_server/common/log"
)
// session相关的...


var MJUSER_SESSION_GAMESTATUS_NOGAME int32 = 1 //没有在游戏中
var MJUSER_SESSION_GAMESTATUS_FRIEND int32 = 2 //朋友桌


var MJSESSION_KEY_PRE = "redis_majiang_session"

func getSessionKey(userId uint32) string {
	idstr, _ := numUtils.Uint2String(userId)
	ret := strings.Join([]string{MJSESSION_KEY_PRE, idstr}, "_")
	return ret
}

func GetSession(userId uint32) *MjSession {
	s := redisUtils.GetObj(getSessionKey(userId), &MjSession{})
	if s != nil {
		return s.(*MjSession)
	} else {
		return nil
	}
}

//更新用户的session信息，具体更新什么信息待定
func UpdateSession(userId uint32, gameStatus int32, roomId int32, deskId int32, deskPassWord string) (*MjSession, error) {
	session := GetSession(userId)
	if session == nil {
		log.T("没有找到user[%v]的session,需要重新申请一个并保存...", userId)
		session = NewMjSession()
	}

	*session.UserId = userId
	*session.DeskId = deskId
	*session.RoomId = roomId
	*session.GameStatus = 0

	//保存session
	log.T("保存到redis的session【%v】", session)
	redisUtils.SetObj(getSessionKey(userId), session)
	return session, nil
}
