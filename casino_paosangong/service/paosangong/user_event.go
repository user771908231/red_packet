package paosangong

import (
	"casino_common/common/consts"
	"casino_common/common/log"
	"casino_common/proto/ddproto"
	"casino_common/proto/funcsInit"
	"casino_common/utils/db"
	"casino_paosangong/conf/config"
	"errors"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"
	"casino_common/common/service/roomAgent"
	"casino_common/common/Error"
	"casino_common/common/service/countService"
	"casino_common/common/userService"
)
//入座
func (u *User) DoSiteDown(site_index int32) {
	//if site_index < 0 || site_index > u.Desk.DeskOption.GetMaxUser() - 1 {
	//	u.SendEnterDeskBcACK(-2, "座位号非法！")
	//	return
	//}

	//if !u.Desk.IsFreeSite(site_index) {
	//	u.SendEnterDeskBcACK(-3, "这个位置已经有人坐了！")
	//	return
	//}

	//if u.Desk.GetIsOnGamming() {
	//	u.SendEnterDeskBcACK(-5, "已开局，无法入座！")
	//	return
	//}

	if u.GetIndex() >= 0 {
		u.SendEnterDeskBcACK(-7, "您已经在座位上了，请不要重复入座！")
		return
	}

	desk := u.Desk

	free_site_index, err := desk.GetFreeSiteIndex()
	if err != nil {
		u.SendEnterDeskBcACK(-4, "已经没有空闲座位了！")
		return
	}

	//设置座位号
	*u.Index = int32(free_site_index)

	//自动准备
	//if u.Desk.DeskOption.GetAutoStartGammer() == 0 || u.Desk.GetOwner() != u.GetUserId() {
	//	if u.Desk.GetCircleNo() == 1 && u.Desk.DeskOption.GetAutoStartGammer() == 0 {
	//		//非房主时，自动准备
	//		if u.GetUserId() != u.Desk.GetOwner() {
	//			defer u.DoReadyFriend()
	//		}
	//	}else {
	//		//自动准备
	//		defer u.DoReadyFriend()
	//	}
	//}

	//如果是代开,且该用户第一个进入房间则设置该用户为房主
	if desk.GetDeskNumber() != "" {
		if desk.GetSitedUserNum() == 1 {
			desk.Owner = proto.Uint32(u.GetUserId())
			desk.CurrBanker = proto.Uint32(u.GetUserId())
		}
		if desk.GetIsDaikai() {
			//同步代开状态
			roomAgent.DoAddUser(desk.GetDaikaiUser(), int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN), desk.GetDeskId(), u.GetNickName())
		}
	}

	//发送入座广播
	u.SendEnterDeskBC()
}

//离座
func (u *User) DoSiteUp() {
	if u.Desk.GetIsStart() {
		u.SendEnterDeskBcACK(-5, "已开局，无法离座！")
		return
	}

	if u.GetIndex() < 0 {
		u.SendEnterDeskBcACK(-4, "你还未入座，离座失败！")
		return
	}

	*u.Index = -1
	*u.IsReady = false

	//发送离座广播
	u.SendEnterDeskBC()
}

//准备
func (user *User) DoReady() {
	if user.Desk.GetDeskNumber() != "" {
		//朋友桌准备
		user.DoReadyFriend()
	}else {
		//金币场准备
		user.DoReadyCoin()
	}
}

//朋友桌准备
func (u *User) DoReadyFriend() {
	log.T("用户%d开始准备。", u.GetUserId())
	if u.GetIndex() == -1 {
		u.SendReadyAck(-7, "请先入座再准备！")
		return
	}
	defer u.Desk.WipeSnapShot()
	//如果是第一局，则等待所有非房主玩家准备，然后房主点击"开始游戏"，开始抢庄或加倍
	if u.Desk.GetCircleNo() == 1 {
		//玩家准备
		if u.Desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_ENTER || u.Desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY || u.Desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START {
			if u.GetIsReady() == true {
				log.E("用户%d重复准备.", u.GetUserId())
				u.SendReadyAck(1, "请不要重复准备！")
				return
			}

			if !u.IsOwner() {
				log.T("用户%d准备成功.", u.GetUserId())
				//如果是非房主准备
				*u.IsReady = true
				u.SendReadyAck(1, "准备成功！")
				u.SendReadyBC()

				//如果达到autoStart条件
				if u.Desk.DeskOption.GetAutoStartGammer() != 0 {
					if err := u.IsAllReadyAutoStart(); err == nil {
						owner,_ := u.Desk.GetUserByUid(u.Desk.GetOwner())
						//更新房主状态
						*owner.IsReady = true
						u.DoStart()
						return
					}
				}

				//如果四个人全部已准备
				if err := u.IsAllReadyExcludeOwner(); err == nil {
					//给房主发送开始overturn
					u.Desk.SendStartOt()
					//将房间状态切换为等待房主开局
					u.Desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START.Enum()
					return
				}
				return
			} else {
				//如果是房主准备
				if u.GetIsReady() == true {
					log.E("房主%d重复开局.", u.GetUserId())
					u.SendReadyAck(1, "请不要重复开局！")
					return
				}
				if err := u.IsAllReadyExcludeOwner(); err == nil {
					log.T("房主%d开局成功.", u.GetUserId())
					*u.IsReady = true
					u.SendReadyAck(1, "发起开局成功！")
					//尝试开局
					u.Desk.DoStart()
					return
				} else {
					log.E("房主%d开局失败，原因：%s", u.GetUserId(), err.Error())
					u.SendReadyAck(-2, "开局失败！因为"+err.Error())
					return
				}
			}
		}else {
			//中途加入，但不在准备阶段
			u.SendReadyAck(-3, "当前不在准备阶段。")
		}
	} else {
		//如果是第二局及以后，则等所有玩家都点"继续游戏(准备协议)" 后，开始抢庄或加倍。
		if u.Desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY || u.Desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_ENTER || u.Desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START {
			if u.GetIsReady() == true {
				log.E("用户%d在第%d圈，重复准备!", u.GetUserId(), u.Desk.GetCircleNo())
				u.SendReadyAck(1, "请不要重复准备！")
				return
			} else {
				log.T("用户%d在第%d圈,准备成功！", u.GetUserId(), u.Desk.GetCircleNo())
				*u.IsReady = true
				u.SendReadyAck(1, "准备成功！")
				u.SendReadyBC()
				//尝试开局
				if err := u.Desk.IsAllReady(); err == nil {
					u.Desk.DoStart()
				}
				return
			}
		}else {
			//中途加入，但不在准备阶段
			u.SendReadyAck(-3, "当前不在准备阶段。")
		}
	}

}

