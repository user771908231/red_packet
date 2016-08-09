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
	"sort"
	"casino_server/conf/casinoConf"
	"github.com/name5566/leaf/db/mongodb"
	"fmt"
	"casino_server/mode"
	"gopkg.in/mgo.v2/bson"
	"casino_server/gamedata"
)
//config

var TH_GAME_SMALL_BLIND int64 = 10                //小盲注的金额
var GAME_THROOM_MAX_COUNT int32 = 500             //一个游戏大厅最多有多少桌德州扑克
var TH_TIMEOUT_DURATION = time.Second * 1500      //德州出牌的超时时间
var TH_TIMEOUT_DURATION_INT int32 = 1500       	  //德州出牌的超时时间
var TH_LOTTERY_DURATION = time.Second * 5         //德州开奖的时间
var TH_DESK_CREATE_DIAMOND int64 = 10;		  //创建牌桌需要的钻石数量


//测试的时候 修改喂多人才可以游戏
var TH_DESK_LEAST_START_USER int32 = 2       //最少多少人可以开始游戏
var TH_DESK_MAX_START_USER int32 = 8         //玩德州扑克,每个房间最多多少人

//德州扑克 玩家的状态
var TH_USER_STATUS_WAITSEAT int32 = 1        //刚上桌子 等待开始的玩家
var TH_USER_STATUS_SEATED int32 = 2          //刚上桌子 游戏中的玩家
var TH_USER_STATUS_BETING int32 = 3          //押注中
var TH_USER_STATUS_ALLINING int32 = 4        //allIn
var TH_USER_STATUS_FOLDED int32 = 5          //弃牌
var TH_USER_STATUS_WAIT_CLOSED int32 = 6     //等待结算
var TH_USER_STATUS_CLOSED int32 = 7          //已经结算
var TH_USER_STATUS_LEAVE int32 = 8           //

var TH_USER_BREAK_STATUS_TRUE int32 = 1      //已经断线
var TH_USER_BREAK_STATUS_FALSE int32 = 0     //没有断线


//德州扑克,牌桌的状态
var TH_DESK_STATUS_STOP int32 = 1            //没有开始的状态
var TH_DESK_STATUS_SART int32 = 2            //已经开始的状态
var TH_DESK_STATUS_LOTTERY int32 = 3         //已经开始的状态

var TH_DESK_ROUND1 int32 = 1         //第一轮押注
var TH_DESK_ROUND2 int32 = 2         //第二轮押注
var TH_DESK_ROUND3 int32 = 3         //第三轮押注
var TH_DESK_ROUND4 int32 = 4         //第四轮押注
var TH_DESK_ROUND_END int32 = 5      //完成押注


//押注的类型
var TH_DESK_BET_TYPE_BET int32 = 1        //押注
var TH_DESK_BET_TYPE_CALL int32 = 2       //跟注,和别人下相同的筹码
var TH_DESK_BET_TYPE_FOLD int32 = 3       //弃牌
var TH_DESK_BET_TYPE_CHECK int32 = 4      //让牌
var TH_DESK_BET_TYPE_RAISE int32 = 5      //加注
var TH_DESK_BET_TYPE_RERRAISE int32 = 6   //再加注
var TH_DESK_BET_TYPE_ALLIN int32 = 7      //全下

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
	RoomStatus     int32 //游戏大厅的状态
	ThDeskBuf      []*ThDesk
	ThRoomSeatMax  int32 //每个房间的座位数目
	ThRoomCount    int32 //房间数目
	Id             int32 //房间的id
	SmallBlindCoin int64 //小盲注的金额
}


//初始化游戏房间
func (r *ThGameRoom) OnInit() {
	r.ThRoomSeatMax = TH_DESK_MAX_START_USER
	r.Id = 0
	r.SmallBlindCoin = TH_GAME_SMALL_BLIND;
}

//run游戏房间
func (r *ThGameRoom) Run() {

}


//判断roomkey是否已经存在了
func (r *ThGameRoom) IsRoomKeyExist(roomkey string) bool {
	ret := false
	for i := 0; i < len(r.ThDeskBuf); i++ {
		d := r.ThDeskBuf[i]
		if d != nil && d.RoomKey == roomkey {
			ret = true
			break
		}
	}
	return ret
}

//创建一个房间
func (r *ThGameRoom) CreateDeskByUserIdAndRoomKey(userId uint32, roomCoin int64, roomkey string, smallBlind int64, bigBlind int64, jucount int32) *ThDesk {

	//1,创建房间
	desk := NewThDesk()
	desk.RoomKey = roomkey
	desk.InitRoomCoin = roomCoin
	desk.deskOwner = userId
	desk.SmallBlindCoin = smallBlind
	desk.BigBlindCoin = bigBlind
	desk.JuCount = jucount
	desk.GetRoomCoin()
	r.AddThRoom(desk)

	//2,创建房间成功之后,扣除user的钻石
	upDianmond := 0-TH_DESK_CREATE_DIAMOND
	remainDiamond := userService.UpdateUserDiamond(userId,upDianmond)
	//3,生成一条交易记录
	err := userService.CreateDiamonDetail(userId,mode.T_USER_DIAMOND_DETAILS_TYPE_CREATEDESK,upDianmond,remainDiamond,"创建房间消耗钻石");
	if err != nil {
		log.E("创建用户的钻石交易记录失败")
		return err
	}
	return desk
}

//增加一个thRoom
func (r *ThGameRoom) AddThRoom(throom *ThDesk) error {
	r.ThDeskBuf = append(r.ThDeskBuf, throom)
	return nil
}


//删除一个throom
func (r *ThGameRoom) RmThroom(number int32) error {

	//第一步找到index
	var index int = -1
	for i := 0; i < len(r.ThDeskBuf); i++ {
		desk := r.ThDeskBuf[i]
		if desk != nil && desk.Number == number {
			index = i
			break
		}
	}

	//判断是否找到对应的desk
	if index == -1 {
		log.E("没有找到对应desk.number[%v]的桌子", number)
		return errors.New("没有找到对应的desk")
	}

	//删除对应的desk
	r.ThDeskBuf = append(r.ThDeskBuf[:index], r.ThDeskBuf[index + 1:]...)
	return nil
}

//通过Id找到对应的桌子
func (r *ThGameRoom) GetDeskById(id int32) *ThDesk {
	var result *ThDesk = nil
	for i := 0; i < len(r.ThDeskBuf); i++ {
		if r.ThDeskBuf[i] != nil && r.ThDeskBuf[i].Id == id {
			result = r.ThDeskBuf[i]
			break
		}
	}
	return result
}

//通过UserId判断是不是重复进入房间
func (r *ThGameRoom) IsRepeatIntoRoom(userId uint32, a gate.Agent) *ThDesk {
	//新的判断的代码
	desk := r.GetDeskByUserId(userId)
	if desk != nil {
		log.T("用户[%v]重新进入房间了", userId)
		u := desk.GetUserByUserId(userId)
		//替换User的agent
		u.agent = a
		u.BreakStatus = TH_USER_BREAK_STATUS_FALSE        //设置没有掉线的情况
		desk.UserCountOnline ++
		//绑定参数
		userAgentData := &gamedata.AgentUserData{}
		userAgentData.UserId = userId
		userAgentData.ThDeskId = desk.Id
		a.SetUserData(userAgentData)
	}
	return desk
}

/**
	通过UserId 找到对应的桌子
 */
