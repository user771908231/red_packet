package majiang

import (
	"strings"
	"casino_server/utils/redisUtils"
	"casino_server/utils/numUtils"
)
// session相关的...

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
	return nil
}
