package friendPlay

import (
	"github.com/name5566/leaf/module"
	"casino_server/common/log"
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/data"
	"casino_majiang/msg/protogo"
	"casino_common/common/Error"
	"casino_common/common/userService"
	"casino_common/common/consts"
	"casino_majianagv2/core/ins/changSha"
	"casino_majiang/msg/funcsInit"
	"casino_common/common/sessionService"
)

type FMJRoom struct {
	*module.Skeleton         //leaf 的骨架
	*skeleton.SkeletonMJRoom //麻将room骨架
	desks []api.MjDesk       //所有的desk
}

//初始化一个朋友桌room
func NewDefaultFMJRoom(s *module.Skeleton) api.MjRoom {
	return &FMJRoom{
		Skeleton: s,
	}
}

//room创建房间
func (r *FMJRoom) CreateDesk(config interface{}) (api.MjDesk, error) {
	c := config.(*data.SkeletonMJConfig)
	//1,找到是否有已经创建的房间
	oldDesk := r.getDeskByOwer(c.Owner)
	if oldDesk != nil && oldDesk.GetStatus().IsNotGaming() {
		//如果房间没有开始游戏..则返回老的房间,否则创建新的房间...
		return oldDesk, nil
	}

	//2,判断房卡是否足够
	createFee, err := r.CalcCreateFee(c.BoardsCout)
	if err != nil {
		log.E("玩家[%v]创建房间的时候出错..传入的局数[%v]有误...", c.BoardsCout)
		return nil, Error.ERR_SYS
	}

	rc := userService.GetUserRoomCard(c.Owner)
	if rc < createFee {
		log.E("玩家[%v]创建房间的时候出错..房卡[%v]不足...", c.Owner, rc)
		return nil, Error.NewError(consts.ACK_RESULT_ERROR, "房卡不足，创建房间失败")
	}

	//3,创建房间
	var desk api.MjDesk
	//根据不同的类型来得到不同地区的麻将
	if c.MjRoomType == int32(mjproto.MJRoomType_roomType_changSha) {
		desk = changSha.NewChangShaFMJDesk(c) //创建长沙麻将朋友桌
	} else {
		desk = NewFMJDesk(c) //创建成都麻将朋友桌
	}
	//4，进入房间
	err = desk.EnterUser(c.Owner)
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
			c := d.GetMJConfig()
			if c.Password == ps {
				return d
			}
		}
	}
	return nil
}

func (r *FMJRoom) getDeskByOwer(userId uint32) api.MjDesk {
	for _, d := range r.desks {
		if d != nil {
			if d.GetMJConfig().Owner == userId {
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

//todo
func (r *FMJRoom) CalcCreateFee(n int32) (int64, error) {
	return 0, nil
}

// room 解散房间...解散朋友桌
func (r *FMJRoom) DissolveDesk(desk api.MjDesk, sendMsg bool) error {
	//清楚数据,1,session相关。2,
	log.T("%v开始解散...", desk.GetMJConfig().DeskId)
	for _, user := range desk.GetUsers() {
		if user != nil {
			sessionService.DelSessionByKey(user.GetUserId(), desk.GetMJConfig().RoomType)
			agent := user.GetAgent()
			if agent != nil {
				agent.SetUserData(nil)
			}
		}
	}
	log.T("开始删除desk[%v]...", desk.GetMJConfig().DeskId)

	//发送解散房间的广播
	rmErr := r.RmDesk(desk)
	if rmErr != nil {
		log.E("删除房间失败,errmsg[%v]", rmErr)
		return rmErr
	}

	//删除reids
	r.DelMjDeskRedis(desk)

	//删除房间
	log.T("删除desk[%v]之后，发送删除的广播...", desk.GetMJConfig().DeskId)
	if sendMsg {
		//发送解散房间的广播
		dissolve := newProto.NewGame_AckDissolveDesk()
		*dissolve.DeskId = desk.GetMJConfig().DeskId
		*dissolve.PassWord = desk.GetMJConfig().Password
		*dissolve.UserId = desk.GetMJConfig().Owner
		desk.BroadCastProto(dissolve)
	}

	return nil

}
