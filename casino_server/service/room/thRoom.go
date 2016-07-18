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
	"time"
	"github.com/nu7hatch/gouuid"
)



//config
var THROOM_SEAT_COUNT int32 = 8               //玩德州扑克,每个房间最多多少人
var GAME_THROOM_MAX_COUNT int32 = 500         //一个游戏大厅最多有多少桌德州扑克
var TH_TIMEOUT_DURATION = time.Second*5	      //德州出牌的超时时间

//测试的时候 修改喂多人才可以游戏
var TH_DESK_LEAST_START_USER int32 = 5        //最少多少人可以开始游戏

//德州扑克 玩家的状态
var TH_USER_STATUS_WAITSEAT int32 = 1        //刚上桌子 等待开始的玩家
var TH_USER_STATUS_SEATED int32 = 2        //刚上桌子 游戏中的玩家
var TH_USER_STATUS_BETING int32 = 3        //押注中
var TH_USER_STATUS_ALLINING int32 = 4        //allIn
var TH_USER_STATUS_FOLDED int32 = 5        //弃牌


//德州扑克,牌桌的状态
var TH_DESK_STATUS_STOP int32 = 1         //没有开始的状态
var TH_DESK_STATUS_SART int32 = 2         //没有已经开始的状态

var TH_DESK_ROUND1 int32 = 1         //第一轮押注
var TH_DESK_ROUND2 int32 = 2         //第二轮押注
var TH_DESK_ROUND3 int32 = 3         //第三轮押注
var TH_DESK_ROUND4 int32 = 4         //第四轮押注


//押注的类型
var TH_DESK_BET_TYPE_BET int32 = 1        //押注
var TH_DESK_BET_TYPE_CALL int32 = 2        //跟注,和别人下相同的筹码
var TH_DESK_BET_TYPE_FOLD int32 = 3        //弃牌
var TH_DESK_BET_TYPE_CHECK int32 = 4        //让牌
var TH_DESK_BET_TYPE_RAISE int32 = 5        //加注
var TH_DESK_BET_TYPE_RERRAISE int32 = 6        //再加注
var TH_DESK_BET_TYPE_ALLIN int32 = 7        //全下

/**
	初始化函数:
		初始化游戏房间
 */

var ThGameRoomIns ThGameRoom        //房间实例,在init函数中初始化

