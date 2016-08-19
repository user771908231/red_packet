package room

import (
	"casino_server/msg/bbprotogo"
	"sync"
	"time"
	"casino_server/conf/intCons"
	"casino_server/service/pokerService"
	"github.com/name5566/leaf/db/mongodb"
	"casino_server/conf/casinoConf"
	"fmt"
	"casino_server/mode"
	"casino_server/service/userService"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"sort"
	"casino_server/common/log"
	"casino_server/msg/bbprotoFuncs"
	"casino_server/utils/numUtils"
	"errors"
	"github.com/name5566/leaf/gate"
	"github.com/golang/protobuf/proto"
	"casino_server/utils/db"
)



/**
	一个德州扑克的房间
 */
type ThDesk struct {
	sync.Mutex
	Id             int32                              //roomid
	MatchId        int32                              //matchId
	DeskOwner      uint32                             //房主的id
	RoomKey        string                             //room 自定义房间的钥匙
	DeskType       int32                              //桌子的类型,1,表示自定义房间,2表示锦标赛的
	InitRoomCoin   int64                              //进入这个房间的roomCoin 带入金额标准是多少
	JuCount        int32                              //这个桌子最多能打多少局
	JuCountNow     int32                              //这个桌子已经玩了多少局
	SmallBlindCoin int64                              //小盲注的押注金额
	BigBlindCoin   int64                              //大盲注的押注金额
	BeginTime      time.Time                          //游戏开始时间
	EndTime              time.Time                    //游戏结束时间

	Dealer               uint32                       //庄家
	BigBlind             uint32                       //大盲注
	SmallBlind           uint32                       //小盲注
	RaiseUserId          uint32                       //加注的人的Id,一轮结束的判断需要按照这个人为准
	NewRoundFirstBetUser uint32                       //新一轮,开始押注的第一个人//第一轮默认是小盲注,但是当小盲注弃牌之后,这个人要滑倒下一家去
	BetUserNow           uint32                       //当前押注人的Id

	GameNumber           int32                        //每一局游戏的游戏编号
	Users                []*ThUser                    //坐下的人
	PublicPai            []*bbproto.Pai               //公共牌的部分
	UserCount            int32                        //玩游戏的总人数
	UserCountOnline      int32                        //在先人数
	Status               int32                        //牌桌的状态
	BetAmountNow         int64                        //当前的押注金额是多少
	RoundCount           int32                        //第几轮
	Jackpot              int64                        //奖金池
	edgeJackpot          int64                        //边池
	MinRaise             int64                        //最低加注金额
	AllInJackpot         []*pokerService.AllInJackpot //allin的标记
}

/**
	新生成一个德州的桌子
 */
func NewThDesk() *ThDesk {
	result := new(ThDesk)
	result.Id = newThDeskId()
	result.UserCount = 0
	result.UserCountOnline = 0
	result.Dealer = 0                //不需要创建  默认就是为空
	result.Status = TH_DESK_STATUS_STOP
	result.BetUserNow = 0
	result.BigBlind = 0
	result.SmallBlind = 0
	result.Users = make([]*ThUser, TH_DESK_MAX_START_USER)
	result.RaiseUserId = 0
	result.RoundCount = 0
	result.NewRoundFirstBetUser = 0
	result.edgeJackpot = 0
	result.SmallBlindCoin = ThGameRoomIns.SmallBlindCoin
	result.BigBlindCoin = 2 * ThGameRoomIns.SmallBlindCoin
	result.Status = TH_DESK_STATUS_STOP        //游戏还没有开始的状态
	result.DeskType = intCons.GAME_TYPE_TH_CS                          //游戏桌子的类型
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
			log.T("当前desk[%v]的user[%v]的状态status[%v],牌的信息[%v]", t.Id, u.UserId, u.Status, u.HandCards)
		}
	}
	log.T("当前desk[%v]的信息的状态,总人数SeatedCount[%v],在线人数[%v]", t.Id, t.UserCount, t.UserCountOnline)
	log.T("当前desk[%v]的信息的状态,小盲注[%v]", t.Id, t.SmallBlind)
	log.T("当前desk[%v]的信息的状态,大盲注[%v]", t.Id, t.BigBlind)
	log.T("当前desk[%v]的信息的状态,压注人[%v]", t.Id, t.BetUserNow)
	log.T("当前desk[%v]的信息的状态,压注轮次[%v]", t.Id, t.RoundCount)
	log.T("当前desk[%v]的信息的状态,NewRoundFirstBetUser[%v]", t.Id, t.NewRoundFirstBetUser)
	log.T("当前desk[%v]的信息的状态,总共押注Jackpot[%v]", t.Id, t.Jackpot)
	log.T("当前desk[%v]的信息的状态,边池bianJackpot[%v]", t.Id, t.edgeJackpot)
	log.T("当前desk[%v]的信息的状态,当前加注的人BetUserRaiseUserId[%v]", t.Id, t.RaiseUserId)

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
			u.IsBreak = false                //设置用户的离线状态
			return true
		}
	}
	return false
}