//离开房间
func (user *User) DoLeaveDesk() error {
	log.T("用户%d开始离开房间%d。", user.GetUserId(), user.Desk.GetDeskNumber())
	msg := &ddproto.CommonAckLeaveDesk{
		Header:     commonNewPorot.NewHeader(),
		UserId:     user.UserId,
		IsExchange: proto.Bool(false),
	}
	if user.Desk.GetDeskNumber() != "" {
		//朋友桌
		if user.Desk.GetCircleNo() == 1 && !user.GetIsOnGamming() {
			//离开房间
			user.Desk.RemoveUser(user.GetUserId())
			return nil
		}else if user.GetIndex() <= -1 {
			//已开局，但是未入座的，允许离开房间
			user.Desk.RemoveUser(user.GetUserId())
			return nil
		}else {
			log.E("用户%d退出房间%d失败，原因：已开局不能离开房间", user.GetUserId(), user.Desk.GetDeskNumber())
			*msg.Header.Code = -1
			*msg.Header.Error = "已开局，不支持离开房间。"
			user.WriteMsg(msg)
			return errors.New("朋友桌不能离开房间。")
		}
	}

	defer user.Desk.WipeSnapShot()

	//金币场
	if user.GetIsOnGamming() {
		//游戏中不能退出，将玩家标记为已退出状态
		//user.IsLeave = proto.Bool(true)
		//log.T("将用户%d标记为离开状态。", user.GetUserId())
		*msg.Header.Code = -2
		*msg.Header.Error = "游戏中不能离开房间，请打完这局再试！"
		user.WriteMsg(msg)
		return nil
	}

	//金币场，未开局直接退出
	user.Desk.RemoveUser(user.GetUserId())

	return nil
}

//开始比赛
func (desk *Desk) DoStart() error {
	//朋友桌状态验证
	if !desk.GetIsCoinRoom() && desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START && desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY && desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_ENTER {
		log.E("DoStart()，该房间%d未在准备状态，请求非法。", desk.GetDeskNumber())
		return errors.New("非法的请求！")
	}

	//朋友桌
	if desk.GetDeskNumber() != ""{
		//房间标记为游戏中
		desk.IsOnGamming = proto.Bool(true)
		//游戏一开始，则关闭自动同意解散的计时器，防止房间被自动解散
		if desk.DissolveTimer != nil {
			desk.DissolveTimer.Stop()
			desk.DissolveTimer = nil
		}
		//关闭readyTimer
		if desk.ReadyTimer != nil {
			desk.ReadyTimer.Stop()
			desk.ReadyTimer = nil
		}
		//所有已准备用户标记为处于游戏中
		for _,u := range desk.Users{
			if !u.GetIsReady() {
				continue
			}
			u.IsOnGamming = proto.Bool(true)

			if desk.DeskOption.GetIsCoinRoom() {
				coin_fee := u.Desk.DeskOption.GetBaseScore()/2
				if coin_fee <= 1 {
					coin_fee = 1
				}
				//扣除房费
				userService.DECRUserCOIN(u.GetUserId(), coin_fee, "牛牛金币场自建房，房费扣除")
				//房费代理结算
				countService.DoCoinFeeAgentBill(u.GetUserId(), coin_fee, int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN))
			}
		}
		//开始统计
		desk.DoCountStart()
		//如果是第一圈
		if desk.GetCircleNo() == 1 {
			log.T("房间%d,BillInit()", desk.GetDeskNumber())
			//朋友桌房间设为已开局状态
			*desk.IsStart = true
			//如果是代开则同步
			if desk.GetIsDaikai() {
				roomAgent.DoStart(desk.GetDaikaiUser(), int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN), desk.GetDeskId())
			}
		}

		//开始游戏
		desk.DoPlay()
	}else {
		//金币场
		if desk.GetIsOnGamming() == true || !desk.CoinHasEoughtUserReady() {
			log.T("尝试开局失败，因为牌桌正在游戏中或者未达到最小开局人数。")
			return errors.New("尝试开局失败，因为牌桌正在游戏中或者未达到最小开局人数。")
		}
		//达到开局条件,首先锁定牌桌状态防止重复开局
		desk.IsOnGamming = proto.Bool(true)

		//等待用户准备
		if desk.ReadyTimer == nil {
			desk.ReadyTimer = time.AfterFunc(3 * time.Second, func() {  //自贡客户要求：等待3秒后，即开始游戏
				defer Error.ErrorRecovery("desk.ReadyTimer()")
				defer func() {
					desk.ReadyTimer = nil
				}()
				if !desk.CoinHasEoughtUserReady() {
					//取消开局锁定
					desk.IsOnGamming = proto.Bool(false)
					log.T("尝试开局失败，因为中途有人退出或者离线。")
					return
				}

				//初始化用户开局状态
				for _,u := range desk.Users{
					if u != nil && u.GetIsReady() && u.GetIsOnline() {
						u.IsOnGamming = proto.Bool(true)
						//扣除房费
						userService.DECRUserCOIN(u.GetUserId(), u.Room.GetCoinFee(), "牛牛金币场，房费扣除")
						//房费代理结算
						countService.DoCoinFeeAgentBill(u.GetUserId(), u.Room.GetCoinFee(), int32(ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN))
					}
				}
				//开始游戏
				desk.DoPlay()
			})
		}

	}

	return nil
}

