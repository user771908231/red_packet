package laowangye

import (
	"casino_common/common/consts"
	"casino_common/common/log"
	"casino_common/proto/ddproto"
	"casino_common/proto/funcsInit"
	"casino_common/utils/db"
	"casino_laowangye/conf/config"
	"errors"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"
	"casino_common/common/service/roomAgent"
	"casino_common/common/Error"
)
//入座
func (u *User) DoSiteDown() {
	//if site_index < 0 || site_index > u.Desk.DeskOption.GetMaxUser() - 1 {
	//	u.SendEnterDeskBcACK(-2, "座位号非法！")
	//	return
	//}

	//if !u.Desk.IsFreeSite(site_index) {
	//	u.SendEnterDeskBcACK(-3, "这个位置已经有人坐了！")
	//	return
	//}

	if u.Desk.GetIsStart() {
		u.SendSiteDownACK(-5, "已开局，无法入座！")
		return
	}

	if u.GetIndex() >= 0 {
		u.SendSiteDownACK(-7, "您已经在座位上了，请不要重复入座！")
		return
	}

	desk := u.Desk

	free_site_index, err := desk.GetFreeSiteIndex()
	if err != nil {
		u.SendSiteDownACK(-4, "已经没有空闲座位了！")
		return
	}

	//设置座位号
	*u.Index = int32(free_site_index)

	//自动准备
	//if u.Desk.DeskOption.GetAutoStartGammer() == 0 || u.Desk.GetOwner() != u.GetUserId() {
	//	u.DoReadyFriend()
	//}

	//如果是代开,且该用户第一个进入房间则设置该用户为房主
	if !desk.GetIsCoinRoom() && desk.GetIsDaikai() {
		if len(desk.Users) == 0 {
			desk.Owner = proto.Uint32(u.GetUserId())
			desk.CurrBanker = proto.Uint32(u.GetUserId())
		}
		//同步代开状态
		go roomAgent.DoAddUser(desk.GetDaikaiUser(), int32(ddproto.CommonEnumGame_GID_LAOWANGYE), desk.GetDeskId(), u.GetNickName())
	}

	//发送入座广播
	u.SendSiteDownBC()
}

//离座
func (u *User) DoSiteUp() {
	if u.Desk.GetIsStart() {
		u.SendSiteDownACK(-5, "已开局，无法离座！")
		return
	}

	if u.GetIndex() < 0 {
		u.SendSiteDownACK(-4, "你还未入座，离座失败！")
		return
	}

	*u.Index = -1
	*u.IsReady = false

	//发送离座广播
	u.SendSiteDownBC()
}

