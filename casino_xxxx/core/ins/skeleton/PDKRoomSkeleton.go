package skeleton

import (
	"casino_paodekuai/core/api"
	"github.com/name5566/leaf/util"
)

//room骨架
type PDKRoomSkeleton struct {
	desks *util.Map //所有的desk
}

func NewPDKRoomSkeleton() *PDKRoomSkeleton {
	return &PDKRoomSkeleton{
		desks: new(util.Map),
	}
}

//增加一个desk
func (rs *PDKRoomSkeleton) AddDesk(desk api.PDKDesk) {
	rs.desks.Set(desk.GetDeskId(), desk)
}

//得到一个desk
func (rs *PDKRoomSkeleton) GetDesk(deskId int32) api.PDKDesk {
	ret := rs.desks.Get(deskId)
	if ret != nil {
		return ret.(api.PDKDesk)
	}
	return nil
}

//删除一个desk
func (rs *PDKRoomSkeleton) RmDesk(deskId int32) error {
	rs.desks.Del(deskId)
	return nil
}
