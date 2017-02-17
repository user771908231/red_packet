package friendPlay

import (
	"casino_common/common/log"
	"casino_majiang/service/majiang"
	"casino_majianagv2/core/api"
	"casino_common/common/consts/tableName"
	"casino_majiang/gamedata/model"
	"casino_common/utils/timeUtils"
	"time"
	"casino_common/utils/db"
	"casino_majiang/msg/funcsInit"
	"casino_common/common/userService"
)

//成都麻将的 结算方式
func (d *FMJDesk) LotteryChengDu() error {
	//结账需要分两中情况
	/**
		1，只剩一个玩家没有胡牌的时候
		2，没有生育麻将的时候.需要分别做处理...
	 */

	//判断是否可以胡牌
	log.T("现在开始处理lottery()的逻辑....")

	//查花猪
	d.ChaHuaZhu()

	//查大叫
	d.ChaDaJiao()

	//1，处理开奖的数据,
	d.DoLottery()

	//发送结束的广播
	d.SendLotteryData()

	//开奖之后 desk需要处理
	d.AfterLottery()

	//判断牌局结束(整场游戏结束)
	if !d.End() {
		//go d.begin()
	}
	return nil
}

//查花猪
/**
	查花猪是查用户是否没有缺
 */
func (d *FMJDesk) ChaHuaZhu() error {
	for _, u := range d.GetFMJUsers() {
		if !d.IsXueLiuChengHe() && u.GetStatus().IsHu() {
			//如果不是血流成河且用户已胡 则不能查u的花猪
			log.T("不查[%v]的花猪", u.GetUserId())
			continue
		}
		if u != nil && u.GetStatus().IsNotHu() {
			if u.IsHuaZhu() {
				log.T("玩家[%v]是花猪", u.GetUserId())
				d.DoHuaZhu(u)
			}
		}
	}
	return nil
}

//花猪玩家需要给封顶分数
func (d *FMJDesk) DoHuaZhu(huazhu *FMJUser) error {
	log.T("开始处理花猪[%v]", huazhu.GetUserId())
	fanTop := d.GetMJConfig().CapMax
	score := d.GetMJConfig().BaseValue * fanTop
	for _, user := range d.GetFMJUsers() {
		if user != nil && !user.IsHuaZhu() {
			//不是花猪的用户都可以收钱
			//判断不是花猪，可以赢钱...
			log.T("DoHuaZhu: 查[%v]的花猪", user.GetUserId())
			user.AddBill(huazhu.GetUserId(), majiang.MJUSER_BILL_TYPE_YING_CHAHUAZHU, "用户查花猪，赢钱", score, nil, d.GetMJConfig().RoomType)
			user.AddStatisticsCountChaHuaZhu(d.GetMJConfig().CurrPlayCount)

			huazhu.AddBill(user.GetUserId(), majiang.MJUSER_BILL_TYPE_SHU_CHAHUAZHU, "用户查花猪，输钱", -score, nil, d.GetMJConfig().RoomType)
			huazhu.AddStatisticsCountBeiChaHuaZhu(d.GetMJConfig().CurrPlayCount)
		}
	}

	return nil
}

//查大叫
/**
	查用户有没有叫
 */
func (d *FMJDesk) ChaDaJiao() error {
	//循环判断谁可以被查叫
	for _, u := range d.GetFMJUsers() {
		if u != nil && u.GetStatus().IsNotHu() && u.GetGameData().IsNotHuaZhu() {
			//用户没有胡 且 不是花猪 可以被查
			jiaoPai := u.GetDesk().GetHuParser().GetJiaoPais(u.GetGameData().HandPai.Pais)
			if jiaoPai == nil || len(jiaoPai) <= 0 {
				d.DoDaJiao(u)
			}
		}
	}
	return nil
}

//获取用户听张的最大番数
func (d *FMJDesk) GetJiaoMaxFan(u api.MjUser) int32 {
	maxFan := 0

	//取得可以叫的infos
	jiaoInfos := d.GetJiaoInfos(u)

	if jiaoInfos == nil {
		return 0
	}
	for _, jiaoInfo := range jiaoInfos {
		paiInfos := jiaoInfo.PaiInfos
		if paiInfos != nil {
			for _, paiInfo := range paiInfos {
				fan := paiInfo.Fan
				if int32(maxFan) < *fan {
					//如果最大番小于当前可叫的番数
					maxFan = int(*fan)
				}
			}
		}
	}
	return int32(maxFan)
}

