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
	desk.HuParser = huParserIns.NewChengDuHuParser(
		desk.GetMJConfig().BaseValue,
		desk.IsNeedZiMoJiaDi(),
		desk.IsNeedYaojiuJiangdui(),
		desk.IsDaodaohu(),
		desk.IsNeedMenqingZhongzhang(),
		desk.IsNeedZiMoJiaFan(),
		desk.GetMJConfig().CapMax,
	)
	return desk
}
