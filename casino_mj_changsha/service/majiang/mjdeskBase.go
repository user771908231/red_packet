package majiang

import (
	mjproto        "casino_mj_changsha/msg/protogo"
	"casino_common/common/Error"
	"casino_common/common/log"
	"sync/atomic"
	"casino_common/utils/numUtils"
	"casino_mj_changsha/msg/funcsInit"
	"fmt"
	"sync"
	"casino_common/common/game"
	"casino_common/utils/rand"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/timer"
	"github.com/name5566/leaf/module"
	"time"
	"casino_common/common/sessionService"
	"casino_common/proto/ddproto"
	"casino_common/common/consts"
)

type MjDesk struct {
	*PMjDesk
	*game.GameDesk
	*module.Skeleton
	Users               []*MjUser
	sync.Mutex
	RobotEnterTimer     *timer.Timer
	overTurnTimer       *time.Timer //等待操作的timer
	lastAcTime          time.Time
	HuParser            HuParser            //胡牌的解析器
	BirdInfo            []*mjproto.BirdInfo //抓鸟的信息
	ChangShaPlayOptions *mjproto.ChangShaPlayOptions
}

///获取用户已知亮出台面的牌 包括自己手牌、自己和其他玩家碰杠牌、其他玩家outPais
//得到麻将牌的总张数
func (d *MjDesk) GetTotalMjPaiCount() int32 {
	return 36 * d.GetFangCountLimit();
}

func (d *MjDesk) GetUsers() []*MjUser {
	return d.Users
}

