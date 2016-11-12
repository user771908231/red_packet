package majiang

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_majiang/service/AgentService"
	"casino_server/conf/intCons"
	"time"
	"casino_majiang/conf/config"
	"casino_server/utils/db"
	"strings"
	"sync/atomic"
	"casino_majiang/gamedata/model"
	"casino_server/utils/numUtils"
	"casino_server/utils/timeUtils"
	"casino_majiang/service/lock"
	"casino_server/utils"
)

//状态表示的是当前状态.
var MJDESK_STATUS_CREATED int32 = 1 //刚刚创建
var MJDESK_STATUS_READY int32 = 2//正在准备
var MJDESK_STATUS_ONINIT int32 = 3//准备完成之后，desk初始化数据
var MJDESK_STATUS_FAPAI int32 = 4//发牌的阶段
var MJDESK_STATUS_EXCHANGE int32 = 5//desk初始化完成之后，告诉玩家可以开始换牌
var MJDESK_STATUS_DINGQUE int32 = 6//换牌结束之后，告诉玩家可以开始定缺
var MJDESK_STATUS_RUNNING int32 = 7 //定缺之后，开始打牌
var MJDESK_STATUS_LOTTERY int32 = 8 //结算
var MJDESK_STATUS_END int32 = 9//一局结束


var OVER_TURN_ACTTYPE_MOPAI int32 = 1; //摸牌的类型...
var OVER_TURN_ACTTYPE_OTHER int32 = 2; //碰OTHER

var MJDESK_ACT_TYPE_MOPAI int32 = 1; ///摸牌
var MJDESK_ACT_TYPE_DAPAI int32 = 2; //打牌
var MJDESK_ACT_TYPE_WAIT_CHECK int32 = 3; //等待check


var DINGQUE_SLEEP_DURATION time.Duration = time.Second * 5        //定缺的延迟
var SHAIZI_SLEEP_DURATION time.Duration = time.Second * 4        //定缺的延迟


//判断是不是朋友桌
func (d *MjDesk) IsFriend() bool {
	return true
}

//朋友桌用户加入房间
/**
return  reconnect,error
 */
func (d *MjDesk) addNewUserFriend(userId uint32, a gate.Agent) (bool, error) {

	// 设置agent
	AgentService.SetAgent(userId, a)

	//1,是否是重新进入
	//user := d.GetUserByUserId(userId)
	//if user != nil {
	//	//是断线重连
	//	*user.IsBreak = false;
	//	return nil
	//}

	//2,是否是离开之后重新进入房间
	userLeave := d.GetUserByUserId(userId)
	if userLeave != nil {
		log.T("玩家[%v]断线重连....", userId)
		*userLeave.IsBreak = false
		*userLeave.IsLeave = false
		return true, nil
	}

	//3,加入一个新用户
	newUser := NewMjUser()
	*newUser.UserId = userId
	*newUser.DeskId = d.GetDeskId()
	*newUser.RoomId = d.GetRoomId()
	*newUser.Coin = d.GetBaseValue()
	*newUser.IsBreak = false
	*newUser.IsLeave = false
	*newUser.IsBanker = false
	*newUser.Status = MJUSER_STATUS_INTOROOM
	newUser.GameData = NewPlayerGameData()

	//设置agent
	err := d.addUser(newUser)
	if err != nil {
		log.E("用户[%v]加入房间[%v]失败,errMsg[%v]", userId, err)
		return false, errors.New("用户加入房间失败")
	} else {
		//加入房间成功
		return false, nil
	}
}

//发送重新连接之后的overTurn
func (d *MjDesk) SendReconnectOverTurn(userId uint32) error {
	log.T("开始处理 sendReconnectOverTurn(%v),当前desk.status(%v),", userId, d.GetStatus())

	//得到玩家
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.T("发送SendReconnectOverTurn(%v)失败,因为没有找到玩家", userId)
		return errors.New("发送SendReconnectOverTurn()失败,因为没有找到玩家")
	}

	//开发发送overTurn
	if d.IsPreparing() && user.IsNotReady() {
		//给玩家发送开始准备的消息...
		log.T("sendReconnectOverTurn，给user[%v]发送准备的消息....", userId)
	} else if d.IsExchange() {
		//给玩家发送换牌的消息...
		log.T("sendReconnectOverTurn，给user[%v]发送换牌的消息....", userId)

	} else if d.IsDingQue() && user.IsNotDingQue() {
		log.T("sendReconnectOverTurn，给user[%v]发送定缺的消息....", userId)

		//给玩家发送定缺的信息...
		beginQue := newProto.NewGame_BroadcastBeginDingQue()
		user.WriteMsg(beginQue)

	} else if d.IsGaming() && user.GetUserId() == d.GetActUser() {

		//游戏中的情况，发送act的消息,这里需要更具当前的状态来发送overTurn
		if d.GetActType() == MJDESK_ACT_TYPE_MOPAI {
			log.T("sendReconnectOverTurn，给user[%v]发送摸牌的消息....", userId)
			overTrun := d.GetMoPaiOverTurn(user, false)                        //重新进入房间之后
			user.WriteMsg(overTrun)
			log.T("玩家重新进入游戏之后 [%v]开始摸牌【%v】...", user.GetUserId(), overTrun)

		} else if d.GetActType() == MJDESK_ACT_TYPE_DAPAI {
			//发送打牌的协议
			log.T("sendReconnectOverTurn，给user[%]发送打牌的消息....", userId)

		} else if d.GetActType() == MJDESK_ACT_TYPE_WAIT_CHECK {
			log.T("sendReconnectOverTurn，给user[%]发送checkCase的消息....", userId)

			caseBean := d.CheckCase.GetBeanByUserIdAndStatus(user.GetUserId(), CHECK_CASE_BEAN_STATUS_CHECKING)
			if caseBean == nil {
				log.E("没有找到玩家[%v]对应的checkBean", user.GetUserId())
				return errors.New("玩家重新进入房间发送check overturn的时候出错")
			}

			//找到需要判断bean之后，发送给判断人	//send overTurn
			overTurn := d.GetOverTurnByCaseBean(d.CheckCase.CheckMJPai, caseBean, OVER_TURN_ACTTYPE_OTHER)        //重新进入游戏
			*overTurn.Time = int32(user.GetWaitTime() - time.Now().Unix())
			///发送overTurn 的信息
			log.T("开始发送玩家[%v]断线重连的overTurn[%v]", user.GetUserId(), overTurn)
			user.WriteMsg(overTurn)
		}
	}

	log.T("开始处理 sendReconnectOverTurn(%v),当前desk.status(%v)----处理完毕...", userId, d.GetStatus())
	return nil
}


//删除已经离开的人的Id todo
func (d *MjDesk) rmLeaveUser(userId uint32) error {
	return nil
}

//通过userId得到user
func (d *MjDesk) GetUserByUserId(userId uint32) *MjUser {
	for _, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == userId {
			return u
		}
	}

	return nil
}

//通过userId得到user
func (d *MjDesk) getLeaveUserByUserId(userId uint32) *MjUser {
	for _, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == userId {
			return u
		}
	}
	return nil
}


//新增加一个玩家
func (d *MjDesk) addUser(user *MjUser) error {
	//找到座位
	seatIndex := -1
	for i, u := range d.Users {
		if u == nil {
			seatIndex = i
			break
		}
	}

	//如果找到座位那么，增加用户，否则返回错误信息
	if seatIndex >= 0 {
		d.Users[seatIndex] = user
		return nil
	} else {
		return errors.New("没有找到合适的位置，加入桌子失败")
	}
}

func (d *MjDesk) GetRoomTypeInfo() *mjproto.RoomTypeInfo {
	typeInfo := newProto.NewRoomTypeInfo()
	*typeInfo.Settlement = d.GetSettlement()
	typeInfo.PlayOptions = d.GetPlayOptions()
	*typeInfo.MjRoomType = mjproto.MJRoomType(d.GetMjRoomType())
	*typeInfo.BaseValue = d.GetBaseValue()
	*typeInfo.BoardsCout = d.GetBoardsCout()
	*typeInfo.CapMax = d.GetCapMax()
	*typeInfo.CardsNum = d.GetCardsNum()
	return typeInfo
}

func (d *MjDesk) GetPlayOptions() *mjproto.PlayOptions {
	o := newProto.NewPlayOptions()
	*o.ZiMoRadio = d.GetZiMoRadio()
	*o.HuRadio = d.GetHuRadio()
	*o.DianGangHuaRadio = d.GetDianGangHuaRadio()
	o.OthersCheckBox = d.GetOthersCheckBox()
	//log.T("回复的时候回复的othersCheckBox[%v]", o.OthersCheckBox)
	return o
}

//广播协议
func (d *MjDesk) BroadCastProto(p proto.Message) error {
	for _, u := range d.Users {
		//user不为空，并且user没有离开，没有短线的时候才能发送消息...
		if u != nil && !u.GetIsBreak() && !u.GetIsLeave() {
			u.WriteMsg(p)
		}
	}
	return nil
}

// 广播 但是不好办 userId
func (d *MjDesk) BroadCastProtoExclusive(p proto.Message, userId uint32) error {
	for _, u := range d.Users {
		if u != nil && u.GetUserId() != userId {
			u.WriteMsg(p)
		}
	}
	return nil
}

//得到deskGameInfo
func (d *MjDesk) GetDeskGameInfo() *mjproto.DeskGameInfo {
	deskInfo := newProto.NewDeskGameInfo()
	//deskInfo.ActionTime
	//deskInfo.DelayTime
	*deskInfo.GameStatus = d.GetClientGameStatus()
	*deskInfo.CurrPlayCount = d.GetCurrPlayCount() //当前第几局
	*deskInfo.TotalPlayCount = d.GetTotalPlayCount()//总共几局
	*deskInfo.PlayerNum = d.GetPlayerNum()        //玩家的人数
	deskInfo.RoomTypeInfo = d.GetRoomTypeInfo()
	*deskInfo.RoomNumber = d.GetPassword()        //房间号码...
	*deskInfo.Banker = d.GetBanker()
	//deskInfo.NRebuyCount
	//deskInfo.InitRoomCoin
	*deskInfo.NInitActionTime = d.GetNInitActionTime()
	//deskInfo.NInitDelayTime
	*deskInfo.ActiveUserId = d.GetActiveUser()
	*deskInfo.RemainCards = d.GetRemainPaiCount()
	return deskInfo
}


