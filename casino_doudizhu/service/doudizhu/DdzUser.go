package doudizhu

import (
	"sync"
	"casino_server/common/log"
	"casino_server/common/Error"
	"sync/atomic"
	"github.com/name5566/leaf/gate"
	"github.com/golang/protobuf/proto"
)

var (
	DDZUSER_STATUS_READY int32 = 2        //已经准备

)

var (
	DDZUSER_QIANGDIZHU_STATUS_NOACT int32 = 0 //没操作
	DDZUSER_QIANGDIZHU_STATUS_JIAO int32 = 1 //抢地主
	DDZUSER_QIANGDIZHU_STATUS_QIANG int32 = 2 //抢地主
	DDZUSER_QIANGDIZHU_STATUS_PASS int32 = 3 //不叫
)

var (
	DDZUSER_JIABEI_STATUS_NOACT int32 = 0 //没操作
	DDZUSER_JIABEI_STATUS_JIABEI int32 = 1 //抢地主
	DDZUSER_JIABEI_STATUS_BUJIABEI int32 = 2 //抢地主
)

type DdzUser struct {
	sync.Mutex
	*PDdzUser
	agent gate.Agent
}

//清楚session
func (u *DdzUser)ClearAgentGameData() {

}

func (u *DdzUser) SetOnline() error {
	*u.IsBreak = false
	return nil
}

func (u *DdzUser) UpdateSession() {
	//更新session的信息
	UpdateSession(u.GetUserId(), u.GetSessionGameStatus(), u.GetRoomId(), u.GetDeskId())
}

//得到游戏状态
func (u *DdzUser) GetSessionGameStatus() int32 {
	return 0
}

//设置庄太
func (u *DdzUser) SetStatus(s int32) {
	*u.Status = s
}

func (u *DdzUser) SetQiangDiZhuStatus(s int32) {
	*u.QiangDiZhuStatus = s
}

//设置连接
func (u *DdzUser) setAgent(a gate.Agent) {
	u.agent = a
}

func (u *DdzUser) IsReady() bool {
	return u.GetStatus() == DDZUSER_STATUS_READY
}

func (u *DdzUser) IsNotReady() bool {
	return !u.IsReady()
}

//是否抢地主
func (u *DdzUser) IsQiangDiZhuQiang() bool {
	return u.GetQiangDiZhuStatus() == DDZUSER_QIANGDIZHU_STATUS_QIANG
}

//是否不叫
func (u *DdzUser) IsQiangDiZhuBuJiao() bool {
	return u.GetQiangDiZhuStatus() == DDZUSER_QIANGDIZHU_STATUS_PASS
}

//抢注的时候还没有操作
func (u *DdzUser) IsQiangDiZhuNoAct() bool {
	return u.GetQiangDiZhuStatus() == DDZUSER_QIANGDIZHU_STATUS_NOACT
}

//没有加倍操作
func (u *DdzUser) IsJiaBeiNoAct() bool {
	return u.GetJiabeiStatus() == DDZUSER_JIABEI_STATUS_NOACT

}

//加倍
func (u *DdzUser) IsJiaBeiJiaBei() bool {
	return u.GetJiabeiStatus() == DDZUSER_JIABEI_STATUS_JIABEI

}

//不加倍
func (u *DdzUser) IsJiaBeiBuJiaBei() bool {
	return u.GetJiabeiStatus() == DDZUSER_JIABEI_STATUS_BUJIABEI

}

func (u *DdzUser) DelHandlPai(pais *PPokerPai) error {
	index := -1
	for i, pai := range u.GameData.HandPokers {
		if pai != nil && pai.GetId() == pais.GetId() {
			index = i
			break
		}
	}
	if index > -1 {
		u.GameData.HandPokers = append(u.GameData.HandPokers[:index], u.GameData.HandPokers[index + 1:]...)
		return nil

	} else {
		log.E("服务器错误：删除手牌的时候出错，没有找到对应的手牌[%v]", pais)
		return Error.NewError(-1, "删除手牌时出错，没有找到对应的手牌...")
	}

}

//增加出牌
func (u *DdzUser) DOPoutPokerPais(out *POutPokerPais) error {
	//1，增加出的牌
	u.GameData.OutPaiList = append(u.GameData.OutPaiList, out)

	//2，删除手牌
	for _, p := range out.PokerPais {
		u.DelHandlPai(p)
	}

	return nil
}

//得到玩家手牌的张数
func (u *DdzUser) GetHandPaiCount() int32 {
	return int32(len(u.GameData.HandPokers))
}

func (u *DdzUser) beginInit() error {
	return nil
}

func (u *DdzUser) AddNewBill(coin int64, winUser, loseUser uint32, des string) {
	bean := NewPDdzBillBean()
	*bean.Coin = coin
	*bean.WinUser = winUser
	*bean.LoseUser = loseUser
	*bean.Desc = des
	//增加账单
	u.Bill.BillBean = append(u.Bill.BillBean, bean)
	//增加输的钱
	atomic.AddInt64(u.Bill.WinCoin, coin)
}

//游戏中的玩家放信息
func (u *DdzUser) WriteMsg(msg proto.Message) {
	agent := u.agent
	if agent == nil {
		log.E("玩家[%v]发送信息失败", u.GetUserId())
	} else {
		agent.WriteMsg(msg)
	}
}