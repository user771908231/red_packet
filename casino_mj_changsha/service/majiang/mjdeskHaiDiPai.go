package majiang

import (
	mjproto        "casino_mj_changsha/msg/protogo"
	"github.com/golang/protobuf/proto"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_common/common/log"
)

func (d *MjDesk) NeedHaidi(userId uint32, need bool) error {
	log.T("锁日志: %v NeedHaidi(%v,%v)的时候等待锁", d.DlogDes(), userId, need)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v NeedHaidi(%v,%v)的时候释放锁", d.DlogDes(), userId, need)
	}()

	user := d.GetUserByUserId(userId)
	if user == nil {
		return Error.NewError(consts.ACK_RESULT_ERROR, "海底捞失败")
	}

	if need {
		user.NeedHaidiStatus = proto.Int32(MJUSER_NEEDHAIDI_STATUS_NEED) //设置为需要海底牌
		user.GameData.HandPai.InPai = d.GetLastMjPai()
		overTrun := d.GetMoPaiOverTurn(user, false) //普通摸牌，用户摸牌的时候,发送一个用户摸牌的overturn
		user.SendOverTurn(overTrun)                 //玩家摸排之后发送overturn
		d.SetNextBanker(userId)                     //谁要海底牌就设置是谁的庄
		log.T("[%v][%v]开始海底捞牌【%v】...", d.DlogDes(), user.GetUserId(), user.UserPai2String(), overTrun)

		//给其他人广播协议
		o2 := &mjproto.Game_OverTurn{
			ActCard: overTrun.ActCard,
			UserId:  proto.Uint32(user.GetUserId()),
		}
		d.BroadCastProtoExclusive(o2, user.GetUserId())

		//如果不能胡牌，直接打出去,如果能乎，让玩家选择胡牌
		if !overTrun.GetCanHu() {
			go d.ActOut(user.GetUserId(), overTrun.GetActCard().GetId(), true) //海底牌不能胡牌，直接打出去...
		}
	} else {
		//玩家不需要海底牌询问下一家是否需要海底牌
		user.NeedHaidiStatus = proto.Int32(MJUSER_NEEDHAIDI_STATUS_REFUSE) //设置为需要海底牌
		//得到下一个玩家
		index := d.getIndexByUserId(userId)
		var needUser *MjUser
		for i := index + 1; i < len(d.GetUsers())+index; i++ {
			tuser := d.GetUsers()[i%len(d.GetUsers())]
			if tuser != nil && tuser.GetNeedHaidiStatus() == MJUSER_NEEDHAIDI_STATUS_DEFAULT {
				log.T("开始判断user[%v]是否判断过要海底牌..", tuser.GetUserId())
				needUser = tuser
				break
			}
		}

		if needUser == nil {
			//没有找到判断海底捞的玩家，游戏结束
			d.LotteryChangSha()
		} else {
			d.enquireHaiDi(needUser)
		}
	}

	return nil
}

//询问玩家是否需要海底牌
func (d *MjDesk) enquireHaiDi(user *MjUser) {
	d.SetAATUser(user.GetUserId(), MJDESK_ACT_TYPE_WAIT_HAIDI) //海底牌
	ack := &mjproto.Game_DealHaiDiCards{
		HaidiCard: d.GetLastMjPai().GetCardInfo(),
		UserId:    proto.Uint32(user.GetUserId()),
	}
	user.WriteMsg(ack)
}
