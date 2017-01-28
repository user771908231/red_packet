package skeleton

import "casino_common/common/log"

//常见的get set 方法 需要放置在这里
func (f *SkeletonMJDesk) GetMJConfig() interface{} {
	log.Debug("玩家[%v]进入fdesk")
	return f.config.Owner
}
