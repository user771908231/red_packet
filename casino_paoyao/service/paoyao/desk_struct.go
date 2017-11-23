package paoyao

import (
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"errors"
	"github.com/name5566/leaf/gate"
	"fmt"
	"sync"
	"time"
	"casino_common/common/userService"
	"casino_common/common/model"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"casino_common/common/service/roomAgent"
	//"casino_common/common/service/whiteListService"
	"casino_common/proto/funcsInit"
	"casino_common/common/service/userGameBillService"
	"casino_common/utils/chessUtils"
	"casino_common/common/log"
)

type Desk struct {
	*Room
	*ddproto.PaoyaoSrvDesk
	Users      []*User
	UserLock   sync.Mutex
	ReqLock    sync.Mutex
	DissolveTimer *time.Timer  //解散倒计时，定时器
	ReadyTimer *time.Timer
	QiangzhuangTimer *time.Timer  //抢庄倒计时，计时器
	JiaBeiTimer *time.Timer  //加倍倒计时，计时器
}

//得到玩家列表，用户message
func (desk *Desk) GetUsers() []*ddproto.PaoyaoSrvUser {
	list := []*ddproto.PaoyaoSrvUser{}
	for _,u := range desk.Users {
		if u != nil {
			list = append(list, u.PaoyaoSrvUser)
		}
	}
	return list
}

