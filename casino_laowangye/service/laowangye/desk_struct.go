package laowangye

import (
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"errors"
	"github.com/name5566/leaf/gate"
	"fmt"
	"sync"
	"time"
	"casino_common/common/userService"
	"casino_common/common/service/whiteListService"
	"casino_common/proto/funcsInit"
)

type Desk struct {
	*Room
	*ddproto.LwySrvDesk
	Users            []*User
	UserLock         sync.Mutex
	ReqLock          sync.Mutex
	DissolveTimer    *time.Timer  //解散倒计时，定时器
	ReadyTimer       *time.Timer
	QiangzhuangTimer *time.Timer //抢庄倒计时，计时器
	YazhuTimer       *time.Timer //押注倒计时，计时器
	YaoshaiziTimer   *time.Timer //摇色子倒计时，计时器
}

//得到玩家列表，用户message
func (desk *Desk) GetUsers() []*ddproto.LwySrvUser {
	list := []*ddproto.LwySrvUser{}
	for _,u := range desk.Users {
		if u != nil {
			list = append(list, u.LwySrvUser)
		}
	}
	return list
}

//往牌桌加入用户
func (desk *Desk) AddUser(user_id uint32, agent gate.Agent) (*User, error) {
	//加锁
	desk.UserLock.Lock()
	defer desk.UserLock.Unlock()

	//如果已经在房间内，则直接加入房间
	ex_user,_ := desk.GetUserByUid(user_id)
	if ex_user != nil {
		return ex_user, nil
	}

	//是否有空余座位
	free_site_index,err := desk.GetFreeSiteIndex()
	if err != nil {
		return nil, errors.New("该房间人数已满！")
	}

	//是否允许中途加入
	if !desk.GetIsCoinRoom() && desk.DeskOption.GetDenyHalfJoin() && desk.GetIsStart() {
		//如果是朋友桌，且选项为不能中途加入，则开局后不能加入
		return nil, errors.New("该房间已开局，进房失败。")
	}

	//获取用户资料
	user_info := userService.GetUserById(user_id)
	if user_info == nil {
		return nil, errors.New("该用户不存在，进房失败。")
	}

	new_user := &User{
		Agent: agent,
		LwySrvUser: &ddproto.LwySrvUser{
			UserId: &user_id,
			Bill: &ddproto.LwyUserBill{
				UserId: proto.Uint32(user_id),
				Score: proto.Int64(0),
			},
			IsOnline: proto.Bool(true),
			Index: proto.Int32(int32(free_site_index)),
			BankerScore: proto.Int64(0),
			IsReady: proto.Bool(false),
			LastScore: proto.Int64(0),
			DissolveState: proto.Int32(0),
			NickName: user_info.NickName,
			IsOnWhiteList: proto.Bool(false),
			WhiteWinRate: proto.Int32(whiteListService.DefaultWinRate),
			IsOnGamming: proto.Bool(false),
			IsLeave: proto.Bool(false),
			IsXuanpai: proto.Bool(false),
			IsRobot: proto.Bool(false),
			AsideWatchTime:proto.Int64(time.Now().Unix()),
			TuizhuScore: proto.Int32(0),
			LastTuizhuCircleNo: proto.Int32(-1),
			Mai7Lost: proto.Int64(0),
			Mai8Lost: proto.Int64(0),
			Mai7: proto.Int64(0),
			Mai8: proto.Int64(0),
			Chi7: proto.Int64(0),
			Chi8: proto.Int64(0),
			ChizhuDetail: []*ddproto.LwySrvChizhuDetailItem{},
		},
		Desk: desk,
	}

	//刷新白名单
	new_user.CheckWhiteList()

	desk.Users =  append(desk.Users, new_user)

	//更新session
	new_user.UpdateSession()

	//保存切片
	desk.WipeSnapShot()
	return new_user, nil
}

//删除用户
func (desk *Desk) RemoveUser(user_id uint32) error {
	//加锁
	desk.UserLock.Lock()
	defer desk.UserLock.Unlock()

	var ex_user *User
	for i, u := range desk.Users {
		if u != nil && u.GetUserId() == user_id {
			u.Users = append(u.Users[:i], u.Users[i+1:]...)
			ex_user = u
			break
		}
	}

	if ex_user == nil {
		return errors.New("未找到该用户！")
	}

	//删除玩家后，清除session
	ex_user.ClearSession()

	//如果是机器人，则释放机器人
	if ex_user.GetIsRobot() {
		RobotManager.ReleaseRobots(ex_user.GetUserId())
	}

	//发送广播
	msg := &ddproto.CommonAckLeaveDesk{
		Header:     commonNewPorot.NewHeader(),
		UserId:     ex_user.UserId,
		IsExchange: proto.Bool(false),
	}
	*msg.Header.Code = 1
	*msg.Header.Error = "退出房间成功！"
	ex_user.WriteMsg(msg)
	ex_user.Desk.BroadCast(msg)

	return nil
}

//通过user_id查找用户
func (desk *Desk) GetUserByUid(user_id uint32) (*User, error) {
	for _,u := range desk.Users {
		if u != nil && u.GetUserId() == user_id {
			return u, nil
		}
	}
	return nil, errors.New("user not found.")
}

