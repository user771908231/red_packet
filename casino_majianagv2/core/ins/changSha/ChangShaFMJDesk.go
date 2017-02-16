package changSha

import (
	"casino_majianagv2/core/ins/skeleton"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_common/common/log"
	"casino_majianagv2/core/majiangv2"
)

//长沙麻将朋友桌

//朋友桌麻将的desk
type ChangShaFMJDesk struct {
	*skeleton.SkeletonMJDesk
}

//创建一个朋友桌的desk
func NewChangShaFMJDesk(config *data.SkeletonMJConfig) api.MjDesk {
	//desk 骨架
	desk := &ChangShaFMJDesk{
		SkeletonMJDesk: skeleton.NewSkeletonMJDesk(config),
	}
	desk.HuParser = NewHuParserChangSha()
	return desk
}

func (d *ChangShaFMJDesk) GetCSUser(u interface{}) *ChangShaMJUser {
	return u.(*ChangShaMJUser)
}

//打牌
func (d *ChangShaFMJDesk) ActOut(userId uint32, cardId int32, bu bool) error {
	return nil
}

//进入游戏
func (d *ChangShaFMJDesk) Leave(userId uint32) error {
	return nil
}

//可以把overturn放在一个地方,目前都是摸牌的时候在用
func (d *ChangShaFMJDesk) GetOverTurnByCaseBean(user api.MjUser, isOpen bool) *mjproto.Game_OverTurn {
	overTurn := d.GetMoPaiOverTurn(user, isOpen)
	//这里需要对长沙麻将做特殊处理(主要是杠，补的处理)
	if overTurn.GetCanGang() {
		overTurn.CanBu = proto.Bool(true)
		overTurn.CanGang = proto.Bool(false)
		overTurn.BuCards = overTurn.GangCards
		overTurn.GangCards = nil
		//判断长沙麻将能不能杠
		for _, g := range overTurn.BuCards {
			cang := d.GetCSUser(user).GetCanChangShaGang(majiangv2.InitMjPaiByIndex(int(g.GetId())))
			log.T("判断玩家[%v]对牌[%v]是否可以长沙杠[%v]", user.GetUserId(), g.GetId(), cang)
			if cang {
				overTurn.CanGang = proto.Bool(true)
				overTurn.GangCards = append(overTurn.GangCards, g)
			}
		}
	}
	return overTurn
}

//可以把overturn放在一个地方,目前都是摸牌的时候在用
func (d *ChangShaFMJDesk) GetMoPaiOverTurn(user api.MjUser, isOpen bool) *mjproto.Game_OverTurn {
	overTurn := d.SkeletonMJDesk.GetMoPaiOverTurn(user, isOpen)
	overTurn.JiaoInfos = d.GetJiaoInfos(user)
	//这里需要对长沙麻将做特殊处理(主要是杠，补的处理)
	if overTurn.GetCanGang() {
		overTurn.CanBu = proto.Bool(true)
		overTurn.CanGang = proto.Bool(false)
		overTurn.BuCards = overTurn.GangCards
		overTurn.GangCards = nil
		//判断长沙麻将能不能杠
		for _, g := range overTurn.BuCards {
			cang := d.GetCSUser(user).GetCanChangShaGang(majiangv2.InitMjPaiByIndex(int(g.GetId())))
			log.T("判断玩家[%v]对牌[%v]是否可以长沙杠[%v]", user.GetUserId(), g.GetId(), cang)
			if cang {
				overTurn.CanGang = proto.Bool(true)
				overTurn.GangCards = append(overTurn.GangCards, g)
			}
		}
	}
	return overTurn
}
