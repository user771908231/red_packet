package dzService

import (
	"casino_server/msg/bbprotogo"
	"github.com/name5566/leaf/gate"
	"casino_server/conf/intCons"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/service/userService"
	"casino_server/common/log"
	"casino_server/service/room"
)

/**
	处理退出或者进入房间的请求
 */
func HandlerThRoom(m *bbproto.ThRoom, a gate.Agent) error {
	var err error
	result := *&bbproto.ThRoom{}
	reqType := m.GetReqType()
	if reqType == intCons.REQ_TYPE_IN {
		//表示进入德州扑克的房间
		err = getIntoRoom(m, a)
	} else {
		//表示退出德州扑克的房间
		err = getOutRoom(m, a)
	}

	if err != nil {
		//表示进入或者退出房间出错了
		errMsg := string(err.Error())
		result.Header = protoUtils.GetErrorHeaderWithMsg(&errMsg)
	} else {
		//退出或者进入房间成功
		result.Header = protoUtils.GetSuccHeader()
	}

	//向客户端返回结果
	a.WriteMsg(result)
	return err
}



/**
	进入房间
	1,修改gameRoom的数据,增加房间人数等
	2,判断那个房间缺人
	3,管理agent

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
	}


	//2,查询那个德州的房间缺人:循环每个德州的房间,然后查询哪个房间缺人
	var intoRoom *room.ThDesk = nil
	if len(room.ThGameRoomIns.ThRoomBuf) > 0 {
		for roomId := range room.ThGameRoomIns.ThRoomBuf {
			log.T("roomid[%v]", roomId)
			intoRoom = room.ThGameRoomIns.ThRoomBuf[roomId]        //通过roomId找到德州的room
			if *intoRoom.SeatedCount < *room.ThGameRoomIns.ThRoomSeatMax {
				break;
			}
		}
	}

	if len(room.ThGameRoomIns.ThRoomBuf) == 0 || intoRoom == nil {
		intoRoom = room.NewThDesk()
		room.ThGameRoomIns.AddThRoom(intoRoom)
	}

	//3,进入房间
	err := intoRoom.AddThUser(userId, a)
	if err != nil {
		log.E("用户上德州扑克的桌子 失败...")
		return err
	}
	*intoRoom.SeatedCount = *intoRoom.SeatedCount + 1


	//上了牌桌之后,如果玩家人数大于1,并且游戏处于stop的状态,则直接开始游戏
	if *intoRoom.SeatedCount >= room.TH_DESK_LEAST_START_USER  && *intoRoom.Status == room.TH_DESK_STATUS_STOP{
		err = intoRoom.Run()
		if err != nil {
			log.E("开始德州扑克游戏的时候失败")
			return nil
		}
	}
	//4,修改gameRoom的状态



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
	return nil
}