//开始玩游戏
func (desk *Desk)DoPlay() {
	//发牌
	err := desk.DoSendPoker()
	if err != nil {
		log.E("牌桌%d发牌失败，原因:%s", desk.GetDeskNumber(), err.Error())
	}

	//确定庄家
	switch desk.GetDeskOption().GetBankRule() {
	case ddproto.NiuniuEnumBankerRule_SUI_JI_ZUO_ZHUANG:
		//更新牌桌状态
		desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_JIABEI.Enum()
		//随机坐庄
		length := len(desk.Users)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		index := r.Intn(length)
		*desk.CurrBanker = desk.Users[index].GetUserId()
		//发送加倍overturn
		desk.SendJiabeiOt()
	case ddproto.NiuniuEnumBankerRule_DING_ZHUANG, ddproto.NiuniuEnumBankerRule_FANGZHU_DINGZHUANG:  //虽然叫定庄，但实际为牛牛换庄
		//更新牌桌状态
		desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_JIABEI.Enum()
		//不处理
		//发送加倍overturn
		desk.SendJiabeiOt()
	case ddproto.NiuniuEnumBankerRule_QIANG_ZHUANG:
		//更新牌桌状态
		desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_QIANGZHUANG.Enum()
		//发起抢庄overturn
		desk.SendQiangzhuangOt()
	case ddproto.NiuniuEnumBankerRule_TONG_BI_NIUNIU:
		//加倍overturn，发牌
		desk.SendJiabeiOt()
		//直接将牌最大的玩家设为庄家
		user_rank := desk.GetUserRankByPoker()
		*desk.CurrBanker = user_rank[0].GetUserId()
		//自动加倍阶段
		desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_JIABEI.Enum()
		for _,u := range desk.Users {
			if !u.GetIsOnGamming() {
				continue
			}

			*u.DoubleScore = 1
			u.SendJiabeiBc()
		}
		//直接比牌
		desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_BIPAI.Enum()
		desk.DoBipai()
	}
}

