package skeleton

import (
	"casino_majiang/msg/protogo"
	"casino_majianagv2/core/majiangv2"
	"github.com/golang/protobuf/proto"
	"reflect"
	"casino_majiang/service/majiang"
	"casino_common/common/log"
)

//玩家发送消息
func (u *SkeletonMJUser) DoRobotAct(msg proto.Message) error {
	if msg == nil {
		log.T("AI收到的msg 是 nil 直接返回...")
		return nil
	}

	t := reflect.TypeOf(msg)
	switch t {
	case reflect.TypeOf(&mjproto.Game_SendGameInfo{}):
		u.DoRobotReady(msg.(*mjproto.Game_SendGameInfo))

	case reflect.TypeOf(&mjproto.Game_BroadcastBeginExchange{}):
		u.DoRobotExcahnge()

	case reflect.TypeOf(&mjproto.Game_BroadcastBeginDingQue{}):
		u.DoRobotDingQue()

	case reflect.TypeOf(&mjproto.Game_OverTurn{}):
		ack := msg.(*mjproto.Game_OverTurn)
		u.DoRobotOverTurn(ack)

	case reflect.TypeOf(&mjproto.Game_AckActPeng{}):
		ack := msg.(*mjproto.Game_AckActPeng)
		u.DoRobotAfterPeng(ack)
	}
	return nil
}

func (u *SkeletonMJUser) DoRobotReady(ack *mjproto.Game_SendGameInfo) {
	if ack.GetSenderUserId() == u.GetUserId() {
		log.T("AI玩家[%v]收到进入成功的ack,现在开始准备", u.GetUserId())
		//HandlerGame_Ready(u.GetUserId(), nil)

		userId := u.GetUserId()
		//开始准备
		err := u.GetDesk().Ready(u.GetUserId()) //ai玩家准备
		if err != nil {
			log.E("AI玩家[%v]准备失败.err %v", userId, err)
			return
		}
	}
}

func (u *SkeletonMJUser) DoRobotDingQue() {
	log.T("AI玩家[%v]收到开始定缺的ack,现在开始定缺", u.GetUserId())
	//开始换牌
	err := u.GetDesk().DingQue(u.GetUserId(), majiangv2.GetDingQUe(u.GameData.HandPai)) //ai玩家定缺
	if err != nil {
		log.E("AI玩家[%v]定缺失败...", u.GetUserId())
	}
}

func (u *SkeletonMJUser) DoRobotExcahnge() {
	log.T("AI玩家[%v]收到开始换三张的ack,现在开始换三张", u.GetUserId())
	//开始换牌
}

func (u *SkeletonMJUser) DoRobotOverTurn(ack *mjproto.Game_OverTurn) {
	if ack.GetUserId() == u.GetUserId() {
		log.T("AI玩家[%v]收到发给自己的overTurn的ack[%v]", u.GetUserId(), ack)
		//1.胡牌
		if ack.GetCanHu() {
			log.T("AI玩家[%v]收到overTurn的ack,开始胡牌", u.GetUserId())
			err := u.GetDesk().ActHu(u.GetUserId()) //AI玩家开始胡牌
			if err != nil {
				log.E("AI 玩家[%v]胡牌失败...", u.GetUserId())
			}
			return
		}

		//2,杠牌
		if ack.GetCanGang() {
			gangCardId := int32(-1)
			if ack.GetActType() == majiang.OVER_TURN_ACTTYPE_MOPAI {
				log.T("AI玩家[%v]收到overTurn的ack,开始暗杠", u.GetUserId())
				gangCardId = ack.GangCards[0].GetId()
			} else {
				log.T("AI玩家[%v]收到overTurn的ack,开始明杠", u.GetUserId())
				gangCardId = ack.ActCard.GetId()
			}

			u.GetDesk().ActGang(u.GetUserId(), gangCardId, false) //AI玩家开始杠牌
			return
		}

		//3,碰牌
		if ack.GetCanPeng() {
			log.T("AI玩家[%v]收到overTurn的ack,开始碰牌", u.GetUserId())
			u.GetDesk().ActPeng(u.GetUserId()) //ai 玩家开始碰牌
			return
		}

		//4，打牌
		u.DoRobotSendOutCard(ack)
	}
}

//打最后一张牌
func (u *SkeletonMJUser) DoRobotAfterPeng(ack *mjproto.Game_AckActPeng) {
	if ack.GetUserIdIn() == u.GetUserId() {
		u.DoRobotSendOutCard(nil)
	}
}

//统一的打牌方法
/**
	打牌的逻辑:
	1,先打缺的牌

 */
func (u *SkeletonMJUser) DoRobotSendOutCard(ack *mjproto.Game_OverTurn) {
	var paiId int32
	//todo 暂时这样处理,如果ack中有叫,直接打，如果没有 自动选择一张打
	if ack != nil && ack.JiaoInfos != nil && len(ack.JiaoInfos) > 0 {
		paiId = ack.JiaoInfos[0].OutCard.GetId() //AI玩家开始打牌
	} else {
		//默认值
		paiId = majiangv2.GetOutPai(u.GameData.HandPai).GetIndex()
	}

	//最后打第一张牌
	log.T("AI玩家[%v]收到overTurn的ack,开始打牌", u.GetUserId())
	u.GetDesk().ActOut(u.GetUserId(), paiId, false) //AI玩家打牌
}
