package majiang

import (
	"casino_majiang/conf/log"
	"casino_server/utils"
	"casino_server/utils/numUtils"
	"casino_majiang/msg/protogo"
	"errors"
	"github.com/name5566/leaf/gate"
)


//普通的麻将房间...
func init() {
	FMJRoomIns = NewMjRoom()
	FMJRoomIns.OnInit()
}

var FMJRoomIns *MjRoom

//初始化
func (r *MjRoom) OnInit() {
	log.T("初始化麻将的room")

}

func (r *MjRoom) CreateDesk(m *mjproto.Game_CreateRoom) *MjDesk {
	//create 的时候，是否需要通过type 来判断,怎么样创建房间

	//1,创建一个房间，并初始化参数
	desk := NewMjDesk()
	*desk.Password = r.RandRoomKey()
	*desk.Owner = m.GetHeader().GetUserId()        //设置房主
	*desk.CardsNum = m.GetRoomTypeInfo().GetCardsNum()
	//把创建的desk加入到room中
	r.AddDesk(desk)
	return desk
}

func (r *MjRoom) RandRoomKey() string {
	a := utils.Rand(100000, 1000000)
	roomKey, _ := numUtils.Int2String(a)
	//1,判断roomKey是否已经存在
	if r.IsRoomKeyExist(roomKey) {
		log.E("房间密钥[%v]已经存在,创建房间失败,重新创建", roomKey)
		return r.RandRoomKey()
	} else {
		log.T("最终得到的密钥是[%v]", roomKey)
		return roomKey
	}
	return ""
}


//判断roomkey是否已经存在了
func (r *MjRoom) IsRoomKeyExist(roomkey string) bool {
	ret := false
	for i := 0; i < len(r.Desks); i++ {
		d := r.Desks[i]
		if d != nil && d.GetPassword() == roomkey {
			ret = true
			break
		}
	}

	return ret
}

//通过房间号码得到desk
func (r *MjRoom) GetDeskByPassword(key string) *MjDesk {
	//如果找到对应的房间，则返回
	for _, d := range r.GetDesks() {
		if d != nil && d.GetPassword() == key {
			return d
		}
	}

	//如果没有找到，则返回nil
	return nil
}

//进入房间
//进入的时候，需要判断牌房间的类型...
func (r *MjRoom) EnterRoom(key string, userId uint32, a gate.Agent) (*MjDesk, error) {
	var desk *MjDesk
	//如果是朋友桌,需要通过房间好来找到desk
	if r.IsFriend() {
		desk = r.GetDeskByPassword(key)
		if desk == nil {
			return nil, errors.New("没有找到对应的desk")
		} else {
			err := desk.addNewUserFriend(userId, a)
			if err != nil {
				//用户加入房间失败...
				return desk, errors.New("用户加入房间失败...")
			}
		}
	}

	//如果是锦标赛

	//返回结果...
	return nil, nil
}

func (r *MjRoom) IsFriend() bool {
	return true
}

func (r *MjRoom) AddDesk(desk *MjDesk) error {
	r.Desks = append(r.Desks, desk)
	return nil
}

func GetMJRoom() *MjRoom {
	//暂时返回朋友桌
	return FMJRoomIns
}
//通过用户的session 找到mjroom
func GetMjroomBySession(userId uint32) *MjRoom {
	session := GetSession(userId)
	if session == nil {
		return nil
	}

	session.GetRoomId()
	session.GetDeskId()

	//目前暂时返回一个房间，方便测试 todo
	return FMJRoomIns

}

func GetMjDeskBySession(userId uint32) *MjDesk {
	//得到session
	session := GetSession(userId)
	if session == nil {
		return nil
	}

	room := GetMjroomBySession(userId)
	for _, d := range room.Desks {
		if d != nil && d.GetDeskId() == session.GetDeskId() {
			return d
		}
	}

	return nil
}




