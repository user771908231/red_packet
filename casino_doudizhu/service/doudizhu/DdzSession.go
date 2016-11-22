package doudizhu

import (
	"casino_server/common/log"
	"casino_server/utils/numUtils"
	"strings"
	"casino_server/utils/redisUtils"
)

func GetSession(userId uint32) *PDdzSession {
	data := redisUtils.GetObj(getSessionKey(userId), &PDdzSession{})
	if data == nil {
		return nil
	} else {
		return data.(*PDdzSession)
	}
}

//斗地主的session
func GetDdzDeskBySession(userId uint32) *DdzDesk {
	session := GetSession(userId)
	if session == nil {
		log.E("玩家[%v]的session没有找到", userId)
		return nil
	}

	room := GetFDdzRoom()
	desk := room.GetDeskByDeskId(session.GetDeskId())
	if desk == nil {
		log.E("通过玩家[%v]的session[%v]没有找到对应的desk.", userId, session)
	}
	return desk
}

var MJSESSION_KEY_PRE = "redis_ddz_session"

func getSessionKey(userId uint32) string {
	idstr, _ := numUtils.Uint2String(userId)
	ret := strings.Join([]string{MJSESSION_KEY_PRE, idstr}, "_")
	return ret
}

//更新session的信息
func UpdateSession(userId uint32, gamestatus int32, roomId int32, deskId int32) (*PDdzSession, error) {
	session := GetSession(userId)
	if session == nil {
		log.T("没有找到user[%v]的session,需要重新申请一个并保存...", userId)
		session = NewPDdzSession()
	}

	*session.UserId = userId
	*session.DeskId = deskId
	*session.RoomId = roomId
	*session.GameStatus = gamestatus

	//保存session
	log.T("保存到redis的斗地主的session【%v】", session)
	redisUtils.SetObj(getSessionKey(userId), session)
	return session, nil
}