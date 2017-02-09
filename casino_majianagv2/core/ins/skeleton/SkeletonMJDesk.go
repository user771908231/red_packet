package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
)

//desk 的骨架,业务逻辑的方法 放置在这里
type SkeletonMJDesk struct {
	config   data.SkeletonMJConfig //这里不用使用指针，此配置创建之后不会再改变
	status   *data.MjDeskStatus    //桌子的所有状态都在这里
	HuParser api.HuPaerApi         //胡牌解析器
}

func NewSkeletonMJDesk(config data.SkeletonMJConfig) *SkeletonMJDesk {
	desk := &SkeletonMJDesk{
		config: config,
	}
	return desk
}

func (f *SkeletonMJDesk) EnterUser(userId uint32) error {
	log.Debug("玩家[%v]进入fdesk")
	return nil
}

//准备
func (f *SkeletonMJDesk) Ready(userId uint32) error {
	return nil
}

//定缺
func (f *SkeletonMJDesk) DingQue(userId uint32, color int32) error {
	return nil
}