func (r *ThGameRoom) GetDeskByUserId(userId uint32) *ThDesk {
	var result *ThDesk
	var breakFlag bool = false
	desks := ThGameRoomIns.ThDeskBuf
	for i := 0; i < len(desks); i++ {
		if breakFlag {
			break
		}
		desk := desks[i]
		if desk != nil {
			users := desk.Users
			for j := 0; j < len(users); j++ {
				u := users[j]

				//查找房间,并且,用户离开的房间是不算的
				if u != nil && u.UserId == userId && u.Status != TH_USER_STATUS_LEAVE {
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
	通过roomKey 找到desk
 */
func (r *ThGameRoom) GetDeskByRoomKey(roomKey string) *ThDesk {
	var result *ThDesk
	desks := ThGameRoomIns.ThDeskBuf
	for i := 0; i < len(desks); i++ {
		desk := desks[i]
		if desk != nil && desk.RoomKey == roomKey {
			result = desk
			break
		}
	}
	return result
}


/**
	给指定的房间增加用户
 */
func (r *ThGameRoom) AddUserWithRoomKey(userId uint32, roomCoin int64, roomKey string, a gate.Agent) (*ThDesk, error) {
	//1,首先判断roomKey 是否喂空
	if roomKey == "" {
		return nil, errors.New("房间密码不应该为空")
	}

	//2,如果roomKey 不是为""
	mydesk := r.GetDeskByRoomKey(roomKey)
	if mydesk == nil {
		return nil, errors.New("没有找到对应的房间")
	}

	//3,判断用户是否是掉线重连
	isRepeat := mydesk.IsrepeatIntoWithRoomKey(userId, a)
	if isRepeat {
		return mydesk, nil
	}

	//4,进入房间
	err := mydesk.AddThUser(userId, roomCoin, a)
	if err != nil {
		log.E("用户上德州扑克的桌子 失败...")
		return nil, err
	}

	mydesk.LogString()        //答应当前房间的信息

	return mydesk, nil

}
//游戏大厅增加一个玩家
func (r *ThGameRoom) AddUser(userId uint32, roomCoin int64, a gate.Agent) (*ThDesk, error) {
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
		r.AddThRoom(mydesk)
	}

	//3,进入房间
	err := mydesk.AddThUser(userId, roomCoin, a)
	if err != nil {
		log.E("用户上德州扑克的桌子 失败...")
		return nil, err
	}

	mydesk.LogString()        //答应当前房间的信息

	return mydesk, nil
}

//退出房间,设置房间状态
func (r *ThGameRoom) LeaveRoom(userId uint32) error {
	desk := r.GetDeskByUserId(userId)
	desk.LeaveThuser(userId)
	return nil
}


//随机返回一个编号
func (t *ThGameRoom) RandDeskNumber() int32 {
	//一般来说,桌子的编号是有用户来创建,如果用户没有创建,则系统帮忙创建一个编号
	return int32(999888)        //测试代码,返回一个测试用编号
}

/**
	正在玩德州的人
 */
type ThUser struct {
	sync.Mutex
	UserId          uint32                //用户id
	Seat            int32                 //用户的座位号
	agent           gate.Agent            //agent
	Status          int32                 //当前的状态
	BreakStatus     int32                 //离线前的状态,当重新连接的时候,需要回复之前的状态
	Cards           []*bbproto.Pai        //手牌
	thCards         *pokerService.ThCards //手牌加公共牌取出来的值,这个值可以实在结算的时候来取
	waiTime         time.Time             //等待时间
	waitUUID        string                //等待标志
	deskId          int32                 //用户所在的桌子的编号
	TotalBet        int64                 //押注总额 todo 注意,目前这个值是用来计算all in 的
	winAmount       int64                 //总共赢了多少钱
	winAmountDetail []int64               //赢钱的详细
	TurnCoin        int64                 //单轮押注(总共四轮)的金额
	HandCoin        int64                 //用户下注多少钱、指单局
					      //Coin            int64                 //用户余额多少钱,总余额  ///*****暂时不用这个字段,用户的余额都从redis取出来
	RoomCoin        int64                 //用户上分的金额
	NickName        string                //用户昵称
}

func (t *ThUser) GetCoin() int64 {
	redu := userService.GetUserById(t.UserId)
	if redu == nil {
		return -1
	} else {
		return redu.GetCoin()
	}
}

func (t *ThUser) GetRoomCoin() int64 {
	return t.RoomCoin
}

//
func (t *ThUser) trans2bbprotoThuser() *bbproto.THUser {

	thuserTemp := &bbproto.THUser{}
	thuserTemp.Status = &(t.Status)        //已经就做
	thuserTemp.User = userService.GetUserById(t.UserId)        //得到user
	thuserTemp.HandPais = t.Cards
	thuserTemp.SeatNumber = new(int32)
	return thuserTemp
}

//等待用户出牌
func (t *ThUser) wait() error {
	//如果不是押注中的状态,不用wait任务
	if t.Status != TH_USER_STATUS_BETING {
		return nil
	}

	ticker := time.NewTicker(time.Second * 1)
	t.waiTime = time.Now().Add(TH_TIMEOUT_DURATION)
	uuid, _ := uuid.NewV4()
	t.waitUUID = uuid.String()                //设置出牌等待的标志
	go func() {
		for timeNow := range ticker.C {
			//表示已经过期了
			bool, err := t.TimeOut(timeNow)
			if err != nil {
				log.E("处理超时的逻辑出现错误,errMsg[%v]", err.Error())
				return
			}

			//判断是否已经超时
			if bool {
				log.E("user[%v]已经超时,结束等待任务", t.UserId)
				return
			}
		}
	}()

	return nil

}

//返回自己所在的桌子
func (t *ThUser) GetDesk() *ThDesk {
	desk := ThGameRoomIns.GetDeskById(t.deskId)
	return desk
}

//用户超时,做处理
func (t *ThUser) TimeOut(timeNow time.Time) (bool, error) {
	t.Lock()
	defer t.Unlock()

	//没有超时标志,直接返回
	if t.waitUUID == "" {
		//不需要等待
		return true, nil
	}

	//如果用户超市,或者用户选择离线,那么直接做弃牌的操作
	if t.waiTime.Before(timeNow) || t.Status == TH_USER_STATUS_LEAVE {
		log.T("玩家[%v]超时,现在做超时的处理", t.UserId)
		//表示已经超时了
		//给玩家发送超时的广播
		err := t.GetDesk().OGBet(t.Seat, TH_DESK_BET_TYPE_FOLD, 0)
		if err != nil {
			log.E("用户[%v]弃牌失败", t.UserId)
		}
		//这里需要设置为弃牌的状态
		log.T("玩家[%v]超时,现在做超时的处理,处理完毕", t.UserId)
		return true, err
	} else {
		//没有超时,继续等待
		log.T("玩家[%v]出牌中还没有超时", t.UserId)
		return false, nil
	}
}

func (t *ThUser) InitWait() {
	t.waitUUID = ""
}

//判断用户是否正在等待出牌
func (t *ThUser) IsWaiting() bool {
	return t.waitUUID != ""
}


//操作押注时的waiting 状态
func (t *ThUser) CheckBetWaitStatus() error {
	if t.IsWaiting() {
		t.InitWait()
		return nil
	} else {
		return errors.New("用户状态错误")
	}
}

func NewThUser() *ThUser {
	result := &ThUser{}
	result.UserId = 0
	result.Status = 0
	result.TurnCoin = 0
	result.TotalBet = 0
	result.winAmount = 0
	result.RoomCoin = 0
	return result
}

/**
	一个德州扑克的房间
 */
type ThDesk struct {
	AgentMap           map[uint32]gate.Agent
	sync.Mutex
	Id                 int32                        //roomid
	Number             int32                        //桌子的编号
	DeskType           int32                        //桌子的类型,1,表示自定义房间,2表示锦标赛的
	Dealer             uint32                       //庄家
	PublicPai          []*bbproto.Pai               //公共牌的部分
	UserCount          int32                        //玩游戏的总人数
	UserCountOnline    int32                        //在先人数
	Users              []*ThUser                    //坐下的人
	Status             int32                        //牌桌的状态
	BigBlindCoin       int64                        //大盲注的押注金额
	SmallBlindCoin     int64                        //小盲注的押注金额
	BigBlind           uint32                       //大盲注
	SmallBlind         uint32                       //小盲注
	NewRoundBetUser    uint32                       //新一轮,开始押注的第一个人//第一轮默认是小盲注,但是当小盲注弃牌之后,这个人要滑倒下一家去
	BetUserRaiseUserId uint32                       //加注的人的Id
	BetUserNow         uint32                       //当前押注人的Id
	RemainTime         int32                        //剩余投资的时间  多少秒
	BetAmountNow       int64                        //当前的押注金额是多少
	RoundCount         int32                        //低几轮
	CheckUserFirst     uint32                       //第一个人让牌的人
	Jackpot            int64                        //奖金池
	bianJackpot        int64                        //边池
	AllInJackpot       []*pokerService.AllInJackpot //allin的标记
	MinRaise           int64                        //最低加注金额
	CanRaise           int32                        //是否能加注
	deskOwner          uint32                       //房主的id
	RoomKey            string                       //room 自定义房间的钥匙
	InitRoomCoin       int64			//进入这个房间的roomCoin 带入金额标准是多少
	JuCount            int32                        //这个桌子最多能打多少局
}

/**
	新生成一个德州的桌子
 */
func NewThDesk() *ThDesk {
	result := new(ThDesk)
	result.AgentMap = make(map[uint32]gate.Agent)
	result.Id = newThDeskId()
	result.UserCount = 0
	result.UserCountOnline = 0
	result.Dealer = 0                //不需要创建  默认就是为空
	result.Status = TH_DESK_STATUS_STOP
	result.BetUserNow = 0
	result.BigBlind = 0
	result.SmallBlind = 0
	result.Users = make([]*ThUser, TH_DESK_MAX_START_USER)
	result.RemainTime = 0
	result.BetUserRaiseUserId = 0
	result.RoundCount = 0
	result.NewRoundBetUser = 0
	result.bianJackpot = 0
	result.Number = ThGameRoomIns.RandDeskNumber()
	result.SmallBlindCoin = ThGameRoomIns.SmallBlindCoin
	result.BigBlindCoin = 2 * ThGameRoomIns.SmallBlindCoin
	result.Status = TH_DESK_STATUS_STOP        //游戏还没有开始的状态
	result.CanRaise = 1
	result.DeskType = 0                           //游戏桌子的类型
	result.JuCount = 0
	return result
}


//获取数据库的id
func newThDeskId() int32 {
	// 获取连接 connection
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)

	id, _ := c.NextSeq(casinoConf.DB_NAME, casinoConf.DBT_T_TH_DESK, casinoConf.DB_ENSURECOUNTER_KEY)
	return int32(id)
}

func (t *ThDesk) LogString() {
	log.T("当前desk[%v]的信息:-----------------------------------begin----------------------------------", t.Id)
	log.T("当前desk[%v]的信息的状态status[%v]", t.Id, t.Status)
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("当前desk[%v]的user[%v]的状态status[%v],牌的信息[%v]", t.Id, u.UserId, u.Status, u.Cards)
		}
	}
	log.T("当前desk[%v]的信息的状态,总人数SeatedCount[%v],在线人数[%v]", t.Id, t.UserCount, t.UserCountOnline)
	log.T("当前desk[%v]的信息的状态,小盲注[%v]", t.Id, t.SmallBlind)
	log.T("当前desk[%v]的信息的状态,大盲注[%v]", t.Id, t.BigBlind)
	log.T("当前desk[%v]的信息的状态,压注人[%v]", t.Id, t.BetUserNow)
	log.T("当前desk[%v]的信息的状态,压注轮次[%v]", t.Id, t.RoundCount)
	log.T("当前desk[%v]的信息的状态,NewRoundBetUser[%v]", t.Id, t.NewRoundBetUser)
	log.T("当前desk[%v]的信息的状态,总共押注Jackpot[%v]", t.Id, t.Jackpot)
	log.T("当前desk[%v]的信息的状态,边池bianJackpot[%v]", t.Id, t.bianJackpot)
	log.T("当前desk[%v]的信息的状态,当前加注的人BetUserRaiseUserId[%v]", t.Id, t.BetUserRaiseUserId)

	log.T("当前desk[%v]的信息:-----------------------------------end----------------------------------", t.Id)
}

//打印测试信息用
func (t *ThDesk) LogStirngWinCoin() {
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("用户[%v]的wincoin[%v],detail[%v]", u.UserId, u.winAmount, u.winAmountDetail)
		}
	}
}