func init() {
	ThGameRoomIns.OnInit()                //初始化房间
	ThGameRoomIns.Run()                   //运行房间
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
func (r *ThGameRoom) AddThRoom(index int, throom *ThDesk) error {
	*throom.Id = uint32(index)
	r.ThDeskBuf[index] = throom
	return nil
}

//通过Id找到对应的桌子
func (r *ThGameRoom) GetDeskById(id uint32) *ThDesk {
	var result *ThDesk = nil
	for i := 0; i < len(r.ThDeskBuf); i++ {
		if r.ThDeskBuf[i] != nil && *r.ThDeskBuf[i].Id == id {
			result = r.ThDeskBuf[i]
			break
		}
	}
	return result
}

//通过UserId判断是不是重复进入房间
func (r *ThGameRoom) IsRepeatIntoRoom(userId uint32) bool {

	//废弃的判断的代码
	/**
	//todo 开发过程中,可能需要关掉这个限制
	var  agentUser *gamedata.AgentUserData
	userData := a.UserData()
	if userData == nil {
		log.T("用户进入德州扑克房间的时候,没有知道对应的agentUserdata")
		//todo 需要删掉的代码
		agentUser = &gamedata.AgentUserData{}
		a.SetUserData(agentUser)
	}else {
		agentUser = userData.(*gamedata.AgentUserData)
	}

	if agentUser.Status == gamedata.AGENT_USER_STATUS_GAMING && agentUser.ThDeskId > 0 {
		log.E("用户已经在房间中了,请不要重复进入")
		return errors.New("玩家已经在房间中了,请不要重复进入")
	}
 */

	//新的判断的代码
	result := r.GetDeskByUserId(userId)
	if result == nil {
		return false
	}else{
		return true
	}
}

/**
	通过UserId 找到对应的桌子
 */
func (r *ThGameRoom) GetDeskByUserId(userId uint32) *ThDesk{
	var result *ThDesk
	var breakFlag bool = false
	desks := ThGameRoomIns.ThDeskBuf
	for i := 0; i < len(desks); i++ {
		if breakFlag {
			break
		}
		desk := desks[i]
		if desk != nil {
			users := desk.users
			for j := 0; j < len(users); j++ {
				u := users[j]
				if u != nil && *u.userId == userId {
					result = desk
					breakFlag = true
					break
				}
			}

		}
	}
	return result
}


/**
	正在玩德州的人
 */
type ThUser struct {
	sync.Mutex
	userId   *uint32               //用户id
	agent    gate.Agent            //agent
	status   *int32                //当前的状态
	cards    []*bbproto.Pai        //手牌
	thCards  *pokerService.ThCards //手牌加公共牌取出来的值,这个值可以实在结算的时候来取
	waiTime  time.Time             //等待时间
	waitUUID string                //等待标志
	deskId   *uint32               //用户所在的桌子的编号
}

//
func (t *ThUser) trans2bbprotoThuser() *bbproto.THUser {

	thuserTemp := &bbproto.THUser{}
	thuserTemp.Status = t.status        //已经就做
	thuserTemp.User = userService.GetUserById(*t.userId)        //得到user
	thuserTemp.HandPais = t.cards
	thuserTemp.SeatNumber = new(int32)
	return thuserTemp
}

//等待用户出牌

func (t *ThUser) wait() error{
	ticker := time.NewTicker(time.Second * 1)
	t.waiTime =  time.Now().Add(TH_TIMEOUT_DURATION)
	uuid,_ := uuid.NewV4()
	t.waitUUID = uuid.String()		//设置出牌等待的标志
	go func() {
		for timeNow := range ticker.C {
			//表示已经过期了
			bool,err := t.TimeOut(timeNow)
			if err != nil {
				log.E("处理超时的逻辑出现错误")
				return
			}

			//判断是否已经超时
			if bool {
				log.E("user[%v]已经超市,结束等待任务",*t.userId)
				return
			}
		}
	}()

	return nil

}

//返回自己所在的桌子
func (t *ThUser) GetDesk() *ThDesk{
	desk := ThGameRoomIns.GetDeskById(*t.deskId)
	return desk

}
//用户超市,做处理
func (t *ThUser) TimeOut(timeNow time.Time) (bool,error){
	t.Lock()
	defer t.Unlock()
	if t.waiTime.Before(timeNow) && t.waitUUID != ""{
		log.T("玩家[%v]超时,现在做超时的处理",*t.userId)
		//表示已经超时了
		t.waitUUID = ""
		//给玩家发送超时的广播

		//玩家自动押注弃牌,这里需要模拟数据
		m := &bbproto.THBet{}
		m.Header = protoUtils.GetSuccHeaderwithUserid(t.userId)
		m.BetType = &TH_DESK_BET_TYPE_FOLD

		a := t.agent
		desk := t.GetDesk()
		err := desk.Bet(m,a)
		//给其他玩家广播
		log.T("玩家[%v]超时,现在做超时的处理,处理完毕",*t.userId)
		return true,err
	}else{
		log.T("玩家[%v]出牌中还没有超时",*t.userId)
		return false,nil
	}

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
	sync.Mutex
	Id                 *uint32        //roomid
	Dealer             *uint32        //庄家
	PublicPai          []*bbproto.Pai //公共牌的部分
	SeatedCount        *int32         //已经坐下的人数
	users              []*ThUser      //坐下的人
	Status             *int32         //牌桌的状态
	BigBlind           *uint32        //第一个押注人的Id
	SmallBlind         *uint32        //第一个押注人的Id
	BetUserRaiseUserId *uint32        //加注的人的Id
	BetUserNow         *uint32        //当前押注人的Id
	RemainTime         *int32         //剩余投资的时间  多少秒
	BetAmountNow       *int32         //挡墙的押注金额是多少
	RoundCount         *int32         //低几轮
	CheckUserFirst     *uint32        //第一个人让牌的人
}

func (t *ThDesk) LogString() {
	log.T("当前desk[%v]的信息:-----------------------------------begin----------------------------------", *t.Id)
	log.T("当前desk[%v]的信息的状态status[%v]", *t.Id, *t.Status)
	log.T("当前desk[%v]的信息的状态users[%v]", *t.Id, t.users)
	log.T("当前desk[%v]的信息的状态,总人数SeatedCount[%v]", *t.Id, *t.SeatedCount)
	log.T("当前desk[%v]的信息的状态,小盲注[%v]", *t.Id, *t.SmallBlind)
	log.T("当前desk[%v]的信息的状态,大盲注[%v]", *t.Id, *t.BigBlind)
	log.T("当前desk[%v]的信息的状态,压注人[%v]", *t.Id, *t.BetUserNow)
	log.T("当前desk[%v]的信息:-----------------------------------end----------------------------------", *t.Id)
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
	thUser.deskId = t.Id		//桌子的id

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
	log.T("玩家加入桌子的结果,", t.users)
	return nil
}

//  用户退出德州游戏的房间,rmUser 需要在事物中进行
func (t *ThDesk) RmThuser(userId uint32) error{
	t.Lock()
	defer  t.Unlock()

	//设置为nil就行了
	index := t.GetUserIndex(userId)		//
	t.users[index] = nil			//设置为nil

	//设置在线的人数 减1
	*t.SeatedCount --
	return nil

}


/**
	开始游戏,开始游戏的时候需要初始化desk
	1,初始化庄
	2,初始化
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

	//设置房间状态
	t.OinitBegin()

	//广播消息
	res := &bbproto.THBetBegin{}
	res.Header = protoUtils.GetSuccHeader()
	res.SmallBlind = t.GetResUserModelById(*t.SmallBlind)
	res.BigBlind = t.GetResUserModelById(*t.BigBlind)
	res.BetUserNow =  t.BetUserNow
	err = t.THBroadcastProto(res, 0)
	if err != nil {
		log.E("广播开始消息的时候出错")
		return err
	}

	//设置第一个押注的人
	return nil
}

/**
	把正在等待的用户安置在座位上
 */
func (t *ThDesk) UserWait2Seat() error {
	//打印移动之前的信息
	//log.T("UserWait2Seat 之前,t.users,[%v]",t.users)
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil {
			t.users[i].status = &TH_USER_STATUS_BETING
		}
	}
	//打印测试消息
	//log.T("UserWait2Seat 之后,t.users,[%v]",t.users)
	return nil
}

