package coinPlay

import (
	"github.com/golang/protobuf/proto"
	"casino_majianagv2/core/ins/skeleton"
)

type CMJRobotUser struct {
	*skeleton.SkeletonMJUser
}

func (u *CMJRobotUser) WriteMsg(msg proto.Message) error {
	return u.DoRobotAct(msg)
}