/**

MJDESK_STATUS_CREATED = 1 //刚刚创建
MJDESK_STATUS_READY = 2//正在准备
MJDESK_STATUS_ONINIT = 3//准备完成之后，desk初始化数据
MJDESK_STATUS_EXCHANGE = 4//desk初始化完成之后，告诉玩家可以开始换牌
MJDESK_STATUS_DINGQUE = 5//换牌结束之后，告诉玩家可以开始定缺
MJDESK_STATUS_RUNNING = 6 //定缺之后，开始打牌
MJDESK_STATUS_LOTTERY = 7 //结算
MJDESK_STATUS_END = 8//一局结束
 */

func (d *MjDesk) GetClientGameStatus() int32 {
	var gameStatus mjproto.DeskGameStatus = mjproto.DeskGameStatus_INIT//默认状态
	switch d.GetStatus() {
	case MJDESK_STATUS_CREATED:
		gameStatus = mjproto.DeskGameStatus_INIT
	case MJDESK_STATUS_READY:
		gameStatus = mjproto.DeskGameStatus_INIT
	case MJDESK_STATUS_ONINIT:
		gameStatus = mjproto.DeskGameStatus_INIT
	case MJDESK_STATUS_FAPAI:
		gameStatus = mjproto.DeskGameStatus_FAPAI
	case MJDESK_STATUS_EXCHANGE:
		gameStatus = mjproto.DeskGameStatus_EXCHANGE
	case MJDESK_STATUS_DINGQUE:
		gameStatus = mjproto.DeskGameStatus_DINGQUE
	case MJDESK_STATUS_RUNNING:
		gameStatus = mjproto.DeskGameStatus_PLAYING
	case MJDESK_STATUS_LOTTERY:
		gameStatus = mjproto.DeskGameStatus_FINISH
	case MJDESK_STATUS_END:
		gameStatus = mjproto.DeskGameStatus_FINISH
	}
	return int32(gameStatus)
}
//返回玩家的数目
func (d *MjDesk) GetPlayerInfo(receiveUserId uint32) []*mjproto.PlayerInfo {
	var players []*mjproto.PlayerInfo
	for _, user := range d.Users {
		if user != nil {
			showHand := (user.GetUserId() == receiveUserId)                //是否需要显示手牌
			isOwner := ( d.GetOwner() == user.GetUserId())                //判断是否是房主
			info := user.GetPlayerInfo(showHand)
			*info.IsOwner = isOwner
			players = append(players, info)
		}
	}
	return players
}

//玩家的人数
func (d *MjDesk) GetPlayerNum() int32 {
	var count int32 = 0
	for _, user := range d.Users {
		if user != nil {
			count ++
		}
	}
	return count
}

// 发送gameInfo的信息
func (d *MjDesk) GetGame_SendGameInfo(receiveUserId uint32) *mjproto.Game_SendGameInfo {
	gameInfo := newProto.NewGame_SendGameInfo()
	gameInfo.DeskGameInfo = d.GetDeskGameInfo()
	*gameInfo.SenderUserId = receiveUserId
	gameInfo.PlayerInfo = d.GetPlayerInfo(receiveUserId)
	return gameInfo
}


//用户准备
func (d *MjDesk) Ready(userId  uint32) error {
	//找到需要准备的user
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("用户[%v]在desk[%v]准备的时候失败,没有找到对应的玩家", userId, d.GetDeskId())
		return errors.New("没有找到用户，准备失败")
	}

	if user.IsReady() {
		log.E("玩家[%v]已经准备好了...请不要重新准备...", userId)
		return errors.New("玩家已经准备了，请不要重复准备...")
	}

	//设置为准备的状态
	user.SetStatus(MJUSER_STATUS_READY)
	*user.Ready = true

	return nil
}

//是不是所有人都准备
func (d *MjDesk) IsAllReady() bool {
	for _, u := range d.Users {
		if u != nil && !u.IsReady() {
			return false
		}
	}
	return true
}

//这里表示 是否是 [正在] 准备中...
func (d *MjDesk) IsPreparing() bool {
	if d.GetStatus() == MJDESK_STATUS_READY {
		return true
	} else {
		return false
	}
}

func (d *MjDesk) IsNotPreparing() bool {
	return !d.IsPreparing()
}

//是否在定缺中
func (d *MjDesk) IsDingQue() bool {
	if d.GetStatus() == MJDESK_STATUS_DINGQUE {
		return true
	} else {
		return false
	}
}

func (d *MjDesk) IsNotDingQue() bool {
	return !d.IsDingQue()
}

//是否处于换牌的阶段
func (d *MjDesk) IsExchange() bool {
	if d.GetStatus() == MJDESK_STATUS_EXCHANGE {
		return true
	} else {
		return false
	}
}

//是否已经开始游戏了...
func (d *MjDesk) IsGaming() bool {
	if d.GetStatus() == MJDESK_STATUS_RUNNING {
		return true
	} else {
		return false
	}
}

func (d *MjDesk) IsNotGaming() bool {
	return !d.IsGaming()
}

//得到当前桌子的人数..
func (d *MjDesk) GetUserCount() int32 {
	var count int32 = 0
	for _, user := range d.Users {
		if user != nil {
			count ++
		}
	}
	//log.T("当前桌子的玩家数量是count[%v]", count)
	return count;

}

//玩家是否足够
func (d *MjDesk) IsPlayerEnough() bool {
	if d.GetUserCount() == 4 {
		return true
	} else {
		return false;
	}
}

//用户准备之后的一些操作
func (d *MjDesk) AfterReady() error {
	//如果所有人都准备了，那么开始游戏
	d.begin()
	return nil
}

//开始游戏
/**
开始游戏需要有几个步骤
1，desk的状态是否正确，现在是否可以开始游戏


 */
func (d *MjDesk) begin() error {
	lock := lock.GetDeskLock(d.GetDeskId())
	lock.Lock()
	defer lock.Unlock()

	//1，检查是否可以开始游戏
	err := d.time2begin()
	if err != nil {
		//log.T("无法开始游戏:err[%v]", err)
		return err
	}

	//2，初始化桌子的状态
	d.beginInit()

	//3，发13张牌
	err = d.initCards()
	if err != nil {
		log.E("初始化牌的时候出错err[%v]", err)
		return err
	}

	//这里需要判断
	if d.IsNeedExchange3zhang() {
		//开始换三张
		err = d.beginExchange()
		if err != nil {
			log.E("发送开始换三张的广播的时候出错err[%v]", err)
			return err
		}
	} else {
		//不用换三张 //直接法开始定缺的广播
		err := d.beginDingQue()
		if err != nil {
			log.E("开始发送定缺广播的时候出错err[%v]", err)
			return err
		}
	}

	return nil
}

//是否需要换三张
func (d *MjDesk) IsNeedExchange3zhang() bool {
	return d.IsOpenOption(mjproto.MJOption_EXCHANGE_CARDS)
}

//是否需要天地胡
func (d *MjDesk) IsNeedTianDiHu() bool {
	return d.IsOpenOption(mjproto.MJOption_TIAN_DI_HU)
}

//是否需要幺九将对
func (d *MjDesk) IsNeedYaojiuJiangdui() bool {
	return d.IsOpenOption(mjproto.MJOption_YAOJIU_JIANGDUI)
}

//是否需要门清中张
func (d *MjDesk) IsNeedMenqingZhongzhang() bool {
	return d.IsOpenOption(mjproto.MJOption_MENQING_MID_CARD)
}

//是否需要自摸加底
func (d *MjDesk) IsNeedZiMoJiaDi() bool {
	if mjproto.MJOption(*d.GetRoomTypeInfo().GetPlayOptions().ZiMoRadio) == mjproto.MJOption_ZIMO_JIA_DI {
		return true
	}
	return false
}

//是否需要自摸加番
func (d *MjDesk) IsNeedZiMoJiaFan() bool {
	if mjproto.MJOption(*d.GetRoomTypeInfo().GetPlayOptions().ZiMoRadio) == mjproto.MJOption_ZIMO_JIA_FAN {
		return true
	}
	return false
}

//是否可以开始
func (d *MjDesk) time2begin() error {
	log.T("检测游戏是否可以开始... ")
	if d.IsAllReady() && d.IsPlayerEnough() && d.IsNotDingQue() {
		return nil
	} else {
		return errors.New("开始游戏失败，因为还有人没有准备")
	}
	return nil
}


/**
1,初始化desk
2,初始化user
 */
func (d *MjDesk) beginInit() error {
	//初始化桌子的信息
	//1,初始化庄的信息,如果目前没有庄，则设置房主为庄,如果有庄，则不用管，每局游戏借宿的时候，会设置下一局的庄
	if d.GetBanker() == 0 {
		*d.Banker = d.GetOwner()
	}

	d.AddCurrPlayCount()        //场次数目加一
	d.SetActiveUser(d.GetBanker())        //设置当前的活动玩家
	*d.GameNumber, _ = db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)        //设置游戏编号
	*d.BeginTime = timeUtils.Format(time.Now())

	//初始化每个玩家的信息
	for _, user := range d.GetUsers() {
		if user != nil && user.CanBegin() {
			user.BeginInit(d.GetCurrPlayCount(), d.GetBanker())
		}
	}
	//发送游戏开始的协议...
	log.T("发送游戏开始的协议..")
	open := newProto.NewGame_Opening()
	d.BroadCastProto(open)
	time.Sleep(SHAIZI_SLEEP_DURATION)
	return nil
}

func (d *MjDesk) AddCurrPlayCount() {
	atomic.AddInt32(d.CurrPlayCount, 1)
}


/**
	初始化牌相关的信息
 */
