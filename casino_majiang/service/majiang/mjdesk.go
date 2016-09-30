package majiang

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"github.com/name5566/leaf/gate"
	"casino_majiang/conf/log"
)

//状态表示的是当前状态.
var MJDESK_STATUS_CREATED = 1 //刚刚创建
var MJDESK_STATUS_READY = 2//正在准备
var MJDESK_STATUS_ONINIT = 3//准备完成之后，desk初始化数据
var MJDESK_STATUS_EXCHANGE = 4//desk初始化完成之后，告诉玩家可以开始换牌
var MJDESK_STATUS_DINGQUE = 5//换牌结束之后，告诉玩家可以开始定缺
var MJDESK_STATUS_RUNNING = 6 //定缺之后，开始打牌
var MJDESK_STATUS_LOTTERY = 7 //结算
var MJDESK_STATUS_END = 8//一局结束


//判断是不是朋友桌
func (d *MjDesk) IsFriend() bool {
	return true
}

//用户加入房间
func (d *MjDesk) addNewUser(userId uint32, a gate.Agent) error {
	if d.IsFriend() {
		return d.addNewUserFriend(userId, a)
	} else {
		return nil
	}
}

func (d *MjDesk) addNewUserFriend(userId uint32, a gate.Agent) error {
	//1,是否是重新进入
	user := d.getUserByUserId(userId)
	if user != nil {
		//是断线重连
		*user.IsBreak = false;
		return nil

	}

	//2,是否是离开之后重新进入房间
	userLeave := d.getLeaveUserByUserId(userId)
	if userLeave != nil {
		err := d.addUser(userLeave)
		if err != nil {
			//加入房间的时候错误
			log.E("已经离开的用户[%v]重新加入desk[%v]的时候出错errMsg[%v]", userId, d.GetDeskId(), err)
			return errors.New("已经离开的用户重新加入房间的时候出错")
		} else {
			*userLeave.IsBreak = false
			*userLeave.IsLeave = false
			d.rmLeaveUser(userLeave.GetUserId())
			return nil
		}
	}

	//加入一个新用户
	newUser := NewMjUser()
	*newUser.DeskId = d.GetDeskId()
	*newUser.RoomId = d.GetRoomId()
	*newUser.Coin = d.GetBaseValue()
	*newUser.IsBreak = false
	*newUser.IsLeave = false
	*newUser.Status = MJUSER_STATUS_INTOROOM
	err := d.addUser(newUser)
	if err != nil {
		log.E("用户[%v]加入房间[%v]失败,errMsg[%v]", )
		return errors.New("用户加入房间失败")
	} else {
		//加入房间成功
		return nil
	}
}


//删除已经离开的人的Id todo
func (d *MjDesk) rmLeaveUser(userId uint32) error {
	return nil
}

//通过userId得到user
func (d *MjDesk) getUserByUserId(userId uint32) *MjUser {
	for _, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == userId {
			return u
		}
	}

	return nil
}

//通过userId得到user
func (d *MjDesk) getLeaveUserByUserId(userId uint32) *MjUser {
	for _, u := range d.GetUsers() {
		if u != nil && u.GetUserId() == userId {
			return u
		}
	}

	return nil
}


//新增加一个玩家
func (d *MjDesk) addUser(user *MjUser) error {
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
		return errors.New("没有找到合适的位置，加入桌子失败")
	}
}

func (d *MjDesk) GetRoomTypeInfo() *mjproto.RoomTypeInfo {
	typeInfo := newProto.NewRoomTypeInfo()
	*typeInfo.Settlement = d.GetSettlement()
	typeInfo.PlayOptions = d.GetPlayOptions()
	*typeInfo.MjRoomType = mjproto.MJRoomType(d.GetMjRoomType())
	*typeInfo.BaseValue = d.GetBaseValue()
	*typeInfo.BoardsCout = d.GetBoardsCout()
	*typeInfo.CapMax = d.GetCapMax()
	*typeInfo.CardsNum = d.GetCardsNum()
	return typeInfo
}

func (d *MjDesk) GetPlayOptions() *mjproto.PlayOptions {
	o := newProto.NewPlayOptions()
	*o.ZiMoRadio = d.GetZiMoRadio()
	*o.HuRadio = d.GetHuRadio()
	*o.DianGangHuaRadio = d.GetDianGangHuaRadio()
	o.OthersCheckBox = d.GetOthersCheckBox()
	return o
}

//广播协议
func (d *MjDesk) BroadCastProto(p proto.Message) error {
	for _, u := range d.Users {
		if u != nil {
			u.WriteMsg(p)
		}
	}
	return nil
}

func (d *MjDesk) GetDeskGameInfo() *mjproto.DeskGameInfo {
	return nil
}

func (d *MjDesk) GetGame_SendGameInfo() *mjproto.Game_SendGameInfo {
	return nil
}


//用户准备
func (d *MjDesk) Ready(userId  uint32) error {
	//找到需要准备的user
	user := d.getUserByUserId(userId)
	if user == nil {
		log.E("用户[%v]在desk[%v]准备的时候失败,没有找到对应的玩家", user.GetUserId(), d.GetDeskId())
		return errors.New("没有找到用户，准备失败")
	}

	//设置为准备的状态
	*user.Status = MJUSER_STATUS_READY        //用户准备

	return nil
}

//是不是所有人都准备
func (d *MjDesk) IsAllReady() bool {
	for _, u := range d.Users {
		if u != nil && !u.IsReady() {
			return false
		}
	}
	return true
}



//用户准备之后的一些操作
func (d *MjDesk) AfterReady() error {

	//如果所有人都准备了，那么开始游戏
	if d.IsAllReady() {
		d.begin()
	}

	return nil
}

//开始游戏
func (d *MjDesk) begin() {
	//1，检查是否可以开始游戏
	//2，初始化user的状态


	//3，初始化桌子的状态
	d.beginInit()


	//4，发13张牌
	d.initCards()


	//5，开始定缺
	d.beginDingQue()

}

func (d *MjDesk) beginInit() error {
	return nil
}

func (d *MjDesk) initCards() error {
	//得到一副已经洗好的麻将
	d.AllMJPai = XiPai()

	//更别给每个人发牌
	for i, u := range d.Users {
		if u != nil && u.IsGaming() {
			u.MJHnadPai.Pais = d.AllMJPai[i * 13: (i + 1) * 13]
			*d.MJPaiNexIndex = int32((i + 1) * 13);
		}
	}

	//发送发牌的广播
	return nil
}

//开始定缺
func (d *MjDesk) beginDingQue() error {
	return nil
}