func (d *MjDesk) GetUsersApi() []game.GameUserApi {
	ret := make([]game.GameUserApi, 0)
	for _, u := range d.GetUsers() {
		if u != nil {
			ret = append(ret, u)
		}
	}

	return ret
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

func (d *MjDesk) SetAATUser(userId uint32, actType int32) {
	*d.ActiveUser = userId
	*d.ActUser = userId
	*d.ActType = actType
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

//通过index得到user
func (d *MjDesk) GetUserByIndex(index int) *MjUser {
	if index >= int(d.GetUserCountLimit()) && index < 0 {
		return nil
	}
	u := d.Users[index]
	if u == nil {
		return nil
	}
	return u
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
func (d *MjDesk) addUserBean(user *MjUser) error {
	//根据房间类型判断人数是否已满
	if d.IsPlayerEnough() {
		return Error.NewFailError("房间已满，加入桌子失败")
	}

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
		return Error.NewFailError("没有找到合适的位置，加入桌子失败")
	}
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
	return count
}

//玩家是否足够
func (d *MjDesk) IsPlayerEnough() bool {
	return d.GetUserCount() == d.GetUserCountLimit()
}

//判断是否是血流成河
func (d *MjDesk) IsXueLiuChengHe() bool {
	return d.GetMjRoomType() == int32(mjproto.MJRoomType_roomType_xueLiuChengHe)
}

//判断是否是两人游戏
func (d *MjDesk) IsLiangRen() bool {
	return d.GetUserCountLimit() == 2
}

//判断是否是两房游戏
func (d *MjDesk) IsLiangFang() bool {
	return d.GetFangCountLimit() == 2
}

//是否是长沙麻将
func (d *MjDesk) IsChangShaMaJiang() bool {
	return d.GetRoomTypeInfo().GetMjRoomType() == mjproto.MJRoomType_roomType_changSha
}

//判断是否是倒倒胡
func (d *MjDesk) IsDaodaohu() bool {
	//倒倒胡，长沙麻将默认为倒倒胡
	if mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_daoDaoHu ||
		mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_changSha {
		return true
	}
	return false
}

//判断是否是两人两房
func (d *MjDesk) IsLiangRenLiangFang() bool {
	if mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_liangRenLiangFang {
		return true
	}
	return false
}

//判断是否是两人三房
func (d *MjDesk) IsLiangRenSanFang() bool {
	if mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_liangRenSanFang {
		return true
	}
	return false
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

func (d *MjDesk) SendGameInfo(userId uint32, reconnect mjproto.RECONNECT_TYPE) {
	gameinfo := d.GetGame_SendGameInfo(userId, reconnect)
	log.T("[%v]用户[%v]进入房间,reconnect type [%v]之后", d.DlogDes(), userId, reconnect)
	d.BroadCastProto(gameinfo)

	//如果是重新进入房间，需要发送重近之后的处理
	if reconnect == mjproto.RECONNECT_TYPE_RECONNECT {
		time.Sleep(time.Second * 3)
		d.SendReconnectOverTurn(userId)
	}
}

func (d *MjDesk) SetNextBanker(userId uint32) {
	*d.NextBanker = userId
}

func (d *MjDesk) SetBanker(userId uint32) {
	*d.Banker = userId

}

func (d *MjDesk) AddCurrPlayCount() {
	atomic.AddInt32(d.CurrPlayCount, 1)
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

//把桌子的数据保存到redis中
/**
	需要调用的地方
	1,新增加一个桌子的时候
	2,
 */
func (d *MjDesk) updateRedis() error {
	err := UpdateMjDeskRedis(d) //保存数据到redis
	if err != nil {
		return err
	} else {
		return nil
	}
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

//判断下一个庄是否已经确定
func (d *MjDesk) IsNextBankerExist() bool {
	return d.GetNextBanker() > 0
}

//剩下的牌的数量
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
	//log.T("leftPai is %v Count is : %v", mjPai.GetDes(), count)
	return count
}

//获取用户已知亮出台面的牌 包括自己手牌、自己和其他玩家碰杠牌、其他玩家outPais
func (d *MjDesk) GetDisplayPais(user *MjUser) []*MJPai {
	//获取所有玩家的亮出台面的牌 outPais + pengPais + gangPais

	displayPais := []*MJPai{}
	for _, user := range d.GetUsers() {
		userHandPai := user.GetGameData().GetHandPai()

		if userHandPai.GetGangPais() != nil && len(userHandPai.GangPais) > 0 {
			displayPais = append(displayPais, userHandPai.GangPais...) //杠的牌
		}
		if userHandPai.GetPengPais() != nil && len(userHandPai.PengPais) > 0 {
			displayPais = append(displayPais, userHandPai.PengPais...) //碰的牌
		}
		if userHandPai.GetChiPais() != nil && len(userHandPai.ChiPais) > 0 {
			displayPais = append(displayPais, userHandPai.ChiPais...) //吃的牌
		}
		if userHandPai.GetOutPais() != nil && len(userHandPai.OutPais) > 0 {
			displayPais = append(displayPais, userHandPai.OutPais...) //打出去的牌
		}
	}

	//在亮出台面的牌中加入用户自己的手牌
	userHandPai := user.GetGameData().GetHandPai()
	displayPais = append(displayPais, userHandPai.InPai)
	displayPais = append(displayPais, userHandPai.Pais...)
	return displayPais
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
	var gameStatus mjproto.DeskGameStatus = mjproto.DeskGameStatus_INIT //默认状态
	switch d.GetStatus() {
	case MJDESK_STATUS_READY:
		gameStatus = mjproto.DeskGameStatus_INIT
	case MJDESK_STATUS_EXCHANGE:
		gameStatus = mjproto.DeskGameStatus_EXCHANGE
	case MJDESK_STATUS_DINGQUE:
		gameStatus = mjproto.DeskGameStatus_DINGQUE
	case MJDESK_STATUS_RUNNING:
		gameStatus = mjproto.DeskGameStatus_PLAYING
	case MJDESK_STATUS_QISHOUHU:
		gameStatus = mjproto.DeskGameStatus_PLAYING
	}
	return int32(gameStatus)
}

//返回玩家的数目
/**
	needInpai ： 是否需要把inpai去得到
 */
func (d *MjDesk) GetPlayerInfo(receiveUserId uint32) []*mjproto.PlayerInfo {
	var players []*mjproto.PlayerInfo = make([]*mjproto.PlayerInfo, len(d.Users))
	for i, user := range d.Users {
		if user != nil {
			showHand := (user.GetUserId() == receiveUserId) //是否需要显示手牌
			isOwner := ( d.GetOwner() == user.GetUserId())  //判断是否是房主
			//定缺的状态，并且用户是 用户是庄，那么就显示inpai
			info := user.GetPlayerInfo(showHand)
			*info.IsOwner = isOwner
			players[i] = info
		} else {
			players[i] = &mjproto.PlayerInfo{}
		}
	}
	return players
}

// 发送gameInfo的信息
func (d *MjDesk) GetGame_SendGameInfo(receiveUserId uint32, isReconnect mjproto.RECONNECT_TYPE) *mjproto.Game_SendGameInfo {
	gameInfo := newProto.NewGame_SendGameInfo()
	gameInfo.DeskGameInfo = d.GetDeskGameInfo()
	*gameInfo.SenderUserId = receiveUserId
	gameInfo.PlayerInfo = d.GetPlayerInfo(receiveUserId)
	*gameInfo.IsReconnect = isReconnect
	return gameInfo

}

//获取桌子的关键信息帮助打日志
func (d *MjDesk) DlogDes() string {
	if d == nil {
		return "[desk==nil]"
	}
	s := fmt.Sprintf("[desk-%v-r-%v]  ", d.GetDeskId(), d.GetCurrPlayCount())
	return s
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

//是否有空位
func (d *MjDesk) IsYoulKongWei() bool {
	return d.GetUserCount() == d.GetUserCountLimit()
}

//判断是不是朋友桌
func (d *MjDesk) IsFriend() bool {
	return d.GetRoomType() == ROOMTYPE_FRIEND
}

//
func (d *MjDesk) CheckActUser(userId uint32, actType int32) error {
	if d.GetActUser() != userId {
		return Error.NewError(consts.ACK_RESULT_ERROR, "不是这个玩家")
	}

	if !(ACTTYPE_PENG == actType && d.GetUserByUserId(userId).ActCheck.canPeng) {
		return Error.NewError(consts.ACK_RESULT_ERROR, "不能操作碰")
	}
	return nil
}

//通过庄来判断骰子的数目
func (d *MjDesk) GetDice1() int32 {
	return rand.Rand(1, 7)
}

func (d *MjDesk) GetDice2() int32 {
	return rand.Rand(1, 7)
}

func (d *MjDesk) ClearBreakUser() error {
	for _, u := range d.Users {
		if u != nil && u.GetIsBreak() {
			log.T("%v 强制踢掉短线的玩家%v", d.DlogDes(), u.GetUserId())
			d.rmUser(u) //金币场，一局结束之后踢掉断线的玩家

		}
	}
	return nil
}

func (d *MjDesk) ClearLeaveUser() error {
	for _, u := range d.Users {
		if u != nil && u.GetIsLeave() {
			log.T("%v 强制踢掉离开的玩家%v", d.DlogDes(), u.GetUserId())
			d.rmUser(u) //金币场，一局结束之后踢掉离开的玩家
		}
	}
	return nil
}

func (d *MjDesk) ClearRobotUser() error {
	for _, u := range d.Users {
		if u != nil && u.GetIsRobot() {
			log.T("%v 强制踢掉机器人的玩家%v", d.DlogDes(), u.GetUserId())
			d.rmUser(u) //金币场，一局结束之后踢掉机器人玩家
		}
	}
	return nil
}

//金币不足的时候，需要把玩家踢掉
func (d *MjDesk) ClearCoinInsufficient() error {
	for _, u := range d.Users {
		if u != nil && u.GetCoin() < d.GetCoinLimit() {
			log.T("%v 强制踢掉金币[%v]不足[%v]的玩家%v", d.DlogDes(), u.GetCoin(), d.GetCoinLimit(), u.GetUserId())
			d.rmUser(u) //金币场，一局结束之后踢掉金币不足的玩家
		}
	}
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

//是否需要自摸加底
func (d *MjDesk) IsNeedZiMoJiaDi() bool {
	return mjproto.MJOption(*d.GetRoomTypeInfo().GetPlayOptions().ZiMoRadio) == mjproto.MJOption_ZIMO_JIA_DI
}

//是否需要自摸加番
func (d *MjDesk) IsNeedZiMoJiaFan() bool {
	return mjproto.MJOption(*d.GetRoomTypeInfo().GetPlayOptions().ZiMoRadio) == mjproto.MJOption_ZIMO_JIA_FAN
}

func (d *MjDesk) BroadCastProto(msg proto.Message) {
	for _, u := range d.Users {
		if u != nil {
			u.WriteMsg(msg)
		}
	}
}

func (d *MjDesk) BroadCastProtoExclusive(msg proto.Message, userId uint32) {
	for _, u := range d.Users {
		if u != nil && u.GetUserId() != userId {
			u.WriteMsg(msg)
		}
	}
}

//获取桌子当前的庄家
func (d *MjDesk) GetBankerUser() *MjUser {
	return d.GetUserByUserId(d.GetBanker())
}

//判断玩家是否是桌子当前的庄家
func (d *MjDesk) IsBanker(u *MjUser) bool {
	banker := d.GetBankerUser()
	if banker == nil {
		return false
	}
	return banker.GetUserId() == u.GetUserId()
}

//是不是全部都定缺了
func (d *MjDesk) AllDingQue() bool {
	for _, user := range d.GetUsers() {
		if user != nil && !user.IsDingQue() {
			log.T("%v用户[%v]还没有缺牌，等待定缺之后庄家开始打牌...", d.DlogDes(), user.GetUserId())
			return false
		}
	}
	return true
}

func (d *MjDesk) GetGamingCount() int32 {
	var gamingCount int32 = 0 //正在游戏中的玩家数量
	for _, user := range d.GetUsers() {
		if user != nil && user.IsGaming() {
			gamingCount ++
		}
	}
	return gamingCount
}

//得到玩家的坐标
func (d *MjDesk) getIndexByUserId(userId uint32) int {
	for i, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == userId {
			return i
		}
	}
	return -1
}

//交换两张牌
func (d *MjDesk) ExchangeMJPai(index1, index2 int) {
	d.AllMJPai[index1], d.AllMJPai[index2] = d.AllMJPai[index2], d.AllMJPai[index1]
}

//得到最后一张麻将pai
func (d *MjDesk) GetLastMjPai() *MJPai {
	return d.AllMJPai[len(d.AllMJPai)-1]
}

//是否还有牌
func (d *MjDesk) HandPaiCanMo() bool {
	if d.GetRemainPaiCount() == 0 {
		return false
	} else {
		return true
	}
}

func (d *MjDesk) getPaiById(paiId int32) *MJPai {
	for _, pai := range d.AllMJPai {
		if pai != nil && pai.GetIndex() == paiId {
			return pai
		}
	}
	return nil

}

func PaiTypeEnum2IntArry(ps []mjproto.PaiType) []int32 {
	var result []int32
	for _, p := range ps {
		result = append(result, int32(p))
	}
	return result
}

func IntArry2PaiTypeEnum(ia []int32) []mjproto.PaiType {
	var result []mjproto.PaiType
	for _, i := range ia {
		result = append(result, mjproto.PaiType(i))
	}
	return result
}

//返回金币bean
func (d *MjDesk) GetUserCoinBeans() []*mjproto.UserCoinBean {
	var ret []*mjproto.UserCoinBean
	for _, u := range d.GetUsers() {
		b := &mjproto.UserCoinBean{
			UserId: u.UserId,
			Coin:   u.Coin,
		}
		ret = append(ret, b)
	}
	return ret
}

//删除玩家的session
func (d *MjDesk) RmUserSession(u *MjUser) error {
	if u == nil {
		return nil
	}
	//开始删除
	sessionService.DelSessionByKey(u.GetUserId(), d.GetRoomType(), int32(ddproto.CommonEnumGame_GID_MAHJONG), d.GetDeskId())
	return nil
}