func (d *MjDesk) initCards() error {
	//得到一副已经洗好的麻将
	d.SetStatus(MJDESK_STATUS_FAPAI)        //发牌的阶段
	d.AllMJPai = XiPai()
	//d.AllMJPai = XiPaiTestHu()
	//给每个人初始化...
	for i, u := range d.Users {
		if u != nil && u.IsReady() {
			//log.T("开始给你玩家[%v]初始化手牌...", u.GetUserId())
			ps := make([]*MJPai, 13)
			copy(ps, d.AllMJPai[i * 13: (i + 1) * 13])                //这里这样做的目的是不能更改base的值
			u.GameData.HandPai.Pais = ps
			*d.MJPaiCursor = int32((i + 1) * 13) - 1;
		}
	}

	//庄需要多发一张牌
	bankUser := d.GetBankerUser()
	bankUser.GameData.HandPai.InPai = d.GetNextPai()

	//发牌的协议game_DealCards  初始化完成之后，给每个人发送牌
	for _, user := range d.Users {
		if user != nil {
			dealCards := d.GetDealCards(user)
			if dealCards != nil {
				log.T("把玩家[%v]的牌[%v]发送到客户端...", user.GetUserId(), dealCards)
				user.WriteMsg(dealCards)
			} else {
				log.E("给user[%v]发牌的时候出现错误..", user.GetUserId())
			}
		}
	}

	//发送发牌的广播
	return nil
}

//发牌的协议
func (d *MjDesk) GetDealCards(user *MjUser) *mjproto.Game_DealCards {
	dealCards := newProto.NewGame_DealCards()
	*dealCards.DealerUserId = d.GetBanker()
	for _, u := range d.GetUsers() {
		if u != nil {
			if u.GetUserId() == user.GetUserId() {
				//表示是自己，可以看到手牌
				pc := u.GetPlayerCard(true)
				if d.GetBanker() == u.GetUserId() {
					pc.HandCard = append(pc.HandCard, u.GameData.HandPai.InPai.GetCardInfo())
				}

				dealCards.PlayerCard = append(dealCards.PlayerCard, pc)
			} else {
				pc := u.GetPlayerCard(false)                                //表示不是自己，不能看到手牌

				if d.GetBanker() == u.GetUserId() {
					pc.HandCard = append(pc.HandCard, NewBackPai())
				}

				dealCards.PlayerCard = append(dealCards.PlayerCard, pc)
			}

		}

	}

	return dealCards
}

func (d *MjDesk) SetStatus(status int32) {
	*d.Status = status
}

//设置当前用户的status
func (d *MjDesk) UpdateUserStatus(status int32) {
	for _, user := range d.GetUsers() {
		if user != nil {
			user.SetStatus(status)
		}
	}

}

//开始换三张
func (d *MjDesk) beginExchange() error {
	time.Sleep(time.Second * 3)
	data := newProto.NewGame_BroadcastBeginExchange()
	d.BroadCastProto(data)
	return nil
}



//开始定缺
func (d *MjDesk) beginDingQue() error {
	//开始定缺，修改desk的状态
	d.SetStatus(MJDESK_STATUS_DINGQUE)


	//给每个人发送开始定缺的信息
	beginQue := newProto.NewGame_BroadcastBeginDingQue()
	log.T("开始给玩家发送开始定缺的广播[%v]", beginQue)
	//
	time.Sleep(DINGQUE_SLEEP_DURATION)
	d.BroadCastProto(beginQue)
	return nil
}

//把桌子的数据保存到redis中
/**
	需要调用的地方
	1,新增加一个桌子的时候
	2,
 */
func (d *MjDesk)updateRedis() error {
	err := UpdateMjDeskRedis(d)        //保存数据到redis
	if err != nil {
		return err
	} else {
		return nil
	}
}

//个人开始定缺
func (d *MjDesk) DingQue(userId uint32, color int32) error {
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("定缺的时候，服务器出现错误，没有找到对应的user【%v】", userId)
		return errors.New("没有找到用户，定缺失败")
	}

	if user.IsDingQue() {
		log.E("玩家[%v]重复定缺.", userId)
		return errors.New("用户已经定缺了，重复定缺....")
	}

	//设置定缺
	*user.DingQue = true
	*user.GameData.HandPai.QueFlower = color
	user.SetStatus(MJUSER_STATUS_DINGQUE)        //设置目前的状态是已经定缺

	//回复定缺成功的消息
	ack := newProto.NewGame_DingQue()
	*ack.Header.UserId = userId
	*ack.Color = -1
	*ack.UserId = userId
	log.T("回复定缺的协议ack[%v]", ack)
	d.BroadCastProto(ack)        //发送定缺成功的广播
	return nil
}

//是不是全部都定缺了
func (d *MjDesk) AllDingQue() bool {
	for _, user := range d.GetUsers() {
		if user != nil && !user.IsDingQue() {
			log.T("用户[%v]还没有缺牌，等待定缺之后庄家开始打牌...", user.GetUserId())
			return false
		}
	}
	return true
}

func (d *MjDesk) GetBankerUser() *MjUser {
	return d.GetUserByUserId(d.GetBanker())
}

//初始化checkCase
//如果出错 设置checkCase为nil
func (d *MjDesk) InitCheckCase(p *MJPai, outUser *MjUser) error {

	checkCase := NewCheckCase()
	*checkCase.UserIdOut = outUser.GetUserId()
	*checkCase.CheckStatus = CHECK_CASE_STATUS_CHECKING        //正在判定
	checkCase.CheckMJPai = p
	checkCase.PreOutGangInfo = outUser.GetPreMoGangInfo()

	//初始化checkbean
	for _, checkUser := range d.GetUsers() {
		//这里要判断用户是不是已经胡牌
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() {

			log.T("用户[%v]打牌，判断user[%v]是否可以碰杠胡.手牌[%v]", outUser.GetUserId(), checkUser.GetUserId(), checkUser.GameData.HandPai.GetDes())
			checkUser.GameData.HandPai.InPai = p
			//添加checkBean
			bean := checkUser.GetCheckBean(p, d.IsXueLiuChengHe())
			if bean != nil {
				checkCase.CheckB = append(checkCase.CheckB, bean)
			}
		}
	}

	log.T("判断最终的checkCase[%v]", checkCase)
	if checkCase.CheckB != nil || len(checkCase.CheckB) > 0 {
		d.CheckCase = checkCase
	} else {
		d.CheckCase = nil
	}

	return nil
}

//判断是否可以initCheckCas

////暂时不用？？ 摸牌之后
//func (d *MjDesk) InitMoPaiCheckCase(p *MJPai, moPaiUser *MjUser) error {
//
//	//初始化参数
//	moPaiUser.GameData.HandPai.InPai = p
//
//	//判断可能性
//	checkCase := NewCheckCase()
//	*checkCase.UserIdOut = moPaiUser.GetUserId()
//	*checkCase.CheckStatus = CHECK_CASE_STATUS_CHECKING        //正在判定
//	checkCase.CheckMJPai = p
//
//	checkCase.PreOutGangInfo = moPaiUser.GetPreMoGangInfo()
//	checkCase.CheckB = append(checkCase.CheckB, moPaiUser.GetCheckBean(p))
//	if checkCase.CheckB == nil || len(checkCase.CheckB) > 0 {
//		d.CheckCase = checkCase
//	} else {
//		d.CheckCase = nil
//	}
//
//	return nil
//}


//执行判断事件
/**

	这里仅仅是用于判断打牌之后别人的碰杠胡

 */
func (d *MjDesk) DoCheckCase(gangUser *MjUser) error {
	//检测参数
	if d.CheckCase == nil || d.CheckCase.GetNextBean() == nil {
		log.T("已经没有需要处理的CheckCase,下一个玩家摸牌...")
		//直接跳转到下一个操作的玩家...,这里表示判断已经玩了...
		d.CheckCase = nil
		//在这之前需要保证 activeUser 是正确的...
		d.SendMopaiOverTurn(gangUser)
		return nil
	} else {
		log.T("继续处理CheckCase,开处理下一个checkBean...")
		//1,找到胡牌的人来进行处理
		caseBean := d.CheckCase.GetNextBean()
		//找到需要判断bean之后，发送给判断人	//send overTurn
		overTurn := d.GetOverTurnByCaseBean(d.CheckCase.CheckMJPai, caseBean, OVER_TURN_ACTTYPE_OTHER)        //别人打牌，判断是否可以碰杠胡

		///发送overTurn 的信息
		log.T("开始发送overTurn[%v]", overTurn)
		d.GetUserByUserId(caseBean.GetUserId()).SendOverTurn(overTurn)
		d.SetActUserAndType(caseBean.GetUserId(), MJDESK_ACT_TYPE_WAIT_CHECK)        //设置当前活动的玩家

		return nil

	}

}

//得到麻将牌的总张数
func (d *MjDesk) GetTotalMjPaiCount() int32 {
	return 108; //暂时返回108张

}


/**
	1，只剩一个玩家没有胡牌
	2, 已经没有牌了...
 */

func (d *MjDesk) Time2Lottery() bool {
	//游戏中的玩家只剩下一个人，表示游戏结束...
	gamingCount := d.GetGamingCount()        //正在游戏中的玩家数量

	log.T("判断是否胡牌...当前的gamingCount[%v],当前的PaiCursor[%v]", gamingCount, d.GetMJPaiCursor())

	//1,只剩下一个人的时候. 表示游戏结束
	if gamingCount == 1 {
		return true
	}


	//2,当牌已经被抹完的时候，表示游戏结束
	if d.GetMJPaiCursor() == (d.GetTotalMjPaiCount() - 1) {
		return true;
	}

	return false
}

func (d *MjDesk) GetGamingCount() int32 {
	var gamingCount int32 = 0        //正在游戏中的玩家数量
	for _, user := range d.GetUsers() {
		if user != nil && user.IsGaming() {
			gamingCount ++
		}
	}
	return gamingCount
}


// 一盘麻将结束....这里需要针对每个人结账...并且对desk和user的数据做清楚...
func (d *MjDesk) Lottery() error {
	//结账需要分两中情况
	/**
		1，只剩一个玩家没有胡牌的时候
		2，没有生育麻将的时候.需要分别做处理...
	 */

	//判断是否可以胡牌
	if !d.Time2Lottery() {
		return errors.New("没有到lottery()的时间...")
	}

	log.T("现在开始处理lottery()的逻辑....")


	//查花猪
	d.ChaHuaZhu()
	//查大叫
	d.ChaDaJiao()
	//
	d.DoLottery()

	//发送结束的广播
	d.SendLotteryData()

	//
	d.AfterLottery()

	//判断牌局结束(整场游戏结束)
	if !d.End() {
		//go d.begin()
	}

	return nil
}

func (d *MjDesk) GetUserIds() string {
	ids := ""
	for _, user := range d.GetUsers() {
		if user != nil {
			idStr, _ := numUtils.Uint2String(user.GetUserId())
			ids = ids + "," + idStr
		}

	}
	return ids

}

