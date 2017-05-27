package skeleton

import (
	"casino_paodekuai/core/api"
	"casino_paodekuai/core/data"
	"casino_common/common/log"
	"fmt"
)

//room骨架
type PDKDeskSkeleton struct {
	DeskId      int32         //桌子的Id
	PokerParser api.PDKParser //扑克解析器
	Users       []api.PDKUser //玩家
}

func NewPDKDeskSkeleton() *PDKDeskSkeleton {
	return &PDKDeskSkeleton{}
}

func (d *PDKDeskSkeleton) GetDeskId() int32 {
	return d.DeskId
}

func (d *PDKDeskSkeleton) Dlog() string {
	return fmt.Sprintf("【desk-%v-%v】", d.GetDeskId(), 0)
}

//出牌的方法
func (d *PDKDeskSkeleton) ActOut(userId uint32, p interface{}) error {
	return nil
}

//准备
func (d *PDKDeskSkeleton) ActReady(userId uint32) error {
	return nil
}

//洗牌
func (d *PDKDeskSkeleton) XiPai() error {
	all := d.PokerParser.XiPai().([]*data.PokerCard)
	log.T("%v 洗牌，所有的牌:%v", d.Dlog(), all)
	return nil
}
