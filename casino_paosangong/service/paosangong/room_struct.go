package paosangong

import (
	"casino_common/proto/ddproto"
	"errors"
	"github.com/golang/protobuf/proto"
	"casino_common/utils/db"
	"casino_paosangong/conf/config"
	"casino_common/utils/chessUtils"
	"sync"
	"casino_common/common/log"
	"casino_common/common/service/roomAgent"
	"casino_common/common/Error"
	"encoding/json"
	"io/ioutil"
)

func init() {

}

type Room struct {
	*ddproto.NiuniuSrvRoom
	Desks []*Desk
	DeskLock sync.Mutex
}

//房间列表
type RoomList []*Room

var Rooms RoomList

//初始化
func InitRoomList() {
	//初始化房间
	Rooms = RoomList{}

	data, err := ioutil.ReadFile("../conf/niuniu_rooms.json")
	if err != nil {
		log.E("read ../conf/niuniu_rooms.json err:%v", err)
		return
	}
	err = json.Unmarshal(data, &Rooms)
	if err != nil {
		log.E("json unmarshal err:%v", err)
		log.E("%v", data)
		return
	}

	//恢复房间内的牌桌数据
	for _,room := range Rooms{
		room.Desks = []*Desk{}
		desk_list := room.GetSnapShotIdList()
		for _, desk_id := range desk_list.GetDeskId() {
			room.RecoveryDeskSnapShot(desk_id)
		}
	}

}

//通过房间id获取房间
func (list RoomList) GetRoomById(room_id int32) (*Room, error) {
	for _,room := range list {
		if room != nil && room.GetRoomId() == room_id {
			return room, nil
		}
	}
	return nil, errors.New("room not found.")
}

//查找用户
func (list RoomList) GetUserById(room_id int32, desk_id int32, user_id uint32) (*User, error) {
	room, err := list.GetRoomById(room_id)
	if err != nil {
		return nil, err
	}
	desk, err := room.GetDeskById(desk_id)
	if err != nil {
		return nil, err
	}
	user, err := desk.GetUserByUid(user_id)
	return user, err
}

//通过房号获得牌桌
func (room *Room) GetDeskByPassword(number string) (*Desk, error) {
	for _, desk := range room.Desks {
		if desk != nil && desk.GetDeskNumber() == number {
			return desk, nil
		}
	}
	return nil, errors.New("desk not found.")
}

//通过房间id获得牌桌
func (room *Room) GetDeskById(desk_id int32) (*Desk, error) {
	for _, desk := range room.Desks {
		if desk != nil && desk.GetDeskId() == desk_id {
			return desk, nil
		}
	}
	return nil, errors.New("desk not found.")
}

//获取房间内所有牌桌游戏的人数
func (room *Room) GetPlayerNum(bankRule ddproto.NiuniuEnumBankerRule) (num int32) {
	for _, desk := range room.Desks {
		num += int32(len(desk.Users))
	}
	return
}

//创建朋友桌房间
func (room *Room) CreateFriendDesk(desk_option *ddproto.NiuniuDeskOption, owner uint32) (*Desk, error) {
	//加锁
	room.DeskLock.Lock()
	defer room.DeskLock.Unlock()

	if desk_option == nil {
		log.E("desk option is nill.")
		return nil, errors.New("desk option is null.")
	}
	if owner <= 0 {
		log.E("create user is nil.")
		return nil, errors.New("user is nil.")
	}
	//验证配置正确性，及设默认值
	if desk_option.GetMinUser() < 2 {
		desk_option.MinUser = proto.Int32(2)
	}

	//自动开相关
	if desk_option.GetAutoStartGammer() >= 2 {
		*desk_option.MinUser = 2
	}

	if desk_option.GetMaxUser() < 2 {
		desk_option.MaxUser = proto.Int32(6)
	}

	if desk_option.GetMaxCircle() < 1 {
		desk_option.MaxCircle = proto.Int32(1)
	}

	if desk_option.GetBaseScore() <= 0 {
		desk_option.BaseScore = proto.Int64(1)
	}

	if desk_option.GetDiFen() <= 0 {
		desk_option.DiFen = proto.Int32(1)
	}

	//强制禁止中途加入
	//desk_option.DenyHalfJoin = proto.Bool(true)

	new_desk_id,err := db.GetNextSeq(config.DBT_NIUNIU_DESK)
	if err != nil {
		return nil, errors.New("get desk seq id fail.")
	}
	new_game_number,err := db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)
	if err != nil {
		return nil, errors.New("get gamenumber seq id fail.")
	}

	//生成房号
	desk_number := ""
	//如果为朋友桌，或者需要房号的金币场房间
	for {
		desk_number = chessUtils.GetRoomPass(int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN))
		//检查房号是否重复
		isExist := false
		for _, d := range room.Desks {
			if d.GetDeskNumber() == desk_number {
				isExist = true
				break
			}
		}
		//如果不存在重复的则跳出循环
		if isExist == false {
			break
		}
	}


	new_desk := &Desk{
		NiuniuSrvDesk: &ddproto.NiuniuSrvDesk{
			DeskId: proto.Int32(new_desk_id),
			DeskNumber: proto.String(desk_number),
			GameNumber: proto.Int32(new_game_number),
			RoomId: proto.Int32(room.GetRoomId()),
			LastWiner: proto.Uint32(0),
			Status: ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_ENTER.Enum(),
			DeskOption: desk_option,
			CircleNo: proto.Int32(1),
			Owner: proto.Uint32(owner),
			CurrBanker: proto.Uint32(owner),
			IsStart: proto.Bool(false),
			IsOnDissolve: proto.Bool(false),
			DissolveTime: proto.Int64(0),
			OneStartTime: proto.Int64(0),
			AllStartTime: proto.Int64(0),
			DaikaiUser: proto.Uint32(0),
			IsDaikai: proto.Bool(false),
			IsOnGamming: proto.Bool(false),
			IsCoinRoom: proto.Bool(desk_option.GetIsCoinRoom()),
			StartTime: proto.Int64(0),
		},
		Room: room,
		Users: []*User{},
	}

	//兼容金币场
	if desk_option.GetIsCoinRoom() {
		new_desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY.Enum()
	}

	room.Desks = append(room.Desks, new_desk)

	//新增快照索引
	new_desk.NewSnapShot()

	return new_desk, nil
}

