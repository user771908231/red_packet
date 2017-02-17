package coinPlay

import (
	"github.com/name5566/leaf/module"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/ins/skeleton"
	"casino_common/common/consts/tableName"
	"casino_majianagv2/core/data"
	"casino_common/common/log"
	"casino_common/utils/db"
	"github.com/name5566/leaf/gate"
)

type CMjRoom struct {
	*module.Skeleton //leaf 的骨架
	*skeleton.SkeletonMJRoom
	RoomLevel int32 //金币场等级
}

func NewDefaultCMjRoom(s *module.Skeleton, l int32) api.MjRoom {
	ret := &CMjRoom{
		Skeleton:       s,
		RoomLevel:      l,
		SkeletonMJRoom: skeleton.NewSkeletonMJRoom(l),
	}
	return ret
}

//room创建房间
func (r *CMjRoom) CreateDesk(config interface{}, a gate.Agent) (api.MjDesk, error) {
	c := config.(*data.SkeletonMJConfig)
	c.DeskId, _ = db.GetNextSeq(tableName.DBT_MJ_DESK)
	//根据不同的类型来得到不同地区的麻将
	desk := NewCMJDesk(c, r.Skeleton) //创建成都麻将朋友桌
	desk.SetRoom(r)
	//4，进入房间
	err := desk.EnterUser(c.Owner, nil)
	if err != nil {
		log.E("玩家%v进入desk失败 err %v ", c.Owner, err)
		return nil, nil
	}
	return desk, nil
}
