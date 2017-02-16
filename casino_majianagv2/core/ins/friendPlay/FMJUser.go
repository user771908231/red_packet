package friendPlay

import (
	"casino_majianagv2/core/ins/skeleton"
	"github.com/golang/protobuf/proto"
)

type FMJUser struct {
	*skeleton.SkeletonMJUser
}

func (u *FMJUser) SendOverTurn(p proto.Message) error {
	//如果是金币场有超时的处理...
	u.WriteMsg(p)
	return nil
}