//判断是否重新进入
func (t *ThDesk) IsrepeatIntoWithRoomKey(userId uint32, a gate.Agent) bool {
	//1,先判断用户是否已经在房间了
	for i := 0; i < len(t.Users); i ++ {
		u := t.Users[i]
		if u != nil && u.UserId == userId {
			//如果u!=nil 那么
			u.agent = a                                                //设置用户的连接
			u.BreakStatus = TH_USER_BREAK_STATUS_FALSE                //设置用户的离线状态
			return true
		}
	}
	return false
}

/**
	为桌子增加一个人
 */
func (t *ThDesk) AddThUser(userId uint32, roomCoin int64, a gate.Agent) error {

	//1,从redis得到redisUser
	redisUser := userService.GetUserById(userId)
	//2,通过userId 和agent 够做一个thuser
	thUser := NewThUser()
	thUser.UserId = userId
	thUser.agent = a
	thUser.Status = TH_USER_STATUS_WAITSEAT        //刚进房间的玩家
	thUser.deskId = t.Id                //桌子的id
	thUser.NickName = *redisUser.NickName
	thUser.RoomCoin = roomCoin
	log.T("初始化thuser的时候coin[%v]:,roomCoin[%v]", thUser.GetCoin(), thUser.GetRoomCoin())

	//3,添加thuser
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] == nil {
			thUser.Seat = int32(i)                //给用户设置位置编号
			t.Users[i] = thUser
			break
		}
		if (i + 1) == len(t.Users) {
			log.E("玩家加入桌子失败")
			return errors.New("加入房间失败")
		}
	}

	//4, 把用户的信息绑定到agent上
	userAgentData := &gamedata.AgentUserData{}
	userAgentData.UserId = userId
	userAgentData.ThDeskId = t.Id
	a.SetUserData(userAgentData)


	//5,等待的用户加1
	t.UserCount ++
	t.UserCountOnline ++
	t.AgentMap[userId] = a
	return nil
}



//  用户退出德州游戏的房间,rmUser 需要在事物中进行
func (t *ThDesk) LeaveThuser(userId uint32) error {
	t.Lock()
	defer t.Unlock()
	user := t.GetUserByUserId(userId)
	user.Status = TH_USER_STATUS_LEAVE        //设置状态为离开
	return nil
}

