package doudizhu

import (
	"sync"
	"sync/atomic"
	"github.com/name5566/leaf/gate"
	"github.com/golang/protobuf/proto"
	"casino_doudizhu/msg/protogo"
	"casino_doudizhu/msg/funcsInit"
	"sort"
	"casino_common/common/log"
	"casino_common/common/Error"
	"casino_common/common/userService"
)

var (
	DDZUSER_QIANGDIZHU_STATUS_NOACT int32 = 0 //没操作
	DDZUSER_QIANGDIZHU_STATUS_JIAO int32 = 1 //叫地主
	DDZUSER_QIANGDIZHU_STATUS_BUJIAO int32 = 2 //叫地主
	DDZUSER_QIANGDIZHU_STATUS_QIANG int32 = 3 //抢地主
	DDZUSER_QIANGDIZHU_STATUS_BUQIANG int32 = 4 //不叫
)

var (
	DDZUSER_JIABEI_STATUS_NOACT int32 = 0 //没操作
	DDZUSER_JIABEI_STATUS_JIABEI int32 = 1 //加倍
	DDZUSER_JIABEI_STATUS_BUJIABEI int32 = 2 //不加倍
)

var (
	DDZUSER_STATUS_DEFAULT int32 = 0        //准备
	DDZUSER_STATUS_READY int32 = 1        //准备
	DDZUSER_STATUS_JIAODIZHU int32 = 2        //叫地主
	DDZUSER_STATUS_JIABEI int32 = 3        //加倍
	DDZUSER_STATUS_GAMING int32 = 4        //游戏中
	DDZUSER_SATTUS_LOTTERY int32 = 5        //
	DDZUSER_STATUS_END int32 = 6                //
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

func (u *DdzUser) SetJiaBeiStatus(s int32) {
	*u.JiabeiStatus = s
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
	return u.GetQiangDiZhuStatus() == DDZUSER_QIANGDIZHU_STATUS_BUJIAO
}

func (u *DdzUser) IsQiangDiZhuJiao() bool {
	return u.GetQiangDiZhuStatus() == DDZUSER_QIANGDIZHU_STATUS_JIAO
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
	u.Bill = nil
	u.GameData = NewPGameData()
	u.SetJiaBeiStatus(DDZUSER_JIABEI_STATUS_NOACT)
	u.SetQiangDiZhuStatus(DDZUSER_QIANGDIZHU_STATUS_NOACT)
	u.SetStatus(DDZUSER_STATUS_DEFAULT)
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

func (u *DdzUser) GetNickName() string {
	return ""
}

func (u *DdzUser) GetSex() int32 {
	return 0
}

func (u *DdzUser) GetBReady() int32 {
	if u.IsReady() {
		return 1
	} else {
		return 0
	}
}

func (u *DdzUser) GetWxInfo() *ddzproto.WeixinInfo {
	user := userService.GetUserById(u.GetUserId())
	wx := newProto.NewWeixinInfo()
	*wx.City = user.GetCity()
	*wx.HeadUrl = user.GetHeadUrl()
	*wx.NickName = user.GetNickName()
	*wx.OpenId = user.GetOpenId()
	return nil
}

func (u *DdzUser) GetGameStatus() int32 {
	return 0
}

func (u *DdzUser) GetPlayerPokers() []*ddzproto.Poker {
	var list []*ddzproto.Poker
	for _, p := range u.GameData.HandPokers {
		if p != nil {
			list = append(list, p.GetClientPoker())
		}
	}
	return list
}

//返回玩家的游戏装套
func (u *DdzUser) GetPlayerGameStatus() *ddzproto.PlayerGameStatus {

	return nil
}

func (u *DdzUser) GetOnlineStatus() int32 {
	return 0
}

func ( u *DdzUser) GetPlayerInfo(desk *DdzDesk) *ddzproto.PlayerInfo {
	info := newProto.NewPlayerInfo()
	*info.IsDiZhu = desk.GetDiZhuUserId() == u.GetUserId()        //是否是地主
	info.PlayerPokers = u.GetPlayerPokers()        //玩家的扑克牌
	*info.Coin = u.GetCoin()        //玩家的coin
	*info.NickName = u.GetNickName()        //玩家的nickName
	*info.Sex = u.GetSex()        //玩家的性别
	*info.UserId = u.GetUserId()        //玩家的id
	*info.IsOwner = desk.GetOwner() == u.GetUserId()        //是否是房主
	*info.BReady = u.GetBReady()        //是否准备
	info.Status = u.GetPlayerGameStatus()        //游戏状态
	info.WxInfo = u.GetWxInfo()        //微信信息
	*info.OnlineStatus = u.GetOnlineStatus() //在线的状态
	return info
}

func (u *DdzUser) GetTransferredHandPokerPais() string {
	ret := ""
	suit := ""
	handPokers := DdzPokerOutList{}
	handPokers = u.GameData.HandPokers
	sort.Sort(handPokers)
	for _, p := range handPokers {
		suit = p.GetLogDes()
		if suit != "" {
			ret = ret + suit + "\t "
		}
	}
	return ret
}

func (u *DdzUser) GetTransferredOutPokerPais() string {
	ret := ""
	suit := ""
	for _, outList := range u.GameData.GetOutPaiList() {
		if outList != nil {
			outLists := DdzPokerOutList{}
			outLists = outList.GetPokerPais()
			sort.Sort(outLists)
			for _, p := range outLists {
				suit = p.GetLogDes()
				if suit != "" {
					ret = ret + suit + "\t "
				}
			}
		}
	}
	return ret
}

func (u *DdzUser) GetTransferredStatus() string {
	ret := ""
	switch u.GetStatus() {
	case DDZUSER_STATUS_JIABEI:
		ret = "加倍"
	case DDZUSER_STATUS_END:
		ret = "已结束"
	case DDZUSER_STATUS_DEFAULT:
		ret = "默认"
	case DDZUSER_STATUS_GAMING:
		ret = "游戏中"
	case DDZUSER_STATUS_JIAODIZHU:
		ret = "叫地主"
	case DDZUSER_STATUS_READY:
		ret = "准备"
	case DDZUSER_SATTUS_LOTTERY:
		ret = "开奖"
	default:
	}
	return ret
}

func (u *DdzUser) GetTransferredJiaBeiStatus() string {
	ret := ""
	switch u.GetJiabeiStatus() {
	case DDZUSER_JIABEI_STATUS_BUJIABEI:
		ret = "不加倍"
	case DDZUSER_JIABEI_STATUS_JIABEI:
		ret = "加倍"
	case DDZUSER_JIABEI_STATUS_NOACT:
		ret = "未设置"
	default:
	}
	return ret
}

func (u *DdzUser) GetTransferredQiangDiZhuStatus() string {
	ret := ""
	switch u.GetQiangDiZhuStatus() {
	case DDZUSER_QIANGDIZHU_STATUS_BUJIAO:
		ret = "不叫"
	case DDZUSER_QIANGDIZHU_STATUS_BUQIANG:
		ret = "不抢"
	case DDZUSER_QIANGDIZHU_STATUS_JIAO:
		ret = "叫地主"
	case DDZUSER_QIANGDIZHU_STATUS_NOACT:
		ret = "未设置"
	case DDZUSER_QIANGDIZHU_STATUS_QIANG:
		ret = "抢地主"
	default:
	}
	return ret
}