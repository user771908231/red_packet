package doudizhu

import (
	"errors"
	"casino_server/common/log"
	"sync"
	"casino_server/common/Error"
	"fmt"
	"sync/atomic"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/gate"
	"casino_doudizhu/msg/protogo"
	"casino_doudizhu/msg/funcsInit"
)

//斗地主的desk
type DdzDesk struct {
	sync.Mutex
	*PDdzDesk
	Users []*DdzUser
}

//斗地主的桌子//把数据同步到redis中去
func (d *DdzDesk) Update2Redis() error {
	bak := NewPDdzbak()
	bak.Desk = d.PDdzDesk
	for _, u := range d.Users {
		if u != nil && u.PDdzUser != nil {
			bak.Users = append(bak.Users, u.PDdzUser)
		}
	}

	//备份desk的数据
	UpdateDesk2Redis(bak)
	return nil
}

//得到一个用户
func (d *DdzDesk) GetUserByUserId(userId uint32) *DdzUser {
	for _, u := range d.Users {
		if u != nil && u.GetUserId() == userId {
			return u
		}
	}

	//哈哈哈
	return nil
}

func (d *DdzDesk) AddCountQiangDiZhu() {
	atomic.AddInt32(d.Tongji.CountQiangDiZhu, 1)
}

//设置低分
func (d *DdzDesk) setBaseValue(value int64) {
	*d.BaseValue = value
}

func (d *DdzDesk) setWinValue(value int64) {
	*d.WinValue = value
}

func (d *DdzDesk) setQingDizhuValue(value int64) {
	*d.QingDizhuValue = value
}

func (d *DdzDesk) addBombTongjiInfo(bomb *POutPokerPais) {
	d.Tongji.Bombs = append(d.Tongji.Bombs, bomb)
}

//添加一个玩家
func (d *DdzDesk) AddUser(userId uint32, a gate.Agent) error {
	user := NewDdzUser()
	user.GameData = NewPGameData()
	*user.UserId = userId
	user.agent = a
	*user.DeskId = d.GetDeskId()
	*user.RoomId = d.GetRoomId()
	user.UpdateSession()
	err := d.AddUserBean(user)

	return err
}

func (d *DdzDesk) AddUserBean(user *DdzUser) error {
	for i := 0; i < len(d.Users); i++ {
		if d.Users[i] == nil {
			d.Users[i] = user
			return nil
		}
	}
	log.E("玩家[%v]加入desk[%v]失败，因为没有合适的座位.", user.GetUserId(), d.GetDeskId())
	return errors.New("加入失败，没有找到合适的座位...")
}

//都准备
func (d *DdzDesk) IsAllReady() bool {
	for _, user := range d.Users {
		if user != nil && user.IsNotReady() {
			return false
		}
	}

	return true
}



//都确认了是否加倍
func (d *DdzDesk) IsAllActJiaBei() bool {
	for _, user := range d.Users {
		if user != nil && user.IsJiaBeiNoAct() {
			return false
		}
	}

	return true
}

func (d *DdzDesk) CheckOutPai(out *POutPokerPais) error {
	right, err := out.GT(d.OutPai)
	if err != nil {
		log.E("出牌的时候，判断牌型的时候失败...")
		return Error.NewError(-1, "比较失败,出牌失败...")
	}

	if right {
		return nil
	} else {
		return Error.NewError(-1, "出的牌比别人的牌小，没有办法出牌")
	}

}

func (d *DdzDesk) GetUserIndexByUserId(userId uint32) int {
	for i, user := range d.Users {
		if user != nil && user.GetUserId() == userId {
			return i
		}
	}
	return -1;
}

func (d *DdzDesk) SetActiveUser(userId uint32) {
	*d.ActiveUserId = userId
}

func (d *DdzDesk) SetDizhu(userId uint32) {
	*d.DiZhuUserId = userId
}

//判断是否是当前活动玩家
func (d *DdzDesk) CheckActiveUser(userId uint32) error {
	if d.GetActiveUserId() == userId {
		return nil
	} else {
		log.E("服务器错误：desk的当前操作玩家[%v]和请求的userid[%v]不一致", d.GetActiveUserId(), userId)
		return Error.NewFailError(fmt.Sprintf("当前活动玩家是[%v]", d.GetActiveUserId()))
	}
}

//判断用户的身份是不是地主
func (d *DdzDesk) IsDiZhuRole(user *DdzUser) bool {
	return d.GetDiZhuUserId() == user.GetUserId()
}

//广播协议
func (d *DdzDesk) BroadCastProto(p proto.Message) error {
	for _, u := range d.Users {
		//user不为空，并且user没有离开，没有短线的时候才能发送消息...
		if u != nil && !u.GetIsBreak() && !u.GetIsLeave() {
			u.WriteMsg(p)
		}
	}
	return nil
}

