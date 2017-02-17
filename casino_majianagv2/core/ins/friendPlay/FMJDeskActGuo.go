package friendPlay

import (
	"casino_common/common/log"
	"errors"
	"casino_majiang/msg/funcsInit"
	"casino_majiang/service/majiang"
)

func (desk *FMJDesk) ActGuoChengdu(userId uint32) error {
	if desk.CheckNotActUser(userId) { //过牌
		log.E("[%v]没有轮到玩家[%v]操作...", userId)
		return errors.New("没哟轮到当前的玩家")
	}

	user := desk.GetFMJUser(userId)
	if desk.CheckCase == nil {
		log.E("%v玩家【%v】过牌的时候出错，因为checkCase为nil", desk.DlogDes(), userId)
		return errors.New("没有找到对应的玩家")
	}

	//停止定时器
	if desk.OverTurnTimer != nil {
		desk.OverTurnTimer.Stop()
	}

	if desk.GetCheckCase() == nil {
		result := newProto.NewGame_AckActGuo()
		result.Header = newProto.SuccessHeader()
		*result.UserId = user.GetUserId()
		user.WriteMsg(result)
		return nil
	} else {
		//两人麻将 剩余四张牌时能胡则不能过
		//todo 这里应该在jiaoInfo 那里就设置canGuo = false
		if desk.GetMJConfig().PlayerCountLimit == 2 {
			canHu, _, _, _, _, _ := desk.HuParser.GetCanHu(user.GetGameData().GetHandPai(), user.GetGameData().GetHandPai().GetInPai(), false, 0)
			if canHu && (desk.GetRemainPaiCount() <= 4) {
				//能胡且剩余牌数小于等于4
				log.E("玩家【%v】过牌失败，因为剩余牌数<=4 且 能胡牌", userId)
				return errors.New("不能过来")
			}
		}

		//添加一个过hu的info,下次init的时候，需要判断是否有这个guohu
		user.AddGuoHuInfo(desk.CheckCase)

		err := desk.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), majiang.CHECK_CASE_BEAN_STATUS_CHECKED) // update checkCase...
		if err != nil {
			log.T("过牌的时候失败，err[%v]", err)
		}

		//返回信息,过 只返回给过的
		result := newProto.NewGame_AckActGuo()
		result.Header = newProto.SuccessHeader()
		*result.UserId = user.GetUserId()
		user.WriteMsg(result)

		//进行下一个判断
		err = desk.DoCheckCase() //过牌之后，处理下一个判定牌
		if err != nil {
			log.E("%v,玩家[%v]过牌之后，docheckcase的时候失败", desk.DlogDes(), userId)
		}
	}
	return nil
}