//设置用户为掉线的状态
func (t *ThDesk) SetOfflineStatus(userId uint32) error {

	u := t.GetUserByUserId(userId)
	//1,设置状态为断线

	//这里需要保存掉线钱的状态
	u.BreakStatus = TH_USER_BREAK_STATUS_TRUE        //设置为掉线的状态
	t.UserCountOnline --
	return nil

}

//初始化前注的信息
func (t *ThDesk) OinitPreCoin() error {
	log.T("开始一局新的游戏,现在开始初始化前注的信息")
	ret := &bbproto.Game_PreCoin{}
	ret.Pool = new(int64)
	//ret.Coin = t.GetCoin()
	ret.Coin = t.GetRoomCoin()
	*ret.Pool = 0
	//ret.Precoin
	t.THBroadcastProtoAll(ret)
	log.T("开始一局新的游戏,现在开始初始化前注的信息完毕....")
	return nil
}

// 盲注开始押注
func (t *ThDesk) InitBlindBet() error {
	log.T("开始一局新的游戏,现在开始初始化盲注的信息")
	//小盲注押注
	t.AddBetCoin(t.SmallBlindCoin)
	t.caclUserCoin(t.SmallBlind, t.SmallBlindCoin)

	//大盲注押注
	t.AddBetCoin(t.BigBlindCoin)
	t.caclUserCoin(t.BigBlind, t.BigBlindCoin)

	//发送盲注的信息
	log.T("开始广播盲注的信息")
	blindB := &bbproto.Game_BlindCoin{}
	blindB.Banker = new(int32)
	blindB.Bigblindseat = new(int32)
	blindB.Smallblindseat = new(int32)

	//初始化默认值
	blindB.Tableid = &t.Id        //deskid
	//blindB.Matchid = &room.ThGameRoomIns.Id //roomId
	*blindB.Banker = t.GetUserByUserId(t.Dealer).Seat        //int32(t.GetUserIndex(t.Dealer))        //庄
	blindB.Bigblind = &t.BigBlindCoin        //大盲注
	blindB.Smallblind = &t.SmallBlindCoin        //小盲注
	*blindB.Bigblindseat = t.GetUserByUserId(t.BigBlind).Seat        //	int32(t.GetUserIndex(t.BigBlind))        //大盲注座位号
	*blindB.Smallblindseat = t.GetUserByUserId(t.SmallBlind).Seat        //int32(t.GetUserIndex(t.SmallBlind))        //小盲注座位号
	//blindB.Coin = t.GetCoin()        //每个人手中的coin
	blindB.Coin = t.GetRoomCoin()
	blindB.Handcoin = t.GetHandCoin()        //每个人下注的coin
	blindB.Pool = &t.Jackpot        //奖池
	t.THBroadcastProto(blindB, 0)
	log.T("广播盲注的信息完毕")

	log.T("开始一局新的游戏,现在开始初始化盲注的信息完毕....")
	return nil
}

/**
	游戏开始的时候,初始化玩家的信息
 */
func (t *ThDesk) InitUserBeginStatus() error {
	log.T("开始一局新的游戏,开始初始化用户的状态")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.BreakStatus == TH_USER_BREAK_STATUS_FALSE {
			u.Status = TH_USER_STATUS_BETING
			u.HandCoin = 0
			u.TurnCoin = 0
			u.winAmount = 0
			u.winAmountDetail = nil
		}
	}

	log.T("开始一局新的游戏,初始化用户的状态完毕")
	return nil
}

/**
	初始化纸牌的信息
 */
func (t *ThDesk) OnInitCards() error {
	log.T("开始一局新的游戏,初始化牌的信息")
	var total = int(2 * TH_DESK_MAX_START_USER + 5); //人数*手牌+5张公共牌
	totalCards := pokerService.RandomTHPorkCards(total)        //得到牌

	//得到所有的牌,前五张为公共牌,后边的每两张为手牌
	t.PublicPai = totalCards[0:5]
	log.T("初始化得到的公共牌的信息:")
	//给每个人分配手牌
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil {
			begin := i * 2 + 5
			end := i * 2 + 5 + 2
			t.Users[i].Cards = totalCards[begin:end]
			log.T("用户[%v]的手牌[%v]", t.Users[i].UserId, t.Users[i].Cards)
			t.Users[i].thCards = pokerService.GetTHPoker(t.Users[i].Cards, t.PublicPai, 5)
			log.T("用户[%v]的:拍类型,所有牌[%v],th[%v]", t.Users[i].UserId, t.Users[i].thCards.ThType, t.Users[i].thCards.Cards, t.Users[i].thCards)
		}
	}
	log.T("开始一局新的游戏,初始化牌的信息完毕...")
	return nil
}

//广播porto消息的通用方法
func (t *ThDesk) THBroadcastProto(p proto.Message, ignoreUserId uint32) error {
	log.Normal("开始广播proto消息【%v】", p)
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]                //给这个玩家发送广播信息
		if u != nil && u.UserId != ignoreUserId && u.Status != TH_USER_STATUS_LEAVE {
			a := t.Users[i].agent
			a.WriteMsg(p)
		}
	}
	return nil
}

//给全部人发送广播
func (t *ThDesk) THBroadcastProtoAll(p proto.Message) error {
	return t.THBroadcastProto(p, 0)
}

func (t *ThDesk) Testb(p proto.Message) {
	log.Normal("给每个房间发送proto 消息%v", p)
	for key := range t.AgentMap {
		log.Normal("开始给%v发送消息", key)
		//首先判断连接是否有断开
		a := t.AgentMap[key]
		a.WriteMsg(p)
		log.Normal("给%v发送消息,发送完毕", key)
	}
}

/**
当有新用户进入房间的时候,为其他人广播新过来的人的信息
 */
func (t *ThDesk) THBroadcastAddUser(newUserId, ignoreUserId uint32) error {
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil && t.Users[i].UserId != ignoreUserId {
			users := t.GetResUserModelClieSeq(t.Users[i].UserId)                //v1 需要更具当前用户进行排序
			broadUsers := &bbproto.THRoomAddUserBroadcast{}
			broadUsers.Header = protoUtils.GetSuccHeaderwithUserid(&(t.Users[i].UserId))
			for i := 0; i < len(users); i++ {
				if users[i] != nil && users[i].User.GetId() == newUserId {
					broadUsers.User = users[i]
					break
				}
			}

			a := t.Users[i].agent
			log.Normal("给userId[%v]发送消息:[%v]", t.Users[i].UserId, broadUsers)
			a.WriteMsg(broadUsers)
		}
	}
	return nil
}


/**
	返回res需要的User实体
 */
func (t *ThDesk) GetResUserModel() []*bbproto.THUser {
	result := make([]*bbproto.THUser, TH_DESK_MAX_START_USER)

	//就坐的人
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil {
			result[i] = t.Users[i].trans2bbprotoThuser()
			*result[i].SeatNumber = int32(i)
		} else {
			result[i] = &bbproto.THUser{}
			result[i].SeatNumber = new(int32)        //设置为0
		}
	}

	return result
}

func (t *ThDesk) GetResUserModelById(userId uint32) *bbproto.THUser {
	var result *bbproto.THUser
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil &&  t.Users[i].UserId == userId {
			result = t.Users[i].trans2bbprotoThuser()
		}
	}
	log.T("通过userId得到的bbproto.THUser的情况【%v】", result)
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

	log.T("得到排序后的bbproto.THUser的情况的信息[%v]", users)
	return users
}




