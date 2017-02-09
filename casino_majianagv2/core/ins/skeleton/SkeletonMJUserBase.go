package skeleton

import (
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
)

func (u *SkeletonMJUser) GetUserId() uint32 {
	return u.userId
}

func (u *SkeletonMJUser) GetStatus() *data.MjUserStatus {
	return u.status
}

//todo
func (u *SkeletonMJUser) GetDesk() api.MjDesk {
	return u.desk
}