//牌桌--10圈比赛打完初始化
func (desk *Desk) DoEnd() error {
	//每一圈结束
	//房间标记为游戏结束
	desk.IsOnGamming = proto.Bool(false)
	//所有用户标记为处于游戏结束
	for _,u := range desk.Users{
		*u.IsReady = false
		*u.IsLiangpai = false
	}

	//结束统计
	desk.DoCountEnd()
	//朋友桌
	if desk.GetDeskNumber() != "" {
		for _,u := range desk.Users{
			u.IsOnGamming = proto.Bool(false)
		}

		if !desk.GetIsCoinRoom() && desk.GetCircleNo() == desk.DeskOption.GetMaxCircle() {
			//关闭托管
			for _,u := range desk.Users {
				if u.GetIsTuoguan() {
					u.DoTuoguan(false, nil)
				}
			}
			//朋友桌最后一圈结束
			desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_RESULT.Enum()

			//销毁房间状态,发送结算广播，结算统计
			err := desk.RemoveFriendDesk(desk.GetDeskId())
			if err != nil {
				log.E("销毁房间%d失败。原因：%s", desk.GetDeskNumber(), err.Error())
				return err
			}
			return nil
		} else {
			*desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY
			//朋友桌非最后一圈结束 、 金币场每一圈结束
			*desk.NiuniuSrvDesk.CircleNo++
			//初始化牌桌状态
			new_game_number, err := db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)
			if err != nil {
				log.E("房间%d更新自增id失败，原因:%s", desk.GetDeskNumber(), err.Error())
				return errors.New("get gamenumber seq id fail.")
			}
			*desk.GameNumber = new_game_number

			//初始化用户状态
			for _, u := range desk.Users {
				if u != nil {
					//初始化手牌
					//u.Pokers = nil  不清理手牌，以便前端断线处理
					*u.BankerScore = 0
					*u.DoubleScore = 0
				}
			}

			//托管-自动准备
			desk.DoTuoguan()

			//7秒后自动准备
			if desk.ReadyTimer != nil {
				desk.ReadyTimer.Stop()
			}
			desk.StartTime = proto.Int64(time.Now().Unix())
			desk.ReadyTimer = time.AfterFunc(10*time.Second, func() {
				defer Error.ErrorRecovery("desk.ReadyTimer2()")
				//容错
				if desk.GetCircleNo() == desk.DeskOption.GetMaxCircle() {
					return
				}
				for _,u := range desk.Users {
					if u== nil || u.GetIndex() <= -1 || u.GetIsReady() {
						continue
					}
					//自动准备
					//u.DoReady()
					//自动进入托管
					u.DoTuoguan(true, &ddproto.NiuTuoguanOption{
						QiangZhuangOpt: ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_BU_QIANG.Enum(),
						YaZhuOpt: ddproto.NiuEnumTuoguanYzopt_NIU_TG_YZ_YA_1.Enum(),
					})
				}
			})
		}
	}else {
		//金币场
		*desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY
		//初始化牌桌状态
		new_game_number, err := db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)
		if err != nil {
			log.E("房间%d更新自增id失败，原因:%s", desk.GetDeskNumber(), err.Error())
			return errors.New("get gamenumber seq id fail.")
		}
		*desk.GameNumber = new_game_number

		//初始化用户状态
		for _, u := range desk.Users {
			if u == nil {
				continue
			}

			//初始化手牌
			//u.Pokers = nil  不清理手牌，以便前端断线处理
			*u.BankerScore = 0
			*u.DoubleScore = 0

			//更新观战时间,并请出超时的用户
			if u.GetIsOnGamming() {
				u.IsOnGamming = proto.Bool(false)
			}

			//如果该用户已离开，则从牌桌删除
			if u.GetIsLeave() {
				u.Desk.RemoveUser(u.GetUserId())
			}

			//机器人准备处理
			if u.GetIsRobot() {
				u.DoRobotReady(false)
			}
		}
	}

	//玩家金币不足则踢出房间
	if desk.GetIsCoinRoom() {
		for _,u := range desk.Users {
			if userService.GetUserCoin(u.GetUserId()) < desk.DeskOption.GetMinEnterScore() {
				u.DoLeaveDesk()
			}
		}
	}

	return nil
}

//设置抢庄
func (u *User) DoQiangzhuang(qiangzhuang_score int64) error {
	log.T("用户%d在房间%d,发起抢庄。", u.GetUserId(), u.Desk.GetDeskNumber())

	if u.Desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_QIANGZHUANG {
		log.E("用户%d在房间%d,抢庄失败，原因：该房间状态不是wait_qiangzhuang,当前状态为%v", u.GetUserId(), u.Desk.GetDeskNumber(), u.Desk.GetStatus())
		u.SendQiangzhuangAck(-1, "该桌面状态下不能抢庄!")
		return errors.New("该桌面状态下不能抢庄!")
	}

	defer u.Desk.WipeSnapShot()

	if u.GetBankerScore() == 0 {
		if qiangzhuang_score > 0 || qiangzhuang_score == -1 {
			log.T("用户%d抢庄成功。", u.GetUserId())
			*u.BankerScore = qiangzhuang_score

			//u.SendQiangzhuangAck(1, "抢庄成功！")
			//把抢庄成功广播出去
			u.SendQiangzhuangBC()
			//刷新旁观timer
			u.RefreshAsideTimer()

			//尝试发起加倍
			if err := u.Desk.IsAllQiangzhuang(); err == nil {
				//清除抢庄计时器
				if u.Desk.QiangzhuangTimer != nil {
					u.Desk.QiangzhuangTimer.Stop()
					u.Desk.QiangzhuangTimer = nil
				}
				//将抢庄分数最高的玩家设为庄家
				var max_score int64 = -1
				for _, user := range u.Desk.Users {
					if user != nil && user.GetIsOnGamming() && user.GetBankerScore() > max_score {
						max_score = user.GetBankerScore()
					}
				}
				max_users := []*User{}
				for _, user := range u.Desk.Users {
					if user != nil && user.GetIsOnGamming() && user.GetBankerScore() == max_score {
						max_users = append(max_users, user)
					}
				}

				var new_banker *User = nil
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				new_banker = max_users[r.Intn(len(max_users))]

				//新庄家下把必定不能推注
				*new_banker.TuizhuScore = 0

				//切换新庄家
				*u.Desk.CurrBanker = new_banker.GetUserId()

				//发送抢庄结果广播
				u.Desk.SendQiangzhuangResBc()

				u.Desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_JIABEI.Enum()
				//开始发送加倍overturn
				u.Desk.SendJiabeiOt()

			}else {
				log.T("IsAllQiangzhuang() err:%v", err)
			}
			return nil
		} else {
			log.E("用户%d抢庄失败,原因：抢庄倍数非法%d", u.GetUserId(), qiangzhuang_score)
			u.SendQiangzhuangAck(-3, "抢庄失败！因为该倍数非法")
			return nil
		}
	} else {
		log.E("用户%d抢庄失败,原因：重复抢庄", u.GetUserId())
		u.SendQiangzhuangAck(-2, "你已经抢过庄了!请不要重复请求")
		return errors.New("你已经抢过庄了!请不要重复请求")
	}

	return nil
}

