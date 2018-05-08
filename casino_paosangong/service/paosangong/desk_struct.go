package paosangong

import (
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"errors"
	"github.com/name5566/leaf/gate"
	"casino_common/utils/chessUtils"
	"fmt"
	"sync"
	"time"
	"casino_common/common/userService"
	"casino_common/common/model"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	//"casino_common/common/service/roomAgent"
	"casino_common/common/log"
	"math/rand"
	"casino_common/common/service/whiteListService"
	"casino_common/proto/funcsInit"
	"casino_common/common/service/userGameBillService"
)

type Desk struct {
	*Room
	*ddproto.NiuniuSrvDesk
	Users      []*User
	UserLock   sync.Mutex
	ReqLock    sync.Mutex
	DissolveTimer *time.Timer  //解散倒计时，定时器
	ReadyTimer *time.Timer
	QiangzhuangTimer *time.Timer  //抢庄倒计时，计时器
	JiaBeiTimer *time.Timer  //加倍倒计时，计时器
	LiangpaiTimer *time.Timer  //亮牌计时，计时器
}

//得到玩家列表，用户message
func (desk *Desk) GetUsers() []*ddproto.NiuniuSrvUser {
	list := []*ddproto.NiuniuSrvUser{}
	for _,u := range desk.Users {
		if u != nil {
			list = append(list, u.NiuniuSrvUser)
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

	//free_site_index,err := desk.GetFreeSiteIndex()
	if user_num >= desk.GetDeskOption().GetMaxUser() {
		return nil, errors.New("该房间人数已满！")
	}

	//是否允许中途加入
	if desk.DeskOption.GetDenyHalfJoin() && desk.GetIsStart() {
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
		NiuniuSrvUser: &ddproto.NiuniuSrvUser{
			UserId: &user_id,
			Bill: &ddproto.NiuniuUserBill{
				UserId: proto.Uint32(user_id),
				Score: proto.Int64(0),
				CountHasNiu: proto.Int32(0),
				CountNoNiu: proto.Int32(0),
				CountWin: proto.Int32(0),
				CountLost: proto.Int32(0),
			},
			IsOnline: proto.Bool(true),
			Index: proto.Int32(int32(-1)),
			Pokers: &ddproto.NiuniuSrvPoker{
				Pais: []*ddproto.CommonSrvPokerPai{},
				Type: ddproto.NiuniuEnum_PokerType_NO_NIU.Enum(),
			},
			BankerScore: proto.Int64(0),
			DoubleScore: proto.Int64(0),
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
			TuizhuScore: proto.Int32(0),
			LastTuizhuCircleNo: proto.Int32(-1),
			IsLiangpai: proto.Bool(false),
			IsTuoguan: proto.Bool(false),
			TuoGuanOpt: nil,
		},
		Desk: desk,
	}

	//刷新白名单
	new_user.CheckWhiteList()

	//如果是代开,且该用户第一个进入房间则设置该用户为房主
	//if !desk.GetIsCoinRoom() && desk.GetIsDaikai() {
	//	if len(desk.Users) == 0 {
	//		desk.Owner = proto.Uint32(new_user.GetUserId())
	//		desk.CurrBanker = proto.Uint32(new_user.GetUserId())
	//	}
	//	//同步代开状态
	//	roomAgent.DoAddUser(desk.GetDaikaiUser(), int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN), desk.GetDeskId(), new_user.GetNickName())
	//}
	desk.Users =  append(desk.Users, new_user)

	//更新session
	new_user.UpdateSession()

	//发送进房广播
	//new_user.SendEnterDeskBC()

	//保存切片
	desk.WipeSnapShot()

	//匹配场刷新旁观timer
	if new_user.GetRoomId() > 0 {
		new_user.RefreshAsideTimer()
	}
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

	//停止旁观timer
	if ex_user.AsideTimer != nil {
		ex_user.AsideTimer.Stop()
	}
	ex_user.AsideTimer = nil

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

//给游戏中用户发牌
func (desk *Desk) DoSendPoker() error {
	// 发牌
	rand_index := chessUtils.Xipai(0, 52)

	//如果没有花牌，则把j、q、k过滤掉（移动到尾部）
	if desk.GetDeskOption().GetHasFlower() == false {
		filter_index := []int32{}
		for _,index := range rand_index {
			switch index {
			case 8, 9, 10, 21, 22, 23, 34, 35, 36, 47, 48, 49:
			default:
				filter_index = append(filter_index, index)
			}
		}
		rand_index = filter_index
	}

	for i,u := range desk.Users {
		if u != nil {
			//如果没在游戏中,则给玩家发一副空牌
			if !u.GetIsOnGamming() {
				u.Pokers = &ddproto.NiuniuSrvPoker{
					Pais: []*ddproto.CommonSrvPokerPai{},
					Type: ddproto.NiuniuEnum_PokerType_NO_NIU.Enum(),
					SelectedId: []int32{},
				}
				continue
			}

			pais_index := rand_index[:5]
			rand_index = rand_index[5:]
			pais := []*ddproto.CommonSrvPokerPai{}
			for _,index := range pais_index {
				pais = append(pais, PokerList[index])
			}
			user_poker := ParseNiuPoker(pais, desk.GetDeskOption())

			//如果出现一条龙、牛牛等牌型则重新洗牌
			rand_num := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
			if poker_rand,ok := Exchange_poker_score_map[user_poker.GetType()];ok && rand_num<poker_rand && len(rand_index)-((len(desk.Users)-i+1)*5)>=2 {
				//需要换牌
				user_poker.Pais[1] = PokerList[int(rand_index[len(rand_index)-1])]
				rand_index = rand_index[:len(rand_index)-1]
				user_poker.Pais[3] = PokerList[int(rand_index[len(rand_index)-1])]
				rand_index = rand_index[:len(rand_index)-1]
				user_poker = ParseNiuPoker(user_poker.Pais, desk.GetDeskOption())
			}

			u.Pokers = user_poker
		}
	}

	//比分均衡策略
	//desk.ScoreBalanceRule()
	desk.WhiteListSendPokerRule()

	return nil
}

//白名单策略
func (desk *Desk) WhiteListSendPokerRule() {
	var white_user *User = nil
	for _,u := range desk.Users {
		if u.GetIsOnWhiteList() {
			white_user = u
			break
		}
	}
	if white_user == nil {
		log.T("没有用户在白名单内。")
		return
	}

	rand_index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)

	log.T("白名单策略开始执行。白名单用户为%d", white_user.GetUserId())

	user_poker_rank := desk.GetUserRankByPoker()

	//如果总分数低于0，则通杀
	if white_user.Bill.GetScore() < 0 {
		log.T("用户%d(胜率%d)，总分数小于0，判胜。", white_user.GetUserId(), white_user.GetWhiteWinRate())
		rand_index = 0
		//poker_rank_1 := user_poker_rank[0].Pokers
		//for _,u := range desk.Users {
		//	if u.Pokers == poker_rank_1 {
		//		u.Pokers = white_user.Pokers
		//		break
		//	}
		//}
		//white_user.Pokers = poker_rank_1
		//return
	}

	//如果为抢庄模式或者通比牛牛，则给他最大的牌
	if desk.DeskOption.GetBankRule() == ddproto.NiuniuEnumBankerRule_QIANG_ZHUANG || desk.DeskOption.GetBankRule() == ddproto.NiuniuEnumBankerRule_TONG_BI_NIUNIU || white_user.IsBanker() {
		bank_mode := ""
		if desk.DeskOption.GetBankRule() == ddproto.NiuniuEnumBankerRule_QIANG_ZHUANG {
			bank_mode = "抢庄模式"
		}else if desk.DeskOption.GetBankRule() == ddproto.NiuniuEnumBankerRule_TONG_BI_NIUNIU {
			bank_mode = "通比牛牛"
		}else {
			bank_mode = "庄家"
		}
		//让他赢
		if rand_index < int(white_user.GetWhiteWinRate()) {
			log.T("用户%d为%s(胜率%d)，判胜。", white_user.GetUserId(), bank_mode, white_user.GetWhiteWinRate())
			//大于4人，发第二大牌
			poker_rank_2_index := 0
			poker_rank_2 := user_poker_rank[poker_rank_2_index].Pokers
			for _,u := range desk.Users {
				if u.Pokers == poker_rank_2 {
					u.Pokers = white_user.Pokers
					break
				}
			}
			white_user.Pokers = poker_rank_2
		}else {
			//让他输
			log.T("用户%d为%s(胜率%d)，判输。", white_user.GetUserId(), bank_mode, white_user.GetWhiteWinRate())
			//倒数第二大牌
			poker_rank_4 := user_poker_rank[(len(desk.Users)/2)+len(desk.Users)%2].Pokers
			for _,u := range desk.Users {
				if u.Pokers == poker_rank_4 {
					u.Pokers = white_user.Pokers
					break
				}
			}
			white_user.Pokers = poker_rank_4
		}
	}else {
		cuur_banker,_ := desk.GetUserByUid(desk.GetCurrBanker())
		if cuur_banker == nil {
			log.E("当前庄家为nil")
			return
		}
		banker_poker := cuur_banker.Pokers
		//让他赢
		if rand_index < int(white_user.GetWhiteWinRate()) {
			log.T("用户%d不为庄家(胜率%d)，判胜。", white_user.GetUserId(), white_user.GetWhiteWinRate())
			//如果不比庄家大
			if !IsBigThanBanker(cuur_banker.Pokers, white_user.Pokers){
				//跟庄家互换牌型
				cuur_banker.Pokers = white_user.Pokers
				white_user.Pokers = banker_poker
			}
		}else {
			log.T("用户%d不为庄家(胜率%d)，判输。", white_user.GetUserId(), white_user.GetWhiteWinRate())
			//让他输
			//如果比庄家大
			if IsBigThanBanker(cuur_banker.Pokers, white_user.Pokers){
				//跟庄家互换牌型
				cuur_banker.Pokers = white_user.Pokers
				white_user.Pokers = banker_poker
			}
		}
	}

	log.T("白名单策略执行结束。白名单用户为%d", white_user.GetUserId())
}

//比分均衡策略
func (desk *Desk) ScoreBalanceRule() {
	log.T("开始换牌")
	if desk.GetCircleNo() <= 1 {
		log.T("第一圈不换牌")
		return
	}

	user_bill_rank := desk.GetUserRankByBill()
	user_poker_rank := desk.GetUserRankByPoker()

	if user_bill_rank[0].Bill.GetScore() < 80 || user_bill_rank[len(user_bill_rank)-1].Bill.GetScore() > -80 {
		return
	}

	log.T("user_bill_rank", user_bill_rank)
	log.T("user_poker_rank", user_poker_rank)

	if len(user_bill_rank) == 0 || len(user_bill_rank) != len(user_poker_rank) {
		log.E("换牌出错！")
		return
	}

	//值拷贝
	poker_rank := make([]*ddproto.NiuniuSrvPoker,len(user_poker_rank))
	for i,p := range user_poker_rank {
		poker_rank[i] = p.Pokers
	}

	//重新分配牌型
	for i,_ := range user_bill_rank {
		y := len(user_bill_rank)-i-1

		user_bill_rank[i].Pokers = poker_rank[y]
	}
	log.T("换牌成功！")
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

//获取牌型大小排名
func (desk *Desk)GetUserRankByPoker() []*User {
	users := []*User{}

	for _,u := range desk.Users {
		if u.GetIsOnGamming() {
			users = append(users, u)
		}
	}

	for i:=0;i<len(users)-1;i++ {
		for j:=i+1;j<len(users);j++ {
			if IsBigThanBanker(users[i].Pokers, users[j].Pokers) {
				tmp := users[i]
				users[i] = users[j]
				users[j] = tmp
			}
		}
	}
	return users
}


//是否都准备了-除了房主
func (desk *Desk) IsAllReadyExcludeOwner() error {
	var i int32 = 0
	if desk.GetStatus() > ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START {
		return errors.New("该房间正在游戏中！")
	}
	for _,u := range desk.Users {
		//如果不是房主
		if u != nil && !u.IsOwner() && u.GetIndex() >= 0 {
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
	if desk.GetStatus() > ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START {
		return errors.New("该房间正在游戏中！")
	}
	for _,u := range desk.Users {
		//如果不是房主
		if u != nil && u.GetIndex() >= 0 {
			if u.GetIsReady() == true {
				i++
			}else {
				return errors.New(fmt.Sprintf("用户%d未准备！", u.GetUserId()))
			}
		}
	}

	if i >= desk.DeskOption.GetMinUser() {
		return nil
	}

	return errors.New(fmt.Sprintf("未达到%d人最小开局条件！", desk.DeskOption.GetMinUser()))
}

//是否达到自动开局条件
func (desk *Desk) IsAllReadyAutoStart() error {
	var i int32 = 0
	if desk.GetStatus() > ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START {
		return errors.New("该房间正在游戏中！")
	}
	for _,u := range desk.Users {
		//如果不是房主
		if u != nil && !u.IsOwner() && u.GetIndex() >= 0 {
			if u.GetIsReady() == true {
				i++
			}else {
				return errors.New(fmt.Sprintf("用户%d未准备！", u.GetUserId()))
			}
		}
	}
	if i >= desk.DeskOption.GetAutoStartGammer() - 1 {
		return nil
	}
	return errors.New(fmt.Sprintf("未达到%d人最小开局条件！", desk.DeskOption.GetMinUser()))
}

//是否都亮牌了
func (desk *Desk) IsAllLiangpai() error {
	for _,u := range desk.Users {
		//如果在游戏中
		if u != nil && u.GetIsOnGamming() {
			if !u.GetIsLiangpai() {
				return errors.New(fmt.Sprintf("用户%d未亮牌！", u.GetUserId()))
			}
		}
	}
	return nil
}

//是否所有人都已经抢庄
func (desk *Desk) IsAllQiangzhuang() error {
	if desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_QIANGZHUANG {
		return errors.New("调用IsAllQiangzhuang非法！")
	}
	var i int32 = 0
	for _,u := range desk.Users {
		//如果不是房主
		if u != nil && u.GetIsOnGamming() {
			if u.GetBankerScore() == 0 {
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
//是否所有人都已经加倍
func (desk *Desk) IsAllJiabeiExcludeBanker() error {
	if desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_JIABEI {
		return errors.New("调用IsAllJiaBei非法！")
	}
	for _,u := range desk.Users {
		//如果不是庄家
		if u != nil && !u.IsBanker() && u.GetIsOnGamming() {
			if u.GetDoubleScore() == 0 {
				return errors.New(fmt.Sprintf("用户%d未加倍！", u.GetUserId()))
			}
		}
	}

	return nil
}


//客户端牌桌
func (desk *Desk) GetClientDesk() *ddproto.NiuniuClientDesk {
	users := []*ddproto.NiuniuClientUser{}

	for _,u := range desk.Users {
		if u == nil {
			continue
		}
		users =append(users, u.GetClientUser())
	}

	client_desk := &ddproto.NiuniuClientDesk{
		DeskId: proto.Int32(desk.GetDeskId()),
		DeskNumber: proto.String(desk.GetDeskNumber()),
		GameNumber: proto.Int32(desk.GetGameNumber()),
		RoomId: proto.Int32(desk.GetRoomId()),
		Status: desk.Status,
		LastWiner: proto.Uint32(desk.GetLastWiner()),
		CircleNo: proto.Int32(desk.GetCircleNo()),
		Owner: proto.Uint32(desk.GetOwner()),
		CurrBanker: proto.Uint32(desk.GetCurrBanker()),
		DeskOption: desk.DeskOption,
		Users: users,
		IsStart: desk.IsStart,
		IsOnDissolve: desk.IsOnDissolve,
		DissolveTime: proto.Int32(0),
		IsDaikai: proto.Bool(desk.GetIsDaikai()),
		DaikaiUser: proto.Uint32(desk.GetDaikaiUser()),
		IsCoinRoom: proto.Bool(desk.GetIsCoinRoom()),
		IsOnGamming: proto.Bool(desk.GetIsOnGamming()),
		SurplusTime: proto.Int32(int32(desk.GetSurplusTime())),
	}

	//投票剩余时间
	if desk.GetIsOnDissolve() && desk.GetDissolveTime() > 0 {
		last_time := time.Now().Unix() - desk.GetDissolveTime()
		client_desk.DissolveTime = proto.Int32(int32(last_time))
	}

	//屏蔽等待进入的状态
	if desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_ENTER {
		client_desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY.Enum()

		//TODO: kory临时新增下面一行（金币场自建房）
		desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY.Enum()
	}

	return client_desk
}

//房费计算
func GetOwnerFee(circle int32) int32 {
	switch  {
	case circle > 0 && circle <= 10 :
		return 1
	case circle > 10 && circle <= 20:
		return 2
	case circle > 20:
		return 3
	default:
		return 0
	}
}

//插入比赛记录表
func (desk *Desk)InsertOneCounter() error {
	users_str := ""
	for _, u := range desk.Users {
		users_str += fmt.Sprintf(",%d", u.GetUserId())
	}

	oneCounter := model.T_niuniu_desk_round{
		BeginTime: time.Unix(desk.GetOneStartTime(), 0),
		DeskId: desk.GetDeskId(),
		UserIds: users_str,
		GameNumber: desk.GetGameNumber(),
		Password: desk.GetDeskNumber(),
		EndTime: time.Now(),
		IsCoinRoom: desk.GetIsCoinRoom(),
		BankRule: desk.DeskOption.GetBankRule(),
	}

	for _, u := range desk.Users {
		//1局记录
		oneCounter.Records = append(oneCounter.Records, model.NiuRecordBean{
			UserId:    u.GetUserId(),
			NickName:  u.GetNickName(),
			WinAmount: int64(u.GetLastScore()),
		})

		//game_bill
		userGameBillService.Insert(userGameBillService.T_user_game_bill{
			UserId: u.GetUserId(),
			GameId: int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN),
			RoomType: int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND),
			DeskId: desk.GetDeskId(),
			Password: desk.GetDeskNumber(),
			GameNumber: desk.GetGameNumber(),
			WinAmount: int64(u.GetLastScore()),
			Balance: int64(u.Bill.GetScore()),
			IsRobot: u.GetIsRobot(),
			EndTime: time.Now(),
		})
	}

	if !desk.GetIsCoinRoom(){
		//插入朋友桌记录
		db.Log(tableName.DBT_NIU_DESK_ROUND_ONE).Insert(oneCounter)
	}else {
		//插入金币场记录
		db.Log(tableName.DBT_NIU_DESK_ROUND_ONE_COIN).Insert(oneCounter)
	}

	return nil
}

//插入10局比赛记录
func (desk *Desk)InsertAllCounter() error {
	users_str := ""
	for _, u := range desk.Users {
		users_str += fmt.Sprintf(",%d", u.GetUserId())
	}

	allCounter := model.T_niuniu_desk_round{
		BeginTime: time.Unix(desk.GetAllStartTime(), 0),
		DeskId: desk.GetDeskId(),
		UserIds: users_str,
		TotalRound: desk.DeskOption.GetMaxCircle(),
		GameNumber: desk.GetGameNumber(),
		Password: desk.GetDeskNumber(),
		EndTime: time.Now(),
	}

	for _, u := range desk.Users {
		//1局记录
		allCounter.Records = append(allCounter.Records, model.NiuRecordBean{
			UserId:    u.GetUserId(),
			NickName:  u.GetNickName(),
			WinAmount: int64(u.Bill.GetScore()),
		})
	}

	return db.Log(tableName.DBT_NIU_DESK_ROUND_ALL).Insert(allCounter)
}

//玩法
func (desk *Desk) GetPlayRuleStr() string {
	rule_type := "牛牛换庄"
	switch desk.DeskOption.GetBankRule() {
	case ddproto.NiuniuEnumBankerRule_DING_ZHUANG:
		rule_type = "牛牛换庄"
	case ddproto.NiuniuEnumBankerRule_FANGZHU_DINGZHUANG:
		rule_type = "牛牛定庄"
	case ddproto.NiuniuEnumBankerRule_QIANG_ZHUANG:
		rule_type = "明牌抢庄"
	case ddproto.NiuniuEnumBankerRule_SUI_JI_ZUO_ZHUANG:
		rule_type = "随机坐庄"
	case ddproto.NiuniuEnumBankerRule_TONG_BI_NIUNIU:
		rule_type = "通比牛牛"
	}
	return rule_type
}

//获取房间简介
func (desk *Desk) GetTips() string {
	room_type := "急速牛牛"
	if desk.DeskOption.GetHasAnimation() {
		room_type = "经典牛牛"
	}
	has_flower := "无花牌"
	if desk.DeskOption.GetHasFlower() {
		has_flower = "有花牌"
	}
	rule_type := desk.GetPlayRuleStr()

	deny_half_join := "允许中途加入"
	if desk.DeskOption.GetDenyHalfJoin() {
		deny_half_join = "禁止中途加入"
	}

	is_flower_play := "普通玩法"
	if desk.DeskOption.GetIsFlowerPlay() {
		is_flower_play = "花式玩法"
	}

	tips := fmt.Sprintf("%s %s %d人%d局 %s %s %s", room_type, rule_type, desk.DeskOption.GetMaxUser(), desk.DeskOption.GetMaxCircle(), has_flower, is_flower_play, deny_half_join)

	return tips
}

//获取空余的位置索引
func (desk *Desk)GetFreeSiteIndex() (int, error)  {
	//座位号从1开始计数
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

//获取已入座用户数
func (desk *Desk) GetSitedUserNum() (num int32) {
	for _,u := range desk.Users {
		if u == nil || u.GetIndex() == -1 {
			continue
		}
		num++
	}
	return
}

//获取剩余秒数
func (desk *Desk) GetSurplusTime() int64 {
	now := time.Now().Unix()
	var surplus_time int64 = 0
	switch desk.GetStatus() {
	case ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_ENTER, ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY, ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START:
		if desk.GetCircleNo() == 1 {
			if desk.GetIsStart() {
				surplus_time = 7 - (now - desk.GetStartTime())
			}else {
				surplus_time = 15
			}
		}else {
			surplus_time = 7 - (now - desk.GetStartTime())
		}
	case ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_QIANGZHUANG:
		surplus_time = 7 - (now - desk.GetStartTime())
	case ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_JIABEI:
		surplus_time = 7 - (now - desk.GetStartTime())
	case ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_LIANGPAI:
		surplus_time = 9 - (now - desk.GetStartTime())
	}

	if surplus_time <= 0 {
		return 1
	}

	return surplus_time
}
