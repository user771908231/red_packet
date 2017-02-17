package skeleton

import (
	"casino_majiang/service/majiang"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_common/common/log"
	"errors"
	"casino_majiang/msg/funcsInit"
	"time"
	"casino_majianagv2/core/api"
	"casino_common/common/Error"
	"github.com/name5566/leaf/gate"
)

func (d *SkeletonMJDesk) TryReEnter(userId uint32, a gate.Agent) (re bool) {
	//1设置agent
	//2,是否是离开之后重新进入房间
	userLeave := d.GetUserByUserId(userId)
	if userLeave != nil {
		log.T("玩家[%v]断线重连....", userId)
		userLeave.ReEnterDesk(a)
		d.SendGameInfo(userId, mjproto.RECONNECT_TYPE_RECONNECT)
		return true
	}
	return false
}

//发送重新连接之后的overTurn
func (d *SkeletonMJDesk) SendReconnectOverTurn(userId uint32) error {
	log.T("[%v]开始处理 sendReconnectOverTurn(%v),当前desk.status(%v),", d.DlogDes(), userId, d.GetStatus())

	//得到玩家
	user := d.GetSkeletonMJUser(d.GetUserByUserId(userId))
	if user == nil {
		log.T("发送SendReconnectOverTurn(%v)失败,因为没有找到玩家", userId)
		return errors.New("发送SendReconnectOverTurn()失败,因为没有找到玩家")
	}

	//开发发送overTurn
	if d.GetStatus().IsPreparing() && !user.GetStatus().IsReady() {
		//给玩家发送开始准备的消息...
		log.T("sendReconnectOverTurn，给user[%v]发送准备的消息....", userId)
	} else if d.GetStatus().IsExchange && user.GetStatus().Exchange {
		//给玩家发送换牌的消息...
		log.T("sendReconnectOverTurn，给user[%v]发送换牌的消息....", userId)
		beginExchange := newProto.NewGame_BroadcastBeginExchange()
		*beginExchange.Reconnect = true
		user.WriteMsg(beginExchange)

	} else if d.IsDingQue() && !user.GetStatus().DingQue {
		log.T("sendReconnectOverTurn，给user[%v]发送定缺的消息....", userId)

		//给玩家发送定缺的信息...
		beginQue := newProto.NewGame_BroadcastBeginDingQue()
		*beginQue.Reconnect = true
		user.WriteMsg(beginQue)

	} else if d.GetStatus().IsGaming() && user.GetUserId() == d.GetMJConfig().ActUser {

		//游戏中的情况，发送act的消息,这里需要更具当前的状态来发送overTurn
		if d.GetMJConfig().ActType == majiang.MJDESK_ACT_TYPE_MOPAI {
			log.T("sendReconnectOverTurn，给user[%v]发送摸牌的消息....", userId)
			overTrun := d.GetMoPaiOverTurn(user, false) //重新进入房间之后
			user.WriteMsg(overTrun)
			log.T("玩家重新进入游戏之后 [%v]开始摸牌【%v】...", user.GetUserId(), overTrun)

		} else if d.GetMJConfig().ActType == majiang.MJDESK_ACT_TYPE_DAPAI {
			//发送打牌的协议:这里只有在碰的时候才会出现这种情况，其他的时候都是摸牌(摸牌，刚拍之后摸牌)
			log.T("sendReconnectOverTurn，给user[%v]发送打牌的消息....", userId)
			overTurn := newProto.NewGame_OverTurn()
			*overTurn.UserId = userId
			*overTurn.CanGang = false
			*overTurn.CanPeng = false
			*overTurn.CanHu = false
			overTurn.CanBu = proto.Bool(false)
			*overTurn.ActType = majiang.OVER_TURN_ACTTYPE_MOPAI
			*overTurn.PaiCount = d.GetRemainPaiCount()
			overTurn.ActCard = majiang.NewBackPai()
			user.WriteMsg(overTurn)
		} else if d.GetMJConfig().ActType == majiang.MJDESK_ACT_TYPE_WAIT_CHECK {
			log.T("sendReconnectOverTurn，给user[%]发送checkCase的消息....", userId)

			caseBean := d.CheckCase.GetBeanByUserIdAndStatus(user.GetUserId(), majiang.CHECK_CASE_BEAN_STATUS_CHECKING)
			if caseBean == nil {
				log.E("没有找到玩家[%v]对应的checkBean", user.GetUserId())
				return errors.New("玩家重新进入房间发送check overturn的时候出错")
			}

			//找到需要判断bean之后，发送给判断人	//send overTurn
			overTurn := d.GetOverTurnByCaseBean(d.CheckCase.CheckMJPai, caseBean, majiang.OVER_TURN_ACTTYPE_OTHER) //重新进入游戏
			*overTurn.Time = int32(user.GetWaitTime() - time.Now().Unix())
			///发送overTurn 的信息
			log.T("开始发送玩家[%v]断线重连的overTurn[%v]", user.GetUserId(), overTurn)
			user.WriteMsg(overTurn)
		} else if d.GetMJConfig().ActType == majiang.MJDESK_ACT_TYPE_WAIT_HAIDI {
			log.T("sendReconnectOverTurn，给user[%]发送判定海底牌的消息....", userId)
			d.enquireHaiDi(user)
		}
	} else if d.GetStatus().S() == majiang.MJDESK_STATUS_QISHOUHU {
		//长沙麻将起手胡牌的阶段...查看玩家是否需要发送起手胡牌
		bean := d.GetCheckCase().GetNextBean()
		if bean != nil && bean.GetUserId() == user.GetUserId() && bean.GetCanHu() {
			//给玩家发送起手胡的ack
			overTurn := &mjproto.Game_ChangshQiShouHuOverTurn{
				Header: &mjproto.ProtoHeader{
					UserId: bean.UserId,
				},
			}
			log.T("短线重连之后，开始发送起手胡overTurn[%v]", overTurn)                             //打日志
			user.SendOverTurn(overTurn)                                               //发送OverTurn
			d.SetActUserAndType(bean.GetUserId(), majiang.MJDESK_ACT_TYPE_WAIT_CHECK) //DoCheckCase 设置当前活动的玩家
		}
	}

	log.T("%v开始处理 sendReconnectOverTurn(%v),当前desk.status(%v)----处理完毕...", d.DlogDes(), userId, d.GetStatus())
	return nil
}