/**
	初始化纸牌的信息
 */
func (t *ThDesk) OnInitCards() error {
	var total = 21;
	totalCards := pokerService.RandomTHPorkCards(total)        //得到牌
	log.T("得到的所有手牌:", totalCards)
	//得到所有的牌,前五张为公共牌,后边的每两张为手牌
	t.PublicPai = totalCards[0:5]
	log.T("得到的公共牌:", t.PublicPai)
	log.T("总人数:", len(t.users))

	//给每个人分配手牌
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil {
			t.users[i].cards = totalCards[i * 2 + 5:i * 2 + 5 + 2]
			log.T("用户[%v]的手牌[%v]", t.users[i].userId, t.users[i].cards)
			t.users[i].thCards = pokerService.GetTHMax(t.users[i].cards, t.PublicPai, 5)
		}
	}

	return nil

}


/**
	德州扑克广播消息
 */
func (t *ThDesk) THBroadcastProto(p proto.Message, ignoreUserId uint32) error {
	log.Normal("给每个房间发送proto 消息%v", p)
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil && *t.users[i].userId != ignoreUserId {
			log.Normal("开始userId[%v]发送消息", *t.users[i].userId)
			a := t.users[i].agent
			a.WriteMsg(p)
			log.Normal("给userId[%v]发送消息,发送完毕", *t.users[i].userId)
		}
	}
	return nil
}

