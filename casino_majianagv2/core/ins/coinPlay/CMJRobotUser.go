package coinPlay

import (
	"github.com/golang/protobuf/proto"
	"casino_majianagv2/core/ins/skeleton"
	"time"
	"casino_common/utils/timeUtils"
	"casino_common/common/Error"
	"casino_common/common/log"
)

type CMJRobotUser struct {
	*skeleton.SkeletonMJUser
}

//发送OverTurn
func (u *CMJRobotUser) SendOverTurn(p proto.Message) error {
	u.WriteMsg(p)
	return nil
}
func (u *CMJRobotUser) WriteMsg(msg proto.Message) error {
	log.T("给ai玩家发送message:%v", msg)
	time.AfterFunc(timeUtils.RandDuration(1, 5), func() {
		//log.T("开始给robot[%v]发送msg %v", u.GetUserId(), msg)
		defer Error.ErrorRecovery("发送信息给机器人")
		u.DoRobotAct(msg) //机器人发送信息
	})
	return nil
}
