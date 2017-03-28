package majiang

import (
	"casino_common/common/log"
	"errors"
	"casino_common/common/consts"
	mjproto        "casino_mj_changsha/msg/protogo"
	"github.com/golang/protobuf/proto"
	"casino_mj_changsha/msg/funcsInit"
)

//吃牌的逻辑
func (d *MjDesk) ActChi(userId uint32, chooseCards []*mjproto.CardInfo) error {
	log.T("锁日志: %v ActChi(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActChi(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//1找到玩家
	user := d.GetUserByUserId(userId)
	if user == nil {
		return errors.New("服务器错误吃失败")
	}

	//检测当前的后动玩家是否正确
	err := d.CheckActUser(userId, ACTTYPE_PENG)
	if err != nil {
		log.E("%v没有轮到当钱的玩家[%v]吃牌", d.DlogDes(), userId)
		return errors.New("没有轮到你操作")
	}

	checkBean := d.CheckCase.GetNextBean()
	//检测是否可以碰牌
	if checkBean == nil {
		log.E("用户[%v]吃牌的时候，没有找到可以吃的牌...", userId)
		return errors.New("服务器错误吃失败")
	}

	//停止定时器
	if d.overTurnTimer != nil {
		d.overTurnTimer.Stop()
	}

	chiPai := checkBean.GetCheckPai() //吃的牌
	choosePai := make([]*MJPai, 3)
	choosePai[0] = InitMjPaiByIndex(int(chooseCards[0].GetId()))
	choosePai[1] = InitMjPaiByIndex(int(chooseCards[1].GetId()))
	choosePai[2] = InitMjPaiByIndex(int(chooseCards[2].GetId()))

	outUser := d.GetUserByUserId(d.CheckCase.GetUserIdOut())
	canChi, _ := user.GameData.HandPai.GetCanChi(chiPai)
	if !canChi {
		//如果不能碰，直接返回
		log.E("玩家[%v]吃牌id[%v]-[%v]的时候，出现错误，吃不了...", userId, chiPai.GetIndex(), chiPai.LogDes())
		return errors.New("服务器出现错误..")
	}

	user.DelGuoHuInfo() //吃牌之后删除过胡的信息（其实吃牌之后可以不用判断，因为吃牌之后和碰牌一样 都会打牌，可以也是打牌之后判断）
	var chiKeys []int32
	for _, cp := range choosePai {
		for _, hp := range user.GetGameData().GetHandPai().GetPais() {
			if cp.GetIndex() == hp.GetIndex() {
				chiKeys = append(chiKeys, cp.GetIndex())
				break
			}
		}
	}

	user.GameData.HandPai.ChiPais = append(user.GameData.HandPai.ChiPais, choosePai...) //碰牌

	//2.2 删除手牌
	for _, key := range chiKeys {
		user.GameData.HandPai.DelHandlPai(key)
	}

	//2.3 删除打牌的out牌
	errDelOut := outUser.GameData.HandPai.DelOutPai(chiPai.GetIndex())
	if errDelOut != nil {
		log.E("吃牌的时候，删除打牌玩家的out牌[%v]...", chiPai)
	}

	//3,生成吃牌的信息
	//user.GameData.

	//4,处理 checkCase
	d.CheckCase.UpdateCheckBeanStatus(user.GetUserId(), CHECK_CASE_BEAN_STATUS_CHECKED)
	d.CheckCase.UpdateChecStatus(CHECK_CASE_STATUS_CHECKED) //碰牌之后，checkcase处理完毕
	d.SetAATUser(user.GetUserId(), MJDESK_ACT_TYPE_DAPAI)   //吃牌之后轮到用户打牌

	//5,发送碰牌的广播
	ack := &mjproto.Game_AckActChi{
		Header: &mjproto.ProtoHeader{
			Code: proto.Int32(consts.ACK_RESULT_SUCC),
		},
		UserIdOut: proto.Uint32(d.CheckCase.GetUserIdOut()),
		UserIdIn:  proto.Uint32(user.GetUserId()),
		ChiCard:   chooseCards, //迟的是哪些牌
		JiaoInfos: d.GetJiaoInfos(user),
	}
	d.BroadCastProto(ack)
	d.CheckCase = nil //设置为nil	吃牌操作之后设置CheckCase 为nil

	if d.IsChangShaMaJiang() {
		//处理after长沙吃杠
		overTurn := d.AfterPengChiChangSha(user) //长沙麻将吃牌之后，发送补 杠的overturn
		user.SendOverTurn(overTurn)              //发送overturn
	}
	return nil
}

//判断授牌是否可以杠(长沙碰，吃之后的操作)
func (d *MjDesk) AfterPengChiChangSha(user *MjUser) *mjproto.Game_OverTurn {

	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = user.GetUserId()         //这个是摸牌的，所以是广播...
	*overTurn.PaiCount = d.GetRemainPaiCount()  //桌子剩余多少牌
	*overTurn.ActType = OVER_TURN_ACTTYPE_MOPAI //摸牌
	*overTurn.Time = 30
	overTurn.ActCard = user.GameData.HandPai.InPai.GetBackPai()
	overTurn.CanHu = proto.Bool(false)
	*overTurn.CanPeng = false //是否可以碰牌

	canGangBool, gangPais := user.GameData.HandPai.GetCanGang(nil, d.GetRemainPaiCount()) //是否可以杠牌
	if canGangBool && gangPais != nil && len(gangPais) > 0 {
		overTurn.CanGang = proto.Bool(canGangBool)
		overTurn.CanGuo = proto.Bool(true)
		for _, g := range gangPais {
			overTurn.GangCards = append(overTurn.GangCards, g.GetCardInfo())
		}
	}

	//对长沙麻将做特殊处理
	overTurn.JiaoInfos = d.GetJiaoInfos(user)
	//这里需要对长沙麻将做特殊处理(主要是杠，补的处理)
	if overTurn.GetCanGang() {
		overTurn.CanBu = proto.Bool(true)
		overTurn.CanGang = proto.Bool(false)
		overTurn.BuCards = overTurn.GangCards
		overTurn.GangCards = nil
		//判断长沙麻将能不能杠
		for _, g := range overTurn.BuCards {
			cang := user.GetCanChangShaGang(InitMjPaiByIndex(int(g.GetId()))) // 摸牌的时候 判断能否gang
			log.T("判断玩家[%v]对牌[%v]是否可以长沙杠[%v]", user.GetUserId(), g.GetId(), cang)
			if cang {
				overTurn.CanGang = proto.Bool(true)
				overTurn.GangCards = append(overTurn.GangCards, g)
			}
		}
	}

	//判断是否可以补或者杠
	if overTurn.GetCanBu() || overTurn.GetCanGang() {
		return overTurn
	} else {
		return nil
	}
}