//查花猪
/**
	查花猪是查用户是否没有缺
 */
func (d *MjDesk) ChaHuaZhu() error {
	for _, u := range d.GetUsers() {
		if u != nil && u.IsNotHu() {
			//开对用户查花猪
			if u.IsHuaZhu() {
				log.T("玩家[%v]是花猪", u.GetUserId())
				d.DoHuaZhu(u)
			}
		}
	}
	return nil
}


//花猪玩家需要给每一个非花猪8倍分
//todo 花猪的情况比较少见，所以可以先不用实现..
func (d *MjDesk) DoHuaZhu(huazhu *MjUser) error {
	log.T("开始处理花猪[%v]", huazhu.GetUserId())
	for _, user := range d.GetUsers() {
		if user != nil && user.IsNotHuaZhu() {
			//判断不是花猪，可以赢钱...

		}

	}

	return nil
}




//查大叫
/**
	查用户有没有叫
 */
func (d *MjDesk) ChaDaJiao() error {
	for _, u := range d.GetUsers() {
		if u != nil && u.IsNotHu() && u.IsNotHuaZhu() {
			//开对用户查花猪
			if !u.ChaJiao() {
				log.T("玩家[%v]没叫", u.GetUserId())
				d.DoDaJiao(u)
			}
		}
	}
	return nil

}

//得到一个canhuinfos
/**
	一次判断打出每一张牌的时候，有哪些牌可以胡，可以胡的翻数是多少
 */
func (d *MjDesk) GetJiaoInfos(user *MjUser) []*mjproto.JiaoInfo {
	jiaoInfos := []*mjproto.JiaoInfo{}

	//获取用户手牌 包括inPai
	var userPais []*MJPai
	userHandPai := *user.GetGameData().HandPai
	userPais = append(userPais, userHandPai.Pais...)
	userPais = append(userPais, userHandPai.InPai)

	//type JiaoInfo struct {
	//	OutCard          *CardInfo      `protobuf:"bytes,1,opt,name=outCard" json:"outCard,omitempty"`
	//	PaiInfos         []*JiaoPaiInfo `protobuf:"bytes,2,rep,name=paiInfos" json:"paiInfos,omitempty"`
	//	XXX_unrecognized []byte         `json:"-"`
	//}


	handPai := NewMJHandPai()
	var canHu, is19 bool

	userForPais := make([]*MJPai, len(userPais))
	copy(userForPais, userPais)

	handPai.GangPais = userHandPai.GangPais
	handPai.PengPais = userHandPai.PengPais

	for i := 0; i < len(userPais); i++ {
		//遍历用户手牌

		//从用户手牌中移除当前遍历的元素
		removedPai := userForPais[i]
		//log.T("removedPai is : %v", removedPai.GetDes())
		userForPais = removePaiFromPais(userForPais, i)
		//log.T("after remove user pais is:%v", userForPais)

		//copy(handPai.Pais, userForPais)
		handPai.Pais = userForPais
		jiaoInfo := NewJiaoInfo()

		//GetJiaoPais()
		for l := 0; l < len(mjpaiMap); l += 4 {

			//遍历未知牌
			//将遍历到的未知牌与用户手牌组合成handPai 去canhu
			mjPai := InitMjPaiByIndex(l)

			//定缺花色不用循环
			if user.GetGameData().GetHandPai().GetQueFlower() == mjPai.GetFlower() {
				//log.T("is ding que continue")
				continue
			}
			mjPaiLeftCount := int32(d.GetLeftPaiCount(user, mjPai)) //该可胡牌在桌面中的剩余数量 注 对于自己而言的剩余
			if mjPaiLeftCount == 0 {
				//log.T("left pai count is 0 continue")
				//剩余数为零不用循环
				continue
			}
			log.T("拿%v尝试胡牌", mjPai.GetDes())
			handPai.InPai = mjPai

			log.T("handPai: %v", handPai.GetDes())
			log.T("inPai: %v", handPai.InPai.GetDes())
			canHu, is19 = handPai.GetCanHu()

			if canHu {
				log.T("可胡")
				//可胡
				jiaoPaiInfo := NewJiaoPaiInfo()

				//计算番数得分胡牌类型
				fan, _, _ := GetHuScore(handPai, false, is19, 0, *d.GetRoomTypeInfo(), d)
				log.T("胡的番数%v", fan)
				//可胡牌的信息
				jiaoPaiInfo.HuCard = mjPai.GetCardInfo()
				*jiaoPaiInfo.Fan = fan //可胡番数
				*jiaoPaiInfo.Count = mjPaiLeftCount
				log.T("可胡%v, 牌剩余%v张", mjPai.GetDes(), *jiaoPaiInfo.Count)
				//打出去的牌信息
				jiaoInfo.OutCard = removedPai.GetCardInfo() //当前打出去的牌

				if jiaoPaiInfo != nil {
					jiaoInfo.PaiInfos = append(jiaoInfo.PaiInfos, jiaoPaiInfo)
					//log.T("可以胡 且 jiaoInfo is %v", jiaoInfo)
				}
			}

		}
		userForPais = addPaiIntoPais(removedPai, userForPais, i) //将移除的牌添加回原位置继续遍历
		//log.T("after add user pais is:%v", userPais)
		//log.T("after add user for pais is:%v", userForPais)
		if jiaoInfo.PaiInfos != nil {
			jiaoInfos = append(jiaoInfos, jiaoInfo)
		}
	}
	//log.T("jiaoInfos is %v", jiaoInfos)
	return jiaoInfos
}

//用户没有叫的处理了
func (d *MjDesk) DoDaJiao(u *MjUser) {
	log.T("开始处理玩家[%v]没叫,开始处理查大叫...", u.GetUserId())
}

//处理lottery的数据

//需要保存到 ..T_mj_desk_round   ...这里设计到保存数据，战绩相关的查询都要从这里查询
func (d *MjDesk) DoLottery() error {
	log.T("desk(%v),gameNumber(%v)处理DoLottery()", d.GetDeskId(), d.GetGameNumber())
	data := model.T_mj_desk_round{}
	data.DeskId = d.GetDeskId()
	data.GameNumber = d.GetGameNumber()
	data.BeginTime = timeUtils.String2YYYYMMDDHHMMSS(d.GetBeginTime())
	data.EndTime = time.Now()
	data.UserIds = d.GetUserIds()

	//一次处理每个胡牌的人
	for _, user := range d.GetUsers() {
		//这里不应该是胡牌的人才有记录...而是应该记录每一个人...
		if user != nil {
			//处理胡牌之后，分数相关的逻辑.
			//这里有一个统计...实在杠牌，或者胡牌之后会更新的数据...结算的时候，数据落地可以使用这个...
			//user.Statisc
			bean := model.MjRecordBean{}
			bean.UserId = user.GetUserId()
			bean.NickName = user.GetNickName()
			bean.WinAmount = user.Bill.GetWinAmount() //	赢了多少钱...

			//添加到record
			data.Records = append(data.Records, bean)
		}
	}

	//保存数据
	err := db.InsertMgoData(config.DBT_MJ_DESK_ROUND, &data)
	if err != nil {
		log.E("dolottery()时保存数据[%v]失败...", data)
	}
	log.T("desk(%v),gameNumber(%v)处理DoLottery(),处理完毕,保存数据data[%v]", d.GetDeskId(), d.GetGameNumber(), data)

	return nil

}

func (d *MjDesk) SendLotteryData() error {
	log.T("desk(%v),gameNumber(%v)SendLotteryData()", d.GetDeskId(), d.GetGameNumber())

	//发送开奖的数据,需要得到每个人的winCoinInfo
	result := newProto.NewGame_SendCurrentResult()
	for _, user := range d.GetUsers() {
		if user != nil {
			result.WinCoinInfo = append(result.WinCoinInfo, d.GetWinCoinInfo(user))
		}
	}

	//开始发送开奖的广播
	log.T("发送lottery的广播[%v]", result)
	d.BroadCastProto(result)
	log.T("desk(%v),gameNumber(%v)SendLotteryData(),处理完毕", d.GetDeskId(), d.GetGameNumber())

	return nil
}

func (d *MjDesk) AfterLottery() error {
	//开奖完成之后的一些处理

	//设置desk为准备的状态
	d.SetStatus(MJDESK_STATUS_READY)


	//设置用户为没有准备
	for _, user := range d.GetUsers() {
		user.AfterLottery()
	}
	return nil

}

func (d *MjDesk) End() bool {
	//判断结束的条件,目前只有局数能判断
	if d.GetCurrPlayCount() < d.GetTotalPlayCount() {
		//表示游戏还没有结束。。。.
		return false;
	} else {
		d.DoEnd()
		return true
	}
}

func (d *MjDesk)DoEnd() error {

	//1
	//首先发送游戏 结束的广播....game_SendEndLottery
	result := newProto.NewGame_SendEndLottery()
	for _, user := range d.GetUsers() {
		if user != nil {
			result.CoinInfo = append(result.CoinInfo, d.GetEndLotteryInfo(user))
		}
	}

	//发送之前需要判断谁是大赢家...这里暂时没有判断...

	//发送游戏结束的结果
	d.BroadCastProto(result)


	//2,清楚数据，解散房间....
	GetFMJRoom().DissolveDesk(d, false)

	return nil
}

//当前指针指向的玩家
func (d *MjDesk) SetActiveUser(userId uint32) error {
	*d.ActiveUser = userId
	return nil
}

//当前操作的玩家
func (d *MjDesk) SetActUserAndType(userId uint32, actType int32) error {
	*d.ActUser = userId
	*d.ActType = actType
	return nil
}

//得到下一个摸牌的人...
func (d *MjDesk) GetNextMoPaiUser() *MjUser {
	log.T("查询下一个玩家...当前的activeUser[%v]", d.GetActiveUser())
	var activeUser *MjUser = nil
	activeIndex := -1
	for i, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == d.GetActiveUser() {
			activeIndex = i
			break
		}
	}
	log.T("查询下一个玩家...当前的activeUser[%v],activeIndex[%v]", d.GetActiveUser(), activeIndex)
	if activeIndex == -1 {
		return nil
	}

	for i := activeIndex + 1; i < activeIndex + int(d.GetUserCount()); i++ {
		user := d.GetUsers()[i % int(d.GetUserCount())]
		log.T("查询下一个玩家...当前的activeUser[%v],activeIndex[%v],循环检测index[%v],user.IsNotHu(%v),user.CanMoPai[%v]", d.GetActiveUser(), activeIndex, i, user.IsNotHu(), user.CanMoPai(d.IsXueLiuChengHe()))
		if user != nil && user.CanMoPai(d.IsXueLiuChengHe()) {
			activeUser = user
			break
		}
	}

	//找到下一个操作的user
	return activeUser

}