/**
	为桌子增加一个人
 */
func (t *ThDesk) AddThUser(userId uint32, roomCoin int64, userStatus int32, a gate.Agent) error {

	//1,从redis得到redisUser
	//redisUser := userService.GetUserById(userId)
	//2,通过userId 和agent 够做一个thuser
	thUser := NewThUser()
	thUser.UserId = userId
	thUser.agent = a
	thUser.Status = userStatus        //刚进房间的玩家
	thUser.deskId = t.Id                //桌子的id
	//thUser.NickName = *redisUser.NickName		//todo 测试阶段,把nickName显示成用户id
	thUser.NickName, _ = numUtils.Uint2String(userId)
	thUser.RoomCoin = roomCoin
	log.T("初始化thuser的时候coin[%v]:,roomCoin[%v]", thUser.GetCoin(), thUser.GetRoomCoin())

	//3,添加thuser
	err := t.addThuserBean(thUser)
	if err != nil {
		log.E("增加user【%v】到desk【%v】失败", thUser.UserId, t.Id)
		return errors.New("增加user失败")
	}

	//4, 把用户的信息绑定到agent上
	thUser.UpdateAgentUserData(a,t.Id,t.MatchId)

	//5,等待的用户加1
	t.UserCount ++
	t.UserCountOnline ++
	return nil
}

//增加一个user实体
func (t *ThDesk) addThuserBean(user *ThUser) error {
	for i := 0; i < len(t.Users); i++ {
		if t.Users[i] == nil {
			user.Seat = int32(i)                //给用户设置位置编号
			t.Users[i] = user
			return nil
		}

		if (i + 1) == len(t.Users) {
			log.E("玩家加入桌子失败")
			return errors.New("加入房间失败")
		}
	}

	return nil
}

//用户准备
func (t *ThDesk) Ready(userId uint32) error {
	user := t.GetUserByUserId(userId)
	user.Status = TH_USER_STATUS_READY
	return nil
}


//  用户退出德州游戏的房间,rmUser 需要在事物中进行
func (t *ThDesk) LeaveThuser(userId uint32) error {
	t.Lock()
	defer t.Unlock()
	user := t.GetUserByUserId(userId)
	user.IsLeave = true     //设置状态为离开
	return nil
}

