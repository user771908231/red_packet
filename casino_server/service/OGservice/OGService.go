package OGservice

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/service/room"
	"casino_server/common/log"
	"casino_server/service/userService"
	"errors"
	"casino_server/mode"
	"casino_server/conf/intCons"
	"casino_server/utils/timeUtils"
	"casino_server/common/Error"
	"casino_server/service/noticeServer"
	"casino_server/mode/dao/TThDeskRecordDao"
)


/**
	用户通过钻石创建游戏房间
	ateDesk(userId,initCoin,roomKey,smallBlind,bigBlind,juCount)
	//新修改:创建房间的时候roomKey 系统指定,不需要用户输入
 */

func HandlerCreateDesk(userId uint32, roomCoin int64, preCoin int64, smallBlind int64, bigBlind int64, jucount int32) (*room.ThDesk, error) {

	//0,首先看用户是否已经创建了房间,如果已经创建了直接返回已经创建的...
	odesk := room.ThGameRoomIns.GetDeskByOwner(userId)
	if odesk != nil {
		return odesk, nil
	}

	//1,得到一个随机的密钥
	roomKey := room.ThGameRoomIns.RandRoomKey()

	//2,开始创建房间
	desk, err := room.ThGameRoomIns.CreateDeskByUserIdAndRoomKey(userId, roomCoin, roomKey, preCoin, smallBlind, bigBlind, jucount);
	if err != nil {
		log.E("用户创建房间失败errMsg[%v]", err.Error())
		return nil, err
	}

	//new : 创建房间之后返回创建房间的结果,不用知己进入房间
	return desk, nil

}

//离开房间
func HandlerLeaveDesk(m *bbproto.Game_LeaveDesk, a gate.Agent) {
	//deskId
	desk := room.GetDeskByAgent(a)
	if desk == nil {
		log.E("HandlerLeaveDesk,失败,因为desk不存在...")
		return
	}
	desk.LeaveThuser(m.GetUserId())
}

///用户发送消息
func HandlerMessage(m *bbproto.Game_Message, a gate.Agent) {
	//deskId
	log.T("用户开始发送消息m[%v]", m)
	result := &bbproto.Game_SendMessage{}
	desk := room.GetDeskByAgent(a)

	if desk != nil {
		result.UserId = m.UserId
		result.Msg = m.Msg
		result.MsgType = m.MsgType
		result.Id = m.Id
		result.Seat = m.Seat
		log.T("开始广播用户[%v]发送的消息[%v]", m.GetUserId(), result)
		desk.THBroadcastProtoAll(result)
	} else {
		log.T("发送广播失败,没有找到对应的desk")
	}
}

//用户准备
func HandlerReady(m *bbproto.Game_Ready, a gate.Agent) error {
	log.T("用户开始准备游戏m[%v]", m)
	//1,找到userId
	result := &bbproto.Game_AckReady{}
	result.Result = new(int32)
	result.Msg = new(string)
	result.SeatId = new(int32)

	userId := m.GetUserId()
	//2,通过userId 找到桌子
	desk := room.GetDeskByAgent(a)
	if desk == nil {
		log.E("用户id[%v]准备的时候,房间不存在", userId)
		return errors.New("房间不存在")
	}

	//3,用户开始准备
	err := desk.Ready(userId)
	if err != nil {
		log.T("用户【%v】准备失败,err[%v]", userId, err.Error())
		*result.Result = Error.GetErrorCode(err)
		*result.Msg = Error.GetErrorMsg(err)
		a.WriteMsg(result)
		return err
	}

	//4,返回准备的结果
	*result.SeatId = desk.GetUserSeatByUserId(userId)
	*result.Result = intCons.ACK_RESULT_SUCC
	desk.THBroadcastProtoAll(result)        //广播用户准备的协议

	//如果全部的人都准备好了,那么可以开始游戏
	//1.1,所有人都准备好了,并且不是第一局的时候,才能开始游戏, 第一句必须要房主点击开始,才能开始
	if desk.JuCountNow > 1 && desk.IsAllReady() {
		//准备之后判断游戏是否开始
		go desk.Run()
	}
	return nil
}

//房主强制开始游戏
func HandlerBegin(m *bbproto.Game_Begin, a gate.Agent) error {
	userId := m.GetUserId()
	desk := room.GetDeskByAgent(a)

	if desk == nil {
		log.E("没有找到桌子，开始游戏失败...")
		return errors.New("没有找到房间")
	}

	//如果是朋友桌的开始方式
	if desk.IsFriend() {
		if desk.DeskOwner != userId {
			log.E("没有找到房主为[%v]的desk", userId)
			return errors.New("没有找到房间")
		} else {
			//开始游戏
			go desk.Run()
			return nil
		}
	}

	//如果是锦标赛的开始方式
	if desk.IsChampionship() {
		go desk.Run()
		return nil
	}

	return errors.New("游戏开始失败")

}

