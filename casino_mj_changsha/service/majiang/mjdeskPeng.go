package majiang

import (
	"errors"
	"casino_common/common/log"
	"casino_mj_changsha/msg/funcsInit"
	"casino_common/common/consts"
)

func (d *MjDesk) ActPeng(userId uint32) error {
	log.T("锁日志: %v ActPeng(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()

	//1，检测玩家是否可以进行碰的操作
	err := d.CheckActUser(userId, ACTTYPE_PENG)
	if err != nil {
		log.E("%v没有轮到当钱的玩家[%v]碰牌", d.DlogDes(), userId)
		return errors.New("没有轮到你操作")
	}

	//1找到玩家
	user := d.GetUserByUserId(userId)
	if user == nil {
		return errors.New("服务器错误碰牌失败")
	}

	//检测是否可以碰牌
	if d.CheckCase == nil {
		log.E("用户[%v]碰牌的时候，没有找到可以碰的牌...", userId)
		return errors.New("服务器错误碰牌失败")
	}

	//new 如果checkCase 里面没有人胡，那么可以直接碰
	if d.GetCheckCase().GetHuBean(CHECK_CASE_BEAN_STATUS_CHECKING) == nil {
		d.checkPeng <- true
	}
	d.Unlock()

	//接收信号，继续操作
	if canCheckPeng := <-d.checkPeng; !canCheckPeng {
		//已经有人操作了，现在不能碰了
		return nil
	}

	d.Lock()
	defer d.Unlock()

	//停止定时器
	if d.overTurnTimer != nil {
		d.overTurnTimer.Stop()
	}

	//2.1开始碰牌的操作
	pengPai := d.CheckCase.CheckMJPai
	outUser := d.GetUserByUserId(d.CheckCase.GetUserIdOut())
	canPeng := user.GameData.HandPai.GetCanPeng(pengPai)
	if !canPeng {
		//如果不能碰，直接返回
		log.E("玩家[%v]碰牌id[%v]-[%v]的时候，出现错误，碰不了...", userId, pengPai.GetIndex(), pengPai.LogDes())
		return errors.New("服务器出现错误..")
	}

	user.GameData.HandPai.InPai = nil
	user.GameData.HandPai.PengPais = append(user.GameData.HandPai.PengPais, pengPai) //碰牌
	user.DelGuoHuInfo()                                                              //碰牌之后删除过胡的信息

	//碰牌之后，需要删掉的自己手里面碰牌的对子..
	var pengKeys []int32
	//pengKeys = append(pengKeys, pengPai.GetIndex())
	for _, pai := range user.GameData.HandPai.Pais {
		if pai != nil && pai.GetClientId() == pengPai.GetClientId() {
			user.GameData.HandPai.PengPais = append(user.GameData.HandPai.PengPais, pai) //碰牌
			pengKeys = append(pengKeys, pai.GetIndex())
			//碰牌只需要拆掉手里的两张牌
			if len(pengKeys) == 2 {
				break;
			}
		}
	}

	//2.2 删除手牌
	for _, key := range pengKeys {
		user.GameData.HandPai.DelHandlPai(key)
	}

	//2.3 删除打牌的out牌
	errDelOut := outUser.GameData.HandPai.DelOutPai(pengPai.GetIndex())
	if errDelOut != nil {
		log.E("碰牌的时候，删除打牌玩家的out牌[%v]...", pengPai)
	}

	//3,生成碰牌信息
	//user.GameData.

	//4,处理 checkCase
	d.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), CHECK_CASE_BEAN_STATUS_CHECKED)
	d.CheckCase.UpdateChecStatus(CHECK_CASE_STATUS_CHECKED) //碰牌之后，checkcase处理完毕
	d.SetAATUser(user.GetUserId(), MJDESK_ACT_TYPE_DAPAI)   //碰牌之后需要设置为activeUser轮到用户打牌

	//5,发送碰牌的广播
	ack := newProto.NewGame_AckActPeng()
	ack.JiaoInfos = d.GetJiaoInfos(user)
	*ack.Header.Code = consts.ACK_RESULT_SUCC
	*ack.UserIdOut = d.CheckCase.GetUserIdOut()
	*ack.UserIdIn = user.GetUserId()
	//组装牌的信息
	for _, ackpai := range user.GameData.HandPai.PengPais {
		if ackpai != nil && ackpai.GetClientId() == pengPai.GetClientId() {
			ack.PengCard = append(ack.PengCard, ackpai.GetCardInfo())
		}
	}

	user.SendOverTurn(ack) //碰牌之后的ACK
	d.BroadCastProtoExclusive(ack, user.GetUserId())
	d.CheckCase = nil //设置为nil

	overTurn := d.AfterPengChiChangSha(user)
	user.SendOverTurn(overTurn) //发送overturn
	return nil
}