//设置用户为掉线的状态
func (t *ThDesk) SetOfflineStatus(userId uint32) error {

	u := t.GetUserByUserId(userId)
	//1,设置状态为断线

	//这里需要保存掉线钱的状态
	u.IsBreak = true        //设置为掉线的状态
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
		//判断用户是否为空
		if u == nil {
			continue
		}

		log.T("用户[%v]的status[%v],BreakStatus[%v],:", u.UserId, u.Status, u.IsBreak)
		//新的一局开始,设置用户的状态
		u.HandCoin = 0
		u.TurnCoin = 0
		u.winAmount = 0
		u.TotalBet4calcAllin = 0
		u.TotalBet = 0                                //新的一局游戏开始,把总的押注金额设置为0
		u.winAmountDetail = nil

		//如果用户的余额不足或者用户的状态是属于断线的状态,则设置用户为等待入座
		if u.RoomCoin <= t.BigBlindCoin || u.IsBreak == true {
			log.T("由于用户[%v] status[%v],的roomCoin[%v] <= desk.BigBlindCoin 所以设置用户为TH_USER_STATUS_WAITSEAT", u.UserId, u.IsBreak, u.RoomCoin, t.BigBlindCoin)
			u.Status = TH_USER_STATUS_WAITSEAT        //只是坐下,没有游戏中
			continue
		}

		//用户不是离线的状态,并且,用户已经准备好了,则可以开始游戏
		if u.IsBreak == false && u.Status == TH_USER_STATUS_READY {
			log.T("由于用户[%v]的status[%v]BreakStatus[%v],所以设置状态为TH_USER_STATUS_BETING", u.UserId, u.Status, u.IsBreak)
			u.Status = TH_USER_STATUS_BETING
		}
	}

	//------------------------------------由于联众前端设计的问题...这里的user需要重新排列user的顺序------------------------------------
	usersTemp := make([]*ThUser, len(t.Users))
	copy(usersTemp, t.Users)
	log.T("原来的thsuers:[%v]", t.Users)
	log.T("复制的thsuers:[%v]", usersTemp)
	log.T("排序之前的thusers【%v】", t.Users)

	//是原来的thusers置为nil
	for i := 0; i < len(t.Users); i++ {
		t.Users[i] = nil
	}

	//排序游戏中的玩家
	for i := 0; i < len(usersTemp); i++ {
		u := usersTemp[i]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.addThuserBean(u)
		}
	}

	for i := 0; i < len(usersTemp); i++ {
		u := usersTemp[i]
		if u != nil && u.Status != TH_USER_STATUS_BETING {
			t.addThuserBean(u)
		}
	}
	log.T("排序之后的thusers【%v】", t.Users)

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
		//只有当用不为空,并且是在游戏中的状态的时候,才可以发牌
		if t.Users[i] != nil && t.Users[i].Status == TH_USER_STATUS_BETING {
			begin := i * 2 + 5
			end := i * 2 + 5 + 2
			t.Users[i].HandCards = totalCards[begin:end]
			log.T("用户[%v]的手牌[%v]", t.Users[i].UserId, t.Users[i].HandCards)
			t.Users[i].thCards = pokerService.GetTHPoker(t.Users[i].HandCards, t.PublicPai, 5)
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
		if u != nil && u.UserId != ignoreUserId && u.IsLeave != true {
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
			if u != nil && u.Status == TH_USER_STATUS_BETING {
				t.Dealer = u.UserId
			}
		}
	}

	//设置小盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.SmallBlind = u.UserId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	//设置大盲注
	for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
		u := userTemp[(i + 1) % len(userTemp)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.BigBlind = u.UserId
			userTemp[(i + 1) % len(userTemp)] = nil
			break
		}
	}

	/**
		设置第一个押注的人
	 */

	if t.UserCountOnline == int32(2) {
		//如果只有两个人,当前押注的人是小盲注
		t.BetUserNow = t.SmallBlind
	} else {
		//设置当前押注的人
		for i := dealerIndex; i < len(userTemp) + dealerIndex; i++ {
			u := userTemp[(i + 1) % len(userTemp)]
			if u != nil && u.Status == TH_USER_STATUS_BETING {
				t.BetUserNow = u.UserId
				userTemp[(i + 1) % len(userTemp)] = nil
				break
			}
		}
	}

	t.RaiseUserId = t.BetUserNow        //第一个加注的人
	t.NewRoundFirstBetUser = t.SmallBlind           //新一轮开始默认第一个押注的人,第一轮默认是小盲注
	t.RoundCount = TH_DESK_ROUND1
	t.BetAmountNow = t.BigBlindCoin                   //设置第一次跟住时的跟注金额应该是多少
	t.MinRaise = t.BigBlindCoin
	t.Jackpot = 0
	t.edgeJackpot = 0
	t.AllInJackpot = nil                          // 初始化allInJackpot 为空

	//本次押注的热开始等待
	waitUser := t.GetUserByUserId(t.BetUserNow)
	log.T("用户开始等待,等待之前桌子的状态[%v]:", t.Status)
	waitUser.wait()

	log.T("初始化游戏之后,庄家[%v]", t.Dealer)
	log.T("初始化游戏之后,小盲注[%v]", t.SmallBlind)
	log.T("初始化游戏之后,大盲注[%v]", t.BigBlind)
	log.T("初始化游戏之后,当前押注Id[%v]", t.BetUserNow)
	log.T("初始化游戏之后,第一个加注的人Id[%v]", t.RaiseUserId)
	log.T("初始化游戏之后,当前轮数Id[%v]", t.RoundCount)
	log.T("初始化游戏之后,当前jackpot[%v]", t.Jackpot)
	log.T("初始化游戏之后,当前bianJackpot[%v]", t.edgeJackpot)
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

	log.T("判断是否应该开奖,打印每个人的信息://1,刚上桌子,2,坐下,3,ready 4 押注中,5,allin,6,弃牌,7,等待结算,8,已经结算,9,离开")
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
	log.T("开奖之前答应每个人的状态,并且修改为等待结算")
	//设置用户的状态为等待开奖
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("用户[%v].nickname[%v]的status[%v]", u.UserId, u.NickName, u.Status)
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

