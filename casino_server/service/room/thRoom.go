package room

import (
	"casino_server/msg/bbprotogo"
	"sync"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"errors"
	"casino_server/service/userService"
	"time"
	"casino_server/mode"
	"casino_server/gamedata"
	"casino_server/utils/numUtils"
	"casino_server/utils"
	"casino_server/conf/intCons"
)
//config

var TH_GAME_SMALL_BLIND 	int64 = 10              //小盲注的金额
var GAME_THROOM_MAX_COUNT 	int32 = 500             //一个游戏大厅最多有多少桌德州扑克
var TH_TIMEOUT_DURATION = time.Second * 1500      	//德州出牌的超时时间
var TH_TIMEOUT_DURATION_INT 	int32 = 1500          	//德州出牌的超时时间
var TH_LOTTERY_DURATION = time.Second * 5         	//德州开奖的时间
var TH_DESK_CREATE_DIAMOND 	int64 = 10; 		//创建牌桌需要的钻石数量


//测试的时候 修改喂多人才可以游戏
var TH_DESK_LEAST_START_USER 	int32 = 2       	//最少多少人可以开始游戏
var TH_DESK_MAX_START_USER 	int32 = 8         	//玩德州扑克,每个房间最多多少人

//德州扑克 玩家的状态
var TH_USER_STATUS_WAITSEAT 	int32 = 1        	//刚上桌子 等待开始的玩家
var TH_USER_STATUS_SEATED 	int32 = 2          	//刚上桌子 但是没有在游戏中
var TH_USER_STATUS_READY 	int32 = 3
var TH_USER_STATUS_BETING 	int32 = 4          	//押注中
var TH_USER_STATUS_ALLINING 	int32 = 5        	//allIn
var TH_USER_STATUS_FOLDED 	int32 = 6          	//弃牌
var TH_USER_STATUS_WAIT_CLOSED 	int32 = 7    		 //等待结算
var TH_USER_STATUS_CLOSED 	int32 = 8          	//已经结算
var TH_USER_STATUS_LEAVE 	int32 = 9           	//离开  ????厉害是否有离开的状态


var TH_USER_BREAK_STATUS_TRUE 	int32 = 1      		//已经断线
var TH_USER_BREAK_STATUS_FALSE 	int32 = 0     		//没有断线


//德州扑克,牌桌的状态
var TH_DESK_STATUS_STOP 	int32 = 1            	//没有开始的状态
var TH_DESK_STATUS_SART 	int32 = 2            	//已经开始的状态
var TH_DESK_STATUS_LOTTERY 	int32 = 3         	//已经开始的状态

var TH_DESK_ROUND1 		int32 = 1        	//第一轮押注
var TH_DESK_ROUND2 		int32 = 2         	//第二轮押注
var TH_DESK_ROUND3 		int32 = 3         	//第三轮押注
var TH_DESK_ROUND4 		int32 = 4         	//第四轮押注
var TH_DESK_ROUND_END 		int32 = 5      		//完成押注


//押注的类型
var TH_DESK_BET_TYPE_BET 	int32 = 1        	//押注
var TH_DESK_BET_TYPE_CALL 	int32 = 2       	//跟注,和别人下相同的筹码
var TH_DESK_BET_TYPE_FOLD 	int32 = 3       	//弃牌
var TH_DESK_BET_TYPE_CHECK 	int32 = 4      		//让牌
var TH_DESK_BET_TYPE_RAISE 	int32 = 5      		//加注
var TH_DESK_BET_TYPE_RERRAISE 	int32 = 6   		//再加注
var TH_DESK_BET_TYPE_ALLIN 	int32 = 7      		//全下

//桌子的类型
var TH_DESK_TYPE_ZIDINGYI 	int32 = 1        	//自定义桌子
var TH_DESK_TYPE_JINBIAOSAI	int32 = 2        	//锦标赛的桌子

/**
	初始化函数:
		初始化游戏房间
 */

var ThGameRoomIns ThGameRoom        	//房间实例,在init函数中初始化


