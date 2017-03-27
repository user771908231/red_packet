package majiang

import (
	"casino_common/common/log"
	"errors"
	"casino_mj_changsha/msg/funcsInit"
	"github.com/golang/protobuf/proto"
)

/**
	只有判断别人打的牌的时候，需要过的时候才会请求这个协议，自己摸牌 需不需要过的时候不需要请求这个协议...
*/

func (desk *MjDesk) ActGuoChangSha(userId uint32) error {

	err := desk.CheckActUser(userId, ACTTYPE_GUO)
	//过牌
	if err != nil {
		log.E("[%v]没有轮到玩家[%v]操作...", userId)
		return errors.New("没哟轮到当前的玩家")
	}

	user := desk.GetUserByUserId(userId)
	if user == nil {
		return errors.New("没哟轮到当前的玩家")
	}

	//停止定时器
	if desk.overTurnTimer != nil {
		desk.overTurnTimer.Stop()
	}

	//起手胡阶段，过的处理
	if desk.GetStatus() == MJDESK_STATUS_QISHOUHU {
		go desk.ActQiShouHu(userId, false)
		return nil
	}

	//别人打牌时候的判断
	if desk.GetCheckCase() != nil {
		user.AddGuoHuInfo(desk.CheckCase)
		bean := desk.GetCheckCase().GetNextBean()
		if bean != nil {
			//设置为已经检测过了
			bean.CheckStatus = proto.Int32(CHECK_CASE_BEAN_STATUS_CHECKED)
		}
		//返回信息,过 只返回给过的
		result := newProto.NewGame_AckActGuo()
		result.Header = newProto.SuccessHeader()
		*result.UserId = user.GetUserId()
		user.WriteMsg(result)

		//进行下一个判断
		err := desk.DoCheckCase() //过牌之后，处理下一个判定牌
		if err != nil {
			log.E("%v,玩家[%v]过牌之后，docheckcase的时候失败", desk.DlogDes(), userId)
		}
		return nil
	}

	//这里需要判断，如果是杠(不是补牌)之后的过牌,并且是别人点炮，那么需要直接去请求打牌
	if user.GetChangShaGangStatus() {
		log.T("%v 长沙麻将玩家[%v]杠牌之后，自摸能胡牌，点击了过", desk.DlogDes(), userId)
		//返回信息,过 只返回给过的
		result := newProto.NewGame_AckActGuo()
		result.Header = newProto.SuccessHeader()
		*result.UserId = user.GetUserId()
		user.WriteMsg(result)
		go desk.ActOut(userId, 0, true) //长沙杠自后，点击过，自动打牌
		return nil
	} else {
		//没有处理直接返回
		result := newProto.NewGame_AckActGuo()
		result.Header = newProto.SuccessHeader()
		*result.UserId = user.GetUserId()
		user.WriteMsg(result)
		return nil
	}
	return nil
}
