package paosangong

import (
	"errors"
	"casino_common/proto/ddproto"
	"time"
	"casino_common/common/Error"
)

//牌桌托管
func (desk *Desk) DoTuoguan() {
	for _,u := range desk.Users {
		if u.GetIsTuoguan() {
			u.DoTuoguanAct()
		}
	}
}

//自动托管
func (user *User) DoTuoguanAct() error {
	defer Error.ErrorRecovery("DoTuoguanAct()")

	if !user.GetIsTuoguan() {
		return errors.New("not tuoguan.")
	}

	//延时1秒
	<-time.After(1 * time.Second)

	switch user.Desk.GetStatus() {
	case ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_ENTER, ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_READY, ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START:
		//准备阶段
		if !user.GetIsReady() {
			if user.Desk.GetCircleNo() == 1 {
				//金币场：直接准备
				if user.Desk.GetDeskNumber() == "" {
					user.DoReady()
				}else {
					if !user.IsOwner() {
						//不是房主,直接准备
						user.DoReady()
					}else {
						//是房主
						if user.Desk.GetStatus() == ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_START {
							user.DoReady()
						}
					}
				}
			}else {
				//非第一局,非最后一局
				//直接准备
				if user.Desk.GetCircleNo() <= user.Desk.DeskOption.GetMaxCircle() {
					user.DoReady()
				}
			}
		}
	case ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_QIANGZHUANG:
		//抢庄阶段
		if user.GetBankerScore() == 0 {
			poker_type := user.Pokers.GetType()
			qz_opt := user.TuoGuanOpt.GetQiangZhuangOpt()

			switch  {
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_BU_QIANG:
				//不抢
				user.DoQiangzhuang(-1)
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_RENYI_QIANG_1:
				//任意牌抢1倍
				user.DoQiangzhuang(1)
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_RENYI_QIANG_2:
				//任意牌抢2倍
				user.DoQiangzhuang(2)
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_NIU8_QIANG_1 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_8:
				//牛8抢1倍
				user.DoQiangzhuang(1)
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_NIU8_QIANG_2 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_8:
				//牛8抢2倍
				user.DoQiangzhuang(2)
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_NIU9_QIANG_1 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_9:
				//牛9抢1倍
				user.DoQiangzhuang(1)
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_NIU9_QIANG_2 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_9:
				//牛9抢2倍
				user.DoQiangzhuang(2)
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_NIUNIU_QIANG_1 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_NIU:
				//牛牛抢1倍
				user.DoQiangzhuang(1)
			case qz_opt == ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_NIUNIU_QIANG_2 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_NIU:
				//牛牛抢2倍
				user.DoQiangzhuang(2)
			default:
				//默认不抢
				user.DoQiangzhuang(-1)
			}
		}
	case ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_JIABEI:
		//加倍阶段
		if user.GetDoubleScore() == 0 && !user.IsBanker() {
			poker_type := user.Pokers.GetType()
			yz_opt := user.TuoGuanOpt.GetYaZhuOpt()

			switch  {
			case yz_opt == ddproto.NiuEnumTuoguanYzopt_NIU_TG_YZ_YA_1:
				//押1分
				user.DoJiabei(1 * int64(user.Desk.DeskOption.GetDiFen()))
			case yz_opt == ddproto.NiuEnumTuoguanYzopt_NIU_TG_YZ_YA_2:
				//押2分
				user.DoJiabei(2 * int64(user.Desk.DeskOption.GetDiFen()))
			case yz_opt == ddproto.NiuEnumTuoguanYzopt_NIU_TG_YZ_NIU8_YA_2 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_8:
				//牛8押2分
				user.DoJiabei(2 * int64(user.Desk.DeskOption.GetDiFen()))
			case yz_opt == ddproto.NiuEnumTuoguanYzopt_NIU_TG_YZ_NIU9_YA_2 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_9:
				//牛9押2分
				user.DoJiabei(2 * int64(user.Desk.DeskOption.GetDiFen()))
			case yz_opt == ddproto.NiuEnumTuoguanYzopt_NIU_TG_YZ_NIUNIU_YA_2 && poker_type >= ddproto.NiuniuEnum_PokerType_NIU_NIU:
				//牛牛押2分
				user.DoJiabei(2 * int64(user.Desk.DeskOption.GetDiFen()))
			default:
				//默认押1
				user.DoJiabei(1 * int64(user.Desk.DeskOption.GetDiFen()))
			}
		}
	case ddproto.NiuniuEnumDeskState_NIU_DESK_STATUS_WAIT_LIANGPAI:
		//亮牌阶段
		user.DoLiangpai()
	}

	return nil
}