//
func HandlerGetGameRecords(m *bbproto.Game_GameRecord, a gate.Agent) {
	userId := m.GetUserId()

	//1,战绩查询
	var deskRecords []mode.T_th_desk_record = TThDeskRecordDao.Find(userId, 20)
	log.T("查询到用户[%v]的战绩[%v]", userId, deskRecords)
	//需要返回的数据
	resDatav2 := bbproto.NewGame_AckGameRecord()
	*resDatav2.UserId = userId
	*resDatav2.Result = intCons.ACK_RESULT_SUCC

	//第一层循环战绩
	for i := 0; i < len(deskRecords); i++ {
		tr := deskRecords[i]
		r := bbproto.NewGame_BeanGameRecord()
		*r.DeskId = tr.DeskId
		log.T("找到的时间[%v]", tr.BeginTime)
		*r.BeginTime = timeUtils.Format(tr.BeginTime)

		//第二层循环人
		for j := 0; j < len(tr.Records); j++ {
			tru := tr.Records[j]

			ru := bbproto.NewGame_BeanUserRecord()
			*ru.UserId = tru.UserId
			*ru.NickName = tru.NickName
			*ru.WinAmount = tru.WinAmount
			r.Users = append(r.Users, ru)
		}

		resDatav2.Records = append(resDatav2.Records, r)                //加入要返回的结果
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


	//进入房间的时候  需要更具游戏类型来做不同的进入逻辑...
 */
func HandlerGameEnterMatch(m *bbproto.Game_EnterMatch, a gate.Agent) error {
	log.T("用户请求进入德州扑克的游戏房间,HandlerGameEnterMatch(m[%v])", m)

	var err error                                   //错误信息
	var mydesk *room.ThDesk                         //用户需要进入的房间
	userId := m.GetUserId()                         //进入游戏房间的user
	roomKey := ""             //房间的roomkey
	passWord := m.GetPassWord()
	if passWord != nil {
		roomKey = string(passWord)
	}
	matchId := m.GetMatchIdInt()                        //进入锦标赛的时候检测锦标赛的matchId


	//1.1 检测参数是否正确,判断userId 是否合法
	userCheck := userService.CheckUserIdRightful(userId)
	if userCheck == false {
		log.E("进入德州扑克的房间的时候,userId[%v]不合法。", userId)
		return errors.New("用户Id不合法")
	}

	//1.2,进入房间,返回房间和错误信息
	if roomKey == "" {
		//这里可以尝试从session 获取游戏类型..
		r := room.GetCSTHroom(matchId)
		if r == nil {
			log.E("用户[%v]请求进入游戏desk的时候失败，因为没有找到matchId[%v]对应的锦标赛", userId, matchId)
			mydesk = nil        //没有找到对应的桌子，进入房间失败...
			err = errors.New("没有找到对应的锦标赛desk")
		} else {
			mydesk, err = room.GetCSTHroom(matchId).AddUser(userId, a)
		}
	} else {
		mydesk, err = room.ThGameRoomIns.AddUserWithRoomKey(userId, roomKey, a)
	}

	//2 判断进入房间是否失败...
	if err != nil || mydesk == nil {
		log.E("用户[%v]进入房间失败,err[%v]", userId, err)
		//返回错误信息
		result := bbproto.NewGame_SendGameInfo()
		*result.Result = Error.GetErrorCode(err)
		a.WriteMsg(result)
		return err
	}

	//进入房间成功之后,发送房间当前的广播
	mydesk.BroadGameInfo(userId)

	//保存信息
	mydesk.UpdateThdeskAndUser2redis(mydesk.GetUserByUserId(userId))

	return nil
}

//处理bbproto.Game_login
func HandlerGameLogin(userId uint32, a gate.Agent) {
	log.T("用户[%v]请求gameLogin", userId)
	session := userService.GetUserSessionByUserId(userId)
	log.T("用户的回话信息:session[%v]", session)

	//返回session之前需要检测session的合法性
	ret := bbproto.NewGame_AckLogin()
	if !room.CheckUserSessionRight(session) {
		log.E("用户[%v]的session[%v]信息有误，请管理查看", userId, session)
		*ret.MatchId = 0
		*ret.TableId = 0
		*ret.GameStatus = room.TH_USER_GAME_STATUS_NOGAME
		*ret.RoomPassword = ""
	} else {
		*ret.MatchId = session.GetMatchId()
		*ret.TableId = session.GetDeskId()
		*ret.GameStatus = session.GetGameStatus()
		*ret.RoomPassword = session.GetRoomKey()
	}

	*ret.Notice = noticeServer.GetNoticeByType(noticeServer.NOTICE_TYPE_GUNDONG).GetNoticeContent()        //滚动信息
	*ret.CostRebuy = int64(1)
	*ret.Championship = false                //锦标赛是否开启
	*ret.Chip = userService.GetUserById(userId).GetDiamond()
	*ret.CostCreateRoom = room.ThdeskConfig.CreateFee        //创建房间的单价

	a.WriteMsg(ret)
}
