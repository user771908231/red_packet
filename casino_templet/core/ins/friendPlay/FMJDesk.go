package friendPlay

import (
	"casino_templet/core/ins/skeleton"
	"casino_templet/core/api"
	"casino_templet/core/data"
)

//朋友桌麻将的desk
type FMJDesk struct {
	*skeleton.SkeletonMJDesk
}

//创建一个朋友桌的desk
func NewFMJDesk(config data.SkeletonMJConfig) api.MjDesk {
	return &FMJDesk{
		SkeletonMJDesk:skeleton.NewSkeletonMJDesk(config),
	}
}

//离开房间
func (d *FMJDesk) Leave(userId uint32) error {
	return nil
}