/**
当有新用户进入房间的时候,为其他人广播新过来的人的信息
 */
func (t *ThDesk) THBroadcastAddUser(newUserId, ignoreUserId uint32) error {
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil && *t.users[i].userId != ignoreUserId {
			users := t.GetResUserModelClieSeq(*t.users[i].userId)
			broadUsers := &bbproto.THRoomAddUserBroadcast{}
			broadUsers.Header = protoUtils.GetSuccHeaderwithUserid(t.users[i].userId)
			for i := 0; i < len(users); i++ {
				if users[i] != nil && users[i].User.GetId() == newUserId {
					broadUsers.User = users[i]
					break
				}
			}

			a := t.users[i].agent
			log.Normal("给userId[%v]发送消息:[%v]", *t.users[i].userId, broadUsers)
			a.WriteMsg(broadUsers)
		}
	}
	return nil

}


/**
	返回res需要的User实体
 */
func (t *ThDesk) GetResUserModel() []*bbproto.THUser {
	result := make([]*bbproto.THUser, THROOM_SEAT_COUNT)

	//就坐的人
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil {
			result[i] = t.users[i].trans2bbprotoThuser()
		} else {
			result[i] = &bbproto.THUser{}
			result[i].SeatNumber = new(int32)        //设置为0
		}
	}

	log.T("得到的User的情况,", result)
	return result
}

func (t *ThDesk) GetResUserModelById(userId uint32) *bbproto.THUser {
	var result *bbproto.THUser
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil &&  *t.users[i].userId == userId {
			result = t.users[i].trans2bbprotoThuser()
		}
	}
	log.T("通过userId得到的bbproto.THUser的情况,", result)
	return result
}

// 返回res需要的User实体 并且排序,排序规则是,当前用户排在第一个
func (t *ThDesk) GetResUserModelClieSeq(userId uint32) []*bbproto.THUser {
	//需要根据当前用户的Userid来进行排序
	users := t.GetResUserModel()
	var userIndex int = 0
	for i := 0; i < len(users); i++ {
		if users[i] != nil && users[i].User != nil && *(users[i].User.Id) == userId {
			userIndex = i
			break
		}
	}

	for i := 0; i < len(users); i++ {
		u := users[(i + userIndex) % len(users)]
		if u != nil {
			*u.SeatNumber = int32(i)
		}
	}

	log.T("得到排序后的User的情况,", users)
	return users
}




// 	初始化第一个押注的人,当前押注的人
func (t *ThDesk) OinitBegin() error {

	//设置德州desk状态//设置状态为开始游戏
	*t.Status = TH_DESK_STATUS_SART

	userTemp := make([]*ThUser, len(t.users))
	copy(userTemp, t.users)
	//这里需要定义一个庄家,todo 暂时默认为第一个,后边再修改
	var dealerIndex int = 0;
	if t.Dealer == nil {
		log.T("游戏没有庄家,现在默认初始化第一个人为庄家")
		t.Dealer = t.users[0].userId
	} else {
		dealerIndex = t.GetUserIndex(*t.Dealer)
		dealerIndex ++
	}


	//设置小盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil {
			t.SmallBlind = u.userId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	//设置大盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil {
			t.BigBlind = u.userId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	if *t.SeatedCount == int32(2) {
		//如果只有两个人,当前押注的人是小盲注
		t.BetUserNow = t.SmallBlind
	} else {
		//设置当前押注的人
		for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
			u := userTemp[(i + 1) % len(userTemp)]
			if u != nil {
				t.BetUserNow = u.userId
				userTemp[(i + 1) % len(userTemp)] = nil
				break
			}
		}
	}

	t.BetUserRaiseUserId = t.BetUserNow        //第一个加注的人
	t.RoundCount = &TH_DESK_ROUND1

	//本次押注的热开始等待
	waitUser := t.users[t.GetUserIndex(*t.BetUserNow)]
	waitUser.wait()

	log.T("初始化游戏之后,庄家[%v]", *t.Dealer)
	log.T("初始化游戏之后,小盲注[%v]", *t.SmallBlind)
	log.T("初始化游戏之后,大盲注[%v]", *t.BigBlind)
	log.T("初始化游戏之后,当前押注Id[%v]", *t.BetUserNow)
	log.T("初始化游戏之后,第一个加注的人Id[%v]", *t.BetUserRaiseUserId)
	log.T("初始化游戏之后,当前轮数Id[%v]", *t.RoundCount)
	return nil
}



