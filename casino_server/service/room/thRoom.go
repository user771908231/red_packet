package room

import (
	"casino_server/msg/bbprotogo"
	"sync"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"errors"
	"casino_server/service/userService"
	"casino_server/mode"
	"casino_server/utils/numUtils"
	"casino_server/utils"
	"casino_server/conf/intCons"
	"casino_server/common/Error"
)

var ThGameRoomIns ThGameRoom                //房间实例,在init函数中初始化

func init() {
	ThGameRoomIns.OnInit()                //初始化房间
}

/**
	德州扑克
 */

//游戏房间
type ThGameRoom struct {
	sync.Mutex
	RoomStatus      int32 //游戏大厅的状态
	ThDeskBuf       []*ThDesk
	ThRoomSeatMax   int32 //每个房间的座位数目
	ThRoomCount     int32 //房间数目
	Id              int32 //房间的id
	SmallBlindCoin  int64 //小盲注的金额
	RebuyCountLimit int32 //重购的次数限制
}


//初始化游戏房间
func (r *ThGameRoom) OnInit() {
	log.T("初始化thgameroom.oninit()")
	r.ThRoomSeatMax = ThdeskConfig.TH_DESK_MAX_START_USER
	r.Id = 0
	r.SmallBlindCoin = ThdeskConfig.TH_GAME_SMALL_BLIND;
	r.RebuyCountLimit = 10000; //朋友桌可以一直重购
}


//随机生成6位数字
func (r *ThGameRoom) RandRoomKey() string {
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
}


//判断roomkey是否已经存在了
func (r *ThGameRoom) IsRoomKeyExist(roomkey string) bool {
	ret := false
	for i := 0; i < len(r.ThDeskBuf); i++ {
		d := r.ThDeskBuf[i]
		if d != nil && d.RoomKey == roomkey {
			ret = true
			break
		}
	}
	return ret
}


//
func (r *ThGameRoom) CalcCreateFee(jucount int32) int64 {
	//通过局数来计算消耗的钻石数量
	return (int64(jucount)) / (int64(ThdeskConfig.CreateJuCountUnit)) * ThdeskConfig.CreateFee
}

//创建一个房间
func (r *ThGameRoom) CreateDeskByUserIdAndRoomKey(userId uint32, roomCoin int64, roomkey string, preCoin int64, smallBlind int64, bigBlind int64, jucount int32) (*ThDesk, error) {

	//1,创建房间成功之后,扣除user的钻石
	upDianmond := r.CalcCreateFee(jucount)
	remainDiamond, err := userService.UpdateUserDiamond(userId, -upDianmond)
	if err != nil {
		log.E("创建房间的时候出错,error", err.Error())
		return nil, err
	}

	//2,创建房间
	desk := NewThDesk()
	desk.RoomKey = roomkey
	desk.CreateFee = upDianmond
	desk.InitRoomCoin = roomCoin
	desk.DeskOwner = userId
	desk.SmallBlindCoin = smallBlind
	desk.BigBlindCoin = bigBlind
	desk.JuCount = jucount
	desk.GetRoomCoin()
	desk.GameType = intCons.GAME_TYPE_TH        //表示是自定义的房间
	desk.PreCoin = preCoin
	r.AddThDesk(desk)

	//3,生成一条交易记录
	err = userService.CreateDiamonDetail(userId, mode.T_USER_DIAMOND_DETAILS_TYPE_CREATEDESK, upDianmond, remainDiamond, "创建房间消耗钻石");
	if err != nil {
		log.E("创建用户的钻石交易记录失败")
		return nil, err
	}
	return desk, nil

}

//增加一个thRoom
func (r *ThGameRoom) AddThDesk(throom *ThDesk) error {
	r.ThDeskBuf = append(r.ThDeskBuf, throom)
	return nil
}


//删除一个throom
func (r *ThGameRoom) RmThroom(desk *ThDesk) error {
	//第一步找到index
	var index int = -1
	for i := 0; i < len(r.ThDeskBuf); i++ {
		desk = r.ThDeskBuf[i]
		if desk != nil && desk.Id == desk.Id {
			index = i
			break
		}
	}

	//修改用户session状态,因为seession中存在desk相关的信息
	for i := 0; i < len(desk.Users); i++ {
		u := desk.Users[i]
		if u != nil {
			agent := u.agent
			//更新agentData的数据
			u.deskId = 0
			u.MatchId = 0
			u.UpdateAgentUserData(agent)
		}
	}


	//删除对应的desk
	r.ThDeskBuf = append(r.ThDeskBuf[:index], r.ThDeskBuf[index + 1:]...)
	return nil
}


