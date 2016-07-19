package dzService

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/conf/intCons"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/common/log"
	"casino_server/service/room"
	"errors"
)

/**

	处理退出或者进入房间的请求
	判断是进入房间还是退出房间是通过reqType 来确定的...
 */
func HandlerThRoom(m *bbproto.ThRoom, a gate.Agent) error {
	//退出房间,和进入房间都需要加锁
	room.ThGameRoomIns.Lock()
	defer room.ThGameRoomIns.Unlock()

	var err error
	//result := *&bbproto.ThRoom{}
	reqType := m.GetReqType()
	if reqType == intCons.REQ_TYPE_IN {
		//表示进入德州扑克的房间
		err = getIntoRoom(m, a)
	} else {
		//表示退出德州扑克的房间
		err = getOutRoom(m, a)
	}

	if err != nil {
		log.E("HandlerThRoom()出错,errMsg[%v]",err.Error())
	}

	return err
}


/**
	进入房间
	1,修改gameRoom的数据,增加房间人数等
	2,判断那个房间缺人
	3,管理agent

	4,如果正在游戏中,需要吧当前的状态放回给玩家,比如公共牌什么的
	如果进入房间之后,游戏正在进行中,则等待游戏完成,进入下一轮
 */
func getIntoRoom(m *bbproto.ThRoom, a gate.Agent) error {
	//定义需要的参数
	userId := m.GetHeader().GetUserId()	//请求的用户id
	result := &bbproto.ThRoom{}		//需要返回的信息

	//1,进入房间,放回房间和错误信息
	mydesk,err := room.ThGameRoomIns.AddUser(userId,a)
	if err != nil || mydesk == nil{
		log.E("进入房间失败,errMsg[%v]",err.Error())
		//这里需要给客户端返回失败的信息
		errMsg := err.Error()
		result.Header = protoUtils.GetErrorHeaderWithMsgUserid(m.GetHeader().UserId,&errMsg)
		a.WriteMsg(result)
		return err
	}

	//4登陆成功的处理,给请求登陆的玩家返回登陆结果的消息
	log.T("开始给客户端返回信息")
	result.Header = protoUtils.GetSuccHeaderwithUserid(m.GetHeader().UserId)
	result.DeskStatus = &(mydesk.Status)		//当前桌子的状态
	result.PublicPais = mydesk.PublicPai		//公共牌
	result.Users = mydesk.GetResUserModelClieSeq(userId)
	log.T("返回信息",result)
	log.T("返回信息Users",result.Users)
	a.WriteMsg(result)

	//5,进入房间的广播,告诉其他人有新的玩家进来了
	mydesk.THBroadcastAddUser(userId,userId)

	//6,最后:确定是否开始游戏, 上了牌桌之后,如果玩家人数大于1,并且游戏处于stop的状态,则直接开始游戏
	//这是游戏刚开始,的处理方式
	if mydesk.userCount >= room.TH_DESK_LEAST_START_USER  && mydesk.Status == room.TH_DESK_STATUS_STOP{
		err = mydesk.Run()
		if err != nil {
			log.E("开始德州扑克游戏的时候失败")
			return nil
		}
	}
	return nil
}

/**
	退出房间
 */
func getOutRoom(m *bbproto.ThRoom, a gate.Agent) error {
	//1,判断参数是否正确
	userId := m.GetHeader().GetUserId()
	//2,找到用户所在的房间
	desk := room.ThGameRoomIns.GetDeskByUserId(userId)
	//3,退出房间
	desk.RmThuser(userId)
	//4,修改thgame的值
	if desk.UserCount == 0 {
		//表示这个房间已经没有人了
		room.ThGameRoomIns.RmThroom(desk.Number)
	}
	
	return nil
}

/**
	处理德州扑克押注的问题
 */
func HandlerTHBet(m *bbproto.THBet, a gate.Agent) error {
	//通过桌子号找到桌子
	desk := room.ThGameRoomIns.GetDeskByUserId(m.GetHeader().GetUserId())
	if desk == nil {
		return errors.New("没有找到id[%v]对应的桌子")
	}

	//开始进行押注
	err := desk.Bet(m,a)
	if err != nil {
		log.E("用户[%v]在桌子[%v]押注的时候出错errMsg[%v]",m.GetHeader().GetUserId(),desk.Id,err.Error())
	}

	//返回错误信息
	return err
}