//判断是否是开奖的时刻
/**
开奖的时候
1,目前已经是第四轮
2,计算出来的下一个押注这和当前押注的是同一个人
3,即使没有到第四轮,但是所有人都all in 了,那么还是要开奖
 */
func (t *ThDesk) Tiem2Lottery() bool {

	//第四轮,并且计算出来的押注人和start是同一个人
	if *t.RoundCount == TH_DESK_ROUND4  && *t.BetUserNow == *t.BetUserRaiseUserId {
		return true
	}

	//如果只有一个人没有all in  或者全部都all in 了也要开牌

	return false
}
//开奖
/**
	开奖的规则
 */
func (t *ThDesk) Lottery() error {

	return nil

}



/**
	押注,押注其实也是桌子在负责
	押注的逻辑说明:
	1,userId必须是当前的UserId
	2,当加注的时候,庄变成加注的那个人
	3,跟注到庄的时候一轮结束
 */
func (t *ThDesk) Bet(m *bbproto.THBet, a gate.Agent) error {
	t.Lock()
	t.Unlock()

	//1,检测押注的参数是否正确
	//1.1 userId是否正确
	userId := m.GetHeader().GetUserId()
	if userId != *t.BetUserNow {
		//如果押注的不是当前用户,则直接返回错误
		log.E("当前应该押注的用户是[%v]二不是[%v]", *t.BetUserNow, userId)
		return errors.New("押注的用户Id不正确")
	}

	//1.2 是否轮到当前押注
	if *t.BetUserNow != userId {
		log.E("还没有轮到当前用户押注")
		return errors.New("用户不合法")
	}


	//2,根据押注的类型来分别处理
	betType := m.GetBetType()
	switch betType {
	case TH_DESK_BET_TYPE_BET:
	//押注:大盲注之后的第一个人,这个类型不会使用到

	case TH_DESK_BET_TYPE_CALL:
		t.BetUserCall(userId, m.GetBetAmount())

	case TH_DESK_BET_TYPE_FOLD:
		t.BetUserFold(userId)

	case TH_DESK_BET_TYPE_CHECK:        //让牌
		t.BetUserCheck(userId)

	case TH_DESK_BET_TYPE_RAISE:        //加注
		t.BetUserRaise(userId)

	case TH_DESK_BET_TYPE_RERRAISE:        //再加注

	case TH_DESK_BET_TYPE_ALLIN:        //全部
		t.BetUserAllIn(userId)
	}

	//押注的人移动到下一个人
	t.NextBetUser()
	//押注之后,desk的状态
	t.LogString()

	//是否可以开奖
	/**

	 */

	if t.Tiem2Lottery() {
		return t.Lottery()
	} else {
		//继续押注

		//4,押注完成之后返回信息,广播给其他玩家的信息玩家
		result := &bbproto.THBetBroadcast{}
		result.Header = protoUtils.GetSuccHeader()
		result.BetType = m.BetType
		result.BetAmount = m.BetAmount
		result.User = t.GetResUserModelById(m.GetHeader().GetUserId())
		result.NextBetUserId = t.BetUserNow

		t.THBroadcastProto(result, userId)

		//给押注的玩家返回押注结果
		betResult := &bbproto.THBet{}
		betResult.Header = protoUtils.GetSuccHeaderwithUserid(m.GetHeader().UserId)
		betResult.BetAmount = m.BetAmount
		betResult.BetType = m.BetType
		betResult.NextBetUser = t.BetUserNow
		a.WriteMsg(betResult)
	}

	return nil
}


