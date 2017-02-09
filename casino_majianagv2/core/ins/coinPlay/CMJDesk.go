package coinPlay

import (
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/ins/huParserIns"
)

//朋友桌麻将的desk
type CMJDesk struct {
	*skeleton.SkeletonMJDesk
}

//创建一个朋友桌的desk
func NewCMJDesk(config data.SkeletonMJConfig) api.MjDesk {
	//判断创建条件：房卡，
	//desk 骨架
	desk := &CMJDesk{
		SkeletonMJDesk: skeleton.NewSkeletonMJDesk(config),
	}
	desk.HuParser = huParserIns.NewChengDuHuParser()
	return desk
}

//离开房间
func (d *CMJDesk) Leave(userId uint32) error {
	return nil
}

//准备
func (d *CMJDesk) Ready(userId uint32) error {
	d.Lock()
	defer d.Unlock()
	d.SkeletonMJDesk.Ready(userId)
	d.initEnterTimer() //房间进入一个人之后开始计划添加机器人
	d.begin()
	return nil
}

//todo
func (d *CMJDesk) initEnterTimer() {

}

//todo
func (d *CMJDesk) begin() {

}
