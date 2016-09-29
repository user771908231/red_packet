package majiang

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
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

//
func (d *MjDesk) addNewUser(userId uint32) error {
	//
	if d.IsFriend() {
		//1,是否是重新进入

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

