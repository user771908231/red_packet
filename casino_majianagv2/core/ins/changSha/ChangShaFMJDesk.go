package changSha

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
func NewChangShaFMJDesk(config *data.SkeletonMJConfig) api.MjDesk {
	//desk 骨架
	desk := &ChangShaFMJDesk{
		SkeletonMJDesk: skeleton.NewSkeletonMJDesk(config),
	}
	desk.HuParser = huParserIns.NewChangShaHuParser()
	return desk
}

//打牌
func (d *ChangShaFMJDesk) ActOut(userId uint32, cardId int32) error {
	return nil
}

//进入游戏
func (d *ChangShaFMJDesk) Leave(userId uint32) error {
	return nil
}