//成都麻将查大叫
func (d *FMJDesk) DoDaJiao(u *FMJUser) {
	//判断谁可以查u的大叫
	//没听牌的玩家(花猪除外)赔给听牌的玩家 按听牌的最大番型给
	//log.T("开始处理玩家[%v]没叫,开始处理查大叫...", u.GetUserId())

	for _, user := range d.GetFMJUsers() {
		if !d.IsXueLiuChengHe() && u.GetStatus().IsNotHu() {
			//如果不是血流成河且用户已胡 则不能查u的大叫
			continue
		}

		//获得听牌的最大番数
		fan := d.GetJiaoMaxFan(user)
		score := d.GetMJConfig().BaseValue * int64(fan)
		//判断是否有叫
		youjiao := false
		jiaoPai := user.GetDesk().GetHuParser().GetJiaoPais(u.GetGameData().HandPai.Pais)
		if len(jiaoPai) > 0 {
			youjiao = true
		}

		//如果looper不是被查大叫的玩家 且 该looper有听牌 且 该looper没有胡 为该玩家增加赢钱的bill
		if (user.GetUserId() != u.GetUserId()) && youjiao && user.GetStatus().IsNotHu() {
			user.AddBill(u.GetUserId(), majiang.MJUSER_BILL_TYPE_YING_DAJIAO, "用户查大叫，赢钱", score, nil, d.GetMJConfig().RoomType)
			user.AddStatisticsCountChaDaJiao(d.GetMJConfig().CurrPlayCount)

			u.AddBill(user.GetUserId(), majiang.MJUSER_BILL_TYPE_SHU_DAJIAO, "用户被查叫，输钱", -score, nil, d.GetMJConfig().RoomType)
			u.AddStatisticsCountBeiChaJiao(d.GetMJConfig().CurrPlayCount)
		}
	}
}

//需要保存到 ..T_mj_desk_round   ...这里设计到保存数据，战绩相关的查询都要从这里查询
func (d *FMJDesk) DoLottery() error {
	log.T("%v,gameNumber(%v)处理DoLottery()", d.DlogDes(), d.GetMJConfig().GameNumber)
	data := model.T_mj_desk_round{}
	data.DeskId = d.GetMJConfig().DeskId
	data.GameNumber = d.GetMJConfig().GameNumber
	data.BeginTime = timeUtils.String2YYYYMMDDHHMMSS(d.GetMJConfig().BeginTime)
	data.EndTime = time.Now()
	data.UserIds = d.GetUserIds()

	//一次处理每个胡牌的人
	for _, user := range d.GetFMJUsers() {
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
			//开奖之后，玩家的状态修改
			user.AfterLottery()
		}
	}

	//保存数据
	err := db.InsertMgoData(tableName.DBT_MJ_DESK_ROUND, &data)
	if err != nil {
		log.E("dolottery()时保存数据[%v]失败...", data)
	}
	log.T("desk(%v),处理DoLottery(),处理完毕,保存数据data[%v]", d.DlogDes(), data)

	return nil

}

func (d *FMJDesk) SendLotteryData() error {
	//发送开奖的数据,需要得到每个人的winCoinInfo
	result := newProto.NewGame_SendCurrentResult()
	for _, user := range d.GetUsers() {
		if user != nil {
			result.WinCoinInfo = append(result.WinCoinInfo, d.GetWinCoinInfo(user))
		}
	}
	result.BridInfo = d.BirdInfo

	//开始发送开奖的广播
	d.BroadCastProto(result)
	log.T("desk[%v]SendLotteryData()", d.DlogDes())
	return nil
}

func (d *FMJDesk) AfterLottery() error {
	//开奖完成之后的一些处理
	if d.OverTurnTimer != nil {
		d.OverTurnTimer.Stop()
	}

	//把信息更新到mgo
	for _, u := range d.GetUsers() {
		if u != nil {
			userService.UpdateUser2MgoById(u.GetUserId())
		}
	}
	//如果是金币场，需要把短线的，离开的，机器人都踢走
	d.GetStatus().SetStatus(majiang.MJDESK_STATUS_READY) //桌子开始ready
	return nil

}

func (d *FMJDesk) End() bool {
	//判断结束的条件,目前只有局数能判断
	log.T("%v游戏是否End() CurrPlayCount[%v], TotalPlayCount[%v]",
		d.DlogDes(), d.GetMJConfig().CurrPlayCount, d.GetMJConfig().TotalPlayCount)
	//朋友桌有整场结束的概念
	if d.GetMJConfig().CurrPlayCount < d.GetMJConfig().TotalPlayCount {
		//表示游戏还没有结束。。。.
		return false
	} else {
		d.DoEnd()
		return true
	}

}
