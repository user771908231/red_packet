package friendPlay

import (
	"casino_common/common/Error"
	"casino_common/common/log"
	"casino_common/common/consts"
	"casino_majiang/msg/funcsInit"
	"casino_majianagv2/core/majiangv2"
)

var ERR_OUTPAI = Error.NewError(consts.ACK_RESULT_ERROR, "")

//打牌
func (d *FMJDesk) ActOut(userId uint32, paiKey int32, auto bool) error {
	defer Error.ErrorRecovery("actOut")
	log.T("锁日志: %v ActOut(%v,%v)的时候等待锁", d.DlogDes(), userId, paiKey)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActOut(%v,%v)的时候释放锁", d.DlogDes(), userId, paiKey, )
	}()

	outUser := d.GetUserByUserId(userId)
	if outUser == nil {
		log.E("[%v]打牌失败，没有找到玩家[%v]", d.DlogDes(), userId)
		return ERR_OUTPAI
	}

	//判断是否轮到当前玩家打牌了...
	if d.CheckNotActUser(userId) { //打牌
		log.E("%v没有轮到当前玩家%v打牌", d.DlogDes(), userId)
		return ERR_OUTPAI
	}

	//判断是不是在游戏中的状态
	if d.GetStatus().IsNotGaming() {
		log.E("%v玩家%v打牌失败，desk不在游戏状态[%v]", d.DlogDes(), userId, d.GetStatus().IsNotGaming())
		return ERR_OUTPAI
	}

	//停止定时器
	if d.OverTurnTimer != nil {
		d.OverTurnTimer.Stop()
	}

	//得到参数
	outPai := majiangv2.InitMjPaiByIndex(int(paiKey))

	/**
		1,如果是碰牌打牌的时候,inpai为nil，不需要增加
		2,如果是摸牌打牌（杠之后也是摸牌，需要增加in牌...）
	 */

	if outUser.GetGameData().HandPai.InPai != nil {
		outUser.GetGameData().HandPai.AddPai(outUser.GetGameData().HandPai.InPai) //把inpai放置到手牌上
	}
	log.T("玩家打牌之前的手牌:%v", majiangv2.ServerPais2string(outUser.GetGameData().HandPai.Pais))
	errDapai := outUser.GetGameData().HandPai.DelHandlPai(outPai.GetIndex()) //删除要打出去的牌
	if errDapai != nil {
		log.E("[%v]打牌的时候出现错误，没有找到要到的牌,id[%v]", d.DlogDes(), paiKey)
		return ERR_OUTPAI
	}
	log.T("玩家打牌之后的手牌:%v", majiangv2.ServerPais2string(outUser.GetGameData().HandPai.Pais))

	outUser.GetGameData().HandPai.OutPais = append(outUser.GetGameData().HandPai.OutPais, outPai) //自己桌子前面打出的牌，如果其他人碰杠胡了之后，需要把牌删除掉...
	outUser.GetGameData().HandPai.InPai = nil                                                     //打牌之后需要把自己的  inpai给移除掉...
	outUser.GetGameData().DelGuoHuInfo()                                                          //删除过胡的信息
	//打牌之后的逻辑,初始化判定事件
	err := d.InitCheckCase(outPai, outUser) //打牌之后
	if err != nil {
		//表示无人需要，直接给用户返回无人需要
		//给下一个人摸排，并且移动指针
		log.E("%v服务器错误，初始化判定牌的时候出错err[%v]", d.DlogDes(), err)
		return ERR_OUTPAI
	}

	log.T("[%v]玩家[%v]打牌之后InitCheckCase之后的checkCase[%v]", d.DlogDes(), userId, d.CheckCase)
	//回复消息,打牌之后，广播打牌的信息...s
	outUser.GetGameData().DelPreMoGangInfo() //清楚摸牌前的杠牌info
	result := newProto.NewGame_AckSendOutCard()
	*result.UserId = userId
	result.Card = outPai.GetCardInfo()
	d.BroadCastProto(result)

	log.T("[%v]用户[%v]已经打牌结束，开始处理下一个checkCase", d.DlogDes(), userId)
	d.DoCheckCase() //打牌之后，别人判定牌
	return nil
}
