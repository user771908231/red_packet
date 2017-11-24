package paoyao

import (
	"casino_common/proto/ddproto"
	"errors"
	"github.com/golang/protobuf/proto"
	"casino_common/utils/db"
	"casino_paoyao/conf/config"
	"casino_common/utils/chessUtils"
	"sync"
	"casino_common/common/userService"
	"casino_common/common/service/countService"
	"casino_common/common/log"
	"casino_common/common/service/roomAgent"
	"casino_common/common/Error"
)

func init() {

}

type Room struct {
	*ddproto.PaoyaoSrvRoom
	Desks []*Desk
	DeskLock sync.Mutex
}

//房间列表
type RoomList []*Room

var Rooms RoomList

//初始化
func InitRoomList() {
	//初始化房间
	Rooms = RoomList{
		//-----------------------------刨幺朋友桌-----------------------------
		&Room{
			PaoyaoSrvRoom: &ddproto.PaoyaoSrvRoom{
				RoomId: proto.Int32(0),
				RoomType: proto.Int32(int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND)),
				RoomLevel: proto.Int32(1),
				RoomTitle: proto.String("刨幺朋友桌"),
				BaseChip: proto.Int32(1),
				EnterCoin: proto.Int32(0),
			},
			Desks: []*Desk{},
		},
		//-------------------------金币场-刨幺抢庄-----------------------------
		&Room{
			PaoyaoSrvRoom: &ddproto.PaoyaoSrvRoom{
				RoomId: proto.Int32(1),
				RoomType: proto.Int32(int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN)),
				RoomLevel: proto.Int32(0),
				RoomTitle: proto.String("刨幺金币场-新手场"),
				BaseChip: proto.Int32(25),
				EnterCoin: proto.Int32(1000),
			},
			Desks: []*Desk{},
		},
		&Room{
			PaoyaoSrvRoom: &ddproto.PaoyaoSrvRoom{
				RoomId: proto.Int32(2),
				RoomType: proto.Int32(int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN)),
				RoomLevel: proto.Int32(1),
				RoomTitle: proto.String("刨幺金币场-低倍场"),
				BaseChip: proto.Int32(50),
				EnterCoin: proto.Int32(1500),
			},
			Desks: []*Desk{},
		},
		&Room{
			PaoyaoSrvRoom: &ddproto.PaoyaoSrvRoom{
				RoomId: proto.Int32(3),
				RoomType: proto.Int32(int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN)),
				RoomLevel: proto.Int32(2),
				RoomTitle: proto.String("刨幺金币场-中倍场"),
				BaseChip: proto.Int32(100),
				EnterCoin: proto.Int32(2500),
			},
			Desks: []*Desk{},
		},
		&Room{
			PaoyaoSrvRoom: &ddproto.PaoyaoSrvRoom{
				RoomId: proto.Int32(4),
				RoomType: proto.Int32(int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN)),
				RoomLevel: proto.Int32(3),
				RoomTitle: proto.String("刨幺金币场-高倍场"),
				BaseChip: proto.Int32(300),
				EnterCoin: proto.Int32(5000),
			},
			Desks: []*Desk{},
		},

	}

	//恢复房间内的牌桌数据
	for _,room := range Rooms{
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
		if desk != nil && desk.GetPwd() == number {
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

//创建朋友桌房间
func (room *Room) CreateFriendDesk(desk_option *ddproto.PaoyaoDeskOption, owner uint32) (*Desk, error) {
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
	switch desk_option.GetGammerNum() {
	case 2, 4:
	default:
		desk_option.GammerNum = proto.Int32(4)
	}

	//圈数
	switch desk_option.GetBoardsCout() {
	case 4, 8, 12:
	default:
		desk_option.BoardsCout = proto.Int32(4)
	}

	//筹码
	if desk_option.GetBaseChip() <= 0 {
		desk_option.BaseChip = proto.Int32(1)
	}

	//转幺模式
	switch desk_option.GetAllyType() {
	case 1, 2:
	default:
		desk_option.AllyType = proto.Int32(1)
	}

	new_desk_id,err := db.GetNextSeq(config.DBT_PAOYAO_DESK)
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
		desk_number = chessUtils.GetRoomPass(int32(ddproto.CommonEnumGame_GID_PAOYAO))
		//检查房号是否重复
		isExist := false
		for _, d := range room.Desks {
			if d.GetPwd() == desk_number {
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
		PaoyaoSrvDesk: &ddproto.PaoyaoSrvDesk{
			DeskId: proto.Int32(new_desk_id),
			Pwd: proto.String(desk_number),
			GameNumber: proto.Int32(new_game_number),
			RoomId: proto.Int32(room.GetRoomId()),
			Status: ddproto.PaoyaoEnumDeskStatus_PAOYAO_DESK_STATUS_WAIT_READY.Enum(),
			CircleNo: proto.Int32(1),
			CurrDeskScore: proto.Int32(0),
			Owner: proto.Uint32(owner),
			IsDaikai: proto.Bool(false),
			DaikaiUser: proto.Uint32(0),
			DeskOption: desk_option,
			LastActUser:proto.Uint32(0),
			LastChupaiUser:proto.Uint32(0),
			IsOnDissolve: proto.Bool(false),
			DissolveTime: proto.Int64(0),
			IsStart: proto.Bool(false),
			OneStartTime:proto.Int64(0),
			AllStartTime:proto.Int64(0),
			DissolveUser:proto.Uint32(0),
			IsCoinRoom: proto.Bool(false),
			SurplusTime: proto.Int32(0),
		},
		Room: room,
		Users: []*User{},
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

	//发送牌局总结算
	desk.SendGameEndResultBc()

	//如果未开局，则返还房主/代开人房费
	if desk.GetCircleNo() == 1 && desk.GetStatus() == ddproto.PaoyaoEnumDeskStatus_PAOYAO_DESK_STATUS_WAIT_READY {
		ownerFee := int64(GetOwnerFee(desk.DeskOption.GetBoardsCout()))
		create_user_id := desk.GetOwner()
		//如果是代开
		if desk.GetIsDaikai() {
			create_user_id = desk.GetDaikaiUser()
			//同步代开状态
			go func() {
				defer Error.ErrorRecovery("RemoveFriendDesk->roomAgent.DoDissolve()")
				err := roomAgent.DoDissolve(create_user_id, int32(ddproto.CommonEnumGame_GID_PAOYAO), desk.GetDeskId())
				log.T("roomAgent.DoDissolve() err:%v", err)
			}()
		}
		userService.INCRUserRoomcard(create_user_id, ownerFee, int32(ddproto.CommonEnumGame_GID_PAOYAO), "刨幺朋友桌，未开始游戏，解散房间归还房主房卡")
	}else {
		//已开局，则统计真实房卡消耗
		countService.AddFriendRoomCardTrueConsume(desk.GetOwner(), int64(GetOwnerFee(desk.DeskOption.GetBoardsCout())), int32(ddproto.CommonEnumGame_GID_PAOYAO))
		if desk.GetIsDaikai() {
			//同步代开状态
			go func() {
				defer Error.ErrorRecovery("RemoveFriendDesk->roomAgent.DoEnd()")
				roomAgent.DoEnd(desk.GetDaikaiUser(), int32(ddproto.CommonEnumGame_GID_PAOYAO), desk.GetDeskId())
			}()
		}
	}

	//中途退出，更新全局统计
	if desk.GetCircleNo() > 1 && desk.GetCircleNo() < desk.DeskOption.GetBoardsCout() {
		//更新全局统计
		go func() {
			defer Error.ErrorRecovery("RemoveFriendDesk->InsertAllCounter()")
			desk.InsertAllCounter()
		}()
	}

	//清除用户session
	for _,u := range desk.Users {
		if u != nil {
			u.ClearSession()
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
			log.T("RemoveFriendDesk() Err:成功从Room删除该牌桌 passwd:%v deskid:%v", desk.GetPwd(), desk.GetDeskId())
			return nil
		}
		i++
	}
	log.E("RemoveFriendDesk()未找到该房间passwd:%v deskid:%v", desk.GetPwd(), desk.GetDeskId())
	return errors.New("未找到该房间！")
}

var UserDeskMap map[uint32]*User = map[uint32]*User{}

//获取牌桌上的用户
func FindUserById(user_id uint32) (*User, error) {
	//先从缓存中去找
	if u, ok := UserDeskMap[user_id]; ok {
		log.T("UserDeskMap[%d]", user_id)
		return u, nil
	}

	//如果缓存中没找到就去desk列表中查找，并更新缓存
	for _, room := range Rooms{
		for _,d := range room.Desks {
			for _,u := range d.Users {
				if u.GetUserId() == user_id {
					//更新缓存
					log.T("d.Users[%d]", user_id)
					UserDeskMap[user_id] = u
					return u, nil
				}
			}
		}
	}

	return nil, errors.New("未找到该用户。")
}