//得到下一张牌...
func (d *MjDesk) GetNextPai() *MJPai {
	*d.MJPaiCursor ++
	//目前暂时是108张牌...
	if d.GetMJPaiCursor() >= 108 {
		log.E("服务器错误:要找的牌的坐标[%v]已经超过整副麻将的坐标了... ", d.GetMJPaiCursor())
		*d.MJPaiCursor --
		return nil
	} else {

		p := d.AllMJPai[d.GetMJPaiCursor()]
		pai := NewMjpai()
		*pai.Des = p.GetDes()
		*pai.Flower = p.GetFlower()
		*pai.Index = p.GetIndex()
		*pai.Value = p.GetValue()
		return pai
	}
}

//发送摸牌的广播
//指定一个摸牌，如果没有指定，则系统通过游标来判断
func (d *MjDesk) SendMopaiOverTurn(user *MjUser) error {
	if user == nil {
		user = d.GetNextMoPaiUser()
	}

	if user == nil {
		log.E("服务器出现错误..没有找到下一个摸牌的玩家...")
		return errors.New("没有找到下一家")
	}

	d.SetActiveUser(user.GetUserId())        //用户摸牌之后，设置前端指针指向的玩家
	d.SetActUserAndType(user.GetUserId(), MJDESK_ACT_TYPE_MOPAI)                //用户摸牌之后，设置当前活动的玩家

	//发送给当事人时候的信息
	nextPai := d.GetNextPai()
	//这里需要判断，如果牌摸完了，需要判断游戏结束
	if nextPai == nil {
		d.Lottery()
		return errors.New("牌摸完了，游戏结束...")
	}

	user.GameData.HandPai.InPai = nextPai

	overTrun := d.GetMoPaiOverTurn(user, false)        //用户摸牌的时候,发送一个用户摸牌的overturn
	user.SendOverTurn(overTrun)
	log.T("玩家[%v]当前的手牌是[%v]开始摸牌【%v】...", user.GetUserId(), user.GameData.HandPai.GetDes(), overTrun)


	//给其他人广播协议
	*overTrun.CanHu = false
	*overTrun.CanGang = false
	overTrun.ActCard = NewBackPai()
	d.BroadCastProtoExclusive(overTrun, user.GetUserId())

	return nil
}

func (d *MjDesk) GetDingQueEndInfo() *mjproto.Game_DingQueEnd {
	end := newProto.NewGame_DingQueEnd()

	for _, u := range d.GetUsers() {
		if u != nil && u.GameData.HandPai != nil {
			bean := newProto.NewDingQueEndBean()
			*bean.UserId = u.GetUserId()
			*bean.Flower = u.GameData.HandPai.GetQueFlower()
			end.Ques = append(end.Ques, bean)
		}
	}
	return end
}

func (d *MjDesk) ActPeng(userId uint32) error {
	//1找到玩家
	user := d.GetUserByUserId(userId)
	if user == nil {
		return errors.New("服务器错误碰牌失败")
	}

	if d.CheckCase == nil {
		log.E("用户[%v]碰牌的时候，没有找到可以碰的牌...", userId)
		return errors.New("服务器错误碰牌失败")

	}

	//2.1开始碰牌的操作
	pengPai := d.CheckCase.CheckMJPai
	canPeng := user.GameData.HandPai.GetCanPeng(pengPai)
	if !canPeng {
		//如果不能碰，直接返回
		log.E("玩家[%v]碰牌id[%v]-[%v]的时候，出现错误，碰不了...", userId, pengPai.GetIndex(), pengPai.LogDes())
		return errors.New("服务器出现错误..")
	}

	user.GameData.HandPai.InPai = nil
	user.GameData.HandPai.PengPais = append(user.GameData.HandPai.PengPais, pengPai)        //碰牌
	user.DelGuoHuInfo()        //删除过胡的信息

	//碰牌之后，需要删掉的pai...
	var pengKeys []int32
	//pengKeys = append(pengKeys, pengPai.GetIndex())
	for _, pai := range user.GameData.HandPai.Pais {
		if pai != nil && pai.GetClientId() == pengPai.GetClientId() {
			user.GameData.HandPai.PengPais = append(user.GameData.HandPai.PengPais, pai)        //碰牌
			pengKeys = append(pengKeys, pai.GetIndex())
			//碰牌只需要拆掉手里的两张牌
			if len(pengKeys) == 2 {
				break;
			}
		}
	}

	//2.2 删除手牌
	for _, key := range pengKeys {
		user.GameData.HandPai.DelHandlPai(key)
	}

	//3,生成碰牌信息
	//user.GameData.


	//4,处理 checkCase
	d.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), CHECK_CASE_BEAN_STATUS_CHECKED)
	d.CheckCase.UpdateChecStatus(CHECK_CASE_STATUS_CHECKED) //碰牌之后，checkcase处理完毕
	d.SetActiveUser(user.GetUserId())
	d.SetActUserAndType(user.GetUserId(), MJDESK_ACT_TYPE_DAPAI) //轮到用户打牌

	//5,发送碰牌的广播
	ack := newProto.NewGame_AckActPeng()
	ack.JiaoInfos = d.GetJiaoInfos(user)
	*ack.Header.Code = intCons.ACK_RESULT_SUCC
	*ack.UserIdOut = d.CheckCase.GetUserIdOut()
	*ack.UserIdIn = user.GetUserId()
	//组装牌的信息
	for _, ackpai := range user.GameData.HandPai.PengPais {
		if ackpai != nil && ackpai.GetClientId() == pengPai.GetClientId() {
			ack.PengCard = append(ack.PengCard, ackpai.GetCardInfo())
		}
	}
	d.BroadCastProto(ack)

	//最后设置checkCase = nil
	d.CheckCase = nil        //设置为nil
	return nil
}

//检测是否轮到当前玩家打牌...
func (d *MjDesk) CheckActive(userId uint32) bool {
	if d.GetActiveUser() == userId {
		return true        //检测通过
	} else {
		//没有轮到当前玩家
		log.E("非法请求，没有轮到当前玩家打牌..")
		return false
	}
}

//用户打一张牌出来
func (d *MjDesk)ActOut(userId uint32, paiKey int32) error {
	log.T("开始处理用户[%v]打牌[%v]的逻辑", userId, paiKey)

	outUser := d.GetUserByUserId(userId)
	if outUser == nil {
		log.E("打牌失败，没有找到玩家[%v]", userId)
		return errors.New("玩家[%v]没有找到，打牌失败...")
	}

	//判断是否轮到当前玩家打牌了...
	if !d.CheckActive(userId) {
		result := newProto.NewGame_AckSendOutCard()
		*result.Header.Code = intCons.ACK_RESULT_ERROR
		*result.Header.Error = "不是打牌的状态"
		outUser.WriteMsg(result)
		return errors.New("没有轮到当前玩家....")
	}

	//得到参数
	outPai := InitMjPaiByIndex(int(paiKey))
	/**
		1,如果是碰牌打牌的时候,inpai为nil，不需要增加
		2,如果是摸牌打牌（杠之后也是摸牌，需要增加in牌...）
	 */
	if outUser.GameData.HandPai.InPai != nil {
		outUser.GameData.HandPai.AddPai(outUser.GameData.HandPai.InPai)        //把inpai放置到手牌上
	}
	errDapai := outUser.DaPai(outPai)
	if errDapai != nil {
		log.E("打牌的时候出现错误，没有找到要到的牌,id[%v]", paiKey)
		return errors.New("用户打牌失败...没有找到牌")
	}

	outUser.GameData.HandPai.OutPais = append(outUser.GameData.HandPai.OutPais, outPai)        //自己桌子前面打出的牌，如果其他人碰杠胡了之后，需要把牌删除掉...
	outUser.GameData.HandPai.InPai = nil        //打牌之后需要把自己的  inpai给移除掉...
	outUser.DelGuoHuInfo()        //删除过胡的信息
	//打牌之后的逻辑,初始化判定事件
	err := d.InitCheckCase(outPai, outUser)
	if err != nil {
		//表示无人需要，直接给用户返回无人需要
		//给下一个人摸排，并且移动指针
		log.E("服务器错误，初始化判定牌的时候出错err[%v]", err)
	}

	log.T("InitCheckCase之后的checkCase[%v]", d.CheckCase)
	//回复消息,打牌之后，广播打牌的信息...s
	outUser.PreMoGangInfo = nil        //清楚摸牌前的杠牌info
	result := newProto.NewGame_AckSendOutCard()
	*result.UserId = userId
	result.Card = outPai.GetCardInfo()
	d.BroadCastProto(result)

	return nil

}

/**
		自摸

		//		出牌之前都有一个杠的状态，出牌之后设置这个杠为nil，胡牌的时候判断是否有杠的状态，有的话就根据杠的状态来判断是怎么胡的
				每次打牌之后，需要清空摸牌前的状态、

				///****

		1，如何判断是杠上花
		2，如何判断是明杠 杠上花
		3，如何判断是暗杠 杠上花
		4，如何判断是巴杠 杠上花

		5，如何判断是 天胡			//胡的时候，没有打过牌，
		6，如何判断是 地胡			//胡的时候，还有没有剩余的麻将


		点炮
		checkCase 需要注明判定盘打之前的杠状态
		1，普通点炮
		2，抢杠
		3，杠上花
	 */
