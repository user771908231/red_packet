package coinPlay

import (
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/ins/huParserIns"
	"github.com/name5566/leaf/timer"
	"github.com/name5566/leaf/module"
)

//朋友桌麻将的desk
type CMJDesk struct {
	*skeleton.SkeletonChengDuMJDesk
	RobotEnterTimer *timer.Timer
	CoinLimit       int64 //金币场的金币最低限制
}

//创建一个朋友桌的desk
func NewCMJDesk(config *data.SkeletonMJConfig, s *module.Skeleton) api.MjDesk {
	//判断创建条件：房卡，
	//desk 骨架
	desk := &CMJDesk{
		SkeletonChengDuMJDesk: skeleton.NewSkeletonChengDuMJDesk(config, s),
	}
	desk.HuParser = huParserIns.NewChengDuHuParser(
		desk.GetMJConfig().BaseValue,
		desk.IsNeedZiMoJiaDi(),
		desk.IsNeedYaojiuJiangdui(),
		desk.IsDaodaohu(),
		desk.IsNeedMenqingZhongzhang(),
		desk.IsNeedZiMoJiaFan(),
		desk.GetMJConfig().CapMax)
	return desk
}

//金币场的玩家
func (d *CMJDesk) getCMJUser(user api.MjUser) *CMJUser {
	return user.(*CMJUser)
}