//跟注:跟注的时候 不需要重新设置押注的人
/**
	只是跟注,需要减少用户的资产,增加奖池的金额
 */
func (t *ThDesk) BetUserCall(userId uint32, coin int32) error {

	//1,增加奖池的金额

	//2,减少用户的金额
	userService.DecreaseUserCoin(userId, coin)

	//3,修改下次押注的人

	return nil
}

//用户弃牌
func (t *ThDesk) BetUserFold(userId uint32) error {
	index := t.GetUserIndex(userId)
	t.users[index].status = &TH_USER_STATUS_FOLDED
	return nil
}

//让牌:只有第一个人才可以让牌
func (t *ThDesk) BetUserCheck(userId uint32) error {
	if userId == *t.BetUserRaiseUserId {
		//第一个人的时候才可以让牌

		//设置喂第一个让牌的人
		if t.CheckUserFirst == nil {
			t.CheckUserFirst = new(uint32)
			*t.CheckUserFirst = userId
		}

		//设置一个押注的人为下一个人
		index := t.GetUserIndex(userId)
		for i := index; i < len(t.users) + index - 1; i++ {
			u := t.users[(i + 1) % len(t.users)]
			if u != nil && *u.status == TH_USER_STATUS_BETING && *u.userId != *t.CheckUserFirst {
				*t.BetUserRaiseUserId = userId
				break
			}
		}

	}
	return nil
}

//用户加注
func (t *ThDesk) BetUserRaise(userId uint32) error {
	return nil
}

//用户AllIn
func (t *ThDesk) BetUserAllIn(userId uint32) error {
	return nil
}

/**
	根据userId 找到在桌子上的index
 */
func (t *ThDesk) GetUserIndex(userId uint32) int {
	var result int = 0
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil && *t.users[i].userId == userId {
			result = i
			break
		}
	}

	return result
}



/**
	初始化下一个押注的人
	初始化下一个人的时候需要一个超时的处理
 */
func (t *ThDesk) NextBetUser() error {


	log.T("当前押注的人是userId[%v]", *t.BetUserNow)
	index := t.GetUserIndex(*t.BetUserNow)
	for i := index; i < len(t.users) + index; i++ {
		u := t.users[(i + 1) % len(t.users)]
		log.T("设置下次押注人 i[%v]", i)

		if u != nil && *u.status == TH_USER_STATUS_BETING {
			log.T("设置betUserNow 为[%v]", *u.userId)
			t.BetUserNow = u.userId
			break
		}
	}


	//如果不是第四局,并且下次押注的人指向同一个人,那么设置下次押注的人是小盲注
	if *t.BetUserRaiseUserId == *t.BetUserNow && *t.RoundCount != TH_DESK_ROUND4 {
		t.BetUserRaiseUserId = t.SmallBlind
		t.BetUserNow = t.SmallBlind
		*t.RoundCount ++
		log.T("设置下次押注的人是小盲注,下轮次[%v]", *t.RoundCount)
	}


	//用户开始等待,如果超时,需要做超时的处理
	waitUser := t.users[t.GetUserIndex(*t.BetUserNow)]
	waitUser.wait()
	log.T("重新计算出来的押注userId是[%v]", *t.BetUserNow)
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
	//result.Dealer = new(uint32)		//不需要创建  默认就是为空
	result.Status = &TH_DESK_STATUS_STOP
	result.BetUserNow = new(uint32)
	result.BigBlind = new(uint32)
	result.SmallBlind = new(uint32)
	result.users = make([]*ThUser, THROOM_SEAT_COUNT)
	result.RemainTime = new(int32)
	result.BetUserRaiseUserId = new(uint32)
	result.RoundCount = new(int32)
	return result
}