//通过index查找用户
func (desk *Desk) GetUserByIndex(site_index int32) *User {
	for _,u := range desk.Users {
		if u != nil && u.GetIndex() == site_index {
			return u
		}
	}
	return nil
}

//获取用户积分排名，从大道小
func (desk *Desk) GetUserRankByBill() []*User {
	users := make([]*User, len(desk.Users))
	copy(users, desk.Users)

	for i:=0;i<len(users)-1;i++ {
		for j:=i+1;j<len(users);j++ {
			if users[j].Bill.GetScore() > users[i].Bill.GetScore() {
				tmp := users[i]
				users[i] = users[j]
				users[j] = tmp
			}
		}
	}
	return users
}

//获取可吃的用户
func (user *User) GetCanChiYours() (list []*User) {
	owner,_ := user.Desk.GetUserByUid(user.Desk.GetOwner())
	list = append(list, owner)
	for _,u := range user.Users {
		if u == nil || !u.GetIsOnGamming() || u.IsOwner() || u.GetUserId() == user.GetUserId() {
			continue
		}
		list = append(list, u)
	}
	return
}

//是否都准备了-除了房主
func (desk *Desk) IsAllReadyExcludeOwner() error {
	var i int32 = 0
	if desk.GetStatus() > ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY {
		return errors.New("该房间正在游戏中！")
	}
	for _,u := range desk.Users {
		//如果不是房主
		if u != nil && !u.IsOwner() {
			if u.GetIsReady() == true {
				i++
			}else {
				return errors.New(fmt.Sprintf("用户%d未准备！", u.GetUserId()))
			}
		}
	}
	if i >= desk.DeskOption.GetMinUser() - 1 {
		return nil
	}
	return errors.New(fmt.Sprintf("未达到%d人最小开局条件！", desk.DeskOption.GetMinUser()))
}

//是否都准备了
func (desk *Desk) IsAllReady() error {
	var i int32 = 0
	if desk.GetStatus() > ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY {
		return errors.New("该房间正在游戏中！")
	}
	for _,u := range desk.Users {
		//如果不是房主
		if u != nil {
			if u.GetIsReady() == true {
				i++
			}else {
				return errors.New(fmt.Sprintf("用户%d未准备！", u.GetUserId()))
			}
		}
	}

	if i >= desk.DeskOption.GetMinUser() - 1 {
		return nil
	}
	return errors.New(fmt.Sprintf("未达到%d人最小开局条件！", desk.DeskOption.GetMinUser()))
}

//是否所有人都已经抢庄
func (desk *Desk) IsAllQiangzhuang() error {
	if desk.GetStatus() != ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_QIANGZHUANG {
		return errors.New("调用IsAllQiangzhuang非法！")
	}
	var i int32 = 0
	for _,u := range desk.Users {
		//如果不是房主
		if u != nil {
			if u.GetIsOnGamming() && u.GetBankerScore() == 0 {
				return errors.New(fmt.Sprintf("用户%d未抢庄！", u.GetUserId()))
			}else {
				i++
			}
		}
	}

	if i <= 1 {
		return errors.New(fmt.Sprintf("牌桌%d未达到抢庄人数下限2人！（防止panic）"))
	}

	return nil
}


//客户端牌桌
func (desk *Desk) GetClientDesk() *ddproto.LwyClientDesk {
	users := []*ddproto.LwyClientUser{}

	for _,u := range desk.Users {
		if u == nil {
			continue
		}
		users =append(users, u.GetClientUser())
	}

	client_desk := &ddproto.LwyClientDesk{
		DeskId: proto.Int32(desk.GetDeskId()),
		Password: proto.String(desk.GetPassword()),
		GameNumber: proto.Int32(desk.GetGameNumber()),
		RoomId: proto.Int32(desk.GetRoomId()),
		Status: desk.Status,
		CircleNo: proto.Int32(desk.GetCircleNo()),
		Owner: proto.Uint32(desk.GetOwner()),
		CurrBanker: proto.Uint32(desk.GetCurrBanker()),
		DeskOption: desk.DeskOption,
		Users: users,
		IsOnDissolve: desk.IsOnDissolve,
		DissolveTime: proto.Int32(0),
		IsDaikai: proto.Bool(desk.GetIsDaikai()),
		DaikaiUser: proto.Uint32(desk.GetDaikaiUser()),
		IsCoinRoom: proto.Bool(desk.GetIsCoinRoom()),
		IsOnGamming: proto.Bool(desk.GetIsOnGamming()),
		SurplusTime: proto.Int32(1),
	}

	//投票剩余时间
	if desk.GetIsOnDissolve() && desk.GetDissolveTime() > 0 {
		last_time := time.Now().Unix() - desk.GetDissolveTime()
		client_desk.DissolveTime = proto.Int32(int32(last_time))
	}

	return client_desk
}

