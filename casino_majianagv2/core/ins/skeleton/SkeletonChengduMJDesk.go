package skeleton

import (
	"casino_majianagv2/core/data"
	"github.com/name5566/leaf/module"
)

//成都麻将的骨架类
type SkeletonChengDuMJDesk struct {
	*SkeletonMJDesk
}

//成都麻将的骨架方法
func NewSkeletonChengDuMJDesk(config *data.SkeletonMJConfig, s *module.Skeleton) *SkeletonChengDuMJDesk {
	return &SkeletonChengDuMJDesk{
		SkeletonMJDesk: NewSkeletonMJDesk(config, s),
	}
}
