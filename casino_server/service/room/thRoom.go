package room

import (
	"casino_server/msg/bbprotogo"
	"sync"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"errors"
	"github.com/golang/protobuf/proto"
	"casino_server/msg/bbprotoFuncs"
)



//config
var THROOM_SEAT_COUNT int32 = 8                //玩德州扑克,每个房间最多多少人
var GAME_THROOM_MAX_COUNT int32 = 500                //一个游戏大厅最多有多少桌德州扑克
var TH_DESK_LEAST_START_USER int32 = 2        //最少多少人可以开始游戏

//德州扑克 玩家的状态
var TH_USER_STATUS_WAITSEAT int32 = 1; //刚上桌子 等待开始的玩家


//德州扑克,牌桌的状态
var TH_DESK_STATUS_STOP int32 = 1; //没有开始的状态
var TH_DESK_STATUS_SART int32 = 1; //没有已经开始的状态


/**
	初始化函数:
		初始化游戏房间
 */

var ThGameRoomIns ThGameRoom        //房间实例,在init函数中初始化

func init() {
	ThGameRoomIns.OnInit()                //初始化房间
	ThGameRoomIns.Run()                //运行房间
}

/**
	德州扑克
 */

//游戏房间
type ThGameRoom struct {
	sync.Mutex
	RoomStatus    *int32 //游戏大厅的状态
	ThRoomBuf     []*ThDesk
	ThRoomSeatMax *int32 //每个房间的座位数目
	ThRoomCount   *int32 //房间数目
}


//初始化游戏房间
func (r *ThGameRoom) OnInit() {
	r.ThRoomBuf = make([]*ThDesk, GAME_THROOM_MAX_COUNT)
	r.ThRoomSeatMax = &THROOM_SEAT_COUNT
}

//run游戏房间
func (r *ThGameRoom) Run() {

}


//增加一个thRoom
func (r *ThGameRoom) AddThRoom(throom *ThDesk) error {
	var result error = nil
	for i := 0; i < int(GAME_THROOM_MAX_COUNT); i++ {
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
	正在玩德州的人
 */
type ThUser struct {
	userId *uint32    //用户id
	agent  gate.Agent //agent
	status *int32     //当前的状态
}

func NewThUser() *ThUser {
	result := &ThUser{}
	result.userId = new(uint32)
	result.status = new(int32)
	return result
}

/**
	一个德州扑克的房间
 */
type ThDesk struct {
	Id          *uint32        //roomid
	Dealer      *uint32        //荷官的UserId
	PublicPai   []*bbproto.Pai //公共牌的部分
	SeatedCount *int32         //已经坐下的人数
	userSeated  []*ThUser      //坐下的人
	userWait    []*ThUser      //等待开始的人
	Status      *int32         //牌桌的状态
}



/**
	为桌子增加一个人
 */
func (t *ThDesk) AddThUser(userId uint32, a gate.Agent) error {

	//通过userId 和agent 够做一个thuser
	thUser := NewThUser()
	*thUser.userId = userId
	thUser.agent = a
	*thUser.status = TH_USER_STATUS_WAITSEAT        //刚进房间的玩家

	//添加thuser
	for i := 0; i < len(t.userWait); i++ {
		if t.userWait[i] == nil {
			t.userWait[i] = thUser
			break
		}
		if (i + 1) == len(t.userWait) {
			return errors.New("加入房间失败")
		}
	}

	return nil
}

/**
	开始游戏
 */
func (t *ThDesk) Run() error {

	//设置房间状态
	*t.Status = TH_DESK_STATUS_SART                //设置状态为开始游戏

	//把正在等待的用户设置可以游戏
	err := t.UserWait2Seat()
	if err != nil {
		log.E("开始游戏失败,errMsg[%v]", err.Error())
		return err
	}

	//初始化牌的信息
	err = t.OnInitPork()
	if err != nil {
		log.E("开始德州扑克游戏,初始化扑克牌的时候出错")
	}

	//广播消息
	res := &bbproto.THBegin{}
	res.Header = protoUtils.GetSuccHeader()
	res.Users = t.GetResUserModel()
	err = t.THBroadcastProto(res, 0)
	if err != nil {
		log.E("广播开始消息的时候出错")
		return err
	}

	return nil
}

/**
	把正在等待的用户安置在座位上
 */
func (t *ThDesk) UserWait2Seat() error {

	//打印移动之前的信息
	log.T("UserWait2Seat 之前,t.userWait,[%v]",t.userWait)
	log.T("UserWait2Seat 之前,t.userSeated,[%v]",t.userSeated)


	for i := 0; i < len(t.userWait); i++ {
		for j := 0; j < len(t.userSeated); j++ {
			if t.userSeated[j] == nil {
				t.userSeated[j] = t.userWait[i]
				break
			}
		}
	}

	//打印测试消息
	log.T("UserWait2Seat 之后,t.userWait,[%v]",t.userWait)
	log.T("UserWait2Seat 之后,t.userSeated,[%v]",t.userSeated)

	return nil
}


/**
	初始化纸牌的信息
 */
func (t *ThDesk) OnInitPork() error {
	return nil

}


/**
	德州扑克广播消息
 */
func (t *ThDesk) THBroadcastProto(p proto.Message, ignoreUserId int32) error {
	log.Normal("给每个房间发送proto 消息%v", p)
	for i := 0; i < len(t.userSeated); i++ {
		log.Normal("开始userId[%v]发送消息", t.userSeated[i].userId)
		a := t.userSeated[i].agent
		a.WriteMsg(p)
		log.Normal("给userId[%v]发送消息,发送完毕", t.userSeated[i].userId)
	}

	return nil
}


/**
	返回res需要的User实体
 */
func (t *ThDesk) GetResUserModel() []*bbproto.THUser {

	return nil
}




/**
	一个德州扑克的座位,座位包含一下信息:
	人
	牌
 */
type ThSeat struct {
	User     *bbproto.User   //座位上的人
	HandPais []*bbproto.Pai  //手牌
	THPork   *bbproto.THPork //德州的牌
}

/**
	新生成一个德州的桌子
 */
func NewThDesk() *ThDesk {
	result := new(ThDesk)
	result.SeatedCount = new(int32)
	result.Dealer = new(uint32)
	result.Status = new(int32)

	return result
}