//设置加倍
func (u *User) DoJiabei(double_score int64) error {
	log.T("用户%d在房间%d发起加倍请求。", u.GetUserId(), u.Desk.GetDeskNumber())
	// 加倍
	if u.Desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_JIABEI {
		log.E("用户%d加倍失败，原因：房间状态为%s", u.GetUserId(), u.Desk.GetStatus().String())
		u.SendJiabeiAck(-1, "该桌面状态下不能加倍!")
		return errors.New("该桌面状态下不能加倍!")
	}
	if !u.GetIsOnGamming() {
		log.E("用户%d加倍失败，原因：isOnGamming状态为%s", u.GetUserId(), u.GetIsOnGamming())
		u.SendJiabeiAck(-4, "您未在游戏中，无法加倍!")
		return errors.New("为在游戏中不能加倍!")
	}

	//庄家不能加倍
	if u.IsBanker() {
		u.SendJiabeiAck(-5, "您是庄家，无法加倍!")
		return errors.New("您是庄家，无法加倍!")
	}

	defer u.Desk.WipeSnapShot()

	if u.GetDoubleScore() == 0 {
		if double_score > 0 {
			log.T("用户%d加倍成功。", u.GetUserId())
			*u.DoubleScore = double_score

			//刷新旁观timer
			u.RefreshAsideTimer()

			//更新推注数据
			if u.GetTuizhuScore() == int32(double_score) {
				*u.TuizhuScore = 0
				*u.LastTuizhuCircleNo = u.Desk.GetCircleNo()
			}

			u.SendJiabeiAck(1, "加倍成功！")
			//发送加倍广播
			u.SendJiabeiBc()
			//尝试比牌
			if err := u.Desk.IsAllJiabeiExcludeBanker(); err == nil {
				//关闭加倍timer
				if u.Desk.JiaBeiTimer != nil {
					u.Desk.JiaBeiTimer.Stop()
					u.JiaBeiTimer = nil
				}
				u.Desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_LIANGPAI.Enum()
				//开始亮牌
				u.Desk.SendLiangpaiOt()
			}
			return nil
		} else {
			log.E("用户%d加倍失败，原因：倍数%d非法", u.GetUserId(), double_score)
			u.SendJiabeiAck(-3, "加倍失败！因为该倍数非法")
			return nil
		}
	} else {
		log.E("用户%d加倍失败，原因：重复加倍", u.GetUserId())
		u.SendJiabeiAck(-2, "你已经加过倍了!请不要重复请求")
		return errors.New("你已经加过倍了!请不要重复请求")
	}

	return nil
}

//亮牌
func (user *User) DoLiangpai() error {
	log.T("用户%d在房间%d发起亮牌请求。", user.GetUserId(), user.Desk.GetDeskNumber())
	// 亮牌
	if user.Desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_LIANGPAI {
		log.E("用户%d亮牌失败，原因：房间状态为%s", user.GetUserId(), user.Desk.GetStatus().String())
		user.SendLiangpaiAck(-1, "该桌面状态下不能加倍!")
		return errors.New("该桌面状态下不能加倍!")
	}

	if user.GetIsLiangpai() {
		user.SendLiangpaiAck(-2, "您已亮牌，请不要重复发起亮牌请求！")
		return errors.New("您已亮牌，请不要重复发起亮牌请求！")
	}

	if !user.GetIsOnGamming() {
		user.SendLiangpaiAck(-3, "您未在游戏中，无法进行亮牌操作！")
		return errors.New("您未在游戏中，无法进行亮牌操作！")
	}

	*user.IsLiangpai = true

	//亮牌成功
	user.SendLiangpaiAck(1, "亮牌成功!")
	//广播出去
	user.SendLiangpaiBc()

	//所有人亮完牌就结算
	if user.Desk.IsAllLiangpai() == nil {
		//关掉亮牌timer
		if user.Desk.LiangpaiTimer != nil {
			user.Desk.LiangpaiTimer.Stop()
			user.Desk.LiangpaiTimer = nil
		}
		user.Desk.Status = ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_BIPAI.Enum()
		//开始比牌
		user.Desk.DoBipai()
	}

	return nil
}

