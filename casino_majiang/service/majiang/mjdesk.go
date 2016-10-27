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

//判断是不是朋友桌
func (d *MjDesk) IsFriend() bool {
	return true
}

//朋友桌用户加入房间
func (d *MjDesk) addNewUserFriend(userId uint32, a gate.Agent) error {

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
		d.SendReconnectOverTurn(userLeave.GetUserId())
		return nil
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
		return errors.New("用户加入房间失败")
	} else {
		//加入房间成功
		return nil
	}
}

//发送重新连接之后的overTurn
func (d *MjDesk) SendReconnectOverTurn(userId uint32) error {

	//得到玩家
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.T("发送SendReconnectOverTurn(%v)失败,因为没有找到玩家", userId)
		return errors.New("发送SendReconnectOverTurn()失败,因为没有找到玩家")
	}

	//开发发送overTurn
	if d.IsReady() && user.IsNotReady() {
		//给玩家发送开始准备的消息...
	} else if d.IsExchange() {
		//给玩家发送换牌的消息...
	} else if d.IsDingQue() && user.IsNotDingQue() {
		//给玩家发送定缺的信息...
		beginQue := newProto.NewGame_BroadcastBeginDingQue()
		user.WriteMsg(beginQue)

	} else if d.IsGaming() && user.GetUserId() == d.GetActUser() {
		//游戏中的情况，发送act的消息,这里需要更具当前的状态来发送overTurn
		if d.GetActType() == MJDESK_ACT_TYPE_MOPAI {
			//发送摸牌的协议
			overTrun := newProto.NewGame_OverTurn()
			*overTrun.UserId = user.GetUserId()                //这个是摸牌的，所以是广播...
			*overTrun.ActType = OVER_TURN_ACTTYPE_MOPAI        //摸牌
			overTrun.ActCard = user.GameData.HandPai.InPai.GetCardInfo()
			*overTrun.CanHu = user.GameData.HandPai.GetCanHu()                //是否可以胡牌
			canGangBool, gangPais := user.GameData.HandPai.GetCanGang(nil)    //是否可以杠牌

			*overTrun.CanGang = canGangBool
			if canGangBool && gangPais != nil {
				for _, g := range gangPais {
					overTrun.GangCards = append(overTrun.GangCards, g.GetCardInfo())
				}
			}
			//是否可以碰牌
			*overTrun.CanPeng = false
			//user.SendOverTurn(overTrun) 这里不能使用sendOverTurn 有wait 的逻辑
			user.WriteMsg(overTrun)
			log.T("玩家重新进入游戏之后 [%v]开始摸牌【%v】...", user.GetUserId(), overTrun)

		} else if d.GetActType() == MJDESK_ACT_TYPE_DAPAI {
			//发送打牌的协议
		} else if d.GetActType() == MJDESK_ACT_TYPE_WAIT_CHECK {
			caseBean := d.CheckCase.GetBeanByUserIdAndStatus(user.GetUserId(), CHECK_CASE_BEAN_STATUS_CHECKING)
			if caseBean == nil {
				log.E("没有找到玩家[%v]对应的checkBean", user.GetUserId())
				return errors.New("玩家重新进入房间发送check overturn的时候出错")
			}

			//找到需要判断bean之后，发送给判断人	//send overTurn
			overTurn := newProto.NewGame_OverTurn()
			*overTurn.UserId = caseBean.GetUserId()
			*overTurn.CanGang = caseBean.GetCanGang()
			*overTurn.CanPeng = caseBean.GetCanPeng()
			*overTurn.CanHu = caseBean.GetCanHu()
			overTurn.ActCard = d.CheckCase.CheckMJPai.GetCardInfo()        //
			*overTurn.ActType = OVER_TURN_ACTTYPE_OTHER
			*overTurn.Time = int32(user.GetWaitTime() - time.Now().Unix())

			///发送overTurn 的信息
			log.T("开始发送玩家[%v]断线重连的overTurn[%v]", user.GetUserId(), overTurn)
			user.WriteMsg(overTurn)
			//这里不需要设置actUser and type 因为在断线之前已经设置过了...
			//d.SetActUserAndType(caseBean.GetUserId(), MJDESK_ACT_TYPE_WAIT_CHECK)        //设置当前活动的玩家
		}
	}

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
	//deskInfo.NInitActionTime
	//deskInfo.NInitDelayTime
	*deskInfo.ActiveUserId = d.GetActiveUser()
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

			//判断是否是房主
			isOwner := false
			if d.GetOwner() == user.GetUserId() {
				isOwner = true
			}

			//得到信息
			if user.GetUserId() == receiveUserId {
				info := user.GetPlayerInfo(true)
				*info.IsOwner = isOwner
				players = append(players, info)
			} else {
				info := user.GetPlayerInfo(false)
				*info.IsOwner = isOwner
				players = append(players, user.GetPlayerInfo(false))
			}
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

