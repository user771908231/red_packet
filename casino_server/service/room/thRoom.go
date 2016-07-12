package room

import "casino_server/msg/bbprotogo"


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
	thRoomBuf map[int32]*ThRoom
}


//初始化游戏房间
func (r *ThGameRoom) OnInit(){
	r.thRoomBuf = make(map[int32]*ThRoom)
}

//run游戏房间
func (r *ThGameRoom) Run(){

}



/**
	一个德州扑克的房间
 */
type ThRoom struct {
	room
	Id	*uint32		//roomid
	Dealer  *uint32		//荷官的UserId
	PublicPai []*bbproto.Pai	//公共牌的部分
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