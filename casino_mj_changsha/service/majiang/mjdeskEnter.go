package majiang

import (
	mjproto        "casino_mj_changsha/msg/protogo"
	"github.com/name5566/leaf/gate"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_common/common/log"
	"casino_common/common/userService"
	"casino_common/common/Error"
	"errors"
	"casino_mj_changsha/msg/funcsInit"
	"casino_common/utils/agentUtils"
)

//朋友桌用户加入房间
/**
return  reconnect,error
 */
func (d *MjDesk) enterUser(userId uint32, a gate.Agent) error {
	log.T("锁日志: %v enterUser(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v enterUser(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//1设置agent
	userLeave := d.GetUserByUserId(userId)
	if userLeave != nil {
		log.T("%v玩家[%v]重新进入房间....", d.DlogDes(), userId)
		userLeave.UpdateAgent(a)
		userLeave.UpdateSession(int32(ddproto.COMMON_ENUM_GAMESTATUS_GAMING))
		d.SendGameInfo(userId, mjproto.RECONNECT_TYPE_RECONNECT) //断线重连
		return nil
	}

	//先判断房间人数是否已经足够,这样就可以不处理下边的逻辑了
	if d.GetUserCount() == d.GetUserCountLimit() {
		return ERR_ENTER_DESK //进入房间失败
	}

	//3,加入一个新用户
	redisUser := userService.GetUserById(userId) //新玩家进入房间的时候
	newUser := NewMjUser()
	newUser.sex = redisUser.GetSex()
	newUser.headUrl = redisUser.GetHeadUrl()
	newUser.openId = redisUser.GetOpenId()
	*newUser.NickName = redisUser.GetNickName()
	newUser.d = d
	*newUser.UserId = userId
	*newUser.RoomId = d.GetRoomId()
	*newUser.IsBanker = false
	*newUser.Status = MJUSER_STATUS_INTOROOM
	newUser.GameData = NewPlayerGameData()
	newUser.UpdateAgent(a)
	*newUser.RoomType = d.GetRoomType()
	*newUser.RoomPassword = d.GetPassword()
	*newUser.WaitTime = 30                                 //todo 现在默认显示30秒
	newUser.IsRobot = proto.Bool(a == nil)                 //是否是机器人
	newUser.ActTimeoutCount = proto.Int32(0)               //设置默认的超时次数为0
	newUser.MjuserChangShaConfig = &MjuserChangShaConfig{} //长沙杠
	newUser.ip = agentUtils.GetIP(a)
	*newUser.Coin = 0

	///加入到desk
	err := d.addUserBean(newUser)
	if err != nil {
		log.E("用户[%v]加入房间[%v]失败,errMsg[%v]", userId, d.GetDeskId(), err)
		return Error.NewFailError(Error.GetErrorMsg(err))
	} else {
		//加入房间成功,更新session  并且发送游戏数据
		newUser.UpdateSession(int32(ddproto.COMMON_ENUM_GAMESTATUS_GAMING))
		d.SendGameInfo(userId, mjproto.RECONNECT_TYPE_NORMAL) //不是断线重连

		//如果是金币长，进入房间之后，设置准备倒计时

		return nil
	}
}

/**
玩家开始进入房间
1,如果是朋友桌，不做操作
2,如果是金币场，到了时间之后增加机器人
 */
func (d *MjDesk) beginEnter() error {
	d.SetStatus(MJDESK_STATUS_READY) //桌子开始ready
	return nil
}

//设置准备倒计时

//强制没有准备的玩家退出
func (d *MjDesk) ForceOutReadyTimeOutUser(userId uint32) {
	defer Error.ErrorRecovery("强制踢出玩家")
	//因为此方法是在afterFun执行，所以需要上锁
	log.T("锁日志: %v ForceOutReadyTimeOutUser(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ForceOutReadyTimeOutUser(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//如果玩家没有准备，强制玩家退出
	readyUser := d.GetUserByUserId(userId)
	if readyUser != nil && !readyUser.GetReady() {
		log.T("%s,玩家[%v]超时没有准备，强制退出", d.DlogDes(), userId)
		err := d.rmUser(readyUser) //准备倒计时结束的时候,强制退出一个玩家
		if err != nil {
			log.E("%v玩家[%v]准备超时，强制退出的时候出错", d.DlogDes(), userId)
		} else {
			//发送强制离开的广播
			ack := new(ddproto.CommonBcKickout)
			ack.UserId = proto.Uint32(readyUser.GetUserId())
			ack.Type = ddproto.COMMON_ENUM_KICKOUT_K_TIMEOUT_NOTREADY_ENTERDESK.Enum()
			ack.Msg = proto.String("准备超时,退出房间")
			readyUser.WriteMsg(ack)
		}
	}
}

//机器人进入房间
func (d *MjDesk) enterRobot() {
	log.T("%v 开是添加机器人", d.DlogDes())
	//1,做异常处理
	defer Error.ErrorRecovery("添加机器人")

	//2,获取机器人
	robot := MjroomManagerIns.RobotManger.ExpropriationRobotByCoin2(d.GetCoinLimit(), d.GetCoinLimitUL())
	if robot == nil {
		log.E("[%v]添加机器人的时候，没有找到合适的机器人...", d.DlogDes())
		return
	}

	//3，加入房间
	err := d.enterUser(robot.GetId(), nil) //机器人进入房间
	if err != nil {
		//用户加入房间失败...
		log.E("机器人玩家[%v]加入房间失败errMsg[%v]", robot.GetId(), err)
		MjroomManagerIns.RobotManger.ReleaseRobots(robot.GetId())
	}

}

//发送重新连接之后的overTurn
func (d *MjDesk) SendReconnectOverTurn(userId uint32) error {
	log.T("[%v]开始处理 sendReconnectOverTurn(%v),当前desk.status(%v),", d.DlogDes(), userId, d.GetStatus())

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
	} else if d.IsExchange() && user.IsNotExchange() {
		//给玩家发送换牌的消息...
		log.T("sendReconnectOverTurn，给user[%v]发送换牌的消息....", userId)
		beginExchange := newProto.NewGame_BroadcastBeginExchange()
		*beginExchange.Reconnect = true
		user.WriteMsg(beginExchange)

	} else if d.IsDingQue() && user.IsNotDingQue() {
		log.T("sendReconnectOverTurn，给user[%v]发送定缺的消息....", userId)

		//给玩家发送定缺的信息...
		beginQue := newProto.NewGame_BroadcastBeginDingQue()
		*beginQue.Reconnect = true
		user.WriteMsg(beginQue)

	} else if d.IsGaming() && user.GetUserId() == d.GetActUser() {

		//游戏中的情况，发送act的消息,这里需要更具当前的状态来发送overTurn
		if d.GetActType() == MJDESK_ACT_TYPE_MOPAI { //断线重连之后，如果是摸牌的状态则开发摸牌
			log.T("sendReconnectOverTurn，给user[%v]发送摸牌的消息....", userId)
			overTrun := d.GetMoPaiOverTurn(user, true) //重新进入房间之后
			user.WriteMsg(overTrun)
			log.T("玩家重新进入游戏之后 [%v]开始摸牌【%v】...", user.GetUserId(), overTrun)

		} else if d.GetActType() == MJDESK_ACT_TYPE_BANK_FIRST_MOPAI {
			log.T("sendReconnectOverTurn，给user[%v]发送摸牌的消息....", userId)
			overTrun := d.GetMoPaiOverTurn(user, true) //重新进入房间之后
			user.WriteMsg(overTrun)
			log.T("玩家重新进入游戏之后 [%v]开始摸牌【%v】...", user.GetUserId(), overTrun)
		} else if d.GetActType() == MJDESK_ACT_TYPE_DAPAI {
			//发送打牌的协议:这里只有在碰的时候才会出现这种情况，其他的时候都是摸牌(摸牌，刚拍之后摸牌)
			log.T("sendReconnectOverTurn，给user[%v]发送打牌的消息....", userId)
			overTurn := newProto.NewGame_OverTurn()
			*overTurn.UserId = userId
			*overTurn.CanGang = false
			*overTurn.CanPeng = false
			*overTurn.CanHu = false
			overTurn.CanBu = proto.Bool(false)
			*overTurn.ActType = OVER_TURN_ACTTYPE_MOPAI
			*overTurn.PaiCount = d.GetRemainPaiCount()
			overTurn.ActCard = NewBackPai()
			user.WriteMsg(overTurn)
		} else if d.GetActType() == MJDESK_ACT_TYPE_WAIT_CHECK {
			log.T("sendReconnectOverTurn，给user[%]发送checkCase的消息....", userId)

			caseBean := d.CheckCase.GetBeanByUserIdAndStatus(user.GetUserId(), CHECK_CASE_BEAN_STATUS_CHECKING)
			if caseBean == nil {
				log.E("没有找到玩家[%v]对应的checkBean", user.GetUserId())
				return errors.New("玩家重新进入房间发送check overturn的时候出错")
			}

			//找到需要判断bean之后，发送给判断人	//send overTurn
			overTurn := d.GetOverTurnByCaseBean(d.CheckCase.CheckMJPai, caseBean, OVER_TURN_ACTTYPE_OTHER) //重新进入游戏
			*overTurn.Time = 15                                                                            //紧急修改
			///发送overTurn 的信息
			log.T("开始发送玩家[%v]断线重连的overTurn[%v]", user.GetUserId(), overTurn)
			user.WriteMsg(overTurn)
		} else if d.GetActType() == MJDESK_ACT_TYPE_WAIT_HAIDI {
			log.T("sendReconnectOverTurn，给user[%]发送判定海底牌的消息....", userId)
			d.enquireHaiDi(user)
		} else if d.GetActType() == MJDESK_ACT_TYPE_WAIT_CHECK_CHANGSHAGANG {
			log.T("sendReconnectOverTurn，给user[%]发送长沙杠之后的消息....", userId)

			inpai1 := user.GetGameData().GetHandPai().GetInPai()
			inpai2 := user.GetGameData().GetHandPai().GetInPai2()

			user.GameData.HandPai.InPai = inpai1
			overTrun1 := d.GetMoPaiOverTurn(user, false) //长沙玩家杠牌之后，用户摸牌的时候,发送一个用户摸牌的overturn
			overTrun1.CanGang = proto.Bool(false)
			overTrun1.GangCards = nil
			overTrun1.CanBu = proto.Bool(false)
			overTrun1.BuCards = nil

			user.GameData.HandPai.InPai2 = inpai1
			user.GameData.HandPai.InPai = inpai2
			overTrun2 := d.GetMoPaiOverTurn(user, false) //长沙玩家杠牌之后，用户摸牌的时候,发送一个用户摸牌的overturn
			overTrun2.CanGang = proto.Bool(false)
			overTrun2.GangCards = nil
			overTrun2.CanBu = proto.Bool(false)
			overTrun2.BuCards = nil

			/**
				注意：杠牌和补牌的差别...如果是杠牌，摸两张，然后自动打出去...如果补，那么和普通麻将的杠是一样的...
			 */
			ack := &mjproto.Game_ChangShaOverTurnAfterGang{}
			ack.GangPai = append(ack.GangPai, user.GetGameData().GetHandPai().GetInPai().GetCardInfo(), user.GetGameData().GetHandPai().GetInPai2().GetCardInfo())
			ack.Header = newProto.NewHeader()
			ack.Header.UserId = proto.Uint32(user.GetUserId())
			if overTrun1.GetCanHu() {
				ack.CanHu = proto.Bool(true)
				ack.CanGuo = proto.Bool(true)
				ack.HuCards = append(ack.HuCards, overTrun1.ActCard)
			}

			if overTrun2.GetCanHu() {
				ack.CanHu = proto.Bool(true)
				ack.CanGuo = proto.Bool(true)
				ack.HuCards = append(ack.HuCards, overTrun2.ActCard)
			}
			log.T("[%v][%v]断线重连，长沙杠之后，开摸的牌【%v】...", d.DlogDes(), user.GetUserId(), user.UserPai2String(), ack)
			d.BroadCastProto(ack)

		}
	} else if d.GetStatus() == MJDESK_STATUS_QISHOUHU {
		//长沙麻将起手胡牌的阶段...查看玩家是否需要发送起手胡牌
		bean := d.GetCheckCase().GetNextBean()
		if bean != nil && bean.GetUserId() == user.GetUserId() && bean.GetCanHu() {
			//给玩家发送起手胡的ack
			overTurn := &mjproto.Game_ChangshQiShouHuOverTurn{
				Header: &mjproto.ProtoHeader{
					UserId: bean.UserId,
				},
			}
			log.T("短线重连之后，开始发送起手胡overTurn[%v]", overTurn)                     //打日志
			user.SendOverTurn(overTurn)                                       //发送OverTurn
			d.SetActUserAndType(bean.GetUserId(), MJDESK_ACT_TYPE_WAIT_CHECK) //断线重连 起手胡牌
		}
	}

	log.T("%v开始处理 sendReconnectOverTurn(%v),当前desk.status(%v)----处理完毕...", d.DlogDes(), userId, d.GetStatus())
	return nil
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
	typeInfo.ChangShaPlayOptions = d.ChangShaPlayOptions
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

//得到deskGameInfo
func (d *MjDesk) GetDeskGameInfo() *mjproto.DeskGameInfo {
	deskInfo := newProto.NewDeskGameInfo()
	//deskInfo.ActionTime
	//deskInfo.DelayTime
	*deskInfo.GameStatus = d.GetClientGameStatus()
	*deskInfo.CurrPlayCount = d.GetCurrPlayCount()   //当前第几局
	*deskInfo.TotalPlayCount = d.GetTotalPlayCount() //总共几局
	*deskInfo.PlayerNum = d.GetUserCount()           //玩家的人数
	deskInfo.RoomTypeInfo = d.GetRoomTypeInfo()
	*deskInfo.RoomNumber = d.GetPassword() //房间号码...
	*deskInfo.Banker = d.GetBanker()
	//deskInfo.NRebuyCount
	//deskInfo.InitRoomCoin
	*deskInfo.NInitActionTime = d.GetNInitActionTime()
	//deskInfo.NInitDelayTime
	*deskInfo.ActiveUserId = d.GetActiveUser()
	*deskInfo.RemainCards = d.GetRemainPaiCount()
	return deskInfo
}
