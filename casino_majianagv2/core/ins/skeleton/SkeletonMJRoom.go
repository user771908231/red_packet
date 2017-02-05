package skeleton

import "casino_majianagv2/core/api"

type SkeletonMJRoom struct {
}

//得到一个desk
func (r *SkeletonMJRoom) GetDesk() api.MjDesk {
	return nil
}

//进入一个User
func (r *SkeletonMJRoom) EnterUser(userId uint32) error {
	return nil
}
