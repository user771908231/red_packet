package room

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
)


//锦标赛
type CSThGameRoom struct {
	ThGameRoom

	//锦标赛房间的专有属性
	matchId	int32		//比赛内容

}

//run游戏房间
func (r *CSThGameRoom) Run() {

}

//游戏大厅增加一个玩家
func (r *CSThGameRoom) AddUser(userId uint32, roomCoin int64, a gate.Agent) (*ThDesk, error) {
	//进入房间的操作需要加锁
	r.Lock()
	defer r.Unlock()
	log.T("userid【%v】进入德州扑克的房间", userId)

	var mydesk *ThDesk = nil                //为用户找到的desk

	//1,判断用户是否已经在房间里了,如果是在房间里,那么替换现有的agent,
	mydesk = r.IsRepeatIntoRoom(userId, a)
	if mydesk != nil {
		return mydesk, nil
	}

	//2,查询哪个德州的房间缺人:循环每个德州的房间,然后查询哪个房间缺人
	for deskIndex := 0; deskIndex < len(r.ThDeskBuf); deskIndex++ {
		tempDesk := r.ThDeskBuf[deskIndex]
		if tempDesk == nil {
			log.E("找到房间为nil,出错")
			break
		}
		if tempDesk.UserCount < r.ThRoomSeatMax {
			mydesk = tempDesk        //通过roomId找到德州的room
			break;
		}
	}

	//如果没有可以使用的桌子,那么重新创建一个,并且放进游戏大厅
	if mydesk == nil {
		log.T("没有多余的desk可以用,重新创建一个desk")
		mydesk = NewThDesk()
		r.AddThDesk(mydesk)
	}

	//3,进入房间,竞标赛进入房间的时候,默认就是准备的状态
	err := mydesk.AddThUser(userId, roomCoin, TH_USER_STATUS_READY, a)
	if err != nil {
		log.E("用户上德州扑克的桌子 失败...")
		return nil, err
	}

	mydesk.LogString()        //答应当前房间的信息

	return mydesk, nil
}
