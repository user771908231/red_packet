package majiang

import (
	"casino_majiang/msg/protogo"
	"casino_common/common/Error"
	"casino_common/common/log"
	"sync/atomic"
	"casino_common/utils/numUtils"
	"casino_majiang/msg/funcsInit"
)

///获取用户已知亮出台面的牌 包括自己手牌、自己和其他玩家碰杠牌、其他玩家outPais
//得到麻将牌的总张数
func (d *MjDesk) GetTotalMjPaiCount() int32 {
	return 36 * d.GetFangCountLimit();
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

//根据房间类型初始化房间玩家数
func (d *MjDesk) InitUsers(mjRoomType mjproto.MJRoomType) {

	switch mjRoomType {
	case mjproto.MJRoomType_roomType_liangRenLiangFang,
		mjproto.MJRoomType_roomType_liangRenSanFang :
		//log.T("两人")
		d.Users = make([]*MjUser, 2)

	case mjproto.MJRoomType_roomType_sanRenSanFang,
		mjproto.MJRoomType_roomType_sanRenLiangFang:
		//log.T("三人")
		d.Users = make([]*MjUser, 3)

	case mjproto.MJRoomType_roomType_xueZhanDaoDi,
		mjproto.MJRoomType_roomType_daoDaoHu,
		mjproto.MJRoomType_roomType_deYangMaJiang,
		mjproto.MJRoomType_roomType_xueLiuChengHe,
		mjproto.MJRoomType_roomType_siRenLiangFang :
		//log.T("四人")
		d.Users = make([]*MjUser, 4)

	default :
		d.Users = make([]*MjUser, 4)
	}
}

//新增加一个玩家
func (d *MjDesk) addUser(user *MjUser) error {
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
	return count;

}

//玩家是否足够
func (d *MjDesk) IsPlayerEnough() (isPlayerEnough bool) {
	switch {
	case d.IsSanRen() && d.GetUserCount() == 3 : //是三人游戏并且玩家数等于3
		isPlayerEnough = true
	case d.IsSiRen() && d.GetUserCount() == 4 : //是四人游戏并且玩家数等于4
		isPlayerEnough = true
	case d.IsLiangRen() && d.GetUserCount() == 2 : //是两人游戏且玩家数等于2
		isPlayerEnough = true
	default:
		isPlayerEnough = false
	}
	return isPlayerEnough
}

//判断是否是血流成河
func (d *MjDesk) IsXueLiuChengHe() bool {
	return d.GetMjRoomType() == int32(mjproto.MJRoomType_roomType_xueLiuChengHe)
}

//判断是否是血战到底
func (d *MjDesk) IsXueZhanDaoDi() bool {
	return d.GetMjRoomType() == int32(mjproto.MJRoomType_roomType_xueZhanDaoDi)
}


//判断是否是三人两房
func (d *MjDesk) IsSanRenLiangFang() bool {
	return mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_sanRenLiangFang
}

//判断是否是三人三房
func (d *MjDesk) IsSanRenSanFang() bool {
	return mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_sanRenSanFang
}



//判断是否是四人两房
func (d *MjDesk) IsSiRenLiangFang() bool {
	return mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_siRenLiangFang
}

//判断是否是四人游戏
func (d *MjDesk) IsSiRen() bool {
	return d.IsSiRenLiangFang() || d.IsXueLiuChengHe() || d.IsXueZhanDaoDi() || d.IsDaodaohu()
}

//判断是否是三人游戏
func (d *MjDesk) IsSanRen() bool {
	return d.IsSanRenLiangFang() || d.IsSanRenSanFang()
}

//判断是否是两人游戏
func (d *MjDesk) IsLiangRen() bool {
	return d.IsLiangRenLiangFang() || d.IsLiangRenSanFang()
}

//判断是否是两房游戏
func (d *MjDesk) IsLiangFang() bool {
	return d.IsSanRenLiangFang() || d.IsSiRenLiangFang() || d.IsLiangRenLiangFang()
}

//判断是否是倒倒胡
func (d *MjDesk) IsDaodaohu() bool {
	if mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_daoDaoHu {
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

//通过房间type初始化人数和麻将房数
func (d *MjDesk) InitUserCountAndFangCountByType() {
	if d.IsSanRenLiangFang() {
		*d.UserCountLimit = 3
		*d.FangCountLimit = 2
	} else if d.IsSiRenLiangFang() {
		*d.UserCountLimit = 4
		*d.FangCountLimit = 2
	} else if d.IsLiangRenLiangFang() {
		*d.UserCountLimit = 2
		*d.FangCountLimit = 2
	} else {
		*d.UserCountLimit = 4
		*d.FangCountLimit = 3
	}

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
	log.T("用户[%v]进入房间,reconnect[%v]之后，返回的数据gameInfo[%v]", userId, reconnect, gameinfo)
	d.BroadCastProto(gameinfo)
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
func (d *MjDesk)updateRedis() error {
	err := UpdateMjDeskRedis(d)        //保存数据到redis
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
	if d.GetNextBanker() > 0 {
		return true
	} else {
		return false
	}
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

		switch {
		case userHandPai.GetGangPais() != nil:
			displayPais = append(displayPais, userHandPai.GangPais...) //杠的牌
		case userHandPai.GetPengPais() != nil:
			displayPais = append(displayPais, userHandPai.PengPais...) //碰的牌
		case userHandPai.GetOutPais() != nil:
			displayPais = append(displayPais, userHandPai.OutPais...) //打出去的牌
		default:

		}
	}

	//在亮出台面的牌中加入用户自己的手牌
	userHandPai := user.GetGameData().GetHandPai()
	displayPais = append(displayPais, userHandPai.InPai)
	displayPais = append(displayPais, userHandPai.Pais...)
	return displayPais
}

func (d *MjDesk) GetTransferredStatus() string {
	ret := ""
	switch d.GetStatus() {
	case MJDESK_STATUS_CREATED:
		ret = "创建成功"
	case MJDESK_STATUS_DINGQUE:
		ret = "开始定缺"
	case MJDESK_STATUS_END:
		ret = "单局结束"
	case MJDESK_STATUS_EXCHANGE:
		ret = "开始换牌"
	case MJDESK_STATUS_FAPAI:
		ret = "开始发牌"
	case MJDESK_STATUS_LOTTERY:
		ret = "开始结算"
	case MJDESK_STATUS_ONINIT:
		ret = "开始初始化数据"
	case MJDESK_STATUS_READY:
		ret = "开始准备"
	case MJDESK_STATUS_RUNNING:
		ret = "定缺后开始打牌"
	default:

	}
	return ret
}

func (d *MjDesk) GetDeskMJInfo() string {
	if d == nil || d.AllMJPai == nil {
		return "暂时没有初始化麻将"
	}
	s := ""
	for i, p := range d.AllMJPai {
		is, _ := numUtils.Int2String(int32(i))
		s = s + " (" + is + "-" + p.LogDes() + ")"
		if (i + 1) % 27 == 0 {

		}
	}
	return s
}

func (d *MjDesk) GetTransferredRoomType() string {
	ret := ""
	switch mjproto.MJRoomType(d.GetMjRoomType()) {
	case mjproto.MJRoomType_roomType_xueZhanDaoDi:
		ret = "血战到底"
	case mjproto.MJRoomType_roomType_sanRenLiangFang:
		ret = "三人两房"
	case mjproto.MJRoomType_roomType_siRenLiangFang:
		ret = "四人两房"
	case mjproto.MJRoomType_roomType_deYangMaJiang:
		ret = "德阳麻将"
	case mjproto.MJRoomType_roomType_daoDaoHu:
		ret = "倒倒胡"
	case mjproto.MJRoomType_roomType_xueLiuChengHe:
		ret = "血流成河"
	case mjproto.MJRoomType_roomType_liangRenLiangFang:
		ret = "两人两房"
	case mjproto.MJRoomType_roomType_liangRenSanFang:
		ret = "两人三房"
	case mjproto.MJRoomType_roomType_sanRenSanFang:
		ret = "三人三房"
	default:

	}
	return ret
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
/**
	needInpai ： 是否需要把inpai去得到
 */
func (d *MjDesk) GetPlayerInfo(receiveUserId uint32, isDingQue bool) []*mjproto.PlayerInfo {
	var players []*mjproto.PlayerInfo
	for _, user := range d.Users {
		if user != nil {
			showHand := (user.GetUserId() == receiveUserId)                //是否需要显示手牌
			isOwner := ( d.GetOwner() == user.GetUserId())                //判断是否是房主

			//定缺的状态，并且用户是 用户是庄，那么就显示inpai
			needInpai := false
			if isDingQue && user.GetUserId() == d.GetBanker() {
				needInpai = true

			}
			info := user.GetPlayerInfo(showHand, needInpai)
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
func (d *MjDesk) GetGame_SendGameInfo(receiveUserId uint32, isReconnect mjproto.RECONNECT_TYPE) *mjproto.Game_SendGameInfo {

	//如果是短线重连，并且玩家还没有换三张，或者处于定缺的状态，那么需要发送庄家的inpai
	isDingQue := false
	recUser := d.GetUserByUserId(receiveUserId)
	/**
		1,当前阶段处于定缺的阶段
		2，用户是庄稼的情况
		3，此时需要发送 inpai
	 */
	if isReconnect == mjproto.RECONNECT_TYPE_RECONNECT &&
		( d.IsExchange() || ( d.IsDingQue() && recUser.IsNotDingQue() )) {
		isDingQue = true
	}

	gameInfo := newProto.NewGame_SendGameInfo()
	gameInfo.DeskGameInfo = d.GetDeskGameInfo()
	*gameInfo.SenderUserId = receiveUserId
	gameInfo.PlayerInfo = d.GetPlayerInfo(receiveUserId, isDingQue)
	*gameInfo.IsReconnect = isReconnect
	return gameInfo

}