//销毁房间
func (room *Room) RemoveFriendDesk(desk_id int32) error {

	desk, _ := room.GetDeskById(desk_id)
	//停止定时器
	if desk.DissolveTimer != nil {
		desk.DissolveTimer.Stop()
		desk.DissolveTimer = nil
	}
	if desk.QiangzhuangTimer != nil {
		desk.QiangzhuangTimer.Stop()
		desk.QiangzhuangTimer = nil
	}
	if desk.JiaBeiTimer != nil {
		desk.JiaBeiTimer.Stop()
		desk.JiaBeiTimer = nil
	}

	if desk.LiangpaiTimer != nil {
		desk.LiangpaiTimer.Stop()
		desk.LiangpaiTimer = nil
	}

	if desk.ReadyTimer != nil {
		desk.ReadyTimer.Stop()
		desk.ReadyTimer = nil
	}

	//发送牌局总结算
	desk.SendGameEndResultBc()

	//如果未开局，则返还房主/代开人房费
	if desk.GetCircleNo() == 1 && desk.GetStatus() <= ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START {
		create_user_id := desk.GetOwner()
		//如果是代开
		if desk.GetIsDaikai() {
			create_user_id = desk.GetDaikaiUser()
			//同步代开状态
			go func() {
				defer Error.ErrorRecovery("RemoveFriendDesk->roomAgent.DoDissolve()")
				err := roomAgent.DoDissolve(create_user_id, int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN), desk.GetDeskId())
				log.T("roomAgent.DoDissolve() err:%v", err)
			}()
		}
	}else {
		if desk.GetIsDaikai() {
			//同步代开状态
			go func() {
				defer Error.ErrorRecovery("RemoveFriendDesk->roomAgent.DoEnd()")
				roomAgent.DoEnd(desk.GetDaikaiUser(), int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN), desk.GetDeskId())
			}()
		}
	}

	//中途退出，更新全局统计
	if desk.GetCircleNo() > 1 && desk.GetCircleNo() < desk.DeskOption.GetMaxCircle() {
		//更新全局统计
		go func() {
			defer Error.ErrorRecovery("RemoveFriendDesk->InsertAllCounter()")
			desk.InsertAllCounter()
		}()
	}

	//清除用户session
	for _,u := range desk.Users {
		if u != nil {
			//清理session
			u.ClearSession()
			//停止旁观timer
			if u.AsideTimer != nil {
				u.AsideTimer.Stop()
			}
			u.AsideTimer = nil
		}
	}

	//加锁
	room.DeskLock.Lock()
	defer room.DeskLock.Unlock()

	//删除快照索引
	desk.RemoveSnapShot()

	i := 0
	for _,desk := range room.Desks {
		if desk != nil && desk.GetDeskId() == desk_id {
			room.Desks = append(room.Desks[:i], room.Desks[i+1:]...)
			log.T("RemoveFriendDesk() Err:成功从Room删除该牌桌 passwd:%v deskid:%v", desk.GetDeskNumber(), desk.GetDeskId())
			return nil
		}
		i++
	}
	log.E("RemoveFriendDesk()未找到该房间passwd:%v deskid:%v", desk.GetDeskNumber(), desk.GetDeskId())
	return errors.New("未找到该房间！")
}

var UserDeskMap sync.Map

//获取牌桌上的用户
func FindUserById(user_id uint32) (*User, error) {
	//先从缓存中去找
	if u, ok := UserDeskMap.Load(user_id); ok {
		//log.T("UserDeskMap[%d]", user_id)
		return u.(*User), nil
	}

	//如果缓存中没找到就去desk列表中查找，并更新缓存
	for _, room := range Rooms{
		for _,d := range room.Desks {
			for _,u := range d.Users {
				if u.GetUserId() == user_id {
					//更新缓存
					//log.T("d.Users[%d]", user_id)
					UserDeskMap.Store(user_id, u)
					return u, nil
				}
			}
		}
	}

	//log.E("朋友桌牌桌列表中未找到未找到用户%d", user_id)
	return nil, errors.New("未找到该用户。")
}
