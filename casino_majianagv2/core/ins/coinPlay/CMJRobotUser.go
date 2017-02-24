package coinPlay

import (
	"github.com/golang/protobuf/proto"
	"casino_majianagv2/core/ins/skeleton"
	"time"
	"casino_common/utils/timeUtils"
	"casino_common/common/Error"
)

type CMJRobotUser struct {
	*skeleton.SkeletonMJUser
}

func (u *CMJRobotUser) WriteMsg(msg proto.Message) error {
	time.AfterFunc(timeUtils.RandDuration(1, 5), func() {
		//log.T("开始给robot[%v]发送msg %v", u.GetUserId(), msg)
		defer Error.ErrorRecovery("发送信息给机器人")
		u.DoRobotAct(msg) //机器人发送信息
	})
	//随机的演示
	return nil
}