//插入比赛记录表
func (desk *Desk)InsertOneCounter() error {
	return nil
	//users_str := ""
	//for _, u := range desk.Users {
	//	users_str += fmt.Sprintf(",%d", u.GetUserId())
	//}
	//
	//oneCounter := model.T_niuniu_desk_round{
	//	BeginTime: time.Unix(desk.GetOneStartTime(), 0),
	//	DeskId: desk.GetDeskId(),
	//	UserIds: users_str,
	//	GameNumber: desk.GetGameNumber(),
	//	Password: desk.GetPassword(),
	//	EndTime: time.Now(),
	//	IsCoinRoom: desk.GetIsCoinRoom(),
	//	BankRule: desk.DeskOption.GetBankRule(),
	//}
	//
	//for _, u := range desk.Users {
	//	//1局记录
	//	oneCounter.Records = append(oneCounter.Records, model.LwyRecordBean{
	//		UserId:    u.GetUserId(),
	//		NickName:  u.GetNickName(),
	//		WinAmount: int64(u.GetLastScore()),
	//	})
	//
	//	//game_bill
	//	userGameBillService.Insert(userGameBillService.T_user_game_bill{
	//		UserId: u.GetUserId(),
	//		GameId: int32(ddproto.CommonEnumGame_GID_LAOWANGYE),
	//		RoomType: int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND),
	//		DeskId: desk.GetDeskId(),
	//		Password: desk.GetPassword(),
	//		GameNumber: desk.GetGameNumber(),
	//		WinAmount: int64(u.GetLastScore()),
	//		Balance: int64(u.Bill.GetScore()),
	//		IsRobot: u.GetIsRobot(),
	//		EndTime: time.Now(),
	//	})
	//}
	//
	//if !desk.GetIsCoinRoom(){
	//	//插入朋友桌记录
	//	db.Log(tableName.DBT_NIU_DESK_ROUND_ONE).Insert(oneCounter)
	//}else {
	//	//插入金币场记录
	//	db.Log(tableName.DBT_NIU_DESK_ROUND_ONE_COIN).Insert(oneCounter)
	//}
	//
	//return nil
}

//插入10局比赛记录
func (desk *Desk)InsertAllCounter() error {
	return nil
	//users_str := ""
	//for _, u := range desk.Users {
	//	users_str += fmt.Sprintf(",%d", u.GetUserId())
	//}
	//
	//allCounter := model.T_niuniu_desk_round{
	//	BeginTime: time.Unix(desk.GetAllStartTime(), 0),
	//	DeskId: desk.GetDeskId(),
	//	UserIds: users_str,
	//	TotalRound: desk.DeskOption.GetBoardsCout(),
	//	GameNumber: desk.GetGameNumber(),
	//	Password: desk.GetPassword(),
	//	EndTime: time.Now(),
	//}
	//
	//for _, u := range desk.Users {
	//	//1局记录
	//	allCounter.Records = append(allCounter.Records, model.LwyRecordBean{
	//		UserId:    u.GetUserId(),
	//		NickName:  u.GetNickName(),
	//		WinAmount: int64(u.Bill.GetScore()),
	//	})
	//}
	//
	//return db.Log(tableName.DBT_NIU_DESK_ROUND_ALL).Insert(allCounter)
}

//获取房间简介
func (desk *Desk) GetTips() string {
	room_type := "老王爷"
	rule_type := "房主定庄"
	switch desk.DeskOption.GetBankRule() {
	case ddproto.LwyEnumBankerRule_LWY_FANGZHU_DINGZHUANG:
		rule_type = "房主定庄"
	case ddproto.LwyEnumBankerRule_LWY_LUN_LIU_ZUO_ZHUANG:
		rule_type = "轮流坐庄"
	case ddproto.LwyEnumBankerRule_LWY_QIANG_ZHUANG:
		rule_type = "自由抢庄"
	case ddproto.LwyEnumBankerRule_LWY_SUI_JI_ZUO_ZHUANG:
		rule_type = "随机坐庄"
	}

	deny_half_join := "允许中途加入"
	if desk.DeskOption.GetDenyHalfJoin() {
		deny_half_join = "禁止中途加入"
	}

	tips := fmt.Sprintf("%s %s %d人%d局 ", room_type, rule_type, desk.DeskOption.GetMaxUser(), desk.DeskOption.GetBoardsCout(), deny_half_join)

	return tips
}

//获取空余的位置索引(座位号从1开始)
func (desk *Desk)GetFreeSiteIndex() (int, error)  {
	for i:=0;i<int(desk.DeskOption.GetMaxUser());i++ {
		has_gammer := false
		for _,u := range desk.Users {
			if int(u.GetIndex()) == i+1 {
				has_gammer = true
				break
			}
		}
		if has_gammer == false {
			return i+1, nil
		}
	}
	return 0, errors.New("not free site")
}

//是否是空位置
func (desk *Desk) IsFreeSite(site_index int32) bool {
	for _,u := range desk.Users {
		if u == nil {
			continue
		}
		if u.GetIndex() == site_index {
			return false
		}
	}

	return true
}

//获取已入座的人数
func (desk *Desk) GetSitedGammerNum() (num int32) {
	for _,u := range desk.Users {
		if u == nil {
			continue
		}
		if u.GetIndex() >= 0 {
			num++
		}
	}
	return
}
