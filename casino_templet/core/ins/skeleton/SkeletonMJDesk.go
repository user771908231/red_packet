package skeleton

import (
	"casino_common/common/log"
	"casino_templet/core/data"
)

//desk 的骨架,业务逻辑的方法 放置在这里
type SkeletonMJDesk struct {
	config data.SkeletonMJConfig //这里不用使用指针，此配置创建之后不会再改变
}

func NewSkeletonMJDesk(config data.SkeletonMJConfig) *SkeletonMJDesk {
	return &SkeletonMJDesk{}
}

func (f *SkeletonMJDesk) EnterUser(userId uint32) error {
	log.Debug("玩家[%v]进入fdesk")
	return nil
}
