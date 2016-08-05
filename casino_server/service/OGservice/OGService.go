package OGservice

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/service/room"
	"casino_server/common/log"
	"casino_server/gamedata"
	"casino_server/service/userService"
	"errors"
)

//联众德州,桌子状态
var (
	GAME_STATUS_READY int32 = 0        //准备
	GAME_STATUS_DEAL_CARDS int32 = 1        //发牌
	GAME_STATUS_PRECHIP int32 = 2        //盲注
	GAME_STATUS_FIRST_TURN int32 = 3        //第一轮
	GAME_STATUS_SECOND_TURN int32 = 4        //第二轮
	GAME_STATUS_THIRD_TURN int32 = 5        //第三轮
	GAME_STATUS_LAST_TURN int32 = 6        //第四轮
	GAME_STATUS_SHOW_RESULT int32 = 7        //完成
)


/**
	用户通过钻石创建游戏房间
 */

func HandlerCreateDesk(userId uint32,diamond int64,roomKey string) error{

	//1,判断roomKey是否已经存在
	if room.ThGameRoomIns.IsRoomKeyExist(roomKey) {
		log.E("房间密钥[%v]已经存在,创建房间失败",roomKey)
		return errors.New("房间密钥已经存在,创建房间失败")
	}

	//2,开始创建房间
	desk := room.ThGameRoomIns.CreateDeskByUserIdAndRoomKey(userId,diamond,roomKey);
	log.T("",desk)
	//3,根据返回的desk返回创建房间的信息
	return nil

}


//处理登录游戏的协议
/**
	1,判断用户是否已经登陆了游戏
	2,如果已经登陆了游戏,替换现有的agent
	3,如果没有登陆游戏,走正常的流程
 */
func HandlerGameEnterMatch(m *bbproto.Game_EnterMatch, a gate.Agent) error {
	log.T("用户请求进入德州扑克的游戏房间,m[%v]",m)

	userId := m.GetUserId()				//进入游戏房间的user
	result := newGame_SendGameInfo()                //需要返回的信息
	roomCoin := int64(1000)				//to do 暂时设置为1000
	//1.1 检测参数是否正确,判断userId 是否合法
	userCheck := userService.CheckUserIdRightful(userId)
	if userCheck == false {
		log.E("进入德州扑克的房间的时候,userId[%v]不合法。", userId)
		return errors.New("用户Id不合法")
	}

	//1,进入房间,返回房间和错误信息
	mydesk, err := room.ThGameRoomIns.AddUser(userId,roomCoin,"", a)
	if err != nil || mydesk == nil {
		errMsg := err.Error()
		log.E("用户[%v]进入房间失败,errMsg[%v]",userId,errMsg)
		a.WriteMsg(result)
		return err
	}

	//2 构造信息并且返回
	initGameSendgameInfoByDesk(mydesk, result,userId)
	log.T("给请求登陆房间的人[%v]回复信息[%v]",userId,result)

	//3 发送进入游戏房间的广播
	mydesk.OGTHBroadAddUser(result)

	//4,最后:确定是否开始游戏, 上了牌桌之后,如果玩家人数大于1,并且游戏处于stop的状态,则直接开始游戏

	//如果是朋友桌的话,需要房主点击开始才能开始...
	go mydesk.OGRun()

	return nil
}


//初始化一个Game_SendGameInfo
func newGame_SendGameInfo() *bbproto.Game_SendGameInfo {
	result := &bbproto.Game_SendGameInfo{}
	result.TablePlayer = new(int32)
	result.Tableid = new(int32)
	result.BankSeat = new(int32)
	result.ChipSeat = new(int32)
	result.ActionTime = new(int32)
	result.DelayTime = new(int32)
	result.GameStatus = new(int32)
	//result.NRebuyCount
	//result.NAddonCount
	result.Pool = new(int64)
	result.MinRaise = new(int64)
	result.NInitActionTime = new(int32)
	result.NInitDelayTime = new(int32)
	result.Handcard = make([]*bbproto.Game_CardInfo, 0)
	result.TurnCoin = make([]int64, 0)
	result.BFold = make([]int32, 0)
	result.BAllIn = make([]int32, 0)
	result.BBreak = make([]int32, 0)
	result.BLeave = make([]int32, 0)
	result.Seat = new(int32)

	return result
}

