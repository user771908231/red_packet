package bbproto

import (
	"casino_server/utils/timeUtils"
	"time"
)

/**
	初始化User登陆转盘奖励是否可以领取的状态
	1.如果最后领取的时间为空,那么表示今天的还没有领取
	2.如果最后一次领取的时间不为今天,那么表示今天的还没有领取
 */
func (u *User) OninitLoginTurntableState(){
	var able bool = false

	if u.GetLoginTurntableTime() == "" {
		able = true
	}

	if u.GetLoginTurntableTime() == timeUtils.Format(time.Now()) {
		able = true
	}

	u.LoginTurntable = &able
}