//朋友桌准备
func (u *User) DoReadyFriend() {
	log.T("用户%d开始准备。", u.GetUserId())
	defer u.Desk.WipeSnapShot()

	if u.GetIndex() <= 0 {
		u.SendReadyAck(-4, "请先入座！")
		return
	}

	//如果是第一局，则等待所有非房主玩家准备，然后房主点击"开始游戏"，开始抢庄或押注
	if u.Desk.GetCircleNo() == 1 {
		//玩家准备
		if u.Desk.GetStatus() == ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY {
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

				//如果四个人全部已准备
				if err := u.IsAllReadyExcludeOwner(); err == nil {
					//给房主发送开始overturn
					u.Desk.SendStartOt()
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
		//如果是第二局及以后，则等所有玩家都点"继续游戏(准备协议)" 后，开始抢庄或押注。
		if u.Desk.GetStatus() == ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY {
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
	log.T("用户%d开始离开房间%d。", user.GetUserId(), user.Desk.GetPassword())
	msg := &ddproto.CommonAckLeaveDesk{
		Header:     commonNewPorot.NewHeader(),
		UserId:     user.UserId,
		IsExchange: proto.Bool(false),
	}
	if !user.Desk.GetIsCoinRoom() {
		//朋友桌
		//if user.Desk.GetStatus() != ddproto.LwyEnumDeskState_NIU_DESK_STATUS_WAIT_START && user.Desk.GetStatus() != ddproto.LwyEnumDeskState_NIU_DESK_STATUS_WAIT_ENTER {
		//	log.E("用户%d退出房间%d失败，原因：房间已经开始比赛", user.GetUserId(), user.Desk.GetPassword())
		//	*msg.Header.Code = -1
		//	*msg.Header.Error = "比赛已经开始，退出房间失败！"
		//	user.WriteMsg(msg)
		//	return errors.New("比赛已经开始，退出房间失败！")
		//}
		//if user.IsOwner() {
		//	log.E("用户%d退出房间%d失败，原因：该用户为房主。", user.GetUserId(), user.Desk.GetPassword())
		//	*msg.Header.Code = -2
		//	*msg.Header.Error = "您是房主，退出房间失败！"
		//	user.WriteMsg(msg)
		//	return errors.New("您是房主，退出房间失败！")
		//}

		log.E("用户%d退出房间%d失败，原因：朋友桌不能离开。", user.GetUserId(), user.Desk.GetPassword())
		*msg.Header.Code = -1
		*msg.Header.Error = "朋友桌不能离开房间。"
		user.WriteMsg(msg)
		return errors.New("朋友桌不能离开房间。")
	}

	defer user.Desk.WipeSnapShot()

	//金币场
	if user.GetIsOnGamming() {
		//游戏中不能退出，将玩家标记为已退出状态
		user.IsLeave = proto.Bool(true)
		log.T("将用户%d标记为离开状态。", user.GetUserId())
		*msg.Header.Code = 2
		*msg.Header.Error = "退出房间成功，系统将帮你打完这一局。"
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
	if !desk.GetIsCoinRoom() && desk.GetStatus() != ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY {
		log.E("DoStart()，该房间%d未在准备状态，请求非法。", desk.GetPassword())
		return errors.New("非法的请求！")
	}

	//朋友桌
	if !desk.GetIsCoinRoom(){
		//房间标记为游戏中
		desk.IsOnGamming = proto.Bool(true)
		//游戏一开始，则关闭自动同意解散的计时器，防止房间被自动解散
		if desk.DissolveTimer != nil {
			desk.DissolveTimer.Stop()
			desk.DissolveTimer = nil
		}
		//所有用户标记为处于游戏中
		for _,u := range desk.Users{
			u.IsOnGamming = proto.Bool(true)
		}
		//开始统计
		desk.DoCountStart()
		//如果是第一圈
		if desk.GetCircleNo() == 1 {
			log.T("房间%d,BillInit()", desk.GetPassword())
			//朋友桌房间设为已开局状态
			*desk.IsStart = true
			//如果是代开则同步
			if desk.GetIsDaikai() {
				roomAgent.DoStart(desk.GetDaikaiUser(), int32(ddproto.CommonEnumGame_GID_LAOWANGYE), desk.GetDeskId())
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
			desk.ReadyTimer = time.AfterFunc(10 * time.Second, func() {
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
						//userService.DECRUserCOIN(u.GetUserId(), u.Room.GetCoinFee(), "老王爷金币场，房费扣除")
						//房费代理结算
						//countService.DoCoinFeeAgentBill(u.GetUserId(), u.Room.GetCoinFee(), int32(ddproto.CommonEnumGame_GID_LAOWANGYE))
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
	//关闭准备倒计时
	if desk.ReadyTimer != nil {
		desk.ReadyTimer.Stop()
		desk.ReadyTimer = nil
	}
	//确定庄家
	switch desk.GetDeskOption().GetBankRule() {
	case ddproto.LwyEnumBankerRule_LWY_LUN_LIU_ZUO_ZHUANG:
		//轮流坐庄
		if desk.GetCircleNo() > 1 {
			banker,err := desk.GetUserByUid(desk.GetCurrBanker())
			if err != nil {
				log.E("banker not found ! err:%v", err)
				return
			}
			next_site_index := (banker.GetIndex()%desk.GetSitedGammerNum())+1
			next_banker := desk.GetUserByIndex(next_site_index)
			if next_banker == nil {
				log.E("banker not found, site_index = %v ! err:%v", next_site_index, err)
				return
			}
			//更新新一轮的庄家
			*desk.CurrBanker = next_banker.GetUserId()
		}
		//更新牌桌状态
		desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_YAZHU.Enum()
		//发送押注overturn
		desk.SendYazhuOt()
	case ddproto.LwyEnumBankerRule_LWY_SUI_JI_ZUO_ZHUANG:
		//随机坐庄
		length := desk.GetSitedGammerNum()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		index := int32(r.Intn(int(length)))
		*desk.CurrBanker = desk.GetUserByIndex(index+1).GetUserId()
		//更新牌桌状态
		desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_YAZHU.Enum()
		//发送押注overturn
		desk.SendYazhuOt()
	case ddproto.LwyEnumBankerRule_LWY_FANGZHU_DINGZHUANG:
		//更新牌桌状态
		desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_YAZHU.Enum()
		//发送押注overturn
		desk.SendYazhuOt()
	case ddproto.LwyEnumBankerRule_LWY_QIANG_ZHUANG:
		//更新牌桌状态
		desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_QIANGZHUANG.Enum()
		//发起抢庄overturn
		desk.SendQiangzhuangOt()
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
	}

	//结束统计
	desk.DoCountEnd()
	//朋友桌
	if !desk.GetIsCoinRoom() {
		for _,u := range desk.Users{
			u.IsOnGamming = proto.Bool(false)
		}

		if !desk.GetIsCoinRoom() && desk.GetCircleNo() == desk.DeskOption.GetBoardsCout() {
			//朋友桌最后一圈结束
			desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_RESULT.Enum()

			//销毁房间状态,发送结算广播，结算统计
			err := desk.RemoveFriendDesk(desk.GetDeskId())
			if err != nil {
				log.E("销毁房间%d失败。原因：%s", desk.GetPassword(), err.Error())
			}
		} else {
			*desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY
			//朋友桌非最后一圈结束 、 金币场每一圈结束
			*desk.LwySrvDesk.CircleNo++
			//初始化牌桌状态
			new_game_number, err := db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)
			if err != nil {
				log.E("房间%d更新自增id失败，原因:%s", desk.GetPassword(), err.Error())
				return errors.New("get gamenumber seq id fail.")
			}
			*desk.GameNumber = new_game_number

			//初始化用户状态
			for _, u := range desk.Users {
				if u != nil {
					*u.BankerScore = 0
					u.YazhuDetail = &ddproto.LwyYazhuDetail{
						UserId: proto.Uint32(u.GetUserId()),
						Mai7: proto.Int64(0),
						Mai8: proto.Int64(0),
						Chi7: proto.Int64(0),
						Chi8: proto.Int64(0),
					}
					*u.Mai7Lost = 0
					*u.Mai8Lost = 0
					u.ChizhuDetail = []*ddproto.LwyChizhuDetailItem{}
				}
			}
		}
	}else {
		//金币场
		*desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY
		//初始化牌桌状态
		new_game_number, err := db.GetNextSeq(config.DBT_T_TH_GAMENUMBER_SEQ)
		if err != nil {
			log.E("房间%d更新自增id失败，原因:%s", desk.GetPassword(), err.Error())
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
			//*u.DoubleScore = 0

			//更新观战时间,并请出超时的用户
			if u.GetIsOnGamming() {
				u.IsOnGamming = proto.Bool(false)
				u.AsideWatchTime = proto.Int64(time.Now().Unix())
			}else {
				if time.Now().Unix() - u.GetAsideWatchTime() > int64(3 * 60) {
					//观战超时
					u.IsLeave = proto.Bool(true)
				}
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

	return nil
}

//设置抢庄
func (u *User) DoQiangzhuang(qiangzhuang_score int64) error {
	log.T("用户%d在房间%d,发起抢庄。", u.GetUserId(), u.Desk.GetPassword())

	if u.Desk.GetStatus() != ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_QIANGZHUANG {
		log.E("用户%d在房间%d,抢庄失败，原因：该房间状态不是wait_qiangzhuang,当前状态为%v", u.GetUserId(), u.Desk.GetPassword(), u.Desk.GetStatus())
		u.SendQiangzhuangAck(-1, "该桌面状态下不能抢庄!")
		return errors.New("该桌面状态下不能抢庄!")
	}

	defer u.Desk.WipeSnapShot()

	if u.GetBankerScore() == 0 {
		if qiangzhuang_score > 0 || qiangzhuang_score == -1 {
			log.T("用户%d抢庄成功。", u.GetUserId())
			*u.BankerScore = qiangzhuang_score
			u.SendQiangzhuangAck(1, "抢庄成功！")
			//尝试发起押注
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

				//切换新庄家
				*u.Desk.CurrBanker = new_banker.GetUserId()

				//更改牌桌状态为押注
				u.Desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_YAZHU.Enum()
				//开始发送押注overturn
				u.Desk.SendYazhuOt()
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

//押注
func (user *User) DoYazhu(yazhuType ddproto.LwyYazhuType, yazhuScore int64) error {
	log.T("用户%d在房间%d发起押注请求。", user.GetUserId(), user.Desk.GetPassword())
	// 押注
	if user.Desk.GetStatus() != ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_YAZHU {
		log.E("用户%d押注失败，原因：房间状态为%s", user.GetUserId(), user.Desk.GetStatus().String())
		user.SendYazhuAck(-1, "该桌面状态下不能押注!")
		return errors.New("该桌面状态下不能押注!")
	}
	if !user.GetIsOnGamming() {
		log.E("用户%d押注失败，原因：isOnGamming状态为%s", user.GetUserId(), user.GetIsOnGamming())
		user.SendYazhuAck(-4, "您未在游戏中，无法押注!")
		return errors.New("为在游戏中不能押注!")
	}
	defer user.Desk.WipeSnapShot()

	if yazhuType == ddproto.LwyYazhuType_LWY_MAI_7 || yazhuType == ddproto.LwyYazhuType_LWY_MAI_8 {
		if can_mai_score := 200 - (user.YazhuDetail.GetMai7() + user.YazhuDetail.GetMai8()); yazhuScore > can_mai_score {
			yazhuScore = can_mai_score
		}
		if yazhuScore <= 0 {
			user.SendYazhuAck(-5, "最多卖200！")
			return nil
		}
	}

	//总吃分
	chi_score := yazhuScore
	//已吃分
	var ex_chi_score int64 = 0
	if yazhuType == ddproto.LwyYazhuType_LWY_CHI_7 || yazhuType == ddproto.LwyYazhuType_LWY_CHI_8 {
		if user.YazhuDetail.GetChi7() + user.YazhuDetail.GetChi8() + chi_score > 200 {
			chi_score = 200 - user.YazhuDetail.GetChi7() - user.YazhuDetail.GetChi8()
		}
		if chi_score <= 0 {
			user.SendYazhuAck(-6, "最多吃200！")
			return nil
		}
	}


	switch yazhuType {
	case ddproto.LwyYazhuType_LWY_MAI_7:
		//卖7
		*user.YazhuDetail.Mai7 += yazhuScore
		*user.Mai7Lost += yazhuScore
	case ddproto.LwyYazhuType_LWY_MAI_8:
		//卖8
		*user.YazhuDetail.Mai8 += yazhuScore
		*user.Mai8Lost += yazhuScore
	case ddproto.LwyYazhuType_LWY_CHI_7:
		//开始吃7
		for _,u := range user.GetCanChiYours() {
			can_chi_score := chi_score
			if u.GetMai7Lost() < chi_score {
				can_chi_score = u.GetMai7Lost()
			}

			if can_chi_score > 0 {
				//更新chizhuDetail
				ex_user := false
				for _,item := range user.ChizhuDetail{
					if item.GetFrom() == u.GetUserId() {
						ex_user = true
						*item.Chi7Score += can_chi_score
					}
				}
				if ex_user == false {
					user.ChizhuDetail = append(user.ChizhuDetail, &ddproto.LwyChizhuDetailItem{
						From: proto.Uint32(u.GetUserId()),
						Chi7Score: proto.Int64(can_chi_score),
						Chi8Score: proto.Int64(0),
					})
				}
				//更新吃7
				*user.YazhuDetail.Chi7 += can_chi_score
				*u.Mai7Lost -= can_chi_score

				//更新chi_score
				chi_score -= can_chi_score
				ex_chi_score += can_chi_score
			}

			if chi_score == 0 {
				break
			}
		}
	case ddproto.LwyYazhuType_LWY_CHI_8:
		//开始吃8
		for _,u := range user.GetCanChiYours() {
			can_chi_score := chi_score
			if u.GetMai8Lost() < chi_score {
				can_chi_score = u.GetMai8Lost()
			}

			if can_chi_score > 0 {
				//更新chizhuDetail
				ex_user := false
				for _,item := range user.ChizhuDetail{
					if item.GetFrom() == u.GetUserId() {
						ex_user = true
						*item.Chi8Score += can_chi_score
					}
				}
				if ex_user == false {
					user.ChizhuDetail = append(user.ChizhuDetail, &ddproto.LwyChizhuDetailItem{
						From: proto.Uint32(u.GetUserId()),
						Chi7Score: proto.Int64(0),
						Chi8Score: proto.Int64(can_chi_score),
					})
				}
				//更新吃8
				*user.YazhuDetail.Chi8 += can_chi_score
				*u.Mai8Lost -= can_chi_score

				//更新chi_score
				chi_score -= can_chi_score
				ex_chi_score += can_chi_score
			}

			if chi_score == 0 {
				break
			}
		}
	}

	if yazhuType == ddproto.LwyYazhuType_LWY_CHI_7 || yazhuType == ddproto.LwyYazhuType_LWY_CHI_8 {
		if ex_chi_score == 0 {
			user.SendYazhuAck(-7, "没有分可吃了！")
			return nil
		}
	}

	//押注广播出去
	user.SendYazhuBC(yazhuType, yazhuScore)
	return nil
}

//摇色子
func (user *User) DoYaoshaizi() {
	if user.Desk.GetStatus() != ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_YAOSHAIZI {
		user.SendYaoshaiziAck(-2, "当前牌桌不在摇色子状态！")
		return
	}

	if !user.IsOwner() {
		user.SendYaoshaiziAck(-3, "您不是房主，无法摇色子！")
		return
	}

	//停止倒计时timer
	if user.Desk.YaoshaiziTimer != nil {
		user.Desk.YaoshaiziTimer.Stop()
		user.Desk.YaoshaiziTimer = nil
	}

	//摇色子
	list := make([]int32, 4)
	for i:=0; i<4; i++ {
		list[i] = rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(6) + 1
	}
	shaizi_type, _ := ParseShaiziType(list)

	msg := &ddproto.LwyGameEndOne{
		Header: commonNewPorot.NewHeader(),
		Type: shaizi_type.Enum(),
		UserScore: []*ddproto.LwyShaiziResultItem{},
	}

	//烂点，则重新摇
	if shaizi_type == ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_LAN_DIAN {
		//广播出去
		user.Desk.BroadCast(msg)

		//重新摇色子
		user.DoYaoshaizi()
		return
	}

	//更改状态为等待准备
	user.Desk.Status = ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_RESULT.Enum()

	//开始算分
	user_map :=  map[uint32]int64{}
	for _,u := range user.Desk.Users {
		user_map[u.GetUserId()] = 0
	}
	for _,u := range user.Users {
		for _, item := range u.ChizhuDetail {
			if shaizi_type == ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_7_DIAN {
				//摇到7点，则吃7输，吃8赢
				user_map[u.GetUserId()] -= item.GetChi7Score()
				user_map[item.GetFrom()] += item.GetChi7Score()

				user_map[u.GetUserId()] += item.GetChi8Score()
				user_map[item.GetFrom()] -= item.GetChi8Score()
			}else {
				//摇到8点，则吃8输，吃7赢
				user_map[u.GetUserId()] += item.GetChi7Score()
				user_map[item.GetFrom()] -= item.GetChi7Score()

				user_map[u.GetUserId()] -= item.GetChi8Score()
				user_map[item.GetFrom()] += item.GetChi8Score()
			}
		}
	}

	//算分结果
	for _,u := range user.Desk.Users {
		if u == nil || !u.GetIsOnGamming() {
			continue
		}

		score := user_map[u.GetUserId()]
		*u.Bill.Score += score

		msg.UserScore = append(msg.UserScore, &ddproto.LwyShaiziResultItem{
			UserId: proto.Uint32(u.GetUserId()),
			Score: proto.Int64(score),
			AllScore: proto.Int64(u.Bill.GetScore()),
		})
	}

	//将结果广播出去
	user.Desk.BroadCast(msg)

	//单局牌局结束
	user.Desk.DoEnd()
}

//发起解散房间
func (user *User) DoApplyDissolve() error {
	log.T("用户%d发起解散房间", user.GetUserId())

	if user.Desk.GetIsOnDissolve() == true {
		log.T("用户%d发起解散房间失败，原因：%s", user.GetUserId(), "正处于解散投票阶段，不要重复申请！")
		return errors.New("正处于解散投票阶段，不要重复申请！")
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
		if u != nil {
			if u.GetDissolveState() == 0 {
				//如果离线，则自动同意
				if u.GetIsOnline() == false {
					*u.DissolveState = 1
				} else {
					log.T("房间%d解散房间失败，原因：用户%d未投票", user.Desk.GetPassword(), user.GetUserId())
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
		log.T("房间%d解散房间成功", user.Desk.GetPassword())
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
	if !desk.GetIsCoinRoom() && desk.GetCircleNo() == desk.DeskOption.GetBoardsCout() {
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