//返回房间的信息 todo 登陆成功的处理,给请求登陆的玩家返回登陆结果的消息
func initGameSendgameInfoByDesk(mydesk *room.ThDesk, result *bbproto.Game_SendGameInfo,myUserId uint32) error {
	//初始化桌子相关的信息
	*result.Tableid = int32(mydesk.Id)        //桌子的Id
	*result.TablePlayer = mydesk.UserCount	  //玩家总人数
	*result.BankSeat =  mydesk.GetUserByUserId(mydesk.Dealer).Seat //int32(mydesk.GetUserIndex(mydesk.Dealer))        //庄家
	*result.ChipSeat =  mydesk.GetUserByUserId(mydesk.BetUserNow).Seat //int32(mydesk.GetUserIndex(mydesk.BetUserNow))//当前活动玩家
	*result.ActionTime = int32(room.TH_TIMEOUT_DURATION_INT)        //当前操作时间,服务器当前的时间
	*result.DelayTime = int32(1000)        //当前延时时间
	*result.GameStatus = deskStatus2OG(mydesk)
	*result.Pool = int64(mydesk.Jackpot)                //奖池
	result.Publiccard = mydesk.ThPublicCard2OGC()        //公共牌...
	*result.MinRaise = int64(mydesk.MinRaise)        //最低加注金额
	*result.NInitActionTime = int32(room.TH_TIMEOUT_DURATION_INT)
	*result.NInitDelayTime = int32(room.TH_TIMEOUT_DURATION_INT)
	result.Handcard = mydesk.GetHandCard()		//用户手牌
	//result.HandCoin = mydesk.GetCoin()	//下注的金额
	result.HandCoin = mydesk.GetRoomCoin()	//带入金额
	result.TurnCoin = getTurnCoin(mydesk)
	*result.Seat	= mydesk.GetUserByUserId(myUserId).Seat	//int32(mydesk.GetUserIndex(myUserId))	//我

	//循环User来处理
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户当前状态
			result.BAllIn = append(result.BAllIn, isAllIn(u))                        //是否已经全下了
			result.BBreak = append(result.BBreak, isBreak(u))                        //是否
			result.BLeave = append(result.BLeave, isLeave(u))
			result.BEnable = append(result.BEnable, isEnable(u))                        //用户是否可以操作,0表示不能操作,1表示可以操作
			result.BFold = append(result.BFold, isFold(u))        //是否弃牌
			//nickName			//seatId
			result.NickName = append(result.NickName, u.NickName)
			result.SeatId = append(result.SeatId, int32(i))

		} else {

		}
	}

	return nil

}

//判断是否allIn
func isAllIn(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_ALLINING {
		return 1
	} else {
		return 0
	}
}


//判断是否allIn
func isBreak(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_BREAK {
		return 1
	} else {
		return 0
	}
}

//判断是否allIn
func isLeave(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_LEAVE {
		return 1
	} else {
		return 0
	}
}

//判断是否allIn
func isFinalAddon(u *room.ThUser) int32 {
	return 0
}

//判断是否allIn
func isEnable(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_BETING {
		//游戏中
		return 1
	} else {
		return 0
	}
}

//判断是否allIn
func isFold(u *room.ThUser) int32 {
	if u.Status == room.TH_USER_STATUS_FOLDED {
		//游戏中
		return 1
	} else {
		return 0
	}
}

//解析TurnCoin
func getTurnCoin(mydesk *room.ThDesk) []int64{
	var result []int64
	for i := 0; i < len(mydesk.Users); i++ {
		u := mydesk.Users[i]
		if u != nil {
			//用户手牌
			result = append(result,int64(u.TurnCoin))
		}
	}
	return result
}





//游戏状态的转换
func deskStatus2OG(desk *room.ThDesk) int32 {
	var result int32 = GAME_STATUS_READY
	status := desk.Status
	round := desk.RoundCount

	if status == room.TH_DESK_STATUS_STOP {
		//没有开始
		result = GAME_STATUS_READY
	} else if status == room.TH_DESK_STATUS_SART {
		switch round {
		case room.TH_DESK_ROUND1:
			result = GAME_STATUS_FIRST_TURN
		case room.TH_DESK_ROUND2:
			result = GAME_STATUS_SECOND_TURN
		case room.TH_DESK_ROUND3:
			result = GAME_STATUS_THIRD_TURN
		case room.TH_DESK_ROUND4:
			result = GAME_STATUS_LAST_TURN
		default:
			result = GAME_STATUS_DEAL_CARDS
		}
	} else if status == room.TH_DESK_STATUS_LOTTERY {
		result = GAME_STATUS_SHOW_RESULT
	}

	return result
}





//通过agent返回UserId
func getUserIdByAgent( a gate.Agent) uint32{
	//获取agent中的userData
	ad := a.UserData()
	if ad == nil {
		log.E("agent中的userData为nil")
		return uint32(0)
	}

	userData := ad.(*gamedata.AgentUserData)
	log.T("得到的UserAgent中的userId是[%v]",userData.UserId)
	//return userData.UserId

	//测试代码,返回10006
	return 10006
}

