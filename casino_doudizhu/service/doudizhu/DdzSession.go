package doudizhu

import "casino_server/common/log"

func GetSession(userId uint32) *PDdzSession {
	return nil
}

//斗地主的session
func GetDdzDeskBySession(userId uint32) *DdzDesk {
	session := GetSession(userId)
	room := GetFDdzRoom()
	desk := room.GetDeskByDeskId(session.GetDeskId())
	if desk == nil {
		log.E("通过玩家[%v]的session[%v]没有找到对应的desk.", userId, session)
	}
	return desk
}