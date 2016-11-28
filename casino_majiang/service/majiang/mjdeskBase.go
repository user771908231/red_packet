package majiang

import (
	"casino_majiang/msg/protogo"
	"casino_common/common/Error"
	"casino_common/common/log"
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
	case mjproto.MJRoomType_roomType_sanRenLiangFang :
		d.Users = make([]*MjUser, 3)
	default :
		d.Users = make([]*MjUser, 4)
	}
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
	case d.IsSanRenLiangFang() && d.GetUserCount() == 3 : //是三人两房并且玩家数等于3
		isPlayerEnough = true
	case !d.IsSanRenLiangFang() && d.GetUserCount() == 4 : //不是三人两房并且玩家数等于4
		isPlayerEnough = true
	default:
		isPlayerEnough = false
	}
	return isPlayerEnough
}


//判断是否是三人两房
func (d *MjDesk) IsSanRenLiangFang() bool {
	return mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_sanRenLiangFang
}

//判断是否是四人两房
func (d *MjDesk) IsSiRenLiangFang() bool {
	if mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_siRenLiangFang {
		return true
	}
	return false
}

//判断是否是倒倒胡
func (d *MjDesk) IsDaodaohu() bool {
	if mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_daoDaoHu {
		return true
	}
	return false
}

//判断是否是倒倒胡
func (d *MjDesk) IsLiangRenLiangFang() bool {
	if mjproto.MJRoomType(d.GetMjRoomType()) == mjproto.MJRoomType_roomType_daoDaoHu {
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

//判断是否是血流成河
func (d *MjDesk) IsXueLiuChengHe() bool {
	return d.GetMjRoomType() == int32(mjproto.MJRoomType_roomType_xueLiuChengHe)
}

//判断是否是血战到底
func (d *MjDesk) IsXueZhanDaoDi() bool {
	return d.GetMjRoomType() == int32(mjproto.MJRoomType_roomType_xueZhanDaoDi)
}

func (d *MjDesk) SendGameInfo(userId uint32, reconnect mjproto.RECONNECT_TYPE) {
	gameinfo := d.GetGame_SendGameInfo(userId, reconnect)
	log.T("用户[%v]进入房间,reconnect[%v]之后，返回的数据gameInfo[%v]", userId, reconnect, gameinfo)
	d.BroadCastProto(gameinfo)
}