// 	初始化第一个押注的人,当前押注的人
func (t *ThDesk) OninitThDeskBeginStatus() error {
	log.T("开始一局游戏,现在初始化desk的信息")
	//设置德州desk状态//设置状态为开始游戏
	t.Status = TH_DESK_STATUS_SART

	userTemp := make([]*ThUser, len(t.Users))
	copy(userTemp, t.Users)
	//这里需要定义一个庄家,todo 暂时默认为第一个,后边再修改
	var dealerIndex int = 0;
	if t.Dealer == 0 {
		log.T("游戏没有庄家,现在默认初始化第一个人为庄家")
		t.Dealer = t.Users[0].UserId
	} else {
		dealerIndex = t.GetUserIndex(t.Dealer)
		for i := dealerIndex; i < len(t.Users); i++ {
			u := t.Users[(i + 1) % len(t.Users)]
			if u != nil {
				t.Dealer = u.UserId
			}
		}
	}


	//设置小盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil {
			t.SmallBlind = u.UserId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	//设置大盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil {
			t.BigBlind = u.UserId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	/**

	 */

	if t.UserCountOnline == int32(2) {
		//如果只有两个人,当前押注的人是小盲注
		t.BetUserNow = t.SmallBlind
	} else {
		//设置当前押注的人
		for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
			u := userTemp[(i + 1) % len(userTemp)]
			if u != nil {
				t.BetUserNow = u.UserId
				userTemp[(i + 1) % len(userTemp)] = nil
				break
			}
		}
	}

	t.BetUserRaiseUserId = t.BetUserNow        //第一个加注的人
	t.NewRoundBetUser = t.SmallBlind           //新一轮开始默认第一个押注的人,第一轮默认是小盲注
	t.RoundCount = TH_DESK_ROUND1
	t.BetAmountNow = t.BigBlindCoin                   //设置第一次跟住时的跟注金额应该是多少
	t.MinRaise = t.BigBlindCoin
	t.Jackpot = 0
	t.bianJackpot = 0
	t.AllInJackpot = nil                          // 初始化allInJackpot 为空

	//本次押注的热开始等待
	waitUser := t.GetUserByUserId(t.BetUserNow)
	waitUser.wait()

	log.T("初始化游戏之后,庄家[%v]", t.Dealer)
	log.T("初始化游戏之后,小盲注[%v]", t.SmallBlind)
	log.T("初始化游戏之后,大盲注[%v]", t.BigBlind)
	log.T("初始化游戏之后,当前押注Id[%v]", t.BetUserNow)
	log.T("初始化游戏之后,第一个加注的人Id[%v]", t.BetUserRaiseUserId)
	log.T("初始化游戏之后,当前轮数Id[%v]", t.RoundCount)
	log.T("初始化游戏之后,当前jackpot[%v]", t.Jackpot)
	log.T("初始化游戏之后,当前bianJackpot[%v]", t.bianJackpot)
	log.T("初始化游戏之后,当前总人数[%v]", t.UserCount)
	log.T("初始化游戏之后,当前在线人数[%v]", t.UserCountOnline)

	log.T("开始一局游戏,现在初始化desk的信息完毕...")

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
	//如果处于押注状态的人只有一个人了,那么是开奖的时刻
	//

	/**
	//var TH_USER_STATUS_WAITSEAT int32 = 1        //刚上桌子 等待开始的玩家
	//var TH_USER_STATUS_SEATED int32 = 2                //刚上桌子 游戏中的玩家
	//var TH_USER_STATUS_BETING int32 = 3                //押注中
	//var TH_USER_STATUS_ALLINING int32 = 4        //allIn
	//var TH_USER_STATUS_FOLDED int32 = 5                //弃牌
	//var TH_USER_STATUS_WAIT_CLOSED int32 = 6                //等待结算
	//var TH_USER_STATUS_CLOSED int32 = 7                //已经结算
	//var TH_USER_STATUS_LEAVE int32 = 8                //
	//var TH_USER_STATUS_BREAK int32 = 9                //已经结算
	 */
	log.T("判断是否应该开奖,打印每个人的信息://1,刚上桌子,2,坐下,3押注中,4,allin,5,弃牌,6,等待结算,7,已经结算,8,离线,9断线")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("[%v]判断是否应该开奖,打印user[%v]的状态[%v]:", i, u.UserId, u.Status)
		}
	}

	var betingCount int = 0
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil && t.Users[i].Status == TH_USER_STATUS_BETING {
			betingCount ++
		}
	}

	log.T("当前处于押注中的人数是[%v]", betingCount)
	//如果押注的人只有一个人了,那么是开奖的时刻
	if betingCount <= 1 {
		log.T("现在处于押注中(beting)状态的人,只剩下一个了,所以直接开奖")
		return true
	}

	//第四轮,并且计算出来的押注人和start是同一个人
	if t.RoundCount == TH_DESK_ROUND_END {
		log.T("现在处于第[%v]轮押注,所以可以直接开奖", t.RoundCount)
		return true
	}

	//如果只有一个人没有all in  或者全部都all in 了也要开牌

	return false
}



//计算牌面是否赢
func (t *ThDesk) CalcThcardsWin() error {
	log.T("开始计算谁的牌是赢牌:")

	log.T("打印每个人牌的信息:")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED {
			log.T("玩家[%v] ", u.UserId)
			log.T("牌的信息:牌类型[%v],", u.thCards.ThType)
			log.T("所有牌[%v]", u.thCards.Cards)
		}
	}

	//1去牌最大的user,需要满足条件1,牌最大,2,状态是等待开奖
	//var userWin *ThUser
	//
	////1.1, 去第一个状态是等待开奖的那个人
	//for i := 1; i < len(t.Users); i++ {
	//	if t.Users[i].Status == TH_USER_STATUS_WAIT_CLOSED {
	//		userWin = t.Users[i]
	//		break;
	//	}
	//}
	//
	////1.2 比较得到牌是最大的,并且状态是等待开奖的那个人
	//for i := 1; i < len(t.Users); i++ {
	//	if t.Less(userWin, t.Users[i]) && t.Users[0].Status == TH_USER_STATUS_WAIT_CLOSED {
	//		userWin = t.Users[i]
	//	}
	//}

	userWin := t.Users[0]                //最大的牌的userId
	for i := 1; i < len(t.Users); i++ {
		if t.Less(userWin, t.Users[i]) {
			userWin = t.Users[i]
		}
	}

	//1.3打印得到那个人的信息
	if userWin == nil {
		log.E("服务器出错,没有找到赢牌的人...")
		return errors.New("没有找到赢牌的人")
	} else {
		log.T("得到的牌最大的user[%v]的信息[%v]", userWin.UserId, userWin)
	}

	//赢牌的人依次置为1,每个用户的牌都需要和最大的用户的牌来想比较
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil &&
		u.Status == TH_USER_STATUS_WAIT_CLOSED  &&
		pokerService.ThCompare(userWin.thCards, u.thCards) == pokerService.THPOKER_COMPARE_EQUALS {
			u.thCards.IsWin = true
		}
	}

	//------------------------------------------------打印测试信息
	log.T("开始计算谁的牌是赢牌,计算出来的结果:")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED {
			log.T("user[%v]的牌 isWin[%v]", u.UserId, u.thCards.IsWin)
		}
	}

	return nil
}

