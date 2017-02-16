package skeleton

import (
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"casino_majiang/service/majiang"
	"sync/atomic"
	"casino_common/common/log"
	"errors"
	"casino_common/common/userService"
	"github.com/golang/protobuf/proto"
)

func (u *SkeletonMJUser) GetUserData() *data.MJUserData {
	return u.UserData
}

func (u *SkeletonMJUser) GetUserId() uint32 {
	return u.userId
}

func (u *SkeletonMJUser) GetStatus() *data.MjUserStatus {
	return u.status
}

//todo
func (u *SkeletonMJUser) GetDesk() api.MjDesk {
	return u.desk
}

/***************************************账单相关***************************************/
func (u *SkeletonMJUser) GetBill() *majiang.Bill {
	return u.Bill
}

func (u *SkeletonMJUser) SubBillAmount(amount int64) {
	atomic.AddInt64(u.GetBill().WinAmount, -amount)
}

func (u *SkeletonMJUser) AddBillAmount(amount int64) {
	atomic.AddInt64(u.GetBill().WinAmount, amount)
}

//删除账单
func (u *SkeletonMJUser) DelBillBean(pai *majiang.MJPai) (error, *majiang.BillBean) {
	var bean *majiang.BillBean
	index := -1
	for i, info := range u.GetBill().GetBills() {
		if info != nil && info.GetPai().GetIndex() == pai.GetIndex() {
			index = i
			bean = info
			break
		}
	}

	if index > -1 {
		u.Bill.Bills = append(u.Bill.Bills[:index], u.Bill.Bills[index+1:]...)
		u.SubBillAmount(bean.GetAmount()) //减去
		return nil, bean
	} else {
		log.E("服务器错误：删除账单 billBean的时候出错，没有找到对应的杠牌[%v]", pai)
		return errors.New("删除手牌时出错，没有找到对应的手牌..."), nil
	}
}

func (u *SkeletonMJUser) GetGameData() *data.MJUserGameData {
	return u.GameData
}

//增加一条账单
func (u *SkeletonMJUser) AddBillBean(bean *majiang.BillBean) error {
	u.Bill.Bills = append(u.Bill.Bills, bean)
	u.AddBillAmount(bean.GetAmount())
	return nil
}

//增加账单
func (u *SkeletonMJUser) AddBill(relationUserid uint32, billType int32, des string, score int64, pai *majiang.MJPai, roomType int32) error {
	//用户赢钱的账户,赢钱的账单
	bill := majiang.NewBillBean()
	*bill.UserId = u.GetUserId()
	*bill.OutUserId = relationUserid
	//*bill.Type = MJUSER_BILL_TYPE_YING_HU
	*bill.Type = billType
	*bill.Des = des
	*bill.Amount = score //杠牌的收入金额
	bill.Pai = pai
	u.AddBillBean(bill)

	//计算账单的地方 来加减用户的coin
	u.AddCoin(score, roomType)    //统计用户剩余多少钱
	u.AddStatisticsWinCoin(score) //统计用户输赢多少钱
	return nil
}

//更新用户金币
func (u *SkeletonMJUser) AddCoin(coin int64, roomType int32) {
	//减少玩家金额
	atomic.AddInt64(&u.GetUserData().Coin, coin) //更新账户余额
	//如果是金币场。需要更新用户的金币余额
	if roomType == majiang.ROOMTYPE_COINPLAY {
		remainCoin, _ := userService.INCRUserCOIN(u.GetUserId(), coin)
		u.GetUserData().Coin = proto.Int64(remainCoin) //增加用户的金币
	}
}

func (u *SkeletonMJUser) GetSkeletonMJUser() *SkeletonMJUser {
	return u
}
