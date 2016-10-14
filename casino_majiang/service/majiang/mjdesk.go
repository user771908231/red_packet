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
	newUser.GameData = NewPlayerGameData()

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
			go u.WriteMsg(p)
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
func (d *MjDesk) GetPlayerInfo(receiveUserId uint32) []*mjproto.PlayerInfo {
	var players []*mjproto.PlayerInfo
	for _, user := range d.Users {
		if user != nil {
			if user.GetUserId() == receiveUserId {
				players = append(players, user.GetPlayerInfo(true))
			} else {
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
			user.GameData = NewPlayerGameData()        //初始化一个空的麻将牌
		}
	}

	//初始化桌子的信息

	//1,初始化庄的信息,如果目前没有庄，则设置房主为庄,如果有庄，则不用管，每局游戏借宿的时候，会设置下一局的庄
	if d.GetBanker() == 0 {
		*d.Banker = d.GetOwner()
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
			u.GameData.HandPai.Pais = d.AllMJPai[i * 13: (i + 1) * 13]
			*d.MJPaiCursor = int32((i + 1) * 13) - 1;
		}
	}

	//庄需要多发一张牌
	bankUser := d.GetBankerUser()
	bankUser.GameData.HandPai.AddPai(d.GetNextPai())

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
	for _, u := range d.GetUsers() {
		if u != nil {
			if u.GetUserId() == user.GetUserId() {
				//表示是自己，可以看到手牌
				dealCards.PlayerCard = append(dealCards.PlayerCard, u.GetPlayerCard(true))
			} else {
				dealCards.PlayerCard = append(dealCards.PlayerCard, u.GetPlayerCard(false))
				//表示不是自己，不能看到手牌
			}

		}

	}

	return dealCards
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
	user.SetStatus(MJUSER_STATUS_DINGQUE)        //设置目前的状态是已经定缺
	*user.GameData.HandPai.QueFlower = color

	//回复定缺成功的消息
	ack := newProto.NewGame_DingQue()
	//*ack.Color = color
	*ack.Color = -1
	*ack.MatchId = d.GetRoomId()
	*ack.TableId = d.GetDeskId()
	*ack.UserId = userId
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
		if checkUser != nil && checkUser.GetUserId() != outUser.GetUserId() {
			bean := checkUser.GetCheckBean(p)
			if bean != nil {
				checkCase.CheckB = append(checkCase.CheckB, bean)
			}
		}
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
		//直接跳转到下一个操作的玩家...,这里表示判断已经玩了...
		d.CheckCase = nil
		d.SendMopaiOverTurn(gangUser)
		return nil
	} else {
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
		d.GetUserByUserId(caseBean.GetUserId()).SendOverTurn(overTurn)

		return nil

	}

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

	//判断是否可以胡牌

	//保存数据

	//发送结束的广播


	//判断牌局结束(整场游戏结束)

	return nil
}

func (d *MjDesk) SetActiveUser(userId uint32) error {
	*d.ActiveUser = userId
	return nil
}