//托管
func (user *User) DoTuoguan(isTuoguan bool, opt *ddproto.NiuTuoguanOption) error {
	if user.Desk.GetRoomId() > 0 {
		user.SendTuoguanAck(-3, "匹配场不支持游戏托管！")
		return errors.New("匹配场不支持游戏托管！")
	}

	if user.GetIndex() == -1 {
		user.SendTuoguanAck(-1, "请先入座，才能选择托管！")
		return errors.New("请先入座，才能选择托管！")
	}

	if user.GetIsTuoguan() == isTuoguan {
		user.SendTuoguanAck(-2, "请勿重复操作！")
		return errors.New("请勿重复操作！")
	}

	if isTuoguan {
		//打开托管
		user.IsTuoguan = proto.Bool(true)
		user.TuoGuanOpt = opt
	}else {
		//关闭托管
		user.IsTuoguan = proto.Bool(false)
		user.TuoGuanOpt = nil
	}

	//托管bc
	user.SendTuoguanBc()

	//开始托管
	go user.DoTuoguanAct()

	return nil
}

//开始比牌
func (desk *Desk) DoBipai() error {
	log.T("牌桌%d开始比牌", desk.GetDeskNumber())
	if desk.GetStatus() != ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_BIPAI {
		log.E("牌桌%d比牌失败，原因：该房间状态为%s", desk.GetDeskNumber(), desk.GetStatus().String())
		return errors.New("未达到比牌条件！")
	}
	//开始比牌
	banker, err := desk.GetUserByUid(desk.GetCurrBanker())
	if err != nil {
		log.E("牌桌%d获取房主%d失败！", desk.GetDeskNumber(), desk.GetCurrBanker())
		return nil
	}

	var banker_score int64 = 0
	var max_user_id uint32 = 0
	var max_user_poker *ddproto.NiuniuSrvPoker = nil

	poker_score_list := []*ddproto.NiuBipaiResultItem{}

	//初始化计分参数
	for _, u := range desk.Users {
		if u.GetBankerScore() == 0 {
			*u.BankerScore = 1
		}
		if u.GetBankerScore() == -1 {
			*u.BankerScore = 1
		}
		if u.GetDoubleScore() == 0 {
			*u.DoubleScore = 1
		}
		if u.GetDoubleScore() == -1 {
			*u.DoubleScore = 1
		}
	}

	for _, u := range desk.Users {
		if u != nil && !u.IsBanker() && u.GetIsOnGamming() {
			var user_score int64 = 0
			//比庄家大
			if IsBigThanBanker(banker.GetPokers(), u.GetPokers()) {
				user_score = banker.GetBankerScore() * u.GetDoubleScore() * GetPokerScore(u.GetPokers(), desk.DeskOption)
				//如果游戏模式为顶庄，且该玩家比庄家牌大，且牌型大于等于牛牛;则下局该玩家成为庄家
				if desk.GetDeskOption().GetBankRule() == ddproto.NiuniuEnumBankerRule_DING_ZHUANG && u.GetPokers().GetType() >= ddproto.NiuniuEnum_PokerType_NIU_NIU {
					if max_user_poker == nil {
						max_user_poker = u.GetPokers()
						max_user_id = u.GetUserId()
					}
					if IsBigThanBanker(max_user_poker, u.GetPokers()) {
						max_user_poker = u.GetPokers()
						max_user_id = u.GetUserId()
					}
				}

				//闲家比庄家大，则下把可推注，不允许连续推注
				if desk.DeskOption.GetTuizhuScore() > 0 && desk.GetCircleNo() - u.GetLastTuizhuCircleNo() > 1 {
					//老版推注分数算法
					//tuizhu_score := desk.DeskOption.GetDiFen() + int32(u.GetDoubleScore())
					//if tuizhu_score > desk.DeskOption.GetTuizhuScore() {
					//	tuizhu_score = desk.DeskOption.GetTuizhuScore()
					//}
					//新版牛元帅推注算法：
					tuizhu_score := desk.DeskOption.GetDiFen() * desk.DeskOption.GetTuizhuScore()
					*u.TuizhuScore = tuizhu_score
				}else {
					//无法推注，则将推注值清零
					*u.TuizhuScore = 0
				}

			} else {
				user_score = banker.GetBankerScore() * u.GetDoubleScore() * GetPokerScore(banker.GetPokers(), desk.DeskOption) * -1
				//牌比庄家小，必定无法推注，则将推注值清零
				*u.TuizhuScore = 0
			}
			//金币场结算
			user_score = user_score * u.DeskOption.GetBaseScore()
			//更新非庄家玩家总积分
			if !desk.GetIsCoinRoom() {
				//朋友桌
				*u.Bill.Score += user_score
			}else {
				//金币场
				banker_coin := userService.GetUserCoin(banker.GetUserId())
				user_coin := userService.GetUserCoin(u.GetUserId())
				var surplus_coin int64 = 0
				if user_score > 0 {
					//玩家赢，庄家输
					if banker_coin < user_score {
						user_score = banker_coin
					}
					surplus_coin,_ = userService.INCRUserCOIN(u.GetUserId(), int64(user_score), "牛牛金币场，单局结算")
				}else if user_score < 0 {
					//玩家输，庄家赢
					if user_coin < int64(-user_score) {
						user_score = -user_coin
					}
					surplus_coin,_ = userService.DECRUserCOIN(u.GetUserId(), int64(-user_score), "牛牛金币场，单局结算")
				}
				*u.Bill.Score = surplus_coin
			}
			banker_score = banker_score - user_score
			//更新非庄家上局得分
			*u.LastScore = user_score
			//非庄家输赢
			poker_score_list = append(poker_score_list, &ddproto.NiuBipaiResultItem{
				Poker:    GetClientPoker(u.GetPokers()),
				Score:    &user_score,
				UserId:   u.UserId,
				AllScore: u.Bill.Score,
			})
			log.T("[比牌结算]:房号:%v 牌桌id:%v 圈数:%v 用户%d 为闲家 QzScore:%d JbScore:%d Score:%d AllScore:%d Poker:%v",u.Desk.GetDeskNumber(), u.GetDeskId(), u.Desk.GetCircleNo(), u.GetUserId(), u.GetBankerScore(), u.GetDoubleScore(), user_score, u.Bill.GetScore(), u.Pokers.GetType())
			//朋友桌总输赢统计
			if !desk.GetIsCoinRoom() {
				//统计非庄家输赢
				if user_score > 0 {
					*u.Bill.CountWin++
				} else {
					*u.Bill.CountLost++
				}
				//统计非庄家是否有牛
				if u.GetPokers().GetType() > ddproto.NiuniuEnum_PokerType_NO_NIU {
					*u.Bill.CountHasNiu++
				} else {
					*u.Bill.CountNoNiu++
				}
			}
		}
	}
	//更新庄家总积分
	if !desk.GetIsCoinRoom() {
		//朋友桌
		*banker.Bill.Score += banker_score
	}else {
		//金币场 庄家剩余金币数
		var surplus_coin int64 = 0
		if banker_score > 0 {
			surplus_coin,_ = userService.INCRUserCOIN(banker.GetUserId(), int64(banker_score), "牛牛金币场，单局结算")
		}else if banker_score < 0 {
			surplus_coin,_ = userService.DECRUserCOIN(banker.GetUserId(), int64(-banker_score), "牛牛金币场，单局结算")
		}
		*banker.Bill.Score = surplus_coin
	}
	//更新庄家上局得分
	*banker.LastScore = banker_score

	//朋友桌总输赢统计
	if !desk.GetIsCoinRoom() {
		//统计庄家输赢
		if banker_score > 0 {
			*banker.Bill.CountWin++
		} else {
			*banker.Bill.CountLost++
		}
		//统计庄家是否有牛
		if banker.GetPokers().GetType() > ddproto.NiuniuEnum_PokerType_NO_NIU {
			*banker.Bill.CountHasNiu++
		} else {
			*banker.Bill.CountNoNiu++
		}
	}

	//庄家输赢
	poker_score_list = append(poker_score_list, &ddproto.NiuBipaiResultItem{
		Poker:    GetClientPoker(banker.GetPokers()),
		Score:    &banker_score,
		UserId:   banker.UserId,
		AllScore: banker.Bill.Score,
	})
	log.T("[比牌结算]:房号:%v 牌桌id:%v 圈数:%v 用户%d 为庄家 QzScore:%d JbScore:%d Score:%d AllScore:%d oker:%v",banker.Desk.GetDeskNumber(), banker.GetDeskId(), banker.Desk.GetCircleNo(), banker.GetUserId(), banker.GetBankerScore(), banker.GetDoubleScore(), banker_score, banker.Bill.GetScore(), banker.Pokers.GetType())
	//插入非游戏中玩家空数据，使其兼容客户端处理
	for _,u := range desk.Users {
		if u.GetIsOnGamming() {
			continue
		}
		poker_score_list = append(poker_score_list, &ddproto.NiuBipaiResultItem{
			Poker:    GetClientPoker(u.GetPokers()),
			Score:    proto.Int64(0),
			UserId:   u.UserId,
			AllScore: u.Bill.Score,
		})
	}

	//如果是定庄模式，则判断是否有新庄家产生
	if desk.GetDeskOption().GetBankRule() == ddproto.NiuniuEnumBankerRule_DING_ZHUANG && max_user_id != 0 {
		*desk.CurrBanker = max_user_id
	}

	//发送比牌结果广播
	bipai_bc := &ddproto.NiuBipaiResultBc{
		UserState: poker_score_list,
	}
	desk.BroadCast(bipai_bc)

	//牌局结束
	desk.DoEnd()
	return nil
}