// 广播 但是不好办 userId
func (d *DdzDesk) BroadCastProtoExclusive(p proto.Message, userId uint32) error {
	for _, u := range d.Users {
		if u != nil && u.GetUserId() != userId {
			u.WriteMsg(p)
		}
	}
	return nil
}

func (d *DdzDesk) IsEnoughUser() bool {
	var count int32 = 0
	for _, user := range d.Users {
		if user != nil {
			count++
		}
	}
	return count == d.GetUserCountLimit()
}

//得到deskInfo
func (d *DdzDesk) GetDdzDeskInfo() *ddzproto.DdzDeskInfo {
	deskInfo := newProto.NewDdzDeskInfo()

	//*deskInfo.ActionTime = d.GetActionTime() //当前操作时间
	//*deskInfo.NInitActionTime = d.GetInitActionTime() //初始操作时间

	*deskInfo.ActiveUserId = d.GetActiveUserId() //当前操作人
	*deskInfo.CurrPlayCount = d.GetCurrPlayCount() //todo desk的游戏局数
	*deskInfo.DiZhuUserId = d.GetDiZhuUserId() //地主id
	//deskInfo.FootPokers = d.GetDiPokerPai() //底牌 todo change type []*PPokerPai to []*Poker
	*deskInfo.FootRate = d.GetFootRate() //todo
	*deskInfo.GameStatus = d.GetGameStatus() //todo
	*deskInfo.InitRoomCoin = d.GetInitRoomCoin() //todo
	*deskInfo.PlayerNum = d.GetPlayerNum() //todo
	*deskInfo.RoomNumber = d.GetKey()
	//deskInfo.RoomTypeInfo = d.GetRoomType()
	*deskInfo.PlayRate = d.GetPlayRate() //todo
	*deskInfo.TotalPlayCount = d.GetTotalPlayCount() //todo
	return deskInfo
}

func (d *DdzDesk) GetDdzSendGameInfo(SenderUserId uint32, isReconnect int32) *ddzproto.DdzSendGameInfo {
	ret := newProto.NewDdzSendGameInfo()
	ret.PlayerInfo = d.GetPlayerInfo()
	ret.DdzDeskInfo = d.GetDdzDeskInfo()
	*ret.SenderUserId = SenderUserId
	*ret.IsReconnect = isReconnect
	return ret
}

func (d *DdzDesk) GetPlayerInfo() []*ddzproto.PlayerInfo {
	var infos []*ddzproto.PlayerInfo
	for _, u := range d.Users {
		if u != nil {
			infos = append(infos, u.GetPlayerInfo(d))
		}
	}

	//返回用户信息
	return infos
}

//是否是四川斗地主
func (d *DdzDesk) IsSiChuanDouDiZhu() bool {
	//todo
	return false
}

func (d *DdzDesk) IsHuanLeDoudDiZhu() bool {
	//todo
	return true
}

func (d *DdzDesk) IsJingDianDouDiZhu() bool {
	//todo
	return true
}

func (d *DdzDesk) IsHadDiZhuUser() bool {
	return d.GetDiZhuUserId() != 0
}

//更具条件得到下一个玩家
func (d *DdzDesk) GetNextUserByPros(preUserId uint32, check func(nu *DdzUser) bool) *DdzUser {
	//查找下一家抢地主的人
	index := d.GetUserIndexByUserId(preUserId)
	if index < 0 {
		log.E("服务器错误，没有找到玩家[%v]的index...", preUserId)
		return nil
	}
	var nextUser *DdzUser
	for i := index + 1; i < len(d.Users) + index; i++ {
		u := d.Users[(i) / len(d.Users)]
		if check(u) {
			nextUser = u
			break
		}
	}
	return nextUser
}


//todo 需要做超时的处理
func (d *DdzDesk) SendChuPaiOverTurn(userId uint32) error {

	//轮到谁出牌时发送的overTurn
	//设置activeUser
	d.SetActiveUser(userId)

	overTurn := newProto.NewDdzOverTurn()
	*overTurn.UserId = userId
	*overTurn.ActType = ddzproto.ActType_T_NORMAL_ACT        //普通出牌
	d.BroadCastProto(overTurn)

	return nil
}

//发送房间的信息
func (d *DdzDesk) SendGameDeskInfo(sendUserId uint32, isReconnect int32) error {
	info := d.GetDdzSendGameInfo(sendUserId, isReconnect)
	d.BroadCastProto(info)
	return nil
}

//每个人都做一件事
/*
	1,user != nil
	2,func(user)
*/
func (d *DdzDesk) EveryUserDoSomething(dos func(user *DdzUser) error) error {
	for _, u := range d.Users {
		if u != nil {
			err := dos(u)
			if err != nil {
				log.E("玩家dos的时候出错err[%v]", err)
			}
		}
	}
	return nil
}


