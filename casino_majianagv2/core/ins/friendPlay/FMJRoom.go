package friendPlay

import (
	"github.com/name5566/leaf/module"
	"casino_common/common/log"
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/data"
	"casino_majiang/msg/protogo"
	"casino_common/common/Error"
	"casino_common/common/userService"
	"casino_common/common/consts"
	"casino_majianagv2/core/ins/changSha"
	"casino_common/utils/chessUtils"
	"casino_common/proto/ddproto"
	"casino_common/common/consts/tableName"
	"casino_common/utils/db"
	"casino_common/common/sessionService"
	"casino_majiang/service/majiang"
	"github.com/name5566/leaf/gate"
)

type FMJRoom struct {
	*module.Skeleton         //leaf 的骨架
	*skeleton.SkeletonMJRoom //麻将room骨架
	desks []api.MjDesk       //所有的desk
}

//初始化一个朋友桌room
func NewDefaultFMJRoom(s *module.Skeleton) api.MjRoom {
	ret := &FMJRoom{
		Skeleton:       s,
		SkeletonMJRoom: skeleton.NewSkeletonMJRoom(0),
	}
	return ret
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
	createFee := r.CalcCreateFee(c.BoardsCout) //开房费用
	rc := userService.GetUserRoomCard(c.Owner)
	if rc < createFee {
		log.E("玩家[%v]创建房间的时候出错..房卡[%v]不足...", c.Owner, rc)
		return nil, Error.NewError(consts.ACK_RESULT_ERROR, "房卡不足，创建房间失败")
	}

	//3,创建房间
	var desk api.MjDesk
	c.Password = r.RandRoomKey()
	c.DeskId, _ = db.GetNextSeq(tableName.DBT_MJ_DESK)

	//根据不同的类型来得到不同地区的麻将
	if c.MjRoomType == int32(mjproto.MJRoomType_roomType_changSha) {
		desk = changSha.NewChangShaFMJDesk(c, r.Skeleton) //创建长沙麻将朋友桌
	} else {
		desk = NewFMJDesk(c, r.Skeleton) //创建成都麻将朋友桌
	}
	desk.SetRoom(r) //给desk设置room
	r.AddDesk(desk) //
	return desk, nil
}

//通过key得到一个desk
func (r *FMJRoom) getDeskByPassword(ps string) api.MjDesk {
	for _, d := range r.desks {
		if d != nil {
			if d.GetMJConfig().Password == ps {
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
func (r *FMJRoom) EnterUser(userId uint32, key string, a gate.Agent) error {
	desk := r.getDeskByPassword(key)
	if desk == nil {
		sessionService.DelSessionByKey(userId, majiang.ROOMTYPE_FRIEND)
		log.E("没有找到对应的desk,进入房间失败")
		return Error.NewError(consts.ACK_RESULT_FAIL, "房间号输入错误")
	}
	err := desk.EnterUser(userId, nil) //朋友桌 进入房间
	if err != nil {
		log.E("进入房间失败...")
	}

	return nil
}

//todo
func (r *FMJRoom) CalcCreateFee(n int32) (int64) {
	return 0
}

func (r *FMJRoom) RandRoomKey() string {
	//金币场没有房间号码
	roomKey := chessUtils.GetRoomPass(int32(ddproto.CommonEnumGame_GID_MAHJONG))
	//1,判断roomKey是否已经存在
	if r.IsRoomKeyExist(roomKey) {
		//log.E("房间密钥[%v]已经存在,创建房间失败,重新创建", roomKey)
		return r.RandRoomKey()
	} else {
		//log.T("最终得到的密钥是[%v]", roomKey)
		return roomKey
	}
	return ""
}

//判断roomkey是否已经存在了
func (r *FMJRoom) IsRoomKeyExist(roomkey string) bool {
	//log.T("测试 r == nil ? %v ", r == nil)
	//log.T("测试 r.Desks == nil ? %v ", r.Desks == nil)
	ret := false
	for i := 0; i < len(r.Desks); i++ {
		d := r.Desks[i]
		if d != nil && d.GetMJConfig().Password == roomkey {
			ret = true
			break
		}
	}
	return ret
}
