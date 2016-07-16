package dzService

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/conf/intCons"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/service/userService"
	"casino_server/common/log"
	"casino_server/service/room"
	"errors"
)

/**
	处理退出或者进入房间的请求
 */
func HandlerThRoom(m *bbproto.ThRoom, a gate.Agent) error {
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

	//if err != nil {
	//	//表示进入或者退出房间出错了
	//	errMsg := string(err.Error())
	//	result.Header = protoUtils.GetErrorHeaderWithMsg(&errMsg)
	//} else {
	//	//退出或者进入房间成功
	//	result.Header = protoUtils.GetSuccHeader()
	//}
	//
	////向客户端返回结果
	//a.WriteMsg(result)
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
	//进入房间需要加锁
	room.ThGameRoomIns.Lock()
	defer room.ThGameRoomIns.Unlock()


	//1,判断参数是否正确
	userId := m.GetHeader().GetUserId()
	userCheck := userService.CheckUserIdRightful(userId)
	if userCheck == false {
		log.E("用户[%v]不合法", userId)
		return errors.New("用户Id不合法")
	}

	log.T("userid【%v】进入德州扑克的房间",userId)


	//判断是否已经在房间:这里可以用过agent user data 来判断
	//关于有没有重复进入房间,通过agentUser来判断可能不太准确
	isRepeat := room.ThGameRoomIns.IsRepeatIntoRoom(userId)
	if isRepeat {
		log.E("重复进入房间了")
		return errors.New("重复进入房间")
	}

	//2,查询哪个德州的房间缺人:循环每个德州的房间,然后查询哪个房间缺人
	var mydesk *room.ThDesk = nil
	var index int = 0
	if len(room.ThGameRoomIns.ThDeskBuf) > 0 {
		//log.T("当前拥有的ThDesk 的数量[%v]",len(room.ThGameRoomIns.ThDeskBuf))
		for  deskIndex := 0; deskIndex < len(room.ThGameRoomIns.ThDeskBuf); deskIndex++ {
			if room.ThGameRoomIns.ThDeskBuf[deskIndex] !=nil {
				mydesk = room.ThGameRoomIns.ThDeskBuf[deskIndex]        //通过roomId找到德州的room
				mydesk.LogString()
				log.T("每个desk限制的最大人数是[%v]",*room.ThGameRoomIns.ThRoomSeatMax)
				if *mydesk.SeatedCount < *room.ThGameRoomIns.ThRoomSeatMax {
					log.T("room.index[%v]有空的座位,", deskIndex)
					break;
				}
			}else{
				mydesk = nil
				index = deskIndex
				log.T("deskId[%v]为nil,直接返回,", deskIndex)
				break
			}

		}
	}


	//如果没有可以使用的桌子,那么重新创建一个,并且放进游戏大厅
	if len(room.ThGameRoomIns.ThDeskBuf) == 0 || mydesk == nil {
		log.T("没有多余的desk可以用,重新创建一个desk")
		mydesk = room.NewThDesk()
		room.ThGameRoomIns.AddThRoom(index,mydesk)
	}

	//3,进入房间
	err := mydesk.AddThUser(userId, a)
	if err != nil {
		log.E("用户上德州扑克的桌子 失败...")
		return err
	}

	mydesk.LogString()	//答应当前房间的信息

	//4,进入房间的玩家的信息
	log.T("开始给客户端返回信息")
	result := &bbproto.ThRoom{}
	result.Header = protoUtils.GetSuccHeaderwithUserid(m.GetHeader().UserId)
	result.DeskStatus = mydesk.Status		//当前桌子的状态
	result.PublicPais = mydesk.PublicPai		//公共牌
	result.Users = mydesk.GetResUserModelClieSeq(userId)

	log.T("返回信息",result)
	log.T("返回信息Users",result.Users)
	a.WriteMsg(result)

	//5,进入房间的广播,告诉其他人有新的玩家进来了
	mydesk.THBroadcastAddUser(userId,userId)

	//6,最后:确定是否开始游戏, 上了牌桌之后,如果玩家人数大于1,并且游戏处于stop的状态,则直接开始游戏
	//这是游戏刚开始,的处理方式
	if *mydesk.SeatedCount >= room.TH_DESK_LEAST_START_USER  && *mydesk.Status == room.TH_DESK_STATUS_STOP{
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
	//2,找到用户所在的房间
	//3,退出房间
	//4,修改德州房间的状态
	//5,修改gameRoom的状态
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
		log.E("用户[%v]在桌子[%v]押注的时候出错errMsg[%v]",m.GetHeader().GetUserId(),*desk.Id,err.Error())
	}

	//返回错误信息
	return err
}