//询问玩家是否需要海底牌
func (d *SkeletonMJDesk) enquireHaiDi(user api.MjUser) {
	d.SetActiveUser(user.GetUserId())
	d.SetActUserAndType(user.GetUserId(), majiang.MJDESK_ACT_TYPE_WAIT_HAIDI) //海底牌
	ack := &mjproto.Game_DealHaiDiCards{
		HaidiCard: d.GetLastMjPai().GetCardInfo(),
		UserId:    proto.Uint32(user.GetUserId()),
	}
	user.WriteMsg(ack)
}

//新增加一个玩家
func (d *SkeletonMJDesk) AddUserBean(user api.MjUser) error {
	//根据房间类型判断人数是否已满
	if d.IsPlayerEnough() {
		return Error.NewFailError("房间已满，加入桌子失败")
	}

	//找到座位
	seatIndex := -1
	for i, u := range d.Users {
		if u == nil {
			seatIndex = i
			break
		}
	}

	//如果找到座位那么，增加用户，否则返回错误信息
	if seatIndex >= 0 {
		d.Users[seatIndex] = user
		return nil
	} else {
		return Error.NewFailError("没有找到合适的位置，加入桌子失败")
	}
}

func (d *SkeletonMJDesk) SendGameInfo(userId uint32, reconnect mjproto.RECONNECT_TYPE) {
	gameinfo := d.GetGame_SendGameInfo(userId, reconnect)
	log.T("[%v]用户[%v]进入房间,reconnect[%v]之后", d.DlogDes(), userId, reconnect)
	d.BroadCastProto(gameinfo)

	//如果是重新进入房间，需要发送重近之后的处理
	if reconnect == mjproto.RECONNECT_TYPE_RECONNECT {
		time.Sleep(time.Second * 3)
		d.SendReconnectOverTurn(userId)
	}
}
