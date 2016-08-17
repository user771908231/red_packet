package OGservice

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/service/room"
	"casino_server/common/log"
	"casino_server/gamedata"
	"casino_server/service/userService"
	"errors"
	"github.com/name5566/leaf/db/mongodb"
	"casino_server/conf/casinoConf"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"casino_server/mode"
	"casino_server/conf/intCons"
	"casino_server/utils/timeUtils"
	"casino_server/utils/numUtils"
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
	ateDesk(userId,initCoin,roomKey,smallBlind,bigBlind,juCount)

	//新修改:创建房间的时候roomKey 系统指定,不需要用户输入
 */

func HandlerCreateDesk(userId uint32,roomCoin int64,smallBlind int64,bigBlind int64,jucount int32) (*room.ThDesk,error){

	//1,得到一个随机的密钥
	roomKey := room.ThGameRoomIns.RandRoomKey()

	//2,开始创建房间
	desk,err := room.ThGameRoomIns.CreateDeskByUserIdAndRoomKey(userId,roomCoin,roomKey,smallBlind ,bigBlind ,jucount );
	if err != nil {
		log.E("用户创建房间失败errMsg[%v]",err.Error())
		return nil,err
	}

	//new : 创建房间之后返回创建房间的结果,不用知己进入房间
	return desk,nil

}

//离开房间
func HandlerLeaveDesk(m *bbproto.Game_Message,a gate.Agent){
	//deskId
	desk := room.ThGameRoomIns.GetDeskByUserId(m.GetUserId())
	desk.LeaveThuser(m.GetUserId())
}

///用户发送消息
func HandlerMessage(m *bbproto.Game_Message,a gate.Agent){
	//deskId
	result := &bbproto.Game_SendMessage{}
	desk := room.ThGameRoomIns.GetDeskByUserId(m.GetUserId())
	if desk == nil {
		//返回错误信息
	}else{
		result.UserId = m.UserId
		result.Msg = m.Msg
		result.MsgType = m.MsgType
		result.Id = m.Id
		desk.THBroadcastProtoAll(result)
	}
}

//用户准备
func HandlerReady(m *bbproto.Game_Ready,a gate.Agent) error{

	//1,找到userId
	userId := m.GetUserId()

	//2,通过userId 找到桌子
	desk := room.ThGameRoomIns.GetDeskByUserId(userId)

	//3,用户开始准备
	desk.Ready(userId)

	//4,返回准备的结果
	return nil
}

//开始游戏
func HandlerBegin(m *bbproto.Game_Begin,a gate.Agent) error{
	userId := m.GetUserId()
	desk := room.ThGameRoomIns.GetDeskByDeskOwner(userId)
	if desk == nil {
		log.E("没有找到房主为[%v]的desk",userId)
		return errors.New("没有找到房间")
	}else{
		//开始游戏
		desk.OGRun()
		return nil
	}
}

//
func HandlerGetGameRecords(m *bbproto.Game_GameRecord,a gate.Agent){
	userId := m.GetUserId()

	//1,获取数据库连接
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()

	s := c.Ref()
	defer c.UnRef(s)

	//v2 战绩查询
	var deskRecords []mode.T_th_desk_record
	queryKey,_ := numUtils.Uint2String(userId)
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TH_DESK_RECORD).Find(bson.M{"userids": bson.RegEx{queryKey,"."}}).Limit(20).All(&deskRecords)
	log.T("查询到用户[%v]的战绩[%v]",userId,deskRecords)

	//需要返回的数据
	resDatav2 := &bbproto.Game_AckGameRecord{}
	resDatav2.UserId = new(uint32)
	resDatav2.Result = new(int32)

	*resDatav2.UserId = userId
	*resDatav2.Result = intCons.ACK_RESULT_SUCC

	//第一层循环战绩
	for i := 0; i < len(deskRecords); i++ {
		tr := deskRecords[i]
		r := &bbproto.Game_BeanGameRecord{}
		r.DeskId = new(int32)
		r.Id = new(int32)
		r.BeginTime = new(string)

		*r.DeskId 	= tr.DeskId
		*r.BeginTime 	= timeUtils.Format(tr.BeginTime)

		//第二层循环人
		for j := 0; j < len(tr.Records); j++ {
			tru := tr.Records[j]

			ru := &bbproto.Game_BeanUserRecord{}
			ru.UserId = new(uint32)
			ru.NickName = new(string)
			ru.WinAmount = new(int64)

			*ru.UserId = tru.UserId
			*ru.NickName = tru.NickName
			*ru.WinAmount = tru.WinAmount
			r.Users = append(r.Users,ru)
		}
		resDatav2.Records =  append(resDatav2.Records,r)		//加入要返回的结果
	}
	//返回查询到的记录
	a.WriteMsg(resDatav2)
}