//比较两张牌的大小
func (t *ThDesk) Less(u1, u2 *ThUser) bool {

	if u1 == nil || u1.Status != TH_USER_STATUS_WAIT_CLOSED {
		return true
	}

	if u2 == nil || u2.Status != TH_USER_STATUS_WAIT_CLOSED {
		return false
	}

	//必将两个人的牌,u1的牌是否大于u2的牌
	ret := pokerService.ThCompare(u1.thCards, u2.thCards)
	if ret == pokerService.THPOKER_COMPARE_BIG {
		return false
	} else {
		return true
	}
}

//设置用户的状态喂等待开奖
func (t *ThDesk) SetStatusWaitClose() error {
	//设置用户的状态为等待开奖
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			u.InitWait()        //不再等待
			if u.Status == TH_USER_STATUS_ALLINING || u.Status == TH_USER_STATUS_BETING {
				//如果用户当前的状态是押注中,或者all in,那么设置用户的状态喂等待结算
				u.Status = TH_USER_STATUS_WAIT_CLOSED
			} else {
				u.Status = TH_USER_STATUS_CLOSED
			}
		}
	}
	return nil

}


//开奖
/**
	开奖的规则
	1,判断是谁赢了
 */

func (t *ThDesk) Lottery() error {
	log.T("现在开始开奖,并且发放奖励....")

	//todo 开奖之前 是否需要把剩下的牌 全部发完**** 目前是不可能
	//设置桌子的状态为开奖中
	t.Status = TH_DESK_STATUS_LOTTERY

	//设置用户的状态都为的等待开奖
	t.SetStatusWaitClose()

	//需要计算本局allin的奖金池
	t.CalcAllInJackpot()

	//todo 做结算是按照奖池来做,还是按照人员来做...
	//测试按照每个奖池来做计算
	for i := 0; i < len(t.AllInJackpot); i++ {
		log.T("现在开始开奖,计算allInJackpot[%v]的奖励....", i)
		a := t.AllInJackpot[i]
		if a != nil {
			//对这个奖池做计算
			/**
				1,几个人的牌是赢牌
				2,user的状态必须是没有结算的状态
			 */
			var winCount int = t.GetWinCount()
			bonus := a.Jackpopt / int64(winCount)        //每个人赢的奖金
			//这里吧奖金发放给每个人之后,需要把这局allin的人排除掉,再来计算剩下的人的将近
			//牌的方式只需要把这个人的状态设置为已经结清就行了
			for j := 0; j < len(t.Users); j++ {
				//todo 这里的将近可以选择使用一个数组存醋,方面clien制作动画
				//todo 目前只是计算总的金额
				u := t.Users[i]
				if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
					//可以发送奖金
					log.T("用户在allin.index[%v]活的奖金[%v]", i, bonus)
					u.winAmount += bonus
					u.RoomCoin += bonus
					u.winAmountDetail = append(u.winAmountDetail, bonus)
					userService.IncreasUserCoin(u.UserId, bonus)
				}

				//如果用户是这个奖金池all in的用户,则此用户设置喂已经结清的状态
				if u.UserId == a.UserId {
					u.Status = TH_USER_STATUS_CLOSED
				}
			}
		}
	}

	//计算边池的奖金	t.bianJackpot,同样需要看是几个人赢,然后评分将近
	log.T("现在开始开奖,计算边池的奖励....")
	bwinCount := t.GetWinCount()

	log.T("获奖的人数:[%v]", bwinCount)
	if bwinCount != 0 {

		var bbonus int64 = t.bianJackpot / int64(bwinCount)
		for i := 0; i < len(t.Users); i++ {
			u := t.Users[i]

			if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
				//
				//对这个用户做结算...
				log.T("现在开始开奖,计算边池的奖励,user[%v]得到[%v]....", u.UserId, bbonus)
				u.winAmount += bbonus
				u.RoomCoin += bbonus
				u.winAmountDetail = append(u.winAmountDetail, bbonus)        //详细的奖励(边池主池分开)
				userService.IncreasUserCoin(u.UserId, bbonus)
			}

			//设置为结算完了的状态
			if u != nil {
				u.Status = TH_USER_STATUS_CLOSED        //结算完了之后需要,设置用户的状态为已经结算
			}

		}
		//log.T	("现在开始开奖,计算奖励之后t.getWinCoinInfo()[%v]", t.getWinCoinInfo())

	}

	//todo 这里需要删除 打印测试信息的代码
	t.LogStirngWinCoin()

	//保存数据到数据库
	t.SaveLotteryData()

	// 新的开奖协议
	result := &bbproto.Game_TestResult{}
	result.Tableid = &t.Id                                //桌子
	result.BCanShowCard = t.GetBshowCard()                //
	result.BShowCard = t.GetBshowCard()                //亮牌
	result.Handcard = t.GetHandCard()                //手牌
	result.WinCoinInfo = t.getWinCoinInfo()
	result.HandCoin = t.GetHandCoin()
	t.THBroadcastProto(result, 0)

	//开奖之后,设置状态为 没有开始游戏
	//
	log.T("开奖结束,设置desk的状态为stop")
	t.Status = TH_DESK_STATUS_STOP        //设置喂没有开始开始游戏

	//开奖时间是在5秒之后开奖
	time.Sleep(TH_LOTTERY_DURATION)
	go t.OGRun()

	return nil
}


//保存数据到数据库
func (t *ThDesk)  SaveLotteryData() error {

	//得到连接
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)

	//循环对每个人做处理
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u == nil {
			continue
		}

		if u.Status != TH_USER_STATUS_CLOSED {
			log.E("保存用户[%v]信息到数据库的时候出错,状态[%v]不正确", u.UserId, u.Status)
			continue
		}

		//1,修改user在redis中的数据
		userService.IncreasUserCoin(u.UserId, u.winAmount)        //更新redis中的数据
		userService.FlashUser2Mongo(u.UserId)                        //刷新redis中的数据到mongo
		//2,保存游戏相关的数据
		//todo  游戏相关的数据结构 还没有建立
		thData := &mode.T_th{}
		thData.Mid = bson.NewObjectId()
		thData.BetAmount = u.TotalBet
		thData.UserId = u.UserId
		thData.DeskId = u.deskId
		thData.WinAmount = u.winAmount
		thData.Blance = *(userService.GetUserById(u.UserId).Coin)
		s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TH).Insert(thData)
	}

	return nil
}

//得到这句胜利的人有几个
func (t *ThDesk) GetWinCount() int {
	t.CalcThcardsWin()        //先计算牌的局面

	var result int = 0
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
			//如果用户不为空,并且状态是等待结算,牌的信息现实的是win 那么,表示一个赢的人
			result ++
		}
	}
	log.T("本局总共有[%v]人是赢牌,", result)
	return result
}


//跟注:跟注的时候 不需要重新设置押注的人
//只是跟注,需要减少用户的资产,增加奖池的金额
func (t *ThDesk) BetUserCall(user  *ThUser) error {
	//log.T("用户[%v]押注coin[%v]", userId, coin)
	followCoin := t.BetAmountNow - user.HandCoin
	if user.RoomCoin <= followCoin {
		//allin
		t.BetUserAllIn(user.UserId, followCoin)
	} else {
		//1,增加奖池的金额
		t.AddBetCoin(followCoin)
		//2,增加用户本轮投注的金额
		t.caclUserCoin(user.UserId, followCoin)
	}
	return nil
}

func (t *ThDesk) AddBetCoin(coin int64) error {
	t.Jackpot += coin                        //底池 增加
	t.bianJackpot += coin                        //边池 增加
	return nil
}

