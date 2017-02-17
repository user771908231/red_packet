package skeleton

import (
	"casino_common/common/log"
	"errors"
	"casino_majiang/msg/funcsInit"
	"casino_common/common/consts"
	"casino_majiang/service/majiang"
)

func (d *SkeletonMJDesk) ActPeng(userId uint32) error {
	log.T("锁日志: %v ActPeng(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActPeng(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//1找到玩家
	user := d.GetUserByUserId(userId)
	if user == nil {
		return errors.New("服务器错误碰牌失败")
	}

	//检测当前的后动玩家是否正确
	if d.CheckNotActUser(userId) { //碰牌
		log.E("%v没有轮到当钱的玩家[%v]碰牌", d.DlogDes(), userId)
		return errors.New("没有轮到你操作")
	}

	//检测是否可以碰牌
	if d.CheckCase == nil {
		log.E("用户[%v]碰牌的时候，没有找到可以碰的牌...", userId)
		return errors.New("服务器错误碰牌失败")
	}

	//停止定时器
	if d.OverTurnTimer != nil {
		d.OverTurnTimer.Stop()
	}

	//2.1开始碰牌的操作
	pengPai := d.CheckCase.CheckMJPai
	outUser := d.GetUserByUserId(d.CheckCase.GetUserIdOut())
	canPeng := user.GetGameData().HandPai.GetCanPeng(pengPai)
	if !canPeng {
		//如果不能碰，直接返回
		log.E("玩家[%v]碰牌id[%v]-[%v]的时候，出现错误，碰不了...", userId, pengPai.GetIndex(), pengPai.LogDes())
		return errors.New("服务器出现错误..")
	}

	user.GetGameData().HandPai.InPai = nil
	user.GetGameData().HandPai.PengPais = append(user.GetGameData().HandPai.PengPais, pengPai) //碰牌
	user.GetGameData().DelGuoHuInfo()                                                          //删除过胡的信息

	//碰牌之后，需要删掉的自己手里面碰牌的对子..
	var pengKeys []int32
	//pengKeys = append(pengKeys, pengPai.GetIndex())
	for _, pai := range user.GetGameData().HandPai.Pais {
		if pai != nil && pai.GetClientId() == pengPai.GetClientId() {
			user.GetGameData().HandPai.PengPais = append(user.GetGameData().HandPai.PengPais, pai) //碰牌
			pengKeys = append(pengKeys, pai.GetIndex())
			//碰牌只需要拆掉手里的两张牌
			if len(pengKeys) == 2 {
				break;
			}
		}
	}

	//2.2 删除手牌
	for _, key := range pengKeys {
		user.GetGameData().HandPai.DelHandlPai(key)
	}

	//2.3 删除打牌的out牌
	errDelOut := outUser.GetGameData().HandPai.DelOutPai(pengPai.GetIndex())
	if errDelOut != nil {
		log.E("碰牌的时候，删除打牌玩家的out牌[%v]...", pengPai)
	}

	//3,生成碰牌信息
	//user.GameData.

	//4,处理 checkCase
	d.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), majiang.CHECK_CASE_BEAN_STATUS_CHECKED)
	d.CheckCase.UpdateChecStatus(majiang.CHECK_CASE_STATUS_CHECKED)      //碰牌之后，checkcase处理完毕
	d.SetActiveUser(user.GetUserId())                                    //碰牌之后需要设置为activeUser
	d.SetActUserAndType(user.GetUserId(), majiang.MJDESK_ACT_TYPE_DAPAI) //轮到用户打牌

	//5,发送碰牌的广播
	ack := newProto.NewGame_AckActPeng()
	ack.JiaoInfos = d.GetJiaoInfos(user)
	*ack.Header.Code = consts.ACK_RESULT_SUCC
	*ack.UserIdOut = d.CheckCase.GetUserIdOut()
	*ack.UserIdIn = user.GetUserId()
	//组装牌的信息
	for _, ackpai := range user.GetGameData().HandPai.PengPais {
		if ackpai != nil && ackpai.GetClientId() == pengPai.GetClientId() {
			ack.PengCard = append(ack.PengCard, ackpai.GetCardInfo())
		}
	}

	user.SendOverTurn(ack)
	d.BroadCastProtoExclusive(ack, user.GetUserId())
	//最后设置checkCase = nil
	d.CheckCase = nil //设置为nil
	return nil
}
