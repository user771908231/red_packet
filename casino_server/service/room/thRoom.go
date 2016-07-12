package room

import (
	"casino_server/msg/bbprotogo"
	"sync"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
)



//config
var THROOM_SEAT_COUNT int32 = 8		//玩德州扑克,每个房间最多多少人
var GAME_THROOM_MAX_COUNT = 500		//一个游戏大厅最多有多少桌德州扑克


/**
	初始化函数:
		初始化游戏房间
 */

var ThGameRoomIns ThGameRoom	//房间实例,在init函数中初始化

func init(){
	ThGameRoomIns.OnInit()		//初始化房间
	ThGameRoomIns.Run()		//运行房间
}

/**
	德州扑克
 */

//游戏房间
type ThGameRoom struct {
	sync.Mutex
	RoomStatus	*int32		//游戏大厅的状态
	ThRoomBuf 	[]*ThDesk
	ThRoomSeatMax 	*int32		//每个房间的座位数目
	ThRoomCount	*int32		//房间数目
}


//初始化游戏房间
func (r *ThGameRoom) OnInit(){
	r.ThRoomBuf = make([]*ThDesk,GAME_THROOM_MAX_COUNT)
	r.ThRoomSeatMax = &THROOM_SEAT_COUNT
}

//run游戏房间
func (r *ThGameRoom) Run(){

}


//增加一个thRoom
func (r *ThGameRoom) AddThRoom(throom *ThDesk) error{
	var result error = nil
	for i := 0; i < GAME_THROOM_MAX_COUNT; i++ {
		if r.ThRoomBuf[i] == nil {
			r.ThRoomBuf[i] = throom
			break;
		}
	}
	if result != nil {
		log.E("增加德州扑克的桌子失败")
	}
	return result

}


/**
	一个德州扑克的房间
 */
type ThDesk struct {
	room
	Id	*uint32		//roomid
	Dealer  *uint32		//荷官的UserId
	PublicPai []*bbproto.Pai	//公共牌的部分
	SeatedCount	*int32	//已经坐下的人数
	Agents  []gate.Agent	//坐下的人
}

/**
	为桌子增加一个人
 */
func (t ThDesk) addAgent() error{
	return nil
}


/**
	一个德州扑克的座位,座位包含一下信息:
	人
	牌
 */
type ThSeat struct {
	User	*bbproto.User		//座位上的人
	HandPais	[]*bbproto.Pai	//手牌
	THPork	*bbproto.THPork		//德州的牌
}

/**
	新生成一个德州的桌子
 */
func NewThDesk() *ThDesk{

	result := new(ThDesk)
	result.SeatedCount = new(int32)
	result.Dealer = new(uint32)
	result.AgentMap = make(map[uint32]gate.Agent)

	return result

}