//得到下一个摸牌的人...
func (d *MjDesk) GetNextMoPaiUser() *MjUser {
	var activeUser *MjUser = nil
	activeIndex := 0
	for i, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == d.GetActiveUser() {
			activeIndex = i
			break
		}
	}

	for i := activeIndex + 1; i < i + int(d.GetUserCount()); i++ {

		user := d.GetUsers()[i / int(d.GetUserCount())]
		if user != nil && user.IsNotHu() {
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
		return d.AllMJPai[0]
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

	d.SetActiveUser(user.GetUserId())        //用户摸牌之后，设置当前活动玩家为摸牌的玩家

	//发送广播
	overTrun := newProto.NewGame_OverTurn()
	*overTrun.UserId = user.GetUserId()                //这个是摸牌的，所以是广播...
	*overTrun.ActType = OVER_TURN_ACTTYPE_MOPAI        //摸牌

	//发送给当事人时候的信息
	overTrun.ActCard = d.GetNextPai().GetCardInfo()
	*overTrun.CanHu = user.GameData.HandPai.GetCanHu()
	*overTrun.CanGang = user.GameData.HandPai.GetCanGang()
	*overTrun.CanPeng = user.GameData.HandPai.GetCanPeng()
	user.SendOverTurn(overTrun)

	//给其他人广播协议
	*overTrun.CanHu = false
	*overTrun.CanPeng = false
	*overTrun.CanGang = false
	overTrun.ActCard = nil
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

	//todo 需要判断是否是可以碰


	//2.1开始碰牌的操作
	var pengKeys []int32
	pengPai := d.CheckCase.CheckMJPai
	user.GameData.HandPai.PengPais = append(user.GameData.HandPai.PengPais, pengPai)        //碰牌
	pengKeys = append(pengKeys, pengPai.GetIndex())

	for _, pai := range user.GameData.HandPai.Pais {
		if pai != nil && pai.GetIndex() == pengPai.GetIndex() {
			user.GameData.HandPai.PengPais = append(user.GameData.HandPai.PengPais, pai)        //碰牌
			pengKeys = append(pengKeys, pai.GetIndex())
		}
	}

	//2.2 删除手牌
	for _, key := range pengKeys {
		user.GameData.HandPai.DelPai(key)
	}

	//3,生成碰牌信息
	//user.GameData.

	//4,处理 checkCase
	d.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), CHECK_CASE_BEAN_STATUS_CHECKED)
	d.CheckCase.UpdateChecStatus(CHECK_CASE_STATUS_CHECKED) //碰牌之后，checkcase处理完毕

	//5,发送碰牌的广播
	ack := newProto.NewGame_AckActPeng()
	*ack.UserIdOut = d.CheckCase.GetUserIdOut()
	//todo 临时处理
	ack.PengCard[0] = d.CheckCase.CheckMJPai.GetCardInfo()
	ack.PengCard[1] = d.CheckCase.CheckMJPai.GetCardInfo()
	ack.PengCard[2] = d.CheckCase.CheckMJPai.GetCardInfo()
	d.BroadCastProto(ack)

	return nil
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

	//玩家胡牌
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("服务器错误：没有找到胡牌的user[%v]", userId)
		return errors.New("服务器错误，没有找到胡牌的user")
	}

	err := user.ActHu(d.CheckCase.CheckMJPai, d.CheckCase.GetUserIdOut())
	if err != nil {
		//如果这里出现胡牌失败，证明是系统有问题...
		log.E("用户[%v]胡牌失败...", userId)
		//result.Header = newProto.ErrorHeader()
		//user.WriteMsg(result)        //返回失败的信息
		return nil
	}

	return nil
}


//杠牌   怎么判断是明杠，暗杠，巴杠...
func (d *MjDesk) ActGang(userId uint32) error {
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("用户[%v]没有找到杠牌失败...", userId)
		return nil
	}

	//todo 判断是不是可以杠，如果可以杠牌，那么就开始杠牌

	//开始杠牌
	err := user.Gang(d.CheckCase.CheckMJPai, d.CheckCase.GetUserIdOut())
	if err != nil {
		//杠牌失败，这里是非法请求，或者服务器错误...
	}

	//todo 返回杠牌成功的逻辑，返回一个摸牌的overTurn
	result := newProto.NewGame_AckActGang()
	*result.GangType = user.PreMoGangInfo.GetGangType()
	*result.UserIdOut = user.PreMoGangInfo.GetSendUserId()

	//todo 暂时着这么处理的，之后需要修改...
	result.GangCard[0] = user.PreMoGangInfo.GetPai().GetCardInfo()
	result.GangCard[1] = user.PreMoGangInfo.GetPai().GetCardInfo()
	result.GangCard[2] = user.PreMoGangInfo.GetPai().GetCardInfo()
	result.GangCard[3] = user.PreMoGangInfo.GetPai().GetCardInfo()

	d.BroadCastProto(result)
	return nil
}