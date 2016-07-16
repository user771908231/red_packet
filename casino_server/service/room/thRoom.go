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
var TH_USER_STATUS_WAITSEAT int32 = 1        //刚上桌子 等待开始的玩家
var TH_USER_STATUS_SEATED int32 = 2        //刚上桌子 游戏中的玩家
var TH_USER_STATUS_BETING int32 = 3        //押注中
var TH_USER_STATUS_ALLINING int32 = 4        //allIn
var TH_USER_STATUS_FOLDED int32 = 5        //弃牌


//德州扑克,牌桌的状态
var TH_DESK_STATUS_STOP int32 = 1         //没有开始的状态
var TH_DESK_STATUS_SART int32 = 2         //没有已经开始的状态
var TH_DESK_STATUS_ROUND1 int32 = 3         //第一轮押注
var TH_DESK_STATUS_ROUND2 int32 = 4         //第二轮押注
var TH_DESK_STATUS_ROUND3 int32 = 5         //第三轮押注
var TH_DESK_STATUS_ROUND4 int32 = 6         //第四轮押注
var TH_DESK_STATUS_ROUND5 int32 = 7         //第五轮押注


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
	r.ThDeskBuf[index] = throom
	return nil
}

//通过Id找到对应的桌子
func (r *ThGameRoom) GetDeskById(id uint32) *ThDesk {
	var result *ThDesk = nil
	for i := 0; i < len(r.ThDeskBuf); i++ {
		if r.ThDeskBuf[i] != nil && *r.ThDeskBuf[i].Id == id{
			result = r.ThDeskBuf[i]
			break
		}
	}
	return result
}


/**
	正在玩德州的人
 */
type ThUser struct {
	userId  *uint32               //用户id
	agent   gate.Agent            //agent
	status  *int32                //当前的状态
	cards   []*bbproto.Pai        //手牌
	thCards *pokerService.ThCards //手牌加公共牌取出来的值,这个值可以实在结算的时候来取
}

//
func (t *ThUser) trans2bbprotoThuser() *bbproto.THUser {
	thuserTemp := &bbproto.THUser{}
	thuserTemp.Status = t.status        //已经就做
	thuserTemp.User = userService.GetUserById(*t.userId)        //得到user
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
	sync.Mutex
	Id           *uint32        //roomid
	Dealer       *uint32        //庄家
	PublicPai    []*bbproto.Pai //公共牌的部分
	SeatedCount  *int32         //已经坐下的人数
	users        []*ThUser      //坐下的人
	Status       *int32         //牌桌的状态
	BigBlind     *uint32        //第一个押注人的Id
	SmallBlind   *uint32        //第一个押注人的Id
	BetUserNow   *uint32        //当前押注人的Id
	RemainTime   *int32         //剩余投资的时间  多少秒
	BetAmountNow *int32         //挡墙的押注金额是多少
}