//计算用户输赢盈眶
func (t *ThDesk) calcUserWinAmount() error {

	//计算每个allin池子的将近分配
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
			log.T("allInJackpot[%v]的总共有[%v]人获得奖励[%v],平均[%v]", i, winCount, a.Jackpopt, bonus)
			//这里吧奖金发放给每个人之后,需要把这局allin的人排除掉,再来计算剩下的人的将近
			//牌的方式只需要把这个人的状态设置为已经结清就行了
			for j := 0; j < len(t.Users); j++ {
				//todo 这里的将近可以选择使用一个数组存醋,方便clien制作动画
				//todo 目前只是计算总的金额
				u := t.Users[j]

				//如果用户nil,则直接开始下一次循环
				if u == nil {
					continue
				}

				//判断用户是否得奖
				if u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
					//可以发送奖金
					log.T("用户[%v].status[%v],iswin[%v]在allin.index[%v]活的奖金[%v]", u.UserId, u.Status, u.thCards.IsWin, i, bonus)
					u.winAmount += bonus
					u.RoomCoin += bonus
					u.winAmountDetail = append(u.winAmountDetail, bonus)
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
		var bbonus int64 = t.edgeJackpot / int64(bwinCount)
		for i := 0; i < len(t.Users); i++ {
			u := t.Users[i]

			if u != nil && u.Status == TH_USER_STATUS_WAIT_CLOSED && u.thCards.IsWin {
				//
				//对这个用户做结算...
				log.T("现在开始开奖,计算边池的奖励,user[%v]得到[%v]....", u.UserId, bbonus)
				u.winAmount += bbonus
				u.RoomCoin += bbonus
				u.winAmountDetail = append(u.winAmountDetail, bbonus)        //详细的奖励(边池主池分开)
			}

			//设置为结算完了的状态
			if u != nil {
				u.Status = TH_USER_STATUS_CLOSED        //结算完了之后需要,设置用户的状态为已经结算
			}

		}
		//log.T	("现在开始开奖,计算奖励之后t.getWinCoinInfo()[%v]", t.getWinCoinInfo())
	}

	return nil

}


//开奖
func (t *ThDesk) Lottery() error {
	log.T("现在开始开奖,并且发放奖励....")

	//todo 开奖之前 是否需要把剩下的牌 全部发完**** 目前是不可能
	//设置桌子的状态为开奖中
	t.Status = TH_DESK_STATUS_LOTTERY

	//设置用户的状态都为的等待开奖
	t.SetStatusWaitClose()

	//需要计算本局allin的奖金池
	t.CalcAllInJackpot()

	//计算用户输赢情况
	t.calcUserWinAmount()

	//todo 这里需要删除 打印测试信息的代码
	t.LogStirngWinCoin()

	//保存数据到数据库
	t.SaveLotteryData()

	//广播开奖的额结果
	t.broadLotteryResult()

	//开奖之后,设置状态为 没有开始游戏
	t.afterLottery()

	//判断是否可以继续啊下次游戏
	if t.canNextRun() {
		time.Sleep(TH_LOTTERY_DURATION)	//开奖的延迟
		go t.Run()
	}else{
		//表示不能继续开始游戏
		t.End()
	}
	return nil
}



//判断开奖之后是否可以继续游戏
func (t *ThDesk) canNextRun() bool{

	if t.DeskType == intCons.GAME_TYPE_TH {
		//如果是自定义房间
		//1,局数
		//2,人数
		if t.JuCountNow > t.JuCount {
			return false
		}
	}else if t.DeskType == intCons.GAME_TYPE_TH_CS{
		//锦标赛
		if ChampionshipRoom.CanNextDeskRun() {
			return true
		}else{
			return false
		}
	}

	return true
}


//广播开奖的结果
func (t *ThDesk) broadLotteryResult() error{
	result := &bbproto.Game_TestResult{}
	result.Tableid = &t.Id                           //桌子
	result.BCanShowCard = t.GetBshowCard()           //
	result.BShowCard = t.GetBshowCard()              //亮牌
	result.Handcard = t.GetHandCard()                //手牌
	result.WinCoinInfo = t.getWinCoinInfo()
	result.HandCoin = t.GetHandCoin()
	result.CoinInfo	= t.getCoinInfo()		//每个人的输赢情况

	//开始广播
	t.THBroadcastProtoAll(result)
	return nil

}

//开奖之后的处理
func (t *ThDesk) afterLottery() error {
	//1,设置游戏竹子的状态
	log.T("开奖结束,设置desk的状态为stop")
	t.Status = TH_DESK_STATUS_STOP        //设置喂没有开始开始游戏
	t.JuCountNow ++

	//2,设置用户的状态
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]

		if u != nil && u.IsBreak == false {
			if t.DeskType == intCons.GAME_TYPE_TH {
				//如果是自定义的房间,设置每个人都是坐下的状态
				u.Status = TH_USER_STATUS_SEATED
			} else {
				//如果是锦标赛的房间
				u.Status = TH_USER_STATUS_READY
			}
		}

	}

	//3,

	return nil
}


