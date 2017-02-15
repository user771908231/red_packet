package friendPlay

import (
	"casino_majianagv2/core/ins/huParserIns"
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
)

//长沙麻将朋友桌

//朋友桌麻将的desk
type ChangShaFMJDesk struct {
	*skeleton.SkeletonMJDesk
}

//创建一个朋友桌的desk
func NewChangShaFMJDesk(config data.SkeletonMJConfig) api.MjDesk {

	//desk 骨架
	desk := &FMJDesk{
		SkeletonMJDesk: skeleton.NewSkeletonMJDesk(config),
	}

	desk.HuParser = huParserIns.NewChangShaHuParser()
	return desk
}