//某人胡牌...
func (d *MjDesk)ActHu(userId uint32) error {

	//对于杠，有摸牌前 杠的状态，有打牌前杠的状态
	//1,胡的牌是当前check里面的牌，如果没有check，则表示是自摸

	//玩家胡牌
	huUser := d.GetUserByUserId(userId)
	if huUser == nil {
		log.E("服务器错误：没有找到胡牌的user[%v]", userId)
		return errors.New("服务器错误，没有找到胡牌的user")
	}
	checkCase := d.GetCheckCase()

	//设置判定牌
	if checkCase != nil {
		huUser.GameData.HandPai.InPai = checkCase.CheckMJPai
	}
	//判断是否可以胡牌，如果不能胡牌直接返回
	canHu, is19 := huUser.GameData.HandPai.GetCanHu()
	if !canHu {
		return errors.New("不可以胡牌...")
	}

	/**
	HuPaiType_H_QiangGang         HuPaiType = 11 todo
	HuPaiType_H_HaiDiLao          HuPaiType = 12 todo
	HuPaiType_H_HaiDiPao          HuPaiType = 13 todo
	HuPaiType_H_HaidiGangShangHua HuPaiType = 14 todo
	HuPaiType_H_HaidiGangShangPao HuPaiType = 15 todo
	 */

	//得到胡牌的信息
	var isZimo bool = false //是否是自摸
	var isGangShangHua bool = false //是否是杠上花
	var isGangShangPao bool = false
	var extraAct mjproto.HuPaiType = 0        //杠上花，
	var outUserId uint32
	var roomInfo mjproto.RoomTypeInfo = *d.GetRoomTypeInfo()  //roomType  桌子的规则
	var outUser *MjUser

	hupai := huUser.GameData.HandPai.InPai
	if checkCase == nil {
		//表示是自摸
		isZimo = true
		outUserId = userId
		if huUser.GetPreMoGangInfo() != nil {
			isGangShangHua = true  //杠上花
			extraAct = mjproto.HuPaiType_H_GangShangHua
		} else {
			extraAct = 0
		}

	} else {
		isZimo = false                //表示是点炮
		outUserId = checkCase.GetUserIdOut()
		if checkCase.GetPreOutGangInfo() != nil {
			//这里需要判断是杠上炮，还是抢杠
			isGangShangPao = true//杠上炮
			extraAct = mjproto.HuPaiType_H_GangShangPao
		} else {

		}
	}

	outUser = d.GetUserByUserId(outUserId)

	log.T("点炮的人[%v],胡牌的人[%v],杠上花[%v],杠上炮[%v],接下来开始getHuScore(%v,%v,%v,%v)", userId, outUserId, isGangShangHua, isGangShangPao,
		huUser.GameData.HandPai, isZimo, extraAct, roomInfo)

	fan, score, huCardStr := GetHuScore(huUser.GameData.HandPai, isZimo, is19, extraAct, roomInfo, d)
	log.T("胡牌(getHuScore)之后的结果fan[%v],score[%v],huCardStr[%v]", fan, score, huCardStr)

	//胡牌之后的信息
	hu := NewHuPaiInfo()
	*hu.GetUserId = huUser.GetUserId()
	*hu.SendUserId = outUserId
	//*hu.ByWho = 打牌的方位，对家，上家，下家？
	*hu.HuType = int32(extraAct)        ////杠上炮 杠上花 抢杠 海底捞 海底炮 天胡 地胡
	*hu.HuDesc = strings.Join(huCardStr, " ");
	*hu.Fan = fan
	*hu.Score = score        //只是胡牌的分数，不是赢了多少钱
	hu.Pai = hupai

	//胡牌之后，设置用户的数据

	huUser.AddHuPaiInfo(hu)

	//统计胡牌的次数
	huUser.StatisticsHuCount(d.GetCurrPlayCount(), huUser.GetUserId(), hu.GetHuType())
	//统计点炮的次数
	if !isZimo {
		outUser.StatisticsDianCount(outUserId, hu.GetHuType())
	}

	//处理抢杠的逻辑
	d.DoQiangGang(hu)

	//胡牌之后计算账单
	d.DoHuBill(hu)

	//点炮之后设置checkCase的状态
	d.DoAfterDianPao(hu)

	//设置下一次的庄
	d.InitNextBanker(hu)

	//发送胡牌成功的回复
	d.SendAckActHu(hu)
	return nil
}

//处理抢杠的逻辑
/**
处理抢杠的逻辑，抢杠的逻辑需要特殊处理...
1,首先是清楚杠牌的info
2,增加碰牌
3,删除杠牌的账单
*/
func (d *MjDesk) DoQiangGang(hu *HuPaiInfo) error {
	if d.CheckCase != nil && d.CheckCase.PreOutGangInfo != nil && d.CheckCase.PreOutGangInfo.GetGangType() == GANG_TYPE_BA {
		log.T("开始处理抢杠的逻辑....")
		dianUser := d.GetUserByUserId(hu.GetSendUserId())

		//1,首先是清楚杠牌的info
		var gangKeys []int32
		for _, pai := range dianUser.GameData.HandPai.GangPais {
			if pai == nil {
				continue
			}
			//需要删除的杠牌
			if pai != nil && pai.GetClientId() == hu.Pai.GetClientId() {
				gangKeys = append(gangKeys, pai.GetIndex())        //需要删除的杠牌
				if pai.GetIndex() != hu.Pai.GetIndex() {
					dianUser.GameData.HandPai.PengPais = append(dianUser.GameData.HandPai.PengPais, pai)        //碰牌
				}
			}
		}

		//删除杠牌的信息
		dianUser.GameData.DelGangInfo(hu.Pai)
		//删除杠牌的账单
		for _, billUser := range d.GetUsers() {
			//处理每一个人的账单,并且减去amount
			billUser.DelBillBean(hu.Pai)
		}
	}

	return nil
}

func (d *MjDesk) DoAfterDianPao(hu *HuPaiInfo) {
	//自摸的不用关
	if hu.GetGetUserId() == hu.GetSendUserId() {
		return
	}

	//点炮胡牌成功之后的处理... 处理checkCase
	d.SetActiveUser(hu.GetGetUserId())        // 胡牌之后 设置当前操作的用户为当前胡牌的人...
	d.CheckCase.UpdateCheckBeanStatus(hu.GetGetUserId(), CHECK_CASE_BEAN_STATUS_CHECKED)        // update checkCase...
	d.CheckCase.UpdateChecStatus(CHECK_CASE_STATUS_CHECKING_HUED)        //已经有人胡了，后边的人就不能碰或者杠了

}

//判断下一个庄是否已经确定
func (d *MjDesk) IsNextBankerExist() bool {
	if d.GetNextBanker() > 0 {
		return true
	} else {
		return false
	}
}




//设置下一次的庄
/**
	1，如果当前的nextBanker 没有值(nextBanker==0)，那代表此人是第一个胡牌的，设置为nextBanekr
	2,如果当前的nextBanker有值(nextBanker > 0 ),那需要判断是不是当前的点炮的人一炮双向
 */
func (d *MjDesk)InitNextBanker(hu *HuPaiInfo) {
	if d.IsNextBankerExist() {
		//已经存在的情况 //有双响就双响点炮的人做庄，不论之前是否有人胡牌  by 亮哥
		//这里可以用过pai 查询点炮账单的个数

		//判断是否是自摸
		isZimo := (hu.GetGetUserId() == hu.GetSendUserId())
		if isZimo {
			//如果是自摸，并且nextBanker已经有值了,那么直接返回不用设置
			return
		}

		//如果是点炮的，需要判断是不是双响
		dianUser := d.GetUserByUserId(hu.GetSendUserId())
		count := 0
		for _, bill := range dianUser.GetBill().Bills {
			if bill != nil && bill.Pai.GetIndex() == hu.GetPai().GetIndex() && bill.GetType() == MJUSER_BILL_TYPE_SHU_DIANPAO {
				count++
			}
		}

		//表示多响
		if count > 1 {
			//设置一炮多响的人为庄
			d.SetNextBanker(hu.GetSendUserId())
		}
	} else {
		d.SetNextBanker(hu.GetGetUserId())
	}
}

func (d *MjDesk) SetNextBanker(userId uint32) {
	*d.NextBanker = userId
}

func (d *MjDesk) SendAckActHu(hu *HuPaiInfo) {
	ack := newProto.NewGame_AckActHu()
	*ack.HuType = hu.GetHuType()        //这里需要判断是自摸还是点炮
	*ack.UserIdIn = hu.GetGetUserId()
	*ack.UserIdOut = hu.GetSendUserId()
	ack.HuCard = hu.Pai.GetCardInfo()
	*ack.IsZiMo = (hu.GetGetUserId() == hu.GetSendUserId())
	log.T("给用户[%v]广播胡牌的ack[%v]", hu.GetGetUserId(), ack)
	d.BroadCastProto(ack)
}

func (d *MjDesk)getPaiById(paiId int32) *MJPai {
	for _, pai := range d.AllMJPai {
		if pai != nil && pai.GetIndex() == paiId {
			return pai
		}
	}
	return nil

}

