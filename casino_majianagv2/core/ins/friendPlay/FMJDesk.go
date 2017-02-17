package friendPlay

import (
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/ins/huParserIns"
	"github.com/name5566/leaf/module"
)

//朋友桌麻将的desk
type FMJDesk struct {
	*skeleton.SkeletonMJDesk
}

//创建一个朋友桌的desk
func NewFMJDesk(config *data.SkeletonMJConfig, s *module.Skeleton) api.MjDesk {
	//desk 骨架
	desk := &FMJDesk{
		SkeletonMJDesk: skeleton.NewSkeletonMJDesk(config, s),
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

func (d *FMJDesk) GetFMJUser(userId uint32) *FMJUser {
	u := d.GetUserByUserId(userId)
	if u != nil {
		return u.(*FMJUser)
	}
	return nil
}

func (d *FMJDesk) GetFMJUsers() []*FMJUser {
	ret := make([]*FMJUser, len(d.GetUsers()))
	for i, u := range d.GetUsers() {
		if u != nil {
			ret[i] = u.(*FMJUser)
		}
	}
	return ret
}
