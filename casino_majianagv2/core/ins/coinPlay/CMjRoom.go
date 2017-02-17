package coinPlay

import (
	"github.com/name5566/leaf/module"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/ins/skeleton"
)

type CMjRoom struct {
	*module.Skeleton //leaf 的骨架
	*skeleton.SkeletonMJRoom
	RoomLevel int32 //金币场等级
}

func NewDefaultCMjRoom(s *module.Skeleton, l int32) api.MjRoom {
	ret := &CMjRoom{
		Skeleton:       s,
		RoomLevel:      l,
		SkeletonMJRoom: skeleton.NewSkeletonMJRoom(l),
	}
	return ret
}