//保存数据到数据库

//这里需要根据游戏类型的不同来分别存醋

func (t *ThDesk)  SaveLotteryData() error {

	if t.DeskType == intCons.GAME_TYPE_TH {
		//自定义房间
		return t.SaveLotteryDatath()
	}else if t.DeskType == intCons.GAME_TYPE_TH_CS{
		//锦标赛
		return t.SaveLotteryDatacsth()
	}

	return nil

}


func (t *ThDesk) SaveLotteryDatath() error{
	log.T("一局游戏结束,开始保存游戏的数据到数据库")

	//为每一局保存一组数据
	deskRecord := &mode.T_th_desk_record{}
	deskRecord.DeskId = t.Id
	deskRecord.BeginTime = t.BeginTime
	deskRecord.UserIds = ""

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
		userService.FlashUser2Mongo(u.UserId)                        //刷新redis中的数据到mongo
		//2,保存游戏相关的数据
		//todo  游戏相关的数据结构 还没有建立,
		thData := &mode.T_th_record{}
		thData.Mid = bson.NewObjectId()
		thData.BetAmount = u.TotalBet
		thData.UserId = u.UserId
		thData.DeskId = u.deskId
		thData.WinAmount = u.winAmount - u.TotalBet
		thData.Blance = *(userService.GetUserById(u.UserId).Coin)
		thData.BeginTime = time.Now()
		thData.GameNumber = t.GameNumber
		db.SaveMgoData(casinoConf.DBT_T_TH_RECORD,thData)

		//获取游戏数据
		userRecord := mode.BeanRecord{}
		userRecord.UserId = u.UserId
		userRecord.NickName = u.NickName
		userRecord.WinAmount = u.winAmount - u.TotalBet

		deskRecord.Records = append(deskRecord.Records,userRecord)
		deskRecord.UserIds = strings.Join([]string{deskRecord.UserIds,u.NickName},",")
	}

	log.T("开始保存DBT_T_TH_DESK_RECORD的信息")
	//保存桌子的用户信息
	db.SaveMgoData(casinoConf.DBT_T_TH_DESK_RECORD,deskRecord)
	return nil

}


func (t *ThDesk) SaveLotteryDatacsth() error{
	log.T("一局游戏结束,开始保存游戏的数据到数据库")

	//为每一局保存一组数据
	deskRecord := &mode.T_th_desk_record{}
	deskRecord.DeskId = t.Id
	deskRecord.BeginTime = t.BeginTime
	deskRecord.UserIds = ""

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
		userService.FlashUser2Mongo(u.UserId)                        //刷新redis中的数据到mongo
		//2,保存游戏相关的数据
		//todo  游戏相关的数据结构 还没有建立,
		thData := &mode.T_th_record{}
		thData.Mid = bson.NewObjectId()
		thData.BetAmount = u.TotalBet
		thData.UserId = u.UserId
		thData.DeskId = u.deskId
		thData.WinAmount = u.winAmount - u.TotalBet
		thData.Blance = *(userService.GetUserById(u.UserId).Coin)
		thData.BeginTime = time.Now()
		thData.GameNumber = t.GameNumber
		db.SaveMgoData(casinoConf.DBT_T_TH_RECORD,thData)

		//获取游戏数据
		userRecord := mode.BeanRecord{}
		userRecord.UserId = u.UserId
		userRecord.NickName = u.NickName
		userRecord.WinAmount = u.winAmount - u.TotalBet

		deskRecord.Records = append(deskRecord.Records,userRecord)
		deskRecord.UserIds = strings.Join([]string{deskRecord.UserIds,u.NickName},",")
	}

	log.T("开始保存DBT_T_TH_DESK_RECORD的信息")
	//保存桌子的用户信息
	db.SaveMgoData(casinoConf.DBT_T_TH_DESK_RECORD,deskRecord)
	return nil

}

//得到这局胜利的人有几个
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


//新用户进入的时候,返回的信息
func (t *ThDesk) GetWeiXinInfos() []*bbproto.WeixinInfo{
	var result []*bbproto.WeixinInfo
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			ru :=userService.GetUserById(u.UserId)
			wxi := &bbproto.WeixinInfo{}
			wxi.HeadUrl = ru.HeadUrl
			wxi.NickName = ru.NickName
			wxi.OpenId = ru.OpenId

			//放在列表中
			result = append(result,wxi)
		}
	}
	return result
}

