package OGservice

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/service/room"
	"casino_server/common/log"
	"casino_server/service/userService"
	"errors"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/mode"
	"casino_server/conf/intCons"
	"casino_server/utils/timeUtils"
	"casino_server/utils/numUtils"
	"casino_server/utils/db"
	"gopkg.in/mgo.v2"
	"casino_server/common/Error"
)


/**
	用户通过钻石创建游戏房间
	ateDesk(userId,initCoin,roomKey,smallBlind,bigBlind,juCount)
	//新修改:创建房间的时候roomKey 系统指定,不需要用户输入
 */

func HandlerCreateDesk(userId uint32, roomCoin int64, preCoin int64,smallBlind int64, bigBlind int64, jucount int32) (*room.ThDesk, error) {

	//1,得到一个随机的密钥
	roomKey := room.ThGameRoomIns.RandRoomKey()

	//2,开始创建房间
	desk, err := room.ThGameRoomIns.CreateDeskByUserIdAndRoomKey(userId, roomCoin, roomKey,preCoin, smallBlind, bigBlind, jucount);
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
	desk.LeaveThuser(m.GetUserId())

	//如果是竞标赛,在线人数-1
	room := room.GetCSTHroom(1)
	if room != nil {
		room.SubOnlineCount()
	}

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
	}
}

//用户准备
func HandlerReady(m *bbproto.Game_Ready, a gate.Agent) error {
	log.T("用户开始准备游戏m[%v]", m)
	//1,找到userId
	result := &bbproto.Game_AckReady{}
	result.Result = new(int32)
	userId := m.GetUserId()

	//2,通过userId 找到桌子
	//desk := room.ThGameRoomIns.GetDeskByUserId(userId)
	//
	desk := room.GetDeskByAgent(a)
	if desk == nil {
		log.E("房间不存在")
		return errors.New("房间不存在")
	}

	//3,用户开始准备
	err := desk.Ready(userId)
	if err != nil {
		*result.Result = Error.GetErrorCode(err)
		a.WriteMsg(result)
		return err
	}

	//4,返回准备的结果
	*result.Result = intCons.ACK_RESULT_SUCC
	a.WriteMsg(result)

	//如果全部的人都准备好了,那么可以开始游戏
	desk.Run()

	return nil
}

//开始游戏
func HandlerBegin(m *bbproto.Game_Begin, a gate.Agent) error {
	userId := m.GetUserId()
	//desk := room.ThGameRoomIns.GetDeskByDeskOwner(userId)
	desk := room.GetDeskByAgent(a)
	if desk == nil || desk.DeskOwner != userId {
		log.E("没有找到房主为[%v]的desk", userId)
		return errors.New("没有找到房间")
	} else {
		//开始游戏
		desk.Run()
		return nil
	}
}

//
func HandlerGetGameRecords(m *bbproto.Game_GameRecord, a gate.Agent) {
	userId := m.GetUserId()

	//1,战绩查询
	var deskRecords []mode.T_th_desk_record
	queryKey, _ := numUtils.Uint2String(userId)
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_TH_DESK_RECORD).Find(bson.M{"userids": bson.RegEx{queryKey, "."}}).Limit(20).All(&deskRecords)
	})

	log.T("查询到用户[%v]的战绩[%v]", userId, deskRecords)
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

		*r.DeskId = tr.DeskId
		*r.BeginTime = timeUtils.Format(tr.BeginTime)

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
 */
func HandlerGameEnterMatch(m *bbproto.Game_EnterMatch, a gate.Agent) error {
	log.T("用户请求进入德州扑克的游戏房间,HandlerGameEnterMatch(m[%v])", m)

	var err error                                   //错误信息
	var mydesk *room.ThDesk                         //用户需要进入的房间
	userId := m.GetUserId()                         //进入游戏房间的user
	roomCoin := int64(1000)                         //to do 暂时设置为1000
	roomKey := string(m.GetPassWord())              //房间的roomkey


	//1.1 检测参数是否正确,判断userId 是否合法
	userCheck := userService.CheckUserIdRightful(userId)
	if userCheck == false {
		log.E("进入德州扑克的房间的时候,userId[%v]不合法。", userId)
		return errors.New("用户Id不合法")
	}

	//1.2,进入房间,返回房间和错误信息
	if roomKey == "" {
		mydesk, err = room.ChampionshipRoom.AddUser(userId, roomCoin, a)
	} else {
		mydesk, err = room.ThGameRoomIns.AddUserWithRoomKey(userId, roomCoin, roomKey, a)
	}

	//2 判断进入房间是否失败...
	if err != nil || mydesk == nil {
		errMsg := err.Error()
		log.E("用户[%v]进入房间失败,errMsg[%v]", userId, errMsg)

		//返回错误信息
		result := &bbproto.Game_SendGameInfo{}
		result.Result = new(int32)
		*result.Result = Error.GetErrorCode(err)
		a.WriteMsg(result)
		return err
	}

	mydesk.OGTHBroadAddUser()


	//5,如果是朋友桌的话,需要房主点击开始才能开始...,如果是锦标赛,则自动开始游戏
	if mydesk.DeskType == intCons.GAME_TYPE_TH_CS {
		go mydesk.Run()
	}

	return nil
}






//通过agent返回UserId
func getUserIdByAgent(a gate.Agent) uint32 {
	//获取agent中的userData
	ad := a.UserData()
	if ad == nil {
		log.E("agent中的userData为nil")
		return uint32(0)
	}

	userData := ad.(*bbproto.ThServerUserSession)
	log.T("得到的UserAgent中的userId是[%v]", userData.UserId)
	//return userData.UserId
	//测试代码,返回10006
	return 10006
}
