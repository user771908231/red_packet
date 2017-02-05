package friendPlay

import (
	"github.com/name5566/leaf/module"
	"casino_server/common/log"
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/data"
)

type FMJRoom struct {
	*module.Skeleton         //leaf 的骨架
	*skeleton.SkeletonMJRoom //麻将room骨架
	desks []api.MjDesk       //所有的desk
}

//初始化一个朋友桌room
func NewDefaultFMJRoom(s *module.Skeleton) api.MjRoom {
	return &FMJRoom{
		Skeleton:s,
	}
}

//room创建房间
func (r *FMJRoom) CreateDesk(config interface{}) (error, api.MjDesk) {
	c := config.(data.SkeletonMJConfig)
	//1,创建房间
	desk := NewFMJDesk(c)

	//2，进入房间
	err := desk.EnterUser(c.Owner)
	if err != nil {
		log.E("进入desk失败")
		return nil, nil
	}
	return nil, nil
}

//通过key得到一个desk
func (r *FMJRoom) getDeskByPassword(ps string) api.MjDesk {
	for _, d := range r.desks {
		if d != nil {
			c := d.GetMJConfig().(data.SkeletonMJConfig)
			if c.Password == ps {
				return d
			}
		}
	}
	return nil
}

//进入desk
func (r *FMJRoom) EnterUser(userId uint32, key string) error {
	desk := r.getDeskByPassword(key)
	if desk == nil {
		log.E("进入房间失败")
		return nil

	}
	err := desk.EnterUser(userId)
	if err != nil {
		log.E("进入房间失败...")
	}

	return nil
}