//是否在准备中
func (d *MjDesk) IsReady() bool {
	return true
}

//是否在定缺中
func (d *MjDesk) IsDingQue() bool {
	return true
}

func (d *MjDesk) IsExchange() bool {
	return true
}

func (d *MjDesk) IsGaming() bool {
	return true
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
	if d.IsAllReady() && d.IsPlayerEnough() {
		d.begin()
	}

	return nil
}

//开始游戏
/**
开始游戏需要有几个步骤
1，desk的状态是否正确，现在是否可以开始游戏


 */
func (d *MjDesk) begin() error {
	//1，检查是否可以开始游戏
	err := d.time2begin()
	if err != nil {
		log.E("无法开始游戏:err[%v]", err)
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


	//4，开始定缺
	err = d.beginDingQue()
	if err != nil {
		log.E("开始发送定缺广播的时候出错err[%v]", err)
		return err
	}

	return nil
}

//是否可以开始
func (d *MjDesk) time2begin() error {
	log.T("检测游戏是否可以开始... ")
	return nil
}


/**
1,初始化desk
2,初始化user
 */
func (d *MjDesk) beginInit() error {

	//初始化每个玩家的信息
	for _, user := range d.GetUsers() {
		if user != nil {
			user.GameData = NewPlayerGameData()        //初始化一个空的麻将牌
			user.SetStatus(MJUSER_STATUS_READY)
		}
	}

	//初始化桌子的信息
	//1,初始化庄的信息,如果目前没有庄，则设置房主为庄,如果有庄，则不用管，每局游戏借宿的时候，会设置下一局的庄
	if d.GetBanker() == 0 {
		*d.Banker = d.GetOwner()
	}

	*d.GetBankerUser().IsBanker = true        //设置庄
	d.SetActiveUser(d.GetBanker())        //设置当前的活动玩家
	*d.GameNumber, _ = db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)        //设置游戏编号
	d.AddCurrPlayCount()        //场次数目加一
	*d.BeginTime = timeUtils.Format(time.Now())

	//发送游戏开始的协议...
	log.T("发送游戏开始的协议..")
	open := newProto.NewGame_Opening()
	d.BroadCastProto(open)
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
	//d.AllMJPai = XiPai()
	d.SetStatus(MJDESK_STATUS_FAPAI)        //发牌的阶段
	d.AllMJPai = XiPaiTestHu()
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

//开始定缺
func (d *MjDesk) beginDingQue() error {
	//开始定缺，修改desk的状态
	d.SetStatus(MJDESK_STATUS_DINGQUE)


	//给每个人发送开始定缺的信息
	beginQue := newProto.NewGame_BroadcastBeginDingQue()
	log.T("开始给玩家发送开始定缺的广播[%v]", beginQue)
	//
	time.Sleep(time.Second * 5)
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
	user.SetStatus(MJUSER_STATUS_DINGQUE)        //设置目前的状态是已经定缺
	*user.GameData.HandPai.QueFlower = color

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
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() &&  d.CanInitCheckCase(checkUser) {
			log.T("用户[%v]打牌，判断user[%v]是否可以碰杠胡...", outUser.GetUserId(), checkUser.GetUserId())
			checkUser.GameData.HandPai.InPai = p
			bean := checkUser.GetCheckBean(p)
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

//判断是否可以initCheckCase

func (d *MjDesk ) CanInitCheckCase(user *MjUser) bool {
	//这里需要判断是否是 血流成河，目前暂时不判断...
	if user.IsNotHu() {
		return true
	} else {
		return false
	}

}

//暂时不用？？ 摸牌之后
func (d *MjDesk) InitMoPaiCheckCase(p *MJPai, moPaiUser *MjUser) error {

	//初始化参数
	moPaiUser.GameData.HandPai.InPai = p

	//判断可能性
	checkCase := NewCheckCase()
	*checkCase.UserIdOut = moPaiUser.GetUserId()
	*checkCase.CheckStatus = CHECK_CASE_STATUS_CHECKING        //正在判定
	checkCase.CheckMJPai = p

	checkCase.PreOutGangInfo = moPaiUser.GetPreMoGangInfo()
	checkCase.CheckB = append(checkCase.CheckB, moPaiUser.GetCheckBean(p))
	if checkCase.CheckB == nil || len(checkCase.CheckB) > 0 {
		d.CheckCase = checkCase
	} else {
		d.CheckCase = nil
	}

	return nil
}


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
		overTurn := newProto.NewGame_OverTurn()
		*overTurn.UserId = caseBean.GetUserId()
		*overTurn.CanGang = caseBean.GetCanGang()
		*overTurn.CanPeng = caseBean.GetCanPeng()
		*overTurn.CanHu = caseBean.GetCanHu()
		overTurn.ActCard = d.CheckCase.CheckMJPai.GetCardInfo()        //
		*overTurn.ActType = OVER_TURN_ACTTYPE_OTHER

		///发送overTurn 的信息
		log.T("开始发送overTurn[%v]", overTurn)
		d.GetUserByUserId(caseBean.GetUserId()).SendOverTurn(overTurn)
		d.SetActUserAndType(caseBean.GetUserId(), MJDESK_ACT_TYPE_WAIT_CHECK)        //设置当前活动的玩家

		return nil

	}

}


/**
	1，只剩一个玩家没有胡牌
	2, 已经没有牌了...
 */

func (d *MjDesk) Time2Lottery() bool {
	//游戏中的玩家只剩下一个人，表示游戏结束...
	gamingCount := d.GetGamingCount()        //正在游戏中的玩家数量

	log.T("判断是否胡牌...但钱的gamingCount[%v]", gamingCount)

	if gamingCount != 1 {
		//正在游戏中的玩家的数量不为1，表示还没有结束
		return false
	}

	//所有的条件都满足，一局麻将结束...
	return true
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

	d.DoLottery()

	//发送结束的广播
	d.SendLotteryData()

	//
	d.AfterLottery()

	//判断牌局结束(整场游戏结束)
	if !d.End() {
		go d.begin()
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
	return nil

}

func (d *MjDesk) End() bool {
	//判断结束的条件,目前只有局数能判断
	if d.GetCurrPlayCount() < d.GetTotalPlayCount() {
		//表示游戏还没有结束。。。.
		return false;
	} else {
		d.DoEnd()
	}

	return true
}

func (d *MjDesk)DoEnd() error {
	//game_SendEndLottery
	result := newProto.NewGame_SendEndLottery()
	for _, user := range d.GetUsers() {
		if user != nil {
			result.CoinInfo = append(result.CoinInfo, d.GetEndLotteryInfo(user))
		}
	}

	//发送游戏结束的结果
	d.BroadCastProto(result)
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
		log.T("查询下一个玩家...当前的activeUser[%v],activeIndex[%v],循环检测index[%v],user.IsNotHu(%v),user.CanMoPai[%v]", d.GetActiveUser(), activeIndex, i, user.IsNotHu(), user.CanMoPai())
		if user != nil && user.CanMoPai() {
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

	//发送广播
	overTrun := newProto.NewGame_OverTurn()
	*overTrun.UserId = user.GetUserId()                //这个是摸牌的，所以是广播...
	*overTrun.ActType = OVER_TURN_ACTTYPE_MOPAI        //摸牌

	//发送给当事人时候的信息
	nextPai := d.GetNextPai()
	user.GameData.HandPai.InPai = nextPai
	overTrun.ActCard = nextPai.GetCardInfo()

	//是否可以胡牌
	*overTrun.CanHu = user.GameData.HandPai.GetCanHu()
	//是否可以杠牌
	canGangBool, gangPais := user.GameData.HandPai.GetCanGang(nil)
	*overTrun.CanGang = canGangBool
	if canGangBool && gangPais != nil {
		for _, g := range gangPais {
			overTrun.GangCards = append(overTrun.GangCards, g.GetCardInfo())
		}
	}

	//是否可以碰牌
	*overTrun.CanPeng = false
	user.SendOverTurn(overTrun)
	log.T("玩家[%v]开始摸牌【%v】...", user.GetUserId(), overTrun)


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
	//todo 需要判断是否是可以碰


	//2.1开始碰牌的操作
	pengPai := d.CheckCase.CheckMJPai
	user.GameData.HandPai.PengPais = append(user.GameData.HandPai.PengPais, pengPai)        //碰牌

	//碰牌之后，需要删掉的pai...
	var pengKeys []int32
	//pengKeys = append(pengKeys, pengPai.GetIndex())
	for _, pai := range user.GameData.HandPai.Pais {
		if pai != nil && pai.GetIndex() == pengPai.GetIndex() {
			user.GameData.HandPai.PengPais = append(user.GameData.HandPai.PengPais, pai)        //碰牌
			pengKeys = append(pengKeys, pai.GetIndex())
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
	*ack.Header.Code = intCons.ACK_RESULT_SUCC
	*ack.UserIdOut = d.CheckCase.GetUserIdOut()
	*ack.UserIdIn = user.GetUserId()

	//todo 临时处理
	ack.PengCard[0] = d.CheckCase.CheckMJPai.GetCardInfo()
	ack.PengCard[1] = d.CheckCase.CheckMJPai.GetCardInfo()
	ack.PengCard[2] = d.CheckCase.CheckMJPai.GetCardInfo()
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
	outUser.GameData.HandPai.AddPai(outUser.GameData.HandPai.InPai)        //把inpai放置到手牌上
	errDapai := outUser.DaPai(outPai)
	if errDapai != nil {
		log.E("打牌的时候出现错误，没有找到要到的牌,id[%v]", paiKey)
		return errors.New("用户打牌失败...没有找到牌")
	}

	outUser.GameData.HandPai.OutPais = append(outUser.GameData.HandPai.OutPais, outPai)        //自己桌子前面打出的牌，如果其他人碰杠胡了之后，需要把牌删除掉...
	outUser.GameData.HandPai.InPai = nil        //打牌之后需要把自己的  inpai给移除掉...

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
	u := d.GetUserByUserId(userId)
	if u == nil {
		log.E("服务器错误：没有找到胡牌的user[%v]", userId)
		return errors.New("服务器错误，没有找到胡牌的user")
	}
	checkCase := d.GetCheckCase()

	//设置判定牌
	if checkCase != nil {
		u.GameData.HandPai.InPai = checkCase.CheckMJPai
	}
	//判断是否可以胡牌，如果不能胡牌直接返回
	canHu := u.GameData.HandPai.GetCanHu()
	if !canHu {
		return errors.New("不可以胡牌...")
	}


	/**
	HuPaiType_H_GangShangHua      HuPaiType = 9
	HuPaiType_H_GangShangPao      HuPaiType = 10
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

	hupai := u.GameData.HandPai.InPai
	if checkCase == nil {
		//表示是自摸
		isZimo = true
		outUserId = userId
		if u.GetPreMoGangInfo() != nil {
			isGangShangHua = true  //杠上花
			extraAct = mjproto.HuPaiType_H_GangShangHua
		}

	} else {
		isZimo = false                //表示是点炮
		outUserId = checkCase.GetUserIdOut()
		if checkCase.GetPreOutGangInfo() != nil {

			//这里需要判断是杠上炮，还是抢杠
			isGangShangPao = true//杠上炮
			extraAct = mjproto.HuPaiType_H_GangShangPao
		}
	}

	log.T("点炮的人[%v],胡牌的人[%v],杠上花[%v],杠上炮[%v],接下来开始getHuScore(%v,%v,%v,%v)", userId, outUserId, isGangShangHua, isGangShangPao,
		u.GameData.HandPai, isZimo, extraAct, roomInfo)

	fan, score, huCardStr := getHuScore(u.GameData.HandPai, isZimo, extraAct, roomInfo)
	log.T("胡牌(getHuScore)之后的结果fan[%v],score[%v],huCardStr[%v]", fan, score, huCardStr)

	//胡牌之后的信息
	hu := NewHuPaiInfo()
	*hu.GetUserId = u.GetUserId()
	*hu.SendUserId = outUserId
	//*hu.ByWho = 打牌的方位，对家，上家，下家？
	*hu.HuType = int32(extraAct)        ////杠上炮 杠上花 抢杠 海底捞 海底炮 天胡 地胡
	*hu.HuDesc = strings.Join(huCardStr, " ");
	*hu.Fan = fan
	*hu.Score = score        //只是胡牌的分数，不是赢了多少钱
	hu.Pai = hupai


	//胡牌之后，设置用户的数据
	u.GameData.HuInfo = append(u.GameData.HuInfo, hu)
	u.GameData.HandPai.HuPais = append(u.GameData.HandPai.HuPais, hu.Pai)        //增加胡牌
	u.SetStatus(MJUSER_STATUS_HUPAI)

	/**
		处理抢杠的逻辑，抢杠的逻辑需要特殊处理...
		1,首先是清楚杠牌的info
		2,增加碰牌
		3,删除杠牌的账单
	 */

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
			if pai != nil && pai.GetClientId() == hupai.GetClientId() {
				gangKeys = append(gangKeys, pai.GetIndex())        //需要删除的杠牌
				if pai.GetIndex() != hupai.GetIndex() {
					dianUser.GameData.HandPai.PengPais = append(dianUser.GameData.HandPai.PengPais, pai)        //碰牌
				}
			}
		}

		//删除杠牌的信息
		dianUser.GameData.DelGangInfo(hupai)
		//删除杠牌的账单
		for _, billUser := range d.GetUsers() {
			//处理每一个人的账单,并且减去amount
			_, bean := billUser.DelBillBean(hupai)
			billUser.SubBillAmount(bean.GetAmount())
		}
	}


	/**
		增加账单
		//todo 这里需要完善算账的逻辑逻辑,目前就自摸和点炮来做
	 */

	log.T("玩家[%v]胡牌，开始处理计算分数的逻辑...", userId)
	if isZimo {
		//如果是自摸的话，三家都需要给钱
		for _, shuUser := range d.GetUsers() {
			if shuUser != nil  && shuUser.IsGaming() {
				//用户赢钱的账户
				bill := NewBillBean()
				*bill.UserId = u.GetUserId()
				*bill.OutUserId = shuUser.GetUserId()
				*bill.Type = 1
				*bill.Des = "用户自摸，获得收入"
				*bill.Amount = hu.GetScore()        //杠牌的收入金额
				bill.Pai = hupai
				u.AddBillBean(bill)
				u.AddBillAmount(bill.GetAmount())

				//用户输钱的账单
				shubill := NewBillBean()
				*shubill.UserId = shuUser.GetUserId()
				*shubill.OutUserId = u.GetUserId()
				*shubill.Type = 1
				*shubill.Des = "用户自摸，输钱"
				*shubill.Amount = -hu.GetScore()       //杠牌的收入金额
				shubill.Pai = hupai
				shuUser.AddBillBean(shubill)
				shuUser.SubBillAmount(shubill.GetAmount())
			}
		}

	} else {
		//点炮胡牌成功之后的处理... 处理checkCase
		d.SetActiveUser(userId)        // 胡牌之后 设置当前操作的用户为当前胡牌的人...
		d.CheckCase.UpdateCheckBeanStatus(userId, CHECK_CASE_BEAN_STATUS_CHECKED)        // update checkCase...
		d.CheckCase.UpdateChecStatus(CHECK_CASE_STATUS_CHECKING_HUED)        //已经有人胡了，后边的人就不能碰或者杠了

		//如果是点炮的话，只有一家需要给钱...
		shuUser := d.GetUserByUserId(outUserId)
		bill := NewBillBean()
		*bill.UserId = u.GetUserId()
		*bill.OutUserId = shuUser.GetUserId()
		*bill.Type = 1
		*bill.Des = "用户自摸，获得收入"
		*bill.Amount = hu.GetScore()        //杠牌的收入金额
		bill.Pai = hupai
		u.AddBillBean(bill)
		u.AddBillAmount(bill.GetAmount())


		//用户输钱的账单
		shubill := NewBillBean()
		*shubill.UserId = shubill.GetUserId()
		*shubill.OutUserId = u.GetUserId()
		*shubill.Type = 1
		*shubill.Des = "用户自摸，输钱"
		*shubill.Amount = -d.GetBaseValue()        //杠牌的收入金额
		shubill.Pai = hupai
		shuUser.AddBillBean(shubill)
		shuUser.SubBillAmount(shubill.GetAmount())
	}

	//发送胡牌成功的回复
	ack := newProto.NewGame_AckActHu()
	*ack.HuType = hu.GetHuType()
	*ack.UserIdIn = hu.GetGetUserId()
	*ack.UserIdOut = hu.GetSendUserId()
	ack.HuCard = hu.Pai.GetCardInfo()
	log.T("给用户[%v]广播胡牌的ack[%v]", hu.GetGetUserId(), ack)
	//u.WriteMsg(ack)
	d.BroadCastProto(ack)

	return nil
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
		gangType = GANG_TYPE_MING        //明杠
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


	if gangType == GANG_TYPE_BA {
		log.T("用户[%v]杠牌是巴杠,现在处理巴杠...", userId)
		//循环碰牌来处理
		var pengKeys []int32
		user.GameData.HandPai.GangPais = append(user.GameData.HandPai.GangPais, gangPai)
		for _, pengPai := range user.GameData.HandPai.PengPais {
			if pengPai != nil && pengPai.GetValue() == gangPai.GetValue() && pengPai.GetFlower() == gangPai.GetFlower() {
				//增加杠牌
				user.GameData.HandPai.GangPais = append(user.GameData.HandPai.GangPais, pengPai)
				pengKeys = append(pengKeys, pengPai.GetIndex())
			}
		}

		//删除碰牌
		//减少手中的杠牌
		for _, key := range pengKeys {
			user.GameData.HandPai.DelPengPai(key)
		}

		//用户ba杠次数+1
		//user.sat

	} else if gangType == GANG_TYPE_MING || gangType == GANG_TYPE_AN {
		log.T("用户[%v]杠牌不是巴杠 是 gangType[%v]...", userId, gangType)

		//杠牌的类型
		var gangKey []int32
		//增加杠牌
		user.GameData.HandPai.Pais = append(user.GameData.HandPai.Pais, user.GameData.HandPai.InPai)
		//如果不是摸的牌，而是手中本来就有的牌，那么需要把他移除
		for _, pai := range user.GameData.HandPai.Pais {
			if pai.GetFlower() == gangPai.GetFlower() && pai.GetValue() == gangPai.GetValue() {
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
	//info.ByWho

	user.GameData.GangInfo = append(user.GameData.GangInfo, info)

	//增加杠牌状态
	user.PreMoGangInfo = info


	//处理账单
	/**
		没有胡牌的人，都需要给钱
	 */

	//现在杠牌的逻辑是,没有胡牌的人都要给钱...
	if info.GetGangType() == GANG_TYPE_BA || info.GetGangType() == GANG_TYPE_AN || info.GetGangType() == GANG_TYPE_MING {
		//暗杠和巴杠的处理方式
		for _, ou := range d.GetUsers() {
			//不为nil 并且不是本人，并且没有胡牌
			if ou != nil && ou.GetUserId() != user.GetUserId() && ou.IsGaming() {

				//用户赢钱的账户
				bill := NewBillBean()
				*bill.UserId = user.GetUserId()
				*bill.OutUserId = ou.GetUserId()
				*bill.Type = MJUSER_BILL_TYPE_YING_GNAG
				*bill.Des = "用户杠牌，获得收入"
				*bill.Amount = d.GetBaseValue()        //杠牌的收入金额
				bill.Pai = gangPai
				user.AddBillAmount(bill.GetAmount())
				user.Bill.Bills = append(user.Bill.Bills, bill)

				//用户输钱的账单
				shubill := NewBillBean()
				*shubill.UserId = ou.GetUserId()
				*shubill.OutUserId = user.GetUserId()
				*shubill.Type = MJUSER_BILL_TYPE_SHU_GNAG
				*shubill.Des = "用户杠牌，获得收入"
				*shubill.Amount = d.GetBaseValue()        //杠牌的收入金额
				shubill.Pai = gangPai
				ou.SubBillAmount(bill.GetAmount())
				ou.Bill.Bills = append(ou.Bill.Bills, bill)
			}
		}

		//} else if info.GetGangType() == GANG_TYPE_MING {
		//明杠的处理方式
	}

	//杠牌之后的逻辑
	//1,设置inpai为nil
	user.GameData.HandPai.InPai = nil

	//todo 返回杠牌成功的逻辑，返回一个摸牌的overTurn
	result := newProto.NewGame_AckActGang()
	*result.GangType = user.PreMoGangInfo.GetGangType()
	*result.UserIdOut = user.PreMoGangInfo.GetSendUserId()
	*result.UserIdIn = user.GetUserId()

	//todo 暂时着这么处理的，之后需要修改...
	result.GangCard[0] = user.PreMoGangInfo.GetPai().GetCardInfo()
	result.GangCard[1] = user.PreMoGangInfo.GetPai().GetCardInfo()
	result.GangCard[2] = user.PreMoGangInfo.GetPai().GetCardInfo()
	result.GangCard[3] = user.PreMoGangInfo.GetPai().GetCardInfo()
	log.T("广播玩家[%v]杠牌[%v]之后的ack[%v]", user.GetUserId(), gangPai, result)
	d.BroadCastProto(result)

	//设置 判断的为nil
	d.CheckCase = nil

	///如果是巴杠，需要设置巴杠的判断
	if gangType == GANG_TYPE_BA {
		d.InitCheckCase(gangPai, user)
		if d.CheckCase == nil {
			log.T("巴杠没有人可以抢杠...")
		}
	}

	return nil
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
	//*win.CardTitle =user.GameData.get // 赢牌牌型信息( 如:"点炮x2 明杠x2 根x2 自摸 3番" )
	//user.Statisc.
	win.Cards = user.GetPlayerCard(true) //牌信息,true 表示要显示牌的信息...
	*win.IsDealer = (d.GetBanker() == user.GetUserId() )        //是否是庄家
	*win.HuCount = 1        //本局胡的次数(血流成河会多次胡)
	return win
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
	return d.GetCurrPlayCount() == 0
}