//往牌桌加入用户
func (desk *Desk) AddUser(user_id uint32, agent gate.Agent) (*User, error) {
	//加锁
	desk.UserLock.Lock()
	defer desk.UserLock.Unlock()

	//如果该用户已在房间中，则让他直接进房
	var user_num int32 = 0
	for _,u := range desk.Users {
		if u != nil {
			if u.GetUserId() == user_id {
				return u, nil
			}else {
				user_num++
			}
		}
	}

	free_site_index,err := desk.GetFreeSiteIndex()
	if user_num >= desk.GetDeskOption().GetGammerNum() || err != nil {
		return nil, errors.New("该房间人数已满！")
	}

	//是否允许中途加入
	if !desk.GetIsCoinRoom() && desk.GetIsStart() {
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
		PaoyaoSrvUser: &ddproto.PaoyaoSrvUser{
			UserId: &user_id,
			Bill: &ddproto.PaoyaoUserBill{
				UserId: proto.Uint32(user_id),
				Score: proto.Int32(0),
				CountWin: proto.Int32(0),
				CountLost: proto.Int32(0),
			},
			IsOnline: proto.Bool(true),
			Index: proto.Int32(int32(free_site_index+1)),
			Pokers: &ddproto.PaoyaoSrvPoker{
				Pais: []*ddproto.CommonSrvPokerPai{},
			},
			OutPai: nil,
			IsPass:proto.Bool(false),
			DeskScore:proto.Int32(0),
			IsReady: proto.Bool(false),
			LastScore: proto.Int32(0),
			WxInfo: &ddproto.WeixinInfo{
				OpenId: user_info.OpenId,
				NickName: user_info.NickName,
				HeadUrl: proto.String(userService.GetUserHeadImg(user_info)),
				Sex:user_info.Sex,
				City: user_info.City,
				UnionId:user_info.UnionId,
			},
			DissolveState: proto.Int32(0),
			IsRobot:proto.Bool(false),
			IsOnWhiteList:proto.Bool(false),
			WhiteWinRate:proto.Int32(0),
			IsLeave:proto.Bool(false),
		},
		Desk: desk,
	}

	//刷新白名单
	new_user.CheckWhiteList()

	//如果是代开,且该用户第一个进入房间则设置该用户为房主
	if !desk.GetIsCoinRoom() && desk.GetIsDaikai() {
		if len(desk.Users) == 0 {
			desk.Owner = proto.Uint32(new_user.GetUserId())
		}
		//同步代开状态
		roomAgent.DoAddUser(desk.GetDaikaiUser(), int32(ddproto.CommonEnumGame_GID_PAOYAO), desk.GetDeskId(), new_user.WxInfo.GetNickName())
	}
	desk.Users =  append(desk.Users, new_user)

	//更新session
	new_user.UpdateSession()

	//发送进房广播
	new_user.SendEnterDeskBC()

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
		//RobotManager.ReleaseRobots(ex_user.GetUserId())
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

//给游戏中用户发牌
func (desk *Desk) DoSendPoker() error {
	//洗牌
	rand_index := chessUtils.Xipai(0, 108)
	for i,v := range rand_index {
		if v > 53 {
			rand_index[i] -= 54
		}
	}

	//发牌
	for _,u := range desk.Users {
		if u == nil {
			continue
		}
		pais_index := rand_index[:27]
		rand_index = rand_index[27:]

		u.Pokers = ParseOutPai(pais_index)
	}

	//发牌overturn
	desk.SendFapaiOt()

	return nil
}

//获取下个说话的人
func (desk *Desk) GetNextChupaiUser() *User {
	//第一局让房主
	if desk.GetCircleNo() == 1 || desk.GetLastActUser() == 0 {
		next_user,_ := desk.GetUserByUid(desk.GetOwner())
		return next_user
	}

	//让下一个玩家出牌
	last_user,_ := desk.GetUserByUid(desk.GetLastActUser())
	next_user := desk.Users[int(last_user.GetIndex()+1)%len(desk.Users)]

	return next_user
}


//是否都准备了-除了房主
func (desk *Desk) IsAllReadyExcludeOwner() error {
	var i int32 = 0
	if desk.GetStatus() != ddproto.PaoyaoEnumDeskStatus_PAOYAO_DESK_STATUS_WAIT_READY {
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
	if i >= 4 {
		return nil
	}
	return errors.New(fmt.Sprintf("未达到%d人最小开局条件！", 4))
}
//是否都准备了
func (desk *Desk) IsAllReady() error {
	var i int32 = 0
	if desk.GetStatus() != ddproto.PaoyaoEnumDeskStatus_PAOYAO_DESK_STATUS_WAIT_READY {
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
	if i >= desk.DeskOption.GetGammerNum() {
		return nil
	}
	return errors.New(fmt.Sprintf("未达到%d人最小开局条件！", 4))
}


//客户端牌桌，参数表示是否显示某人的牌
func (user *User) GetClientDesk() *ddproto.PaoyaoClientDesk {
	if user == nil {
		log.E("user is nil")
		return nil
	}
	users := []*ddproto.PaoyaoClientUser{}

	for _,u := range user.Users {
		if u == nil {
			continue
		}
		new_user := u.GetClientUser()
		//显示某个人的牌
		if  u.GetUserId() == user.GetUserId() {
			new_user.Pokers = GetClientPoker(u.Pokers)
		}
		users =append(users, new_user)
	}

	client_desk := &ddproto.PaoyaoClientDesk{
		DeskId:       proto.Int32(user.GetDeskId()),
		Pwd:          proto.String(user.GetPwd()),
		GameNumber:   proto.Int32(user.GetGameNumber()),
		RoomId:       proto.Int32(user.GetRoomId()),
		Status:       user.Status,
		CircleNo:     proto.Int32(user.GetCircleNo()),
		Owner:        proto.Uint32(user.GetOwner()),
		DeskOption:   user.DeskOption,
		PlayerInfo:   users,
		IsStart:      user.IsStart,
		IsOnDissolve: user.IsOnDissolve,
		DissolveTime: proto.Int64(0),
		IsDaikai:     proto.Bool(user.GetIsDaikai()),
		DaikaiUser:   proto.Uint32(user.GetDaikaiUser()),
		IsCoinRoom:   proto.Bool(user.GetIsCoinRoom()),
		SurplusTime:  proto.Int32(1),
	}

	//投票剩余时间
	if user.GetIsOnDissolve() && user.GetDissolveTime() > 0 {
		last_time := time.Now().Unix() - user.GetDissolveTime()
		client_desk.DissolveTime = proto.Int64(int64(last_time))
	}

	return client_desk
}

//房费计算
func GetOwnerFee(circle int32) int32 {
	//4局 1张 8局 2张 12局 3张
	switch  {
	case circle > 0 && circle <= 4 :
		return 1
	case circle > 4 && circle <= 8:
		return 2
	case circle > 8:
		return 3
	default:
		return 0
	}
}

//插入比赛记录表
func (desk *Desk)InsertOneCounter() error {
	return nil
	users_str := ""
	for _, u := range desk.Users {
		users_str += fmt.Sprintf(",%d", u.GetUserId())
	}

	oneCounter := model.T_niuniu_desk_round{
		BeginTime: time.Unix(desk.GetOneStartTime(), 0),
		DeskId: desk.GetDeskId(),
		UserIds: users_str,
		GameNumber: desk.GetGameNumber(),
		Password: desk.GetPwd(),
		EndTime: time.Now(),
	}

	for _, u := range desk.Users {
		//1局记录
		oneCounter.Records = append(oneCounter.Records, model.NiuRecordBean{
			UserId:    u.GetUserId(),
			NickName:  u.WxInfo.GetNickName(),
			WinAmount: int64(u.GetLastScore()),
		})

		//game_bill
		userGameBillService.Insert(userGameBillService.T_user_game_bill{
			UserId: u.GetUserId(),
			GameId: int32(ddproto.CommonEnumGame_GID_PAOYAO),
			RoomType: int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND),
			DeskId: desk.GetDeskId(),
			Password: desk.GetPwd(),
			GameNumber: desk.GetGameNumber(),
			WinAmount: int64(u.GetLastScore()),
			Balance: int64(u.Bill.GetScore()),
			IsRobot: u.GetIsRobot(),
			EndTime: time.Now(),
		})
	}

	db.Log(tableName.DBT_NIU_DESK_ROUND_ONE).Insert(oneCounter)

	return nil
}

//插入10局比赛记录
func (desk *Desk)InsertAllCounter() error {
	return nil
	users_str := ""
	for _, u := range desk.Users {
		users_str += fmt.Sprintf(",%d", u.GetUserId())
	}

	allCounter := model.T_niuniu_desk_round{
		BeginTime: time.Unix(desk.GetAllStartTime(), 0),
		DeskId: desk.GetDeskId(),
		UserIds: users_str,
		TotalRound: desk.DeskOption.GetBoardsCout(),
		GameNumber: desk.GetGameNumber(),
		Password: desk.GetPwd(),
		EndTime: time.Now(),
	}

	for _, u := range desk.Users {
		//1局记录
		allCounter.Records = append(allCounter.Records, model.NiuRecordBean{
			UserId:    u.GetUserId(),
			NickName:  u.WxInfo.GetNickName(),
			WinAmount: int64(u.Bill.GetScore()),
		})
	}

	return db.Log(tableName.DBT_NIU_DESK_ROUND_ALL).Insert(allCounter)
}

//获取房间简介
func (desk *Desk) GetTips() string {
	room_type := "急速刨幺"
	if desk.DeskOption.GetHasAnimation() {
		room_type = "经典刨幺"
	}

	tips := fmt.Sprintf("%s %d人%d局", room_type, desk.DeskOption.GetBoardsCout(), desk.DeskOption.GetBoardsCout())

	return tips
}

//获取空余的位置索引
func (desk *Desk)GetFreeSiteIndex() (int, error)  {
	for i:=0;i<int(desk.DeskOption.GetGammerNum());i++ {
		has_gammer := false
		for _,u := range desk.Users {
			if int(u.GetIndex()) == i+1 {
				has_gammer = true
				break
			}
		}
		if has_gammer == false {
			return i, nil
		}
	}
	return 0, errors.New("not free site")
}
