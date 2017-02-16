package friendPlay

import (
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/ins/huParserIns"
)

//朋友桌麻将的desk
type FMJDesk struct {
	*skeleton.SkeletonMJDesk
}

//创建一个朋友桌的desk
func NewFMJDesk(config *data.SkeletonMJConfig) api.MjDesk {

	//desk 骨架
	desk := &FMJDesk{
		SkeletonMJDesk: skeleton.NewSkeletonMJDesk(config),
	}

	//胡牌的解析器
	desk.HuParser = huParserIns.NewChengDuHuParser()
	return desk
}

//离开房间
func (d *FMJDesk) Leave(userId uint32) error {
	return nil
}

func (d *FMJDesk) Ready(userId uint32) error {
	d.SkeletonMJDesk.Ready(userId)
	//如果人数还是不够，就需要在计划增加机器人
	d.begin()
	return nil
}

//开始
func (d *FMJDesk) begin() {

}