//发起解散房间
func (user *User) DoApplyDissolve() error {
	log.T("用户%d发起解散房间", user.GetUserId())

	if user.Desk.GetIsOnDissolve() == true {
		log.T("用户%d发起解散房间失败，原因：", user.Desk.GetDeskNumber(), "正处于解散投票阶段，不要重复申请！")
		return errors.New("正处于解散投票阶段，不要重复申请！")
	}

	//第一局必须是房主才能解散
	if user.Desk.GetCircleNo() == 1 {
		owner := user.Desk.GetOwner()
		if user.Desk.GetIsDaikai() {
			owner = user.Desk.GetDaikaiUser()
		}

		if user.GetUserId() == owner {
			return user.Desk.Room.RemoveFriendDesk(user.Desk.GetDeskId())
		}
	}

	//是否入座
	if user.Desk.GetCircleNo() > 1 && user.GetIndex() == -1 {
		log.T("用户%d发起解散房间失败，原因：%s", user.GetUserId(), "未入座！")
		return errors.New("未入座，不能解散房间！")
	}

	//更改状态
	*user.Desk.IsOnDissolve = true
	user.Desk.DissolveUser = proto.Uint32(user.GetUserId())
	*user.Desk.DissolveTime = time.Now().Unix()

	//5分钟后强制解散房间
	if user.Desk.DissolveTimer != nil {
		user.Desk.DissolveTimer.Stop()
		user.Desk.DissolveTimer = nil
	}
	user.Desk.DissolveTimer = time.AfterFunc(consts.APPLYDISSOLVE_DURATION, func() {
		defer Error.ErrorRecovery("desk.DessolveTimer()")
		if user.Desk.GetIsOnDissolve() == true {
			for _, u := range user.Desk.Users {
				if u != nil && u.GetDissolveState() == 0 {
					u.DoDissolveBack(true)
				}
			}
		}
	})

	//发起解散房间广播
	user.SendApplyDissolveBc()

	//自动投票
	user.DoDissolveBack(true)

	return nil
}

