package majiang

import (
	"casino_common/common/log"
)

//处理账单
/**
	没有胡牌的人，都需要给钱  ,目前不是承包的方式...
 */

func (d *MjDesk) DoGangBill(info *GangPaiInfo) {
	gangType := info.GetGangType()
	gangUser := d.GetUserByUserId(info.GetGetUserId())
	gangPai := info.GetPai()

	if gangType == GANG_TYPE_AN {
		//处理暗杠的账单
		score := d.GetBaseValue() * 2 //暗杠的分数
		for _, ou := range d.GetUsers() {
			//不为nil 并且不是本人，并且没有胡牌
			if ou != nil && (ou.GetUserId() != gangUser.GetUserId()) && ou.IsGaming() && ou.IsNotHu() {
				gangUser.AddBill(ou.GetUserId(), MJUSER_BILL_TYPE_YING_AN_GNAG, "用户暗杠，收入", score, gangPai, d.GetRoomType()) //用户赢钱的账户
				ou.AddBill(gangUser.GetUserId(), MJUSER_BILL_TYPE_SHU_AN_GNAG, "用户暗杠，输钱", -score, gangPai, d.GetRoomType()) //用户输钱的账单

				ou.AddStatisticsCountBeiAnGang(d.GetCurrPlayCount()) //被暗杠用户的统计信息
			} else if ou != nil && (ou.GetUserId() == gangUser.GetUserId()) && ou.IsGaming() && ou.IsNotHu() {
				gangUser.AddStatisticsCountAnGang(d.GetCurrPlayCount()) //暗杠用户的统计信息
			}
		}

	} else if gangType == GANG_TYPE_DIAN {
		//处理点杠点账单
		score := d.GetBaseValue() * 2 //点杠的分数
		shuUser := d.GetUserByUserId(info.GetSendUserId())
		gangUser.AddBill(shuUser.GetUserId(), MJUSER_BILL_TYPE_YING_GNAG, "用户点杠，收入", score, gangPai, d.GetRoomType()) //用户赢钱的账户
		shuUser.AddBill(gangUser.GetUserId(), MJUSER_BILL_TYPE_SHU_GNAG, "用户点杠，输钱", -score, gangPai, d.GetRoomType()) //用户输钱的账单

		gangUser.AddStatisticsCountMingGang(d.GetCurrPlayCount()) //明杠用户的统计信息
		shuUser.AddStatisticsCountDianGang(d.GetCurrPlayCount())  //点杠用户的统计信息

	} else if gangType == GANG_TYPE_BA {
		//处理巴杠的账单
		score := d.GetBaseValue() //巴杠的分数
		for _, ou := range d.GetUsers() {
			if ou != nil && (ou.GetUserId() != gangUser.GetUserId()) && ou.IsGaming() && ou.IsNotHu() {
				//账单多次添加
				gangUser.AddBill(ou.GetUserId(), MJUSER_BILL_TYPE_YING_BA_GANG, "用户巴杠，收入", score, gangPai, d.GetRoomType()) //用户赢钱的账户
				ou.AddBill(gangUser.GetUserId(), MJUSER_BILL_TYPE_SHU_BA_GANG, "用户巴杠，输钱", -score, gangPai, d.GetRoomType()) //用户输钱的账单

				ou.AddStatisticsCountBeiBaGang(d.GetCurrPlayCount()) //被巴杠用户的统计信息

			} else if ou != nil && (ou.GetUserId() == gangUser.GetUserId()) && ou.IsGaming() && ou.IsNotHu() {
				gangUser.AddStatisticsCountBaGang(d.GetCurrPlayCount()) //巴杠用户的统计信息
			}
		}
	}
}

//巴杠之后需要初始化抢杠的CheckCase
func (d *MjDesk) InitCheckCaseAfterGang(gangType int32, gangPai *MJPai, user *MjUser) {
	d.CheckCase = nil //设置 判断的为nil
	///如果是巴杠，需要设置巴杠的判断  initCheckCase
	if gangType == GANG_TYPE_BA {
		d.InitCheckCase(gangPai, user, true) //杠牌之后 初始化checkCase
		if d.CheckCase == nil {
			log.T("%v 玩家[%v]巴杠没有人可以抢杠...", d.DlogDes(), user.GetUserId())
		}

	}

}
