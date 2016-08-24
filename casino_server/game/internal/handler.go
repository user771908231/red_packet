package internal

import (
	"github.com/name5566/leaf/gate"
	"reflect"
	"casino_server/common/log"
	"casino_server/service/fruitService"
	"casino_server/msg/bbprotogo"
	"casino_server/service/zjhService"
	"casino_server/conf/intCons"
	"casino_server/service/room"
	"casino_server/service/OGservice"
	"casino_server/service/CSTHService"
	"casino_server/service/userService"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	//水果机相关的业务
	handler(&bbproto.GetIntoRoom{}, handlerGetIntoRoom)
	handler(&bbproto.Shuiguoji{}, handlerShuiguoji)
	handler(&bbproto.ShuiguojiHilomp{}, handlerShuiguojiHilomp)

	//扎金花相关的业务
	handler(&bbproto.ZjhRoom{}, handlerZjhRoom)
	handler(&bbproto.ZjhLottery{}, handlerZjhLottery)
	handler(&bbproto.ZjhMsg{}, handlerZjhMsg)
	handler(&bbproto.ZjhBet{}, handlerZjhBet)
	handler(&bbproto.ZjhReqSeat{}, handlerZjhReqSeat)
	handler(&bbproto.ZjhQueryNoSeatUser{}, handlerZjhQueryNoSeatUser)

	//德州扑克相关的业务
	handler(&bbproto.ThRoom{}, handlerThRoom)
	handler(&bbproto.THBet{}, handlerThBet)                //押注的协议

	//联众的德州扑克
	handler(&bbproto.Game_LoginGame{}, handlerGameLoginGame)          //登陆游戏
	handler(&bbproto.Game_EnterMatch{}, handlerGameEnterMatch)        //进入房间
	handler(&bbproto.Game_FollowBet{}, handlerFollowBet)              //处理押注的请求
	handler(&bbproto.Game_RaiseBet{}, handlerRaise)                   //处理加注的请求
	handler(&bbproto.Game_FoldBet{}, handlerFoldBet)                  //处理弃牌的请求
	handler(&bbproto.Game_CheckBet{}, handlerCheckBet)                //处理让牌的请求
	handler(&bbproto.Game_CreateDesk{}, handlerCreateDesk)            //创建房间
	handler(&bbproto.Game_DissolveDesk{}, handlerDissolveDesk)        //解散房间


	handler(&bbproto.Game_Ready{}, handlerReady)                      //准备游戏
	handler(&bbproto.Game_Begin{}, handlerBegin)                      //开始游戏

	handler(&bbproto.Game_Message{}, handlerGameMessage)
	handler(&bbproto.Game_LeaveDesk{}, handlerLeaveDesk)              //离开桌子

	handler(&bbproto.Game_GameRecord{}, handlerGetGameRecords)        //查询战绩的接口

	handler(&bbproto.Game_TounamentBlind{}, handlerGame_TounamentBlind)
	handler(&bbproto.Game_TounamentRewards{}, handlerGame_TounamentRewards)
	handler(&bbproto.Game_TounamentRank{}, handlerGame_TounamentRank)
	handler(&bbproto.Game_TounamentSummary{}, handlerGame_TounamentSummary)

	handler(&bbproto.Game_MatchList{}, handlerGame_MatchList)         //锦标赛列表

}

/**
	请求进入游戏房间
	1,分配房间(根据游戏类型)
	2,proto中的标志 如果in=1表示进入房间,其他则表示退出房间

 */
func handlerGetIntoRoom(args []interface{}) {
	log.T("进入到 game.handlerGetIntoRoom()")
	m := args[0].(*bbproto.GetIntoRoom)                //请求体
	a := args[1].(gate.Agent)                //连接
	log.T("agent:", &a)
	log.T("请求进入房间的user %v ,in:%v\n", m.GetUserId(), m.GetIn())
	if m.GetIn() == intCons.REQ_TYPE_IN {
		room.SGJRoom.AddAgent(m.GetUserId(), a)
	} else {
		room.SGJRoom.RemoveAgent(m.GetUserId())
	}
}

/**
处理水果机的业务
 */
func handlerShuiguoji(args []interface{}) {
	log.T("进入到 game.handlerShuiguoji()")
	//检测参数是否正确
	m := args[0].(*bbproto.Shuiguoji)                //请求体
	a := args[1].(gate.Agent)
	result, err := fruitService.HandlerShuiguoji(m, a)
	if err != nil {
		log.E(err.Error())
	}

	//给客户端返回数据
	log.N("给客户端返回的数据%v", result)
	a.WriteMsg(result)
}


/**
	处理水果机比大小的业务
 */
func handlerShuiguojiHilomp(args []interface{}) {
	log.T("进入到 game.handlerShuiguojiHilomp()")
	//检测参数是否正确
	m := args[0].(*bbproto.ShuiguojiHilomp)                //请求体
	a := args[1].(gate.Agent)
	result, err := fruitService.HandlerShuiguojiHilomp(m)
	if err != nil {
		log.E(err.Error())
	}
	a.WriteMsg(result)
	log.N("%v", result)

}