//杠牌   怎么判断是明杠，暗杠，巴杠...
func (d *MjDesk) ActGang(userId uint32, paiId int32) error {

	//检测参数是否正确
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("用户[%v]没有找到杠牌失败...", userId)
		return nil
	}

	gangPai := d.getPaiById(paiId);
	if gangPai == nil {
		log.E("用户[%v]没有找到杠牌,id[%v]，杠牌失败...", userId, paiId)
		return errors.New("服务器错误,杠牌失败..")
	}


	//判断是否可以杠牌

	var gangType int32 = 0
	var sendUserId uint32 = 0        //打出牌的人，暗杠的话 就表示是自己..
	var canGang bool = false

	if d.CheckCase != nil {
		//表示是明港
		gangType = GANG_TYPE_DIAN        //明杠
		sendUserId = d.CheckCase.GetUserIdOut()
		canGang, _ = user.GameData.HandPai.GetCanGang(gangPai)

	} else {
		canGang, _ = user.GameData.HandPai.GetCanGang(nil)
		//如果碰牌中有这张牌表示是巴杠 //如果碰牌中没有这张牌，表示是暗杠
		isBaGang := user.GameData.HandPai.IsExistPengPai(gangPai)
		if isBaGang {
			gangType = GANG_TYPE_BA        //巴杠
		} else {
			gangType = GANG_TYPE_AN        //暗杠
		}
		sendUserId = userId
	}


	//判断是否可以杠牌
	if !canGang {
		log.E("玩家[%v]杠牌[%v],牌id[%v]失败", userId, gangPai.LogDes(), paiId)
		return errors.New("用户杠牌失败..")
	}



	/**
		根据杠牌的类型不同，处理的方式也不同
		1,如果是巴杠，移除碰牌中的牌和碰牌info , 并且生成杠牌，和杠牌info
		2,如果是明杠或者暗杠，需要把所有的牌放在杠牌中，不用处理碰牌
	 */

	//首先把需要杠的牌放在手中
	user.GameData.HandPai.Pais = append(user.GameData.HandPai.Pais, user.GameData.HandPai.InPai)

	if gangType == GANG_TYPE_BA {
		log.T("用户[%v]杠牌是巴杠,现在处理巴杠...", userId)
		//循环碰牌来处理
		user.GameData.HandPai.GangPais = append(user.GameData.HandPai.GangPais, gangPai)

		var pengKeys []int32
		for _, pengPai := range user.GameData.HandPai.PengPais {
			if pengPai != nil && pengPai.GetClientId() == gangPai.GetClientId() {
				//增加杠牌
				user.GameData.HandPai.GangPais = append(user.GameData.HandPai.GangPais, pengPai)
				pengKeys = append(pengKeys, pengPai.GetIndex())
			}
		}

		//删除碰牌,手中的杠牌
		for _, key := range pengKeys {
			log.T("巴杠删除手牌..index[%v]", key)
			user.GameData.HandPai.DelPengPai(key)
		}

		//删除手牌
		user.GameData.HandPai.DelHandlPai(gangPai.GetIndex())//

	} else if gangType == GANG_TYPE_DIAN || gangType == GANG_TYPE_AN {
		log.T("用户[%v]杠牌不是巴杠 是 gangType[%v]...", userId, gangType)

		//杠牌的类型
		var gangKey []int32
		//增加杠牌
		//如果不是摸的牌，而是手中本来就有的牌，那么需要把他移除
		for _, pai := range user.GameData.HandPai.Pais {
			if pai.GetClientId() == gangPai.GetClientId() {
				//增加杠牌
				user.GameData.HandPai.GangPais = append(user.GameData.HandPai.GangPais, pai)
				gangKey = append(gangKey, pai.GetIndex())
			}
		}

		log.T("用户杠牌[%v]之后移除需要移除的手牌id数组[%v]", userId, gangKey)
		//减少手中的杠牌
		for _, key := range gangKey {
			log.T("用户杠牌[%v]之后移除手牌id[%v]", userId, key)
			user.GameData.HandPai.DelHandlPai(key)
		}
	}

	//增加杠牌info
	info := NewGangPaiInfo()
	*info.GetUserId = user.GetUserId()
	*info.SendUserId = sendUserId
	*info.GangType = gangType
	info.Pai = gangPai
	user.GameData.GangInfo = append(user.GameData.GangInfo, info)
	user.PreMoGangInfo = info        //增加杠牌状态
	user.GameData.HandPai.InPai = nil        //1,设置inpai为nil
	user.StatisticsGangCount(d.GetCurrPlayCount(), gangType)        //处理杠牌的账单
	user.DelGuoHuInfo()        //删除过胡的信息

	d.DoGangBill(info);
	d.DoCheckCaseAfterGang(gangType, gangPai, user)

	//杠牌之后的逻辑
	result := newProto.NewGame_AckActGang()
	*result.GangType = user.PreMoGangInfo.GetGangType()
	*result.UserIdOut = user.PreMoGangInfo.GetSendUserId()
	*result.UserIdIn = user.GetUserId()

	//组装杠牌的信息
	for _, ackpai := range user.GameData.HandPai.GangPais {
		if ackpai != nil && ackpai.GetClientId() == gangPai.GetClientId() {
			result.GangCard = append(result.GangCard, ackpai.GetCardInfo())
		}
	}
	log.T("广播玩家[%v]杠牌[%v]之后的ack[%v]", user.GetUserId(), gangPai, result)
	d.BroadCastProto(result)
	return nil
}


//处理账单
/**
	没有胡牌的人，都需要给钱  ,目前不是承包的方式...
 */


func (d *MjDesk) DoGangBill(info *GangPaiInfo) {
	gangType := info.GetGangType()
	gangUser := d.GetUserByUserId(info.GetGetUserId())
	gangPai := info.GetPai()

	if gangType == GANG_TYPE_AN {
		//处理暗杠的账单
		score := d.GetBaseValue() * 2        //暗杠的分数
		for _, ou := range d.GetUsers() {
			//不为nil 并且不是本人，并且没有胡牌
			if ou != nil && ou.GetUserId() != gangUser.GetUserId() && ou.IsGaming() {
				gangUser.AddBill(ou.GetUserId(), MJUSER_BILL_TYPE_YING_GNAG, "用户暗杠，收入", score, gangPai)        //用户赢钱的账户
				ou.AddBill(gangUser.GetUserId(), MJUSER_BILL_TYPE_SHU_GNAG, "用户暗杠，输钱", -score, gangPai)        //用户输钱的账单
			}
		}

	} else if gangType == GANG_TYPE_DIAN {
		//处理点杠点账单
		score := d.GetBaseValue() * 2        //点杠的分数
		shuUser := d.GetUserByUserId(info.GetSendUserId())
		gangUser.AddBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_GNAG, "用户点杠，收入", score, gangPai)        //用户赢钱的账户
		shuUser.AddBill(gangUser.GetUserId(), MJUSER_BILL_TYPE_SHU_GNAG, "用户点杠，输钱", -score, gangPai)        //用户输钱的账单

	} else if gangType == GANG_TYPE_BA {
		//处理巴杠的账单
		score := d.GetBaseValue()        //巴杠的分数
		for _, ou := range d.GetUsers() {
			if ou != nil && ou.GetUserId() != gangUser.GetUserId() && ou.IsGaming() {
				gangUser.AddBill(ou.GetUserId(), MJUSER_BILL_TYPE_YING_GNAG, "用户巴杠，收入", score, gangPai)        //用户赢钱的账户
				ou.AddBill(gangUser.GetUserId(), MJUSER_BILL_TYPE_SHU_GNAG, "用户巴杠，输钱", -score, gangPai)        //用户输钱的账单

			}
		}
	}
}

//计算胡牌的账单
/**
	增加账单
	//todo 这里需要完善算账的逻辑逻辑,目前就自摸和点炮来做
 */
func (d *MjDesk)DoHuBill(hu *HuPaiInfo) {
	isZimo := (hu.GetGetUserId() == hu.GetSendUserId())
	outUser := hu.GetSendUserId()
	huUser := d.GetUserByUserId(hu.GetGetUserId())

	log.T("玩家[%v]胡牌，开始处理计算分数的逻辑...", huUser.GetUserId())
	if isZimo {
		//如果是自摸的话，三家都需要给钱
		for _, shuUser := range d.GetUsers() {
			if shuUser != nil  && shuUser.IsGaming() && shuUser.GetUserId() != huUser.GetUserId() {

				//赢钱的账单
				huUser.AddBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_HU, "用户自摸，获得收入", hu.GetScore(), hu.Pai)

				//输钱的账单
				shuUser.AddBill(huUser.GetUserId(), MJUSER_BILL_TYPE_SHU_ZIMO, "用户自摸，输钱", -hu.GetScore(), hu.Pai)
			}
		}
	} else {

		//如果是点炮的话，只有一家需要给钱...
		shuUser := d.GetUserByUserId(outUser)

		//赢钱的账单
		huUser.AddBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_HU, "点炮胡牌，获得收入", hu.GetScore(), hu.Pai)

		//输钱的账单
		shuUser.AddBill(huUser.GetUserId(), MJUSER_BILL_TYPE_SHU_DIANPAO, "用户点炮，输钱", -hu.GetScore(), hu.Pai)
	}

}

func (d *MjDesk)DoCheckCaseAfterGang(gangType int32, gangPai *MJPai, user *MjUser) {
	d.CheckCase = nil        //设置 判断的为nil
	///如果是巴杠，需要设置巴杠的判断  initCheckCase
	if gangType == GANG_TYPE_BA {
		d.InitCheckCase(gangPai, user)
		if d.CheckCase == nil {
			log.T("巴杠没有人可以抢杠...")
		}
	}

}

//设置用户的状态为离线
func (d *MjDesk) SetOfflineStatus(userId uint32) {
	user := d.GetUserByUserId(userId)
	*user.IsBreak = true
}

//返回一个牌局结果

func (d *MjDesk) GetWinCoinInfo(user *MjUser) *mjproto.WinCoinInfo {
	win := newProto.NewWinCoinInfo()
	*win.NickName = user.GetNickName()
	*win.UserId = user.GetUserId()
	*win.WinCoin = user.Bill.GetWinAmount()        //本次输赢多少(负数表示输了)
	*win.Coin = user.GetCoin()        // 输赢以后，当前筹码是多少
	*win.CardTitle = d.GetCardTitle4WinCoinInfo(user)// 赢牌牌型信息( 如:"点炮x2 明杠x2 根x2 自摸 3番" )
	win.Cards = user.GetPlayerCard(true) //牌信息,true 表示要显示牌的信息...
	*win.IsDealer = user.GetIsBanker()      //是否是庄家
	*win.HuCount = user.Statisc.GetCountHu()        //本局胡的次数(血流成河会多次胡)
	return win
}

//得到这个人的胡牌描述
func (d *MjDesk) GetCardTitle4WinCoinInfo(user *MjUser) string {
	var huDesk string = ""                //胡牌的描述...
	//目前暂时返回hu的信息
	if user.GameData.HuInfo != nil && len(user.GameData.HuInfo) > 0 {
		huDesk = user.GameData.HuInfo[0].GetHuDesc()
	}
	return huDesk
}

//得到EndLotteryInfo结果...
func (d *MjDesk)GetEndLotteryInfo(user *MjUser) *mjproto.EndLotteryInfo {
	end := newProto.NewEndLotteryInfo()
	*end.UserId = user.GetUserId()
	*end.BigWin = false             //是否是大赢家...
	*end.CountAnGang = user.Statisc.GetCountAnGang()            //暗杠的次数
	*end.CountChaJiao = user.Statisc.GetCountChaJiao()          //查叫的次数..
	*end.CountDianGang = user.Statisc.GetCountDianGang()          // 点杠的次数
	*end.CountDianPao = user.Statisc.GetCountDianPao()          //点炮的次数
	*end.CountHu = user.Statisc.GetCountHu()//胡牌的次数
	*end.CountZiMo = user.Statisc.GetCountZiMo()              //自摸的次数
	*end.WinCoin = user.Statisc.GetWinCoin()//赢了多少钱
	return end
}

