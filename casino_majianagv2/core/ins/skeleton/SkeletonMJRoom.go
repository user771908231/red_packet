package skeleton

import (
	"casino_majianagv2/core/api"
	"errors"
	"casino_common/utils/redisUtils"
	"strings"
	"casino_common/utils/numUtils"
	"casino_majiang/service/majiang"
	"sync"
	"casino_common/common/sessionService"
	"casino_common/common/log"
	"casino_majiang/msg/funcsInit"
)

type SkeletonMJRoom struct {
	RoomId      int32
	Desks       []api.MjDesk
	RoomMnanger api.MjRoomMgr
	sync.Mutex
}

func NewSkeletonMJRoom(id int32) *SkeletonMJRoom {
	return &SkeletonMJRoom{
		RoomId: id,
	}
}

func (r *SkeletonMJRoom) GetRoomId() int32 {
	return r.RoomId

}
func (r *SkeletonMJRoom) GetRoomLevel() int32 {
	return 0
}

//得到一个desk
func (r *SkeletonMJRoom) GetDesk(id int32) api.MjDesk {
	for _, d := range r.Desks {
		if d != nil {
			log.T("开始通过id %v 查找deskId:%v", id, d.GetMJConfig().DeskId)
			if d.GetMJConfig().DeskId == id {
				return d
			}
		}
	}
	return nil
}

//进入一个User
func (r *SkeletonMJRoom) EnterUser(userId uint32, key string) error {
	return nil
}

func (r *SkeletonMJRoom) RmDesk(desk api.MjDesk) error {
	index := -1
	for i, d := range r.Desks {
		if d != nil && d.GetMJConfig().DeskId == desk.GetMJConfig().DeskId {
			index = i
			break
		}
	}

	if index >= 0 {
		r.Desks = append(r.Desks[:index], r.Desks[index+1:]...)
		return nil
	} else {
		return errors.New("删除失败，没有找到对应的desk")
	}

}

//删除redis中麻将的数据 应该通过游戏编号来删除????
func (r *SkeletonMJRoom) DelMjDeskRedis(desk api.MjDesk) error {

	//删除desk
	key := r.GetDeskRedisKey(desk.GetMJConfig().DeskId)
	redisUtils.Del(key)

	//删除running key
	r.DelRunningDeskKey(desk.GetMJConfig().DeskId)

	return nil
}
func (r *SkeletonMJRoom) GetDeskRedisKey(id int32) string {
	idStr, _ := numUtils.Int2String(id)
	return strings.Join([]string{majiang.REDIS_KEY_MJ_DESK, idStr}, "_")

}

//删除对应的key
func (r *SkeletonMJRoom) DelRunningDeskKey(deskId int32) error {
	//删除和增加的时候 需要同步
	r.Lock()
	defer r.Unlock()

	runningKeys := r.GetRunningDeskeys()
	if runningKeys == nil {
		return nil
	}

	//删除key
	delIndex := -1
	for index, k := range runningKeys.Keys {
		if k == deskId {
			delIndex = index
		}
	}

	//删除对应的key
	if delIndex >= 0 {
		runningKeys.Keys = append(runningKeys.Keys[:delIndex], runningKeys.Keys[delIndex+1:]...)
	}
	r.SaveRunningDeskKeys(runningKeys)
	return nil
}

func (r *SkeletonMJRoom) GetRunningDeskeys() *majiang.RunningDeskKeys {
	keys := redisUtils.GetObj(majiang.REDIS_KEY_MJ_RUNNING, &majiang.RunningDeskKeys{})
	if keys == nil {
		return nil
	} else {
		return keys.(*majiang.RunningDeskKeys)
	}
}

func (r *SkeletonMJRoom) SaveRunningDeskKeys(keys *majiang.RunningDeskKeys) {
	redisUtils.SetObj(majiang.REDIS_KEY_MJ_RUNNING, keys)
}

func (r *SkeletonMJRoom) GetRoomMgr() api.MjRoomMgr {
	return r.RoomMnanger
}

//计算创建房间需要使用的费用
func (r *SkeletonMJRoom) CalcCreateFee(boardsCout int32) (int64) {
	var fee int64 = 0
	if boardsCout == 4 {
		fee = 2
	} else if boardsCout == 8 {
		fee = 3
	} else if boardsCout == 12 {
		fee = 5
	} else {
		return 5
	}
	return fee
}

//给room增加desk
func (r *SkeletonMJRoom) AddDesk(desk api.MjDesk) error {
	r.Desks = append(r.Desks, desk)
	return nil
}

// room 解散房间...解散朋友桌
func (r *SkeletonMJRoom) DissolveDesk(desk api.MjDesk, sendMsg bool) error {
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