/**
	进入扎金花的房间
 */
func handlerZjhRoom(args []interface{}) {
	log.T("进入到扎金花的房间 game.handlerZjhRoom()")
	//检测参数是否正确
	m := args[0].(*bbproto.ZjhRoom)                //请求体
	a := args[1].(gate.Agent)

	//通过serVice来处理
	result, err := zjhService.HandlerZjhRoom(m, a)
	if err != nil {
		log.E(err.Error())
	}

	//处理返回信息
	log.T("得到的结果:", result)

}


/**
	扎金花奖励

 */
func handlerZjhLottery(args []interface{}) {
	log.T("进入到扎金花的房间 game.handlerZjhRoom()")
}


/**
	扎金花房间消息
 */
func handlerZjhMsg(args []interface{}) {
	log.T("进入到扎金花的房间 game.handlerZjhMsg()")
}


/**
扎金花 押注
 */
func handlerZjhBet(args []interface{}) {
	log.T("进入到扎金花的房间 game.handlerZjhBet()")
	//检测参数是否正确
	m := args[0].(*bbproto.ZjhBet)                //请求体
	a := args[1].(gate.Agent)
	zjhService.HandlerZjhBet(m, a)
}


/**
	扎金花请求座位
 */
func handlerZjhReqSeat(args []interface{}) {
	log.T("进入到扎金花的房间 game.handlerZjhReqSeat()")
}


/*
	扎金花请求没有作为的玩家
 */


func handlerZjhQueryNoSeatUser(args []interface{}) {
	log.T("进入到扎金花的房间 game.handlerZjhQueryNoSeatUser()")
}


/**
	退出或者进入德州扑克的房间
 */
func handlerThRoom(args []interface{}) {

}


/**
	处理德州扑克押注
 */
func handlerThBet(args []interface{}) {
}

//登录游戏
func handlerGameLoginGame(args []interface{}) {
	log.T("快速登陆德州扑克游戏")
	m := args[0].(*bbproto.Game_LoginGame)
	log.T("收到的数据[%v]", m)

	a := args[1].(gate.Agent)
	result := &bbproto.Game_LoginGame{}
	result.Result = new(int32)                //默认是0表示成功
	a.WriteMsg(result)
}


//用户创建一个房间
func handlerCreateDesk(args []interface{}) {
	m := args[0].(*bbproto.Game_CreateDesk)
	a := args[1].(gate.Agent)
	log.T("玩家请求创建房间m[%v]", m)

	//需要返回的信息
	result := &bbproto.Game_AckCreateDesk{}
	result.Result = new(int32)
	result.Password = new(string)
	result.DeskId = new(int32)
	result.CreateFee = new(int64)
	result.UserBalance = new(int64)

	//开始创建房间
	desk, err := OGservice.HandlerCreateDesk(m.GetUserId(), m.GetInitCoin(),m.GetPreCoin(), m.GetSmallBlind(), m.GetBigBlind(),m.GetInitCount())
	if err != nil {
		log.E("创建房间失败 errmsg [%v]", err)
		*result.Result = int32(bbproto.DDErrorCode_ERRORCODE_CREATE_DESK_DIAMOND_NOTENOUGH)
	} else {
		*result.Result = 0
		*result.DeskId = desk.Id
		*result.Password = desk.RoomKey
		*result.CreateFee = desk.CreateFee
		*result.UserBalance = userService.GetUserById(desk.DeskOwner).GetDiamond()		//得到用户的余额
	}

	//返回信息
	log.T("创建房间成功,返回的数据[%v]", result)
	a.WriteMsg(result)
}


//解散房间..这个协议只有自定义的房间可以使用
func handlerDissolveDesk(args []interface{}) {
	//解散房间
	m := args[0].(*bbproto.Game_DissolveDesk)
	a := args[1].(gate.Agent)
	log.T("解散房间的请求参数[%v]", m)
	room.ThGameRoomIns.DissolveDeskByDeskOwner(m.GetUserId(), a)
}


// 处理请求进入游戏房间
func handlerGameEnterMatch(args []interface{}) {
	m := args[0].(*bbproto.Game_EnterMatch)
	a := args[1].(gate.Agent)
	//返回房间的信息
	OGservice.HandlerGameEnterMatch(m, a)
}

//处理离开房间的请求
func handlerLeaveDesk(args []interface{}) {
	m := args[0].(*bbproto.Game_LeaveDesk)
	a := args[1].(gate.Agent)
	OGservice.HandlerLeaveDesk(m, a)
}

func handlerGameMessage(args []interface{}) {
	m := args[0].(*bbproto.Game_Message)
	a := args[1].(gate.Agent)
	OGservice.HandlerMessage(m, a)
}


