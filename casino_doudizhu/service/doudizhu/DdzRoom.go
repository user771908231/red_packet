package doudizhu

import (
	"sync"
	"casino_server/utils/numUtils"
	"casino_server/utils"
	"casino_server/common/log"
	"casino_doudizhu/service/lock"
	"errors"
)

//初始化一个斗地主的房间实例
var FDdzRoomIns DdzRoom = new(DdzRoom)

//斗地主的房间
type  DdzRoom struct {
	sync.Mutex
	Desks []*DdzDesk
}

func (room *DdzRoom) CreateDesk() *DdzDesk {
	//创建一个desk

	//1,得到一个key
	key := room.NewRoomKey()

	//2, newDesk and 赋值
	desk := NewDdzDesk()
	desk.key = key

	return desk

}

//得到一个roomKey
func (r *DdzRoom) NewRoomKey() string {
	a := utils.Rand(100000, 1000000)
	roomKey, _ := numUtils.Int2String(a)
	//1,判断roomKey是否已经存在
	if r.IsRoomKeyExist(roomKey) {
		//log.E("房间密钥[%v]已经存在,创建房间失败,重新创建", roomKey)
		return r.NewRoomKey()
	} else {
		//log.T("最终得到的密钥是[%v]", roomKey)
		return roomKey
	}
	return ""
}

//判断roomkey是否已经存在了
func (r *DdzRoom) IsRoomKeyExist(roomkey string) bool {
	ret := false
	for i := 0; i < len(r.Desks); i++ {
		d := r.Desks[i]
		if d != nil && d.key == roomkey {
			ret = true
			break
		}
	}
	return ret
}



// room 解散房间...
func (r *DdzRoom)DissolveDesk(desk *DdzDesk, sendMsg bool) error {
	//清楚数据,1,session相关。2,
	log.T("开始解散desk[%v]...", desk.DeskId)
	log.T("开始解散desk[%v]user的session数据...", desk.DeskId)
	for _, user := range desk.Users {
		if user != nil {
			user.ClearAgentGameData()
		}
	}

	log.T("开始删除desk[%v]...", desk.DeskId)

	//发送解散房间的广播
	rmErr := r.RmDesk(desk)
	if rmErr != nil {
		log.E("删除房间失败,errmsg[%v]", rmErr)
		return rmErr
	}

	//删除锁
	lock.DelDeskLock(desk.DeskId)
	//删除reids
	DelMjDeskRedis(desk)

	//删除房间
	log.T("删除desk[%v]之后，发送删除的广播...", desk.DeskId)
	if sendMsg {
		//发送解散房间的广播
	}

	return nil

}

func (r *DdzRoom) RmDesk(desk *DdzDesk) error {
	index := -1
	for i, d := range r.Desks {
		if d != nil && d.DeskId == desk.DeskId {
			index = i
			break
		}
	}

	if index >= 0 {
		r.Desks = append(r.Desks[:index], r.Desks[index + 1:]...)
		return nil
	} else {
		return errors.New("删除失败，没有找到对应的desk")
	}

}

func GetFDdzRoom() *DdzRoom {
	//暂时返回朋友桌
	return FDdzRoomIns
}
//通过用户的session 找到mjroom
func GetMjroomBySession(userId uint32) *DdzRoom {
	session := GetSession(userId)
	if session == nil {
		return nil
	}

	//目前暂时返回一个房间，方便测试 todo
	return FDdzRoomIns

}

func (r *DdzRoom)GetDeskByDeskId(deskId int32) *DdzDesk {
	return nil
}