func (t *ThDesk) LogString() {
	log.T("当前desk[%v]的信息:-----------------------------------begin", t.Id)
	log.T("当前desk[%v]的信息的状态status[%v]", *t.Id, *t.Status)
	log.T("当前desk[%v]的信息的状态users[%v]", *t.Id, t.users)
	log.T("当前desk[%v]的信息的状态SeatedCount[%v]", *t.Id, *t.SeatedCount)
	log.T("当前desk[%v]的信息:-----------------------------------end", t.Id)
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
	log.T("玩家加入桌子的结果,", t.users)
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
			t.users[i].status = &TH_USER_STATUS_SEATED
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
	result := make([]*bbproto.THUser, THROOM_SEAT_COUNT)

	//就坐的人
	for i := 0; i < len(t.users); i++ {
		if t.users[i] != nil {
			result[i] = t.users[i].trans2bbprotoThuser()
		} else {
			result[i] = &bbproto.THUser{}
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
		if users[i] != nil && *(users[i].User.Id) == userId {
			userIndex = i
			break
		}
	}

	result := make([]*bbproto.THUser, len(users))
	for i := 0; i < len(users); i++ {
		result[i] = users[(i + userIndex) % len(users)]
	}

	log.T("得到排序后的User的情况,", result)
	return result
}




// 	初始化第一个押注的人,当前押注的人
func (t *ThDesk) OinitBegin() error {

	//设置德州desk状态//设置状态为开始游戏
	*t.Status = TH_DESK_STATUS_SART

	userTemp := make([]*ThUser ,len(t.users))
	copy(userTemp,t.users)
	//这里需要定义一个庄家,todo 暂时默认为第一个,后边再修改
	var dealerIndex int = 0;
	if t.Dealer == nil {
		t.Dealer = t.users[0].userId
	} else {
		dealerIndex = t.GetUserIndex(*t.Dealer)
		dealerIndex ++
	}


	//设置小盲注
	for i := dealerIndex; i < len(userTemp)+dealerIndex; i++ {
		if userTemp[i%len(userTemp)] != nil {
			t.SmallBlind = userTemp[i%len(userTemp)].userId
			userTemp[i%len(userTemp)] = nil
			break
		}
	}

	//设置大盲注
	for i := dealerIndex; i < len(userTemp)+dealerIndex; i++ {
		if userTemp[i%len(userTemp)] != nil {
			t.BigBlind = userTemp[i%len(userTemp)].userId
			userTemp[i%len(userTemp)] = nil
			break
		}
	}

	if t.SeatedCount == 2 {	//如果只有两个人,当前押注的人是小盲注
		t.BetUserNow = t.SmallBlind
	}else{
		//设置当前押注的人
		for i := dealerIndex; i < len(userTemp)+dealerIndex; i++ {
			if userTemp[i%len(userTemp)] != nil {
				t.BetUserNow = userTemp[i%len(userTemp)].userId
				userTemp[i%len(userTemp)] = nil
				break
			}
		}
	}


	log.T("初始化游戏之后,庄家[%v]",*t.Dealer)
	log.T("初始化游戏之后,大盲注[%v]",*t.BigBlind)
	log.T("初始化游戏之后,小盲注[%v]",*t.SmallBlind)
	log.T("初始化游戏之后,当前押注Id[%v]",*t.BetUserNow)
	return nil
}


//	开奖
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
		t.BetUserCall(userId,m.GetBetAmount())

	case TH_DESK_BET_TYPE_FOLD:
		t.BetUserFold(userId)

	case TH_DESK_BET_TYPE_CHECK:	//让牌
		t.BetUserCheck(userId)

	case TH_DESK_BET_TYPE_RAISE:	//加注
		t.BetUserRaise(userId)

	case TH_DESK_BET_TYPE_RERRAISE:	//再加注

	case TH_DESK_BET_TYPE_ALLIN:	//全部
		t.BetUserAllIn(userId)
	}


	//3,处理之后,设置desk的状态
	userIndex := t.GetUserIndex(userId)
	for i := userIndex; i < userIndex + len(t.users); i++ {
		nextUser := t.users[(i + 1) % len(t.users)]
		if nextUser != nil && *nextUser.status != TH_USER_STATUS_BETING {
			//用户不为nil ,并且状态是押注中的才可以押注
			t.BetUserNow = nextUser.userId
			break
		}
	}


	//4,押注完成之后返回信息
	result := &bbproto.THBetBroadcast{}
	result.Header = protoUtils.GetSuccHeader()
	result.BetType = m.BetType
	result.BetAmount = m.BetAmount
	result.User = t.GetResUserModelById(m.GetHeader().GetUserId())
	result.NextBetUserId = m.GetHeader().UserId

	//广播给每个玩家
	t.THBroadcastProto(result, userId)

	//给押注的玩家返回押注结果
	betResult := &bbproto.THBet{}
	betResult.Header = protoUtils.GetSuccHeaderwithUserid(m.GetHeader().UserId)
	betResult.BetAmount = m.BetAmount
	betResult.BetType = m.BetType
	a.WriteMsg(betResult)

	return nil
}


//跟注:跟注的时候 不需要重新设置押注的人
/**
	只是跟注,需要减少用户的资产,增加奖池的金额
 */
func (t *ThDesk) BetUserCall(userId uint32,coin int32) error{

	//1,增加奖池的金额

	//2,减少用户的金额
	userService.DecreaseUserCoin(userId,coin)

	//3,修改下次押注的人

	return nil
}

//用户弃牌
func (t *ThDesk) BetUserFold(userId uint32) error{
	return nil
}


func (t *ThDesk) BetUserCheck(userId uint32) error{
	return nil
}

//用户加注
func (t *ThDesk) BetUserRaise(userId uint32) error{
	return nil
}

//用户AllIn
func (t *ThDesk) BetUserAllIn(userId uint32) error{
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
 */
func (t *ThDesk) NextBetUser() error{
	index := t.GetUserIndex(*t.BetUserNow)
	for i := index; i < len(t.users)+index;i++ {
		u :=  t.users[i % len(t.users)]
		if u !=nil && *u.status ==  TH_USER_STATUS_BETING {
			t.BetUserNow = u.userId
			break
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
	//result.Dealer = new(uint32)		//不需要创建  默认就是为空
	result.Status = &TH_DESK_STATUS_STOP
	result.BetUserNow = new(uint32)
	result.BigBlind = new(uint32)
	result.SmallBlind = new(uint32)
	result.users = make([]*ThUser, THROOM_SEAT_COUNT)
	result.RemainTime = new(int32)
	return result
}