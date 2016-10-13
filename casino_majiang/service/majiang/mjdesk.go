package majiang

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_majiang/service/AgentService"
)

//状态表示的是当前状态.
var MJDESK_STATUS_CREATED = 1 //刚刚创建
var MJDESK_STATUS_READY = 2//正在准备
var MJDESK_STATUS_ONINIT = 3//准备完成之后，desk初始化数据
var MJDESK_STATUS_EXCHANGE = 4//desk初始化完成之后，告诉玩家可以开始换牌
var MJDESK_STATUS_DINGQUE = 5//换牌结束之后，告诉玩家可以开始定缺
var MJDESK_STATUS_RUNNING = 6 //定缺之后，开始打牌
var MJDESK_STATUS_LOTTERY = 7 //结算
var MJDESK_STATUS_END = 8//一局结束


var OVER_TURN_ACTTYPE_MOPAI int32 = 1; //摸牌的类型...
var OVER_TURN_ACTTYPE_OTHER int32 = 2; //碰OTHER

//判断是不是朋友桌
func (d *MjDesk) IsFriend() bool {
	return true
}

//用户加入房间
func (d *MjDesk) addNewUser(userId uint32, a gate.Agent) error {
	if d.IsFriend() {
		return d.addNewUserFriend(userId, a)
	} else {
		return nil
	}
}

