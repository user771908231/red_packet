package doudizhu

import "casino_server/common/log"

type Session struct {
	UserId uint32
	RoomId int32
	DeskId int32
}

func GetSession(userId uint32) *Session {
	return nil
}

func GetMjDeskBySession(userId uint32) *DdzDesk {
	//得到session
	session := GetSession(userId)
	if session == nil {
		return nil
	}
	log.T("得到用户[%v]的session[%v]", userId, session)

	//得到room
	room := GetMjroomBySession(userId)
	if room == nil {
		return nil
	}

	//返回desk
	return room.GetDeskByDeskId(session.DeskId)
}