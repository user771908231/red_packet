package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
)

//常见的get set 方法 需要放置在这里
func (f *SkeletonMJDesk) GetMJConfig() data.SkeletonMJConfig {
	log.Debug("玩家[%v]进入fdesk")
	return f.config
}

func (r *SkeletonMJDesk) GetStatus() *data.MjDeskStatus {
	return r.status
}