//如果弃牌的人是 t.NewRoundBetUser ,需要重新设置值
func (t *ThDesk) NextNewRoundBetUser() error {
	index := t.GetUserIndex(t.NewRoundBetUser)
	for i := index + 1; i < len(t.Users) + index; i++ {
		u := t.Users[i % len(t.Users)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.NewRoundBetUser = u.UserId
			break
		}
		//如果没有找到,那么返回失败
		if i == (len(t.Users) + index - 1) {
			return errors.New("没有找到下一个默认开始的押注的人")
		}
	}

	return nil

}

//让牌:只有第一个人才可以让牌
func (t *ThDesk) BetUserCheck(userId uint32) error {
	if userId == t.BetUserRaiseUserId {
		//第一个人的时候才可以让牌

		//设置喂第一个让牌的人
		if t.CheckUserFirst == 0 {
			t.CheckUserFirst = 0
			t.CheckUserFirst = userId
		}

		//设置一个押注的人为下一个人
		index := t.GetUserIndex(userId)
		for i := index; i < len(t.Users) + index - 1; i++ {
			u := t.Users[(i + 1) % len(t.Users)]
			if u != nil && u.Status == TH_USER_STATUS_BETING && u.UserId != t.CheckUserFirst {
				t.BetUserRaiseUserId = userId
				break
			}
		}

	}
	return nil
}

//用户加注
func (t *ThDesk) BetUserRaise(user *ThUser, coin int64) error {
	t.BetAmountNow += coin
	betCoin := t.BetAmountNow - user.HandCoin

	if betCoin == user.RoomCoin {
		// allin
		t.BetUserAllIn(user.UserId, betCoin)
	} else {
		//1,增加奖池的金额
		t.AddBetCoin(betCoin)                                //desk-coin
		//2,减少用户的金额
		t.caclUserCoin(user.UserId, betCoin)                        //thuser
		userService.DecreaseUserCoin(user.UserId, betCoin)        //redis-user
		//3,设置状态:设置为第一个加注的人,如果后边所有人都是跟注,可由这个人判断一轮是否结束
		t.BetUserRaiseUserId = user.UserId
	}
	return nil
}

//用户AllIn
func (t *ThDesk) BetUserAllIn(userId uint32, coin int64) error {
	//1,增加奖池的金额
	t.AddBetCoin(coin)

	//2,减少用户的金额
	t.caclUserCoin(userId, coin)

	//3,增加all in的状态
	allinpot := &pokerService.AllInJackpot{}
	allinpot.UserId = userId
	allinpot.Jackpopt = 0
	allinpot.ThroundCount = t.RoundCount
	allinpot.AllInAmount = t.GetUserByUserId(userId).TotalBet
	t.AllInJackpot = append(t.AllInJackpot, allinpot)        //增加一个池子
	log.T("用户[%v] all in 的时候,allin的值是[%v]", allinpot.UserId, allinpot.AllInAmount)
	return nil
}

/**
	根据userId 找到在桌子上的index
 */
func (t *ThDesk) GetUserIndex(userId uint32) int {
	var result int = 0
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] != nil && t.Users[i].UserId == userId {
			result = i
			break
		}
	}
	return result
}

//通过UserId找到User
func (t *ThDesk) GetUserByUserId(userId uint32) *ThUser {
	index := t.GetUserIndex(userId)
	return t.Users[index]
}

func (t *ThDesk) caclUserCoin(userId uint32, coin int64) error {
	user := t.GetUserByUserId(userId)
	user.TurnCoin += coin
	user.HandCoin += coin
	user.TotalBet += coin
	user.RoomCoin -= coin                //这里暂时不处理roomCoin,roomCoin是在每一轮结束的时候来结算
	userService.DecreaseUserCoin(userId, coin)
	return nil
}


/**
	初始化下一个押注的人
	初始化下一个人的时候需要一个超时的处理
 */
func (t *ThDesk) NextBetUser() error {

	log.T("当前押注的人是userId[%v]", t.BetUserNow)
	index := t.GetUserIndex(t.BetUserNow)
	for i := index; i < len(t.Users) + index; i++ {
		u := t.Users[(i + 1) % len(t.Users)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			log.T("设置betUserNow 为[%v]", u.UserId)
			t.BetUserNow = u.UserId
			break
		}
	}

	//设置新一轮
	if t.BetUserRaiseUserId == t.BetUserNow {
		//处理allin 奖金池分割的问题
		t.CalcAllInJackpot()
		t.BetUserRaiseUserId = t.NewRoundBetUser
		t.BetUserNow = t.NewRoundBetUser
		t.BetAmountNow = 0        //下一句重新开始的时候,设置当前押注的人为0
		t.RoundCount ++

		log.T("设置下次押注的人是小盲注,下轮次[%v]", t.RoundCount)
	}

	//打印当前桌子的信息
	t.LogString()
	return nil

}

//下一轮
func (t *ThDesk) nextRoundInfo() {

	if !t.isNewRound() {
		return
	}

	log.T("本次设置的押注人和之前的是同一个人,所以开始第[%v]轮的游戏", t.RoundCount)
	//一轮完之后需要发送完成的消息
	sendData := NewGame_SendOverTurn()
	*sendData.Tableid = t.Id
	*sendData.MinRaise = t.MinRaise
	*sendData.NextSeat = t.GetUserByUserId(t.BetUserNow).Seat        //int32(t.GetUserIndex(t.BetUserNow))
	sendData.Handcoin = t.GetHandCoin()
	//sendData.Coin = t.GetCoin()
	sendData.Coin = t.GetRoomCoin()
	*sendData.Pool = t.Jackpot
	sendData.SecondPool = t.GetSecondPool()
	t.THBroadcastProto(sendData, 0)

	//清空每个人的handCoin
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			u.HandCoin = 0
			u.TurnCoin = 0
		}
	}

	switch t.RoundCount {
	case TH_DESK_ROUND2:
		//发三张公共牌
		t.sendFlopCard()
	case TH_DESK_ROUND3:
		//发第四章牌
		t.sendTurnCard()
	case TH_DESK_ROUND4:
		//发第五章牌
		t.sendRiverCard()
	}

	//做redis的持久化
	/**
		需要持久化的数据
		1,当前桌子的状态
		2,当鸭玩家的状态
	 */

}




//判断是否是新的一局
func (t *ThDesk) isNewRound() bool {
	//
	log.T("判断是否是新的一轮t.BetUserRaiseUserId[%v],t.BetUserNow(%v),t.status[%v].//status:1,stop,2,start,3,lottery", t.BetUserRaiseUserId, t.BetUserNow, t.Status)
	if t.BetUserRaiseUserId == t.BetUserNow &&  t.Status == TH_DESK_STATUS_SART {
		log.T("t.BetUserRaiseUserId[%v] == t.BetUserNow[%v],新的一局开始", t.BetUserRaiseUserId, t.BetUserNow)
		return true
	} else {
		return false
	}
}


//发三张公共牌
func (t *ThDesk) sendFlopCard() error {
	log.T("发送三张公共牌begin")
	//申明
	result := &bbproto.Game_SendFlopCard{}
	result.Tableid = &t.Id
	result.Card0 = ThCard2OGCard(t.PublicPai[0])
	result.Card1 = ThCard2OGCard(t.PublicPai[1])
	result.Card2 = ThCard2OGCard(t.PublicPai[2])

	//广播消息
	t.THBroadcastProto(result, 0)
	log.T("发送三张公共牌end")

	return nil
}