func init() {
	ThGameRoomIns.OnInit()                //初始化房间
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


//随机生成6位数字
func (r *ThGameRoom) RandRoomKey() string {
	a := utils.Rand(100000, 1000000)
	roomKey, _ := numUtils.Int2String(a)
	//1,判断roomKey是否已经存在
	if r.IsRoomKeyExist(roomKey) {
		log.E("房间密钥[%v]已经存在,创建房间失败,重新创建", roomKey)
		return r.RandRoomKey()
	} else {
		log.T("最终得到的密钥是[%v]", roomKey)
		return roomKey
	}
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
func (r *ThGameRoom) CreateDeskByUserIdAndRoomKey(userId uint32, roomCoin int64, roomkey string, smallBlind int64, bigBlind int64, jucount int32) (*ThDesk, error) {

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
	upDianmond := 0 - TH_DESK_CREATE_DIAMOND
	remainDiamond := userService.UpdateUserDiamond(userId, upDianmond)
	//3,生成一条交易记录
	err := userService.CreateDiamonDetail(userId, mode.T_USER_DIAMOND_DETAILS_TYPE_CREATEDESK, upDianmond, remainDiamond, "创建房间消耗钻石");
	if err != nil {
		log.E("创建用户的钻石交易记录失败")
		return nil, err
	}
	return desk, nil

}

//增加一个thRoom
func (r *ThGameRoom) AddThRoom(throom *ThDesk) error {
	r.ThDeskBuf = append(r.ThDeskBuf, throom)
	return nil
}


//删除一个throom
func (r *ThGameRoom) RmThroom(id int32) error {

	//第一步找到index
	var index int = -1
	for i := 0; i < len(r.ThDeskBuf); i++ {
		desk := r.ThDeskBuf[i]
		if desk != nil && desk.Id == id {
			index = i
			break
		}
	}

	//判断是否找到对应的desk
	if index == -1 {
		log.E("没有找到对应desk.id[%v]的桌子", id)
		return errors.New("没有找到对应的desk")
	}

	//删除对应的desk
	r.ThDeskBuf = append(r.ThDeskBuf[:index], r.ThDeskBuf[index + 1:]...)
	return nil
}


//通过房主id解散房间
func (r *ThGameRoom) DissolveDeskByDeskOwner(userId uint32, a gate.Agent) error {

	result := &bbproto.Game_AckDissolveDesk{}
	result.Result = new(int32)
	result.UserId = new(uint32)
	result.DeskId = new(int32)
	result.PassWord = new(string)

	//1,找到桌子
	desk := r.GetDeskByDeskOwner(userId)        //
	//2,解散桌子的条件,如果正在游戏中,是否能解散?
	if desk.Status != TH_DESK_STATUS_STOP {
		*result.Result = intCons.ACK_RESULT_ERROR
		a.WriteMsg(result)
		return errors.New("游戏正在进行中,不能解散")
	}

	//3,发送解散的广播

	*result.Result = intCons.ACK_RESULT_SUCC
	*result.UserId = desk.deskOwner
	*result.PassWord = desk.RoomKey

	desk.THBroadcastProtoAll(result)

	//4,解散
	r.RmThroom(desk.Id)
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

//通过房主来找到房间
func (r *ThGameRoom) GetDeskByDeskOwner(userId uint32) *ThDesk {
	for i := 0; i < len(ThGameRoomIns.ThDeskBuf); i++ {
		desk := ThGameRoomIns.ThDeskBuf[i]
		if desk != nil && desk.deskOwner == userId {
			return desk
		}
	}
	return nil
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
	log.T("玩家[%v]通过roomkey[%v]进入房间", userId, roomKey)
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
	err := mydesk.AddThUser(userId, roomCoin, TH_USER_STATUS_SEATED, a)
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

	//3,进入房间,竞标赛进入房间的时候,默认就是准备的状态
	err := mydesk.AddThUser(userId, roomCoin, TH_USER_STATUS_READY, a)
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