//通过房主id解散房间
func (r *ThGameRoom) DissolveDeskByDeskOwner(userId uint32, a gate.Agent) error {

	result := &bbproto.Game_AckDissolveDesk{}
	result.Result = new(int32)
	result.UserId = new(uint32)
	result.DeskId = new(int32)
	result.PassWord = new(string)

	//1,找到桌子
	desk := GetDeskByAgent(a)

	if desk == nil {
		return errors.New("房间已经解散了")
	}

	if desk.DeskOwner != userId {
		return errors.New("不是房主,没有权限解散房间")
	}

	//2,解散桌子的条件,如果正在游戏中,是否能解散?
	if desk.Status != TH_DESK_STATUS_STOP {
		*result.Result = intCons.ACK_RESULT_ERROR
		a.WriteMsg(result)
		return errors.New("游戏正在进行中,不能解散")
	}

	//3,解散
	RmThdesk(desk)

	//4,发送解散的广播
	*result.Result = intCons.ACK_RESULT_SUCC
	*result.UserId = desk.DeskOwner
	*result.PassWord = desk.RoomKey
	desk.THBroadcastProtoAll(result)

	return nil
}



//通过UserId判断是不是重复进入房间
func (r *ThGameRoom) IsRepeatIntoRoom(userId uint32, a gate.Agent) *ThDesk {
	//1,取回话信息,如果回话信息为nil,直接返回nil
	userData := userService.GetUserSessionByUserId(userId)
	if userData == nil {
		return nil
	}

	//2,取桌子的信息,如果桌子为nil,则直接返回nil
	desk := GetDeskByIdAndMatchId(userData.GetDeskId(), userData.GetMatchId())
	if desk == nil {
		return nil
	}

	//3,重新设置用户的信息
	desk.GetUserByUserId(userId).UpdateAgentUserData(a)
	desk.AddUserCountOnline()

	log.T("用户[%v]重新进入房间了", userId)
	return desk
}

/**
	通过roomKey 找到desk
 */
func (r *ThGameRoom) GetDeskByRoomKey(roomKey string) *ThDesk {
	var result *ThDesk
	desks := ThGameRoomIns.ThDeskBuf
	for i := 0; i < len(desks); i++ {
		desk := desks[i]
		if desk != nil && desk.RoomKey == roomKey {
			result = desk
			break
		}
	}
	return result
}


/**
	给指定的房间增加用户
 */
func (r *ThGameRoom) AddUserWithRoomKey(userId uint32, roomKey string, a gate.Agent) (*ThDesk, error) {
	log.T("玩家[%v]通过roomkey[%v]进入房间", userId, roomKey)
	//1,首先判断roomKey 是否喂空
	if roomKey == "" {
		return nil, errors.New("房间密码不应该为空")
	}

	//2,如果roomKey 不是为""
	mydesk := r.GetDeskByRoomKey(roomKey)
	if mydesk == nil {
		return nil, Error.NewError(int32(bbproto.DDErrorCode_ERRORCODE_INTO_DESK_NOTFOUND), "没有找到对应的房间")
	}

	//3,判断用户是否是掉线重连
	isRepeat := mydesk.IsrepeatIntoWithRoomKey(userId, a)
	if isRepeat {
		return mydesk, nil
	}

	//4,进入房间
	_, err := mydesk.AddThUser(userId, TH_USER_STATUS_SEATED, a)
	if err != nil {
		log.E("用户上德州扑克的桌子 失败...")
		return nil, err
	}

	mydesk.LogString()        //答应当前房间的信息

	return mydesk, nil

}

//得到锦标赛房间
func GetCSTHroom(matchId int32) *CSThGameRoom {
	if matchId < 0 {
		return nil
	} else {
		return &ChampionshipRoom
	}
}



//通过Id找到对应的桌子
func (r *ThGameRoom) GetDeskById(id int32) *ThDesk {
	var result *ThDesk = nil
	for i := 0; i < len(r.ThDeskBuf); i++ {
		if r.ThDeskBuf[i] != nil && r.ThDeskBuf[i].Id == id {
			result = r.ThDeskBuf[i]
			break
		}
	}
	return result
}


//通过肘子的类型和Match得到thdesk
func GetDeskByIdAndMatchId(deskId int32, matchId int32) *ThDesk {
	//1,把type 转义
	if matchId > 0 {
		//返回锦标赛的房间
		return ChampionshipRoom.GetDeskById(deskId)
	} else if matchId == 0 {
		//返回自定义房间里面的desk
		return ThGameRoomIns.GetDeskById(deskId)

	} else {
		return nil
	}
}

//通过连接得到桌子
func GetDeskByAgent(a gate.Agent) *ThDesk {
	//得到用户数据
	var userData *bbproto.ThServerUserSession
	agentData := a.UserData()
	if agentData == nil {
		return nil
	}

	userData = agentData.(*bbproto.ThServerUserSession)

	//得到桌子
	deskId := userData.GetDeskId()
	matchId := userData.GetMatchId()
	//gameStatus := userData.GetGameStatus()

	//返回数据
	desk := GetDeskByIdAndMatchId(deskId, matchId)
	log.T("通过agent.userData()[%v]得到thdesk[%v]", userData, desk)

	return desk
}

//删除房间
func RmThdesk(desk *ThDesk) error {
	log.T("开始解散房间id[%v],desk[%v]", desk.Id, desk)
	if desk.GameType == intCons.GAME_TYPE_TH_CS {
		//锦标赛
		ChampionshipRoom.RmThroom(desk)
	} else if desk.GameType == intCons.GAME_TYPE_TH {
		ThGameRoomIns.RmThroom(desk)
	}
	return nil
}

