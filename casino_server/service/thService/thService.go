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

	//判断是否已经在房间:这里可以用过agent user data 来判断
	//agentUser := a.UserData().(*gamedata.AgentUserData{})
	//if agentUser.Status == gamedata.AGENT_USER_STATUS_GAMING && agentUser.ZhDeskId > 0 {
	//	log.E("用户已经在房间中了,请不要重复进入")
	//	return errors.New("玩家已经在房间中了,请不要重复进入")
	//}




	//2,查询哪个德州的房间缺人:循环每个德州的房间,然后查询哪个房间缺人
	var mydesk *room.ThDesk = nil
	var index int = 0
	if len(room.ThGameRoomIns.ThDeskBuf) > 0 {
		log.T("当前拥有的ThDesk 的数量[%v]",len(room.ThGameRoomIns.ThDeskBuf))
		for  deskIndex := 0; deskIndex < len(room.ThGameRoomIns.ThDeskBuf); deskIndex++ {
			if room.ThGameRoomIns.ThDeskBuf[deskIndex] !=nil {
				mydesk = room.ThGameRoomIns.ThDeskBuf[deskIndex]        //通过roomId找到德州的room
				if *mydesk.SeatedCount < *room.ThGameRoomIns.ThRoomSeatMax {
					log.T("roomid[%v]有空的座位,", deskIndex)
					break;
				}
			}else{
				index = deskIndex
				log.T("deskId[%v]为nil,直接返回,", deskIndex)
				break
			}

		}
	}

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
	*mydesk.SeatedCount = *mydesk.SeatedCount + 1

	//上了牌桌之后,如果玩家人数大于1,并且游戏处于stop的状态,则直接开始游戏
	if *mydesk.SeatedCount >= room.TH_DESK_LEAST_START_USER  && *mydesk.Status == room.TH_DESK_STATUS_STOP{
		err = mydesk.Run()
		if err != nil {
			log.E("开始德州扑克游戏的时候失败")
			return nil
		}
	}
	//4,修改gameRoom的状态
	result := &bbproto.ThRoom{}
	result.Header = protoUtils.GetSuccHeaderwithUserid(m.GetHeader().UserId)
	result.DeskStatus = mydesk.Status		//当前桌子的状态
	result.PublicPais = mydesk.PublicPai		//公共牌
	result.Users = mydesk.GetResUserModel()

	log.T("返回信息",result)
	a.WriteMsg(result)

	//5,返回结果
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
	//找到游戏的桌子号




	return nil
}