//同意、拒绝解散房间
func (user *User) DoDissolveBack(isAgree bool) error {
	log.T("用户%d同意或拒绝解散房间,是否同意：%v", user.GetUserId(), isAgree)
	if user.Desk.GetIsOnDissolve() == false {
		log.E("用户%d解散房间投票失败，原因：当前房间未处于解散投票阶段，投票失败！", user.GetUserId())
		return errors.New("当前房间未处于解散投票阶段，投票失败！")
	}

	if user.GetDissolveState() != 0 {
		log.E("用户%d解散房间投票失败，原因：请不要重复投票！", user.GetUserId())
		return errors.New("请不要重复投票！")
	}

	//没入座，不能同意或拒绝
	if user.GetIndex() == -1 {
		log.E("用户%d解散房间投票失败，原因：没入座没有资格投票！", user.GetUserId())
		return errors.New("没入座不能投票！")
	}

	//更新用户解散投票状态
	if isAgree == true {
		*user.DissolveState = 1
	} else {
		*user.DissolveState = -1
	}

	//发送投票广播
	user.SendDissolveBackBc(isAgree)

	if isAgree == false {
		//如果有人拒绝解散，则初始化房间和用户状态
		user.Desk.SendDissolveDoneBc(isAgree)
		if user.Desk.DissolveTimer != nil {
			user.Desk.DissolveTimer.Stop()
			user.Desk.DissolveTimer = nil
		}
		//初始化桌面解散状态
		*user.Desk.IsOnDissolve = false
		*user.Desk.DissolveTime = 0
		for _, u := range user.Desk.Users {
			*u.DissolveState = 0
		}
		return nil
	}

	//确认所有人都已投票
	for _, u := range user.Desk.Users {
		if u != nil && u.GetIndex() != -1 {
			if u.GetDissolveState() == 0 {
				//如果离线，则自动同意
				if u.GetIsOnline() == false {
					*u.DissolveState = 1
				} else {
					log.T("房间%d解散房间失败，原因：用户%d未投票", user.Desk.GetDeskNumber(), user.GetUserId())
					//有人未投票
					return nil
				}
			}
		}
	}

	// 当所有人都确认解散,发送解散成功或失败广播:
	// 如果解散失败则发协议，解散成功则直接发送牌局结束协议。
	if isAgree == false {
		user.Desk.SendDissolveDoneBc(isAgree)
	}

	if user.Desk.DissolveTimer != nil {
		user.Desk.DissolveTimer.Stop()
		user.Desk.DissolveTimer = nil
	}
	//初始化桌面解散状态
	*user.Desk.IsOnDissolve = false
	*user.Desk.DissolveTime = 0
	for _, u := range user.Desk.Users {
		if u != nil {
			*u.DissolveState = 0
		}
	}

	if isAgree == true {
		//删除牌桌状态
		user.Desk.Room.RemoveFriendDesk(user.Desk.GetDeskId())
		log.T("房间%d解散房间成功", user.Desk.GetDeskNumber())
	}

	return nil
}

//统计开局
func (desk *Desk) DoCountStart() {
	//更新统计时间
	*desk.OneStartTime = time.Now().Unix()

	if desk.GetCircleNo() == 1 {
		*desk.AllStartTime = time.Now().Unix()
	}
}

//统计结束
func (desk *Desk) DoCountEnd() {
	//朋友桌，插入10局记录
	if !desk.GetIsCoinRoom() && desk.GetCircleNo() == desk.DeskOption.GetMaxCircle() {
		go func(){
			defer Error.ErrorRecovery("DoCountEnd()->all")
			desk.InsertAllCounter()
		}()
	}

	//插入1局记录
	go func() {
		defer Error.ErrorRecovery("DoCountEnd()->one")
		desk.InsertOneCounter()
	}()
}