//跟注:跟注的时候 不需要重新设置押注的人
//只是跟注,需要减少用户的资产,增加奖池的金额
func (t *ThDesk) BetUserCall(user  *ThUser) error {
	log.T("用户[%v],nikename[%v],t.BetAmountNow[%v],user.HandCoin[%v]", user.UserId, user.NickName, t.BetAmountNow, user.HandCoin)
	followCoin := t.BetAmountNow - user.HandCoin
	if user.RoomCoin <= followCoin {
		//allin
		t.BetUserAllIn(user.UserId, user.RoomCoin)
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
	t.edgeJackpot += coin                        //边池 增加
	return nil
}

//如果弃牌的人是 t.NewRoundBetUser ,需要重新设置值
func (t *ThDesk) NextNewRoundBetUser() error {
	index := t.GetUserIndex(t.NewRoundFirstBetUser)
	for i := index + 1; i < len(t.Users) + index; i++ {
		u := t.Users[i % len(t.Users)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			t.NewRoundFirstBetUser = u.UserId
			break
		}
		//如果没有找到,那么返回失败
		if i == (len(t.Users) + index - 1) {
			log.T("没有找到下一个默认开始押注的人")
			return errors.New("没有找到下一个默认开始的押注的人")
		}
	}

	return nil

}

//让牌:只有第一个人才可以让牌
func (t *ThDesk) BetUserCheck(userId uint32) error {
	return nil
}

//用户加注
func (t *ThDesk) BetUserRaise(user *ThUser, coin int64) error {
	log.T("开始[%v],nickname[%v]操作用户加注[%v]的操作,user.RoomCoin[%v],user.handcoin[%v].t.BetAmountNow[%v]", user.UserId, user.NickName, coin, user.RoomCoin, user.HandCoin, t.BetAmountNow)

	/**
		判断allin的条件
		1,加注的金额coin 和受伤的余额一样多的情况	coin == user.RoomCoin
	 */

	if (coin > user.RoomCoin) {
		return errors.New("用户加注的金额,比自己手中的金额还要多,出错!")
	}

	//all-in的情况
	if coin == user.RoomCoin {
		//如果加注的金额和手上的余额一样大的话,all in
		if (coin + user.HandCoin) < t.BetAmountNow {
			//表示钱不够,被迫all-in,此时的betamountNow 不会改变
		} else {
			//主动设置下家需要跟注的金额
			t.MinRaise = user.HandCoin + coin - t.BetAmountNow
			t.BetAmountNow = user.HandCoin + coin
			t.RaiseUserId = user.UserId
		}
		//开始allin
		t.BetUserAllIn(user.UserId, user.RoomCoin)

	} else {
		t.MinRaise = user.HandCoin + coin - t.BetAmountNow
		//普通加注的情况
		t.BetAmountNow = user.HandCoin + coin
		//1,增加奖池的金额
		t.AddBetCoin(coin)                                //desk-coin
		//2,减少用户的金额
		t.caclUserCoin(user.UserId, coin)                        //thuser
		//3,设置状态:设置为第一个加注的人,如果后边所有人都是跟注,可由这个人判断一轮是否结束
		t.RaiseUserId = user.UserId
	}

	return nil
}

//用户AllIn
func (t *ThDesk) BetUserAllIn(userId uint32, coin int64) error {
	log.T("用户[%v]开始allin[%v]", userId, coin)
	//1,增加奖池的金额
	t.AddBetCoin(coin)

	//2,减少用户的金额
	t.caclUserCoin(userId, coin)

	//3,设置用户的状态
	t.GetUserByUserId(userId).Status = TH_USER_STATUS_ALLINING        //设置用户的状态为all-in
	if t.NewRoundFirstBetUser == userId {
		//如果用户是第一个押注的人,all-in之后第一个押注的是往下滑
		t.NextNewRoundBetUser()
		log.T("重新设置了t.NewRoundFirstBetUser[%v]", t.NewRoundFirstBetUser)
	}

	//4,增加all in的状态
	allinpot := &pokerService.AllInJackpot{}
	allinpot.UserId = userId
	allinpot.Jackpopt = 0
	allinpot.ThroundCount = t.RoundCount
	allinpot.AllInAmount = t.GetUserByUserId(userId).TotalBet4calcAllin
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

// 用户加注,跟住,allin 之后对他的各种余额属性进行计算
func (t *ThDesk) caclUserCoin(userId uint32, coin int64) error {
	user := t.GetUserByUserId(userId)
	user.TurnCoin += coin
	user.HandCoin += coin
	user.TotalBet4calcAllin += coin
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
	t.BetUserNow = 0        //这里设置为-1是为了方便判断找不到下一个人的时候,设置为新的一局
	for i := index; i < len(t.Users) + index; i++ {
		u := t.Users[(i + 1) % len(t.Users)]
		if u != nil && u.Status == TH_USER_STATUS_BETING {
			log.T("设置betUserNow 为[%v]", u.UserId)
			t.BetUserNow = u.UserId
			break
		}
	}

	log.T("判断是否是下一轮,计算出来的t.BetUserNow[%v],t.BetUserRaiseUserId[%v],t.NextNewRoundBetUser[%v]", t.BetUserNow, t.RaiseUserId, t.NewRoundFirstBetUser)
	//设置新一轮
	if t.BetUserNow == 0 || t.RaiseUserId == t.BetUserNow {
		//处理allin 奖金池分割的问题
		t.CalcAllInJackpot()
		t.RaiseUserId = t.NewRoundFirstBetUser
		t.BetUserNow = t.NewRoundFirstBetUser
		t.BetAmountNow = 0
		t.RoundCount ++
		t.MinRaise = t.BigBlindCoin        //第二轮开始的时候,最低加注金额设置喂大盲注

		////这里吧TurnCoin设置0 是为了是getMinRaise的值正确,目前turnCoin的作用只是在这里
		for i := 0; i < len(t.Users); i++ {
			u := t.Users[i]
			if u != nil {
				u.TurnCoin = 0
			}
		}

		log.T("设置下次押注的人是小盲注,下轮次[%v]", t.RoundCount)
	}

	//如果第一个押注的人,已经弃牌了,那么BetUserRaiseUserId 需要滑向下一个人
	if t.GetUserByUserId(t.RaiseUserId).Status == TH_USER_STATUS_FOLDED {
		log.T("第一个押注的人,弃牌了,需要把t.BetUserRaiseUserId 设置为下一个人")
		for i := t.GetUserIndex(t.RaiseUserId); i < len(t.Users) + index; i++ {
			u := t.Users[(i + 1) % len(t.Users)]
			if u != nil && u.Status == TH_USER_STATUS_BETING {
				log.T("设置betUserNow 为[%v]", u.UserId)
				t.RaiseUserId = u.UserId
				break
			}
		}
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
	*sendData.MinRaise = t.GetMinRaise()
	*sendData.NextSeat = t.GetUserByUserId(t.BetUserNow).Seat        //int32(t.GetUserIndex(t.BetUserNow))
	sendData.Handcoin = t.GetHandCoin()
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
	log.T("判断是否是新的一轮t.BetUserRaiseUserId[%v],t.BetUserNow(%v),t.status[%v].//status:1,stop,2,start,3,lottery", t.RaiseUserId, t.BetUserNow, t.Status)
	if t.RaiseUserId == t.BetUserNow &&  t.Status == TH_DESK_STATUS_SART {
		log.T("t.BetUserRaiseUserId[%v] == t.BetUserNow[%v],新的一局开始", t.RaiseUserId, t.BetUserNow)
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
		all := t.AllInJackpot[i]
		log.T("第[%v]次循环的时候,allinlist[%v]", i, all)

		//如果这个池子为nil ,则跳过这个循环,如果这个allIn 不是本轮的,则把之前的allin 的jackpot 累加起来
		if all == nil || all.ThroundCount != t.RoundCount {
			continue
		}

		log.T("开始计算用户[%v]allIn.index[%v] allin.amount[%v]计算all in 时的池子金额", all.UserId, i, all.AllInAmount)
		//每个allin计算金额
		for n := 0; n < len(t.Users); n++ {
			u := t.Users[n]
			if u != nil {
				log.T("用户[%v]押注的总金额是[%v]", u.UserId, u.TotalBet4calcAllin)
				if u.TotalBet4calcAllin > all.AllInAmount {
					all.Jackpopt += all.AllInAmount
					u.TotalBet4calcAllin -= all.AllInAmount
					log.T("用户[%v]押注加入all的金额是[%v]", u.UserId, all.AllInAmount)
				} else {
					all.Jackpopt += u.TotalBet4calcAllin
					u.TotalBet4calcAllin = 0
					log.T("用户[%v]押注加入all的金额是[%v]", u.UserId, u.TotalBet4calcAllin)
				}

			}
		}
		log.T("计算出来用户[%v]allIn.index[%v] allin.amount[%v]计算all in 的池子总金额[]", all.UserId, i, all.AllInAmount, all.Jackpopt)

		//之后的allinamount - 当前allin
		for k := i; k < len(t.AllInJackpot); k++ {
			allk := t.AllInJackpot[k]
			if allk != nil {
				allk.AllInAmount -= all.AllInAmount
			}
		}
		t.edgeJackpot -= all.Jackpopt
		log.T("开始给allIn[%v]计算all in 时的池子金额---------------------------------end---------------", i)
		log.T("目前t.bianJackPot 的剩余值是[%v]", t.edgeJackpot)
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
		log.E("用户userId[%v]name[%v]的带入金额小于0,所以不能押注或者投注了", user.UserId, user.NickName)
		return false
	}

	//用户合法,设置等待状态
	user.InitWait()
	return true
}

//是不是可以开始游戏了
func (t *ThDesk) IsTime2begin() bool {

	log.T("判断是否可以开始一局新的游戏")


	//1, 锦标赛开始的逻辑限制
	if t.DeskType == intCons.GAME_TYPE_TH_CS {
		//如果是锦标赛:1,房间的状态是正在游戏中
		if !ChampionshipRoom.CanNextDeskRun() {
			return false
		}
	}



	//2,通用的开始逻辑限制
	/**
		开始游戏的要求:
		1,[在线]用户的人数达到了最低可玩人数
		2,当前的状态是游戏停止的状态
	 */


	log.T("当前玩家的状态://1,等待开始,2,游戏中,3,押注中,5,allin,6,弃牌,7,等待结算,8,已经结算,9,裂开,10,掉线")
	for i := 0; i < len(t.Users); i++ {
		u := t.Users[i]
		if u != nil {
			log.T("用户[%v].seat[%v]的状态是[%v]", u.UserId, u.Seat, u.Status)
		}
	}

	//todo 金钱大于大盲注的人数必须要大于最低人数才可以玩
	log.T("当前在线玩家的数量是[%v],当前desk的状态是[%v],1未开始,2游戏中,3,开奖中", t.UserCountOnline, t.Status)

	if t.UserCountOnline >= TH_DESK_LEAST_START_USER  && t.Status == TH_DESK_STATUS_STOP{
		log.T("游戏到了开始的时候----begin----")
		return true
	} else {
		log.T("游戏还不到开始的时候")
		return false
	}
}

//开始游戏
func (mydesk *ThDesk) Run() error {
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
	initCardB.MinRaise = new(int64)

	//设置初始化值
	*initCardB.Tableid = int32(mydesk.Id)
	initCardB.HandCard = mydesk.GetHandCard()
	initCardB.PublicCard = mydesk.ThPublicCard2OGC()
	*initCardB.MinRaise = mydesk.GetMinRaise()
	*initCardB.NextUser = mydesk.GetUserByUserId(mydesk.BetUserNow).Seat                //	int32(mydesk.GetUserIndex(mydesk.BetUserNow))
	*initCardB.ActionTime = TH_TIMEOUT_DURATION_INT
	//initCardB.Seat = &mydesk.UserCount
	mydesk.THBroadcastProto(initCardB, 0)
	log.T("广播Game_InitCard的信息完毕")

	log.T("\n\n开始一局新的游戏,初始化完毕\n\n")
	return nil
}

//表示游戏结束
func (t *ThDesk) End(){
	//广播结算的信息
	result := &bbproto.Game_SendDeskEndLottery{}
	result.Result = &intCons.ACK_RESULT_SUCC

	for i:=0 ;i <len(t.Users); i++{
		u := t.Users[i]
		if u != nil {
			//
			gel := &bbproto.Game_EndLottery{}
			gel.Coin = new(int64)
			gel.BigWin = new(bool)
			gel.Owner = new(bool)
			gel.Rolename = new(string)

			*gel.Coin = u.winAmount

			if t.DeskOwner == u.UserId {
				*gel.Owner = true
			}else{
				*gel.Owner = false
			}

			*gel.Rolename = u.NickName
			result.CoinInfo = append(result.CoinInfo,gel)
		}
	}

	//设置bigWin

	//广播消息
	t.THBroadcastProtoAll(result)
}


//手牌转换为OG可以使用的牌
func (t *ThDesk) ThPublicCard2OGC() []*bbproto.Game_CardInfo {
	result := make([]*bbproto.Game_CardInfo, len(t.PublicPai))
	for i := 0; i < len(t.PublicPai); i++ {
		result[i] = ThCard2OGCard(t.PublicPai[i])
	}
	return result
}

//canrase
func (t *ThDesk) GetCanRise() int32 {
	if t.Tiem2Lottery() {
		return 0
	} else {
		return 1
	}
}