//朋友桌用户加入房间
func (d *MjDesk) addNewUserFriend(userId uint32, a gate.Agent) error {
	//1,是否是重新进入
	user := d.GetUserByUserId(userId)
	if user != nil {
		//是断线重连
		*user.IsBreak = false;
		AgentService.SetAgent(userId, a)
		return nil

	}

	//2,是否是离开之后重新进入房间
	userLeave := d.getLeaveUserByUserId(userId)
	if userLeave != nil {
		err := d.addUser(userLeave)
		if err != nil {
			//加入房间的时候错误
			log.E("已经离开的用户[%v]重新加入desk[%v]的时候出错errMsg[%v]", userId, d.GetDeskId(), err)
			return errors.New("已经离开的用户重新加入房间的时候出错")
		} else {
			*userLeave.IsBreak = false
			*userLeave.IsLeave = false
			d.rmLeaveUser(userLeave.GetUserId())
			return nil
		}
	}

	//加入一个新用户
	newUser := NewMjUser()
	*newUser.UserId = userId
	*newUser.DeskId = d.GetDeskId()
	*newUser.RoomId = d.GetRoomId()
	*newUser.Coin = d.GetBaseValue()
	*newUser.IsBreak = false
	*newUser.IsLeave = false
	*newUser.Status = MJUSER_STATUS_INTOROOM
	newUser.MJHandPai = NewMJHandPai()

	//设置agent
	AgentService.SetAgent(userId, a)
	err := d.addUser(newUser)
	if err != nil {
		log.E("用户[%v]加入房间[%v]失败,errMsg[%v]", userId, err)
		return errors.New("用户加入房间失败")
	} else {
		//加入房间成功
		return nil
	}
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
		if u != nil {
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
	//deskInfo.ActiveSeat
	//deskInfo.DelayTime
	//deskInfo.GameStatus
	*deskInfo.CurrPlayCount = d.GetCurrPlayCount() //当前第几局
	*deskInfo.TotalPlayCount = d.GetTotalPlayCount()//总共几局
	*deskInfo.PlayerNum = d.GetPlayerNum()        //玩家的人数
	deskInfo.RoomTypeInfo = d.GetRoomTypeInfo()
	//deskInfo.NRebuyCount
	//deskInfo.InitRoomCoin
	//deskInfo.NInitActionTime
	//deskInfo.NInitDelayTime
	return deskInfo
}

//返回玩家的数目
func (d *MjDesk) GetPlayerInfo() []*mjproto.PlayerInfo {
	var players []*mjproto.PlayerInfo
	for _, user := range d.Users {
		if user != nil {
			players = append(players, user.GetPlayerInfo())
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
func (d *MjDesk) GetGame_SendGameInfo() *mjproto.Game_SendGameInfo {
	gameInfo := newProto.NewGame_SendGameInfo()
	gameInfo.DeskGameInfo = d.GetDeskGameInfo()
	gameInfo.PlayerInfo = d.GetPlayerInfo()
	//gameInfo.SenderUserId   发起请求的人 ... agent 返回信息的时候 取userId

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

	//设置为准备的状态
	*user.Status = MJUSER_STATUS_READY        //用户准备

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


//得到当前桌子的人数..
func (d *MjDesk) GetUserCount() int32 {
	var count int32 = 0
	for _, user := range d.Users {
		if user != nil {
			count ++
		}
	}
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
			user.MJHandPai = NewMJHandPai()        //初始化一个空的麻将牌
		}
	}

	//发送游戏开始的协议...
	log.T("发送游戏开始的协议..")
	open := newProto.NewGame_Opening()
	d.BroadCastProto(open)
	return nil
}

/**
	初始化牌相关的信息
 */
func (d *MjDesk) initCards() error {
	//得到一副已经洗好的麻将
	d.AllMJPai = XiPai()

	//给每个人初始化...
	for i, u := range d.Users {
		if u != nil && u.IsGaming() {
			log.T("开始给你玩家[%v]初始化手牌...", u.GetUserId())
			u.MJHandPai.Pais = d.AllMJPai[i * 13: (i + 1) * 13]
			*d.MJPaiNexIndex = int32((i + 1) * 13);
		}
	}


	//发牌的协议game_DealCards  初始化完成之后，给每个人发送牌
	for _, user := range d.Users {
		if user != nil {
			dealCards := user.GetDealCards()
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

//开始定缺
func (d *MjDesk) beginDingQue() error {
	//给每个人发送开始定缺的信息
	beginQue := newProto.NewGame_BroadcastBeginDingQue()
	log.T("开始给玩家发送开始定缺的广播[%v]", beginQue)
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

	//设置定缺
	*user.DingQue = true
	*user.MJHandPai.DingQueColor = color

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
func (d *MjDesk) InitCheckCase(p *MJPai) error {
	return nil
}

//执行判断事件
/**

	这里仅仅是用于判断打牌之后别人的碰杠胡

 */
func (d *MjDesk) DoCheckCase() error {

	//检测参数
	if d.CheckCase == nil || d.CheckCase.IsChecked() {
		//直接跳转到下一个操作的玩家...,这里表示判断已经玩了...
		return errors.New("")
	}

	//switch d.CheckCase.GetCheckStatus() {
	//case CHECK_CASE_STATUS_0:
	//
	///**
	//	0 表示没有进行判断过，优先选择胡牌的case来进行判断，如果没有找到胡牌的case,那么找到其他的bean来进行判断
	// */
	//case CHECK_CASE_STATUS_1:
	//
	///**
	//	1 表示已经进行过胡牌的事件进行来判断，接下来进行另一个胡牌的判断，其他的不判断了
	// */
	//
	//case CHECK_CASE_STATUS_2:
	//
	///**
	//	2 表示全部都已经判断完了...s
	// */
	//
	//
	//}

	//1,找到胡牌的人来进行处理
	var caseBean *CheckBean
	for _, bean := range d.CheckCase.CheckB {
		if bean != nil && !bean.IsChecked() && bean.GetCanHu() {
			caseBean = bean
			break
		}
	}

	//如果这里的caseBean ！=nil 表示还有可以胡牌的人没有进行判定
	if caseBean == nil {
		for _, bean := range d.CheckCase.CheckB {
			if bean != nil && !bean.IsChecked() && !bean.GetCanHu() {
				caseBean = bean
				break
			}
		}
	}

	if caseBean == nil {
		log.E("服务器错误....这里不应该出现的...checkCae[%v]", d.CheckCase)
		return errors.New("已经没有需要处理的了")

	}

	//找到需要判断bean之后，发送给判断人	//send overTurn
	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = caseBean.GetUserId()
	*overTurn.CanGang = caseBean.GetCanGang()
	*overTurn.CanPeng = caseBean.GetCanPeng()
	*overTurn.CanHu = caseBean.GetCanHu()
	overTurn.ActCard = d.CheckCase.CheckMJPai.GetCardInfo()        //
	*overTurn.ActType = OVER_TURN_ACTTYPE_OTHER

	///发送overTurn 的信息
	d.GetUserByUserId(caseBean.GetUserId()).SendOverTurn(overTurn)

	return nil
}


/**
	1，只剩一个玩家没有胡牌
	2, 已经没有牌了...
 */

func (d *MjDesk) Time2Lottery() bool {
	return false
}

// 一盘麻将结束....这里需要针对每个人结账...并且对desk和user的数据做清楚...
func (d *MjDesk) Lottery() error {
	//结账需要分两中情况
	/**
		1，只剩一个玩家没有胡牌的时候
		2，没有生育麻将的时候.需要分别做处理...
	 */

	return nil
}

func (d *MjDesk) SetNestUserCursor(userId uint32) error {
	*d.NextUserCursor = userId
	return nil
}

//得到下一个摸牌的人...
func (d *MjDesk) GetNextMoPaiUser() *MjUser {

	return nil
}

//得到下一张牌...
func (d *MjDesk) GetNextPai() *MJPai {
	return nil
}


//发送摸牌的广播
//指定一个摸牌，如果没有指定，则系统通过游标来判断
func (d *MjDesk) SendMopaiOverTurn(user *MjUser) error {
	if user == nil {
		user = d.GetNextMoPaiUser()
	}
	overTrun := newProto.NewGame_OverTurn()
	*overTrun.UserId = user.GetUserId()                //这个是摸牌的，所以是广播...
	*overTrun.ActType = OVER_TURN_ACTTYPE_MOPAI        //摸牌
	*overTrun.CanHu = false        //这里需要判断之后得到...
	*overTrun.CanPeng = false
	*overTrun.CanGang = false
	overTrun.ActCard = d.GetNextPai().GetCardInfo() //得到下一张牌
	return nil
}

func (d *MjDesk) GetDingQueEndInfo() *mjproto.Game_DingQueEnd {
	end := newProto.NewGame_DingQueEnd()

	for _, u := range d.GetUsers() {
		if u != nil && u.MJHandPai != nil {
			bean := newProto.NewGame_DingQueEndBean()
			*bean.UserId = u.GetUserId()
			*bean.Flower = u.MJHandPai.GetDingQueColor()
			end.Ques = append(end.Ques, bean)
		}
	}
	return end
}

//某人胡牌...
func (d *MjDesk)ActHu(userId uint32) error {

	//对于杠，有摸牌前 杠的状态，有打牌前杠的状态

	//1,胡的牌是当前check里面的牌，如果没有check，则表示是自摸
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


	/**



	 */


	return nil
}