//处理准备游戏
func handlerReady(args []interface{}) {
	m := args[0].(*bbproto.Game_Ready)
	a := args[1].(gate.Agent)
	OGservice.HandlerReady(m, a)
}


//开始游戏的请求
func handlerBegin(args []interface{}) {
	m := args[0].(*bbproto.Game_Begin)
	a := args[1].(gate.Agent)
	OGservice.HandlerBegin(m, a)
}


//处理跟注
func handlerFollowBet(args []interface{}) {
	m := args[0].(*bbproto.Game_FollowBet)
	a := args[1].(gate.Agent)
	seatId := m.GetSeat()
	desk := room.GetDeskByAgent(a)
	desk.DDBet(seatId, room.TH_DESK_BET_TYPE_CALL, 0)
}

// 处理加注
func handlerRaise(args []interface{}) {
	m := args[0].(*bbproto.Game_RaiseBet)
	a := args[1].(gate.Agent)
	seatId := m.GetSeat()
	coin := m.GetCoin()
	desk := room.GetDeskByAgent(a)

	desk.DDBet(seatId, room.TH_DESK_BET_TYPE_RAISE, coin)
}

// 处理弃牌
func handlerFoldBet(args []interface{}) {
	m := args[0].(*bbproto.Game_FoldBet)
	a := args[1].(gate.Agent)
	seatId := m.GetSeat()
	desk := room.GetDeskByAgent(a)

	desk.DDBet(seatId, room.TH_DESK_BET_TYPE_FOLD, 0)
}

// 处理让牌
func handlerCheckBet(args []interface{}) {
	m := args[0].(*bbproto.Game_CheckBet)
	a := args[1].(gate.Agent)

	seatId := m.GetSeat()
	desk := room.GetDeskByAgent(a)

	desk.DDBet(seatId, room.TH_DESK_BET_TYPE_CHECK, 0)
}

//获得个人的战绩,并且按照时间排序
func handlerGetGameRecords(args []interface{}) {
	m := args[0].(*bbproto.Game_GameRecord)
	a := args[1].(gate.Agent)
	OGservice.HandlerGetGameRecords(m, a)
}

func handlerGame_TounamentBlind(args []interface{}) {
	m := args[0].(*bbproto.Game_TounamentBlind)
	a := args[1].(gate.Agent)
	d1 := bbproto.NewGame_TounamentBlindBean()
	*d1.Ante = 10
	*d1.SmallBlind = 10
	m.Data = append(m.Data, d1)

	d2 := bbproto.NewGame_TounamentBlindBean()
	*d2.Ante = 20
	*d2.SmallBlind = 20
	m.Data = append(m.Data, d2)

	d3 := bbproto.NewGame_TounamentBlindBean()
	*d3.Ante = 75
	*d3.SmallBlind = 150
	m.Data = append(m.Data, d3)

	a.WriteMsg(m)
}

func handlerGame_TounamentRewards(args []interface{}) {
	m := args[0].(*bbproto.Game_TounamentRewards)
	a := args[1].(gate.Agent)

	d1 := bbproto.NewGame_TounamentRewardsBean()
	*d1.IconPath = "1"
	*d1.Rewards = "100元红包"
	m.Data = append(m.Data, d1)

	d2 := bbproto.NewGame_TounamentRewardsBean()
	*d2.IconPath = "2"
	*d2.Rewards = "50元红包"
	m.Data = append(m.Data, d2)

	d3 := bbproto.NewGame_TounamentRewardsBean()
	*d3.IconPath = "3"
	*d3.Rewards = "30元红包"
	m.Data = append(m.Data, d3)
	a.WriteMsg(m)
}

func handlerGame_TounamentRank(args []interface{}) {
	m := args[0].(*bbproto.Game_TounamentRank)
	a := args[1].(gate.Agent)

	d1 := bbproto.NewGame_TounamentRankBean()
	*d1.PlayerName = "playerName"
	*d1.PlayerImage = ""
	*d1.Coin = 999
	*d1.Place = 1

	m.Data = append(m.Data, d1)
	a.WriteMsg(m)
}

func handlerGame_TounamentSummary(args []interface{}) {
	m := args[0].(*bbproto.Game_TounamentSummary)
	a := args[1].(gate.Agent)

	log.T("用户请求handlerGame_TounamentSummarym[%v]",m)

	m.Coin = new(string)
	m.Fee = new(string)
	m.PersonCount = new(string)
	m.Time = new(string)

	*m.Fee = "免费"
	*m.Time = "20分钟一场"
	*m.PersonCount = "65-1000人"
	*m.Coin = "5000"
	a.WriteMsg(m)
}

//竞标赛列表
func handlerGame_MatchList(args []interface{}) {
	a := args[1].(gate.Agent)
	data := CSTHService.GetGameMatchList()
	log.T("得到的锦标赛列表[%v]", data)
	a.WriteMsg(data)
}