/**
	房间是否开始游戏
	1,开始:游戏已经开始
	2，游戏并没有开始，round==0
 */
func (d *MjDesk) IsBegin() bool {
	return d.GetCurrPlayCount() > 0
}

//剩余牌的数量
func (d *MjDesk) GetRemainPaiCount() int32 {
	//todo 几门牌? 这里的107需要通过有几门牌来确定...
	return 107 - d.GetMJPaiCursor()
}

func (d *MjDesk) GetByWho() {

}

//判断是否是血流成河
func (d *MjDesk) IsXueLiuChengHe() bool {
	return d.GetMjRoomType() == int32(mjproto.MJRoomType_roomType_xueLiuChengHe)
}

//换三张
func (d *MjDesk) DoExchange(userId uint32, exchangeNum int32, cards []*mjproto.CardInfo) error {
	//换三张需要同步
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("换三张失败,因为没有找到对应的玩家[%v]", userId)
		return errors.New("换三张失败")
	}

	//判断如参是否正确
	if cards == nil || int32(len(cards)) != exchangeNum {
		return errors.New("换三张失败...s")
	}

	for _, card := range cards {
		pai := InitMjPaiByIndex(int(card.GetId()))
		user.ExchangeCards = append(user.ExchangeCards, pai)        //增加需要换的牌
		user.GameData.HandPai.DelHandlPai(card.GetId())        //删除手牌
	}
	//设置已经换了
	*user.Exchanged = true

	//返回结果
	result := newProto.NewGame_AckExchangeCards()
	*result.Header.Code = intCons.ACK_RESULT_SUCC
	*result.UserId = user.GetUserId()
	d.BroadCastProto(result)

	//之后判断
	go d.ExchangeEnd()

	return nil
}

//开始完成换三张
func (d *MjDesk) ExchangeEnd() error {
	lock := lock.GetDeskLock(d.GetDeskId())
	lock.Lock()
	defer lock.Unlock()

	//首先
	for _, user := range d.GetUsers() {
		if user == nil || !user.GetExchanged() {
			return errors.New("还有人没有请求换三张")
		}
	}

	time.Sleep(time.Second * 5)

	//开始换牌
	exchangeType := utils.Rand(0, 3)
	if exchangeType == int32(mjproto.ExchangeType_EXCHANGE_TYPE_DUIJIA) {
		exchangeCards(d.Users[0], d.Users[2])
		exchangeCards(d.Users[1], d.Users[3])
		exchangeCards(d.Users[2], d.Users[0])
		exchangeCards(d.Users[3], d.Users[1])

	} else if exchangeType == int32(mjproto.ExchangeType_EXCHANGE_TYPE_SHUNSHIZHEN) {
		exchangeCards(d.Users[0], d.Users[3])
		exchangeCards(d.Users[1], d.Users[0])
		exchangeCards(d.Users[2], d.Users[1])
		exchangeCards(d.Users[3], d.Users[2])

	} else if exchangeType == int32(mjproto.ExchangeType_EXCHANGE_TYPE_NISHIZHEN) {
		exchangeCards(d.Users[0], d.Users[1])
		exchangeCards(d.Users[1], d.Users[2])
		exchangeCards(d.Users[2], d.Users[3])
		exchangeCards(d.Users[3], d.Users[0])
	}

	//最后三张表示是已经换了的牌
	for _, user := range d.GetUsers() {
		result := newProto.NewGame_ExchangeCardsEnd()
		*result.Header.Code = intCons.ACK_RESULT_SUCC
		*result.Header.UserId = user.GetUserId()
		*result.ExchangeType = exchangeType
		paiCount := len(user.GameData.HandPai.Pais)
		ps := user.GameData.HandPai.Pais[paiCount - 3:paiCount]
		for _, rp := range ps {
			c := rp.GetCardInfo()
			result.ExchangeInCards = append(result.ExchangeInCards, c)
		}
		//给用户发送换牌之后的信息
		user.WriteMsg(result)
	}


	//延时之后发送开始定缺的广播

	time.Sleep(time.Second * 3)
	//开始定缺
	err := d.beginDingQue()
	if err != nil {
		log.E("开始发送定缺广播的时候出错err[%v]", err)
		return err
	}

	return nil
}

func exchangeCards(u1 *MjUser, u2 *MjUser) {
	//换三张的账户
	count := 3
	cars := make([]*MJPai, count)
	copy(cars, u2.ExchangeCards)        //copy ，防止出错
	u1.GameData.HandPai.Pais = append(u1.GameData.HandPai.Pais, cars...)

}

//判断是否开启房间的某个选
func (d *MjDesk) IsOpenOption(option mjproto.MJOption) bool {
	for _, opt := range d.GetOthersCheckBox() {
		if opt == int32(option) {
			return true
		}
	}
	return false

}

//可以把overturn放在一个地方,目前都是摸牌的时候在用
func (d *MjDesk) GetMoPaiOverTurn(user *MjUser, isOpen bool) *mjproto.Game_OverTurn {
	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = user.GetUserId()                     //这个是摸牌的，所以是广播...
	*overTurn.PaiCount = d.GetRemainPaiCount()                //桌子剩余多少牌
	*overTurn.ActType = OVER_TURN_ACTTYPE_MOPAI                //摸牌
	*overTurn.Time = 30
	if isOpen {
		overTurn.ActCard = user.GameData.HandPai.InPai.GetBackPai()
	} else {
		overTurn.ActCard = user.GameData.HandPai.InPai.GetCardInfo()
	}
	*overTurn.CanHu, _ = user.GameData.HandPai.GetCanHu()        //是否可以胡牌
	*overTurn.CanPeng = false        //是否可以碰牌

	//处理杠牌的时候
	/**
		1，血战到底：用户胡牌之后是不会进入到这个方法的
		2，血流成河：用户已经胡牌，那么杠牌之后，胡牌不会改变的情况下，才可以杠 // todo
	 */
	canGangBool, gangPais := user.GameData.HandPai.GetCanGang(nil)    //是否可以杠牌
	*overTurn.CanGang = canGangBool
	if canGangBool && gangPais != nil {
		if user.IsHu() && d.IsXueLiuChengHe() {
			//血流成河，胡牌之后 杠牌的逻辑
			jiaoPais := user.GetJiaoPaisByHandPais(); //得到杠牌之前的可以胡的叫牌
			for _, g := range gangPais {
				//判断杠牌之后的叫牌是否和杠牌之前一样
				if user.AfterGangEqualJiaoPai(jiaoPais, g) {
					overTurn.GangCards = append(overTurn.GangCards, g.GetCardInfo())
				}
			}

		} else {
			//没有胡牌之前，杠牌的逻辑....
			for _, g := range gangPais {
				overTurn.GangCards = append(overTurn.GangCards, g.GetCardInfo())
			}
		}
	}

	//最后判断是否可以杠牌
	if overTurn.GangCards == nil || len(overTurn.GangCards) <= 0 {
		*overTurn.CanGang = false;
	}

	//
	overTurn.JiaoInfos = d.GetJiaoInfos(user)

	return overTurn
}

//通过checkCase 得到一个OverTurn
func (d *MjDesk) GetOverTurnByCaseBean(checkPai *MJPai, caseBean *CheckBean, actType int32) *mjproto.Game_OverTurn {
	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = caseBean.GetUserId()
	*overTurn.CanGang = caseBean.GetCanGang()
	*overTurn.CanPeng = caseBean.GetCanPeng()
	*overTurn.CanHu = caseBean.GetCanHu()
	*overTurn.PaiCount = d.GetRemainPaiCount()        //剩余多少钱
	overTurn.ActCard = checkPai.GetCardInfo()        //
	*overTurn.ActType = actType
	*overTurn.Time = 30
	return overTurn
}

func (d *MjDesk) GetLeftPaiCount(user *MjUser, mjPai *MJPai) int {
	var count int = 0
	displayPais := d.GetDisplayPais(user)
	//for i := 0; i < len(displayPais); i++ {
	//	log.T("用户%v已知的牌是:%v", user.GetUserId(), displayPais[i].GetDes())
	//}
	for i := 0; i < len(displayPais); i++ {
		if (displayPais[i].GetValue() == mjPai.GetValue()) && (displayPais[i].GetFlower() == mjPai.GetFlower()) {
			count++
		}
	}
	count = 4 - count
	if count < 0 {
		count = 0
	}
	log.T("leftPai is %v Count is : %v", mjPai.GetDes(), count)
	return count
}

//获取用户未知的牌 即未出现在台面上的牌
func (d *MjDesk) GetHiddenPais(user *MjUser) []*MJPai {

	//获取已知亮出台面的牌
	displayPais := d.GetDisplayPais(user)
	displayPaiCounts := GettPaiStats(displayPais)

	//获取一副完整的牌
	totalPais := XiPai()
	totalPaiCounts := GettPaiStats(totalPais)

	for i := 0; i < len(displayPaiCounts); i++ {
		totalPaiCounts[i] -= displayPaiCounts[i]
		if totalPaiCounts[i] < 0 {
			//最低为零
			totalPaiCounts[i] = 0
		}
	}

	return GetPaisByCounts(totalPaiCounts)
}

//获取用户已知亮出台面的牌 包括自己手牌、自己和其他玩家碰杠牌、其他玩家outPais
func (d *MjDesk) GetDisplayPais(user *MjUser) []*MJPai {
	//获取所有玩家的亮出台面的牌 outPais + pengPais + gangPais
	users := d.GetUsers()
	displayPais := []*MJPai{}
	for i := 0; i < len(users); i++ {
		displayPais = append(displayPais, users[i].GetGameData().GetHandPai().OutPais...) //打出去的牌

		userHandPai := users[i].GetGameData().GetHandPai()
		displayPais = append(displayPais, userHandPai.GangPais...) //杠的牌
		displayPais = append(displayPais, userHandPai.PengPais...) //碰的牌
	}

	//在亮出台面的牌中加入用户自己的手牌
	userHandPai := user.GetGameData().GetHandPai()
	displayPais = append(displayPais, userHandPai.InPai)
	displayPais = append(displayPais, userHandPai.Pais...)

	return displayPais
}