//处理登录游戏的协议
/**
	1,判断用户是否已经登陆了游戏
	2,如果已经登陆了游戏,替换现有的agent
	3,如果没有登陆游戏,走正常的流程

	//错误码的说明:result



 */
func HandlerGameEnterMatch(m *bbproto.Game_EnterMatch, a gate.Agent) error {
	log.T("用户请求进入德州扑克的游戏房间,m[%v]",m)

	var err error					//错误信息
	var mydesk *room.ThDesk				//用户需要进入的房间
	userId := m.GetUserId()				//进入游戏房间的user
	roomCoin := int64(1000)				//to do 暂时设置为1000
	roomKey := string(m.GetPassWord())		//房间的roomkey
	result := newGame_SendGameInfo()                //需要返回的信息


	//1.1 检测参数是否正确,判断userId 是否合法
	userCheck := userService.CheckUserIdRightful(userId)
	if userCheck == false {
		log.E("进入德州扑克的房间的时候,userId[%v]不合法。", userId)
		return errors.New("用户Id不合法")
	}

	log.T("用户请求进入德州扑克的游戏房间,password[%v]",m.GetPassWord())
	//1,进入房间,返回房间和错误信息
	if roomKey == "" {
		mydesk, err = room.ThGameRoomIns.AddUser(userId,roomCoin, a)
	}else {
		mydesk, err = room.ThGameRoomIns.AddUserWithRoomKey(userId,roomCoin,roomKey, a)
	}

	if err != nil || mydesk == nil {
		errMsg := err.Error()
		log.E("用户[%v]进入房间失败,errMsg[%v]",userId,errMsg)
		*result.Result = intCons.ACK_RESULT_ERROR
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
	if mydesk.DeskType == room.TH_DESK_TYPE_JINBIAOSAI {
		go mydesk.OGRun()
	}

	return nil
}

//处理登录游戏的协议
/**
	1,判断用户是否已经登陆了游戏
	2,如果已经登陆了游戏,替换现有的agent
	3,如果没有登陆游戏,走正常的流程
 */
func HandlerGameEnterMatchWithRoomKey(m *bbproto.Game_EnterMatch, a gate.Agent) error {
	log.T("用户请求进入德州扑克的游戏房间,m[%v]",m)

	log.T("用户请求进入德州扑克的游戏房间,m[%v]",m)

	var err error					//错误信息
	var mydesk *room.ThDesk				//用户需要进入的房间
	userId := m.GetUserId()				//进入游戏房间的user
	roomCoin := int64(1000)				//to do 暂时设置为1000
	roomKey := string(m.GetPassWord())		//房间的roomkey
	result := newGame_SendGameInfo()                //需要返回的信息


	//1.1 检测参数是否正确,判断userId 是否合法
	userCheck := userService.CheckUserIdRightful(userId)
	if userCheck == false {
		log.E("进入德州扑克的房间的时候,userId[%v]不合法。", userId)
		return errors.New("用户Id不合法")
	}


	//1,进入房间,返回房间和错误信息
	mydesk, err = room.ThGameRoomIns.AddUserWithRoomKey(userId,roomCoin,roomKey, a)
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
	//这里需要特别处理


	//go mydesk.OGRun()


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
	result.TurnMax = new(int64)
	result.Result = new(int32)
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
	*result.MinRaise = mydesk.GetMinRaise()        //最低加注金额
	*result.NInitActionTime = int32(room.TH_TIMEOUT_DURATION_INT)
	*result.NInitDelayTime = int32(room.TH_TIMEOUT_DURATION_INT)
	result.Handcard = mydesk.GetHandCard()		//用户手牌
	//result.HandCoin = mydesk.GetCoin()	//下注的金额
	result.HandCoin = mydesk.GetRoomCoin()	//带入金额
	result.TurnCoin = getTurnCoin(mydesk)
	*result.Seat	= mydesk.GetUserByUserId(myUserId).Seat	//int32(mydesk.GetUserIndex(myUserId))	//我
	result.SecondPool = mydesk.GetSecondPool()
	*result.TurnMax = mydesk.BetAmountNow
	result.WeixinInfos = mydesk.GetWeiXinInfos()

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


//判断是否是掉线
func isBreak(u *room.ThUser) int32 {
	if u.BreakStatus == room.TH_USER_BREAK_STATUS_TRUE {
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

