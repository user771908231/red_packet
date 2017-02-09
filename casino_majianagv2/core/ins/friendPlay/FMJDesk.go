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
func NewFMJDesk(config data.SkeletonMJConfig) api.MjDesk {
	//判断创建条件：房卡，
	//desk 骨架
	desk := &FMJDesk{
		SkeletonMJDesk: skeleton.NewSkeletonMJDesk(config),
	}
	desk.HuParser = huParserIns.NewChengDuHuParser()
	return desk
}

//离开房间
func (d *FMJDesk) Leave(userId uint32) error {
	return nil
}