//发送第四张牌
func (t *ThDesk) sendTurnCard() error {
	log.T("发送第四张公共牌begin")

	result := &bbproto.Game_SendTurnCard{}
	result.Tableid = &t.Id
	result.Card = ThCard2OGCard(t.PublicPai[3])

	t.THBroadcastProto(result, 0)
	log.T("发送第四张公共牌end")

	return nil
}

//发送低五张牌
func (t *ThDesk) sendRiverCard() error {
	log.T("发送第五张公共牌begin")

	result := &bbproto.Game_SendRiverCard{}
	result.Tableid = &t.Id
	result.Card = ThCard2OGCard(t.PublicPai[4])

	t.THBroadcastProto(result, 0)
	log.T("发送第五张公共牌end")

	return nil
}

//计算奖金池的划分的问题
func (t *ThDesk) CalcAllInJackpot() error {
	log.T("开始计算allin将近池子begin")
	//1,对allin 进行排序,排序之后,可以对奖金池做划分,得到当前未做处理的all和边池的值
	var list pokerService.AllInJackpotList = t.AllInJackpot
	sort.Sort(list)

	for i := 0; i < len(t.AllInJackpot); i++ {
		log.T("第[%v]次循环的时候,allinlist[%v]", i, t.AllInJackpot)
		all := t.AllInJackpot[i]

		//如果这个池子为nil ,则跳过这个循环,如果这个allIn 不是本轮的,则把之前的allin 的jackpot 累加起来
		if all == nil || all.ThroundCount != t.RoundCount {
			continue
		}

		log.T("开始计算用户[%v]allIn.index[%v] allin.amount[%v]计算all in 时的池子金额", all.UserId, i, all.AllInAmount)
		//每个allin计算金额
		for n := 0; n < len(t.Users); n++ {
			u := t.Users[n]
			if u != nil {
				log.T("用户[%v]押注的总金额是[%v]")
				if u.TotalBet > all.AllInAmount {
					all.Jackpopt += all.AllInAmount
					u.TotalBet -= all.AllInAmount
					log.T("用户[%v]押注加入all的金额是[%v]", u.UserId, all.AllInAmount)
				} else {
					all.Jackpopt += u.TotalBet
					//*u.roundBet = 0
					log.T("用户[%v]押注加入all的金额是[%v]", u.UserId, u.TotalBet)
				}

			}
		}
		log.T("计算出来用户[%v]allIn.index[%v] allin.amount[%v]计算all in 的池子总金额", all.UserId, i, all.Jackpopt)

		//之后的allinamount - 当前allin
		for k := i; k < len(t.AllInJackpot); k++ {
			allk := t.AllInJackpot[k]
			if allk != nil {
				allk.AllInAmount -= all.AllInAmount
			}
		}
		t.bianJackpot -= all.Jackpopt
		log.T("开始给allIn[%v]计算all in 时的池子金额---------------------------------end---------------", i)
		log.T("目前t.bianJackPot 的剩余值是[%v]", t.bianJackpot)
	}
	log.T("计算出来的allIn:【%v】", t.AllInJackpot)
	log.T("开始计算allin将近池子end")
	return nil

}
func (t *ThDesk) CheckBetUserBySeat(user *ThUser) bool {
	//2,判断押注的用户是否是当前的用户
	if t.BetUserNow != user.UserId {
		return false
	}

	//3, 判断用户的状态是否正确
	err := user.CheckBetWaitStatus()
	if err != nil {
		return false
	}

	if user.RoomCoin <= 0 {
		log.E("用户的带入金额小于0,所以不能押注或者投注了")
		return false
	}

	//用户合法,设置等待状态
	user.InitWait()
	return true
}

//是不是可以开始游戏了
func (t *ThDesk) IsTime2begin() bool {

	//todo 这里需要增加锦标赛的逻辑

	log.T("判断是否可以开始一局新的游戏")
	/**
		开始游戏的要求:
		1,[在线]用户的人数达到了最低可玩人数
		2,当前的状态是游戏停止的状态
	 */


	log.T("当前玩家的状态://1,等待开始,2,游戏中,3,押注中,4,allin,5,弃牌,6,等待结算,7,已经结算,8,裂开,9,掉线")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("用户[%v].seat[%v]的状态是[%v]", u.UserId, u.Seat, u.Status)
		}
	}

	log.T("当前在线玩家的数量是[%v],当前desk的状态是[%v],1未开始,2游戏中,3,开奖中", t.UserCountOnline, t.Status)
	if t.UserCountOnline >= TH_DESK_LEAST_START_USER  && t.Status == TH_DESK_STATUS_STOP {
		log.T("游戏到了开始的时候----begin----")
		return true
	} else {
		log.T("游戏还不到开始的时候")
		return false
	}
}

//开始游戏
func (mydesk *ThDesk) OGRun() error {
	mydesk.Lock()
	defer mydesk.Unlock()

	log.T("\n\n开始一局新的游戏\n\n")
	//1,判断是否可以开始游戏
	if !mydesk.IsTime2begin() {
		return nil
	}

	//2,初始化玩家的信息
	err := mydesk.InitUserBeginStatus()
	if err != nil {
		log.E("开始游戏失败,errMsg[%v]", err.Error())
		return err
	}

	//2,初始化牌的信息
	err = mydesk.OnInitCards()
	if err != nil {
		log.E("开始德州扑克游戏,初始化扑克牌的时候出错")
		return err
	}

	//3,初始化游戏房间的状态
	err = mydesk.OninitThDeskBeginStatus()
	if err != nil {
		log.E("开始德州扑克游戏,初始化房间的状态的时候报错")
		return err
	}


	//3, 初始化前注的信息
	err = mydesk.OinitPreCoin()
	if err != nil {
		log.E("开始德州扑克游戏,初始化房间的状态的时候报错")
		return err
	}

	//4 初始化盲注开始押注
	err = mydesk.InitBlindBet()
	if err != nil {
		log.E("盲注下注的时候出错errMsg[%v]", err.Error())
		return err
	}

	log.T("广播Game_InitCard的信息")
	initCardB := &bbproto.Game_InitCard{}

	//设置默认值
	initCardB.Tableid = new(int32)
	initCardB.ActionTime = new(int32)
	initCardB.DelayTime = new(int32)
	initCardB.NextUser = new(int32)

	//设置初始化值
	*initCardB.Tableid = int32(mydesk.Id)
	initCardB.HandCard = mydesk.GetHandCard()
	initCardB.PublicCard = mydesk.ThPublicCard2OGC()
	initCardB.MinRaise = &mydesk.MinRaise
	*initCardB.NextUser = mydesk.GetUserByUserId(mydesk.BetUserNow).Seat                //	int32(mydesk.GetUserIndex(mydesk.BetUserNow))
	*initCardB.ActionTime = TH_TIMEOUT_DURATION_INT
	//initCardB.Seat = &mydesk.UserCount
	mydesk.THBroadcastProto(initCardB, 0)
	log.T("广播Game_InitCard的信息完毕")

	log.T("\n\n开始一局新的游戏,初始化完毕\n\n")
	return nil
}


//手牌转换为OG可以使用的牌
func (t *ThDesk) ThPublicCard2OGC() []*bbproto.Game_CardInfo {
	result := make([]*bbproto.Game_CardInfo, len(t.PublicPai))
	for i := 0; i < len(t.PublicPai); i++ {
		result[i] = ThCard2OGCard(t.PublicPai[i])
	}
	return result
}