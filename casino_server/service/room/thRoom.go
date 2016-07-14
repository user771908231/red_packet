package room

import (
	"casino_server/msg/bbprotogo"
	"sync"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"errors"
	"github.com/golang/protobuf/proto"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/service/userService"
	"casino_server/service/pokerService"
)



//config
var THROOM_SEAT_COUNT int32 = 8               //玩德州扑克,每个房间最多多少人
var GAME_THROOM_MAX_COUNT int32 = 500         //一个游戏大厅最多有多少桌德州扑克
var TH_DESK_LEAST_START_USER int32 = 2        //最少多少人可以开始游戏

//德州扑克 玩家的状态
var TH_USER_STATUS_WAITSEAT 	int32 = 1; //刚上桌子 等待开始的玩家
var TH_USER_STATUS_SEATED 	int32 = 2; //刚上桌子 等待开始的玩家



//德州扑克,牌桌的状态
var TH_DESK_STATUS_STOP int32 = 1; //没有开始的状态
var TH_DESK_STATUS_SART int32 = 2; //没有已经开始的状态


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
	ThDeskBuf     []*ThDesk
	ThRoomSeatMax *int32 //每个房间的座位数目
	ThRoomCount   *int32 //房间数目
}


//初始化游戏房间
func (r *ThGameRoom) OnInit() {
	r.ThDeskBuf = make([]*ThDesk, GAME_THROOM_MAX_COUNT)
	r.ThRoomSeatMax = &THROOM_SEAT_COUNT
}

//run游戏房间
func (r *ThGameRoom) Run() {

}


//增加一个thRoom
func (r *ThGameRoom) AddThRoom(index int,throom *ThDesk) error {
	r.ThDeskBuf[index] = throom
	return nil
}


/**
	正在玩德州的人
 */
type ThUser struct {
	userId *uint32    //用户id
	agent  gate.Agent //agent
	status *int32     //当前的状态
	cards  []*bbproto.Pai
}

//
func (t *ThUser) trans2bbprotoThuser() *bbproto.THUser{
	thuserTemp := &bbproto.THUser{}
	thuserTemp.Status = t.status	//已经就做
	thuserTemp.U =userService.GetUserById(*t.userId) 	//得到user
	thuserTemp.HandPais = t.cards
	return thuserTemp
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
	users	   []*ThUser      //坐下的人
	Status      *int32         //牌桌的状态
	BetUserStart	*uint32    //第一个押注人的Id
	BetUserNow	*uint32	   //当前押注人的Id
}


func (t *ThDesk) LogString(){
	log.T("当前desk[%v]的信息:-----------------------------------begin",t.Id)
	log.T("当前desk[%v]的信息的状态status[%v]",*t.Id,*t.Status)
	log.T("当前desk[%v]的信息的状态users[%v]",*t.Id,t.users)
	log.T("当前desk[%v]的信息的状态SeatedCount[%v]",*t.Id,*t.SeatedCount)

	log.T("当前desk[%v]的信息:-----------------------------------end",t.Id)
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
	for i := 0; i < len(t.users); i++ {
		if t.users[i] == nil {
			t.users[i] = thUser
			break
		}
		if (i + 1) == len(t.users) {
			log.E("玩家加入桌子失败")
			return errors.New("加入房间失败")
		}
	}

	//等待的用户加1
	*t.SeatedCount = (*t.SeatedCount) + 1
	log.T("玩家加入桌子的结果,",t.users)
	return nil
}

/**
	开始游戏
 */
func (t *ThDesk) Run() error {

	//把正在等待的用户设置可以游戏
	err := t.UserWait2Seat()
	if err != nil {
		log.E("开始游戏失败,errMsg[%v]", err.Error())
		return err
	}

	//初始化牌的信息
	err = t.OnInitCards()
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

	//设置房间状态
	*t.Status = TH_DESK_STATUS_SART                	//设置状态为开始游戏
	t.OinitBetUserStar()				//设置第一个押注的人
	return nil
}

/**
	把正在等待的用户安置在座位上
 */
func (t *ThDesk) UserWait2Seat() error {
	//打印移动之前的信息
	log.T("UserWait2Seat 之前,t.users,[%v]",t.users)
	for i := 0; i < len(t.users); i++ {
		if t.users[i] !=nil {
			t.users[i].status = &TH_USER_STATUS_SEATED
		}
	}
	//打印测试消息
	log.T("UserWait2Seat 之后,t.users,[%v]",t.users)
	return nil
}

/**
	初始化纸牌的信息
 */
func (t *ThDesk) OnInitCards() error {
	var total = 21;
	totalCards  := pokerService.RandomTHPorkCards(total)	//得到牌
	log.T("得到的所有手牌:",totalCards)
	//得到所有的牌,前五张为公共牌,后边的每两张为手牌
	t.PublicPai = totalCards[0:5]
	log.T("得到的公共牌:",t.PublicPai)
	log.T("总人数:",len(t.users))

	//给每个人分配手牌
	for i := 0; i < len(t.users); i++ {
		if t.users[i] !=nil && userService.CheckUserIdRightful(*t.users[i].userId)  {
			t.users[i].cards = totalCards[i*2+5:i*2+5+2]
			log.T("用户[%v]的手牌[%v]",t.users[i].userId,t.users[i].cards)
		}
	}

	return nil

}


/**
	德州扑克广播消息
 */
func (t *ThDesk) THBroadcastProto(p proto.Message, ignoreUserId int32) error {
	log.Normal("给每个房间发送proto 消息%v", p)
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil {
			log.Normal("开始userId[%v]发送消息", t.users[i].userId)
			a := t.users[i].agent
			a.WriteMsg(p)
			log.Normal("给userId[%v]发送消息,发送完毕", t.users[i].userId)
		}
	}
	return nil
}


/**
	返回res需要的User实体
 */
func (t *ThDesk) GetResUserModel() []*bbproto.THUser {
	result := make([]*bbproto.THUser,THROOM_SEAT_COUNT)

	//就坐的人
	for i := 0; i< len(t.users);i++  {
		if  t.users[i] != nil {
			result[i] = t.users[i].trans2bbprotoThuser()
		}else{
			result[i] = &bbproto.THUser{}
		}
	}

	log.T("得到的User的情况,",result)
	return result
}

// 返回res需要的User实体 并且排序,排序规则是,当前用户排在第一个
func (t *ThDesk) GetResUserModelClieSeq(userId uint32) []*bbproto.THUser {
	//需要根据当前用户的Userid来进行排序
	users := t.GetResUserModel()
	var userIndex int = 0
	for i := 0; i < len(users); i++ {
		if users[i] !=nil && *(users[i].U.Id) == userId {
			userIndex = i
			break
		}
	}

	result := make([]*bbproto.THUser,len(users))
	for i := 0; i < len(users); i++ {
		result[i] = users[(i+userIndex)%len(users)]
	}

	log.T("得到排序后的User的情况,",result)
	return result
}




// 	初始化第一个押注的人
func (t *ThDesk) OinitBetUserStar() error{
	users := t.users
	for i := 0; i < len(users); i++ {
		if users[i] !=nil {
			t.BetUserStart = users[i].userId
			t.BetUserNow   = users[i].userId
		}
	}

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
	result.Id = new(uint32)
	result.SeatedCount = new(int32)
	result.Dealer = new(uint32)
	result.Status = &TH_DESK_STATUS_STOP
	result.BetUserNow = new(uint32)
	result.BetUserStart = new(uint32)
	result.users = make([]*ThUser,THROOM_SEAT_COUNT)
	return result
}